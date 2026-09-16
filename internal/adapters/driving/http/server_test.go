package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeRegisterConsumerUseCase trägt eine In-Memory-Fälschung des Inbound
// Ports (`ADR-0030`): Whitebox-Test des Adapters ohne reale Persistenz —
// die Registrierungslogik selbst ist bereits am Use Case getestet
// (`internal/application/usecase/register`), dieser Test belegt nur die
// Übersetzung HTTP↔Use-Case.
type fakeRegisterConsumerUseCase struct {
	registered map[model.ConsumerID]string
}

func newFakeRegisterConsumerUseCase() *fakeRegisterConsumerUseCase {
	return &fakeRegisterConsumerUseCase{registered: map[model.ConsumerID]string{}}
}

func (f *fakeRegisterConsumerUseCase) Register(_ context.Context, command inbound.RegisterConsumerCommand) (inbound.RegisterConsumerResult, error) {
	consumer, err := model.NewConsumer(command.Consumer, command.Name)
	if err != nil {
		return inbound.RegisterConsumerResult{}, err
	}
	_, already := f.registered[consumer.ID]
	f.registered[consumer.ID] = consumer.Name
	return inbound.RegisterConsumerResult{Consumer: consumer, AlreadyRegistered: already}, nil
}

// fakeFailingUseCase trägt einen nicht-domänenspezifischen Fehler (z. B.
// eine Speicherstörung) — die Grenze zwischen `400` (Domänen-Invariante)
// und `500` (unerwarteter interner Fehler).
type fakeFailingUseCase struct{ err error }

func (f fakeFailingUseCase) Register(context.Context, inbound.RegisterConsumerCommand) (inbound.RegisterConsumerResult, error) {
	return inbound.RegisterConsumerResult{}, f.err
}

const (
	testReaderToken = "reader-token"
	testAdminToken  = "admin-token"
)

func newTestServer(t *testing.T, useCase inbound.RegisterConsumerUseCase) *httptest.Server {
	t.Helper()
	srv := New(Config{
		Addr:             "unused:0",
		TokenReader:      testReaderToken,
		TokenAdmin:       testAdminToken,
		RegisterConsumer: useCase,
	})
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

func postConsumer(t *testing.T, ts *httptest.Server, token, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/consumers", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Request bauen: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Request senden: %v", err)
	}
	return resp
}

// TestRegisterConsumerOhneTokenEndetMit401 trägt die Fitness Function aus
// `ADR-0057` Teilfrage 3: ein Aufruf ohne Bearer-Token endet mit `401`.
func TestRegisterConsumerOhneTokenEndetMit401(t *testing.T) {
	ts := newTestServer(t, newFakeRegisterConsumerUseCase())
	resp := postConsumer(t, ts, "", `{"consumer_id":"c1","name":"Consumer 1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401)", resp.StatusCode)
	}
}

// TestRegisterConsumerUnbekannterTokenEndetMit401 trägt denselben Pfad
// für ein Token, das keiner konfigurierten Klasse entspricht — nicht nur
// ein fehlender Header.
func TestRegisterConsumerUnbekannterTokenEndetMit401(t *testing.T) {
	ts := newTestServer(t, newFakeRegisterConsumerUseCase())
	resp := postConsumer(t, ts, "unbekannt", `{"consumer_id":"c1","name":"Consumer 1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401)", resp.StatusCode)
	}
}

// TestRegisterConsumerReaderTokenEndetMit403 trägt die zweite Hälfte der
// Fitness Function: ein gültiges `reader`-Token erreicht den schreibenden
// Endpunkt `RegisterConsumer` nicht.
func TestRegisterConsumerReaderTokenEndetMit403(t *testing.T) {
	ts := newTestServer(t, newFakeRegisterConsumerUseCase())
	resp := postConsumer(t, ts, testReaderToken, `{"consumer_id":"c1","name":"Consumer 1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Status: %d (Erwartung: 403)", resp.StatusCode)
	}
}

// TestRegisterConsumerAdminTokenRegistriert trägt den Happy-Path
// (`LH-FA-SST-006`, `LH-FA-CON-001`): ein gültiges `admin`-Token
// registriert den Consumer, die Antwort trägt die übersetzten Felder.
func TestRegisterConsumerAdminTokenRegistriert(t *testing.T) {
	ts := newTestServer(t, newFakeRegisterConsumerUseCase())
	resp := postConsumer(t, ts, testAdminToken, `{"consumer_id":"c1","name":"Consumer 1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Status: %d (Erwartung: 201)", resp.StatusCode)
	}
	var decoded registerConsumerResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.ConsumerID != "c1" || decoded.Name != "Consumer 1" || decoded.AlreadyRegistered {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestRegisterConsumerAdminTokenIstIdempotent trägt `LH-FA-CON-001`s
// Boundary: eine bereits registrierte Kennung meldet das Ergebnis, ohne
// den Aufruf abzulehnen.
func TestRegisterConsumerAdminTokenIstIdempotent(t *testing.T) {
	useCase := newFakeRegisterConsumerUseCase()
	ts := newTestServer(t, useCase)
	first := postConsumer(t, ts, testAdminToken, `{"consumer_id":"c1","name":"Consumer 1"}`)
	first.Body.Close()

	resp := postConsumer(t, ts, testAdminToken, `{"consumer_id":"c1","name":"Consumer 1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Status: %d (Erwartung: 201)", resp.StatusCode)
	}
	var decoded registerConsumerResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if !decoded.AlreadyRegistered {
		t.Fatalf("Antwort: %+v (Erwartung: already_registered=true)", decoded)
	}
}

// TestRegisterConsumerUngueltigesJSONEndetMit400 trägt die Formgrenze des
// Request-Bodys.
func TestRegisterConsumerUngueltigesJSONEndetMit400(t *testing.T) {
	ts := newTestServer(t, newFakeRegisterConsumerUseCase())
	resp := postConsumer(t, ts, testAdminToken, `{nicht-json`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// TestRegisterConsumerLeereKennungEndetMit400 trägt die Domänen-Invariante
// (`domainerrors.ErrEmptyIdentifier`) als `400`, nicht als `500`.
func TestRegisterConsumerLeereKennungEndetMit400(t *testing.T) {
	ts := newTestServer(t, newFakeRegisterConsumerUseCase())
	resp := postConsumer(t, ts, testAdminToken, `{"consumer_id":"","name":""}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// TestRegisterConsumerInternerFehlerEndetMit500 trägt die Abgrenzung zur
// Domänen-Invariante: ein nicht-domänenspezifischer Fehler (z. B. eine
// Speicherstörung) endet als `500`, nicht als `400`.
func TestRegisterConsumerInternerFehlerEndetMit500(t *testing.T) {
	ts := newTestServer(t, fakeFailingUseCase{err: errors.New("speicher nicht erreichbar")})
	resp := postConsumer(t, ts, testAdminToken, `{"consumer_id":"c1","name":"Consumer 1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Status: %d (Erwartung: 500)", resp.StatusCode)
	}
}

// TestClassifyTokenLeereKonfigurationTrifftKeinToken trägt das §6-Risiko
// dieses Slice: Ein leer konfiguriertes Token darf kein Aufruf-Token
// treffen — sonst würde ein fehlender Header (leerer Token) gegen eine
// ebenfalls ungesetzte Token-Klasse eine dritte, implizite Rechtsklasse
// eröffnen.
func TestClassifyTokenLeereKonfigurationTrifftKeinToken(t *testing.T) {
	if got := classifyToken("", "", ""); got != roleNone {
		t.Fatalf("classifyToken(\"\", \"\", \"\") = %v (Erwartung: roleNone)", got)
	}
	if got := classifyToken("beliebig", "", ""); got != roleNone {
		t.Fatalf("classifyToken(\"beliebig\", \"\", \"\") = %v (Erwartung: roleNone)", got)
	}
}

// TestClassifyTokenAdminDecktReaderAb trägt die Hierarchie aus `ADR-0057`
// Teilfrage 3 (analog `ADR-0047`).
func TestClassifyTokenAdminDecktReaderAb(t *testing.T) {
	if got := classifyToken(testAdminToken, testReaderToken, testAdminToken); got != roleAdmin {
		t.Fatalf("classifyToken(admin) = %v (Erwartung: roleAdmin)", got)
	}
	if roleAdmin < roleReader {
		t.Fatalf("roleAdmin trägt nicht die höhere Rechtsklasse")
	}
}

// TestServerHandlerOhneStartLiefertRouting trägt den Grundzug aus
// `New`/`Handler` ohne einen realen Socket zu öffnen: der verdrahtete
// Handler antwortet bereits vor jedem `Start`-Aufruf — `Handler()` trägt
// keinen Seiteneffekt.
func TestServerHandlerOhneStartLiefertRouting(t *testing.T) {
	srv := New(Config{
		Addr:             "127.0.0.1:0",
		TokenReader:      testReaderToken,
		TokenAdmin:       testAdminToken,
		RegisterConsumer: newFakeRegisterConsumerUseCase(),
	})
	req := httptest.NewRequest(http.MethodPost, "/consumers", strings.NewReader(`{"consumer_id":"c1","name":"Consumer 1"}`))
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401 — kein Token gesetzt)", rec.Code)
	}
}

// TestServerStartMeldetBindfehler trägt den Bind-Fehlerpfad des
// Lebenszyklus (`SPEC-008`): eine unbrauchbare Horch-Adresse endet als
// Fehler von `Start`, nicht in einem stillen Lauf ohne Listener. `net.Listen`
// weist die Adresse ohne Port zurück, bevor irgendein Socket entsteht — der
// Test läuft deshalb netzlos.
// Rot färbende Mutation: in `Start` das `err != nil` verwerfen und immer
// `nil` zurückgeben.
func TestServerStartMeldetBindfehler(t *testing.T) {
	srv := New(Config{Addr: "127.0.0.1", TokenReader: testReaderToken, TokenAdmin: testAdminToken})
	err := srv.Start()
	if err == nil {
		t.Fatal("Start mit unbrauchbarer Horch-Adresse liefert keinen Fehler")
	}
	if !strings.Contains(err.Error(), "127.0.0.1") {
		t.Fatalf("Fehler nennt die Horch-Adresse nicht: %v", err)
	}
}

// TestServerShutdownVorStartIstRegulaererAusgang trägt die zweite Hälfte des
// Lebenszyklus: `http.ErrServerClosed` ist der reguläre Ausgang eines
// geordneten `Shutdown`, keine Fehlerklasse (`SPEC-008`). Ein vor dem Start
// beendeter Server kehrt aus `Start` deshalb ohne Fehler zurück — der Test
// braucht dafür keinen Socket und keinen zweiten Lauf.
// Rot färbende Mutation: in `Start` `err != nil` statt
// `err != http.ErrServerClosed` prüfen — dann wird `ErrServerClosed` als
// Fehler zurückgegeben.
func TestServerShutdownVorStartIstRegulaererAusgang(t *testing.T) {
	srv := New(Config{Addr: "127.0.0.1:0", TokenReader: testReaderToken, TokenAdmin: testAdminToken})
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown ohne laufenden Server: %v", err)
	}
	if err := srv.Start(); err != nil {
		t.Fatalf("Start nach Shutdown: %v (Erwartung: kein Fehler — ErrServerClosed ist kein Fehlerausgang)", err)
	}
}
