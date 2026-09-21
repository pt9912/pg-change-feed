# Welle sdk-kotlin-lh-fa-sst-009: Drittes SDK-Package (Kotlin/GitHub Packages, `pgchangefeed-kotlin`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, direkt aus
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
§Konsequenzen Folgepflicht 1 geschnitten). **Datum:** 2026-09-20.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
(Accepted, 2026-09-20) entscheidet Form, Ort und Umfang des dritten
offiziellen Client-Bibliothek-Packages für
[`LH-FA-SST-009`](../../../../spec/lastenheft.md): Kotlin, GitHub Packages
(Gradle-/Maven-Registry dieses Repositories), ein Package (Arbeitsname
`pgchangefeed-kotlin`, Koordinate `io.github.pt9912:pgchangefeed-kotlin`) mit
HTTP-API (`SPEC-018`) und gRPC-Stream (`SPEC-020`) — wie beim C#-Package,
**nicht** größer, obwohl für Kotlin bereits alle fünf Zustellwege ein
funktionierendes Beispielprogramm haben (`ADR-0109` §Kontext „Was das
ändert"). Neuer Baum `sdks/kotlin/pgchangefeed-kotlin/`, eigenständiges
SemVer ab `0.x.y`, Docker-only-Bau (`make sdk-pack-kotlin`, Gradle/JDK 21) und
ein separater Netz-Workflow für die GitHub-Packages-Veröffentlichung —
**ohne** externes Repository-Secret (`GITHUB_TOKEN` mit `packages: write`
genügt, der zentrale Unterschied zu den beiden vorigen SDK-Wellen). Diese ADR
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
Dieselbe Fünf-Slice-Form wie bei der C#-Welle
([welle-sdk-csharp-lh-fa-sst-009](welle-sdk-csharp-lh-fa-sst-009.md)) —
nicht die vier Slices der Python-Welle
([welle-sdk-python-lh-fa-sst-009](welle-sdk-python-lh-fa-sst-009.md)),
weil v1 hier wie bei C# **zwei** Draht-Flächen deckt (HTTP + gRPC), nicht nur
eine.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  ist `Accepted` (bereits erfüllt, 2026-09-20).
- Kein weiterer Trigger nötig — die Roadmap führt aktuell keine offene
  Welle (`in-progress/roadmap.md` §Offene Wellen: „Nichts in Arbeit") und
  *Nächste Wellen* ist leer; die Welle kann sofort eröffnet werden.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle fünf Slices in `done/`.
- `make gates` grün (netzloser Gate-Satz, unverändert Docker-only).
- Ein real gebautes `.jar` (`make sdk-pack-kotlin`) als Smoke-Beleg — nicht
  nur ein grüner `docker build`, sondern das Artefakt selbst existiert und
  trägt die erwartete Version (`0.1.0` oder gleichwertig). Ob ein
  Sources-/Javadoc-Jar zusätzlich real nötig ist, entscheidet
  `slice-sdk-kotlin-pack-werkzeug` anhand der realen GitHub-Packages-Regeln
  (kein Vorgriff hier).
- Der GitHub-Packages-Publish-Workflow wurde real gegen
  `https://maven.pkg.github.com/pt9912/pg-change-feed` versucht (Tag
  `sdk-kotlin-v0.1.0` oder gleichwertig) — Erfolg **oder** ein benannter
  roter Befund mit Folgemaßnahme; nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 bleibt dieses Risiko bis zum
  realen Post-Push-Lauf offen, unabhängig vom Rest der Welle.
- Closure-Notiz in `welle-sdk-kotlin-lh-fa-sst-009-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-kotlin-projektgeruest | SDK-Projektgerüst (`sdks/kotlin/pgchangefeed-kotlin/`, Gradle-Projekt, Docker-Bau, leeres API-Skelett) | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |
| slice-sdk-kotlin-http-client-flaeche | HTTP-Client-Fläche (`SPEC-018`, neun Port-gedeckte Fähigkeiten + Changes-Lesen `SPEC-022`, Bearer-Auth, eigene Tests) | [`LH-FA-SST-009`](../../../../spec/lastenheft.md), [`LH-FA-SST-006`](../../../../spec/lastenheft.md) |
| slice-sdk-kotlin-grpc-client-flaeche | gRPC-Client-Fläche (`SPEC-020`, `StreamChanges`, `.proto`-Bezug über Zusatzkontext, Metadata-Auth) | [`LH-FA-SST-009`](../../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../../spec/lastenheft.md) |
| slice-sdk-kotlin-pack-werkzeug | `make sdk-pack-kotlin`, Trägernachzug `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |
| slice-sdk-kotlin-publish-workflow | `.github/workflows/sdk-kotlin-release.yml`, Tag-Präfix `sdk-kotlin-v*`, `GITHUB_TOKEN` (`packages: write`, kein externes Secret) | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |

**Reihenfolge (Abhängigkeitskette):**

1. `slice-sdk-kotlin-projektgeruest` zuerst — legt
   `sdks/kotlin/pgchangefeed-kotlin/` samt Gradle-Projekt/Docker-Bau an; ohne
   dieses Gerüst haben die beiden Fläche-Slices keinen Ort, in den sie
   schreiben.
2. `slice-sdk-kotlin-http-client-flaeche` und `slice-sdk-kotlin-grpc-client-flaeche`
   — beide hängen nur vom Projektgerüst ab, nicht voneinander
   (unterschiedliche Dateien, unterschiedliche Draht-Verträge); sie sind
   **parallelisierbar** (dieselbe Struktur wie bei der C#-Welle).
3. `slice-sdk-kotlin-pack-werkzeug` danach — `make sdk-pack-kotlin` baut,
   testet und paketiert **beide** Flächen in einem Package (`ADR-0109`
   Festlegung 1, ein Package); es braucht deshalb beide Vorgänger-Slices
   fertig.
4. `slice-sdk-kotlin-publish-workflow` zuletzt — der Publish-Workflow
   validiert die Tag-Version gegen die in `build.gradle.kts` geführte
   Version und veröffentlicht das vom Pack-Werkzeug erzeugte Artefakt; ohne
   ein bereits real paketierbares Package liefe er ins Leere.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle — die Roadmap führt aktuell keine
  parallele offene Welle.
- Wird blockiert von: keiner anderen Welle;
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  ist bereits `Accepted`.
- Intern (siehe §4 Reihenfolge): `slice-sdk-kotlin-http-client-flaeche` und
  `slice-sdk-kotlin-grpc-client-flaeche` hängen beide von
  `slice-sdk-kotlin-projektgeruest` ab; `slice-sdk-kotlin-pack-werkzeug`
  hängt von beiden Flächen-Slices ab; `slice-sdk-kotlin-publish-workflow`
  hängt vom Pack-Werkzeug-Slice ab.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **SSE- (`SPEC-021`) und NATS-Vollinhalts-Abdeckung (`SPEC-024`) im
  Package** — obwohl für Kotlin bereits real funktionierende Beispiele
  existieren (`examples/kotlin/sse-client`, `examples/kotlin/nats-stream-client`),
  grenzt [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  Festlegung 1 v1 ausdrücklich auf HTTP-API und gRPC-Stream ein — dieselbe
  Begründung wie bei C# (Abhängigkeits-Footprint, kein zusätzlicher
  fachlicher Nutzen ggü. gRPC für denselben Consumer). Beide bleiben über
  `examples/kotlin/sse-client`/`nats-stream-client` als Vorbild einsehbar,
  bis ein Folge-Package sie deckt (§Re-Evaluierungs-Trigger 1/2 der ADR).
- **Eine vierte SDK-Sprache oder ein vierter Vertriebsweg** —
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  entscheidet ausdrücklich nur die dritte Sprache/den dritten Vertriebsweg;
  eine vierte bleibt eine eigene, künftige ADR (§Re-Evaluierungs-Trigger 1).
- **Ein Wechsel von GitHub Packages zu Maven Central** —
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  §Verglichene Alternativen B2 wägt Maven Central bereits vollständig ab und
  verwirft es (Nutzerentscheidung); ein Wechsel bleibt
  §Re-Evaluierungs-Trigger 5 vorbehalten (Auth-Pflicht erweist sich als
  echtes Consumer-Problem).
- **Aufspaltung in mehrere Packages** (ein Package je Zustellweg) — von
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  Alternative C5 ausdrücklich verworfen.
- **Umbau, Migration oder Extraktion von `examples/kotlin/`** — die ADR
  verbietet das ausdrücklich (§Entscheidung Festlegung 3, Alternative D2
  gewählt); die Beispiele bleiben unverändert Vorbild, kein Belegträger
  (`SPEC-023`).
- **`spec/architecture.md`** bleibt unberührt — `sdks/**` ist außerhalb der
  Server-Komponentensicht (`ADR-0109` §Entscheidung Festlegung 6).
- **`.a-check.yml`** bleibt unberührt — a-check liest kein Kotlin
  (`ADR-0109` §Kontext Bindung 4).
- **Kein neues Gate** — `make sdk-pack-kotlin` und der Publish-Workflow
  bleiben Werkzeuge, kein Gate-Eintrag in `GATE_CHECKS` (`ADR-0109`
  Festlegung 5).
- **`docs/user/version.md` (Server-SemVer) und die C#-/Python-SDK-Versionsräume**
  bleiben unberührt — das Kotlin-Package führt sein eigenes, unabhängiges
  SemVer-2.0-Schema (`ADR-0109` Festlegung 4).
- **Ein Wechsel von JDK 21 auf JDK 25 / Gradle ≥ 9.1.0** —
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  Festlegung 5 bindet die JDK-Basis bewusst an die bereits erprobte
  `examples/kotlin/`-Kombination; ein Wechsel bleibt
  §Re-Evaluierungs-Trigger 3 vorbehalten.

**Sechster Slice — bewusst nicht geschnitten.** Ein separater
„Handbuch-Nachzug"-Slice (`ADR-0109` §Konsequenzen Folgepflicht 4,
`docs/user/benutzerhandbuch.md`-Hinweis je gedeckter Oberfläche **inklusive**
eines expliziten PAT-Hinweises für den Bezug) entfällt als eigener Schnitt —
dieselbe Begründung wie bei der C#-Welle: Die zugrunde liegende Beobachtung
(`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
Ausgang bereits **verkörpert**) trägt die Pflicht bereits an zwei Stellen —
der diff-skopierten Selbstprüf-Instruktion in
`.claude/commands/implement-slice.md` Schritt 17 und dem eigenen HIGH-Punkt
in `.harness/skills/reviewer.md`. Der Handbuch-Hinweis je Oberfläche gehört
deshalb **in** `slice-sdk-kotlin-http-client-flaeche` und
`slice-sdk-kotlin-grpc-client-flaeche` selbst (DoD-Punkt „Doku-Update …
falls öffentlicher Vertrag berührt"), nicht in einen eigenen Slice.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
vollständig durchgesehen, mit besonderem Blick auf die aus den beiden
vorigen SDK-Wellen bekannten Klassen:

- `BEO-PGC/report-nackte-id-ohne-link` — bereits **verkörpert**
  (`AGENTS.md` §3.9), Zähler real ausgezählt bei 8×
  (`ls .../evidence | wc -l`). Bleibt als bereits verkörperte Regel für
  Reviewer-/Verifier-Berichte dieser Welle unverändert gültig — jeder
  Bericht dieser Welle prüft seine eigenen nackten Kennungen selbst nach.
- `BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor` — offen, Zähler real
  1× (`slice-sdk-python-pack-werkzeug`), **unter** der 3×-Schwelle, aber mit
  real nachgewiesener Fernwirkung (ein fehlendes schließendes Backtick
  verschiebt die Codespan-Parität für den Rest eines Dokuments und kann
  dadurch eine an ihrer eigenen Stelle korrekt gebacktickte Kennung an
  anderer Stelle vor `d-check` verbergen). Direkt einschlägig für diese
  Welle: **vor jedem Commit** dieser Welle ein Backtick-Paritäts-Check je
  geänderter Markdown-Datei (Gesamtzahl der Backtick-Zeichen der Datei
  zählen — `wc -l` nach einem `grep -o`-Extraktionsmuster genügt — und
  gegen eine gerade Zahl prüfen) — als eigener Hinweis in jedem der fünf
  Slice-Pläne dieser Welle aufgenommen (§6 Risiken bzw. §3 Plan-Hinweis),
  nicht erst nach einem Reviewer-Finding.
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` — bereits **verkörpert**
  (`AGENTS.md` §3.13), Zähler real 19× (`ls .../evidence | wc -l`). Beide
  vorigen SDK-Wellen fanden hier real neue Belege (README/`pyproject.toml`
  behaupteten „folgt erst" für eine bereits gelieferte Fläche) — der
  Träger-Nachzug-Suchlauf dieser Welle deckt deshalb explizit **beide**
  Achsen ab, die die Python-Welle als unvollständig fand: welche Dateien
  (nicht nur `README.md`, auch `build.gradle.kts`/`gradle.properties`) und
  welche Formulierungen (deutsch **und** englisch, „folgt erst"/„follow-up
  release"/„added by").
- `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` — offen,
  Zähler real 2× (`slice-104`, `slice-sdk-csharp-projektgeruest`), weiter
  **unter** der 3×-Schwelle. Betrifft direkt die Situation dieser Welle
  (Docker-only-Bau eines dritten Sprach-Ökosystems in einer neuen
  Sub-Area) — als wiederkehrender Hinweis in `slice-sdk-kotlin-projektgeruest`
  §6 aufgenommen: vor dem ersten Test-/Bau-Lauf einen Suchlauf im bereits
  bestehenden Bestand (`sdks/csharp/**`, `sdks/python/**`, **und**
  `examples/kotlin/**` — hier zusätzlich zu den beiden Vorbild-Sprachen
  bereits eine reale Kotlin/Gradle-Referenz im selben Repo) auf ein bereits
  gelöstes, analoges Problem prüfen, statt ein Workaround-Muster neu zu
  erfinden. Ein drittes Auftreten dieser Welle würde die 3×-Schwelle
  erreichen; diese Welle ist bewusst so geplant, dass es dazu nicht kommt.
- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` — offen,
  Zähler real 1× (`slice-sdk-csharp-publish-workflow`). Die Python-Welle
  vermied ein zweites Auftreten bereits durch einen vorab eingeplanten
  DoD-Punkt; `slice-sdk-kotlin-publish-workflow` übernimmt dieselbe
  Vorsicht und trägt den `docs/user/releasing.md`-Nachzug (SDK-Release-Weg-
  Abschnitt, analog den bestehenden C#-/Python-Einträgen — hier zusätzlich
  mit der Auth-Pflicht-Anmerkung: kein externes Secret, aber
  `packages: write`-Permission) bereits **vorab als eigenen DoD-Punkt**,
  nicht erst nach einem Reviewer-Finding.
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` — offen, Zähler real 1×
  (`slice-sdk-python-projektgeruest`), unter der 3×-Schwelle. Betrifft
  Plan-Nachzüge, bei denen ein neuer Absatz einen dadurch überholten
  Nachbarabsatz stehen lässt — als Hinweis in jedem Slice-Plan dieser Welle
  aufgenommen: ein Nachzug prüft explizit auch den unmittelbar
  angrenzenden Text auf Widerspruchsfreiheit.

Kein weiterer Treffer (`grep`-Suche nach `Secret`/`GitHub Packages`/
`kotlin`/`gradle`/`examples-kotlin`/`maven.pkg.github` über alle
`observation.md`-Dateien blieb ohne zusätzlichen Treffer über die sechs
oben genannten hinaus). Kein Treffer ist damit selbst die Feststellung
dieser Eröffnung, keine Auslassung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [welle-sdk-kotlin-lh-fa-sst-009-results.md](welle-sdk-kotlin-lh-fa-sst-009-results.md)
Zähler: [../observations/](../observations/)
