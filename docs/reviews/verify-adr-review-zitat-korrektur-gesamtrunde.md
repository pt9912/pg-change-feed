# Verifikationsbericht: ADR-Review-Zitat-Korrektur-Gesamtrunde — 2026-09-18

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen das
eigentliche Ziel der Runde: entsperrt sie `archive-welle` für die 25
korrigierten ADRs? Geprüft gegen den vorausgehenden Architect-Verdikt
zur ADR-Review-Zitat-Korrektur
(§8) und [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(Zitat-Korrektur-Grenze). **Nicht** Gegenstand: der Diff selbst (Reviewer-Aufgabe,
bereits erledigt — das Review zur ADR-Review-Zitat-Korrektur-Gesamtrunde,
0 HIGH) — und **nicht** realer Bedarf (Validator, nicht ausgelöst).

**Gegenstand:** Diff `df673da..HEAD` (9 Commits: `9543b13`, `5679f57`,
`9d485eb`, `7809a9c`, `c2bc868`, `442ce2e`, `96047eb`, `248b1fc`, `1fdb868`
— der letzte ist der bereits reviewte Review-Report-Commit selbst).

**Frischer Kontext.** Jede Zahl in diesem Bericht stammt aus einem hier
selbst gefahrenen Lauf. Der Reviewer-Report war Kontext (Prüfumfang,
bereits erledigte Negativbefunde), nicht übernommen — eigener `diff`, eigener
`archive-welle --vorschau`-Vergleich Vorher/Nachher, eigene ADR-Stichprobe
(vier ANDERE ADRs als der Reviewer), eigene `make gates`-Ausführung.

---

## 1. Eigene Messungen dieses Laufs

Jeder Lauf ungepiped, Exit-Code direkt aus einem eigenen, abgeschlossenen
Schritt gelesen (`AGENTS.md` §3.9).

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `git worktree add … df673da`, Binary `.harness/state/bin/ai-harness-init` hineinkopiert (gitignored, lokal gebaut, durch diesen reinen Doku-Diff unverändert) | — | „Vorher"-Stand für den Vorschau-Vergleich hergestellt |
| 2 | `ai-harness-init archive-welle --vorschau welle-d-check` auf `df673da` | — | Sperren: 2 (`[untergrenze]`, `[haenger]`); Hänger-Liste 217 Kanten, davon **20** mit `docs/plan/adr/*.md` als Quelle |
| 3 | dieselbe `--vorschau`-Zeile auf `HEAD` | — | Sperren: 2 (`[untergrenze]`, `[haenger]`); Hänger-Liste 208 Kanten, davon **0** mit `docs/plan/adr/*.md` als Quelle |
| 4 | `comm -23`/`comm -13` sortierter Hänger-Listen (vorher/nachher) | — | 20 Kanten verschwunden (alle `docs/plan/adr/*.md → docs/reviews/*.md`), 11 neu hinzugekommen (alle mit dem Architect-Verdikt zur ADR-Review-Zitat-Korrektur als Quelle, Ziel je eine `docs/reviews/*.md`-Datei) |
| 5 | `git diff c2bc868~1 c2bc868 -- docs/plan/adr/0053-*.md docs/plan/adr/0066-*.md docs/plan/adr/0085-*.md docs/plan/adr/0091-*.md` (eigene Stichprobe, andere 4 ADRs als der Reviewer-Report) | — | in allen vier Fällen ausschließlich Pfad/Link → Kennung ersetzt; §Entscheidung, §Konsequenzen, Zahlen, Findings, Datum/Autor unverändert |
| 6 | `grep -c PENDING_COMMIT docs/plan/adr/*.md` (repo-weit) | — | 0 Treffer — der Platzhalter ist überall durch den echten Hash `c2bc868` ersetzt (auch in den vier Stichproben-ADRs bestätigt) |
| 7 | `grep -rln "reviews/review-\|reviews/verify-\|reviews/architect-verdict-\|reviews/architect-review-" docs/plan/adr/*.md \| wc -l` (eigener Nachvollzug der Architect-Verdikt-Zahl) | — | 0 — kein einziges lebendes Adress-Zitat mehr im gesamten ADR-Bestand |
| 8 | `make gates` (eigener Lauf, ungepiped, Log in Datei, Exit separat gelesen) | **0** | alle sechs Gates grün; `docs-check` 906 Datei(en), 0 Befund(e); `coverage-gate` 83.40 % ≥ 80 %; `commit-traceability` OK (5 Commits `HEAD~5..HEAD`); `generated-sync` OK; `a-check` 0 Befunde; `baseline-verify` v6.9.0 OK |
| 9 | `make doc-commits RANGE=df673da..HEAD` | **0** | 906 Datei(en) geprüft, 0 Befund(e) (Modul `commits`) |
| 10 | `make doc-immutable RANGE=df673da..HEAD` | **0** | 906 Datei(en) geprüft, 0 Befund(e) (Modul `vcs`) |
| 11 | `grep -n "review" .d-check.yml` | — | `review`-Klasse (`paths: ["docs/reviews/*.md"]`, `token`), Regel `{from: adr, to: review, allow: false}`, `matrix.exempt-paths` trägt `docs/reviews/*.md` zusätzlich zu den zwei Grandfather-Einträgen — Struktur bestätigt wie im Reviewer-Report behauptet |

---

## 2. Check 1 — Architect-Verdikt §8 Punkt 7 gelesen

Der offene Rest ist explizit benannt: „Diese Fixrunde entsperrt
`archive-welle` für die 25 ADRs, aber nicht notwendig repo-weit — andere
Träger (`done/`-Records, `docs/reviews/**` selbst, Nicht-Markdown) können
denselben Review-Basisnamen noch erwähnen und bleiben außerhalb dieses
Verdikts." Diese Aussage ist **präzise zutreffend** — Lauf 4 findet genau
diese Klasse real vor (elf neue Kanten, deren Quelle eine
`docs/reviews/**`-Datei ist, nicht eine ADR).

## 3. Check 2 — Vorher/Nachher-Diff der Hänger-Liste

**Die im Architect-Verdikt namentlich genannten ADR→Review-Kanten sind
vollständig verschwunden.** Von den ursprünglich zwölf genannten ADRs
(`ADR-0062, -0069, -0070, -0075, -0083, -0084, -0086, -0087, -0089, -0091,
-0092, -0093`) trugen zehn davon real eine oder mehrere Hänger-Kanten in
`welle-d-check`s Vorschau (0062: 1, 0069: 1, 0070: 2, 0075: 2, 0083: 3,
0084: 1, 0086: 5, 0087: 1, 0089: 1, 0091: 1, 0092: 1, 0093: 1 = 20 Kanten
insgesamt) — **alle zwanzig sind nach der Korrektur weg** (Lauf 4, Spalte
„verschwunden"). Kein einziges der 25 korrigierten ADRs erscheint in der
Nachher-Liste noch als Quelle einer `[haenger]`-Kante.

**Die Sperre insgesamt bleibt trotzdem bestehen** — erwartungsgemäß, exakt
wie im Architect-Verdikt §8 Punkt 7 und im Reviewer-Negativbefund vermerkt.
Zwei unabhängige Ursachen:

1. **Vorbestehend, außerhalb dieser Runde:** Der Altbestand an
   `done/`-Records und `BEO-PGC/evidence/*.md`-Dateien, die
   Review-Basisnamen bare erwähnen (z. B.
   `docs/plan/planning/done/slice-028-*.md` zitiert das Review zu `slice-023`,
   `docs/plan/planning/observations/BEO-PGC/**/evidence/*.md → review-slice-*.md`)
   — unverändert vorher wie nachher, nie Gegenstand dieser Runde.
2. **Neu, selbst-erzeugt durch diese Runde:** Der Architect-Verdikt dieser
   Runde selbst (der Architect-Verdikt zur ADR-Review-Zitat-Korrektur)
   zitiert in seiner eigenen Bestands-Tabelle (§1) mehrere Review-Basisnamen
   als bare Text-Erwähnung (Reviews zu `slice-078`, `slice-081`,
   Verifikationsbericht zu `slice-096`, etc.) — das ist eine `docs/reviews/** → docs/reviews/**`-Kante,
   für die `.d-check.yml`s neue `review`-Klasse bewusst **keine** Matrix-Regel
   trägt (§7.3 des Verdikts: „review → review" ist als unproblematisch
   eingestuft, weil beide Dateien gemeinsam archiviert werden könnten) —
   aber der `Haenger`-Scan von `ai-harness-init` kennt diese Ausnahme nicht,
   er zählt jeden Basisnamen-Teilstring repo-weit. Dieser zweite Punkt ist
   exakt die Risiko-Klasse, die `slice-archive-altbestand-vollzug.md` §4
   („in-progress → open") bereits vorab benennt: „ein zwischenzeitlich
   entstandener zweiter `[haenger]`-Fund außerhalb der extern bereinigten
   Menge — dann Blocker, Carveout prüfen". Kein neuer, unerwarteter Fund —
   ein bereits antizipierter Fall, der real eingetreten ist.

**Zahlen im Überblick:**

| Messgröße | Vorher (`df673da`) | Nachher (`HEAD`) | Differenz |
|---|---|---|---|
| Hänger-Kanten gesamt | 217 | 208 | −9 |
| davon Quelle `docs/plan/adr/*.md` | 20 | 0 | −20 |
| davon Quelle: der Architect-Verdikt zur ADR-Review-Zitat-Korrektur | 0 | 11 | +11 |
| Sperren (`archive-welle --vorschau welle-d-check`) | 2 (`[untergrenze]`, `[haenger]`) | 2 (`[untergrenze]`, `[haenger]`) | unverändert |

## 4. Check 3 — `AGENTS.md` §3.12 (Herkunft von Aussagen)

- **Die 25er-Zahl** trägt ihren Ursprung im Architect-Verdikt als
  ausdrücklich benannter Messbefehl mit Ergebnis (§1: „liefert **25** Treffer
  (Zahl bestätigt)"); eigener Nachvollzug (Lauf 7) bestätigt sie erneut und
  zusätzlich, dass die Zahl nach der Korrektur auf 0 gefallen ist (kein
  Träger behauptet mehr „25" als aktuellen Zustand nach der Korrektur — der
  Reviewer-Report nennt „0 Treffer über den gesamten ADR-Bestand", korrekt
  als Ist-Zustand nach dem Fix formuliert, nicht als Wiederholung der
  25er-Zahl).
- **`make gates`-Zahlen** im Reviewer-Report (83.40 % Coverage, „905
  Datei(en), 0 Befund(e)", Exit 0) sind durchgängig als „real ausgeführt
  (Exit direkt geprüft, kein Pipe)" gekennzeichnet — gemessen, mit Methode
  genannt (§3.12 Instanz A erfüllt). Die Dateizahl ist inzwischen auf 906
  gestiegen (Lauf 8) — reine Bestandsverschiebung durch den seither
  hinzugekommenen Review-Report-Commit (`1fdb868`), keine Drift-Verletzung:
  der Reviewer maß zu seinem Zeitpunkt korrekt 905, dieser Bericht misst zu
  seinem eigenen Zeitpunkt korrekt 906 — beide Zahlen tragen ihren Lauf.
- **Die Hänger-Zahlen** dieses Berichts (217/208/20/11/9, Tabelle oben)
  tragen ihren Ursprung als eigene, hier gefahrene Läufe (Lauf 2–4),
  inklusive des Commits, gegen den je gemessen wurde (`df673da` bzw.
  `HEAD`).

Keine Verletzung von §3.12 gefunden.

## 5. Check 4 — Eigene Stichprobe gegen `ADR-0073`s Grenze (vier andere ADRs)

`ADR-0053`, `ADR-0066`, `ADR-0085`, `ADR-0091` (bewusst disjunkt von den
zehn im Reviewer-Report bereits gediffeten ADRs `0044, 0047, 0062, 0065,
0069, 0072, 0073, 0083, 0086, 0087, 0092, 0093`, um die Abdeckung über die
25 zu maximieren — zusammen mit dem Reviewer-Sample sind damit 14 der 25
individuell gegen ihre Vorversion gediffed).

In allen vier Fällen (Lauf 5): ausschließlich Pfad/Link durch Kennung
ersetzt (z. B. `ADR-0053`: „Architect-Verdikt" mit eigenem Link auf den
Verdikt zur Retention-Löschausführung → „der
vorausgehende Architect-Verdikt zur Retention-Löschausführung"; `ADR-0066`:
Link auf das Review zu `slice-069` → „Review zu `slice-069`"). §Entscheidung,
§Konsequenzen, §Verglichene Alternativen, §Status, Datum/Autor und jede
übernommene Zahl/Finding-Aussage bleiben wortgleich. Jede der vier ADRs
trägt zusätzlich die neue `§Geschichte`-Zeile „Zitat-Korrektur —
`docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`)" mit dem echten
Commit-Hash `c2bc868` (Lauf 6 bestätigt: kein `PENDING_COMMIT`-Rest im
gesamten Bestand). **Keine Grenzverletzung gefunden** — deckt sich mit dem
Reviewer-Befund, unabhängig erhoben.

## 6. Check 5 — `make gates` selbst reproduziert

Lauf 8: Exit `0`, ungepiped, direkt geprüft. Alle sechs Gates grün. Zusätzlich
`make doc-commits`/`make doc-immutable` über denselben Diff-Bereich
(`df673da..HEAD`, Läufe 9/10) — je 0 Befunde.

---

## 7. Verdikt

**Ziel teilweise erreicht — präzise im vom Architect-Verdikt selbst
gezogenen Rahmen.**

- **Erreicht:** Die 25 korrigierten ADRs sind als Quelle einer
  `[haenger]`-Kante vollständig verschwunden (20 von 20 real gemessenen
  Kanten weg, 0 neue ADR-Kanten). Die im Architect-Verdikt namentlich
  benannten zwölf ADRs (`0062, 0069, 0070, 0075, 0083, 0084, 0086, 0087,
  0089, 0091, 0092, 0093`) tragen keine einzige Hänger-Kante mehr. Die
  Zitat-Korrektur-Grenze aus `ADR-0073` ist in der eigenen Stichprobe (vier
  andere ADRs als der Reviewer) durchgängig eingehalten. `make gates` läuft
  eigenständig grün.
- **Nicht erreicht — aber auch nicht behauptet:** `archive-welle --vorschau
  welle-d-check` bleibt gesperrt (`Sperren: 2`, unverändert vor/nach). Der
  Grund ist zur Hälfte vorbestehend (Altbestand in `done/`/`BEO-PGC`-Evidenz,
  außerhalb dieser Runde) und zur Hälfte ein **Selbst-Effekt** dieser Runde:
  ihr eigener Architect-Verdikt-Report zitiert Review-Basisnamen bare in
  seiner Bestandstabelle und erzeugt dadurch elf neue Hänger-Kanten. Dieser
  Effekt ist im vorausgehenden Verdikt (§2 Punkt 4, §8 Punkt 7) und im
  Reviewer-Negativbefund bereits als offener Rest benannt, und im
  nachfolgenden Slice `slice-archive-altbestand-vollzug.md` §4 bereits als
  möglicher Rückführungs-Trigger vorgesehen — kein unerwarteter Fund, kein
  neuer Blocker-Typ.

**Netto:** `welle-archive-altbestand`s Slices sind ihrer eigenen Ausführung
einen Schritt näher — die spezifische Sperren-Klasse, die diese Runde
beheben sollte, ist real behoben (20 → 0 ADR-Kanten) —, aber
`slice-archive-altbestand-vollzug`s LP1 (erneuter `--vorschau`-Lauf zeigt
„Sperren: keine") ist **noch nicht** erreicht: Die Gesamtzahl der
`[haenger]`-Kanten sank von 217 auf 208, der `[haenger]`-Befund selbst
bleibt aktiv. Ein Folge-Vorgang, der die verbleibenden 208 Kanten (Altbestand
plus die elf selbst-erzeugten) adressiert, bleibt vor dem realen,
schreibenden `archive-welle`-Lauf nötig.

**Nicht Gegenstand dieser Verifikation:** Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst); der Diff selbst (Reviewer-Aufgabe,
bereits erledigt, 0 HIGH).
