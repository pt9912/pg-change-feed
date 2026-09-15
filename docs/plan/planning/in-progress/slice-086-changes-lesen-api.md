# Slice slice-086: Changes über die HTTP-API lesen — `GET /changes`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (ein Endpunkt
und sein Vertrag), kein repo-weites *Mehr*.

**Bezug:** [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)
(der Endpunkt, der Inbound Port, die Filterachse, die Rechtsklasse — er führt
`ADR-0057`s **Option C** aus) · [`ADR-0057`](../../adr/0057-http-grpc-api.md)
(die API und ihre Token-Klassen; in zwei Klauseln teil-abgelöst) ·
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (die Anforderung, deren
**Beschreibung** das Lesen über einen Netzwerkweg verlangt) ·
[`LH-FA-REA-001`](../../../../spec/lastenheft.md) ff. (der Lesezugriff, den dieser
Endpunkt über den Netzwerkweg zugänglich macht) ·
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md) (der
View-Direktzugriff des SQL-Kanals — **unberührt**) ·
[`ADR-0042`](../../adr/0042-transport-typen-am-port.md)
(Transport-Typen am Port).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md) (die
Endpunkt-Familie, deren Satz das Changes-Lesen ausgrenzt) · **neu `SPEC-022`**
(die Ausgestaltung dieses Endpunkts).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die HTTP-API bekommt das **Changes-Lesen**: `GET /changes` in der
`reader`-Rechtsklasse, als **dünne lesende Fassade über den bestehenden
`ChangeStorePort`**. Damit wird `LH-FA-SST-006`s **Beschreibung** eingelöst (sie
verlangt das Lesen über einen Netzwerkweg) und `ADR-0057`s **Option C**
ausgeführt, die dort seit ihrer Entscheidung als *vertagt* geführt wurde.

**Kein zweiter Lesepfad.** Derselbe Port, dieselbe Delegation, dieselbe
Sortierung — der Endpunkt entscheidet nichts, was der View-Direktzugriff nicht
schon entscheidet; `ADR-0046` bleibt unberührt. Der Lesepfad zieht den
`cdc.source_table`-Join nach, damit die Antwort die Tabellen-Identität in
**Klartext** trägt — die opake `SourceTableID` ist nach `ADR-0056` Festlegung 1
kein Draht-Bezeichner, und ohne Namen wäre die Antwort nicht deckungsgleich mit
`cdc.changes`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Diagnose/Health über die API.** `ADR-0057` hat beide gemeinsam vertagt, aber
  sie brauchen **je** eine eigene Port-Design-Entscheidung; `ADR-0081` entscheidet
  nur die Lese-Hälfte. Eigene Folge-ADR.
- **Ein Default-Limit oder eine harte Obergrenze.** `ADR-0081` lässt beides
  bewusst weg — eine still eingebaute Grenze machte die Antwort nicht mehr
  deckungsgleich mit dem View-Direktzugriff. Als Re-Evaluierungs-Trigger benannt.
- **Paginierungs-Kursor und weitere Filter** jenseits Quelle/Schema/Tabelle/
  Bereich/Limit — Trigger benannt, eigener Beschluss.
- **Ein Schema-/DDL-Eingriff** — dieser Slice **liest**, er schreibt nicht.
- **Die Wegwerf-Clients und `examples/**`.** Dieser Slice **entblockt**
  `slice-083`, zieht es aber nicht mit: dessen Arbeitsstand liegt auf seinem
  eigenen Branch und ist laut `ADR-0081` nach diesem Bau **unverändert
  lauffähig**.

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

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — der Lesepfad trägt die Tabellen-Identität.**

- [x] Die Filterachse des Leseports ist **eine Form**, nicht zwei: `Schema`/
      `Table` als **Klartext** ersetzt `Table *model.SourceTableID`.
- [x] `SelectChanges` zieht den `cdc.source_table`-Join nach; `mapper.ToChange`
      setzt `Schema`/`Table`, sodass die Antwort die Identität trägt.
- [x] Port-/Mapper-Tests decken die neue Achse (einzeln, kombiniert, leer).

**Liefer-Punkt 2 — der Use Case und sein Port.**

- [x] `ReadChangesUseCase` als Inbound Port, Transport-Typen **am Port**
      (`ADR-0042`); der Inbound Port importiert **nicht** outbound.
- [x] Die Filter: `source` **Pflicht**, `schema`/`table` optional und
      **unabhängig**, `from` inklusiv / `to` exklusiv, `limit` optional —
      **kein** Default-Limit.
- [x] Netzlose Tests am Use Case: Filterwirkung, Bereichs-Grenzen, Leerfall.

**Liefer-Punkt 3 — der Endpunkt und sein Vertrag.**

- [x] `GET /changes` in der `reader`-Rechtsklasse, kollisionsfrei neben
      `GET /changes/stream`; Fehler-Mapping `ErrNonPositiveLimit`/
      `ErrRangeInverted` → **400**, **unbekannter Query-Parameter → 400**
      (strenger als die neun Bestandsendpunkte — ein unbekannter *Filter*
      änderte den Ergebnisstand still), **kein Treffer → 200 mit
      `{"changes": []}`**, nie 404.
- [x] Der Vertrag ist fortgeschrieben: `SPEC-018`s Satz nachgezogen, **`SPEC-022`**
      angelegt (Endpunkt, Parameter, JSON-Schema in der Form der neun
      bestehenden), und das Handbuch nennt den Endpunkt in der
      Fähigkeits-Tabelle samt Querverweis aus seinem Lese-Abschnitt.
- [x] `make gates` grün (Exit direkt, ungepiped).

- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-086.md` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Fixrunde gelaufen: 0 HIGH; F-2…F-5 geschlossen; F-1 mit `slice-087`
      adressiert.
- [x] Verifikation durchgeführt, Report unter
      `docs/reviews/verify-slice-086.md` liegt vor (Modul 11, frischer Kontext) —
      **closure-fähig**; LP1-K1…LP3-K3 einzeln bestätigt, `tools/schema/**`
      unberührt, §3.10 **geprüft und nicht ausgelöst** (kein
      `.github/workflows/**` im Diff). Eigene Gegenproben des Verifiers: beide
      Schwellen auf 99 → beide Gates rot; zwei eigene Mutationen (Projektions-
      Ordnung, `source`-Pflicht) rot gesehen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register (`../reconciliation.md`) — **entfällt**: dieses
      Repo führt die Datei nicht (Greenfield-Bootstrap, kein Inventur-Fund).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — ein Beleg
      ergänzt (`generierte-artefakte-ohne-sync-sensor`, damit **4×**) und ein
      Verzeichnis **neu angelegt** (`negativtest-ohne-bindung-an-seine-eingabe`,
      **1×**), **kein Zähler gesetzt**.
- [x] Jedes Risiko aus §6 trägt einen Ausgang — alle drei *entfallen,
      gestrichen mit Begründung*.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) — dieses Repo führt
      **Wellen-Betrieb**; die Prüfung fällt der `welle-20`-Closure zu, hier nicht
      geprüft und hier nicht fällig.

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| der Leseport (`internal/application/port/outbound/**`) · `internal/adapters/driven/postgresstorage/**` (`SelectChanges`, Mapper) | refactor | die Filterachse wird **eine** Form: `Schema`/`Table` als **Klartext** statt `Table *model.SourceTableID`; der Join auf `cdc.source_table` zieht nach, damit die Antwort die Identität trägt |
| `internal/application/port/inbound/**` · `internal/application/usecase/readchanges/**` | neu | `ReadChangesUseCase` samt Transport-Typen **am Port** (`ADR-0042`) |
| `internal/adapters/driving/http/**` | update | `GET /changes`, Query-Parameter, Fehler-Mapping, Adapter-Tests |
| `spec/pflichtenheft.md` · `docs/user/benutzerhandbuch.md` | update | der Vertrag: `SPEC-018`s Satz, **neu `SPEC-022`**, die Handbuch-Fähigkeitszeile samt Querverweis |

**Der Schnitt geht über die Application-Schicht hinaus in den Lesepfad** — das
ist **kein** Schicht-Schnitt, sondern der Lieferwert: ein Leseport **ohne
Gegenstand** wäre ein Zombie-Slice (die Antwort trüge die Tabellen-Identität
nicht), und die Schnittform `…-store`/`…-usecase`/`…-http` ist die von Modul 5
ausdrücklich verworfene. `ADR-0081` Teilfrage 3 hat die Alternative (nur
Application + HTTP, Schema/Tabelle gar nicht in der Antwort) verworfen.

**Nicht in dieser Liste:** `tools/schema/**` (kein DDL, kein Schema-Eingriff);
`examples/**` und die Wegwerf-Clients (dieser Slice **entblockt** `slice-083`,
zieht es aber nicht mit); `.a-check.yml` (die neuen Dateien liegen in
bestehenden Schichten — `app→ports`, `adapters→ports`, `adapters→domain`
existieren).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)
liegt `Accepted` vor — **erfüllt** —, `Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich der
  Lesepfad-Umbau als eigener Gegenstand (die Filterachse berührt den
  View-Direktzugriff), gehört er zurück — **Store zuerst**, dann der Endpunkt.
  **Benannter Preis:** der erste Slice allein ist dann nicht nutzersichtbar.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass die
  Klartext-Identität ohne Eingriff in die **View** selbst nicht zu haben ist,
  ist das ein Blocker mit Entscheidung — `ADR-0046`/`ADR-0056` sind dann
  berührt, und dieser Slice bleibt bei seinem Zuschnitt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `GET /changes` ist real über die API erreichbar
(Adapter-Test) **und** der Vertrag ist fortgeschrieben (`SPEC-018`, `SPEC-022`,
Handbuch) **und** `make gates` grün **und** die Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Es könnte ein zweiter Lesepfad entstehen.** Der Endpunkt darf nichts
  entscheiden, was die View nicht schon entscheidet. Wächter: derselbe Port,
  dieselbe Delegation, dieselbe Sortierung — das Review prüft es.
  — **Ausgang:** *entfallen — gestrichen mit Begründung*: es gibt **eine**
  Kette (`http → inbound → readchanges → ChangeStorePort → SelectChanges`); die
  Sortierung ist unverändert, der Handler und der Use Case sortieren nichts,
  und der Port wird **genau einmal** gerufen (`TestReadChangesTranslatesQueryToPort`).
  Der Wächter hat getragen: **fünf** Mutationen an genau diesen Zusagen wurden
  rot gesehen — Sortierung (Reviewer), Bereichsende inklusiv (Reviewer),
  Tabellenfilter an Schema gekoppelt (Reviewer), Projektions-Ordnung vertauscht
  (Verifier) — und die Projektion trägt **13 Spalten in der Ordnung der View**.
- **Die Klartext-Identität könnte den View-Vertrag berühren.** `cdc.changes`
  bleibt **unverändert**; der Join lebt im **Lesepfad**, nicht in der View
  (§4 nennt den Blockerfall). — **Ausgang:** *entfallen — gestrichen mit
  Begründung*: `git diff b2d1cc2..HEAD -- tools/schema/` ist **leer** (Verifier
  und Reviewer unabhängig), die View ist unberührt, und
  `TestE2EChangesViewMatchesReadChanges` stellt den View-Zugriff gegen die
  Store-Lesung — grün.
- **Der Endpunkt könnte still begrenzen.** Kein Default-Limit, keine harte
  Obergrenze (`ADR-0081`) — eine eingebaute Grenze wäre ein eigener Beschluss.
  — **Ausgang:** *entfallen — gestrichen mit Begründung*: weder Default noch
  Obergrenze sind gebaut; die Gegenprobe „Default 500 im Use Case" wurde **rot**
  gesehen (Implementer), und der Verifier hat beide **Schwellen** auf 99 gesetzt
  und beide Gates **rot** gesehen — die Grenzen prüfen noch.

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

- **Was hat funktioniert:** **Der Vertrag zuerst.** Der Befund („die API hat
  keinen Lese-Endpunkt") ging zum Architect, und `ADR-0081` hat nicht nur den
  Endpunkt entschieden, sondern die **Formen** fixiert — Port, Filterachse,
  Projektion, Sortierung. Dadurch war die tragende Zusage („**kein zweiter
  Lesepfad**") überhaupt prüfbar: nicht als Versprechen, sondern als Diff und
  als Mutation. Der Verifier hat bestätigt, was zählt: **eine** Kette,
  `tools/schema/**` leer, Projektion in der Spaltenordnung der View.
  Und **die Rollen haben gehalten, wo der Auftrag falsch war**: der Reviewer hat
  eine Zusage in `ADR-0081` gefunden, die keinen Träger hat; der Implementer hat
  in der Fixrunde die **Richtung begründet** statt der Anweisung zu folgen —
  F-2 (die Lücke saß im **Test**, nicht in der Ablehnung) und F-4 (die **Spec**
  wurde gezogen, nicht der Code).
- **Was ging anders als geplant:** (1) Der Slice ist **einmal über seine
  §3-Liste hinausgewachsen** — vier Dateien (`wiring.go`,
  `test/integration/integration_test.go`, `docs/user/e2e-abdeckung.md`,
  `harness/image-hash.txt`), jede begründet und vom Reviewer einzeln geprüft.
  **Benannt, nicht still.** (2) Der Review hat ein Loch gefunden, das **kein
  Lesen** gezeigt hätte: `?limit=0` war nicht an seine Eingabe gebunden — die
  Mutation ließ die Suite grün, obwohl der Endpunkt danach **still unbegrenzt**
  liest. (3) `ADR-0081`s Fitness-Function-Zeile hat **keinen Träger** → eigene
  Adresse `slice-087`. (4) **Meine eigenen Fehler in diesem Zug:** die Linktiefe
  auf `spec/` zweimal falsch (`slice-083`, `slice-086`) und ein Vorlagen-Rest in
  **drei** neuen Slices — gefangen von `d-check` bzw. von meinem Prüfmuster,
  das ich nach dem ersten Fehlschlag auf das allgemeine `<…>`-Muster umgestellt
  habe.
- **Lerneintrag (geschärfte Regel):** *Ein Negativtest, der seine Ausgabe nicht
  an seine **Eingabe** bindet, ist grün ohne Aussage.* Der Test `?limit=0 → 400`
  stellte gegen einen Fake, der den Sentinel unabhängig von der Abfrage lieferte
  — der Status hing an **keinem** Parameterwert. Sichtbar wurde das **nur** durch
  Mutieren der **Eingabeseite**. Regel für den nächsten Fall: **jede Zusage wird
  an ihrer Eingabeseite mutiert**, nicht nur an ihrer Ausgabeseite.
  **Verkörperung offen** (Träger käme im Reviewer-Skill in Frage); als Kandidat
  geführt, nicht behauptet.
- **Steering-Loop-Eintrag:** **keine Verkörperung durch diesen Slice.** Zwei
  Registerbewegungen: `BEO-PGC/generierte-artefakte-ohne-sync-sensor` bekommt
  einen weiteren Beleg (dieser Slice **erzeugt** die E2E-Abdeckungstabelle neu
  und nimmt `tools/schema/plan.yaml` wieder von Hand zurück) → **4×**; **neu
  angelegt** `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (**1×**),
  belegt durch Review F-2.
- **Beobachtungs-Register (`../observations/`):** ein Beleg ergänzt
  (`generierte-artefakte-ohne-sync-sensor/evidence/slice-086.md`, Zähler **4×**)
  und ein Verzeichnis **neu angelegt**
  (`negativtest-ohne-bindung-an-seine-eingabe/`, **1×**). **Kein Zähler
  gesetzt** — er folgt aus den Dateien. **Benannt, nicht gezählt:** die Reibung,
  dass ein frisch geschriebener `SPEC-*`-Abschnitt die tragende ADR nicht
  zitieren darf (`matrix` verbietet `spec → adr`) — **beide** Läufe sind
  darüber gestolpert. Das ist die **beabsichtigte** Referenzrichtung und kein
  Defekt; sie kostet jeden Spec-Autor aber einen roten Lauf.
- **Folge-Slices:** `slice-087` (E2E-Beleg für `GET /changes`) — die Adresse für
  `ADR-0081`s Fitness-Function-Zeile (Review F-1, Verifikation V-1); liegt als
  Datei in `open/`.
- **Risiken aus §6:** alle drei *entfallen, gestrichen mit Begründung* — siehe §6.
- **Drei Paarungen:** nicht hier — dieses Repo führt **Wellen-Betrieb**, die
  Prüfung fällt der `welle-20`-Closure zu (Modul 6 Schritt 3c, auch für Slices
  ohne Wellen-Zugehörigkeit).

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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt Ports, Use Case, HTTP-Adapter
und die Spec-/Handbuch-Pflege in **einem** Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen; die Zählerstände sind am Register
**nachgezählt** (`ls evidence/`), nicht aus diesem Text übernommen:

- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (1×, offen): **kein**
  Treffer — sie betrifft Zahlen in Trägern, die gegen die Messung driften;
  dieser Slice ändert keinen solchen Träger.
- `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (3×, offen,
  Schwelle erreicht — Ausgang beim Lese-Schritt der `welle-20`-Closure):
  **kein** Treffer für diesen Plan; seine Zahlen sind benannt oder aus
  `ADR-0081` zitiert.

**Ergebnis** (Stand: Anlage dieses Plans — kein Versprechen über den Lauf):
kein Eintrag der Sichtung rückt mit diesem Slice auf 3×.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (`ADR-0057` und `SPEC-018` binden die API-Form, die
Port-und-Adapter-Ordnung führt `spec/architecture.md`, `.a-check.yml` prüft die
Kanten), **Phase-Reife** hoch, **Evidenz-/Diskrepanz-Risiko** **niedrig** — der
Port existiert, die Delegation ist bekannt, `ADR-0081` hat die Filterachse
entschieden —, **Reconciliation-Aufwand** null.
