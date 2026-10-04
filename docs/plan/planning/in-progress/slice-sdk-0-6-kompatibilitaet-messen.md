# Slice sdk-0-6-kompatibilitaet-messen: Binär-/Quellkompatibilität der 0.6.0-SDKs und ihr Verhalten gegen einen Server vor 0.6.0 messen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD dieses
Slice (gemessen am Planungsstand: `ls docs/plan/planning/*.md` nennt nur
`README.md`, die Roadmap führt unter *Nächste Wellen* keine Zeile), siehe
Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht
(Modul 6).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (SDK-Packages,
Scope),
[`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
(Festlegung 4: „jede 0.5.x-Signatur bleibt binär und quellseitig erhalten“;
Konsequenz: der Code gilt nur dort, wo der Server einen setzt — beide Aussagen
*hergeleitet*, nicht gemessen),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
(Server-Seite der Meldungscodes: Release 0.6.0),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Festlegung 2: Mechanik der SDK-Realserver-Tiers),
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Festlegung 2: Index-Digest als Pin-Form, je Image ein Wert),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Ursprung je
Zahl und Aussage),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine
interne Kennung unter `sdks/`); Herkunft der zwei Messaufträge: der
Verifikationsbericht
[`verifikation-slice-sdk-meldungscodes-in-fehlertypen`](../../../reviews/verifikation-slice-sdk-meldungscodes-in-fehlertypen.md)
§6 („Binärkompatibilität … *hergeleitet* (Gast-Assembly nicht gemessen)“,
„Server alt/neu … *hergeleitet* (am alten Server nicht gelaufen)“) und der
Review-Befund F-1 in
[`review-slice-sdk-meldungscodes-in-fehlertypen`](../../../reviews/review-slice-sdk-meldungscodes-in-fehlertypen.md).

**Berührte Spec-Stellen:** — (der Slice misst und ändert weder Wire noch API;
zu belegen durch den Diff: keine Datei unter `spec/`).

**Verantwortlich:** —

<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** Planner-Agent, Auftrag des Aufrufers (kein Architect: der Slice trifft
keine Entscheidung; eine Berührung von
[`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) durch
einen Befund wird in §7 nur benannt). **Datum:** 2026-10-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass.** Zwei Aussagen der Entscheidung
[`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
(`Accepted`) und der Closure des Slice
[`slice-sdk-meldungscodes-in-fehlertypen`](../done/slice-sdk-meldungscodes-in-fehlertypen.md)
stehen als *hergeleitet*, nicht als gemessen; die SDK-Packages 0.6.0 sind seit
der Veröffentlichung öffentlich (Tags `sdk-csharp-v0.6.0`, `sdk-python-v0.6.0`,
`sdk-kotlin-v0.6.0`, gemessen: `git tag` am Planungsstand):

- **(A) Kompatibilität.** Ein gegen 0.5.0 gebautes Gast-Programm läuft ohne
  Neukompilierung gegen 0.6.0 (C#: Assembly, Überladungen je Fehlertyp; Kotlin:
  Class-Files, `@JvmOverloads`; Python: keine Binärkompatibilität, aber
  Aufruf-Kompatibilität, `message_code` keyword-only), und die Quellseite trägt
  — auch am `null`-Literal der protected Basis-Konstruktoren der C#-HTTP-Hierarchie
  (Review-Befund F-1: `base(status, text, null)` in einer fremden Unterklasse
  ist hergeleitet mehrdeutig, nicht kompiliert).
- **(B) Server vor 0.6.0.** Die 0.6.0-SDKs gegen einen Server ohne `code` im
  Fehlerkörper, ohne `ErrorInfo` im gRPC-Statusdetail und ohne `error_code` im
  Diagnose-`HeartbeatStatus`: die Eigenschaft bleibt leer, Fehlertext und
  Fehlertyp bleiben, kein Absturz.

**Ziel:** Beide Aussagen sind je Sprache (C#, Kotlin, Python) an einem realen Lauf
gemessen — jede mit gedruckter Zeile, Exit und Ursprungsangabe, jede mit einer
Probe, die rot wird, wenn die Aussage falsch ist — oder als Befund für den
Planner benannt; ohne Produktänderung und ohne Eingriff in `make gates`.

**Messweg — vom Planner an der Quelle geprüft (2026-10-04, Parent
`ff97a91059f6f706b16cb16b0681a1563a064d83`, Messhost `x86_64`); die Werte sind
Momentaufnahmen, der Implementer misst sie neu
([`AGENTS.md`](../../../../AGENTS.md) §3.12):**

| Gegenstand | Messung (Befehl) | Ergebnis |
|---|---|---|
| Server-Image 0.5.0 existiert | `docker buildx imagetools inspect ghcr.io/pt9912/pg-change-feed:0.5.0` | Index-Digest `sha256:f99a77ff335a5fa2e84937771337d03711b6ad4b4ac9fb038dd6e310778707ce`, Mitglieder `linux/amd64` und `linux/arm64` (gemessen) |
| Image anonym ladbar | `docker pull ghcr.io/pt9912/pg-change-feed:0.5.0@<Index-Digest>` | `Status: Downloaded newer image`, Entrypoint `/pg-change-feed`, Architektur `amd64` (gemessen, `docker image inspect`); das Label `org.opencontainers.image.version` ist leer, die Version ergibt sich allein aus Tag und Digest |
| Image 0.6.0 zum Vergleich | `docker buildx imagetools inspect ghcr.io/pt9912/pg-change-feed:0.6.0 --format '{{.Manifest.Digest}}'` | `sha256:dd50ed55eb5ac9bb29cdc3df5ece621ab0b1a4078822835605ee23ad86c11ec3` (gemessen); die Tags der Registry nennen `0.6.0`, `latest` und die Vorgänger bis `0.1.2` (`gh api /users/pt9912/packages/container/pg-change-feed/versions`) |
| Server 0.5.0 trägt kein `ErrorInfo` | `git grep -n -E 'errdetails\|ErrorInfo' v0.5.0 -- internal cmd` (ohne Testdateien) | 0 Zeilen; am Tag `v0.6.0` 2 Dateien (`internal/adapters/driving/grpc/errors.go`, `…/administration.go`) — gemessen |
| Gast-Quellen 0.5.0 sind öffentlich | NuGet: `https://api.nuget.org/v3-flatcontainer/pgchangefeed.client/index.json`; PyPI: `https://pypi.org/pypi/pgchangefeed/0.5.0/json`; Cloudsmith: `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/io/github/pt9912/pgchangefeed-kotlin/0.5.0/pgchangefeed-kotlin-0.5.0.jar` (README-Pfad des Kotlin-SDK; der Pfad `maven.cloudsmith.io` antwortet 401, der README-Pfad `dl.cloudsmith.io` 200) | NuGet nennt die Versionen 0.1.0 bis 0.6.0, PyPI HTTP 200, Cloudsmith HTTP 200 für 0.5.0 und 0.6.0 — je Abfrage in einem Container (`golang:1.27-alpine` mit `apk add curl`, Muster `tools/harness/lib-github-api.sh`); GitHub Packages (Maven) nennt 0.2.0 bis 0.6.0 und verlangt Anmeldung (`gh api`) — nicht verwendet |
| Arbeitsstand der SDKs gleicht dem Tag | `git diff --stat sdk-<sprache>-v0.6.0 HEAD -- sdks/<sprache>` | je Sprache 2 Dateien, 6 bzw. 7 Zeilen (gemessen; Inhalt: Pin-Zeilen der Dockerfiles, nicht gelesen) — der Implementer prüft, dass keine Quelldatei der Bibliothek abweicht |
| Das Runner-Gerüst hängt an `:dev` | `compose.yaml` Zeile 72 `image: ghcr.io/pt9912/pg-change-feed:dev`; `tools/harness/run-sdk-csharp-integration-tests.sh` Zeile 74 `COMPOSE=${COMPOSE:-docker compose -f compose.yaml}`, Zeile 112 Prüfung auf das `:dev`-Image | gelesen; `COMPOSE` ist übersteuerbar, das Image nicht ohne Override-Datei |
| Die Fehlercode-Phase der drei Runner erwartet den Code | `tools/harness/run-sdk-*-integration-tests.sh` Phase „Fehlercode-Fläche“, Erwartung `RECEIVED code=PCF-E[0-9]{4} http=404 …` | gelesen: gegen einen Server vor 0.6.0 ist diese Phase **absichtlich rot**; die vorhandenen Runner können Messung B nicht unverändert tragen (neuer Runner, §3) |
| Schema-Unterschied der zwei Stände | `git diff --stat v0.5.0 v0.6.0 -- tools/schema proto` | `process_heartbeat.error_code` (nullable Spalte), View `cdc.heartbeat` mit `error_code` als letzter Spalte, `HeartbeatStatus.error_code` (Feld 4), sonst Rollout-Mechanik (`rollout.sh`, Dockerfile) — gemessen; dass der Server 0.5.0 gegen das **aktuelle** Schema läuft, ist *hergeleitet* (additive nullable Spalte) und wird in Messung B0 gemessen |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein SDK-Release, kein Version-Bump, kein Tag:** der Slice misst die bereits
  veröffentlichten Stände; jede Veröffentlichung braucht eine eigene
  Nutzerfreigabe, und die Versionsdateien (`.csproj`, `pyproject.toml`,
  `build.gradle.kts`) bleiben unberührt.
- **Keine Produktänderung:** `internal/`, `cmd/`, `proto/` und die
  Bibliotheks-Quellen der SDKs (`sdks/*/…/src/main`, `sdks/csharp/PgChangeFeed.Client/` ohne
  Tests, `sdks/python/pgchangefeed/src/`) bleiben unberührt. Neu sind nur
  Mess-Quellen: die Gast-Programme unter `tools/harness/sdk-kompat/`, der Runner
  und je SDK eine **Test**-Klasse im bestehenden Integrations-Quellsatz. Ein Fund
  (zum Beispiel eine mehrdeutige Überladung) ist ein Befund für den Planner, kein
  stiller Fix.
- **Kein Nachzug in `Accepted` Records:**
  [`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md),
  der Review, der Verifikationsbericht und der Plan in `done/` halten den Stand
  ihrer Messung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); der Befund steht im Befund-Nachtrag
  in §7 dieses Plans, und der Suchlauf in §3 hält die ADR-Dateien unverändert.
- **Kein Eingriff in die drei bestehenden Realserver-Runner, in `compose.yaml`
  und in `make gates`:** der Runner-Vertrag bleibt; die neuen Ziele stehen
  nicht in `GATE_CHECKS` (sie brauchen Netz und Docker-Pulls).
- **Kein Messen anderer Alt-Stände:** ein Server vor 0.5.0 und SDKs vor 0.5.0
  sind nicht Gegenstand; die Aussagen von
  [`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
  nennen 0.5.x gegen 0.6.0. Ein weiterer Stand wäre ein Parameter des neuen
  Runners (`SDK_ALTSERVER_IMAGE`), kein eigener Slice.
- **Kein Fehlerweg der gRPC-Stream-Clients und kein NATS-Fehler:** sie bleiben
  roh bzw. ohne Code
  ([`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
  Festlegung 5); es gibt dort nichts, das alt gegen neu verschieden wäre.
- **Keine Neubewertung der Entscheidung:** die Re-Evaluierung von
  [`ADR-0064`](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md) (echter
  Alt-Image-gegen-Neu-Image-Vergleich) gehört dem Architect; der Slice liefert
  dazu nur die Beobachtungs-Evidenz (§2) und benennt die Berührung in §7.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Jede gedruckte Zahl und jede Aussage der Belege trägt ihren Ursprung
(*gemessen* mit Lauf und gedruckter Zeile · *übernommen* · *abgeleitet*,
[`AGENTS.md`](../../../../AGENTS.md) §3.12); eine nicht gefahrene Verallgemeinerung
(zum Beispiel „gilt für Kotlin und Python wie für C#“) steht als *hergeleitet*.

- [x] **Messung A — Kompatibilität (Liefer-Punkt 1).** *(Belege und Abweichungen von A3/A4: §7, „Messung am Implementer-Stand“ und „Befund-Nachtrag (A)“.)* `make test-sdk-kompat`
      (Docker-only, Netz für NuGet/PyPI/Cloudsmith/Maven Central; kein Gate) fährt je
      Sprache (C#, Kotlin, Python) ein **Gast-Programm**, das jede öffentliche
      Konstruktor-/Aufrufform der 0.5.x-Fehlertypen beider Hierarchien (HTTP, gRPC)
      benutzt. Gedruckt je Sprache und Schritt eine Zeile
      `KOMPAT <sprache> <schritt>: <n> Aufrufe ok` bzw. die Ausnahme-Zeile, Exit 0 genau bei
      Erwartung. Schritte: **A1** Gast gegen das veröffentlichte 0.5.0 gebaut und
      gegen 0.5.0 gelaufen (Grundlinie: beweist, dass der Gast läuft); **A2** dieselben
      Gast-Binärdateien **ohne Neubau** gegen die Bibliothek 0.6.0 (Austausch der
      Bibliotheksdatei; Quelle der 0.6.0-Bibliothek per `SDK_KOMPAT_NEU=dist`, Default:
      die Artefakte von `make sdk-pack-*`, und einmal `SDK_KOMPAT_NEU=registry`: die
      veröffentlichten 0.6.0-Pakete — beide Läufe stehen in §7); **A3** Negativprobe
      „Gegenrichtung“: ein gegen 0.6.0 gebauter Gast (nutzt Konstruktor mit Code bzw.
      `message_code`) gegen die 0.5.0-Bibliothek muss mit `MissingMethodException`
      (C#), `NoSuchMethodError` (Kotlin) bzw. `TypeError` (Python) scheitern —
      beweist, dass der Austausch die Bindung wirklich prüft; **A4** Mutationsprobe in
      **C#**: eine Kopie der C#-Bibliothek im Scratchpad ohne den alten
      Zwei-Argument-Konstruktor **eines** Blatt-Typs
      (`PgChangeFeedBadRequestException(int, string)`), gebaut über die Stufe `build` von
      `sdks/csharp/Dockerfile`, lässt den A2-Gast an dieser Stelle scheitern
      (`MissingMethodException`); die Übertragung auf Kotlin (`@JvmOverloads` entfernt) und
      Python (`message_code` ohne Default) steht als *hergeleitet*, solange sie nicht
      gefahren ist (der Implementer fährt sie, wenn der Aufwand klein ist, und nennt
      dann Stelle, Instanz und Farbe). **A5** Quellseite: derselbe Gast-**Quelltext** wird
      gegen 0.5.0 und gegen 0.6.0 übersetzt (C#, Kotlin) bzw. unter `inspect.signature`
      gelesen (Python); dazu in C# die **`null`-Matrix** von F-1: je eine Fremd-Unterklasse
      der abstrakten Basis `PgChangeFeedException` (HTTP) und `PgChangeFeedGrpcException`
      mit `base(…, null)` in jeder Aufrufform, die unter 0.5.0 übersetzt, gegen 0.6.0 — die
      Erwartung (hergeleitet) ist `CS0121` an `(int, string, null)` der HTTP-Basis;
      Python: ein positionales drittes Argument am Fehlertyp muss unter 0.6.0
      `TypeError` liefern (keyword-only gemessen). Ein Ergebnis gegen die Erwartung ist
      ein **Befund** (§7), kein Anlass, die Probe anzupassen.
- [x] **Messung B — Server vor 0.6.0 (Liefer-Punkt 2).** *(Belege: §7, „Messung am Implementer-Stand“ und „Befund-Nachtrag (B)“.)* `make test-sdk-altserver`
      (neuer Runner `tools/harness/run-sdk-altserver-tests.sh`, Docker-only, Netz für
      Pulls und NuGet/PyPI/Maven-Restore; kein Gate) fährt dieselbe Compose-Umgebung wie
      die drei Realserver-Runner, **nur der Feed-Container ist das Image
      `ghcr.io/pt9912/pg-change-feed:0.5.0@<Index-Digest>`** (Override-Datei im
      Temp-Verzeichnis über `COMPOSE`, `compose.yaml` bleibt unverändert; `SDK_ALTSERVER_IMAGE`
      übersteuert das Image). Gedruckt und mit Exit ausgewertet: **B0** Schema-Vorbedingung —
      der Server 0.5.0 läuft gegen den Rollout des **Arbeitsbaums** (`make schema-rollout`,
      wie die anderen Runner): Slot da, Health `healthy`, eine eingefügte Zeile einer
      aktivierten Tabelle liegt über `cdc.changes` (die Zeile
      `ALTSERVER B0: healthy, Slot, change_id=<…>`); scheitert B0, endet der Slice mit einem
      Befund zur Schema-Kompatibilität, nicht mit einer Anpassung des Servers. **B1**
      Rohdraht-Probe der HTTP-Seite: der Fehlerkörper der Aktivierung einer fehlenden
      Tabelle am Server 0.5.0 trägt **kein** Feld `code` (gedruckt der rohe Körper; Abruf aus
      einem Container im Compose-Netz). **B2** je Sprache (die SDK-Quelle des Arbeitsstands,
      gebaut über die bestehende `integration`-Stufe) eine neue Test-Klasse
      (`ErrorCodeAltServerTests` C#, `ErrorCodeAltServerTest` Kotlin,
      `test_error_code_altserver.py` Python) gegen den Server 0.5.0: HTTP und gRPC enden mit
      dem typisierten `NotFound`-Fehler, `MessageCode`/`messageCode`/`message_code` ist
      `null`/`None`, der Fehlertext ist nicht leer und gleich dem rohen Text aus B1 (HTTP);
      das Reader-Token endet mit 403 bzw. `PermissionDenied` ohne Code; `Diagnose` (gRPC)
      liefert den Bericht ohne Absturz und `error_code` leer — in Normalbetrieb und in einem
      direkt geschriebenen Fehlerzustand (`error_class = 'schema'` in
      `cdc.process_heartbeat` **ohne** `error_code`, Muster
      `tools/harness/run-integration-tests.sh` §diagnose; der Alt-Server schreibt die
      Spalte nicht). Gedruckt je Sprache eine Zeile
      `RECEIVED code=none http=404 grpc=NOT_FOUND text=<n> diag_error_code=leer`. **B3**
      Negativprobe: `SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev` (der geladene
      0.6.0-Bau, `make image` vorher) färbt den Lauf rot (Exit ≠ 0 an B1, weil der Körper dort
      `code` trägt) — beweist, dass B2 nicht trivial grün ist; die Gegenseite
      (SDK 0.6.0 liest den Code am 0.6.0-Server) trägt unverändert die Fehlercode-Phase der
      drei bestehenden Runner und wird hier nicht wiederholt. Der Runner schreibt **nicht**
      nach `docs/user/sdk-e2e-abdeckung.md` (zu belegen: `git status --porcelain` vor und
      nach dem Lauf gleich). Die gRPC-Seite hat keine Rohdraht-Probe: die Abwesenheit des
      `ErrorInfo` am Server 0.5.0 ist aus dem Quelltext **abgeleitet** (§1, 0 Zeilen) und
      durch das Paar B2/B3 (SDK liest am 0.6.0-Server den Code, am 0.5.0-Server nicht)
      getragen.
- [x] **Verträge, Verdrahtung und Befund (Liefer-Punkt 3).**
      `harness/targets/sdk-kompat.md` und
      `harness/targets/sdk-altserver.md` (Form der Nachbarn unter `harness/targets/`:
      Aufruf, Eingaben, Ausgänge, Exit-Tabelle, Overrides, Grenze) tragen je Target den
      Vertrag; `harness/README.md` führt je Target eine Zeile in der Werkzeug-Tabelle
      (Form der Zeile zu `make test-sdk-python-integration`); `harness/mk/sdk.mk` trägt die
      zwei Ziele mit `.PHONY` und `##`-Hilfetext; kein neues `@sha256:`-Literal außer dem
      einen Image-Default des Runners (Suchlauf §3). Der Befund-Nachtrag in §7 nennt zu
      jeder der zwei Aussagen den Ausgang (**trägt** · **trägt nicht** · **trägt mit
      Einschränkung**) mit den gedruckten Zeilen, und bei einem Befund die berührte
      Entscheidung (nur benannt, kein ADR-Entwurf).

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0,
      `make sdk-public-doc-check` Exit 0 (die neuen Test-Klassen unter `sdks/` tragen keine
      interne Kennung), `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0,
      `make kommentar-kennungen DIFF=<Parent>` ohne Kandidat in den neuen Skripten.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: kein öffentlicher Vertrag berührt (die Handbücher und SDK-READMEs
      bleiben unverändert, zu belegen: der Diff trägt keine Datei unter `docs/user/`
      und keine README unter `sdks/`); `harness/README.md` und die zwei
      Target-Verträge sind Teil von Liefer-Punkt 3.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben: eine weitere Datei
      `evidence/slice-sdk-0-6-kompatibilitaet-messen.md` in
      [`BEO-PGC/kein-echter-versionswechsel-upgrade-test`](../observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md)
      (Zähler damit 2×, gezählt am Planungsstand: 1 Datei — der Slice liefert den ersten
      echten Alt-Server-gegen-Neu-SDK-Vergleich); dazu je nach Befund eine weitere Datei in
      [`BEO-PGC/adr-aussage-breiter-als-ihre-messung`](../observations/BEO-PGC/adr-aussage-breiter-als-ihre-messung/state.md)
      (11 Dateien am Planungsstand), wenn die Messung eine Aussage von
      [`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
      widerlegt oder einschränkt. **Kein Zähler wird gesetzt**, er folgt aus den Dateien.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt und „die
      nächste Welle-Closure“ damit keine Adresse ist. Ein Befund geht an den Planner (ein
      Folge-Slice oder eine ADR-Berührung wird in der Closure benannt, nicht vorab
      angelegt).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Der Implementer erweitert die Liste in seinem ersten Lauf; die Pfade der
Gast-Dateien sind Vorschläge, der Zuschnitt (Verzeichnis je Sprache) ist frei,
solange alles unter `tools/harness/sdk-kompat/` liegt.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/sdk-kompat/csharp/` (Gast-Projekt, Dockerfile) | neu | A1–A5 C#: Gast mit `PackageReference` auf `PgChangeFeed.Client` 0.5.0, Austausch der Assembly vor dem Lauf; `null`-Matrix als eigene Übersetzungs-Fälle (`LH-FA-SST-009`) |
| `tools/harness/sdk-kompat/kotlin/` (Gradle-Projekt zur Auflösung der Laufzeit-Klassen, Gast-Quelle, Dockerfile) | neu | A1–A5 Kotlin: Gast-Class-Files gegen das 0.5.0-Jar, Klassenpfad-Austausch; `javap`-Zeile belegt die gebundenen Deskriptoren |
| `tools/harness/sdk-kompat/python/` (Gast-Skript, Dockerfile) | neu | A1–A5 Python: Gast-Skript in einer virtuellen Umgebung mit `pgchangefeed==0.5.0`, danach 0.6.0 über derselben Umgebung |
| `tools/harness/run-sdk-kompat-tests.sh` | neu | Treiber von `make test-sdk-kompat`: baut die Basis aus den `build`-Stufen der SDK-Dockerfiles, fährt die Schritte A1 bis A5 je Sprache, druckt die `KOMPAT`-Zeilen, endet mit Exit ≠ 0 bei jeder Abweichung von der Erwartung des Schritts |
| `tools/harness/run-sdk-altserver-tests.sh` | neu | Treiber von `make test-sdk-altserver`: Bring-up wie `run-sdk-csharp-integration-tests.sh` (Zeilen 74–170: Compose, Rollout, Tabellen, Feed, Health), Override-Datei im Temp-Verzeichnis, B0 bis B3, Aufräumen im `trap` |
| `sdks/csharp/PgChangeFeed.Client.Integration/ErrorCodeAltServerTests.cs` | neu | B2 C# (Test-Quelle, nicht Teil des Packages; Muster `ErrorCodeRealserverTests.cs`, Phasen-Umgebung `PhaseEnvironment`) |
| `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/ErrorCodeAltServerTest.kt` | neu | B2 Kotlin (Muster `ErrorCodeRealserverTest.kt`) |
| `sdks/python/pgchangefeed/integration/test_error_code_altserver.py` | neu | B2 Python (Muster `test_error_code_realserver.py`) |
| `harness/mk/sdk.mk` | update | zwei Ziele `test-sdk-kompat`, `test-sdk-altserver` mit `.PHONY` und `##`-Text, nicht in `GATE_CHECKS` |
| `harness/targets/sdk-kompat.md`, `harness/targets/sdk-altserver.md` | neu | Verträge der Targets (Form der Nachbarn `harness/targets/image-mutation.md`, `harness/targets/bench-backfill.md`) |
| `harness/README.md` | update | je Target eine Zeile in der Werkzeug-Tabelle, Bindung auf den Vertrag und `ADR-0145` |
| `tools/harness/sdk-kompat/csharp/Gast/` (Projekt, Quelltexte `alt` und `neu`), `NullMatrix/`, `nuget-registry.config`, `nuget-dist.config`, `run.sh` | neu (Zuschnitt des Implementers) | Gast-Projekt mit zwei Quelltexten (0.5.x-Formen · Formen mit Code), die null-Matrix als reine Übersetzung, zwei NuGet-Konfigurationen (Modus `registry` und `dist` mit Quell-Zuordnung), Treiber der Schritte im Container (`LH-FA-SST-009`) |
| `tools/harness/sdk-kompat/kotlin/` (`settings.gradle.kts`, `build.gradle.kts`, `src/alt`, `src/neu`, `run.sh`, Dockerfile) und `tools/harness/sdk-kompat/python/` (`gast_alt.py`, `gast_neu.py`, `quelle.py`, `run.sh`, Dockerfile) | neu (Zuschnitt des Implementers) | Gast-Projekte der beiden anderen Sprachen und ihre Treiber; Python trägt die Quellseite A5 als eigenes Skript (`quelle.py`, `inspect.signature`) |
| `tools/harness/sdk-kompat/altserver/rohdraht.py` | neu (Zuschnitt des Implementers) | B1: die Rohdraht-Probe der HTTP-Seite, per `python -` im Python-Image des SDK |
| `sdks/csharp/PgChangeFeed.Client.Integration/PhaseEnvironment.cs`, `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/PhaseEnvironment.kt` | update | je eine Zeile für den rohen HTTP-Fehlertext aus B1 (`PGCHANGEFEED_ALTSERVER_HTTP_TEXT`); Test-Quelle, kein Package-Inhalt |
| `tools/harness/run-sdk-altserver-tests.sh` — Diagnose-Schalter `SDK_ALTSERVER_WEITER` | neu (über den Plan hinaus) | lässt B3 die Test-Klassen von B2 am `:dev`-Server zeigen (B1 bricht sonst vor B2 ab) |
| `docs/plan/planning/observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/evidence/slice-sdk-0-6-kompatibilitaet-messen.md` | neu | Register-Fortschreibung (§2) |
| `docs/plan/planning/observations/BEO-PGC/adr-aussage-breiter-als-ihre-messung/evidence/slice-sdk-0-6-kompatibilitaet-messen.md` | neu | Register-Fortschreibung bei Befund gegen `ADR-0145` (§2: „quellseitig erhalten“ trägt an einem Fall nicht) |

- **Messort ist der Container, der Host liefert nur Docker, `make`, `bash`, `git`
  und `mktemp`** ([`AGENTS.md`](../../../../AGENTS.md) §3.1). Keine
  `dotnet`-/`java`-/`python`-Aufrufe auf dem Host; Textänderungen an
  Repo-Dateien über Edit/Write, die Mutation von A4 ausschließlich an einer
  **Kopie im Scratchpad** (Edit/Write auf der Kopie, Bau über die Stufe `build` von
  `sdks/csharp/Dockerfile` mit der Kopie als Kontext und `--build-context proto=proto`;
  `make image` mit mutiertem Arbeitsbaum bleibt verboten).
- **Basis-Images ohne neuen Pin.** Die Gast-Dockerfiles bauen auf einem lokal
  getaggten Image der `build`-Stufe des jeweiligen SDK-Dockerfiles auf
  (`docker build --target build -t <lokaler Tag> --build-context proto=proto sdks/<sprache>`,
  Basis per `ARG`/`FROM` oder `--build-context …=docker-image://…`) und tragen **kein
  eigenes `@sha256:`-Literal**: ein weiterer Träger desselben Digests widerspräche
  [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  (je Image ein Wert an allen Trägern). Der Suchlauf unten zählt die
  `@sha256:`-Zeilen: am Diff genau eine mehr, der Image-Default des Altserver-Runners
  (hergeleitet; `make pin-stale-all` prüft ihn gegen den Index-Digest des Tags
  `0.5.0`, ein unveränderlicher Release-Tag — am Diff zu lesen, die Zeile
  `OK` ist Teil der Nachmessung in §7).
- **Austausch statt Neubau (A2).** C#: `PgChangeFeed.Client.dll` im Ausgabeordner des
  Gasts durch die Datei aus dem 0.6.0-`.nupkg` (`lib/net10.0/`) ersetzen, dann
  `dotnet <Gast>.dll` ohne `dotnet build`; Kotlin: das 0.5.0-Jar im Klassenpfad durch
  das 0.6.0-Jar ersetzen, Laufzeit-Abhängigkeiten aus der Auflösung der 0.6.0-Version
  (der Implementer stellt fest und druckt, ob die beiden Abhängigkeitsmengen gleich sind,
  [`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
  Festlegung 3 sagt „keine neue“ — *hergeleitet* aus dem Diff, hier gemessen);
  Python: `uv pip install pgchangefeed==0.6.0` (oder das Wheel aus `sdks/python/dist/`) in
  dieselbe Umgebung, derselbe Gast.
- **Ausführung per `docker run`, nie in einer `RUN`-Schicht.** Bau (`docker build`) und
  Lauf der Gast-Programme sind getrennt; ein Lauf innerhalb des Baus käme aus dem
  Schicht-Cache still durch
  ([`BEO-PGC/docker-cache-ueberspringt-tests-still`](../observations/BEO-PGC/docker-cache-ueberspringt-tests-still/state.md)).
  Jeder Schritt druckt seine Zeile bei jedem Lauf, ein Lauf ohne gedruckte Zeile ist rot.
- **Befund-Regel.** Zeigt A2, A5 oder B2 ein Ergebnis gegen die Erwartung, läuft die
  Messung zu Ende (alle Sprachen, alle Schritte), die Zeile steht ungeglättet in §7, und
  der Slice endet **ohne** Änderung an Produkt-Quellen; die Entscheidung über Fix,
  Folge-Slice oder ADR liegt beim Planner/Architect. Eine Probe wird nie angepasst, damit
  sie grün wird.

**Suchlauf** (§3.13 von [`AGENTS.md`](../../../../AGENTS.md); bewegte Eigenschaft: die
Aussagen „binär und quellseitig erhalten“ / „kein Code ohne Server-Code“ und die
Digest-Literale). Der Parent ist `ff97a91059f6f706b16cb16b0681a1563a064d83`
(`git rev-parse HEAD` am Planungsstand, vor dem Plan-Commit; nie `HEAD`). Suchraum:
ganzer Baum ohne `docs/reviews`, `done/`, `observations/`, `open/` (die Pläne tragen die
Suchmuster selbst) und `.harness/baseline`. Die Parent-Zeilen sind gemessen (Soll =
Trefferzeilen am Parent); die `diff`-Zeilen sind die **Erwartung** (hergeleitet, noch nicht
gemessen), der Implementer misst nach und trägt Gefundenes und Nichtgefundenes in §7 ein:

```suchlauf
ff97a91059f6f706b16cb16b0681a1563a064d83 4 -n -i -E 'binärkompat|binär und quellseitig|binary compat' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
ff97a91059f6f706b16cb16b0681a1563a064d83 2 -n -E 'quellseitig' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
ff97a91059f6f706b16cb16b0681a1563a064d83 85 -n -E '@sha256:' -- . ':!docs' ':!.harness'
ff97a91059f6f706b16cb16b0681a1563a064d83 0 -n -E 'test-sdk-(kompat|altserver)' -- . ':!docs/plan/planning/open'
ff97a91059f6f706b16cb16b0681a1563a064d83 4 -n -i -E 'binärkompat|binär und quellseitig|binary compat' -- docs/plan/adr
ff97a91059f6f706b16cb16b0681a1563a064d83 2 -n -E 'quellseitig' -- docs/plan/adr
diff 4 -n -i -E 'binärkompat|binär und quellseitig|binary compat' -- docs/plan/adr
diff 2 -n -E 'quellseitig' -- docs/plan/adr
diff 86 -n -E '@sha256:' -- . ':!docs' ':!.harness'
diff 6 -n -i -E 'binärkompat|binär und quellseitig|binary compat' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 4 -n -E 'quellseitig' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
```

  Die zwei letzten `diff`-Zeilen sind vom Implementer ergänzt: der neue Träger
  `harness/targets/sdk-kompat.md` zitiert die Aussage der ADR zweimal, der ganze
  Baum trägt damit 6 bzw. 4 Trefferzeilen (Parent 4 bzw. 2, gemessen).

  **Erwartung am `diff`-Stand** (hergeleitet): die Treffer der ersten zwei
  Parent-Zeilen liegen in den `Accepted` ADRs
  ([`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md) 2,
  [`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) 2);
  die `diff`-Zeilen mit dem Pfad `docs/plan/adr` halten fest, dass **kein ADR
  nachgezogen wird** (Soll gleich dem am Parent mit demselben Pfad gemessenen Wert
  4 bzw. 2). Die
  Zeile mit `@sha256:` erwartet genau **einen** neuen Treffer (der Image-Default des
  Altserver-Runners). Neue Träger der Aussagen (README-Zeilen, Target-Verträge) liegen
  außerhalb dieser Muster, soweit sie „Binärkompatibilität“ nicht als Aussage
  wiederholen — wer sie dort nennt, zählt die Zeile am Diff nach.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1; am Planungsstand trägt es nur `roadmap.md`), und am Arbeitsstand des
Implementers sind die Messwege aus §1 erneut geprüft: `docker buildx imagetools inspect
ghcr.io/pt9912/pg-change-feed:0.5.0` antwortet, `docker pull` des Digests gelingt und die drei
Gast-Quellen (NuGet, PyPI, Cloudsmith) liefern 0.5.0. Fehlt eine davon, wird der Slice nicht
gestartet (zurück nach `open/`, Grund in §7).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Mess-Stufen
  wachsen über die zwei Targets hinaus (zum Beispiel ein dritter Alt-Stand oder eine Änderung am
  bestehenden Runner nötig, damit B trägt), oder der Aufwand je Sprache verlangt eine eigene
  Zerlegung (A und B je ein Slice, A in die drei Sprachen). Zerlegungslinie (hergeleitet aus
  dem Zuschnitt in §3, nicht erprobt): A braucht kein Compose, B keinen Gast-Austausch; der
  Schnitt wäre sauber.
- `in-progress` → `open` (blockiert — Carveout?): ein Registry- oder Netzausfall
  (Docker-Hub-Abruflimit für `postgres`/`nats`, ghcr, NuGet, PyPI, Cloudsmith) hindert die
  Messung; kein Carveout, weil kein Gate rot ist. Gleiches gilt, wenn B0 (der Server 0.5.0 läuft
  nicht gegen das aktuelle Schema) eintritt: das ist ein Befund über die Schema-Kompatibilität
  von Server und Rollout, der den Slice auf den Teil A begrenzt und B an den Planner
  zurückgibt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der Closure, beide
neuen Targets haben je einen Lauf mit gedruckten Zeilen und Exit 0 (Erwartungsschritte) bzw.
Exit ≠ 0 (Negativ- und Mutationsproben) in §7, und die Closure-Notiz in §7 trägt den
Befund-Nachtrag zu beiden Aussagen und einen Lerneintrag (geschärfte Regel, neuer Sensor oder
benannte Spec-Lücke; die zwei neuen Targets sind Sensoren im Sinn des Lerneintrags, wenn eine
Aussage damit von *hergeleitet* auf *gemessen* geht). Ein Befund gegen die Erwartung
verhindert die Closure nicht, solange er benannt, belegt und an den Planner adressiert ist;
ein Gate, das am Stand der Closure rot ist, geht nur mit dokumentiertem Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen (Grund) ·
weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur Closure **offen**
(Platzhalter `Ausgang: offen bis Closure`).

- **Die Messung ist trivial grün** (der Gast bindet nicht an die alte Signatur, der
  Austausch ersetzt nichts, der Alt-Server ist in Wahrheit der `:dev`-Bau). — **Ausgang:**
  offen bis Closure; Gegenmittel: A3 (Gegenrichtung), A4 (Mutation), B1 (Rohdraht), B3 (`:dev`
  muss rot färben), `javap`-Zeile des Kotlin-Gasts.
- **Der Befund widerlegt eine Aussage von
  [`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)** (am
  wahrscheinlichsten: die `null`-Mehrdeutigkeit der protected Basis, F-1). — **Ausgang:**
  offen bis Closure; der Befund steht ungeglättet in §7, die berührte Entscheidung wird
  benannt; Änderung an `Accepted` ADR nur als neue ADR mit `Supersedes`
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5), nicht in diesem Slice.
- **Der Server 0.5.0 läuft nicht gegen das aktuelle Schema** (B0). — **Ausgang:** offen bis
  Closure; die Aussage „additive nullable Spalte“ ist *hergeleitet*; tritt es ein, ist es ein
  Befund über die Rollout-Reihenfolge Schema vor Server, kein Fehler dieses Slice.
- **Die Dependency-Mengen der Kotlin-Auflösung unterscheiden sich zwischen 0.5.0 und 0.6.0**
  (der Austausch des Jars lädt dann gegen eine Teilmenge). — **Ausgang:** offen bis Closure;
  der Implementer druckt beide Mengen (Schritt A2) und löst für den Lauf gegen 0.6.0 die
  Mengen der 0.6.0-Version auf.
- **Registry-/Netzgrenzen** (Docker-Hub-Abruflimit für `postgres:18-alpine` und
  `nats:2-alpine`, das der Altserver-Runner wie jeder Compose-Runner zieht; ghcr, NuGet, PyPI,
  Cloudsmith, Maven Central). — **Ausgang:** offen bis Closure; erste Linie: Images im lokalen
  Cache (am Planungsstand `:dev` und 0.5.0 geladen), UNBESTIMMT ist kein Erfolg.
- **Zwei Runner teilen Container-Namen und Netz** (`cdc-test-feed`, `cdc-feed-test`): ein
  gleichzeitiger Lauf mit `make test-integration` oder einem SDK-Runner kollidiert. —
  **Ausgang:** offen bis Closure; der Runner bricht mit klarer Meldung ab, wenn die Namen
  belegt sind (Aufräumen nur im eigenen `trap`).
- **Der Altserver-Runner überschreibt `docs/user/sdk-e2e-abdeckung.md`** (die drei
  Realserver-Runner schreiben ihren Abschnitt dorthin). — **Ausgang:** offen bis Closure;
  DoD verlangt `git status --porcelain` vor und nach dem Lauf gleich.
- **Das `.dockerignore`-Verhalten der neuen Bau-Kontexte** (Beobachtung
  [`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`](../observations/BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad/state.md)):
  die Gast-Verzeichnisse sind eigene Kontexte; `sdks/*/` führt am Planungsstand keine
  `.dockerignore` (gemessen: `ls -a` je Verzeichnis, kein Treffer). — **Ausgang:** offen bis
  Closure.
- **Interne Kennung in öffentlichem Text:** die Test-Klassen unter `sdks/` und ihre
  Fehlertexte erreichen keinen Anwender (kein Package-Inhalt), `make sdk-public-doc-check`
  prüft sie dennoch. — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Wird mit der Closure gefüllt (Inhalt, dann `git mv`, dann Häkchen der Paarungs-Zeile —
[`AGENTS.md`](../../../../AGENTS.md) §3.3). Hier stehen: der **Befund-Nachtrag** zu den
zwei Aussagen von
[`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) (je
Aussage und Sprache: Ausgang *trägt* / *trägt nicht* / *trägt mit Einschränkung*, die
gedruckten Zeilen mit Lauf und Parent-Kennung, Ursprung gemessen/übernommen/abgeleitet),
die Messzeilen A1 bis A5 und B0 bis B3, die Nachmessung des Suchlaufs (Gefundenes und
Nichtgefundenes) und, falls ein Befund eine Entscheidung berührt, deren Benennung.

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** wird mit der Closure gefüllt.
- **Was ging anders als geplant:** wird mit der Closure gefüllt.
- **Messung am Implementer-Stand** (Läufe am Arbeitsbaum auf dem Parent `5b0ad92a`, nach den Commits `f238bc3e` und `6aa987db`, Messhost `x86_64`; jede Zeile *gemessen*, gedruckt von `make test-sdk-kompat` bzw. `make test-sdk-altserver`, Exit direkt gelesen):
  - `make test-sdk-kompat` (Modus `dist`, Standard) endet mit Exit 0, Schlusszeile `run-sdk-kompat-tests: Kompatibilitätsmessung (dist) grün für: csharp kotlin python`; `SDK_KOMPAT_NEU=registry` dasselbe (Exit 0, `… (registry) grün für: csharp kotlin python`).
  - A1/A2 (Zeilen `KOMPAT <sprache> A1|A2: <n> Aufrufe ok`): C# 53 Aufrufe an beiden Schritten (`A2: Bibliothek 0.6.0 (Artefakt PgChangeFeed.Client.0.6.0.nupkg) f8e6415b23e5 ersetzt 0.5.0 6a4992f3f08f` — die Dateien sind verschieden, der Austausch geschah; im Modus `registry` trägt die 0.6.0-Datei dieselbe Kurz-Prüfsumme `f8e6415b23e5`), Kotlin 45 (`A2: Bibliothek 0.6.0 (Artefakt pgchangefeed-kotlin-0.6.0.jar) 2f600c1e671c ersetzt 0.5.0 4380d0d0664b`; die `javap`-Zeile des Gasts nennt den gebundenen Deskriptor `PgChangeFeedBadRequestException."<init>":(ILjava/lang/String;)V`, den die Bibliothek 0.6.0 neben `(ILjava/lang/String;Ljava/lang/String;)V` weiter trägt), Python 45 (`A2: Bibliothek pgchangefeed 0.6.0 (Artefakt pgchangefeed-0.6.0-py3-none-any.whl) ersetzt 0.5.0`, dieselbe virtuelle Umgebung).
  - A3: Grundlage `A3-Grundlage: 2 Aufrufe ok` je Sprache; Gegenrichtung gegen 0.5.0 scheitert — C# `Ausnahme FileNotFoundException: Could not load file or assembly 'PgChangeFeed.Client, Version=0.6.0.0 …'`, Kotlin `Ausnahme NoSuchMethodError: 'void io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException.<init>(int, java.lang.String, java.lang.String)'`, Python `Ausnahme TypeError: PgChangeFeedError.__init__() got an unexpected keyword argument 'message_code'`.
  - A5: C# und Kotlin übersetzen den Gast-Quelltext (0.5.x-Formen) gegen 0.5.0 und gegen 0.6.0 (`A5 quelle 0.5.0|0.6.0: … übersetzt`); die Abhängigkeitsmengen sind gleich (C# `A5 abhaengigkeiten: 0.5.0 und 0.6.0 gleich (3 Zeilen)`, Kotlin `… gleich (28 Jars)`); Python liest 15 Typen unter 0.5.0 und unter 0.6.0 (`A5 signatur 0.6.0 <Typ>: 0.5.x-Form bindet=True, drittes Positional bindet=False, message_code keyword-only=True, Standard None=True wie-erwartet`). Die C#-null-Matrix (14 Fälle, je Version): 13 Fälle `0.5.0 ok, 0.6.0 ok`, der Fall `http-basis-3-argumente-null-literal` `0.5.0 ok, 0.6.0 CS0121`.
  - `make test-sdk-altserver` endet mit Exit 0: `ALTSERVER B0: healthy, Slot slot_pgc_e2e, change_id=812-1, Feed-Image ghcr.io/pt9912/pg-change-feed:0.5.0@sha256:f99a77ff…` (Digest im Lauf voll gedruckt, hier gekürzt: dieser Plan trägt kein zweites Digest-Literal), B1 `BODY {"error":"Tabelle existiert nicht an der Quelle: public.sdk_error_code_missing_table"}`, `STATUS 404`; B2 je Sprache `RECEIVED code=none http=404 grpc=NOT_FOUND text=74 diag_error_code=leer (Exit 0)` (C#, Kotlin, Python). `git status --porcelain` ist vor und nach dem Lauf gleich (gemessen: zwei Ausgaben, `cmp` ohne Abweichung).
  - B3: `SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev make test-sdk-altserver` endet mit Exit ≠ 0 (über `make` 2): `ALTSERVER B1: BODY {"error":"Tabelle existiert nicht an der Quelle: public.sdk_error_code_missing_table","code":"PCF-E8025"}`, `run-sdk-altserver-tests: B1 ROT — der Fehlerkörper der Aktivierung einer fehlenden Tabelle ist nicht der eines Servers ohne Meldungscodes (Ausgang 1)`; mit `SDK_ALTSERVER_WEITER=1` zusätzlich `B2 csharp ROT`, `B2 kotlin ROT`, `B2 python ROT` (je „Marker NORMAL_DONE blieb aus“) und Exit ≠ 0.
  - `make pin-stale-all`: die Zeile `OK          tools/harness/run-sdk-altserver-tests.sh:35 (ghcr.io/pt9912/pg-change-feed:0.5.0) == sha256:f99a77ff…` (gemessen); der Gesamtlauf endet mit Exit ≠ 0 wegen eines fremden `DRIFT` (`sdks/python/Dockerfile:34 (python:3.14-slim)`, Pin unverändert durch diesen Slice, `15 OK, 1 DRIFT, 0 UNBESTIMMT`) — Hinweis an den Planner, kein Befund dieses Slice.
- **Mutationen (Zusage · mutierte Eingabe · gesehenes Rot; Instanz je Zeile; Kopien im Scratchpad, Bau über `docker build --target pack-export` der Kopie, Messung über `SDK_KOMPAT_DIST_<SPRACHE>`):**

  | Zusage | mutierte Eingabe | Instanz | gesehenes Rot |
  |---|---|---|---|
  | A4 C#: der 0.5.x-Konstruktor bleibt binär erhalten | `PgChangeFeedBadRequestException(int, string)` aus der Kopie der Bibliothek entfernt (Test der Kopie angepasst) | C#-Bibliothek 0.6.0, Kopie von `sdks/csharp` | `KOMPAT csharp A2: Ausnahme MissingMethodException: Method not found: 'Void PgChangeFeed.Client.Http.PgChangeFeedBadRequestException..ctor(Int32, System.String)'`, Exit 2 über `make` |
  | A4 Kotlin: die JVM-Signatur `(int, String)` bleibt erhalten | `@JvmOverloads` an `PgChangeFeedBadRequestException` entfernt | Kotlin-Bibliothek 0.6.0, Kopie von `sdks/kotlin` | `KOMPAT kotlin A2: Ausnahme IllegalAccessError: class kompat.GastKt tried to access private method 'void io.github.pt9912.pgchangefeed.http.PgChangeFeedException.<init>(int, java.lang.String)'`, Exit 2 |
  | A4 Python: die 0.5.x-Aufrufform bleibt gültig | `PgChangeFeedBadRequestError.__init__` verlangt `message_code` ohne Standard | Python-Bibliothek 0.6.0, Kopie von `sdks/python` | `KOMPAT python A2: Ausnahme TypeError: PgChangeFeedBadRequestError.__init__() missing 1 required keyword-only argument: 'message_code'` und `A5 signatur 0.6.0 PgChangeFeedBadRequestError: 0.5.x-Form bindet=False … ROT`, Exit 2 |
  | B1: der Server trägt kein `code` | Feed-Image `:dev` statt 0.5.0 | Server mit Meldungscodes | `ALTSERVER B1 FEHLER: der Fehlerkörper trägt das Feld code`, Exit 2 |
  | B2: die Eigenschaft bleibt leer, in allen drei Sprachen | Feed-Image `:dev`, `SDK_ALTSERVER_WEITER=1` | Server mit Meldungscodes | B2 C#, Kotlin, Python je ROT, Exit 2 |

  Menge: je eine Mutation je Sprache und eine Server-Mutation, je einmal gefahren; die Übertragung auf andere Blatt-Typen und auf die gRPC-Hierarchie ist *hergeleitet*. Die erste Python-Mutation (Pflicht-`message_code` an der Basis) färbte schon den Unit-Lauf im Docker-Bau rot (`TypeError … missing 1 required keyword-only argument` in `http_client.py`), bevor der Gast lief; sie wurde auf das Blatt `PgChangeFeedBadRequestError` verlegt. Die Mutation der Kotlin-Bibliothek endet mit `IllegalAccessError`, nicht mit dem erwarteten `NoSuchMethodError`.
- **Suchlauf nachgemessen** (`make suchlauf-nachmessen PLAN=<diese Datei>`, Exit 0, `suchlauf-nachmessen: 11 Zeilen stimmen`, am Stand nach `6aa987db` mit den zwei ergänzten `diff`-Zeilen): `@sha256:` am `diff` 86 (erwartet 86, genau ein neuer Treffer — `tools/harness/run-sdk-altserver-tests.sh`); ADR-Pfad `docs/plan/adr` 4 und 2 (unverändert); die Neuzeile der Digest-Referenz in `harness/targets/`: 0 Treffer (die Verträge nennen den Tag). **Gefunden:** der neue Träger der Aussage „binär und quellseitig erhalten“ ist `harness/targets/sdk-kompat.md` (2 Zeilen, ein Zitat der Aussage der ADR und ein Verweis auf sie; am ganzen Baum ohne den Plan 6 bzw. 4 Trefferzeilen); **nicht gefunden:** weitere Träger der Aussage in `sdks/*/README`, `docs/user/` und `spec/` (die Treffer des ganzen Baums liegen in `docs/plan/adr/` und in `harness/targets/sdk-kompat.md`).
- **Befund-Nachtrag (A, Kompatibilität):**
  - *„binär erhalten“* — **trägt**, gemessen in drei Sprachen: dieselben Gast-Binärdateien (C#-Assembly, Kotlin-Class-Files, Python-Skript) laufen ohne Neubau gegen die Bibliothek 0.6.0 (A2, Zeilen oben, beide Modi); die Mutation einer entfernten Signatur färbt A2 rot (A4, drei Sprachen). Für Python gilt Aufruf-Kompatibilität (keine Binärdatei).
  - *„quellseitig erhalten“* — **trägt mit Einschränkung:** C# und Kotlin übersetzen den Gast-Quelltext (0.5.x-Formen) gegen beide Versionen, Python bindet die 0.5.x-Form; **eine** Form trägt nicht: ein `base(status, text, null)` einer fremden Unterklasse der abstrakten HTTP-Basis `PgChangeFeedException` übersetzt unter 0.5.0 und scheitert unter 0.6.0 mit `CS0121` (gedruckt `A5 null-matrix http-basis-3-argumente-null-literal: 0.5.0 ok, 0.6.0 CS0121`; Befund `F-1` des Reviews zum SDK-Slice, jetzt gemessen). Die Umgehung mit Cast oder benanntem Argument übersetzt unter beiden (Zeilen `…-null-cast-exception`, `…-benannt-innerException-null`); die öffentlichen Blatt-Konstruktoren und die gRPC-Basis sind nicht betroffen. **Berührte Entscheidung:** `ADR-0145` Festlegung 4 (Satz „quellseitig erhalten“, `Accepted`); nur benannt — Fix, Folge-ADR (`Supersedes`) oder Annahme der Einschränkung liegt beim Planner/Architect, dieser Slice ändert weder ADR noch SDK.
  - *Abweichungen von den Erwartungen des Plans* (ungeglättet): (1) A3 in C# endet mit `FileNotFoundException`, nicht mit der erwarteten `MissingMethodException` — die Laufzeit weist die Assembly 0.5.0 wegen ihrer kleineren Version ab, bevor sie die Signatur sucht; die `MissingMethodException` zeigt erst die Mutation A4 bei gleicher Version. Das Kriterium von A3 steht deshalb als Alternative (`MissingMethodException|FileNotFoundException|FileLoadException`); der Plan-Satz „muss mit `MissingMethodException` scheitern“ gilt nicht wörtlich. (2) A4 in Kotlin endet mit `IllegalAccessError`, nicht mit `NoSuchMethodError`. (3) Die Übersetzung des Kotlin-Gasts braucht `io.grpc:grpc-api` ausdrücklich: die Fehlertypen tragen `io.grpc.Status` in ihrer Signatur, die Bibliothek veröffentlicht die gRPC-Abhängigkeiten nur für die Laufzeit (`implementation`) — dieselbe Eigenschaft in 0.5.0 und 0.6.0, kein Befund der Versionen, aber eine Information für Anwender (Hinweis an den Planner). (4) Das Kotlin-Jar der Registry und das Artefakt aus `make sdk-pack-kotlin` haben verschiedene Prüfsummen (`660d7b9db9ba` und `2f600c1e671c`, gemessen in den zwei Modi), das C#-Paket dieselbe: A2 trägt in beiden Modi.
- **Befund-Nachtrag (B, Server vor 0.6.0):** **trägt**, in drei Sprachen gemessen: gegen das veröffentlichte Image 0.5.0 (Index-Digest aus dem Plan, im Lauf gedruckt) enden HTTP und gRPC mit dem typisierten `NotFound`-Fehler, der Meldungscode ist `null`/`None`, der Fehlertext ist da (HTTP: gleich dem rohen Text aus B1, Länge 74), das Reader-Token endet mit 403 bzw. `PermissionDenied` ohne Code, die Diagnose liefert den Bericht ohne Absturz mit leerem `error_code` im Normalbetrieb und im Fehlerzustand (Zeilen oben). B0 trägt: der Server 0.5.0 läuft gegen das Schema des Arbeitsbaums (die Aussage „additive nullable Spalte“ ist damit für diese Paarung gemessen, nicht mehr hergeleitet). Die gRPC-Seite hat keine Rohdraht-Probe: die Abwesenheit des `ErrorInfo` am Server 0.5.0 ist aus dem Quelltext abgeleitet (`git grep -n -E 'errdetails|ErrorInfo' v0.5.0 -- internal cmd ':!*_test.go'` druckt 0 Zeilen, gemessen) und durch das Paar B2/B3 getragen. Der Plan nennt am Tag `v0.6.0` 2 Dateien; die Nachmessung (ohne Test-Dateien) trifft 3 (`internal/adapters/driving/grpc/administration.go`, `…/errors.go`, `internal/domain/messagecode/codes.go`) — kein Einfluss auf die Aussage.
- **Steering-Loop-Eintrag:** wird mit der Closure gefüllt (Kandidat: die zwei Targets als
  neue Sensoren, die zwei *hergeleiteten* Aussagen auf *gemessen* heben; liegt in dem
  Makefile-Ziel und dem Vertrag unter `harness/targets/`).
- **Beobachtungs-Register (`../observations/`):** wird mit der Closure gefüllt.
- **Folge-Slices:** wird mit der Closure gefüllt (nur bei Befund).
- **Risiken aus §6:** wird mit der Closure gefüllt — jedes mit genau einem Ausgang, siehe §6.
- **Drei Paarungen:** wird mit der Closure gefüllt — Anker · Folge-Slice · Register,
  getragen von der Slice-Closure selbst.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `tools/harness/` (neue
Runner und Gast-Quellen), `sdks/*` (nur Test-Quellen der Integrations-Quellsätze) und
`harness/` (Mk-Datei, Verträge, README-Zeilen) — alle unter der Default-Sub-Area `*`
(`PGC`, Greenfield; `harness/conventions.md` §Modus-Deklaration pro Sub-Area, gelesen: eine
Zeile). Die Schwelle ≥ 2 von 3 Achsen ist nicht gemessen und braucht keine Aufteilung:
der Slice ändert keine Produkt-Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) am Planungsstand durchgegangen; Treffer mit
Zähler-Stand (Zahl der Dateien unter `evidence/`):
[`kein-echter-versionswechsel-upgrade-test`](../observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md)
1× (der Slice ist ein echter Alt-Server-Vergleich und liefert die zweite Datei, §2);
[`adr-aussage-breiter-als-ihre-messung`](../observations/BEO-PGC/adr-aussage-breiter-als-ihre-messung/state.md)
11× (Zustand verkörpert; eine weitere Datei nur bei einem Befund gegen
[`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md));
[`docker-cache-ueberspringt-tests-still`](../observations/BEO-PGC/docker-cache-ueberspringt-tests-still/state.md)
2× (die Gast-Läufe geschehen per `docker run`, nie in einer `RUN`-Schicht des Baus — sonst
übersprünge der Cache die Messung still; eine dritte Datei entstünde, wenn das
trotzdem eintritt);
[`integrationsprojekt-uebersetzt-nicht-unbemerkt`](../observations/BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt/state.md)
1× (die neuen C#-Test-Klassen werden erst vom Runner übersetzt; die Übersetzung ist
Teil des Laufs von B, nicht des Gates);
[`dockerignore-default-deny-blockiert-neuen-pfad`](../observations/BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad/state.md)
1× (§6). Weitere Einträge treffen keine berührte Sub-Area.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

