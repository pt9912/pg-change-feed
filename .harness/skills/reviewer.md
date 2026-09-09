# Reviewer-Skill — PG Change Feed

* Status: Accepted
* Bezug: `AGENTS.md` §3 Hard Rules und §5 Dokumentations-Regeln ·
  `harness/conventions.md` (MR-000 ID-Schema, MR-001 Pflichtenheft-Umbenennung) ·
  <!-- d-check:ignore (Kurs-/ADR-Referenzen; Anker gelten im Ziel-Repo) -->
* Gilt für: kein eigenes Make-Target — getragen vom Reviewer-Agenten
  (`.claude/agents/reviewer.md`); Reports nach `docs/reviews/`

## Kontext-Eingang (Pflicht)

Was der Reviewer *immer* mitbringt, bevor er den Diff liest:

- Diff des PR
- `spec/lastenheft.md` (für referenzierte `LH-*`-IDs)
- `spec/pflichtenheft.md` (für Verfeinerungen `LH-*.<a>` und Festlegungen
  `SPEC-*`)
- ADRs, deren ID im PR oder in der Commit-Message vorkommt
- `AGENTS.md` §"Hard Rules"
- `harness/conventions.md` (MR-000/MR-001 — ID-Schema und Datei-Namen)
- vorherige Findings am gleichen Modul (letzte ~5 PRs)

Ohne diesen Block sieht der Reviewer den Code, aber nicht *die Verträge, gegen
die er prüft*.

## Klassifikation

Jeder Anker HIGH/MEDIUM/LOW hat eine *konkrete* Liste — nicht generisch. INFO ist
bewusst kurz (Ergänzungs-Kanal, nicht Hauptkanal).

**HIGH** — eines der folgenden:
- ADR-Verstoß (Layer, Tool, Hard Rule)
- Sicherheits-Anti-Pattern (Injection, fehlende Auth-Prüfung)
- Korrektheitsfehler im *kritischen* Pfad (Persist-before-ACK-Verletzung,
  Retention löscht benötigte Changes, Lesen verändert Positionen)
- Suppression eines Gates (`#noqa`, `//nolint`, `[SuppressMessage]`) ohne ADR
- **Norm nur im Template-Kommentar** — eine Regel steht im `<!-- -->`-Block
  eines `.template.md` und nirgends sonst. Sie ist beim Adopter weg, sobald er
  die Kommentare entfernt. Kein Gate fängt das (siehe Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Template-Schichtung)
- **Kommentar trägt keine der Kommentar-Klassen** — ein Kommentar in Code, Config
  oder Skript beschreibt die verworfene Alternative („Ohne X wäre …"), einen
  abwesenden Text („früher stand hier …") oder bricht mitten im Satz ab, weil
  eine Teilersetzung den Rest stehen ließ. Kein Gate fängt das (siehe
  Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein Kommentar trägt)
- **Zustandsfeld trägt Chronik** — eine `Stand`-/`Status`-Zelle (Roadmap,
  Beobachtungs-Register, Meilenstein) erzählt, wie der Zustand entstand, statt
  Zustand und Beleg als Anker zu nennen; oder ein Drift-Log protokolliert
  Schließungen und erreichte Meilensteine. Kein Gate fängt das (siehe
  Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein Kommentar trägt,
  *Dieselbe Regel für Zustandsfelder*)
- **Traceability-/ID-Schema-Verstoß** — Commit oder PR nennt keine
  `LH-*`- oder `ADR-*`-Kennung; oder eine Kennung nutzt ein Präfix, das MR-000
  nicht deklariert (`LH-FA/QA-<BEREICH>-<NNN>`, `SPEC-<NNN>`, `ARC-<NNN>`,
  `ADR-<NNNN>`, `CO-<NNN>`, `slice-<NNN>`, `MR-<NNN>`, `BEO-<KUERZEL>/<slug>`,
  `RC-<NNN>`). Erstes Auftreten dieser Klasse: Review F-2 (2026-09-09,
  `PH-*`/`TST-*` im Traceability-Beispiel).
- **Spec-Stratum-Verstoß** — das Technik-Stratum erweitert, wo es nur
  präzisieren darf („präzisieren ja, erweitern nie“): das Pflichtenheft führt
  eine neue bindende Anforderung ein, statt eine bestehende `LH-*`-ID zu
  schärfen; oder das Lastenheft wird geändert, ohne dass die Änderung in
  einem eigenen Commit vor dem umsetzenden Slice liegt. Erstes Auftreten
  dieser Klasse: Review F-3 (2026-09-09, Zurückstellung bindender
  Anforderungen im Sammelabschnitt).
- **Docker-only-Verstoß** — ein Build-/Test-/Betriebs-Skript installiert ein
  lokales Toolchain (venv, SDK, Paketmanager) statt über `make` zu laufen
  (`AGENTS.md` §3.1).
- **Zwei-Quellen-Drift** — derselbe Zustand wird in zwei Dateien geführt,
  ohne dass der Gewinner deklariert ist (z. B. Anforderungstext doppelt in
  Lastenheft und Pflichtenheft; Zustand in Verzeichnis *und* Status-Feld).

**MEDIUM** — eines der folgenden:
- unklare Fehlerbehandlung am Rand des Spec-Bereichs
- fehlende Negativtests bei neuem öffentlichem Vertrag
- Wiederholung eines Musters, das schon zweimal LOW war
- **Delegation ohne Entsprechung** — das Lastenheft delegiert eine Festlegung
  an das Pflichtenheft (`LH-QA-PER-002`, `LH-QA-PER-004`, `LH-QA-POR-001`),
  dort steht sie nicht, auch nicht als offene Festlegung mit Schlusspunkt.
  Erstes Auftreten: Review F-4 (2026-09-09).
- **GWT-Pfad unvollständig** — ein Akzeptanzkriterium fehlt ganz (`—` ohne
  die `—`-Regel aus Lastenheft §3), bricht das Given/When/Then-Muster oder
  verfehlt die drei Pfade Happy/Boundary/Negative. Erstes Auftreten:
  Review F-5 (2026-09-09).
- **MVP-Kennzeichnung inkonsistent** — `MVP: ja`-Marker und die
  MVP-Abnahme-Mapping-Tabelle (Lastenheft §1) weichen voneinander ab.
  Erstes Auftreten: Review F-1 (2026-09-09).

**LOW** — stilistisch unschön ohne semantische Auswirkung, einmalige Tippfehler,
unbenutzte Imports.

**INFO** — Hinweis ohne erwartete Aktion (z. B. „diese Stelle hat ein passendes
ArchUnit-Pendant, das du nicht kennst“; „Metrik deckt die Anforderung nur
indirekt ab“ — erstes Auftreten: Review F-9, 2026-09-09).

> **Pflicht beim Ausfüllen (Modul 10 §Übungen):** Die HIGH-Liste muss mindestens
> *zwei* repo-spezifische Regeln nennen, die ein generischer Skill nicht abdeckt.
> Ist der Skill ohne sie, ist er noch nicht scharf genug — dann kommt bei einem
> Lauf auf einem realen Diff keines deiner Repo-HIGHs zur Anwendung.
> **Erfüllt seit 2026-09-09:** vier repo-spezifische HIGH-Regeln
> (Traceability/ID-Schema, Spec-Stratum, Docker-only, Zwei-Quellen-Drift).

## Was dieser Skill NICHT macht

- Keine Lösungsvorschläge („schreib das so“) — Reviewer kategorisiert,
  Implementer entscheidet.
- Kein Refactoring-Vorschlag, der über den Diff hinausgeht.
- Keine Verifikation gegen DoD — das ist Verifier-Aufgabe (Modul 11).
- Keine Validation gegen reale Bedürfnisse — das ist Validator-Aufgabe.

Wenn etwas auffällt, das in diese Kategorien gehört: ein INFO-Finding mit Verweis
auf die zuständige Rolle.

## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: ADR-ID, `LH-*`-ID, Hard-Rule-Name oder „Maintainability“
- `pfad`: Datei:Zeile
- `befund`: 1–2 Sätze, beobachtbar, ohne Lösungsvorschlag
- `verifizierbar`: ja/nein — gibt es einen Gate-Lauf, der es bestätigen würde?
- `klasse`: stabile Kurz-Bezeichnung des Fehlermusters, z. B. „Delegation ohne
  Entsprechung“ — speist den Steering-Loop-Zähler (siehe §Pflege)

Zusätzlich am Ende: eine Zeile „geprüft, ohne Befund“ pro betrachtetem
Verzeichnis (Negativbefund-Zeile — sonst ist „keine Findings“ nicht von „nicht
geprüft“ unterscheidbar). Report-Gerüst für den ganzen Lauf:
`docs/reviews/review-report.template.md`, ein Report pro Lauf, Folgeläufe als
neue Datei statt Überschreibung.

## Pflege (Steering-Loop)

Bei dreimaligem Auftreten desselben Findings:

- ist die Kategorie noch richtig? → Klassifikation schärfen
- gibt es einen ADR/`AGENTS.md`-Eintrag, der das verhindert hätte?
  → Folge-ADR oder `AGENTS.md`-Update
- gibt es eine Fitness Function, die das prüfen würde? → Modul 13, Gate hinzufügen

Diese Skill-Datei wird **nicht** überschrieben, sondern versioniert
(ADR-Hard-Rule, Modul 4). Erste Schärfung 2026-09-09 aus dem
Eröffnungs-Review `docs/reviews/review-lastenheft-pflichtenheft.md`.