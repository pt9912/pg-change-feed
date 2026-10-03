package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// errorBody liest den Fehlerkörper einer Antwort als Schlüssel-Wert-Abbildung;
// ein Schlüssel, der nicht im Körper steht, fehlt auch in der Abbildung.
func errorBody(t *testing.T, resp *http.Response) map[string]string {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	var decoded map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Fehlerkörper dekodieren: %v", err)
	}
	return decoded
}

// TestFehlerkoerperTraegtDenCodeJeFehlerstelle trägt: jede Fehlerstelle des
// HTTP-Adapters schreibt den Code ihrer Ursache in das Feld `code` neben
// `error`, der Statuscode bleibt unverändert. Eingabeseite: die Anfrage, die
// die Stelle erreicht. Rot färbende Mutation je Zeile: an der Stelle den
// Konstanten-Wert gegen eine andere Konstante tauschen.
func TestFehlerkoerperTraegtDenCodeJeFehlerstelle(t *testing.T) {
	tests := []struct {
		name       string
		cfg        Config
		method     string
		path       string
		token      string
		body       string
		wantStatus int
		wantCode   messagecode.Code
	}{
		{"RegisterConsumer ungültiges JSON", Config{}, "POST", "/consumers", testAdminToken, "{", 400, messagecode.RejectedBodyInvalid},
		{"RegisterConsumer leere Kennung", Config{RegisterConsumer: fakeFailingUseCase{err: domainerrors.ErrEmptyIdentifier}}, "POST", "/consumers", testAdminToken, `{}`, 400, messagecode.RejectedRequiredField},
		{"RegisterConsumer klassifizierte Ursache", Config{RegisterConsumer: fakeFailingUseCase{err: fmt.Errorf("registrieren: %w", outbound.ErrConsumerStateStorage)}}, "POST", "/consumers", testAdminToken, `{"consumer_id":"c","name":"n"}`, 500, messagecode.ConsumerStateFailed},
		{"RegisterConsumer unklassifizierter Fehler", Config{RegisterConsumer: fakeFailingUseCase{err: errors.New("speicher nicht erreichbar")}}, "POST", "/consumers", testAdminToken, `{"consumer_id":"c","name":"n"}`, 500, messagecode.InternalFallback},
		{"RegisterConsumer Ablehnung außerhalb der Pflichtfelder bleibt 500", Config{RegisterConsumer: fakeFailingUseCase{err: domainerrors.ErrPositionRegression}}, "POST", "/consumers", testAdminToken, `{"consumer_id":"c","name":"n"}`, 500, messagecode.InternalFallback},
		{"AcknowledgeConsumer ungültiges JSON", Config{}, "POST", "/consumers/acknowledge", testAdminToken, "{", 400, messagecode.RejectedBodyInvalid},
		{"AcknowledgeConsumer Position ohne Offset", Config{AcknowledgeConsumer: fakeAcknowledgeFailingUseCase{err: domainerrors.ErrInvalidPosition}}, "POST", "/consumers/acknowledge", testAdminToken, `{"consumer_id":"c","source_id":"s","offset":0}`, 400, messagecode.RejectedPositionInvalid},
		{"AcknowledgeConsumer Position rückt nicht vor", Config{AcknowledgeConsumer: fakeAcknowledgeFailingUseCase{err: domainerrors.ErrPositionRegression}}, "POST", "/consumers/acknowledge", testAdminToken, `{"consumer_id":"c","source_id":"s","offset":1}`, 400, messagecode.RejectedPositionRegressed},
		{"AcknowledgeConsumer Position einer anderen Quelle", Config{AcknowledgeConsumer: fakeAcknowledgeFailingUseCase{err: domainerrors.ErrSourceMismatch}}, "POST", "/consumers/acknowledge", testAdminToken, `{"consumer_id":"c","source_id":"s","offset":1}`, 400, messagecode.RejectedPositionSource},
		{"GetConsumerPosition ohne consumer_id", Config{}, "GET", "/consumers/position", testReaderToken, "", 400, messagecode.RejectedRequiredField},
		{"GetConsumerPosition klassifizierte Ursache", Config{GetConsumerPosition: fakeGetConsumerPositionFailingUseCase{err: outbound.ErrConsumerStateStorage}}, "GET", "/consumers/position?consumer_id=c", testReaderToken, "", 500, messagecode.ConsumerStateFailed},
		{"RemoveConsumer ungültiges JSON", Config{}, "POST", "/consumers/remove", testAdminToken, "{", 400, messagecode.RejectedBodyInvalid},
		{"RemoveConsumer leere Kennung", Config{RemoveConsumer: fakeRemoveConsumerFailingUseCase{err: domainerrors.ErrEmptyIdentifier}}, "POST", "/consumers/remove", testAdminToken, `{}`, 400, messagecode.RejectedRequiredField},
		{"RunRetention ungültiges JSON", Config{}, "POST", "/retention/run", testAdminToken, "{", 400, messagecode.RejectedBodyInvalid},
		{"RunRetention negative Dauer", Config{RunRetention: fakeRunRetentionUseCase{}}, "POST", "/retention/run", testAdminToken, `{"source":"s","min_age_nanos":-1}`, 400, messagecode.RejectedValueInvalid},
		{"EnableTable ungültiges JSON", Config{}, "POST", "/tables/enable", testAdminToken, "{", 400, messagecode.RejectedBodyInvalid},
		{"EnableTable Tabelle fehlt an der Quelle", Config{EnableTable: fakeEnableTableUseCase{tableExists: false, enabled: map[string]bool{}}}, "POST", "/tables/enable", testAdminToken,
			`{"source":"s","schema":"public","table":"t","table_id":"tid","schema_version_id":"v","version":1,"publication":"p"}`, 404, messagecode.RejectedTableMissing},
		{"EnableTable Version kleiner 1", Config{EnableTable: fakeEnableTableUseCase{tableExists: true, enabled: map[string]bool{}}}, "POST", "/tables/enable", testAdminToken,
			`{"source":"s","schema":"public","table":"t","table_id":"tid","schema_version_id":"v","version":0,"publication":"p"}`, 400, messagecode.RejectedValueInvalid},
		{"DisableTable ungültiges JSON", Config{}, "POST", "/tables/disable", testAdminToken, "{", 400, messagecode.RejectedBodyInvalid},
		{"GetStatus ohne Pflichtparameter", Config{}, "GET", "/tables/status?source=s", testReaderToken, "", 400, messagecode.RejectedRequiredField},
		{"ListTables ohne Pflichtparameter", Config{}, "GET", "/tables?source=s", testReaderToken, "", 400, messagecode.RejectedRequiredField},
		{"ReadChanges unbekannter Parameter", Config{}, "GET", "/changes?source=s&frm=1", testReaderToken, "", 400, messagecode.RejectedParameterUnknown},
		{"ReadChanges ohne source", Config{}, "GET", "/changes", testReaderToken, "", 400, messagecode.RejectedRequiredField},
		{"ReadChanges from keine Ganzzahl", Config{}, "GET", "/changes?source=s&from=abc", testReaderToken, "", 400, messagecode.RejectedValueInvalid},
		{"ReadChanges limit keine Ganzzahl", Config{}, "GET", "/changes?source=s&limit=abc", testReaderToken, "", 400, messagecode.RejectedValueInvalid},
		{"ReadChanges from unter 1", Config{}, "GET", "/changes?source=s&from=0", testReaderToken, "", 400, messagecode.RejectedPositionInvalid},
		{"ReadChanges limit unter 1", Config{ReadChanges: &fakeReadChangesUseCase{err: outbound.ErrNonPositiveLimit}}, "GET", "/changes?source=s&limit=0", testReaderToken, "", 400, messagecode.RejectedValueInvalid},
		{"ReadChanges Endposition vor Startposition", Config{ReadChanges: &fakeReadChangesUseCase{err: outbound.ErrRangeInverted}}, "GET", "/changes?source=s&from=5&to=2", testReaderToken, "", 400, messagecode.RejectedRangeInverted},
		{"Diagnose unbekannter Parameter", Config{}, "GET", "/diagnose?source=s&x=1", testReaderToken, "", 400, messagecode.RejectedParameterUnknown},
		{"Diagnose ohne source", Config{}, "GET", "/diagnose", testReaderToken, "", 400, messagecode.RejectedRequiredField},
		{"Diagnose klassifizierte Ursache", Config{Diagnose: &fakeDiagnoseUseCase{err: outbound.ErrDiagnosticsStorage}}, "GET", "/diagnose?source=s", testReaderToken, "", 500, messagecode.DiagnosticsReadFailed},
		{"SSE ohne Broadcaster", Config{}, "GET", "/changes/stream", testReaderToken, "", 503, messagecode.WiringPrecondition},
		{"SSE unbekannter Parameter", Config{Subscriber: newFakeChangeSubscriber(0)}, "GET", "/changes/stream?x=1", testReaderToken, "", 400, messagecode.RejectedParameterUnknown},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newDefaultTestServer(t, tc.cfg)
			resp := doRequest(t, ts, tc.method, tc.path, tc.token, tc.body)
			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("Status: %d (Erwartung: %d)", resp.StatusCode, tc.wantStatus)
			}
			decoded := errorBody(t, resp)
			if decoded["error"] == "" {
				t.Fatalf("Körper trägt keinen Klartext: %v", decoded)
			}
			if decoded["code"] != string(tc.wantCode) {
				t.Fatalf("code = %q, Erwartung %q (Körper %v)", decoded["code"], tc.wantCode, decoded)
			}
		})
	}
}

// TestFehlerkoerperOhneCodeZuordnungTraegtKeinFeldCode trägt: ein
// fehlender oder unbekannter Token (`401`) und eine unzureichende Rechtsklasse
// (`403`) tragen keinen Code — die Tabelle der Codes führt für sie keinen
// Eintrag; das Feld `code` fehlt im Körper, es steht nicht leer darin.
func TestFehlerkoerperOhneCodeZuordnungTraegtKeinFeldCode(t *testing.T) {
	ts := newDefaultTestServer(t, Config{RegisterConsumer: newFakeRegisterConsumerUseCase()})
	for _, tc := range []struct {
		name       string
		token      string
		wantStatus int
	}{
		{"kein Token", "", http.StatusUnauthorized},
		{"unbekannter Token", "unbekannt", http.StatusUnauthorized},
		{"Reader-Token gegen Admin-Endpunkt", testReaderToken, http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, ts, "POST", "/consumers", tc.token, `{"consumer_id":"c","name":"n"}`)
			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("Status: %d (Erwartung: %d)", resp.StatusCode, tc.wantStatus)
			}
			decoded := errorBody(t, resp)
			if _, present := decoded["code"]; present {
				t.Fatalf("Körper trägt ein Feld code: %v", decoded)
			}
			if decoded["error"] == "" {
				t.Fatalf("Körper trägt keinen Klartext: %v", decoded)
			}
		})
	}
}

// noFlushWriter trägt einen `http.ResponseWriter` ohne `http.Flusher`.
type noFlushWriter struct{ http.ResponseWriter }

// TestStreamOhneFlusherTraegtDenRueckfallDerKlasseInternal trägt: ein
// Antwort-Writer ohne `http.Flusher` endet mit `500` und dem Rückfall der
// Klasse `internal`. Eingabeseite: der Writer.
func TestStreamOhneFlusherTraegtDenRueckfallDerKlasseInternal(t *testing.T) {
	srv := New(Config{Addr: "unused:0", TokenReader: testReaderToken, TokenAdmin: testAdminToken, Subscriber: newFakeChangeSubscriber(0)})
	req := httptest.NewRequest(http.MethodGet, "/changes/stream", nil)
	req.Header.Set("Authorization", "Bearer "+testReaderToken)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(noFlushWriter{rec}, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Status: %d (Erwartung: 500)", rec.Code)
	}
	var decoded map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("Körper dekodieren: %v", err)
	}
	if decoded["code"] != string(messagecode.InternalFallback) {
		t.Fatalf("code = %q, Erwartung %q", decoded["code"], messagecode.InternalFallback)
	}
}

// TestInternerFehlerWarntMitDemCodeDerAnfrage trägt: jede der drei
// Warnstellen eines gescheiterten Aufrufs (Use-Case-Fehler in einem Handler,
// Fehler der Registrierung) schreibt das Attribut `code` mit `PCF-W4008`;
// nur ein Fehler mit `500` warnt, eine Ablehnung mit `400` nicht.
func TestInternerFehlerWarntMitDemCodeDerAnfrage(t *testing.T) {
	tests := []struct {
		name      string
		cfg       func(log *warnCodeLog) Config
		method    string
		path      string
		body      string
		wantCodes []string
	}{
		{"writeDomainError", func(l *warnCodeLog) Config {
			return Config{Log: l, GetConsumerPosition: fakeGetConsumerPositionFailingUseCase{err: errors.New("speicher")}}
		}, "GET", "/consumers/position?consumer_id=c", "", []string{string(messagecode.WarnAPIRequestFailed)}},
		{"RegisterConsumer", func(l *warnCodeLog) Config {
			return Config{Log: l, RegisterConsumer: fakeFailingUseCase{err: errors.New("speicher")}}
		}, "POST", "/consumers", `{"consumer_id":"c","name":"n"}`, []string{string(messagecode.WarnAPIRequestFailed)}},
		{"Ablehnung warnt nicht", func(l *warnCodeLog) Config {
			return Config{Log: l, GetConsumerPosition: fakeGetConsumerPositionFailingUseCase{err: inbound.ErrSourceTableMissing}}
		}, "GET", "/consumers/position?consumer_id=c", "", nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			log := &warnCodeLog{}
			ts := newDefaultTestServer(t, tc.cfg(log))
			token := testReaderToken
			if tc.method == "POST" {
				token = testAdminToken
			}
			resp := doRequest(t, ts, tc.method, tc.path, token, tc.body)
			_ = resp.Body.Close()
			got := log.snapshot()
			if fmt.Sprint(got) != fmt.Sprint(tc.wantCodes) {
				t.Fatalf("Attribut code der Warnungen = %q, Erwartung %q", got, tc.wantCodes)
			}
		})
	}
}
