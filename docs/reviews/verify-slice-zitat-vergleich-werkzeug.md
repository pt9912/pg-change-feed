# Verifikationsbericht: slice-zitat-vergleich-werkzeug — 2026-10-06

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“. Geprüft
wird gegen die DoD von `slice-zitat-vergleich-werkzeug` (§2, Liefer-Punkte 1 bis
3 und die Gate-Pflicht). Geprüft wird außerdem gegen
[`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
(Semantik), [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md),
das [Architect-Verdikt](architect-verdict-zitat-vergleich-werkzeug.md) §1 bis §5
und [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7 und §3.12. Den Diff als
Maintainability-Frage prüft der Reviewer:
[`review-slice-zitat-vergleich-werkzeug`](review-slice-zitat-vergleich-werkzeug.md)
und die Re-Reviews
[`-fixrunde`](review-slice-zitat-vergleich-werkzeug-fixrunde.md),
[`-fixrunde-2`](review-slice-zitat-vergleich-werkzeug-fixrunde-2.md) und
[`-fixrunde-3`](review-slice-zitat-vergleich-werkzeug-fixrunde-3.md). Der
reale Bedarf ist Sache des Validators; dies ist kein MVP-Slice.

**Gegenstand:** Diff `2d21f9ab..d25c79c3` mit 14 Commits, darunter der
Architect-Commit `728b75e3` (Verdikt) und die Code-Commits `f13b7f9d`,
`adf9c0c1`, `33e1a81e` und `967564e6`. Der Slice liegt in `in-progress/`. Die
vier Closure-Punkte der DoD stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat Folgendes gelesen: den Plan, das
Verdikt, den letzten Re-Review-Report, das Skript, den Vertrag
`harness/targets/zitat-vergleich.md` und den Diff der Träger. Keine Behauptung
wurde übernommen, außer wo „übernommen“ steht. Gemessen wurde am Stand
`d25c79c3` bei sauberem Arbeitsbaum. Die Exit-Codes sind direkt und ohne Pipe
gesichert (`AGENTS.md` §3.9). Die Probe-Commits und die Mutanten liegen nur im
Scratchpad (`<Scratchpad>/verify-werkzeug/`); die Probe-Commits entstanden in
einem `git clone`. Dateien wurden dort per `sed … > <Datei>` und `mv` geändert.
Außer diesem Bericht wurde keine Repo-Datei geschrieben.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile, gekürzt) | Exit |
|---|---|---|
| `make test-zitat-vergleich` | `run-zitat-vergleich-tests: 297 Fälle bestanden (je Runde 99, Runden: ohne Option, nullglob, failglob)`, 0 Zeilen `FEHLER` | 0 |
| Pflichtproben LP1 mit `make zitat-vergleich` am Klon (Abschnitt 2) | wie Soll | je unten |
| Referent-Messung am realen Bump `5d8855d9` (Abschnitt 3) | drei `cmp 0`, `MR-001` keine Einheit | je unten |
| Zwei eigene Mutationen, Tabellentest mit `PROG=<Kopie>` (Abschnitt 4) | beide rot | 1 / 1 |
| Gleichstand ADR-Form gegen Skript, eigene Aufzählung (Abschnitt 5) | `abweichend=4`, dieselben vier Anker wie im Plan | — |
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 22 Zeilen stimmen` | 0 |
| `make docs-check` | `d-check: 1783 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=2d21f9ab..HEAD` | `0 Befund(e)` | 0 |
| `make doc-immutable RANGE=2d21f9ab..HEAD` | `0 Befund(e)`; die Range enthält keinen Pin-Commit (kein Commit ändert eine MR-Datei) | 0 |
| `make kommentar-kennungen PATHS="tools/harness/zitat-vergleich.sh tools/harness/run-zitat-vergleich-tests.sh Makefile"` | kein Kandidat | 0 |
| `git diff 2d21f9ab..HEAD -- docs/plan/adr` | 0 Zeilen: `ADR-0158`, `ADR-0159` und der Index sind unberührt (Verdikt §5) | 0 |
| `make gates` nach dem Commit dieses Berichts | Abschnitt 8 | Abschnitt 8 |

---

## 2. Liefer-Punkt 1 — Werkzeug und Pflichtproben

Ich habe die Pflichtproben selbst gefahren, am Klon zum Stand `d25c79c3`, also
nach allen drei Fixrunden. Der Plan belegt sie am Stand `f13b7f9d`. Für die
Probe gibt es zwei Commits, die nur im Klon existieren:

- `4926780e`: drei Zeilen vor Zeile 41 von `F`.
- `d7438d64`: `M` Zeile 272, „demselben Steering-Loop“ wird zu „demselbigen
  Steering-Loop“.

Dabei gilt `F` = `docs/plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md`
und `M` = `.harness/baseline/v6.14.1/regelwerk/modul-13-quality-gates.md`.

| Probe | Aufruf (`ARGS=`) | gedruckte Zeile (gekürzt) | Soll |
|---|---|---|---|
| Anker-Wechsel fällt | `d25c79c3 $F '#teilfrage-2--subjekt--und-nachrichtenschema' d25c79c3 $F '#teilfrage-3--zustellsemantik-core-nats-vs-jetstream'` | `vergleich roh: …#teilfrage-2-… <-> …#teilfrage-3-… cmp 1` | ja |
| Gegenprobe: gleicher Anker, Zeilen davor eingefügt | `… '#teilfrage-2-…' 4926780e $F '#teilfrage-2-…'` | `cmp 0` | ja |
| verschobener Lokator besteht | `d25c79c3 $F L139-145 4926780e $F L142-148` | `vergleich roh: …mdL139-145 <-> …mdL142-148 cmp 0` | ja |
| Lokator nicht nachgezogen | `d25c79c3 $F L139-145 4926780e $F L139-145` | `cmp 1` | ja |
| `#guard-haertung` mit geändertem Wort fällt | `d25c79c3 $M '#guard-haertung' d7438d64 $M '#guard-haertung'` | `vergleich roh: …#guard-haertung <-> …#guard-haertung cmp 1` | ja |
| Gegenprobe: Zeile der `id` allein | `… $M L266-266 d7438d64 $M L266-266` | `cmp 0` | ja |
| Gegenprobe: Änderung nur in `F` | `… $M '#guard-haertung' 4926780e $M '#guard-haertung'` | `cmp 0` | ja |
| ungequotetes `#…` in `ARGS` | `#guard-haertung d25c79c3 $M` | `Aufruf: … — 0 Argumente, erwartet 6 oder 7; …, Exit 2` | ja |

Der Exit von `make` war bei jeder Zeile mit `cmp 1` oder `Exit 2` gleich 2,
sonst 0. Das entspricht Vertrag §Ausgabe: die Farbe steht in der Zeile.
Bei `cmp 1` druckt `cmp` vor der Zeile die Stelle des ersten Unterschieds
(`/dev/fd/63 /dev/fd/62 differ: …`). Der Vertrag nennt das.

**Form gegen Verdikt §5:**

- Das Skript liegt unter `tools/harness/zitat-vergleich.sh` und enthält die
  Funktionen `einheit`, `tagnorm`, `vergleich`, `baumda` und `main`.
- `make zitat-vergleich` bricht ohne `ARGS` mit `$(error …)` ab.
- Beide Ziele stehen neben `suchlauf-nachmessen`. Keines steht in
  `GATE_CHECKS`; per `grep` gelesen.
- Die Argumentzahl ist auf 6 oder 7 begrenzt.
- `LC_ALL=C.UTF-8` und die Fähigkeitsprobe stehen vor der ersten Messung.
- Ein gawk-eigenes Konstrukt (`gensub`, `PROCINFO`, `asorti`) gibt es im
  Skript nicht; gelesen.

**Verdikt LP1:** bestätigt.

---

## 3. Referent-Messung am realen Bump `5d8855d9`

Gemessen im Repo mit `make zitat-vergleich`. `R0` steht für
`.harness/baseline/v6.14.0/regelwerk`, `R1` für
`.harness/baseline/v6.14.1/regelwerk`. Verglichen wird `5d8855d9~1` gegen
`5d8855d9`. Die Verweise sind dem Diff des Pin-Commits entnommen.

| MR | Anker | gedruckte Zeile (gekürzt) |
|---|---|---|
| `MR-001` | `R*/grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument` | `einheit: leere Einheit 5d8855d9~1:…` / `vergleich: 5d8855d9~1:…#spec-straten-mehr-als-ein-spec-dokument keine Einheit, Exit 2` — nicht messbar, wie `ADR-0159` Entscheidung 2 |
| `MR-002` | `R*/grundlagen-source-precedence.md#vergabe-woher-die-nächste-kennung-kommt` | `vergleich roh: … cmp 0` |
| `MR-003` | `R*/grundlagen-durchsetzungsschicht.md#grenzen--ehrlich-benannt` | `vergleich roh: … cmp 0` |
| `MR-004` | `R*/modul-13-quality-gates.md#guard-haertung` | `vergleich roh: … cmp 0` |

Dazu `R0/grundlagen-begriffe.md` am Stand `11a5bac5` gegen
`R1/grundlagen-begriffe.md` am Stand `625ddbef`, je als ganze Datei:

- roh: `cmp 1`;
- mit `v6.14.0:v6.14.1`: `vergleich norm v6.14.0:v6.14.1: … cmp 0`;
- mit `v0.79.0:v0.80.0`:
  `vergleich: Tag-Paar v0.79.0:v0.80.0, 11a5bac5 trägt .harness/baseline/v0.79.0 nicht, Exit 2`.

Alle Zeilen stimmen mit dem Plan überein. Einen weiteren nicht messbaren
Referenten als `MR-001` gibt es nicht. `MR-001` ist bereits gemeldet und
entschieden (Closure von `slice-zitat-korrektur-vergleichseinheit`).

---

## 4. Liefer-Punkt 2 — Tabellentest und eigene Mutationen

Der Tabellentest ergibt 297 Fälle, je Runde 99, alle grün (Abschnitt 1). Die
Fallgruppen F-1, F-2, F-3, F-5 und F-6 sowie die Zusatzfälle aus Verdikt §5
stehen im Test; per `grep` stichprobenartig gelesen. Zwei eigene Mutationen der
Eingabeseite habe ich gefahren, je auf einer Kopie im Scratchpad. Instanz ist
der Tabellentest mit `PROG=<Kopie>`. Keine der beiden steht in den
Mutationstabellen des Plans.

| Mutation (Stelle in `tools/harness/zitat-vergleich.sh`) | gesehene Farbe |
|---|---|
| V1: Tag-Paar-Form ohne `$` am Ende der Regex (Z. 179), damit wird ein Suffix hinter dem neuen Tag angenommen | Lauf Exit 1; „Tag-Paar mit Suffix“ `FEHLER` in allen drei Runden |
| V2: `tagnorm` ohne `g` (Z. 161), es wird nur das erste Segment je Zeile normalisiert | Lauf Exit 1; „Bump mit Tag-Paar“ und „Bump, ganze Datei mit Tag-Paar“ `FEHLER` (Exit 1 statt 0) in allen drei Runden |

**Zu V1:** Der Exit des Mutanten bleibt 2, weil `baumda` den Tag `v6.14.1x` am
Stand nicht findet. Rot wird der Fall allein über den gebundenen Meldungstext.
Die Prüfung der Form hat also eine fail-closed Rückfallebene. Ihre eigene
Bindung hängt am Text der Meldung, nicht am Exit. Das ist kein Mangel, nur eine
Beobachtung (INFO).

Die Mutationstabellen des Plans (§2, §3 Fixrunde bis Fixrunde 3) sind
**übernommen** und nicht einzeln nachgefahren. Gemessen sind in diesem Lauf nur
V1 und V2.

**Verdikt LP2:** bestätigt.

---

## 5. Gleichstand mit der ADR-Form

Die ADR-Form habe ich selbst gezogen: `sed -n '163,218p'` aus `ADR-0159`, die
Einrückung von drei Zeichen entfernt. Das ergibt 56 Zeilen, `bash -n` endet mit
Exit 0. Die Anker habe ich **eigenständig** aufgezählt: je Heading-Slug (mit
Dublettensuffix, ohne Rücksicht auf Fences) und je `<a id="…"` in allen
getrackten `.md` unter `.harness/baseline/v6.14.1/regelwerk/`, `docs/plan/adr/`,
`harness/` und `AGENTS.md`. Je Anker habe ich Ausgabe samt stderr und Exit von
`einheit` beider Formen verglichen, je per `source` geladen, unter
`LC_ALL=C.UTF-8`. Gedruckt:

```text
stand=967564e6 anker=2133 abweichend=4 leer_alt=24 leer_neu=28
stand=d25c79c3 anker=2133 abweichend=4 leer_alt=24 leer_neu=28
```

Die vier Abweichungen sind an beiden Ständen dieselben wie im Plan:

- `ADR-0159#[^`
- `harness/conventions.md#mr-<NNN>`
- `harness/targets/zitat-vergleich.md#X`
- `harness/targets/zitat-vergleich.md#x`

Alle vier gehören zu den gewollten Punkten des Vertrags §Abweichungen:
eingerückter Fence, `id` im HTML-Kommentar, `id` im Code-Span mit Text davor.

Meine Absolutzahlen unterscheiden sich von denen des Plans
(`anker=2116 … leer_alt=7 leer_neu=11`), weil meine Aufzählung auch Headings in
Fences als Kandidaten nimmt. Die Zahlen des Plans sind deshalb **übernommen**.
Gemessen und gleich ist die Menge der Abweichungen.

**Risiko „Zwei Träger“ (§6):** Der Ausgang *entfallen* nach Verdikt §2 trägt.
Außerhalb der gewollten Punkte gibt es keine Abweichung.

---

## 6. Liefer-Punkt 3 — Träger

Gelesen im Diff, je gegen die ADR und das Verdikt:

| Träger | Befund |
|---|---|
| `AGENTS.md` §3.5 | Der Messsatz nennt `make zitat-vergleich`, den Vertrag und Verdikt §3 als Auslegung. Der Beleg-Satz nennt die „gedruckte Zeile von `make zitat-vergleich`“. Die Semantik bleibt bei `ADR-0158`/`ADR-0159`. Konform mit Verdikt §2 und §5. |
| `.claude/agents/verifier.md` | Die Referent-Messung läuft „mit `make zitat-vergleich`“, mit Vertrag und „Einheit nach `ADR-0159`“. Der Zeiger auf die Befehlsform ist entfernt. Konform. |
| `.claude/agents/implementer.md` | Die Referent-Messung läuft mit `make zitat-vergleich`, mit `ADR-0159` Entscheidung 2 und dem Vertrag. Konform. |
| `harness/targets/pin-stale.md` §MR-Pins | Nennt das Ziel und den Vertrag; Beleg ist die gedruckte Zeile. Konform. |
| `harness/README.md` | Zwei Zeilen in der Werkzeug-Tabelle (`make zitat-vergleich`, `make test-zitat-vergleich`), beide „kein Gate“. Konform mit `AGENTS.md` §4, denn beide Ziele existieren im `Makefile`. |
| Vertrag `harness/targets/zitat-vergleich.md` | Er nennt: Aufruf mit Quotierung, Ausgänge, Farbe in der Zeile, die zehn Abweichungen vom Block der ADR, die Host-Werkzeuge (`bash`, `git`, `awk`, `sed`, `cmp`; im Test zusätzlich `env`, `grep`, `mktemp`, `mv`) und die Grenzen. Die Host-Werkzeuge liegen in der Klasse von `AGENTS.md` §3.1. |

Kennungen im Kommentar (`AGENTS.md` §3.7): Der Kopf des Skripts trägt eine
Kennung (`ADR-0159`), `make kommentar-kennungen` findet keinen Kandidaten.
Die Plan-Aussagen tragen Anker oder Ursprung (§3.12). Die Mutationstabellen
nennen Stelle, Instanz und Farbe. „*hergeleitet*“ steht an der Locale-Frage
(§6) und an der Verallgemeinerung über die mutierte Stelle hinaus.

Der Suchlauf (§3.13) ist nachgemessen: 22 Zeilen stimmen. Eine Meldung an den
Planner steht im Plan mit Frist zur Closure:

- `docs/plan/adr/README.md` Z. 183;
- die zwei `state.md`-Einträge im Beobachtungs-Register.

**Verdikt LP3:** bestätigt.

---

## 7. Plan gegen Code-Diff

| Plan §3 | Code-Diff | Abgleich |
|---|---|---|
| Architect-Zug, falls nötig | `728b75e3`, Verdikt ohne ADR | wie geplant |
| `tools/harness/zitat-vergleich.sh` | neu, 213 Zeilen | wie geplant; die Erweiterungen über den Plan hinaus sind in der Nachzug-Tabelle und den Fixrunden benannt |
| `tools/harness/run-zitat-vergleich-tests.sh` | neu, 314 Zeilen | wie geplant |
| `Makefile` | +9 Zeilen, zwei Ziele | wie geplant, nicht in `GATE_CHECKS` |
| `harness/targets/zitat-vergleich.md` | neu, 268 Zeilen | wie geplant |
| Träger (5 Dateien) | geändert | wie geplant |
| `ADR-0159`, ADR-Index | unberührt | wie im Verdikt §5 |

Es gibt keine Datei im Diff ohne Plan-Zeile. Die Review-Reports und das Verdikt
sind Rollen-Artefakte. Kein Produkt-Code ist berührt; das entspricht der
Abgrenzung in §1.

---

## 8. Gate-Lauf

`make gates` ist ungefiltert nach dem Commit dieses Berichts gefahren worden.
Der Exit ist direkt gesichert. Das Ergebnis steht in der Rückmeldung an den
Auftraggeber. Dieser Bericht friert vor dem Lauf ein, und eine Nachtragung
wäre ein zweiter Commit ohne neuen Inhalt.

---

## 9. Abweichungen und Übergabe an den Planner

- **Keine DoD-Verletzung.** Die Liefer-Punkte 1 bis 3 sind bestätigt. Die
  Review-Pflicht ist erfüllt: das letzte Re-Review zeigt 0 HIGH und 0 MEDIUM.
- **Offen für die Closure, nicht vom Verifier zu schließen:**
  - (a) Closure-Notiz §7 mit Lerneintrag.
  - (b) Beobachtungs-Register.
  - (c) Ausgänge beider Risiken aus §6. Für „Zwei Träger“ trägt *entfallen*
    (Abschnitt 5). Für „Locale am Host“ ist der Fall L gemessen; ein Host ohne
    `C.utf8` ist *hergeleitet*.
  - (d) Drei Paarungen.
  - (e) LOW F-1 aus dem Re-Review zur Fixrunde 3 (escapter Backtick
    verschiebt die Parität, fail-open, ohne Fundstelle) als benannte Grenze
    in Vertrag §Grenzen oder mit Ausgang.
  - (f) Die gemeldeten fremden Träger (Abschnitt 6).
- **INFO:** V1 zeigt, dass die Bindung der Tag-Paar-Form am Meldungstext
  hängt, weil `baumda` den Exit ohnehin auf 2 hält (Abschnitt 4).
- **INFO:** Die Absolutzahlen des Gleichstands hängen an der Aufzählung.
  Belastbar ist die Menge der Abweichungen (Abschnitt 5).
