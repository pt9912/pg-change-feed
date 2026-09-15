# Slice slice-083: NATS-Beispielclient — Wecksignal lauschen, Änderung über HTTP holen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — `examples/` bekommt je Client einen eigenen Schnitt;
die Closure-Bedingung dieses Slice ist die eigene DoD (ein Beispiel, das
übersetzt und beide Schritte geht), kein repo-weites *Mehr*.

**Bezug:** [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(die Beispiel-Clients unter `examples/` — er legt den Umfang auf **drei** fest;
diesen vierten nimmt eine Folge-ADR auf) ·
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md) (Core NATS,
kein JetStream) · [`ADR-0056`](../../adr/0056-nats-tabellen-granulares-subjekt.md)
(das tabellen-granulare Subjekt, Korrektur von `ADR-0055` Punkt 2).

**Berührte Spec-Stellen:** [`SPEC-017`](../../../../spec/pflichtenheft.md) —
Subjekt- und Nachrichtenform des Wecksignals (`spec/pflichtenheft.md` §2): das
Subjekt-Schema und der **leere Payload** sind der Vertrag, den das Beispiel
benutzt.

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein **eigenständiges CLI-Beispiel** `examples/nats-client`, das den
Wecksignal-Weg **vollständig** geht: es lauscht auf das tabellen-granulare
Subjekt `cdc.changes.<source_id>.<schema>.<table>`
([`SPEC-017`](../../../../spec/pflichtenheft.md),
[`ADR-0056`](../../adr/0056-nats-tabellen-granulares-subjekt.md)) und holt beim
Weckruf die Änderung **selbst** über die HTTP-API — denn das Signal trägt
**keinen** Payload ([`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md)).
Es ist damit das einzige der Beispiele, das die **nicht-offensichtliche**
Nutzung zeigt: nicht „empfange die Änderung", sondern „erfahre, dass du
nachsehen musst".

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein payload-lesendes Beispiel.** Das Signal trägt per Vertrag einen leeren
  Payload ([`SPEC-017`](../../../../spec/pflichtenheft.md)); ein Client, der dort
  Daten erwartet, wäre kein Vorbild,
  sondern ein Fehler mit Kommentar.
- **Eine Änderung an den drei beschlossenen Clients**
  ([`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)):
  HTTP-, SSE- und gRPC-Beispiel bleiben, wie sie entschieden sind; dieser Slice
  fügt hinzu, er schneidet nicht um.
- **Eine zweite NATS-Abhängigkeit.** `github.com/nats-io/nats.go` ist bereits
  Modul-Abhängigkeit (`natsnotify`-Adapter, `go.mod`); das Beispiel benutzt
  dieselbe.
- **Ein neuer Lauf-Beleg-Träger.** [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  hat für alle Beispiele entschieden: **Doku mit Kompilier-Bindung, kein
  Lauf-Beleg** — ein realer Lauf gegen die Compose-Umgebung duplizierte den
  E2E-Beleg, statt ihn zu ergänzen.
- **Der `.a-check.yml`-Vertragsumzug** (`ADR-0076`s erster Zerlegungs-Slice):
  dieses Beispiel braucht **keine** `contract`-Kante — es importiert weder
  `gen/**` noch `internal/**`.
- **Die Coverage-Rampen-Frage** (der Blocker von `slice-081`): eigener Vorgang,
  eigene Entscheidung.

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

**Liefer-Punkt 1 — das Beispiel existiert und taugt als Vorbild.**

- [ ] `examples/nats-client` ist ein **eigenständiges** CLI-Programm (eigenes
      `main`, wie die drei anderen) und benutzt ausschließlich den öffentlichen
      Draht-Vertrag: **kein** Import aus `/internal/`.
- [ ] Es benutzt `github.com/nats-io/nats.go` (bereits Modul-Abhängigkeit) und
      die Standardbibliothek — **keine** neue Abhängigkeit.

**Liefer-Punkt 2 — beide Schritte sind vollständig.**

- [ ] Das **Lauschen**: das Subjekt wird aus `<source_id>.<schema>.<table>`
      **abgeleitet**, nicht handgetippt
      ([`SPEC-017`](../../../../spec/pflichtenheft.md)); der Client bezieht
      **keine** Daten aus dem Payload.
- [ ] Die **nachfolgende HTTP-Abfrage**: beim Weckruf holt der Client die
      Änderung real über die HTTP-API und gibt sie aus. Ist `CDC_HTTP_ADDR`
      ungesetzt (API deaktiviert), **scheitert er sichtbar**, statt still
      nichts zu tun.
- [ ] Die netzlos prüfbare Hälfte — Subjekt-Ableitung und Aufbau der Abfrage —
      liegt als **reine Funktion** mit eigenen Tests vor.

**Liefer-Punkt 3 — es ist zitierbar.**

- [ ] `docs/user/benutzerhandbuch.md` trägt einen Abschnitt zum Zugriff über das
      NATS-Wecksignal und **nennt das Beispiel beim Namen** — in der Form der
      bestehenden Zugriffs-Abschnitte der drei anderen Schnittstellen.
- [ ] Die Folge-ADR zu [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
      liegt `Accepted` vor und ist im ADR-Index eingetragen.
- [ ] `make gates` grün (Exit direkt, ungepiped).

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

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/nats-client/main.go` | neu | das eigenständige CLI: Verbindung, Subjekt-Ableitung, Subscription, Weckruf → HTTP-Abfrage, Ausgabe |
| `examples/nats-client/subject.go` + `_test.go` | neu | die netzlos prüfbare Hälfte: Subjekt-Ableitung und Aufbau der Abfrage als **reine Funktionen** |
| `docs/user/benutzerhandbuch.md` | update | der zitierbare Zugriffs-Abschnitt; das Beispiel wird beim Namen genannt |
| `docs/plan/adr/00NN-…` + `docs/plan/adr/README.md` | neu/update | die Folge-ADR, die `ADR-0076`s Drei-Client-Umfang um diesen vierten erweitert |

**Der genaue Datei-Zuschnitt entsteht im ersten Implementer-Lauf** — die Liste
nennt die Träger. Wer sie erweitert, prüft die Größenregel (≤ 3 Liefer-Punkte,
**eine** Schicht: das Beispiel; die Handbuch-Zeile ist Doku, die ADR ist
Entscheidung).

**Nicht in dieser Liste:** `internal/**` (das Beispiel ist kein Harness-Client),
`tools/harness/**` (die Wegwerf-Clients bleiben, `ADR-0076`), `.a-check.yml`
(keine neue Kante nötig — das Beispiel importiert weder `gen/**` noch
`internal/**`), `spec/**`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): die Folge-ADR zu [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
liegt `Accepted` vor (sie nimmt den vierten Client in den Umfang auf),
`Verantwortlich:` gesetzt, WIP-Limit frei (siehe Hinweis). **Nicht**
vorausgesetzt: der `contract`-Umzug aus `ADR-0076`s Zerlegung — dieses Beispiel
braucht ihn nicht.

**WIP-Hinweis:** Das Limit ist **1 Slice pro Rolleninhaber** (Modul 5). Dieser
Slice darf deshalb erst beansprucht werden, wenn `slice-081` seinen Platz in
`in-progress/` geräumt hat — auch wenn seine eigenen Vorbedingungen früher
erfüllt sind.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich, dass der
  Weckruf→Abfrage-Weg eine eigene Zustandsmaschine braucht (Reconnect,
  Deduplizierung, Rückstand), gehört er zurück zur Zerlegung — **Lauschen** und
  **Abholen** als zwei Slices —, nicht als wachsender Sammelclient.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass die
  HTTP-Abfrage einen Vertrag braucht, den es noch nicht gibt (etwa weil der
  Weckruf-Bezug über die vorhandenen Endpunkte nicht auflösbar ist), ist das ein
  Blocker mit Entscheidung — die Zusage dieses Slice ist „zwei Schritte, kein
  neuer Vertrag".

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** das Beispiel übersetzt im Kompilierpfad des Repos
**und** die netzlos prüfbare Hälfte läuft real grün **und** `make gates` grün
**und** die Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Das Beispiel könnte übersetzen, ohne zu laufen.** [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  hat entschieden, dass Beispiele **keinen** Lauf-Beleg tragen
  (Kompilier-Bindung statt E2E-Lauf). Ein Beispiel, das übersetzt, aber den
  Weckruf→Abfrage-Weg falsch geht, bliebe damit unbemerkt. — **Ausgang:** <bei Closure>
- **Der Client könnte den Payload für Daten halten.** Das Signal ist per Vertrag
  leer ([`SPEC-017`](../../../../spec/pflichtenheft.md)); die Verwechslung ist
  der naheliegendste Fehler und machte
  das Beispiel zum Anti-Vorbild. Wächter: §1 schließt es aus, das Review prüft
  es. — **Ausgang:** <bei Closure>
- **Er könnte die drei anderen Clients nachziehen.** `ADR-0076`s Entscheidung
  begründet deren Form; ein „angleichender" Umbau wäre eine Änderung an
  entschiedener Sache. — **Ausgang:** <bei Closure>
- **Die Folge-ADR könnte ausbleiben.** Ohne sie stünde ein vierter Client neben
  einer `Accepted`-Entscheidung, die drei sagt. — **Ausgang:** <bei Closure>

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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt `examples/**`, `docs/user/**`
und die ADR-Ablage in **einem** Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand. Die Zählerstände
sind am Register **nachgezählt** (`ls evidence/`), nicht aus diesem Text
übernommen:

- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (3×, **verkörpert**): **kein** Treffer — die Klasse ist „eine **neue**
  Betreiber-Oberfläche bleibt im Handbuch unerwähnt"; dieser Slice **zieht
  nach**, er löst sie nicht neu aus. Die offene Handbuch-Remediation bleibt ihr
  eigener Träger.
- `BEO-PGC/adapter-fehler-ausgang` (2×, offen): **kein** Treffer — die
  Adapter-Grenze, nicht Beispiel-Clients.

**Ergebnis** (Stand des Plan-Zeitpunkts, kein Versprechen über den Lauf): kein
Eintrag steht über der Schwelle, und keiner rückt mit diesem Slice auf 3×.
Erreicht einer sie doch, steht er in §7, und sein Ausgang fällt dem Lese-Schritt
der laufenden Welle-Closure zu (Modul 6), nicht dieser Closure.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die Beispiel-Form ist über `ADR-0076` verankert,
die benutzte Schnittstelle über [`SPEC-017`](../../../../spec/pflichtenheft.md),
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md) und
[`ADR-0056`](../../adr/0056-nats-tabellen-granulares-subjekt.md)),
**Phase-Reife** hoch, **Evidenz-/Diskrepanz-Risiko** **niedrig** — der Vertrag
ist dokumentiert, sein Kern (leerer Payload) sogar ausdrücklich —,
**Reconciliation-Aufwand** null.
