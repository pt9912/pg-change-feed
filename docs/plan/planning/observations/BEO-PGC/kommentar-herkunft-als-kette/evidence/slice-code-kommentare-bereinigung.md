**Vorgang:** slice-code-kommentare-bereinigung (Bereinigung der Go-Kommentare in acht Tranchen; zwei Reviews mit Fixrunden)

**Fund:** Zwei neue Klassen der Lese-Handlung (kein Werkzeug-Kandidat, kein Gate fängt sie):

1. „Kürzung lässt hängende Herkunfts-Referenz zurück“ — **4×** (B1: 1 HIGH + 2 MEDIUM; B2: 1 MEDIUM;
   dazu 2 LOW als milde Formen derselben Klasse). Das Kürzen auf einen Anker entfernte den tragenden
   Anker und ließ Teil-Referenzen stehen, die kein Leser auflösen kann: „die Decision Festlegung 2“
   (`internal/bootstrap/config_file.go:6`), „(Teilfrage 4/5)“
   (`tools/harness/natssub/main.go:10`), „(Festlegung 2)“ (`examples/nats-client/main.go:62`),
   „(Konvergenz)“ (`test/integration/integration_test.go:1008`); milde Formen: „die Konvergenz“ und
   „der korrigierte Ursprungstext“ hinter einem resolvierten Anker. `make kommentar-kennungen` meldet
   keine dieser Stellen — das Werkzeug zählt Kennungen, nicht hängende Verweise.
2. „ungenau umgeformte Spec-Wiedergabe“ — **1×** (B1 F-4, `translate.go`): die Umformulierung nannte
   eine „Spalten-Ordnung der Antrags-Tabelle“, die die Switch-Reihenfolge nicht wieder gibt — das
   Umformulieren selbst ist eine Lese-Handlung gegen Schema und Spec, nicht nur das Kürzen.

Die bereinigte Restmenge ist strukturell: **14** Kandidaten, alle Godocs von `func TestE2E*` in
`test/integration`, deren LH-/SPEC-Kennungs-Menge der Generator der E2E-Abdeckungstabelle maschinell
liest (`abdeckungsAdressiert`) — Kürzen auf einen Anker ließe Abdeckungszeilen entfallen. Das ist die
bejahte offene Frage „Grenzfälle mit zwei Ankern“ dieses Eintrags: der Fall ist real, die
Konkretisierung von `AGENTS.md` §3.7 (Erzeugnis-Eingabe ist Abdeckung, nicht Herkunft) ist im
Steering-Loop des Slice verkörpert.

Quelle: Review-Reports `review-slice-code-kommentare-bereinigung-b1` (F-1/F-2/F-3/F-4) und
`review-slice-code-kommentare-bereinigung-b2` (F-1/F-2/F-3/F-5) · Fixrunden `89d347db`/`933ea5c0` ·
Verifikations-Report `verify-slice-code-kommentare-bereinigung` (§1 Restmenge 14, §3 Fixrunden).
