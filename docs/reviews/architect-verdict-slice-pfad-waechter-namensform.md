# Architect-Verdikt: Slice-Pfad-Wächter auf die Namens-Kennung — keine ADR

**Rolle:** Architect (Modul 8)

**Anlass:** Liefer-Punkt 3 von `slice-abgeleitete-dokumente-vorlagen-nachzug`
(§2, §3) verlangt vor der Umsetzung eine Übergabe an den Architect: Braucht die
Erweiterung des `forbid-pattern` der zwei `structure`-Regeln gegen
Slice-Pfad-Links in `.d-check.yml` von der Nummern-Form `slice-[0-9]{3}` auf die
Namens-Form eine ADR ([`AGENTS.md`](../../AGENTS.md) §3.6,
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`)?

**Rolleninhaber:** pt9912 (Architect-Zug; anderer Eingabe-Kontext als der
Planner-Lauf, der den Slice angelegt hat).

**Datum:** 2026-10-06

**Bezug:** [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md)
(Bezug des Slice) · [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(Verfeinerung gegen Erweiterung) ·
[`ADR-0099`](../plan/adr/0099-slice-welle-review-regel-zurueckgenommen.md)
(zweiter Beleg des Register-Eintrags) ·
[`MR-002`](../../harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md)
(Kennungsform) · [`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md)
§Grenze Punkt 6 und §Bindung · `BEO-PGC/slice-pfad-als-link-in-berichten`
(Träger der Regeln) · `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`.

---

## Frage

Ist die Ausdehnung des Musters auf `…/(open|next|in-progress)/slice-<Name>.md`
eine Änderung, die nach §3.6 oder nach der Klasse des Register-Eintrags
`gate-scope-erweiterung-ohne-adr-traeger` eine ADR braucht?

## Befund — wo die Regeln ihre Bindung tragen

Gemessen am Stand `187596e3`:

- **Entstehung.** `git log -S'(open|next|in-progress)/slice-[0-9]{3}' -- .d-check.yml`
  nennt genau einen Commit, `7ef760b5` (2026-09-15), der `.d-check.yml` und
  `harness/sensors/docs-check.md` ändert — keine ADR im Commit.
- **Keine ADR trägt die Regeln.** `git grep -nE 'slice-pfad-als-link|\(open\|next\|in-progress\)|slice-\[0-9\]\{3\}' -- docs/plan/adr`
  liefert 0 Treffer. Keine `Accepted`-ADR legt das Muster fest oder schließt die
  Namens-Form aus.
- **Bindung.** Die Regeln sind die Verkörperung von
  `BEO-PGC/slice-pfad-als-link-in-berichten` (Ausgang *verkörpert*,
  `· seit slice-075`, Kommentar in `.d-check.yml` und `state.md` des Eintrags);
  `harness/sensors/docs-check.md` §Bindung führt sie unter genau diesem Anker.
  Gegenstand der Regel ist laut `state.md` und `hint` **ein Markdown-Link auf
  einen Slice-Plan in einem Lifecycle-Verzeichnis** — nicht eine bestimmte
  Kennungsform. Die Nummern-Form ist die Abbildung des damals einzigen
  ID-Schemas; dass sie die Namens-Form seit `MR-002` nicht mehr abdeckt, steht
  als Grenze (5) bzw. Punkt 6 im Sensor-Vertrag — eine benannte Lücke, keine
  Entscheidung.

## Abgrenzung zum Register-Eintrag

Die Klasse von `gate-scope-erweiterung-ohne-adr-traeger` ist in beiden Belegen
eine andere:

| Beleg | Richtung | Was die ADR auslöste |
|---|---|---|
| `slice-071` | Lockerung: `composition_root` um `tools/**` gab **Importrecht** | [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md) Festlegung 3: ADR-pflichtig ist die Aufnahme eines Bereichs **mit Import-Berechtigung**; das Verfeinern bestehender Regeln bleibt ADR-frei |
| `ADR-0099` | Verschärfung, aber **im Widerspruch** zu `ADR-0094`/`ADR-0097`, die genau diese Regel bewusst ausgelassen hatten | die Kollision mit `Accepted`-ADRs (§3.5), nicht die Verschärfung als solche — §Entscheidung von `ADR-0099` stellt keine allgemeine ADR-Pflicht für Verschärfungen auf |

Der vorliegende Fall ist keines von beiden: Er gibt keine Berechtigung, er senkt
keine Schwelle, und keine ADR hat die Namens-Form bewusst ausgelassen. Er bildet
eine bestehende, ADR-frei verkörperte Regel auf das geltende ID-Schema ab — die
Verfeinerung im Sinn von `ADR-0068`, übertragen auf `.d-check.yml`.

## Verdikt — (a) keine ADR

**Gewählt: (a).** Die Erweiterung ist eine Verschärfung innerhalb der bestehenden
Bindung. Die **Bindung** ist `BEO-PGC/slice-pfad-als-link-in-berichten`
(*verkörpert*, `· seit slice-075`) samt `harness/sensors/docs-check.md`
§Bindung; die Kennungsform liefert
[`MR-002`](../../harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md).
[`AGENTS.md`](../../AGENTS.md) §3.6 verlangt eine ADR für die Senkung einer
Schwelle; hier wird nichts gesenkt.

**Liefer-Punkt 3 bleibt im Slice**; die Rückführung `in-progress → next` aus §4
des Plans tritt nicht ein.

**Register:** kein Beleg für `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` —
dieser Zug ist kein Auftreten der Klasse (keine Lockerung, kein Widerspruch zu
einer `Accepted`-ADR), der Zähler bleibt bei 2×. Akzeptiertes Negativ: ob die
Unterscheidung Verfeinerung/Erweiterung für `.d-check.yml` eine eigene Regel
bekommt, entscheidet der Lese-Schritt bei 3×; ein Einzelfall ohne Kollision
rechtfertigt keinen neuen Vorgang.

**Was die Umsetzung trägt (Constraint für den Implementer):**

1. Beide `forbid-pattern` werden gleich erweitert; erprobte Form (siehe unten):
   `'\]\([^)]*(open|next|in-progress)/slice-[a-z0-9]'` — sie enthält die
   Nummern-Form als Teilmenge.
2. Grenze (5) im Kommentar von `.d-check.yml` und der entsprechende Satz in
   `harness/sensors/docs-check.md` §Grenze Punkt 6 werden gestrichen; der
   Kommentar-Anker `· seit slice-075` bleibt.
3. Die übrigen Grenzen (1) bis (4) bleiben unverändert.

## Wirkungsmessung

Lesend am Stand `187596e3`, über die zwei Dateiklassen der Regeln:

```text
git grep -nE '\]\([^)]*(open|next|in-progress)/slice-' -- 'docs/reviews/**/*.md' 'docs/plan/planning/observations/**/observation.md'
```

**1 Treffer** (Textform): `BEO-PGC/slice-pfad-als-link-in-berichten/observation.md`
Zeile 10, `` `[…](../plan/planning/in-progress/slice-NNN-….md)` `` — in
Inline-Code. Die Nummern-Form trifft in derselben Messung 0 Zeilen.

Ob der Treffer die Regel färbt, ist **gemessen**, nicht hergeleitet: d-check
(`ghcr.io/pt9912/d-check@sha256:b4b8756b…`, der Pin aus `d-check.mk`) lief
gegen einen Klon von `187596e3` im Scratchpad:

| Lauf | Konfiguration | Inhalt | Exit | `section-forbidden` |
|---|---|---|---|---|
| alt | Muster `slice-[0-9]{3}` | Bestand | 0 | 0 |
| weit | Muster `slice-[a-z0-9]` | Bestand | 0 | 0 |
| Mutation alt | Muster `slice-[0-9]{3}` | je ein Link auf `…/in-progress/slice-abgeleitete-dokumente-vorlagen-nachzug.md` in einem `docs/reviews/`-Bericht und in einer `observation.md` | 0 | 0 |
| Mutation weit | Muster `slice-[a-z0-9]` | dieselben zwei Links | 1 | 2 (je Datei einer) |

**Ergebnis:** 0 Bestands-Treffer, die die Regel sehen würde — der eine
Text-Treffer liegt in Inline-Code und ist für die Regel unsichtbar (Grenze (2)).
Im selben Slice ist kein Bestand zu beheben. Die Mutation bestätigt die Lücke
(alt grün) und ihre Schließung (weit rot, an beiden Stellen); erprobt an beiden
Dateiklassen, je eine Instanz. Dass die Links im Mutationslauf alt auch
`links` nicht rot färbten, ist gemessen: das Ziel existiert am Ist-Ort, der
Bruch käme erst mit dem nächsten `git mv`.

Die Mutation des Implementers (DoD Liefer-Punkt 3) bleibt eine Erwartung an die
Umsetzung; dieser Lauf ersetzt sie nicht.
