# ADR-0160: Leer-Test der Teil-Range am Pin-Commit und Home-relative Pfade in der `hostpaths`-Regel

**Status:** Accepted — Supersedes [`ADR-0157`](0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(nur: in Entscheidung 4 der erste Spiegelstrich „Teil-Ranges … jede endet mit
Exit 0“; unter §Konsequenzen der Satz „Ab dem nächsten Bump ist
`make doc-immutable` im Verifier-Lauf grün“; in der zweiten
Fitness-Function-Zeile die Leerfall-Hälfte „Range `B..P~1` bei leerer Range
Exit 0“. Der Pin-Commit, seine Disziplin und der normalisierte `cmp` aus
Entscheidung 4 bleiben in Kraft, ebenso alles, was
[`ADR-0158`](0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) und
[`ADR-0159`](0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
an `ADR-0157` ändern oder stehen lassen) · Supersedes
[`ADR-0075`](0075-hostpaths-reichweite-und-wortlaut.md) (nur Entscheidung 2,
der Wortlaut von `AGENTS.md` §3.11 in Fassung 2; Entscheidung 1 und 3 bleiben
in Kraft, ebenso [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
Entscheidung 1, 2 und 4)

**Datum:** 2026-10-06

**Autor:** pt9912 (Architect-Rolle, Modul 8 §Konflikt-Pfad; anderer Kontext als
der Implementer- und der Reviewer-Lauf von `slice-dcheck-v0-82-0`, deren
Findings F-1 bis F-3 diese ADR beantwortet; „eine ADR für beide Fragen“ hat
der Auftraggeber entschieden)

**Bezug:** [`ADR-0157`](0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
Entscheidung 4 · [`ADR-0075`](0075-hostpaths-reichweite-und-wortlaut.md)
Entscheidung 1 und 2 · [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
Entscheidung 2 (kein Ventil) · [`ADR-0051`](0051-cicd-pipeline-github-actions.md)
Entscheidung 7 (Pin P7) · [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md)
(gemessen/hergeleitet) · `AGENTS.md` §3.5, §3.11, §3.12 · Review
`review-slice-dcheck-v0-82-0` F-1 bis F-6 · d-check v0.82.0
(`DC-FA-VCS-002`, Ventil `hostpaths.exempt-targets`).

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie `ADR-0075` und
`ADR-0157`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Der Pin von d-check steigt mit `slice-dcheck-v0-82-0` von v0.79.0 auf v0.82.0.
Zwei Änderungen aus d-check v0.80.0 treffen Entscheidungen dieses Repos.

**Erste Frage — die leere Teil-Range.** `ADR-0157` Entscheidung 4 lässt
`make doc-immutable` bei einem Baseline-Bump in den Teil-Ranges `B..P~1` und
`P..H` um den Pin-Commit `P` laufen, „jede endet mit Exit 0“. Ist `P` der erste
Commit nach `B`, ist `B..P~1` leer. Die Fitness-Function-Zeile von `ADR-0157`
hält dafür „Exit 0“ fest; gemessen wurde das mit v0.79.0. Seit v0.80.0
(`DC-FA-VCS-002`) endet eine leere Range mit Exit 2. Gemessen in dieser ADR
(Fitness Function): `make doc-immutable RANGE=B..P~1` mit v0.82.0 druckt
„Range-Leerfall … Basis und Spitze benennen denselben Commit“, make-Exit 2.
Die Entscheidung verlangt damit einen Lauf, der am gepinnten Stand nie grün
werden kann.

Die Träger des Slice haben das mit „eine leere Teil-Range läuft nicht, wenn
`git rev-list --count <range>` `0` druckt“ ersetzt. Ohne Architect-Artefakt
ändert das eine `Accepted`-Entscheidung (Review F-1). Außerdem trägt der
Beleg den Satz nicht (Review F-2): `git rev-list --count` druckt `0` auch für
eine umgekehrte Range, zwischen deren Enden ungeprüfte Commits liegen.
Gemessen hier: `git rev-list --count H..B` druckt `0`. d-check selbst trennt
beide Fälle („Basis und Spitze benennen denselben Commit“ gegen „0 Commits“).

**Zweite Frage — Home-relative Pfade.** Seit v0.80.0 meldet `hostpaths` auch
Home-relative Pfade, also Tilde, Schrägstrich und ein Segment. Die Regel
`AGENTS.md` §3.11 im Wortlaut von `ADR-0075` Entscheidung 2 (Fassung 2) nennt
nur „host-lokale absolute Pfade“ (Wurzel-Segment aus `hostpaths.prefixes`,
Windows-/UNC-Muster). Das Gate blockiert damit eine Form, die keine Regel
verbietet (Review F-3). Der Auftraggeber nimmt die erweiterte Reichweite an.
Offen ist, wo sie entschieden wird.

- `ADR-0072` Entscheidung 1 und 2 regelt die Aktivierung und das Fehlen eines
  Ventils. Beides bleibt richtig.
- `ADR-0075` Entscheidung 1 regelt die Reichweite über die Markdown-Fläche
  (Fences eingeschlossen, verbotene Form nur als Platzhalter). Sie sagt nichts
  über die Pfad-Klasse und bleibt richtig.
- `ADR-0075` Entscheidung 2 legt den Wortlaut der Aussage fest, „host-lokaler
  absoluter Pfad (Wurzel-Segment … oder Windows-/UNC-Muster)“. Genau dieser
  Satz ändert sich.

Am bestehenden Repo zeigt das neue Modul keinen Befund (übernommen aus dem
Plan von `slice-dcheck-v0-82-0`, Liefer-Punkt 2 (a): `make docs-check` mit
v0.82.0 druckt `0 Befund(e)`; der Gate-Lauf am Ende dieses Zugs misst es
erneut).

## Entscheidung

1. **Leer-Test der Teil-Range (ersetzt in `ADR-0157` Entscheidung 4 den ersten
   Spiegelstrich).** Für jede Teil-Range `<basis>..<spitze>` um einen
   Pin-Commit `P` — `B..P~1` und `P..H` — gilt genau eines:
   - (a) **Leer** genau dann, wenn `git rev-parse --verify` für `<basis>` und
     `<spitze>` denselben Commit druckt. Dann läuft `make doc-immutable` für
     diese Teil-Range nicht. Der Beleg ist die gedruckte Zeile des Tests.
   - (b) **Nicht leer**, wenn `git merge-base --is-ancestor <basis> <spitze>`
     gilt und `git rev-list --count <basis>..<spitze>` größer 0 ist. Dann läuft
     `make doc-immutable RANGE=<basis>..<spitze>` und endet mit Exit 0.
   - (c) **Falsch gebildet** in jeder anderen Lage: Eine Seite löst nicht auf,
     oder `<basis>` ist kein Vorfahr von `<spitze>`. Die umgekehrte Range
     gehört dazu. Der Test endet mit Exit 2, und das ist rot. Der Messende
     bildet die Range neu. Er lässt sie nicht aus.

   `git rev-list --count` allein entscheidet „leer“ nicht. Der Befehl, wörtlich
   kopierbar:

   ```bash
   teilrange() {  # teilrange <basis> <spitze> — Leer-Test, danach ggf. der Lauf
     local b s n
     b=$(git rev-parse --verify -q "$1^{commit}") || { echo "teilrange: Basis $1 löst nicht auf, Exit 2"; return 2; }
     s=$(git rev-parse --verify -q "$2^{commit}") || { echo "teilrange: Spitze $2 löst nicht auf, Exit 2"; return 2; }
     if [ "$b" = "$s" ]; then echo "teilrange: $1..$2 leer (Basis = Spitze = $b), kein Lauf"; return 0; fi
     git merge-base --is-ancestor "$b" "$s" || { echo "teilrange: $1 ist kein Vorfahr von $2, Exit 2"; return 2; }
     n=$(git rev-list --count "$b..$s")
     [ "$n" -gt 0 ] || { echo "teilrange: $1..$2 enthält 0 Commits, Exit 2"; return 2; }
     echo "teilrange: $1..$2 enthält $n Commit(s), Lauf"
     make doc-immutable RANGE="$b..$s"
   }
   teilrange "$B" "$P~1"; echo "Exit $?"
   teilrange "$P" "$H"; echo "Exit $?"
   ```

   Ein falsch bestimmtes `P` fällt so nicht still durch. Liegt das gewählte `P`
   hinter dem echten Pin-Commit, enthält `B..P~1` den echten Pin-Commit und
   wird rot (gemessen, Fitness Function). Liegt es vor `B`, ist `B` kein
   Vorfahr von `P~1`, und der Test endet mit Exit 2.

2. **Die Teile von `ADR-0157`, die bleiben.** Es bleiben der Pin-Commit, der
   nur MR-Dateien ändert und `ADR-0073` nennt, und der normalisierte `cmp` je
   MR-Datei am Pin-Commit aus Entscheidung 4. Die Konsequenz „ab dem nächsten
   Bump grün“ liest sich ab hier so: Jede Teil-Range ist entweder leer nach
   (a) oder grün nach (b). Ein Exit 2 aus (c) oder aus `make doc-immutable` ist
   ein Befund und kein Leerfall. Die Leerfall-Hälfte der Fitness-Function-Zeile
   von `ADR-0157` ist der Messwert mit v0.79.0. Am Pin v0.82.0 gilt sie nicht
   mehr; die gültige Zeile steht in dieser ADR. Die Konfiguration in
   `.d-check.yml` bleibt unverändert.

3. **Home-relative Pfade gehören zur Regel (ersetzt `ADR-0075`
   Entscheidung 2, Wortlaut von `AGENTS.md` §3.11, Fassung 3).**

   Gegenstand der Regel sind:
   - der host-lokale **absolute** Pfad wie bisher (Wurzel-Segment aus
     `hostpaths.prefixes`, Windows-Laufwerks-/UNC-Muster);
   - der **Home-relative** Pfad: eine Tilde, ein Schrägstrich und ein erstes
     Segment ohne führenden Punkt. Das umfasst die Formen `~/<Verzeichnis>/…`,
     eine Datei direkt unter der Tilde (`~/<Datei>`) und ein einzelnes Segment
     ohne abschließenden Schrägstrich (`~/<Verzeichnis>`);
   - die **Tilde mit Benutzername** (`~<Benutzer>/…`). Sie nennt ein Konto
     auf einem Rechner und ist damit Maschinen-Layout.

   Kein Gegenstand der Regel, mit Grund:
   - die **Werkzeug-Konvention mit Punkt-Segment** (`~/.config/…`). Sie
     beschreibt, wo ein Werkzeug auf jedem Rechner liest, und kein Layout
     eines bestimmten Rechners;
   - die **nackte Tilde** (`~`). Sie ist kein Pfad;
   - die **Tilde in einem URL-Pfad**. Sie ist Teil einer Adresse im Netz und
     kein Pfad auf einem Rechner.

   `ADR-0075` Entscheidung 1 gilt für die neue Klasse unverändert. Die Regel
   deckt die ganze Markdown-Fläche, Fences eingeschlossen. Eine verbotene Form
   steht nur als Platzhalter da: `<Host-Wurzel>` für das Wurzel-Segment,
   `~/<Verzeichnis>/…` für den Home-relativen Pfad.

   Die Entscheidung ist eine Teil-Ablösung von `ADR-0075` und keine eigene
   Entscheidung daneben. Sie ändert den Satz, den `ADR-0075` Entscheidung 2 als
   Aussage beschließt. Eine zweite Entscheidung neben einer unveränderten
   Fassung 2 ließe zwei Wortlaute derselben Regel stehen.

   **Wortlaut der Aussage (Fassung 3):**

   > Kein Markdown-Dokument dieses Repos nennt an irgendeiner Stelle — Prosa,
   > Inline-Code oder Fenced-Code-Block — einen host-lokalen Pfad. Das ist ein
   > absoluter Pfad mit dem Wurzel-Segment eines Entwicklerrechners
   > (Präfixliste `hostpaths.prefixes`) oder mit einem Windows-Laufwerks- oder
   > UNC-Muster, oder ein Home-relativer Pfad: Tilde, Schrägstrich und ein
   > Segment ohne führenden Punkt, oder Tilde mit Benutzername. Ein
   > Schwester-Artefakt wird in der Hausform zitiert. Kein Gegenstand der Regel
   > sind die Werkzeug-Konvention mit Punkt-Segment, die nackte Tilde und die
   > Tilde in einem URL-Pfad.

   Das Formzitat bleibt wie in Fassung 2, ergänzt um den Platzhalter
   `~/<Verzeichnis>/…`.

   **Was der Sensor deckt — und was nicht (Fassung 3).** Das Modul `hostpaths`
   deckt `.md`-Dateien unter `scan.roots` in Prosa und Inline-Code. Dort meldet
   es absolute Pfade und Home-relative Pfade ohne Punkt-Segment. Es deckt
   **nicht**:
   - Fenced-Code-Blöcke;
   - die Tilde mit Benutzername (still gemessen; die Regel deckt sie, das Modul
     nicht);
   - relative Pfade;
   - Nicht-Markdown (`Makefile`, `tools/**`, `harness/mk/**`);
   - Dateien unter `scan.ignore`.

   Bei Fences und bei der Tilde mit Benutzername ist die Regel strenger als der
   Sensor. Der Wächter dort ist das Review.

4. **Das Ventil bleibt ungenutzt.** d-check bietet seit v0.80.0 das Ventil
   `hostpaths.exempt-targets` (Globs über den gemeldeten Pfad). Es ist
   verfügbar, und dieses Repo setzt es nicht. `ADR-0072` Entscheidung 2
   („keine Ausnahme, kein Ventil“) gilt weiter, und `.d-check.yml` bleibt ohne
   Block `hostpaths:`. Wer das Ventil setzen will, braucht eine Folge-ADR zu
   `ADR-0072` Entscheidung 2. Einen Fund behebt man, man nimmt ihn nicht aus.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

Erste Frage (leere Teil-Range):

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; `ADR-0157` verlangt weiter Exit 0 | kein Eingriff | ein Lauf, der am gepinnten Stand nie grün wird; die Träger hätten keine Norm |
| B — Trägerwortlaut als Ausführungsregel, Leer-Test `git rev-list --count = 0` | eine Zeile, kein neuer Befehl | ändert eine `Accepted`-Entscheidung ohne Träger (Review F-1); `0` auch für eine umgekehrte Range, gemessen an `H..B` (Review F-2) — nimmt die laute Meldung von `DC-FA-VCS-002` zurück |
| C — Teil-Range `B..P~1` durch `B..P` ersetzen und den Pin-Commit mitprüfen | nie leer | `vcs` meldet am Pin-Commit genau die vier `core-drift-vcs`, die Entscheidung 4 umgehen soll (gemessen an `B..H`) |
| D — Leerfall im Verfahren als `--staged` oder eigenes Make-Ziel kapseln | ein Aufruf | ein neues Ziel mit Vertrag und Test für einen Fall, der nur bei Bumps auftritt; der Shell-Befehl trägt dasselbe |
| **E — Leer-Test über Commit-Gleichheit, sonst Vorfahr und Zählung, sonst Exit 2 (gewählt)** | die drei Lagen sind getrennt; die umgekehrte Range bleibt rot; ein Befehl zum Kopieren, kein neues Ziel | ein Verfahrensschritt mehr für den Verifier; der Test ist Shell im Träger, kein Sensor |

Zweite Frage (Home-relative Pfade):

| Option | Pro | Contra |
|---|---|---|
| F — nichts tun; Regel nennt nur absolute Pfade | kein Eingriff | ein rotes Gate ohne Norm (Review F-3) |
| G — Home-relative Pfade per `hostpaths.exempt-targets` ausnehmen und die Regel lassen | Regel und Gate wieder gleich weit | widerspricht `ADR-0072` Entscheidung 2; der Auftraggeber nimmt die Reichweite an |
| H — eigene Entscheidung neben `ADR-0075` Fassung 2 | `ADR-0075` bleibt ganz in Kraft | zwei Wortlaute derselben Aussage in zwei `Accepted`-ADRs; der Leser muss sie zusammenlegen |
| **I — Teil-Supersede von `ADR-0075` Entscheidung 2, Fassung 3 (gewählt)** | ein Wortlaut; Entscheidung 1 (Fences, Platzhalter) gilt unverändert für die neue Klasse | die Regel ist an der Tilde mit Benutzername strenger als ihr Sensor; Wächter ist das Review |

**Fazit:** E und I.

## Konsequenzen

- Positiv: Der Verifier eines Bump-Slice hat ein Verfahren, das am gepinnten
  Stand grün oder rot ist und keinen Leerfall still durchlässt. Die laute
  Meldung von `DC-FA-VCS-002` bleibt für die falsch gebildete Range wirksam.
- Positiv: Gate und Regel decken die Home-relativen Pfade gleich weit, außer an
  den benannten Rändern.
- Negativ mit Grenze: Den Leer-Test führt die Rolle aus, kein Sensor. Wer den
  Befehl nicht fährt und nur zählt, fällt in Option B zurück. Wächter ist der
  Reviewer, der die gedruckte Zeile des Tests im Slice-Plan sucht.
- Negativ mit Grenze: Bei der Tilde mit Benutzername ist die Regel strenger als
  das Modul. Ändert d-check das, schließt Re-Evaluierungs-Trigger (c) die Lücke.
- Folgepflicht (Implementer, `slice-dcheck-v0-82-0`):
  1. `AGENTS.md` §3.11 auf Fassung 3: Überschrift „Kein host-lokaler
     absoluter oder Home-relativer Pfad in der Doku“, die Aussage aus
     Entscheidung 3, der Satz zum Platzhalter um `~/<Verzeichnis>/…` ergänzt,
     „Was der Sensor deckt“ nach Entscheidung 3, Träger-Zeile um diese ADR
     ergänzt. Die Überschrift hat keinen eingehenden Anker-Link (gemessen mit
     `git grep -n "311-kein"`, kein Treffer).
  2. `harness/sensors/docs-check.md` §Grenze Punkt 8: die Formen nach
     Entscheidung 3 (auch Datei direkt unter der Tilde und ein einzelnes
     Segment, Review F-6). Der Satz „meldet damit mehr, als der Wortlaut von
     `AGENTS.md` §3.11 nennt“ wird ersetzt durch „Regel und Modul sind bei den
     Home-relativen Pfaden gleich weit, außer bei der Tilde mit Benutzername;
     `~/.config/…`, nackte Tilde und URL sind kein Gegenstand der Regel“. Das
     Ventil steht als verfügbar und nicht gesetzt da (Entscheidung 4).
     §Bindung: „kein host-lokaler absoluter oder Home-relativer Pfad in der
     Doku“ mit dieser ADR neben `ADR-0075`.
  3. `.claude/agents/verifier.md`, `.claude/agents/implementer.md` und
     `harness/targets/pin-stale.md` §Bump-Ablauf: Die Klammer zur leeren
     Teil-Range wird durch den Leer-Test aus Entscheidung 1 ersetzt, mit
     Verweis auf den Befehl in dieser ADR. Der Satz „`git rev-list --count`
     druckt `0`“ als Leer-Kriterium entfällt.
  4. Im Slice-Plan die Ausführungsregel durch einen Verweis auf diese ADR
     ersetzen, dazu die Nachbarn aus Review F-4 und F-5.
- Folgepflicht (Planner, Closure von `slice-dcheck-v0-82-0`): keine
  Register-Pflicht aus dieser ADR. Ob „Werkzeug-Bump ändert eine
  `Accepted`-Fitness-Function“ eine Beobachtung ist, entscheidet die Closure.

## Fitness Function (falls maschinell prüfbar)

Gemessen am 2026-10-06 mit d-check v0.82.0
(`sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`,
dem Pin aus `d-check.mk`) an einem Wegwerf-Klon von `1474274e` im Scratchpad.
`B` = `1474274e`; `P` = `689b88d6` (stellt in `MR-001` bis `MR-004` den Pfad
`.harness/baseline/v6.14.1/` auf `v6.14.9/` um, nur MR-Dateien, Message mit
`ADR-0073`); `H` = `82ed7f01` (ein Commit an `README.md`); `SEITE` = `552b8e39`
(ein Commit auf einem Seitenzweig ab `B~1`).

| Tooling | Regel | Make-Target |
|---|---|---|
| `teilrange` (Entscheidung 1), erprobt am Text dieser ADR (Funktion per `awk` aus dem `bash`-Block gezogen) | leer: `teilrange B P~1` druckt „… leer (Basis = Spitze = 1474274e…), kein Lauf“, Exit 0 — grün. Zum Vergleich `make doc-immutable RANGE=B..P~1` direkt: „Range-Leerfall … Basis und Spitze benennen denselben Commit“, make-Exit 2. Nicht leer: `teilrange P H` druckt „… enthält 1 Commit(s), Lauf“ und `d-check: 1795 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — grün; `teilrange B~3 P~1` „… enthält 3 Commit(s), Lauf“, `0 Befund(e)`, Exit 0 — grün. Umgekehrt: `teilrange H B` druckt „… ist kein Vorfahr von …, Exit 2“, Exit 2 — rot (`git rev-list --count H..B` druckt dort `0`, `make doc-immutable RANGE=H..B` „Range-Leerfall … 0 Commits“, Exit 2). Seitenzweig: `teilrange SEITE P~1` Exit 2 — rot, obwohl `git rev-list --count SEITE..P~1` `1` druckt. Unbekannte Basis: „Basis gibtsnicht löst nicht auf, Exit 2“ — rot. Falsches `P` (`H` statt `P`): `teilrange B H~1` läuft und meldet `core-drift-vcs` an den MR-Dateien, make-Exit 2 — rot. Simulierter Bump, volle Range `make doc-immutable RANGE=B..H`: `1795 Datei(en) geprüft, 4 Befund(e)`, je MR-Datei ein `core-drift-vcs`, make-Exit 2; normalisierter `cmp` aus `ADR-0157` am Pin-Commit: je MR-Datei `0`. Alle Fälle an einer Instanz (dieser Klon, eine Pin-Commit-Lage); dass jede andere Bump-Lage gleich fällt, ist *hergeleitet* aus den drei Zweigen des Tests | kein Make-Target (Verifier-Lauf) |
| `d-check` Modul `hostpaths` (Entscheidung 3), Wegwerf-Verzeichnis mit `modules: [hostpaths]` und einer `probe.md` | v0.82.0: Exit 1, `d-check: 1 Datei(en) geprüft, 4 Befund(e)`. Gemeldet sind ein Home-relativer Pfad mit zwei Segmenten in Prosa und derselbe in Inline-Code, eine Datei direkt unter der Tilde und ein einzelnes Segment. Still bleiben `~/.config/…`, die nackte Tilde, die Tilde mit Benutzername, die Tilde im URL-Pfad, der Fence und die Platzhalter `~/<Verzeichnis>/<Datei>` (Prosa) und `` `~/<Verzeichnis>/…` `` (Inline). Dieselbe Datei mit v0.79.0 (`sha256:b4b8756b…`): Exit 0, `0 Befund(e)`. Am Repo: `make docs-check` im `make gates`-Lauf dieses Zugs `0 Befund(e)` (Lauf-Beleg im Bericht des Zugs) | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger:
- **(a)** `d-check` bietet im Modul `vcs` eine Normalisierung im Core
  (`ADR-0157` Trigger (b)). Dann entfallen Teil-Ranges und Leer-Test.
- **(b)** `d-check` macht den Leerfall wieder zu Exit 0 oder bietet einen
  Schalter dafür. Dann ersetzt der Schalter den Leer-Test (a), (b) und (c)
  bleiben.
- **(c)** `hostpaths` meldet die Tilde mit Benutzername oder Fences. Dann
  sind Regel und Sensor dort gleich weit, und die benannte Lücke entfällt.
- **(d)** Ein Dokument braucht eine Werkzeug-Konvention ohne Punkt-Segment
  unter der Tilde als Beispiel, und der Platzhalter trägt sie nicht. Dann
  ist eine eng benannte Folge-ADR zu `ADR-0072` Entscheidung 2 fällig, kein
  stilles Ventil.

Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Accepted — Architect-Zug zu Review `review-slice-dcheck-v0-82-0` F-1 bis F-3 (Konflikt-Pfad); eine ADR für beide Fragen nach Entscheidung des Auftraggebers | Review `review-slice-dcheck-v0-82-0` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0160` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
