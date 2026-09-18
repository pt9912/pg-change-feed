# Review-Report: ADR-Review-Zitat-Korrektur-Gesamtrunde — 2026-09-18

**Review-Art:** Code/Design — geprüft gegen Plan (die zwei vorausgehenden
Architect-Verdikte), `AGENTS.md` §3.5/§3.6/§3.7/§3.11/§3.12/§3.13 und
`harness/conventions.md` (MR-000). Kein Slice-Plan zugrunde — Auslöser war
eine direkte Nutzer-Direktive plus zwei Architect-Verdikte; DoD-Konformität
ist deshalb nicht Gegenstand (kein DoD-Träger vorhanden), das bleibt der
Verifier-Rolle für die Slices vorbehalten, die diesen Zug später
konsumieren (`slice-archive-altbestand-vollzug`).

**Gegenstand:** Diff `df673da..HEAD`, 8 Commits (`9543b13`, `5679f57`,
`9d485eb`, `7809a9c`, `c2bc868`, `442ce2e`, `96047eb`, `248b1fc`).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (dieses Repo, kein Tag-Bezug
— repo-lokale Skill-Datei, keine Baseline-Ressource).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- der vorausgehende Architect-Verdikt zur ADR-Review-Zitat-Korrektur
  (2026-09-18) — 25 `Accepted`-ADRs, Zielform, Verdikt „Zitat-Korrektur
  zulässig für alle 25"
- [ADR-0094](../plan/adr/0094-review-matrixklasse-kennung-statt-adresse.md)
  (Review-Matrixklasse, ergänzt [ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md))
- der Architect-Verdikt zur `matrix.status`-Ausnahme für die `review`-Klasse
  (2026-09-18, Folgefrage zu [ADR-0094](../plan/adr/0094-review-matrixklasse-kennung-statt-adresse.md))
- [ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md)
  (Review-Klasse — Status-Ausnahme, ergänzt [ADR-0094](../plan/adr/0094-review-matrixklasse-kennung-statt-adresse.md))
- `AGENTS.md` §3.5 (ADR-Immutabilität/Zitat-Korrektur), §3.6
  (Gate-Lockerung braucht ADR), §3.7 (Kommentar-Klassen), §3.11
  (Hostpaths-Verbot), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug)
- `harness/conventions.md` MR-000 (ID-Schema)

---

## Findings

### F-1 — `ADR-0080` trägt eine bare Review-Kennung außerhalb der 25-Liste, ohne Adress-Risiko

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md:71`
- `befund`: Der Satz „… und die Verengung fallen gelassen (Review
  `review-slice-081` <!-- d-check:status-provenance --> F-2, Plan §1)"
  nennt eine Basisform ohne `docs/reviews/`-Präfix und ohne `.md`-Endung —
  weder vom `matrix`-Token `docs/reviews/[\w.-]+\.md` (`ADR-0094`) noch vom
  Hänger-Scan (`filepath.Base`-Teilstring-Vergleich, der die volle
  `.md`-Endung braucht) erfasst, und außerhalb der ursprünglich
  gezählten 25 ADRs (diese ADR ist vom 2026-09-18-Verdikt nicht
  betrachtet worden, weil ihr Basisname kein `reviews/`-Präfix trägt). Kein
  Hänger-Risiko, keine Regelverletzung — nur nicht ganz die neu etablierte
  Kennung-Form („Review zu `slice-081`, Finding F-2").
- `verifizierbar`: ja — `grep -n "review-slice-081" docs/plan/adr/0080-*.md`
- `klasse`: Kennung-Form uneinheitlich außerhalb des geprüften Bestands

## Negativbefunde

- geprüft, ohne Befund: Vollständigkeit der 25 ADR-Korrekturen — eigener
  `grep -rln "reviews/review-\|reviews/verify-\|reviews/architect-verdict-\|reviews/architect-review-" docs/plan/adr/*.md`
  liefert 0 Treffer über den **gesamten** ADR-Bestand (nicht nur die 25
  gelisteten); keine 26. ADR mit demselben Fehler entstanden, `ADR-0094`
  und `ADR-0095` selbst eingeschlossen (beide zitieren den vorausgehenden
  Architect-Verdikt ausschließlich über Kennung/Thema-Prosa, nie über
  seinen Datei-Basisnamen)
- geprüft, ohne Befund: Zitat-Korrektur-Grenze (`AGENTS.md` §3.5) — 10 der 25
  ADRs im Detail gegen ihre Vorversion (`c2bc868~1` vs. `c2bc868`) gediffed
  (`0044`, `0047`, `0062`, `0065`, `0066`, `0069`, `0072`, `0073`, `0083`,
  `0086`, `0087`, `0091`, `0092`, `0093`); in keinem Fall wurde
  §Entscheidung, §Konsequenzen (der Aussage nach), §Verglichene
  Alternativen, §Status, eine übernommene Zahl oder ein Finding-Inhalt
  angefasst — ausschließlich Pfad/Link durch Kennung ersetzt
- geprüft, ohne Befund: die zwei im Verdikt benannten Fußnoten — `ADR-0062`
  hält die Unterscheidung zwischen dem „ersten Architect-Verdikt zur
  Commit-Traceability (kein Vorab-Hook)" und „der Architect-Verdikt dieses
  Zugs (Gegenprüfung)" korrekt als zwei verschiedene Kennungen fest (keine
  Glättung zu einer einzigen); `ADR-0072`s komprimierte Notation
  „`review-slice-036/039/049/073.md`" wurde exakt wie im Verdikt
  vorgegeben durch „Reviews zu `slice-036`, `-039`, `-049`, `-073`" ersetzt
- geprüft, ohne Befund: §Geschichte-Belege — alle 25 ADRs tragen nach
  `442ce2e` eine neue Geschichte-Zeile „Zitat-Korrektur — …" mit dem echten
  Commit-Hash `c2bc868` als Verweis (kein `PENDING_COMMIT`-Platzhalter mehr
  im Bestand, `grep -c PENDING_COMMIT docs/plan/adr/*.md` → 0); der Hash
  stimmt mit dem tatsächlichen Korrektur-Commit überein
- geprüft, ohne Befund: §3.13-Träger-Nachzug bei neu exponierten
  `slice-NNN`-Token — Diff `c2bc868~1..c2bc868` skriptgestützt ausgewertet:
  74 hinzugefügte vs. 61 entfernte `<!-- d-check:status-provenance -->`
  -Marker-Zeilen; jede hinzugefügte Marker-Zeile trägt ein `slice-\d{3}`
  -Token (0 verdächtige/dekorative Marker ohne Token), 16 konkrete
  Vorher/Nachher-Paare händisch verifiziert (u. a. `ADR-0044`, `ADR-0047`,
  `ADR-0069`, `ADR-0075`) — durchgängig Marker dort neu gesetzt, wo das
  Delinken einen zuvor in `[]()` verpackten `slice-NNN`-Basisnamen als
  bare Inline-Code-Token neu exponiert hat, kein Fall eines verlorenen,
  aber weiterhin nötigen Markers gefunden
- geprüft, ohne Befund: [ADR-0094](../plan/adr/0094-review-matrixklasse-kennung-statt-adresse.md)/[ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md) selbst — beide zitieren den
  vorausgehenden Architect-Verdikt konsequent über Kennung/Thema-Prosa
  („der vorausgehende Architect-Verdikt zur ADR-Review-Zitat-Korrektur
  (2026-09-18)"; „Architect-Verdikt zur `matrix.status`-Ausnahme für die
  `review`-Klasse (2026-09-18)"), nie über den Datei-Basisnamen; `ADR-0094`
  wurde nach seiner Erstellung (`7809a9c`) durch keinen der Folge-Commits
  mehr angefasst (`git log --oneline -- docs/plan/adr/0094-*.md` → ein
  einziger Commit) — `ADR-0095`s Behauptung „`ADR-0094`s Klausel bleibt
  real unverändert" ist damit verifizierbar wahr, kein `Supersedes` nötig
- geprüft, ohne Befund: `.d-check.yml`-Struktur — `matrix.exempt-paths` ist
  ein direktes Geschwister von `matrix.status` (nicht dessen
  Unter-Schlüssel), Einrückung geprüft; die zwei bestehenden
  Grandfather-Einträge (`0039-*.md`, `0041-*.md`) unverändert, nur
  `docs/reviews/*.md` als dritter Eintrag ergänzt; `matrix.classes`/`rules`
  tragen exakt den in `ADR-0094` abgedruckten Diff
- geprüft, ohne Befund: `make docs-check` isoliert real ausgeführt (Exit
  direkt geprüft, kein Pipe) → „905 Datei(en) geprüft, 0 Befund(e)" — 0
  `matrix-forbidden` und 0 `matrix-inactive`
- geprüft, ohne Befund: `make gates` vollständig real ausgeführt (Exit
  direkt geprüft) → Exit 0; baseline-verify, coverage-gate (83.40 % ≥
  80 %), commit-traceability, generated-sync, a-check, docs-check
  allesamt grün
- geprüft, ohne Befund: der bewusst offene Rest (Basisnamen-Erwähnungen
  außerhalb der 25 ADRs, z. B. in `done/`-Records oder `docs/reviews/**`
  selbst) verschwindet nicht spurlos — `docs/plan/planning/open/slice-archive-altbestand-vollzug.md`
  benennt ihn explizit als Blocker-Kandidat („ein zwischenzeitlich
  entstandener zweiter `[haenger]`-Fund außerhalb der extern bereinigten
  Menge — dann Blocker, Carveout prüfen") und `welle-archive-altbestand.md`
  §5 trägt die noch nachzutragende Kennung dieses Vorgangs als offene
  Abhängigkeit
- geprüft, ohne Befund: Commit-Traceability aller 8 Commits im
  Diff-Bereich — `tools/harness/commit-traceability.sh "df673da..HEAD"` →
  „OK — 8 Commit(s) …, Betreffs ohne Struktur-ID"; jeder Betreff trägt
  eine `ADR-*`-Kennung, keiner eine `LH-*`/`SPEC-*`/`ARC-*`-Kennung im
  Betreff
- geprüft, ohne Befund: Hostpaths-Verbot (`AGENTS.md` §3.11) in den beiden
  neuen Verdikt-Dokumenten und den zwei neuen Slice-Plänen — kein
  `<Host-Wurzel>`-Präfix mehr vorhanden (Commit `9d485eb` hat die
  ursprünglich im vorausgehenden Architect-Verdikt zur
  ADR-Review-Zitat-Korrektur vorhandenen Host-Pfade bereits selbst auf
  Hausform korrigiert)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Kennung-Form uneinheitlich außerhalb des
geprüften Bestands

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 1 INFO ohne erwartete
Aktion. Die gesamte Runde (Bestandsaufnahme → Zitat-Korrektur der 25 ADRs →
Geschichte-Beleg-Nachzug → Mechanisierung via `ADR-0094` →
Nebenwirkungs-Korrektur via `ADR-0095` → `.d-check.yml`-Aktivierung) hält
die Zitat-Korrektur-Grenze aus `AGENTS.md` §3.5 durchgängig ein, die zwei
Fußnoten wurden korrekt behandelt, die §Geschichte-Belege tragen echte
Commit-Hashes, der §3.13-Träger-Nachzug bei den neu exponierten
`slice-NNN`-Token ist sauber (keine fehlenden, keine dekorativen
Marker), `ADR-0094`/`ADR-0095` verstoßen nicht gegen ihre eigene neue
Regel, die `.d-check.yml`-Struktur ist korrekt, und `make gates` läuft
real grün (Exit 0, direkt geprüft, `AGENTS.md` §3.9).

**Übergabe:** Kein Fixrunden-Pfeil an einen Implementer-Lauf — das eine
INFO-Finding (F-1) braucht keine Aktion (Skill: „Hinweis ohne erwartete
Aktion"). Da kein Slice-Plan dieser Runde zugrunde liegt, entfällt der
DoD-Checkbox-Nachzug des Reviewer-Skills (kein DoD-Träger vorhanden, den
man nachziehen könnte) — die Konsumenten dieser Runde
(`slice-archive-altbestand-vollzug`) tragen ihre eigene DoD/Review-Zeile
bei ihrer eigenen Closure. Dieser Report selbst ist ein Lauf-Beleg (Audit:
dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und wird über
Läufe hinweg nicht wieder gelesen.
