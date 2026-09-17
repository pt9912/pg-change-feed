# Verifikationsbericht: slice-096 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-096` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(Diskriminator, Feldmenge, Durchleitung, Fehlerform) und
[ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)
(die Folge-ADR: die Feldmengen-Paarung ist Review-Prüfpflicht, kein Sensor) —
sowie die Hard Rules `AGENTS.md` §3.7, §3.9, §3.11, §3.12, §3.13. **Nicht**
gegen den Diff als solchen (Reviewer-Aufgabe) und **nicht** gegen realen Bedarf
(Validator, nicht ausgelöst).

**Frischer Kontext.** Der Slice-Plan wurde am Stand `HEAD` vollständig gelesen
(§1–§8), dazu beide ADRs, der Review-Report, der berührte Sensor
`docs-check.md`, `config_file.go`, die zwei `config_file`-Testdateien,
`wiring.go` (Config, `changeStreamEnabled`, `Run`s Start-Zweige), das Handbuch
§5 samt Änderungshistorie, `SPEC-016` und der Register-Eintrag
`BEO-PGC/arbeit-ueberholt-stehenden-traeger`. Review-Report und Commit-Texte
waren **Kontext**, ihre Zahlen **nicht** übernommen: jede Zahl dieses Berichts
stammt aus einem hier selbst gefahrenen Lauf. Die **Findings des Reviews habe
ich einzeln nachgemessen** — F-1, F-2 und F-3 halten; F-5 war an seinem
Mess-Stand wahr und ist bereits geschlossen (V-6); F-4 ist eine Träger-INFO.

**Gegenstand.** `HEAD` = `a04eb68`, Zweig `main`, Baum sauber. Der Vorgang sind
**neun** Commits über `3132d86`: `0823e8d` (Plan), `043c9bc`/`65d01ef`/
`3c8b6ff` (Lifecycle + Ruhe-Marker), `b8fd926` (Arbeit, 4 Dateien
+359/−65), `383e924` (Review — und die Korrektur der §3-Zeile im selben
Commit), `b261dc3` (F-3), `51a6e67` (`ADR-0089` + vier Zitat-Korrekturen),
`a04eb68` (Geschichte-Zelle). Der Slice liegt in `in-progress/`; der `git mv`
nach `done/` ist **nicht** erfolgt.

**Beleg-Lage.** Jeder Gate-Exit ist **ungepiped** ermittelt und in einem
**eigenen**, abgeschlossenen Schritt aus einer separaten Datei gelesen
(`AGENTS.md` §3.9); Gate-Lauf und Auswertung waren getrennt beauftragt. Die
Mutations- und Probe-Läufe liefen auf Arbeitsbaum-Kopien **außerhalb** des
Repos (`git archive`, netzlos im gepinnten Toolchain-Container); die
Repo-Dateien wurden **nicht** angefasst. Zum Zeitpunkt der Gate-Läufe (#1–#4)
war der Baum sauber (`git status --porcelain` leer, `HEAD` = `a04eb68`); **während**
dieser Verifikation erschien eine **fremde**, ungetrackte Datei im Baum
(`docs/plan/adr/0090-beispiel-clients-volle-matrix.md`, mtime 07:35:27 Ortszeit,
`??` in `git status`) — eine **unausgefüllte Kopie der ADR-Vorlage**. Sie gehört
**nicht** zu diesem Vorgang, ich habe sie **nicht** angefasst; ihr Effekt auf das
Doku-Gate ist gemessen (V-10). Kein host-lokaler absoluter Pfad in diesem
Bericht; die Arbeitskopien stehen als `<Arbeitskopie>`.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` (gedruckte `total:`-Zeile desselben Laufs: `83.4%`) · `d-check: 790 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)` · `generated-sync: OK` (beide `.pb.go`) · `a-check: gesamt: 0 Befund(e)`; Baum danach leer |
| 2 | `make test` (netzlos, `-race`; Exit aus eigener Datei) | **0** | **35 × `ok`, 0 × `FAIL`, 0 × `DATA RACE`**, 7 × `[no test files]` |
| 3 | `make doc-commits RANGE=3132d86..a04eb68` | **0** | 790 Dateien, 0 Befunde — jeder der neun Commits trägt seine Kennung |
| 4 | `make doc-immutable RANGE=3132d86..a04eb68` | **0** | 790 Dateien, 0 Befunde (Scope-Notiz: das Ziel fährt `--disable immutable --enable vcs`; die Hälfte trägt dort das `vcs`-Modul) |
| 5 | Feldzählung an beiden Ständen (awk über die Struktur-/Funktionsgrenzen) | 0 | `Config` **15** Felder; `mergeConfig` am Parent `3c8b6ff` erreicht **10** verschiedene Felder, am Stand `a04eb68` **15** (all 15); `ConfigFromEnv` liest **13**; die Differenz Parent→heute ist **genau** `APITokenAdmin` `APITokenReader` `GRPCAddr` `HTTPAddr` `NatsURL` |
| 6 | Arbeitskopie-Kontrolle `go test ./internal/bootstrap/...` (unmutiert, netzlos) | **0** | `ok` — Bezugsgrün für alle Mutationsproben |
| 7 | **M1** `changeStreamEnabled`s Körper → `return true` | **1** | **rot**: `TestChangeStreamEnabled/beide_leer` und `TestMergeConfigOberflaechenUnterDateiAktiv/keine_Adresse_in_beiden_Quellen` |
| 8 | **M2** `Run`s HTTP-Grenze `if cfg.HTTPAddr != ""` → `!= "disabled"` | **0** | **grün** über die ganze `bootstrap`-Suite |
| 9 | **M3** Zuweisung `cfg.NatsURL = getenv(envNatsURL)` in `mergeConfig` gestrichen | **1** | rot: `TestMergeConfigOberflaechenVariablenAusEnvUnterDatei/CDC_NATS_URL` + `…UnterDateiAktiv/Adressen aus der Datei …` |
| 10 | **M4** beide Token-Zuweisungen gestrichen | **1** | rot: beide Token-Subtests + der Oberflächen-Fall |
| 11 | **M5** `HTTPAddr`/`GRPCAddr` ohne `overrideString` (nur `getenv`) | **1** | rot: beide `AdressFelderPrecedence`-Fälle + `changeStreamEnabled("", "")=false` |
| 12 | **M6** `api_token_reader` aus `forbiddenFileCredentialKeys` | **1** | rot: `TestConfigFromFileLehntZugangsdatenAb/api_token_reader` — Log zeigt die generische Meldung `field api_token_reader not found in type bootstrap.fileConfig` |
| 13 | **M7** Tippfehler-Schlüssel `http_adr` **in** die Klasse aufgenommen | **1** | rot: `TestConfigFromFileUnbekannterSchluesselOhneZugangsdatenGrund` — die Gegenprobe bindet |
| 14 | **M8** `Run`s NATS-Grenze `!= ""` → `!= "disabled"` | **0** | grün |
| 15 | **M9** `Run`s gRPC-Grenze `!= ""` → `!= "disabled"` | **0** | grün |
| 16 | Eigene Probe-Testdatei in der unmutierten Arbeitskopie (`-v`, netzlos) | **0** | sechs Klassenschlüssel → je `… trägt den Schlüssel "<key>" — Zugangsdaten bleiben env-var-exklusiv (ADR-0088 Festlegung 1, SPEC-016)`; `http_adr` → `field http_adr not found…`; Handbuch-Beispiel-YAML lädt (`HTTPAddr=":8090" GRPCAddr=":9090" Source="quelle-1" changeStream=true`); Datei **plus alle fünf** Variablen → alle fünf getragen, `HTTPAddr=":18090"` (Env schlägt Datei) |
| 17 | d-check gegen eine **mutierte Träger-Kopie** (`grpc_addr`→`grpc_addr_x` in der `SPEC-016`-Tabelle, `nats_url` aus der Handbuch-Klassenliste, `http_addr`/`grpc_addr` aus dem §5-Beispiel entfernt), Bündel-Modulsatz | **0** | 790 Dateien, **0 Befunde** — der genannte Lauf liefert mit und ohne die Paarung dasselbe |
| 18 | dieselbe mutierte Kopie mit den optionalen, nicht git-gebundenen Modulen (`citations`, `sources`, `targets`, `spans`, `codepaths`, `diagrams`, `pins`, `workflows`) | **1** | 790 Dateien, **137 Befunde** = **133 × `codepath-missing`** + **4 × `span-unclosed`**; **kein** Befund nennt `SPEC-016`, das Handbuch oder einen Feldnamen |
| 19 | d-check gegen die Kopie von `b261dc3` | **0** | **789** Dateien, 0 Befunde — reproduziert die Zahl der [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md); `a04eb68` zählt **790** (+1 Datei = die Folge-ADR selbst) |
| 20 | Stufen-Rezept der `coverage`-Stufe auf dem Parent `3c8b6ff` nachgefahren (Arbeitskopie; **nicht** `docker build --target coverage`) | **0** | gedruckte `total: (statements) **83.1%**` |
| 21 | Umfang und Zustands-Prüfungen: `git diff --name-status 3132d86..a04eb68`; `-- harness/sensors/coverage-gate.md`; `-- docs/plan/planning/observations/`; Zeilenzahl der `wiring.go`-Hunks; §2-Häkchen; §6-Ausgänge; `docs/plan/planning/reconciliation.md` | 0 | **10** Pfade (4 Code/Doku + Plan + Review + 2 ADR + Index + Roadmap); `coverage-gate.md` **0** Treffer, Register **0** Treffer; **26** geänderte `wiring.go`-Zeilen, **alle** Kommentarzeilen; **0** von 11 Häkchen gesetzt; **4 von 4** §6-Risiken auf `<…>`; `reconciliation.md` existiert **nicht** |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### Liefer-Punkt 1 — die fünf Variablen werden durchgereicht und wirken

| Kriterium (§2) | Befund |
|---|---|
| „Unter gesetzter `CDC_CONFIG_FILE` trägt `Config` die Werte aus den fünf Variablen“ | **erfüllt** — Datei + alle fünf Variablen gesetzt: die fünf Felder tragen die Env-Werte (#16); die Zählung stützt es strukturell (#5: `mergeConfig` erreicht 15 von 15, Parent 10) |
| „je aus ihrer Env-Herkunft, und **Env schlägt Datei Feld für Feld**“ | **erfüllt** — `HTTPAddr=":18090"` gegen Datei `:8090`, `GRPCAddr=":19090"` gegen `:9090` (#16); M5 rot, wenn die Datei-Basis wegfällt (#11) |
| „Je Variable ein Test“ | **erfüllt** — `TestMergeConfigOberflaechenVariablenAusEnvUnterDatei` mit **fünf** je einzeln gesetzten Variablen; `TestMergeConfigAdressFelderPrecedence` mit **vier** Vorrang-Fällen |
| „und ein rot gesehenes Gegenbeispiel (Mutation am Durchreichen)“ | **erfüllt — selbst gesehen** — **drei** eigene Mutationen (#9, #10, #11) decken **alle fünf** Variablen ab und sind je rot |
| „plus der Nachweis, dass die Oberflächen unter geladener Datei **wirklich an** sind“ | **erfüllt an den Prädikaten, mit benannter Grenze** — `changeStreamEnabled` ist **real gebunden** (M1 → Exit 1, #7); die drei `!= ""`-Grenzen in `Run` sind es **nicht** (M2/M8/M9 → je Exit 0, #8/#14/#15) und werden im Test als eigene Ausdrücke über `cfg` ausgewertet. Beide Grenzen nennt der Testkommentar (nach F-3); **keine** §2-/§5-Zeile nennt sie → **V-3** |

### Liefer-Punkt 2 — Feldmenge und Zugangsdaten-Klasse

| Kriterium (§2) | Befund |
|---|---|
| „`http_addr` und `grpc_addr` sind Datei-Felder (`host:port`, nicht credential-tragend)“ | **erfüllt** — `fileConfig` deklariert **9** `yaml`-Tags, darunter beide; `forbiddenFileCredentialKeys` führt **6** Schlüssel, keiner davon eine Adresse; die Datei-Herkunft lädt real (#16) |
| „`nats_url` und die zwei Tokens sind **env-exklusiv** und enden in `ErrConfiguration` mit einer **eigenen, den Grund nennenden Zeile** — nicht als ‚unbekannter Schlüssel‘“ | **erfüllt** — je Schlüssel gemessen (#16): `… trägt den Schlüssel "<key>" — Zugangsdaten bleiben env-var-exklusiv (ADR-0088 Festlegung 1, SPEC-016)`. Der Gegenzustand ist selbst gemessen: M6 zeigt für denselben Schlüssel ohne Klasseneintrag die generische Meldung (#12) |
| „Je Klasse ein Test“ | **erfüllt** — `TestConfigFromFileLehntZugangsdatenAb` (sechs Subtests, prüft Schlüssel **und** Grund), `TestConfigFromFileUnbekannterSchluesselOhneZugangsdatenGrund` (Gegenprobe: der Tippfehler trägt den Grund **nicht**) — die Gegenprobe ist gebunden: M7 färbt sie rot (#13) |
| Die Klasse ist **vollständig** in der Datei unzulässig | **erfüllt** — gemessen **6** Einträge (`capture_dsn`, `admin_dsn`, `reader_dsn`, `api_token_reader`, `api_token_admin`, `nats_url`), mengengleich mit `SPEC-016`, dem Handbuch §5.2 und [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) Festlegung 1 — **nicht** aber mit Festlegung 2, die „fünf“ zählt → **V-1** |

### Liefer-Punkt 3 — der Träger zieht nach

| Kriterium (§2) | Befund |
|---|---|
| „§5.2 führt die zwei neuen Felder“ | **erfüllt** — Prosa-Absatz plus YAML-Beispiel mit `http_addr: ":8090"`/`grpc_addr: ":9090"`; das abgedruckte Beispiel lädt **real** durch den Loader (#16) |
| „und die env-exklusiven Namen mit ihrer Begründung“ | **erfüllt** — alle sechs Namen genannt, mit der Form-Begründung (`host:port` kann keine Zugangsdaten tragen, `nats_url` kann Benutzer/Passwort einbetten) und der Wirkungsaussage für den Datei-Pfad |
| „und die Versionshistorie bekommt ihre Zeile“ | **erfüllt** — `Version: 1.18` und die Zeile `1.18` liegen **im selben** Commit `b8fd926`; die verkörperte Regel `BEO-PGC/handbuch-versionshistorie-uebersprungen` ist seit `slice-053` verkörpert (3×) |
| Die Feldmenge der drei Träger ist gepaart | **erfüllt für die tragende Hälfte** — `SPEC-016` und der Code führen **mengengleich 9** zulässige Schlüssel und **6** Klassenschlüssel (gemessen, sortierte Mengen); das Handbuch führt **7 von 9** zulässigen und **6 von 6** Klassenschlüsseln → **V-2** |

### Die Closure-Pflichten aus §2

| Kriterium (§2) | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#1, Exit **0**, sechs Checks) |
| Review durchgeführt, Report unter `docs/reviews/` | **erfüllt** — `review-slice-096.md` (`383e924`), 1 HIGH / 2 LOW / 2 INFO; die Findings sind einzeln nachgemessen (§7) |
| Doku-Update für `<Schnittstelle X>` falls öffentlicher Vertrag berührt | **nicht erfüllbar wie geschrieben** — die Zeile ist die Vorlagenzeile, unausgefüllt → **V-4** |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** — §7 trägt in allen Inhaltszeilen Platzhalter (#21); die Zahlen dafür liefert §5 |
| Reconciliation-Register fortgeschrieben *(`falls …`)* | **entfällt nachweislich** — `docs/plan/planning/reconciliation.md` existiert nicht (#21) |
| Beobachtungs-Register fortgeschrieben — **kein Zähler wird gesetzt** | **offen** — im Vorgang ist **keine** Registerdatei geändert (#21). Die Substanz liefert §6; `arbeit-ueberholt-stehenden-traeger` steht bei 3× und ist seit `welle-20` verkörpert — dieser Vorgang ist sein **erster Fall nach der Verkörperung** |
| Jedes Risiko aus §6 trägt einen Ausgang | **offen** — **4 von 4** stehen auf `<…>` (#21) |
| Die drei Paarungen sind getragen | **offen** — Repo **mit** Wellen-Betrieb: fällt der nächsten Welle-Closure zu (`welle-20` ist geschlossen, dieser Slice ist wellenlos) |

**Ergebnis:** Die **drei Liefer-Punkte** und die Gate-Zeile **tragen** — gemessen,
nicht gelesen. Offen sind die regulären Closure-Pflichten und zwei benannte
Kanten (V-3 in §2/§5, V-4 in §2). Der eine Befund, der eine
**Entscheidung**-Trägerin trifft, ist V-1.

---

## 3. Entscheidungs-Konformität — hält der Vorgang, was die ADRs zusagen?

| Zusage | Befund |
|---|---|
| [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) Festlegung 1 (Diskriminator) | **eingehalten** — die Klasse ist im Code **eine** Liste, Form-geprüft, vor dem Decoding; `host:port`-Felder sind ihr nicht zugeordnet (#16). Die *Zahl* der Klasse trägt Festlegung 2 falsch → **V-1** |
| Festlegung 2 (Feldmenge: `http_addr`/`grpc_addr` neu zulässig; „neun zulässige Felder“) | **eingehalten** — neun `yaml`-Tags, mengengleich mit `SPEC-016`; die zwei neuen Felder in beiden Vorrang-Richtungen geprüft (#5, #11, #16) |
| Festlegung 3 (Durchleitung der fünf; `overrideString`-Precedence für die zwei Adressfelder; env-exklusiv auf beiden Pfaden) | **eingehalten** — kein zweiter Mechanismus (`overrideString`), die drei env-exklusiven kommen auf dem Datei-Pfad direkt aus `getenv` (#5, #16); M3/M4/M5 rot (#9–#11) |
| Festlegung 4 (Fehlerform: eigene, den Grund nennende Zeile; **kein** neuer Fehlertyp) | **eingehalten** — `ErrConfiguration` bleibt der Sentinel; die sechs Meldungen gemessen (#16), der Gegenfall ebenso (#12, #13) |
| Festlegung 5 / „Was diese ADR nicht ändert“ (kein `--config`, kein neues Format, keine Zulassung von Zugangsdaten, keine Env-only-Vorgabe geändert) | **eingehalten** — nur `CDC_CONFIG_FILE` als Zugang; kein CLI-Flag im Diff; die DSN-Vorbedingungen unberührt; `compose.yaml` im Vorgang nicht angefasst (#21) |
| §Konsequenzen, Folgepflicht (Code-Zug): `fileConfig`-Erweiterung, `mergeConfig`-Durchleitung, Klassen-Generalisierung samt Fehlerzeile, Tests je Feld | **eingehalten** — alle vier Teile im Diff; „je Feld“ durch fünf Subtests und drei Mutationen gedeckt (#9–#11, #16) |
| §Konsequenzen, Folgepflicht (Träger-Zug): Handbuch §5.2 **und** seine Änderungshistorie | **eingehalten für dieses Zug** — `b8fd926`; die Spec-Hälfte war mit demselben Nachzug bereits in `3132d86` gelandet (Feldtabelle, Klassen-Form, Änderungshistorie-Zeile gemessen) → „derselbe Zug“ ist über zwei Commits desselben Vorgangs geteilt, nicht offen |
| §Konsequenzen, „Negativ mit benannter Wirkung“: die Verhaltensänderung gehört in die Closure-Notiz | **offen** — §7 leer (#21); die Substanz (Oberflächen unter Datei **an**) ist gemessen (#16) und im Handbuch §5.2 beschrieben |
| Fitness Function, Zeile 1/2 (`go test`, kein Gate) | **eingehalten** — beide Zeilen sind real gefahren: Durchleitung/Precedence (#9–#11, #16), Zugangsdaten-Klasse (#12, #13, #16); `make test` **Exit 0**, 35 × ok, 0 × `DATA RACE` (#2) |
| Fitness Function, Zeile 3 (`make docs-check` als Träger der Feldmengen-Paarung) | **superseded durch [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)** — die Zeile steht unverändert (§3.5); die Supersession ist gedeckt: der genannte Lauf ist mit **und** ohne die Paarung identisch (#17), auch mit allen nicht git-gebundenen optionalen Modulen nennt kein Befund die drei Träger (#18) |
| [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) Festlegung 1/2 (Ersatztext; die Code-Hälfte bleibt bei Zeile 1/2) | **eingehalten** — Ersatztext steht in §Entscheidung Festlegung 2, nennt Review **und** Verifier als Wächter und das Fehlen eines Sensors; die Code-Hälfte trägt unverändert Zeile 1/2 (gemessen §2/LP1/LP2) |
| [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) Festlegung 3 (Grenze **an** der Zeile, die gelesen wird) | **eingehalten** — die Prüfpflicht steht in der Zeile selbst und nicht nur in `docs-check.md`; meine eigene Nachzählung ist genau diese Pflicht (§5, §6) |
| [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) §Fitness Function: Make-Target **—** (kein Sensor) | **eingehalten** — beide Zeilen führen `—`, die Zeile „kein Sensor“ nennt Begründung und Wächter; **kein** Make-Target, `.d-check.yml`-Modul oder Gate wurde gebaut | 
| [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) „Was diese ADR nicht ändert“ (kein Sensor abgeschaltet, keine Schwelle, kein Code-Eingriff) | **eingehalten** — `make gates` unverändert sechs Checks (#1); `harness/mk`, `Makefile`, `Dockerfile`, `tools/` im Vorgang **nicht** berührt (#21) |
| [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) §Kontext (1)–(4): die drei Mess-Zeilen | **nachgemessen** — 789 Dateien/0 Befunde an `b261dc3` (#19); `.d-check.yml` führt `[links, anchors, ids, matrix, versions, structure, hostpaths]`, kein Modul liest Go-Quelltext (gelesen); `structure`-Regeln adressieren je **eine** Datei (gelesen); „Kein Modul prüft einen Symbol- oder Funktionsnamen“ steht in §Grenze 7, und die Grenzen 8/9 existieren (#18) |
| [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) §Kontext (2)/Festlegung 2: „das Handbuch §5.2 [nennt] dieselben“ | **nicht eingehalten** — gemessen nennt §5.2 sieben der neun zulässigen Schlüssel (die zwei `wal_retention_*`-Overrides fehlen in Prosa **und** Beispiel; kein Handbuch-Treffer für diese zwei Namen) → **V-2** |
| `AGENTS.md` §3.5 (Accepted-ADRs immutable) | **eingehalten** — [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) wurde **nicht** in-place inhaltlich geändert: §Entscheidung, §Konsequenzen, §Verglichene Alternativen, §Status und die `Supersedes`-Kette sind unberührt; die vier Zitat-Korrekturen sind Form (Symbolname/Lokatoren) mit einer §Geschichte-Zeile, die F-1-Korrektur ist eine Folge-ADR |
| `AGENTS.md` §3.6 (Schwellen nur per ADR) | **eingehalten** — `THRESHOLD ?= 80` unberührt (#21) |
| `AGENTS.md` §3.7 (Ist-Zustand, keine Chronik) in den neuen Sätzen | **eingehalten** — die neuen Kommentare beschreiben den geltenden Zustand („unter geladener Datei“, „env-var-exklusiv“), ohne Vorher/Nachher-Wendung; `wiring.go`s 26 geänderte Zeilen sind **alle** Kommentar (gemessen) und tragen den Träger-Nachzug |
| `AGENTS.md` §3.11 / §3.2 | **eingehalten** — `make docs-check` grün über den ganzen Vorgang (#1, #3, #4); **0** `//nolint`-Treffer in den geänderten Tests; dieser Bericht ohne host-lokalen Pfad |
| `AGENTS.md` §3.12 in den geänderten Trägern | **eingehalten** — die neuen Zahlen des Handbuchs sind qualitativ („neun“/„sechs“-frei formuliert), die ADRs nennen ihre Läufe; die eine driftende Zahl steht in [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) Festlegung 2 → **V-1** |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen; Found und Not-found berichten) | **teilweise** — die vier Zitat-Korrekturen sind vollzogen und **alle auflösend** (§6), aber der Lauf war unvollständig und sein Bericht hat keinen Repo-Träger → **V-5** |

### Ist F-1 wirklich geschlossen?

**Ja — in der Substanz selbst nachgemessen, nicht übernommen.** Der Befund war:
die dritte Fitness-Function-Zeile nennt `make docs-check` als Träger der
Feldmengen-Paarung. Gemessen: der Bündel-Modulsatz führt **kein** Modul, das
Go-Quelltext liest, und keine `structure`-Regel stellt zwei Dokumente einander
gegenüber (#18, gelesen); eine Kopie mit **mutierten** Feldnamen in `SPEC-016`
und im Handbuch liefert **0 Befunde** (#17) — derselbe Lauf, mit und ohne die
geprüfte Eigenschaft. Die Zeile selbst **steht unverändert** — das ist richtig:
`ADR-0073` §Entscheidung 1 nimmt die Fitness-Function-Regeln ausdrücklich von
der Zitat-Korrektur aus, also Folge-ADR. Die Folge-ADR ist `Accepted`, nennt
die ersetzte Klausel wörtlich, trägt den Ersatztext, führt in ihrer eigenen
§Fitness Function das Make-Target `—` und steht mit `(Supers. ADR-0088, teilw.)`
im ADR-Index. **Kein In-Place-Rückzeiger in [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)**
ist die Repo-Norm, nicht eine Lücke: gemessen führen **9** teil-superseded ADRs
und **keiner** einen Nachfolger-Rückzeiger (der Index trägt ihn je Zeile).
**Rest:** der Ersatztext ist für zwei der neun Namen zu weit (**V-2**) — dieselbe
Klasse an derselben Stelle, eine Nummer kleiner.

---

## 4. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `HEAD`; die Lieferung ist
`b8fd926` (+ `b261dc3` als Kommentar-Nachzug aus F-3).

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `internal/bootstrap/config_file.go` — update | zwei `yaml`-Felder, Durchleitung der fünf, generalisierte Klasse samt Fehlerzeile | **Plan eingehalten** |
| `internal/bootstrap/config_file_internal_test.go` — Test neu/update | +265 Zeilen in `b8fd926`, weitere 14/−13 in `b261dc3` (F-3) | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` §5.2 — update | zwei neue Felder, Klassensprache auf sechs Schlüssel, Versionshistorie `1.18` | **Plan eingehalten** |
| `harness/sensors/coverage-gate.md` — **nicht** | **0** Treffer im Vorgang (#21) | **Plan eingehalten** |
| `spec/pflichtenheft.md` — **nicht** („`SPEC-016` ist mit `ADR-0088` bereits nachgezogen“) | nicht angefasst; die Feldtabelle, die Klassen-Form und die Zeile `2026-09-17` in der Änderungshistorie stehen aus `3132d86` | **Plan eingehalten** — die Begründung hält der Messung stand (§5) |
| — | `internal/bootstrap/wiring.go` (26 geänderte Zeilen, **alle** Kommentar) | **kein Plan-Bruch, aber unlisted**: §3 nennt `wiring.go` nicht; die Änderung ist Träger-Nachzug (`ADR-0088`-Bezüge in den `Config`-Feldkommentaren), kein Verhalten |
| — | `docs/plan/adr/0088…` (M), `0089…` (A), `…/adr/README.md` (M) | **kein Plan-Bruch**: Architect-Zug nach F-1/F-2 (Modul 8 §Konflikt-Pfad, Verdikt 1) |
| — | `docs/reviews/review-slice-096.md` (A) | **kein Plan-Bruch**: Reviewer-Übergabe-Artefakt |
| — | `docs/plan/planning/in-progress/roadmap.md` (M) | **kein Plan-Bruch**: Planner-Lifecycle (`3c8b6ff` entfernt den Ruhe-Marker, weil `in-progress/` den Slice trägt) |
| §3-Kandidatenliste selbst | `383e924` berichtigt die Testdatei-Zeile (`config_file_test.go` → `config_file_internal_test.go`) **im Review-Commit** | **kein Liefer-Defekt** — F-5 war an seinem Stand wahr und ist geschlossen → **V-6** |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts. Der Plan nennt
zwei „**nicht**“-Zeilen (Sensor-Doku, Pflichtenheft) und beide sind pfadmäßig
unberührt (#21). **Was der Diff enthält, obwohl der Plan es nicht nennt:** die
`wiring.go`-Kommentarzeilen — dieselbe Klasse wie F-5, hier **nicht** berichtigt.

---

## 5. Die Zahlen — selbst gemessen (`AGENTS.md` §3.12), und was in §7 gehört

| Zahl | Wo sie steht | Mein Lauf | Befund |
|---|---|---|---|
| **15 / 10 / 5** | Commit `b8fd926` („15 Config-Felder, am Parent 10, jetzt 15 von 15“), Review #1 | **15** Felder, Parent `mergeConfig` **10**, heute **15**; `ConfigFromEnv` **13**; die Differenz ist **genau** die fünf (#5) | **hält punktgenau** |
| **83,40 %** (vorher **83,10**) | Commit `b8fd926`; Review-Tabelle | gedruckte Stufe: `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` (#1); Parent-Rezept: `total: (statements) 83.1%` (#20) | **hält** — beide Zahlen sind **gedruckte Zeilen eines konkreten Laufs**; der Parent-Wert ist **nachgefahren**, nicht als `docker build --target coverage` |
| Feldmengen-Paarung, **9 zulässige** | [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) Festlegung 2, `SPEC-016` | **9** `yaml`-Tags ↔ **9** `SPEC-016`-Tabellenzeilen, mengengleich (sortierte Mengen) | **hält** (`SPEC-016` ↔ Code) |
| Klasse, **6 Schlüssel** | `SPEC-016`, Handbuch, Code, [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) | **6** in allen vier Trägern (#16, gelesen) | **hält** — außer in [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) Festlegung 2 („fünf“) → **V-1** |
| Handbuch-Feldmenge | [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) Kontext (2): „dieselben“ | **7** der neun zulässigen Schlüssel (Beispiel **und** Prosa); **6 von 6** Klassenschlüsseln | **hält nur für die Klasse** → **V-2** |
| **789 / 0** an `b261dc3` | [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) Kontext (1) | **789** Dateien, **0** Befunde am Stand `b261dc3`; **790** an `a04eb68` (#19) | **hält** — die Abweichung von 1 Datei ist die Folge-ADR selbst |
| **137 / 133 / 4** (mutierte Kopie, alle optionalen Module) | [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) Kontext (1) | **137** = **133 × `codepath-missing`** + **4 × `span-unclosed`**, Exit 1; **0** Treffer auf die drei Träger oder einen Feldnamen (#18) | **hält, Zahl für Zahl** — Scope: die git-gebundenen Module (`tracked`/`immutable`/`vcs`) brachen in meiner Kopie ohne `.git` mit Exit 2 ab; die 137 ist mit den übrigen optionalen Modulen reproduziert |

**Was in §7 gehört** (`AGENTS.md` §3.12, Instanz A/§5 des Plans):

1. **Zustandsgrößen** (hängen am Code-Stand): `Config` **15** Felder; `fileConfig`
   **9** zulässige `yaml`-Schlüssel; die Klasse **6** Schlüssel.
2. **Der Zuwachs als Messwert mit seinem Lauf:** `mergeConfig` erreichte am
   Parent `3c8b6ff` **10** Felder, am Stand `a04eb68` **15** — **+5**, und die
   fünf sind namentlich genau die der [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) Festlegung 3.
3. **Die erreichte Quote mit ihrem Lauf:** gedruckt **83.40 %** (`make gates`,
   Lauf `slice-096`, Stand `a04eb68`); Parent **83.1 %** (Stufen-Rezept auf
   `3c8b6ff`, nachgefahren) — als **gemessen**, der Zuwachs als **abgeleitet**.
4. **Die Verhaltensänderung als Satz** (die [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
   §Konsequenzen ausdrücklich in die Closure-Notiz verlangt): unter geladener
   Datei tragen die fünf Felder ihre Env-Herkunft, `changeStreamEnabled` steht
   auf `true` — gemessen für den Stand `a04eb68` (#16).
5. **Die benannten Grenzen** (V-3): gebunden ist `changeStreamEnabled`s Körper,
   **nicht** `Run`s drei Start-Zweige; kein Lauf dieses Repos setzt
   `CDC_CONFIG_FILE` (gemessen: `compose.yaml`, `tools/`, `test/`, `harness/`
   ohne Treffer, #21).

---

## 6. Der Suchlauf nach der bewegten Eigenschaft (`AGENTS.md` §3.13) — eigener Lauf

Die bewegte Eigenschaft dieses Vorgangs ist: **welche Variablen unter gesetzter
Konfigurationsdatei wirken**. Eigener `grep` über die Träger, beide Stände
gelesen:

**Gefunden — Träger, die die Eigenschaft beschreiben, alle jetzt richtig:**

| Träger | Befund |
|---|---|
| `spec/pflichtenheft.md` `SPEC-016` | Feldtabelle **9** Schlüssel, Klassen-Form **6**, Nachzug-Zeile `2026-09-17` — deckungsgleich mit dem Code |
| `docs/user/benutzerhandbuch.md` §5.1/§5.2 | §5.1 führt die drei Oberflächen-Variablen samt Aktivierungs-Semantik; §5.2 nennt die zwei Datei-Felder und die sechs env-exklusiven (7 von 9 zulässigen → V-2) |
| [ADR-0087](../plan/adr/0087-beispiel-clients-csharp-kotlin.md) | „Beispiele lesen `CDC_CONFIG_FILE` **nicht**“ — bleibt wahr, unberührt |
| [ADR-0052](../plan/adr/0052-optionale-yaml-konfigurationsdatei.md) | die Feld-Aufzählung ist per [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) §Status abgelöst und steht als Messung des damaligen Stands — **nicht** zu ändern |
| [ADR-0055](../plan/adr/0055-nats-change-notification-wecksignal.md)/[ADR-0057](../plan/adr/0057-http-grpc-api.md)/[ADR-0060](../plan/adr/0060-grpc-streaming-mechanismus.md)/[ADR-0061](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) | die „additiv, No-Op bei fehlender Adresse“-Aussagen gelten für den Env-only-Betrieb unverändert; die Datei-Herkunft fügt hinzu, statt zu ersetzen |
| `compose.yaml`, `tools/`, `test/`, `harness/` | **kein** Treffer für `CDC_CONFIG_FILE` — die Betriebsform ist unberührt (bestätigt die Grenze im Testkommentar) |

**Nicht gefunden — und das ist die Antwort:** kein weiterer Träger außerhalb des
Diffs beschreibt die bewegte Eigenschaft falsch. Insbesondere **kein** Träger
behauptet, die drei Oberflächen-Variablen seien unter einer Datei unwirksam.

**Und die Vollständigkeit des gemeldeten Laufs (Auftrag aus dem Register):**
Der Implementer-Lauf hat **einen** fremden Fund gemeldet — den alten
Symbolnamen `forbiddenFileDSNKeys` in [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(richtig, fremder Träger gemeldet statt still geändert). Gemessen am Stand
`b261dc3` waren in **derselben** Datei **drei** Anker gebrochen: derselbe
Symbolname **an zwei Stellen** und **zwei Zeilen-Lokatoren**
(`wiring.go:276-280`, `config_file.go:148-193` — beide am Parent `3c8b6ff`
**exakt**, gemessen: die fünf `cfg.X = getenv(...)`-Zeilen bzw. `func mergeConfig`
bis `}`). Die zwei Lokatoren fand erst der **Review** (F-2), nicht der Lauf —
der Suchlauf greift Symbolnamen, Zahlen fallen durch. **Der Bericht des Laufs
hat keinen Repo-Träger**: kein Commit des Bereichs trägt ein
Gefunden/Nicht-gefunden-Ergebnis; rekonstruierbar ist der Lauf nur aus seinen
Folgen. Alle **vier** Korrektur-Stellen lösen heute auf (§3.12-Prüfung, #21):
`forbiddenFileCredentialKeys` existiert, `wiring.go:282-286` = die fünf
Zuweisungen, `config_file.go:167-222` = `mergeConfig`. Eine **dritte**
Fundstelle des alten Symbolnamens bleibt bewusst stehen — als Zeiger in der
`Verweis`-Spalte der `Accepted`-Zeile ([ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
§Geschichte); die Begründung des Reviews (Parent-Stand als Messung) trägt für
die **Ereignis**-Spalte, für die Zeiger-Spalte ist sie eine benannte Kante
(V-7).

---

## 7. Findings

### V-1 — [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) Festlegung 2 zählt die Zugangsdaten-Klasse als „**fünf**“; gemessen sind es **sechs**

- `kategorie`: **MEDIUM**
- `quelle`: `AGENTS.md` §3.12 Instanz A (Zahl in einem Träger) ·
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- `pfad`: `docs/plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md`,
  §Entscheidung — Festlegung 2, Satz „… dazu die **fünf** env-exklusiven
  Zugangsdaten-Schlüssel“
- `befund`: Der Satz sagt von sich „die Liste ist danach **vollständig**“ und
  zählt die Klasse **eine zu klein**. Gemessen: Festlegung **1 desselben
  Dokuments** enumeriert sechs (drei DSN, zwei Token, `nats_url`); `SPEC-016`
  („Die Klasse umfasst die drei DSN-Schlüssel, die zwei Token-Schlüssel und
  `nats_url`“), das Handbuch §5.2 (sechs Namen), `forbiddenFileCredentialKeys`
  (**6** Einträge, gemessen) und [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)
  Kontext (2) („**sechs** env-exklusive“) sagen sechs. Keine Lesart trägt die
  Fünf: die Klasse ist 3 + 2 + 1 = **6**; die **Fünf** der Festlegung 3 ist die
  Durchreichungs-Menge (`HTTPAddr`, `GRPCAddr`, `NatsURL`, beide Token) — ein
  **anderes** Feld. Die Zahlen der Nachbarsätze („neun zulässige“, „7 Felder,
  3 Schlüssel“ als Parent-Messung) sind richtig; nur dieser eine zählt falsch.
- `verifizierbar`: **ja** — die vier Träger und die Liste im Code
  (`forbiddenFileCredentialKeys`, Länge **6**, in der eigenen Probe ausgegeben)
- `urteil`: **Der ausgelieferte Stand ist richtig** (6) — der Defekt sitzt
  allein im Entscheidungstext, und kein Sensor sieht ihn (`make docs-check` grün
  über den ganzen Vorgang, gemessen). Wirkung: wer die Klasse **erweitert**
  (Re-Evaluierungs-Trigger 1 der ADR) oder auditiert und die Zahl nimmt, nimmt
  eine zu kleine; die Zusage „vollständig“ deckt dann fünf statt sechs ab.
  **Weg:** Architect-Zug (Modul 8 §Konflikt-Pfad, Verdikt 1) — §Entscheidung ist
  von der Zitat-Korrektur **ausgenommen** ([ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1), also **Folge-ADR mit `Supersedes`** in enger Klausel, wie
  F-1. **Nicht** Implementer-Pfad, **nicht** still in-place.

### V-2 — der Ersatztext der Feldmengen-Paarung ist für zwei der neun Namen zu weit

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.12 Instanz B (eine als Beleg gelesene Aussage nennt
  ihren Anker oder ist als *erwartet* formuliert) · Review F-1 (dieselbe Klasse,
  an derselben Zeile)
- `pfad`: `docs/plan/adr/0089-…md` §Kontext (2) („das Handbuch §5.2 dieselben“)
  und §Entscheidung Festlegung 2, Ersatztext („nennen dieselben Namen“)
- `befund`: Gemessen nennt `SPEC-016` **neun** zulässige Schlüssel, der Code
  **neun** `yaml`-Tags (mengengleich), das Handbuch §5.2 **sieben**: die zwei
  `wal_retention_warn_bytes`/`wal_retention_error_bytes` fehlen im Beispiel
  **und** in der Prosa (kein Handbuch-Treffer für diese Namen). Die **Klasse**
  nennen alle drei vollständig (6 von 6).
- `verifizierbar`: **ja** — sortierte Mengen aus `fileConfig`, der
  `SPEC-016`-Tabelle und dem §5.2-Beispiel; `grep` nach den zwei Namen im
  Handbuch
- `urteil`: Die **tragende** Hälfte der Paarung hält (`SPEC-016` ↔ Code
  mengengleich; die Klasse in allen drei), und die Übereinstimmung ist
  **nachgezählt**, nicht angenommen — die Zeile leistet also, was sie verspricht,
  für jeden Namen, den zwei Träger führen. Zu weit ist die **Formulierung**
  „dieselben“ — um zwei Namen, an genau der Zeile, die F-1s Zeile beerbt.
  Ein-Zeilen-Nachzug („`SPEC-016` und der Code führen neun; das Handbuch die
  sieben, die es führt“) oder §5.2 zieht die zwei Namen nach. **Kein**
  Liefer-Punkt hängt daran (LP3 verlangt die zwei **neuen** Felder und die
  Klasse — beide sind da).

### V-3 — §2/LP1 („wirklich an“) und §5 („wirksam nachgewiesen“) nennen ihre Grenze nicht; sie steht nur im Testkommentar

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.12 Instanz B · Review F-3 (die Kommentar-Hälfte) ·
  [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)
  §Festlegung 3 (die Prüfpflicht gehört **in** die Zeile, die gelesen wird)
- `pfad`: `slice-096-…md` §2 LP1 („der Nachweis, dass die Oberflächen unter
  geladener Datei **wirklich an** sind“), §5 („die fünf Variablen sind unter
  geladener Datei **wirksam nachgewiesen** (nicht nur ‚im Struct gesetzt‘)“)
  gegen `internal/bootstrap/config_file_internal_test.go`, Kommentar über
  `TestMergeConfigOberflaechenUnterDateiAktiv`
- `befund`: Gemessen ist `changeStreamEnabled`s Körper **real gebunden**
  (M1 → Exit **1**, #7) — und `Run`s drei Start-Grenzen (`NatsURL`, `HTTPAddr`,
  `GRPCAddr`; gemessen bei `wiring.go:607`, `:704`, `:759`) sind es **nicht**:
  M2, M8, M9 lassen die ganze `bootstrap`-Suite **grün** (#8, #14, #15). Der
  Test wertet die drei Leer-Prüfungen als **eigene** Ausdrücke über `cfg` aus.
  Der Testkommentar nennt nach `b261dc3` beide Grenzen ehrlich; **keine** Zeile
  in §2 oder §5 tut es.
- `verifizierbar`: **ja** — vier eigene Mutationsläufe, je mit eigenem Exit
- `urteil`: LP1 ist **erfüllt an den Prädikaten** — und die Substanz dieses
  Slice hält: der Defekt war „die Felder **kommen nicht an**“, und das ist
  behoben (Durchleitung dreifach rot bewiesen, #9–#11) **und** an
  `changeStreamEnabled` gebunden (#7). Zu stark ist die **Wortwahl**: „wirklich
  an“ liest sich wie ein laufender Server, „wirksam nachgewiesen“ wie ein
  Beweis über den Verzweigungspfad — beides endet an den Prädikaten. Fix ist ein
  Halbsatz in §2/§5 (oder ein Zeiger auf den Kommentar), kein Code-Nachzug; §7
  kann es führen.

### V-4 — §2 trägt eine unausgefüllte Vorlagenzeile als DoD-Kriterium

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.7 (ein Träger beschreibt, was da ist) · Vorlagen-Form
  (vendored `slice.template.md` führt die Zeile als Beispiel)
- `pfad`: `slice-096-…md` §2, Zeile „- [ ] Doku-Update für `<Schnittstelle X>`
  falls öffentlicher Vertrag berührt.“
- `befund`: Gemessen tragen **2** der **96** Slice-Pläne diese Zeile
  (`slice-096`, `slice-097`), **94** nicht. Als Kriterium ist sie nicht
  entscheidbar — sie nennt keine Schnittstelle. Der Slice **berührt** mit
  `SPEC-016` und dem Handbuch sehr wohl einen öffentlichen Vertrag; der Nachzug
  ist erfolgt (LP3 bzw. `3132d86`).
- `verifizierbar`: **ja** — `grep` über `docs/plan/planning/`; die Vorlage
- `urteil`: Ein DoD-Punkt, der weder abgehakt noch verworfen werden kann.
  **Vor** dem `git mv` streichen oder auf die zwei berührten Träger füllen. Kein
  Liefer-Defekt, keine Wirkung auf LP1–LP3.

### V-5 — der §3.13-Lauf war unvollständig, und sein Bericht hat keinen Repo-Träger

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.13 („Sein Ergebnis steht im Bericht des Slice — mit
  dem, was er fand, **und** mit dem, was er nicht fand“; „der Reviewer und der
  Verifier prüfen das **berichtete Ergebnis**, nicht die Behauptung, gesucht zu
  haben“) · Review F-2
- `pfad`: `docs/plan/adr/0088-…md` (§Bezug und §Kontext (3)/(4)) gegen die
  Commits des Bereichs
- `befund`: Am Stand `b261dc3` waren **drei** Anker gebrochen — ein Symbolname
  (an **zwei** Stellen) und **zwei** Zeilen-Lokatoren (`wiring.go:276-280`,
  `config_file.go:148-193`; beide am Parent **exakt**, gemessen). Gemeldet wurde
  der Symbolname; die zwei Lokatoren fand erst der Review. **Kein** Commit des
  Bereichs trägt ein Ergebnis des Laufs (gefunden/nicht gefunden). Die vier
  Korrektur-Stellen lösen heute auf (§6).
- `verifizierbar`: **ja** — `grep` über beide Stände, die Commit-Texte, die
  Zitat-Korrektur-Zeile in [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
  §Geschichte
- `urteil`: **Kein Liefer-Defekt** — die Substanz ist getragen, und die
  Reparatur ist vollständig (vier Stellen, alle auflösend). Die Lehre ist
  trotzdem eine: der Suchlauf greift Symbolnamen, Zahlen fallen durch (F-2), und
  ohne Repo-Träger ist sein Ergebnis für den **nächsten** Prüfer nicht
  nachvollziehbar. Der Ort, sie zu führen, ist der §7-Eintrag der Closure bzw.
  das Register (`BEO-PGC/arbeit-ueberholt-stehenden-traeger` steht bei **3×**,
  verkörpert seit `welle-20`; dieser Vorgang ist sein erster Fall **nach** der
  Verkörperung — die Zähl-Entscheidung fällt im Lese-Schritt, nicht hier).

### V-6 — INFO: F-5 war wahr an seinem Mess-Stand und ist im Review-Commit selbst geschlossen — die §3-Berichtigung ritt im Review-Commit

- `kategorie`: **INFO** · `quelle`: Plan-vs-Code-Diff · Modul 8
  §Rollen-Sequenz (Planner-Artefakt vs. Reviewer-Artefakt)
- `pfad`: `slice-096-…md` §3, Zeile „`internal/bootstrap/config_file_test.go`“
  (Stand `0823e8d`/`b8fd926`) gegen `383e924`
- `befund`: Der Plan nannte in §3 die Datei `config_file_test.go`, die es im
  Baum nicht gibt (`config_file_internal_test.go`,
  `config_file_rest_internal_test.go`); `383e924` berichtigt genau diese Zeile
  (1 Zeile) **im selben Commit wie den Review-Report**. Zum Stand `a04eb68` nennt
  §3 die existierende Datei; **nichts** ist offen.
- `urteil`: F-5 ist damit **geschlossen** und war **nicht** falsch — die
  Nachmessung an beiden Ständen bestätigt es. Benannte Kante: die
  §3-Kandidatenliste ist seither nicht mehr der Stand, aus dem gearbeitet wurde,
  und die Berichtigung ritt in einem **Review**-Commit (Präzedenz `slice-094`:
  eigener Planner-Commit `3b7f8c7`). Ohne Wirkung auf Lieferung oder DoD.

### V-7 — INFO: der alte Symbolname bleibt einmal als **Zeiger** stehen

- `kategorie`: **INFO** · `quelle`: `AGENTS.md` §3.12 · Review F-2
- `pfad`: [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
  §Geschichte, `Accepted`-Zeile, Spalte `Verweis`
- `befund`: `forbiddenFileDSNKeys` steht dort in einer Zeigerliste, deren
  Nachbareinträge (`ConfigFromEnv`, `Config`, `fileConfig`, `mergeConfig`) heute
  alle existieren; die **Ereignis**-Spalte derselben Zeile zitiert den
  Parent-Stand als Messung (dort ist der alte Name richtig). Der Review hat das
  ausdrücklich entschieden („nicht korrigieren“).
- `urteil`: **Benannte Kante, kein Defekt** — der Zitat-Korrektur-Lauf hat die
  Zeiger-Spalte nicht als `Parent-Stand`-Zitat gekennzeichnet, sodass ein Leser
  des Verweises ins Leere greift. Wenn sie geführt werden soll, ist es eine
  Halbzeile; sonst bleibt sie wie die Parent-Messungen daneben stehen.

### V-8 — INFO: der Commit-Text zu `51a6e67` rechnet „drei Zeilen-Lokatoren“; gemessen sind es zwei

- `kategorie`: **INFO** · `quelle`: `AGENTS.md` §3.12 (Geltungsbereich:
  Git-Historie liegt außerhalb — `ADR-0083`)
- `befund`: Commit `51a6e67` sagt „VIER Anker … und DREI Zeilen-Lokatoren“, die
  §Geschichte-Zeile sagt „vier gebrochene Anker“ und zählt vier **Stellen** auf.
  Gemessen an `b261dc3`: **zwei** Lokatoren + **ein** Symbolname an **zwei**
  Stellen = drei Anker, vier Stellen.
- `urteil`: Kein Träger-Fehler (Git-Historie), aber es hängt an der
  Vollständigkeits-Frage von V-5: die Zahl ist ein Nebenprodukt des Laufs, nicht
  sein Bericht.

### V-9 — INFO: [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) führt **drei** (nicht zwei) veraltete Lokatoren

- `kategorie`: **INFO** · `quelle`: Review F-2 (Begründung der LOW-Stufe)
- `befund`: Gemessen: `wiring.go:981` liegt heute auf einer Kommentarzeile
  (`runWALRetentionCheck` bei **987**), `wiring.go:414` auf einer Kommentarzeile
  (`Run` bei **420**), `store.go:35` auf `func New` (`pool.Ping` bei **41**).
  Der Review schreibt „`ADR-0082` führt heute zwei ebenso veraltete Lokatoren“.
- `urteil`: Die Klasse ist **stärker** getragen als der Review gemessen hat —
  F-2s LOW-Stufe bleibt begründet, eher stärker. Kein Gegenstand dieses Slice.

### V-10 — INFO: eine fremde, unausgefüllte ADR-Vorlagenkopie im Baum macht `make docs-check` rot

- `kategorie`: **INFO** · `quelle`: `AGENTS.md` §3.12 (Beleg), §4 (Gate-Lage) ·
  `.d-check.yml` (`ignore: **/*.template.md`)
- `pfad`: `docs/plan/adr/0090-beispiel-clients-volle-matrix.md` (**ungetrackt**,
  nicht Teil dieses Vorgangs)
- `befund`: Die Datei ist eine **unausgefüllte** Kopie der ADR-Vorlage
  (`# ADR-NNNN: <Titel der Entscheidung>`, Template-Hinweis, `<Platzhalter>`).
  Sie liegt unter ihrem **Zielnamen**, und der Ignore des Gates greift nur auf
  `**/*.template.md` — **isoliert gemessen** (Kopie des Baums mit dieser Datei,
  Bündel-Modulsatz): **Exit 1**, Befunde u. a. `anchor-missing` auf `#<anker>`
  und `target-missing` auf `../../../spec/spezifikation.md#<anker>`, Ziel `NNNN`.
- `urteil`: **Nicht Gegenstand dieses Slice, nicht von diesem Lauf erzeugt** —
  ich melde sie, statt sie anzufassen (§3.13-Disziplin auf einen fremden Träger
  angewandt). Für die Closure bedeutet sie: solange sie **unausgefüllt** im Baum
  liegt, ist `make gates` **rot** — die sechs Checks aus #1 gelten für den Stand
  **vor** ihrem Erscheinen. Wer sie anlegt, füllt sie aus oder legt sie außerhalb
  von `docs/plan/adr/` ab.

---

## 8. Negativbefunde

- **geprüft, ohne Befund: LP1 trägt, und die Durchleitung ist dreifach rot
  bewiesen.** Datei + alle fünf Variablen → alle fünf Felder getragen, Env
  schlägt Datei feldweise (`:18090`/`:19090` gegen `:8090`/`:9090`) (#16); M3,
  M4, M5 je Exit **1** (#9–#11) und decken alle fünf Variablen ab.
- **geprüft, ohne Befund: LP2 trägt, beide Klassen und die Gegenprobe.** Sechs
  Klassenschlüssel → je die den Grund nennende Zeile; der Tippfehler-Schlüssel →
  die generische Meldung, **ohne** die Klassen-Begründung (#16); M6 und M7 zeigen
  beide Richtungen der Bindung (#12, #13).
- **geprüft, ohne Befund: F-3 hält in beiden Hälften** — `changeStreamEnabled`
  gebunden (M1 rot), `Run`s drei Grenzen nicht gebunden (M2/M8/M9 grün). Der
  Kommentar in `b261dc3` sagt genau das; die Aussage ist **nicht** über den
  Beleg hinaus.
- **geprüft, ohne Befund: das Handbuch-Beispiel ist kein Dekor.** Der
  abgedruckte YAML-Block lädt real durch `ConfigFromEnvAndFile` →
  `HTTPAddr=":8090"`, `GRPCAddr=":9090"`, `Source="quelle-1"`,
  `changeStreamEnabled=true` (#16).
- **geprüft, ohne Befund: LP3 hält exakt.** `Version: 1.18` **und** die
  Versionshistorie-Zeile `1.18` im selben Commit `b8fd926`; die zwei neuen Felder
  in Prosa **und** Beispiel; die sechs env-exklusiven Namen mit
  Form-Begründung; die verkörperte Regel `handbuch-versionshistorie-uebersprungen`
  ist seit `slice-053` verkörpert.
- **geprüft, ohne Befund: F-1s Substanz** — der Bündel-Modulsatz liest keinen
  Go-Quelltext, keine `structure`-Regel stellt zwei Dokumente gegenüber (#17,
  #18); die Supersession in
  [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)
  ist damit **begründet, nicht behauptet**.
- **geprüft, ohne Befund: die vier Zitat-Korrekturen lösen auf**, und die alten
  Lokatoren waren am Parent **exakt** (§6) — die Korrektur ist Form bei
  unverändertem Referenten (`ADR-0073`).
- **geprüft, ohne Befund: kein Produktionsverhalten außerhalb
  `config_file.go`s.** `wiring.go` trägt **26** geänderte Zeilen, **alle**
  Kommentar (#21); `harness/mk`, `Makefile`, `Dockerfile`, `tools/`, `compose.yaml`
  sind nicht im Vorgang; `THRESHOLD ?= 80` unberührt.
- **geprüft, ohne Befund: die Gates tragen den ganzen Vorgang.** `make gates`
  Exit **0** mit sechs Checks (#1), `make test` Exit **0** (#2), `doc-commits`
  und `doc-immutable` über die **volle** Range je Exit **0**, 790 Dateien,
  0 Befunde (#3, #4) — nötig, weil das Standing-Gate nur die letzten **5**
  Commits liest.
- **geprüft, ohne Befund: die Repo-Norm für superseded ADRs** — **9**
  teil-superseded ADRs, **keiner** mit In-Place-Rückzeiger auf seinen Nachfolger;
  der ADR-Index trägt ihn je Zeile (`ADR-0089 | … (Supers. ADR-0088, teilw.)`).
- **geprüft, ohne Befund: der eigene §3.13-Lauf über die Träger** (§6) — gefunden
  die Trägerliste, **nicht gefunden** ein weiterer falsch beschreibender Träger.
- **geprüft, ohne Befund: `AGENTS.md` §3.2, §3.9, §3.11 in diesem Vorgang und in
  diesem Bericht.** Kein `//nolint`; alle Gate-Exits ungepiped und aus separaten
  Dateien gelesen; kein host-lokaler absoluter Pfad im Bericht; `docs-check`
  **grün mit diesem Bericht** — isoliert gemessen an einer Kopie des Stands
  `a04eb68` **plus** diesem Bericht: 791 Dateien, 0 Befunde, Exit 0. Im
  **Repo**-Baum liegt derzeit zusätzlich die fremde Vorlagenkopie (V-10), die den
  Lauf rot macht; die ist **nicht** dieser Bericht.

---

## 9. Was ich nicht prüfen konnte

- **Der §3.13-Bericht des Implementers selbst.** Er hat keinen Repo-Träger; ich
  habe den Lauf aus seinen Folgen nachgemessen (V-5) und die Substanz unabhängig
  geprüft (§6).
- **Die Mutationen des Implementers** („acht“). Nicht einsehbar; ich habe
  **neun** eigene gefahren (#7–#15) — kein Gegenbeispiel hat sich gezeigt.
- **Der volle optionale `d-check`-Modulsatz inklusive der git-gebundenen
  Module.** Meine Arbeitskopie trägt kein `.git`; `tracked`/`immutable`/`vcs`
  brechen dort mit Exit 2 ab. Die Zahl **137** ist mit den übrigen optionalen
  Modulen reproduziert (#18) — dieselbe Zahl, dieselben Klassen.
- **Der Parent-Lauf als `docker build --target coverage`.** Ich habe das
  Stufen-Rezept in einem Container derselben gepinnten Toolchain nachgefahren
  (#20) — die Zahl ist **gemessen, nicht übernommen**, aber die Form des Laufs
  ist eine andere.
- **`make test-store`, `-replication`, `-notify`, `-integration`, `make image`,
  `make bench`.** Kein Gate und kein Bestandteil der Lieferung; kein
  Container-Vertrag berührt (`compose.yaml` ohne `CDC_CONFIG_FILE`, gemessen).
- **Der reale Post-Push-Lauf.** Der Vorgang ändert **keinen** Workflow —
  `AGENTS.md` §3.10 greift dem Buchstaben nach nicht.
- **Der Fremd-Zustand des Baums nach dem Mess-Stand.** Die ungetrackte
  ADR-Vorlagenkopie (V-10) lag beim Mess-Stand der Gate-Läufe noch nicht vor; ob
  und wie sie vor der Closure gefüllt, entfernt oder committet wird, ist nicht
  meine Entscheidung. Der Effekt auf `make docs-check` ist isoliert gemessen
  (V-10), der Effekt auf die übrigen fünf Checks **nicht**.
- **Die Register- und §6-Entscheidungen.** Ob `arbeit-ueberholt-stehenden-traeger`
  um einen Beleg wächst und welchen Ausgang die vier §6-Risiken bekommen, ist
  Lese-Schritt- bzw. Closure-Arbeit (Modul 6/5); §5/§6 dieses Berichts liefern
  die Messungen, nicht die Entscheidung.

---

## 10. Verdikt

**Die drei Liefer-Punkte aus §2 tragen — gemessen, nicht gelesen.**

- **LP1:** Die fünf Felder werden unter gesetzter Datei durchgereicht (#16), Env
  schlägt Datei feldweise, **drei** eigene Mutationen decken alle fünf ab und
  sind je rot (#9–#11); der stille Defekt ist **behoben** — `ConfigFromEnv`s 13
  Felder minus die 10 des Parents sind genau die fünf (#5). „Wirklich an“ gilt an
  den Prädikaten; `changeStreamEnabled` ist gebunden (#7), `Run`s drei Grenzen
  sind es nicht (#8, #14, #15) → **V-3**.
- **LP2:** Der Diskriminator ist exakt umgesetzt — **9** zulässige `yaml`-Felder
  (mengengleich mit `SPEC-016`), **6** env-exklusive Schlüssel, je mit eigener,
  den Grund nennenden Fehlerzeile (#16); die Gegenprobe bindet (#13), und der
  Gegenzustand ist selbst gemessen (#12).
- **LP3:** §5.2 führt die zwei Felder und die sechs env-exklusiven Namen mit
  Begründung, das Beispiel lädt real (#16), Version `1.18` **und** die
  Historie-Zeile im selben Commit.

**Entscheidungs-Konformität:** Diskriminator, Feldmenge, Durchleitung,
Fehlerform und alle „Was diese ADR nicht ändert“-Sätze der
[ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
sind eingehalten; die [ADR-0089](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)
ist in allen Festlegungen eingehalten (kein Sensor, Make-Target `—`,
Ersatztext vorhanden, Code-Hälfte bei Zeile 1/2). **F-1 ist geschlossen** —
die Zeile ist durch eine `Accepted`-Folge-ADR mit wörtlich benannter Klausel
ersetzt, und die Supersession ist **nachgemessen** (#17, #18). Der
Plan-vs-Code-Diff ergibt **keine unbegründete Abweichung**; die
`wiring.go`-Kommentarzeilen sind ein legitimer Träger-Nachzug, der in §3 fehlt.

**`done/`-fähig ist der Slice, sobald zwei Dinge erledigt sind:**

1. **V-1 hat einen Ausgang.** Der Entscheidungstext zählt die Klasse als „fünf“,
   gemessen sind es sechs — dieselbe Bauform wie F-1 (ein Träger der
   [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
   vom Vorgang nachgewiesen, aber nicht getragen). Wie F-1 gehört er in einen
   **Architect-Zug** (Folge-ADR mit `Supersedes` in der Klausel) — oder, wenn er
   offen bleibt, in §7/§6 mit Adresse (Kennung), nicht still. **Das ist der
   einzige Punkt dieser Liste, der eine Zusage berührt**, die über diesen Slice
   hinausreicht.
2. **Die regulären Closure-Pflichten:** §7 (7 Zeilen Platzhalter, darunter die
   Zahlen aus §5 mit Lauf und die Verhaltensänderung), die **vier §6-Ausgänge**,
   die Register-Entscheidung, die Häkchen — und **V-4** vor dem `git mv`
   (eine Vorlagenzeile mit `<Schnittstelle X>`). Dazu **V-10**: der Baum trägt
   eine fremde, unausgefüllte ADR-Vorlagenkopie, die `make docs-check` rot macht;
   sie ist vor der Closure-Forderung „`make gates` grün“ zu füllen oder zu
   entfernen — nicht von diesem Lauf erzeugt, aber im selben Baum.

**V-2, V-3 und V-5 sind Ein-Bis-Zwei-Zeilen-Nachzüge** ohne Wirkung auf
Lieferung, Gate oder DoD: eine zu weite Formulierung im Ersatztext (V-2), eine
Grenze, die §2/§5 nicht nennen (V-3), eine benannte Kante aus dem Suchlauf
(V-5). **Kein Liefer-Defekt, kein rotes Gate:** `make gates` **Exit 0**
(sechs Checks), `make test` **Exit 0** (35 × ok, 0 × `DATA RACE`),
`doc-commits`/`doc-immutable` über die volle Range je **Exit 0**.

---

**Beleg-Lage dieses Berichts:** jede Zahl stammt aus einem der Läufe in §1, je in
eigener Werkzeug-Beauftragung gefahren; die Gate-Läufe (#1–#4) und ihre
Auswertung waren **zwei** Schritte, ihre Exit-Codes wurden aus separaten Dateien
gelesen, nie durch eine Pipe (`AGENTS.md` §3.9). Die Mutations- und Probe-Läufe
liefen auf Arbeitskopien **außerhalb** des Repos (`git archive`), netzlos im
gepinnten Toolchain-Container; die Repo-Dateien wurden **nicht** angefasst. Der
Baum trug an seinem Mess-Stand der Gate-Läufe allein den Slice-Stand, danach
diesen Bericht — **kein** Commit, keine Änderung an Artefakten des Slice, an
`THRESHOLD`, an Produktcode oder an einem Träger außerhalb dieses Berichts. Die
während der Verifikation erschienene fremde ADR-Vorlagenkopie (V-10) habe ich
**nicht** angefasst und **nicht** entfernt: sie ist nicht Gegenstand dieses
Befunds, und ihr Effekt auf das Doku-Gate steht gemessen in V-10.
