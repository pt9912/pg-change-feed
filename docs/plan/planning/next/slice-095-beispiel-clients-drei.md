# Slice slice-095: Beispiel-Clients unter `examples/` — die drei fehlenden

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf Nutzerentscheidung, ein einzelner Slice
ohne Closure-Bedingung jenseits seiner DoD (Modul 6 §Wann Arbeit eine Welle
braucht). `welle-20` ist davon unberührt.

**Bezug:** [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(die **drei** Clients, ihre Verhaltensweisen, die Startform und der Doc-Kommentar) ·
[`ADR-0079`](../../adr/0079-nats-beispielclient-vierter-examples-client.md)
(der vierte — gebaut, dient als **Form-Vorbild**) ·
[`ADR-0057`](../../adr/0057-http-grpc-api.md) und
[`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md) (die zwei
HTTP-Oberflächen) · [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(der gRPC-Server-Stream) · `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(**3×**, verkörpert) und `BEO-PGC/handbuch-versionshistorie-uebersprungen`
(**3×**, verkörpert) — beide treffen diesen Slice.

**Berührte Spec-Stellen:** `LH-FA-SST-006` (die HTTP-API), `LH-FA-SST-008` (der
Live-Change-Stream über gRPC und SSE) — dieser Slice **zeigt** sie, er ändert
sie nicht.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.


**Ziel:** `examples/` wird **vollständig**: die drei Client-Programme, die
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
§Entscheidung benennt, entstehen — `examples/http-client/` (ein echter
Anfrage/Antwort-Aufruf gegen die Verwaltungs-API, `GET /tables` mit dem
`reader`-Token), `examples/sse-client/` (öffnet `GET /changes/stream`, gibt
jedes Event aus) und `examples/grpc-client/` (öffnet
`ChangeStream/StreamChanges`, gibt jede Nachricht aus) — je als eigenes
`main`-Paket, mit der Startform `go run ./examples/<name>`, den Umgebungsvariablen
des Handbuchs (`CDC_HTTP_ADDR` bzw. `CDC_GRPC_ADDR`, `CDC_API_TOKEN_READER`)
und der Flag-Übersteuerung `-addr`/`-token`.

**Warum sie fehlen:** `slice-083` hat den **vierten** gebaut (NATS,
[`ADR-0079`](../../adr/0079-nats-beispielclient-vierter-examples-client.md)) —
der `examples/`-Ordner trägt seither **ein** Programm statt vier, und die zwei
Träger, die die Clients nennen, nennen nur dieses eine.

**Was dieser Slice liefert:** drei Programme, ihre Unit-Tests und die zwei
Träger. Er ändert **keinen** Produkt-Code und **keine** Zusage.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der NATS-Client.** Er existiert (`slice-083`) und ist das **Form-Vorbild**;
  ihn umzubauen wäre Arbeit an Bestand ohne Adresse.
- **Ein neuer Belegträger für die Draht-Verträge.** Die Beispiele sind
  **Vorbilder**, keine Belegträger — das ist `ADR-0076`s Entscheidung. Die
  E2E-Belege bleiben bei den Wegwerf-Clients unter `tools/harness/`
  (`httpclient`, `sseclient`, `grpcclient`, `natssub`).
- **Ein eigenes Gate für die Beispiele.** `make test` fährt `go test -race ./...`
  und **kompiliert sie damit mit**; ein zweites Gate wäre Zeremonie.
- **Eine Änderung an der API oder am Stream-Vertrag.** Die Clients **benutzen**
  die Oberflächen; wenn einer sie nicht bedienen kann, ist das ein Befund, kein
  Umbau.
- **Ein Beispiel für die DB-Adapter oder die SQL-Administration.** [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  schneidet `examples/` auf die **öffentlichen Draht-Oberflächen**; ein
  SQL-Beispiel wäre ein anderer Vorgang.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.


- [ ] **LP1 — die zwei HTTP-Familien-Clients.** `examples/http-client/` und
      `examples/sse-client/` liegen vor, sind über `make test` **kompiliert und
      getestet**, und ihr Verhalten trägt, was
      [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
      §Entscheidung zusagt (echter Aufruf gegen die Verwaltungs-API bzw. offenes
      `GET /changes/stream`, jedes Event ausgegeben).
- [ ] **LP2 — der gRPC-Client.** `examples/grpc-client/` liegt vor, ist über
      `make test` kompiliert und getestet, und öffnet den **Server-Stream**
      `ChangeStream/StreamChanges`.
- [ ] **LP3 — die zwei Träger sind nachgezogen.** `.a-check.yml` nennt **alle
      vier** Clients (heute nur den NATS-Client), und
      `docs/user/benutzerhandbuch.md` führt sie an ihren
      Zugriffs-Abschnitten **namentlich** — mit der Startform und der
      Versionshistorie-Zeile. Auslöser: die zwei Register-Einträge, die mit
      **je 3×** genau diesen Nachzug tragen.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor — **weist er eine Fixrunde aus, deckt ein Delta-Review sie ab.**
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-095.md` liegt vor (Modul 11, frischer Kontext).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) — *(entfällt: Greenfield-Bootstrap.)*
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — **kein Zaehler wird gesetzt.**
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen — dieses Repo führt Wellen-Betrieb; die Prüfung fällt der nächsten Welle-Closure zu.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.


| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/http-client/**` | neu | Der Anfrage/Antwort-Client (`ADR-0076`), Form-Vorbild `examples/nats-client` (`slice-083`). |
| `examples/sse-client/**` | neu | Der SSE-Client (`ADR-0076`, `ADR-0061`). |
| `examples/grpc-client/**` | neu | Der gRPC-Server-Stream-Client (`ADR-0076`, `ADR-0060`). |
| `.a-check.yml` | update | Der Kommentar der `examples`-Gruppe nennt heute **nur** den NATS-Client; er trägt die vier und ihre Gemeinsamkeit (nur Standardbibliothek und öffentliche Fremdmodule, keine Kante — `ADR-0079` Festlegung 1). |
| `docs/user/benutzerhandbuch.md` | update | Die Zugriffs-Abschnitte nennen die Clients namentlich (heute nur `examples/nats-client`) **und** die Versionshistorie bekommt ihre Zeile. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.


**Start** (`next` → `in-progress`): [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
ist `Accepted`, `examples/` trägt **ein** Programm, und der Nutzer hat den
Nachzug beauftragt. Ohne Rückfrage feststellbar.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): wenn ein Client **mehr als ein** Paket
  braucht oder eine Zusage der API sich als nicht bedienbar erweist. Dann ist
  der Schnitt falsch — ein Beispiel, das den Vertrag ändern müsste, ist ein
  anderer Vorgang.
- `in-progress` → `open` (blockiert): wenn ein Client **ohne** einen laufenden
  Dienst nicht einmal **kompilierbar** ist. Dann gehört die Frage nach der Form
  (Build-Tag? eigenes Modul?) in eine Entscheidung, nicht in diesen Slice.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.


**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — die drei Programme
sind über `make test` kompiliert und getestet, und **beide** Träger nennen alle
vier — **und** `make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Naheliegender Kandidat: die zwei Register-Einträge, die mit **je 3×** den
Handbuch-Nachzug tragen — ob dieser Slice sie einlöst oder erneut trifft,
entscheidet der Lauf.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.


- **Ein Client braucht mehr als die Standardbibliothek** (z. B. eine SSE-Bibliothek).
  — **Ausgang:** <…>
- **Ein Client ist ohne laufenden Dienst nicht kompilierbar** (Import einer
  internen Schnittstelle, die den Composition Root zieht). — **Ausgang:** <…>
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Einträgen; sie ist der einzige Teil dieses Slice, den **kein** Kompilat
  erzwingt. — **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, 3×). — **Ausgang:** <…>

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).


- **Nachtrag zur Anlage (Selbst-Fund):** Diese Datei wurde nach dem `cp` aus der
  Vorlage **vollständig überschrieben** statt in-place gefüllt — die Skill-Regel
  nennt das „derselbe Verstoß, weil der `cp` verworfen wird". Aufgefallen an
  einem harten Merkmal: die Datei trug **0** der **8** „Regeln dieser Sektion"-
  Zeilen, die die Vorlage und die Nachbar-Slices tragen. Nachgezogen. Der Rest
  des Berichts folgt bei der Closure.
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.


**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `examples/**`,
`.a-check.yml` und `docs/user/benutzerhandbuch.md` — die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Vier
Treffer: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `handbuch-versionshistorie-uebersprungen` (**3×**, verkörpert) —
**beide als LP3 in die DoD gezogen**, weil sie genau diesen Nachzug tragen;
`dod-begruendung-unzutreffende-tatsachenbehauptung` (**4×**, verkörpert) und
`arbeit-ueberholt-stehenden-traeger` (**3×**, offen) als Risiko in §6.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**.
