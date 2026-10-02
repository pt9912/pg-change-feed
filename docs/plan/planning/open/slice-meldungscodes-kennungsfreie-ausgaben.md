# Slice meldungscodes-statt-interner-kennungen: Nutzerseitige Meldungscodes statt interner Kennungen in Programm-Ausgaben

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
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) und
[`ADR-0023`](../../adr/0023-fehlerklassifikation.md) (Fehlerklassen, die der Slice
nicht ändert), [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
(Herkunft von Aussagen), [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md)
(Vorbild des Wächters: keine interne Kennung unter `sdks/`).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklassen-Tabelle; gelesen, die Klassen bleiben) — ob eine Meldungscode-Festlegung
im Pflichtenheft eine eigene Stelle braucht, entscheidet die ADR des Architects (§4).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `f8aca691`, 2026-10-02; Befehle im Feld §3).**
Der Auftraggeber hat festgelegt: Interne Kennungen (`LH-*`, `ADR-*`, `SPEC-*`,
`ARC-*`) gehören auch nicht in die Programm-Ausgabe; für Fehler und Warnungen
braucht der Betreiber eigene Nummern, Vorbild ist das Schwester-Repo d-migrate.
Die Handbuch-Bereinigung (`7d1c4611`) hat die Kennungen aus dem Benutzerhandbuch
entfernt und die Lücke offen gelassen: das Handbuch erklärt in einem Hinweissatz, die
Ausgabe von `diagnose` trage „zusätzlich eine interne Anforderungskennung in
Klammern“. Nachgemessen am Code:

- **Produktions-Code, Ausgabe-Strings:** 14 Zeilen in drei Dateien tragen eine
  Kennung vor jedem `//`. Zwölf davon sind die Diagnose-Ausgabe in
  `internal/bootstrap/wiring.go` (Zeilen 2232–2270, Beispiel
  `Betriebsstatus (LH-FA-ADM-002): …`, `Fehlerzustand (LH-FA-ADM-003): …`,
  `CDC-Abstand cdc_capture_lag (LH-FA-ADM-004): …`,
  `Blockierender Consumer (LH-FA-RET-005): …`,
  `Speicherverbrauch cdc_storage_bytes (LH-FA-RET-006): …`,
  `Backfill je Tabelle (LH-FA-CAP-009, …)`); eine ist ein Fehlertext der
  Konfigurationsdatei (`internal/bootstrap/config_file.go:114`, trägt
  `ADR-0088 Festlegung 1, SPEC-016` im Text); eine ist die Begründung der
  Rollout-Wache (`tools/schema/rolloutguard/guard.go:148`, `ADR-0043`/`ADR-0114`).
- **Skripte, die der Betreiber ausführt:** `tools/schema/rollout.sh` druckt in drei
  `echo`-Zeilen Kennungen (`Vorlauf (ADR-0114) - View-Signatur-Aenderung …`,
  `bekannte Fremdobjekt-Blocker (ADR-0043) …`, `… Erstlieferung des
  d-migrate-Einbaus (ADR-0043) …`).
- **Datenbank-Metadaten:** `tools/schema/schema.yaml` trägt in 16 `description:`-Zeilen
  Kennungen (`SPEC-001`, `LH-FA-DAT-002`, …); ob sie als `COMMENT ON` in der
  Datenbank ankommen und der Betreiber sie mit `\d+` liest, ist **nicht gemessen**
  (Frage an den Implementer; §3).
- **Nicht betroffen (gemessen):** `sdks/` trägt 0 Zeilen (`make sdk-public-doc-check`
  deckt das); die SQL-Funktionen `cdc.*` tragen in `RAISE EXCEPTION`-Texten 0
  Kennungen; HTTP-/gRPC-Fehlertexte tragen außer dem Konfigurationstext oben keine.
  Die Quelle der Fehlertexte trägt bereits ein einheitliches Präfix `Fehlerklasse
  <klasse>: …` (28 Zeilen in Produktions-Go, z. B.
  `internal/adapters/driving/replication/mapper/mapper.go:30`) — die Klasse ist heute
  die einzige stabile Maschinen-Kennung einer Fehlermeldung, eine Nummer je Meldung
  fehlt.
- **Tests und Test-Läufer, die den Text erwarten:** `internal/bootstrap/diagnose_test.go`
  (7 Zeilen), `tools/harness/run-integration-tests.sh` (12 Zeilen, `grep`-Erwartungen
  an die Diagnose-Ausgabe, Suchlauf-Zeile 5 in §3),
  `tools/harness/run-schema-rollout-guard-test.sh` (Erwartung an `Vorlauf (ADR-0114)`).
  Die übrigen `*_test.go`-Treffer (142 Zeilen mit Kennung im Literal, gemessen mit
  `git grep` über `*_test.go`) sind **nicht** nach Ausgabe-Erwartung und Testname getrennt (ungemessen, nur die
  drei genannten Dateien sind als Erwartung gelesen); der Implementer trennt beide (§3).

**Muster im Schwester-Repo d-migrate (gelesen, Dateien relativ zu dessen Wurzel):**

- *Aufbau.* Ein Code ist `E` oder `W` plus drei Ziffern (`E053`, `W116`); der Buchstabe
  trägt die Schwere (Fehler/Warnung), die Nummer einen Bereich
  (`spec/cli-spec.md` §4.1: `E001–E099` Validierung, `E100–E199` Verbindung,
  `E200–E299` Migration, `E300–E399` KI, `E400–E499` Konfiguration; `W001–W099`
  Validierungswarnungen, `W100–W199` Kompatibilität, `W200–W299` Performance). Die
  Praxis ist über die Bereiche hinausgewachsen und führt sie in
  `spec/ledger.md` §Code-Nummernbereiche als Tabelle (z. B. `E052–E056`
  Dialekt-Inkompatibilitäten, `W114–W117` Sequence-Emulation); neue Codes werden am
  Ende des Bereichs angefügt, Lücken sind reserviert.
- *Träger im Code.* `ValidationError(code, message, objectPath)` und
  `ValidationWarning(code, message, objectPath)`
  (`hexagon/core/src/main/kotlin/dev/dmigrate/core/validation/ValidationResult.kt`):
  der Code ist ein eigenes Feld, nicht Teil des Textes. Die CLI formatiert ihn aus
  Message-Bundles
  (`adapters/driving/cli/src/main/resources/messages/messages.properties`, `…_de.properties`):
  `Error [{0}]: {1}`, `Warning [{0}]: {1}`, `[ERROR] {0}`.
- *Dokumentation und Stabilität.* Jeder nutzersichtbare Code ist in einem **Ledger**
  registriert (`ledger/error-code-ledger-<version>.yaml`, Schema
  `ledger/code-ledger-<version>.schema.json`; Konventionen `spec/ledger.md`): je Code
  `level` (error/warning), `status` (`active` | `not_applicable` | `reserved`),
  `test_path` und `evidence_paths` (Produktionsquelle). Ein Code wird nie neu
  belegt: `E137` ist als „Zurückgezogen“ geführt, „die Kennung bleibt vergeben“, weil
  sie in ausgelieferten Berichten stand. Benutzerseitig nennt
  `docs/user/api-referenz.md` §2.4 das Namensschema und eine Auswahl, vollständig
  steht es in `spec/cli-spec.md` §4.
- *Test.* `hexagon/core/src/test/kotlin/dev/dmigrate/core/validation/CodeLedgerValidationTest.kt`
  liest das Ledger, prüft Status, Level und Pflichtfelder und dass jeder Code Test und
  Evidenz hat — „kein Code ohne Test und Evidence“ (`spec/ledger.md` §Zweck). Eine
  Offene-Punkte-Datei hält die Lücken (`docs/planning/open/warn-code-ledger-completeness.md`).
- *Grenze der Übertragbarkeit.* d-migrate ist ein Batch-Werkzeug mit Validierungsergebnissen;
  pg-change-feed ist ein Dauerprozess mit Log, Diagnose, HTTP, gRPC, SQL-Funktionen und
  einer eigenen Fehlerklassen-Achse (`transient|permission|configuration|schema|storage|replication|internal`).
  Ob Klasse und Code zwei Achsen oder eine sind, ist die Architect-Frage (A) — der
  Plan nimmt sie nicht vorweg.

**Ziel:** Programm-Ausgaben des Servers und seiner Betriebs-Skripte tragen keine interne
Kennung mehr; Fehler und Warnungen tragen stattdessen einen stabilen, im Handbuch
katalogisierten Meldungscode nach einer Architektur-Entscheidung, und ein Wächter hält
den Zustand.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Release oder Versionsänderung der Pakete** (C#, Python, Kotlin) — jedes Release ist
  eine Freigabe des Auftraggebers; die SDKs tragen keine Kennung (gemessen, 0 Zeilen) und
  reichen Servertexte durch (Risiko §6), ändern sich in diesem Slice aber nicht.
- **Server-Release** — separater Vorgang mit eigener Freigabe; der Slice setzt kein Tag.
- **Änderung der Fehlerklassen-Semantik** — Klassen, ihr Verhalten und Exit-Codes sind
  [`ADR-0023`](../../adr/0023-fehlerklassifikation.md)/[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md);
  ein Meldungscode ergänzt die Klasse, ersetzt sie nicht. Ändert die Architect-Entscheidung
  die Klassen, ist das eine eigene ADR mit `Supersedes`, kein Teil dieses Slice.
- **Quellcode-Kommentare** — [`AGENTS.md`](../../../../AGENTS.md) §3.7 regelt sie
  (ein Anker je Kommentar, `make kommentar-kennungen`); der Wächter aus (C) liest
  Ausgabe-Literale, keine Kommentare.
- **Test-Läufer-Ausgaben** — `tools/harness/run-*.sh` drucken Kennungen als Fortschritt
  und Abdeckungs-Deklaration für Entwickler (z. B. `abdeckung_declare` und
  Fortschrittszeilen in `run-integration-tests.sh`). Sie sind keine Ausgabe an
  den Betreiber; nur ihre **Erwartungen an Produkt-Ausgaben** (12 + 2 Zeilen, s. o.) ziehen
  mit. Ob die Läufer-Ausgabe selbst kennungsfrei werden soll, ist ein anderer Vorgang
  (Folge-Slice nur auf Auftraggeber-Wunsch).
- **Die Spec-Dokumente und ADRs** — sie führen Kennungen als ihre Adresse
  ([`AGENTS.md`](../../../../AGENTS.md) §3.4, §3.5); der Slice berührt sie nur, wenn die
  ADR eine Spec-Stelle festlegt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(A) Architect-Entscheidung liegt vor.** Eine ADR (`Accepted`, im
      [ADR-Index](../../adr/README.md) eingetragen) legt das Schema nutzerseitiger
      Meldungscodes fest (Teilfragen §4) und nennt für jede Teilfrage mindestens drei
      Optionen mit Pro/Contra. *Zu belegen durch:* die ADR-Datei, der Index-Eintrag,
      `make docs-check` Exit 0. Die Umsetzung (B) beginnt erst danach.
- [ ] **(B) Umsetzung.** Kein Programm-Ausgabe-Literal trägt eine interne Kennung;
      Fehler und Warnungen tragen den Meldungscode nach der ADR. Der Katalog steht im
      Benutzerhandbuch (kennungsfrei: Code, Bedeutung, Maßnahme); die Beispiele des
      Handbuchs sind die echte Ausgabe, der Hinweissatz „zusätzlich eine interne
      Anforderungskennung in Klammern“ und der Platzhalter `<interne Kennung>` im
      Rollout-Beispiel entfallen. *Zu belegen durch:* Suchlauf §3 (Zeilen der
      Produktion und des Handbuchs am Diff gleich 0 bzw. nur Meldungscodes); die
      angepassten Tests grün (`make test`); `make test-integration` und
      `make test-store` grün nach `make image` (Läufer-Erwartungen an die
      Diagnose-Ausgabe, Rollout-Wache `tools/harness/run-schema-rollout-guard-test.sh`);
      ein Vergleich der Handbuch-Beispiel-Ausgabe mit einem realen `diagnose`-Lauf
      (die gedruckten Zeilen im Closure-Bericht, **gemessen**, nicht übernommen).
- [ ] **(C) Wächter.** Ein Sensor lässt eine Kennung in einem Ausgabe-Literal
      scheitern. *Zu belegen durch:* Tabellentest mit Treffer-, Nicht-Treffer- und
      Randfällen (Kommentar mit Kennung, Testdatei, Dokumenten-Beispiel), eine
      **Mutationsprobe** (eine Kennung in `wiring.go` auf einer Kopie im Scratchpad
      eingefügt, Wächter rot gesehen, Rücknahme per `cp`), und der Eintrag in
      `harness/README.md` §Sensors samt Sensor-Vertrag. Ob er ein Gate in
      `make gates` wird, entscheidet die ADR bzw. der Architect nach der Messung der
      Falsch-Positiv-Quote (Option, §3 Wächter).
- [ ] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-statt-interner-kennungen.md`
      endet mit Exit 0.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten, nicht die Antwort. Der Implementer erweitert
die Liste in seinem ersten Lauf; Teil (B) wird erst nach der ADR (A) konkret.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/<nächste freie Nummer>-…md`, `docs/plan/adr/README.md` | neu / update | (A), Autor Architect; nicht vom Planner vorweggenommen. |
| `internal/bootstrap/wiring.go` (Diagnose-Ausgabe, Z. 2232–2270) | update | zwölf Ausgabe-Zeilen ohne Kennung, ggf. mit Meldungscode nach ADR. |
| `internal/bootstrap/config_file.go` (Z. 114) | update | Fehlertext der Konfigurationsdatei ohne `ADR-0088`/`SPEC-016`. |
| `tools/schema/rollout.sh`, `tools/schema/rolloutguard/guard.go` (Z. 148), ggf. `report.go` | update | Rollout-Meldungen ohne Kennung. |
| `tools/schema/schema.yaml` (16 `description:`-Zeilen) | prüfen / update | erst messen, ob sie als `COMMENT ON` beim Betreiber ankommen; dann kennungsfrei oder bewusst als Bestand begründet. |
| Fehlertexte mit Präfix `Fehlerklasse <klasse>:` (28 Zeilen in Produktions-Go) | prüfen / update | Träger des Meldungscodes, falls die ADR ihn dort verankert. |
| `internal/bootstrap/diagnose_test.go`, `internal/bootstrap/*_test.go` | update | Erwartungen an die Ausgabe; Happy/Boundary/Negative je neuem Code nach der ADR. |
| `tools/harness/run-integration-tests.sh`, `tools/harness/run-schema-rollout-guard-test.sh` | update | `grep`-Erwartungen an Produkt-Ausgaben (nur diese Zeilen, nicht die Läufer-Ausgabe). |
| `docs/user/benutzerhandbuch.md` | update | Katalog der Meldungscodes (kennungsfrei), Diagnose- und Rollout-Beispiele auf die echte Ausgabe, Hinweissatz entfällt. Parallel arbeiten Implementer an dieser Datei: der Implementer dieses Slice rebasiert auf den Stand von `main` und committet nur seine Abschnitte. |
| Wächter-Skript unter `tools/harness/`, `harness/sensors/<name>.md`, `harness/README.md` §Sensors, `Makefile`/`harness/mk/*.mk` | neu / update | (C), Form siehe unten. |

**Option für den Wächter (C) — Machbarkeit, keine Entscheidung.** Die Hauptzahl misst
dieser Plan: ein Muster `^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]` über `*.go` ohne `*_test.go`
trifft 14 Zeilen, **alle** echte Ausgabe-Literale (Suchlauf-Zeile 1) — der Ausdruck
`^[^/]*` schließt Kommentarzeilen und nachgestellte Kommentare aus, weil dort ein `/`
vor der Kennung steht. Grenzen, die der Implementer vermisst, bevor der Wächter ein Gate
wird:

- *Falsch-Negative:* ein Literal mit `/` vor der Kennung (URL, Pfad im Text) und ein
  mehrzeiliges Literal; ein Go-AST-Wächter (Programm wie
  `tools/harness/kommentar-kennungen/`) liest String-Literale statt Zeilen und hat diese
  Lücke nicht, kostet aber ein Docker-Programm.
- *Falsch-Positive:* `*_test.go`-Dateien (142 Zeilen, ausgenommen), Struct-Tags, Metrik-
  und Label-Namen, Dokumentations-Konstanten, die einer Kennung ähneln; Skripte brauchen
  eine Unterscheidung `echo`/`printf` gegen Kommentar (`rollout.sh`: 3 echo-Zeilen,
  Suchlauf-Zeile 2 in §3).
- *Reichweite:* Go-Produktion + Betriebs-Skripte (`tools/schema/`, `examples/`) + die
  `description:`-Felder des Schemas; Test-Läufer ausgenommen (§1). Das Vorbild ist
  `make sdk-public-doc-check` ([`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md)):
  reines `grep`, netzlos, eigene Ausnahmen namentlich im Skript.
- *Gate oder Werkzeug:* ein Gate nur, wenn der Ist-Stand nach (B) 0 Treffer hat und die
  Ausnahmen benannt sind (AGENTS §3.6: Gates werden nicht ohne ADR gelockert, ein neues
  Gate braucht eine Zeile in `harness/README.md` §Sensors).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „Programm-Ausgaben tragen
interne Kennungen“ und das Handbuch beschreibt das; Parent ist `f8aca691`; die `diff`-Zeilen
und die Befunde trägt der Implementer ein; neue Dateien sind für den Stand `diff` mit
`git add` im Index):**

```suchlauf
f8aca691 14 -n -P '^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]' -- '*.go' ':!*_test.go'
f8aca691 3 -n -E 'echo.*(ADR|LH|SPEC)-[A-Z0-9]' -- tools/schema/rollout.sh
f8aca691 13 -n -E '\((LH|ADR)-' -- internal/bootstrap/wiring.go
f8aca691 7 -n -E '\((LH|ADR)-' -- internal/bootstrap/diagnose_test.go
f8aca691 12 -n -E '(Betriebsstatus|Fehlerzustand|Blockierender Consumer|CDC-Abstand cdc_capture_lag|Speicherverbrauch cdc_storage_bytes) \\*\(LH' -- tools/harness/run-integration-tests.sh
f8aca691 7 -n -E 'Vorlauf \(ADR|Fremdobjekt-Blocker \(ADR' -- tools harness
f8aca691 1 -n -F 'interne Anforderungskennung' -- docs/user
f8aca691 1 -n -F '<interne Kennung>' -- docs/user
f8aca691 0 -n -E '(SPEC|ADR|ARC)-[0-9]+|LH-(FA|QA)-' -- sdks
f8aca691 16 -n -E 'description:.*(LH|ADR|SPEC|ARC)-' -- tools/schema/schema.yaml
f8aca691 28 -n -E '"Fehlerklasse [a-z]+: ' -- '*.go' ':!*_test.go'
```

| Träger | Messung am Parent (`f8aca691`, gemessen am 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| Ausgabe-Literale im Go-Produktionscode | Zeile 1: 14 Zeilen (12 `wiring.go`, 1 `config_file.go`, 1 `guard.go`) | Implementer: Diff-Stand soll 0 sein |
| Rollout-Meldungen | Zeile 2: 3 Zeilen; Zeile 6: 7 Zeilen (Meldungstexte samt Erwartungen und Vertragsdatei `harness/targets/schema-rollout.md`) | Implementer: Quelle, Erwartungen im Läufer und Vertragsdatei zusammen nachziehen |
| Diagnose-Ausgabe und ihre Erwartungen | Zeile 3 (13, davon 12 Literale + 1 Kommentar), Zeile 4 (7), Zeile 5 (12 Erwartungen in `run-integration-tests.sh`) | Implementer: die `diff`-Zeilen ergänzen, Soll 0 |
| Handbuch | Zeilen 7 und 8 (je 1): Hinweissatz und Platzhalter | Diff-Stand 0 beider; der Katalog-Abschnitt ist eine neue Fläche, kein Treffer |
| SDKs | Zeile 9: 0 | Nichtgefunden belegt; SDKs reichen Servertexte durch, kein Nachzug |
| Datenbank-Metadaten | Zeile 10: 16 | offen: Messung `COMMENT ON` am Rollout-Ergebnis, dann Entscheidung im Bericht |
| Fehlertext-Präfix `Fehlerklasse` | Zeile 11: 28 | Träger-Kandidat des Codes (ADR); Diff-Stand je nach Entscheidung |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt, und die **Architect-Entscheidung (A) liegt `Accepted` vor** — der
Slice beginnt mit ihr. Der Planner nimmt sie nicht vorweg. Fragen an den Architect (eine
Entscheidung, mehrere Teilfragen; die Optionen aus d-migrate stehen in §1 und sind **eine**
Option, kein Vorschlag):

1. **Schema:** Präfix (eigenes Projekt-Präfix? gemeinsamer Buchstabe `E`/`W`/`I` wie
   d-migrate?), Nummernraum und seine Gliederung (nach Fehlerklasse? nach Komponente?),
   Zählweise für Warnungen und Hinweise (Info), Länge der Nummer.
2. **Beziehung zu den sieben Fehlerklassen:** Code **in** der Klasse (Klasse bleibt
   Maschinen-Achse, Code feiner), Code **neben** der Klasse (zwei Achsen), oder Klasse aus
   dem Code ableitbar. Wirkung auf `cdc.process_heartbeat.error_class`, das Fehlerklassen-Label
   der Metriken (Metriken-Vertrag nicht brechen) und den Text `Fehlerklasse
   <klasse>: …`.
3. **Abbildung je Weg:** Log-Zeile, `diagnose`-Ausgabe, Exit-Ausgabe des Prozesses,
   HTTP-Fehlerkörper, gRPC-Status-Details, Fehlertexte der SQL-Funktionen und der
   Antrags-Zeile (`error_message`), Rollout-Skript. Gilt der Code als eigenes Feld (wie
   d-migrate) oder als Präfix im Text?
4. **Stabilität:** Zusage (ein Code wird nie neu belegt, Zurückziehen statt Löschen — d-migrate
   `E137`), Verhältnis zur Text-Stabilität (Text ist nicht Vertrag, Code ist es?), was ein
   parsender Betreiber behalten darf.
5. **Dokumentation und Prüfung:** Katalog im Handbuch (kennungsfrei) als einzige Quelle oder
   zusätzlich eine maschinenlesbare Registry (Ledger-Muster) und ein Test, der „kein Code
   ohne Test und Katalog-Eintrag“ prüft; Ort der Registry.
6. **Spec-Stelle:** Gehört die Festlegung in `SPEC-008` / ein neues `SPEC-*` /
   `LH-QA-OPS-001.a`, oder bleibt sie ADR-gebunden?
7. **Reichweite:** welche Ausgaben sind Meldungen mit Code (nur Fehler/Warnungen, wie der
   Auftraggeber sagt) und welche reine Ausgabe (Diagnose-Zeilen ohne Code, nur ohne
   Kennung).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die ADR verlangt eine Registry
  samt Test **und** die Umstellung aller Wege (Log, HTTP, gRPC, SQL) — dann trennt sich
  (B) in `…-diagnose-und-rollout` (Ausgaben dieses Plans) und `…-codes-http-grpc-sql`; die
  ADR, der Katalog und der Wächter (C) bleiben im ersten.
- `in-progress` → `open` (blockiert): die ADR kommt nicht zustande, oder die Messung zeigt,
  dass Betreiber die Ausgabe maschinell parsen und eine Kompatibilitätszusage fehlt
  (Auftraggeber-Entscheidung); die Änderung geht nie als Anpassung einer Erwartung durch.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration`, `make test-store` grün nach `make image`, die
Handbuch-Beispiel-Ausgabe mit einem realen `diagnose`-Lauf verglichen, Wächter-
Mutationsprobe rot gesehen, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **Ausgaben-Stabilität für Betreiber, die die Ausgabe parsen.** Die Zeilen der
  Diagnose-Ausgabe und der Rollout-Meldungen sind heute Text; wer sie mit `grep`
  auswertet, bricht an jeder Änderung. Das Handbuch sagt nicht ausdrücklich, ob der Text
  Vertrag ist (zu messen: Handbuch-Abschnitte Diagnose, Fehlerklassen, Exit-Codes).
  Gegenmittel: die ADR legt die Stabilitätszusage fest (Teilfrage 4), das Handbuch nennt
  den Code als stabilen Anker. — **Ausgang:** (bei Closure)
- **Tests erwarten Texte.** `diagnose_test.go` (7), `run-integration-tests.sh` (12 +
  Diagnose-Querabgleiche über HTTP/gRPC-CLI), `run-schema-rollout-guard-test.sh`: ein
  vergessener Läufer wird erst beim `make test-integration`-Lauf rot, der nicht in
  `make gates` liegt. Gegenmittel: Suchlauf-Zeilen je Träger, vollständiger Lauf vor
  Closure. — **Ausgang:** (bei Closure)
- **SDKs und Clients reichen Servertexte durch.** Ein SDK, das den Servertext einer Ablehnung
  in einer Ausnahme zeigt, zeigt künftig den Code mit; SDK-Tests, die auf den Wortlaut
  prüfen, wären betroffen. Gemessen ist nur, dass `sdks/` keine Kennung trägt (0), nicht
  dass kein Test einen Servertext erwartet — der Implementer misst das (`git grep` der
  Server-Fehlertexte unter `sdks/`). Keine SDK-Änderung ohne Freigabe (§1). — **Ausgang:** (bei Closure)
- **Die Datenbank-Metadaten (`COMMENT ON`) tragen Kennungen.** Ungemessen, ob der Betreiber
  sie sieht; eine Änderung der `description:`-Felder ändert das Rollout-Ergebnis und die
  Schema-Identität (Rollout-Wache, `make test-store`). — **Ausgang:** (bei Closure)
- **Wächter mit Falsch-Positiven oder -Negativen.** Eine zeilenbasierte Regel liest keine
  mehrzeiligen Literale; ein AST-Wächter kostet ein Docker-Programm (§3 Option). Ein
  zu eng gelesenes Muster erzeugt falsche Sicherheit — der Sensor-Vertrag nennt seine
  Grenze ausdrücklich ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)). — **Ausgang:** (bei Closure)
- **Kollision mit parallelen Arbeiten am Handbuch.** Mehrere Implementer ändern
  `docs/user/benutzerhandbuch.md` und Harness-Skripte; der Katalog-Abschnitt und die
  Beispiel-Anpassung konkurrieren um Zeilen. Gegenmittel: Rebase auf `main`, nur eigene
  Abschnitte committen. — **Ausgang:** (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird die Default-Sub-Area `*` (Kürzel
`PGC`, Modus Greenfield, [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); sie umfasst Produktions-Code, Skripte und Doku gleichermaßen, eine
feinere Aufteilung ist für diesen Slice nicht nötig — der Schnitt läuft über Ausgabe-Wege,
nicht über Pfade.

**Vorgelagert — offene Beobachtungen sichten:** das Register `../observations/BEO-PGC/` ist
nicht Inhalt dieses Plans; der Implementer sichtet es vor dem Start auf Treffer zu
„Kennung in Ausgabe“ und „Text-Erwartung im Läufer“ (nicht gemessen, kein Treffer
behauptet).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
