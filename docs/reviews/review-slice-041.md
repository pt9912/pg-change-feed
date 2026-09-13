# Review-Report: slice-041 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-041`, §1/§2/§3/§6/§8) und
`ADR-0052` (Accepted, alle sieben Entscheidungen zentraler Prüfmaßstab
dieses Laufs) sowie `AGENTS.md` §3 Hard Rules (§3.1, §3.5, §3.7).

**Gegenstand:** Commits `af94b22` (`feat(bootstrap): optionale
YAML-Konfigurationsdatei ergänzt Env-Vars (ADR-0052)`), `deeb54f`
(`docs(planning): slice-041 Plan-Nachzug + DoD-Belege (ADR-0052)`),
`d90d680` (`docs(planning): slice-041 Plan-Nachzug — Doku-Update-Ort
begründet (ADR-0052)`) — neu: `internal/bootstrap/config_file.go`,
`internal/bootstrap/config_file_internal_test.go`; geändert:
`cmd/pg-change-feed/main.go`, `internal/bootstrap/wiring.go`, `go.mod`,
`spec/pflichtenheft.md` (neu `SPEC-016`), `docs/user/benutzerhandbuch.md`,
`docs/plan/planning/in-progress/slice-041-yaml-konfigurationsdatei.md`
(DoD-Häkchen, Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-041-yaml-konfigurationsdatei.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug, §4
  Trigger, §6 Risiken, §8)
- `docs/plan/adr/0052-optionale-yaml-konfigurationsdatei.md` (vollständig
  — zentraler Prüfmaßstab, alle sieben Entscheidungen, Verglichene
  Alternativen, Konsequenzen, Re-Evaluierungs-Trigger)
- `AGENTS.md` §3.1 (Docker-only), §3.5 (ADR-Immutabilität), §3.7
  (Kommentar-Disziplin)
- `spec/pflichtenheft.md` §`SPEC-016` (neue Feldform-Verfeinerung, vor und
  nach dem Diff)
- `.a-check.yml` (Composition-Root-Freistellung für den neuen
  `yaml.v3`-Import)
- `.d-check.yml` (`matrix`-Regeln — kein Spec-Stratum nennt ADR/Slice)
- Alle neuen/geänderten Dateien vollständig gelesen, nicht nur die
  Implementer-Zusammenfassung: `config_file.go`,
  `config_file_internal_test.go`, `main.go`-Diff, `wiring.go`-Diff,
  `spec/pflichtenheft.md`-Diff, `docs/user/benutzerhandbuch.md`-Diff,
  `go.mod`-Diff, alle drei Commit-Diffs vollständig
- `docs/plan/planning/observations/BEO-PGC/*` (vollständiger Bestand,
  22 Einträge) — auf thematische Nähe zu diesem Slice geprüft
- `docs/reviews/review-slice-039.md` (Format-Vorlage)
- Real ausgeführt in dieser Sitzung: `make test` (grün, inkl.
  `internal/bootstrap`), `make docs-check` (0 Befunde), `make a-check`
  (0 Befunde), `make commit-traceability` (OK)

---

## Findings

### F-1 — Slice-Chronik in Go-Quellcode-Kommentar

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Disziplin), Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- `pfad`: `internal/bootstrap/config_file.go:186–192` (Doc-Kommentar von
  `mergeTables`)
- `befund`: Der Kommentar lautet: „`mergeTables` trägt die
  Implementer-Entscheidung zu §6 Risiko 1 des Slice-Plans (`slice-041`,
  `ADR-0052` entscheidet die Frage nicht explizit): …". Der Kommentar
  benennt explizit die Slice-Kennung `slice-041` als Begründungsquelle
  einer Code-Entscheidung — genau das Muster, das für dieses Repo bereits
  einmal korrigiert wurde (Go-Quellcode-Kommentare tragen keine
  Slice-Referenzen, nur den Ist-Zustand; Ausnahme sind ausschließlich die
  harness-eigenen Herkunfts-Anker in `AGENTS.md`/`harness/README.md`/
  `.harness/skills/*.md`/Beobachtungs-Register — `config_file.go` gehört
  zu keiner dieser Klassen). Der Slice-Plan ist zudem kein dauerhaftes
  Artefakt: Bei einer künftigen Welle-Closure wird er nach `done/`
  verschoben und kann später archiviert werden (`archiv.zip` +
  gekürzter Stub) — der Verweis referenziert dann eine Adresse, deren
  Inhalt sich ändert oder verschwindet, während der Code bestehen
  bleibt. `ADR-0052` selbst wäre der stabile, kanonische Verweis
  gewesen (Rang-Zeiger); die Slice-Kennung trägt hier keine der
  zulässigen Kommentar-Klassen (Zusage · Kopplung · Abgrenzung ·
  Rang-Zeiger · Grenze), sondern Chronik.
- `verifizierbar`: ja — `grep -n "slice-[0-9]" internal/bootstrap/*.go`
  (in dieser Sitzung ausgeführt, ein Treffer)
- `klasse`: „Slice-Chronik in Go-Quellcode-Kommentar"

### F-2 — Tables-Merge-Precedence als Implementer-Entscheidung ohne Folge-ADR in SPEC-016 verankert

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Modul 8 §Rollen-Regeln — „Implementer darf
  höchstens Folge-ADR vorschlagen, niemals stillschweigend einer ADR
  widersprechen" bzw. eine von der ADR offen gelassene Entscheidung
  eigenständig treffen und als bindend fortschreiben)
- `pfad`: `internal/bootstrap/config_file.go:189–201` (`mergeTables`),
  `spec/pflichtenheft.md:246–249` (`SPEC-016`, `tables`-Ersetzungsregel),
  Plan-Nachzug „Implementierungsentscheidung zu §6 Risiko 1"
- `befund`: `ADR-0052` benennt selbst, dass die `tables`-Merge-Semantik
  bei gleichzeitig gesetzter `CDC_TABLES` und Datei-`tables` „nicht
  explizit" entschieden ist (§6 Risiko 1 des Slice-Plans zitiert das).
  Es handelt sich um eine Precedence-Entscheidung mit echten,
  gegeneinander abzuwägenden Alternativen (Gesamt-Ersetzung vs.
  Element-weise Vermischung je Tabelle) — strukturell dieselbe Art
  Entscheidung, die `ADR-0052` selbst für die übrigen Felder trifft. Der
  Implementer trifft diese Entscheidung selbst, dokumentiert sie
  begründet im Plan-Nachzug und schreibt sie unmittelbar als bindende
  `SPEC-016`-Festlegung fest — ohne Architect-Verdikt oder Folge-ADR. Die
  Begründung ist inhaltlich nachvollziehbar (Konsistenz zur
  Ganzwert-Precedence der Stringfelder) und dadurch nicht willkürlich,
  aber die Entscheidung verlässt damit den Bereich „Implementierungsdetail
  innerhalb der ADR" und wird zu einer eigenständigen, extern sichtbaren
  Verhaltenszusage, für die die ADR selbst keine Abwägung trägt.
- `verifizierbar`: nein (kein Gate prüft, ob eine SPEC-Festlegung durch
  einen Architect-Zug oder eine ADR gedeckt ist — reines Prozess-Urteil)
- `klasse`: „Implementer-Entscheidung mit Architektur-Charakter ohne
  Folge-ADR"

### F-3 — Doppelte DSN-Ablehnung im aktuellen Kontrollfluss nicht unterscheidbar

- `kategorie`: LOW
- `quelle`: Maintainability (Klarheit der Fehlerpfad-Dokumentation)
- `pfad`: `internal/bootstrap/config_file.go:99–113` (`ConfigFromFile`)
- `befund`: Der Kommentar über `ConfigFromFile` beschreibt zwei
  „Verteidigungslinien" gegen DSN-Schlüssel in der Datei: den expliziten
  Schlüssel-Check auf dem roh eingelesenen `map[string]any` und die
  implizite Ablehnung durch `KnownFields(true)`, weil die drei
  DSN-Schlüssel im `fileConfig`-Typ nicht deklariert sind. Der explizite
  Check läuft im Code jedoch **immer zuerst** und kehrt bei einem Treffer
  sofort zurück, bevor `decoder.Decode` erreicht wird — die „implizite"
  Linie ist unter dem aktuellen Kontrollfluss für diese drei Schlüssel
  nicht beobachtbar/erreichbar, solange der explizite Check unverändert
  bleibt (per statischer Analyse bestätigt: `fileConfig` trägt keines der
  drei verbotenen Felder, aber der Codepfad dorthin ist durch das frühe
  `return` blockiert). `TestConfigFromFileLehntDSNAb` bestätigt das:
  Es prüft nur die Fehlerklasse (`ErrConfiguration`), nicht welcher der
  beiden Pfade sie geliefert hat, und würde unverändert grün bleiben,
  wenn ausschließlich `KnownFields` die Ablehnung trüge. Kein
  Funktionsfehler — die Redundanz ist als Zukunftsabsicherung begründet
  —, aber die Kommentar-Formulierung suggeriert zwei aktuell gleichzeitig
  wirksame Linien, wo real nur eine beobachtbar ist.
- `verifizierbar`: ja (Code-Inspektion in dieser Sitzung; ein isolierter
  Test, der den expliziten Check laufzeit-lokal umginge, würde das
  bestätigen — nicht Teil der bestehenden Suite)
- `klasse`: „Verteidigungslinie im aktuellen Kontrollfluss unerreichbar"

### F-4 — Testkommentar referenziert „Slice-Plan" statt kanonischer Quelle

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Disziplin), verwandt zu F-1
- `pfad`: `internal/bootstrap/config_file_internal_test.go:227–231`
  (Doc-Kommentar von `TestMergeConfigTabellenCDCTablesSchlaegtDatei`)
- `befund`: „trägt die Implementer-Entscheidung zu §6 Risiko 1 des
  Slice-Plans: …" — milder als F-1 (keine explizite Slice-Kennung, keine
  Chronik-Sprache wie „jetzt"/„nur noch"), aber derselbe Verweis-Typ: Ein
  Testkommentar begründet sich über ein Planungsartefakt ohne dauerhafte
  Identität statt über `ADR-0052`/`SPEC-016`, die beide kanonisch und
  stabil sind.
- `verifizierbar`: ja — `grep -n "Slice-Plan" internal/bootstrap/*_test.go`
- `klasse`: „Slice-Chronik in Go-Quellcode-Kommentar" (verwandtes,
  schwächeres Vorkommen derselben Klasse wie F-1)

### F-5 — `main.go`-Rewiring: Scope-Erweiterung sachlich gerechtfertigt

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form:
  Slice / Modul 8 §Konflikt-Pfad als Rollen-Sequenz
- `pfad`: `cmd/pg-change-feed/main.go` (alle fünf Aufrufstellen),
  Plan-Nachzug „Dateien — Abweichung von der Tabelle oben"
- `befund`: Die ursprüngliche §3-Dateiliste nannte `cmd/pg-change-feed/main.go`
  nicht; der Implementer hat die Erweiterung vor dem Review in einem
  eigenen Plan-Nachzug-Commit (`deeb54f`) dokumentiert und begründet.
  **Urteil:** Die Erweiterung ist gerechtfertigt und kein
  Rollen-Konflikt. Drei Gründe: (1) `ADR-0052` Entscheidung 5 setzt
  ausdrücklich voraus, dass `CDC_CONFIG_FILE` „im Container" wirkt —
  ohne das Rewiring bliebe die Variable im echten Binary tot, die ADR
  wäre nur auf Paket-Ebene umgesetzt, nicht in ihrer eigentlichen
  Absicht. (2) Keiner der vier §1-Ausschlusspunkte erfasst dies: Der
  ausgeschlossene `--config`-CLI-Flag ist ein anderer Zugriffsweg
  (Kommandozeile) als das Einlesen der bereits als Env-Var-Vertrag
  bestehenden `CDC_CONFIG_FILE` (Entscheidung 5 selbst). (3) Die Änderung
  bleibt innerhalb der Composition Root (`cmd/**` und
  `internal/bootstrap/**` sind laut `.a-check.yml` beide Teil davon),
  berührt keine neue Schicht und ist rein mechanisch (Funktionsaufruf
  ausgetauscht, kein neues Verhalten in `main.go` selbst). Eine
  Vorab-Planänderung *vor* der Umsetzung wäre sauberer gewesen, aber der
  hier genutzte Plan-Nachzug ist in diesem Repo ein etablierter,
  akzeptierter Mechanismus für genau diese Klasse von Abweichung
  (`BEO-PGC/plan-nachzug`) und wurde korrekt angewendet. Keine
  Architect-Sequenz nach Modul 8 erforderlich — kein HIGH, kein
  bestrittener Rollen-Widerspruch.
- `verifizierbar`: ja (`.a-check.yml` `composition_root`, Diff-Inspektion)
- `klasse`: „Plan-Nachzug für Scope-Erweiterung — Urteil: gerechtfertigt"
  (kein Fehlermuster, dokumentierter Positiv-Befund)

## Negativbefunde

- geprüft, ohne Befund: **`ConfigFromEnv` bleibt exakt unverändert.**
  `git show af94b22 -- internal/bootstrap/wiring.go` ändert
  ausschließlich den Datei-Kopf-Kommentar (Zeilen 7–15); der Funktionskörper
  von `ConfigFromEnv` (Zeile 156ff.) ist byteidentisch zum Vorzustand.
  `TestConfigFromEnvAndFileLeereEnvVariable` vergleicht beide Pfade
  Feld-für-Feld real.
- geprüft, ohne Befund: **DSN-Felder bleiben env-var-exklusiv.**
  `fileConfig` trägt keine Felder für `capture_dsn`/`admin_dsn`/
  `reader_dsn`; `mergeConfig` liest alle drei DSNs ausschließlich über
  `getenv(...)`, nie aus `file`.
- geprüft, ohne Befund: **`.a-check.yml`-Konformität.** Der neue
  `yaml.v3`-Import und die Importe von `mapper`/`model` in
  `config_file.go` liegen vollständig in `composition_root`
  (`internal/bootstrap/**`); `make a-check` real ausgeführt: 0 Befunde.
- geprüft, ohne Befund: **`spec/pflichtenheft.md` `SPEC-016` verletzt
  keine Matrix-Regel.** Kein `ADR-*`- oder `slice-*`-Verweis innerhalb
  von `SPEC-016` oder sonst im Diff der Datei (`grep -n "ADR-\|slice-"
  spec/pflichtenheft.md` — die einzigen Treffer liegen in
  unveränderten, allgemeinen Abschnitten außerhalb des Diffs); `make
  docs-check` real ausgeführt: 0 Befunde über 322 Dateien.
- geprüft, ohne Befund: **Docker-only (`AGENTS.md` §3.1).** Kein
  lokales Toolchain-Setup im Diff; alle real ausgeführten Prüfungen
  dieser Sitzung liefen über `make test`/`make docs-check`/`make
  a-check`/`make commit-traceability` (Docker-Container).
- geprüft, ohne Befund: **Kein Accepted-ADR verändert (`AGENTS.md`
  §3.5).** `ADR-0052` bleibt `Accepted` und inhaltlich unverändert; der
  Diff setzt sie um, ändert sie nicht.
- geprüft, ohne Befund: **`compose.yaml`/`harness/README.md` bleiben
  unverändert**, wie in §1 Out-of-Scope und im zweiten Plan-Nachzug-Commit
  (`d90d680`) begründet — kein Diff auf beiden Dateien in diesem
  Commit-Bereich.
- geprüft, ohne Befund: **11 Testfälle in
  `config_file_internal_test.go` decken die im Plan versprochenen Fälle
  ab** (striktes Decoding, DSN-Ablehnung ×3, nicht lesbare Datei, leere
  Datei, Tabellen-Mapping, Env-only-Regression, Datei-nicht-ladbar,
  Einzelfeld-Override, `CDC_TABLES`-Vorrang, fehlende Pflichtfelder,
  `log_level`/WAL-Retention aus der Datei) — alle real grün
  (`make test`, in dieser Sitzung ausgeführt). Die von F-3 benannte
  Einschränkung betrifft nur die Unterscheidbarkeit zweier
  DSN-Ablehnungspfade, nicht die Testabdeckung insgesamt.
- geprüft, ohne Befund: **`docs/user/benutzerhandbuch.md`-Update
  vollständig und konsistent.** Neue Unterüberschrift „Optionale
  YAML-Konfigurationsdatei", `CDC_CONFIG_FILE`-Zeile in der
  Env-Var-Tabelle, `CDC_TABLES`-Pflichtangabe korrekt auf „falls keine
  Konfigurationsdatei dieselbe Aktivierung trägt" präzisiert,
  Versionshistorie-Zeile 1.6 mit `ADR-0052`/`SPEC-016`/`slice-041`-Bezug
  (Versionshistorie-Tabellen sind der etablierte Ort für
  Slice-Zeitbezüge in diesem Dokument, siehe Zeile 1.5 — keine Abweichung
  von der in `feedback_no_chronicle_comments`-Regel benannten Ausnahme
  für Doku-Changelogs).
- geprüft, ohne Befund: **`go.mod` korrekt auf direkte Abhängigkeit
  umgestellt**, keine Versionsänderung von `gopkg.in/yaml.v3` (bleibt
  `v3.0.1`, war bereits transitiv vorhanden — kein neues
  Supply-Chain-Risiko).
- geprüft, ohne Befund: **Beobachtungs-Register-Sichtung im Slice-Plan
  (§8) korrekt.** Eigenständig gegen den vollständigen `BEO-PGC/*`-Bestand
  (22 Einträge) geprüft: kein Eintrag ist zu `CDC_CONFIG_FILE`/YAML-
  Konfigurationsdatei thematisch einschlägig; die im Plan genannten zwei
  Treffer (`github-actions-unverifizierbar-lokal`,
  `verwaltung-keine-sql-administration`) sind korrekt als nicht
  einschlägig eingestuft.
- geprüft, ohne Befund: **`make commit-traceability` real grün** über
  alle drei Commits dieses Slice (`ADR-0052` im Betreff jeder der drei
  Commit-Messages, kein `SPEC-*`/`ARC-*` im Betreff).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Slice-Chronik in
Go-Quellcode-Kommentar" (2×, F-1 und F-4 — F-1 als vollständige, F-4 als
schwächere Ausprägung derselben Klasse; erstes Auftreten dieser Klasse
*in einem Review-Report* dieses Repos, auch wenn dasselbe Muster einmal
zuvor außerhalb des Review-Prozesses per Nutzer-Korrektur behoben wurde)
· „Implementer-Entscheidung mit Architektur-Charakter ohne Folge-ADR" (1×,
erstes Auftreten) · „Verteidigungslinie im aktuellen Kontrollfluss
unerreichbar" (1×, erstes Auftreten) · „Plan-Nachzug für
Scope-Erweiterung — Urteil: gerechtfertigt" (1×, Positiv-Befund, keine
Fehler-Klasse) — kein Steering-Loop-Eintrag zwingend fällig (keine Klasse
erreicht 3×), aber F-1/F-4 zusammen sind ein Wiederholungssignal: Dasselbe
Muster (Slice-Referenz statt ADR-Referenz in Code-Kommentaren) trat bereits
einmal außerhalb dieses Review-Prozesses auf (Nutzer-Korrektur, nicht im
Beobachtungs-Register erfasst). Empfehlung an die Planner-Rolle bei
Closure: prüfen, ob `BEO-PGC/…` einen neuen Eintrag für diese Klasse
braucht, damit künftige Wiederholungen gezählt statt jedes Mal neu entdeckt
werden.

## Verdikt

**Merge-blockierend:** ja im Sinne der Closure — der Commit ist bereits
gepusht (Repo-Konvention: Push vor Review üblich), aber der Slice darf
wegen F-1 (HIGH) nicht nach `done/` übergehen, bevor der Kommentar
korrigiert ist (Slice-Referenz durch eine Ist-Zustands-Beschreibung mit
`ADR-0052`/`SPEC-016`-Verweis ersetzt). F-2 (MEDIUM) blockiert die
Closure nicht zwingend, sollte aber vor Closure entschieden werden:
entweder der Implementer/Planner akzeptiert die Entscheidung als
hinreichend begründetes Implementierungsdetail (dann bleibt es so, mit
Vermerk in der Closure-Notiz), oder sie wird per Folge-ADR bestätigt.
Kein Rollen-Widerspruch liegt vor (der Implementer hat die offene Frage
selbst benannt, nicht bestritten) — die Architect-Sequenz aus Modul 8
greift damit nicht zwingend, ist aber die sauberere Option, falls die
Precedence-Frage später erneut strittig wird.

**Zur zentralen Prüffrage dieses Laufs (`ADR-0052`-Umsetzungstreue):**
Bestätigt in allen sieben Entscheidungen — striktes Decoding (1), Env
schlägt Datei Feld für Feld (2), Env-only-Pfad real unverändert (3),
Parser-Ort Composition Root (4), `CDC_CONFIG_FILE`-Zugriffsweg ohne
CLI-Flag (5), DSN-Ablehnung env-var-exklusiv (6, mit der in F-3
benannten Klarheits-Einschränkung), Fehlerklasse `configuration`
durchgängig (7). **Zur Scope-Frage `main.go`:** gerechtfertigt (F-5).
**Zur `tables`-Merge-Semantik:** vertretbar, aber prozessual lückenhaft
dokumentiert (F-2).

**Übergabe:** Ein HIGH-Finding (F-1, Kommentar-Korrektur vor Closure
zwingend), ein MEDIUM-Finding (F-2, Entscheidung des Implementers/
Planners vor Closure), zwei LOW-Findings (F-3, F-4), ein INFO-Urteil
(F-5, keine Aktion nötig). Der Implementer entscheidet über Annahme oder
Begründung der MEDIUM-/LOW-Findings; F-1 ist keine Ermessensfrage
(etablierte, bereits einmal korrigierte Konvention). Dieser Report ist
ein Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen — die
Summary-Zeile speist bei Bedarf den Closure-Eintrag (Modul 5). Er ersetzt
keine Verifikation — DoD-/Spec-Konformität (inkl. der übrigen offenen
DoD-Punkte: Review-Häkchen selbst, Closure-Notiz, Register,
Risiko-Ausgänge, drei Paarungen) prüft der Verifier separat.
