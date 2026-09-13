# Architect-Review welle-1 — Slice-Pläne gegen die Entscheidungslage

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-09.
**Eingang:** drei Slice-Pläne aus `open/` (slice-001…003, welle-1).
**Ausgang:** Übergabe-Artefakt an Planner — ADR-Bezüge bestätigt/ergänzt/beanstandet.
**Harte Regel eingehalten:** keine Accepted-ADR wurde in-place geändert; keine
Folge-ADR-Datei angelegt. Folge-ADR-Vorschläge stehen am Ende und warten auf Vorgabe.

**Verdikt je Slice:** slice-001 **mit Ergänzungen lieferbar** · slice-002 **mit
Ergänzungen lieferbar** · slice-003 **beanstandet — vor `next`-Übergang
nachzuziehen (fehlender Port, fehlende Kennung)**. Kein Plan widerspricht einer
Accepted-ADR inhaltlich; alle Befunde sind Bezugs-Lücken oder Form-Defekte, kein
Fall „Plan behauptet eine Lockerung".

---

## slice-001 — Go-Modul-Bootstrap und Gate-Aktivierung

### ADR-Bezüge bestätigt

- **ADR-0039** trägt `cmd/pg-change-feed/main.go` und die `internal/`-Form
  (ADR-0039 §Entscheidung, Baum). Der Plan baut genau diesen Baum nicht an, aber
  die Struktur, die er später befüllt, hält der Deklaration. Anchor korrekt.
- **ADR-0010 / ADR-0006** als Out-of-Scope-Verweis für Store-/Stream-Adapter:
  korrekt — beide Accepted, beide treffen genau die ausgeschlossenen Adapter.
- **CGO aus** (slice-001:78): korrekt; Anker ist das CGO-Ziel, das ADR-0039
  (Kontext, „CGO-Ziel … bleiben unverändert bestehen") aus dem superseded
  ADR-0038 fortgilt. Empfehlung: im DoD/§3 als **ADR-0039** ankern, nicht auf
  ADR-0038 verweisen — ADR-0038 ist `Superseded` und trägt keine Bindung mehr.
- `.a-check.yml` (Schichten-Globs, Edges) hält gegen den ADR-0039-Baum:
  `domain`/`ports`/`app`/`adapters`/`composition_root` decken
  `internal/domain/**`, `internal/application/**`, `internal/adapters/**`,
  `internal/bootstrap/**`, `cmd/**` ab — konsistent. Kein Widerspruch.

### Beanstandet / ergänzt

1. **ADR-0036-Bezug fehlt** (slice-001:109–111, Kopf-Bezug Zeile 12). Der Slice
   aktiviert genau das Gate, dessen Existenz ADR-0036s Hochschalt-Trigger ist
   („Import-Linting-Gate existiert und grün läuft" — ADR-0036
   §Re-Evaluierungs-Trigger). Ohne Bezug sieht der Slice die Entscheidung nicht,
   die er vollzieht — und die welle-1-Closure (Trigger-Audit) hat später keinen
   Beleg-Anker für die Hochschaltung. **Ergänzung:** `ADR-0036` in Kopf-Bezug
   und in die `.a-check.mk`/`.a-check.yml`-Zeile von §3.
2. **[`SPEC-012`](../../../spec/pflichtenheft.md) falsch zitiert** (slice-001:107): „Go 1.26 ([`SPEC-012`](../../../spec/pflichtenheft.md)-Umfeld)" —
   [`SPEC-012`](../../../spec/pflichtenheft.md) ist `PG_MAJOR_VERSIONS` (PostgreSQL 17/18,
   spec/pflichtenheft.md §3), hat nichts mit der Go-Toolchain zu tun. Die
   Go-Version trägt heute keinen Spec-Anker (ADR-0038/0039 sagen Sprache, nicht
   Version). **Ergänzung:** [`SPEC-012`](../../../spec/pflichtenheft.md)-Zitation streichen; Go-Version entweder
   ohne ID lassen oder eine ADR-/SPEC-Verankerung nachziehen.
3. **Aktivierungs-Bedingung vs. Plan-Baum** (Makefile:16–22 vs. slice-001
   §3): Der Makefile-Kommentar nennt als beobachtbare Bedingung „erste
   `.go`-Datei unter `internal/`"; slice-001 legt aber nur `cmd/pg-change-feed/
   main.go` an — `internal/` entsteht erst mit slice-002. Der aktivierte
   a-check-Lauf ist damit über die Layer-Globs leer und nur über
   `composition_root: cmd/**` nichtleer (`.a-check.yml:24`). Das ist kein
   leerer Gate im Modul-13-Sinn, aber der Kommentar im Makefile ist mit der
   Aktivierung mitzuziehen (Zeile 18–19), sonst beschreibt er nach der
   Aktivierung einen abwesenden Zustand (AGENTS.md §3.7). **Ergänzung:** die
   Makefile-Kommentar-Zeile ausdrücklich in die `Makefile`-Zeile von §3
   aufnehmen.
4. **Platzhalter in Pflichtzeile** (slice-001:88): „Doku-Update für
   <Schnittstelle X>" — Vorlagen-Rest in einem geplanten DoD. Streichen oder
   konkretisieren (hier: `harness/README.md`-Update ist als Plan-Zeile schon
   drin — Zeile 111 deckt es; dann entfällt der Punkt ersatzlos).

### DoD-Realismus

Beobachtbar und prüfbar: `--version`-Beleg ([`LH-QA-POR-003`](../../../spec/lastenheft.md) Teil-Beleg sauber
markiert), `image-hash.txt` ([`LH-QA-OPS-001`](../../../spec/lastenheft.md) Teil-Beleg), a-check-Aktivierung.
`make gates` und `make a-check` sind getrennte Targets — die welle-1-Closure
fordert beide, konsistent mit der Fragment-Struktur (a-check hängt nicht an
`GATE_CHECKS`).

**Nebenbefund (nicht Slice, aber welle-1-Berührung):** `harness/README.md:127`
bindet `make image` an **ADR-0038** — diese ist `Superseded`; der gültige Anker
ist ADR-0039. Das ist eine Index-Zeile, kein ADR-Inhalt — im slice-001-
Doku-Update mitziehen.

---

## slice-002 — Domänenkern: Modelle, Invarianten, ClockPort

### ADR-Bezüge bestätigt

- **ADR-0029** (Invarianten in Konstruktoren) korrekt als Träger der
  Konstruktions-Verpflichtung; **ADR-0040** (ClockPort, `Now() domain.TimePoint`,
  Fake Clock in Tests) deckt Port-Typ und Test-Ansatz; die Gefahr aus
  slice-002 §6 (Risiko 2) ist durch die ADR-0040-Fitness-Funktion bereits
  gerahmt: `domain` und `usecase/*` importieren `time` nicht — ein `time.Time`
  an der Port-Kante (`internal/application/port/outbound/clock.go`) ist von der
  Regel **nicht** erfasst und damit zulässig. Das Risiko kann bei Closure
  sauber „entfallen" tragen, mit diesem Anker.
- **ADR-0032** als Begründung der Technologieunabhängigkeit im Out-of-Scope:
  korrekt; **ADR-0025** (Domain Events) als Ausschluss-Begründung korrekt —
  die ADR erlaubt, fordert nicht.
- **[`LH-FA-CAP-004`](../../../spec/lastenheft.md)/-005 als Teil-Beleg** markiert (slice-002:88–89): sauber —
  die Ordnungs-/Zusammengehörigkeits-Aussage ist erst am Lese-/Store-Pfad voll
  belegbar, das Modell ist der Teil.

### Beanstandet / ergänzt

1. **ADR-0005 fehlt im Bezug** (slice-002:10–17). SourcePosition ist das
   Zentralobjekt dieses Slices; ADR-0005 (Accepted, schärft [`SPEC-003`](../../../spec/pflichtenheft.md)) ist die
   Entscheidung, die das Positionstyp-Design bindet. Das Risiko in §6
   (slice-002:162–166) sucht die Semantik-Quelle nur bei ADR-0029 und
   [`LH-FA-REA-004.a`](../../../spec/pflichtenheft.md) — die Invariante (0029) sagt die Ordnung, **ADR-0005** sagt
   den Typ. **Ergänzung:** ADR-0005 in Kopf-Bezug und in die Risiko-Zeile.
2. **Ziel-Liste divergiert von [`ARC-001`](../../../spec/architecture.md) und ADR-0004** (slice-002:43–47): Die
   Aufzählung nennt 7 Modelle und lässt **Source und SourceTable** weg — beide
   stehen in [`ARC-001`](../../../spec/architecture.md) (Domain-Core-Zeile, spec/architecture.md §1) und in
   ADR-0004 §Entscheidung („Kernobjekte: Source, SourceTable, …"). Der DoD
   („Domänenmodelle je [`ARC-001`](../../../spec/architecture.md)", slice-002:83) verspricht damit mehr als das
   Ziel listet. Zwei saubere Ausgänge: (a) Source/SourceTable ins Ziel (dann
   gehört [`LH-FA-DAT-002`](../../../spec/lastenheft.md) — Quelltabelle identifizierbar — in den Bezug), oder
   (b) expliziter Out-of-Scope-Punkt Klasse 1/2 mit Kennung. Nichts tun ist
   keine der beiden — der DoD-Text je [`ARC-001`](../../../spec/architecture.md) prüft sonst gegen eine Liste, die
   das Ziel nicht hält.
   *(Randnotiz aus der Entscheidungslage, keine Slice-Frage: ADR-0039s
   model/-Kommentar führt dieselbe 7er-Liste und lässt Source/SourceTable
   ebenfalls weg — ADR-0039 ist Struktur-ADR, seine Liste ist illustrativ
   (`#`-Kommentar) und hebt die [`ARC-001`](../../../spec/architecture.md)-Aufzählung nicht auf; [`ARC-001`](../../../spec/architecture.md) steht im
   Stabilitäts-Rang über der ADR. Kein Folge-ADR nötig, solange der Slice die
   [`ARC-001`](../../../spec/architecture.md)-Liste als maßgeblich liest.)*
3. **Platzhalter in Pflichtzeile** (slice-002:96) — wie slice-001.

### DoD-Realismus

Fake Clock ohne Verbraucher: In diesem Slice konsumiert kein Use Case den Port
(der erste ist der Capture-Pfad); der DoD-Punkt „mit Fake Clock in den Tests"
ist deshalb nur als Typ-/Konstruktions-Beleg realistisch, nicht als
Verhaltens-Beleg. Nicht beanstandet, aber der Closure-Eintrag sollte den Beleg
als Konstruktions-Beleg kennzeichnen, sonst meldet der Verifier einen leeren
Test. Der §6-Risiko-Ausgang („entfallen", slice-002:167–169) trägt das schon.

---

## slice-003 — Capture-Persist-Pfad: Ports, Service, Persist-before-ACK

### ADR-Bezüge bestätigt

- **ADR-0011 / ADR-0027 / ADR-0009** korrekt: Ordnung ([`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md)), Ordnung
  im Application Layer, ChangeStorePort als Outbound Port. Die Teil-Beleg-Kenn-
  zeichnungen im DoD ([`LH-QA-REL-002`](../../../spec/lastenheft.md), [`LH-FA-CAP-006.a`](../../../spec/pflichtenheft.md)) sind sauber und realistisch
  für Fake-Port-Tests.
- **ADR-0021 als Proposed korrekt behandelt** (Out-of-Scope, slice-003:59–61) —
  keine Bindung aus einer Proposed-ADR beansprucht, Entscheidung auf Welle 2
  verschoben. Korrekt.
- **ADR-0028 / ADR-0039** Pfad-Kandidaten deckungsgleich mit dem
  ADR-0039-Baum (`port/inbound/capture.go`, `port/outbound/`, `usecase/capture/`
  mit Command/Result).

### Beanstandet

1. **`ReplicationAckPort` fehlt im Plan** (slice-003 §3, Zeilen 116–121). Die
   Persist-before-ACK-Orchestrierung ist per **ADR-0027** „für jeden Aufruf
   über zwei Ports (ChangeStorePort, ReplicationAckPort) hinweg entworfen";
   ADR-0007 macht den ACK zum Outbound Port; ADR-0039 listet
   `ReplicationAckPort` im outbound-Bestand; die Sequenz in
   spec/architecture.md §4 ([`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md)) zeigt App → RAP → AA. Ohne das
   Port-Interface ist DoD-Punkt 1 („ACK nur nach dauerhafter Persistenz")
   gegen Fake Ports **nicht prüfbar** — es fehlt die Fake-Gegenstelle. Das ist
   der Port aus Prüffrage 4, den der Slice schon berührt.
   **Ergänzung:** Zeile `internal/application/port/outbound/replicationack.go`
   (neu, ADR-0007) in §3. Umfang wächst nicht über die drei Liefer-Punkte —
   ein Interface, kein Adapter.
2. **Out-of-Scope Klasse 1 ohne Kennung** (slice-003:56–58): „Reale
   PostgreSQL-Adapter … — ein Folge-Slice übernimmt es (Welle 2 …)". „Welle 2"
   ist kein Lifecycle-Verweis; die Regel verlangt eine Kennung, die die Sendung
   annimmt. slice-002 erfüllt dieselbe Pflicht korrekt („slice-003 setzt die
   Ports"). **Ergänzung:** Slice-Kennung benennen (Welle-2-Eröffnung legt sie
   an) oder die Klasse wechseln („anderer Vorgang" — dann trägt die
   Welle-2-Nennung die Begründung, aber der Verweis heißt so).
3. **Idempotenz-Risiko zu früh „entfallen"** (slice-003:171–173): ADR-0011
   macht idempotente Persistenz zur „Pflicht, kein Optimismus"; ein Fake-Port-
   Test, der die Persistenz wiederholt, belegt nur die **Aufruf-Ordnung**, nicht
   die Deduplizierung ([`SPEC-002`](../../../spec/pflichtenheft.md) `change_id` wird erst vom realen Store
   durchgesetzt). Der Ausgang „entfallen" ist an der Form der drei Ausgänge
   korrekt besetzt, trägt aber nicht — der richtige Ausgang ist **weiter
   offen** (ins Register, falls er 3× auffällt) oder ein Ausgang mit
   Welle-2-Bezug. Der Plan-Zeile „Fake-Port-Tests: … Idempotenz"
   (slice-003:121) ist entsprechend zu lesen: Sie prüft das Wiederholungs-
   verhalten des Services, nicht die Deduplizierungs-Garantie.
4. **Kopf-Bezug unvollständig** (slice-003:14–18): ADR-0030 (Testpyramide) und
   ADR-0039 sind nur inline in Ziel/§3 zitiert; **ADR-0012** (At-Least-Once)
   fehlt ganz, obwohl DoD-Punkt 2 genau seine Konsequenz prüft („Crash zwischen
   Persistenz und ACK erzeugt höchstens erneute Verarbeitung" — ADR-0012
   §Kontext). **Ergänzung:** ADR-0012 in den Bezug; ADR-0030/0039 nachziehen,
   damit der Reviewer gegen den Kopf prüfen kann.

### DoD-Realismus

Nach Ergänzung von (1) prüfbar. Ohne ReplicationAckPort wäre DoD-Punkt 1 ein
Kriterium, dessen Beleg-Artefakt der Plan gar nicht vorsieht.

---

## Welle-1-Kontext (kein Beanstandung, zwei Hinweise)

- **Start-Trigger** (welle-1:44–49): ADR-0039 und ADR-0040 sind `Accepted`
  (ADR-Index-Zeilen 52–53) — der Trigger ist eingetreten und beobachtbar am
  Index. Korrekt platziert (kein Ergebnis dieser Welle).
- **ADR-0036 in der Closure:** Der Closure-Trigger fordert `make a-check` grün
  (welle-1:61–63) — genau ADR-0036s Hochschalt-Bedingung. Der Trigger-Audit der
  Welle-Closure (Modul 6, Schritt 2) muss die Hochschaltung Proposed→Accepted
  oder die Entscheidung dafür dann mit Übergabe-Artefakt tragen
  (Planner → Architect → Planner). Die welle-1-Datei muss das nicht
  vorwegnehmen; der Architect-Übergang bei Closure ist der Träger. Nicht
  nachtragen — nur merken.
- **Out-of-Scope der Welle** (welle-1:103–113) deckt sich mit den drei
  Slice-Ausschlüssen; kein Konflikt zwischen Welle und Slices.

---

## Folge-ADR-Vorschläge (nur Vorschlag — keine Datei angelegt)

| Nr. | Anlass | Form |
|---|---|---|
| 1 | Das Architektur-Gate ist als **a-check** (`.a-check.yml`, gepinntes Release) implementiert, während ADR-0036 (Proposed) und die Fitness-Tabelle ADR-0039s `depguard` nennen. Beide ADRs tragen dieselbe Regel in zwei Werkzeug-Sprachen. | Nach dem ersten grünen `make a-check` (welle-1-Closure): Folge-ADR mit `Supersedes ADR-0036`, das a-check als Maschinenform der §2-Constraints festlegt — inklusive der benannten **Deckungsgrenze**, dass a-check driving/driven als eine Schicht prüft (`.a-check.yml:15–23`) und die ADR-0039-Fitness-Differenz (driving nur inbound+domain, driven nur outbound) Review-Prüfpflicht bleibt. |
| 2 | *(nur falls der Implementer Source/SourceTable aus `internal/domain/model/` halten will)* | Klärungs-ADR gegen die [`ARC-001`](../../../spec/architecture.md)/ADR-0004-Liste — dann **mit neuer Evidenz** (warum die zwei Objekte keinen Core-Typ brauchen), nicht gegen die 7er-Liste des ADR-0039-Kommentars argumentiert. Vorerst nicht nötig: [`ARC-001`](../../../spec/architecture.md) ist maßgeblich, slice-002 ergänzt oder schließt explizit aus. |

Keiner der Vorschläge ist eine In-place-Änderung; ADR-0036 ist `Proposed`,
wird aber aus derselben Disziplin heraus per `Supersedes` abgelöst, nicht
überarbeitet.

---

## Beleg-Anker je Bestätigung (Kurzfassung)

| Slice | Aussage | Tragende ADR |
|---|---|---|
| slice-001 | cmd/ + internal/-Baum | ADR-0039 |
| slice-001 | CGO-Ziel (Fortgeltung) | ADR-0039 (trägt ADR-0038 §CGO fort) |
| slice-001 | Store/Stream-Adapter später | ADR-0010 · ADR-0006 |
| slice-002 | Modelle + Invarianten in Konstruktoren | ADR-0004 · ADR-0029 |
| slice-002 | SourcePosition-Modell | ADR-0005 (ergänzt) · [`SPEC-003`](../../../spec/pflichtenheft.md) |
| slice-002 | ClockPort, `domain.TimePoint`, Fake | ADR-0040 |
| slice-003 | Ordnung Receive→…→ACK | [`ADR-0011`](README.md) · [`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md) |
| slice-003 | Orchestrierung im Application Layer | ADR-0027 |
| slice-003 | ChangeStorePort, ReplicationAckPort als Outbound | ADR-0009 · ADR-0007 |
| slice-003 | At-Least-Once / Crash-Wiederholung | ADR-0012 (ergänzt) |
| slice-003 | Fake Ports vor realem Adapter | ADR-0030 · ADR-0027 §Folgepflicht |