# Slice slice-078: host-lokale absolute Pfade entfernen, Modul aktivieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine beobachtbare Closure-Bedingung, die von
der DoD dieses Slice verschieden wäre: die Entfernung ist an zwei Greps
messbar und das Gate danach am aktivierten Modul (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
(aktiviert das Modul ohne Ausnahme und trägt den Entwurf der Hard Rule),
[`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(die Zitat-Korrektur-Klasse, die Accepted ADRs korrigierbar macht),
[`ADR-0074`](../../adr/0074-zitationsform-schwester-repo-hausform.md)
(Hausform der Zitation; ersetzt `ADR-0072` in genau einer Klausel);
`AGENTS.md` §3.6 (Verschärfung) und §3.7 (Kommentar/Zustandsfeld).

**Berührte Spec-Stellen:** — (kein `SPEC-NNN` und kein `ARC-NNN`: eine
Gate-Aktivierung, eine Hard Rule und Zitat-Korrekturen, ohne Spec-Stratum).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Kein host-lokaler absoluter Pfad steht mehr im Repo. Gemessen sind
**42 Vorkommen in 15 Dateien** (§3); sie werden **entfernt**, nicht ausgenommen
— und das `hostpaths`-Modul des Doku-Gates wird aktiviert, damit die Regel ab
dann mechanisch greift.

Zwei Formen tragen die Korrektur (beide am Bestand belegt, `ADR-0074`):
**Zitatstellen** verlieren den Host-Präfix und nennen **Anzahl und
repo-relativen Dateinamen** — der Modul lässt Fenced-Blöcke frei, die Regel
dieses Repos nicht: „keine absoluten Pfade" heißt „sie stehen nicht mehr da",
nicht „das Gate meldet sie nicht". **Verweise auf ein Schwester-Repo** tragen
die Hausform `` `d-check`s `tools/coverage-gate.sh` `` — Repo als
Inline-Code-Wort, Pfad relativ darin.

**Abnahme (Messwert, keine Zusage):** der Modul-Aufruf über das Repo ergibt
**0** Befunde, `make gates` ist mit aktivem Modul **ohne** Ausschlussblock grün,
und die zwei Stellen, die der Modul nicht liest (`Makefile`,
`tools/coverage-gate.sh`), sind im Diff nachweislich korrigiert.

**Die Präfix-Liste wird hier bewusst nicht wiederholt** — sie steht an genau
einer Stelle, in den Vorgaben des Moduls (`hostpaths.prefixes`), und die Abnahme
zieht sie von dort. Eine Regel gegen host-lokale Pfade, die ihre Präfixe
aufzählt, verstößt gegen sich selbst; „eine Quelle, alle übrigen zeigen darauf"
ist dieselbe Disziplin, die dieser Slice für den beweglichen Wert der
Coverage-Rampe schon trägt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Rangordnung zwischen Hard Rules.** Der Auftraggeber hat gefragt; die
  Antwort ist nein, und sie steht im Verdikt (`ADR-0072`): der Konflikt ist
  **Umfang, nicht Rang** — die klassenweise Eingrenzung von `AGENTS.md` §3.5
  löst ihn, eine Rangfolge wäre Seniorität in Tabellenform (Baseline-Regelwerk
  `modul-08-agentenrollen.md` §Konflikt-Pfad verbietet Seniorität als Argument).
- **Eine Ausnahme im Scope des Moduls** (`scope.ignore`, `exempt-paths`,
  Marker). Der Ausgang ist 42 → 0; jede Ausnahme verfehlt ihn. Das Modul kennt
  ohnehin keinen Marker und kein `exempt-paths` — geprüft, nicht angenommen.
- **Eine Folge-ADR für `ADR-0051`/`0054`.** Sie würde die Pfade **nicht**
  entfernen: der alte Text bleibt im Scan stehen, die neue Fassung träte daneben.
  Die Korrektur läuft über die Zitat-Klasse (`ADR-0073`).
- **Änderungen an Entscheidungen** der berührten Accepted ADRs
  (§Entscheidung, §Konsequenzen, §Verglichene Alternativen, §Status,
  `Supersedes`-Kette) — die Klasse deckt ausschließlich das **Zitatgerüst**.
- **Ein Wächter für die Stellen, die der Modul nicht liest** (`Makefile`,
  Skriptkommentare unter `tools/**`). Sie werden **hier** korrigiert; einen
  Sensor dafür zu bauen wäre ein anderer Vorgang (der Modul liest `.md`), und
  die Lücke wird in `harness/sensors/docs-check.md` **benannt** statt
  verschwiegen.
- **Produktionscode unter `internal/**`, Schemamigrationen unter
  `tools/schema/**`, Docker-Stage, Gate-Skript-Logik** — Schicht-Abgrenzung:
  dieser Slice ändert Doku, zwei Kommentarzeilen und eine Config-Zeile.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — kein Host-Pfad mehr im Baum.** Die 42 Vorkommen in 15
Dateien (§3) sind entfernt; Zitatstellen nennen Anzahl und Datei, Verweise auf
Schwester-Repos die Hausform.

- [x] **Abnahme: 0** — der Modul-Aufruf über das Repo (Präfix-Liste aus dem
      Modul gezogen, hier nicht wiederholt) als eigener, ungepipeter Schritt
      (`AGENTS.md` §3.9), plus die zwei korrigierten Stellen außerhalb seiner
      Reichweite.
- [x] Beide Formen sind belegt: kein Fence-Zitat trägt mehr einen Host-Präfix,
      und jeder Schwester-Verweis nennt das Repo als Inline-Code-Wort.
- [x] Die zwei Stellen, die der Modul **nicht** liest (`Makefile`,
      `tools/coverage-gate.sh`), sind mitgezogen.

**Liefer-Punkt 2 — das Modul greift.** `hostpaths` steht in `modules` **ohne**
Ausschlussblock; die Sensors-Doku trägt den erweiterten Vertrag.

- [x] `.d-check.yml`: `modules: […, hostpaths]`, **kein** `scope`, **kein**
      `ignore`, **kein** `exempt-paths`; der Modul-Aufruf liefert **0**.
- [x] `harness/README.md` §Sensors-Zeile `make docs-check` nennt `hostpaths` —
      sie listet die Module auf und behauptete sonst einen anderen Vertrag als
      das Gate hat.
- [x] `harness/sensors/docs-check.md` benennt die Grenzen: Fenced-Blöcke frei,
      Windows-Muster fest, **relative** Pfade ungeprüft, `Makefile`, `tools/**`
      und `harness/mk/**` ungescannt.

**Liefer-Punkt 3 — die Regel ist verkörpert.** `AGENTS.md` §3.11 trägt die Hard
Rule nach dem Entwurf in `ADR-0072` — und `AGENTS.md` §3.5 trägt die Ausnahme
der Zitat-Klasse, damit die Hard Rule nicht gegen die Korrektur steht, die
dieselbe Entscheidung sanktioniert.

- [x] §3.11 steht mit Aussage, Falsch/Richtig an **realen** Beispielen,
      Begründung, Grenzen und Rang-Zeiger auf das Modul.
- [x] §3.5 nennt die Zitat-Korrektur als Ausnahme zu „Korrekturen entstehen als
      neue ADR mit `Supersedes`"
      ([`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)),
      und der Kopf-Satz des ADR-Index führt dieselbe Ausnahme. **Nachtrag zum
      Plan:** beide Stellen standen nicht in der Fassung, mit der implementiert
      wurde — sie sind die Folgepflicht aus `ADR-0073`, und ihr Fehlen war ein
      Plan-Versäumnis, kein Auftragsverzicht.

- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Inventur (gemessen, ohne die Fundstellen zu wiederholen):** 42 Vorkommen in
15 Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0054-…md` (10) · `docs/plan/adr/0051-…md` (1) | update (Zitatgerüst) | Accepted → Korrektur ausschließlich über die Klasse aus `ADR-0073`; Entscheidungstext unberührt |
| `docs/plan/adr/0072-…md` (4) | update (Zitatgerüst) | dieselbe Klasse, auf die ADR angewandt, die sie definiert: ihre Inventur zitiert die Pfade |
| `docs/reviews/architect-verdict-hostpaths-…md` (5) | update | die Inventur des Verdikts; Zitat wird Anzahl + Datei |
| `docs/plan/planning/done/{slice-049, slice-050, welle-14, welle-14-results}` (11) | update | Zeitdokumente: ihre Einfrierung ist zeitlich, nicht status-basiert (`ADR-0073`) |
| `docs/reviews/{review-slice-036, 039, 049, 073}` (7) | update | Lauf-Belege; dieselbe Klasse, Commit ist der Beleg |
| `harness/sensors/coverage-gate.md` (2) | update | lebende Doku → Hausform |
| `Makefile` (1) · `tools/coverage-gate.sh` (1) | update | vom Modul **nicht** gelesen; trotzdem korrigiert, damit die Zusage dieses Repos trägt |
| `.d-check.yml` | update | `hostpaths` in `modules`, **ohne** Ausschlussblock |
| `AGENTS.md` | update | Hard Rule §3.11 nach dem Entwurf in `ADR-0072`; **und** §3.5 trägt die Ausnahme der Zitat-Klasse (`ADR-0073`) |
| `docs/plan/adr/README.md` | update | der Kopf-Satz nennt dieselbe Ausnahme — sonst lehrt der Index die abgelöste Fassung |
| `harness/README.md` | update | §Sensors-Zeile `make docs-check` nennt das Modul |
| `harness/sensors/docs-check.md` | update | Grenzen des Moduls, Ist-Zustand der Modulliste |

Kein weiterer Ort führt einen Host-Pfad; `spec/**`, `internal/**`, `cmd/**`,
`tools/schema/**` und die übrigen `docs/`-Bäume sind befundfrei (gemessen).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): die drei ADRs (`0072`, `0073`, `0074`) liegen
`Accepted` vor, `Verantwortlich:` gesetzt, WIP-Limit frei — `in-progress/` trägt
keinen anderen Slice (`slice-077` wurde für diesen Vorgang zurückgestellt).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt sich, dass die
  Zitat-Korrektur mehr tragen müsste als das **Zitatgerüst** — dass an einer
  Stelle der Entscheidungstext selbst geändert werden müsste, um den Pfad
  loszuwerden —, ist die Klasse zu eng und der Schnitt gehört zurück zur
  Zerlegung; dann ist es je Dokument eine inhaltliche Korrektur und braucht eine
  Folge-ADR.
- `in-progress` → `open` (blockiert — Carveout?): erweist sich ein Vorkommen
  ohne eine Entscheidung als nicht entfernbar, die dieser Slice nicht treffen
  darf, ist das ein Blocker mit Carveout-Entscheidung — **nicht** eine stille
  Ausnahme im Modul-Scope.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** die Abnahme-Messung ergibt **0** (Modul-Aufruf über das
Repo, plus die zwei Stellen außerhalb seiner Reichweite) **und** `make gates`
grün mit aktivem `hostpaths` **ohne** Ausschlussblock **und** die Closure-Notiz
geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Zitat-Korrektur an Accepted ADRs könnte als Einfallstor gelesen
  werden** — sie ist die einzige Stelle, an der dieser Slice ein immutables
  Dokument anfasst. Die Grenze steht in `ADR-0073`: nur das Zitatgerüst, nie
  Entscheidung, Konsequenzen, Alternativen, Status oder `Supersedes`-Kette; der
  Beleg ist die Commit-Kennung **plus** eine §Geschichte-Zeile je betroffener
  ADR. — **Ausgang:** <bei Closure>
- **Die Form „Anzahl + Dateiname" verliert die konkrete Fundstelle.** Beabsichtigt
  — der Präfix ist Maschinen-Layout, keine Information — und re-derivierbar: das
  aktivierte Modul zeigt die Fundstellen jederzeit wieder. — **Ausgang:**
  <bei Closure>
- **Beim Korrigieren neue erzeugen:** die Beschreibung der Funde zitiert die
  Funde. Die Regel „Zitat = Anzahl + Datei" verhindert es, der Abnahme-Grep
  fängt es. — **Ausgang:** <bei Closure>
- **Zwei Stellen liegen außerhalb der Modul-Reichweite** (`Makefile`,
  `tools/coverage-gate.sh`): sie werden korrigiert, aber **kein** Sensor hält
  sie künftig — die Lücke wird in `harness/sensors/docs-check.md` benannt. —
  **Ausgang:** <bei Closure>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md`
§Modus-Deklaration pro Sub-Area) — sie deckt `docs/`, `harness/`, `AGENTS.md`,
`.d-check.yml`, `Makefile` und `tools/**` in **einem** Kürzel. Eine feinere
Sub-Area ist nicht zu bilden: „Doku-Gate" ist keine deklarierte Sub-Area, und
die Registereinträge zu Gate- und Prozessfragen führen dasselbe Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; drei Treffer:

- `BEO-PGC/slice-pfad-als-link-in-berichten` (als Regel **verkörpert**): dieselbe
  Familie — eine Textform, die vor dem nächsten `git mv` bricht —, aber eine
  **andere** Klasse: dort der Link auf ein wanderndes Artefakt, hier ein
  Host-Pfad. Der Slice berührt die verkörperte Regel nicht; **kein** neuer Beleg.
- `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (2×, offen): **kein** Treffer —
  dieser Slice führt kein Erzeugnis ein, er korrigiert Text.
- `BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (1×,
  `verkörpert`, `ADR-0071`): **kein** Treffer — anderer Gegenstand (Zielzahl vs.
  Messumfang).

**Vorschlag für die Closure:** ein neuer Eintrag
`BEO-PGC/gate-modul-abgeschaltet-trotz-regel` (1×, `verkörpert` — `ADR-0072`).
Die Klasse tritt hier zum ersten Mal auf: die Regel bestand als Absicht, ihr
Sensor war **abgeschaltet**, und der Bestand sammelte 42 Verstöße, ohne dass
etwas sie hätte melden können. Kein Registereintrag erreicht mit diesem Slice
die Schwelle 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`) — kein Modus-Begründungsblock.
Die vier Pflichtkriterien tragen hier dennoch: **Konventionen-Dichte** hoch
(`.d-check.yml`-Module, `harness/sensors/`, das Hard-Rule-System),
**Phase-Reife** hoch (das Coverage-/Doku-Gate läuft seit `slice-049`, 78 Slices
sind geschlossen), **Evidenz-/Diskrepanz-Risiko** niedrig — die Aussage wird
nicht gegen einen Code-Bestand inventarisiert, sondern **aus** ihm gemessen, und
die Messung ist jederzeit wiederholbar —, **Reconciliation-Aufwand** null, kein
Brownfield-Bestand.
