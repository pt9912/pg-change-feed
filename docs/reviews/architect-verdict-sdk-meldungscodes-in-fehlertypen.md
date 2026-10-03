# Architect-Verdikt — Meldungscodes in den Fehlertypen der drei SDK-Packages (Fragen A1 bis A3)

**Datum:** 2026-10-03 · **Stand:** `8b30f071` (Plan `sdk-meldungscodes-in-fehlertypen` in `open/`) ·
**Rolle:** Architect (frischer Kontext; Ursprung der Aussagen je Zeile genannt: *gelesen*,
*übernommen*, *hergeleitet*) · **Anlass:** die drei Architect-Fragen in §4 des Plans, bevor der Slice
von `open/` nach `next/` geht.

**Bezug:** [`LH-FA-SST-009`](../../spec/lastenheft.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md),
[`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) (neu, trägt die
Entscheidung), [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
[`AGENTS.md`](../../AGENTS.md) §3.5, §3.8, §3.12.

---

## 0. Ergebnis in einem Blick

| Frage | Verdikt |
|---|---|
| A1 Statusdetail | **Eigenlesung, keine neue Abhängigkeit** (Plan-Vorschlag bestätigt) |
| A2 API-Form | Name `MessageCode`/`messageCode`/`message_code` **bestätigt**; Konstruktoren: **Überladung statt optionalem Parameter** (C#, Kotlin `@JvmOverloads`), Python keyword-only — **weicht vom Plan ab**, Bruchfläche null, „Upgrading 0.6.0“ entfällt |
| A3 Umfang | **bestätigt**: nur HTTP, SSE (Öffnen) und gRPC-Administrations-Client; Streams und NATS außerhalb |
| ADR nötig? | **Ja** — [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) (öffentliche API-Form und Lieferketten-Entscheidung; ein Namenswechsel nach Veröffentlichung ist ein Bruch) |

## 1. A1 — Statusdetail ohne neue Abhängigkeit

**Entscheid:** je Sprache eine interne Hilfsfunktion liest `grpc-status-details-bin` mit der
vorhandenen Protobuf-Laufzeit (C# `CodedInputStream`, Kotlin `CodedInputStream`, Python `google.protobuf`
mit Varint-Leser).

**Begründung.** Drei Felder rechtfertigen keine neue Laufzeit-Abhängigkeit in einem
veröffentlichten Package: Version, Lizenz, Pin und Folgepflege je Sprache (Geist von
[`AGENTS.md`](../../AGENTS.md) §3.8), und drei verschiedene API-Formen (`CommonProtos`,
`grpcio-status`, `StatusProto`) machten die gemeinsame Tabelle des Plans zu drei Prüfungen. Auch
wenn Kotlin transitiv bereits eine fertige Klasse hätte (im Plan *nicht gemessen*), bleibt es bei
der Eigenlesung, damit die drei Parser dieselbe Tabelle prüfen.

**Wire-Feldnummern (*gelesen* am 2026-10-03 in einer lokalen Kopie der googleapis-Include-Protos eines
anderen Projekts, Version der Kopie nicht festgestellt — Protobuf-Feldnummern sind nicht änderbar):**
`google.rpc.Status`: 1 `code`, 2 `message`, 3 `repeated Any details`; `google.protobuf.Any`: 1
`type_url`, 2 `value`; `google.rpc.ErrorInfo`: 1 `reason`, 2 `domain`, 3 `metadata`. Der Trailer-Name
`grpc-status-details-bin` ist aus der gRPC-Spezifikation *übernommen* und wird im Realserver-Fall am
Server-Wire belegt. Der Plan nannte die Nummern „übernommen“; das ist jetzt *gelesen*, der
Realserver-Fall bleibt trotzdem Pflicht (Zusammenspiel, nicht Nummern).

**Zusatz gegenüber dem Plan:** der Parser wertet ein `Any` nur aus, wenn `type_url` auf
`google.rpc.ErrorInfo` endet (deckt Tabellenzeile 8 `RetryInfo`).

## 2. A2 — API-Form

**Name.** `MessageCode`, `messageCode`, `message_code` bestätigt. Die Kollision in Python ist real
(gelesen in `exceptions.py`: `PgChangeFeedGrpcError.code` ist `grpc.StatusCode`); `StatusCode`/`statusCode`
liegen in C#/Kotlin daneben. Ein anderer Name wäre nicht besser.

**Konstruktor — Abweichung vom Plan.** Die Blatt-Typen in C# haben einen *öffentlichen* Konstruktor
(gelesen); ein optionaler Parameter ändert dessen binäre Signatur. Eine **zusätzliche Überladung**
(alter Konstruktor bleibt, neuer mit `string? messageCode` hinten) erhält Binärkompatibilität in C#
vollständig; in Kotlin erzeugt `@JvmOverloads` am Primär-Konstruktor die alte JVM-Signatur weiter;
Python: `*, message_code=None` keyword-only. Kosten: zwei Konstruktoren je Typ in C# (gleiche Zeilenzahl
wie die Parameter-Variante in der Praxis, nur mechanischer). Folge: die Bruchfläche ist null, der
Eintrag „Upgrading 0.6.0“ in den READMEs von C# und Kotlin und der Neukompilierungs-Hinweis in der
Release-Notiz **entfallen**. Der Präzedenzfall 0.5.0 (README „Upgrading“) beweist nur, dass der Bruch
dort toleriert wurde, nicht dass er nötig ist. Hinweis für C#: an der Basis stehen dann
`(int, string, Exception)` und `(int, string, string?)` nebeneinander; ein `null`-Literal wäre
mehrdeutig — der Implementer nutzt benannte Argumente (`messageCode: null`) oder gar keinen Aufruf
der Überladung mit `null`.

**Semantik** (bestätigt, in der ADR festgelegt): leer = `null`/`None`; unverändert durchgereicht, keine
`PCF-…`-Formatprüfung; im HTTP-Körper zählt nur ein JSON-String (Zahl → leer, Text bleibt); gRPC erster
`ErrorInfo` mit `domain == "pg-change-feed"`; der Parser wirft nie; statusunabhängig (Zeile 12). Bei
`MalformedResponse` und NATS-Fehlern ist die Eigenschaft leer.

## 3. A3 — Umfang

Bestätigt: HTTP-Client, SSE-Client beim Öffnen und gRPC-Administrations-Client typisiert; Stream-Clients
und NATS bleiben roh bzw. außerhalb. README sagt, dass der Code an den Stream-Clients über die native
Statusdetail-API des gRPC-Clients lesbar bleibt (Aussage *hergeleitet*: der Server sendet `ErrorInfo`
bei `ChangeStream` ohne Broadcaster laut Plan; am Realserver nicht gemessen — kein Beleg verlangt, die
README-Zeile sagt „lesbar über die native API“, nicht mehr).

## 4. Ist eine ADR nötig?

Ja, [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md): öffentliche API-Form
über drei Packages und Wahl „keine Abhängigkeit“ sind Entscheidungen mit Re-Evaluierungs-Trigger. Der
Versionsschritt (additiv, Minor 0.6.0) ist dort als Festlegung 6 genannt. Kein Spec-Nachzug
(`SPEC-026` bis `SPEC-028` nennen keine Fehlertypen, Plan-Messung); `ADR-0144` bleibt unverändert, der
Satz „die SDKs ändern sich nicht“ ist dort nur eine Abgrenzung.

## 5. Herkunft der Plan-Aussagen (§3.12) — Unstimmigkeiten

| Plan-Stelle | Befund |
|---|---|
| §3 „Feldnummern … *übernommen*“ | jetzt *gelesen* ([`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Kontext 4); Plan-Satz nachziehen |
| §3 „Binär: … ändern die Konstruktor-Signaturen … *hergeleitet*“ und §6 „Binärkompatibilität“ | durch A2 hinfällig; Risiko entfällt |
| §3 „Kotlin: transitiv bereits im Klassenpfad … nicht gemessen“ | bleibt unerheblich (Eigenlesung in allen Sprachen) |
| §3 Realserver „*Erwartung, nicht gemessen:* `PCF-E8025`“ | korrekt als Erwartung gekennzeichnet; nicht beanstandet |
| §6 „Server alt/neu … am echten alten Server **nicht gelaufen**“ | korrekt gekennzeichnet |
| Gemessene Zahlen (47/28/13/0-Suchläufe, Handbuchzeile 1938, `sdk-public-doc-check`-Muster) | **von mir nicht nachgemessen**; der Reviewer misst am Diff, `make suchlauf-nachmessen` ist der Beleg |
| §3 Konstruktorform „letzter, optionaler Parameter“ (C#, Kotlin), DoD (A) „Abschnitt Upgrading 0.6.0“ | durch A2 geändert |

Keine Aussage steht als geprüft, die ungeprüft wäre.

## 6. Plan-Nachzugsliste (für den Planner; ich ändere den Plan nicht)

1. §1 Bezug und §4: `ADR-0145` nennen; A1 bis A3 als entschieden (Verdikt verlinken), „Rückführung neue
   Abhängigkeit“ entfällt.
2. §3 „API-additiv“: Überladung statt optionalem Parameter (C# zweiter Konstruktor je Typ inkl. Basis,
   Kotlin `@JvmOverloads`, Python keyword-only `message_code=None`); den Absatz zur Binärkompatibilität
   und „Wer die binäre Kompatibilität erhalten will …“ streichen; Hinweis auf die `null`-Mehrdeutigkeit
   an der C#-Basis.
3. §3 Statusdetail: `type_url` endet auf `google.rpc.ErrorInfo`; Feldnummern als *gelesen* ([`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Kontext 4).
4. DoD (A): „README (Fehlerbehandlung, Abschnitt Upgrading 0.6.0)“ → ohne Upgrading-Eintrag; DoD (B) analog.
   §3 Tabellenzeile `sdks/…/README.md`: „Eintrag unter Upgrading“ streichen.
5. §6: Risiken „Binärkompatibilität“ und „Neue Abhängigkeit“ als *entfallen* mit Verweis auf [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md);
   „Handgeschriebener Wire-Parser“ mit dem Realserver-Fall unverändert offen bis Closure.
6. §5: Release-Notiz ohne Neukompilierungs-Hinweis; der Folgeschritt 0.6.0 bleibt Hauptlauf.
7. DoD ergänzen: ein Test, dass der alte Konstruktor weiter existiert und `MessageCode` leer liefert
   (C# und Kotlin; Python: Aufruf ohne Keyword) — belegt die Binär-/Quellkompatibilität statt sie zu behaupten.
8. Mutationen: [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) nennt sie *hergeleitet*; der Bericht des Implementers nennt Stelle, Instanz, Farbe.

## 7. Empfehlungen an den Auftraggeber

- Der Slice kann nach dem Planner-Nachzug nach `next/`; es gibt keinen blockierenden Architect-Befund.
- Die Release-Freigabe 0.6.0 (drei externe Konten) ist weiter eine eigene Freigabe.
- Akzeptiertes Negativ ohne Folgevorgang: der handgeschriebene Wire-Leser (drei Felder); der Realserver-Fall
  und die vom Parser unabhängigen Unit-Fixtures (Plan §6) tragen ihn.
