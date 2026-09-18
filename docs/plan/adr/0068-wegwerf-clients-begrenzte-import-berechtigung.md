# ADR-0068: Wegwerf-Harness-Clients — begrenzte Import-Berechtigung auf den Adapter-Stub

**Status:** Accepted — Supersedes [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md)
(nur deren **Änderungs-Ausnahmeklausel**: den Satz „Änderungen an
Schichten-Globs und Edges sind Änderungen dieser Datei, keine neuen ADRs,
solange die §2-Constraints unverändert abgebildet bleiben" in §Entscheidung
und den ihn wiederholenden Schlusssatz im §Re-Evaluierungs-Trigger. Die
übrige [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md) —
a-check als Maschinenform, `GATE_CHECKS += a-check`, die Ablösung von
depguard, §Kontext, §Verglichene Alternativen, die Fitness-Function-Zeile und
der `time`-Import-Trigger — ist **hiermit bestätigt** und bleibt unverändert;
sie wird hier nicht wiederholt)

**Datum:** 2026-09-14

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Implementer-Lauf von `slice-071` <!-- d-check:status-provenance -->,
der `.a-check.yml` im selben Slice erweiterte, und als der Reviewer-Lauf, der
den Befund als Entscheidung weiterreichte; Modul 8 §Konflikt-Pfad)

**Bezug:** [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md)
(§Entscheidung, §Fitness Function, §Re-Evaluierungs-Trigger — in der
Ausnahmeklausel superseded), [`ADR-0026`](0026-composition-root.md)
(Composition Root), [`ADR-0002`](0002-abhaengigkeitsrichtung.md)
(Abhängigkeitsrichtung), [`spec/architecture.md §2`](../../../spec/architecture.md)
(Zeilen [`ARC-005`](../../../spec/architecture.md), [`ARC-007`](../../../spec/architecture.md)),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (der gRPC-Stream, dessen
Protokoll-Stub der Client konsumiert),
Review zu `slice-071` <!-- d-check:status-provenance --> (Anlass:
F-1, HIGH), `docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md` <!-- d-check:status-provenance -->
(§3 Plan-Nachzug), `.a-check.yml` (`composition_root`, `layers`, `edges`),
`tools/harness/grpcclient/main.go` (importiert genau ein Paket),
`test/integration/integration_test.go` (importiert sieben),
`harness/sensors/a-check.md`

**Schärft:** — (Prozess-ADR mit Architektur-Geltung, wie
[`ADR-0063`](0063-lh-fa-sch-003-testform-korrektur.md)/[`ADR-0064`](0064-lh-qa-ops-005-testansatz-korrektur.md)/[`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md);
sie ändert keine Zusage des Lastenhefts oder Pflichtenhefts und keinen Satz
der Sicht — sie entscheidet den Zuschnitt der Maschinenform der
§2-Constraints)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der E2E-Beleg für [`LH-FA-SST-008`](../../../spec/lastenheft.md) braucht einen
Wegwerf-Client, der den gRPC-Server-Stream real über das Netz anspricht
(`slice-071` <!-- d-check:status-provenance -->, [`ADR-0060`](0060-grpc-streaming-mechanismus.md)).
Der Client liegt in `tools/harness/grpcclient/` und importiert **genau ein**
Paket: den erzeugten Protokoll-Stub
`internal/adapters/driving/grpc/streamv1`. Ohne Eintrag in `.a-check.yml`
fällt der Client unter `wrong-direction: (ohne Schicht) -> adapters` — der
Baum war rot.

Der Implementer-Lauf hat daraufhin `tools/**` in `composition_root`
aufgenommen (`.a-check.yml`, Commit `b835dde`). Der Review dieses Slices
(Review zu `slice-071` <!-- d-check:status-provenance -->, F-1, HIGH)
hat die Aufnahme an drei Stellen als nicht gedeckt befundet:

1. **Die Rolle trägt nicht.** `composition_root` bedeutet in diesem Repo
   „die Verdrahtung konkreter Adapter" (`.a-check.yml`-Kommentar,
   [`ADR-0026`](0026-composition-root.md)). Der Client verdrahtet nichts; er
   konsumiert den Protokoll-Stub über das Netz. Seine Nachbarn
   `tools/harness/httpclient`/`natssub` importieren nicht einmal ein Paket
   unter `internal/`.
2. **Die Weite trägt nicht.** `tools/**` nimmt den **gesamten** Werkzeug-Baum
   auf — auch `tools/schema/**` (die d-migrate-Werkzeugkette), das der
   hinzugefügte Kommentar („die testseitigen Wegwerf-Clients des Harness")
   nicht meint.
3. **Die Begründung ist sachlich falsch.** Der Kommentar setzt den Client mit
   `test/integration/**` gleich („derselbe testseitige Verdrahtungs-Konsument");
   `test/integration/*.go` importiert **sieben** interne Pakete und verdrahtet
   real, der Client kein einziges.

**Was die Aufnahme zusätzlich bewirkt — real gemessen.** `composition_root`
ist die **unbeschränkteste** Rolle des Modells: ein Pfad dort darf *alles*
importieren. Der Review hat eine Probe-Datei unter `tools/harness/grpcclient/`
mit Import eines Produktions-Use-Cases laufen lassen — sie blieb grün. Dieser
Zug hat denselben Lauf reproduziert: mit `tools/**` in `composition_root`
Exit 0, **ohne jeden Hinweis**. Die Aufnahme tauscht damit ein sichtbares
„ungeprüft" (der Abdeckungs-Hinweis `harness/sensors/a-check.md` Grenze 1)
gegen ein stilles „deklariert".

**Was `ADR-0041` wörtlich über Änderungen sagt.** §Entscheidung: „Die
`.a-check.yml` ist ihr deklarativer Stand (Schichten und Edges,
`composition_root`); Änderungen an Schichten-Globs und Edges sind Änderungen
dieser Datei, keine neuen ADRs, solange die §2-Constraints unverändert
abgebildet bleiben." §Re-Evaluierungs-Trigger: „Änderungen an Schichten-Globs
und Edgen sind `.a-check.yml`-Änderungen, keine ADRs." `composition_root`
steht dort **zweimal in der Definition des Standes und zweimal außerhalb der
ADR-freien Änderungsklasse**. Die Aufnahme ist damit keine zulässige
Änderung des Standes allein.

**Und die Klausel greift an der anderen Seite ebenso wenig.** Sie nimmt
„Schichten-Globs und Edges" pauschal aus, ohne **Verfeinerung** von
**Erweiterung** zu trennen: das Verfeinern bestehender Layer-Globs
(`internal/**`-Unterteilung, neue Adapter-Pakete) bildet §2 unverändert ab;
das Aufnehmen eines Pfad-Bereichs, den §2 **nicht als Komponente führt**,
**mit einer Import-Berechtigung** erweitert die Architekturregel selbst. Beide
Formen stehen heute unter demselben ADR-freien Wortlaut. Diese Unterscheidung
ist der Gegenstand dieser ADR.

## Entscheidung

Wir wählen: **Die Import-Berechtigung der Wegwerf-Harness-Clients wird als
begrenzte Kante im Schichten-Modell ausgedrückt, nicht als unbeschränkte
Aufnahme in den `composition_root`. Und die ADR-freie Änderungsklasse der
`.a-check.yml` wird auf die Verfeinerung eingegrenzt.**

Vier Festlegungen:

1. **`composition_root` erhält keinen Werkzeug-Pfad.** `tools/**` und
   `tools/harness/**` bleiben draußen. `composition_root` trägt weiter genau
   die drei Verdrahtungs-Träger: `internal/bootstrap/**`, `cmd/**` und
   `test/integration/**`. Die Rolle des Schlüssels ist „verdrahtet konkret"
   ([`ADR-0026`](0026-composition-root.md), [`ARC-007`](../../../spec/architecture.md));
   ein Client, der nichts verdrahtet, gehört nicht hinein — unabhängig von
   der Glob-Weite.

2. **Die Clients bekommen eine eigene, begrenzte Gruppe.** `.a-check.yml`
   führt eine Schicht `tooling: ["tools/harness/**"]` **und** genau eine
   Kante `{from: tooling, to: adapters}`. Damit darf der Client den
   Protokoll-Stub importieren; ein Import aus `app`, `ports` oder `domain`
   in `tools/harness/**` bleibt `wrong-direction` (real belegt:
   `wrong-direction: tools -> app`, Exit 2). `tooling` ist **keine**
   Komponente der Sicht — es ist eine Gate-Scope-Gruppe, im Kommentar der
   Datei als solche benannt. Die drei `.go`-Dateien der Clients liegen damit
   in einer Schicht; der Abdeckungs-Hinweis „Dateien ohne Schicht" verstummt,
   weil nichts mehr ungedeckt ist — nicht, weil es stillgestellt wurde.

3. **Die ADR-freie Änderungsklasse wird geschärft.** ADR-frei sind
   `.a-check.yml`-Änderungen, die den **§2-Komponentensatz und seine
   Constraints unverändert abbilden** — das Verfeinern bestehender
   Layer-Globs und bestehender Edges. **Eine ADR braucht**, wer einen
   Pfad-Bereich, den `spec/architecture.md §2` nicht als Komponente führt,
   **mit einer Import-Berechtigung** in die Datei aufnimmt — über `layers`
   +`edges` **oder** über `composition_root`. Das ist eine Erweiterung der
   Architekturregel, keine Abbildung von §2; sie fällt unter
   [`AGENTS.md`](../../../AGENTS.md) §3.6 („Jede Schwellen-Senkung … ist ein
   ADR"). Diese Schärfung ersetzt die betreffende Klausel der
   [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md); der
   [ADR-Index](README.md) trägt die Nachfolge-Linie.

4. **Der Präzedenzfall `test/integration/**` wird bestätigt, nicht
   fortgeschrieben.** Seine Aufnahme war eine §2-Erweiterung, die vor dieser
   Klausel ohne ADR vollzogen wurde (`8d23352`). Diese ADR bestätigt sie als
   **zutreffend**: die Verdrahtungs-Tests tragen denselben Aufbau wie der
   Bootstrap und importieren dafür breit — sie verdrahten real. Damit ist die
   stille Fortschreibung aus `slice-007`<!-- d-check:status-provenance -->
   beendet: der Eintrag hat jetzt einen Träger.

### Was diese ADR nicht ändert

- **`internal/**` bleibt unberührt.** Kein Code-Eingriff; der Client, die
  Clients-Nachbarn und der gRPC-Adapter bleiben, wo sie sind. Der Client
  importiert den Stub weiterhin — die Aufgabe „Import des erzeugten Stubs
  vermeiden" (Handrollen der Protobuf-Form) ist damit nicht gestellt.
- **`test/integration/**` bleibt `composition_root`-Mitglied** (Festlegung 4).
- **`ADR-0041`s Kern gilt unverändert:** a-check ist die Maschinenform der
  §2-Constraints, hängt an `GATE_CHECKS`, und depguard bleibt draußen.
- **Keine neue Abdeckungs- oder Auflösungs-Ausnahme.** Die Schärfung macht
  die Änderungsklasse **enger**, nicht weiter: sie nimmt `composition_root`
  ausdrücklich aus der ADR-freien Klasse, statt ihn (wie in einer der
  geprüften Alternativen) hineinzuziehen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: `tools/**` bleibt in `composition_root`, der Kommentar bleibt | kein Eingriff; der E2E-Beleg läuft grün | die Aufnahme hat in `ADR-0041` keine Deckung und fällt unter [`AGENTS.md`](../../../AGENTS.md) §3.6; dem **ganzen** Werkzeug-Baum ist unbeschränktes Importrecht gegeben (Probe: Produktions-Use-Case-Import in `tools/harness/grpcclient/` bleibt grün); der Kommentar zitiert eine Aussage, die dort nicht steht — genau die stille Gate-Lockerung, die der Review als HIGH führt |
| B — `composition_root` auf `tools/harness/**` verengen | kleinster Eingriff; deckt sich mit dem Kommentar-Wortlaut | verengt den **Baum**, nicht das **Recht**: die Use-Case-Probe bleibt grün (der unbeschränkte Import ist die Eigenschaft des Schlüssels, nicht der Glob-Weite); die Rolle „verdrahtet konkret" trägt weiter nicht; `composition_root` verliert für jeden Leser seine Bedeutung (Test-Verdrahtung *und* Werkzeug-Zugriff in einer Liste) |
| C — `composition_root` behalten, Glob auf den konkreten Client (`tools/harness/grpcclient/**`) schneiden | minimale Blast-Radius-Reduktion | Per-Client-Ausnahme proliferiert mit jedem weiteren Client; die Rolle bleibt falsch und das Recht unbeschränkt; behebt keinen der drei Review-Punkte |
| D — `tools/**` zurücknehmen; den Client unter `test/integration/**` verschieben | kein `.a-check.yml`-Eingriff; nutzt den bereits bestätigten `composition_root`-Eintrag | versteckt den Adapter-Import unter einer bestehenden Ausnahme, statt ihn zu benennen — dieselbe stille Lockerung in anderer Form; bricht die Ablage-Konsistenz mit `tools/harness/httpclient`/`natssub` (alle Wegwerf-Clients an einem Ort); `test/integration/` ist der Ort der Go-**Tests**, nicht eines ausführbaren Client-Binaries |
| E — Schicht `tooling: ["tools/harness/**"]` + Kante `{from: tooling, to: adapters}` **(gewählt)** | begrenzt das Recht auf `adapters` statt „alles"; die Use-Case-Probe wird rot (`wrong-direction: tools -> app`, real belegt); die Abhängigkeit steht als **Kante** im Modell, nicht als Ausnahme davon; alle `.go`-Dateien liegen in einer Schicht (kein „ungeprüft"-Hinweis mehr) | führt eine Gruppe in `layers` ein, die keine Komponente der Sicht ist (im Kommentar als Gate-Scope benannt); die Kante ist grob — sie erlaubt die ganze `adapters`-Schicht, nicht nur den Stub (benannte Grenze, s. §Konsequenzen) |
| F — den generierten Stub aus `adapters` herauslösen (eigener Vertrags-Bereich, z. B. `internal/contracts/**`) | trennt Protokoll-Vertrag von Adapter-Implementierung — die genaueste Antwort auf die Ursache | größerer Umbau: Generierungs-Pfade, Server-Importe und die `adapters`-Globs ändern sich; eigener Slice mit eigener Spec-Klärung; für den akuten Befund unverhältnismäßig — diese ADR hält die Tür offen (Re-Evaluierungs-Trigger) |

**Zur Alternative D im Verhältnis zu E:** D wäre konform und billig, verlagert
das Problem aber in die Ablage statt in die Entscheidung. Der Client **ist**
ein Test-Werkzeug und gehört zu seinen Nachbarn nach `tools/harness/`; die
Frage, die D umgeht — „darf ein Werkzeug außerhalb des Hexagons einen Adapter
importieren, und wie eng?" — ist genau die, die eine Entscheidung braucht.

## Konsequenzen

- Positiv: Die Import-Berechtigung ist **begrenzt** — nur `adapters`, kein
  `app`/`ports`/`domain`. Der Review-Befund „unbeschränktes Importrecht für
  den ganzen Baum" ist ausgeräumt und **rot-fähig** belegt (Probe).
- Positiv: `composition_root` behält seine eine Bedeutung (Verdrahtung) und
  seine drei Träger; ein Leser muss nicht zwischen Test-Verdrahtung und
  Werkzeug-Zugriff unterscheiden.
- Positiv: Die Klausel unterscheidet jetzt **Verfeinerung** (ADR-frei) von
  **Erweiterung** (ADR-pflichtig). Der Präzedenzfall `test/integration/**`
  ist bestätigt, nicht länger still fortgeschrieben.
- Positiv: **kein Code-Eingriff, kein Umzug** — der Client bleibt an seinem
  Ort und behält den Stub-Import.
- Negativ: `.a-check.yml` liest sich um eine Gruppe schwerer: `layers` führt
  mit `tooling` einen Eintrag, der keine Komponente der Sicht ist. Benannt
  im Kommentar der Datei und hier — die Sicht selbst (`spec/architecture.md`)
  bleibt unberührt und trägt keine `tooling`-Zeile.
- Negativ mit Grenze: Die Kante `tooling → adapters` ist **grob** — sie
  erlaubt die ganze `adapters`-Schicht, nicht nur den Protokoll-Stub
  (`streamv1`). Ein Werkzeug, das einen *driven* Adapter importiert, bliebe
  grün. Feiner ginge es nur, indem man `adapters` in Unter-Schichten
  zerlegt (fragmentiert §2) oder den Stub herauslöst (Option F, eigener
  Slice). Grenze hier benannt, nicht still.
- Folgepflicht: `.a-check.yml` trägt die Rücknahme von `tools/**`, die
  Schicht `tooling` + Kante und den korrigierten Kommentar (dieser Zug).
  `harness/sensors/a-check.md` zieht die Bindung und die Grenze-1-Aussage
  nach (dieser Zug).
- Folgepflicht: `slice-071` <!-- d-check:status-provenance --> (Implementer-/
  Planner-Zug) — §3 (Tabellenzeile `.a-check.yml`) und der Absatz
  *Implementer-Entscheidungen* verweisen auf die zurückgenommene
  `composition_root`-Erweiterung und die falsche `test/integration`-Gleichsetzung;
  sie werden auf diese ADR und die Kante `tooling → adapters` nachgezogen.
  Der Slice-Plan ist Planner-/Implementer-Arbeit; diese ADR benennt den
  Nachzug, ändert ihn nicht.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check (Modulzeile `wrong-direction`) | Die Gruppe `tooling` (`tools/harness/**`) importiert **nur** `adapters`: ein Import aus `app`, `ports` oder `domain` in `tools/harness/**` ist `wrong-direction` (real belegt: `wrong-direction: tools -> app` beim Import `internal/application/usecase/capture`, Exit 2). Zugleich bleibt der Stub-Import (`internal/adapters/driving/grpc/streamv1`) grün | `make a-check` (im Gate-Bündel) |
| Review-Prüfpflicht (nicht maschinell) | `composition_root` führt keinen Pfad außerhalb der drei Verdrahtungs-Träger `internal/bootstrap/**`, `cmd/**`, `test/integration/**`. a-check prüft den **Inhalt** des Schlüssels nicht gegen eine feste Liste — eine neue Werkzeug-Aufnahme dort ist an der Form erkennbar (Review), nicht am Lauf | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Zwei benannte Trigger, sonst permanent:

1. **Ein zweiter Werkzeug-Baum außerhalb `tools/harness/**` braucht denselben
   Stub-Zugriff.** Dann entscheidet eine Folge-ADR, ob die Gruppe `tooling`
   auf einen weiteren Glob gezogen wird oder ob der neue Bereich seine eigene
   Kante braucht — die Kante wird nicht stillschweigend erweitert.
2. **Der generierte Protokoll-Stub wandert aus `internal/adapters/**` heraus**
   (etwa in einen eigenen Vertrags-Bereich, Option F). Dann wird die Kante auf
   diesen Bereich umgestellt; die Berechtigung wird feiner, nicht breiter.

Sonst permanent — die Rolle „testseitiger Protokoll-Konsument, der
`adapters` erreicht und sonst nichts" gilt unabhängig von der Zahl der
Werkzeug-Clients.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: Review `slice-071`<!-- d-check:status-provenance --> F-1 (HIGH) belegt, dass die Aufnahme `tools/**` in `composition_root` in [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md) keine Deckung hat und dem Werkzeug-Baum unbeschränktes Importrecht gibt. Unabhängiger Architect-Zug nimmt die Zeile zurück, führt die begrenzte Gruppe `tooling` + Kante `tooling → adapters` ein und schärft die ADR-freie Änderungsklasse (Modul 8 §Konflikt-Pfad: Lockerung legitim, aber falsch zugeschnitten → ADR, die den Gegenstand der Ausnahme benennt) | Review zu `slice-071` |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | PENDING_COMMIT |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0068` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
