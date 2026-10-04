# Verifikations-Report: Slice api-token-mehrfach-konfiguration ([`LH-FA-SST-012`](../../spec/lastenheft.md), [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md)) — 2026-10-04

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD- und
Entscheidungs-Konformität, Plan-vs-Code-Diff und eigene Sensor-Läufe in
frischem Kontext. Der Verifier repariert nichts und setzt keine DoD-Häkchen.

**Gegenstand:** `git diff 430cc97f HEAD` — Code `38adefd4`, Runner-Phase
`425682ec`, Handbuch/Plan/README `88209377`, Suchlauf `be0f4f9d`, Review
`845b32b8` ([`review-slice-api-token-mehrfach-konfiguration`](review-slice-api-token-mehrfach-konfiguration.md):
0 HIGH, 0 MEDIUM, F-1 LOW, F-2 bis F-4 INFO), Fix `4ac7faf4` (F-1), Review-Haken
`ba491f8c`. Plan:
[`slice-api-token-mehrfach-konfiguration`](../plan/planning/done/slice-api-token-mehrfach-konfiguration.md);
Spec-Stellen [`SPEC-035`](../../spec/pflichtenheft.md),
[`SPEC-016`](../../spec/pflichtenheft.md); Entscheidungen
[`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md),
[`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md).

**Verdikt: bestanden.** Alle drei Liefer-Punkte der DoD sind an meinen eigenen
Läufen belegt, alle gefahrenen Gates enden mit Exit 0, die vier nachgefahrenen
Mutationen sind rot. Zwei Bedingungen für die Closure (Abschnitt 6), keine
blockierende.

---

## 1. Eigene Sensor-Läufe (Exit direkt gelesen, nie durch eine Pipe, [`AGENTS.md`](../../AGENTS.md) §3.9)

Reihenfolge am Stand `ba491f8c`, Arbeitsbaum sauber vor und nach dem Lauf
(`git status --short` leer).

| Sensor | Exit | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make image` | 0 | Image gebaut (Build-Kontext-Dateien geändert) |
| `make test` | 0 | `ok` für `…/internal/adapters/driving/grpc`, `…/driving/http`, `…/application/port/apiauth`, `…/internal/bootstrap`, `…/internal/domain/messagecode` |
| `make test-store` | 0 | `ok …/internal/bootstrap`; `db-coverage: OK — DB-Adapter-Coverage 83.03% erfuellt Schwelle 80%` |
| `make test-integration` | 0 | Zeile `API-Token-Wechsel (LH-FA-SST-012) belegt — Neustart 1 mit Singular und Liste: …; Neustart 2 ohne die Singular-Token: die alten Token endeten mit 401 bzw. Unauthenticated, die Listen-Token wurden bedient; eine Liste mit leerem Element und eine mit Leerraum beendeten den Start mit Ausgang 2, Fehlerklasse configuration und PCF-E2008 im Log ohne Token-Wert; nach der Wiederherstellung galten wieder die Token der Compose-Umgebung`; `E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`; `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 55 Bash-Zeilen`; 0 `FAIL`-Treffer im Log; danach keine `cdc-*`-Container und -Netze, Arbeitsbaum sauber |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`; `d-check: 1658 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — Coverage 82.50% erfüllt Schwelle 80%`; `commit-traceability: OK — 5 Commit(s)`; `generated-sync: OK`; `sdk-public-doc-check: keine interne Kennung`; `a-check: gesamt: 0 Befund(e)` |
| `make docs-check` | 0 | `d-check: 1658 Datei(en) geprüft, 0 Befund(e)` |
| `make fmt-check` | 0 | `337 Go-Dateien geprüft, alle formatiert` |
| `make a-check` | 0 | `gesamt: 0 Befund(e)` |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make ausgabe-kennungen-check` | 0 | `keine interne Kennung in Ausgabe-Literalen von 129 Go-Dateien und 4 Skripten` |
| `make meldungscodes-check` | 0 | `98 Codes in Tabelle und Katalog gleich` |
| `make doc-trace` | 0 | `83 Anforderung(en), 2 Waise(n).`; Zeile `LH-FA-SST-012 … E2E … ok`; Waisen `LH-FA-SST-010`, `LH-FA-SST-011` |
| `make doc-immutable RANGE=430cc97f..HEAD` | 0 | `1658 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=430cc97f..HEAD` | 0 | `1658 Datei(en) geprüft, 0 Befund(e)` |
| `make commit-traceability` | 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make kommentar-kennungen DIFF=430cc97f` | 0 | kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | 0 | `suchlauf-nachmessen: 10 Zeilen stimmen` |

Docker-Hub-Abruflimit: im Lauf nicht aufgetreten. Nicht gelaufen:
`make test-replication`, `make test-notify`, die SDK-Tiers und
`make examples-*` (nicht Gegenstand des Slice; `git diff 430cc97f HEAD --stat -- sdks examples compose.yaml .a-check.yml` ist leer).

## 2. DoD-Zeilen gegen Beleg

| DoD-Zeile | Befund |
|---|---|
| Liefer-Punkt 1 (a): ein Klassifikator | bestätigt: `internal/application/port/apiauth/apiauth.go`; `git grep -n classifyToken -- '*.go'` ohne Treffer; Suchlauf-Zeilen 6 und 7 (`diff 0`) stimmen |
| (b) Fälle je Tabellentest, auch HTTP und gRPC | bestätigt über die Mutationsfarben unten (die Fälle hängen an der Eingabe) |
| (c) beide Zugriffswege, `PCF-E2008`, leere Zeichenkette | bestätigt: Mutation 5 rot; Container-Lauf der Phase endet mit Ausgang 2 und der Zeile samt Variable und Elementnummer |
| (d) Klasse neun, Mengengleichheit | bestätigt: Mutation 4 rot an `TestZugangsdatenKlasseCodeUndTestSindMengengleich` und `TestConfigFromFileLehntZugangsdatenAb` |
| (e) Mutationsprobe vier Stellen | Instanz und Farben unten; Plan nennt Zahl der Mutationen ehrlich (acht), Verallgemeinerung als hergeleitet |
| Liefer-Punkt 2: Runner-Phase, vier Schritte, `abdeckung_declare` | bestätigt am Lauf; die Phase liest Ausgang und Log-Zeile am Container (`docker inspect`, `docker logs`), prüft, dass kein Token-Wert im Log steht, und stellt den Container ohne Override wieder her; Folge-Phasen liefen grün (Exit 0 des Gesamtlaufs) |
| Liefer-Punkt 3: Handbuch 1.95 | bestätigt: Abschnitt „Token-Wechsel in zwei Neustarts“ (drei Schritte, Wertregeln) stimmt mit Neustart 1 und 2 der Runner-Phase überein; Katalogzeile `PCF-E2008`; `make handbuch-public-doc-check` Exit 0 |
| Gate- und Lauf-Pflichten, Review | bestätigt (Abschnitt 1; Review liegt vor, Haken gesetzt) |
| Doku-Update `harness/README.md`, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | im Plan **offen** (Planner/Implementer bei Closure); kein Häkchen von mir |

## 3. Mutationen — selbst nachgefahren

Instanz: `go test` im Toolchain-Container wie das Makefile-Ziel `test`
(`-race`, `--network none`, Kopie des Baums als Read-only-Mount im Scratchpad;
Mutation per `sed … > Kopie`, kein `-i`). Die Echt-Repo-Dateien blieben
unberührt (`git status --short` leer).

| Nr. | Mutation | Gesehene Farbe |
|---|---|---|
| 2 | Wache `token == ""` aus `Classify` entfernt | rot: `TestClassify` (apiauth), `TestLeereTokenKonfigurationLaesstKeinenAufrufDurch` (HTTP), `TestStreamChangesLeereTokenKonfigurationEndetMitUnauthenticated` (gRPC) |
| 4 | `api_tokens_admin` aus `forbiddenFileCredentialKeys` gestrichen | rot: `TestZugangsdatenKlasseCodeUndTestSindMengengleich`, `TestConfigFromFileLehntZugangsdatenAb` |
| 5 | `mergeConfig` ohne `applyAPITokens` (nur Singular-Zuweisungen) | rot: `TestAPITokenListenWerdenGelesen`, `TestAPITokenListeUngueltigEndetMitConfiguration` |
| 7 | gRPC-`New` ohne `cfg.TokensReader` | rot: `TestStreamChangesListenTokenOeffnetStream` |
| 8 | Leerraum-Prüfung entfernt, Mutations-Image über `make image-mutation` / `make image-mutation-rm` | Mutations-Image mit `CDC_API_TOKENS_ADMIN="a, b"` läuft bis zum Datenbankfehler (Ausgang 1, kein `PCF-E2008`); das reguläre `:dev`-Image endet mit Ausgang 2 und `PCF-E2008` — beide Beobachtungen, die `tw_expect_start_refused` verlangt; Image danach entfernt |

Nicht nachgefahren, **übernommen** aus dem Plan: Nr. 1 (admin gewinnt), Nr. 3
(Leerraum im Unit-Test, 16 Fälle), Nr. 6 (HTTP-`New`). Verallgemeinerung auf
alle Stellen: hergeleitet.

## 4. Sicherheitsaussagen

- **Zeitkonstanz.** Gelesen: `Classify` weist das leere Token zuerst ab (kein
  Zeitsignal über Inhalt, nur über die Leere, die der Aufrufer selbst
  bestimmt), hasht das Aufruf-Token einmal und ruft `matchAny` für beide
  Mengen; `matchAny` durchläuft immer die ganze Menge und verknüpft
  `subtle.ConstantTimeCompare` auf SHA-256-Werten gleicher Länge mit `|=`,
  kein früher Abbruch, kein Auswerten der Admin-Menge nur bei Reader-Miss.
  Der Kommentar nennt die Grenze (Laufzeit hängt von der Zahl der Token) und
  führt die Zeitkonstanz als Erwartung; kein Test misst Laufzeit. Als
  *Erwartung* korrekt gekennzeichnet; ich habe keine Laufzeit gemessen. Der
  Kommentar-Befund F-2 des Reviews (INFO) bleibt dessen Sache.
- **Token-Lecks.** `git grep` nach `slog.*Token` in `internal`/`cmd` ohne
  Treffer; die Fehlertexte der Liste nennen Variable und Elementnummer, nie den
  Wert (am Binary gesehen, Abschnitt 5). Die Runner-Phase prüft das Log gegen
  alle Test-Token und fand keinen; im Gesamtlog des Runners kommt kein
  Phasen-Token vor.
- **Schichten.** `.a-check.yml` unverändert; `apiauth` unter
  `internal/application/port/apiauth`; beide `classifyToken` entfernt;
  `make a-check` Exit 0.

## 5. Rückwärtskompatibilität und Nebenwirkung

- Singular-Token: Neustart 1 der Phase belegt `e2e-reader-token` und
  `e2e-admin-token` (Singular) neben der Liste auf HTTP und gRPC; alle übrigen
  Phasen des Gesamtlaufs laufen mit der unveränderten Compose-Umgebung
  (nur Singular) grün. SDKs, Beispiele, `compose.yaml`: Diff leer.
- Ungültige Liste, selbst gefahren am `:dev`-Image (`--network none`, Variable
  `CDC_API_TOKENS_ADMIN="a, b"`): `--healthcheck`, `diagnose` und
  `register-consumer c` enden mit Ausgang **1**, der reguläre Lauf mit Ausgang
  **2**; die Zeile lautet in allen vier Fällen
  `Fehlerklasse configuration [PCF-E2008]: API-Token-Liste ungültig: CDC_API_TOKENS_ADMIN: Element 2 enthält Leerraum`.
  Das bestätigt F-3 des Reviews: Plan und Handbuch nennen die Sondermodi nicht;
  das Verhalten entspricht dem bestehenden Muster der anderen
  Konfigurationsfehler (gleiche `main.go`-Zweige). Die Handbuch-Fehlersuche
  „Container startet nicht“ nennt die Meldung und den Code; die Ausgänge der
  Sondermodi stehen dort nicht. Kein Muss, wenn die Closure es als
  gleichbleibende Eigenschaft benennt.
- Klasse neun gegen elf: der Plan benennt den Umsetzungsrückstand (Code neun,
  `SPEC-016` elf, bis der OTLP-Slice schließt) in §1 ausdrücklich;
  Handbuch und Code tragen neun.

## 6. Plan-vs-Code-Diff, Bedingungen, offene Punkte für den Planner

Abweichungen vom Plan, im Plan selbst benannt und vertretbar: `Config` führt
Singular und Liste nebeneinander statt je Klasse eine Menge (Testmigration
entfällt), `apitokens_test.go` und der Modus `probe` der Wegwerf-Clients sind
Zusätze über den Plan hinaus. Ich fand keine unbenannte Abweichung.

**Bedingungen für die Closure:**

1. **Link-Nachzug beim `git mv` nach `done/`.** Die beiden Folge-Slices
   [`slice-tls-http-grpc-server`](../plan/planning/done/slice-tls-http-grpc-server.md)
   und [`slice-otlp-metrik-export`](../plan/planning/in-progress/slice-otlp-metrik-export.md)
   verweisen auf `../in-progress/slice-api-token-mehrfach-konfiguration.md`
   (aktuell auflösbar, `make docs-check` Exit 0); beim Verschieben brechen
   sie. Der Nachzug gehört in die Closure (Frist laut
   [`AGENTS.md`](../../AGENTS.md) §3.13).
2. **§6-Ausgänge und Doku-Update** sind offen (Vorschlag unten); die Zeile
   `make test-integration` in [`harness/README.md`](../../harness/README.md)
   ist laut Diff gesetzt, das DoD-Häkchen dafür steht noch aus.

**Vorschlag §6-Ausgänge (Planner entscheidet):**

| Risiko | Vorschlag |
|---|---|
| Zeitkonstanz nicht belegt | weiter offen, Register: die Erwartung bleibt ohne Messung; Leser sind Reviewer und Verifier (gelesen, kein früher Abbruch) |
| Halbänderung Klasse neun/elf | weiter offen bis zum OTLP-Slice; Code und Test sind gleich (Mengengleichheit-Test), Handbuch und Code tragen neun, nur `SPEC-016` elf; Ausgang „eingetreten“ wäre falsch, nichts ist halb geblieben |
| Spec-Lücke Plural-Variable | weiter offen: Lerneintrag als benannte Spec-Lücke (leere Variable gleich ungesetzt; Komma im Token nicht ausdrückbar; der Singular wird nicht auf Leerraum geprüft, siehe F-3/F-4 des Reviews) — eine Klarstellung ist ein Zug des Planners/Architects |
| Zwei Konfigurationswege | entfallen: Mutation 5 rot an beiden Wegen, die Phase liest über `ConfigFromEnvAndFile` |
| Runner-Phase stört Folge-Phasen | entfallen: Gesamtlauf Exit 0, Wiederherstellung belegt, `TW_TMP` im Trap (F-1, Fix `4ac7faf4`) |
| Handbuch zieht nicht mit | entfallen: Abschnitt, Tabelle, Katalogzeile, Historienzeile 1.95 vorhanden, `make handbuch-public-doc-check` Exit 0 |

**Register-/Lerneintrag-Vorschläge (nur Vorschläge):**

- Beobachtung
  [`ready-ist-nicht-verbunden`](../plan/planning/observations/BEO-PGC/ready-ist-nicht-verbunden/observation.md):
  **kein** Zähler-Zuwachs. Die Phase belegt ihr Negativ über den Status der
  Aufrufe (401/Unauthenticated je Token), nicht über Abwesenheit in einem
  Stream-Fenster.
- Lerneintrag-Kandidat: der Mengengleichheit-Test als Sensor gegen die
  Halbänderung einer Listen-Klasse (Plan §5).
- Beobachtung zu `inplace-textwerkzeug`: nicht aufgetreten in meinem Lauf.
- Release-Folge: keine; kein Release ohne Freigabe.

**Folge-Slices:** beide Dateien liegen in `open/`, ihre Verweise lösen auf
(`make docs-check` Exit 0).
