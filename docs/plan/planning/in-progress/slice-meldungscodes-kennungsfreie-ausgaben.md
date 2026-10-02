# Slice meldungscodes-kennungsfreie-ausgaben: Programm-Ausgaben ohne interne Kennungen, mit Gate (Teil 1 von 4, ohne Codes)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md) (Betriebsfähigkeit,
Haupt-Bezug), [`LH-QA-REL-003`](../../../../spec/lastenheft.md) (Fehlerklassen),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (Fehlerzustand erkennbar),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Festlegung 9 T1,
Festlegung 10 Wächter; Entscheidungsgrundlage),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) und
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(Vorbilder des Wächters).
Verdikt: [`architect-verdict-meldungscodes-statt-interner-kennungen`](../../../reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md).

**Berührte Spec-Stellen:** keine (T1 führt keine Codes ein; der Spec-Nachzug von
[`SPEC-008`](../../../../spec/pflichtenheft.md) gehört zu T2,
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 8).

**Teil-Slices der Umsetzung** ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 9): **T1 (dieser Slice)** → [T2 `meldungscodes-registry-fehlerkopf`](../open/slice-meldungscodes-registry-fehlerkopf.md)
→ [T3 `meldungscodes-warnungen-heartbeat-diagnose`](../open/slice-meldungscodes-warnungen-heartbeat-diagnose.md)
und [T4 `meldungscodes-http-grpc-fehlerkoerper`](../open/slice-meldungscodes-http-grpc-fehlerkoerper.md)
(T3 und T4 untereinander unabhängig).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `ba60c7bc`, 2026-10-02; Befehle im Feld §3).**
Der Auftraggeber hat festgelegt: Interne Kennungen (`LH-*`, `ADR-*`, `SPEC-*`,
`ARC-*`) gehören nicht in die Programm-Ausgabe; für Fehler und Warnungen braucht
der Betreiber eigene Nummern. Die Entscheidung dazu ist
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md); sie teilt die
Umsetzung in vier Slices. **T1 nimmt die Kennungen aus der Ausgabe, ohne Codes
einzuführen** — der sichtbare Teil der Vorgabe ist damit erfüllt, bevor die
Code-Tabelle existiert, und das Gate hält den Ausgangszustand 0, bevor T2 die
Fehlertext-Literale anfasst.

Nachgemessen am Code:

- **Produktions-Go, Ausgabe-Literale:** 14 Zeilen in drei Dateien tragen eine Kennung
  vor jedem `//`. Zwölf davon sind die Diagnose-Ausgabe in
  `internal/bootstrap/wiring.go` (Beispiel `Betriebsstatus (LH-FA-ADM-002): …`,
  `Fehlerzustand (LH-FA-ADM-003): …`, `Blockierender Consumer (LH-FA-RET-005): …`);
  eine ist ein Fehlertext der Konfigurationsdatei
  (`internal/bootstrap/config_file.go`), eine die Begründung der Rollout-Wache
  (`tools/schema/rolloutguard/guard.go`).
- **Skripte, die der Betreiber ausführt:** `tools/schema/rollout.sh` druckt in drei
  `echo`-Zeilen Kennungen (`Vorlauf (ADR-0114) …`, `bekannte Fremdobjekt-Blocker
  (ADR-0043) …`, `… Erstlieferung des d-migrate-Einbaus (ADR-0043) …`).
- **Datenbank-Metadaten:** `tools/schema/schema.yaml` trägt in 16 `description:`-Zeilen
  Kennungen. Ob sie als `COMMENT ON` beim Betreiber ankommen, ist **nicht gemessen**:
  `git grep -n 'COMMENT ON' -- tools/schema` liefert 0 Zeilen (Architect-Messung),
  ob d-migrate sie erzeugt, ist offen. Die Entscheidungsregel steht in
  [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 9 T1:
  Sichtbar beim Betreiber → kennungsfrei; sonst Quelltext-Dokumentation und bleibt.
- **Handbuch:** ein Hinweissatz („zusätzlich eine interne Anforderungskennung in
  Klammern“) und der Platzhalter `<interne Kennung>` im Rollout-Beispiel (je 1 Treffer).
- **Nicht betroffen (gemessen):** `sdks/` 0 Zeilen; die SQL-Funktionen `cdc.*` 0
  `RAISE EXCEPTION`.
- **Läufer und Tests, die den Text erwarten:** `internal/bootstrap/diagnose_test.go`,
  `tools/harness/run-integration-tests.sh`, `tools/harness/run-schema-rollout-guard-test.sh`,
  `harness/targets/schema-rollout.md`, `tools/schema/rolloutguard/guard_test.go`
  (Zahlen in §3). Die Zählung der übrigen Test-Zeilen mit Kennung ist **nicht**
  eindeutig: der Plan des Vorläufers nannte 142, der Architect summierte mit dem Muster
  `^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]` über `*_test.go` 97; der Implementer klärt die
  Differenz mit dem Muster des Plans und trennt Ausgabe-Erwartung von Testname und
  Kommentar (Suchlauf-Zeile 12; Zahl am Parent `ba60c7bc` in §3).

**Ziel:** Kein Ausgabe-Literal des Servers und seiner Betriebs-Skripte (Go-Produktion,
`echo`/`printf` in `tools/schema/` und `examples/`) trägt eine interne Kennung; das Gate
`ausgabe-kennungen-check` in `make gates` hält den Zustand ab Ist-Stand 0.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Meldungscodes, Tabelle, Katalog** — T2 bis T4
  ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 9);
  T1 braucht keine Code-Entscheidung.
- **Die 28 `Fehlerklasse <klasse>: `-Literale** — sie tragen keine Kennung und gehören
  zu T2 (dort wird der Kopf `Fehlerklasse <klasse> [<code>]: …`).
- **Release oder Versionsänderung der Pakete / Server-Release** — jeweils Freigabe des
  Auftraggebers; der Slice setzt kein Tag, die SDKs bleiben unberührt.
- **Änderung der Fehlerklassen-Semantik** — [`ADR-0023`](../../adr/0023-fehlerklassifikation.md)/
  [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) bleiben.
- **Quellcode-Kommentare** — [`AGENTS.md`](../../../../AGENTS.md) §3.7 regelt sie; das
  Gate liest Ausgabe-Literale, keine Kommentare.
- **Test-Läufer-Ausgaben** — `tools/harness/run-*.sh` drucken Kennungen als Fortschritt
  für Entwickler; nur ihre **Erwartungen an Produkt-Ausgaben** ziehen mit
  ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 10:
  Läufer sind ausgenommen).
- **Spec-Dokumente und ADRs** — sie führen Kennungen als Adresse
  ([`AGENTS.md`](../../../../AGENTS.md) §3.4, §3.5).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(A) Ausgaben kennungsfrei.** Die 14 Go-Zeilen und die 3 `echo`-Zeilen tragen
      keine Kennung mehr; die Erwartungen in `diagnose_test.go`, `run-integration-tests.sh`,
      `run-schema-rollout-guard-test.sh`, `guard_test.go` und `harness/targets/schema-rollout.md`
      sind nachgezogen; im Benutzerhandbuch entfallen der Hinweissatz und der Platzhalter
      `<interne Kennung>` (Beispiele = echte Ausgabe); die `description:`-Felder von
      `tools/schema/schema.yaml` sind nach der Regel der ADR entschieden (Messung im Bericht).
      *Zu belegen durch:* Suchlauf §3 (Diff-Zeilen 0 bzw. begründet); `make test` grün;
      `make test-integration` und `make test-store` grün nach `make image` (Läufer-Erwartungen,
      Rollout-Wache `tools/harness/run-schema-rollout-guard-test.sh`); die gedruckten Zeilen
      eines realen `diagnose`-Laufs im Closure-Bericht (**gemessen**, nicht übernommen) gegen
      das Handbuch-Beispiel; die Messung, ob die `description:`-Texte als Datenbank-Kommentar
      ankommen (z. B. `\d+` an einer gerollten Wegwerf-Datenbank, Befehl und Ausgabe im Bericht).
- [ ] **(B) Gate `ausgabe-kennungen-check`.** Ein netzloses Gate (`bash`/`git`/`grep`, ERE,
      kein `-P`) mit dem zweistufigen Muster und der Reichweite aus
      [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 10 steht
      in `GATE_CHECKS` und endet am Ist-Stand mit Exit 0; Tabellentest `make test-ausgabe-kennungen-check`
      (Werkzeug ohne Gate-Aufnahme, Muster [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md)
      Teilfrage 2); Sensor-Vertrag unter `harness/sensors/`; Zeilen in `harness/README.md` §Sensors
      (Gate und Test-Ziel; auch die Gates-Zeile `make gates` nennt es). *Zu belegen durch:*
      Tabellentest mit Treffer (Literal, Literal mit `/` davor, Literal mit URL, Raw-String in
      Backticks), Nicht-Treffer (Kommentarzeile, nachgestellter Kommentar ohne Anführungszeichen,
      Struct-Tag) und dem benannten Falsch-Positiv (nachgestellter Kommentar mit zitiertem Literal);
      eine **Mutationsprobe** auf einer Kopie im Scratchpad (Kennung in einem `fmt.Println`-Literal
      von `wiring.go` → rot, dieselbe Kennung im Kommentar → grün; Stelle, Instanz und gesehene
      Farbe im Bericht — in der ADR nur *hergeleitet*); `make gates` Exit 0 mit dem neuen Gate.
      Die Grenzen des Sensors (mehrzeiliges Raw-String-Literal wird nicht gelesen; nachgestelltes
      Kommentar-Zitat trifft) stehen im Vertrag ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)).
- [ ] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-kennungsfreie-ausgaben.md`
      endet mit Exit 0.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

Die Gate-Lockerung ist nicht berührt: das Gate ist neu, kein bestehendes wird
gesenkt; [`AGENTS.md`](../../../../AGENTS.md) §3.6 ist durch
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 10 gedeckt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten, nicht die Antwort. Der Implementer erweitert
die Liste in seinem ersten Lauf.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` (Diagnose-Ausgabe) | update | zwölf Ausgabe-Zeilen ohne Kennung. |
| `internal/bootstrap/config_file.go` | update | Fehlertext der Konfigurationsdatei ohne `ADR-0088`/`SPEC-016`. |
| `tools/schema/rollout.sh`, `tools/schema/rolloutguard/guard.go` | update | Rollout-Meldungen ohne Kennung. |
| `tools/schema/schema.yaml` (16 `description:`-Zeilen) | messen / update | erst messen, ob sie als `COMMENT ON` ankommen; dann kennungsfrei oder als Bestand begründet (Regel der ADR). |
| `internal/bootstrap/diagnose_test.go`, `tools/schema/rolloutguard/guard_test.go` | update | Erwartungen an die Ausgabe; Happy/Boundary/Negative bleiben. |
| `tools/harness/run-integration-tests.sh`, `tools/harness/run-schema-rollout-guard-test.sh`, `harness/targets/schema-rollout.md` | update | Erwartungen an Produkt-Ausgaben (nur diese Zeilen, nicht die Läufer-Ausgabe). |
| `docs/user/benutzerhandbuch.md` | update | Hinweissatz und Platzhalter entfallen, Beispiele = echte Ausgabe; Änderungshistorie in Betreibersicht ohne Kennung. Parallel arbeiten Implementer an der Datei: Rebase auf `main`, nur eigene Abschnitte committen. |
| Wächter-Skript unter `tools/harness/`, Tabellentest-Läufer, `harness/sensors/ausgabe-kennungen-check.md`, `harness/README.md` §Sensors, `harness/mk/doc-gate.mk`, `GATE_CHECKS` | neu / update | (B); Heimat des Targets analog [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 5 (`meldungscodes-check` liegt in `doc-gate.mk`), der Implementer bestätigt oder begründet eine andere. |

**Wächter-Form (entschieden, nicht mehr Option).** Zweistufige ERE nach
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 10: eine
Zeile trifft, wenn eine Kennung `(LH|ADR|SPEC|ARC)-[A-Z0-9]` hinter einem ungeschlossenen
Anführungszeichen steht und die Zeile nicht mit `//` beginnt. Gegenstand: Produktions-Go
unter `internal/`, `cmd/`, `tools/schema/` (ohne `*_test.go`, ohne erzeugten Code) und
`echo`/`printf`-Zeilen der Skripte unter `tools/schema/` und `examples/`; ausgenommen die
Läufer unter `tools/harness/`. Die Suchlauf-Zeilen unten bleiben `-P` (ein Suchlauf ist kein
Gate), das Gate selbst nutzt kein `-P`. Die Frage nach `grep`-Falsch-Positiven am Bestand ist
von der ADR gemessen (0 bei 14 Treffern); der Implementer misst sie am Diff erneut.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „Programm-Ausgaben tragen interne
Kennungen“ und das Handbuch beschreibt das; Parent ist `ba60c7bc`; die `diff`-Zeilen und die
Befunde trägt der Implementer ein; neue Dateien sind für den Stand `diff` mit `git add` im
Index):**

```suchlauf
ba60c7bc 14 -n -P '^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]' -- '*.go' ':!*_test.go'
ba60c7bc 3 -n -E 'echo.*(ADR|LH|SPEC)-[A-Z0-9]' -- tools/schema/rollout.sh
ba60c7bc 13 -n -E '\((LH|ADR)-' -- internal/bootstrap/wiring.go
ba60c7bc 7 -n -E '\((LH|ADR)-' -- internal/bootstrap/diagnose_test.go
ba60c7bc 9 -n -E '(Betriebsstatus|Fehlerzustand|Blockierender Consumer|CDC-Abstand cdc_capture_lag|Speicherverbrauch cdc_storage_bytes) \(LH' -- tools/harness/run-integration-tests.sh
ba60c7bc 7 -n -E 'Vorlauf \(ADR|Fremdobjekt-Blocker \(ADR' -- tools harness
ba60c7bc 1 -n -F 'interne Anforderungskennung' -- docs/user
ba60c7bc 1 -n -F '<interne Kennung>' -- docs/user
ba60c7bc 0 -n -E '(SPEC|ADR|ARC)-[0-9]+|LH-(FA|QA)-' -- sdks
ba60c7bc 16 -n -E 'description:.*(LH|ADR|SPEC|ARC)-' -- tools/schema/schema.yaml
ba60c7bc 28 -n -E '"Fehlerklasse [a-z]+: ' -- '*.go' ':!*_test.go'
ba60c7bc 97 -n -P '^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]' -- '*_test.go'
diff 0 -n -P '^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]' -- '*.go' ':!*_test.go'
diff 0 -n -E 'echo.*(ADR|LH|SPEC)-[A-Z0-9]' -- tools/schema/rollout.sh
diff 1 -n -E '\((LH|ADR)-' -- internal/bootstrap/wiring.go
diff 0 -n -E '\((LH|ADR)-' -- internal/bootstrap/diagnose_test.go
diff 0 -n -E '(Betriebsstatus|Fehlerzustand|Blockierender Consumer|CDC-Abstand cdc_capture_lag|Speicherverbrauch cdc_storage_bytes) \(LH' -- tools/harness/run-integration-tests.sh
ba60c7bc 3 -n -E '\\\((LH|ADR)-' -- tools/harness/run-integration-tests.sh
diff 0 -n -E '\\\((LH|ADR)-' -- tools/harness/run-integration-tests.sh
diff 3 -n -E 'Vorlauf \(ADR|Fremdobjekt-Blocker \(ADR' -- tools harness
diff 0 -n -F 'interne Anforderungskennung' -- docs/user
diff 0 -n -F '<interne Kennung>' -- docs/user
diff 0 -n -E '(SPEC|ADR|ARC)-[0-9]+|LH-(FA|QA)-' -- sdks
diff 16 -n -E 'description:.*(LH|ADR|SPEC|ARC)-' -- tools/schema/schema.yaml
diff 28 -n -E '"Fehlerklasse [a-z]+: ' -- '*.go' ':!*_test.go'
diff 90 -n -P '^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]' -- '*_test.go'
```

Die Sollzahlen sind am Parent `ba60c7bc` mit `make suchlauf-nachmessen` gemessen (12 Zeilen,
Exit 0). Zeile 5 liefert 9 statt der 12 des Vorläufer-Plans (Stand `f8aca691`); das Muster
wurde dort mit einem zusätzlichen Escape (`\\*`) geführt, die Ursache der Differenz ist nicht
geklärt — der Implementer klärt sie, bevor er sich auf die Zahl der Läufer-Erwartungen stützt.

| Träger | Messung am Parent (`ba60c7bc`, gemessen am 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| Ausgabe-Literale im Go-Produktionscode | Zeile 1: 14 (12 `wiring.go`, 1 `config_file.go`, 1 `guard.go`) | Gefunden und nachgezogen: 14 Zeilen in `wiring.go`, `config_file.go`, `guard.go`; Diff-Stand 0 (erste `diff`-Zeile des Blocks). Nichtgefunden: kein weiteres Ausgabe-Literal mit Kennung im Produktions-Go |
| Rollout-Meldungen | Zeile 2: 3; Zeile 6: Meldungstexte samt Erwartungen und Vertragsdatei `harness/targets/schema-rollout.md` | Gefunden und nachgezogen: 3 `echo`-Zeilen in `rollout.sh`, 2 Erwartungen in `run-schema-rollout-guard-test.sh`, 1 Meldungsnennung in `harness/targets/schema-rollout.md`; Diff-Stand Zeile 6 = 3, alle drei bewusst stehen geblieben: ein Läufer-Kommentar (`run-schema-rollout-guard-test.sh:12`), ein Test-Kommentar (`guard_test.go:123`), eine Fixture des neuen Tabellentests; Zeile 2 Diff-Stand 0 |
| Diagnose-Ausgabe und ihre Erwartungen | Zeile 3 (13: 12 Literale + 1 Kommentar), Zeile 4 (7), Zeile 5 (Erwartungen in `run-integration-tests.sh`) | Gefunden und nachgezogen: 12 Literale; 7 Erwartungen in `diagnose_test.go` (Diff 0), 9 Erwartungszeilen in `run-integration-tests.sh` (Diff 0); Zeile 3 Diff-Stand 1 = der Kommentar mit `ADR-0132` (Quellcode-Kommentar, außerhalb der Reichweite). Differenz 9 gegen 12 des Vorläufer-Plans geklärt: 12 Erwartungszeilen an die Produkt-Ausgabe, davon 9 mit der Klammer `(LH-…)` im Klartext (Suchlauf-Zeile 5, 3 in der Prüfliste, 6 einzelne `grep -qF`) und 3 als `grep -qE` mit maskierter Klammer `\(LH-…\)` (eigene Suchlauf-Zeile mit dem maskierten Muster: `Speicherverbrauch` zweimal, `Blockierender Consumer` einmal), die das Muster der Zeile 5 nicht trifft; alle 12 nachgezogen. Der erste `make test-integration`-Lauf nach den neun Zeilen endete rot an der Speicherverbrauch-Erwartung, genau an diesen drei; die weiteren Kennungen in der Datei (`echo`-Läufer-Ausgaben, Kommentare) sind Läufer-Ausgabe und bleiben |
| Handbuch | Zeilen 7 und 8 (je 1) | Gefunden und nachgezogen: Hinweissatz und Platzhalter entfernt (Diff 0 beider); Version 1.89 mit Historienzeile; zusätzlich die Beispielzeile „Blockierender Consumer“ an die echte Ausgabe (`<Name> (<Kennung>)`) angeglichen |
| SDKs | Zeile 9: 0 | Nichtgefunden belegt (Diff-Stand 0); kein Nachzug |
| Datenbank-Metadaten | Zeile 10: 16 | Gemessen (Wegwerf-PostgreSQL 18, `bash tools/schema/apply-rollout.sh`, danach `SELECT count(*) FROM pg_description d JOIN pg_class c ON c.oid = d.objoid WHERE c.relnamespace = 'cdc'::regnamespace` → 0; mit Kennungsfilter → 0; weder `\d+ cdc.source` noch `plan.yaml`, `down.sql`, `rollout-precheck.yaml` oder das Rollout-Log tragen eine `description:`): die Texte kommen nicht als Datenbank-Kommentar an, kein `COMMENT ON` im Rollout. Entscheidung nach der Regel der ADR: Quelltext-Dokumentation des Schema-YAML, bleibt (Diff-Stand 16) |
| `Fehlerklasse`-Präfix | Zeile 11: 28 | Diff-Stand 28, unverändert (gehört zu T2) |
| Test-Zeilen mit Kennung | Zeile 12: 97 (Muster des Architects; Vorläufer-Plan nannte 142, nicht nachgemessen) | Diff-Stand 90 (−7: die Erwartungen in `diagnose_test.go`). Die 97 sind 42 Fixtures in `tools/harness/kommentar-kennungen/main_test.go`, 18 in `roles_rollout_file_internal_test.go` und sonst `t.Fatalf`-Fehlertexte und Testnamen; Ausgabe-Erwartung an die Produkt-Ausgabe waren nur die 7 in `diagnose_test.go`. Die 142 sind nicht reproduziert: das Muster `"[^"]*(LH|ADR|SPEC|ARC)-[A-Z0-9]` liefert 147 Zeilen am Parent (mit Kommentaren, die ein Anführungszeichen tragen), das Muster ohne Schrägstrich-Ausschluss 1029 |

**Belege des Implementers (gemessen am Arbeitsbaum nach `ba60c7bc`, 2026-10-03):**

- **Realer `diagnose`-Lauf** (`ghcr.io/pt9912/pg-change-feed:dev` nach `make image`, gegen eine Wegwerf-PostgreSQL 18 nach `apply-rollout.sh`, Quelle `demo`): ohne Lebenszeichen `Betriebsstatus: kein Lebenszeichen — Instanz hat noch nie geschlagen` / `Fehlerzustand: unbekannt (kein Lebenszeichen)`; mit Lebenszeichen und Consumer `Betriebsstatus: Lebenszeichen vor 0.415s`, `Fehlerzustand: keiner (Normalbetrieb)`, `Blockierender Consumer: cli-consumer (c1), bestätigte Position 5, Rückstand unbekannt (…)`; mit `error_class = 'schema'` `Fehlerzustand: schema`. Keine Zeile trägt eine Kennung. Der Fehlertext der Konfigurationsdatei lautet `… trägt den Schlüssel "capture_dsn" — Zugangsdaten bleiben env-var-exklusiv`.
- **Mutationsprobe des Gates** (Kopie von `internal/`, `cmd/`, `examples/`, `tools/schema/` im Scratchpad, Wurzel als erstes Argument): Ausgangslage Exit 0; `fmt.Println("  Fehlerzustand (LH-FA-ADM-003): keiner (Normalbetrieb)")` in `wiring.go` → **Exit 1** mit `internal/bootstrap/wiring.go:2237:…`; dieselbe Kennung als Kommentarzeile davor → **Exit 0**; `echo "schema-rollout: Vorlauf (ADR-0114) - …"` in `rollout.sh` → **Exit 1**. Mutationen am Wächter selbst (Kopie in einem Scratchpad-Repo, Tabellentest dagegen): `skip_go` wirkungslos → Fall „eingerückte Kommentarzeile“ rot; Backtick aus dem Go-Muster → Fall „Raw-String in Backticks“ rot; Test-Ausnahme aus dem `find` → Fall „Kennung in Test-Datei“ rot; Lesefehler-Zweig von `grep` abgeschaltet → die Lesefehler-Fälle rot.

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt. Die Architect-Entscheidung liegt vor
([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md), `Accepted`); T1
braucht keine weitere Entscheidung und keine Bestätigung des Code-Präfixes.

**Rückführung:** `in-progress` → `open` (blockiert), wenn die Messung zeigt, dass die
`description:`-Felder als Datenbank-Kommentar ankommen **und** eine Änderung das Schema der
Rollout-Wache bricht (Rollout-Identität); dann ist das eine Auftraggeber-Frage, keine
stille Anpassung einer Erwartung. Das Gate geht nie rot in `main`: ist der Ist-Stand nach (A)
nicht 0, wird das Gate nicht aufgenommen, bis er es ist.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration`, `make test-store` grün nach `make image`, die
Handbuch-Beispiel-Ausgabe mit einem realen `diagnose`-Lauf verglichen, Wächter-Mutationsprobe
rot gesehen, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.
Die Fortschreibung des Gate-Index (`harness/README.md` §Sensors) ist Teil des Diffs.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **Ausgaben-Stabilität für Betreiber, die die Ausgabe parsen.** Die Diagnose-Zeilen und
  Rollout-Meldungen sind Text; wer sie mit `grep` auswertet, bricht an jeder Änderung.
  Entschieden in [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
  Konsequenzen: keine Übergangsregel, der Text ist nicht Vertrag. Für T1 betrifft es nur die
  entfallende Klammer (`(LH-FA-ADM-002)`), keinen Anker, den das Handbuch zusagt; die
  Handbuch-Abschnitte Diagnose und Rollout liest der Implementer dennoch gegen eine Textzusage
  (ADR: im Handbuch nicht nachgemessen) — findet er eine, geht der Slice zurück an den
  Auftraggeber. — **Ausgang:** (bei Closure)
- **Tests erwarten Texte.** `diagnose_test.go`, `run-integration-tests.sh`,
  `run-schema-rollout-guard-test.sh`, `guard_test.go`: ein vergessener Läufer wird erst beim
  `make test-integration`-Lauf rot, der nicht in `make gates` liegt. Gegenmittel:
  Suchlauf-Zeilen je Träger, vollständiger Lauf vor Closure. — **Ausgang:** (bei Closure)
- **Datenbank-Metadaten (`COMMENT ON`) tragen Kennungen.** Offen bis T1: ungemessen, ob der
  Betreiber sie sieht; eine Änderung der `description:`-Felder ändert das Rollout-Ergebnis und
  die Schema-Identität (Rollout-Wache, `make test-store`). Regel der ADR: sichtbar → kennungsfrei,
  sonst bleibt. — **Ausgang:** (bei Closure)
- **Wächter mit Falsch-Positiven oder -Negativen.** Eine zeilenbasierte Regel liest keine
  mehrzeiligen Literale und trifft ein nachgestelltes Kommentar-Zitat; beides ist als Grenze
  benannt und im Sensor-Vertrag zu nennen. Ein zu eng gelesenes Muster erzeugt falsche
  Sicherheit ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)); die
  Mutationsprobe der ADR ist *hergeleitet*, der Implementer fährt sie. — **Ausgang:** (bei Closure)
- **Kollision mit parallelen Arbeiten am Handbuch.** Gegenmittel: Rebase auf `main`, nur
  eigene Abschnitte committen. — **Ausgang:** (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** T2 `meldungscodes-registry-fehlerkopf`.
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird die Default-Sub-Area `*` (Kürzel
`PGC`, Modus Greenfield, [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); sie umfasst Produktions-Code, Skripte und Doku gleichermaßen, eine
feinere Aufteilung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register `../observations/BEO-PGC/` ist
nicht Inhalt dieses Plans; der Implementer sichtet es vor dem Start auf Treffer zu
„Kennung in Ausgabe“ und „Text-Erwartung im Läufer“ (nicht gemessen, kein Treffer behauptet).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
