# Review-Report: slice-api-token-mehrfach-konfiguration — 2026-10-04

**Review-Art:** Code — geprüft gegen Plan, ADRs und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice
[`slice-api-token-mehrfach-konfiguration`](../plan/planning/done/slice-api-token-mehrfach-konfiguration.md)
(wellenlos), Diff `430cc97f..HEAD` (`HEAD` = `be0f4f9d`, 25 Dateien, +1095/−264). Implementer-Lauf: `38adefd4` (Code),
`425682ec` (Runner-Phase, Wegwerf-Clients), `88209377` (Handbuch, Plan, `harness/README.md`, Erzeugnis
`docs/user/e2e-abdeckung.md`), `be0f4f9d` (Suchlauf-Zeilen).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-/MEDIUM-Klassen ergänzt).
**Modell:** claude-sonnet-5 · **Datum:** 2026-10-04.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Mutationen liefen an Kopien im
Scratchpad (`sed … Datei > Kopie`, Ausgabe nach stdout; nie `sed -i`, nie eine Umleitung auf eine Repo-Datei). Go lief im
gepinnten Toolchain-Image (`golang:1.27@sha256:e0174e51…`, dasselbe wie `make test`) mit
`docker run --rm --network none -v <Kopie>:/src:ro -v pg-change-feed-gomodcache:/go/pkg/mod -w /src go test …`; die Kopie
enthielt `go.mod`, `go.sum`, `internal/`, `gen/`, `cmd/`. Es wurde kein Image gebaut. `git status --short` am Echtrepo nach
den Mutationen: leer.

**Eingangs-Kontext:** Plan (§1–§3, §6), [`LH-FA-SST-012`](../../spec/lastenheft.md),
[`SPEC-035`](../../spec/pflichtenheft.md), [`SPEC-016`](../../spec/pflichtenheft.md),
[`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md), [`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), [`ADR-0088`](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md); `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13;
`harness/conventions.md`.

**Eigene Messungen (Exit je Befehl direkt gelesen):**

| Lauf | Ergebnis |
|---|---|
| `make gates` | Exit 0 (Schlusszeile `gesamt: 0 Befund(e)`, `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung`) |
| `make test` | Exit 0; `ok` für `…/apiauth` und `…/internal/bootstrap` |
| `make fmt-check` | Exit 0, „337 Go-Dateien geprüft, alle formatiert“ |
| `make handbuch-public-doc-check` | Exit 0, „keine interne Kennung in 3 Nutzerdokumenten“ |
| `make ausgabe-kennungen-check` | Exit 0, „129 Go-Dateien und 4 Skripten“ |
| `make meldungscodes-check` | Exit 0, „98 Codes in Tabelle und Katalog gleich“ |
| `make kommentar-kennungen DIFF=430cc97f` | Exit 0, kein Kandidat (Probe, kein Beleg) |
| `make suchlauf-nachmessen PLAN=…` | Exit 0, „10 Zeilen stimmen“ |
| `make a-check` | Teil von `make gates`: `gesamt: 0 Befund(e)`; `.a-check.yml` im Diff nicht berührt |
| `git grep -n classifyToken` | nur noch [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) und `ADR-0150` (Accepted, ausgenommen) und der Plan; in `internal/` 0 Treffer |
| `make test-integration`, `make image`, `make test-store` | **nicht gefahren**, Aussagen des Plans (§2 „Belege des Implementers“) **übernommen** |

**Mutationen (an Scratchpad-Kopien, jede einzeln, Instanz `go test` im Toolchain-Container):**

1. Reader vor Admin am Ende von `Classify` (die zwei `if`-Zweige vertauscht) — rot: `TestClassify`
   (`apiauth`), `TestMehrereTokenJeKlasse` (HTTP), `TestAuthInterceptorMehrereTokenJeKlasse` (gRPC). Entspricht
   Mutation 1 des Plans, **nachgefahren, gleiche Farbe**.
2. `unicode.IsSpace`-Prüfung in `parseAPITokenList` auf `if false && …` gesetzt — rot:
   `TestAPITokenListeUngueltigEndetMitConfiguration`; gezählt 17 `--- FAIL`-Zeilen (1 Eltern-Test, 16 Unterfälle). Der Plan
   nennt 16 Fälle (Mutation 3): **gleich**.

## Findings

### F-1 — Das Temp-Verzeichnis der Runner-Phase steht nicht im `cleanup`-Trap

- `kategorie`: LOW
- `quelle`: Maintainability (Muster der Nachbar-Phasen)
- `pfad`: `tools/harness/run-integration-tests.sh:4440` (`TW_TMP=$(mktemp -d)`), `:4476` (`rm -rf "$TW_TMP"` nur auf dem
  Erfolgspfad), `:256-264` (`cleanup` löscht nur `${WAL_TMP:-}`)
- `befund`: `tw_recreate` schreibt die Override-Datei nach `$TW_TMP`. Bricht `bf_fail` die Phase ab, bleibt das
  Verzeichnis liegen; der Trap entfernt `WAL_TMP`, aber nicht `TW_TMP`. Der Container-Zustand selbst ist unkritisch:
  `cleanup` fährt `compose down -v`, die Wiederherstellung am Phasen-Ende ist nur für die Folge-Phasen eines
  erfolgreichen Laufs nötig. Der Plan sagt „am Phasen-Ende ohne Override wiederhergestellt“ und meint den Erfolgspfad.
- `verifizierbar`: ja — Lesen von `cleanup` und der Phase; ein Fehlschlag der Phase hinterlässt `/tmp/tmp.*`
- `klasse`: Temp-Verzeichnis ohne Trap-Bindung

### F-2 — Der Kommentar von `Classify` formuliert die Grenze der Zeitkonstanz etwas weiter, als der Code sie trägt

- `kategorie`: INFO
- `quelle`: [`SPEC-035`](../../spec/pflichtenheft.md) (Zeitkonstanz als Erwartung), `AGENTS.md` §3.7 (Klasse Grenze)
- `pfad`: `internal/application/port/apiauth/apiauth.go:61-65`
- `befund`: Der Kommentar nennt als Grenze „die Laufzeit hängt von der Zahl der konfigurierten Token ab, nicht vom Inhalt
  des Aufruf-Tokens“. Gelesen am Code: das leere Aufruf-Token kehrt vor dem Hash zurück (schneller Pfad, Zeile 67), und
  `sha256.Sum256` läuft proportional zur Länge des Aufruf-Tokens. Beides verrät nur Leere und Länge des Werts, den der
  Aufrufer selbst sendet, kein Geheimnis; die Schlussklausel „kein Test misst sie“ ist wahr.
- `verifizierbar`: nein — Lese-Handlung
- `klasse`: Grenz-Kommentar knapper als der Code

### F-3 — Eine ungültige Token-Liste beendet auch die Sondermodi, bevor sie die Datenbank erreichen; Plan und Handbuch nennen es nicht

- `kategorie`: INFO
- `quelle`: Plan §3 (Tabelle `wiring.go`), `AGENTS.md` §3.13
- `pfad`: `cmd/pg-change-feed/main.go:45`, `:63`, `:87`, `:102` (alle fünf Modi rufen `ConfigFromEnvAndFile`),
  `internal/bootstrap/wiring.go` (`applyAPITokens` in beiden Zugriffswegen)
- `befund`: `--healthcheck`, `diagnose`, `register-consumer` und `acknowledge-consumer` enden bei einer ungültigen Liste mit
  Exit 1 und der Zeile `PCF-E2008`, der reguläre Lauf mit Exit 2 (gelesen in `main.go`). Das ist keine Regression: dieselbe
  Vorbedingungsprüfung teilen die Modi schon heute (`validateNatsStreamTokenRequiresURL`, Kommentar in `main.go:38-44`), und
  die Sondermodi laufen per `docker exec` in der Umgebung des Containers, der mit einer ungültigen Liste nicht gestartet
  ist. Sichtbar wird es nur, wenn ein Betreiber einen Sondermodus mit anderer Umgebung als der des Containers fährt. Der
  Handbuch-Schritt 3 der Fehlersuche (Token-Listen prüfen) deckt den Container-Start, nicht diese Aufrufform.
- `verifizierbar`: ja — `main.go` lesen; die Exit-Werte sind aus dem Quelltext gelesen, nicht gefahren
- `klasse`: Nebenwirkung einer geteilten Vorbedingung nicht benannt

### F-4 — Der Adapter-`Config` führt Singular und Liste nebeneinander; die Vereinigung steht an einer Stelle

- `kategorie`: INFO
- `quelle`: Plan §3 (benannte Abweichung vom Wortlaut „je Klasse die Menge“), Skill „Zwei-Quellen-Drift“
- `pfad`: `internal/adapters/driving/http/server.go:36-39,89`, `internal/adapters/driving/grpc/server.go:48-51,95`,
  `internal/application/port/apiauth/apiauth.go:41-46`
- `befund`: Je Adapter vier Felder (`TokenReader`/`TokenAdmin` + `TokensReader`/`TokensAdmin`). Beide Adapter bilden sie
  über dieselbe Funktion `apiauth.FromConfig` ab (Vereinigung, `append` auf frischem Slice, kein Aliasing); ein Wert in
  Singular und Liste gleichzeitig ist ein doppelter Digest, wirkungslos. Eine zweite Quelle der Wahrheit entsteht daraus
  nicht, weil die Auswertung an genau einer Stelle steht. Die Abweichung vom Plan-Wortlaut ist im Plan (§3) benannt und
  begründet (Migration von 99 Trefferzeilen vermieden).
- `verifizierbar`: ja — Lesen; Mutationen 6 und 7 des Plans decken die Weitergabe je Adapter (übernommen)
- `klasse`: Plan-Abweichung benannt, ohne Drift

## Negativbefunde

**Schwerpunkt 1 — Zeitkonstanz** (`apiauth.go`): gelesen, ohne Befund über F-2 hinaus. `New` hält SHA-256-Werte
(`[32]byte`), `Classify` berechnet `matchAny` für **beide** Klassen unbedingt, bevor irgendeine Verzweigung kommt
(`isReader`/`isAdmin` stehen vor den `if`s); `matchAny` akkumuliert mit `found |= subtle.ConstantTimeCompare(...)` ohne
`return`/`||`/`break` im Schleifenkörper. Kein Längen-Leak vor dem Hash außer dem leeren Token (F-2). Die Reihenfolge
„Admin zuerst“ betrifft nur die Auswertung **nach** beiden Durchläufen und ändert die Laufzeit nicht (Mutation 1 ändert
das Ergebnis, nicht die Arbeitsmenge). „Leeres Token matcht nie“ hat genau eine Wache (Zeile 67); enthält die Menge ein
leeres Element, lehnt der Parser es auf beiden Zugriffswegen ab, und der unbesetzte Singular (`""`) steckt über
`FromConfig` immer als Digest von `""` in der Menge — dort trägt allein die Wache; Tabellenfälle „leerer Wert gegen leeres
Element der Reader-/Admin-Klasse“ binden sie (Mutation 2 des Plans, Farbe übernommen). Eine Zweitwache an der
Menge-Konstruktion gibt es nicht (einfache Tiefe, nicht Defekt).

**Schwerpunkt 2 — Konfiguration:** `parseAPITokenList`/`applyAPITokens` rufen beide Zugriffswege (`wiring.go:337`,
`config_file.go:223`); gelesen. Leere Zeichenkette liefert `nil, nil`; Element unverändert (kein `TrimSpace`); der Singular
wird nie zerlegt (Fall `a,b` im Singular bleibt ein Token, Test `Singular mit Komma bleibt ein Token`). Token-Lecks:
`git grep` nach Formatierungen mit Token-Bezug in `internal/bootstrap` (`Errorf`/`Sprintf`/`slog`/`log.`) findet nur die
zwei Stellen in `parseAPITokenList` (Variablenname und Elementnummer) und die NATS-Zeile (Variablenname); kein `%v`/`%q`
auf `Config`. Der Test liest den Fehlertext gegen `geheim` (nicht enthalten); die Runner-Phase liest das Container-Log
gegen die Token-Werte. `apiTokenListError.Unwrap() []error` ist Go-1.20-Muster (`go.mod`: 1.27); `messagecode.From`
behandelt `Unwrap() []error` (`messagecode.go:196`). Bestandsmuster für Config-Fehler ist `fmt.Errorf("%w: …", ErrConfiguration)`;
der Typ ist neu, weil Text und Klasse zugleich getragen werden müssen — ohne Befund.

**Schwerpunkt 3 — Zugangsdaten-Klasse ([`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md)):** `forbiddenFileCredentialKeys` trägt neun Schlüssel (drei DSN, vier API-Token, `nats_url`,
`nats_stream_token`); [`SPEC-016`](../../spec/pflichtenheft.md) nennt elf, der Rückstand bis zum OTLP-Slice ist im Plan §1
und §6 benannt (und im Handbuch richtig als Ist-Zustand „neun“). `TestZugangsdatenKlasseCodeUndTestSindMengengleich`
vergleicht beide Richtungen, Duplikate und die Zahl 9; gelesen. Er färbt an der Mutation „Schlüssel aus dem Code“ (Plan
Mutation 4, übernommen), nicht an einer Mutation der Test-Liste, die dieselbe Menge ändert (beide Seiten werden zugleich
gepflegt, das ist die Natur der Halbänderungs-Wache).

**Schwerpunkt 4 — Klassifikator:** `classifyToken` entfernt (0 Treffer in `internal/`); HTTP (`withToken`) und gRPC
(Stream- und Unary-Interceptor, `administrationRPCRoles`) rufen `apiauth.Classifier`; unbekannter Methodenname fällt
fail-closed auf `Admin`. `.a-check.yml` unverändert, `make a-check` 0 Befunde. F-4 zum `Config`-Wachstum.

**Schwerpunkt 6 — Runner und Clients:** `probe` in `httpclient` und `grpcadminclient` druckt je Token nur Index und
Statuswerte, nie den Token-Wert; HTTP wiederholt nur Transportfehler (30 s), gRPC wartet per `WaitForReady` auf die
Verbindung und liest den Status nur einmal. Die Override-Datei liegt im Temp-Verzeichnis (kein Repo-Text per Umleitung,
§3.1 gewahrt: Heredoc nach `$TW_TMP`). Abdeckungszeile `abdeckung_declare` für `LH-FA-SST-012` vorhanden;
`docs/user/e2e-abdeckung.md` ist Erzeugnis (gleiche Ort-Zeilen verschoben, eine neue Zeile) und trägt keine Lauf-Aussage.
Wiederherstellung am Ende: auf dem Erfolgspfad vorhanden und mit `tw_expect_classes` belegt; Fehlschlag: `compose down -v`
im Trap (F-1 nur zum Temp-Verzeichnis). Docker-only: keine neue Host-Toolchain, `go run` im Toolchain-Container.

**Schwerpunkt 7 — Mutationen:** zwei von acht nachgefahren (oben), beide in der Farbe des Plans. Mutation 8 (Image) und
die Klassen-Aussage der Runner-Phase gegen ein Mutations-Image stehen im Plan als *nicht erprobt* beziehungsweise
*hergeleitet*; die Kennzeichnung ist korrekt (der Runner mountet `:dev`, `make image` mit mutiertem Baum ist verboten).
Mutation 8 selbst nicht nachgefahren: **übernommen**.

**Schwerpunkt 8 — Handbuch:** Version 1.95 und Stand hochgezählt, Änderungshistorie-Zeile vorhanden, kennungsfrei
(`make handbuch-public-doc-check` Exit 0), in Betreibersicht; Abschnitt „API-Token in zwei Neustarts wechseln“ stimmt
mit den drei Schritten der Runner-Phase überein (Neustart 1 Singular + Liste, Neustart 2 ohne Singular,
`401`/`Unauthenticated`; Leerraum-Fehler `PCF-E2008`). Katalogzeile `PCF-E2008` vorhanden, in `codes.go` Klasse
`configuration`. „Zeilenumbruch“ als Leerraum deckt `unicode.IsSpace` und der Testfall `Zeilenumbruch am Ende`. Aussagen zu
Beispielen und SDKs fehlen im Handbuch: ohne Befund (Beispiele und SDKs senden ein Token je Aufruf; `git grep -i 'ein
Token|one token|single token'` in `sdks/`, `examples/`, `README.md`: kein Träger nennt „ein Token je Klasse“ als
Servergrenze).

**Schwerpunkt 9 — Träger-Nachzug (§3.13):** Zeilen 1 bis 5 des Suchlaufs stimmen an beiden Ständen (10 Zeilen).
`git grep -n -i 'sieben'` außerhalb ausgenommener Räume: keine Aussage zur Zugangsdaten-Klasse mehr (nur Fehlerklassen und
Antragsarten). Die 15 Treffer von „beiden/zwei Token-Klassen“ im Diff-Stand sind gelesen: sie zählen Klassen, nicht
Token, und bleiben wahr (`compose.yaml:131`, Beispiele, Kommentare der Adapter, `spec/pflichtenheft.md:1023`).
Rückverweise der zwei offenen Folge-Pläne zeigen auf `in-progress/`. Der Suchlauf-Zeile 4 (`nats_stream_token` im Muster)
zählt den Wegfall von „sieben“ nicht allein; die Lesung per `git grep` oben schließt das.

**Schwerpunkt 10 — Hard Rules:** kein Docker-only-, §3.1-, §3.9-Verstoß im Diff gefunden; Kommentare in neuem Code tragen je
höchstens eine Kennung (`kommentar-kennungen` Exit 0) und beschreiben den Ist-Zustand; ein Kommentar im Konjunktiv über eine
verworfene Alternative, eine Slice-Chronik oder ein Vorher/Nachher in Produktionscode ist nicht gefunden. Traceability: die
vier Commit-Betreffe nennen `LH-FA-SST-012` und `ADR-0150`, keine `SPEC-*`-/`ARC-*`-Kennung im Betreff.

Verzeichnisse ohne Befund: `internal/adapters/driving/http`, `internal/adapters/driving/grpc`, `internal/domain/messagecode`,
`internal/application/port/apiauth`, `internal/bootstrap` (Produktionscode und Tests), `tools/harness/httpclient`,
`tools/harness/grpcadminclient`, `docs/user`, `docs/plan/planning/open` (Rückverweise).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (F-1) |
| INFO | 3 (F-2, F-3, F-4) |

## Verdikt

**Nicht merge-blockierend.** 0 HIGH, 0 MEDIUM. Der Vergleich ist zeitkonstant im Sinn der Erwartung von
[`SPEC-035`](../../spec/pflichtenheft.md) (alle Token beider Klassen ohne Abbruch, SHA-256-Werte gleicher Länge); die
Zeitkonstanz selbst bleibt Erwartung (Lesung, kein Test). F-1 ist ein Wunsch an den Runner (Temp-Verzeichnis in den
Trap), F-3 eine Nennung im Handbuch oder im Plan-Lerneintrag, F-2 und F-4 ohne erwartete Aktion.

**Fixrunde:** nicht nötig. Die DoD-Zeile „Review durchgeführt“ des Plans habe ich auf Anweisung des Auftrags **nicht**
gezogen (Skill §DoD-Checkbox-Nachzug sähe es vor); der Planner zieht sie mit diesem Report nach.

**Übergabe:** F-1 an den Implementer (optional, nächster Runner-Zug); F-3 an den Planner (Satz in Handbuch oder
Closure-Notiz). An den Verifier, ungeprüft von mir: `make test-integration` (volle Phase, Zeile
`API-Token-Wechsel (LH-FA-SST-012) belegt`), `make image`, `make test-store`, Mutation 8 (Image) — alles **übernommen**.
Der Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
