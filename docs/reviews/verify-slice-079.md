# Verifikationsbericht: slice-079 — 2026-09-15

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-079` §2) und die bindenden Entscheidungen
([`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkte 1 und 2, §Fitness Function; [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a); `AGENTS.md` §3.6, §3.7, §3.9). **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe — die drei Review-Reports wurden als Kontext gelesen, nicht
als Beleg übernommen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan
(`in-progress/`, Stand `HEAD`), die beiden ADRs, die drei Review-Reports, den
Diff `35b9a00..HEAD` und die berührten Artefakte. **Alle** Zahlen dieses
Berichts stammen aus eigenen Läufen und eigener Profil-Deduplizierung; kein
Beleg des Implementers, des Reviewers oder des Planner-Nachzugs wurde
übernommen. Exit-Codes je in eigenem, ungepiptem Schritt gelesen
(`AGENTS.md` §3.9); Gate-Lauf und Folgehandlung getrennt beauftragt.

**Gegenstand:** `slice-079`. Der Anteil am Range sind `ca9864f` (der Schnitt),
`3827188` (DoD-Nachzug), `a34f7c6`/`c1a3353`/`b25775e` (die drei
Review-Reports), `0aa4d37` (Plan-Fix), `5d3ac6d`/`65aead2` (Fixrunden 1/3),
`9c92d90`/`9671b8f` (Register-Belege), `0dd531d` (Zählbasis-Nachzug,
Planner-Zug). Der reine `next → in-progress`-Move und die `open → next`-Commits
liegen **vor** `35b9a00` und sind nicht Gegenstand.

**Messbedingung.** Arbeitsbaum sauber (`git status --porcelain` leer) vor und
nach allen Läufen. Profil- und Rechenartefakte lagen außerhalb des
Arbeitsbaums (`/tmp`); der Vorher-Stand wurde **ohne** Baumänderung erzeugt —
durch einen zweiten `go test`-Lauf im bereits gebauten `coverage`-Image gegen
den **vollen** `-coverpkg` (dieselbe Toolchain, derselbe Quellstand, netzlos
`--network none`), nicht durch eine Arbeitsbaum-Kopie. Rücknahme entfällt
damit; es blieb kein Artefakt im Baum.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Zeile (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1a | `-coverpkg` **und** Testpaket-Liste führen die drei Pakete nicht mehr | **erfüllt** | `Dockerfile` (Stufe `coverage`): `pkgs=( $(go list ./internal/... ./cmd/... \| grep -vE '(^\|/)(postgresstorage\|postgresack\|replication/receive)$') )`; **beide** Verbraucher lesen dieselbe Variable — `-coverpkg="$(IFS=,; echo "${pkgs[*]}")"` und `"${pkgs[@]}"`. Selbst gemessen: die Testlauf-Ausgabe listet nur die **29** verbleibenden Pakete, `postgresstorage`/`postgresack`/`replication/receive` kommen weder als `ok`/`?` noch im `coverage: … of statements in …`-Kopf vor; `postgresstorage/mapper` **ist** enthalten. |
| 1b | Das Profil führt real **1679** Statements als Nenner (dedupliziert über die Block-Position) | **erfüllt** | Eigene Deduplizierung des gezogenen Profils (`unique Block-Position`, „gedeckt = mindestens ein Vorkommen `count > 0`"): **1104** Positionen, **1679** Statements, **1171** gedeckt = **69,7439 %**. Das rohe Profil trägt **26 598** Zeilen — jede Position **24,1×** (je Testbinary die volle `-coverpkg`-Menge); die Deduplizierung ist damit nicht optional. |
| 1c | Die drei Ausprägungen sind gegen die **Eigenschaft** geprüft, nicht abgeschrieben | **erfüllt, mit benannter Grenze** (V-2) | Eigene Vorher-Messung (voller `-coverpkg`): die drei tragen **788** Statements, davon nur **44** gedeckt — ihre Tests überspringen netzlos, das Paket erfüllt die Eigenschaft also **auf Paket-Ebene** (8 von 10 `postgresstorage`-Testdateien skippen). Genau diese Granularität hat `ADR-0071` Punkt 1 entschieden; die Aussage hält **auf Paket-Ebene**. Was fehlt, ist der *Beleg* der Teil-Erfüllung im Slice-Artefakt → V-2. |
| 2a | Der Effekt ist **beziffert** (vor/nach vergleichbar, die wegfallenden Statements benannt) | **erfüllt** | Vorher/Nachher in **derselben** Messung: **2467 → 1679** Statements, **1215 → 1171** gedeckt (49,25 % → 69,74 %); wegfallend **788** Statements, davon **44** gedeckt (`postgresack` 23/2, `replication/receive` 155/11, `postgresstorage` 610/31). **Entscheidender eigener Nachweis:** die Menge der Block-Positionen **außerhalb** der drei ist vor und nach dem Schnitt **byte-identisch** (1104 Positionen, 1679 Statements; `diff` leer) — es fiel ausschließlich Ausgeschlossenes weg. |
| 2b | Ist-Stand real gemessen, Stufe = abgerundete volle 5-%-Stufe, Endstufe bleibt 80 | **erfüllt** | Ist-Stand **69,74 %** (eigener Lauf) → `floor(69,74/5)*5 = 65`; `harness/mk/coverage.mk` liest `THRESHOLD ?= 65`; Endstufe **80** in `coverage.mk`-Kopf und Sensor-Doku §Kalibrierungs-Bindung unverändert. **Keine** Senkung: **65 > 40** (der Wert aus `slice-076`), Endstufe unverändert (§3.6 unberührt). |
| 2c | **Ein** beweglicher Ort; die übrigen nennen Rampe bzw. Verweis; Sensor-Doku nennt den **Messgegenstand** | **erfüllt** | Repo-weiter `grep` (ohne vendored Baseline): **genau eine** Zuweisung `THRESHOLD ?= 65` (`harness/mk/coverage.mk:15`). Die übrigen Träger nennen die Rampe („Einstieg 65 % → Endstufe 80 %" in `AGENTS.md` §4, `harness/README.md` §Sensors) oder verweisen (`coverage.mk`-Kopf, Sensor-Doku §Geltende Stufe). Die Sensor-Doku §Vertrag nennt die **Eigenschaft** plus Ausprägung und `mapper` ausdrücklich — nicht nur die Zahl. |
| 3a | Grün auf der neuen Stufe, Exit direkt gelesen | **erfüllt** | `make coverage-gate` → **Exit 0** (ungepiped), `coverage-gate: OK — Coverage 69.70% erfüllt Schwelle 65%`; Stage druckt `total: (statements) 69.7%`. |
| 3b | Rot unmittelbar über: `THRESHOLD=75` → Exit ≠ 0 | **erfüllt** | `make coverage-gate THRESHOLD=75` → **Exit 2** (`make`), Gate-Skript **Exit 1**; `coverage-gate: FAIL — Coverage 69.90% unter Schwelle 75%`. Abstand zum Ist-Stand ≈ 5,2 pp — außerhalb der dokumentierten Lauf-zu-Lauf-Schwankung (s. V-3). |
| 3c | `make gates` grün | **erfüllt** | `make gates` → **Exit 0**: baseline-verify `v6.5.0` OK · d-check **650** Dateien / **0** Befunde · commits-Modul `HEAD~5..HEAD` 0 Befunde · `commit-traceability: OK — 5 Commit(s)` · a-check 0 Befunde · coverage-gate OK 69,70 %. |
| 3d | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt (Häkchen-Nachzug offen)** | Drei Reports liegen real vor (`review-slice-079.md`, `-fixrunde.md`, `-fixrunde-2.md`), Verdikt je **0 HIGH / 0 MEDIUM**; das §2-Häkchen ist noch **nicht** gesetzt — Buchhaltungs-Nachzug der Closure, kein Substanz-Mangel. |
| 4 | Reconciliation-Register — **entfällt**, falls kein Inventur-Fund | **erfüllt (Entfall trägt)** | `docs/plan/planning/reconciliation.md` existiert real **nicht**; Repo durchgehend GF (`harness/conventions.md` §Modus-Deklaration `*`/`PGC`). |
| 5 | Beobachtungs-Register fortgeschrieben | **erfüllt** | `BEO-PGC/regel-weiter-als-ihr-sensor/evidence/slice-079.md` liegt real vor, `state.md` trägt **2×** (mit `slice-078`); kein Zähler-Feld gesetzt (abgeleitet). Vorbehalt: §8 des Plans hatte für diesen Eintrag „**kein** Treffer" notiert → V-4. |
| 6 | Jedes §6-Risiko trägt genau einen Ausgang | **korrekt offen** | Alle drei §6-Einträge tragen wörtlich `<bei Closure>`; keines vorzeitig geschlossen. Bewertung der Implementer-Empfehlungen in §5. |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 ist reiner Platzhalter (gelesen). Planner-Arbeit. |
| 8 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen, hier nicht zuständig** | Das Repo **hat** Wellen (`welle-20` ist offen, ihre Datei flach unter `docs/plan/planning/`), also trägt die nächste Wellen-Closure die Paarungen (§2 letzte Zeile) — nicht diese Slice-Closure. |

**Ergebnis §1:** Die **acht** gesetzten Zeilen (1a–3d) sind real erfüllt; die
Belege 1a–3c in dieser Sitzung **selbst** gefahren, die Nenner- und
Bezifferungs-Aussage über eine **eigene** Deduplizierung nachgerechnet. Die
drei Planner-Posten (6, 7, 8) sind korrekt offen. Die Entfall-Zeile (4) trägt.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make coverage-gate` (Default 65) | **0** | `OK — Coverage 69.70% erfüllt Schwelle 65%` — Grün-Beleg der neuen Stufe |
| `make coverage-gate THRESHOLD=75` | **2** (`make`), Skript **1** | `FAIL — Coverage 69.90% unter Schwelle 75%` — Rot-Beleg |
| `make gates` | **0** | alle fünf inneren Gates grün, Nachweis-Stempel zuletzt |
| `make doc-commits RANGE=35b9a00..HEAD` | **0** | 650 Dateien / 0 Befunde — jeder Commit des Fensters kennungstragend |
| `make doc-immutable RANGE=35b9a00..HEAD` | **0** | 650 Dateien / 0 Befunde — keine `Accepted`-ADR überschrieben |
| Profil ziehen (`docker run … cat /out/coverage.out`) | 0 | 26 598 Profilzeilen, 1104 unique Positionen (aus dem **eigenen** Grün-Lauf) |
| Vorher-Profil (`go test -coverpkg=./internal/...,./cmd/...` im selben Image, `--network none`) | 0 | 2467 Statements / 1215 gedeckt / 49,3 % — der Stand **vor** dem Schnitt, ohne Baumänderung |

> **Aufrufform-Hinweis (kein Befund):** `make doc-immutable` **ohne**
> `RANGE=` endet Exit 2 (`d-check: error: flag needs an argument: --range`) —
> die Range-Quelle ist in diesem Target nicht vorbelegt. Mit `RANGE=` grün.

## 3. Die Zahlbasis — der wunde Punkt dieses Slice (abschließend)

**(a) Reproduziert die Sensor-Doku [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)?**
**Ja.** Die Sensor-Doku §Zählbasis sagt „1679 Statements, davon 1171 gedeckt =
69,74 %". Meine eigene Deduplizierung desselben Gegenstands ergibt **1171 /
1679 = 69,7439 %**. Der Vorher-Stand ist ebenso reproduziert: **1215 / 2467 =
49,2501 %** (`ADR-0071` §Kontext: „1215 (49,25 %)"). Die drei Paket-Zahlen
stimmen zeichengenau: `postgresack` **23/2**, `replication/receive` **155/11**,
`postgresstorage` **610/31**, Summe **788/44** — und `postgresstorage/mapper`
**15/12** bleibt im Gegenstand.

> **V-1 (Beobachtung, Entscheidungs-Ebene):** [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
> §Entscheidung Punkt 2 sagt „Dieselbe gedeckte Menge (**1215** Statements)
> über den kleineren Nenner ergibt **69,74 %**". 1215/1679 = **72,36 %**; nur
> **1171**/1679 = 69,74 % — und 1215 − 44 = 1171. Die Klammer nennt die
> Vorher-Zahl, wo die Nachher-Zahl die Rechnung trägt. Die **Sensor-Doku**
> führt den richtigen Wert (1171); der Fehler liegt nur in der `Accepted`-ADR
> (immutabel, §3.5) und ist **kein** DoD-Punkt dieses Slice.

**(b) Ist die gedruckte Prozentzeile keine eigene Größe, sondern dieselbe
Messung in anderer Ausgabepräzision?** **Ja** — der Mechanismus ist von mir
selbst nachgefahren:

| Rechnung am **selben** Profil | Ergebnis |
|---|---|
| je Profilzeile naiv summiert (Duplikate als eigene Blöcke) | 1599 / 40 431 = **3,95 %** |
| dedupliziert über die Block-Position | 1171 / 1679 = **69,74 %** |
| `go tool cover -func` auf dem **Original**-Profil | **69,7 %** |
| `go tool cover -func` auf einem auf die Position deduplizierten Profil | **69,7 %** — identisch |

Die gedruckte Zahl ruht also auf **genau** der deduplizierten Basis: `go tool
cover` behandelt identische Block-Positionen als **einen** Block (die Zählungen
werden zusammengeführt), „gedeckt = mindestens ein Vorkommen `count > 0`" ist
dazu äquivalent, und das Verhältnis ist damit **invariant gegen die
Duplikate**. Die zwei genannten Figuren `69,7 %` (Tool, eine Nachkommastelle)
und `69,70 %` (`tools/coverage-gate.sh`, `%.2f`) sind **eine** Messung. Die
korrigierte §Zählbasis-Aussage (Commit `0dd531d`) ist damit **gemessen richtig**;
die zuvor dort erklärte Mechanik („jede Block-Position so oft, wie sie
vorkommt") hätte 3,95 % ergeben und war falsch.

**(c) Sind §Grenze Punkt 4 (68,92 % · 52,51 % · 64,45 %) aus der genannten
Basis herleitbar?** **Ja**, jede Rückrechnung aus `1171 / 1679` plus den
Paket-Zahlen:

| Rücknahme | Rechnung | Ergebnis | Abstand zu 65 |
|---|---|---|---|
| `postgresack` zurück | (1171 + 2) / (1679 + 23) = 1173 / 1702 | **68,92 %** | +3,9 pp (grün) |
| `postgresstorage` zurück | (1171 + 31) / (1679 + 610) = 1202 / 2289 | **52,51 %** | −12,5 pp (rot) |
| `replication/receive` zurück | (1171 + 11) / (1679 + 155) = 1182 / 1834 | **64,45 %** | −0,6 pp (rot) |

Alle drei Werte und alle drei Abstände sind mit meinen **eigenen** Paket-Zahlen
exakt reproduziert; die Aussage „nur die Prozent-Schwelle trägt die
Gegenstands-Hälfte nicht vollständig gewächtert" trägt.

## 4. Entscheidungs- und Spec-Konformität

1. **`ADR-0071` Punkt 1 ist umgesetzt** — `-coverpkg` und Testpaket-Liste ohne
   die drei namentlich genannten Pakete, `mapper` im Gegenstand, Nenner 1679.
2. **`ADR-0071` Punkt 2 ist umgesetzt** — Endstufe **80 %** unverändert, Rampe
   nach dem unveränderten Mechanismus (`ADR-0054` §(a): Ist-Stand, abgerundet
   auf die volle 5-%-Stufe) auf **65** gesetzt; der Wert steigt real
   **49,25 % → 69,74 %**.
3. **Keine Schwellen-Senkung** (`AGENTS.md` §3.6 unberührt): 65 > 40; Endstufe
   unverändert; die Senkungs-Grenze („nur per ADR") steht im `coverage.mk`-Kopf.
4. **`spec/**` unberührt**; die Messmethode des Lastenhefts bleibt unangetastet
   (Slice-Plan §Berührte Spec-Stellen: `—`).
5. **§1-Abgrenzungen halten am Diff.** Die Dateiliste `35b9a00..HEAD` umfasst
   **kein** `spec/**`, **kein** `tools/coverage-gate.sh`, **kein** `Makefile`,
   **kein** `internal/**`, **kein** `test/**` — genau die im Plan §1
   ausgeschlossenen Gegenstände. Der Gate-Mechanismus (Skript, Stage-Reihenfolge)
   ist unverändert; geändert ist ausschließlich der **Umfang**.
6. **Das Gate-Skript-Verhalten ist unberührt** — `tools/coverage-gate.sh` ist
   nicht im Diff; seine Ausgabeform und Exit-Codes (0/1/2) sind unverändert
   (Rot-/Grün-Lauf belegen 1 bzw. 0).

## 5. Plan-vs-Code-Diff (beide Richtungen)

**Plan → Code (jede Behauptung geprüft):**

| Plan-Behauptung | Befund |
|---|---|
| §3 `Dockerfile` (Stufe `coverage`): Umfang der Messung geschnitten | **trägt** — der Schnitt steht in der Stufe, nicht im Gate-Skript |
| §3 `harness/mk/coverage.mk`: neu kalibrierte `THRESHOLD`, Kopf auf Ist-Zustand, **ein** beweglicher Ort | **trägt** — `THRESHOLD ?= 65`, einziger Wert-Träger (repo-weiter `grep`) |
| §3 `harness/sensors/coverage-gate.md`: **Messgegenstand** statt nur Zahl, neue Beleg-Zeile, benannte Folge | **trägt** — §Vertrag nennt Eigenschaft + Ausprägung + `mapper`; §Grenze Punkt 3 verweist auf `ADR-0071` Punkt 3 |
| §3 `harness/README.md` §Sensors: Umfang + geltende Stufe | **trägt** — nennt Umfang und Rampe, verweist für den Wert auf `coverage.mk` |
| §3 `AGENTS.md` §4: update, **nur falls** die Zeile den Umfang nennt | **trägt** — die Zeile nannte den Umfang, also nachgezogen; Rampe 65 → 80, kein beweglicher Wert |
| §1 „die Endstufe 80 % steht" | **trägt** |
| §1 „Die Tests selbst — dieser Slice fügt keinen Test hinzu" | **trägt** — `test/**` nicht im Diff |
| §6 Risiko 1: „rund 44 gedeckte Statements fallen mit" | **trägt exakt** — 44 (31 + 2 + 11) |

**Code → Plan (Zustände, die der Plan nicht ausspricht):**

- **Eine abgeleitete Liste, zwei Verbraucher.** Plan §3 sagt „der `-coverpkg`-Umfang
  und die Testpaket-Liste"; der Code zieht beide aus **einer** Variable
  (`go list` + Filter) — ein künftiges Paket wandert in beide zugleich. Der Plan
  nennt diese Kopplung nicht.
- **Das Filter-Muster ist suffix-verankert** (`(^|/)(postgresstorage|postgresack|replication/receive)$`),
  nicht baum-anfangs-verankert: ein gleichnamiges Paket **irgendwo** im Baum
  fiele still mit heraus (heute existiert keines — die Differenz `go list` 32 →
  29 ist genau die drei). Der Plan spricht nur von „den drei namentlich".
- **Die Sensor-Doku ist über den Plan §3 hinaus um eine §Zählbasis-Sektion
  gewachsen** — begründet durch die drei Review-Runden (fünf verschiedene
  „gedeckt"-Werte), aber im Plan §3 nicht angelegt.
- **Der Register-Beleg steht unter einem Eintrag, den §8 als „kein Treffer"
  geführt hatte** → V-4.

## 6. Urteil zur Planner-Korrektur `0dd531d` (offengelegter Vorgang)

Der letzte Commit korrigiert in `harness/sensors/coverage-gate.md` die
§Zählbasis-Mechanik und zieht die Register-Zahlen nach (`52,3 % → 52,51 %`,
`64,2 % → 64,45 %`). **Urteil: inhaltlich tragend, als Closure-seitiger
Wortlaut-Nachzug vertretbar.**

- **Inhaltlich richtig** — ich habe den Mechanismus **selbst** nachgefahren
  (§3(b)): Original- und dedupliziertes Profil liefern **denselben** `go tool
  cover -func`-Wert; die per-Zeile-Summe (3,95 %) ist nicht das, was das
  Werkzeug tut. Die neue Aussage („dieselbe Messung in anderer
  Ausgabepräzision") ist damit gemessen wahr, die alte war falsch.
- **Kein Verhaltens-Change, Zahlen unverändert** — der Commit berührt **zwei
  Dokumentationsdateien**; `THRESHOLD`, Profil und Gate-Skript sind unberührt,
  die operativen Zahlen (Einstieg 65, Endstufe 80, 1679, 1171) bleiben. Der
  Nachzug hat die **Abnahme nicht verändert**: Grün- und Rot-Beleg dieses
  Berichts sind mit dem Commit gefahren.
- **Vertretbar als Closure-seitiger Zug, mit einem Vorbehalt.** Die Korrektur
  landet in einem **Implementer-Artefakt** (`harness/sensors/coverage-gate.md`)
  **nach** dem letzten Review-Lauf — der Bestätigungslauf hatte sie als
  **nicht blockierend** bewertet. Formal liegt damit die korrigierte Bullet
  außerhalb jedes Review-Fensters; der einzige unabhängige Beleg für sie ist
  **dieser** Verifikations-Lauf (gemessen), was das Vorgehen trägt. Ich halte
  die Sache für **richtig**, nicht für einen V-Befund.

## 7. Offene Closure-Obliegenheiten (benannt, nicht ausgeführt)

Planner-Arbeit nach diesem Bericht:

1. **§6-Risiko-Ausgänge (drei)** — je **genau einer** aus der geschlossenen
   Menge *eingetreten / entfallen / weiter offen*; „nicht eingetreten" ist
   **kein** Mengen-Element und beim Nachzug auf `entfallen` zu ziehen.
   Bewertung der Implementer-Empfehlungen:
   - **Risiko 1** („der neue Nenner nimmt gedeckte Statements mit") — die
     Wirkung ist **eingetreten und beziffert**: 44 gedeckte Statements fielen
     mit (1215 → 1171). Die **Empfehlung „eingetreten/wie geplant" ist in
     dieser Form nicht formrein**: `Modul 5` verlangt bei *eingetreten* einen
     **Carveout oder eine Folge-Slice-ID** — beides existiert nicht (die
     vorgeschriebene Antwort *Bezifferung* ist geliefert, es bleibt nichts
     auszulagern). Empfehlung: **`entfallen` mit Begründung** (die Bezifferung
     ist gebaut, die befürchtete Überzeichnung ausgeschlossen). Der Rest — 31
     netzlos gedeckte Statements in `postgresstorage` — ist die Substanz von
     V-2.
   - **Risiko 2** („die drei Pakete stehen danach ohne jede Zahl da") —
     **Empfehlung `entfallen` trägt**: meine Vorher-Messung zeigt, dass das
     Gate sie **nie gemessen** hat (788 Statements, davon 44 gedeckt = 5,6 %;
     die Tests skippen netzlos). Der Schnitt entfernt ein *Mit-Zählen*, keine
     Messung. Der Rest — `ADR-0071` Punkt 3 — ist eine **eigene Folgepflicht
     ohne Adresse** (s. 6.).
   - **Risiko 3** („die neue Stufe könnte unter der bisherigen liegen") —
     **Empfehlung `nicht eingetreten` trägt inhaltlich**: 65 > 40, Endstufe
     unverändert, keine Senkung. Nur die **Form** ist auf `entfallen` zu
     ziehen.
2. **Register-Beleg** — liegt vor (`evidence/slice-079.md`, Zähler **2×**,
   weiter unter der Schwelle, kein Ausgang fällig). **Achtung V-4:** §8 des
   Plans notierte für denselben Eintrag „kein Treffer"; §7 sollte die
   Identitäts-Erweiterung benennen (die Beobachtung ist als **Klasse** geführt:
   „Zusage weiter als ihr Sensor"), sonst steht §8 gegen den Beleg.
3. **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag und
   Register-Zeile; die wiederkehrenden Finding-Klassen der drei Review-Reports
   sind dort der Zähler-Übergang.
4. **`git mv` nach `done/`** — erst Inhalt (Häkchen, Closure-Notiz), dann der
   reine Move (`AGENTS.md` §3.3); das §2-Häkchen „Review durchgeführt" ist
   dabei nachzuziehen.
5. **Drei Paarungen** — von der **Wellen-Closure** getragen (das Repo hat
   Wellen; §2 des Plans), nicht von dieser Slice-Closure.
6. **Zwei Folgepflichten ohne Adresse:** `ADR-0071` Punkt 3 (DB-Adapter-Coverage)
   und Punkt 5 (Executor-Naht) sind eigene Vorgänge; `open/` und `next/` sind
   real **leer**, es existiert also **keine** Slice-Datei. Benannt als
   „verwandt, aber nicht Teil" in `welle-20` §5 — der Adress-Punkt bleibt offen
   (Klasse `BEO-PGC/aufschub-adresse-verfaellt`).
7. **Besonderheit dieses Slice — der Wellen-Start.** `welle-20` §2 führt
   „`slice-079` liegt in `done/`" als **Start-Trigger**, noch **nicht erfüllt**.
   Mit dem `git mv` in `done/` ist der Trigger erfüllt; dann kann der **erste
   Test-Slice** nach dem Größenmaß aus `welle-20` §4 geschnitten werden
   (`internal/bootstrap` 591/327 ungedeckt, `cmd/pg-change-feed` 49/49,
   Rest-Tail ~130). Bis dahin misst die Welle korrekt gegen den neuen Nenner.

## 8. V-Befunde (Wortlaut-/Entscheidungs-Ebene, **keine** DoD-Zeile verletzend)

- **V-1 — `ADR-0071` §Entscheidung Punkt 2 nennt 1215, wo 1171 die Rechnung
  trägt** (`1215/1679 = 72,36 %` ≠ 69,74 %). Entscheidungs-Ebene, nicht
  Slice-DoD; die ADR ist `Accepted` (immutabel) → Planner-Notiz oder, falls
  sie als Irreführung empfunden wird, ein Folge-ADR-Weg (`Supersedes`). Die
  Sensor-Doku führt den richtigen Wert.
- **V-2 — die Eigenschaft ist gröber als ihre Ausprägung, und die Teil-Erfüllung
  steht in keinem Slice-Artefakt.** `postgresstorage` erfüllt die Eigenschaft
  nur **teilweise** (8 von 10 Testdateien skippen, 2 nicht; **31/610 netzlos
  gedeckt**), `internal/bootstrap` (44,3 %) und `natsnotify` (84,8 %) führen
  ebenfalls skippende Testdateien und bleiben im Gegenstand. Plan §4 hatte
  **genau diesen Fall** als `in-progress → open`-Bedingung vorab benannt; sie
  wurde nicht ausgelöst. Sie ist auch **nicht** zu erzwingen: `ADR-0071`
  Punkt 1 hat die **Paket-Granularität** entschieden (bindend), eine Rücknahme
  wäre eine Änderung dieser Entscheidung (§3.5). Was fehlt, ist allein die
  *Benennung*: §Grenze Punkt 1 sagt „setzen einen externen Dienst voraus" ohne
  die 5,1 %-Nuance. Empfehlung: als §Grenze-Punkt oder Closure-Notiz nachziehen.
- **V-3 — die Beleg-Zeile der Sensor-Doku fixiert einen schwankenden Wert.**
  Der Rot-Beleg (`THRESHOLD=75`) ist dort mit `69.70 %` zitiert; mein eigener
  Rot-Lauf druckte `69.90 %`, mein Grün-Lauf `69.70 %` — die dokumentierte
  Lauf-zu-Lauf-Schwankung (§Grenze Punkt 4, wenige Statements) wirkt auf beide
  Rollen. Nicht DoD-relevant (Exit ≠ 0 hält bei beiden Werten); die Zeile
  könnte die Schwankung statt eines Festwerts nennen.
- **V-4 — §8 Sichtung gegen den Register-Beleg.** §8 sagt „`kein` Treffer — die
  host-lokale Pfad-Regel, nicht die Coverage", während die Closure
  `evidence/slice-079.md` **unter genau diesem Eintrag** ablegt (→ 2×). Beides
  kann richtig sein (die Beobachtung ist eine **Klasse**), aber die beiden
  Stellen müssen in §7 versöhnt werden.

## 9. Negativbefunde

- geprüft, ohne Befund: `Dockerfile` (nur Stufe `coverage`) — eine Variable,
  zwei Verbraucher, `mapper`/`queries` im Gegenstand; `SHELL … pipefail`
  unverändert; der Kommentar trägt Ist-Zustand und Kopplung, keine
  Slice-/Wellen-Chronik
- geprüft, ohne Befund: `harness/mk/coverage.mk` — `THRESHOLD ?= 65` als
  **einziger** Wert-Träger (repo-weiter `grep`), Kommentar-Klassen Zusage/
  Kopplung/Grenze, `## help` nennt nur die Rampe
- geprüft, ohne Befund: `harness/sensors/coverage-gate.md` §Vertrag und
  §Kalibrierungs-Bindung — Messgegenstand benannt, Endstufe 80 fest, Verweis
  auf den beweglichen Ort; Exit-Tabelle unverändert (Skript-Exits)
- geprüft, ohne Befund: `harness/README.md` §Sensors und `AGENTS.md` §4 —
  Umfang plus Rampe, kein beweglicher Wert
- geprüft, ohne Befund: `docs/plan/adr/**` — im Range **unverändert**,
  `ADR-0071` bleibt `Accepted`, kein `Supersedes`
- geprüft, ohne Befund: §1-Abgrenzungen — **kein** `spec/**`, `internal/**`,
  `test/**`, `tools/**`, `Makefile`
- geprüft, ohne Befund: Commit-Fenster — jeder Betreff nennt `ADR-0071` oder
  `LH-FA-CFG-005`, **keine** `SPEC-*`/`ARC-*` im Betreff, **kein**
  Attributions-Trailer (`make doc-commits` Exit 0)
- geprüft, ohne Befund: Mess-Rückstände — Profil- und Rechenschritte lagen
  außerhalb des Baums, der Vorher-Lauf lief im bestehenden Image; Baum nach
  allen Läufen sauber
- geprüft, ohne Befund: die Register-Form — `evidence/slice-079.md` ist
  formgebunden (`<vorgangs-id>.md`), `state.md` trägt Zustand und auflösbaren
  Anker, kein Zähler-Feld

## Verdikt

**DoD-Konformität:** **bestätigt** — die acht gesetzten Zeilen (Liefer-Punkte
1–3) sind real erfüllt; Grün-/Rot-Beleg, `make gates` und die Traceability-
Sensoren in dieser Sitzung **selbst** gefahren, der Nenner **1679** und die
gedeckte Zahl **1171** über eine **eigene** Profil-Deduplizierung
nachgerechnet, die Statement-Menge außerhalb der drei **byte-identisch**
vor/nach dem Schnitt. Die drei Planner-Posten (§6-Ausgänge, §7-Notiz, drei
Paarungen) sind korrekt offen, das Reconciliation-Item entfällt real.

**Zahlbasis:** **abschließend aufgeklärt** — (a) 1171/1679 = 69,74 %
reproduziert; (b) die gedruckte Prozentzeile ist **keine** eigene Größe
(Original- und dedupliziertes Profil ergeben **denselben** `go tool cover`-Wert;
die 3,95 % entstehen nur bei falscher per-Zeile-Summierung); (c) 68,92 % /
52,51 % / 64,45 % sind aus `1171/1679` plus den Paket-Zahlen **exakt**
herleitbar. Der wunde Punkt des Slice ist damit geschlossen; die vier
V-Befunde liegen auf **Wortlaut-/Entscheidungs-Ebene** und setzen **keine**
DoD-Zeile rot.

**Entscheidungs-/Spec-Konformität:** **bestätigt** — `ADR-0071` Punkte 1 und 2
umgesetzt, Endstufe **80** unverändert, **keine** Senkung (65 > 40),
`spec/**` und der Gate-Mechanismus unberührt, alle §1-Abgrenzungen halten am
Diff.

**Plan-vs-Code-Diff:** in der Hauptrichtung deckungsgleich; in der
Gegenrichtung drei ungenannte Zustände (eine abgeleitete Liste mit zwei
Verbrauchern, das suffix-verankerte Filter-Muster, die §Zählbasis-Sektion) —
ohne DoD-Wirkung.

**Kein Merge-Blocker aus Verifikations-Sicht.** Die Planner-Korrektur
`0dd531d` trägt; der `git mv` nach `done/`, die Closure-Notiz, die drei
Risiko-Ausgänge und die drei Paarungen bleiben Planner-Arbeit. Mit dem
`done/` ist der **Start-Trigger von `welle-20`** erfüllt.
