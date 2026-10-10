# Verifikationsbericht: slice-spec-festlegungen-doku-gates, Nachprüfung — 2026-10-10

**Rolle:** Verifier (Modul 11), Nachprüfung. Grundlage ist der erste Bericht
[`verify-slice-spec-festlegungen-doku-gates`](verify-slice-spec-festlegungen-doku-gates.md)
(`e4f609d8`, Verdikt *nicht bestanden* wegen A-1). Der erste Bericht ist ein
Lauf-Beleg und bleibt unverändert. Diese Datei ist der Folgelauf.

**Gegenstand:** Commit `bf8d26c2`. Er berichtigt den Absatz „Ausgänge“ von
[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
die Historie-Zeile, den Vertrag `harness/sensors/docs-check.md` §Ausgabe und im Plan
von `slice-spec-festlegungen-doku-gates` die Messungen M34 bis M43, die
Anschluss-Frage, den Beleg von LP 2 und den Block „Verifikation“. Er behebt A-1,
A-2 und A-5. Bezug ist
[`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
Entscheidung 1, wonach die Festlegung das Werkzeug wiedergibt.

**Umfang:** Nur die drei Punkte, alles andere gilt aus dem ersten Bericht weiter.
Gemessen habe ich mit d-check v0.82.0 über den Digest aus `d-check.mk`, und zwar
direkt per `docker run --rm --network none -v <Klon>:/repo:ro <Digest>` sowie über
`make docs-check`. Grundlage war ein Klon von `bf8d26c2` im Scratchpad, je Mutation
mit `git reset --hard && git clean -fd` zurückgesetzt, und `git status --porcelain`
war am Ende leer. Exit-Codes habe ich direkt gesichert (`AGENTS.md` §3.9).

---

## 1. Messungen

| # | Lauf | Ergebnis |
|---|---|---|
| N1 | Basislauf direkt, Ströme getrennt | Exit 0; stdout leer; stderr `d-check: 1835 Datei(en) geprüft, 0 Befund(e)` |
| N2 | drei fehlende Linkziele in `README.md`, direkt | Exit 1; stdout drei Zeilen mit vier Feldern (`README.md:65⇥n1.md⇥target-missing⇥Linkziel existiert nicht` …); stderr `… 3 Befund(e)` |
| N3 | N2 über `make docs-check > Datei 2>&1`, 20 Läufe (8 plus 12) | je Make-Exit 2, und je letzte Zeile `make: *** [d-check.mk:20: docs-check] Fehler 1`. **Die Summenzeile steht in 14 Läufen vor dem ersten Befund und in 6 Läufen hinter dem letzten** (erste 8 Läufe: Läufe 3, 6 und 8 mit Summe an Zeile 5 nach den Befunden an Zeile 2 bis 4) |
| N3b | N2 direkt `docker run … 2>&1`, 5 Läufe | Summenzeile 4-mal an Zeile 1 und 1-mal an Zeile 4 (nicht Gegenstand der Festlegung, nur zur Einordnung) |
| N4 | `~/<Verzeichnis>/x` in Inline-Code, direkt | Exit 1; stdout `README.md:66⇥<Pfad>⇥hostpath-forbidden`, also drei Felder ohne Text |
| N5 | unbekannter Schlüssel auf oberster Ebene von `.d-check.yml`, direkt | Exit 2; stdout leer; stderr `d-check: error: .d-check.yml: yaml: unmarshal errors:` und `  line 349: field unbekannt not found` |
| N5b | N5 über `make docs-check 2>&1` | Make-Exit 2; die beiden Fehlerzeilen, danach `make: *** … Fehler 2` |
| N6 | unbekanntes Modul `gibtsnicht` in `modules:`, direkt | Exit 2; stderr eine Zeile `d-check: error: unbekanntes Modul "gibtsnicht" in der Konfiguration — gültig: …` |
| N7 | ein fehlendes Linkziel über `make docs-check 2>&1` | Make-Exit 2; letzte Zeile ist die Fehlerzeile von `make` |

## 2. Text gegen Messung

| Satz (Stand `bf8d26c2`) | Ergebnis |
|---|---|
| `SPEC-040` §Ausgänge: Befundzeilen auf stdout mit vier Feldern, `hostpaths` mit drei Feldern | stimmt (N2, N4) |
| Summenzeile auf stderr; ein Lauf, der nicht prüft, schreibt eine Fehlermeldung nach stderr, die mit `d-check: error:` beginnt; ein YAML-Fehler trägt eine zweite eingerückte Zeile `  line <n>: …` | stimmt (N1, N5, N6) |
| Exit-Tabelle, Zeile 2 mit „ein unbekanntes Modul in `modules:`“ | stimmt (N6). A-5 ist erledigt |
| Über `make` kommt jeder Exit ungleich 0 als Make-Exit 2 an, und die letzte Zeile ist die Fehlerzeile von `make` | stimmt (N3, N5b, N7) |
| **„Über `make` mit beiden Strömen in einer Ausgabe (`2>&1`) steht die Summenzeile vor den Befundzeilen.“** | **stimmt nicht als feste Aussage.** In 6 von 20 Läufen steht sie dahinter (N3). Die Reihenfolge zweier Ströme in einer Ausgabe legt das Werkzeug nicht fest |
| Schlusssatz: „Ob ein Befund oder ein Fehler vorliegt, sagt … nicht die letzte Zeile, sondern die Summenzeile bzw. die Zeile `d-check: error: …`“ | stimmt, unabhängig von der Reihenfolge |
| Historie-Zeile: „über `make` mit `2>&1` die Summenzeile vor den Befunden“ | dieselbe Abweichung wie oben |
| Vertrag §Ausgabe: „ein Beleg zitiert die Summenzeile bzw. die Zeile `d-check: error: …`, nicht die letzte Zeile“ | stimmt |
| Plan M37: „je Zeile 2 die Summenzeile … fünf Läufe“; Anschluss-Frage: „über `make` mit `2>&1` die Summe vor den Befunden, in fünf Läufen gleich (M37)“ | Die gedruckte Messung habe ich nicht nachgefahren und führe sie als *übernommen*. Die Verallgemeinerung aus fünf Läufen trägt nicht (N3: 6 von 20 anders). Das ist dieselbe Klasse wie `AGENTS.md` §3.12 („Aussage breiter als ihre Messung“) |
| Plan M34 bis M36, M38 bis M43 | stimmen mit N1, N2, N4 bis N7 überein. Zeilennummern und die Dateizahl 1835 hängen am Stand und decken sich an `bf8d26c2` |
| Plan LP 2 Beleg (A-2): „Befund-Code-Tabelle …, einzelne Grund-Codes nennt er weiter als Kontext“ | stimmt. A-2 ist erledigt |

## 3. Abweichung

| ID | Klasse | Befund | Bewertung |
|---|---|---|---|
| A-6 | DoD (LP 1), die Festlegung behauptet mehr als das Werkzeug | Die neue Aussage, die Summenzeile stehe über `make` mit `2>&1` vor den Befundzeilen, trifft in 14 von 20 Läufen zu, nicht immer. Betroffen sind `SPEC-040` §Ausgänge (ein Satz), die Historie-Zeile vom 2026-10-10 und im Plan die Anschluss-Frage (Zeile „Befund-Zeilen und Summenzeile“) samt M37 | Klein, und der tragende Schlusssatz bleibt richtig. Eine mögliche Fassung: „Über `make` mit `2>&1` ist die Reihenfolge von Summen- und Befundzeilen nicht festgelegt.“ M37 bekommt die Messung N3 als Gegenbefund |

A-1 ist bis auf A-6 erledigt: Ströme, Feldzahl bei `hostpaths`, die zweizeilige YAML-Meldung und die letzte Zeile über `make` stimmen. A-2 und A-5 sind erledigt. A-3 und A-4 gehen weiter an die Closure.

## Endverdikt

**Nicht bestanden, wegen A-6.** Ein Satz der berichtigten Festlegung (und seine
Spiegelung in Historie und Plan) legt eine Reihenfolge fest, die das Werkzeug nicht
garantiert. Alles Übrige von A-1, A-2 und A-5 ist am Werkzeug bestätigt. Nach der
Berichtigung des einen Satzes, der Historie-Zeile und der zwei Plan-Stellen genügt
eine Lesung dieser Stellen. Eine weitere Messung braucht es nicht, N3 ist der
Beleg.

**Übergabe:** A-6 an den Implementer über den Planner; A-3 und A-4 an die
Closure. Dieser Bericht ist ein Lauf-Beleg.
