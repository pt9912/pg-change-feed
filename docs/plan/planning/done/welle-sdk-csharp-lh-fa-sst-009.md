# Welle sdk-csharp-lh-fa-sst-009: Erstes C#/NuGet-SDK-Package (`PgChangeFeed.Client`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, direkt aus
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md) §Konsequenzen
Folgepflicht 1 geschnitten). **Datum:** 2026-09-19.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md) (Accepted,
2026-09-19) entscheidet Form, Ort und Umfang des ersten offiziellen
Client-Bibliothek-Packages für
[`LH-FA-SST-009`](../../../../spec/lastenheft.md): C#/.NET, NuGet, ein Package
(`PgChangeFeed.Client`) mit HTTP-API (`SPEC-018`) und gRPC-Stream
(`SPEC-020`), in einem neuen Baum `sdks/csharp/`, mit eigenständigem
SemVer ab `0.1.0`, Docker-only-Bau (`make sdk-pack-csharp`) und einem
separaten Netz-Workflow für die NuGet.org-Veröffentlichung. Diese ADR
entscheidet ausdrücklich nur die *Form* — die Umsetzung ist Folgepflicht 1
dieser ADR, kein Vorgriff.

Das *Mehr* gegenüber fünf isolierten Slice-DoDs: `LH-FA-SST-009`s
Happy-Path-AC („Consumer bindet ein Package über den Paketmanager ein und
empfängt Changes, ohne das Protokoll selbst zu implementieren") ist erst
erfüllt, wenn Projektgerüst, HTTP-Fläche, gRPC-Fläche, Pack-Werkzeug und
Publish-Weg **zusammen** ein real baubares, real paketierbares Package
ergeben — ein einzelner fertiger Slice (z. B. nur das Projektgerüst) belegt
die Anforderung nicht; erst der volle Satz macht aus „Sprache und
Vertriebsweg entschieden" ein „Package existiert und ist konsumierbar".

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md) ist
  `Accepted` (bereits erfüllt, 2026-09-19).
- Kein weiterer Trigger nötig — die Roadmap führt aktuell keine offene
  Welle (`in-progress/roadmap.md` §Offene Wellen: „Nichts in Arbeit") und
  *Nächste Wellen* ist leer; die Welle kann sofort eröffnet werden.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle fünf Slices in `done/`.
- `make gates` grün (netzloser Gate-Satz, unverändert Docker-only).
- Ein real gebautes `.nupkg` (`make sdk-pack-csharp`) als Smoke-Beleg —
  nicht nur ein grüner `docker build`, sondern das Artefakt selbst
  existiert und trägt die erwartete `<Version>` (`0.1.0`).
- Der NuGet-Publish-Workflow wurde real gegen NuGet.org versucht (Tag
  `sdk-csharp-v0.1.0` oder gleichwertig) — Erfolg **oder** ein benannter
  roter Befund mit Folgemaßnahme; nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 bleibt dieses Risiko bis zum
  realen Post-Push-Lauf offen, unabhängig vom Rest der Welle.
- Closure-Notiz in `welle-sdk-csharp-lh-fa-sst-009-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-csharp-projektgeruest | SDK-Projektgerüst (`sdks/csharp/PgChangeFeed.Client/`, `.csproj`, Docker-Bau, leeres API-Skelett) | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |
| slice-sdk-csharp-http-client-flaeche | HTTP-Client-Fläche (`SPEC-018`, neun Port-gedeckte Fähigkeiten, Bearer-Auth, eigene Tests) | [`LH-FA-SST-009`](../../../../spec/lastenheft.md), [`LH-FA-SST-006`](../../../../spec/lastenheft.md) |
| slice-sdk-csharp-grpc-client-flaeche | gRPC-Client-Fläche (`SPEC-020`, `StreamChanges`, `.proto`-Bezug über Zusatzkontext, `authorization`-Metadata) | [`LH-FA-SST-009`](../../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../../spec/lastenheft.md) |
| slice-sdk-csharp-pack-werkzeug | `make sdk-pack-csharp`, Trägernachzug `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |
| slice-sdk-csharp-publish-workflow | `.github/workflows/sdk-csharp-release.yml`, Tag-Präfix `sdk-csharp-v*`, `NUGET_API_KEY` | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |

**Reihenfolge (Abhängigkeitskette):**

1. `slice-sdk-csharp-projektgeruest` zuerst — legt `sdks/csharp/PgChangeFeed.Client/`
   samt `.csproj`/Docker-Bau an; ohne dieses Gerüst haben die beiden
   Fläche-Slices keinen Ort, in den sie schreiben.
2. `slice-sdk-csharp-http-client-flaeche` und `slice-sdk-csharp-grpc-client-flaeche`
   — beide hängen nur vom Projektgerüst ab, nicht voneinander
   (unterschiedliche Dateien, unterschiedliche Draht-Verträge); sie sind
   **parallelisierbar**.
3. `slice-sdk-csharp-pack-werkzeug` danach — `make sdk-pack-csharp` baut,
   testet und paketiert **beide** Flächen in einem Package (`ADR-0106`
   Festlegung 1, ein Package); es braucht deshalb beide Vorgänger-Slices
   fertig.
4. `slice-sdk-csharp-publish-workflow` zuletzt — der Publish-Workflow
   validiert die Tag-Version gegen die `<Version>` des `.csproj` und
   veröffentlicht das vom Pack-Werkzeug erzeugte `.nupkg`-Muster; ohne ein
   bereits real paketierbares Package liefe er ins Leere.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle — die Roadmap führt aktuell keine
  parallele offene Welle.
- Wird blockiert von: keiner anderen Welle;
  [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md) ist bereits
  `Accepted`.
- Intern (siehe §4 Reihenfolge): `slice-sdk-csharp-http-client-flaeche` und
  `slice-sdk-csharp-grpc-client-flaeche` hängen beide von
  `slice-sdk-csharp-projektgeruest` ab; `slice-sdk-csharp-pack-werkzeug`
  hängt von beiden Flächen-Slices ab; `slice-sdk-csharp-publish-workflow`
  hängt vom Pack-Werkzeug-Slice ab.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **SSE- (`SPEC-021`) und NATS-Vollinhalts-Abdeckung (`SPEC-024`) im
  Package** — [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
  Festlegung 1 grenzt v1 ausdrücklich auf HTTP-API und gRPC-Stream ein;
  beide bleiben über `examples/csharp/sse-client`/`nats-stream-client` als
  Vorbild einsehbar, bis ein Folge-Package sie deckt
  (§Re-Evaluierungs-Trigger 1/2 der ADR).
- **Eine zweite SDK-Sprache** (Go, TypeScript, Kotlin/Maven, …) —
  [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md) entscheidet
  ausdrücklich nur die erste Sprache; eine zweite bleibt eine eigene,
  künftige ADR (§Re-Evaluierungs-Trigger 1).
- **Ein zweiter Vertriebsweg** für C# (z. B. ein privater Feed statt
  NuGet.org) — nicht Teil der ADR-Entscheidung.
- **Aufspaltung in mehrere Packages** (ein Package je Zustellweg) — von
  [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md) Alternative
  B4 ausdrücklich verworfen; bleibt Re-Evaluierungs-Trigger 3.
- **Umbau oder Migration von `examples/csharp/`** — die ADR verbietet das
  ausdrücklich (§Entscheidung Festlegung 2, §Entscheidung Festlegung 5);
  die Beispiele bleiben unverändert Vorbild, kein Belegträger (`SPEC-023`).
- **`spec/architecture.md`** bleibt unberührt — `sdks/**` ist außerhalb der
  Server-Komponentensicht (`ADR-0106` Festlegung 5).
- **`.a-check.yml`** bleibt unberührt — a-check liest kein C#
  (`ADR-0106` §Kontext Bindung 4).
- **Kein neues Gate** — `make sdk-pack-csharp` und der Publish-Workflow
  bleiben Werkzeuge, kein Gate-Eintrag in `GATE_CHECKS` (`ADR-0106`
  Festlegung 4).
- **`docs/user/version.md` (Server-SemVer)** bleibt unberührt — das SDK
  führt sein eigenes, unabhängiges SemVer (`ADR-0106` Festlegung 3).

**Sechster Slice — bewusst nicht geschnitten.** Ein separater
„Handbuch-Nachzug"-Slice (`ADR-0106` §Konsequenzen Folgepflicht 4,
`docs/user/benutzerhandbuch.md`-Hinweis je gedeckter Oberfläche) entfällt
als eigener Schnitt: Die zugrunde liegende Beobachtung
(`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
Ausgang bereits **verkörpert**) trägt die Pflicht bereits an zwei Stellen —
der diff-skopierten Selbstprüf-Instruktion in
`.claude/commands/implement-slice.md` Schritt 17 und dem eigenen HIGH-Punkt
in `.harness/skills/reviewer.md`. Der Handbuch-Hinweis je Oberfläche gehört
deshalb **in** `slice-sdk-csharp-http-client-flaeche` und
`slice-sdk-csharp-grpc-client-flaeche` selbst (DoD-Punkt „Doku-Update …
falls öffentlicher Vertrag berührt"), nicht in einen eigenen Slice — ein
zusätzlicher Slice würde dieselbe Pflicht doppelt tragen, ohne einen
eigenen Lieferwert zu haben.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
vollständig durchgesehen. Ein Treffer, der bereits 3× erreicht und noch
offen ist und diese Welle beträfe (Docker-only-Bau eines neuen
Sprach-Ökosystems, Netz-Workflow/Secret-Muster), existiert **nicht** — die
`grep`-Suche nach `Secret`/`NuGet`/`dotnet`/`examples-csharp`/
`examples-kotlin` über alle `observation.md`-Dateien blieb ohne Treffer.
Der einzige inhaltlich einschlägige Eintrag,
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`, ist
bereits **verkörpert** (siehe oben) und braucht keinen neuen Slice. Kein
Treffer ist damit selbst die Feststellung dieser Eröffnung, keine Auslassung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [welle-sdk-csharp-lh-fa-sst-009-results.md](welle-sdk-csharp-lh-fa-sst-009-results.md)
Zähler: [../observations/](../observations/)
