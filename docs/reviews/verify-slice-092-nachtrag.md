# Verifikations-Nachtrag: slice-092 — Abschluss-Prüfung nach dritter Runde und Closure · 2026-09-16

**Rolle:** Verifier (Modul 11), **frischer Kontext** für **einen** Durchgang.
Nachtrag zu [`verify-slice-092.md`](verify-slice-092.md) (Stand `cc37ddc`,
HEAD war `ad95444`). Auftrag: die dritte Runde (`1d0d11e`), der Plan, die
Closure-Inhalte (`6371a28`) und ein letztes Verdikt. **Nicht** gegen realen
Bedarf (Validator, nicht ausgelöst).

**Stand:** `HEAD` = `6371a28`, Zweig `main`, Baum sauber
(`git status --porcelain` leer — vor und nach jedem Lauf). Neue Commits seit
§9 des Hauptberichts: `cc37ddc` (dieser Bericht), `a81a6ea` (V-1: die
abgeleitete Zahl trägt ihren Ursprung), `da0ffe3` (**Delta-Review**
[`review-slice-092-delta.md`](review-slice-092-delta.md)), `1d0d11e` (dritte
Runde, D-1…D-5), `6371a28` (Closure: §6, §7, Häkchen, Register). Umfang seit
dem Test-Commit: **10** Pfade, **alle Markdown**; Gesamtvorgang **24** Pfade
(**13** `*_test.go`, **11** Markdown) — **kein** Produktcode, `THRESHOLD`
unberührt (L-1.2, L-1.12). Exit-Codes sind je **ungepiped** und in eigenem
Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und Auswertung waren getrennt
beauftragt. Logs und Arbeitsbaum-Kopien liegen **außerhalb** des Baums.

---

## L-1 Eigene Messungen dieses Nachtrags

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| L-1.1 | `make gates` am Stand `6371a28` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | sechs Checks: `baseline-verify v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 78.60% erfüllt Schwelle 70%` · `d-check: 762 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)`; danach Status leer |
| L-1.2 | Umfang: `git diff --name-status 3de9547..HEAD` · `git diff --name-only 3de9547..HEAD \| grep -v '\.md$'` | **0** / **1** | **10** Pfade seit dem Test-Commit, alle Markdown; die zweite Abfrage **leer** — kein Code geändert, die Messungen aus §1 des Hauptberichts gelten für diesen Stand |
| L-1.3 | Profil am **Band-Stand unten** `3d5b60d` (= `slice-085`), netzlos, Zählbasis der `coverage`-Stufe | **0** | Nenner **1903** · gedeckt **1369** · ungedeckt **534**; gedruckt `71.9%` |
| L-1.4 | Profil am **Band-Stand oben** `f90c3f4` (= `slice-091`), netzlos | **0** | Nenner **1903** · gedeckt **1469** · ungedeckt **434**; gedruckt `77.2%` |
| L-1.5 | `git diff --name-only 3d5b60d..f90c3f4 -- internal cmd` (Produktcode-Claim der dritten Runde) | **0** | ohne Testdateien und ohne `replication/receive`: **leer**; über `7835b6e..f90c3f4` genau die **drei** Dateien in `replication/receive` — die ausgenommene Gruppe, wie der Satz sie nennt |
| L-1.6 | Paarungen der zwei Bänder (**abgeleitet**) | — | `1468 − 1371 = 97` · `1468 − 1369 = 99` · `1471 − 1371 = 100` · `1471 − 1369 = 102` → **97…102** |
| L-1.7 | Plan am Stand `HEAD`: Platzhalter/Häkchen · Zähl-Wort · die vier Zahlen | **1** (leer) | **keine** `<…>` und **keine** leeren `[ ]` mehr; Zähl-Wort `**Vier Zahlen, vier Dinge**`; Platzierungen: **13** (Glob, §8), **10** (§1), **12** (§1, §2/LP1, §3), **11** (§1, §5, §6, §8) — je an der richtigen Aussage |
| L-1.8 | Register: `ls evidence/ \| wc -l` je Eintrag | **0** | `zahl-in-traeger-driftet-gegen-die-messung` **7** · `negativtest-ohne-bindung-an-seine-eingabe` **6** · `arbeit-ueberholt-stehenden-traeger` **1**; Klassen mit ≥ 5×: **7 · 6 · 6 · 6 · 5 · 5** |
| L-1.9 | **P-9/P-10** — `includecolumn`: feste Spalte in der Weitergabe, **Parent** vs. **HEAD** | **0** / **1** | Parent **grün** (`ok`, EXIT 0), HEAD **rot** (`Fehler = <nil>, wollen ErrSourceColumnMissing`); Kontrolle am unmutierten HEAD **Exit 0** — die zweite der zwei reparierten Spaltenbindungen ist damit **gemessen** (für `excludecolumn` schon §6/P-6, P-7) |
| L-1.10 | Träger-Zählung des Glob-Fehlers: naiver `grep -l "10 Pakete\|zehn Pakete"` über `*.md` am Parent `2b7c6ec`, dann je Datei mit Kontext | **0** | **4** Dateien — aber nur **zwei** tragen die Fehlzählung (Slice-Plan: **5** Stellen `:49`, `:99`, `:140`, `:186`, `:210`; `welle-20.md`: **1** Stelle `:100`); die zwei anderen (`verify-slice-005`, `verify-slice-007`) sagen „**10 Pakete `ok`**" über einen **Testlauf**, nicht über den D1-Umfang |
| L-1.11 | §7 auf Lauf/Größe: `grep -E "Lauf …\|am Stand …"` und `grep -E "1493\|1495\|78,5\|78,6"` über den **ganzen** Plan | **1** / **1** | **kein** Lauf-Marker, **keine** erreichte Zahl — im ganzen Plan nicht |
| L-1.12 | Gesamtumfang: `git diff --name-only 3de9547^..HEAD` | **0** | **24** Pfade: **13** `*_test.go` · **11** Markdown (2 Träger, 3 Berichte, 5 Register-Dateien); kein Nicht-Markdown außerhalb der Tests |

---

## L-2 Die dritte Runde — trägt die neue Fassung in `harness/sensors/coverage-gate.md`?

**Drei Fragen, drei gemessene Antworten — alle drei: ja.**

1. **Herkunft richtig gesetzt?** Ja. Der Satz nennt die vier gedeckten Zahlen
   als die **gemessenen** („**1369** von **1903**", „**1371**", „**1468**",
   „**1471**", jede mit ihrem Lauf-Marker) und die **Differenz** ausdrücklich als
   „**abgeleitet**, nicht gemessen" — genau die Umkehrung, die der Delta-Review
   als D-2 verlangt hat. §3.12 Instanz A ist erfüllt: gemessen = Zahl,
   abgeleitet = Differenz.
2. **Satz geschlossen?** Ja. Der Hauptsatz trägt sein Verb im eigenen Satz
   („Über **acht** Läufe desselben Produktionsstands (`go test …`, Auswertung
   über die Block-Position, Lauf `slice-091`) **lag** die gedeckte Zahl zwischen
   **1468** und **1471**"); die Beleg-Klammer steht bei ihrem Satz, die
   Bänder-Aussage ist ein **eigener Absatz** — D-3 ist behoben.
3. **Steht „zwei **Test**-Stände **desselben** Produktionsstands" gegen die
   Messung?** Nein — die Aussage **hält**, und meine eigenen Profile stützen sie
   an beiden Enden: der Nenner ist an beiden Band-Ständen **1903** (L-1.3,
   L-1.4), und zwischen ihnen hat sich im Messgegenstand **kein** Produktcode
   bewegt (L-1.5). Die Verschiebung um ~100 liegt damit in den **gedeckten**
   Zahlen, also am **Test**-Baum. Die dritte Runde hat meinen V-1-Befund zur
   Sache (der Satz band die Differenz an „zwei verschiedene Code-Stände")
   **korrigiert**, nicht wiederholt.

**Und: hat sie etwas Neues falsch gemacht?** In `coverage-gate.md` **nichts** —
die vier Zahlen des neuen Absatzes sind einzeln nachgemessen (L-1.3, L-1.4), die
`git diff`-Klausel mit **zwei** Ranges geprüft (L-1.5), die Paarungen
nachgerechnet (L-1.6), die Etiketten gegen §3.12 gehalten, der Satz geschlossen.
**Der neue Fehler dieses Zugs liegt auf der Plan-Seite** und ist nicht von ihr
erzeugt, sondern von der **Closure** danach sichtbar gemacht worden: die
Rang-Klausel, die dieselbe Runde als D-5-Korrektur geschrieben hat, ist nach der
Closure stale (**N-2**, LOW). Der Zug, der den Träger angefasst hat, hat damit
keinen neuen Zähl-, Herkunfts- oder Satzfehler erzeugt — die Reihe (Glob-Fehler
→ vier Zahlen → Zähl-Wort → invertiertes Etikett) endet hier.

---

## L-3 Der Plan — die vier Zahlen und das Zähl-Wort

**Beides hält.** Der Kopf der Passage sagt jetzt `**Vier Zahlen, vier Dinge**`
und führt darunter **vier** Größen; die Platzierungen sind je Aussage richtig
(L-1.7): **13** führt der Glob, **10** trugen ungedeckte Statements (Planungs-
Stand), **12** haben Tests bekommen, **11** haben **bewegte** Deckung — die drei
Bewegungs-Aussagen in §5/§6/§8 stehen sämtlich auf **11**. Der Widerspruch aus
F-1 (13 gegen 11) und der Zählfehler aus D-1 („Drei" über vier) sind **weg**, und
die dritte Runde hat an dieser Stelle nichts Neues erzeugt.

---

## L-4 Die Closure-Inhalte

| Gegenstand | Befund |
|---|---|
| **R1** (die `22`/`2` seien eine Über-Schätzung) → *entfallen* | **wahr** — beide Zahlen hielten exakt: **22** in den Use-Cases, **2** in `domain/model`, **+24** als Block-Positionen, **Netto +22** durch den einen Flake-Block (§1 Nr. 8 und Nr. 10 des Hauptberichts) |
| **R2** (Negativtest bindet an den Fake) → *eingetreten — und behoben* | **wahr, beide Hälften jetzt gemessen** — `excludecolumn` Parent **grün** / HEAD **rot** (§6 P-6, P-7) und `includecolumn` Parent **grün** / HEAD **rot** (L-1.9). Die Aussage „die Mutation färbt sie jetzt rot, am Parent blieben sie grün (nachgemessen)" trägt wörtlich |
| **R3** (Coverage-Theater) → *entfallen* | **wahr in der Rechnung** — **21** + **32** + **6** = **59** Proben; die **6** des Verifiers sind meine eigenen an **neuen** Tests (richtig abgegrenzt: P-1…P-6; P-7 und P-8 liegen außerhalb dieser Zahl). **Grenze:** „alle **24** neu gedeckten Block-Positionen sind über je eine rot färbende Zusage gebunden" ruht auf den zwei Berichten (21 + 32); ich habe **6** nachgefahren und kein Gegenbeispiel gefunden |
| **R4** (überholter Träger) → *eingetreten — und begrenzt* | **in der Richtung wahr, in der Zahl nicht** — die Prüfung hat getragen (der Glob-Fehler wurde gefunden und berichtigt, und der Sensor-Träger wird **nicht** zahl-falsch: §5 des Hauptberichts hält). Die Zahl **„vier Träger"** ist contra-gemessen (L-1.10) → **N-3** |
| **Register-Zähler** | **stimmen** — `negativtest` **6×**, `zahl-in-traeger` **7×**, `arbeit-ueberholt` **1×** (L-1.8), deckungsgleich mit den `state.md`-Zeilen und mit den zwei neuen `evidence/slice-092.md`. Die Ablehnung des `arbeit-ueberholt`-Kandidaten steht in dessen `state.md`, der Zähler bleibt **1×** — benannt und begründet, wie der Lese-Schritt es verlangt |
| **Häkchen, §6-Ausgänge, §7** | **vollständig** — **11 von 11** Häkchen gesetzt, **keine** `<…>`-Platzhalter mehr, vier Risiken mit je einem Ausgang (L-1.7) |
| **§7: Zahl ohne Lauf oder ohne Band?** | **Ja — ohne Lauf**, und die erreichte Quote fehlt ganz → **N-1** (L-1.11) |

---

## L-5 Findings des Nachtrags

### N-1 — §7 nennt den Zuwachs ohne seinen Lauf, und die erreichte Quote steht in keinem Träger; LP1 ist auf `[x]` gesetzt

- `kategorie`: **MEDIUM**
- `quelle`: §2/LP1 des Slice-Plans („**Der Zuwachs wird als Zahl mit ihrem Lauf
  genannt**, nicht als ‚deutlich besser'") · `AGENTS.md` §3.12 Instanz A ·
  Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**7×** — in
  diesem Vorgang bereits bedient)
- `pfad`: `docs/plan/planning/in-progress/slice-092-coverage-cluster-d1.md`
  §7 („Was hat funktioniert", „Was ging anders als geplant")
- `befund`: §7 nennt die Zahlen (Zuwachs **24** Statements, **22 + 2**, die
  Planungs-Größen **13**/**10**/**22**, die Proben **21**/**32**/**6**) — aber
  **kein Lauf-Marker** steht irgendwo im Plan (L-1.11: beide Greps leer), und die
  **erreichte** Quote (das gedruckte Band `78,5 %`/`78,6 %`; `make gates` druckt
  **78.60 %**, L-1.1) fehlt im ganzen Plan. Die Zahl, die §5 des Hauptberichts
  ausdrücklich für §7 vorgesehen hatte — die Lauf-Größe mit ihrem Lauf —, hat
  damit keinen Träger; sie lebt allein in den Commit-Messages (git-Historie,
  laut `ADR-0083` §Geltungsbereich **kein** Doku-Träger) und in den
  Lauf-Belegen. Die **Zustands**-Zahl (**24**, benannte Block-Positionen) ist
  richtig und steht in §7 — der Mangel ist der **Anker**, nicht der Wert.
- `verifizierbar`: ja — `grep -nE "Lauf \`?slice-092|am Stand" <Plan>` → kein
  Treffer; `grep -nE "1493|1495|78,5|78,6" <Plan>` → kein Treffer; L-1.1
- `urteil`: **Ein-Zeilen-Nachtrag — vor dem `git mv`**, denn die Plan-Datei
  wandert nach `done/` und trägt die Notiz mit. Präzedenz: `verify-slice-091`
  N-3 (dort LOW, weil die Zahlen selbst trugen und nur der Marker fehlte) — hier
  **MEDIUM**, weil die verlangende Klausel **im DoD-Kriterium selbst** steht und
  dieses Häkchen bereits gesetzt ist. Empfohlener Zusatz: „… der Zuwachs ist
  **24** Statements (22 + 2, im Lauf dieses Vorgangs, Block-Position gegen den
  Parent), Nenner **1903** unverändert; gedruckte Gate-Zeile: **78,60 %** (Band
  `78,5 %`–`78,6 %`; die Schwankung ist der eine Block
  `internal/bootstrap/wiring.go:991.5,992.13`)."

### N-2 — Die Rang-Klausel in §2/LP2 ist nach der Closure stale, weil die Closure denselben Zähler erhöht hat

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz A (bewegliche Zahl ohne Zeitpunkt) ·
  dieselbe Klasse wie N-1
- `pfad`: `slice-092-coverage-cluster-d1.md:116-118` (§2/LP2) gegen `:297` (§7)
  und `:352` (§8)
- `befund`: §2/LP2 sagt „`negativtest-ohne-bindung-an-seine-eingabe` (**5×** —
  unter den sieben in §8 genannten die zweithäufigste; über das **ganze** Register
  **geteilt-vierte von sechs Klassen mit 5× oder mehr**)". Bei der Niederschrift
  (`1d0d11e`) war das richtig (drei Klassen bei 6×, drei bei 5× → die Klasse auf
  5× war geteilt-vierte). Die Closure `6371a28` hat **denselben** Zähler auf
  **6×** gesetzt und schreibt das auch hin (§7: „`negativtest…` → **6×**") — im
  selben Dokument stehen damit **5×** (§2, §6, §8) und **6×** (§7) für **eine**
  Klasse. Gemessen (L-1.8) sind die Klassen mit ≥ 5× heute **7 · 6 · 6 · 6 · 5 ·
  5**; die Klasse steht bei **6×** und ist damit **geteilt-zweite**, nicht
  vierte. Der Datierungs-Satz des Plans („Die Zähler-Stände sind der Stand **bei
  dieser Planung**") steht in **§8** — nicht in §2/LP2; §8 selbst ist damit
  gedeckt. Es ist dieselbe Teil-Behebung wie in F-1: ein Wort fehlt an **einer**
  der Stellen.
- `verifizierbar`: ja — `grep -n "5×\|6×\|7×" <Plan>`;
  `ls evidence/ | wc -l` je Eintrag (L-1.8)
- `urteil`: **ein Wort** („… zum Stand dieser Planung") **oder** die Klausel
  fallenlassen. Kein Liefer-Defekt; der Zähler selbst ist überall richtig, nur
  sein Bezugs-Zeitpunkt fehlt an einer Stelle.

### N-3 — „vier Träger mit einer Fehlzählung des Globs" ist contra-gemessen: es sind zwei Träger, und zwei der vier Grep-Treffer tragen eine andere Aussage

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B (Tatsachenbehauptung ohne Beleg-Anker) ·
  dieselbe Klasse wie N-1
- `pfad`: `slice-092-coverage-cluster-d1.md:232-238` (§6 R4) und `:304-311`
  (§7, „Risiken aus §6")
- `befund`: Gemessen am Parent `2b7c6ec` (L-1.10) tragen **zwei** Dokumente die
  Fehlzählung des Globs — der Slice-Plan (**5** Stellen) und `welle-20.md`
  (**1** Stelle) —, zusammen **sechs** Stellen. Ein naiver
  `grep -l "10 Pakete|zehn Pakete"` über die `.md` liefert **vier** Dateien; die
  zwei zusätzlichen sind `verify-slice-005` und `verify-slice-007`, und deren
  „**10 Pakete `ok`**" ist eine Aussage über einen **Testlauf**, keine über den
  D1-Umfang — sie sind **nicht** falsch und waren nicht zu berichtigen. Die Zahl
  „vier Träger" ist damit entweder die Grep-Trefferzahl (dann ist sie keine
  Träger-Zahl) oder eine Verwechslung von **Stelle** und **Träger**. Ironie der
  Stelle: der Ausgang, der von einer **Fehlzählung** handelt, trägt selbst eine.
- `verifizierbar`: ja — `git grep -l -E "10 Pakete|zehn Pakete" 2b7c6ec -- '*.md'`
  → 4 Dateien; dieselben Muster je Datei mit Kontext (L-1.10)
- `urteil`: **ein Wort** — „vier **Stellen** in **zwei** Trägern" (oder „vier
  Grep-Treffer"). Kein Liefer-Defekt; die Richtung des Ausgangs (Prüfung hat
  getragen, Sensor-Dokument bleibt zahl-richtig) ist gemessen und wahr.

---

## L-6 Negativbefunde

- **geprüft, ohne Befund: V-2 des Hauptberichts ist erfüllt.** Der
  Delta-Review-Report liegt vor (0 HIGH / 2 MEDIUM / 1 LOW / 2 INFO) und erklärt
  **beide** Fix-Commits (`ad95444` **und** `a81a6ea`) zu seinem Gegenstand
  (L-1.2: er ist im Umfang). Der Präzedenzweg — ein eigener Delta-Report — ist
  gegangen; V-2 hat keinen Rest.
- **geprüft, ohne Befund: V-1 des Hauptberichts ist beantwortet und in der Sache
  geheilt.** Die §3-Bedingung ist in §7 **benannt** („Der Fix ist damit
  **innerhalb** der Absicht, aber **außerhalb** des Buchstabens der Zeile;
  benannt, nicht still") — genau die Option, die V-1 verlangt hat; die Ausweitung
  ist nicht mehr still. Der zweite Teil (die abgeleitete Zahl ohne
  Herkunfts-Marker) ist mit `a81a6ea` gesetzt und mit `1d0d11e` korrekt
  **umgekehrt** worden (L-2).
- **geprüft, ohne Befund: die dritte Runde hat im Sensor-Träger nichts Neues
  falsch gemacht** — vier Zahlen nachgemessen (L-1.3, L-1.4), die
  `git diff`-Klausel mit zwei Ranges (L-1.5), die Paarungen nachgerechnet
  (L-1.6), Etiketten gegen §3.12 gehalten, Satz geschlossen (L-2).
- **geprüft, ohne Befund: die vier Plan-Zahlen und das Zähl-Wort** (L-1.7) —
  **13** · **10** · **12** · **11**, je an der richtigen Aussage, unter der
  Überschrift „**Vier** Zahlen, vier Dinge".
- **geprüft, ohne Befund: die Register-Zähler stimmen mit `ls evidence/`**
  (L-1.8): **7×** · **6×** · **1×**; die zwei neuen Beleg-Dateien
  (`evidence/slice-092.md` in beiden Einträgen) existieren, die `state.md`-Zeilen
  nennen dieselben Zahlen, und der abgelehnte `arbeit-ueberholt`-Kandidat ist
  **begründet** und **nicht gezählt** (Modul 6: „ein Vorgang zählt einmal";
  Ablehnung ist benannt, nicht still).
- **geprüft, ohne Befund: R1, R2 (beide Hälften) und R3 (in der Rechnung) sind
  wahr** — L-4; für R2 habe ich die zweite Hälfte (`includecolumn`) mit
  **P-9/P-10** selbst gemessen. Die Grenze von R3 („alle 24 Positionen gebunden"
  ruht auf den zwei Berichten) ist in §8 des Hauptberichts benannt und hier
  unverändert.
- **geprüft, ohne Befund: Umfang, Gate und Hygiene.** **24** Pfade im Vorgang,
  **13** davon `*_test.go`, **11** Markdown, **kein** Produktcode, kein
  `THRESHOLD`, keine ADR, kein Workflow (L-1.12); `make gates` am Closure-Stand
  **Exit 0**, sechs Checks grün, `d-check` **762** Dateien / **0** Befunde
  (L-1.1); Baum danach leer. Dieser Nachtrag liegt unter demselben Gate wie die
  Träger, über die er spricht.
- **geprüft, ohne Befund — mit einer verworfenen eigenen Probe:** mein erster
  `includecolumn`-Versuch traf die falsche Zeile (Syntax-Fehler → `build
  failed`, **kein** Test-Rot); er ist **nicht** gezählt und mit der richtigen
  Zeile neu gefahren (L-1.9). Dieselbe Disziplin, die der Hauptbericht in §6 für
  P-8 anwendet: ein Rot, das ein Compiler-Fehler ist, ist kein Beleg.

---

## L-7 Verdikt des Nachtrags

**Die dritte Runde trägt, die Closure-Inhalte sind wahr, und die Lieferung ist
unberührt.** `harness/sensors/coverage-gate.md` ist an allen vier geprüften
Punkten richtig (Herkunft · Satzschluss · Aussage gegen die Messung · nichts
Neues falsch); die vier Plan-Zahlen und ihr Zähl-Wort stimmen; **11 von 11**
Häkchen, **vier** §6-Ausgänge (drei wörtlich wahr, einer in der Zahl zu hoch),
die Register-Zähler deckungsgleich mit `ls evidence/`; `make gates` **Exit 0** am
Closure-Stand; **kein** Produktcode im ganzen Vorgang. **V-2 ist erfüllt, V-1 ist
beantwortet** — die zwei Findings des Hauptberichts sind geschlossen.

**`done/`-fähig: ja — nach drei Textnachträgen in §7/§2/§6** (**N-1** MEDIUM:
der Lauf-Marker und die erreichte Quote in §7; **N-2** und **N-3** LOW: je ein
Wort an einer Zahl-Stelle). Keiner davon ist ein Liefer-Defekt, keiner rührt an
einer Zahl, die trägt, und keiner verlangt eine neue Messung — **aber alle drei
sitzen in `slice-092-coverage-cluster-d1.md`, und diese Datei wandert beim
`git mv` nach `done/` mit**: sie gehören **vor** den Move (Modul 5: erst der
Inhalt, dann der reine `git mv`). Danach ist der Slice zu schließen.

**Offen über diesen Slice hinaus** (unverändert aus §8 des Hauptberichts): der
reale Post-Push-Lauf (`AGENTS.md` §3.10: kein Workflow geändert) und die
Stabilität der Zahl bei `THRESHOLD=80` — die `welle-20`-Closure misst dort,
nicht hier.

---

**Beleg-Lage:** jede Zahl stammt aus L-1.1…L-1.12, je in eigener
Werkzeug-Beauftragung; der Gate-Lauf (L-1.1) und seine Auswertung waren **zwei**
Schritte, sein Exit-Code wurde aus einer separaten Datei gelesen, nie durch eine
Pipe (§3.9). Arbeitsbaum-Kopien und Logs liegen **außerhalb** des Baums; der Baum
ist sauber, **kein** Commit, keine Änderung an einem Träger außer diesem Bericht.
