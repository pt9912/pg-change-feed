# Review-Report: Zitat-Korrektur Beobachtungs-Register, Batch 2 (Commit `af9e7af`) — 2026-09-18

**Review-Art:** Code (Doku-Diff) — geprüft gegen `ADR-0073` (Zitat-Korrektur-Klasse),
`ADR-0097` (`observation`-Matrixklasse) und `AGENTS.md` §3.5.

**Gegenstand:** Commit `af9e7af` — 41 Dateien unter
`docs/plan/planning/observations/BEO-PGC/**`, zweiter Batch nach Commit
`faccc97` (42 Dateien, führte zugleich `ADR-0097` ein).

**Skill:** `.harness/skills/reviewer.md` @ Commit `7e598bb` (Stand zu Review-Beginn).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `ADR-0073` (Zitat-Korrektur an immutablen/historischen Dokumenten)
- `ADR-0097` (`observation`-Matrixklasse, Folgepflicht aus Batch 1)
- `ADR-0094`/`ADR-0095` (Vorläufer-Klasse `review`-Matrix)
- `AGENTS.md` §3.5 (ADR-Immutabilität, Zitat-Korrektur-Ausnahme)
- Baseline-Regelwerk `modul-06-roadmap.md` §Beobachtungs-Register
  (`evidence/<vorgangs-id>.md` unveränderlich ab Merge)
- `.harness/skills/reviewer.md` §Output-Schema, §Klassifikation

---

## Vorgehen

1. `git show af9e7af` vollständig gelesen (nicht nur Stichprobe — 41 Hunks,
   überschaubarer Umfang) statt einer Teilmenge von 8–10 Dateien: jeder Hunk
   ersetzt einen Datei-Basisnamen-Zitat (`docs/reviews/review-slice-NNN.md`,
   `docs/reviews/verify-slice-NNN.md`, `docs/reviews/architect-verdict-*.md`)
   durch eine Prosa-Adressierung („Review zu `slice-NNN`",
   „Verifikationsbericht zu `slice-NNN`", „der Architect-Verdikt zu …").
   Referent und Aussage sind in jedem geprüften Hunk unverändert — auch dort,
   wo die Umformulierung über eine reine 1:1-Ersetzung hinausgeht (z. B.
   `zitat-nennt-die-falsche-stelle/evidence/welle-d-check-verkoerperung.md`:
   „Fund (Review-Finding, HIGH, `docs/reviews/review-welle-d-check-…`)" →
   „Fund (Review-Finding, HIGH, im Delta-Review dieser
   Verkörperungs-Fixrunde)" — derselbe Referent, dieselbe Aussage, nur
   syntaktisch umgehängt).
2. Eigene Gegenprobe: `grep -l` mit
   `docs/reviews/|review-slice-|verify-slice-|architect-verdict-|architect-review-`
   über exakt die 41 von `af9e7af` geänderten Dateien (per
   `git show --name-only`) — **0 Treffer** (`grep`-Exit `1`).
3. `.d-check.yml`-Diff von `af9e7af` und zwischen `faccc97`..`af9e7af`
   geprüft: **keine** Änderung an `.d-check.yml` in diesem Commit — insbesondere
   keine neue `exempt-paths`-Ausnahme der Art, die Batch 1 in `matrix`
   (Zeile `exempt-paths: [..., "docs/reviews/*.md"]`) und `versions`
   (Zeile `exempt-paths: ["harness/conventions/done/**", "docs/reviews/**"]`)
   trägt. Diese zwei bekannten Ausnahmen bleiben unverändert und sind von
   diesem Batch nicht betroffen.
4. `make gates` selbst reproduziert (ungepiped, Exit-Code separat geprüft,
   `AGENTS.md` §3.9): **Exit 0**. `d-check` meldet „911 Datei(en) geprüft,
   0 Befund(e)" (Haupt- und Traceability-Lauf), `a-check` „gesamt: 0
   Befund(e)", `generated-sync: OK`.

## Beobachtung außerhalb des unmittelbaren Diff-Skopus (INFO)

`ADR-0097` (eingeführt in Batch 1, Commit `faccc97`) verlangt unter
§Konsequenzen „Folgepflicht (Implementer-Zug)" drei Schritte: (1) betroffene
Beobachtungs-Register-Dateien umschreiben, (2) `.d-check.yml` um die
`observation`-Matrixklasse/-Regel ergänzen, (3) `make gates` grün. Zum
Zeitpunkt dieses Reviews trägt `.d-check.yml` **noch keine**
`observation`-Klasse/-Regel (Stand geprüft: letzte Änderung an der Datei ist
Commit `248b1fc`, vor `faccc97`/`af9e7af`). Schritt (2) der ADR-0097-Folgepflicht
ist damit noch offen — das ist **kein** HIGH gegen diesen Commit (`af9e7af`
selbst berührt `.d-check.yml` nicht und behauptet auch nicht, die Folgepflicht
abzuschließen), aber ein Hinweis für die laufende Serie: Batch 1 und Batch 2
zusammen liefern Schritt (1) der Folgepflicht; Schritt (2) fehlt noch in
diesem Arbeitsbaum-Stand.

## Findings

Keine HIGH/MEDIUM/LOW-Findings gegen Commit `af9e7af`.

### F-1 — `ADR-0097`-Folgepflicht Schritt 2 noch offen

- `kategorie`: INFO
- `quelle`: `ADR-0097` §Konsequenzen
- `pfad`: `.d-check.yml` (unverändert seit Commit `248b1fc`)
- `befund`: Die `observation`-Matrixklasse/-Regel aus `ADR-0097` ist zum
  Zeitpunkt dieses Reviews noch nicht in `.d-check.yml` verdrahtet; Commit
  `af9e7af` behauptet das auch nicht, betrifft nur Schritt (1) der
  Folgepflicht (Register-Dateien umschreiben).
- `verifizierbar`: ja — `git log -1 --format=%H -- .d-check.yml` gegen
  `grep observation .d-check.yml`.
- `klasse`: Folgepflicht-Schritt offen (kein Wiederholungsmuster, erstes
  Auftreten in dieser Serie)

**geprüft, ohne Befund:** `docs/plan/planning/observations/BEO-PGC/**` (die
41 von `af9e7af` geänderten Dateien, vollständig gelesen) · `.d-check.yml`
(Diff-Skopus dieses Commits: leer) · `make gates`-Reproduktion.

## Verdikt

0 HIGH, 0 MEDIUM, 0 LOW. Ein INFO-Hinweis außerhalb des Commit-Skopus (siehe
oben), keine Fixrunde am Implementer dieses Batches nötig. Kein Slice-Plan
ist an diesen Commit gebunden (Beobachtungs-Register-Pflege ist wellenlos,
Modul 6 „Träger im Repo ohne Wellen") — daher kein DoD-Checkbox-Nachzug
gemäß Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde anwendbar.
