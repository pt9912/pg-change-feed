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


- [x] **LP1 — die zwei HTTP-Familien-Clients.** `examples/http-client/` und
      `examples/sse-client/` liegen vor, sind über `make test` **kompiliert und
      getestet**, und ihr Verhalten trägt, was
      [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
      §Entscheidung zusagt (echter Aufruf gegen die Verwaltungs-API bzw. offenes
      `GET /changes/stream`, jedes Event ausgegeben).
- [x] **LP2 — der gRPC-Client.** `examples/grpc-client/` liegt vor, ist über
      `make test` kompiliert und getestet, und öffnet den **Server-Stream**
      `ChangeStream/StreamChanges`.
- [x] **LP3 — die zwei Träger sind nachgezogen.** `.a-check.yml` nennt **alle
      vier** Clients (heute nur den NATS-Client), und
      `docs/user/benutzerhandbuch.md` führt sie an ihren
      Zugriffs-Abschnitten **namentlich** — mit der Startform und der
      Versionshistorie-Zeile. Auslöser: die zwei Register-Einträge, die mit
      **je 3×** genau diesen Nachzug tragen.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor — **weist er eine Fixrunde aus, deckt ein Delta-Review sie ab.** ([`review-slice-095.md`](../../../reviews/review-slice-095.md), 0 HIGH/MEDIUM, keine Fixrunde)
- [x] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-095.md` liegt vor (Modul 11, frischer Kontext). ([`verify-slice-095.md`](../../../reviews/verify-slice-095.md), Verdikt: konform)
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) — *(entfällt: Greenfield-Bootstrap.)*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — **kein Zaehler wird gesetzt.**
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen — im Repo ohne Wellen-Betrieb hier geprüft, nach dem `git mv` nach `done/`.

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
| `.a-check.yml` | update | Der Kommentar der `examples`-Gruppe nennt heute **nur** den NATS-Client; er trägt die vier und ihre Gemeinsamkeit (nur Standardbibliothek und öffentliche Fremdmodule). Der gRPC-Client führt zusätzlich eine neue Kante `{from: examples, to: contract}` — seit `slice-097`s Umzug auf den öffentlichen Pfad `gen/cdc/stream/v1` nötig und mit ihm begründet; die drei anderen Clients bleiben ohne Kante (`ADR-0079` Festlegung 1). **Korrigiert bei Closure (2. Durchlauf):** diese Zeile behauptete bis dahin für alle vier Clients „keine Kante" — stehen geblieben aus dem ersten `in-progress`-Durchlauf, vor `slice-097`s Umzug (Review-Finding F-1, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`). |
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

  **Eingetreten — der Übergang ist vollzogen.** Zum Zeitpunkt dieser
  Rückführung war der **gRPC-Client (LP2) nicht lieferbar**, und beide Routen
  waren real gemessen: der Import des internen Stubs
  (`internal/adapters/driving/grpc/streamv1`) war `wrong-direction`
  (`make a-check` **Exit 2**), und die **öffentliche Vertragsfläche**
  (`gen/cdc/stream/v1`) **existierte nicht** (`no required module provides
  package …`, **Exit 1**). Die Voraussetzung war
  [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §3 samt seiner Schnitt-Empfehlung **1** („**Der Umzug zuerst** … er soll
  allein stehen") — dieser Vorgänger war zum Zeitpunkt der Rückführung **nie
  geschnitten**, und dieser Plan hatte die Empfehlungen **2 und 3** gebündelt
  und **1 ausgeschlossen**. Der Schnitt war also falsch, nicht die Arbeit: LP1
  und LP3 sind geliefert (Commit mit den zwei HTTP-Familien-Clients und den
  zwei Trägern) und tragen dem Nachfolger. `docs/plan/planning/observations/BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung/evidence/slice-083.md`
  hielt den fehlenden Vorgänger bereits fest.

  **Inzwischen aufgelöst.** `slice-097` hat den Umzug geliefert: die
  öffentliche Vertragsfläche `gen/cdc/stream/v1` existiert seit `1c1d937`
  real, `make generated-sync` bleibt grün. Die zweite gemessene Route
  (`wrong-direction` beim internen Stub) ist damit gegenstandslos — der
  Import läuft jetzt über den öffentlichen Pfad. Diese Zeile bleibt als
  **historischer Beleg** für den `in-progress`→`next`-Übergang stehen
  (Modul 5: der *Grund* wird beim Übergang selbst nachgetragen, nicht
  rückwirkend gelöscht); der **Start**-Trigger unten war zum Zeitpunkt dieser
  Korrektur bereits erfüllt. Fund und Korrektur:
  [`review-slice-097.md`](../../../reviews/review-slice-097.md) (Beobachtung
  außerhalb des Diffs), `BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-097.md`.
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
  — **Ausgang: entfallen.** Der gRPC-Client importiert real
  `google.golang.org/grpc` und dessen Unterpakete — der literale Import ist
  also da, aber kein Risiko im Sinn dieser Zeile: `ADR-0076` §2 benennt
  öffentliche Fremdmodule wie `google.golang.org/grpc` namentlich als
  zulässige Klasse, das Modul war bereits seit `ADR-0060` gepinnt, und
  `go.mod`/`go.sum` sind in diesem Diff unverändert (Verifier #10, Reviewer
  Negativbefund). Es entstand kein *ungeplanter* Bedarf, den die Zeile
  meinte — die Ausnahme war bereits vor diesem Slice entschieden.
- **Ein Client ist ohne laufenden Dienst nicht kompilierbar** (Import einer
  internen Schnittstelle, die den Composition Root zieht). — **Ausgang:
  entfallen.** `make test` kompiliert und testet `examples/grpc-client` ohne
  laufenden Dienst (Verifier #1, Reviewer-Negativbefund); die
  Streamöffnung scheitert — falls überhaupt — erst zur Laufzeit, nicht beim
  Build. Der Import zeigt auf den öffentlichen, statisch vorhandenen Stub
  `gen/cdc/stream/v1`, nicht auf einen internen, Bootstrap-ziehenden Pfad.
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Einträgen; sie ist der einzige Teil dieses Slice, den **kein** Kompilat
  erzwingt. — **Ausgang: entfallen.** Reviewer und Verifier bestätigen
  unabhängig voneinander: der `**Beispiel:**`-Absatz für `grpc-client` steht
  an der erwarteten Stelle im Handbuch, die Versionshistorie ist lückenlos
  (`…1.16, 1.17, 1.18, 1.19`), der Versionskopf ist mitgezogen. Beide
  Register-Einträge (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  `BEO-PGC/handbuch-versionshistorie-uebersprungen`) waren bereits vor diesem
  Slice bei 3× verkörpert; dieser Slice **löst** die Regel ein, statt sie zu
  verletzen — keine neue Evidenzdatei (siehe §7).
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, 3×). — **Ausgang:
  eingetreten, direkt korrigiert.** Der `§3.13`-`grep`-Suchlauf für
  interne Pfad-Referenzen (Implementer/Reviewer) fand nichts Neues — alle
  verbleibenden `internal/adapters/driving/grpc/streamv1`-Erwähnungen sind
  historisch korrekt formuliert. Der Reviewer fand aber unabhängig davon
  (F-1) einen zweiten, engeren Fall derselben Klasse: Die eigene Plan-Zeile
  §3 dieser Datei begründete die `.a-check.yml`-Änderung weiterhin mit
  „keine Kante — `ADR-0079` Festlegung 1", während dieser Durchlauf real
  eine neue Kante liefert (`{from: examples, to: contract}`, seit dem
  öffentlichen Pfad ab `slice-097` nötig und korrekt begründet). Die
  Plan-Zeile war bei ihrer ursprünglichen Niederschrift wahr und wurde durch
  die eigene, spätere Lieferung dieses Slice falsch, ohne dass diese
  Lieferung die Zeile selbst anfasste — derselbe Mechanismus wie bei
  `slice-097`s Fund am eigenen `next/slice-095`-Träger, hier
  selbstreferentiell innerhalb desselben Plan-Dokuments. In diesem Commit
  direkt korrigiert (siehe §3); kein Carveout, kein Folge-Slice nötig, weil
  sofort behebbar (Formvorbild `slice-097` §6). Neue Evidenzdatei
  `arbeit-ueberholt-stehenden-traeger/evidence/slice-095.md` (siehe §7).

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
- **Was hat funktioniert:** Der zweite Anlauf nach der Rückführung
  (`in-progress` → `next` → `in-progress`) lief sauber durch die volle
  Rollen-Sequenz mit je einem Übergabe-Artefakt: Implementer (`d953356`
  Code, `1812e36` Träger-Nachzug) → Reviewer (`117b9d9`, 0 HIGH/MEDIUM,
  1 LOW ohne Fixrunde, eigener Nachvollzug statt Übernahme) → Verifier
  (`d4f4356`, frischer Kontext, eigene Läufe, Verdikt konform) → Planner
  (diese Closure). Beide fälligen Register-Einträge (Handbuch-Nachzug,
  Versionshistorie) wurden **eingelöst**, nicht erneut verletzt — der
  zuvor 3×-verkörperte Mechanismus hat hier sichtbar gegriffen.
- **Was ging anders als geplant:** Die Rückführung hinterließ zwei stehen
  gebliebene Textstellen im Plan-Dokument selbst, die kein Kompilat und kein
  DoD-Häkchen erzwang: §2s Paarungs-Zeile behauptete weiterhin „dieses Repo
  führt Wellen-Betrieb" — falsch seit Anlage (der Kopf dieser Datei nennt
  „Welle: ohne Welle", `harness/conventions.md` deklariert keinen
  Wellen-Betrieb), vermutlich ein unangepasster Vorlagen-Standardtext aus dem
  in der Selbst-Fund-Notiz oben bereits vermerkten `cp`-Überschreiben. Und
  §3s Begründung für die `.a-check.yml`-Änderung blieb bei „keine Kante",
  obwohl der zweite Durchlauf real eine neue Kante liefert (Review-Finding
  F-1) — hier war der Text bei ursprünglicher Niederschrift korrekt und
  wurde erst durch `slice-097`s zwischenzeitlichen Umzug und die eigene
  Lieferung dieses Durchlaufs falsch. Beide Stellen sind bei dieser Closure
  korrigiert (§2, §3).
- **Steering-Loop-Eintrag:** Kein neuer, dritter Beobachtungs-Eintrag für
  „Rückführung/Wiederaufnahme braucht vollständigen Textnachzug" — geprüft,
  aber bewusst nicht angelegt: Beide real gefundenen Stellen fallen bereits
  in **existierende**, längst verkörperte Klassen (siehe unten), und ein
  neuer Eintrag würde dieselbe Beobachtung ein drittes Mal benennen, ohne
  einen neuen Mechanismus zu tragen — die Klassen selbst decken den Fall
  „stehen gebliebener Text nach einer Rückführung" bereits ab (§2 als
  ungeprüfte Vorlagen-Übernahme, §3 als durch eigene spätere Arbeit
  überholter Träger).
- **Beobachtungs-Register (`../observations/`):** Zwei Bewegungen, zwei
  geprüfte Nicht-Treffer.
  1. **`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung/evidence/slice-095.md`
     neu angelegt** — die §2-Paarungs-Zeile behauptete „dieses Repo führt
     Wellen-Betrieb", eine ungeprüfte, aus der Vorlage übernommene
     Tatsachenbehauptung in einem DoD-Kriterium selbst, die dem eigenen
     Slice-Kopf (`Welle: ohne Welle`) und `harness/conventions.md`
     widersprach. Zähler damit **5×**, weiterhin `verkörpert` (kein neuer
     Schwellen-Übertritt).
  2. **`BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-095.md`
     neu angelegt** — Review-Finding F-1: die §3-Begründung „keine Kante"
     wurde durch die eigene, spätere LP3-Lieferung dieses Slice (die neue
     Kante `{from: examples, to: contract}`) falsch, ohne dass diese
     Lieferung die Begründungszeile selbst anfasste — derselbe Mechanismus
     wie bei `slice-097`s Fund am `next/slice-095`-Träger, hier
     selbstreferentiell innerhalb desselben Dokuments. Zähler damit **6×**,
     weiterhin `verkörpert` (`AGENTS.md` §3.13, kein neuer
     Schwellen-Übertritt).
  3. **Geprüft, keine neue Instanz:** `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
     und `BEO-PGC/handbuch-versionshistorie-uebersprungen` (beide bereits vor
     diesem Slice bei 3× verkörpert, als LP3 in die DoD gezogen) — dieser
     Slice **löst** den Nachzug ein (Handbuch nennt alle vier Clients,
     Versionshistorie lückenlos `…1.18, 1.19`), verletzt die Regel nicht
     erneut. Beide Zähler bleiben bei 3×, keine neue Evidenzdatei.
- **Folge-Slices:** keine neuen — die sechs Matrix-Slices (`slice-098`–`103`)
  sind bereits als eigene Pläne in `open/` angelegt und referenzieren
  `slice-095`/`ADR-0090` selbst.
- **Risiken aus §6:** vier Risiken, vier Ausgänge — *entfallen* (Fremdmodul-
  Bedarf durch `ADR-0076` §2 bereits gedeckt, kein ungeplanter Bedarf) ·
  *entfallen* (ohne laufenden Dienst kompilierbar, gemessen) · *entfallen*
  (Handbuch-Nachzug nicht vergessen, gemessen lückenlos) · *eingetreten,
  direkt korrigiert* (überholter Träger — die eigene §3-Begründung, Review
  F-1, in dieser Closure behoben). Siehe §6 für die volle Begründung je
  Zeile.
- **Drei Paarungen:** Repo ohne Wellen-Betrieb — geprüft nach dem `git mv`
  nach `done/` (siehe Commit-Historie). **Anker-Paarung:** entfällt — kein
  `liegt in`-Feld in dieser Notiz, mit diesem Slice wurde keine neue Regel
  verkörpert (beide berührten Klassen waren bereits vor diesem Slice
  verkörpert). **Folge-Slice-Paarung:** entfällt — keine Folge-Slices
  genannt. **Register-Paarung:** grün — alle in dieser Notiz zitierten
  Verzeichnisse existieren mit nicht leerem `evidence/`:
  `dod-begruendung-unzutreffende-tatsachenbehauptung/evidence/`
  (`slice-036.md`, `slice-081.md`, `slice-082.md`, `slice-083.md`,
  `slice-095.md`), `arbeit-ueberholt-stehenden-traeger/evidence/`
  (`slice-091.md`, `slice-093.md`, `slice-094.md`, `slice-096.md`,
  `slice-097.md`, `slice-095.md`),
  `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/evidence/`
  und `handbuch-versionshistorie-uebersprungen/evidence/` (unverändert,
  je 3 Dateien).

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
