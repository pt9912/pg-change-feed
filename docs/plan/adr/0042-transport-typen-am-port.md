# ADR-0042: Transport-Typen am Port

**Status:** Accepted — Supersedes [`ADR-0039`](0039-paketstruktur-detaillierung-go.md)

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`ADR-0039`](0039-paketstruktur-detaillierung-go.md) (abgelöst;
Struktur-Regeln gelten als Rest unverändert fort),
[`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md)
(Maschinenform der Fitness)

**Schärft:** [`ARC-003`](../../../spec/architecture.md),
[`ARC-005`](../../../spec/architecture.md) (Inbound-Ports und
Driving-Adapter der Sicht,
[architecture.md §1](../../../spec/architecture.md)) — die
Paketstruktur-Detaillierung aus
[`ADR-0039`](0039-paketstruktur-detaillierung-go.md) gilt als Rest
unverändert fort; diese ADR schärft nur die Definitions-Stelle der
Transport-Typen.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Review des Capture-Persist-Pfads (Inbound-Port, ChangeStore- und
Ack-Port, Use Case) meldete die Platzierung der Transport-Typen:
`CaptureCommand` und `CaptureResult` sind am Inbound-Port definiert
(`internal/application/port/inbound/capture.go`), das Use-Case-Paket
(`internal/application/usecase/capture/service.go`) führt sie als
Typ-Aliase unter den in der Baumskizze von
[`ADR-0039`](0039-paketstruktur-detaillierung-go.md) skizzierten Namen.
Die **Text-Hälfte** von
[`ADR-0039`](0039-paketstruktur-detaillierung-go.md) — „Jeder Use Case
trägt sein Service- und Transport-Typ-Paar (Command/Query, Result) im
eigenen Paket" — legt die Definition an den Use Case; die
**Fitness-Hälfte** derselben ADR („`internal/adapters/driving` importiert
nur `internal/application/port/inbound` und `internal/domain`") verlangt
sie effektiv am Port: Ein Command, das nur im Use-Case-Paket definiert
wäre, wäre für den Driving-Adapter nur über einen
`driving → usecase`-Import erreichbar, den keine
[`.a-check.yml`](../../../.a-check.yml)-Kante deklariert und den die
Schichten-Sicht ([`ARC-005`](../../../spec/architecture.md) darf
Application-Interna nicht importieren) ausschließt. Die zwei Hälften
des Accepted-ADR tragen also nicht zusammen; die implementierte
Lösung trägt beide (a-check 0 Befunde, keine neue Kante entstanden),
aber die Definitions-Stelle widerspricht der Wortlaut-Lesart — der
Defekt war der fehlende Verdikt, nicht der Code. Der Konflikt lief als
Konflikt-Sequenz (Modul 8) über den Architect; dieses Verdikt ist ihr
Übergabe-Artefakt.

## Entscheidung

Transport-Typen eines Ports — Command/Query und Result — werden **am
Port** definiert; der Use Case führt sie als **Typ-Aliase** im eigenen
Paket und benennt sie wie in der
[`ADR-0039`](0039-paketstruktur-detaillierung-go.md)-Baumskizze
(`CaptureService`, `CaptureCommand`, `CaptureResult`). Kontrakt-Typen
eines Outbound-Ports leben analog am Port. Damit tragen beide Hälften
wörtlich: die Fitness (der Driving-Adapter erreicht Command/Result über
`port/inbound`, ohne das Use-Case-Paket zu importieren — die
Fitness-Funktion bleibt ohne neue Kante erfüllt) und die
Use-Case-Lesart (das Use-Case-Paket führt die Typen unter ihren Namen,
als Aliase). Die übrigen Regeln von
[`ADR-0039`](0039-paketstruktur-detaillierung-go.md) — Paketbaum,
Domain-Events in `domain/event/`, Mapper-Verantwortung der Adapter,
Row-Typen im Driven-Adapter — gelten unverändert fort; die Struktur
wird nicht umgeplant. Wir wählen **Typen am Port, Aliase am Use Case**.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Typen ausschließlich im Use-Case-Paket (ADR-0039-Wortlaut wörtlich) | wörtliche Text-Treue; das Use-Case-Paket ist die einzige Import-Adresse | bricht die Fitness-Hälfte: der Driving-Adapter müsste `usecase/**` importieren — eine Kante, die `.a-check.yml` nicht deklariert und die Schichten-Sicht ausschließt |
| B — Typen in der Domain (`domain/model`) | Driving und Application importieren die Domain ohnehin; ein Ort für alle Typen | ließe `Result` tot: ein CaptureResult trägt Persist- und ACK-Ausgang, keine Domänen-Invariante — die Domain trüge Transport-Syntax, die der Use Case verbraucht |
| C — nichts tun; die Spannung als Kommentar tragen | kein neuer ADR, kein Textaufwand | die Spannung zwischen zwei Hälften eines Accepted-ADR bleibt stehen; jeder Leser löst sie neu, die Definitions-Stelle bleibt gegen den Wortlaut unentschieden |
| **D — Typen am Port, Typ-Aliase am Use Case (gewählt)** | beide Hälften tragen wörtlich; kein Umbau — die bestehende Platzierung ist bereits die Ziel-Form; a-check bleibt ohne neue Kante grün | zwei Namen für einen Typ (Alias-Indirektion); die Alias-Pflicht ist eine Konventions-Zusage, die a-check pfadgetrieben nicht prüft |

## Konsequenzen

- Positiv: Definitions-Stelle und Maschinenform stimmen überein — der
  Driving-Adapter erreicht die Transport-Typen über `port/inbound`, das
  Use-Case-Paket liest sie unter den skizzierten Namen; die
  Paketstruktur aus [`ADR-0039`](0039-paketstruktur-detaillierung-go.md)
  bleibt unverändert, keine Neu-Implementierung.
- Negativ: zwei Namen für denselben Typ — der Alias trägt die
  Use-Case-Lesart, die Port-Definition die Import-Kante; wird der Alias
  zur Neuedefinition, entstehen zwei Quellen für denselben Vertrag.
- Folgepflicht: Use-Case-Pakete **aliasieren** ihre Port-Typen statt sie
  neu zu definieren (Konventions-Zusage, Review-Prüfpflicht — a-check
  prüft Kanten, nicht Alias-Herkunft); die Maschinenform bleibt
  [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | `internal/application/usecase/**` referenziert Transport-Typen über `internal/application/port/inbound` (Alias, Kante `app → ports`); keine Kante `adapters → app`; `adapters/driving` bleibt auf `port/inbound` und `domain` beschränkt | `make a-check` (im Gate-Bündel) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbarer Trigger: **ein Port trägt einen Transport-Typ, den sein
Use Case nicht aliasiert** (mechanisch sichtbar als `grep`-Treffer gegen
`internal/application/usecase/**`, der eine Port-Signatur ohne Alias
führt) — dann ist die Alias-Pflicht gegen den tatsächlichen Bedarf zu
prüfen; die Definition zurück an den Use Case wäre nur mit einer neuen
Kante tragbar und wäre eine neue Entscheidung. Andernfalls `permanent` —
die Platzierung wächst mit den Ports.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted — Supersedes ADR-0039 (Transport-Typ-Platzierung am Port geschärft, übrige Struktur-Regeln unverändert fort; Anlass: Konflikt-Sequenz der Plan-Nachzug-Klasse, drittes Auftreten — Findings F-1/F-2 des Capture-Reviews) | [`ADR-0039`](0039-paketstruktur-detaillierung-go.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).