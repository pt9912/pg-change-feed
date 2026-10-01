# Review-Report: slice-routing-antragsweg, Fixrunde 1 — 2026-10-01

**Review-Art:** Code, Re-Review einer Fixrunde (Modul 10 §Drei Review-Arten). Geprüft
wird die Fixrunde gegen die Findings des Hauptreports
[`review-slice-routing-antragsweg.md`](review-slice-routing-antragsweg.md) (F-1 bis F-10,
A-1, A-2) und gegen die Befunde V-1 bis V-7 des Verifikations-Reports
[`verifikation-slice-routing-antragsweg.md`](verifikation-slice-routing-antragsweg.md).
Kein DoD-Abgleich (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff cc653a85~1 cc653a85` (14 Dateien, +236/−60): Parser
`jsonRouteOrder` in `internal/domain/model/routespec.go`, Tests in vier Paketen,
Dekorator-Test in `internal/bootstrap/assemblersync_internal_test.go`, Kommentare in
`queries.go`, `administrationrequest.go`, `wiring.go`, Slice-Plan und die Pläne
`slice-routing-betriebsdoku` und `slice-routing-e2e`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt). **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Eingangs-Kontext:** Hauptreport und Verifikations-Report (siehe oben);
[`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md);
[`SPEC-032`](../../spec/pflichtenheft.md) (Regelform, `order`: „JSON-Zahl, positive ganze
Zahl“) und [`SPEC-019`](../../spec/pflichtenheft.md) (Fehlertext-Tabelle);
[`LH-FA-CFG-008`](../../spec/lastenheft.md); `AGENTS.md` (§3.7, §3.12, §3.13).

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus Bericht oder Verifikation
übernommen):

- Probe-Test auf einer Kopie im Scratchpad (`git archive HEAD`, gepinntes
  `TOOLCHAIN_RACE_IMAGE`, `go test -race`): `jsonRouteOrder` gegen 25 Eingaben, jede wie
  erwartet. Angenommen: `10`, `10.0`, `1e1`, `1E+1`, `0.1e2`, `100e-2` (Wert 1), `1e0` (Wert
  1), `2147483647.000000000000000000000`. Abgelehnt: `15e-1`, `0.0`, `-0`, `-0.0e5`, `-10.0`,
  `1e30`, `1e999999`, `1e1000001`, `1e-999999`, `"5"`, `true`, `null`, `2147483648`,
  `1.0000000000000000001`, eine Ziffernfolge aus 200.000 Neunen, `1` mit 100.000 Nullen,
  `0.` mit 5.000 Nullen und `1`. Keine Eingabe lief über 200 ms. `ParseRouteSpec` meldet bei
  `1.5` `rule_spec ist ungültig: order fehlt oder ist keine ganze Zahl von 1 bis 2147483647`;
  der Use-Case-Text `rule_spec ist ungültig: <Adresse>` steht wörtlich in
  `TestSetRouteRejectsWithTheSpecTexts` und in `TestProcessAdministrationRequestsRouteOrderNotation`.
- Store-Test `TestAdministrationRequestSetRouteOrderLiteralsReachTheParser` gegen eine
  Wegwerf-PostgreSQL 18 (gepinnter Digest, Schema-Stand über `tools/schema/apply-rollout.sh`
  auf der Kopie, eigener Läufer): grün im Basislauf.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-antragsweg.md`
  → Exit 0, „22 Zeilen stimmen“; `make kommentar-kennungen DIFF=cc653a85~1` → Exit 0;
  `make fmt-check` → Exit 0.

**Mutationen, selbst gefahren** (Kopie im Scratchpad, Änderung als `sed … > Datei-in-der-Kopie`,
kein `-i`, je eine Kopie je Mutation):

| Nr | Stelle | Mutation | Ergebnis |
|---|---|---|---|
| R1 | `jsonRouteOrder` | `value.IsInt()` entfernt | rot: Modell-Ablehnungstest, Use-Case-Tabellentest, Store-Test (Fall `bruchteil`: `1.5` wird als 3 angenommen) |
| R2 | `jsonRouteOrder` | `order > MaxRouteOrder` zu `>=` | rot im Modell-Annahmetest (Fall „an der Obergrenze“); Use-Case-Paket grün (dort kein Grenzwert-Fall) |
| R3 | `jsonRouteOrder` | `big.Rat` durch `strconv.ParseFloat` + `SetFloat64` ersetzt | rot allein im Modell-Ablehnungstest (Fall `1.0000000000000000001`); Use-Case-Paket und Store-Test grün, weil dort kein Fall im Rundungsbereich liegt |
| R4 | `jsonRouteOrder` | Rückführung auf die Ziffernform (`strconv.ParseInt` des Rohtexts) | rot im Modell-Annahmetest, im Use-Case-Test `TestSetRouteAcceptsEveryNotationOfAPositiveWholeOrder` und in `TestProcessAdministrationRequestsRouteOrderNotation`; im Store-Test rot **allein** im Fall `mit_nachkommastelle` (`10.0`), `exponent` und `exponent_gross_mit_vorzeichen` bleiben grün |
| R5 | `syncAssemblerAddBinding` | Regelstand-Lesung durch eine leere Karte ersetzt (Wirkung der Mutation „leerer Routing-Port am Dekorator“) | rot: `TestEnableTableWithAssemblerSyncAddsBindingOnSuccess`, Meldung „Ziel nach Enable = "", wollen alle“ |

---

## Status der Findings des Hauptreports

| Finding | Status | Beleg |
|---|---|---|
| F-1 (HIGH) Plan-Satz „`1e1` endet …; `jsonb` bewahrt die Schreibweise“ | **geschlossen** | Der Plan-Satz ist ersetzt: „`jsonb` bewahrt `10.0` und normalisiert `1e1`/`1E+1` zu `10`; die Annahme hängt an keiner der beiden Schreibweisen“. Das deckt sich mit meiner Messung im Hauptreport. Der Parser nimmt jede Schreibweise an (25 Eingaben, R1, R4). Rest: F-N1 |
| F-2 (MEDIUM) Schreibweise enger als `SPEC-032`; Aufschub-Adresse nimmt sie nicht an | **geschlossen** | Der Code folgt dem Wortlaut der Spec. `slice-routing-betriebsdoku` trägt Obergrenze und Annahmemenge als committeten Text (`git grep 2147483647` im Plan der Adresse: Treffer). Der Entscheid steht ohne Architect-Artefakt (INFO F-N4) |
| F-3 (LOW) Begründung der Obergrenze nennt einen Verbraucher, den es nicht gibt | **geschlossen** | Godoc und Plan nennen die Grenze „eine Setzung“ (der größte Wert einer 32-Bit-Ganzzahl); der Godoc-Zusatz „begrenzt die Eingabe, die `ParseRouteSpec` annimmt“ trägt im Code (R2 färbt den Grenzwert-Test rot). Der Plan-Satz „weit über jeder Zahl von Regeln einer Tabelle“ ist eine Größenordnungsaussage ohne Anker, als Grund der Setzung ausgewiesen, nicht als Messung |
| F-4 (MEDIUM) API-Aktivierung trägt den Regelstand ungebunden | **teilweise geschlossen, Rest gedeckt** | Die Weitergabe im Dekorator ist gebunden (R5 rot). Die Feldbelegung `routing: activation` in `wiring.go` bindet kein Unit-Test; der Rest steht als eigener DoD-Punkt in `slice-routing-e2e` („API-Aktivierung“, aus Review-Frage A-2, *Zu belegen durch* `make test-integration` mit gedruckter Zeile je Weg). Ich habe den Punkt im Plan selbst gelesen: er nennt Pfad (HTTP und gRPC `EnableTable`), Gegenprobe (Tabelle ohne Regel) und die Grenze des Dekorator-Tests. A-2 ist damit beantwortet |
| F-5 (LOW) Wiederholungs-Zusage in `wiring.go` ohne Marker und Test | **geschlossen** | Der Kommentar sagt jetzt „Hergeleitet aus dem Code, ohne Wiederholungs-Test: R1 liest nur `applied`-Zeilen, und `Assembler.SetRoute` ersetzt nach dem Regelnamen“; die zweite, vormals unmarkierte Grenzbeschreibung ist gestrichen |
| F-6 (LOW) Interpunktion und Umbruch nach Teilersetzung | **weitgehend geschlossen** | Komma und Umbruch in `queries.go` (Aufzählung, „sieben übrigen Antragsarten“) und in `administrationrequest.go` behoben. Rest: F-N2 |
| F-7 (LOW) „alle sieben Funktionen“ im Godoc | **geschlossen** | `git grep -E 'alle sieben Funktionen\|sieben Antrags-Funktionen' -- internal` zählt 0 (Parent: 2). „sieben der neun Funktionen“ stimmt: der Test fährt `backfill_table`, `remove_`/`set_transformation`, `exclude_`/`include_column`, `disable_`/`enable_table`; „neun Antrags-Funktionen“ in `administration_endtoend_test.go` stimmt mit den neun `CREATE OR REPLACE FUNCTION cdc.` (Suchlauf-Zeile 9) |
| F-8 (INFO) äquivalente und ungebundene Mutanten | unverändert | keine Aktion erwartet |
| F-9 (INFO) `jsonb` in der ADR, `json` in Spec und Bestand | unverändert | keine Aktion erwartet |
| F-10 (INFO) Umfang | unverändert | die Fixrunde umfasst 14 Dateien (+236/−60, `git show --stat`); der Parser-Teil ist eine Funktion, die Tests sind je Paket klein, in einem Zug lesbar |
| A-1 | beantwortet | `SPEC-032` gilt wörtlich; `jsonRouteOrder` und das Handbuch-Aufschubfeld sind nachgezogen (siehe F-2, F-N4) |
| A-2 | beantwortet | siehe F-4 |

## Status der Befunde des Verifiers

| Befund | Behandlung |
|---|---|
| V-1 Fixrunde von keinem Reviewer gelesen | Mit diesem Report behoben: die Fixrunde `cc653a85` ist gelesen und mutationsgeprüft |
| V-2 Feldbelegung `routing: activation` ungebunden | wie F-4: Rest bei `slice-routing-e2e`, im Plan der Adresse gelesen |
| V-3 Aussage zum Abbruch von `run-store-tests.sh` nicht reproduzierbar | berührt die Fixrunde nicht; der Store-Test lief bei mir einzeln gegen eine eigene PostgreSQL (Basislauf und drei Mutationen), nicht über `run-store-tests.sh` |
| V-4 Kommentarzeile `queries.go:328` | bestätigt, siehe F-N2 |
| V-5 Rollen-Tests in `internal/bootstrap` bleiben bei weggelassenem `REVOKE` von `set_route` grün | bewertet als INFO F-N3 |
| V-6, V-7 | keine Fixrunden-Berührung; V-7 bleibt **übernommen** (am Endstand nicht ablesbar) |

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-N1 — Store-Test benennt die Schreibweise `1e1` als gebunden, bindet sie aber nicht an den Parser

- `kategorie`: LOW
- `quelle`: Hard-Rule-Klasse „Beleg trägt seinen Satz nicht“ (Reviewer-Skill), milde Form;
  `AGENTS.md` §3.12 Instanz B
- `pfad`: `internal/adapters/driven/postgresstorage/administrationrequest_routing_test.go:371-436`
  (Godoc und Tabelle von `TestAdministrationRequestSetRouteOrderLiteralsReachTheParser`)
- `befund`: Der Godoc sagt, `10`, `10.0`, `1e1` und `1E+1` gälten als die Ordnung 10, und nennt als
  rot färbende Mutation allein die Rückführung auf die Ziffernform (`10.0` endet abgelehnt). Die
  Mutation R4 färbt den Test nur im Fall `10.0` rot; `1e1` und `1E+1` erreichen den Parser als `10`,
  weil `jsonb` den Exponenten normalisiert, und beide Fälle sind dort an der Schreibweise
  ungebunden. Der Godoc sagt das nicht; der Plan nennt die Normalisierung. Außerdem trägt der Fall
  `{"null", 0}` den Namen `null` für das Literal `0`. Die Exponent-Schreibweise selbst ist an
  anderer Stelle gebunden (R4 rot im Modell-, Use-Case- und Bootstrap-Test), der Satz des Godocs
  ist also wahr, seine Stütze im Store-Test nur zur Hälfte.
- `verifizierbar`: ja — R4 auf der Kopie, Store-Test mit `-run` gegen eine Wegwerf-PostgreSQL.
- `klasse`: Beleg trägt seinen Satz nicht

### F-N2 — Rest von F-6: eine Kommentarzeile von 98 Zeichen in `queries.go`

- `kategorie`: LOW
- `quelle`: Maintainability (Rest von F-6 des Hauptreports, V-4 des Verifiers)
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:328`
- `befund`: Die Zeile „Prozessstart zum selben Stand. `requested_at` ist der Aufrufzeitpunkt der
  schreibenden Funktion“ misst 98 Zeichen im sonst umbrochenen Block; sie trägt die Spur der
  Teilersetzung. Inhalt und Interpunktion sind richtig.
- `verifizierbar`: ja — `awk 'length>80 && /^\/\//' queries.go` meldet die Zeile.
- `klasse`: Umbruch/Interpunktion nach Teilersetzung

### F-N3 — Rollen-Tests in `internal/bootstrap` binden das `REVOKE` von `set_route` nicht (V-5)

- `kategorie`: INFO
- `quelle`: Maintainability; Einordnung der Verifier-Frage V-5
- `pfad`: `internal/bootstrap/administration_roles_internal_test.go:285-292` (`createRequest`
  ruft die Funktionen über den Superuser-Pool `admin`), `internal/adapters/driven/postgresstorage/administrationrequest_routing_test.go:224-234`
- `befund`: Der Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins` belegt nach seinem Godoc, dass
  jede Anweisung des **Verarbeitungspfads** (Goroutine unter `cdc_admin`/`cdc_capture`) ihre Rolle
  trägt; die Anträge legt er über den Superuser an. Das Ausführungsrecht der Aufrufer gehört nicht
  zu seiner Zusage. Diese trägt `TestAdministrationRequestRoutingFunctionsRequireCdcAdminMembership`
  (Store, rot unter der `REVOKE`-Mutation, vom Verifier gesehen, von mir nicht nachgefahren) und
  `TestAdministrationDateiTraegtDieFunktionsRechte` (Text der Rollout-Datei, nennt `set_route` und
  `remove_route` in seiner Mutationsliste). Die Transformations-Funktionen stehen im selben Schnitt.
  Keine Lücke des Rollen-Tests, sondern ein Schnitt zwischen zwei Testebenen; beide Ebenen tragen.
- `verifizierbar`: ja — `REVOKE`-Mutation, Läufe `Role|Route` je Paket.
- `klasse`: Testebenen-Schnitt (keine Lücke)

### F-N4 — Entscheid zu A-1 steht als „Entscheid des Hauptlaufs“, ohne Verdikt-Artefakt

- `kategorie`: INFO
- `quelle`: Konflikt-Pfad (Modul 8): „Kein Pfeil ohne benennbares Artefakt“
- `pfad`: `docs/plan/planning/in-progress/slice-routing-antragsweg.md:200`
- `befund`: Der Plan nennt den Entscheid zu A-1 „Entscheid des Hauptlaufs“; ein Architect-Verdikt
  unter `docs/reviews/` liegt nicht vor (`ls docs/reviews | grep -i routing`: keine Verdikt-Datei
  zu `order`). Die gewählte Richtung ist die Spec-Fassung und schließt die Abweichung, statt eine
  neue zu setzen; ein Widerspruch zwischen Rollen besteht nicht. Der Übergang ist nur nicht
  beschriftbar.
- `verifizierbar`: nein
- `klasse`: Übergabe ohne Artefakt (Entscheid, kein Konflikt)

---

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/model/routespec.go` (`jsonRouteOrder`) — Korrektheit
  gegen 25 Eingaben (siehe oben); keine Gleitkomma-Rundung (R3 rot im Modell-Test, Fall
  `1.0000000000000000001`); `raw[0]` lässt nur `-` und Ziffern zu (der Rohtext ist gültiges JSON,
  `big.Rat.SetString` kennt dort keine Bruchform `a/b`); Exponenten bis 10^6 rechnet `big.Rat` in
  unter 200 ms durch und lehnt darüber ab; eine Ziffernfolge von 200.000 Stellen endet abgelehnt
  ohne messbare Last (unter 200 ms); `-0` und `0.0` enden bei `order < 1`; Fehlertext
  `rule_spec ist ungültig` wörtlich in den Tests.
- geprüft, ohne Befund: Godoc von `MaxRouteOrder` und `ParseRouteSpec` Punkt c — behaupten
  nichts, was der Code nicht trägt (F-3 geschlossen).
- geprüft, ohne Befund: Tests der Eingabe-Seite — `routespec_test.go`, `setroute/service_test.go`,
  `routing_internal_test.go`: die Eingabe jedes Falls ist die Verletzung (das Literal von
  `order`), keine Nachbarfälle; R1, R3, R4 färben sie rot. Lücke ohne Wirkung: Das
  Use-Case-Paket hat keinen Grenzwert-Fall (R2 grün dort), der Modell-Test trägt ihn.
- geprüft, ohne Befund: `internal/bootstrap/assemblersync_internal_test.go` — der Test trägt eine
  Regel im Port der Quelle und liest das Ziel am Assembler (R5 rot); der Kommentar benennt, dass
  die Belegung in `Run` offen bleibt.
- geprüft, ohne Befund: Kommentare des Diffs — `make kommentar-kennungen DIFF=cc653a85~1` ohne
  Kandidat, keine Chronik, kein Konjunktiv über verworfene Alternativen; `make fmt-check` grün.
- geprüft, ohne Befund: Träger-Nachzug — Suchlauf-Feld des Plans (22 Zeilen stimmen, darunter die
  neue Zeile für „sieben Funktionen“ mit Stand am Parent und am Diff);
  `slice-routing-betriebsdoku` trägt Obergrenze und Annahmemenge; `slice-routing-e2e` trägt den
  Punkt „API-Aktivierung“.
- geprüft, ohne Befund: Commit-Traceability `cc653a85` — `ADR-0137` und `LH-FA-CFG-008` im Betreff,
  keine Struktur-Kennung.
- nicht geprüft (Grenze): ob der Heredoc-Verstoß aus dem Hauptlauf (V-7) im Endstand Spuren
  trägt — **übernommen**, nicht gemessen; ein Hinweis darauf liegt in diesem Diff nicht vor.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht · Umbruch/Interpunktion nach
Teilersetzung · Testebenen-Schnitt (keine Lücke) · Übergabe ohne Artefakt (Entscheid, kein Konflikt)

## Verdikt

**Merge-blockierend:** nein. Das HIGH F-1 und die MEDIUM F-2 und F-4 des Hauptreports sind
geschlossen beziehungsweise mit benanntem, im Plan der Adresse gelesenem Rest bei
`slice-routing-e2e` gedeckt; die Fixrunde führt kein neues HIGH oder MEDIUM ein. Der Parser ist
über fünf eigene Mutationen und 25 Eingaben belegt.

**Übergabe:** F-N1 und F-N2 (beide LOW) gehen an den Implementer und brauchen **keine** eigene
Fixrunde: sie können mit dem nächsten Berühren der Dateien oder durch den Planner im
Closure-Zug gezogen werden. F-N3 und F-N4 sind ohne Aktion. Die Finding-Klassen gehen in die
Slice-Closure §7. **DoD-Zeile „Review durchgeführt, Report liegt vor“:** mit diesem Report liegt
der Artefakt-Beleg der Fixrunde vor (V-1); dieser Lauf hat kein Edit-Werkzeug und zieht die
Checkbox im Slice-Plan deshalb nicht selbst nach — der Planner zieht sie mit Verweis auf
diesen Report und den Hauptreport. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine
Verifikation (Modul 11).
