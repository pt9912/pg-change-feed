# Architect-Verdikt: Gate für kennungsfreie Nutzerdokumentation, Skill, Umstellung der Prüfpunkte

**Rolle:** Architect (Modul 8)
**Anlass:** Startbedingung des Slice
[`handbuch-public-doc-check-gate-und-skill`](../plan/planning/in-progress/slice-handbuch-public-doc-check-gate-und-skill.md)
(Liefer-Punkt 0; Planner-Fragen 1 bis 6)
**Datum:** 2026-10-02
**Bezug:** [`LH-QA-OPS-001`](../../spec/lastenheft.md),
[`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(Ergebnis dieses Verdikts), [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
(Muster), [`ADR-0087`](../plan/adr/0087-beispiel-clients-csharp-kotlin.md),
[`ADR-0088`](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
`AGENTS.md` §3.1, §3.5, §3.6, §3.12

---

## Verdikt in einem Satz

Die Entscheidung ist ein Gate `handbuch-public-doc-check` in `GATE_CHECKS` mit einer
benannten Liste plus Vollständigkeitsprüfung, ein kleiner Skill statt einer Rolle und
eine Ergänzung (nicht Abschwächung) der Historie-Pflicht; die Accepted-ADRs bleiben
unberührt, die Schärfung trägt die neue
[`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md).
Der Plan hat nichts falsch behauptet; er braucht Nachzüge (unten).

## Entscheidungen Q1 bis Q6

| Frage | Entscheidung | Kern der Begründung |
|---|---|---|
| Q1 Reichweite | **Hybrid:** geprüft `benutzerhandbuch.md`, `benutzerhandbuch-standard.md`, `version.md`; ausgenommen `releasing.md` und die vier `*-abdeckung.md`; jede andere `*.md` unter `docs/user/` und jede genannte, fehlende Datei färbt den Lauf rot (Exit 2). `version.md` gehört dazu | Fehlerbild der reinen Liste: neue Nutzerdatei **lautlos unbewacht**. Fehlerbild von „Verzeichnis minus Ausnahmeliste“: neue Erzeugnis-Datei **falsch rot**, `releasing.md` als dauerhafte Schuld in der Liste. Die Klassifikationsprüfung (etwa zehn Zeilen) macht beide Fehler laut und erzwingt die Entscheidung bei jeder neuen Datei |
| Q2 `releasing.md` | **dauerhaft außerhalb, kein Folge-Slice** | Maintainer-Doku mit bewussten Entscheidungs-Links; eine Bereinigung wäre Inhaltsänderung ohne Betreiber-Nutzen. Zahl 26 gemessen am Stand `cfb68a31` |
| Q3 Ausnahmen | **keine**, auch nicht im Fence; Opt-out-Marker nur per eigener ADR | Ein Beispiel mit Kennung ist selbst der Fehler; gemessen 0 Treffer, Tag-Muster wie `sdk-csharp-v*` treffen nicht. Folge: der Standard-Satz zur Kennungsfreiheit **nennt keine Beispiel-Kennung**, sonst färbt er das Gate rot |
| Q4 Muster | ERE (kein `-P`, nicht in der Host-Klasse): `\b(LH-(FA\|QA\|RB)-\|(ADR\|SPEC\|ARC)-[0-9]\|BEO-[A-Za-z]\|(MR\|CO)-[0-9]{3})\|(^\|[^[:alnum:]_-])(slice\|welle)-[a-z0-9]`; Links `docs/(reviews\|plan)/\|\.\./(reviews\|plan)/` | Wortrand für `slice-`/`welle-` über Zeilenanfang oder Nicht-Wortzeichen-und-nicht-Bindestrich (ein bloßes `\b` trifft `byte-slice-x`); `MR-`/`CO-` nur mit drei Ziffern |
| Q5 Schärfung | `Schärft:` [ADR-0087](../plan/adr/0087-beispiel-clients-csharp-kotlin.md)/-0088/-0090 im Wortlaut der Historie-Zeile („ohne Kennung, Betreibersicht“); keine Supersession | Die drei fordern die Zeile, nicht Kennungen darin; die Forderung bleibt wahr. Pflichtenheft-Zeile 1075 bleibt wahr |
| Q6 Namen und Orte | Gate/Skript/Sensor `handbuch-public-doc-check`; Tabellentest `run-handbuch-public-doc-check-tests.sh` / `make test-handbuch-public-doc-check` (Werkzeug); Heimat `harness/mk/doc-gate.mk`; Skill `.harness/skills/nutzerdoku-schreiben.md` | `sdk.mk` ist die Heimat der SDK-Ziele, `doc-gate.mk` die der Doku-Gates. Planner-Arbeitsname bleibt: `nutzerdoku-…` wäre ein Missname, weil `releasing.md` ebenfalls Nutzerdoku ist |

## Messungen (Ursprung nach `AGENTS.md` §3.12)

- **Gemessen** am Stand `cfb68a31`, 2026-10-02, `git grep -cE "<Muster>" -- <Datei>`:
  Handbuch 0, Standard 0, `version.md` 0, `releasing.md` 26, `bench-abdeckung.md` 4,
  `ci-matrix-abdeckung.md` 3, `e2e-abdeckung.md` 74, `sdk-e2e-abdeckung.md` 21. Das
  geschärfte Muster und das Planner-Muster liefern dieselben acht Zahlen.
- **Gemessen**, Link-Muster: `releasing.md` 7, `bench-abdeckung.md` 1,
  `ci-matrix-abdeckung.md` 1, alle anderen 0; das Handbuch verlinkt nur
  `../../examples/README.md`.
- **Gemessen**, Wortrand, `grep -nE` an `printf`-Zeilen: Treffer bei `slice-foo`,
  `welle-x`, `(slice-1)`, `a ADR-0134`, `CO-123`, `MR-001`, `LH-FA-SST-009`, `LH-FA-X`,
  `ARC-001`, `BEO-PGC/x` (10); kein Treffer bei `byte-slice-x`, `my-slice-1`,
  `CO-2 Emissionen`, `ALCO-123`, `x_CO-123`, `ZMR-100`, `sdk-csharp-v1`, `NR-0001`,
  `Slice-1`, `ADR-` (10). Beim ersten Versuch mit `grep -P` und Lookbehind ohne `-`
  traf `byte-slice-x` noch: korrigiert, die ERE-Form trägt.
- **Hergeleitet, nicht erprobt:** Skript-Mutation (Kopie mit eingefügter Kennung),
  Klassifikationsprüfung, Gate-Verdrahtung. Der Implementer fährt sie; das ist eine
  Erwartung (ADR Fitness Function).
- **Gemessen:** `.harness/skills/` trägt genau `reviewer.md` und
  `closure-note-reviewer.md`; `implement-slice.md` hat 373 Zeilen.

## Plan-Nachzugsliste (der Planner; der Architect ändert den Plan nicht)

1. §1/DoD Liefer-Punkt 0: ADR-Nummer **0143** eintragen, Status `Accepted` (Datei und
   Index-Zeile liegen vor); Startbedingung §4 ist damit erfüllt.
2. §1 Vorschlag 1 und DoD Liefer-Punkt 1: Reichweite ist der **Hybrid** (Listen plus
   Vollständigkeitsprüfung, Exit 2 für unklassifiziert/fehlend); der Tabellentest
   deckt zusätzlich: unklassifizierte `.md`, fehlende genannte Datei, `version.md`,
   `releasing.md` und die vier Erzeugnisse als ausgenommen (Exit 0), `byte-slice`,
   `CO-2 Emissionen`, `Slice-1` (kein Treffer).
3. Suchlauf §3: das Muster der Zeilen 1 bis 8 und 9 durch `P`/`L` aus der ADR ersetzen
   (Zahlen unverändert: 0/0/0/26/4/3/74/21); Soll für die ADR-Zeile (`Änderungshistorie`
   in `docs/plan/adr`) ist 8 plus Treffer der neuen ADR (gemessen 7 Zeilen, Summe 15 mit `git add -N`;
   der Implementer misst am Diff-Stand nach, vorher
   `git add`).
4. §1 Vorschlag 3: „keine Ausnahme“ ist entschieden; Hinweis an den Implementer, im
   Standard keine Beispiel-Kennung zu nennen.
5. §1 Vorschlag 7 / Liefer-Punkt 2: Skill-Ort bestätigt; der Skill nennt `P`/`L` nicht
   wörtlich (er liegt außerhalb des Gate-Gegenstands, aber das Muster gehört in den
   Sensor-Vertrag).
6. §6 Risiko „Gate-Fläche“: Ausgang jetzt *entschieden* (Hybrid), am Closure
   *eingetreten/entfallen* gegen die Probe; Risiko „Falsch-positive“: Wortrand-Regel
   steht, Restgrenze Groß-/Kleinschreibung im Sensor-Vertrag führen.
7. §2 „Nur diese Pfade“: die neue ADR und ihr Index-Eintrag gehören zum Diff des
   Architect-Zugs, nicht zu dem des Implementers.
8. DoD Liefer-Punkt 1: zu `harness/README.md` auch die Zeile des Tabellentests
   (Werkzeug) und die `make gates`-Aufzählung.

## Empfehlungen an den Auftraggeber

- Keine Rückfrage nötig. Wenn `releasing.md` irgendwann bereinigt werden soll, ist das
  ein eigener Inhalts-Slice plus Verschiebung in die Liste „geprüft“ ([ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
  Festlegung 9: der Weg von „ausgenommen“ nach „geprüft“ ist eine Verschärfung und
  braucht keine neue ADR; der umgekehrte Weg wäre eine Milderung).
- Die Vollständigkeitsprüfung macht jede künftige `docs/user/*.md` zu einem
  bewusst klassifizierten Eintrag; das ist gewollter, kleiner Aufwand.
