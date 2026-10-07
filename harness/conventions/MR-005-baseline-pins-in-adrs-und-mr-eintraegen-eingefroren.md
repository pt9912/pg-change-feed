# MR-005 — Baseline-Pins in ADRs und MR-Einträgen bleiben auf dem Stand ihrer Abfassung

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger;
`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen
früheren ablöst. `Ersetzt-Baseline-Regel` nennt **genau eine** Regel der
Baseline, an deren Stelle dieser Eintrag tritt — als Link mit
Abschnitts-Anker in die vendored Fassung; ein Datei-Link benennt keine Regel.

- **Datum:** 2026-10-07
- **Geltungsbereich:** `.d-check.yml` Block `versions:` (`exempt-paths`) und
  Block `vcs:`, die Baseline-Pins in `docs/plan/adr/[0-9]*.md` und
  `harness/conventions/MR-[0-9]*.md`, der Bump-Ablauf in
  `harness/targets/pin-stale.md`, [`AGENTS.md`](../../AGENTS.md) §3.5 Absatz
  „Beleg“.
- **Ersetzt-Baseline-Regel:** [`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln](../../.harness/baseline/v6.16.0/regelwerk/grundlagen-traceability.md#herkunfts-anker-für-steering-loop-regeln),
  Spiegelstrich „Zwei Rot-Quellen, ein Prinzip“, die Versions-Hälfte („ein
  Versions-Sensor prüft jeden Pin dagegen“).
- **Adaption:** Baseline-Pins in ADR-Dateien und MR-Einträgen nennen den Stand
  ihrer Abfassung; ein Bump zieht sie nicht nach, und der Versions-Sensor nimmt
  beide Pfade aus. Die Existenzprüfung der Links und Anker bleibt; bricht ein
  Link, weil ein altes Tag-Verzeichnis gelöscht wird, wird er einmal zu
  Inline-Code mit Pfad, alter Version und Anker, an MR-Einträgen in einem
  eigenen Form-Commit vor dem Lösch-Commit. Ob sich die ersetzte Regel eines
  MR-Eintrags im Bump bewegt hat, misst der Adaptions-Durchgang des Bumps
  zwischen den Tags des Bumps. Entscheidung, Grenze und Messungen:
  [`ADR-0161`](../../docs/plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md).
- **Begründung:** ADRs ab `Accepted` und MR-Einträge ab ihrer `Datum`-Zeile
  sind immutabel. Ein `version-stale` in ihnen entsteht durch einen Bump
  außerhalb der Datei, nicht durch ihren Text; bei geändertem Referenten oder
  bei einem Pin, der einen Messbeleg trägt, lässt er sich nicht per
  Zitat-Korrektur nachziehen (Anlass: Bump auf v6.16.0, `MR-001`, `ADR-0095`
  Zeile 107, `ADR-0160` Zeile 278). Ein Pin in einer ADR nennt die Fassung,
  die die Entscheidung las.
- **Auflösungs-Trigger:** permanent — Re-Evaluierungs-Trigger (a) von
  `ADR-0161`: d-check bietet für `versions` eine Bedingung auf den Inhalt
  einer Datei; Trigger (c) macht den Eintrag im Adaptions-Durchgang
  gegenstandslos.
