# Verifikationsbericht: slice-041 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen
Plan (`slice-041` §1/§2 DoD/§3/§6/§8) und `ADR-0052`, nicht gegen Diff
(Reviewer-Aufgabe, drei Runden bereits abgeschlossen) und nicht gegen
realen Bedarf (Validator, hier nicht ausgelöst — kein
MVP-Meilenstein-Slice).

**Frischer Kontext:** Diese Prüfung liest den vollständigen, aktuellen
Slice-Plan, `ADR-0052` vollständig, alle drei Review-Reports und den
tatsächlichen Code selbst — keine Behauptung aus einem Bericht wird
ungeprüft übernommen; jeder unten genannte Sensor-/Werkzeuglauf wurde in
dieser Sitzung **selbst** ausgeführt.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-041-yaml-konfigurationsdatei.md`
zum Stand `HEAD = 27f34a7`. Commits: `eb1bb24` (`next→in-progress`),
`af94b22` (Implementierung), `deeb54f`/`d90d680` (Plan-Nachzüge),
`f98b0b3` (Review: 1 HIGH, 1 MEDIUM, 2 LOW, 1 INFO), `35a5279` (Fixrunde:
F-1/F-3/F-4 behoben, F-2 dokumentiert akzeptiert), `25b7a6c`
(Review-Bestätigung, drittes unabhängiges Chronik-Vorkommen neu
gefunden — HIGH), `8ceee6b` (drittes Vorkommen behoben), `7fc8b90`
(Beobachtungs-Register-Eintrag `BEO-PGC/slice-chronik-in-code-kommentar`,
2×, unter Schwelle), `27f34a7` (finale Review-Bestätigung). Lifecycle
verifiziert (`git log --follow`): `angelegt → Verantwortlich → open→next
→ next→in-progress → Implementierung → Review → Fix → Fixrunde →
Fix → Fixrunde-Bestätigung`, keine Rolle springt rückwärts ohne
Übergabe-Artefakt, kein Rollback (`in-progress→next`/`open`). Die
thematisch fremden Zwischen-Commits (`4e4c7bb` Docs-Check-Fix am
Review-Report, `083caf3` `BEO-PGC/commit-traceability-kein-vorab-hook`,
sowie `de58be8`/`6fc4d31`/`0a19983` für `slice-042`) liegen außerhalb
dieses Gegenstands und wurden nicht geprüft.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `ConfigFromFile`: striktes Decoding (`KnownFields(true)`), unbekannter Schlüssel → `ErrConfiguration`, DSN-Schlüssel in Datei → `ErrConfiguration` | **erfüllt** | `internal/bootstrap/config_file.go:83-109` gelesen: `forbiddenFileDSNKeys`-Check auf rohem `map[string]any` vor dem Decode, danach `decoder.KnownFields(true)`. Eigene Testläufe: `TestConfigFromFileStriktesDecoding`, `TestConfigFromFileLehntDSNAb` (alle drei Schlüssel einzeln) grün in `make test` |
| 2 | Merge-Funktion: Env schlägt Datei Feld für Feld, real gegen Einzelfeld-Override getestet | **erfüllt** | `mergeConfig` (Z. 148-193) nutzt `overrideString` je String-Feld; `TestMergeConfigEinzelnesFeldPerEnv` setzt nur `CDC_SLOT`, prüft `Source`/`Publication` bleiben aus Datei, `Slot` aus Env — Test gelesen und grün |
| 3 | `CDC_CONFIG_FILE` verdrahtet: leer → exakt heutiger `ConfigFromEnv`-Pfad, Regressionstest | **erfüllt** | `ConfigFromEnvAndFile` (Z. 118-128): `path == ""` → `return ConfigFromEnv(getenv)`, keine Umleitung. `git diff eb1bb24..27f34a7 -- internal/bootstrap/wiring.go` zeigt: einzige Änderung ist der Datei-Kopf-Kommentar (Z. 7-13), `ConfigFromEnv`s Funktionskörper unverändert. `TestConfigFromEnvAndFileLeereEnvVariable` vergleicht beide Pfade Feld für Feld — real grün |
| 4 | Tabellen-Aktivierung als YAML-Mapping, eigener Test mit mehreren Tabellen | **erfüllt** | `fileConfig.Tables map[string]fileTableBinding`; `TestConfigFromFileTabellenMapping` lädt zwei Tabellen (`public.t1`, `public.t2`), prüft beide — grün |
| 5 | `spec/pflichtenheft.md` trägt `SPEC-<NNN>`-Verfeinerung | **erfüllt** | `SPEC-016` neu (`spec/pflichtenheft.md:210-249`): Schlüsseltabelle, YAML-Beispiel, `tables`-Ersetzungsregel, `CDC_CONFIG_FILE`-Semantik. Kein `ADR-*`/`slice-*`-Token **innerhalb** der `SPEC-016`-Sektion (`grep -n "ADR-\|slice-" spec/pflichtenheft.md` — Treffer liegen außerhalb des Diffs) — `matrix`-Regel `spec→adr: false`/`spec→slice: false` real eingehalten, durch `make gates`/`d-check` (334 Dateien, 0 Befunde) bestätigt |
| 6 | `make gates` grün, `make test` grün | **erfüllt, selbst reproduziert** | Eigener Lauf gegen `HEAD=27f34a7`: `baseline-verify` (v6.5.0, 54 Dateien) OK, `d-check` (334 Dateien, 0 Befunde, inkl. `commits`-Range-Lauf `HEAD~5..HEAD`), `commit-traceability` OK, `a-check` 0 Befunde. `make test` (Docker, `-race`, `golang:1.27`): alle Pakete grün inkl. `internal/bootstrap` (1.052s) |
| 7 | Review durchgeführt, Report liegt vor | **inhaltlich erfüllt, Formular-Diskrepanz — siehe Finding V-1** | Drei Reports vollständig gelesen (`review-slice-041.md`, `-fixrunde.md`, `-fixrunde-2.md`); letzter Stand: „Reviewer-seitig keine Einwände gegen einen Übergang nach Verifikation." Bedingung tatsächlich erfüllt — aber Checkbox in §2 (Zeile 119) steht weiterhin auf `- [ ]` (`git log -p --follow` bestätigt: seit Anlage nie auf `[x]` gesetzt) |
| 8 | Doku-Update `docs/user/benutzerhandbuch.md` | **erfüllt** | `git diff eb1bb24..27f34a7 -- docs/user/benutzerhandbuch.md` gelesen: neue Unterüberschrift „Optionale YAML-Konfigurationsdatei" (§5), `CDC_CONFIG_FILE`-Zeile in der Env-Var-Tabelle, `CDC_TABLES`-Pflichtangabe korrekt auf „falls keine Konfigurationsdatei …" präzisiert, Versionshistorie-Zeile 1.6. `harness/README.md` unverändert (`git diff` leer) — deckt sich mit der im Plan-Nachzug begründeten Doku-Update-Ort-Entscheidung |
| 9 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit (Modul 8); §7 ist noch die Bedienhinweis-Vorlage |
| 10 | Reconciliation-Register | **entfällt strukturell** | `docs/plan/planning/reconciliation.md` existiert nicht (`find` ohne Treffer); Repo führt laut `harness/conventions.md` ausschließlich Greenfield, kein Brownfield-Bootstrap — Item nennt diese Bedingung selbst |
| 11 | Beobachtungs-Register fortgeschrieben | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit. Zur Einordnung selbst geprüft: `BEO-PGC/slice-chronik-in-code-kommentar/` (2 Evidence-Dateien) und `BEO-PGC/commit-traceability-kein-vorab-hook/` existieren bereits — beide sind Nebenprodukte der Fixrunden, nicht Gegenstand dieser Prüfung, aber ihre bloße Existenz widerspricht dem „offen"-Status des DoD-Punkts nicht: Der Punkt bleibt unbeansprucht, weil kein Häkchen gesetzt und keine Closure-Notiz geschrieben ist |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — korrekt unbeansprucht** | Beide Risiken in §6 tragen noch `<bei Closure einzutragen>` |
| 13 | Drei Paarungen | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit, wellenlos hier statt bei einer Welle-Closure fällig, erst nach dem `git mv` nach `done/` sinnvoll prüfbar |

## 2. Finding V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen

- **Klasse:** Verifier-only — für Tests und Review unsichtbar, weil beide
  Rollen ihre eigene Arbeit erledigt haben; nur ein Blick auf den
  *Formular-Zustand nach* der dreirundigen Review-Sequenz deckt die Lücke
  auf. Dieselbe Finding-Klasse wie `verify-slice-039.md` V-1 — zweites
  Auftreten in diesem Repo (Beleg für den Planner: ggf. ein Kandidat für
  das Beobachtungs-Register, sollte ein drittes Auftreten folgen).
- **Befund:** DoD-Punkt 7 in
  `slice-041-yaml-konfigurationsdatei.md:119` steht auf `- [ ]`, obwohl
  die Bedingung — Review durchgeführt, Report liegt vor, alle Findings
  bestätigt behoben — seit Commit `27f34a7` faktisch und vollständig
  erfüllt ist. Keiner der drei Review-/Fixrunden-Commits hat die Checkbox
  im Plan-Dokument nachgezogen.
- **Einordnung:** kein inhaltlicher Mangel an Implementierung oder
  Review-Substanz — beide sind, wie oben belegt, real und reproduzierbar
  erfüllt (inkl. eines vollständigen dritten Fixzyklus für ein
  eigenständig nachgefundenes HIGH-Finding). Es ist eine Diskrepanz
  zwischen dem DoD-Formular und der tatsächlichen Sachlage.
- **Erwartete Korrektur:** Checkbox auf `[x]` setzen, in einem eigenen,
  kleinen Commit vor dem `git mv` nach `done/` (Inhalt vor Move, Modul 5
  §git mv + Inhaltsänderung). Kein Rollback, keine Rückführung.

## 3. Chronik-Freiheit der Kommentare — eigenständiger Vollständigkeits-Grep

Nicht auf die drei Fixrunden-Berichte vertraut — eigener Grep gegen ein
Vokabular, das den in den Reports genannten Suchbegriffen entspricht
(`slice-[0-9]`, `Slice-Plan`, `vor diesem`, `seit `, `nur noch`, `jetzt `,
`heutigen`, `nicht mehr`, `zuvor`, `früher`, `vormals`, `inzwischen`,
`mittlerweile`), über **alle vier** von diesem Slice geänderten
Go-Dateien (`git diff eb1bb24..27f34a7 --name-only -- '*.go'`):
`internal/bootstrap/config_file.go`,
`internal/bootstrap/config_file_internal_test.go`,
`cmd/pg-change-feed/main.go`, `internal/bootstrap/wiring.go`.

- `config_file.go` / `config_file_internal_test.go` / `main.go`: **kein
  Treffer.**
- `wiring.go`: drei Treffer (Zeilen 85, 578, 794) — alle drei **außerhalb
  des Diffs dieses Slice** (`git blame`: Commits `8270fe41` 2026-09-10,
  `29a49ead` 2026-09-12, `2a771bc2` 2026-09-13, alle vor `af94b22`); der
  einzige von diesem Slice geänderte Bereich in `wiring.go` ist der
  Datei-Kopf-Kommentar (Zeilen 7-13), der selbst chronikfrei ist. Diese
  drei Bestandstreffer sind nicht Gegenstand von F-1/F-4 oder dieser
  Prüfung (vorbestehender Code, nicht durch diesen Slice berührt).

**Ergebnis:** In den tatsächlich von `slice-041` geänderten Dateien ist
jeder Kommentar frei von Slice-Chronik — das dritte, eigenständig von der
Fixrunde-1-Review gefundene Vorkommen (`config_file.go:115`,
`ConfigFromEnvAndFile`) ist mit `8ceee6b` real behoben und durch die
Fixrunde-2-Review sowie diesen eigenen Lauf bestätigt.

## 4. `ADR-0052` — alle sieben Entscheidungen einzeln im Code bestätigt

| # | Entscheidung | Beleg |
|---|---|---|
| 1 | Format YAML, striktes Decoding | `gopkg.in/yaml.v3` jetzt direkte Abhängigkeit (`go.mod`-Diff: `+ gopkg.in/yaml.v3 v3.0.1`, Version unverändert gegenüber der bisherigen transitiven — kein neues Supply-Chain-Risiko); `decoder.KnownFields(true)` (`config_file.go:101`) |
| 2 | Env schlägt Datei, Feld für Feld | `overrideString` je String-Feld (`Source`/`Publication`/`Slot`); `tables` als Ganz-Feld ersetzt statt vermischt — Abweichung explizit im Plan-Nachzug begründet und im Review als F-2 disponiert (siehe §5 unten) |
| 3 | Abwärtskompatibilität, Env-only unverändert | `ConfigFromEnv`-Funktionskörper byteidentisch (`git diff` bestätigt), `ConfigFromEnvAndFile` delegiert bei leerem Pfad vollständig; Regressionstest grün |
| 4 | Parser-Ort Composition Root | `config_file.go` liegt in `internal/bootstrap` (`composition_root: ["internal/bootstrap/**", "cmd/**", "test/integration/**"]` in `.a-check.yml`); `make a-check` selbst ausgeführt: 0 Befunde |
| 5 | `CDC_CONFIG_FILE`-Zugriffsweg, kein CLI-Flag | `envConfigFile = "CDC_CONFIG_FILE"` (`config_file.go:27`); `cmd/pg-change-feed/main.go` importiert `flag` weiterhin nicht (`grep -n "flag\." main.go` → kein Treffer); alle fünf Aufrufstellen rufen jetzt `ConfigFromEnvAndFile` (`git diff` gelesen, alle fünf Stellen einzeln geprüft) |
| 6 | Secrets env-var-exklusiv | `forbiddenFileDSNKeys` (drei Schlüssel), Check vor dem Decode; `fileConfig` deklariert keine DSN-Felder; `mergeConfig` liest alle drei DSNs ausschließlich über `getenv(...)`; `TestConfigFromFileLehntDSNAb` deckt alle drei Schlüssel einzeln ab |
| 7 | Fehlerklasse `configuration`, kein neuer Typ | `ErrConfiguration`-Sentinel aus `wiring.go:73` unverändert wiederverwendet, `%w`-Wrapping durchgängig in allen neuen Fehlerpfaden (`config_file.go`) |

## 5. F-2 — Rollen-Einordnung der `tables`-Merge-Precedence, eigenständig geprüft

`ADR-0052` lässt die `tables`-Merge-Semantik bei gleichzeitig gesetzter
`CDC_TABLES` und Datei-`tables` offen (§6 Risiko 1 des Slice-Plans
zitiert das korrekt). Der Reviewer stufte die Implementer-Entscheidung
(Gesamt-Ersetzung statt Element-Vermischung) im ersten Report als
MEDIUM ein — Entscheidung mit Architektur-Charakter, aber **kein**
bestrittener Rollen-Widerspruch. Der Plan-Nachzug trägt seit `35a5279`
den Absatz „Rollen-Einordnung dieser Entscheidung" (§3, Zeilen 222-236):
er übernimmt das Reviewer-Verdikt unverfälscht, benennt die Entscheidung
ausdrücklich als vom Implementer akzeptiertes Implementierungsdetail
innerhalb der bestehenden ADR (Begründung: Konsistenz zur
Ganzwert-Precedence der übrigen Felder) und hält den Folge-ADR-Pfad für
den Fall offen, dass die Frage künftig erneut strittig wird. Die
Fixrunde-1-Review bestätigt das als „prozessual geschlossen, kein neuer
Fehler" — eigene Prüfung deckt sich damit: Der Absatz ist vollständig,
unverfälscht und dokumentiert klar, warum kein Architect-Zug (Modul 8)
zwingend ausgelöst wurde. `SPEC-016` trägt dieselbe Regel konsistent
(„`tables` als Feld wird als Ganzes ersetzt …").

## 6. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

Eigenständig geprüft:

- Kein `--config`-CLI-Flag: `cmd/pg-change-feed/main.go` importiert
  `flag` weiterhin nicht.
- `CDC_TABLES`/`parseTables` nicht zurückgebaut:
  `func parseTables(raw string) …` weiterhin vorhanden
  (`wiring.go:215`), `envTables = "CDC_TABLES"` weiterhin verdrahtet
  (`wiring.go:60`).
- `compose.yaml` unverändert (`git diff eb1bb24..27f34a7 -- compose.yaml`
  → leer).
- `spec/pflichtenheft.md`-Verfeinerung als Teil **dieses** Slice
  geliefert, nicht als separater vorgeschalteter Slice (§`SPEC-016` im
  selben Commit-Bereich).

Alle vier Ausschlusspunkte real eingehalten; das dokumentierte
`main.go`-Rewiring (nicht in der ursprünglichen §3-Tabelle) ist im
Plan-Nachzug begründet und vom Review als F-5 (INFO, gerechtfertigt)
disponiert — eigene Prüfung deckt sich mit dieser Einordnung: Die
Erweiterung bleibt innerhalb der Composition Root und setzt nur
`ADR-0052` Entscheidung 5 tatsächlich im Binary um.

## 7. Weitere Hard-Rule-Prüfungen (eigenständig)

- **`AGENTS.md` §3.1 (Docker-only):** Alle real ausgeführten Prüfungen
  dieser Sitzung liefen über `make gates`/`make test` (Docker-Container),
  kein lokales Toolchain-Setup im Diff.
- **`AGENTS.md` §3.5 (ADR-Immutabilität):** `ADR-0052` seit seiner
  Erstellung nicht mehr verändert — dieser Slice setzt sie um, ändert sie
  nicht.
- **`AGENTS.md` §3.6 (keine Gate-Lockerung ohne ADR):** Keine
  Schwellen-Änderung in diesem Diff; `make gates` unverändert
  konfiguriert.
- **`AGENTS.md` §3.7 (Kommentar-Disziplin):** siehe §3 oben —
  vollständig erfüllt in allen von diesem Slice geänderten Dateien.
- **Commit-Traceability:** alle Commits dieses Slice tragen `ADR-0052`
  im Betreff, kein `SPEC-*`/`ARC-*`-Struktur-ID-Treffer im Betreff
  (`make gates`-Lauf, Modul `commits`, real grün).

## 8. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Wie in Auftrag benannt: Closure-Notiz (§7), Reconciliation-/
Beobachtungs-Register-Fortschreibung, Risiko-Ausgänge (§6), die drei
Paarungen (Anker · Folge-Slice · Register) — alle vier zugehörigen
§2-Häkchen sind **korrekt unbeansprucht**, das ist der vorgesehene
Zustand vor dem nächsten Rollenwechsel an den Planner (Modul 8 §Rollen-
Sequenz für einen Slice: Verifier→Planner-Übergabe, nicht die
Planner-Closure-Arbeit selbst). Auch nicht Gegenstand: Validierung gegen
realen Bedarf (kein MVP-Meilenstein-Slice, kein Validator-Zug
ausgelöst); der thematisch fremde Commit `083caf3`
(`BEO-PGC/commit-traceability-kein-vorab-hook`) und der aus dem
Prüffenster gerutschte Traceability-Verstoß in `4e4c7bb`.

## Verdikt

**DoD-Konformität: bestätigt**, mit einer benannten, nicht
merge-blockierenden Formular-Diskrepanz (V-1: Checkbox „Review
durchgeführt" nicht nachgezogen — inhaltlich seit `27f34a7` vollständig
erfüllt). Alle acht substanziellen Implementer-DoD-Punkte (1-6, 8) sind
durch eigene, unabhängige Reproduktion gedeckt: striktes YAML-Decoding,
DSN-Ablehnung (alle drei Schlüssel), Feld-für-Feld-Precedence,
Env-only-Regression, Tabellen-Mapping, `SPEC-016`, `make gates`/`make
test -race` selbst grün gelaufen, Doku-Update verifiziert. Die fünf
Planner-Closure-Punkte (9, 11, 12, 13, sowie Punkt 10 strukturell
entfallend) sind korrekt offen.

**Zu `ADR-0052` (zentraler Prüfmaßstab, alle sieben Entscheidungen):**
bestätigt — jede Entscheidung einzeln im Code nachgewiesen (§4 oben),
inklusive der `tables`-Merge-Precedence als bewusst akzeptiertes
Implementierungsdetail (F-2, §5 oben) statt eines stillen ADR-Verstoßes.

**Zur Kommentar-Disziplin (`AGENTS.md` §3.7):** Der dreirundige
Review-Prozess (drei unabhängige HIGH-Funde derselben Klasse in
derselben Datei) hat vollständig gegriffen — eigener Vollständigkeits-
Grep über alle vier geänderten Go-Dateien findet keinen verbleibenden
Chronik-Treffer innerhalb des Diffs dieses Slice.

**Keine Rückführung nötig.** Weder `in-progress→next` (Slice ist nicht
zu groß — fünf Liefer-Punkte aus §2 sauber erfüllt innerhalb der
ursprünglichen zwei Schichten `internal/bootstrap`/`cmd`, keine vierte
Schicht entstanden) noch `in-progress→open` (kein Blocker; §4 des Plans
nennt „kein bekannter Blocker", das gilt weiterhin). Vor dem `git mv`
nach `done/` ist lediglich die Checkbox-Korrektur aus V-1 fällig — ein
Ein-Zeilen-Commit, kein Zerlegungs- oder Blocker-Fall.

**Übergabe an Planner:** Dieser Bericht bestätigt DoD-Konformität für die
Closure-Entscheidung, mit dem Hinweis V-1 zur Nachbesserung vor dem
`git mv`. Die Planner-Rolle prüft zusätzlich, ob die im Fixrunde-1-Review
empfohlene Steering-Loop-Einordnung (`BEO-PGC/slice-chronik-in-code-
kommentar`, aktuell 2× unter Schwelle) bei Closure fortzuschreiben ist —
das ist keine Verifikations-, sondern eine Closure-Entscheidung. Kein
Validator-Zug ausgelöst — `slice-041` ist kein MVP-Meilenstein-Slice im
Sinn von Modul 8.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
