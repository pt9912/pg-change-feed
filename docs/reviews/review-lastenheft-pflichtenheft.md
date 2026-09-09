# Review-Report: Spec-Docs Lastenheft + Pflichtenheft — 2026-09-09

**Review-Art:** Plan-Review (Spec-Dokumente gegen Ziel-Form-Vorlagen, MR-000/MR-001,
Hard Rules) — *wogegen*: Vorlagen-Konformität, ID-Schema, Querverweise,
Dokument-Konsistenz. Keine DoD-Prüfung (Verifier, Modul 11).

**Gegenstand:** HEAD-Stand `dff9b2f` — `spec/lastenheft.md` (0.3.0, Draft;
Entstehung 577d510 → 893fc95, Nachzug 40c8c43) und `spec/pflichtenheft.md`
(265a5f8 → 288b3c5 reiner Rename → 40c8c43; Vorläufer c70ee1c, zurückgezogen
in 5a6f8ea)

**Skill:** `.harness/skills/reviewer.md` — trägt noch Vorlagen-Platzhalter
(repo-spezifisch #1/#2 unausgefüllt); dieser Lauf lief auf
**Baseline-Regeln**: `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Ziel-Form: Reviewer-Skill

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- `spec/lastenheft.md`, `spec/pflichtenheft.md`, `spec/architecture.md` (untracked, Vorlage)
- Ziel-Form: `v6.5.0` · `templates/spec/lastenheft.template.md`, `templates/spec/spezifikation.template.md`
- `AGENTS.md` §3 (Hard Rules, insb. §3.7), §5; `harness/README.md`; `harness/conventions.md` inkl. MR-000/MR-001
- Entstehungs-Diffs 577d510, 893fc95, c70ee1c, 5a6f8ea, 265a5f8, 288b3c5, 40c8c43
- ADRs: keine im Verzeichnis (`docs/plan/adr/` leer); keine ADR-ID in den Commits

> **Standort-Hinweis (dieser Lauf):** Baseline `v6.5.0` ·
> `regelwerk/modul-10-review-harness.md` — „Abgelegt wird ein Report pro Lauf
> unter `docs/reviews/`" (dort existiert bereits eine leere Ablage mit
> `.gitkeep`). Dieser Report liegt auf Anweisung des Auftrags unter
> `docs/plan/reviews/`; die Pfad-Frage ist vom Implementer/Planner zu
> vereinheitlichen (ein Ort, nicht zwei).

---

## Findings

### F-1 — MVP-Marker-Menge und MVP-Abnahme-Mapping divergieren

- `kategorie`: MEDIUM
- `quelle`: `spec/lastenheft.md` §1 *MVP-Schnitt* („Er umfasst die Anforderungen
  mit der Kennzeichnung **MVP: ja**") gegen die Mapping-Tabelle ebenda
- `pfad`: `spec/lastenheft.md:108-116` (sowie `:256`, `:701`, `:1038`)
- `befund`: [`LH-FA-RET-001`](../../spec/lastenheft.md) trägt `**MVP: ja.**` (Zeile 701), fehlt aber im
  MVP-Abnahme-Mapping. Umgekehrt verweisen die Mapping-Zeilen 4 und 6 auf
  [`LH-QA-REL-001`](../../spec/lastenheft.md) und [`LH-FA-CFG-006`](../../spec/lastenheft.md), die keinen MVP-Marker tragen
  ([`LH-QA-POR-003`](../../spec/lastenheft.md), Zeile 7, ist als einziger QA-Eintrag markiert und gemappt).
- `verifizierbar`: ja — Mengenvergleich {Anforderungen mit MVP-Marker} gegen
  {im Mapping referenzierte IDs} ist mechanisch prüfbar
- `klasse`: MVP-Kennzeichnung und Abnahme-Mapping inkonsistent

### F-2 — Traceability-Muster nennt nicht deklarierte ID-Formen

- `kategorie`: MEDIUM
- `quelle`: `harness/conventions.md` MR-000 (ID-Schema-Deklaration: `LH-*`,
  `SPEC-*`, `ARC-*`, `ADR-*`, … — weder `PH-*` noch `TST-*` deklariert);
  MR-001 (Begründung zitiert das `PH-*`-Muster als Umbenennungs-Grund)
- `pfad`: `spec/lastenheft.md:157-163`
- `befund`: Das Traceability-Beispiel führt `PH-FA-CAP-002` und `TST-CAP-002`
  als Spiegel-IDs. Das Pflichtenheft vergibt keine `PH-*`-IDs — es verfeinert
  über `LH-FA-CAP-002.a` und adressiert Festlegungen über `SPEC-*`. `TST-*`
  ist nirgends deklariert. Das Beispiel stammt unverändert aus v0.2
  (`LH-FUN-CAP-002` / `PH-FUN-CAP-002`, 577d510) und wurde beim ID-Wechsel
  auf `LH-FA-<BEREICH>-<NNN>` nicht mitgezogen.
- `verifizierbar`: ja — grep der Beispiel-Präfixe gegen MR-000 und die im
  Pflichtenheft vergebenen IDs
- `klasse`: Traceability-Beispiel nennt nicht deklarierte ID-Präfixe

### F-3 — Zurückstellungsliste in §5 widerspricht bindenden Anforderungen

- `kategorie`: MEDIUM
- `quelle`: Maintainability — Vertragskonsistenz innerhalb des Lastenhefts
- `pfad`: `spec/lastenheft.md:1192-1203` gegen `:746-748` ([`LH-FA-RET-004`](../../spec/lastenheft.md)),
  `:792-794` ([`LH-FA-SCH-001`](../../spec/lastenheft.md)), `:1123 ff.` ([`LH-QA-OPS-001`](../../spec/lastenheft.md) ff.)
- `befund`: Die „zukünftigen Erweiterungen — bewusst nicht im ersten Stand"
  listen „Erweiterte Retention ([`LH-FA-RET-004`](../../spec/lastenheft.md))", „Schema Evolution
  ([`LH-FA-SCH-001`](../../spec/lastenheft.md) ff.)" und „Produktionsreife Observability ([`LH-QA-OPS-001`](../../spec/lastenheft.md)
  ff.)" — während die referenzierten Anforderungen bindend formuliert sind
  („müssen"/„soll") und an der Anforderung selbst keine Zurückstellung tragen.
  Nur die Consumer-Zeile trägt die tragfähige Lesart („Fähigkeit ist
  angefordert; hier steht ihre Produktionsreife aus") — die anderen drei
  Zeilen unterscheiden Fähigkeit und Produktionsreife nicht.
- `verifizierbar`: nein — Urteil, kein Gate
- `klasse`: Zurückstellung im Sammelabschnitt ohne Marker an der Anforderung

### F-4 — Drei Lastenheft-Delegationen ans Pflichtenheft unerfüllt

- `kategorie`: MEDIUM
- `quelle`: [`LH-QA-PER-002`](../../spec/lastenheft.md), [`LH-QA-PER-004`](../../spec/lastenheft.md), [`LH-QA-POR-001`](../../spec/lastenheft.md)
- `pfad`: `spec/lastenheft.md:1079-1080`, `:1092-1093`, `:1158-1159` gegen
  `spec/pflichtenheft.md:248` und `:238-240`
- `befund`: Das Lastenheft legt dreimal fest, etwas werde „in
  `spec/pflichtenheft.md` festgelegt". Im Pflichtenheft stehen zu
  [`LH-QA-POR-001`](../../spec/lastenheft.md) nur „konkrete Versionen noch festzulegen" ([`SPEC-010`](../../spec/pflichtenheft.md)), zu
  [`LH-QA-PER-004`](../../spec/lastenheft.md) nur „Warn- und Fehlerschwellen … sind konfigurierbar" (keine
  festgelegten Werte), und zu [`LH-QA-PER-002`](../../spec/lastenheft.md) (Volumina) fehlt jede Stelle.
- `verifizierbar`: ja — Textvergleich der delegierten Stellen
- `klasse`: Lastenheft-Delegation ohne Entsprechung im Pflichtenheft

### F-5 — Akzeptanzkriterien verletzen die Drei-Pfad-/GWT-Form

- `kategorie`: MEDIUM
- `quelle`: Skill-Klassifikation MEDIUM „fehlende Negativtests bei neuem
  öffentlichem Vertrag"; Abschnittsregel §3 („Jede Anforderung trägt drei
  Pfade — Happy · Boundary · Negative"); Ziel-Form
  `v6.5.0` · `templates/spec/lastenheft.template.md` §3
- `pfad`: `spec/lastenheft.md:234`, `:269`, `:337`, `:355`, `:373`, `:419-422`,
  `:616`, `:876-882`, `:950`, `:1014-1015` (u. a.)
- `befund`: Mehrere Anforderungen tragen `—` statt eines Pfades (z. B.
  [`LH-FA-CFG-004`](../../spec/lastenheft.md) Negative, [`LH-FA-CFG-006`](../../spec/lastenheft.md) Negative, [`LH-FA-CAP-004`](../../spec/lastenheft.md)/005/006
  Negative, [`LH-FA-CON-001`](../../spec/lastenheft.md) Negative, [`LH-FA-SST-005`](../../spec/lastenheft.md) Boundary/Negative), obwohl
  die Abschnittsregel drei Pfade fordert — eine Regel, wann `—` zulässig ist,
  fehlt. Zudem sind dutzende Kriterien als „Given …, then …" ohne When-Teil
  formuliert (z. B. [`LH-FA-DAT-001`](../../spec/lastenheft.md) Happy/Boundary, [`LH-FA-ADM-001`](../../spec/lastenheft.md) Happy/Boundary,
  [`LH-FA-SST-001`](../../spec/lastenheft.md) Happy/Boundary).
- `verifizierbar`: ja — formal gegen die Kriterien-Form prüfbar
- `klasse`: GWT-Pfad unvollständig oder leer

### F-6 — Historie-Verweis „Pflichtenheft v0.2" löst nicht auf

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Zustandsfelder: Zustand + auflösbarer Anker)
- `pfad`: `spec/pflichtenheft.md:261`
- `befund`: Die Historie-Zeile referenziert „Pflichtenheft v0.2" — den
  zurückgezogenen Entwurf (c70ee1c), der *denselben Pfad*
  `spec/pflichtenheft.md` trug wie das heutige Dokument. Das Dokument führt
  selbst keine Versionsnummer, der Anker löst also in der Datei nicht auf;
  über `git log --follow spec/pflichtenheft.md` sieht ein Leser zwei
  Dokument-Lineages auf einem Pfad mit unterschiedlichen ID-Schemata
  (`PH-*`-Entwurf vs. `LH-*.<a>`/`SPEC-*`).
- `verifizierbar`: nein
- `klasse`: Historie-Anker löst nicht auf

### F-7 — SPEC-010 Vertrag-Zeiger zeigt auf unveröffentlichte Vorlage

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `spec/pflichtenheft.md:248`
- `befund`: Die Vertrag-Datei-Spalte verweist auf `spec/architecture.md`
  („Schnittstelle zur Quelle, [`LH-FA-SST-001`](../../spec/lastenheft.md)") — diese Datei ist aktuell die
  leere Ziel-Form-Vorlage (Platzhalter, `ARC-007` = `<…>`); der Zeiger löst
  inhaltlich nicht auf.
- `verifizierbar`: ja — Dateiinhalt gegen Zeiger prüfbar
- `klasse`: Vertrag-Zeiger auf leere Vorlage

### F-8 — Verfeinerungs-ID-Schema im Pflichtenheft generalisiert

- `kategorie`: INFO
- `quelle`: `v6.5.0` · `templates/spec/spezifikation.template.md` §1
  (`<PREFIX>-FA-<NN>.<Buchstabe>`) gegen `harness/conventions.md` MR-000
- `pfad`: `spec/pflichtenheft.md:16-19`, `:21`
- `befund`: Die Abschnittsregel erweitert das Vorlagen-Schema zu
  `LH-<STRATUM>-<BEREICH>-<NN>.<Buchstabe>` und nutzt es für `LH-QA-REL-001.a`.
  Das ist mit MR-000 (`LH-QA-*`) konsistent und sinnvoll — kein Widerspruch;
  die Stelle wird genannt, damit ein späteres Schema-Audit sie als entschieden
  lesen kann, nicht als Abweichung.
- `verifizierbar`: nein
- `klasse`: ID-Schema-Generalisierung im Technik-Stratum

### F-9 — LH-QA-OPS-003 „Consumer-Positionen" nur als Lag-Metrik abgedeckt

- `kategorie`: INFO
- `quelle`: [`LH-QA-OPS-003`](../../spec/lastenheft.md)
- `pfad`: `spec/pflichtenheft.md:233`
- `befund`: Die Metrik-Tabelle bildet Consumer-Positionen ausschließlich als
  `cdc_consumer_lag` („Rückstand je Consumer") ab; eine direkte
  Positions-Metrik fehlt. Die Position ist aus Rückstand und Bestand
  ableitbar — ob das die Anforderung deckt, ist eine Lesart-Frage.
- `verifizierbar`: nein
- `klasse`: Metrik-Deckung indirekt

### F-10 — Reviewer-Skill ungeschärft (Vorlagen-Platzhalter)

- `kategorie`: INFO
- `quelle`: `v6.5.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form:
  Reviewer-Skill (HIGH-Liste braucht mindestens zwei repo-spezifische Regeln)
- `pfad`: `.harness/skills/reviewer.md:45-47`
- `befund`: Die Skill-Datei trägt noch die Platzhalter „repo-spezifisch #1/#2"
  und den Template-Kopf; sie ist nach eigener Pflicht-Notiz „noch nicht
  scharf genug". Dieser Lauf ist auf die Baseline-Klassifikation
  ausgewichen — kein Befund gegen den geprüften Diff, sondern gegen das
  Werkzeug; Zuständigkeit: Architect/Implementer des Skills.
- `verifizierbar`: nein
- `klasse`: Skill-Datei mit Vorlagen-Platzhaltern

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Kategorie** über beide Dokumente — kein
  ADR-Verstoß (keine Accepted-ADR im Verzeichnis), keine Gate-Suppression,
  kein Sicherheits- oder Korrektheits-Finding, keine Norm nur im
  Template-Kommentar, kein Chronik-tragendes Zustandsfeld
- geprüft, ohne Befund: **Lastenheft-Vorlagen-Form, Struktur** — alle sieben
  Abschnitte vorhanden (§1–§7), Kopf mit Version/Status/Autor/Datum, Glossar
  und Historie-Regeln inkl. Decken-Regel übernommen
- geprüft, ohne Befund: **ID-Schema Lastenheft** — `LH-FA-<BEREICH>-<NNN>` /
  `LH-QA-<BEREICH>-<NN>` durchgängig, keine Lücke, keine Doppelvergabe,
  keine Wieder-verwendung
- geprüft, ohne Befund: **SPEC-Zählung Pflichtenheft** — [`SPEC-001`](../../spec/pflichtenheft.md)…011
  fortlaufend je Datei, jede Struktur/Festlegung in §2–§6 trägt eine ID
- geprüft, ohne Befund: **Decken-Regel beider Historien** — kein ADR-, Slice-,
  Carveout- oder Welle-Verweis in irgendeiner Spalte
- geprüft, ohne Befund: **Verfeinerungen Pflichtenheft („präzisieren ja,
  erweitern nie")** — alle sechs `LH-*.<a>`-Elemente verfeinern ihre genannte
  Lastenheft-ID; alle zitierten LH-IDs existieren; kein Erweiterungs-Befund
- geprüft, ohne Befund: **Umbenennungs-Nachzug** — keine Rest-Verweise auf
  `spec/spezifikation.md` außerhalb des vendored Baseline-Bestands
  (case-insensitive grep über `spec/`, `AGENTS.md`, `README.md`,
  `harness/README.md`, `harness/conventions.md`); verbleibende
  „Spezifikation"-Nennungen sind Baseline-Modul-Titel-Zitate bzw. der
  „vormals"-Kontext in MR-001
- geprüft, ohne Befund: **Inhaltliche GWT-Widersprüche zwischen den Pfaden** —
  Intervall-Semantik REA-001 `[p1, p2)` und REA-002 „ab `p` ausschließlich"
  konsistent; CAP-006/CAP-007 konsistent; die Verkettungen
  CFG-005 ↔ DAT-005 ↔ SCH-003/SCH-004 und RET-004 ↔ CON-006 widersprechen
  einander nicht

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** MVP-Kennzeichnung und Abnahme-Mapping
inkonsistent · Traceability-Beispiel nennt nicht deklarierte ID-Präfixe ·
Zurückstellung im Sammelabschnitt ohne Marker an der Anforderung ·
Lastenheft-Delegation ohne Entsprechung im Pflichtenheft · GWT-Pfad
unvollständig oder leer · Historie-Anker löst nicht auf · Vertrag-Zeiger auf
leere Vorlage · ID-Schema-Generalisierung im Technik-Stratum · Metrik-Deckung
indirekt · Skill-Datei mit Vorlagen-Platzhaltern

## Verdikt

**Merge-blockierend:** nein für den Draft-Betrieb (beide Dokumente sind
Status Draft; nichts blockiert das Weiterarbeiten) — **aber
`Accepted`-blockierend**: F-1 bis F-5 betreffen Vertrags-Kernstellen
(MVP-Schnitt, Traceability, Zurückstellungen, Delegationen,
Akzeptanzkriterien-Form) und müssen vor jeder `Accepted`-Setzung des
Lastenhefts geklärt sein; ab `Accepted` wird jede dieser Stellen zur
Vertragsänderung.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan bei
Plan-Defekt: F-2 stammt als Beispiel unverändert aus v0.2 — Plan-Korrektur,
kein Code-Fix). Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7
und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD-/Spec-
Konformität prüft der Verifier separat (`v6.5.0` ·
`regelwerk/modul-11-*.md`).