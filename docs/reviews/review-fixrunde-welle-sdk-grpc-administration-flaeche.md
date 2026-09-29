# Review-Report: Fixrunde zu welle-sdk-grpc-administration-flaeche — 2026-09-29

**Review-Art:** Code — geprüft gegen die beiden vorangegangenen
Review-Reports (`docs/reviews/review-sdk-kotlin-grpc-administration-flaeche.md`,
`docs/reviews/review-sdk-python-grpc-administration-flaeche.md`), die
Slice-Pläne derselben Namen und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten).

**Gegenstand:** zwei Fixrunden-Commits auf `main`, Range `70e19f51..48e04899`:

- `bd10c391` — fix(sdk): Fixrunde zu beiden Reviews — Kotlin-Fixtures
  übersetzt, Python-Suchlauf-Zahl korrigiert (`LH-FA-SST-009`)
- `48e04899` — fix(sdk): C#-Error-Mapping-Fixtures ins Englische übersetzt
  (`LH-FA-SST-009`)

**Skill:** `.harness/skills/reviewer.md` (Arbeitsbaum-Stand beim Review-Lauf)
**Modell:** glm-5.3-flash · **Datum:** 2026-09-29

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/reviews/review-sdk-kotlin-grpc-administration-flaeche.md` (F-1 HIGH,
  F-2/F-3 LOW, F-4 INFO)
- `docs/reviews/review-sdk-python-grpc-administration-flaeche.md` (F-1 HIGH)
- `docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-administration-flaeche.md`
- `docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md`
- `internal/adapters/driving/grpc/interceptor.go` (realer Server-Fehlertext,
  Gegenprobe zur Fixture-Übersetzung)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.13
- `LH-FA-SST-009`

---

## Findings

**Keine.** Der Fixrunden-Diff trägt 0 HIGH, 0 MEDIUM, 0 LOW, 0 INFO.

Die beiden Ausgangs-HIGHs sind real aufgelöst (Nachweis unter „Negativbefunde"
und „Auflösungsnachweis"), die Fixrunde selbst führt kein neues Finding ein.

## Negativbefunde

- geprüft, ohne Befund: **Auflösung Kotlin F-1** — `git show bd10c391` zeigt
  die Übersetzung beider Strings in
  `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientErrorMappingTest.kt`
  (Zeile 51: `"Rechtsklasse unzureichend für diese RPC"` →
  `"insufficient role for this rpc"`; Zeile 73: `"interner Fehler"` →
  `"internal error"`); `grep -n "Rechtsklasse\|interner Fehler" <Datei>` liefert
  0 Treffer (Exit 1). Die Datei (92 Zeilen, vollständig gelesen) ist damit
  durchgängig englisch — konsistent mit den vier übrigen Fixture-Strings und
  der KDoc.
- geprüft, ohne Befund: **Auflösung Python F-1** — `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md`
  real ausgeführt: **6 von 6 Zeilen OK, Exit 0** (Zeile 2 meldet
  `OK soll=15 ist=15`, gemessen am Arbeitsbaum `diff`); die korrigierte
  Plan-Zeile trägt `diff 15 -n -F 'PgChangeFeedAdministrationClient' --
  sdks/python`. Kein Rest-Vorkommen der alten Zahl (`soll=11`, `0 -> 11`) im
  Plan. Der „Zählwort"-Absatz („elf RPCs der Tabelle oben") misst ein
  anderes Symbol und stimmt unverändert (Zeilen 3–4: `soll=1 ist=1`,
  `soll=3 ist=3`) — kein „Nachzug widerspricht dem Nachbarn im selben
  Träger".
- geprüft, ohne Befund: **C#-Folge-Aktion (Commit `48e04899`)** — dieselben
  zwei Strings in
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/PgChangeFeedAdministrationClientErrorMappingTests.cs`
  (Zeile 44, 64) auf dieselbe Ziel-Form übersetzt wie Kotlin (wortgleich
  `"insufficient role for this rpc"` / `"internal error"`); `grep` nach den
  deutschen Fragmenten in `sdks/` liefert 0 Treffer. Der vom Kotlin-Review
  F-1 benannte Altbestand ist damit bereinigt — die dort für die Welle-Closure
  offene Anmerkung ist erledigt.
- geprüft, ohne Befund: **Korrektheit der übersetzten Fixtures** — beide
  geänderten Tests (Kotlin Zeile 54–56, 76–78; C# Zeile 46–48, 66–68)
  assertieren ausschließlich den Statuscode, nie den Beschreibungstext; die
  Aussage in der Commit-Message „Der Text wird in keinem Test assertiert" ist
  am Quelltext belegt. (Zum Abgleich: nur der `InvalidArgument`-Test
  assertiert `ex.message`/`ex.Message` — mit dem englischen String
  „source must not be empty", unverändert.) Die übersetzten Stubs verhalten
  sich testseitig identisch; dass der reale Server weiter die deutschen
  Fehlertexte sendet (`internal/adapters/driving/grpc/interceptor.go:147`,
  `administration.go:64`), ist die bewusste Betreiber-Sprache des Servers und
  widerspricht der Sprachreinheit der ausgelieferten SDK-Quelltexte nicht.
- geprüft, ohne Befund: **Commit-Umfang und Kennzeichnung** — `bd10c391`
  fasst genau zwei Dateien (Kotlin-Testdatei, Python-Slice-Plan), beide im
  Betreff und im Commit-Text je eigenen Review zugeordnet; `48e04899` fasst
  genau eine Datei (C#-Testdatei) und nennt ihre Herkunft (Folge-Aktion aus
  Kotlin-Review F-1). Kein fremder Inhalt unbenannt mitgeführt — der
  F-4-Verdikt des Kotlin-Reviews (unbenannte Fremdübernahme) greift hier
  nicht.
- geprüft, ohne Befund: Kommentar-/Anker-Form (`AGENTS.md` §3.7) — der Diff
  ändert keinen Kommentar, kein Docstring, keinen Anker; nur String-Literale
  und eine Zahl.
- geprüft, ohne Befund: Docker-only (`AGENTS.md` §3.1) — kein Skript, keine
  Umleitung, kein Host-Werkzeug im Diff; Build/Test liefen laut Bericht über
  `make sdk-pack-kotlin`/`make sdk-pack-csharp` (Exit 0) im gepinnten Image.
- geprüft, ohne Befund: Traceability — beide Betreffs nennen
  `LH-FA-SST-009`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — von der Fixrunde
  nicht berührt (keine Betreiber-Oberfläche geändert, Handbuch-Zug nicht
  ausgelöst).
- geprüft, ohne Befund: `make gates` real ausgeführt, Exit-Code direkt
  geprüft (nicht durch Pipe, `AGENTS.md` §3.9) — Exit 0 am finalen Baum
  dieses Reports (d-check: 0 Befunde, commit-traceability: OK,
  generated-sync: OK, a-check: 0 Befunde).

## Auflösungsnachweis der Ausgangs-Findings

| Ausgangs-Finding | Kategorie | Auflösung | Beleg |
|---|---|---|---|
| Kotlin F-1 (Form-Vorbild-Kopie, deutsche Fixtures) | HIGH | gelöst durch `bd10c391` | `git show bd10c391`; grep 0 Treffer; Datei durchgängig englisch |
| Python F-1 (Suchlauf-Zahl gegen die Messung) | HIGH | gelöst durch `bd10c391` | `make suchlauf-nachmessen` 6/6 OK, Exit 0; `soll=15 ist=15` |
| Kotlin F-1-Altbestand (dieselben Strings in C#) | (Hinweis) | gelöst durch `48e04899` | grep in `sdks/` 0 Treffer |
| Kotlin F-2/F-3 (LOW), F-4 (INFO) | LOW/INFO | stehen bewusst als Hinweise — laut Review keine Fixrunden-Pflicht | Verdikt des Kotlin-Reviews |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** — (keine; die Klassen der Ausgangs-Reviews
sind dort gezählt).

## Verdikt

**Merge-blockierend:** nein — die Fixrunde löst beide HIGH-Findings real und
nachprüfbar auf, ohne neue Befunde. Die Reviews beider Slices sind
abgeschlossen; `make suchlauf-nachmessen` läuft grün; die SDK-Paket-Bauten
lauten Exit 0.

Gemäß Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde (keine weitere
Fixrunde nötig) zieht dieser Report-Commit die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" in **beiden**
Slice-Plänen selbst auf `[x]` nach, mit Verweis auf die Report-Pfade. Die
übrigen offenen DoD-Zeilen des Python-Slices (Closure-Notiz,
Beobachtungs-Register, Risiko-Ausgänge) sind nicht Teil dieses Nachzugs —
sie bleiben bei der Slice-Closure.

**Übergabe:** Die Finding-Klassen der Ausgangs-Reviews gehen — wie dort
berichtet — in die Slice-Closure §7 der Welle und von dort in den
Steering-Loop-Zähler; dieser Report zählt keine neue Klasse. Der Report ist
ein Lauf-Beleg; er ersetzt keine Verifikation (Modul 11, Verifier-Aufgabe).
