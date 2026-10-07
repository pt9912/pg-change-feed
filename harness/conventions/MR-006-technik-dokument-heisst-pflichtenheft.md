# MR-006 — Technik-Dokument heißt Pflichtenheft

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger;
`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen
früheren ablöst. `Ersetzt-Baseline-Regel` nennt **genau eine** Regel der
Baseline, an deren Stelle dieser Eintrag tritt — als Link mit
Abschnitts-Anker in die vendored Fassung; ein Datei-Link benennt keine Regel.

- **Datum:** 2026-10-07
- **Geltungsbereich:** `spec/pflichtenheft.md` (vormals
  `spec/spezifikation.md`), `harness/README.md` §Source precedence und
  §Guides, [`AGENTS.md`](../../AGENTS.md) §2 und §5, `spec/lastenheft.md`
  (Verweise auf das Technik-Dokument)
- **Ersetzt-Baseline-Regel:** [`grundlagen-referenz-richtung.md` §Spec-Straten](../../.harness/baseline/v6.16.0/regelwerk/grundlagen-referenz-richtung.md#spec-straten-mehr-als-ein-spec-dokument)
  — der Datei-Name des Rang-2-Dokuments (`spec/spezifikation.md`, Spalte
  „Default-Datei“ der Straten-Tabelle)
- **Adaption:** Das Rang-2-Dokument heißt `spec/pflichtenheft.md` statt
  `spec/spezifikation.md`. Inhalt und Struktur folgen der
  Spezifikations-Vorlage der Baseline: Abschnitte 1–8, darunter §7
  „Festlegungen der Harness-Werkzeuge“ und §8 Historie; IDs `LH-*.<a>` für
  Verfeinerungen, `SPEC-<NNN>` für Festlegungen. Die Rang-Ordnung bleibt
  neunstellig. Die `SPEC-<NNN>`-Präfixe bleiben bestehen — sie kodieren das
  Stratum und sind laut Baseline fest; nur Datei-Name und Dokument-Titel
  ändern sich.
- **Begründung:** V-Modell-Terminologie — das Pflichtenheft ist das
  Antwort-Dokument des Auftragnehmers auf das Lastenheft; genau diese Rolle
  hat das Technik-Stratum hier. Das Traceability-Muster des Lastenhefts
  (§3) adressiert die Verfeinerungen dieses Dokuments als `LH-*.<a>` bzw.
  `SPEC-*`.
- **Auflösungs-Trigger:** permanent.
- **Löst auf:** [`MR-001`](../conventions.md#mr-001) — dessen Satz
  „Inhalt und Struktur (Abschnitte 1–7 …) sind unverändert“ trifft nicht mehr
  zu, seit das Pflichtenheft §7 „Festlegungen der Harness-Werkzeuge“ trägt
  und die Historie §8 ist; die Adaption selbst, der Datei-Name, gilt weiter.
- **Ausgelöst durch Baseline-Stand:** v6.16.0
