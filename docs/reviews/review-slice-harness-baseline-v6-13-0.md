# Review-Report: slice-harness-baseline-v6-13-0 — 2026-09-29

**Review-Art:** Code — geprüft wird der Diff `7103ad59..caf96849` (8 Commits)
wgegen Slice-Plan §1–§6, Drift-Audit, `AGENTS.md` Hard Rules
(§3.1, §3.5, §3.7, §3.12, §3.13), `harness/conventions.md` (MR-000).

**Gegenstand:** `slice-harness-baseline-v6-13-0`, Range `7103ad59..caf96849`

**Skill:** `.harness/skills/reviewer.md` @ `caf96849` ·
**Modell:** Claude Code (Reviewer-Rolle, Modul 10) · **Datum:** 2026-09-29

**Eingangs-Kontext:**

- Slice-Plan `slice-harness-baseline-v6-13-0` (§1 Ziel/Abgrenzung, §2 DoD,
  §3 Tabelle, §4 Trigger/Rückführung, §6 Risiken, Suchlauf-Feld)
- Drift-Audit `docs/reviews/audit-baseline-v6-13-0-drift.md` (`ad530688`)
- [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md) (Bezug),
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  (Zitat-Korrektur, Record-Grenze),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
  (Herkunft von Zahlen)
- `AGENTS.md` §3 Hard Rules · `harness/conventions.md` (MR-000,
  §Baseline) · Baseline `v6.13.0` · `regelwerk/modul-10-review-harness.md`

---

## Eigene Messungen des Laufs (Nachmessen, nicht Übernahme)

- `make baseline-verify`: **Exit 0**, „54 Dateien (Integritaet +
  Vollstaendigkeit, netzlos)".
- Asset-Verifikation (Docker-gekapselt, dasselbe gepinnte Toolchain-Image
  `golang:1.27-alpine@sha256:cf6fca66…3125`): `lab-regelwerk.zip` des
  Releases `v6.13.0` geladen, sha256 **gemessen**
  `b5151e77807e2affebb25cfab9be24b88cf43afc42c075db0b982a2ff1052b96` —
  byte-gleich zur Angabe in Audit §1 und Commit-Message `d443ee39`. Entpackt:
  Inhalt **byte-identisch** zum committeten Baum
  `.harness/baseline/v6.13.0/`, einzige Differenz `SHA256SUMS` (nicht im
  Asset — bestätigt Audit §1 „Bootstrap-eigenes Erzeugnis").
- `SHA256SUMS`: 54 Einträge, `sha256sum -c` direkt: 54/54 OK.
- Suchlauf-Feld im Plan: `make suchlauf-nachmessen PLAN=…` — **6/6 OK,
  Exit 0** (Stände `ad530688` und `diff` gegen `caf96849`).
- Bundle-Vergleich `v6.9.0` (`7103ad59`) ↔ `v6.13.0` (`d443ee39`) per
  Blob-Vergleich und `diff -rq`: 55 Pfade je Baum, Pfad-Menge identisch,
  **35 Dateien inhaltlich unterschiedlich** (26 `regelwerk/` + **8
  Templates** + 1 `SHA256SUMS`), **20 Templates byte-gleich** — darunter
  beide Reviewer-Skills (Bestätigung der Audit-Aussage). `diff -ru`: 768
  Zeilen (`regelwerk/`), 285 Zeilen (`templates/`) — beide Audit-Zahlen
  bestätigt. Neue `##`/`###`-Header in `regelwerk/`: genau einer
  („Nachzug ist keine Überschreibung", Modul 4) — bestätigt; 0 entfernte
  Header, alle `<a id>`-Anker identisch — bestätigt.
- `make gates`: **Exit 0** (ungepiped gesichert, §3.9).

## Findings

### F-1 — Audit §2: Diff-Umfang-Zahlen driften gegen die Messung

- `kategorie`: HIGH
- `quelle`: [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
  (§3.12 — „gemessen"-Label) / Maintainability
- `pfad`: docs/reviews/audit-baseline-v6-13-0-drift.md:32-36
- `befund`: Audit §2 nennt „36 Dateien — 26 `regelwerk/`-Dateien, 9
  Templates, 1× `SHA256SUMS`" und „19 Templates sind byte-gleich"; gemessen
  per Blob-Vergleich und `diff -rq` sind es **35 / 8 / 20**. Der Audits
  eigene §3-Fundtabelle benennt mit „übrige 5" + drei genannte Templates
  ebenfalls **8** geänderte Templates — §2 widerspricht §3. Alle operativen
  Teile des Audits (§3 Fundtabelle, §4 Nachzugs-Entscheidungen, §5
  Empfehlung) sind von dem Zähler-Defekt unabhängig und wurden gegen eigene
  Messung bestätigt.
- `verifizierbar`: ja — Blob-Vergleich (`git ls-tree` beider Bundle-Stände)
  bzw. `diff -rq`; kein Gate, der Audit-Bericht liegt in der gate-exempten
  Fläche
- `klasse`: „Zahl im Träger driftet gegen die Messung"

**F-1-Verifikation nach Fixrunde:** `9f1eb320` korrigiert §2 auf **35 / 8 /
20** und 768/285 — byte-gleich zur eigenen Messung dieses Laufs; der Ursprung
ist jetzt im Bericht deklariert (Blob-Vergleich gegen `d443ee39`, die
Fehlzählung 36/9/19 der Erstfassung benannt). F-4-Lesart des Architects
(Kurzform trägt, Spanne Abschnitte-Liste vs Kurzform als Beobachtungs-Eintrag
in die Closure, engere Lesart für neue Fälle): im Audit §4.2/§7 nachgezogen,
[`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md) selbst
unberührt — F-1 und F-4 sind damit geschlossen.

### F-2 — Plan-Suchlauf-Deutung: Pin-Muster trifft einen der vier Treffer

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12/§3.13
- `pfad`: docs/plan/planning/in-progress/slice-harness-baseline-v6-13-0.md:148-153
- `befund`: Der Plan begründet das Belassen der vier Rest-Treffer mit
  „`versions.pin-pattern` trifft keine der vier" — falsch für einen Treffer:
  `exemption-ohne-reifegrenze/evidence/slice-105.md:3` trägt
  `.harness/baseline/v6.5.0/`, was das generische Pin-Muster matcht; der
  Befund wird durch den `d-check:ignore`-Marker derselben Zeile unterdrückt,
  nicht durch Pin-Abwesenheit. Die Schlussfolgerung („`make docs-check`
  meldet 0 Befunde"; Records bleiben) ist über die zweite, gemessene Stütze
  weiterhin getragen — Gates grün.
- `verifizierbar`: ja — `git grep -E '\.harness/baseline/v[0-9.]+/'` gegen
  die vier Dateien; `make docs-check` (grün) für die Unterdrückung
- `klasse`: „Beleg-Befehl trägt seinen Satz nicht" (abgeschwächt — Schluss
  über zweite Stütze getragen)

### F-3 — DoD-2-Text „in einen Commit" — Umsetzung in zwei Commits

- `kategorie`: LOW
- `quelle`: Maintainability (Plan-Text vs Umsetzung)
- `pfad`: docs/plan/planning/in-progress/slice-harness-baseline-v6-13-0.md:59-62
- `befund`: DoD 2 legt „Mechanik und Entfernung in einen Commit" fest;
  umgesetzt als `d443ee39` (Bundle) + `88cea828` (Entfernung) — mit
  dokumentiertem, erwartbarem rot-Zustand von `make baseline-verify` am
  Zwischenstand (Commit-Message `d443ee39`, Audit §7 Risiko 1). Die
  Abweichung ist in Commit-Message und Audit begründet, aber der Plan-Text
  wurde nicht nachgezogen; ein Intermediate-Commit mit rotem Sensor bleibt
  die Kosten (Rename-Detection unberührt, wie der Plan selbst feststellt).
- `verifizierbar`: ja — Commit-Liste der Range; kein Gate
- `klasse`: „Plan-Abweichung ohne Nachzug im Plan-Träger"

### F-4 — Zitat-Korrektur an
[ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md) berührt
§Verglichene Alternativen

- `kategorie`: LOW
- `quelle`: [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1 / `AGENTS.md` §3.5
- `pfad`: docs/plan/adr/0095-review-klasse-exempt-status-check.md:107
- `befund`: Die Pfad-Korrektur sitzt in der Options-Tabelle von
  §Verglichene Alternativen — einem Abschnitt, den sowohl
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1 als auch `AGENTS.md` §3.5 als unberührbar listen (ohne
  den Zusatz „der Aussage nach", den §Konsequenzen trägt). Nach der
  Kurzform von
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  („das Gerüst darf sich ändern, die Aussage nie; der Referent bleibt
  derselbe") ist die Änderung zulässig — Referent unverändert **gemessen**
  (die `v6.13.0`-Vorlage `templates/.d-check.yml` zeigt `matrix.status`
  weiterhin nur klassen-übergreifend, Zeile 80), Commit nennt
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
  §Geschichte-Zeile gesetzt — aber die Spanne zwischen Abschnitte-Liste
  und Kurzform bliebe zu klären, sonst liest der nächste strenge Lauf den
  Commit als §3.5-Verstoß.
- `verifizierbar`: ja — `git show 00d96eb7` gegen die §3.5-Liste;
  Referent-Check: `grep status` auf
  `v6.13.0` · `templates/.d-check.yml`
- `klasse`: „Zitat-Korrektur-Reichweite: Abschnitte-Liste vs Kurzform"

### F-5 — §4-Rückführungs-Trigger feuerte; Hold-Entscheidung ohne Ausgangs-Feld

- `kategorie`: LOW
- `quelle`: Maintainability (Slice-Plan §4)
- `pfad`: docs/plan/planning/in-progress/slice-harness-baseline-v6-13-0.md:108-112
- `befund`: Der Rückführungs-Trigger „zu groß" (Nachzüge an mehr als den
  beiden §3-Dateien) feuerte real — das Audit zählt neun Nachzugs-Dateien.
  Das Halten des Slices ist inhaltlich gut begründet und committet
  (Audit §5: alle neun mechanische Ein-Zeilen-Bumps, ohne sie bleibt
  `make gates` rot; „Empfehlung: Slice halten"), aber die Ausführung der
  Abweichung hängt im committeten Träger allein an der Architect-Empfehlung
  („Die Entscheidung liegt beim Planner"); der Plan trägt für §4 kein
  Ausgangs-Feld. Die Entscheidung muss in der §7-Closure-Notiz namentlich
  als Planner-Entscheidung festgehalten werden, sonst überlebt sie nur als
  Empfehlungssprache.
- `verifizierbar`: nein — Abwägung und Closure-Feld sind Lese-Handlung
- `klasse`: „Rückführungs-Entscheidung ohne Ausgang im Plan"

## Negativbefunde

- geprüft, ohne Befund: `.harness/baseline/v6.13.0/**` — 55 Pfade,
  Pfad-Menge identisch zum Vorgänger, 54/54 Summen OK, Asset byte-gleich
  reproduziert, 0 Struktur-Abweichungen (Header/Anker)
- geprüft, ohne Befund: Suchlauf-Zahlen im Plan-Feld — `make
  suchlauf-nachmessen` 6/6 OK, die vier Rest-Treffer einzeln identifiziert
  und als Records (historische Vorgangs-Angaben `slice-105` + Dateiname +
  Register-Vergleichsnennung) bestätigt
- geprüft, ohne Befund: verkörperte Form — `AGENTS.md` §1,
  `harness/conventions.md` §Baseline + MR-000-RB-Ergänzung, MR-001…MR-004,
  `.claude/agents/*`, `.harness/skills/closure-note-reviewer.md`,
  `harness/sensors/baseline-verify.md`: durchweg Ein-Zeilen-Bumps der
  audit-forcierten Art, keine unterbliebenen Regel-Nachzüge (Audit §4.3
  Gegenprobe gelesen)
- geprüft, ohne Befund: §3.5-Disziplin außer F-4 — §Entscheidung,
  §Konsequenzen, §Status und `Supersedes`-Kette der
  [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md)
  unberührt; §Geschichte-Zeile mit Datum, Ereignis, Commit-Kennung gesetzt;
  beide Commits nennen
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
- geprüft, ohne Befund: §3.7 Kommentar-Klassen — der Diff trägt keine
  Code-, Config- oder Skript-Änderungen; keine Kommentar-Fundstelle
- geprüft, ohne Befund: §3.1 Docker-only — kein Host-Toolchain-Aufruf, kein
  Host-Schreibweg auf Repo-Dateien im Diff; Asset-Download als
  Container-Lauf reproduziert
- geprüft, ohne Befund: Suppression — keine neuen Suppress-Marker
  (`//nolint`, `noqa`, …) im Diff
- geprüft, ohne Befund: Traceability — alle 8 Commits der Range nennen
  [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
  (beide [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md)-Commits
  zusätzlich
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md));
  kein `SPEC-*`/`ARC-*` im Subject
- geprüft, ohne Befund: Spec-Straten — `spec/**` unberührt; „ausdrücklich
  NICHT" (Sensor-/Gate-Verträge): `.d-check.yml` unberührt, die
  `baseline-verify.md`-Zeile ist Referenz-Bump, kein Vertrags-Text

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Zahl im Träger driftet gegen die Messung ·
Beleg-Befehl trägt seinen Satz nicht · Plan-Abweichung ohne Nachzug im
Plan-Träger · Zitat-Korrektur-Reichweite: Abschnitte-Liste vs Kurzform ·
Rückführungs-Entscheidung ohne Ausgang im Plan

## Verdikt

**Merge-blockierend:** nein — F-1 ist mit `9f1eb320` behoben und gegen die
eigene Messung dieses Laufs verifiziert (35/8/20, 768/285, Ursprung
nachgetragen); F-4 ist im selben Zug mit der Architekt-Lesart geschlossen.
F-2/F-3/F-5 gehen als Closure-/Plan-Punkte an den Planner (F-2-Hinweis zur
`d-check:ignore`-Wirkung, F-3 Plan-Text-Nachzug, F-5 Ausgang im Plan bzw.
§7-Closure-Notiz).

**Übergabe:** Findings gehen an die Haupt-Rolle zur Weiterleitung (F-2/F-3/F-5
an den Planner für die Closure-Notiz; die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler; die
`ADR-0073`-Klärung ist als Beobachtungs-Eintrag (`BEO-PGC`) bei der Closure
zu öffnen). DoD-Checkbox „Review durchgeführt" wurde nach der Skill-Regel
(keine Fixrunde am Implementer nötig — F-1 über den Architect behoben)
vom Reviewer selbst nachgezogen. DoD-/Spec-Konformität prüft der Verifier
separat (Modul 11).
