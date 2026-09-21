# ADR-0110: Python-SDK-Umfang erweitert auf gRPC/SSE/NATS-Vollinhalt (Supersedes ADR-0107, teilweise)

**Status:** Proposed

**Datum:** 2026-09-21

**Autor:** pt9912 (Architect-Rolle, ausgelöst durch einen Nutzerauftrag: alle
drei bestehenden SDKs — C#, Python, Kotlin — auf volle Abdeckung aller vier
Zugriffswege (`LH-FA-SST-008`) erweitern, dieselbe Matrix, die die
Beispiel-Clients unter `examples/csharp/`/`examples/kotlin/` bereits real
abdecken. Für C# ([`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)) und
Kotlin ([`ADR-0109`](0109-kotlin-github-packages-drittes-sdk-package.md))
deckt ein normaler Implementierungs-Slice diese Erweiterung — beide ADRs
haben SSE/NATS als „Folge-Package" bereits ausdrücklich antizipiert und die
Struktur-Entscheidung „einem Folge-Zug" überlassen, ohne eine neue ADR zu
verlangen. Für Python ([`ADR-0107`](0107-python-pypi-zweites-sdk-package.md))
trägt nicht dieselbe Abkürzung: Jene ADR wählte den kleineren
Erst-Scope **explizit wegen** einer Risiko-Asymmetrie — kein
`examples/python/`-Referenzclient existierte, das Wire-Verständnis-Risiko sei
für gRPC ungleich höher als für HTTP ohne eine solche Gegenprobe
(`ADR-0107` §Kontext „Was das ändert"). Diese ADR prüft real, ob sich das
geändert hat, und trifft — weil es sich **nicht** geändert hat — die
Entscheidung, das Risiko trotzdem jetzt bewusst zu tragen, statt es
stillschweigend zu übergehen, Modul 8 §Rollen-Regeln: „Architect schreibt",
§Konflikt-Pfad ist eine Sequenz, keine Seniorität)

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (Client-Bibliotheken/
SDKs), [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Live-Change-Stream,
gRPC/SSE/NATS-Vollinhalt), [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
(teilweise superseded — siehe Kopf), [`ADR-0108`](0108-python-sdk-uv-statt-build-twine.md)
(Bau-/Publish-Frontend — von dieser ADR **unberührt**, andere Klausel),
[`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) und
[`ADR-0109`](0109-kotlin-github-packages-drittes-sdk-package.md)
(Kontrastfolie: dort war derselbe Scope-Sprung bereits ADR-seitig antizipiert,
hier nicht — der Grund, warum Python eine eigene ADR braucht und C#/Kotlin
keine), [`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md) §2
`SPEC-018`/`SPEC-020`/`SPEC-021`/`SPEC-024` (die vier Draht-Festlegungen),
`SPEC-027` (das bestehende Python-Package), [`AGENTS.md`](../../../AGENTS.md)
§3.1 (Docker-only) §3.5 (Accepted-ADR-Immutabilität — der Grund für diese
Datei statt einer In-Place-Korrektur an `ADR-0107`) §3.6 (Gates nur per ADR
gelockert — hier irrelevant, kein Gate berührt) §3.13 (Träger-Nachzug)

**Schärft:** [`LH-FA-SST-009.a`](../../../spec/pflichtenheft.md) — dieselbe
Pflichtenheft-Stelle, die [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
für Python/PyPI beantwortet hat, bleibt für den **Umfang** eines
Python-Folge-Release weiterhin über diese ADR-Kette (jetzt:
[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) +
[`ADR-0108`](0108-python-sdk-uv-statt-build-twine.md) + diese ADR)
beantwortet; an der Sprach-/Vertriebsweg-Antwort selbst (Python/PyPI als
zweite Sprache) ändert diese ADR nichts.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Die Lücke.** [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
§Entscheidung Festlegung 1 begrenzte den Python-SDK-Erst-Scope bewusst auf
HTTP-API allein — **explizit begründet** mit einer Risiko-Asymmetrie: Ohne
ein laufendes `examples/python/`-Referenzprogramm ist das
Wire-Verständnis-Risiko für gRPC (Streaming-Lebenszyklus, generierter Stub,
Channel-Credentials) real höher als für ein zustandsloses REST/JSON-Protokoll.
Der einzige explizit **antizipierte** nächste Schritt war **ausschließlich
gRPC** („Ein Folge-Release trägt gRPC, sobald ein erstes Python-Package real
existiert und der HTTP-Teil real geprüft ist"), nicht SSE oder
NATS-Vollinhalt — beide blieben ohne konkreten Folgeplan, nur als „nicht in
v1" genannt. `ADR-0107` §Re-Evaluierungs-Trigger 3 benennt selbst die
Bedingung, unter der ein **größerer** Erst-/Folge-Scope ohne die dort
geltende Vorsicht zu rechtfertigen wäre: „Ein Python-Referenz-Client entsteht
nachträglich (… `examples/python/` …) — dann verliert die … begründete
Risiko-Asymmetrie ihre Grundlage."

**Der Ist-Stand — gemessen, nicht erinnert (heutiger Zug).**

| Prozedur | Ergebnis |
|---|---|
| `find examples -maxdepth 1 -type d` | `examples/csharp`, `examples/kotlin`, die flache Go-Wurzel (`http-client`, `grpc-client`, `sse-client`, `nats-client`, `nats-stream-client`) — **kein** `examples/python`, genau wie zum Zeitpunkt von `ADR-0107` |
| `find examples -iname "*python*"` | kein Treffer |
| `git log --oneline --all -- examples/python` | leer — es gab **nie** einen Commit, der diesen Pfad berührt |
| `git tag \| grep sdk-python` | `sdk-python-v0.1.0` — das Python-Package ist real auf PyPI-Release-Pfad veröffentlicht (`docs(releasing): sdk-python-v0.1.0-Realisierung nachgezogen`) |
| `git log --oneline --all \| grep sdk-python` | zeigt eine vollständig durchlaufene Welle (`welle-sdk-python-lh-fa-sst-009`, Projektgerüst → HTTP-Client-Fläche → Pack-Werkzeug → Publish-Workflow, jede Stufe review- und verifikationsbelegt, nach `done/`) — der „HTTP-Teil real geprüft"-Teil von `ADR-0107` Festlegung 1s eigenem Anschluss-Kriterium ist damit erfüllt |
| `find sdks/python -maxdepth 3 -type f` | nur `Dockerfile`, `README.md`, `pgchangefeed/pyproject.toml`, `pgchangefeed/tests/**` — kein gRPC-/SSE-/NATS-Code |

**Was das für diese ADR bedeutet.** `ADR-0107` §Re-Evaluierungs-Trigger 3 ist
**nicht eingetreten** — die real begründete Risiko-Asymmetrie (kein
Python-Referenzclient) besteht unverändert fort. Der Auftrag dieses Zuges
verlangt trotzdem **mehr** als das, was `ADR-0107` Festlegung 1 selbst als
nächsten Schritt vorsah (nur gRPC): volle Vier-Wege-Parität, also zusätzlich
SSE und NATS-Vollinhalt, ohne dass die eine Bedingung eintrat, die laut
`ADR-0107` selbst eine größere Vorsicht-Lockerung tragen würde. Eine stille
Erweiterung wäre deshalb keine Umsetzung von `ADR-0107`, sondern eine
Korrektur seiner §Entscheidung/§Verglichene-Alternativen-Kernaussage
(Tabelle C, Option C1 gegen C2/C3/C4 explizit **wegen** der
Risiko-Asymmetrie entschieden) — `AGENTS.md` §3.5 verlangt dafür eine neue
ADR mit `Supersedes`, keine Nachbesserung im laufenden Implementierungs-Zug.

**Was sich seit `ADR-0107` real geändert hat — und was nicht.** Kein
Python-Referenzclient (Trigger 3, unverändert offen). Aber zwei Tatsachen,
die `ADR-0107` beim Schreiben zwar kannte, aber nicht als Risiko-Mitigation
für Python **selbst** bewertete, weil sie dort nicht die Frage war:

1. **Alle vier Zugriffswege haben inzwischen funktionierende
   Referenzimplementierungen in zwei Fremdsprachen (C#, Kotlin) plus Go**,
   real gebaut und getestet gegen denselben, unveränderten Draht-Vertrag
   (`SPEC-018`/`SPEC-020`/`SPEC-021`/`SPEC-024`) — ein Python-Team kann die
   Protokoll-**Semantik** (Stream-Lebenszyklus, Metadata-Header-Form,
   Nachrichtenschema) an einer laufenden Gegenprobe verifizieren, auch ohne
   dass diese Gegenprobe in Python selbst geschrieben ist. Das ist ein
   **schwächerer** Risiko-Minderer als ein sprachnatives Beispiel (er nimmt
   Python-Bibliotheks-/Idiom-Fragen nicht ab), aber ein **realer**, der bei
   `ADR-0107`s Erst-Scope-Entscheidung nicht als Mitigation herangezogen
   wurde.
2. **`make test-integration`** bietet eine reale, laufende Server-Instanz
   (Compose-Umgebung) gegen die ein Python-Client seine eigenen
   Protokoll-Annahmen ausführbar prüfen kann — kein Python-Beispielprogramm,
   aber eine ausführbare Gegenprobe, die zur Zeit von `ADR-0107` ebenfalls
   schon existierte, aber nicht als Teil der Scope-Begründung herangezogen
   wurde.

Diese zwei Tatsachen genügen **nicht**, um Trigger 3 als erfüllt zu erklären
— sie sind keine Python-Referenzimplementierung. Sie senken das reale Risiko
aber spürbar gegenüber dem Zustand, den `ADR-0107` bewertete (zum
damaligen Zeitpunkt existierte noch kein einziges funktionierendes
Fremdsprachen-Vorbild für den **NATS-Vollinhalts**-Weg außerhalb von Go —
`examples/csharp/nats-stream-client` und `examples/kotlin/nats-stream-client`
sind seither hinzugekommen). Diese ADR wählt deshalb bewusst, das
verbleibende, real benannte Risiko zu tragen, statt es zu ignorieren oder auf
einen hypothetischen künftigen Trigger zu warten (§Entscheidung).

## Entscheidung

Wir wählen: **Der Python-Package-Umfang wird für ein Folge-Release auf die
volle Vier-Wege-Matrix erweitert — gRPC (`SPEC-020`), SSE (`SPEC-021`) und
NATS-Vollinhalt (`SPEC-024`), zusätzlich zum bereits bestehenden
HTTP-API-Umfang (`SPEC-018`) — bei explizit fortbestehender, real weiter
höherer Risikolage ggü. C#/Kotlin.** Diese ADR ersetzt ausschließlich
[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) §Entscheidung
Festlegung 1 (Umfang) und den zugehörigen Teil von §Verglichene
Alternativen Tabelle C. Alles Übrige von `ADR-0107` bleibt **hiermit
bestätigt** und wird nicht wiederholt: Festlegung 2 (Vertriebsweg PyPI,
Arbeitsname `pgchangefeed`), Festlegung 3 (Ort `sdks/python/`,
Import-Grenze), Festlegung 4 (Versionierung PEP 440, Start `0.x.y`),
Festlegung 5 in der durch [`ADR-0108`](0108-python-sdk-uv-statt-build-twine.md)
bereits geänderten Fassung (Docker-only, `uv build`/`uv publish`, Secret-Klasse
`PYPI_API_TOKEN`, Tag-Namensraum `sdk-python-v*`), §Kontext-Messung (die
Ist-Stand-Prozeduren von damals bleiben historisch korrekt), §Verglichene
Alternativen Tabelle A/B/D/E vollständig, alle Punkte in §Entscheidung
Festlegung 6 („Was diese ADR nicht ändert"), soweit sie mit der hier
getroffenen Erweiterung vereinbar sind, §Re-Evaluierungs-Trigger 1/2/4 und
§Geschichte.

### 1 — Neuer Umfang: gRPC + SSE + NATS-Vollinhalt als ein gemeinsames Folge-Release

Anders als bei C#/Kotlin (ein Erst-Release deckt HTTP+gRPC, SSE/NATS bleiben
offen) und anders als `ADR-0107`s eigener, engerer Anschluss-Plan (nur
gRPC als nächster Schritt) deckt das Python-Folge-Release **alle drei**
verbleibenden Wege in einem Zug. Begründung:

- **Der Auftrag verlangt Vier-Wege-Parität über alle drei SDKs**, nicht eine
  isolierte, weiter gestufte Ausrollung für Python allein — ein
  Zwischenschritt „nur gRPC jetzt, SSE/NATS erst nach Nutzungsdaten"
  (`ADR-0107`s eigener, ursprünglicher Plan) würde die drei Sprachen
  ungleich weit bringen, ohne einen fachlichen Grund, der diese Ungleichheit
  rechtfertigt.
- **Das Restrisiko ist für gRPC, SSE und NATS keine grundsätzlich
  unterschiedliche Klasse mehr**, seit alle drei Wege reale
  Fremdsprachen-Referenzen tragen (§Kontext „Was sich geändert hat")
  — die ursprüngliche Abstufung in `ADR-0107` (gRPC zuerst, weil dort das
  höchste Risiko lag) verliert an Trennschärfe, wenn ohnehin alle drei neu
  hinzukommen.
- **Ein einziges Folge-Release statt drei aufeinanderfolgender** hält die
  Paket-Design-Kosten (öffentliche API-Form, Versionsgrenze, Doku, Tests als
  Vertrag — fallen laut `ADR-0106`/`ADR-0107` je Zustellweg an, unabhängig
  vom Vorbild) in einem Release-Zyklus statt in drei.

### 2 — Die Risiko-Asymmetrie besteht fort und wird nicht verschwiegen

`ADR-0107` §Re-Evaluierungs-Trigger 3 (Entstehen eines
`examples/python/`-Referenzclients) ist **nicht** eingetreten (§Kontext
Ist-Stand, heute real gemessen). Diese ADR verschweigt das nicht und
behandelt es nicht als erfüllt:

- **Kein Python-Referenzclient wird durch diese ADR verlangt oder
  angelegt.** `examples/python/` bleibt weiterhin nicht existent — diese
  ADR mandatiert **keine** neue Beispiel-Client-Fläche; das wäre eine
  andere, größere Entscheidung (`SPEC-023`s Umfang) ohne fachlichen
  Zusammenhang zu dieser.
- **Die Mitigation ist schwächer als ein sprachnatives Vorbild, aber real:**
  funktionierende Fremdsprachen-Referenzen (C#, Kotlin, Go) gegen denselben,
  unveränderten Draht-Vertrag plus eine reale, ausführbare Server-Gegenprobe
  (`make test-integration`) senken das Protokoll-Verständnis-Risiko
  spürbar, heben es aber nicht auf das Niveau von C#/Kotlin (die volle
  Vorarbeits-Parität, `ADR-0109` §Kontext).
- **Folgepflicht für den umsetzenden Zug** (§Konsequenzen): jede neue
  Zustellweg-Fläche (gRPC-Client, SSE-Client, NATS-Vollinhalts-Client) prüft
  ihre Protokoll-Annahmen **zusätzlich** zu Unit-Tests gegen eine reale,
  laufende Server-Instanz — dieselbe Disziplin wie
  `test-integration`/`test-notify`/`test-replication` es für den Go-Server
  bereits vormachen, nicht nur gegen die Klartext-Spezifikation gelesen.

### 3 — Struktur (ein Package vs. Folge-Package) bleibt Sache des umsetzenden Zuges

Wie bei C# ([`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 1, letzter Absatz) und Kotlin
([`ADR-0109`](0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 1) entscheidet diese ADR **nicht**, ob das Folge-Release eine
`v2` desselben `pgchangefeed`-Packages ist oder ein eigenständiges Package
wird — dieselbe Delegation an „einen Folge-Zug", kein Vorgriff hier.

### 4 — Was diese ADR nicht ändert

- **Kein Produktionscode, kein Eingriff in `internal/**`/`cmd/**`.**
- **`ADR-0106`/`ADR-0108`/`ADR-0109` und `sdks/csharp/**`/`sdks/kotlin/**`
  bleiben unberührt** — diese ADR betrifft ausschließlich den
  Python-Umfang.
- **Vertriebsweg (PyPI), Ort (`sdks/python/`), Versionierung (PEP 440),
  Bau-/Publish-Frontend (`uv`, `ADR-0108`) bleiben unverändert.**
- **`examples/python/` wird durch diese ADR nicht angelegt** (§Entscheidung
  Festlegung 2).
- **Kein neues Gate, keine Schwellen-Senkung.**
- **`.a-check.yml` bleibt unberührt** — a-check liest kein Python.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### C — Umfang des Python-Folge-Release — Fortsetzung von `ADR-0107` Tabelle C

| Option | Pro | Contra |
|---|---|---|
| C1 — nichts tun: bei `ADR-0107`s Plan bleiben (nur gRPC als nächster Schritt, SSE/NATS erst nach Trigger 2/3) | kein neuer ADR-Aufwand; folgt der ursprünglichen, konservativen Reihenfolge exakt | lässt Python hinter C#/Kotlin fachlich unbegründet zurück, wenn der Auftrag ausdrücklich Drei-SDK-Parität verlangt; verzögert SSE/NATS auf einen Trigger, der von realen Download-/Issue-Daten abhängt, die für ein Erst-Release absehbar lange ausbleiben können |
| C2 — nur gRPC jetzt (Festlegung-1-Plan von `ADR-0107` wortgetreu umgesetzt), SSE/NATS bleiben offen | am nächsten an der ursprünglichen, bereits durchdachten Reihenfolge; kleinerer Einzelschritt | erfüllt den aktuellen Auftrag (Vier-Wege-Parität) nur teilweise; SSE/NATS blieben ein zweites, separates ADR-Bedürfnis später, wenn dieselbe Risikofrage erneut gestellt würde |
| **C3 — gRPC + SSE + NATS-Vollinhalt in einem Folge-Release (gewählt)** | erreicht Vier-Wege-Parität mit C#/Kotlin in einem Zug; adressiert die Risikofrage einmal, ehrlich, statt sie dreimal einzeln neu zu stellen; nutzt die inzwischen real vorhandene Fremdsprachen-Referenz (C#/Kotlin/Go) für alle drei Wege gleichermaßen | größerer Einzelschritt ohne Python-natives Vorbild; die Risiko-Asymmetrie aus `ADR-0107` besteht real fort und wird hier bewusst in Kauf genommen, nicht aufgelöst |
| C4 — voller Umfang, zusätzlich eine neue `examples/python/`-Beispiel-Fläche als Vorarbeit anlegen, bevor das SDK erweitert wird | würde Trigger 3 real erfüllen und die Risikolage auf C#/Kotlin-Niveau heben | deutlich größerer Aufwand für dieselbe Ziel-Fläche zweimal (Beispiel **und** SDK) bei einer Anforderung, die ausschließlich das SDK verlangt (`LH-FA-SST-009`, nicht `SPEC-023`); widerspricht der Instruktion, die Lösung mit dem geringsten realen Aufwand zu bevorzugen, wenn eine schlankere Mitigation (§Entscheidung Festlegung 2) bereits ausreicht |

**Fazit:** C3. C1/C2 erfüllen den aktuellen Auftrag nicht vollständig, C4
verdoppelt den Aufwand für eine Mitigation, die schlanker (Festlegung 2)
ebenfalls trägt.

## Konsequenzen

- **Positiv:** Python erreicht dieselbe Vier-Wege-Zielmatrix wie C#/Kotlin in
  einem koordinierten Folge-Release, ohne die reale Risikolage zu
  verschweigen; die Entscheidung ist jetzt explizit dokumentiert statt einer
  stillen Erweiterung im Implementierungs-Zug, die `AGENTS.md` §3.5 verletzt
  hätte.
- **Negativ:** Das Python-Folge-Release trägt ein real höheres
  Protokoll-Verständnis-Risiko als das C#-/Kotlin-Äquivalent (kein
  sprachnatives Vorbild) — mitigiert, nicht aufgehoben, durch
  Fremdsprachen-Referenzen und Server-Integrationstests (§Entscheidung
  Festlegung 2).
- **Folgepflicht:**
  1. Ein Planner-Zug schneidet die Umsetzung (voraussichtlich mehrere
     Slices: gRPC-Client-Fläche, SSE-Client-Fläche, NATS-Vollinhalts-Client-
     Fläche, jeweils mit einem Protokoll-Test gegen eine reale
     Server-Instanz, nicht nur gegen die Spezifikation gelesen) — diese ADR
     entscheidet die Form, nicht den Schnitt.
  2. `spec/pflichtenheft.md` `SPEC-027` und `LH-FA-SST-009.a` brauchen einen
     Träger-Nachzug, sobald das Folge-Release real existiert (`AGENTS.md`
     §3.13) — Sache des umsetzenden Zuges.
  3. `docs/user/benutzerhandbuch.md` bekommt einen SDK-Hinweis je neu
     gedeckter Oberfläche, im selben Zug wie das jeweilige Client-Programm.
  4. Diese ADR wird erst durch eine explizite Annahme-Entscheidung
     (`Status: Accepted`) wirksam — sie bleibt bis dahin `Proposed` und
     entscheidet nichts verbindlich (dieser Zug nimmt keine ADR selbst an,
     Nutzerentscheidung ist ein separater Schritt).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `a-check` | keine — `sdks/python/**` ist keine Go-Quelle, a-check liest sie nicht | — |
| Review (kein Sensor) | jede neue gRPC-/SSE-/NATS-Client-Fläche in `sdks/python/pgchangefeed/**` trägt mindestens einen Test gegen eine reale, laufende Server-Instanz (nicht ausschließlich Unit-Tests gegen einen Mock) — §Entscheidung Festlegung 2, Folgepflicht 1 | — |
| Review (kein Sensor) | `sdks/python/**` importiert weiterhin ausschließlich die Python-Standardbibliothek und öffentliche PyPI-Pakete, keinen privaten Baum dieses Repos (unverändert aus `ADR-0107`) | — |

## Re-Evaluierungs-Trigger

1. **Ein Python-Referenz-Client entsteht nachträglich** (`ADR-0107`
   §Re-Evaluierungs-Trigger 3, unverändert gültig) — löst die hier bewusst
   in Kauf genommene Risiko-Asymmetrie endgültig auf, kein neuer Trigger
   nötig, aber ein sinnvoller Anlass, die Test-Disziplin aus
   §Entscheidung Festlegung 2 gegen das dann existierende Vorbild
   abzugleichen.
2. **Ein reales Protokoll-Missverständnis wird im Python-Folge-Release
   gefunden** (z. B. ein Fehlschlag im in Festlegung 2 verlangten
   Server-Integrationstest, der auf eine falsch verstandene
   Stream-Semantik zurückgeht) — Beleg dafür, dass die hier akzeptierte
   Risikolage real eingetreten ist; keine neue ADR nötig für die Korrektur
   selbst (ein Bugfix ohne Entscheidungsänderung), aber ein Anlass, die
   Test-Pflicht aus Festlegung 2 zu verschärfen, falls sich das Muster
   wiederholt.
3. **`ADR-0107`s Trigger 1/2/4 treten ein** (dritte Sprache/dritter
   Vertriebsweg verlangt; reale PyPI-Nutzungsdaten legen eine andere
   Priorisierung nahe; Trusted-Publishing-Umstellung) — unverändert gültig,
   von dieser ADR nicht berührt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-21 | Proposed | dieser Architect-Zug (kein vorausgehender Slice-Plan — ADR vor Implementierung, Modul 8); ausgelöst durch den Auftrag, alle drei SDKs auf volle Vier-Wege-Parität zu erweitern |

Diese ADR ist **`Proposed`**, nicht `Accepted` — die Annahme ist ein
separater, expliziter Nutzerentscheidungs-Schritt, den dieser Zug nicht
vorwegnimmt. Nach einer künftigen `Accepted`-Setzung wird diese Datei nicht
mehr inhaltlich überschrieben; spätere Korrekturen entstehen als neue ADR mit
`Supersedes ADR-0110` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für
Accepted-ADRs).
