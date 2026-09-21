# Welle sdk-csharp-vollabdeckung: C#/NuGet-SDK auf volle Vier-Wege-Parität erweitern

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, direkt aus
[`ADR-0106`](../adr/0106-csharp-nuget-erstes-sdk-package.md) §Entscheidung
Festlegung 1, letzter Absatz geschnitten — SSE/NATS-Vollinhalt als
Folge-Package war dort bereits ausdrücklich antizipiert, keine neue ADR
nötig). **Datum:** 2026-09-21.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`PgChangeFeed.Client` (`SPEC-026`) deckt bislang HTTP-API (`SPEC-018`) und
gRPC-Stream (`SPEC-020`) — [`ADR-0106`](../adr/0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 1 grenzte v1 bewusst so ein, ließ SSE (`SPEC-021`) und
NATS-Vollinhalt (`SPEC-024`) aber ausdrücklich als „Folge-Package" offen —
„entweder als v2 desselben Packages oder als eigenständiges Package" (dort
wörtlich, Delegation an einen Folge-Zug ohne neue ADR). Diese Welle löst
diese Delegation ein: dasselbe Package (`sdks/csharp/PgChangeFeed.Client/`)
bekommt zwei weitere Client-Flächen, damit erreicht C# dieselbe
Vier-Wege-Matrix, die `examples/csharp/{http,grpc,sse,nats-stream}-client`
bereits real belegen.

Das *Mehr* gegenüber zwei isolierten Slice-DoDs:
[`LH-FA-SST-009`](../../../spec/lastenheft.md)s AC „Consumer bindet ein
Package über den Paketmanager ein, ohne das Protokoll selbst zu
implementieren" ist für die volle Matrix erst erfüllt, wenn SSE- **und**
NATS-Vollinhalts-Fläche zusammen mit der bereits bestehenden HTTP-/
gRPC-Fläche **in einem realen, neu paketierten** `.nupkg` stehen — ein
einzelner fertiger Flächen-Slice (z. B. nur SSE) belegt „volle Parität"
nicht.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf
erwähnt werden, aber nie Trigger sein.

- [`ADR-0106`](../adr/0106-csharp-nuget-erstes-sdk-package.md) ist
  `Accepted` (bereits erfüllt, 2026-09-19) und sein letzter
  Festlegung-1-Absatz erlaubt die Erweiterung ausdrücklich ohne neue ADR.
- `welle-sdk-csharp-lh-fa-sst-009` liegt in `done/` (bereits erfüllt,
  2026-09-19) — das Package existiert real und ist paketierbar
  (`PgChangeFeed.Client.0.1.0.nupkg`) und real auf NuGet.org veröffentlicht
  (`git tag -l "sdk-csharp-v*"` zeigt `sdk-csharp-v0.1.0`, seit der
  csharp-Welle-Closure als externe Betreiber-Handlung erfolgt).
- Kein weiterer Trigger nötig — die Welle kann sofort eröffnet werden.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen.

- Beide Slices in `done/`.
- `make gates` grün (netzloser Gate-Satz, unverändert Docker-only — dieser
  Baum bleibt außerhalb der Go-Coverage-Fläche und außerhalb von
  `.a-check.yml`).
- Ein real neu gebautes `.nupkg` (`make sdk-pack-csharp`) mit der
  gehobenen `<Version>` als Smoke-Beleg — nicht nur ein grüner
  `docker build`, sondern das Artefakt selbst existiert und trägt alle
  vier Client-Flächen.
- `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-026` tragen den
  Träger-Nachzug (`AGENTS.md` §3.13) — der Satz „deckt HTTP-API und
  gRPC-Stream" wird durch diese Welle falsch.
- Ein zweiter realer `sdk-csharp-v<Version>`-Tag-Push bleibt **außerhalb**
  dieser Welle (`AGENTS.md` §3.10, irreversible Betreiber-Handlung, wie
  bei der ersten csharp-Welle) — die Welle liefert bewusst nur das
  paketierbare Ergebnis, kein neuer Publish-Versuch.
- Closure-Notiz in `welle-sdk-csharp-vollabdeckung-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-csharp-sse-client-flaeche | SSE-Client-Fläche (`SPEC-021`) im bestehenden Package | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md) |
| slice-sdk-csharp-nats-stream-client-flaeche | NATS-Vollinhalts-Client-Fläche (`SPEC-024`), Version-Hebung, Träger-Nachzug | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md) |

**Reihenfolge:** Bewusst **sequentiell**, nicht parallelisierbar wie
HTTP/gRPC in der ersten csharp-Welle — beide Flächen berühren
unterschiedliche Dateien (`Sse/`, `Nats/`) und wären technisch unabhängig,
aber der Version-Bump und der Träger-Nachzug in
`spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md` sollen **einmal**,
für die volle Matrix zusammen, geschehen — nicht zweimal mit dem Risiko,
denselben Satz zweimal in unterschiedlichem Zwischenzustand umzuschreiben.
`slice-sdk-csharp-nats-stream-client-flaeche` läuft deshalb zuletzt und
trägt diese gebündelte Nacharbeit.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keiner anderen Welle;
  [`ADR-0106`](../adr/0106-csharp-nuget-erstes-sdk-package.md) ist bereits
  `Accepted`, keine neue ADR nötig.
- **Geschwister-Wellen, keine Abhängigkeit:** `welle-sdk-kotlin-vollabdeckung`
  und `welle-sdk-python-vollabdeckung` verfolgen dasselbe Ziel für ihre
  jeweilige Sprache — eigener Zählraum (`sdks/csharp/**` vs.
  `sdks/kotlin/**` vs. `sdks/python/**`), eigene Sub-Area, unabhängig
  voneinander priorisierbar und schließbar. Das WIP-Limit-1 gilt nur für
  `in-progress/` (je Slice, nicht je Welle) — drei flache, gleichzeitig
  offene Welle-Dateien sind kein Verstoß, solange zu jedem Zeitpunkt
  höchstens ein Slice in `in-progress/` liegt.
- Intern (siehe §4 Reihenfolge): `slice-sdk-csharp-nats-stream-client-flaeche`
  läuft nach `slice-sdk-csharp-sse-client-flaeche`.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1.

- **Aufspaltung in ein eigenständiges Folge-Package** — `ADR-0106`
  überlässt diese Struktur-Frage ausdrücklich „einem Folge-Zug"; diese
  Welle wählt die einfachere Variante (v2 desselben Packages, siehe §8 der
  Slices), ohne die Alternative (eigenständiges Package) neu
  aufzurollen — das bleibt ein künftiger, eigenständiger Entscheid, falls
  Nutzungsdaten dafür sprechen.
- **Ein realer NuGet.org-Tag-Push der neuen Version** — bleibt
  Betreiber-Entscheidung nach `AGENTS.md` §3.10, wie bei der ersten
  csharp-Welle.
- **Eine vierte SDK-Sprache oder ein zweiter C#-Vertriebsweg** —
  unverändert eine eigene, künftige ADR (`ADR-0106`
  §Re-Evaluierungs-Trigger 1).
- **Umbau oder Migration von `examples/csharp/`** — bleibt unverändert
  Vorbild, kein Belegträger (`SPEC-023`), unberührt von dieser Welle.
- **`spec/architecture.md`/`.a-check.yml`** bleiben unberührt —
  `sdks/csharp/**` liegt außerhalb der Server-Komponentensicht und wird
  von a-check nicht gelesen.
- **Kein neues Gate** — `make sdk-pack-csharp` bleibt Werkzeug.
- **Retrofit eines Integrationstests gegen eine reale, laufende
  Server-Instanz für die neuen Flächen** — anders als bei der
  Python-Welle (`ADR-0110` Folgepflicht 1, dort wegen der fehlenden
  Referenz-Vorarbeit verschärft) verlangt weder `ADR-0106` noch diese
  Welle das für C#: Die Vorarbeits-Parität (vier reale
  Beispiel-Clients gegen denselben Draht) bleibt bei C# unverändert
  bestehen, die HTTP-/gRPC-Fläche des Packages selbst wurde bereits ohne
  einen solchen Realserver-Test abgenommen (`slice-sdk-csharp-http-client-flaeche`,
  `slice-sdk-csharp-grpc-client-flaeche`) — dieselbe, bereits akzeptierte
  Teststrategie (Unit-Tests gegen Fakes) gilt für SSE/NATS-Vollinhalt
  fort, kein Sprung in eine strengere Klasse ohne fachlichen Grund.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
vollständig durchgesehen. Relevante Treffer:

- `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (offen, 2×,
  unter der Schwelle) — betrifft genau diesen Baum:
  `sdks/csharp/Dockerfile`s einzige `build`-Stufe macht
  `--build-context proto=proto` für **jeden** Bau zwingend, nicht nur für
  den gRPC-Pfad; `make sdk-pack-csharp` trägt den Flag bereits
  unconditional (`tools/harness/sdk-pack-csharp.sh`). Beide Flächen-Slices
  dieser Welle benennen das explizit in ihrem DoD-Wortlaut, statt „SSE
  braucht kein proto" zu behaupten — ein dritter Treffer würde die
  3×-Schwelle erreichen.
- `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (offen, 3×,
  Schwelle bereits erreicht, Embodiment ist Architect-Entscheidung, keine
  Planner-/Slice-Aufgabe) — gesichtet, kein Handlungsbedarf in dieser
  Welle; die nächste reguläre Architect-Sichtung bleibt zuständig.
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert) — trägt bereits den DoD-Punkt „Doku-Update" in beiden
  Slices, kein neuer Slice nötig.
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert seit
  `AGENTS.md` §3.13) — Suchlauf-Pflicht gilt für beide Slices
  unverändert, insbesondere gegen `sdks/csharp/README.md` §Status
  (bereits zweimal Fundort dieser Klasse in der ersten csharp-Welle).
- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
  (offen, 1×) — nicht einschlägig: diese Welle ändert den
  Publish-Mechanismus selbst nicht (`sdk-csharp-release.yml` bleibt
  unverändert, er validiert weiterhin Tag-gegen-`.csproj`-Version).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten werden erst bei Closure (unmittelbar vor dem `git mv`
nach `done/`) als auflösbare Links eingetragen. Diese Welle ist noch
offen — die Zeiger stehen deshalb als Platzhalter, keine Links.

Ergebnis: `welle-sdk-csharp-vollabdeckung-results.md` (noch nicht angelegt)
Zähler: `docs/plan/planning/observations/` (Beobachtungs-Register)
