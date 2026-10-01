# Verifikations-Report: slice-routing-spec-nachzug — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)) und
Plan-vs-Code-Diff. Review-Artefakt:
[`review-slice-routing-spec-nachzug.md`](review-slice-routing-spec-nachzug.md). Formvorbild:
[`verifikation-slice-capture-retry-aufbau-frist-bindung.md`](verifikation-slice-capture-retry-aufbau-frist-bindung.md).

**Gegenstand:** Slice-Plan [`slice-routing-spec-nachzug`](../plan/planning/in-progress/slice-routing-spec-nachzug.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), Welle
[`welle-routing`](../plan/planning/welle-routing.md)), Diff `b8085839..HEAD` (`e2f9244c`): sechs
Commits, neun Dateien, davon `spec/pflichtenheft.md` und `spec/architecture.md`; **kein Code**,
[`spec/lastenheft.md`](../../spec/lastenheft.md) unberührt (`git diff --stat` leer). Dieser Lauf ändert
weder Spec noch Plan (keine DoD-Häkchen); er schreibt nur diesen Report.

## 1. Eigene Sensor-Belege (ungefiltert, Exit-Code einzeln gesichert, `AGENTS.md` §3.9)

| Sensor | Exit | Ausgabe |
|---|---|---|
| `make docs-check` | 0 | `d-check: 1471 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make gates` | 0 | `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks`, `a-check gesamt: 0 Befund(e)` (Tail; Exit in Datei gesichert) |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-spec-nachzug.md` | 0 | `16 Zeilen stimmen` (acht am Parent `30fd6cb5`, acht am Arbeitsbaum `diff`) |
| `make commit-traceability RANGE=b8085839..HEAD` | 0 | `6 Commit(s) … Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=b8085839..HEAD` | 0 | `0 Befund(e)` |
| `make doc-immutable RANGE=b8085839..HEAD` | 0 | `0 Befund(e)` (ohne `RANGE` Exit 2, `--range` braucht ein Argument — Aufruffehler, kein Befund) |

Nicht gefahren: `make test*` (kein Code im Diff). Zusätzlich selbst gemessen:

- `git diff b8085839 HEAD -- spec | grep '^+' | grep -c -E 'ADR-[0-9]|slice-|welle-|Slice |Welle '` → `0`
  (weder `spec/architecture.md` noch `spec/pflichtenheft.md` tragen im Diff einen ADR-, Slice-,
  Wellen- oder Commit-Hash-Bezug; Hex-Muster ebenfalls ohne Treffer). Treffer von `git grep` in
  `spec/architecture.md` sind ausschließlich der Kopf-Text der Hard Rule (Zeilen 8, 9, 21, 65, Bestand).
  Die Treffer `welle-`/`slice-` in `spec/pflichtenheft.md` liegen in alten Historie-Zeilen (Bestand).
- Zählwörter am Diff-Stand: `request_kind` neun Werte; `column_name` „sieben übrigen" (9 − 2), `rule_name`
  „fünf übrigen" (9 − 4), `rule_spec` „sieben übrigen" (9 − 2); „Die neun SQL-Funktionen" in
  [`SPEC-019`](../../spec/pflichtenheft.md); Architektur „neun Arten" bei neun Tabellenzeilen
  (selbst gezählt). Die Nachrichten-Zählwörter „zehn Felder"/„dreizehn Felder": 22 Zeilen an beiden
  Ständen (unverändert, richtig: das Label ist kein Nachrichtenfeld).

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | [`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md) beantwortet, Zusagen in Zukunfts-Form, Abhilfe als Zusage | **getragen** | Überschrift ohne „offen"; Absätze Zielmodell bis „Nicht anwendbare Regel" gelesen (`spec/pflichtenheft.md` Z. 385–440); DELETE-Aussage als „hergeleitet … nicht gemessen" gekennzeichnet; `make docs-check` Exit 0 |
| 2 | Datenstrukturen und Schnittstellen stehen ([`SPEC-019`](../../spec/pflichtenheft.md), Regelform, [`SPEC-001`](../../spec/pflichtenheft.md)/[`SPEC-002`](../../spec/pflichtenheft.md), [`SPEC-008`](../../spec/pflichtenheft.md), [`SPEC-029`](../../spec/pflichtenheft.md), [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md), [`SPEC-020`](../../spec/pflichtenheft.md)/[`SPEC-021`](../../spec/pflichtenheft.md)/[`SPEC-022`](../../spec/pflichtenheft.md)/[`SPEC-031`](../../spec/pflichtenheft.md)/[`SPEC-024`](../../spec/pflichtenheft.md)); neue Kennung ist die nächste freie; §7 Historie je Änderung | **getragen mit zwei LOW** | alle Stellen gelesen (§3); [`SPEC-032`](../../spec/pflichtenheft.md) ist die nächste freie nach [`SPEC-031`](../../spec/pflichtenheft.md); Historie-Zeilen ohne ADR-/Slice-Bezug (3 + 1 Zeilen vom 2026-10-01); LOW V-1, V-2 |
| 3 | `spec/architecture.md`: neun Arten, Aussage zum Ziel vor der Persistierung, ohne ADR-/Slice-/Wellen-Bezug | **getragen** | Tabelle neun Zeilen; Absatz „Zustellziel der Change"; Bezug-Grep §1 = 0 |
| 4 | `make gates` grün, Exit-Code ungefiltert gesichert | **getragen** | §1, eigener Lauf, Exit 0 |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM | **teilweise — Prozessbefund V-3** | Report liegt vor, 0 HIGH, 2 MEDIUM (F-1, F-2); beide inhaltlich geschlossen (§4), aber kein Folge-Review/Re-Check in einem Artefakt; das Review selbst hält die Zeile bis dahin ausdrücklich offen |
| 6 | §3.13-Suchlauf: Feld trägt Gefundenes und Nichtgefundenes, beide Stände; Nachmessen Exit 0 | **getragen** | Exit 0, 16 Zeilen; die Befund-Spalte nennt je Träger „nicht gefunden"-Aussagen; Träger außerhalb des Diffs (`harness/targets/schema-rollout.md`, Handbuch) sind an benannte Adressen gemeldet |
| 7 | Doku-Update entfällt | **getragen** | Slice ist das Doku-Update; `docs/user` im Diff unberührt |
| 8–12 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §6/§7 tragen Platzhalter (erwartet vor `done/`); die Paarungen hängen an der Closure der Welle |

Die DoD-Häkchen im Plan stehen unverändert alle auf `[ ]`; ich setze keine.

## 3. Plan-vs-Code-Diff

Plan-Tabelle §3 gegen Diff, Zeile für Zeile:

- **Im Plan, im Diff:** [`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md), [`SPEC-019`](../../spec/pflichtenheft.md),
  [`SPEC-032`](../../spec/pflichtenheft.md) (nächste freie Kennung am Parent),
  [`SPEC-001`](../../spec/pflichtenheft.md)/[`SPEC-002`](../../spec/pflichtenheft.md) (`route_target`, letzte Spalte der View),
  [`SPEC-020`](../../spec/pflichtenheft.md)/[`-021`](../../spec/pflichtenheft.md)/[`-022`](../../spec/pflichtenheft.md)/[`-024`](../../spec/pflichtenheft.md),
  [`SPEC-031`](../../spec/pflichtenheft.md) Zeile `ReadChanges` (`target`, Feldnummer 7),
  [`SPEC-008`](../../spec/pflichtenheft.md) Zeile `schema`, [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md),
  Historie, `spec/architecture.md`. [`SPEC-029`](../../spec/pflichtenheft.md) ist im Diff nicht berührt; der Run-Satz
  steht in [`SPEC-008`](../../spec/pflichtenheft.md) und [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) — der Review
  (Negativbefund) hat das bereits als widerspruchsfrei bewertet ([`SPEC-029`](../../spec/pflichtenheft.md) verweist
  generisch auf [`SPEC-008`](../../spec/pflichtenheft.md)); die Plan-Kopfzeile nennt es als zu ändernd, ohne dass es
  geschah (LOW V-4, Plan-Wortlaut, kein Spec-Fehler).
- **Im Diff, über den Plan hinaus, im Plan ausgewiesen:** [`ARC-001`](../../spec/architecture.md)-Komponentenliste,
  Backfill-Absatz, Fehlermodell-Zeile der Sicht; Fixrunden-Zeilen. Kein Mehrumfang ohne Plan-Zeile.
- **Nicht im Diff, aber behauptet:** [`SPEC-021`](../../spec/pflichtenheft.md) in der Fixrunde (Plan-Zeile 199,
  Historie-Zeile „… auch mit U+0000"). Der Fixrunden-Diff (`4b6dc973..HEAD`) ändert die Zeile nicht; ihre Zeile
  „Query-Parameter" trägt weder die Leer-Aussage noch U+0000, sondern verweist auf „dieselbe Kombinatorik wie
  `SPEC-020`s gRPC-Request" — die Aussage gilt dort allenfalls durch Verweis (LOW V-1).
- **Code:** keiner im Diff; Lastenheft unverändert.

## 4. Review-Findings F-1 bis F-4 — am Text geprüft, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg am Text |
|---|---|---|
| F-1 (MEDIUM) Regelstand im Run | **geschlossen** | [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) „Fail-closed vor dem Commit" nennt den Routing-Regelstand, Klasse `configuration`, einmal nach dem Öffnen des Snapshots gelesen, Anwendbarkeit (Klasse `schema`), jeder weitere Block und vor dem Commit Mengenvergleich, nicht lesbarer Stand = Abweichung, Set-und-Rücknahme als benannte Grenze; „Regelstand zum Run" im Punkt „Ziel der Backfill-Changes" definiert (Z. 247–273) |
| F-2 (MEDIUM) `target` außerhalb des Alphabets | **geschlossen mit LOW V-1** | [`SPEC-020`](../../spec/pflichtenheft.md), [`SPEC-022`](../../spec/pflichtenheft.md), [`SPEC-031`](../../spec/pflichtenheft.md) nennen U+0000, den Mechanismus (Prüfung im gemeinsamen Use Case vor dem Speicherzugriff bzw. Vergleich im Speicher) und die Ausnahme `cdc.changes`; [`SPEC-021`](../../spec/pflichtenheft.md) nur durch Verweis |
| F-3 (LOW) `order`-Grenze | **geschlossen mit LOW V-2** | [`SPEC-032`](../../spec/pflichtenheft.md) Zeile `order`: „JSON-Zahl, positive ganze Zahl"; `2147483647`/„Exponent" nirgends mehr; Rest: [`SPEC-019`](../../spec/pflichtenheft.md) Fehlertext-Zeile (Z. 926) sagt weiter „im Wertebereich" |
| F-4 (LOW) Umbruch | **geschlossen** | `spec/architecture.md` (Z. 417–420) und [`SPEC-032`](../../spec/pflichtenheft.md) (`when.equals`-Absatz) umbrochen; im Fixrunden-Diff gelesen |

F-5 und F-6 (INFO) sind unverändert und tragen keine Zusage.

## 5. Entscheidungs-Konformität

**[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 1 (A-1)** — übernommen: einmal lesen nach dem Öffnen des Snapshots, Anwendbarkeit `schema`
vor der ersten Zeile, danach Mengenvergleich je Block und unmittelbar vor dem Commit, Abweichung =
`failed`/`configuration`, Vergleichsstand weder Annahme noch je Block, Grenze (Set-und-Rücknahme)
benannt. Nicht wörtlich in der Spec, aber folgerichtig: die Reihenfolge `schema` vor `configuration`
steht nur implizit („einmal nach dem Öffnen … geprüft", „jeder weitere Block"); die Abhilfe „neuer
Antrag" steht dort nicht (für den Run im Bestandsabsatz „Unterbrechung und Neubeginn" abgedeckt) —
beides INFO, keine Abweichung.

**Festlegung 2 (A-2)** — übernommen für gRPC-Stream, `GET /changes`, `ReadChanges`: leer, kein Fehler,
auch U+0000, Prüfung im gemeinsamen Use Case; SQL-Zugriff ausgenommen. SSE: siehe V-1.

**`AGENTS.md` §3.12 (hergeleitet nicht als erprobt):** Das PostgreSQL-Verhalten gegenüber U+0000 steht
nirgends in der Spec (die Spec macht den Satz davon unabhängig, wie der Plan sagt); der Mechanismus
„Use Case prüft vor dem Speicherzugriff" ist eine Zusage an die Umsetzung und steht in der ADR als
*erwartet* — die Spec führt ihn nicht als gemessen. DELETE gegen PostgreSQL 17/18 und die NATS-Last sind
als „hergeleitet"/„nicht gemessen" gekennzeichnet. Keine Aussage steht als erprobt, die eine ADR als
hergeleitet führt. Befund keiner.

**[`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)/[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md):**
Alphabet, Regelschlüssel, R1–R6, erster Treffer, `NULL` ohne Standardziel, nicht rückwirkend,
`target = 7`, Subjekt `cdc.route.<source_id>.<ziel>`, Run-Satz (`schema`, einmal, run-lokal) stimmen mit
den ADR-Texten überein (Stichprobe der Fehlertext-Tabelle, Beispiel in [`SPEC-032`](../../spec/pflichtenheft.md),
Auswahl-statt-Ausblendung-Absatz).

## 6. Offene Punkte — ehrlich geführt, DoD-unabhängig

- **A-3** (Lastenheft-Klarstellung „Zustellung = Abruf/Abonnement"): in
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  („Offen, nicht entschieden"), in `welle-routing` und in der Plan-Zeile 203 als offen geführt; die
  Spec trägt die Lesart als Zusage. Keine DoD-Zeile hängt von ihr ab. Der Slice-Plan §6 führt A-3
  nicht als eigenes Risiko (nur in der Tabellenzeile) — Empfehlung, es bei der Closure dort einzutragen.
- **V3** (Erreichbarkeit der Nichtanwendbarkeit): in [`SPEC-008`](../../spec/pflichtenheft.md) und
  [`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md) als Zusage/„nicht gemessen"; Plan §6 führt es als offen;
  keine DoD-Zeile hängt davon ab.

## 7. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | LOW | Die Fixrunde nennt [`SPEC-021`](../../spec/pflichtenheft.md) in Plan-Zeile 199 und Historie-Zeile, ändert die Zeile aber nicht; die „leer, kein Fehler"-/U+0000-Aussage steht für SSE nur durch Verweis auf [`SPEC-020`](../../spec/pflichtenheft.md). [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Folgepflicht nennt [`SPEC-021`](../../spec/pflichtenheft.md) ausdrücklich. Plan/Historie behaupten mehr als der Diff | `spec/pflichtenheft.md` Z. 1000 gegen Z. 1575 |
| V-2 | LOW | Rest des zurückgenommenen Wertebereichs: die Fehlertext-Zeile sagt „keine positive Ganzzahl **im Wertebereich**", obwohl [`SPEC-032`](../../spec/pflichtenheft.md) keinen Wertebereich mehr führt | `spec/pflichtenheft.md` Z. 926 |
| V-3 | Prozess | DoD „kein offenes HIGH/MEDIUM": die beiden MEDIUM sind am Text geschlossen (§4), aber kein Reviewer hat die Fixrunde gelesen; der Report selbst hält die Zeile bis dahin offen | DoD-Zeile 5 |
| V-4 | LOW | Plan-Kopf „Berührte Spec-Stellen" nennt [`SPEC-029`](../../spec/pflichtenheft.md) als zu ändernd; der Diff berührt es nicht (begründet erlaubt, s. §3); Kopf und Ergebnis weichen ab | Plan Kopf, Z. 39 |

## 8. Verdikt

**DoD getragen: ja für die Zeilen 1, 3, 4, 6, 7; Zeile 2 mit zwei LOW (V-1, V-2); Zeile 5 nicht
abschließend (V-3); Zeilen 8–12 offen (Planner, erwartet).** Kein HIGH, kein MEDIUM, kein Blocker für
Folge-Slices: die Spec trägt den beschlossenen Stand, A-1 und A-2 sind korrekt aus
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
übernommen, `AGENTS.md` §3.4 und die Stratum-Regeln sind eingehalten, Lastenheft unverändert.

**Nötiger Nachzug:**

1. **Spec (Implementer, ein kleiner Zug):** V-1 — [`SPEC-021`](../../spec/pflichtenheft.md) Zeile
   „Query-Parameter" die Leer-Aussage samt U+0000 und Ausnahme `cdc.changes` tragen lassen (oder
   Plan-/Historie-Zeile auf „durch Verweis auf [`SPEC-020`](../../spec/pflichtenheft.md)" zurücknehmen);
   V-2 — „im Wertebereich" in Z. 926 streichen.
2. **Reviewer:** kurzer Re-Check der Fixrunde (V-3), damit die DoD-Zeile 5 an einem Artefakt hängt.
3. **Planner:** DoD-Zeilen 1–4, 6, 7 abhaken (Belege: §1–§4) nach Nachzug 1; V-4 im Plan-Kopf
   angleichen; Closure-Notiz mit Lerneintrag (nahe liegend: „eine Fixrunde, die eine Zeile nennt,
   ändert sie auch" — Behauptung gegen Diff), Register-Vermerk, Risiko-Ausgänge §6 inklusive A-3.
