# Welle sdk-python-vollabdeckung: Python/PyPI-SDK auf volle Vier-Wege-Parität erweitern (`ADR-0110`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, direkt aus
[`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Konsequenzen Folgepflicht 1 geschnitten — **anders als bei C#/Kotlin
braucht Python diese neue ADR**, weil `ADR-0107` den kleineren Erst-Scope
**explizit wegen** einer Risiko-Asymmetrie wählte, siehe `ADR-0110`
§Kontext). **Datum:** 2026-09-21.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`pgchangefeed` (`SPEC-027`) deckt bislang ausschließlich HTTP-API
(`SPEC-018`) — [`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md)
Festlegung 1 grenzte v1 bewusst darauf ein, **wegen** der fehlenden
Python-Referenz-Vorarbeit (kein `examples/python/`). [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Accepted, 2026-09-21) entscheidet: Das Restrisiko bleibt real bestehen,
wird aber jetzt bewusst getragen — ein Folge-Release deckt alle drei
verbleibenden Wege (gRPC `SPEC-020`, SSE `SPEC-021`, NATS-Vollinhalt
`SPEC-024`) **in einem Zug**, mit einer **verschärften Test-Pflicht**:
jede neue Zustellweg-Fläche prüft ihre Protokoll-Annahmen zusätzlich zu
Unit-Tests gegen eine reale, laufende Server-Instanz (`ADR-0110`
§Entscheidung Festlegung 2, Folgepflicht 1) — eine strengere Auflage als
bei den bisherigen SDK-Slices aller drei Sprachen (C#/Kotlin/Python-HTTP
testeten bislang ausschließlich netzlos gegen Fakes).

Das *Mehr* gegenüber drei isolierten Slice-DoDs:
[`LH-FA-SST-009`](../../../spec/lastenheft.md)s AC ist für die volle
Matrix erst erfüllt, wenn gRPC-, SSE- **und** NATS-Vollinhalts-Fläche
zusammen mit der bestehenden HTTP-Fläche in einem realen, neu
paketierten Artefakt-Paar (`.whl`+`.tar.gz`) stehen **und** jede der drei
neuen Flächen einen realen Server-Rundlauf-Beleg trägt — ein einzelner
fertiger Flächen-Slice ohne diesen Beleg erfüllt `ADR-0110`s
Folgepflicht 1 nicht.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln.

- [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) ist
  `Accepted` (bereits erfüllt, 2026-09-21).
- `welle-sdk-python-lh-fa-sst-009` liegt in `done/` (bereits erfüllt,
  2026-09-19) — das Package existiert real und ist paketierbar
  (`pgchangefeed-0.1.0-py3-none-any.whl`, `pgchangefeed-0.1.0.tar.gz`) und
  real auf PyPI-Release-Pfad getaggt (`git tag -l "sdk-python-v*"` zeigt
  `sdk-python-v0.1.0`, `ADR-0110` §Kontext Ist-Stand).
- Kein weiterer Trigger nötig — die Welle kann sofort eröffnet werden.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

- Alle drei Slices in `done/`.
- `make gates` grün (netzloser Gate-Satz, unverändert Docker-only).
- Ein real neu gebautes Artefakt-Paar (`make sdk-pack-python`) mit der
  gehobenen Version als Smoke-Beleg — alle vier Client-Flächen im selben
  Artefakt.
- **Für jede der drei neuen Flächen** (gRPC, SSE, NATS-Vollinhalt) ein
  realer, grüner Lauf des in `slice-sdk-python-grpc-client-flaeche`
  eingeführten Integrationstest-Werkzeugs (`make test-sdk-python-integration`)
  gegen eine reale, laufende Server-Instanz — `ADR-0110` §Entscheidung
  Festlegung 2/Folgepflicht 1, kein Gate, aber Pflichtbeleg dieser Welle.
- `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-027` tragen den
  Träger-Nachzug (`AGENTS.md` §3.13) — der Satz „deckt HTTP-API" wird
  durch diese Welle falsch.
- Ein realer `sdk-python-v<Version>`-Tag-Push bleibt **außerhalb** dieser
  Welle (`AGENTS.md` §3.10, irreversible Betreiber-Handlung).
- Closure-Notiz in `welle-sdk-python-vollabdeckung-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine.

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-python-grpc-client-flaeche | gRPC-Client-Fläche (`SPEC-020`); führt das Realserver-Integrationstest-Werkzeug ein | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) |
| slice-sdk-python-sse-client-flaeche | SSE-Client-Fläche (`SPEC-021`); erweitert das Integrationstest-Werkzeug | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) |
| slice-sdk-python-nats-stream-client-flaeche | NATS-Vollinhalts-Client-Fläche (`SPEC-024`); erweitert das Werkzeug; Version-Hebung, Träger-Nachzug | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) |

**Reihenfolge:** Sequentiell, in dieser Reihenfolge (gRPC → SSE →
NATS-Vollinhalt) — dieselbe Reihenfolge, die `ADR-0107` Festlegung 1
bereits für gRPC allein antizipiert hatte. Zwei Gründe für die strikte
Sequenz, stärker als bei C#/Kotlin:

1. Das Realserver-Integrationstest-Werkzeug (`tools/harness/run-sdk-python-integration-tests.sh`,
   `make test-sdk-python-integration`) entsteht **einmal**, im ersten
   Slice (gRPC) — SSE und NATS-Vollinhalt **erweitern** es um ihre eigene
   Testdatei, statt es unabhängig neu zu bauen. Eine parallele Bearbeitung
   würde zwei gleichzeitige Änderungen an derselben neuen
   Infrastruktur-Datei riskieren.
2. Wie bei C#/Kotlin sollen Version-Bump und Träger-Nachzug einmal, für
   die volle Matrix zusammen, geschehen — gebündelt im letzten Slice
   (NATS-Vollinhalt).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keiner anderen Welle;
  [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) ist
  bereits `Accepted`.
- **Geschwister-Wellen, keine Abhängigkeit:** `welle-sdk-csharp-vollabdeckung`
  und `welle-sdk-kotlin-vollabdeckung` verfolgen dasselbe Ziel für ihre
  jeweilige Sprache — eigener Zählraum (`sdks/python/**`), eigene
  Sub-Area, unabhängig priorisierbar und schließbar. Das WIP-Limit-1 gilt
  nur für `in-progress/` (je Slice), nicht für die Zahl gleichzeitig
  offener, flacher Welle-Dateien.
- Intern (siehe §4 Reihenfolge): strikt sequentiell, gRPC → SSE → NATS.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1.

- **Ein `examples/python/`-Referenz-Client** — `ADR-0110` §Entscheidung
  Festlegung 2 mandatiert ausdrücklich **keinen** neuen
  Beispiel-Client-Baum; das wäre eine andere, größere Entscheidung
  (`SPEC-023`s Umfang) ohne fachlichen Zusammenhang zu dieser Welle.
- **Aufspaltung in ein eigenständiges Folge-Package** — `ADR-0110`
  §Entscheidung Festlegung 3 überlässt das ausdrücklich einem Folge-Zug;
  diese Welle wählt v2 desselben Packages.
- **Ein realer PyPI-Tag-Push der neuen Version** — Betreiber-Entscheidung
  nach `AGENTS.md` §3.10.
- **Eine dritte SDK-Sprache über C#/Python/Kotlin hinaus, oder ein
  zweiter Python-Vertriebsweg** — unverändert `ADR-0107`
  §Re-Evaluierungs-Trigger 1.
- **Ein Retrofit desselben Realserver-Integrationstests für die
  bestehende HTTP-Fläche** — `ADR-0110` §Entscheidung Festlegung 2/
  Folgepflicht 1 verlangt die verschärfte Test-Pflicht ausdrücklich nur
  für die **neuen** Flächen (gRPC/SSE/NATS); die HTTP-Fläche bleibt bei
  ihrer bereits abgenommenen Teststrategie (Unit-Tests gegen
  `httpx.MockTransport`).
- **Ein neues Gate** — `make test-sdk-python-integration` bleibt
  Werkzeug wie `make test-integration` selbst (braucht DB-Zugang/Docker/
  Netz, kein netzloser Lauf).
- **Trusted Publishing (PyPI OIDC) statt `PYPI_API_TOKEN`** — unverändert
  `ADR-0107` §Re-Evaluierungs-Trigger 4.
- **`spec/architecture.md`/`.a-check.yml`** bleiben unberührt.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
vollständig durchgesehen. Relevante Treffer:

- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert seit
  `AGENTS.md` §3.13, 19×) — dieser Eintrag trägt bereits einen echten
  Python-Beleg derselben Sub-Area
  (`evidence/slice-sdk-python-http-client-flaeche.md`: ein Träger-Nachzug
  übersah die `pyproject.toml` und eine deutsche Formulierung „folgt
  erst") — der Suchlauf in jedem der drei Slices dieser Welle prüft
  explizit auch `sdks/python/pgchangefeed/src/pgchangefeed/options.py`s
  Docstring („no gRPC/SSE/NATS surface exists yet"), der durch diese
  Welle Slice für Slice falsch wird.
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert) — DoD-Punkt „Doku-Update" liegt gebündelt im letzten
  Flächen-Slice.
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, unter der
  Schwelle) — betrifft `run-integration-tests.sh`s `-run`-Musterabdeckung
  für `test/integration/integration_test.go`; **nicht** einschlägig für
  das hier neu eingeführte, eigenständige
  `run-sdk-python-integration-tests.sh` (anderes Skript, andere
  Testquelle) — explizit geprüft und verworfen, kein dritter Treffer.
- `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (offen, 3×,
  Architect-Entscheidung) — gesichtet, kein Handlungsbedarf hier.
- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
  (offen, 1×) — nicht einschlägig, `sdk-python-release.yml` bleibt
  unverändert.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**.
Die beiden Zeiger unten sind als auflösbare Links eingetragen: der
Ergebnis-Zeiger trägt die Ruheort-Datei in `done/` direkt, der
Zähler-Zeiger das Register; der `git mv` der Welle-Datei nach `done/`
folgt als eigener Commit (`AGENTS.md` §3.3) und der
Reconciliations-Commit danach stellt beide auf die Ruheort-Tiefe.

Ergebnis: [welle-sdk-python-vollabdeckung-results.md](done/welle-sdk-python-vollabdeckung-results.md)
Zähler: [observations/](observations/README.md) (Beobachtungs-Register)
