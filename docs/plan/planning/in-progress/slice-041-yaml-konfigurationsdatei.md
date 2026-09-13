# Slice slice-041: YAML-Konfigurationsdatei (additiv zu Env-Vars)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist vollständig durch die
eigene DoD gedeckt (Datei-Loader, Merge-Logik, DSN-Ablehnung, Tests); kein
repo-weiter Beleg geht darüber hinaus (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht). Orthogonal zu
`welle-12` (CDC-Administration) und zur CI/CD-Arbeit (`slice-039`/`slice-040`),
blockiert keine davon.

**Bezug:**
[`LH-FA-SST-003`](../../../../spec/lastenheft.md) (CLI/Installation),
[`LH-FA-SST-002`](../../../../spec/lastenheft.md) (SQL-Zugriffe, Kontext),
[`LH-QA-SEC-001`](../../../../spec/lastenheft.md),
[`LH-QA-SEC-002`](../../../../spec/lastenheft.md) (Least-Privilege/
Secret-Trennung, bindet die DSN-Ablehnung),
[`ADR-0052`](../../../../docs/plan/adr/0052-optionale-yaml-konfigurationsdatei.md)
(bindende Architektur-Entscheidung — nur umgesetzt, nicht geändert).

**Berührte Spec-Stellen:** — (`ADR-0052`s Folgepflicht benennt
`spec/pflichtenheft.md`s konkrete Feldform als eigenständige Aufgabe dieses
Slice, siehe §1/§3 unten — bei Umsetzung entsteht dort eine neue `SPEC-<NNN>`-
Stelle, die zum Zeitpunkt der Planung noch nicht existiert).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `ADR-0052`s Entscheidung real umsetzen: `internal/bootstrap`
bekommt einen `ConfigFromFile`-Ladepfad (striktes YAML-Decoding,
`KnownFields(true)`) und eine Merge-Funktion, die `ConfigFromEnv`s Ergebnis
Feld für Feld über die Datei-Werte legt (Env-Var schlägt Datei). Die neue,
optionale Umgebungsvariable `CDC_CONFIG_FILE` trägt den Dateipfad; ist sie
leer, bleibt der heutige Env-only-Pfad (`ConfigFromEnv` unverändert)
exakt wie vorher. Trägt die Datei einen der drei DSN-Schlüssel
(`capture_dsn`/`admin_dsn`/`reader_dsn` oder gleichwertig benannt), ist
das ein `ErrConfiguration`-Fehler beim Laden. Die Tabellen-Aktivierung
(`tables`) wird in der Datei als YAML-Mapping strukturiert, nicht als die
heutige `CDC_TABLES`-Zeichenkette.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein `--config`-CLI-Flag** — `ADR-0052` Entscheidung 5 benennt das
  ausdrücklich als eigene Folgepflicht, gebunden an eine künftige
  `LH-FA-SST-003`-Install-CLI, die es heute noch nicht gibt
  (`cmd/pg-change-feed/main.go` parst Argumente ohne das `flag`-Paket).
  Der Zugriffsweg dieses Slice ist ausschließlich `CDC_CONFIG_FILE`.
- **Rückbau von `CDC_TABLES`/`parseTables`** — `ADR-0052`s Konsequenzen
  benennen ausdrücklich, dass beide Repräsentationen der
  Tabellen-Aktivierung nebeneinander bestehen bleiben, „bis (falls je
  gewünscht) eine spätere Entscheidung eine der beiden Formen zurückbaut;
  das ist hier nicht entschieden" — ein anderer, hier nicht angefragter
  Vorgang.
- **`compose.yaml`-Anpassung** — `ADR-0052` Konsequenz benennt das als
  eigene Folgepflicht „sobald der Slice sie nutzt"; dieser Slice liefert
  den Mechanismus, ohne die bestehende, weiterhin env-var-basierte
  Compose-Umgebung produktiv auf die Datei umzustellen (Bestand bleibt
  bewusst stehen — kein Deployment-Wechsel als Nebeneffekt eines
  Bootstrap-Slice).
- **`spec/pflichtenheft.md`-Verfeinerung als separater Schritt** — die
  konkrete YAML-Feldform entsteht als Teil dieses Slice (Implementer
  dokumentiert sie im Plan-Nachzug und trägt sie ins Pflichtenheft nach,
  siehe §2 DoD), nicht als eigener vorgeschalteter Slice — `ADR-0052`
  nennt das als Gegenstand „des umsetzenden Slices", also dieses hier.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `ConfigFromFile` (`internal/bootstrap`) lädt eine YAML-Datei striktes
      Decoding (`KnownFields(true)`); ein unbekannter Schlüssel liefert
      `ErrConfiguration`; einer der drei DSN-Schlüssel in der Datei liefert
      `ErrConfiguration` (`ADR-0052` Entscheidung 1/6, `LH-QA-SEC-001`/`002`).
- [x] Merge-Funktion überschreibt Datei-Werte Feld für Feld mit gesetzten
      Env-Vars (`ADR-0052` Entscheidung 2) — real gegen mindestens einen
      Fall getestet, in dem nur ein einzelnes Feld per Env-Var überschrieben
      wird, während die übrigen aus der Datei stammen.
- [x] Neue Env-Var `CDC_CONFIG_FILE` verdrahtet: leer/unbenannt → exakt der
      heutige `ConfigFromEnv`-Pfad, unverändert (`ADR-0052` Entscheidung 3) —
      ein Regressionstest bestätigt, dass ein bestehender Env-only-Aufruf
      ohne `CDC_CONFIG_FILE` identisches Verhalten zu vor diesem Slice zeigt.
- [x] Tabellen-Aktivierung in der Datei als YAML-Mapping (nicht die
      `CDC_TABLES`-Zeichenkettenform) — mit eigenem Test, der eine Datei mit
      mehreren Tabellen-Aktivierungen erfolgreich lädt.
- [x] `spec/pflichtenheft.md` trägt die neue Feldform als `SPEC-<NNN>`-
      Verfeinerung (`ADR-0052`s Folgepflicht) — Schlüsselnamen, YAML-Struktur
      der Tabellen-Aktivierung, `CDC_CONFIG_FILE`-Semantik.
- [x] `make gates` grün, `make test` grün (Whitebox-Tests in
      `internal/bootstrap`, analog zu den bestehenden Config-Tests).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `harness/README.md`/`docs/user/benutzerhandbuch.md`
      (falls vorhanden) nennt `CDC_CONFIG_FILE` und die Datei-Feldform,
      da ein öffentlicher Konfigurationsvertrag entsteht.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` | update | `ConfigFromFile`, Merge-Funktion, neue Env-Var `CDC_CONFIG_FILE`, DSN-Ablehnung |
| `internal/bootstrap/*_test.go` | update/neu | Whitebox-Tests: striktes Decoding, Feld-für-Feld-Merge, DSN-Ablehnung, Env-only-Regression |
| `go.mod`/`go.sum` | update | `gopkg.in/yaml.v3` von transitiv auf direkt |
| `spec/pflichtenheft.md` | update | neue `SPEC-<NNN>`-Verfeinerung: Datei-Feldform |
| `harness/README.md`/`docs/user/benutzerhandbuch.md` | update | `CDC_CONFIG_FILE`-Vertrag dokumentieren |

### Plan-Nachzug (Implementer, nach Umsetzung)

Regeln dieser Sektion: `implement-slice` §Pre-completion-Checkliste — jede
Implementierungsentscheidung, die über den Plan oben hinausgeht oder von
ihm abweicht, wird hier mit Begründung nachgetragen, statt nur im Diff
sichtbar zu sein.

**Dateien — Abweichung von der Tabelle oben:**

- **`internal/bootstrap/config_file.go` (neu)**, statt den Ladepfad/die
  Merge-Funktion direkt in `wiring.go` einzufügen: `wiring.go` trägt bereits
  1140 Zeilen; ein neuer, in sich abgeschlossener Belang (Datei-Loader +
  Merge) bekommt eine eigene Datei im selben Paket, `wiring.go` selbst
  bekommt nur einen aktualisierten Datei-Kommentar (Zeile 10–13, AGENTS.md
  §3.7: der Kommentar beschreibt den jetzigen Zustand, nicht mehr nur den
  Env-only-Stand).
- **`internal/bootstrap/config_file_internal_test.go` (neu)**, `package
  bootstrap` (echtes Whitebox, nicht `bootstrap_test` wie das bestehende
  `wiring_test.go`) — `ADR-0052`s Folgepflicht-Absatz nennt „Whitebox,
  package bootstrap" ausdrücklich; die drei bestehenden
  `*_internal_test.go`-Dateien im Paket folgen derselben Konvention.
- **`cmd/pg-change-feed/main.go` (update, nicht in der ursprünglichen
  Tabelle)**: alle fünf Aufrufstellen (`--healthcheck`,
  `register-consumer`, `acknowledge-consumer`, `diagnose`, Hauptlauf) rufen
  jetzt `bootstrap.ConfigFromEnvAndFile(os.Getenv)` statt
  `bootstrap.ConfigFromEnv(os.Getenv)`. **Begründung:** `CDC_CONFIG_FILE`
  ist laut ADR ein Umgebungsvariablen-Zugriffsweg wie die DSNs — ohne
  diese Rewiring bliebe die neue Variable im tatsächlichen Binary
  wirkungslos (nur über die Paket-API testbar, nie über einen echten
  Container-Start erreichbar), und DoD-Punkt 3 („Neue Env-Var
  `CDC_CONFIG_FILE` **verdrahtet**") wäre für das Binary nicht wahr. Keiner
  der vier Out-of-Scope-Punkte aus §1 schließt das aus — ausgeschlossen ist
  ein `--config`-CLI-Flag, nicht das Einlesen der bereits bestehenden
  Env-Var-Vorbedingung über einen neuen Namen. `ConfigFromEnv` selbst bleibt
  unverändert und weiterhin exportiert (Regressionstest,
  `TestConfigFromEnvAndFileLeereEnvVariable`).

**Exakte YAML-Feldnamen** (`fileConfig` in `config_file.go`, dokumentiert
in `spec/pflichtenheft.md` `SPEC-016`): `source_id`, `publication`, `slot`,
`tables` (Mapping `<schema.tabelle>: {table_id, schema_version}`),
`log_level`, `wal_retention_warn_bytes`, `wal_retention_error_bytes`. Die
drei verbotenen Schlüssel: `capture_dsn`, `admin_dsn`, `reader_dsn`.

**Struktur der Merge-Funktion:** `ConfigFromEnvAndFile(getenv)` ist der neue
Einstiegspunkt — `CDC_CONFIG_FILE` leer → delegiert unverändert an
`ConfigFromEnv(getenv)` (identisches Verhalten, siehe Regressionstest);
sonst `ConfigFromFile(path)` (striktes YAML-Decoding + expliziter
DSN-Schlüssel-Check *vor* dem strikten Decoding, eigene Fehlerzeile statt
generischem „unbekannter Schlüssel") gefolgt von `mergeConfig(file,
getenv)`: DSNs werden unverändert direkt aus der Umgebung gelesen und
geprüft (env-var-exklusiv, kein Datei-Gegenstück); `Source`/`Publication`/
`Slot` je über `overrideString(fileWert, envWert)` (env gewinnt, wenn
gesetzt) und danach identisch zu `ConfigFromEnv` auf Nicht-Leere geprüft;
`Tables` über `mergeTables` (siehe Risiko-Ausgang unten); `LogLevel` über
dieselbe Präzedenz wie die drei String-Felder, geparst mit dem
bestehenden `parseLogLevel`; `WALRetentionWarnBytes`/`WALRetentionErrorBytes`
haben kein Env-Gegenstück und kommen ausschließlich aus der Datei (0
bleibt 0, `Run`s bestehender Fallback auf die [`SPEC-013`](../../../../spec/pflichtenheft.md)-Startwerte bleibt
unverändert wirksam).

**Implementierungsentscheidung zu §6 Risiko 1 (`tables`-Merge-Semantik):**
Eine gesetzte `CDC_TABLES` schlägt die gesamte Datei-`tables`-Mapping
**vollständig** — keine Vermischung einzelner Tabellen aus beiden Quellen
innerhalb derselben Liste (`mergeTables` in `config_file.go`). Begründung:
`ADR-0052` Entscheidung 2 formuliert die Precedence als „Feld für Feld";
`tables` ist im `Config`-Typ ein einzelnes Feld (eine Map), keine Menge
unabhängig geführter Einzelfelder — dieselbe Behandlung wie bei
`Source`/`Publication`/`Slot` (ganzer Wert ersetzt, nicht anteilig
gemischt) vermeidet zusätzlich eine unklare dritte Präzedenz-Ebene
(„welche Tabelle gewinnt, wenn `public.t1` in beiden Quellen mit
unterschiedlicher `table_id` vorkommt?"), die `ADR-0052` nicht anspricht.
Ein anders lautendes künftiges Bedürfnis (additive Mischung) ist eine neue
Entscheidung, kein stiller Fortschritt dieser. Real gegen
`TestMergeConfigTabellenCDCTablesSchlaegtDatei` getestet (rot bei
Sabotage der Implementierung, siehe Bericht an den Reviewer/Verifier).

**Ort der DSN-Ablehnung:** doppelt abgesichert — ein expliziter Check der
drei verbotenen Schlüssel auf einem roh eingelesenen `map[string]any`
*vor* dem typisierten, strikten Decoding (eigene, den Secret-Grund
nennende Fehlerzeile, `ADR-0052` Entscheidung 6: „nicht stillschweigend
ignoriert") **und** implizit durch `KnownFields(true)` selbst (die drei
Schlüssel sind im `fileConfig`-Typ nicht deklariert, ein generischer
„Feld nicht gefunden"-Fehler träte auch ohne den expliziten Check ein).
Der explizite Check bleibt die tragende Instanz für die Fehlermeldung, die
implizite Deckung ist eine Verteidigungslinie unabhängig von künftigen
Struct-Änderungen.

**Reihenfolge im Fehlerpfad `ConfigFromFile`:** Datei nicht lesbar →
ungültiges YAML → verbotener DSN-Schlüssel → striktes Decoding
(unbekannter Schlüssel/Typfehler). Eine leere Datei (kein YAML-Dokument,
`io.EOF` beim Decode) ist **kein** Fehler — sie liefert
`fileConfig{}` zurück, äquivalent zu einer Datei ohne jedes Feld.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `ADR-0052` liegt `Accepted` vor,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Datei-Loader, Merge-Logik und Pflichtenheft-Verfeinerung zusammen mehr
  als drei Liefer-Punkte oder mehr als zwei Schichten in einer
  Review-Sitzung nicht mehr prüfbar machen, gehört das zurück zur
  Zerlegung (z. B. Pflichtenheft-Verfeinerung als eigener Folge-Slice).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test` grün **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Merge-Semantik „Env-Var schlägt Datei, Feld für Feld" könnte für das
  `tables`-Feld unklar werden, wenn sowohl `CDC_TABLES` als auch die
  Datei-`tables`-Mapping gleichzeitig gesetzt sind — welches „Feld" gilt
  hier als Einheit (die ganze Tabellenliste, oder je Tabelle einzeln)?
  `ADR-0052` entscheidet das nicht explizit. **Ausgang:** <bei Closure
  einzutragen>
- Ein bestehendes Env-only-Deployment (`compose.yaml`) könnte durch die
  neue `ConfigFromFile`-Codepfad-Verzweigung unbeabsichtigt einen anderen
  Fehlerpfad durchlaufen, selbst wenn `CDC_CONFIG_FILE` leer bleibt.
  **Ausgang:** <bei Closure einzutragen>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC`: `BEO-PGC/github-actions-unverifizierbar-lokal` (1×,
thematisch nicht berührt — reine CI/CD-Beobachtung),
`BEO-PGC/verwaltung-keine-sql-administration` (0×, thematisch nicht
berührt — CDC-Tabellen-Administration, nicht Tool-Prozesskonfiguration).
Keine der beiden erreicht mit diesem Slice 3×; keine ist thematisch
einschlägig.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
