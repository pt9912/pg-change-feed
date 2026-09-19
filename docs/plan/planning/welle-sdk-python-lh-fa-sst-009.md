# Welle sdk-python-lh-fa-sst-009: Erstes Python/PyPI-SDK-Package (`pgchangefeed`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, direkt aus
[`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md) §Konsequenzen
Folgepflicht 1 geschnitten). **Datum:** 2026-09-19.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md) (Accepted,
2026-09-19) entscheidet Form, Ort und Umfang des zweiten offiziellen
Client-Bibliothek-Packages für
[`LH-FA-SST-009`](../../../spec/lastenheft.md): Python, PyPI, ein Package
(Arbeitsname `pgchangefeed`) mit **ausschließlich** HTTP-API (`SPEC-018`,
inklusive Changes-Lesen `SPEC-022`) — anders als beim ersten (C#/NuGet)
Package **kein** gRPC in v1, in einem neuen Baum `sdks/python/`, mit
eigenständigem PEP-440-Versionsschema ab `0.x.y`, Docker-only-Bau
(`make sdk-pack-python`, `build`+`twine`) und einem separaten Netz-Workflow
für die PyPI-Veröffentlichung. Diese ADR entscheidet ausdrücklich nur die
*Form* — die Umsetzung ist Folgepflicht 1 dieser ADR, kein Vorgriff.

Das *Mehr* gegenüber vier isolierten Slice-DoDs: `LH-FA-SST-009`s
Happy-Path-AC („Consumer bindet ein Package über den Paketmanager ein und
empfängt Changes, ohne das Protokoll selbst zu implementieren") ist erst
erfüllt, wenn Projektgerüst, HTTP-Fläche, Pack-Werkzeug und Publish-Weg
**zusammen** ein real baubares, real paketierbares Package ergeben — ein
einzelner fertiger Slice (z. B. nur das Projektgerüst) belegt die
Anforderung nicht; erst der volle Satz macht aus „Sprache und Vertriebsweg
entschieden" ein „Package existiert und ist konsumierbar". Anders als bei
der C#-Welle ([welle-sdk-csharp-lh-fa-sst-009](done/welle-sdk-csharp-lh-fa-sst-009.md))
gibt es hier nur **eine** Draht-Fläche (HTTP), kein zweiter,
parallelisierbarer Flächen-Slice (`ADR-0107` §Entscheidung Festlegung 1 —
gRPC bleibt Folge-Release).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md) ist
  `Accepted` (bereits erfüllt, 2026-09-19).
- Kein weiterer Trigger nötig — die Roadmap führt aktuell keine offene
  Welle (`in-progress/roadmap.md` §Offene Wellen: „Nichts in Arbeit") und
  *Nächste Wellen* ist leer; die Welle kann sofort eröffnet werden.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle vier Slices in `done/`.
- `make gates` grün (netzloser Gate-Satz, unverändert Docker-only).
- Ein real gebautes `.whl`/`.tar.gz` (`make sdk-pack-python`) als
  Smoke-Beleg — nicht nur ein grüner `docker build`, sondern die
  Artefakte selbst existieren und tragen die erwartete Version (`0.1.0`
  oder gleichwertig, PEP 440).
- Der PyPI-Publish-Workflow wurde real gegen PyPI versucht (Tag
  `sdk-python-v0.1.0` oder gleichwertig) — Erfolg **oder** ein benannter
  roter Befund mit Folgemaßnahme; nach
  [`AGENTS.md`](../../../AGENTS.md) §3.10 bleibt dieses Risiko bis zum
  realen Post-Push-Lauf offen, unabhängig vom Rest der Welle.
- Closure-Notiz in `welle-sdk-python-lh-fa-sst-009-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-python-projektgeruest | SDK-Projektgerüst (`sdks/python/`, `pyproject.toml`, Docker-Bau, minimales Konfigurations-Skelett) | [`LH-FA-SST-009`](../../../spec/lastenheft.md) |
| slice-sdk-python-http-client-flaeche | HTTP-Client-Fläche (`SPEC-018`, neun Port-gedeckte Fähigkeiten + Changes-Lesen `SPEC-022`, Bearer-Auth, eigene Tests, **ohne** Referenz-Client) | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-006`](../../../spec/lastenheft.md) |
| slice-sdk-python-pack-werkzeug | `make sdk-pack-python`, Trägernachzug `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 | [`LH-FA-SST-009`](../../../spec/lastenheft.md) |
| slice-sdk-python-publish-workflow | `.github/workflows/sdk-python-release.yml`, Tag-Präfix `sdk-python-v*`, `PYPI_API_TOKEN`, `docs/user/releasing.md`-Nachzug im selben Slice | [`LH-FA-SST-009`](../../../spec/lastenheft.md) |

**Reihenfolge (rein sequenziell — kein Parallelisierungspotential, anders
als bei der C#-Welle, weil es nur eine Draht-Fläche gibt, `ADR-0107`
§Entscheidung Festlegung 1):**

1. `slice-sdk-python-projektgeruest` zuerst — legt `sdks/python/pgchangefeed/`
   samt `pyproject.toml`/Docker-Bau an; ohne dieses Gerüst hat der
   Fläche-Slice keinen Ort, in den er schreibt.
2. `slice-sdk-python-http-client-flaeche` danach — die einzige Draht-Fläche
   dieser Welle; hängt nur vom Projektgerüst ab. Kein zweiter,
   parallelisierbarer Flächen-Slice (Unterschied zur C#-Welle: dort liefen
   HTTP- und gRPC-Fläche parallel).
3. `slice-sdk-python-pack-werkzeug` danach — `make sdk-pack-python` baut,
   testet und paketiert die HTTP-Fläche; es braucht den Vorgänger-Slice
   fertig.
4. `slice-sdk-python-publish-workflow` zuletzt — der Publish-Workflow
   validiert die Tag-Version gegen die in `pyproject.toml` geführte
   Version und veröffentlicht das vom Pack-Werkzeug erzeugte
   `.whl`/`.tar.gz`-Paar; ohne ein bereits real paketierbares Package liefe
   er ins Leere.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle — die Roadmap führt aktuell keine
  parallele offene Welle.
- Wird blockiert von: keiner anderen Welle;
  [`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md) ist bereits
  `Accepted`.
- Intern (siehe §4 Reihenfolge): `slice-sdk-python-http-client-flaeche`
  hängt von `slice-sdk-python-projektgeruest` ab;
  `slice-sdk-python-pack-werkzeug` hängt von der Fläche ab;
  `slice-sdk-python-publish-workflow` hängt vom Pack-Werkzeug-Slice ab —
  eine reine Kette, keine Verzweigung.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **gRPC- (`SPEC-020`), SSE- (`SPEC-021`) und NATS-Vollinhalts-Abdeckung
  (`SPEC-024`) im Package** —
  [`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md) Festlegung 1
  grenzt v1 ausdrücklich auf HTTP-API ein — **anders als beim C#-Package**,
  das gRPC bereits in v1 trägt; Begründung ist die fehlende
  Python-Referenz-Client-Vorarbeit (§Kontext „Was das ändert" der ADR,
  siehe auch die Risiko-Benennung in
  `slice-sdk-python-http-client-flaeche`). gRPC/SSE/NATS bleiben für
  Python-Consumer über die entsprechenden Go-Server-Endpunkte direkt
  ansprechbar (Negative-AC von [`LH-FA-SST-009`](../../../spec/lastenheft.md)),
  bis ein Folge-Release sie deckt (§Re-Evaluierungs-Trigger 2/3 der ADR).
- **Eine dritte SDK-Sprache** (TypeScript/npm, Kotlin/Maven, Go, …) —
  [`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md) entscheidet
  ausdrücklich nur die zweite Sprache; eine dritte bleibt eine eigene,
  künftige ADR (§Re-Evaluierungs-Trigger 1).
- **Ein zweiter Vertriebsweg** für Python (z. B. Conda(-Forge)) — von der
  ADR ausdrücklich verworfen (§Verglichene Alternativen B2), nicht Teil
  dieser Welle.
- **Aufspaltung in mehrere Packages** (ein Package je Zustellweg) — von
  [`ADR-0107`](../adr/0107-python-pypi-zweites-sdk-package.md) Alternative
  C4 ausdrücklich verworfen.
- **Ein `examples/python/`-Referenz-Client** — existiert nicht und wird
  von dieser Welle nicht angelegt; die ADR stellt sein Fehlen fest
  (§Kontext Ist-Stand), verlangt aber nicht seine Nachlieferung. Entstünde
  er später doch, wäre das ein eigener Re-Evaluierungs-Anlass
  (§Re-Evaluierungs-Trigger 3 der ADR), nicht Teil dieser Welle.
- **Umbau, Zusammenlegung oder Migration von `sdks/csharp/**`** — diese
  ADR bleibt dafür ausdrücklich unberührt (`ADR-0107` §Bezug: „von ihr
  **unberührt**"); die beiden Sprach-SDKs bleiben getrennte Bäume mit
  eigenem Versionsraum.
- **`spec/architecture.md`** bleibt unberührt — `sdks/**` ist außerhalb
  der Server-Komponentensicht (`ADR-0107` §Entscheidung Festlegung 6).
- **`.a-check.yml`** bleibt unberührt — a-check liest kein Python
  (`ADR-0107` §Kontext Bindung 4).
- **Kein neues Gate** — `make sdk-pack-python` und der Publish-Workflow
  bleiben Werkzeuge, kein Gate-Eintrag in `GATE_CHECKS` (`ADR-0107`
  Festlegung 5, §Kontext Bindung 3).
- **`docs/user/version.md` (Server-SemVer) und der C#-SDK-Versionsraum**
  bleiben unberührt — das Python-Package führt sein eigenes, unabhängiges
  PEP-440-Schema (`ADR-0107` Festlegung 4).

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
vollständig durchgesehen — insbesondere die aus der C#-Welle bekannten
Einträge:

- `BEO-PGC/report-nackte-id-ohne-link` — bereits **verkörpert**
  (`AGENTS.md` §3.9), Zähler real ausgezählt bei 7× (`ls .../evidence | wc
  -l`); kein neuer Lese-Schritt fällig, bleibt als bereits verkörperte
  Regel für Reviewer/Verifier-Berichte dieser Welle unverändert gültig.
- `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` — offen,
  Zähler 2× (`slice-104`, `slice-sdk-csharp-projektgeruest`), **unter** der
  3×-Schwelle. Betrifft direkt die Situation dieser Welle (Docker-only-Bau
  eines neuen Sprach-Ökosystems in einer neuen Sub-Area) — als Risiko in
  `slice-sdk-python-projektgeruest` §6 aufgenommen: vor dem ersten
  Test-/Bau-Lauf einen Suchlauf im bereits bestehenden Bestand
  (`sdks/csharp/**`) auf ein bereits gelöstes, analoges Problem prüfen,
  statt ein Workaround-Muster neu zu erfinden. Erreicht mit dieser Welle
  **nicht** die 3×-Schwelle von selbst (Zähler wächst nur bei einem
  tatsächlichen dritten Auftreten, nicht durch bloße Erwähnung hier).
- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` — offen,
  Zähler 1× (`slice-sdk-csharp-publish-workflow`, F-1 dort: ein neuer
  Release-Mechanismus zog `docs/user/releasing.md` nicht automatisch nach).
  Direkt einschlägig: `slice-sdk-python-publish-workflow` trägt deshalb den
  `docs/user/releasing.md`-Nachzug (SDK-Release-Weg-Abschnitt, analog dem
  bestehenden C#-Eintrag) bereits **vorab als eigenen DoD-Punkt** — nicht
  erst nach einem Reviewer-Finding wie bei der C#-Welle. Ein drittes
  Auftreten dieser Klasse würde die 3×-Schwelle erreichen; diese Welle
  ist bewusst so geplant, dass es dazu **nicht** kommt.
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` —
  bereits **verkörpert** (Selbstprüf-Instruktion
  `.claude/commands/implement-slice.md` Schritt 17 + Reviewer-HIGH-Punkt);
  der Handbuch-Hinweis für die Python-HTTP-Oberfläche gehört deshalb als
  DoD-Punkt **in** `slice-sdk-python-http-client-flaeche` selbst, kein
  eigener Slice (dieselbe Begründung wie in der C#-Welle §6).

Kein weiterer Treffer (`grep`-Suche nach `Secret`/`PyPI`/`python`/`twine`/
`examples-csharp`/`examples-kotlin` über alle `observation.md`-Dateien
blieb ohne zusätzlichen Treffer über die vier oben genannten hinaus). Kein
Treffer ist damit selbst die Feststellung dieser Eröffnung, keine
Auslassung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-sdk-python-lh-fa-sst-009-results.md` (Geschwister im
Ruheort `done/`, noch nicht angelegt — erst bei Closure).
Zähler: `docs/plan/planning/observations/` (erst bei Closure als Link
einzutragen).
