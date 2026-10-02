# Architect-Verdikt: Meldungscodes statt interner Kennungen in Programm-Ausgaben

**Rolle:** Architect (Modul 8)
**Anlass:** Startbedingung des Slice
[`meldungscodes-statt-interner-kennungen`](../plan/planning/done/slice-meldungscodes-kennungsfreie-ausgaben.md)
(Liefer-Punkt A; Planner-Fragen 1 bis 7 und die Zusatzfragen des Auftraggebers)
**Datum:** 2026-10-02
**Bezug:** [`LH-QA-OPS-001`](../../spec/lastenheft.md),
[`LH-QA-REL-003`](../../spec/lastenheft.md),
[`LH-FA-ADM-003`](../../spec/lastenheft.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Ergebnis dieses Verdikts),
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md),
[`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md),
[`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md),
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
`AGENTS.md` §3.5, §3.6, §3.11, §3.12

---

## Verdikt in einem Satz

Die Entscheidung ist ein Meldungscode `PCF-<E|W|I><4 Ziffern>` je Fehler- und
Warnursache, mit der Fehlerklasse in der ersten Ziffer (Klasse bleibt die Maschinen-Achse),
einer Tabelle im Code als Quelle, einem Katalog im Handbuch und zwei netzlosen Gates; die
Umsetzung läuft in vier Teil-Slices, deren erster ohne Codes die Kennungen aus der Ausgabe
nimmt. Der Plan hat nichts falsch behauptet; er braucht Nachzüge (unten). Die zwei
Accepted-ADRs bleiben unberührt, die Schärfung trägt die neue
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md).

## Entscheidungen Q1 bis Q7

| Frage | Entscheidung | Kern der Begründung |
|---|---|---|
| Q1 Schema | **`PCF-` + Schwere (`E`/`W`/`I`) + 4 Ziffern**. Fehler: erste Ziffer 1 bis 7 = Klasse (Reihenfolge der `SPEC-008`-Tabelle), 8 = Ablehnung einer Aufrufer-Eingabe, 9 reserviert; Warnung: erste Ziffer = Bereich (Erfassung, Backfill, Retention/Speicher, Verwaltung, Konfiguration/Start). `…000` = Rückfall der Klasse. `I` reserviert | Das Präfix, weil d-migrate im Rollout in dasselbe Terminal `E`/`W`-Codes druckt (hergeleitet, nicht nachgemessen); vier Ziffern, weil d-migrate bei drei seine Bereiche sprengen musste; der Rückfall `…000` garantiert, dass kein klassifizierter Fehler ohne Code bleibt |
| Q2 Klasse und Code | **Code in der Klasse**: Klasse folgt aus der ersten Ziffer, `error_class`, das Metrik-Label und das Wort `Fehlerklasse` bleiben; **kein Code als Metrik-Label** | Zwei unabhängige Achsen könnten sich widersprechen; „Klasse aus dem Code, `error_class` entfällt“ bricht Heartbeat, Metrik und jeden lesenden Betreiber |
| Q3 Abbildung | **Feld, wo der Weg eines hat; sonst Kopf `Fehlerklasse <klasse> [<code>]: …`.** Log-Warnung: Attribut `code`. Ausgang 1 bleibt, kein numerischer Ausgang je Code. `error_message` von Run und Antrag: `<klasse> [<code>]: …` bzw. `abgelehnt [<code>]: …`. Heartbeat: additive Spalte `error_code`. HTTP: Feld `code`. gRPC: `ErrorInfo.reason`. SQL-Funktionen: kein Code (0 `RAISE EXCEPTION`, gemessen). Rollout: `FEHLER [<code>]` | Die Klammer-Form lehnt sich an d-migrate (`Error [E053]: …`) an und lässt Lesern von `Fehlerklasse schema` (ohne Doppelpunkt) die Prüfung; jeder Weg bekommt seine natürliche Form |
| Q4 Stabilität | Code nie neu belegt; Zurückziehen (`withdrawn`) statt Löschen; Bedeutung nicht erweitert, Klasse eines Codes ändert sich nie; **Code, Klasse, Wort `Fehlerklasse`, Ausgang 1 stabil, Text nicht** | Texte tragen Laufzeitdetails (Regelname, Spalte) und können nie stabil sein; der Betreiber braucht einen Anker, der es ist |
| Q5 Dokumentation und Prüfung | **Tabelle im Go-Code (`internal/domain`) als Quelle + Handbuch-Katalog + Gate `meldungscodes-check` (Mengengleichheit Quelltext/Tabelle/Katalog) + Go-Test (Format, Ziffer gegen Klasse, Dopplung, Rückfall)**; kein YAML-Ledger | d-migrate braucht das Ledger, weil seine Codes über Kotlin-Module und Dokumente verteilt sind; hier liegen Code und Test im selben Modul, eine dritte Stelle wäre Pflege ohne Gewinn. Abgrenzung zu `ADR-0143`: `PCF-E4003`, `PCF-W2001`, `PCF-I1001` und `E053` treffen das Muster `P` nicht (je 0, gemessen) |
| Q6 Spec-Stelle | **`SPEC-008` um einen Absatz erweitern**, kein neues `SPEC-*`, keine Lastenheft-Änderung | §4 heißt bereits „Fehler-Codes und Logging-Felder“; eine Lastenheft-Anforderung wäre Auftraggeber-Entscheidung und ist für den Betrieb nicht nötig |
| Q7 Reichweite | **Fehler, Ablehnungen von Aufrufer-Eingaben und Warnungen**; Diagnose-Zeilen (außer „Fehlerzustand“), Fortschrittszeilen und Info-Log tragen keinen Code, aber auch keine Kennung | Zustandsberichte brauchen keinen Anker je Zeile; Rauschen |
| Zusatz Umfang | **vier Teil-Slices:** T1 `meldungscodes-kennungsfreie-ausgaben` (ohne Codes, mit Gate), T2 `meldungscodes-registry-fehlerkopf`, T3 `meldungscodes-warnungen-heartbeat-diagnose`, T4 `meldungscodes-http-grpc-fehlerkoerper`; Reihenfolge T1, T2, danach T3 und T4 unabhängig | T1 erfüllt den sichtbaren Teil der Vorgabe sofort und braucht keine Code-Entscheidung; Spalte `error_code` (Rollout, `make test-store`) und Netz-Wege sind eigene Diffs |
| Zusatz Kompatibilität | **keine Übergangsregel** (kein Doppelformat, kein Schalter): Text ist nicht Vertrag; die Änderungshistorie des Handbuchs nennt es in Betreibersicht | der Server ist vor 1.0; ein Doppelformat verdoppelte die Pflege für eine Zusage, die es nie gab (der Plan nennt sie nicht, das Handbuch ist darauf nicht nachgemessen: Prüfauftrag an T2, Rückfall an den Auftraggeber, wenn eine Textzusage gefunden wird) |
| Zusatz Wächter | **Gate `ausgabe-kennungen-check` in `make gates`, ab Ist-Stand 0 (mit T1)**; ERE (kein `-P`), zweistufig, Gegenstand Produktions-Go unter `internal/`, `cmd/`, `tools/schema/` und `echo`/`printf` in `tools/schema/`, `examples/` | 14 Treffer, alle gelesen und echte Literale, 0 Falsch-Positive am Bestand; die Grenzen (nachgestellter Kommentar mit zitiertem Literal trifft; mehrzeiliges Raw-String-Literal wird nicht gelesen) stehen in der ADR |

## Messungen (Ursprung nach `AGENTS.md` §3.12)

- **Gemessen** am Stand `abf71ace`, 2026-10-02:
  `git grep -n -P '^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]' -- '*.go' ':!*_test.go'` liefert **14** Zeilen
  (12 `wiring.go`, 1 `config_file.go`, 1 `guard.go`, alle gelesen). Die zweistufige ERE-Form
  (`git grep -n -E '^[^"`]*("[^"]*"[^"`]*)*["`][^"`]*(LH|ADR|SPEC|ARC)-[A-Z0-9]'` abzüglich
  Kommentarzeilen) liefert ebenfalls **14**. `echo`/`printf`-Zeilen mit Kennung in
  `tools/schema` und `examples`: **3** (`rollout.sh`, drei Zeilen).
  `git grep -n -E '"Fehlerklasse [a-z]+: ' -- '*.go' ':!*_test.go'`: **28** Literale.
  `git grep -n 'RAISE EXCEPTION' -- tools/schema internal/adapters/driven/postgresstorage`: **0**.
- **Gemessen**, Eingabezeilen: acht `printf`-Zeilen an der zweistufigen Pipeline — Treffer bei
  Literal, Literal mit `/` davor, Literal mit URL, nachgestelltem Kommentar mit zitiertem
  Literal (Falsch-Positiv); kein Treffer bei Kommentarzeilen, nachgestelltem Kommentar und
  Struct-Tag in Backticks. Beim ersten Versuch lieferte die Schleife (Eingabe über `stdin`
  des Schleifenkörpers) falsche Treffer; die korrigierte Form (Auswertung per
  Befehlsersetzung) trägt die Zahlen in der ADR.
- **Gemessen**, Code gegen Muster `P` (`ADR-0143`): `PCF-E4003`, `PCF-W2001`, `PCF-I1001`,
  `E053` je 0 Treffer; `LH-FA-X`, `ADR-0049`, `CO-123` je 1 Treffer.
- **Gemessen**, Klassifikation: `classifyRunError` (`internal/bootstrap/wiring.go`) bildet
  Sentinels über `errors.Is` auf sieben Klassen ab; der Backfill-Dienst entfernt den Kopf
  `Fehlerklasse <Klasse>: ` per Textersetzung (`internal/application/usecase/backfill/service.go`):
  die Umsetzung (T2) ersetzt diese Textchirurgie durch den Code des Fehlerwerts.
- **Hergeleitet, nicht erprobt:** Mutation des Wächters, Registry-Test,
  `meldungscodes-check`, Gate-Verdrahtung, Verhalten von `make handbuch-public-doc-check` mit
  einem Katalog, Verwechslung mit d-migrate-Codes im Rollout, Toleranz der SDKs gegen
  zusätzliche JSON-Felder. Der Implementer fährt sie; das ist eine Erwartung (ADR Fitness Function).
- **Nicht nachgemessen:** `description:`-Felder von `tools/schema/schema.yaml`: ob sie als
  Datenbank-Kommentar beim Betreiber ankommen (kein `COMMENT ON` im Quelltext gefunden,
  `git grep -n 'COMMENT ON' -- tools/schema` 0 Zeilen; ob d-migrate sie erzeugt, ist offen) —
  Prüfauftrag in T1.

## Plan-Nachzugsliste (der Planner; der Architect ändert den Plan nicht)

1. **Kopf und §2 (Liefer-Punkte):** (A) ist mit `ADR-0144` erfüllt (Index-Eintrag vorhanden).
   Der Slice wird zu **T1** `meldungscodes-kennungsfreie-ausgaben` zugeschnitten (Liefer-Punkte:
   Kennungen aus den 14 Go-Zeilen und 3 `echo`-Zeilen, Läufer-Erwartungen, Handbuch-Hinweissatz und
   Platzhalter, `description:`-Entscheidung, Gate `ausgabe-kennungen-check`); T2 bis T4 sind
   neue Pläne mit den Namen aus der ADR. Der Name des vorhandenen Slice ist dann zu
   wechseln oder zu begründen (`MR-002`).
2. **§3 Wächter-Option:** das `-P`-Muster `^[^/]*…` ist durch die ERE-Form aus der ADR
   (Festlegung 10) ersetzt (`grep -P` gehört nicht zur Host-Klasse, `AGENTS.md` §3.1); die Frage
   „Gate oder Werkzeug“ ist entschieden: Gate ab Ist-Stand 0; Reichweite und Ausnahmen stehen
   in der ADR. Die Suchlauf-Zeilen bleiben `-P` (ein Suchlauf ist kein Gate).
3. **§3 Tabelle „Fehlertexte mit Präfix“:** die 28 Literale gehören zu **T2**, nicht zu T1.
4. **§4 Start:** die sieben Fragen sind beantwortet; die Rückführung „zu groß“ ist
   mit der Aufteilung der ADR (Festlegung 9) eingetreten. Der Hinweis „Registry samt Test
   und Umstellung aller Wege“ trifft zu und ist aufgeteilt.
5. **§6 Risiken:** „Ausgaben-Stabilität für Betreiber“ — Ausgang: ADR Konsequenz (keine Übergangsregel,
   Text nicht Vertrag); der Handbuch-Abschnitt dazu wird in T2 gemessen. „Datenbank-Metadaten“ —
   Prüfauftrag in T1. „SDKs reichen Servertexte durch“ — gilt für T2 bis T4, nicht für T1.
6. **§1 Zählung der Test-Dateien:** der Plan nennt 142 Zeilen in `*_test.go`; mit
   `git grep -c -P '^[^/]*(LH|ADR|SPEC|ARC)-[A-Z0-9]' -- '*_test.go'` summiere ich **97** Zeilen
   (anderes Muster als im Plan, die 142 sind nicht nachgemessen); die Differenz klärt der Implementer
   von T1 mit dem Muster aus dem Plan.

## Spec-Nachzug (Planner/Auftraggeber; der Architect ändert `spec/` nicht)

- `spec/pflichtenheft.md` §4: ein Absatz „Meldungscode“ unter `SPEC-008` (Form, Nummernraum,
  Stabilität, Kopf); der Satz „der Fehlertext beginnt mit der Klasse (`schema: `)“ in der Prosa
  zur Nichtanwendbarkeit wird zum Kopf aus `ADR-0144` Festlegung 3 — **mit T2**.
- Mit dem Teil-Slice des jeweiligen Wegs: `SPEC-029` und `SPEC-019` (`error_message`, T2),
  Tabelle von `cdc.process_heartbeat` (`error_code`, T3), `SPEC-018` (HTTP-Fehlerkörper, T4),
  `SPEC-031` (gRPC-Fehler, T4).
- Lastenheft: keine Änderung nötig. Eine Anforderung („Fehler und Warnungen tragen stabile
  Codes“) unter `LH-QA-REL-003` oder `LH-QA-OPS-001` wäre möglich und ist eine
  Auftraggeber-Entscheidung.

## Empfehlungen an den Auftraggeber

1. **Präfix `PCF-` bestätigen** (oder anderes Kürzel nennen); nach T2 ist es nur noch per
   Folge-ADR änderbar.
2. **T1 zuerst und allein ausliefern** — die Kennungen verschwinden aus der Ausgabe, ohne
   dass über Codes entschieden sein muss; die Codes folgen in T2 bis T4.
3. **Server-Release:** T2 ändert den Wortlaut des Fehlerkopfs sichtbar; ob das ein Release
   mit eigener Freigabe braucht, entscheidet der Auftraggeber (die ADR setzt kein Tag).
4. **Lastenheft-Anforderung:** nur wenn gewünscht (siehe Spec-Nachzug); der Betrieb braucht sie nicht.

## Akzeptierte Negative

- Das Gate `ausgabe-kennungen-check` liest keine mehrzeiligen Literale (alle 14 Treffer des
  Bestands sind einzeilig; ein AST-Wächter wäre ein Docker-Programm ohne bisherigen Fund).
- `meldungscodes-check` vergleicht Mengen, nicht Sinn; der Reviewer liest den Katalog gegen den Code.

Weder Produktionscode noch Plan noch Spec noch `AGENTS.md` wurden im Rahmen dieses Verdikts
geändert. Geändert wurden ausschließlich `docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md`,
`docs/plan/adr/README.md` (eine Index-Zeile) und diese Datei.
