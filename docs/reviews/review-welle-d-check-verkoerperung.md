# Review-Report: Verkörperungs-Fixrunde `welle-d-check` (Commit `4da1ad8`) — 2026-09-17

**Review-Art:** Code — geprüft gegen **Plan, Entscheidungen und Hard Rules**
(Maintainability). DoD-/Spec-Konformität ist **nicht** Gegenstand dieses
Laufs (Verifier-Frage, Modul 11).

**Gegenstand:** `4da1ad8` (Parent `7e598bb`) — **2** Dateien,
+49/−16: `.harness/skills/reviewer.md`, `AGENTS.md`.

**Skill:** `.harness/skills/reviewer.md` @ `4da1ad8` (der zu prüfende Commit
selbst — die geschärfte Fassung des HIGH-Punkts „Beleg trägt seinen Satz
nicht" wird hier reflexiv auf den eigenen Commit angewendet).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17.

**Eingangs-Kontext:**

- der Architect-Verdikt zum Lese-Schritt der `d-check`-Welle (der
  beauftragende Architect-Zug, §1.2/§2.3/§3 „Folgearbeit")
- `AGENTS.md` §3.13 (Vorzustand und Nachzustand)
- `.harness/skills/reviewer.md` HIGH-Punkt „Beleg trägt seinen Satz nicht"
  (Vorzustand und Nachzustand)
- das Review zu `slice-096` F-2 (das zitierte Original)
- `docs/plan/planning/observations/BEO-PGC/regel-weiter-als-ihr-sensor/evidence/slice-096.md`
  und `docs/plan/planning/done/slice-096-konfigurationsdatei-nachzug.md` §7
  (die Herkunft der Formulierung „vierter Anker")
- `harness/sensors/coverage-gate.md` §Grenze Punkt 4
- `AGENTS.md` §3.9 (Exit-Code-Disziplin für den eigenen `make gates`-Lauf)

---

## Findings

### F-1 — Der neue Grenz-Absatz in `AGENTS.md` §3.13 zitiert das Review zu `slice-096` F-2 mit einer Zählung und einem Begriff, die das Original nicht trägt

- `kategorie`: **HIGH**
- `quelle`: Skill HIGH „Beleg trägt seinen Satz nicht" (Nachbarform „Verweis
  auf eine Stelle eines anderen Dokuments", genau die Erweiterung, die dieser
  Commit selbst einführt) · `BEO-PGC/zitat-nennt-die-falsche-stelle`
- `pfad`: `AGENTS.md:443-446` gegen Zeilen 49-61 des Reviews zu `slice-096` (F-2)
- `befund`: Der neue Satz lautet: „real bereits so geschlossen:
  das Review zu `slice-096` F-2 fand über den vierten Anker **einen**
  Zeilen-Lokator, den der Implementer-Suchlauf übersehen hatte." Das Original
  (F-2, Titel: „**Drei** Anker in `ADR-0088` lösen … nicht mehr auf") berichtet
  explizit **zwei** Zeilen-Lokatoren („Nicht gemeldet: **die zwei**
  Zeilen-Lokatoren derselben ADR"), nicht einen. Der Begriff „vierter Anker"
  kommt in F-2 selbst **nicht** vor — er stammt aus
  `BEO-PGC/regel-weiter-als-ihr-sensor/evidence/slice-096.md:17` („Gefunden
  hat sie der Review, **den vierten Anker der Architect**") und wird dort dem
  **Architect**-Zug zugeschrieben, nicht dem Review. Die Closure-Notiz von
  `slice-096` (§7, Zeile 255-257) liest sich wiederum anders („der Review zwei
  Lokatoren **und einen vierten Anker**, und der Architect-Zug fand die zwei
  falschen Zahlen in den ADRs selbst") — die drei Quelldokumente widersprechen
  sich bereits untereinander darüber, wer den „vierten Anker" gefunden hat.
  Der neue `AGENTS.md`-Satz übernimmt diese Unschärfe, ohne sie aufzulösen,
  und schreibt sie zusätzlich in einer Form fest, die dem zitierten Original
  (F-2) eine Zählung („einen") zuschreibt, die es nicht enthält (F-2 zählt
  „zwei", nicht „einen"). Die Kern-Aussage des Satzes — dass der Reviewer eine
  Lücke schloss, die der Implementer-Suchlauf offen ließ — ist in der Sache
  richtig; die **Stütze**, die für sie zitiert wird, trägt die exakte Form der
  Aussage nicht.
- `verifizierbar`: ja — das Review zu `slice-096` F-2 im Original
  aufschlagen (nicht in einer Zusammenfassung), Zählung „zwei" vs. „einen"
  liegt offen; `evidence/slice-096.md:17` und die Closure-Notiz `§7`
  gegeneinanderhalten zeigt den Attributions-Widerspruch.
- `klasse`: „Zitat nennt die falsche Stelle" (`BEO-PGC/zitat-nennt-die-falsche-stelle`
  — mutmaßlich viertes Vorkommen dieser bereits dreifach belegten Klasse, in
  eben dem Commit, der die Klasse als HIGH-Nachbarform verkörpert)

## Was der Diff richtig macht (nachgemessen, nicht übernommen)

- **`AGENTS.md` §3.13 — Grenz-Absatz.** Trägt alle drei geforderten Elemente:
  die Symbolname-vs-Zahlen-Lücke ist explizit benannt („trifft zuverlässig
  Symbolnamen … trifft nicht zuverlässig Zahlen"), der Reviewer ist als
  Schließer benannt, und der Absatz verlangt ausdrücklich **keine**
  Erweiterung der Suchform selbst („Diese Regel verlangt deshalb keine
  Erweiterung der Suchform selbst"). Platzierung stimmt mit dem Verdikt
  überein (unmittelbar nach „Warum diese Regel einen Träger braucht und
  keinen Vorsatz").
- **`AGENTS.md` §3.13 — „Wer sie liest".** Trägt jetzt die Pflicht zu einem
  committeten Feld im Slice-Plan vor der Closure-Notiz („der Implementer hält
  Gefundenes und Nichtgefundenes zusätzlich in einem committeten Feld des
  Slice-Plans selbst fest, bevor er die Closure-Notiz (§7) schreibt"),
  inhaltlich deckungsgleich mit Verdikt §1.2 Punkt 2 (Wortlaut-Treue ohne
  Copy-Paste-Zwang zulässig — Verdikt nennt eine DoD-Zeile nur als
  „naheliegend", nicht als Pflichtform; keine Template-Änderung verlangt,
  keine vorgenommen).
- **`AGENTS.md` §3.13 — Herkunfts-Anker.** Trägt jetzt `seit welle-20 ·
  geschärft seit welle-d-check (BEO-PGC/regel-weiter-als-ihr-sensor, 3×,
  slice-078/slice-079/slice-096; Architect-Verdikt)` — Zählung und
  Slice-Liste stimmen mit der Tabelle in §1.1 des Verdikts überein.
- **`.harness/skills/reviewer.md` — Beleg-Formen-Aufzählung.** Trägt jetzt
  „einen Verweis auf eine Stelle eines anderen Dokuments (eine
  Abschnittsnummer, eine ADR-Festlegung, eine Slice-/Welle-Kennung)" —
  wortgleich mit Verdikt §2.3 Punkt 1.
- **`.harness/skills/reviewer.md` — Imperativ-Satz.** Trägt den geforderten
  zweiten Halbsatz „— oder schlägt die genannte Stelle im Original auf, nicht
  in einer Zusammenfassung (die Kopfzeile eines Berichts, die eigene
  Erinnerung)" — wortgleich mit Verdikt §2.3 Punkt 2.
- **`.harness/skills/reviewer.md` — Herkunfts-Zeile.** Trägt jetzt
  `BEO-PGC/zitat-nennt-die-falsche-stelle` (3×,
  `slice-090`/`slice-102`/`slice-d-check-tracked-modul`) mit Anker `· seit
  welle-d-check`; die ergänzte Klassifikations-Drift-Notiz
  („uneinheitlich als MEDIUM/INFO/HIGH klassifiziert") stimmt mit der
  Kategorie-Spalte der Tabelle in Verdikt §2.1 überein (MEDIUM, INFO, HIGH in
  genau dieser Reihenfolge).
- **Out-of-Scope-Disziplin.** `git show 4da1ad8 --stat` zeigt genau zwei
  Dateien; keine Änderung an `harness/sensors/coverage-gate.md`, `AGENTS.md`
  §3.11, `.d-check.yml`, einer ADR oder einem Beobachtungs-Register-Eintrag
  (`state.md`/`observation.md`) — deckungsgleich mit Verdikt §4 „Was nicht
  getan wurde".
- **ADR-Immutabilität / §3.6.** Keine `Accepted`-ADR berührt, keine neue ADR,
  keine Schwellen-Senkung — Verdikt §3 bestätigt dies als reine Textpflege.
- **`make gates`.** Eigener Lauf, Exit-Code direkt und ungepiped geprüft:
  `EXIT_CODE=0` (Coverage 83.40 % ≥ 80 %, `d-check` 0 Befunde in beiden
  Modulen, `commit-traceability` OK, `generated-sync` OK, `a-check` 0
  Befunde).

## Negativbefunde

- geprüft, ohne Befund: `harness/sensors/coverage-gate.md` (unverändert, wie
  gefordert)
- geprüft, ohne Befund: `AGENTS.md` §3.11 (unverändert, wie gefordert)
- geprüft, ohne Befund: `docs/plan/adr/` (kein neuer/geänderter ADR-Eintrag)
- geprüft, ohne Befund: `.d-check.yml` (unverändert)
- geprüft, ohne Befund: `docs/plan/planning/observations/BEO-PGC/regel-weiter-als-ihr-sensor/`
  und `.../zitat-nennt-die-falsche-stelle/` (`state.md`/`observation.md`
  unverändert — Ausgangs-Zuweisung bleibt korrekt beim nachfolgenden
  Planner-Zug, wie Verdikt §4 vorsieht)
- geprüft, ohne Befund: `docs/plan/planning/*.template.md` (keine
  DoD-Zeile ergänzt — Verdikt nennt das nur als „naheliegend", nicht als
  Pflicht dieser Fixrunde; korrekt unterlassen)
- geprüft, ohne Befund: Commit-Message-Form (`ADR-0045` im Betreff, keine
  `SPEC-*`/`ARC-*`-Kennung — `make commit-traceability` bestätigt)
- geprüft, ohne Befund: übrige Prosa der beiden geänderten Absätze (Grammatik,
  Platzierung, Abgrenzung zu §3.12/§3.7) — inhaltlich vollständig und an der
  richtigen Stelle eingefügt

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **1** |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Zitat nennt die falsche Stelle"

## Verdikt

**Merge-blockierend: ja** (ein HIGH-Finding). Die Verkörperung trifft beide
vom Architect-Verdikt beauftragten Textstellen inhaltlich vollständig und
wortlauttreu im Sinne von „Wortlaut-Treue ohne Copy-Paste-Zwang" — die
Out-of-Scope-Disziplin ist eingehalten (keine Änderung an
`coverage-gate.md`, `AGENTS.md` §3.11 oder einer ADR), `make gates` läuft
grün. Der eine gefundene HIGH-Defekt liegt ausgerechnet in dem Satz, der die
Beleg-Genauigkeit des neuen HIGH-Punkts illustrieren soll: die zitierte
Stelle (Review zu `slice-096` F-2) trägt die genannte Zählung („einen"
Zeilen-Lokator) nicht — sie berichtet zwei — und der verwendete Begriff
„vierter Anker" wird in den vorhandenen Quelldokumenten uneinheitlich
zugeschrieben (Review vs. Architect). Das ist inhaltlich keine Katastrophe
(die Kern-Aussage des Satzes bleibt wahr), aber exakt die Fehlerklasse, die
dieser Commit selbst als HIGH einführt — angewendet auf sich selbst.

**Übergabe:** Rückkante an den Implementer für eine punktuelle
Text-Korrektur des einen Satzes in `AGENTS.md:443-446` (Zählung korrigieren,
„vierter Anker"-Zuschreibung entweder auflösen oder fallen lassen). Kein
Architect-Pfad nötig — der Defekt liegt vollständig im Diff selbst, nicht in
einem Träger außerhalb davon, und es besteht kein Rollen-Widerspruch. Die
Finding-Klasse geht zusätzlich in den Steering-Loop-Zähler von
`BEO-PGC/zitat-nennt-die-falsche-stelle` (jetzt potenziell 4×, unabhängig zu
bestätigen bei der nächsten Wellen-Closure).

**Nicht gefahren:** `make a-check` separat, `make baseline-verify`, `make
image`, `make test-store`/`-replication`/`-notify`/`-integration` — keine
Architektur-Kante, keine Baseline, kein Build-Kontext, keine Naht und kein
Container-Vertrag berührt (reiner Prosa-Diff in zwei Markdown-Dateien);
`make gates` deckt `a-check` bereits mit ab.
