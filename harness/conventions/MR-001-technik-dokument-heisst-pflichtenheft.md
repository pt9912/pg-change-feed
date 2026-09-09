# MR-001 — Technik-Dokument heißt Pflichtenheft

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger.

- **Datum:** 2026-09-09
- **Geltungsbereich:** `spec/pflichtenheft.md` (vormals
  `spec/spezifikation.md`), `harness/README.md` §Source precedence und
  §Guides, `AGENTS.md` §2 und §5, `spec/lastenheft.md` (Verweise auf das
  Technik-Dokument)
- **Ersetzt-Baseline-Regel:** [grundlagen-source-precedence.md
  §Spec-Straten](../../.harness/baseline/v6.5.0/regelwerk/grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument)
  — der Datei-Name des Rang-2-Dokuments (`spec/spezifikation.md`)
- **Adaption:** Das Rang-2-Dokument heißt
  `spec/pflichtenheft.md` statt `spec/spezifikation.md`. Inhalt und Struktur
  (Abschnitte 1–7, IDs `LH-*.<a>` für Verfeinerungen, `SPEC-<NNN>` für
  Festlegungen) sind unverändert; die Rang-Ordnung bleibt neunstellig. Die
  `SPEC-<NNN>`-Präfixe bleiben bestehen — sie kodieren das Stratum und sind
  laut Baseline fest; nur Datei-Name und Dokument-Titel ändern sich.
- **Begründung:** V-Modell-Terminologie — das Pflichtenheft ist das
  Antwort-Dokument des Auftragnehmers auf das Lastenheft; genau diese Rolle
  hat das Technik-Stratum hier. Das Traceability-Muster des Lastenhefts
  (§3) adressiert die Verfeinerungen dieses Dokuments als `LH-*.<a>` bzw.
  `SPEC-*`.
- **Auflösungs-Trigger:** permanent.