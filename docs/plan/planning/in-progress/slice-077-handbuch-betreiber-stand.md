# Slice slice-077: Benutzerhandbuch auf den aktuellen Betreiber-Stand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist ausschließlich die eigene
DoD (eine Doku-Lücke aufholen, kein repo-weites *Mehr* wie bei
`welle-13`/`welle-14`/`welle-15`; Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Spaltenausschluss),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-/JSON-API),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream),
[`ADR-0057`](../../adr/0057-http-grpc-api.md) (HTTP-Adresse und
Token-Klassen), [`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md)
(Spaltenausschluss-Mechanik), [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
und [`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
(gRPC-Adresse und Zustellsemantik), [`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
(der **dauerhafte** Träger des Ausschlussstandes: der Prozessstart leitet ihn
aus den `applied`-Zeilen der Spalten-Antragsarten ab).

**Berührte Spec-Stellen:** `SPEC-018` (HTTP-Endpunkte, Token-Header-Form),
`SPEC-020` (gRPC-Stream: Dienst, RPC, Metadata-Wertform, Zustellsemantik) —
beide werden hier **beschrieben**, nicht geändert.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` auf den Stand der ausgelieferten
Betreiber-Oberfläche bringen: die drei fehlenden Umgebungsvariablen-Gruppen in
§5 (HTTP/JSON-API, gRPC-Stream) und die zwei fehlenden Aufgaben in §4
(Spaltenausschluss, Netzwerk-Zugriffswege), samt Versionshistorie.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Beispiel-Clients und E2E-Rundläufe der neuen Zugriffswege** —
  `slice-071` (gRPC) und `slice-072` (HTTP/SSE) liefern die Wegwerf-Clients
  und die realen Belege; dieses Dokument beschreibt die *Oberfläche*, nicht
  die Testwerkzeuge.
- **Den dauerhaften Träger des Ausschlussstandes zu ändern** — der Stand ist
  seit `slice-075` dauerhaft: der Prozessstart leitet ihn aus den
  `applied`-Zeilen der Spalten-Antragsarten ab
  ([`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)). Hier
  wird er **beschrieben**; seine Mechanik zu ändern wäre ein anderer Vorgang.
- **Änderungen an `SPEC-018`/`SPEC-020`** — die Verträge stehen bereits; ein
  Doku-Slice ändert keine Spec-Stelle, er beschreibt sie.

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

- [ ] §5 („Umgebungsvariablen des Feed-Containers") trägt die HTTP-Gruppe
      (`CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`) und
      `CDC_GRPC_ADDR` — je mit Aktivierungs-/No-Op-Semantik, geprüft gegen
      `internal/bootstrap/wiring.go` (nicht aus dem Gedächtnis).
- [ ] §4 trägt den Aufgaben-Abschnitt „Spalte vom Ausschluss konfigurieren"
      (`cdc.exclude_column`/`cdc.include_column`, Antrags-Queue mit
      `status = 'applied'`-Poll wie bei `cdc.enable_table`) samt dem
      **dauerhaften** Träger des Ausschlussstandes (`ADR-0065`): der
      Prozessstart leitet ihn aus den `applied`-Zeilen der Spalten-Antragsarten
      ab — ein Neustart verliert ihn **nicht mehr**. **Nachtrag zum Plan:** die
      Fassung, mit der priorisiert wurde, verlangte hier eine „benannte
      Dauerhaftigkeitsgrenze"; die ist seit `slice-075` behoben, und der Text
      beschreibt den Ist-Zustand, nicht die überholte Grenze.
- [ ] §4 trägt die **drei** Netzwerk-Zugriffswege (HTTP/JSON-API, gRPC-Stream,
      HTTP/Server-Sent-Events `GET /changes/stream` — `SPEC-018`, `SPEC-020`,
      `SPEC-021`):
      Erreichbarkeit, Authentifizierung (`Authorization: Bearer` bzw.
      Metadata), Zustellsemantik und der Hinweis, dass die Nachvollziehbarkeit
      beim Lesezugriffsweg bleibt.
- [ ] Versionshistorie fortgeschrieben (Version und Changelog-Zeile).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — hier
      geprüft **von der Slice-Closure selbst**: die Roadmap führt keine offene
      Welle, es gibt also keine Welle-Closure, die sie einsammeln könnte
      (Baseline-Regelwerk `modul-06-roadmap.md` §Was der wellenlose Betrieb
      selbst auslöst).

## 3. Plan (vor Code)


Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | die drei Liefer-Punkte oben, je in einem eigenen Abschnitt |

## 4. Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-059`/`slice-066`/`slice-069`
liegen in `done/` (die drei Oberflächen existieren), `Verantwortlich:`
gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  die HTTP- und die gRPC-Beschreibung zusammen mit dem Spaltenausschluss mehr
  als drei Liefer-Punkte oder mehr als eine Schicht berühren, wird nach
  Zugriffsweg neu geschnitten.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker —
  alle drei Oberflächen sind ausgeliefert und im Code nachlesbar.

## 5. Closure-Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Das Handbuch könnte die Aktivierungs-Semantik der neuen Variablen falsch
  beschreiben (z. B. No-Op bei fehlender Adresse vs. Start-Abbruch bei
  fehlgeschlagener Verbindung) — dieselbe Verwechslung, die `slice-053` für
  `CDC_NATS_URL` behandelt hat. — **Ausgang:** <bei Closure>
- Der Ausschlussstand könnte im Handbuch als **prozesslebensdauer-gebunden**
  beschrieben werden. Er ist dauerhaft: der Prozessstart leitet ihn aus den
  `applied`-Zeilen der Spalten-Antragsarten ab
  ([`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)). Der
  Text muss den **Ist-Zustand** nennen. —
  **Ausgang:** <bei Closure>

## 7. Closure-Notiz


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
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md` führt keine
feinere Sub-Area für die Betreiber-Dokumentation).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Drei Treffer: `BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert
seit `slice-053`, verpflichtet dieses Vorhaben zur Versionsfortschreibung);
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (3×,
dieser Slice ist sein Träger); `BEO-PGC/slice-pfad-als-link-in-berichten`
(als Regel **verkörpert** seit `slice-075` — die Verweisform auf
Lifecycle-wandernde Artefakte; hier einschlägig, weil dieser Plan auf
Slice-Kennungen verweist) und `BEO-PGC/aufschub-adresse-verfaellt`
(1×, **offen**) — **Treffer**: §2 dieses Plans verwies die Paarungen an „die
nächste Welle-Closure", die es nicht mehr gibt; beim Priorisieren auf die
Slice-Closure als Träger nachgezogen. Beleg bei Closure:
`evidence/slice-077.md` — Zähler dann **3×** (`slice-074`, `slice-076`,
`slice-077`), die Schwelle ist damit erreicht und der Lese-Schritt der Closure
weist den Ausgang zu.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.
