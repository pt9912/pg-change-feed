# Verifikationsbericht: slice-089 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-089` §2) und die im Slice referenzierten Entscheidungen
([`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
§Entscheidung 1–6, §Die benannte Grenze, §Konsequenzen ·
[`architect-verdict-negativtest-eingabeseite-4x.md`](architect-verdict-negativtest-eingabeseite-4x.md)
§Verdikt) sowie die Hard Rules `AGENTS.md` §3.5, §3.9, §3.10, §3.11.
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe, mit
[`review-slice-089.md`](review-slice-089.md) und
[`review-slice-089-delta.md`](review-slice-089-delta.md) abgeschlossen) und
**nicht** gegen realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan am Stand `HEAD` gelesen,
dazu `ADR-0083`, das Architect-Verdikt, die beiden Review-Reports und alle
berührten Träger. Die Reports wurden als **Kontext** gelesen und ihre Findings
**nicht** nachgeprüft (Auftrag); jede Zahl dieses Berichts stammt aus eigenen,
hier gefahrenen Läufen oder aus einem Zitat, das ich gegen seine Quelle gehalten
habe. Exit-Codes je ungepiped und in eigenem Schritt gelesen (`AGENTS.md` §3.9);
Gate-Lauf und Folgehandlung getrennt beauftragt.

**Gegenstand:** `slice-089`, Stand `HEAD` = `f492b0c`, Zweig `main`, Baum sauber
(`git status --porcelain` leer). Der Slice liegt in `in-progress/`; der `git mv`
nach `done/` ist **nicht** erfolgt. Der Diff ist **Doku + Rollen-Skill +
Workflow-Schritt**: keine Testdatei, kein Produktionscode, keine
`.github/workflows/**`-Datei, `harness/README.md` **unberührt**.

**Nachtrag:** §10 beurteilt den Korrektur-Commit `7079484` zu **V-1** und
**V-2** — **kein** neuer Vollauf; §1–§9 stehen als Stand `f492b0c`.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `git diff --name-status c9c3963~1..f492b0c` | **0** | 18 Dateien: 6 `M` Doku/Träger, 1 `M` Roadmap, 1 `M` Slice-Plan, 2 `M` `state.md`, 1 `M` `welle-20.md`, 5 `A` (Register-Eintrag + Evidence), 2 `A` Reports — **keine** `D`, kein Produktionscode |
| 2 | Wörtlich-Probe: `sed -n '143,167p' ADR-0083 \| grep '^>'` → **21** Zeilen; `grep -n '^>' AGENTS.md` über §3.12 (Z. 366–423) → **21** Zeilen; `diff` | **0** | **leerer `diff`** — die beiden `> `-Blöcke sind byte-identisch (die „20 vs. 20" der Übergabe ist um die zwei leeren `>`-Zeilen daneben; in keinem Artefakt) |
| 3 | `make gates` (> Log, Exit **danach** in eigener Datei gelesen) | **0** | `baseline-verify v6.5.0 OK` (54 Dateien) · `d-check: 733 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)` · `coverage-gate: OK — Coverage 74.70% erfüllt Schwelle 70%` · `gesamt: 0 Befund(e)` |
| 4 | `make docs-check` (eigener Aufruf) | **0** | `d-check: 733 Datei(en) geprüft, 0 Befund(e)` |
| 5 | `make commit-traceability RANGE=960fee4..f492b0c` | **0** | `6 Commit(s) in "960fee4..f492b0c", Betreffs ohne Struktur-ID` |
| 6 | Proben A–D (Zahl ↔ Ursprung, §4 unten) | **0** | Ursprung ist die tragende Hälfte, die Zahl nicht — an den neuen Sätzen selbst gezeigt |
| 7 | `grep -nE` nach `62\|53 ?[–-] ?55\|229\|231\|60\|336\|308\|26–28\|198\|1903` in `ADR-0082` gegen `welle-20.md:110–129` | **0** | **jede** Zuordnung des Herkunfts-Absatzes löst auf: Cluster-Tabelle (`258–264`: `≈54`, `60`, `62`, `53 – 55`, Summe `229 – 231`), §Kontext (2) `1903` (`:62`), (3) `0 von 198` (`:80`), (4) `336`/`308`/`26–28` (`:112`, `:120–122`), (5) `1523`/`154`/`13`/`63`/`3,3 pp`/`81,2 %`/`84,1 %` (`:137–148`); die zweite Spalte `netzlos erreichbar` steht mit `60` in derselben Zeile (`:261`) |
| 8 | `grep "go list\|TestGo"` in `verify-slice-084.md`, `review-slice-084.md`, `review-slice-085.md` | **0** | **kein** Treffer für den `go list`-Beleg-Befehl; `verify-slice-084` **V-1** betrifft einen **anderen** Befehl (`git diff --name-only …` im Plan `slice-084-postgresack-naht.md` §3(b)) |
| 9 | `git show 65aead2 -- harness/sensors/coverage-gate.md` | **0** | `65aead2` („slice-079 dritte Fixrunde") **führt** den Satz „Fünf Pakete … (`go list -f '{{len .TestGoFiles}}'` …)" ein — der Beleg-Befehl ist Alt-Bestand aus `slice-079` |
| 10 | `git show b6b2ce4 --stat` | **0** | **eine** Datei: der Slice-Plan (3+/1-) — `b6b2ce4` hat `welle-20.md` **nicht** angefasst |
| 11 | `git log --oneline 93d64dc..HEAD -- welle-20.md` | **0** | **`f492b0c`** — der Träger wurde **nach** dem Delta-Review geändert |
| 12 | `ls …/BEO-PGC/<slug>/evidence/` (drei berührte Einträge) | **0** | `beleg-befehl-traegt-seinen-satz-nicht` = **2** (`slice-084.md`, `slice-085.md`) · `mechanismus-erklaerung-ohne-werkzeugbeleg` = **3** · `zahl-in-traeger-driftet-gegen-die-messung` = **5**; Register gesamt **60** Verzeichnisse |
| 13 | `grep "71.30\|THRESHOLD=75\|THRESHOLD=70" verify-slice-081.md` | **0** | beide Kommandos real gefahren (`FAIL — Coverage 71.30% unter Schwelle 75%` / `OK — Coverage 71.30% …`) — der §Ausgabe-Anker in `coverage-gate.md` löst auf |
| 14 | `.github/workflows/**` im Diff / `harness/README.md` im Diff / hinzugefügte Zeile in `AGENTS.md` §4 | **0** | 0 / 0 / keine (§4 nur Kontext) |
| 15 | `git check-ignore -v .harness/state/gates-passed.diffsha` | **0** | ignoriert (`.harness/.gitignore:6: state/`) — der Nachweis-Stempel bleibt untracked, der Baum bleibt sauber |

---

## 2. DoD-Konformität, Kriterium für Kriterium

**Liefer-Punkt 1 — die Herkunfts-Regel steht.**

| Kriterium (§2) | Befund |
|---|---|
| `AGENTS.md` trägt neuen **§3.12** mit **beiden** Instanzen (Zahlenwert · Tatsachenbehauptung), **wörtlich** wie in `ADR-0083` entschieden | **bestätigt** — §3.12 steht (`:366`), Instanz A (`:371–384`) und Instanz B (`:386–396`); Messung 2: 21 zu 21 Zeilen, **leerer `diff`** |
| Der Abschnitt nennt die **Grenze**: kein Sensor (Formpflicht auf Prosa erzeugt Pflichterfüllung); Falsifikation bleibt die **Messung** | **bestätigt** — `:398–409` („Was diese Regel nicht hat — die benannte Grenze"), beide Elemente wörtlich aus `ADR-0083` §Entscheidung 4 |
| Der Träger-Anker steht: `· seit slice-089` — dieselbe Form wie §3.11 | **bestätigt** — `:422` `… · seit slice-089.`; §3.11 endet mit `· seit slice-078.` (gleiche Hausform) |

**Liefer-Punkt 2 — die Mutations-Richtung steht an ihren zwei Trägern.**

| Kriterium (§2) | Befund |
|---|---|
| `.harness/skills/reviewer.md` trägt einen **HIGH-Unterpunkt** zur Eingabeseite | **bestätigt** — zwei neue Unterpunkte in der HIGH-Liste; der zweite („Zusage ohne Bindung an ihre Eingabeseite — ‚grün ohne Aussage'") trägt die Richtung im Wortlaut des Verdikts |
| `.claude/commands/implement-slice.md` **Schritt 19** trägt dieselbe Richtung | **bestätigt** — Schritt 19 (`:181–196`) trägt Richtung **und** Enumerations-Pflicht |
| **Kein** neuer Sensor, **keine** ADR | **bestätigt** — der Diff enthält keine neue ADR, keine `Makefile`-/`harness/mk/**`-/Sensor-Neuanlage; Verdikt und `ADR-0083` lehnen beides ausdrücklich ab |

**Liefer-Punkt 3 — die zwei bekannten Träger sind regelkonform.**

| Kriterium (§2) | Befund |
|---|---|
| `harness/sensors/coverage-gate.md` §Grenze Punkt 1 **und die zwei §Ausgabe-Abschnitte** tragen **Ursprung und Lauf** | **bestätigt** — §Grenze 1: `20 Statements, 20 gedeckt, Lauf slice-089` (`:97`), `86/61` mit `70,9 %` **als abgeleitet** gekennzeichnet (`:110–111`), `49` mit `Lauf slice-089` (`:114–115`); §Ausgabe `coverage-gate.md`: `Lauf slice-081`, beide Kommandos real gefahren (`:166–174`); §Ausgabe `db-adapter-coverage.md`: `477 von 650 = 73,38 %`, `Lauf slice-081` (Kalibrierungs-Zelle **und** Rot/Grün-Beleg) |
| `welle-20.md` unter derselben Regel geprüft — weitere bewegliche Zahlen neben dem Nenner datiert | **bestätigt** — der Herkunfts-Absatz (`:110–129`) klassifiziert jede Zahlengruppe (`übernommen`/`gemessen`/`gerechnet`); Messung 7: **alle** Locatoren lösen auf. Hinweis: dieser Absatz wurde in **`f492b0c`** geändert, also **nach** dem Delta-Review (§4) |
| `make gates` grün (Exit direkt, ungepiped) | **bestätigt** — eigener Lauf, Exit **0** (Messung 3), aus separater Datei gelesen; `make docs-check` als eigener Aufruf ebenfalls Exit **0** (Messung 4) |

**Die Closure-Pflichten aus §2 (die fünf unteren Zeilen).**

| Kriterium (§2) | Befund |
|---|---|
| Review durchgeführt, Report unter `docs/reviews/`; kein Self-Review | **bestätigt** — `review-slice-089.md` (1 HIGH/2 MEDIUM/2 INFO) und der Delta-Nachlauf liegen vor; die Zeile nennt Report, Funde und Fixrunde zutreffend |
| Reconciliation-Register fortgeschrieben *(entfällt …)* | **entfällt nachweislich** — `docs/plan/planning/reconciliation.md` existiert nicht (Greenfield-Bootstrap), wie die Zeile selbst sagt |
| Beobachtungs-Register fortgeschrieben | **Substanz bestätigt, Häkchen offen, ein Beleg fehlerhaft** — `f492b0c` hat das neue Verzeichnis angelegt, zwei `evidence/`-Dateien ergänzt und zwei `state.md` nachgezogen (Messung 12); die Kopf-/Zähler-/Beleglisten-Prüfung steht in §6, der Defekt in §4 (**V-1**) |
| Jedes Risiko aus §6 trägt einen Ausgang | **nicht erfüllt** — alle drei Risiken stehen noch auf `— **Ausgang:** <bei Closure>` (Plan `:235`, `:239`, `:243`) |
| Die drei Paarungen sind getragen | **nicht erfüllt** — laufen nach dem `git mv` (`:141`, §7 Platzhalter) |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt** — §7 ist vollständig Platzhalter (`<…>`) |

**Ergebnis:** Die **drei Liefer-Punkte und die Review-/Gate-Zeile tragen**; die
vier verbleibenden Häkchen sind die regulären Closure-Pflichten des Planners.
Zwei von ihnen (Register, Paarungen) sind **nicht** bloß Formalien: der
Register-Punkt hat mit **V-1** einen inhaltlichen Defekt, und die Paarung (c)
läuft genau gegen das Verzeichnis, das V-1 betrifft.

---

## 3. „Wörtlich" — die tragende Zusage

**Trägt, mechanisch nachgeprüft (Messung 2).** Meine eigene Probe gegen
`ADR-0083:143–167` und gegen die `> `-Zeilen aus §3.12 ergibt **21 zu 21
Zitierzeilen** und einen **leeren `diff`**. Die im Auftrag genannte „20 vs. 20"
ist eine Zählung der zwei **leeren** `>`-Zeilen (die Absatztrenner) und liegt in
der Übergabe, in keinem Artefakt.

**Die zwei Arrangements — Form, keine Inhaltsänderung:**

1. **Die Regie-Zeile fällt weg** („Der Unterabschnitt trägt:", `ADR-0083:142`) —
   sie adressiert die ADR-Leserin, sie ist kein Regel-Satz.
2. **Die Grenze ist aus zwei ADR-Stellen zusammengezogen** (§Entscheidung 4 +
   §Die benannte Grenze). Der Vergleich zeigt **drei** Abweichungen der Prosa,
   keine davon ändert eine Festlegung: `dieser Entscheidung` → `der
   Entscheidung` (Bezugspunkt statt ADR-Selbstbezug), `die er prüft` → `die er
   prüfen soll` (Grammatik), und der **gestrichene Halbsatz** „und wird nicht als
   künftiger bestellt" samt der Beleg-Klausel zu `slice-081/-084/-085`. Der
   gestrichene Halbsatz ist eine **Auslassung**, kein widersprechender Satz: der
   stehengebliebene Satz („Ein Sensor ist **nicht** Teil der Entscheidung …")
   trägt die Festlegung, und der Verweis auf `ADR-0083` steht im selben
   Abschnitt. Die Begründung des Reviews, der Halbsatz widerspreche
   `ADR-0083` §Re-Evaluierungs-Trigger (a), ist dafür **nicht nötig** — Trigger
   (a) macht einen Sensor bei geänderter Vorbedingung *möglich* („möglich, nicht
   Pflicht") und steht neben einer Nicht-Bestellung nicht in Widerspruch. Die
   Substanz trägt unabhängig davon: die **gebundene** Hälfte des DoD sind die
   zwei Wortlaute, und die sind byte-exakt.

**Ein Punkt zur Vollständigkeit des Trägers (keine DoD-Verletzung):** `ADR-0083`
§Entscheidung 4 nennt für Instanz B **zwei** durchsetzende Hälften — den
**Verifier** *und* den **Planner als Verfasser** (Tabelle, `:184`). §3.12 nannte
für Instanz B nur den Verifier (`:407–408`). Die DoD-Zeile verlangt das nicht
(„beide Instanzen, wörtlich" + Grenze), und die vollständige Fassung bleibt
einen Link entfernt; als **Träger einer Hard Rule** ist die Aufzählung der
durchsetzenden Leser dort aber unvollständig, und die Verfasser-Seite ist genau
die, die `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (4×)
adressiert. — **Erledigt in `7079484`** (§10).

---

## 4. Die Regel gegen sich selbst — mein eigener Suchlauf

Geprüft wurden **alle** von diesem Slice berührten Träger: `AGENTS.md` §3.12,
beide Skill-Unterpunkte, Schritt 19, die drei Sensor-Dokumente, `welle-20.md`
§4, der Slice-Plan und die Register-Artefakte. Die Probe ist die des `§3.12`
selbst: **Zahl ohne Ursprung** bzw. **Behauptung ohne tragenden Beleg-Anker**.

**Was trägt** (Stichprobe, jede Stelle gelesen): die neuen Herkunfts-Zeilen in
`reviewer.md` (Eintrags-Name + vier Vorgangs-IDs), in Schritt 19
(`review-slice-088.md` F-1 + `evidence/slice-088.md`), in `coverage-gate.md`
(`Lauf slice-084`/`slice-085`, `slice-081`, `slice-089`; `70,9 %` als
`abgeleitet aus 61/86` gekennzeichnet; die Rückrechnungen als „kein eigener
Lauf"), in `db-adapter-coverage.md` (`Lauf slice-081` in Zelle **und** Beleg) und
in §3.12 („die vier belegten Fälle stehen in `…/BEO-PGC/zahl-in-traeger-…`").
Der Review-Skill nennt zwei Herkunfts-Zahlen (`4×` je Unterpunkt) mit
Eintrags-Pfad und Vorgangs-Liste — der Stand ist damit datiert/belegt; die
Register-Zahl des ersten Eintrags ist inzwischen **5×** (§6), der Skill nennt
also den **schreibenden** Stand, nicht den geltenden. Das ist die Grenze der
Regel („Prosa bleibt Prosa"), benannt, kein Fund.

**Proben A–D (die tragende Hälfte ist der Ursprung, nicht die Zahl):**

| Probe | Text | Ergebnis |
|---|---|---|
| A | `coverage-gate.md:97` **ohne** die Zahl → „(…, Lauf `slice-089`)" | **hält** — der Ursprung allein trägt die Aussage |
| B | dieselbe Stelle **ohne** `, Lauf slice-089` → „(20 Statements, 20 gedeckt)" | **bricht** — eine Messung ohne Lauf ist genau der §3.12-Verstoß |
| C | `AGENTS.md:405` **ohne** den Register-Pfad, „die vier …" bleibt | **bricht** — eine Zahl ohne Ursprung |
| D | dieselbe Stelle **ohne** „vier", Anker bleibt | **hält** — die Zahl ist entbehrlich, der Anker nicht |

Der Slice fügt **keinen** Test hinzu (Vorgabe: Doku + Skill + Workflow-Schritt);
eine Mutation an einer Eingabeseite gibt es hier also nicht — die vier Proben
sind die anwendbare Form, und sie zeigen: der Träger verlangt den **Ursprung**,
nicht die Zahl.

**Darüber hinaus gefunden — zwei Stellen, die die Regel in ihrer eigenen
Einführung nicht durchsetzt** (beide in `f492b0c`, beide **nach** dem
Delta-Review, also von keiner Review gesehen):

### V-1 — Der neue Register-Eintrag belegt seinen ersten Vorgang mit einem Fund, der dort nicht steht (die Belegliste trägt den Zähler nicht)

- `pfad`: `docs/plan/planning/observations/BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht/observation.md:15–24`
  (Kopf) gegen `…/evidence/slice-084.md` und
  [`verify-slice-084.md`](verify-slice-084.md) V-1
- `befund`: Der Kopf sagte „Belegt an zwei abgeschlossenen Vorgängen" und nannte
  als erste Zeile **`slice-079`** mit „(Verifikation `verify-slice-084` V-1)";
  das `evidence/` führte stattdessen **`slice-084.md`** — und diese Datei
  behauptete ihrerseits („Vorgang: `slice-084` … Gefunden hat es der Verifier als
  `verify-slice-084` **V-1**"), die `go list`-Stelle sei die V-1 jener
  Verifikation. **Sie ist es nicht:** `verify-slice-084` V-1
  (`:301–317`) betrifft einen **anderen** Befehl in einem **anderen** Dokument
  (`git diff --name-only` im Plan `slice-084-postgresack-naht.md` §3(b));
  Messung 8 findet in `verify-slice-084.md`, `review-slice-084.md` und
  `review-slice-085.md` **keinen** Treffer für `go list`/`TestGoFiles`. Der
  `go list`-Fund ist in **`verify-slice-085` §5(b)/V-3** belegt und dort
  selbst als **Alt-Bestand aus `slice-079` (`65aead2`)** eingeordnet — der Satz
  wurde in `slice-079` *eingeführt* (Messung 9), dort aber **nicht** gefunden;
  `slice-079` ist damit der **Ursprung** der Stelle, kein Vorgang, in dem die
  Klasse beobachtet wurde. Zwei Folgen: (a) **Kopf und Belegliste stimmten
  nicht überein** — `{slice-079, slice-085}` gegen `{slice-084.md,
  slice-085.md}`; (b) der Eintrag las sich durch seine **eigene** Formulierung
  („derselbe Satz, derselbe Befehl") als Klasse *über die `go list`-Stelle* — in
  dieser Lesart hätte er **einen** tragenden Beleg (`slice-085`), nicht zwei, und
  stünde mit `2×` über der Schwelle, die er nicht erreicht.
- `verifizierbar`: ja — `grep -n "go list\|TestGo" verify-slice-084.md
  review-slice-084.md review-slice-085.md`; `sed -n '296,318p'
  verify-slice-084.md`; `sed -n '355,362p' verify-slice-085.md`;
  `ls …/beleg-befehl-traegt-seinen-satz-nicht/evidence/`
- `urteil`: **vor der Closure zu beheben.** Die **Form** des Belegs ist in
  Ordnung (Verzeichnis existiert, `evidence/` ist nicht leer, Dateinamen sind
  Vorgangs-Kennungen) — die **Paarung (c)** läuft also grün; was nicht trug,
  ist die **Zuordnung**, und die entscheidet der Lese-Schritt als **Urteil**
  (Modul 6). Auflösung: lautet die Identität des Eintrags „Beleg-Befehl trägt
  seinen Satz nicht" als **Klasse**, dann zählt `verify-slice-084` V-1 als
  *erste* Instanz (der `git diff`-Befehl) und `verify-slice-085` V-3 als zweite
  — dann müssen aber **beide** Beschreibungen (Kopf *und*
  `evidence/slice-084.md`) auf **diesen** Fund umgeschrieben werden, und
  `slice-079` darf nicht als Beleg-Vorgang dastehen. — **So umgesetzt in
  `7079484`** (§10).

### V-2 — Die neue Evidence-Datei nennt für die `welle-20`-Berichtigung einen Commit, der sie nicht trägt

- `pfad`: `docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-089.md`
  (Quellen-Zeile) gegen `git show b6b2ce4 --stat`
- `befund`: Die Quelle-Zeile sagte „`docs/plan/planning/welle-20.md` (berichtigt
  in **`b6b2ce4`**)". `b6b2ce4` hat **genau eine** Datei angefasst — den
  Slice-Plan (Messung 10); `welle-20.md` wurde in `93d64dc` berichtigt (F-2) und
  in `f492b0c` nochmals (Messung 11). Der genannte Anker trug den Satz nicht —
  die Klasse, die der Slice verkörpert, in ihrer eigenen Register-Arbeit
  (Instanz B: „nennt den **Beleg-Anker**, an dem sie geprüft wurde").
  Der zweite Teil derselben Zeile (`slice-089`-Plan §8, „berichtigt in
  `c206762`") **hält** (`c206762` = Planner-F-3). — **Behoben in `7079484`**
  (§10).
- `verifizierbar`: ja — `git show b6b2ce4 --stat`; `git log --oneline
  93d64dc..HEAD -- docs/plan/planning/welle-20.md`

**Nachrangige Beobachtung (kein Fund, für den Planner):** `welle-20.md` §4
ist in **`f492b0c`** geändert worden — einem Commit, dessen Betreff nur die
Register-Arbeit nennt. Der Delta-Review hatte für seinen Gegenstand „vier
Dateien" festgestellt und `welle-20.md` auf dem Stand `93d64dc` beurteilt; die
seitdem gelesene Fassung (Cluster-Tabelle auf §*Was daraus für die Slices folgt*
verortet, `60 von 60` als zweite Spalte erklärt) ist damit **unreviewed**. Ich
habe sie selbst gegen die Anker geprüft (Messung 7) — sie **trägt**; sie behebt
sogar genau den Locator-Rest, den der Delta-Review als **D-1** (INFO) geführt
hatte. Festgehalten, weil die Aussage „verifikationsreif" des Delta-Reports für
den Stand `b6b2ce4` galt, nicht für `f492b0c`. — **Reist als §7-Notiz mit
(Planner-Entscheidung).**

**Nachrangige Beobachtung 2:** `welle-20.md:125–129` sagt zu den zwei Werten des
widerlegten Vorschlags „`591`, `327` … **nicht** als Messung;
[`ADR-0082`](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
widerlegt sie". `327` steht dort (`:277`), `591` **nicht** — widerlegt ist der
Vorschlag als Ganzes. Die beiden Werte tragen ihren Ursprung als **Zitat**
(dieselbe Sektion nennt ihre Herkunft), die Behauptung „widerlegt sie" trägt für
`591` keinen Beleg-Anker. Kein DoD-Bezug; als Instanz-B-Rest benannt. —
**Reist als §7-Notiz mit (Planner-Entscheidung).**

**Nachrangige Beobachtung 3:** die Sichtung in §8 des Plans („Register
durchgegangen; die Zähler sind am Register nachgezählt und mit ihrem Stand
benannt") führt **fünf** der **60** Einträge, ohne das Auswahlkriterium zu
nennen — und zwei Einträge, die dieser Slice tatsächlich berührt
(`mechanismus-erklaerung-ohne-werkzeugbeleg`, `beleg-befehl-…`), stehen nicht in
der Liste. Die fünf genannten Stände stimmen mit `ls evidence/` überein; die
Auswahl ist für einen Dritten nicht reproduzierbar. — **Reist als §7-Notiz mit
(Planner-Entscheidung).**

---

## 5. Die zwei Erweiterungen über die §3-Liste

**`.harness/skills/reviewer.md` mit zwei Unterpunkten — trägt, kein Überlauf.**
`ADR-0083` §Entscheidung 4 nennt als Zielort für Instanz A ausdrücklich
„`AGENTS.md` §3.12 **+ ein benannter HIGH-Unterpunkt in
`.harness/skills/reviewer.md`**" (`:183`), §Konsequenzen wiederholt es als
Folgepflicht, und die ADR bestellt ausdrücklich **keinen** Sensor — ohne diesen
Unterpunkt hätte die Zahlen-Hälfte keinen Leser. Der Plan führt den Nachzug in §3
(„Nachzug im ersten Implementer-Lauf … drei Punkte über die erste Liste hinaus")
und ist damit plan-konform. Die Mutations-Hälfte stammt aus dem Verdikt, das
`reviewer.md` und Schritt 19 als **beide** Träger bestimmt.

**`harness/sensors/db-adapter-coverage.md` mitgezogen — trägt, kein Überlauf.**
`verify-slice-085` V-1 §4.2 nennt die zwei §Ausgabe-Abschnitte **namentlich mit
ihren Werten** (`73,38 %` in `db-adapter-coverage.md`, `71,30 %` in
`coverage-gate.md`). Die Lesart, die V-1 §4.2 ganz abdeckt, ist die einzige, die
beide Dokumente erfasst; ohne den Nachzug bliebe einer der zwei namentlich
genannten Belege stehen. Die LP3-Zeile („§Grenze Punkt 1 **und die zwei
§Ausgabe-Abschnitte**") ist damit **erfüllt**, nicht übererfüllt.

---

## 6. Das Register und die Klassen-Naht

**Das neue Verzeichnis ist vollständig** (Messung 12): `observation.md`,
`state.md` und `evidence/` mit **zwei** Dateien (`slice-084.md`, `slice-085.md`)
— Paarung (c) erfüllt (Existenz + nicht-leeres `evidence/`). Inhaltlich war die
Zuordnung der ersten Datei fehlerhaft (V-1; in `7079484` berichtigt, §10).

**Die zwei nachgezogenen Zähler — Kopf, Zähler-Zeile und Belegliste stimmen
überein** (jede Zahl gegen `ls evidence/`):

| Eintrag | Kopf (`state.md`) | Zähler-Zeile | `evidence/` | Befund |
|---|---|---|---|---|
| `mechanismus-erklaerung-ohne-werkzeugbeleg` | `offen (3×, Schwelle erreicht)` | `2× …, dazu der dritte … — **3×**` | 3 | **stimmig** |
| `zahl-in-traeger-driftet-gegen-die-messung` | `offen — 5× erreicht` | `**5×** (…, slice-089.md)` | 5 | **stimmig** |
| `beleg-befehl-traegt-seinen-satz-nicht` (neu) | `offen (**2×**)` | `**2×** (evidence/slice-084.md, evidence/slice-085.md)` | 2 | Zahlen stimmen mit der Liste; die **Liste** selbst trug V-1 (in `7079484` berichtigt, §10) |

Der Zähler ist in allen drei Einträgen als **abgeleitet** gekennzeichnet und
steht unmittelbar neben seiner Belegliste — die Form vermeidet die „zweite
Quelle" (Modul 6: der Zähler wird abgeleitet, nicht geführt). Die *immutable*
`observation.md` des `mechanismus`-Eintrags sagt weiter „Belegt zweimal" und
nennt die zwei Vorgänge — sie ist ab Anlage unveränderlich und trägt ihren
Ursprung; der geltende Stand lebt in `state.md`. Kein Fund.

**Die Klassen-Naht (Delta-Review D-2) — die Aufteilung trägt.** F-1 (die
`71.94 %`-Zeile) ist `mechanismus-erklaerung-ohne-werkzeugbeleg` zugeordnet,
F-2/F-3 dem `zahl-in-traeger-…`. Ich beurteile das als die **saubere**
Auflösung:

- **F-1 gehört zu genau einem Eintrag.** Der Fund ist eine **Mechanismus-Aussage
  über ein Werkzeug** („`go tool cover` druckt eine Nachkommastelle, das Skript
  formatiert `%.2f` — daraus werden `71.9%` und `71.94%`"), die nicht am
  Werkzeug geprüft war. Genau das ist die Definition des `mechanismus`-Eintrags,
  und es ist **dieselbe §Zählbasis-Stelle**, die dieser Eintrag schon mit
  `slice-079` belegt („die gedruckte Prozentzeile als **eigene Größe**"). Die
  Zuordnung dorthin ist damit die präzisere und vermeidet die Doppelzählung.
- **F-2 und F-3 tragen die `zahl-in-traeger`-Hälfte eigenständig** (beide
  MEDIUM, beide in Trägern dieses Diffs): F-2 nennt Werte als Werte des
  Abschnitts, die die Datei nicht führt, F-3 zählt gegen die eigene Liste — die
  Klasse in ihrer reinsten Form, und F-3 ist die Wiederholung der Form, die
  `slice-088` in denselben Eintrag eingetragen hat.
- **Kein Fall zählt unter zwei Klassen:** `slice-085` V-3 (der `go list`-Befehl)
  ist dem **Beleg-Befehl**-Eintrag zugeordnet und steht **nicht** im
  `mechanismus`-Eintrag (dessen dritte Evidence ist `slice-089`, nicht
  `slice-085`). Die drei Einträge teilen sich die vier Funde ohne Überlappung —
  **außer** in V-1, wo die *Beschreibung* einen Fund in zwei Vorgänge legte, die
  ihn nicht tragen (in `7079484` berichtigt, §10).

---

## 7. Kein Gate, das es nicht gibt — und `AGENTS.md` §3.10

- **`AGENTS.md` §4 und `harness/README.md` §Sensors sind unberührt:** die
  Datei-Liste des Diff-Ranges (Messung 1) führt `harness/README.md` nicht, und
  der einzige `AGENTS.md`-Hunk (`@@ -363,6 +363,64 @@`) endet mit der Anker-Zeile
  `· seit slice-089` plus Leerzeile — im **§4**-Bereich ist **keine** Zeile
  hinzugefügt oder geändert (Messung 14). Das Sync-Gate aus `ADR-0084` wird in
  keiner der beiden Stellen genannt — die Sensors-Zeile entsteht mit dem Gate.
  (Der Hunk in `7079484` liegt in §3.12, nicht in §4 — §10.)
- **§3.10 greift nicht:** der Diff berührt **keine** `.github/workflows/**`-Datei
  (Messung 14). Kein Post-Push-Lauf nötig, kein Workflow-§6-Risiko.
- **Die im Rollen-Beschrieb genannten Sensoren `make doc-commits` und
  `make doc-immutable` existieren in diesem Repo nicht** (kein Target im
  `Makefile`/`harness/mk/**`; die `vcs`-Modul-Konfiguration in `.d-check.yml`
  wird von `make docs-check` nicht eingeschaltet). Ich habe **kein** Target
  aufgerufen, das es nicht gibt; die Traceability je Commit trägt
  `make commit-traceability RANGE=…` (Messung 5), und `make doc-immutable`
  bleibt eine benannte, nicht gefahrene Lücke (§8).

---

## 8. Was ich nicht prüfen konnte

- **Die Mutations-Richtung wirkt nicht hier.** Der Slice fügt keinen Test hinzu;
  die Wächter der neuen Regel sind Lese-Proben, und `ADR-0083` lehnt einen Sensor
  ausdrücklich ab. Es gibt daher **keine** Eingabeseiten-Mutation zu fahren. Was
  ich stattdessen gefahren habe, sind die Proben A–D (§4) — sie zeigen, dass der
  Ursprung und nicht die Zahl die tragende Hälfte ist. **Ob** die neue
  Formulierung in der Praxis greift (ob ein künftiger Reviewer die Klasse findet
  und ein künftiger Verifier die Beleg-Anker prüft), bleibt der Beobachtung
  überlassen; sie ist als `ADR-0083` §Re-Evaluierungs-Trigger (b) benannt.
- **`make doc-immutable` (MR-Immutabilität):** kein Target in diesem Repo, kein
  Sensor für MR-Dateien — ich habe die MR-Kerne daher **nicht** gegen ihren
  Commit-Range geprüft (der Slice berührt ohnehin keine `MR-*`-Datei).
- **Die „Wörtlich"-Probe ist eine Textprobe.** Sie zeigt Byte-Gleichheit der
  Zitate, nicht die Fortgeltung der Entscheidung in jeder Auslegung; ob die
  zusammengezogene Grenze in einem späteren Streitfall als „unvollständig"
  gelesen wird, ist eine Urteils-, keine Messfrage.
- **Die Geltung des `Lauf slice-089`-Ankers** hängt daran, dass der Slice
  schließt; heute ist er ein laufender Vorgang (der Anker ist benannt und die
  Werte sind durch die Review-Messungen gedeckt, aber noch nicht durch einen
  abgeschlossenen Vorgang).
- **Die Wirkung der Skill-Unterpunkte auf künftige Review-Läufe** ist
  naturgemäß nicht prüfbar; `ADR-0083` §Re-Evaluierungs-Trigger (b) benennt
  genau diese Lücke.

---

## 9. Verdikt (Stand `f492b0c`)

**Die drei Liefer-Punkte tragen.** §3.12 steht mit **byte-identischen** Wortlauten
beider Instanzen (Messung 2), die Grenze ist benannt, der Träger-Anker gesetzt;
die Mutations-Richtung steht an **beiden** Trägern mit der Hausform des Verdikts;
die zwei bekannten nicht-regelkonformen Träger sind nachgezogen, `welle-20.md`
ist herkunftsklar (`make gates` **Exit 0**, `make docs-check` **Exit 0**,
`make commit-traceability` über die sechs Commits **Exit 0**). Die zwei
Erweiterungen über die §3-Liste tragen (ADR-Folgepflicht bzw. V-1 §4.2), sie sind
**kein** Überlauf. §4 und `§Sensors` sind unberührt, §3.10 ist nicht ausgelöst.

**Zwei Stellen aber setzen die Regel in ihrer eigenen Einführung nicht durch:**
**V-1** (der neue Register-Eintrag belegte seinen ersten Vorgang mit einem Fund,
der in `verify-slice-084` V-1 **nicht** steht, und Kopf und Belegliste nannten
verschiedene Vorgänge — die Belegliste trug den `2×`-Zähler nicht) und **V-2**
(eine neue Evidence-Datei nannte `b6b2ce4` als Commit der `welle-20`-Berichtigung,
der ihn nicht trug). Beide lagen in **`f492b0c`**, also in Trägern, die nach
dem Delta-Review entstanden sind und die **keine** Review gesehen hat — die Art
Verstoß, die für Tests und für ein diff-skopiertes Review unsichtbar ist. **Beide
sind in `7079484` berichtigt und in §10 beurteilt.**

**Closure-Fähigkeit:** Die **Liefer-Punkte** sind verifiziert und die
DoD-Häkchen LP1–LP3, Gate und Review zu Recht gesetzt. Offen waren zwei Punkte:
**V-1/V-2** (in §10 abgeschlossen bzw. mit einem Rest) und die vier
**Closure-Pflichten** (Closure-Notiz §7, die drei §6-Risiko-Ausgänge, die drei
Paarungen, das Register-Häkchen). Nach Behebung des V-3-Restes aus §10 und den
Closure-Pflichten ist der Slice `done/`-fähig.

---

## 10. Delta-Nachlauf — die Korrekturen aus `7079484` (2026-09-16)

**Kein neuer Vollauf:** beurteilt wurden **nur** die sechs Dateien des
Korrektur-Commits `7079484` (fünf Register-/Regel-Artefakte, `AGENTS.md`)
gegen **V-1** und **V-2** aus §4. §1–§9 bleiben der Stand `f492b0c`.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| D-1 | `git diff --name-only fb6adf6..4035ee7` | **0** | **5 Pfade**: der Plan, `harness/image-hash.txt`, die drei `postgresack`-Dateien — die Beschreibung der neuen `evidence/slice-084.md` trifft wörtlich |
| D-2 | `git diff --name-only fb6adf6..4035ee7 -- internal/ ':!internal/adapters/driven/postgresack/'` | **0** | **leer** (0 Pfade) — „die Gegenrichtung, leer" trägt |
| D-3 | `git show a3cb2c4 --stat` | **0** | 1 Datei, `slice-084-postgresack-naht.md`, Betreff „Beleg-Range in §3 nachgetragen" — „Range nachgetragen, Pathspec fehlt" trägt |
| D-4 | `git log --oneline -S':!internal/adapters/driven/postgresack/' --all -- 'docs/plan/planning/*/slice-084-postgresack-naht.md'` | **0** | **`4387523`** — die §3(b)-Berichtigung kam im **slice-084-Closure-Commit**; **`11ba45b`** ist es **nicht** |
| D-5 | `git show 11ba45b --stat` | **0** | **eine** Datei: `slice-088-coverage-tail-uebersetzung.md` (33+/24-) — ein slice-088-Commit |
| D-6 | Wörtlich-Probe am neuen Stand (Zitat-Zeilen aus §3.12 gegen `ADR-0083:143–167`) | **0** | **21 zu 21, leerer `diff`** — die Planner-Zusage steht **außerhalb** der Zitate, beide Wortlaute bleiben byte-identisch |
| D-7 | `make gates` (> Log, Exit in eigener Datei) | **0** | `baseline-verify` OK · d-check **734** Dateien / **0** Befunde · commit-traceability OK · coverage-gate OK 74.70% |
| D-8 | `ls …/beleg-befehl-…/evidence/` + Kopf/Zähler-Zeile | **0** | 2 Dateien; `Zustand: offen (**2×**)` und `Zähler (abgeleitet): **2×**` — **stimmig** |

### V-1 — **trägt**, mit einer Ausnahme in derselben Datei

Die neuen Beschreibungen sind gegen die Quellen richtig: beide Instanzen sind
jetzt **zwei verschiedene Befehle** (`git diff` ohne Pathspec | `go list` ohne
zweites Test-Datei-Feld), jede mit ihrem eigenen Fund-Beleg (D-1, D-2, D-3), und
**Ursprung ≠ Vorkommen** ist ausdrücklich benannt (`slice-079`/`65aead2` als
Ursprung, `slice-085` als Vorkommen). **Der `2×`-Zähler hält damit**: er ruht
auf zwei Befehlen in zwei abgeschlossenen Vorgängen, nicht auf einer doppelten
Beschreibung desselben Fundes. Kopf, Zähler-Zeile und Belegliste sind stimmig
(D-8). Die Entry-Identität ist dabei zulässig angefasst: `Bezeichnung` und
`Sub-Area` — die nach Modul 6 unveränderliche Hälfte — sind unberührt.

**V-3 (offen): ein falscher Commit-Anker in der korrigierten `evidence/slice-084.md`.**
Die Quelle-Zeile nennt „`docs/plan/planning/done/slice-084-postgresack-naht.md`
§3(b) (**berichtigt in `11ba45b`**, nachgeprüft)". `11ba45b` hat **genau eine**
Datei angefasst — den `slice-088`-Plan (D-5); die §3(b)-Berichtigung (Pathspec
**und** Gegenrichtung) kam in **`4387523`**, dem slice-084-Closure-Commit (D-4).
Es ist **dieselbe Klasse wie V-2** (Beleg-Anker trägt seinen Satz nicht), eine
Datei weiter — in dem Artefakt, das V-1 behebt. Ein-Zeilen-Korrektur:
`11ba45b` → `4387523`; der Rest der Zeile hält (Zielort existiert im `done/`,
V-1-Zitat stimmt).

**Klein, nicht blockierend (beides optional):** `evidence/slice-085.md` sagt, der
Implementer habe den Fund „mit **Adresse**" liegen gelassen — `verify-slice-085`
(`:252`) nennt den **Weg** („eigener kleiner Zug"), die konkrete Kennung entstand
erst in `b6b2ce4` (der Review `review-slice-089` F-5 hielt fest, dass eine
Adresse **fehlte**); „Adresse" ist also einen Grad zu stark. Und in
`observation.md` steht der Gegenrichtungs-Befehl verkürzt
(`-- internal/:!…/postgresack/`) — als Zitat nicht wörtlich ausführbar; die volle
Form trägt die Evidence-Datei.

### V-2 — **trägt**

Die Quelle-Zeile nennt jetzt `welle-20.md` „(berichtigt in `93d64dc` und
`f492b0c`)" — beide Commits haben den Träger real angefasst, der zweite ist der
hier selbst gemessene (Messung 11). Der Beleg-Anker trägt seinen Satz.

### Der §3.12-Nachtrag — **trägt**, und mein §3-Punkt ist geschlossen

Der Abschnitt nennt für Instanz B jetzt auch den **Planner als Verfasser**
(„der sie beim Schreiben an ihren Anker bindet") — genau die zweite Hälfte aus
`ADR-0083` §Entscheidung 4 Tabelle. Die Ergänzung liegt **außerhalb** der zwei
Zitate (D-6: 21 zu 21, leerer `diff`), das DoD-Kriterium „wörtlich" ist damit
unberührt, und der Hunk liegt in §3.12, nicht in §4 (§7 bleibt gültig).

**Verdikt des Deltas:** **V-1 und V-2 tragen**; die vom Korrektur-Commit
behauptete Zuordnung ist an ihren eigenen Befehlen nachgemessen und richtig, und
der `2×`-Zähler des neuen Eintrags hält. **Offen bleibt allein V-3** (der
`11ba45b`-Anker) plus die vier Closure-Pflichten und die zwei vereinbarten
§7-Notizen. Nach V-3 ist der Slice `done/`-fähig.
