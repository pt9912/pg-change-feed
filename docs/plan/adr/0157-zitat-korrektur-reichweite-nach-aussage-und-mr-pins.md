# ADR-0157: Zitat-Korrektur — Reichweite nach Aussage statt nach Abschnitt; MR-Einträge und ihre Baseline-Pins

**Status:** Accepted — Supersedes [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(nur die Abschnitte-Liste in deren Entscheidung 1 und den dort zitierten
Wortlaut für `AGENTS.md` §3.5 in deren Entscheidung 2; die Klasse selbst, die
Belegform (Entscheidung 3), die Record-Regel (Entscheidung 4) und die Ablehnung
von `Supersedes` als Korrekturweg für reine Zitat-Fehler (Entscheidung 5)
bleiben unverändert in Kraft und werden hier nicht wiederholt)

**Datum:** 2026-10-06

**Autor:** pt9912 (Architect-Rolle, Modul 8 §Konflikt-Pfad; anderer Kontext als
der Implementer- und der Reviewer-Lauf von `slice-harness-baseline-v6-14-1`,
deren Findings F-2 und F-3 diese ADR beantwortet)

**Bezug:** [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(die präzisierte Klasse) · [`ADR-0045`](0045-commit-traceability-standing-gate.md)
(Commit-Kennung als Beleg) · [`ADR-0051`](0051-cicd-pipeline-github-actions.md)
Entscheidung 7 (Pin-Inventar, Achse P8 Kurs-Baseline) ·
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) (gemessen/hergeleitet)
· [`ADR-0095`](0095-review-klasse-exempt-status-check.md) und
[`ADR-0123`](0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) (die korrigierten
Fälle) · `AGENTS.md` §3.5, §3.6, §3.12 · `.d-check.yml` Block `vcs:` ·
`.claude/agents/verifier.md` (Lauf `make doc-immutable`) ·
`harness/targets/pin-stale.md` §`make pin-stale-baseline` (Bump-Ablauf) ·
Baseline-Vorlage `templates/harness/conventions/MR-NNN-titel.template.md`
(Feld `Ersetzt-Baseline-Regel`) · `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`
(3×) · Review `review-slice-harness-baseline-v6-14-1` F-2, F-3.

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0073` und `ADR-0156`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Erste Frage — die Reichweite.** `ADR-0073` trägt die Reichweite der
Zitat-Korrektur in zwei Formen, die nicht dieselbe Menge benennen.
§Entscheidung 1 listet Abschnitte einer ADR, an denen eine Änderung „keine
Zitat-Korrektur" ist (§Entscheidung, §Konsequenzen *der Aussage nach*,
§Verglichene Alternativen, §Status, `Supersedes`-Kette …). Die Kurzform sagt:
„das Gerüst darf sich ändern, die Aussage nie; der Referent bleibt derselbe".
`AGENTS.md` §3.5 übernimmt die Liste mit „nie eine Zitat-Korrektur". Eine
Korrektur, die nur ein Versions-Segment oder den Ort eines Referenten in einem
dieser Abschnitte bewegt, ist nach der Kurzform zulässig und nach der Liste
verboten. Das ist dreimal aufgetreten (übernommen aus den Belegen des
Beobachtungs-Eintrags): `ADR-0095` §Verglichene Alternativen beim Bump auf
v6.13.0 (`00d96eb7`), `ADR-0123` §Entscheidung und §Konsequenzen
(`slice-maintainer-ordner-releasing-verschieben`), `ADR-0095` erneut beim Bump
auf v6.14.1 (`eadf3054`). Das Architect-Audit `9f1eb320`
(Audit `audit-baseline-v6-13-0-drift` §4.2) hat die Kurzform-Lesart
gezogen, für neue Fälle aber die engere verlangt, bis ein eigener Träger klärt.
Diese ADR ist dieser Träger.

Die Liste selbst ist schon in `ADR-0073` asymmetrisch: bei §Konsequenzen steht
„*der Aussage nach*", bei §Entscheidung und §Verglichene Alternativen nicht. In
allen drei Abschnitten ist die geschützte Sache dieselbe — die normative
Aussage. Ein Pfadsegment in einer Options-Tabelle ist kein Teil der
Abwägung.

**Zweite Frage — die MR-Einträge.** `.d-check.yml` Block `vcs:` hält
`harness/conventions/**/MR-[0-9]*.md` append-only: Ab der Zeile
`- **Datum:**` darf sich der Core über eine Commit-Range nicht ändern. Der Core
ist nach der Spezifikation des Moduls der rohe Inhalt ohne
`vcs.exclude-sections` und ohne Kopf-Status-Zeile; eine Normalisierung von
Teilzeichenketten gibt es nicht (übernommen aus `d-check`s
`spec/spezifikation.md` §DC-FA-VCS-001.a Schritt 4, gelesen am Stand v0.80.0
des Schwester-Repos; das gepinnte Image ist v0.79.0). Die Baseline-Vorlage für
einen MR-Eintrag verlangt dagegen im Feld `Ersetzt-Baseline-Regel` einen Link
**in die vendored Fassung mit Tag im Pfad** und sagt selbst, dass dieser Tag
bei jedem Bump wandert. Folge: Jeder Bump ändert den Core der vier aktiven
MR-Dateien. Gemessen: `make doc-immutable RANGE=11a5bac5..HEAD` (Stand
`ba0080d2`) Exit 2, vier `core-drift-vcs` an `MR-001` bis `MR-004`; der vorige
Bump (`484d20ec`) zeigt dasselbe (übernommen aus Review F-3). Schon `94de7bb5`
(Bump auf v6.9.0) hat die Pin-Umstellung in `MR-001` als Zitat-Korrektur nach
`ADR-0073` geführt; ob MR-Dateien zur Klasse gehören, steht aber nirgends.

## Entscheidung

1. **Die Reichweite richtet sich nach der Aussage, nicht nach dem Abschnitt.**
   Eine Zitat-Korrektur nach `ADR-0073` ist in **jedem** Abschnitt einer
   `Accepted`-ADR zulässig, auch in §Entscheidung, §Konsequenzen und
   §Verglichene Alternativen, wenn alle drei Bedingungen gelten:
   (a) geändert wird nur die Form eines Verweises: Pfad, Versions-Segment,
   Linkziel, Anker, Zeilen-Lokator, die Form einer gebrochenen Referenz;
   (b) der Referent ist **gemessen** derselbe: `cmp` des Zielinhalts an alter
   und neuer Adresse, und der Beleg nennt, ob roh oder normalisiert gemessen
   wurde (`AGENTS.md` §3.12);
   (c) kein Wort der Aussage ändert sich, auch keines, das den Verweis umgibt.
   Ohne Ausnahme unberührbar bleiben die Status-Zeile, die `Supersedes`-Kette,
   `Datum` und `Autor`; sie sind Aussage, kein Gerüst. Die Abschnitte-Liste
   von `ADR-0073` §Entscheidung 1 liest sich ab hier als Liste der Abschnitte,
   deren **Aussage** geschützt ist. Ändert sich der Referent selbst (anderer
   Inhalt an der neuen Adresse), ist das keine Zitat-Korrektur, sondern
   Folge-ADR.

2. **`AGENTS.md` §3.5 bekommt den passenden Wortlaut** (Folgepflicht). Der Satz
   „**Unberührbar** bleiben §Entscheidung, §Konsequenzen, §Verglichene
   Alternativen, §Status und die `Supersedes`-Kette; ihre Änderung ist eine
   neue ADR mit `Supersedes ADR-NNNN`, nie eine Zitat-Korrektur." wird ersetzt
   durch: „**Unberührbar** bleibt die Aussage von §Entscheidung, §Konsequenzen
   und §Verglichene Alternativen sowie §Status und die `Supersedes`-Kette;
   ihre Änderung ist eine neue ADR mit `Supersedes ADR-NNNN`. Das Zitatgerüst
   in diesen Abschnitten darf eine Zitat-Korrektur ändern, wenn der Referent
   gemessen gleich bleibt (`ADR-0157`)." In `AGENTS.md` steht `ADR-0157` als
   Link auf diese Datei, wie die übrigen ADR-Verweise dort.

3. **MR-Einträge gehören zur Klasse und folgen der Record-Belegform.** Ein
   Eintrag unter `harness/conventions/` (aktiv oder in `done/`) ist ab seiner
   `Datum`-Zeile immutabel; eine Zitat-Korrektur nach Entscheidung 1 ist
   zulässig. Weil die MR-Vorlage keine §Geschichte trägt, gilt die Belegform
   der Records aus `ADR-0073` Entscheidung 3/4: Die Commit-Message nennt
   `ADR-0073`, eine §Geschichte-Zeile entfällt. Die Pin-Umstellung eines
   Baseline-Bumps ist der Regelfall dieser Klasse.

4. **Ein Bump bewegt die MR-Pins in einem eigenen Commit, und der
   `doc-immutable`-Lauf umgeht genau diesen Commit.** Der Bump-Ablauf stellt
   die Baseline-Pins der MR-Dateien in einem Commit um, der **nur** MR-Dateien
   ändert und dessen Message `ADR-0073` nennt (der *Pin-Commit*). Der Verifier
   prüft eine Slice-Range `B..H`, die Pin-Commits enthält, in zwei Teilen:
   - **Teil-Ranges:** `make doc-immutable` je Abschnitt zwischen den
     Pin-Commits, also für einen Pin-Commit `P` die Ranges `B..P~1` und
     `P..H`; jede endet mit Exit 0.
   - **Pin-Commit:** je MR-Datei, die `P` ändert, `cmp` nach Normalisierung
     des Tags auf beiden Seiten; Exit 0 je Datei:

     ```bash
     norm() { sed -E 's#\.harness/baseline/v[0-9]+\.[0-9]+\.[0-9]+/#.harness/baseline/<tag>/#g'; }
     for f in $(git diff --name-only "$P~1" "$P" -- 'harness/conventions/MR-[0-9]*.md' 'harness/conventions/**/MR-[0-9]*.md'); do
       cmp <(git show "$P~1:$f" | norm) <(git show "$P:$f" | norm) || echo "kein reiner Pin: $f"
     done
     ```

   Der `cmp` ersetzt am Pin-Commit die Prüfung, die `vcs` dort nicht leisten
   kann: Er ist grün genau dann, wenn sich außer dem Tag nichts bewegt hat.
   Das ist eine Eingrenzung des `vcs`-Laufs um einzelne Commits; `AGENTS.md`
   §3.6 ist durch diese ADR erfüllt. Die Konfiguration in `.d-check.yml`
   bleibt unverändert.

5. **Für `slice-harness-baseline-v6-14-1` gilt Entscheidung 4 rückwirkend,
   mit einer benannten Lücke.** Der Pin-Commit `5d8855d9` ändert neben den vier
   MR-Dateien acht weitere Dateien, und seine Message nennt nur `ADR-0051`.
   Beides ist im festgeschriebenen Commit nicht mehr zu ändern; Historie wird
   nicht umgeschrieben. Die Teil-Ranges und der `cmp` tragen trotzdem, weil
   `vcs` nur MR-Dateien prüft und der `cmp` nur die MR-Dateien des Commits
   liest (gemessen, siehe Fitness Function). Die fehlende `ADR-0073`-Kennung
   holt der Slice-Plan als committete Zeile nach, die `5d8855d9` als
   Zitat-Korrektur an `MR-001` bis `MR-004` nach `ADR-0073` und dieser ADR
   ausweist.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

Erste Frage (Reichweite):

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; Liste und Kurzform stehen nebeneinander | kein Eingriff | derselbe Commit liest sich je nach Träger als zulässig oder als Verstoß; dreimal belegt, bei jedem Bump wieder |
| B — die Liste gilt wörtlich: jede Gerüst-Änderung in diesen Abschnitten braucht eine Folge-ADR | streng und ohne Urteil am Abschnitt ablesbar | eine Folge-ADR, die einen Tag im Pfad neu fasst, ändert nichts Normatives; `ADR-0073` Option C hat genau diesen Weg als falsches Werkzeug verworfen; jeder Bump erzeugte eine ADR je betroffener ADR |
| **C — die Aussage ist geschützt, das Gerüst nicht; Referent gemessen gleich (gewählt)** | eine Lesart für alle Träger; deckt sich mit der Kurzform, der Asymmetrie „der Aussage nach" in `ADR-0073` und dem Audit `9f1eb320` | der Kanal in `Accepted`-ADRs reicht jetzt in alle Abschnitte; Wächter bleibt das Urteil am Diff plus der gemessene Referent |

Zweite Frage (MR-Pins und `doc-immutable`):

| Option | Pro | Contra |
|---|---|---|
| D — nichts tun; `doc-immutable` ist bei jedem Bump rot, der Verifier liest die Befunde | kein Eingriff | ein Lauf, der regelmäßig rot sein darf, wird überlesen; ein echter Eingriff in einen MR-Eintrag ginge im selben Rot unter |
| E — MR-Dateien zitieren die Baseline-Regel als Text ohne vendored Pfad (Hausform wie im Index von `harness/conventions.md`) | Bumps berühren MR-Dateien nie wieder; `doc-immutable` bleibt ohne Verfahren grün | weicht von der Baseline-Vorlage ab, die den Link mit Tag ausdrücklich verlangt (eine weitere Adaption `MR-<NNN>`); verliert die Link-/Anker-Prüfung, die beim Bump meldet, wenn die ersetzte Baseline-Regel umbenannt wurde; die Umstellung selbst ist einmal rot |
| F — `vcs` aus `.d-check.yml` streichen oder den Glob leeren | grün ohne Verfahren | gibt den einzigen maschinellen Schutz gegen ein Überschreiben eines MR-Eintrags auf; Gate-Lockerung ohne Ersatz |
| G — `d-check` um eine Normalisierung im Core von `vcs` erweitern (etwa ein Regex, der wie `versions.pin-pattern` den Tag ausblendet) | sauberste Lösung, ein Lauf, grün | Arbeit in einem anderen Repo und ein neues Release, bevor hier irgendetwas grün wird |
| **H — Pin-Commit isolieren, Teil-Ranges plus normalisierter `cmp` (gewählt)** | keine neue Datei, kein neues Werkzeug, keine Konfigurationsänderung; jeder Lauf ist grün oder ein echter Befund; der `cmp` fängt einen Wortlaut-Eingriff, der im Pin-Commit mitreist (gemessen) | ein Verfahrensschritt mehr für den Verifier bei Bump-Slices; die Disziplin „Pin-Commit ändert nur MR-Dateien" ist eine Regel des Bump-Ablaufs, kein Sensor |

**Fazit:** C und H. G bleibt der Weg, falls `d-check` eine Normalisierung
liefert (Re-Evaluierungs-Trigger b).

## Konsequenzen

- Positiv: Ein Träger statt zweier Lesarten; der Beobachtungs-Eintrag
  `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` hat mit dieser ADR
  und dem Nachzug von `AGENTS.md` §3.5 seinen Ausgang *verkörpert*.
- Positiv: Die Zitat-Korrekturen an `ADR-0095` (`00d96eb7`, `eadf3054`) und
  `ADR-0123` sind nach Entscheidung 1 regelkonform; keine von ihnen wird
  zurückgenommen.
- Positiv: Ab dem nächsten Bump ist `make doc-immutable` im Verifier-Lauf grün,
  ohne dass `.d-check.yml` gelockert wird.
- Negativ mit Grenze: Der Kanal in `Accepted`-ADRs reicht jetzt in jeden
  Abschnitt. Ob nur das Gerüst bewegt wurde, entscheidet weiter das Urteil am
  Diff; maschinell gemessen ist nur der Referent (`cmp`). `ADR-0073`
  Re-Evaluierungs-Trigger (a) gilt unverändert.
- Negativ: Der Bump-Ablauf bekommt eine Disziplin („MR-Pins in einem eigenen
  Commit mit `ADR-0073`"). Bricht sie, merkt das der Verifier: ein Pin-Commit
  mit anderem Inhalt fällt im `cmp` rot aus. Ein Pin-Commit, der zusätzlich
  Nicht-MR-Dateien ändert, bleibt unerkannt und ist harmlos, weil `vcs` diese
  Dateien nicht prüft.
- Negativ: Ein `git mv` eines MR-Eintrags nach `done/` ist für `vcs` ebenfalls
  ein `core-drift-vcs` (D/R einer immutablen BASE, nach der Spezifikation des
  Moduls, hergeleitet, nicht gefahren). Diese ADR regelt das nicht; bisher gab
  es keine solche Umbenennung (`git log --diff-filter=R -- 'harness/conventions/done/*'`
  ohne Treffer). Tritt der Fall ein, gilt Entscheidung 4 sinngemäß: der
  Umzugs-Commit wird wie ein Pin-Commit umgangen und per `git diff -M` als
  reine Umbenennung belegt.
- Folgepflicht (Implementer, `slice-harness-baseline-v6-14-1`):
  (1) `AGENTS.md` §3.5 nach Entscheidung 2;
  (2) `.claude/agents/verifier.md`, Punkt `make doc-immutable`: bei
  Pin-Commits nach Entscheidung 4 verfahren;
  (3) `harness/targets/pin-stale.md`, Bump-Ablauf: die MR-Pins in einem
  eigenen Commit, der nur MR-Dateien ändert und `ADR-0073` nennt;
  (4) im Slice-Plan die Zeile nach Entscheidung 5 und die Messung der
  Teil-Ranges (Befehl und gedruckte Zeile).
  Der ADR-Index trägt diese ADR und den Verweis in §Konventionen bereits.
- Folgepflicht (Planner, Closure dieses Slice): Ausgang *verkörpert* des
  Beobachtungs-Eintrags mit Zielort `AGENTS.md` §3.5 und Anker
  `seit slice-harness-baseline-v6-14-1`, dritte Evidence-Datei.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | Entscheidung 1 (nur Gerüst, Aussage gleich) ist ein Urteil am Diff, kein Sensor; maschinell ist nur der gemessene Referent (`cmp`) | — |
| `d-check` Modul `vcs` (Image v0.79.0, Digest aus `d-check.mk`) | Teil-Ranges um den Pin-Commit sind grün. **Gemessen** am Repo, Stand `ba0080d2`: `RANGE=11a5bac5..5d8855d9~1` Exit 0, „1752 Datei(en) geprüft, 0 Befund(e)"; `RANGE=5d8855d9..HEAD` Exit 0, dieselbe Zeile; die volle Range `11a5bac5..HEAD` Exit 2, vier `core-drift-vcs`. **Erprobt** an einem Klon im Scratchpad mit einem simulierten nächsten Bump (Pin-Commit v6.14.1 → v6.14.2 an allen vier MR-Dateien, danach ein weiterer Commit): volle Range Exit 2 mit vier Befunden, Range `P..H` Exit 0 mit 0 Befunden, Range `B..P~1` bei leerer Range Exit 0 mit 0 Befunden | `make doc-immutable RANGE=…` |
| `cmp` nach Normalisierung (Befehl in Entscheidung 4) | reiner Pin-Commit ergibt `cmp` 0 je MR-Datei. **Gemessen** an `5d8855d9` und an `484d20ec`: je vier Dateien `cmp` 0. **Mutation** an einer Stelle (`MR-001`, im Klon: Pin v6.14.2 → v6.14.3 und im selben Commit `Auflösungs-Trigger: permanent` → `nie`), Instanz die Shell-Schleife aus Entscheidung 4: `cmp` 1, rot. Dass jede andere Wortänderung ebenso fällt, ist *hergeleitet* aus der Byte-Gleichheit des `cmp` | kein Make-Target (Verifier-Lauf) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** eine nach Entscheidung 1 deklarierte Korrektur
erweist sich als inhaltlich (der Kanal hat Substanz durchgelassen) — dann
Folge-ADR, die die Abschnitte wieder schließt; **(b)** `d-check` bietet im
Modul `vcs` eine Normalisierung oder Ausblendung von Teilzeichenketten im Core
— dann Option G per Konfiguration statt des Verfahrens aus Entscheidung 4;
**(c)** die Baseline-Vorlage für MR-Einträge verlangt den Tag im Link nicht
mehr — dann entfällt der Anlass von Entscheidung 4; **(d)** ein Pin-Commit
fällt im `cmp` rot aus, ohne dass jemand den Wortlaut ändern wollte (die
Normalisierung greift zu kurz) — dann die Normalisierung anpassen.
Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Accepted — Architect-Verdikt zu Review F-2 und F-3 von `slice-harness-baseline-v6-14-1` (Konflikt-Pfad, drittes Auftreten von `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`) | Review `review-slice-harness-baseline-v6-14-1` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0157` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
