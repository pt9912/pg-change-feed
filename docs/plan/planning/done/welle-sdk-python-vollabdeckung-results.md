# Welle welle-sdk-python-vollabdeckung — Closure-Notiz

**Welle:** welle-sdk-python-vollabdeckung
**Abschluss:** 2026-09-23
**Verantwortlich:** pt9912

## Was wurde geliefert?

[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2/Folgepflicht 1 eingelöst — `pgchangefeed`
erreicht die volle Vier-Wege-Parität mit der verschärften Test-Pflicht
(jede neue Fläche prüft ihre Protokoll-Annahmen zusätzlich gegen eine
reale, laufende Server-Instanz) in drei strikt sequentiellen Slices:

- **`slice-sdk-python-grpc-client-flaeche`**: öffentliche gRPC-Client-Fläche
  (`SPEC-020`) — `PgChangeFeedGrpcStreamClient`, alle zehn
  `SPEC-021`-Nachrichtenfelder; führt zugleich das Realserver-
  Integrationstest-Werkzeug ein (`tools/harness/run-sdk-python-integration-tests.sh`,
  `make test-sdk-python-integration`).
- **`slice-sdk-python-sse-client-flaeche`**: öffentliche SSE-Client-Fläche
  (`SPEC-021`) — `PgChangeFeedSseClient`, eigener Frame-Parser, erweitert
  das Werkzeug um die zweite Phasen-Fläche.
- **`slice-sdk-python-nats-stream-client-flaeche`**: öffentliche
  NATS-Vollinhalts-Client-Fläche (`SPEC-024`) —
  `PgChangeFeedNatsStreamClient`, Subjekt-Namensraum
  `cdc.stream.<source_id>.>`, alle zehn `SPEC-024`-Nachrichtenfelder;
  bündelt zusätzlich die Version-Hebung (`0.1.0` → `0.2.0`, PEP-440-Minor-Inkrement)
  und den Träger-Nachzug ([`AGENTS.md`](../../../../AGENTS.md) §3.13) in
  `spec/pflichtenheft.md` (`LH-FA-SST-009.a`, `SPEC-027` — §1, §6
  Vertragszeile, §7 Historie) und `docs/user/benutzerhandbuch.md` (1.43)
  für **alle drei** neu gelieferten Flächen.

Ein real neu gebautes Artefakt-Paar trägt alle vier Client-Flächen in
einem Package: `pgchangefeed-0.2.0-py3-none-any.whl` (20503 Bytes) und
`pgchangefeed-0.2.0.tar.gz` (24655 Bytes, `sdks/python/dist/`, Stempel
2026-09-23 08:50) — frischer Neubau nach der zweiten Fixrunde, der
Wheel-Inhalt gegen den HEAD-Quellstand byte-gleich gehalten
(`nats_stream_client.py`, `exceptions.py`, `options.py` je diff-leer,
METADATA `Version: 0.2.0`); die Abweichung V-1 der Verifikation
(Artefakt-Frische) ist damit real gelöst. [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
gilt für Python damit als vollständig für die Vier-Wege-Matrix erfüllt —
kein isolierter Flächen-Slice allein hätte das belegt (Welle-Ziel §1).

`make gates` grün auf dem Endstand (921 Dateien, 0 `docs-check`-Befunde,
Coverage 82,70 % ≥ 80 %-Endstufe, `a-check` 0 Befunde,
`generated-sync`/`commit-traceability`/`baseline-verify` je ohne Befund),
ungepiped geprüft nach jedem Commit dieser Closure ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
Kein realer `sdk-python-v0.2.0`-Tag-Push (`git tag -l "sdk-python-v*"`
liefert ausschließlich `sdk-python-v0.1.0` — eigenständig gemessen zu
dieser Closure): diese Welle liefert bewusst nur das paketierbare
Ergebnis, wie im Closure-Trigger (§3 der Welle-Datei) vorab festgelegt.

## Was hat funktioniert?

- Die strikte Sequenz trug sich: das Realserver-Integrationstest-Werkzeug
  entstand **einmal** im gRPC-Slice und wurde von SSE und NATS-Vollinhalt
  um ihre eigene Phasen-Fläche erweitert (`run_surface_phase`-Muster, je
  Phase eigene Testdatei, eigener Sentinel, eigener ID-Bereich, eigene
  Reject-Form) — keine parallele zweite Änderung an derselben
  Infrastruktur-Datei, kein zweiter Bau-Mechanismus.
- Die verschärfte Test-Pflicht aus
  [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 2/Folgepflicht 1 ist real eingelöst: ein grüner
  `make test-sdk-python-integration`-Lauf deckt alle drei Flächen mit
  je einem Empfang am Wire und je einer SQL-Gegenprüfung gegen
  `cdc.changes`; ein frischer, eigenständiger Lauf zu dieser Closure
  (Exit 0) belegt gRPC `change_id=804-1`, SSE `808-1` und NATS `810-1`.
- Die Review-Ketten fingen Träger-Lücken, bevor sie in `done/` landeten:
  jede der drei Flächen lief durch Haupt-Review + Fixrunde; der NATS-Slice
  mit einer zweiten Fixrunde (Re-Review R-1…R-4) — die Suchlauf-Lücken
  (F-4, R-1/R-2) wurden real gezogen, das Artefakt-Staleness (V-1) vom
  Verifier gefunden und durch den Neubau gelöst.
- Das Drei-Commit-Move-Muster ([`AGENTS.md`](../../../../AGENTS.md) §3.3) hielt
  über alle drei Slices und die Welle-Datei sauber; die
  Link-Reconciliationen nach den Moves blieben auf Referenz-Korrekturen
  beschränkt (Muster `a7723fa7`, `83021b39`).

## Was ging anders als geplant?

- Der Träger-Nachzug ([`AGENTS.md`](../../../../AGENTS.md) §3.13) zeigte in
  allen drei Slices dieselbe Lücken-Struktur in Variation: beim gRPC-Slice
  drifteten drei Zahlen-/Stand-Werte in derselben Korrektur-Kette (F-1,
  V-1, F-5 — jeweils vom nächsten Leser nachgezählt), beim SSE-Slice war
  der Suchlauf je Slice statt je bewegter Eigenschaft gebunden
  (Haupt-Review F-2, Re-Review FR-1, Verifikation V-1), beim NATS-Slice
  verfehlte die schmale Plan-Mustersatz-Grep-Form alle fünf Fundstellen
  in ihrer Schreibform und das committete Suchlauf-Feld behauptete zwei
  Nachzüge, die der Baum widerlegte (F-4, R-1/R-2). Die Lücken-Struktur
  ist damit dreifach belegt — siehe Steering-Loop-Einträge unten.
- Das erste 0.2.0-Artefakt-Paar (07:47) war gegen den nach zwei
  Fixrunden weitergeänderten Quellstand byte-stale (V-1 der
  Verifikation) — gelöst durch einen frischen Neubau nach der zweiten
  Fixrunde (08:50), der Wheel-Inhalt gegen den HEAD-Stand byte-gleich
  gehalten. Die Smoke-Aussage (Version + vier Flächen in einem Artefakt)
  war auch am stale Paar wahr; die Byte-Frische erst nach dem Neubau.
- Die zweite Fixrunde des NATS-Slices lief ohne Re-Review-Stufe (V-2);
  die Verifikation schloss die drei Residuen durch eigenes Nachmessen am
  HEAD und trug die Lücke als Prozess-Notiz — der Regelfall für künftige
  Fixrunden bleibt eine Reviewer-Gegenprüfung oder eine explizit
  getragene Verifier-Schließung.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier steht der Lese-Schritt
dieser Closure (Modul 6).

- **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** (verkörpert,
  Anker [`AGENTS.md`](../../../../AGENTS.md) §3.13): zwei neue Belege in
  dieser Welle — der zweiundzwanzigste (`evidence/slice-sdk-python-sse-client-flaeche.md`,
  SSE-Slice: Suchlauf je Slice statt je bewegter Eigenschaft, Kette
  F-2/FR-1/V-1) und der dreiundzwanzigste
  (`evidence/slice-sdk-python-nats-stream-client-flaeche.md`,
  NATS-Slice: Mustersatz-Form verfehlt die Schreibformen, das Suchlauf-Feld
  selbst wird zum Träger, Kette F-4/R-1/R-2). **Lese-Schritt: keine
  `AGENTS.md`-Schärfung nötig.** §3.13s Wortlaut trägt die je-Eigenschaft-Form
  bereits („eine Eigenschaft eines Gegenstands", „`grep` über die Träger
  nach der bewegten Eigenschaft, beide Stände gemessen"); beide
  Welle-Funde sind Anwendungs-Verfehlungen der Regel **als geschrieben**,
  gefangen von den bestehenden Lesern (Reviewer, Fixrunden-Reviewer,
  Verifier) — kein Lücke-tragender Satz, keine doppelte Verkörperung,
  Register-Evidence trägt die Schärfung. Zähler 23× (Datei-Anzahl, real
  ausgezählt), Ausgang bleibt **verkörpert**.
- **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** (verkörpert,
  Anker [`AGENTS.md`](../../../../AGENTS.md) §3.12 + Reviewer-Skill): ein
  neuer Beleg in dieser Welle — der dreizehnte
  (`evidence/slice-sdk-python-grpc-client-flaeche.md`: drei Zahl-/Stand-Driften
  in derselben gRPC-Korrektur-Kette, je vom nächsten Leser nachgezählt).
  Der NATS-F-1-Fund („aktuell `0.1.0`" in der §6-Vertragszeile) trägt
  nach der Klassen-Grenze des Eintrags bei
  `arbeit-ueberholt-stehenden-traeger` (die Zahl war bei ihrer letzten
  Niederschrift wahr und wurde durch diese Arbeit falsch — dual benannt
  im NATS-Review F-1, dort geführt). Zähler 13× (Datei-Anzahl), Ausgang
  bleibt **verkörpert**, kein neuer Handlungsbedarf.
- **`BEO-PGC/zitat-nennt-die-falsche-stelle`** (verkörpert): ein neuer
  Beleg (F-3, `evidence/slice-sdk-python-nats-stream-client-flaeche.md`
  dort, siebte Datei) — kein neuer Regelschärfungs-Anlass, die
  Verkörperung trägt unverändert.
- **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** (verkörpert): ein
  neuer Beleg (F-6/R-1, `evidence/slice-sdk-python-nats-stream-client-flaeche.md`
  dort) — kein neuer Regelschärfungs-Anlass.
- **`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`**
  (verkörpert): geprüft, kein Vorkommen — der Handbuch-Hinweis lag in
  jedem Flächen-Slice als eigener DoD-Punkt und wurde real geliefert
  (1.40/1.41, 1.42, 1.43).

**Unter der Schwelle, unverändert, kein Verkörperungsbedarf:**

- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` —
  weiterhin 1×, kein zweites Auftreten (kein realer Tag-Push,
  Publish-Mechanismus unverändert in dieser Welle).
- `BEO-PGC/test-runner-stiller-ausschluss` — weiterhin 2×, kein
  Vorkommen: der CMD-Guard des SSE-Slices deckt die Runner-Mechanik, die
  NATS-Phase folgt demselben Muster (je Phase eine explizite Testdatei).
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als
  [`AGENTS.md`](../../../../AGENTS.md) §3.10) — kein neuer Beleg: kein
  Workflow dieses Repos wurde in dieser Welle geändert, kein realer
  Tag-Push stattfand.
- `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` — weiterhin 3×,
  kein neuer Beleg in dieser Welle.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/). Kein Eintrag
erreichte in dieser Welle **neu** die 3×-Schwelle — der Lese-Schritt oben
bestätigt für `arbeit-ueberholt-stehenden-traeger`,
`zahl-in-traeger-driftet-gegen-die-messung`,
`zitat-nennt-die-falsche-stelle` und
`beleg-befehl-traegt-seinen-satz-nicht` jeweils bereits verkörperte Stände
mit neuen Belegen, die ohne erneute Verkörperung zählen.

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine
belegte Feststellung, jede Zeile der Welle-§3 mit eigener Messung
geprüft.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` enthält
  ausschließlich `.gitkeep`, kein Carveout referenziert `ADR-0110` oder
  einen der drei Slices dieser Welle (gemessen).
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate`
  steht unverändert bei der 80 %-Endstufe (real gemessen 82,70 %) und
  wurde von keinem der drei Slices berührt (reiner Python-Baum, außerhalb
  der netzlos prüfbaren Go-Fläche `./internal/...`+`./cmd/...`+`./gen/...`);
  `.a-check.yml` `languages: go` liest `sdks/python/**` strukturell nicht
  (real bestätigt: `a-check` meldet 0 Befunde auf dem Endstand).
- **Entscheidung/ADR —
  [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Re-Evaluierungs-Trigger, durch diese Welle **nicht** ausgelöst:**
  - Trigger 1 (ein Python-Referenz-Client entsteht nachträglich) — nicht
    eingetreten: `examples/python/` existiert nicht (Verzeichnis gemessen,
    Verifikation §2.6); `ADR-0110` Festlegung 2 mandatiert keinen
    Beispiel-Baum.
  - Trigger 2 (ein reales Protokoll-Missverständnis im Folge-Release) —
    nicht eingetreten: der rote s3e-Lauf des NATS-Slice ging auf eine
    falsch gebundene Test-Assertion (Ursachen-Klasse), nicht auf eine
    falsch verstandene Stream-Semantik; die `SPEC-024`-Konformität hielt
    Feld für Feld (Verifikation §4).
  - Trigger 3 (`ADR-0107` Triggers 1/2/4: dritte Sprache/dritter
    Vertriebsweg verlangt; reale PyPI-Nutzungsdaten; Trusted-Publishing-
    Umstellung) — nicht eingetreten: kein realer Tag-Push
    (`git tag -l "sdk-python-v*"` liefert ausschließlich
    `sdk-python-v0.1.0` — gemessen), also keine Nutzungsdaten; keine
    dritte Sprache; `.github/workflows/sdk-python-release.yml` wurde von
    keinem der drei Slices verändert.

  [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  bleibt `Accepted` und inhaltlich unberührt durch diese Welle.

## Folge-Slices

Keine. [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
ist mit dieser Welle für die volle Python-Vier-Wege-Matrix vollständig
umgesetzt. Was ansteht, sind keine weiteren Slices dieser Welle, sondern:

- Die Vorschau-Zeile `welle-sdk-reale2e` in der Roadmap — ihr Trigger
  („`welle-sdk-python-vollabdeckung` liegt in `done/`") ist mit dieser
  Closure erfüllt; die Eröffnung ist ein künftiger Planner-Zug, kein
  Folge-Slice dieser Welle.
- Ein realer `sdk-python-v0.2.0`-Tag-Push — irreversible, extern
  sichtbare Betreiber-Handlung nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 (wie bei den vorigen
  SDK-Wellen).

## Verifikation

- `docs/reviews/review-slice-sdk-python-grpc-client-flaeche.md`
  (F-1…F-6) + Fixrunde (`b7a993fe`; Fixrunden-Report F-1…F-7, Rest
  in `5d5fa3eb` gelöst) +
  `docs/reviews/verifikation-slice-sdk-python-grpc-client-flaeche.md`
  (DoD erfüllt mit Auflagen V-1 bis V-3).
- `docs/reviews/review-slice-sdk-python-sse-client-flaeche.md`
  (F-1…F-9) + Fixrunde 1 (`3c941b0b`) + Fixrunde 2 (`beeddc3c`, FR-1) +
  `docs/reviews/review-slice-sdk-python-sse-client-flaeche-fixrunde.md` +
  `docs/reviews/verifikation-slice-sdk-python-sse-client-flaeche.md`
  (DoD erfüllt mit Auflagen V-1/V-2).
- `docs/reviews/review-slice-sdk-python-nats-stream-client-flaeche.md`
  (F-1…F-11) + Fixrunde 1 (`b4d1d352`) +
  `docs/reviews/review-slice-sdk-python-nats-stream-client-flaeche-fixrunde.md`
  (R-1…R-4) + Fixrunde 2 (`a3e9d64e`) +
  `docs/reviews/verifikation-slice-sdk-python-nats-stream-client-flaeche.md`
  (DoD-Substanz erfüllt, Abweichungen V-1/V-2 — beide vor bzw. bei dieser
  Closure gelöst: V-1 durch den frischen Neubau, V-2 als Prozess-Notiz).
- `make gates`: grün nach jedem Commit dieser Closure (ungepiped,
  Exit-Code direkt geprüft, [`AGENTS.md`](../../../../AGENTS.md) §3.9) —
  921 Dateien, 0 Befunde; Coverage 82,70 % ≥ 80 %;
  `a-check`/`generated-sync`/`commit-traceability`/`baseline-verify` je
  ohne Befund.
- Realer, grüner `make test-sdk-python-integration`-Lauf zu dieser
  Closure (Exit 0): alle drei Flächen gegen den laufenden Feed-Container —
  gRPC `change_id=804-1`, SSE `808-1`, NATS-Vollinhalt `810-1`, je
  SQL-Gegenprüfung gegen `cdc.changes`; Ablehnungs-Belege je Phase
  (gRPC `Unauthenticated`, HTTP `401`, NATS-Verbindungsablehnung).
- Frisches Artefakt-Paar `0.2.0` (20503/24655 Bytes, Neubau nach der
  zweiten Fixrunde; Wheel-Inhalt byte-gleich gegen den HEAD-Quellstand).
- `git tag -l "sdk-python-v*"`: ausschließlich `sdk-python-v0.1.0` —
  bestätigt, dass kein realer `0.2.0`-Veröffentlichungsversuch stattfand
  (Welle-Datei §3 vorab festgelegt).
- `docs/plan/carveouts/`: nur `.gitkeep` — 0 offene Carveouts dieser Welle.

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente
(kein `archiv`-Ziel in `Makefile`/`harness/mk/*.mk`, real geprüft per
`grep` — dieselbe Feststellung wie bei den vorigen SDK-Wellen) — die
Bedingung für Schritt 4 der Closure-Prozedur ist nicht eingetreten, keine
Handarbeit als Ersatz.