# Review-Report: slice-079 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität, §6-Risiko-Ausgänge, Beobachtungs-Register und die drei
Paarungen sind **nicht** Gegenstand dieses Reports (Verifier bzw.
Planner-Closure, Modul 11/6).

**Gegenstand:** `slice-079`, Commit `ca9864f` (der Schnitt) und `3827188`
(DoD-Nachzug), Diff `35b9a00..3827188`; `3827188` ist `HEAD`. Sechs Dateien:
`Dockerfile` (nur die `coverage`-Stufe), `harness/mk/coverage.mk`,
`harness/sensors/coverage-gate.md`, `harness/README.md`, `AGENTS.md` §4 und der
§2-Häkchen-Nachzug im Slice-Plan. Kein Produktionscode, keine Tests, keine
Spec-Datei.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14) · **Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:**
2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-079-coverage-scope-schnitt.md`
  vollständig (§1–§8), einschließlich des §2-Nachtrags `3827188`
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  (Punkte 1–3, Fitness Function, Re-Evaluierungs-Trigger),
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
  §(a) (Rampe, Scope-Bullet vor dem Supersede),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md) (Unit- vs. E2E-Tier)
- `AGENTS.md` §3.1 (Docker-only), §3.2 (Suppression-Verbot), §3.6
  (Schwellen-Senkung nur per ADR), §3.7 (Ist-Zustand), §3.9 (Exit-Code),
  §4 (Gate-Tabelle), §5 (Doku-Regeln)
- `harness/conventions.md` (MR-000 ID-Schema), `.d-check.yml`
  (`reviews`-/`ids`-/`hostpaths`-Regeln)
- Vorherige Läufe am gleichen Gegenstand: das Review zu `slice-076`
  (derselbe Kalibrierungs-Mechanismus, 3 LOW) und
  der Architect-Verdikt zum Coverage-Gate-Messgegenstand (Anlass des
  Schnitts)

---

## Eigene Messungen dieses Laufs

Alle Zahlen dieses Reports sind mit **eigenen** Läufen erzeugt (Docker-only,
Exit-Code ungepiped gelesen, `AGENTS.md` §3.9). Der Vorher-Stand wurde über
eine Arbeitsbaum-Kopie bei `35b9a00` mit dem **alten** `Dockerfile` gebaut
(`--build-arg COVERAGE_THRESHOLD=0`), die Mess-Artefakte lagen außerhalb des
Arbeitsbaums; `git status --porcelain` ist nach allen Läufen leer.

| Lauf | Exit | Ergebnis |
|---|---|---|
| `make coverage-gate` (Stufe 65, `HEAD`) | **0** | `coverage-gate: OK — Coverage 69.90% erfüllt Schwelle 65%` |
| `make coverage-gate THRESHOLD=75` | **2** (Skript Exit 1) | `coverage-gate: FAIL — Coverage 69.70% unter Schwelle 75%` — Abstand 5,30 pp |
| `make gates` | **0** | baseline-verify `v6.5.0` OK (54 Dateien); d-check 646 Dateien/0 Befunde; commit-traceability 5 Commits; a-check 0 Befunde; coverage-gate OK 69,70 % |
| `coverage`-Stufe bei `35b9a00` (alter Umfang) | **0** | `Coverage 49.30% erfüllt Schwelle 0%` |

**Nachbau der Bezifferung** (Deduplizierung über die Block-Position, wie in
`ADR-0071` beschrieben):

| Größe | vorher (`35b9a00`) | nachher (`3827188`) |
|---|---|---|
| Blöcke / Statements im Profil | 1668 / **2467** | 1104 / **1679** |
| davon gedeckt | 1215 (49,25 %) | 1173 (69,86 %; Gate druckt 69,9 %) |
| Statements in den drei ausgenommenen Paketen | **788**, davon **44** gedeckt | 0 |
| davon `postgresack` | 23 / 2 | — |
| davon `replication/receive` | 155 / 11 | — |
| davon `postgresstorage` (ohne `mapper`) | 610 / 31 (5,1 %) | — |
| `postgresstorage/mapper` | 15 / 12 (80,0 %) | 15 / 12 — **bleibt im Gegenstand** |
| `postgresstorage/queries` | im Gegenstand, `[no test files]` | ebenso |

**Vergleich der beiden Profile je Datei** (der Kern des Liefer-Punkts 1):
keine Datei existiert nur im Nachher-Profil; **keine** Statement-Zahl-Differenz
außerhalb der drei Pakete; die einzige Differenz überhaupt ist die
gedeckte Zahl in `internal/bootstrap/wiring.go` (193 → 195 Statements,
Lauf-zu-Lauf-Schwankung). Die drei `[no test files]`-Pakete druckt nur der
Lauf mit `-coverpkg`; `go test ./cmd/...` druckt `[no test files]` für
`cmd/pg-change-feed` (49 Statements, `main.go` mit `count = 0` im Profil).

**Filter-Probe in der `coverage`-Stufe:** `go list ./internal/... ./cmd/...`
liefert 32 Pakete, der Filter lässt 29 durch; die Differenz ist genau
`postgresstorage`, `postgresack`, `replication/receive`;
`postgresstorage/mapper` und `postgresstorage/queries` sind beide enthalten.

**Eigenschafts-Probe an den Paketen mit überspringenden Tests:**
`postgresstorage` hat 10 Testdateien, davon 8 mit `t.Skip` und 2 ohne
(`administration_backoff_internal_test.go`, `identifier_test.go`);
`internal/bootstrap` (591 Statements, 262 gedeckt = 44,3 %) und
`internal/adapters/driven/natsnotify` (33 / 28) führen ebenfalls
überspringende Testdateien, laufen aber netzlos.

---

## Findings

### F-1 — Die §Grenze-Aussage „kein Paket im Gegenstand ist ganz ohne Testdatei“ ist falsch; `cmd/pg-change-feed` fehlt in der Aufzählung

- `kategorie`: LOW
- `quelle`: Maintainability · Slice-Plan §3 (die Sensor-Doku trägt den
  Messgegenstand) · `AGENTS.md` §3.7 (eine Aussage beschreibt, was da ist) ·
  Vorläufer derselben Wortlaut-Stelle: `slice-049`, dort mitgeprüft und
  mitgetragen
- `pfad`: `harness/sensors/coverage-gate.md:60-64`
- `befund`: Die Zeile behauptet, kein Paket im Gegenstand sei ganz ohne
  Testdatei, und nennt die drei `[no test files]`-Pakete. `cmd/pg-change-feed`
  liegt im Gegenstand (in `-coverpkg` **und** in der Testpaket-Liste), hat
  keine Testdatei und trägt 49 Statements, die im Profil mit `count = 0`
  stehen. Die Aussage ist damit unrichtig und die Aufzählung unvollständig —
  die drei genannten Pakete tragen keine ausführbaren Statements, das
  fehlende vierte trägt 49. Dass nur drei im Log stehen, ist ein Artefakt der
  Aufrufform: mit `-coverpkg` druckt `go test` für dieses Paket
  `coverage: 0.0% of statements` statt `[no test files]`.
- `verifizierbar`: nein — kein Gate prüft Doku-Aussagen; mechanisch
  nachprüfbar mit `go test ./cmd/...` in der gepinnten Toolchain-Stage
- `klasse`: „Messgegenstands-Aussage übertrifft die reale Paketlage“

### F-2 — Die zwei benannten Grenzen des Schnitts stehen in keinem Artefakt; die §Grenze nennt sie nicht

- `kategorie`: LOW
- `quelle`: Maintainability · [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Fitness Function (erste Hälfte: „keine Block-Position im Profil liegt in
  …“) · `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das
  Beobachtungs-Register („Wer schreibt, wer liest“) · `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Closure- und Lerneintrag-Regeln
  („der Report selbst ist Lauf-Beleg und wird über Läufe hinweg nicht
  gelesen“)
- `pfad`: `harness/sensors/coverage-gate.md:49-70` (§Grenze, drei Punkte) ·
  `Dockerfile:65-66` (Filter) · Beleg der Arithmetik: das Profil dieses Laufs
- `befund`: (a) Wird `postgresack` in den Gegenstand zurückgenommen, bleibt
  die Stufe grün — die Prozent-Hälfte fängt die Rücknahme nicht ((1173 + 2) /
  (1679 + 23) = 69,04 % ≥ 65; aus dem Implementer-Lauf 68,68 %);
  (b) die Testpaket-Liste ist nur Disziplin, kein Sensor. Beides benennt der
  Implementer als Lücke, im Diff hat sie keinen Träger: die §Grenze listet
  drei andere Punkte, und `ADR-0071`s Re-Evaluierungs-Trigger (a) greift nur
  bei Kommen/Gehen eines Pakets, nicht bei einer Rücknahme per `-coverpkg`.
  Dieser Report ist Lauf-Beleg und wird über Läufe hinweg nicht gelesen —
  §7/Beobachtungs-Register wären die dauerhafte Adresse; §7 ist noch leer.
- `verifizierbar`: ja — aus dem Profil nachrechenbar
  (`postgresack`-Statements aus dem Vorher-Profil zurückaddieren)
- `klasse`: „Benannte Sensor-Grenze ohne Träger im Artefakt“

### F-3 — Das Plan-Artefakt trägt zwei Kopier-Reste (außerhalb des `35b9a00`-Diffs)

- `kategorie`: LOW
- `quelle`: Maintainability · `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  §Ziel-Form: Slice (ein Plan ist die Grenze, an der sich ein wachsender Slice
  messen lässt — kein Vorlagen-Text) · Vorlage
  `.harness/baseline/v6.5.0/templates/docs/plan/planning/slice.template.md:195`
- `pfad`: `docs/plan/planning/in-progress/slice-079-coverage-scope-schnitt.md:63-90`
  und `:288-302`
- `befund`: Der Block „**Keine Mindestzahl.** … Was hier steht, ist die Grenze
  …“ steht in §1 **zweimal wörtlich**; §8 führt zusätzlich den
  Vorlagen-Satz „**Modus-Begründungsblock — Umfang.** Pflicht, sobald
  mindestens eine berührte Sub-Area BF oder Hybrid ist …“ verbatim aus der
  Vorlage, neben der repo-eigenen Antwort. Beide Reste stammen aus den
  Planungs-Commits `bd9abd5`/`78f84dd` und liegen damit **vor** dem
  Diff-Fenster dieses Laufs — sie sind Plan-Defekt, nicht Implementer-Arbeit.
- `verifizierbar`: ja — Diff des Plans gegen die Vorlage
- `klasse`: „Template-/Kopier-Rest im Plan-Artefakt“

### F-4 — Das Filter-Muster ist suffix-verankert und würde ein gleichnamiges Paket anderswo still mitnehmen

- `kategorie`: INFO
- `quelle`: Maintainability · [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Entscheidung Punkt 1 (Eigenschaft vor Liste)
- `pfad`: `Dockerfile:65-66`
- `befund`: `grep -vE '(^|/)(postgresstorage|postgresack|replication/receive)$'`
  ist am Zeilenende verankert, nicht am Baum-Anfang: jedes künftige Paket mit
  diesem Pfad-Suffix, irgendwo im Baum, fällt mit heraus, ohne dass eine
  Zeile sich ändert. Heute ist kein solcher Fall vorhanden — die Differenz
  zwischen `go list` (32) und dem Filter-Ergebnis (29) ist genau die drei.
- `verifizierbar`: ja — `go list` gegen die gefilterte Liste im Container
- `klasse`: „Denylist-Muster suffix-verankert“

### F-5 — Die Eigenschaft ist gröber als ihre Ausprägung: nach Paket-Granularität bleiben zwei Pakete mit überspringenden Tests im Gegenstand

- `kategorie`: INFO
- `quelle`: Maintainability · [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Entscheidung Punkt 1 („Ein Paket, dessen Testlauf einen externen Dienst
  voraussetzt, ist nicht Gegenstand dieses Gates“) · Slice-Plan §1
  („Die drei Pakete sind die Ausprägung dieser Eigenschaft“) und §4
  (Rückführung `in-progress → open` bei „nur teilweise erfüllt“)
- `pfad`: `Dockerfile:65-66` · `harness/sensors/coverage-gate.md:15-21` ·
  Belege: die Vorher-Profil-Zahlen dieses Reports
- `befund`: `postgresstorage` erfüllt die Eigenschaft nur zu 5,1 % (2 von 10
  Testdateien ohne Skip, 31 netzlos gedeckte Statements) und ist ausgenommen;
  `internal/bootstrap` (44,3 % netzlos gedeckt) und
  `internal/adapters/driven/natsnotify` (84,8 %) führen ebenfalls
  überspringende Testdateien und bleiben im Nenner. Ein **viertes** Paket, das
  die Eigenschaft **vollständig** erfüllt, gibt es nicht — die Kandidatenprobe
  ist vollständig. Wird die Eigenschaft feiner gelesen (Testdatei- statt
  Paket-Granularität), kippt die Grenze an drei Stellen gleichzeitig; die
  `postgresstorage`-Zahl 5,1 % zeigt, um wie wenig es dabei geht.
- `verifizierbar`: ja — Testdatei-Skip-Liste und Profil-Zahlen pro Paket
- `klasse`: „Eigenschaft feiner als ihre Ausprägung“

### F-6 — Der Gate-Lauf führt die Testdateien der drei ausgenommenen Pakete nicht mehr aus

- `kategorie`: INFO
- `quelle`: Maintainability · `AGENTS.md` §4 (Gate-Tabelle: `make gates` führt
  fünf Gates) · `.github/workflows/ci.yml` (`make gates` + `make test`)
- `pfad`: `Dockerfile:68-72`
- `befund`: Vorher liefen die drei Pakete (mit realen Skips) im
  `coverage`-Lauf mit; jetzt stehen sie weder in `-coverpkg` noch in der
  Testpaket-Liste. Ein Übersetzungsfehler in ihren Testdateien erreicht damit
  `make gates` nicht mehr; `make test` (`go test -race ./...`, **kein** Teil
  von `make gates`) führt sie weiterhin aus, und die CI fährt beide Ziele.
  Das ist die im Plan verlangte Änderung („die Testpaket-Liste … führt die
  Pakete nicht mehr“), in §6 aber nicht als Nebenwirkung notiert.
- `verifizierbar`: ja — `make test` gegen `make gates` über dieselben Pakete
- `klasse`: „Gate-Messlauf lässt Testpakete fallen (benannt)“

### F-7 — `AGENTS.md` §3.2 ist nicht falsch, aber seine Antwortmenge ist gegenüber der getroffenen Entscheidung unvollständig

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.2 (Suppression-Verbot) ·
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Entscheidung Punkt 2 („das ist **keine** Schwellen-Senkung: der Gegenstand
  wird präzise“)
- `pfad`: `AGENTS.md:93-103`
- `befund`: Die Aussage „keine Zeile lässt sich davon ausnehmen“ gilt weiter —
  der Paket-Ausschluss ist keine Zeilen-Ausnahme *innerhalb* des Gegenstands,
  und „keine stille Ausnahmeliste“ bleibt gewahrt, weil er ADR-gedeckt,
  uniform und in beiden Trägern benannt ist. Die Passage nennt als Antwort auf
  eine unerreichbare Schwelle aber nur Carveout oder bootstrap-aware Gate; die
  hier gewählte dritte Antwort (Gegenstand präzisieren, Entscheidung per ADR)
  steht dort nicht, obwohl sie jetzt zweimal inhaltlich vorliegt.
- `verifizierbar`: nein — Regel-Text, kein Gate-Gegenstand
- `klasse`: „Regel-Antwortmenge unvollständig gegenüber getroffener Entscheidung“

### F-8 — Der DoD-Wortlaut „der Lauf nennt real 1679 Statements“ verspricht eine Ausgabe, die der Lauf nicht druckt

- `kategorie`: INFO
- `quelle`: Maintainability · Slice-Plan §2 Liefer-Punkt 1 (Verifier-lesbarer
  DoD-Wortlaut) · Rollen-Zeiger: die Bewertung der DoD-Erfüllung liegt beim
  **Verifier** (Modul 11)
- `pfad`: `docs/plan/planning/in-progress/slice-079-coverage-scope-schnitt.md:104-109`
- `befund`: Der Lauf druckt `total: (statements) 69.9%` — die Zahl 1679 steht
  in keiner Ausgabe der `coverage`-Stufe; sie stammt aus der Deduplizierung
  des Profils über die Block-Position. Der Wert selbst ist exakt reproduziert
  (eigener Nachbau: 1679, vorher 2467); nur der Wortlaut „der Lauf nennt“
  beschreibt eine Ausgabe, die nicht existiert.
- `verifizierbar`: ja — im Log der `coverage`-Stufe
- `klasse`: „DoD-Wortlaut verspricht eine Ausgabe, die der Lauf nicht druckt“

## Negativbefunde

- geprüft, ohne Befund: `Dockerfile` (nur `coverage`-Stufe) — der Filter steht
  in der Stufe, nicht im Gate-Skript (Plan §3); das Rezept enthält kein `$$`
  (BuildKit-Entwertung damit ohne Objekt); das Shell-Konstrukt ist real
  ausgeführt — `$COVERAGE_THRESHOLD` wird aus der ENV expandiert (Beleg: die
  Läufe melden „Schwelle 65 %“ bzw. „Schwelle 75 %“), `pkgs=( … )` und
  `"${pkgs[@]}"` tragen die 29 Pakete, `-coverpkg` das Komma-Join derselben
  Liste; `SHELL … pipefail` unverändert. Der neue Kommentar trägt Ist-Zustand
  und Kopplung (Eigenschaft vor Liste, `go list` zieht neue Pakete mit), keine
  Slice-/Wellen-Chronik (§3.7).
- geprüft, ohne Befund: `harness/mk/coverage.mk` — **ein** beweglicher Ort
  (`THRESHOLD ?= 65`); repo-weiter `grep` über `*.mk`, `*.md`, `*.yml`, `*.sh`
  findet keinen zweiten Träger der *geltenden* Stufe (die übrigen Stellen
  nennen Rampe 65 → 80 oder verweisen); Einstieg 65 ist gegenüber 40 aus
  `slice-076` **keine** Senkung (§3.6 unberührt); Kommentar-Klassen Zusage,
  Kopplung, Grenze — keine Stufen-Chronik.
- geprüft, ohne Befund: `harness/sensors/coverage-gate.md` §Vertrag und
  §Kalibrierungs-Bindung — der Vertrag nennt den **Messgegenstand**
  (Eigenschaft plus Ausprägung, `mapper` ausdrücklich im Gegenstand) und die
  benannte Folge (`ADR-0071` Punkt 3); die Bindung führt die Rampe und
  verweist für den beweglichen Wert auf `coverage.mk`; die Exit-Tabelle ist
  unverändert und meint die Skript-Exits; die Beleg-Zeile zitiert konkrete
  Läufe (beide Prozentwerte dieses Reports sind real vorgekommen: 69,70 % und
  69,90 %), was die Sektion selbst als „Beleg, der nicht wandert“ deklariert.
- geprüft, ohne Befund: `harness/README.md` §Sensors-Zeile und `AGENTS.md` §4
  — beide nennen Umfang plus Rampe und keinen beweglichen Wert; der Edit ist
  durch Plan §3 gedeckt (die alte Zeile nannte den Umfang ausdrücklich), und
  die Form folgt `slice-076`.
- geprüft, ohne Befund: `AGENTS.md` §3.2 — von `ADR-0071` **nicht** berührt;
  Begründung in F-7.
- geprüft, ohne Befund: Kalibrierungs-Rechnung — 69,70 % (Rot-Lauf und
  `make gates`) und 69,90 % (Grün-Lauf) runden beide auf die volle 5-%-Stufe
  65 ab; die Wahl ist damit unempfindlich gegen die real gemessene
  Lauf-zu-Lauf-Schwankung (± 3 Statements), und der Rot-Beleg liegt mit
  5,30 pp deutlich außerhalb dieses Bereichs (Plan §2 Liefer-Punkt 3).
- geprüft, ohne Befund: §1-Abgrenzung eingehalten — die Dateiliste des Diffs
  ist sechs Dateien; **kein** `tools/coverage-gate.sh`, **kein** `Makefile`,
  **kein** `internal/**`, **kein** `test/**`, **keine** Spec-Datei. Die
  Schicht-Abgrenzung („Produktionscode unter `internal/**`“) ist gehalten.
- geprüft, ohne Befund: `docs/plan/adr/**` — im Range unverändert,
  `ADR-0071` bleibt `Accepted`, keine zweite Kalibrierung mitgenommen.
- geprüft, ohne Befund: Commit-Range — beide Betreffe nennen `ADR-0071`, keine
  `SPEC-*`/`ARC-*` im Betreff, kein Attributions-Trailer.

## Was ausdrücklich trägt (gegen die Prüf-Schwerpunkte)

- **Der Schnitt-Mechanismus.** Die Paketliste ist abgeleitet, nicht
  hartkodiert: `go list` + Suffix-Filter, derselbe Ausdruck für `-coverpkg`
  **und** die Testpaket-Liste — ein künftiges Paket zieht damit in beide mit
  (F-4 nennt die eine Ausnahme: gleicher Pfad-Suffix anderswo).
  `postgresstorage/mapper` (15 / 12, 80,0 %) und `postgresstorage/queries`
  bleiben erwiesen im Gegenstand.
- **Die Bezifferung — die stärkste Zusage des Slice hält.** Der Nachbau
  reproduziert „nur Ausgeschlossenes fiel weg“ exakt: 788 ausgeschlossene
  Statements, davon 44 gedeckt; die Statement-Menge **außerhalb** der drei
  Pakete ist in beiden Läufen je Datei **identisch** (keine Datei nur
  nachher, keine Statement-Differenz). Die Rest-Differenz ist die gedeckte
  Zahl einer einzigen Bootstrap-Datei (193 → 195) — dieselbe
  Lauf-zu-Lauf-Klasse, die `ADR-0071` dokumentiert.
- **Die Eigenschafts-Prüfung und die ausgebliebene Rückführung.** Die
  Beobachtung ist zutreffend (F-5 belegt sie unabhängig), und die Entscheidung,
  `in-progress → open` **nicht** auszulösen, hält: die Rücknahme eines
  ausgenommenen Pakets wäre eine Änderung der von `ADR-0071` getroffenen
  Entscheidung und damit Architect-Sache (`AGENTS.md` §3.5, Folge-ADR), keine
  stille Paketwahl; Plan §6 Risiko 1 hat genau diesen Fall (44 gedeckte
  Statements fallen mit) vorweggenommen und die Bezifferung als Antwort
  vorgeschrieben — die liegt vor. Der Rückführungs-Trigger greift daher nicht;
  was bleibt, ist die Adresse für die Beobachtung (F-2/F-5).
- **Die Neukalibrierung.** 69,70 % → 65 über den unveränderten Mechanismus
  (`ADR-0054` §(a): abrunden auf die volle 5-%-Stufe), Endstufe 80 bleibt,
  65 > 40 ⇒ keine Senkung; ein beweglicher Ort, eigener `grep` bestätigt.
- **Die `AGENTS.md`-Frage.** §4 nennt Umfang und Rampe ohne beweglichen Wert;
  §3.2 ist nicht berührt (F-7: Antwortmenge, nicht Wahrheitswert).
- **Die Lücken (a) und (b).** Beide sind zutreffend beobachtet und
  arithmetisch bzw. strukturell bestätigt; sie sind benannte Grenzen — ihr
  offener Punkt ist allein die Adresse (F-2).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** „Messgegenstands-Aussage übertrifft die reale
Paketlage“ · „Benannte Sensor-Grenze ohne Träger im Artefakt“ ·
„Template-/Kopier-Rest im Plan-Artefakt“ · „Denylist-Muster suffix-verankert“ ·
„Eigenschaft feiner als ihre Ausprägung“ · „Gate-Messlauf lässt Testpakete
fallen (benannt)“ · „Regel-Antwortmenge unvollständig gegenüber getroffener
Entscheidung“ · „DoD-Wortlaut verspricht eine Ausgabe, die der Lauf nicht
druckt“

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Kein Finding bestreitet eine
Regel oder ein Verdikt; es gibt keinen Konflikt-Pfad (Modul 8), und keiner der
Befunde ist wegen eines Implementer-Widerspruchs herabgestuft.

**Übergabe:** F-1 und F-2 nehmen die **Rückkante Reviewer → Implementer** —
beide ändern, *was die Sensor-Doku über den Messgegenstand und seine Grenzen
aussagt*, und die Sensor-Doku ist der Träger dieses Schnitts (Plan §3); der
Umfang ist zwei Stellen in `harness/sensors/coverage-gate.md`. F-3 nimmt die
**Rückkante Review → Plan** (Plan-Defekt, Planner). F-4 bis F-8 gehen **ohne
Rückkante** in die Closure; ihre Klassen sind der Übergabepunkt in den
Steering-Loop-Zähler. F-8 ist ein Rollen-Zeiger: die DoD-Bewertung selbst
liegt beim Verifier (Modul 11).

**DoD-Checkbox-Nachzug:** **nein** — mit F-1/F-2 ist eine (kleine) Fixrunde am
Implementer verbunden, also greift
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde nicht; die
Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor“ bleibt in
§2 des Slice-Plans offen und wird regulär bei Schritt 21 des
Implementer-Workflows nachgezogen. Nimmt der Implementer F-1/F-2 stattdessen an
oder begründet sie (isolierte LOWs, Modul 8) und wird kein Fix beauftragt,
greift der Nachzug in dem Lauf, der den DoD-Nachtrag ohnehin fährt.

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-Konformität, die
§6-Risiko-Ausgänge und die drei Paarungen prüft der Verifier bzw. die
Planner-Closure separat (Modul 11/6/5).
