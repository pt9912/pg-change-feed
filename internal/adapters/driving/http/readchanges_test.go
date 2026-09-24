package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeReadChangesUseCase trägt eine In-Memory-Fälschung des Inbound Ports
// (`ADR-0030`): Whitebox-Test des Adapters ohne reale Persistenz. Die
// Fehler-Hälfte des Port-Kontrakts (`outbound.ErrNonPositiveLimit`,
// `outbound.ErrRangeInverted`) liefert der Fake als `err` — dass der reale
// Lesepfad diese Sentinels für `limit < 1`/`from > to` hervorbringt, trägt
// `internal/application/usecase/readchanges` (dort gegen den Port-Kontrakt
// geprüft); hier steht die Mapping-Hälfte des Adapters.
type fakeReadChangesUseCase struct {
	result  inbound.ReadChangesResult
	err     error
	queries []inbound.ReadChangesQuery
}

func (f *fakeReadChangesUseCase) ReadChanges(_ context.Context, query inbound.ReadChangesQuery) (inbound.ReadChangesResult, error) {
	f.queries = append(f.queries, query)
	if f.err != nil {
		return inbound.ReadChangesResult{}, f.err
	}
	return f.result, nil
}

func newReadChangesServer(t *testing.T, useCase inbound.ReadChangesUseCase) *httptest.Server {
	t.Helper()
	srv := New(Config{
		Addr:        "unused:0",
		TokenReader: testReaderToken,
		TokenAdmin:  testAdminToken,
		ReadChanges: useCase,
	})
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

func getChanges(t *testing.T, ts *httptest.Server, token, path string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, ts.URL+path, nil)
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

func body(t *testing.T, resp *http.Response) string {
	t.Helper()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Body lesen: %v", err)
	}
	return string(raw)
}

// TestReadChangesOhneTokenEndetMit401 trägt die Authn-Grenze des Endpunkts
// (`SPEC-022`, `ADR-0057` Teilfrage 3): ohne Bearer-Token kein Lesezugriff.
func TestReadChangesOhneTokenEndetMit401(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	resp := getChanges(t, ts, "", "/changes?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401)", resp.StatusCode)
	}
}

// TestReadChangesUnbekannterTokenEndetMit401 trägt denselben Pfad für ein
// Token, das keiner konfigurierten Klasse entspricht.
func TestReadChangesUnbekannterTokenEndetMit401(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	resp := getChanges(t, ts, "unbekannt", "/changes?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401)", resp.StatusCode)
	}
}

// TestReadChangesReaderTokenLiestChanging trägt die Rechtsklasse
// (`ADR-0081` Teilfrage 2): das `reader`-Token genügt für den lesenden
// Endpunkt; die Antwort trägt die übersetzten Felder samt
// Klartext-Identität der Tabelle.
func TestReadChangesReaderTokenLiestChanging(t *testing.T) {
	committedAt := time.Date(2026, 9, 15, 10, 30, 0, 123456789, time.UTC)
	position, err := model.NewSourcePosition("src-1", 4711)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	change, err := model.NewChange(
		model.ChangeID("c-1"), model.TransactionID("tx-1"), model.SourceTableID("tbl-1"),
		3, model.OperationUpdate, []byte(`{"id":1}`), []byte(`{"id":1,"name":"neu"}`), model.SchemaVersionID("sv-1"),
	)
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	useCase := &fakeReadChangesUseCase{result: inbound.ReadChangesResult{Changes: []inbound.ReadChange{{
		Position:    position,
		Change:      change,
		CommittedAt: model.NewTimePoint(committedAt.UnixNano()),
	}}}}
	ts := newReadChangesServer(t, useCase)

	resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1&schema=public&table=orders&from=100&to=500&limit=10")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200), Body: %s", resp.StatusCode, body(t, resp))
	}
	if got := len(useCase.queries); got != 1 {
		t.Fatalf("Use-Case-Aufrufe = %d, wollen 1", got)
	}
	query := useCase.queries[0]
	if query.Source != "src-1" || query.Schema != "public" || query.Table != "orders" {
		t.Fatalf("Abfrage = %+v, wollen Quelle public/orders", query)
	}
	if query.Start == nil || query.Start.Offset != 100 || query.End == nil || query.End.Offset != 500 {
		t.Fatalf("Bereich = %+v, wollen [100,500)", query)
	}
	if query.Limit == nil || *query.Limit != 10 {
		t.Fatalf("Limit = %v, wollen 10", query.Limit)
	}

	var decoded readChangesResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if len(decoded.Changes) != 1 {
		t.Fatalf("Changes = %+v, wollen einen Eintrag", decoded.Changes)
	}
	entry := decoded.Changes[0]
	if entry.CommitPosition != 4711 || entry.ChangeID != "c-1" || entry.TransactionID != "tx-1" ||
		entry.SourceTableID != "tbl-1" || entry.Schema != "public" || entry.Table != "orders" ||
		entry.Sequence != 3 || entry.Operation != "UPDATE" || entry.SchemaVersion != "sv-1" {
		t.Fatalf("Antwort-Eintrag = %+v", entry)
	}
	if string(entry.OldImage) != `{"id":1}` || string(entry.NewImage) != `{"id":1,"name":"neu"}` {
		t.Fatalf("Row Images = %s / %s", entry.OldImage, entry.NewImage)
	}
	if entry.CommittedAt != committedAt.Format(time.RFC3339Nano) {
		t.Fatalf("committed_at = %q, wollen %q", entry.CommittedAt, committedAt.Format(time.RFC3339Nano))
	}
}

// TestReadChangesAdminTokenLiest trägt die Abdeckung der niedrigeren Klasse
// durch das `admin`-Token (`ADR-0057` Teilfrage 3).
func TestReadChangesAdminTokenLiest(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	resp := getChanges(t, ts, testAdminToken, "/changes?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
}

// TestReadChangesKeinTrefferEndetMit200Leer trägt die Boundary
// (`LH-FA-REA-006`, `ADR-0081` Teilfrage 4): kein Treffer ist kein Fehler —
// `200` mit leerer, gesetzter Liste, nie `404`, nie `null`.
func TestReadChangesKeinTrefferEndetMit200Leer(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1&table=fehlt")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200, kein 404)", resp.StatusCode)
	}
	if got := body(t, resp); got != "{\"changes\":[]}\n" {
		t.Fatalf("Body = %q, wollen die leere, gesetzte Liste", got)
	}
}

// TestReadChangesUnbekannterParameterEndetMit400 trägt die Verschärfung
// gegenüber den neun Bestandsendpunkten (`ADR-0081` Teilfrage 4): ein
// unbekannter **Filter** änderte sonst den Ergebnisstand still.
func TestReadChangesUnbekannterParameterEndetMit400(t *testing.T) {
	useCase := &fakeReadChangesUseCase{}
	ts := newReadChangesServer(t, useCase)
	resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1&tabell=orders")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
	raw := body(t, resp)
	if !strings.Contains(raw, "tabell") {
		t.Fatalf("Fehlertext nennt den Parameter nicht: %s", raw)
	}
	if len(useCase.queries) != 0 {
		t.Fatalf("Use-Case-Aufrufe = %d, wollen 0", len(useCase.queries))
	}
}

// TestReadChangesFehlendeQuelleEndetMit400 trägt das Pflichtfeld
// (`ADR-0081` Teilfrage 4): ohne `source` liest der Aufruf nicht.
func TestReadChangesFehlendeQuelleEndetMit400(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	resp := getChanges(t, ts, testReaderToken, "/changes?table=orders")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// TestReadChangesUnlesbareZahlEndetMit400 trägt die Formgrenze der
// Zahlen-Parameter (`ADR-0081` Teilfrage 4).
func TestReadChangesUnlesbareZahlEndetMit400(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	for _, path := range []string{
		"/changes?source=src-1&from=abc",
		"/changes?source=src-1&to=abc",
		"/changes?source=src-1&limit=abc",
	} {
		resp := getChanges(t, ts, testReaderToken, path)
		status := resp.StatusCode
		resp.Body.Close()
		if status != http.StatusBadRequest {
			t.Fatalf("%s: Status %d (Erwartung: 400)", path, status)
		}
	}
}

// TestReadChangesPositionUnterEinsEndetMit400 trägt die Positions-Grenze
// (`SPEC-003`): eine Position 0 existiert nicht.
func TestReadChangesPositionUnterEinsEndetMit400(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	for _, path := range []string{
		"/changes?source=src-1&from=0",
		"/changes?source=src-1&to=0",
		"/changes?source=src-1&from=-5",
	} {
		resp := getChanges(t, ts, testReaderToken, path)
		status := resp.StatusCode
		resp.Body.Close()
		if status != http.StatusBadRequest {
			t.Fatalf("%s: Status %d (Erwartung: 400)", path, status)
		}
	}
}

// TestReadChangesNichtpositivesLimitEndetMit400 trägt die zwei Hälften des
// Adapter-Anteils an `LH-FA-REA-003` Negative. Ein nicht positives Limit
// wird **unverändert** an den Port weitergereicht — der Adapter verwirft es
// nicht, sonst läse der Aufruf still unbegrenzt statt zu enden; die Grenze
// selbst trägt der Port-Kontrakt (`outbound.ChangeQuery.Validate`), dessen
// Sentinel `outbound.ErrNonPositiveLimit` der Adapter auf `400` abbildet.
func TestReadChangesNichtpositivesLimitEndetMit400(t *testing.T) {
	for _, tc := range []struct {
		path  string
		limit int
	}{
		{path: "/changes?source=src-1&limit=0", limit: 0},
		{path: "/changes?source=src-1&limit=-5", limit: -5},
	} {
		useCase := &fakeReadChangesUseCase{err: outbound.ErrNonPositiveLimit}
		ts := newReadChangesServer(t, useCase)
		resp := getChanges(t, ts, testReaderToken, tc.path)
		status := resp.StatusCode
		resp.Body.Close()
		if status != http.StatusBadRequest {
			t.Fatalf("%s: Status %d (Erwartung: 400)", tc.path, status)
		}
		if len(useCase.queries) != 1 {
			t.Fatalf("%s: Use-Case-Aufrufe = %d, wollen 1", tc.path, len(useCase.queries))
		}
		got := useCase.queries[0].Limit
		if got == nil {
			t.Fatalf("%s: Limit ist nicht gesetzt — der Adapter verwarf den Wert, statt ihn dem Port-Kontrakt zu übergeben", tc.path)
		}
		if *got != tc.limit {
			t.Fatalf("%s: Limit = %d, wollen %d", tc.path, *got, tc.limit)
		}
	}
}

// TestReadChangesInvertierterBereichEndetMit400 trägt dieselbe
// Mapping-Hälfte für den Bereich (`LH-FA-REA-001` Negative):
// `outbound.ErrRangeInverted` endet als `400`.
func TestReadChangesInvertierterBereichEndetMit400(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{err: outbound.ErrRangeInverted})
	resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1&from=500&to=100")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// TestReadChangesSpeicherfehlerEndetMit500 trägt die Abgrenzung zum
// Port-Kontrakt (`ADR-0081` Teilfrage 4): ein Store-Fehler der Klasse
// `storage` endet als `500`, nicht als `400`.
func TestReadChangesSpeicherfehlerEndetMit500(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{err: outbound.ErrStorage})
	resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Status: %d (Erwartung: 500)", resp.StatusCode)
	}
}

// TestReadChangesFehlendesBildIstNull trägt die Bildform (`LH-FA-CAP-008`
// Boundary): ein fehlendes Row Image steht als `null`, nicht als leerer
// String.
func TestReadChangesFehlendesBildIstNull(t *testing.T) {
	position, err := model.NewSourcePosition("src-1", 1)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	change, err := model.NewChange(
		model.ChangeID("c-1"), model.TransactionID("tx-1"), model.SourceTableID("tbl-1"),
		1, model.OperationDelete, []byte(`{"id":1}`), nil, model.SchemaVersionID("sv-1"),
	)
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{result: inbound.ReadChangesResult{
		Changes: []inbound.ReadChange{{Position: position, Change: change, CommittedAt: model.NewTimePoint(1)}},
	}})
	resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1")
	defer resp.Body.Close()

	raw := body(t, resp)
	if !strings.Contains(raw, `"old_image":{"id":1}`) || !strings.Contains(raw, `"new_image":null`) {
		t.Fatalf("Body = %s, wollen old_image gesetzt und new_image null", raw)
	}
}

// TestReadChangesUndStreamKollidierenNicht trägt die Koexistenz der beiden
// Formen desselben Gegenstands (`ADR-0081` Teilfrage 2): `GET /changes/stream`
// bleibt erreichbar — ohne verdrahteten Broadcaster mit `503` — und
// `GET /changes` läuft daneben.
func TestReadChangesUndStreamKollidierenNicht(t *testing.T) {
	ts := newReadChangesServer(t, &fakeReadChangesUseCase{})
	stream := getChanges(t, ts, testReaderToken, "/changes/stream")
	defer stream.Body.Close()
	if stream.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("/changes/stream: Status %d (Erwartung: 503 ohne Broadcaster)", stream.StatusCode)
	}
	read := getChanges(t, ts, testReaderToken, "/changes?source=src-1")
	defer read.Body.Close()
	if read.StatusCode != http.StatusOK {
		t.Fatalf("/changes: Status %d (Erwartung: 200)", read.StatusCode)
	}
}

// TestReadChangesLeererParameterBestandIstGueltig trägt die Trennschärfe
// der Parameter-Prüfung: ein Aufruf ohne die optionalen Filter ist gültig
// — der unbekannte Parameter ist die Grenze, nicht die Auslassung.
func TestReadChangesLeererParameterBestandIstGueltig(t *testing.T) {
	useCase := &fakeReadChangesUseCase{}
	ts := newReadChangesServer(t, useCase)
	resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	if got := useCase.queries[0]; got.Schema != "" || got.Table != "" || got.Start != nil || got.End != nil || got.Limit != nil {
		t.Fatalf("Abfrage = %+v, wollen leere, ungesetzte Filter", got)
	}
}

// TestReadChangesOhneVerdrahtungEndetVorDemHandler hält fest, dass der
// Endpunkt ohne verdrahteten Use Case an der Middleware endet (`401`),
// bevor der Handler läuft: die Verdrahtung trägt den Use Case, dieser Test
// belegt nur die Reihenfolge.
func TestReadChangesOhneVerdrahtungEndetVorDemHandler(t *testing.T) {
	srv := New(Config{Addr: "127.0.0.1:0", TokenReader: testReaderToken, TokenAdmin: testAdminToken})
	req := httptest.NewRequest(http.MethodGet, "/changes?source=src-1", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401 — kein Token gesetzt)", rec.Code)
	}
}

// Der Fake erfüllt den Inbound Port; die Zusicherung hält die Fake-Seite
// gegen den Port-Vertrag.
var _ inbound.ReadChangesUseCase = (*fakeReadChangesUseCase)(nil)

// TestReadChangesTraegtOriginAlsLetztesFeld trägt das Antwortfeld `origin`
// (`SPEC-022`, `LH-FA-CAP-009`): `wal` und `backfill` stehen so in der
// Antwort, wie der Change sie trägt; ein fehlender Wert (ein Change ohne
// gesetzte Herkunft) liest als `wal`. Das Feld steht als letztes — die
// Reihenfolge der übrigen Felder bleibt.
func TestReadChangesTraegtOriginAlsLetztesFeld(t *testing.T) {
	position, err := model.NewSourcePosition("src-1", 1)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	base, err := model.NewChange(
		model.ChangeID("c-1"), model.TransactionID("tx-1"), model.SourceTableID("tbl-1"),
		1, model.OperationInsert, nil, []byte(`{"id":1}`), model.SchemaVersionID("sv-1"),
	)
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	backfill, err := base.WithOrigin(model.ChangeOriginBackfill)
	if err != nil {
		t.Fatalf("WithOrigin: %v", err)
	}
	missing := base
	missing.Origin = ""

	for _, tc := range []struct {
		name   string
		change model.Change
		want   string
	}{
		{"wal", base, "wal"},
		{"backfill", backfill, "backfill"},
		{"fehlender Wert liest als wal", missing, "wal"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ts := newReadChangesServer(t, &fakeReadChangesUseCase{result: inbound.ReadChangesResult{
				Changes: []inbound.ReadChange{{Position: position, Change: tc.change, CommittedAt: model.NewTimePoint(1)}},
			}})
			resp := getChanges(t, ts, testReaderToken, "/changes?source=src-1")
			defer resp.Body.Close()

			raw := strings.TrimSpace(body(t, resp))
			if suffix := `,"origin":"` + tc.want + `"}]}`; !strings.HasSuffix(raw, suffix) {
				t.Fatalf("Body = %s, will Endung %s (origin als letztes Feld)", raw, suffix)
			}
			if !strings.Contains(raw, `"schema_version":"sv-1","committed_at":"`) {
				t.Fatalf("Body = %s, will die Feldreihenfolge schema_version, committed_at, origin", raw)
			}
		})
	}
}
