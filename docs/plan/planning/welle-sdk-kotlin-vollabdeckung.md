# Welle sdk-kotlin-vollabdeckung: Kotlin/GitHub-Packages-SDK auf volle Vier-Wege-Parität erweitern

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, direkt aus
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
§Entscheidung Festlegung 1, letzter Absatz geschnitten — SSE/NATS-Vollinhalt
als Folge-Package war dort ausdrücklich antizipiert, keine neue ADR nötig).
**Datum:** 2026-09-21.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`pgchangefeed-kotlin` (`SPEC-028`) deckt bislang HTTP-API (`SPEC-018`) und
gRPC-Stream (`SPEC-020`) —
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 1 grenzte v1 bewusst so ein („wie bei C#, nicht größer"), ließ
SSE (`SPEC-021`) und NATS-Vollinhalt (`SPEC-024`) aber ausdrücklich als
„Folge-Package" offen — dieselbe Formulierung, „jetzt zum dritten Mal
wiederholt", wie in `ADR-0106`. Diese Welle löst diese Delegation ein:
dasselbe Package (`sdks/kotlin/pgchangefeed-kotlin/`) bekommt zwei weitere
Client-Flächen, damit erreicht Kotlin dieselbe Vier-Wege-Matrix, die
`examples/kotlin/{http,grpc,sse,nats-stream}-client` bereits real belegen.

Das *Mehr* gegenüber zwei isolierten Slice-DoDs: dieselbe Begründung wie
bei `welle-sdk-csharp-vollabdeckung` §1 — „volle Parität" ist erst erfüllt,
wenn SSE- **und** NATS-Vollinhalts-Fläche zusammen mit der bereits
bestehenden HTTP-/gRPC-Fläche **in einem realen, neu paketierten** `.jar`
stehen.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln.

- [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  ist `Accepted` (bereits erfüllt, 2026-09-20) und sein letzter
  Festlegung-1-Absatz erlaubt die Erweiterung ausdrücklich ohne neue ADR.
- `welle-sdk-kotlin-lh-fa-sst-009` liegt in `done/` (bereits erfüllt,
  2026-09-21) — das Package existiert real und ist paketierbar
  (`pgchangefeed-kotlin-0.1.0.jar`); kein realer `sdk-kotlin-v*`-Tag wurde
  bislang gepusht (`git tag -l "sdk-kotlin-v*"` leer, Stand dieser
  Eröffnung).
- Kein weiterer Trigger nötig — die Welle kann sofort eröffnet werden.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

- Beide Slices in `done/`.
- `make gates` grün (netzloser Gate-Satz, unverändert Docker-only).
- Ein real neu gebautes `.jar` (`make sdk-pack-kotlin`) mit der gehobenen
  `version` als Smoke-Beleg — alle vier Client-Flächen im selben Artefakt.
- `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-028` tragen den
  Träger-Nachzug (`AGENTS.md` §3.13).
- Ein realer `sdk-kotlin-v<Version>`-Tag-Push bleibt **außerhalb** dieser
  Welle (`AGENTS.md` §3.10, irreversible Betreiber-Handlung).
- Closure-Notiz in `welle-sdk-kotlin-vollabdeckung-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine.

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-kotlin-sse-client-flaeche | SSE-Client-Fläche (`SPEC-021`) im bestehenden Package | [`LH-FA-SST-009`](../../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../../spec/lastenheft.md) |
| slice-sdk-kotlin-nats-stream-client-flaeche | NATS-Vollinhalts-Client-Fläche (`SPEC-024`), Version-Hebung, Träger-Nachzug | [`LH-FA-SST-009`](../../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../../spec/lastenheft.md) |

**Reihenfolge:** Sequentiell, dieselbe Begründung wie bei
`welle-sdk-csharp-vollabdeckung` §4 — der Version-Bump und der
Träger-Nachzug sollen einmal, für die volle Matrix zusammen, geschehen.
`slice-sdk-kotlin-nats-stream-client-flaeche` läuft zuletzt.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keiner anderen Welle;
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  ist bereits `Accepted`, keine neue ADR nötig.
- **Geschwister-Wellen, keine Abhängigkeit:** `welle-sdk-csharp-vollabdeckung`
  und `welle-sdk-python-vollabdeckung` verfolgen dasselbe Ziel für ihre
  jeweilige Sprache — eigener Zählraum (`sdks/kotlin/**`), eigene
  Sub-Area, unabhängig priorisierbar und schließbar. Das WIP-Limit-1 gilt
  nur für `in-progress/` (je Slice), nicht für die Zahl gleichzeitig
  offener, flacher Welle-Dateien.
- Intern (siehe §4 Reihenfolge): `slice-sdk-kotlin-nats-stream-client-flaeche`
  läuft nach `slice-sdk-kotlin-sse-client-flaeche`.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1.

- **Aufspaltung in ein eigenständiges Folge-Package** — `ADR-0109`
  überlässt diese Struktur-Frage ausdrücklich einem Folge-Zug; diese Welle
  wählt v2 desselben Packages.
- **Ein realer GitHub-Packages-Tag-Push der neuen Version** —
  Betreiber-Entscheidung nach `AGENTS.md` §3.10.
- **Eine vierte SDK-Sprache oder ein zweiter Kotlin-Vertriebsweg (z. B.
  Maven Central)** — unverändert eine eigene, künftige ADR (`ADR-0109`
  §Re-Evaluierungs-Trigger 1/5).
- **Ein Wechsel des JDK-Basis-Images (`eclipse-temurin:21-jdk` →
  `25-jdk`)** — `ADR-0109` Festlegung 5 bindet das ausdrücklich an einen
  künftigen Gradle-Versionssprung auf ≥ 9.1.0, der von dieser Welle nicht
  ausgelöst wird (§Re-Evaluierungs-Trigger 3).
- **Umbau oder Migration von `examples/kotlin/`** — bleibt unverändert
  Vorbild, kein Belegträger (`SPEC-023`).
- **`spec/architecture.md`/`.a-check.yml`** bleiben unberührt.
- **Kein neues Gate** — `make sdk-pack-kotlin` bleibt Werkzeug.
- **Retrofit eines Integrationstests gegen eine reale, laufende
  Server-Instanz für die neuen Flächen** — dieselbe Begründung wie bei
  `welle-sdk-csharp-vollabdeckung` §6: die Vorarbeits-Parität (fünf reale
  Kotlin-Beispiel-Clients, `ADR-0109` §Kontext) senkt das Wire-Risiko
  bereits; die bestehende HTTP-/gRPC-Fläche wurde ebenfalls ohne einen
  solchen Realserver-Test abgenommen, kein Sprung in eine strengere
  Klasse ohne fachlichen Grund. **Anders als bei Python** (`ADR-0110`)
  gibt es für Kotlin keine ADR, die eine verschärfte Test-Pflicht
  verlangt.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
vollständig durchgesehen. Relevante Treffer:

- `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` — **Nachtrag
  `AGENTS.md` §3.13, real geprüft bei Slice-Priorisierung
  `slice-sdk-kotlin-sse-client-flaeche`:** Diese Zeile trug bei
  Welle-Eröffnung (2026-09-21) „offen, 2×, unter der Schwelle" — dieselben
  zwei Belege `slice-102` (C#) und `slice-103` (Kotlin), beide
  `examples/**`. Die parallele Geschwister-Welle
  `welle-sdk-csharp-vollabdeckung` hat während ihres SSE-Slice
  (`slice-sdk-csharp-sse-client-flaeche`) inzwischen einen **dritten**
  Beleg erzeugt — Registerstand jetzt real **3×, Schwelle erreicht**
  (`../observations/BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut/state.md`).
  Der dritte Beleg trifft erstmals einen SDK-Package-Baum
  (`sdks/csharp/**`) statt `examples/**`, mit strukturell anderer Ursache
  (ein einziges Package/`.csproj` statt vermeidbar geteilter
  Docker-Stufen) als die ersten beiden — der Ausgang (ob/wie die Klasse
  künftig gefangen wird) bleibt eine offene Architect-Entscheidung, die
  bisherige Aussage „ein weiterer Kotlin-Treffer hier würde die
  3×-Schwelle erreichen" ist damit gegenstandslos: Die Schwelle ist
  bereits erreicht, unabhängig von einem Kotlin-Treffer. Betrifft
  `sdks/kotlin/Dockerfile`s einzige `build`-Stufe unverändert genauso wie
  bei C#. Beide Flächen-Slices dieser Welle benennen das weiterhin
  explizit; ein Kotlin-Treffer hier wäre ein weiterer Beleg derselben,
  bereits über der Schwelle liegenden Beobachtung, kein
  Schwellen-Übertritt mehr.
- `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (offen, 3×,
  Schwelle bereits erreicht — der Lese-Schritt der Kotlin-Erstwelle hat
  das bereits vermerkt, Embodiment bleibt Architect-Entscheidung) —
  gesichtet, kein Handlungsbedarf in dieser Welle.
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert) — DoD-Punkt „Doku-Update" liegt gebündelt im Folge-Slice.
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert seit
  `AGENTS.md` §3.13) — Suchlauf-Pflicht gilt unverändert.
- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
  (offen, 1×) — nicht einschlägig, kein neuer Publish-Mechanismus in
  dieser Welle.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/`
auflösen, nicht vom Schreibort.

Ergebnis: [welle-sdk-kotlin-vollabdeckung-results.md](welle-sdk-kotlin-vollabdeckung-results.md)
Zähler: [../observations/](../observations/)
