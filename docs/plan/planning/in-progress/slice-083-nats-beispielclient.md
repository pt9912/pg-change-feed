# Slice slice-083: NATS-Beispielclient — Wecksignal lauschen, Änderung über HTTP holen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — `examples/` bekommt je Client einen eigenen Schnitt;
die Closure-Bedingung dieses Slice ist die eigene DoD (ein Beispiel, das
übersetzt und beide Schritte geht), kein repo-weites *Mehr*.

**Bezug:** [`ADR-0079`](../../adr/0079-nats-beispielclient-vierter-examples-client.md)
(der **vierte** Beispiel-Client — er hebt `ADR-0076`s Festlegung auf drei
Clients auf und nimmt `examples/nats-client` in unveränderter Form der drei
anderen auf: nur öffentlicher Draht-Vertrag, kein `/internal/`-Import, keine
neue `.a-check.yml`-**Kante**, Doku mit `make test`-Bindung, kein Lauf-Beleg) ·
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(die Beispiel-Clients unter `examples/` — Form, Kompilier-Bindung, die
Gate-Scope-Gruppen) ·
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md) (Core NATS,
kein JetStream) · [`ADR-0056`](../../adr/0056-nats-tabellen-granulares-subjekt.md)
(das tabellen-granulare Subjekt, Korrektur von `ADR-0055` Punkt 2).

**Berührte Spec-Stellen:** [`SPEC-017`](../../../../spec/pflichtenheft.md) —
Subjekt- und Nachrichtenform des Wecksignals (`spec/pflichtenheft.md` §2): das
Subjekt-Schema und der **leere Payload** sind der Vertrag, den das Beispiel
benutzt.

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
  HTTP-, SSE- und gRPC-Beispiel **bleiben, wie sie entschieden sind** — und sie
  sind **entschieden, nicht gebaut** (Review `review-slice-083` F-1); dieser
  Slice legt sie **nicht** an und schneidet nichts um. Ihre Slices sind nicht
  geschnitten; wer sie will, schneidet sie (§Zerlegung der
  [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)).
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
- **Die Coverage-Rampen-Frage** ([`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md),
  teilweise abgelöst durch [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)):
  eigener Vorgang, dort entschieden — dieses Beispiel bewegt keinen Messgegenstand.

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

- [x] `examples/nats-client` ist ein **eigenständiges** CLI-Programm (eigenes
      `main`) und benutzt ausschließlich den öffentlichen Draht-Vertrag:
      **kein** Import aus `/internal/`. **Es ist das erste Beispiel im Repo** —
      die drei von [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
      beschlossenen Clients (`http-client`, `sse-client`, `grpc-client`) sind
      **entschieden, aber nicht gebaut**; ihre Slices sind nicht geschnitten
      (Review `review-slice-083` F-1).
- [x] Es benutzt `github.com/nats-io/nats.go` (bereits Modul-Abhängigkeit) und
      die Standardbibliothek — **keine** neue Abhängigkeit.

**Liefer-Punkt 2 — beide Schritte sind vollständig.**

- [x] Das **Lauschen**: das Subjekt wird aus `<source_id>.<schema>.<table>`
      **abgeleitet**, nicht handgetippt
      ([`SPEC-017`](../../../../spec/pflichtenheft.md)); der Client bezieht
      **keine** Daten aus dem Payload.
- [x] Die **nachfolgende HTTP-Abfrage**: beim Weckruf holt der Client die
      Änderung real über die HTTP-API und gibt sie aus. Ist `CDC_HTTP_ADDR`
      ungesetzt (API deaktiviert), **scheitert er sichtbar**, statt still
      nichts zu tun.
- [x] Die netzlos prüfbare Hälfte — Subjekt-Ableitung und Aufbau der Abfrage —
      liegt als **reine Funktion** mit eigenen Tests vor.

**Liefer-Punkt 3 — es ist zitierbar.**

- [x] `docs/user/benutzerhandbuch.md` trägt einen Abschnitt zum Zugriff über das
      NATS-Wecksignal und **nennt das Beispiel beim Namen** — in der Form der
      bestehenden Zugriffs-Abschnitte der drei anderen Schnittstellen.
- [x] `make gates` grün (Exit direkt, ungepiped).

- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-083.md` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      0 HIGH; das MEDIUM liegt im Plan-Text (F-1), nicht im Code.
- [x] Verifikation durchgeführt, Report unter
      `docs/reviews/verify-slice-083.md` liegt vor (Modul 11, frischer Kontext) —
      **closure-fähig in der Sache**; LP1…LP3 einzeln bestätigt, die drei
      `ADR-0076`-Clients existieren **wirklich nicht** (F-1-Text wahr), der
      Vertrag deckungsgleich mit `server.go:102`/`SPEC-022`/`ADR-0081`, §3.10
      **geprüft und nicht ausgelöst**. Der Verifier hat zwei eigene
      Gegenproben an unberührten Zusagen gefahren (Token-Herkunft,
      `CDC_HTTP_ADDR`-Wächter — beide **ungebunden**, als Grenze benannt) und
      eine Positiv-Kontrolle (`Scheme` → rot). Seine **V-1** (überholte
      Blocker-Erzählung in §4/§6) ist in `56f5aee` bzw. `§6` berichtigt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register (`../reconciliation.md`) — **entfällt**: dieses
      Repo führt die Datei nicht (Greenfield-Bootstrap, kein Inventur-Fund).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — **zwei**
      Belege ergänzt (`dod-begruendung-unzutreffende-tatsachenbehauptung`,
      damit **4×**; `negativtest-ohne-bindung-an-seine-eingabe`, damit **3×**),
      **kein** neues Verzeichnis, **kein Zähler gesetzt**.
- [x] Jedes Risiko aus §6 trägt einen Ausgang — R1 *eingetreten und behoben*;
      R2/R3/R4 *entfallen, gestrichen mit Begründung*.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) — dieses Repo führt
      **Wellen-Betrieb**; die Prüfung fällt der `welle-20`-Closure zu, hier nicht
      geprüft und hier nicht fällig.

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
| `.a-check.yml` | update | die Gate-Scope-Gruppe `examples: ["examples/**"]` aus [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) Festlegung 4 — **ohne Kante**: ihr Objekt ist dieses Beispiel, und eine Kante entsteht dort erst mit ihrem eigenen Objekt (die `contract`-Kante erst mit dem gRPC-Beispiel). Ohne die Gruppe sind die Dateien „ohne Schicht" — nur ein Abdeckungs-Hinweis, kein Exit-Wechsel, aber die errichtete Abschottung fehlt |

**Der genaue Datei-Zuschnitt entsteht im ersten Implementer-Lauf** — die Liste
nennt die Träger. Wer sie erweitert, prüft die Größenregel (≤ 3 Liefer-Punkte,
**eine** Schicht: das Beispiel; die Handbuch-Zeile ist Doku, die ADR ist
Entscheidung).

**Befund des Implementer-Laufs — der zweite Schritt hat keinen Vertrag.** Das
Beispiel ist gebaut (CLI, Subjekt-Ableitung, Aufbau der HTTP-Abfrage, Tests;
`make test` und `make gates` grün), aber sein zweiter Schritt ist gegen die
**heutige** API nicht einlösbar: die HTTP-API führt **keinen**
Changes-Lese-Endpunkt. Das ist **kein** Defekt dieses Slice und **kein**
Implementer-Fehler — es ist die Source-Precedence-Frage dahinter:

- `LH-FA-SST-006` (Lastenheft, **Rang 1**) verlangt in seiner Beschreibung, die
  bestehenden Lese- und Verwaltungsfähigkeiten — „u. a. … **Changes lesen** …" —
  über einen Netzwerkzugriffsweg bereitzustellen. Seine **Akzeptanzkriterien**
  sind demgegenüber **exemplarisch** („eine unterstützte Fähigkeit (**z. B.**
  Consumer-Registrierung)").
- `ADR-0057` (**Rang 4**) hat Option B gewählt und Changes-Lesen „**bewusst
  ausgeschlossen**" — mit benannter Folgepflicht: „Erweiterung um Changes-Lesen
  oder Diagnose/Health braucht [eine eigene ADR]".
- `SPEC-018` (**Rang 2**) hat die Verengung übernommen („Changes-Lesen … bleiben
  außerhalb").

**Die Frage gehört deshalb nicht hierher, sondern zum Architect**: ein ADR darf
die Spezifikation schärfen, **nie das Lastenheft**; ob die Beschreibung oder die
Akzeptanzkriterien binden, ist eine Entscheidung, keine Auslegung durch den
Implementer. `ADR-0057` hat den Weg dorthin als **Option C** bereits
vorgezeichnet (ein `ReadChangesUseCase` samt Inbound Port).

**Der Arbeitsstand liegt auf dem Branch `slice-083-nats-beispielclient`** (lokal,
zwei Commits, kein Push). Er wird **nicht** verworfen: entscheidet die Folge-ADR,
dass die API das Lesen bekommt, ist dieses Beispiel bereits richtig. Bis dahin
ist es das Beispiel für einen Ablauf, den das System noch nicht kann.

### Plan-Nachzug dieses Laufs — der Blocker ist gelöst

[`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md) hat den
Changes-Lese-Endpunkt entschieden, `slice-086` hat ihn geliefert: `GET /changes`
in der Rechtsklasse `reader`, Filter `source`/`schema`/`table`/`from`/`to`/
`limit` ([`SPEC-022`](../../../../spec/pflichtenheft.md)). Der Branch-Stand ist
**unverändert** übernommen — `Subject`, `ChangesURL` und die Token-Quelle des
Clients stimmen mit dem ausgelieferten Endpunkt überein (Pfad `/changes`,
Query-Parameter `source`/`schema`/`table`, `Authorization: Bearer` mit
`CDC_API_TOKEN_READER`). Nachzuziehen waren allein zwei Stellen:
`docs/user/benutzerhandbuch.md` trägt die ausdrückliche Nennung des Endpunkts im
Zugriffs-Abschnitt samt Versionshistorie (1.16). Die `.a-check.yml`-Gruppe
`examples: ["examples/**"]` ist wie geplant **ohne** Kante; ein neuer Träger
oder eine Nicht-Realisierung gegenüber dieser Tabelle liegt nicht vor.

**Nicht in dieser Liste:** `internal/**` (das Beispiel ist kein Harness-Client),
`tools/harness/**` (die Wegwerf-Clients bleiben, `ADR-0076`), `spec/**`
(`SPEC-017` wird **benutzt**, nicht geändert) — und **keine neue `.a-check.yml`-Kante**:
das Beispiel importiert weder `gen/**` noch `internal/**`, `nats.go` ist ein
öffentliches Fremdmodul. Die `tooling→adapters`-Rücknahme und der `contract`-Umzug
aus `ADR-0076`s Zerlegung bleiben **außerhalb**: sie haben hier kein Objekt.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0079`](../../adr/0079-nats-beispielclient-vierter-examples-client.md)
liegt `Accepted` vor — **erfüllt** —, `Verantwortlich:` gesetzt, WIP-Limit frei
(siehe Hinweis). **Nicht** vorausgesetzt: der `contract`-Umzug aus `ADR-0076`s
Zerlegung — dieses Beispiel braucht ihn nicht.

**WIP-Hinweis:** Das Limit ist **1 Slice pro Rolleninhaber** (Modul 5) — es gilt
beim Beanspruchen, nicht beim Schneiden. `slice-081` liegt in `done/`,
`in-progress/` ist leer: der Platz ist frei. Wer später einen zweiten Slice
daneben beansprucht, hat keine Lifecycle, sondern ein Buffet.

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
  **EINGETRETEN — der Grund, mit dem dieser Slice zurückging.** Der Fall ist
  nicht „der Bezug ist nicht auflösbar", sondern größer: **die HTTP-API hat
  keinen Changes-Lese-Endpunkt.** `internal/adapters/driving/http/server.go`
  registriert **zehn** Routen — neun port-gedeckte plus `GET /changes/stream`;
  ein `GET /changes` gibt es **nicht**. Der zweite Schritt des Beispiels zielt
  damit auf einen Vertrag, den es nicht gibt; real antwortet der Abruf `404`.
  Der Implementer hat **keinen** Endpunkt erfunden und `internal/**` unberührt
  gelassen — richtig so: §1 schließt „ein neuer Vertrag" aus, und `AGENTS.md`
  §3.6/Modul 8 lassen den Implementer keine Schnittstelle beschließen.
  **Und wieder aufgelöst:** [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)
  hat den Vertrag entschieden, **`slice-086`** hat ihn gebaut; der Slice wurde
  danach **neu beansprucht** und schließt über `done/`. Die Rückführung war eine
  **Station**, kein Ende — wer §4 als Vorhersage liest, liest sie richtig; wer
  sie als Zustand liest, liest sie falsch (§3 trägt den aktuellen Stand,
  Verifikation `verify-slice-083` V-1).

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
  Weckruf→Abfrage-Weg falsch geht, bliebe damit unbemerkt. — **Ausgang:**
  *eingetreten — und behoben.* Der scharfe Teil ist eingetreten und **nicht
  mehr**: der zweite Schritt hatte **keinen Vertrag** (§4, *EINGETRETEN*), der
  Slice ging deshalb nach `open/` zurück; [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)
  hat den Vertrag entschieden, **`slice-086`** hat ihn gebaut. Was **bleibt**,
  ist die **entschiedene** Grenze aus
  [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  und `ADR-0079` Festlegung 3 — Kompilier-Bindung statt Lauf-Beleg; der Slice
  **weist sie als Grenze aus** und verkauft sie nicht als Beleg (Verifikation
  `docs/reviews/verify-slice-083.md`).
- **Der Client könnte den Payload für Daten halten.** Das Signal ist per Vertrag
  leer ([`SPEC-017`](../../../../spec/pflichtenheft.md)); die Verwechslung ist
  der naheliegendste Fehler und machte
  das Beispiel zum Anti-Vorbild. Wächter: §1 schließt es aus, das Review prüft
  es. — **Ausgang:** *entfallen — gestrichen mit Begründung*: der Client liest
  den Payload **nicht aus** — er zählt nur seine Länge (`len(msg.Data)`), und
  die Verifikation hat es am Code bestätigt.
- **Er könnte Arbeit an den drei von [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  beschlossenen Clients mitnehmen.** Diese drei (`http-client`, `sse-client`,
  `grpc-client`) existieren **nicht** — dieser Slice darf sie weder anlegen
  (eigene Slices, nicht geschnitten) noch ihre Form „angleichen". Das Risiko
  zielt damit auf einen **Bestand, den es noch nicht gibt**; es ist als
  Abgrenzung formuliert, nicht als Gefahr. — **Ausgang:** *entfallen —
  gestrichen mit Begründung*: der Diff setzt ausschließlich
  `examples/nats-client/**` hinzu; `examples/` führt nach dem Zug genau dieses
  eine Programm.
- **Die Folge-ADR könnte ausbleiben.** Ohne sie stünde ein vierter Client neben
  einer `Accepted`-Entscheidung, die drei sagt. — **Ausgang:** *entfallen —
  gestrichen mit Begründung*: [`ADR-0079`](../../adr/0079-nats-beispielclient-vierter-examples-client.md)
  liegt `Accepted` vor und hebt die Drei-Klausel der `ADR-0076` auf.

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

- **Was hat funktioniert:** **Der zweiseitige Aufbau hat die Frage erzwungen,
  die zählte.** Weil das Beispiel „Weckruf → Abfrage" geht, fiel auf, dass die
  Abfrage **keinen Vertrag** hatte — sichtbar, statt still: der Implementer hat
  abgebrochen, statt einen Endpunkt zu erfinden. Und der Slice wurde danach
  **wiederverwendet statt neu geschrieben**: der Branch-Stand von vor dem Bau
  ließ sich konfliktfrei rebasieren, und die Behauptung „unverändert lauffähig"
  wurde gegen den **ausgelieferten** Endpunkt geprüft (`server.go:102`,
  `ADR-0081`, `SPEC-022`) — nicht gegen den Bericht.
  Und **die Kette hat gehalten**: der Reviewer hat einen **Planner**-Fehler
  gefunden (F-1: der Plan berief sich auf Beispiel-Clients, die es nicht gibt),
  der Verifier einen zweiten (V-1: §4/§6 trugen nach dem Unblock die überholte
  Blocker-Erzählung).
- **Was ging anders als geplant:** (1) Der Slice wurde **blockiert**, nach
  `open/` **zurückgeführt** und nach `ADR-0081`/`slice-086` **neu beansprucht** —
  die Rückführung war eine Station, kein Ende (§4). (2) **Mein Fehler (F-1):**
  §1 und §2 verwiesen auf „die drei anderen" Beispiel-Clients, als gäbe es sie;
  `ADR-0076` hat sie **entschieden, nicht gebaut**. Das Häkchen stützte sich
  darauf. (3) **Mein Fehler (V-1):** nach dem Unblock blieb die Blocker-Erzählung
  in §4/§6 stehen. (4) **F-2:** `subject_test.go` trägt einen Kommentar, der
  eine Eigenschaft zusagt, die die Funktion nicht hat (gemessen:
  `ChangesURL("https://feed:8080", …)` → `http://https:%2F%2Ffeed:8080/…`), und
  der Test übt den Fall nicht aus. (5) Der Verifier hat **zwei weitere
  ungebundene Zusagen** gefunden (Token-Herkunft, `CDC_HTTP_ADDR`-Wächter:
  mutiert bleibt die Suite grün) — **kein** DoD-Verstoß, aber als Grenze benannt.
- **Lerneintrag (geschärfte Regel):** *Eine Aussage, die nicht an ihre
  **Eingabeseite** gebunden ist, ist grün ohne Aussage.* **Dritte Gelegenheit in
  drei Tagen, drei Gegenstände:** ein Negativtest (`slice-086`), ein E2E-Beleg
  (`slice-087`), ein Kommentar samt Test (`slice-083`, F-2). Die Regel steht:
  **jede Zusage wird an ihrer Eingabeseite mutiert**, nicht nur an ihrer
  Ausgabeseite. **Verkörperung offen** (Träger: der Reviewer-Skill, wo die
  Mutations-Pflicht schon steht, ohne diese Richtung); als Kandidat geführt.
  **Zweite Regel (F-1, meine):** *ein Plan darf entschiedene, aber nicht gebaute
  Artefakte nicht wie vorhandene behandeln.* Wer Nachbarn zitiert, prüft, dass
  sie da sind — `ls` genügt.
- **Steering-Loop-Eintrag:** **keine Verkörperung durch diesen Slice.** Zwei
  Registerbewegungen, und beide **erreichen die Schwelle**:
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` +`slice-083`
  → **4×**; `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` +`slice-083`
  → **3×** (**Schwelle erreicht** — der Ausgang gehört dem Lese-Schritt der
  laufenden `welle-20`-Closure, Modul 6).
- **Beobachtungs-Register (`../observations/`):** zwei Belege ergänzt, **kein**
  neues Verzeichnis, **kein Zähler gesetzt**. **Benannt, nicht gezählt:** die
  zwei ungebundenen Zusagen des Verifiers (Token-Herkunft, `CDC_HTTP_ADDR`)
  — sie liegen im selben Vorgang wie F-2 und sind damit *eine* Gelegenheit, und
  sie sind **keine** DoD-Verstöße, sondern Grenzen.
- **Folge-Slices:** keiner. Die drei von [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  beschlossenen Clients (`http-client`, `sse-client`, `grpc-client`) sind
  **entschieden, aber nicht geschnitten** — wer sie will, schneidet sie; sie
  sind hier **nicht** angelegt (F-1).
- **Risiken aus §6:** R1 *eingetreten — und behoben*; R2, R3 und R4 *entfallen,
  gestrichen mit Begründung* — siehe §6.
- **Drei Paarungen:** nicht hier — dieses Repo führt **Wellen-Betrieb**, die
  Prüfung fällt der `welle-20`-Closure zu (Modul 6 Schritt 3c, auch für Slices
  ohne Wellen-Zugehörigkeit).
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
