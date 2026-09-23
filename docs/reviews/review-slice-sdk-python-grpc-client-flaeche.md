# Review-Report: slice-sdk-python-grpc-client-flaeche — 2026-09-23

**Review-Art:** Code-Review — geprüft gegen den Slice-Plan (§2 DoD, §3 Plan +
Plan-Nachzug), [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Accepted), [ADR-0107](../plan/adr/0107-python-pypi-zweites-sdk-package.md),
[ADR-0108](../plan/adr/0108-python-sdk-uv-statt-build-twine.md),
[ADR-0090](../plan/adr/0090-beispiel-clients-volle-matrix.md),
[ADR-0044](../plan/adr/0044-image-beleg-semantik.md),
[`SPEC-020`](../../spec/pflichtenheft.md), [`SPEC-027`](../../spec/pflichtenheft.md),
[`LH-FA-SST-008`](../../spec/lastenheft.md),
[`LH-FA-SST-009`](../../spec/lastenheft.md), `AGENTS.md` §3 (Hard Rules),
`harness/conventions.md` (MR-000). **Nicht geprüft gegen die DoD-Substanz** —
das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** `git diff a13c1d03..HEAD` — Implementations-Commit
`f1ce9be4` (feat(sdk-python): gRPC-Stream-Client-Flaeche +
Realserver-Integrationstest-Werkzeug); davor liegen nur reine Move-Commits des
Slice-Plans (`slice-sdk-python-grpc-client-flaeche`, Lifecycle
`in-progress/`).

**Skill:** `.harness/skills/reviewer.md` (Stand der Schärfung 2026-09-09, inkl.
`AGENTS.md` §3.12-Instanz-A-HIGH, §3.13-Träger-Nachzug) · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** glm-5.3-flash (Claude-Agent-SDK, Reviewer-Rolle) · **Datum:**
2026-09-23

**Eingangs-Kontext** (Verträge, gegen die geprüft wurde):

- `Slice-Plan` `slice-sdk-python-grpc-client-flaeche`
  (§2 DoD, §3 Plan + Plan-Nachzug, §6 Risiken)
- [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  (Umfang, Festlegung 2 — verschärfte Test-Pflicht, Import-Grenze),
  [ADR-0107](../plan/adr/0107-python-pypi-zweites-sdk-package.md)
  (Ort, Import-Grenze, unverändert gültig),
  [ADR-0108](../plan/adr/0108-python-sdk-uv-statt-build-twine.md)
  (Bau-/Publish-Frontend, unberührt),
  [ADR-0090](../plan/adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 2 (benannter Bau-Kontext),
  [ADR-0044](../plan/adr/0044-image-beleg-semantik.md) (Image-Beleg)
- [`SPEC-020`](../../spec/pflichtenheft.md) (gRPC-Live-Change-Stream —
  Nachrichtenschema, RPC-Name, Stream-Semantik),
  [`SPEC-027`](../../spec/pflichtenheft.md) (Python-Package),
  [`LH-FA-SST-008`](../../spec/lastenheft.md),
  [`LH-FA-SST-009`](../../spec/lastenheft.md)
- `AGENTS.md` §3.7/§3.9/§3.10/§3.11/§3.12/§3.13, `harness/conventions.md`
  (MR-000 ID-Schema)
- Struktur-Vorbilder: `sdks/csharp/Dockerfile`,
  `tools/harness/sdk-pack-csharp.sh`,
  `tools/harness/run-integration-tests.sh` (gRPC-Stream-Rundlauf, Zeile
  ~2260 ff.), `compose.yaml` (Container-Vertrag), `Makefile` Zeile 224 ff.
  (`SCHEMA_TARGET`/`SCHEMA_ROLLOUT_NETWORK`)
- Vorherige Findings am selben Modul:
  [review-slice-sdk-python-pack-werkzeug.md](review-slice-sdk-python-pack-werkzeug.md),
  [review-slice-sdk-python-projektgeruest.md](review-slice-sdk-python-projektgeruest.md)

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — DoD-Zahl „5 neue gRPC-Tests" driftet gegen die Messung (real 6)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 Instanz A — Zahl im Träger driftet gegen die
  Messung (Herkunft der Klasse:
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, 4×)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md:79`
- `befund`: Der DoD-Sensor-Beleg behauptet „29 Unit-Tests grün, darunter **5
  neue** gRPC-Tests `tests/test_grpc_client.py`"; die Datei
  `sdks/python/pgchangefeed/tests/test_grpc_client.py` trägt real **sechs**
  `def test_`-Funktionen (`grep -c "^def test_"` = 6), und die 29 ergeben
  sich nur mit sechs Neuen (20 + 3 + 6) — die Zahl 5 ist gegen dieselbe
  Messung, aus der die 29 stammt, falsch.
- `verifizierbar`: ja — `grep -c "^def test_"
  sdks/python/pgchangefeed/tests/*.py` (6/20/3) bzw. die pytest-Zählung des
  `make sdk-pack-python`-Laufs
- `klasse`: „Zahl im Träger ohne Ursprung — oder gegen die Messung driftend"

### F-2 — Feldvollständigkeits-Assertion ist für 3 der 10 Felder vakuum (bytes/int gegen `""`)

- `kategorie`: HIGH
- `quelle`: Beleg trägt seinen Satz nicht
  (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` — Assertion, die den
  behaupteten Satz nicht misst); Verwandtschaft: „Zusage ohne Bindung an
  ihre Eingabeseite"
- `pfad`: `sdks/python/pgchangefeed/integration/test_grpc_realserver.py:59-71`
- `befund`: Der Testkommentar behauptet „Feldvollständigkeit am realen
  Wire-Image ([`SPEC-020`](../../spec/pflichtenheft.md)): alle zehn Felder getragen, Identitäten non-empty";
  die Assertion `assert getattr(received, field) != ""` kann für die drei
  Nicht-String-Felder `old_image`/`new_image` (bytes) und `sequence` (int)
  **nie rot werden** — `b"" != ""` und `0 != ""` sind in Python stets wahr,
  die Behauptung „alle zehn Felder" wird von dieser Assertion für 3 der 10
  Felder nicht getragen (die sieben String-Identitäten sind korrekt
  gebunden; `new_image` trägt der Runner-seitige Sentinel-Grep separat,
  `sequence` und `old_image` am INSERT nirgends).
- `verifizierbar`: ja — Mutationsprobe am Wire-Image (leeres `new_image` /
  `sequence=0` in der empfangenen Message) lässt die Assertion grün; die
  Probe ist die Mutation selbst (kein Mutations-Harness im Repo)
- `klasse`: „Beleg trägt seinen Satz nicht" (Assertion, die ihren Satz nicht
  trägt)

### F-3 — Handbuch-Absatz zum Python-SDK bleibt falsch stehen (gebündelter Aufschub vs. ADR-0110 Folgepflicht 3)

- `kategorie`: MEDIUM
- `quelle`: [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Konsequenzen Folgepflicht 3 („SDK-Hinweis je neu gedeckter Oberfläche, im
  selben Zug wie das jeweilige Client-Programm") · `AGENTS.md` §3.13
- `pfad`: `docs/user/benutzerhandbuch.md:682-684`; Abgrenzungs-Quelle:
  `docs/plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md:57-59`
- `befund`: Der Python-SDK-Absatz im Handbuch trägt weiter die Satzform
  „gRPC, SSE und der NATS-Vollinhalts-Stream bleiben für dieses Package
  vorerst außerhalb (`ADR-0107` Festlegung 1)" — genau die Satzform, die
  dieser Slice falsch macht, und sie beruft sich auf `ADR-0107` Festlegung 1,
  die [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  teilweise superseded hat; der Aufschub ist **benannt mit Adresse** (letzter
  Flächen-Slice, Slice-Plan §1 + Welle-Plan §4), aber die C#/Kotlin-Präzedenz
  dieser Klasse korrigierte den Handbuch-Absatz **im selben Slice**
  (Änderungshistorie 1.34/1.37 — jedes Mal mit der nachträglichen
  Korrektur-Begründung „behauptete fälschlich …"), und der
  Welle-Closure-Trigger (`welle-sdk-python-vollabdeckung.md` §3) nennt nur
  `SPEC-027`/`LH-FA-SST-009.a` als Träger-Nachzug, nicht das Handbuch —
  der gebündelte Aufschub hängt damit nur am Slice-Plan des künftigen
  NATS-Slices.
- `verifizierbar`: ja — `grep -n "vorerst" docs/user/benutzerhandbuch.md`
  (Python-Absatz unverändert; beide Stände gemessen: Parent `a13c1d03` und
  HEAD identisch falsch)
- `klasse`: „Träger-Nachzug gebündelt statt im Zug" (Familie
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  hier mit benanntem Aufschub — deshalb MEDIUM, nicht HIGH)

### F-4 — Öffentlicher Parameter `stream_changes(timeout=…)` ohne Bindung: Fake verwirft `timeout`, kein Test kann eine Weiterleitung sehen

- `kategorie`: MEDIUM
- `quelle`: fehlende Negativtests bei neuem öffentlichem Vertrag; Verwandtschaft:
  „Zusage ohne Bindung an ihre Eingabeseite" (der Docstring verspricht
  „raises `DEADLINE_EXCEEDED` once it lapses" — kein Test trägt das)
- `pfad`: `sdks/python/pgchangefeed/tests/test_grpc_client.py:659-663`
  (Fake-`call` nimmt `timeout` an, zeichnet ihn nicht auf) ·
  `sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py:58`
- `befund`: Der Plan-Nachzug führt `stream_changes(timeout=None)` als
  öffentlichen Parameter ein; kein Unit-Test bindet ihn an die Eingabeseite —
  der Fake-Channel akzeptiert `timeout` in `call(...)`, legt ihn aber nicht
  in `invocations` ab, und der Integrationstest (eigene
  `time.monotonic()`-Frist) kann eine weggeworfene `timeout`-Weiterleitung
  nicht von einer geborgenen unterscheiden: eine Regression, die `timeout`
  still fallen lässt, bleibt grün.
- `verifizierbar`: ja — Mutationsprobe: `timeout=timeout` im Aufruf
  `grpc_client.py:66` streichen, alle Tests bleiben grün
- `klasse`: „Fehlende Negativtests bei neuem öffentlichem Vertrag" (einziger
  ungebundener Parameter der neuen Fläche)

### F-5 — Planzeile des Folge-Slices bereits durch diesen Slice überholt (options.py-Docstring)

- `kategorie`: INFO
- `quelle`: Maintainability (Zwei-Quellen-Drift, harmlose Form — Plan-Stand)
- `pfad`: `docs/plan/planning/open/slice-sdk-python-nats-stream-client-flaeche.md:127`
- `befund`: Der offene Plan des NATS-Folge-Slices trägt eine Planzeile
  „options.py — Docstring-Satz ‚no gRPC/SSE/NATS surface exists yet'
  entfernen/korrigieren (`AGENTS.md` §3.13)"; dieser Slice hat die Korrektur
  bereits real ausgeführt (`sdks/python/pgchangefeed/src/pgchangefeed/options.py:1-13`)
  — die Zeile des künftigen Slices ist überholt und müsste beim nächsten
  Planner-Zug reduziert werden, damit der Folgelauf nicht einen Träger
  „korrigiert", den es in dieser Form nicht mehr gibt.
- `verifizierbar`: ja — `grep -n "options.py"
  docs/plan/planning/open/slice-sdk-python-nats-stream-client-flaeche.md`
  gegen den HEAD-Stand der Docstring-Datei
- `klasse`: „Arbeit überholt stehenden Träger" (Planform)

### F-6 — Pfad-Angabe im Plan-Nachzug zeigt auf eine nicht existierende Datei

- `kategorie`: LOW
- `quelle`: Maintainability (Tippfehler-Klasse)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md:162`
- `befund`: Die Plan-Nachzug-Zeile nennt `sdks/python/pyproject.toml` als
  Träger der `protobuf>=6`-Abhängigkeit; die reale Datei (auch in der §3-Planzeile
  korrekt genannt) ist `sdks/python/pgchangefeed/pyproject.toml` — der
  Nachzug-Pfad existiert nicht.
- `verifizierbar`: nein — reine Pfad-Probe (`ls sdks/python/pyproject.toml` →
  fehlt)
- `klasse`: Tippfehler (Pfad)

---

## Negativbefunde (geprüft, ohne Befund)

- **`sdks/python/pgchangefeed/src/pgchangefeed/`** — Import-Grenze
  ([ADR-0107](../plan/adr/0107-python-pypi-zweites-sdk-package.md)
  Festlegung 3, [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 4) real geprüft: `grep -rn "internal/\|cmd/\|gen/" sdks/python/`
  liefert ausschließlich Prosakommentare (pyproject-Kommentar, `grpc_gen`
-Docstring, `.gitignore`-Kommentar, Dockerfile-Kommentar, `.gitignore`
  -Muster) und **keinen** Python-Import — die strenge Import-Zeilen-Prüfung
  (so wie sie der DoD selbst beschreibt) hat keinen Treffer. Kein
  `examples/python/` angelegt (Festlegung 2: `find examples -iname
  "*python*"` leer, unverändert zum Parent). `.a-check.yml` unberührt.
- **`sdks/python/pgchangefeed/src/pgchangefeed/` (Kommentar-Klassen)** —
  neue Kommentare in `grpc_client.py`, `options.py`, `__init__.py`,
  `grpc_gen/__init__.py` tragen Ist-Zustand + ADR-Verankerung, keine
  Vorher/Nachher-Sprache, keine Slice-Nummer als *Begründung* (Slice-IDs
  stehen nur neben ADR-Ankern als Herkunft, konsistent zur bestehenden
  Datei-Konvention, `AGENTS.md` §3.7 — kein Chronik-Befund).
- **`sdks/python/pgchangefeed/tests/`** — alle sechs gRPC-Unit-Tests sind
  netzlos (Fake-Channel, kein Socket) und an der Eingabeseite gebunden:
  RPC-Pfad (`/cdc.stream.v1.ChangeStream/StreamChanges`), Bearer-Metadata,
  filterleerer Request (`SerializeToString() == b""`), unmapped Yield
  (Identitäts-Typ), Unauthenticated-Durchreichung vom Iterator — jede
  Aussage kann durch eine Client-Verhaltensänderung rot werden. Die
  Deskriptor-Assertion bindet den Test an den **zur Bauzeit erzeugten**
  Stub (Proto-Drift bricht den Bau). Die im Implementer-Bericht genannten
  zwei gefahrenen Mutationen sind in keinem committbaren Artefakt
  nachgelesen — die Bindung hier durch eigenes Nachlesen re-hergeleitet
  (Ergebnis: gebunden, siehe F-4 für die eine Lücke).
- **`sdks/python/pgchangefeed/integration/`** — reale Server-Zustellung ist
  am Wire belegt (Tabelle/Operation/Sentinel) und über die `change_id` gegen
  den Lesezugriffsweg `cdc.changes` gehalten (Runner-Block „Unabhängiger
  SQL-Beleg"); die Negative-Prüfung öffnet den Stream **ohne** Metadata und
  prüft `UNAUTHENTICATED` — reale Server-Instanz statt ausschließlich
  Unit-Mock ([ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 2 erfüllt, soweit aus dem Diff und den Runner-Assertionen
  lesbar). Einziger Befund der Datei: F-2.
- **`sdks/python/Dockerfile`, `sdks/python/.gitignore`, `pyproject.toml`** —
  der benannte Bau-Kontext `proto` ist zwingend, der Abbruchpunkt
  (`COPY --from=proto`) ist der einzige, kein stiller Fallback; Stub-Module
  uncommittet (`.gitignore`-Muster + Negation für den Package-Marker — die
  Negation ist korrekt formbar, weil das Muster die Dateien, nicht das
  Verzeichnis ausschließt); `protobuf>=6` als Laufzeitabhängigkeit trägt
  eine Begründung am Ort (Stub-Import), `grpcio-tools` bleibt Test-Extra;
  `testpaths = ["tests"]` trennt Unit-/Integrationssammelung sauber — kein
  stiller Ausschluss (Negativprobe gegen
  `BEO-PGC/test-runner-stiller-ausschluss`: der blanke `pytest`-Lauf
  sammelt `tests/` vollständig, die `integration`-Stufe ruft ihren Pfad
  explizit, keine `--ignore`-Tricks); der `sed`-Schritt hebt den
  generierten flachen Geschwister-Import auf den Package-Import (Docker-ONLY,
  dokumentierte Anpassung generierter Stubs).
- **`tools/harness/run-sdk-python-integration-tests.sh`** — `set -euo
  pipefail`; alle Gate-relevanten Exit-Codes werden ungepiped gelesen
  (`AGENTS.md` §3.9); `|| true` nur an Stellen, deren Ausfall anschließend
  explizit geprüft wird (`captured`-Zahl, `grpc_change_id`-Extraktion); die
  Vertragswerte stimmen mit `compose.yaml` (Container-Namen `cdc-test-postgres`
  / `cdc-test-feed`, Netz `cdc-feed-test`, `CDC_TABLES`-Trio, Slot
  `slot_pgc_e2e`, Reader-Token) und mit dem `Makefile`-Vertrag
  (`SCHEMA_TARGET`/`SCHEMA_ROLLOUT_NETWORK`, Zeile 224 ff.) überein; das
  Skript bleibt eigenständig — `run-integration-tests.sh` ist unverändert
  (geprüft, nicht im Diff); die Insert-ID-Folge 300 ff. kollidiert nicht mit
  dem 260er-Bereich des Server-E2E-Runners.
- **`tools/harness/sdk-pack-python.sh`** — der `pack-export`-Aufruf trägt
  `--build-context proto=proto` zwingend (Abbruchpunkt wie im Dockerfile),
  Muster `tools/harness/sdk-pack-csharp.sh`; Extraktions-Pipe unverändert.
- **`harness/mk/sdk.mk`, `harness/README.md`** — Target-Konvention
  (`.PHONY`, `@bash tools/harness/…`, Hilfe-Zeile) konsistent; die neue
  `make test-sdk-python-integration`-Zeile steht in §Werkzeuge mit
  Bindung-Spalte (kein Gate, [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  · seit slice-sdk-python-grpc-client-flaeche) und das Target existiert real
  in `sdk.mk`; das inline-code-Zitat `ADR-0090` ohne Link in der neuen Zeile
  ist bestehende Datei-Konvention (Zeile 155 macht dasselbe) und kein
  Linkpflicht-Verstoß.
- **Plan-vs-Code (`slice-sdk-python-grpc-client-flaeche.md` §3 + Plan-Nachzug)** —
  alle neun §3-Zeilen und alle sieben Nachzug-Zeilen sind im Diff vertreten;
  **keine still gestrichenen Planpunkte**; die Abweichungen
  (`integration/`-Ordner statt `tests/integration/`, `protobuf>=6`,
  Package-Marker, `.gitignore`, Skript-Kontext, `timeout`-Parameter,
  README/Docstring-Nachzug) sind sämtlich als Plan-Nachzug deklariert.
- **§3.13 Träger-Nachzug (beide Stände gemessen: Parent `a13c1d03` und
  HEAD)** — Suchlauf über `sdks/python/**`, `docs/**`, `harness/**`,
  `spec/**` nach der bewegten Eigenschaft („gRPC/SSE/NATS bleiben außerhalb
  des Python-Packages"): nachgezogen in `sdks/python/README.md` (§Status),
  `__init__.py`-Docstring, `options.py`-Docstring; **nicht** nachgezogen:
  `docs/user/benutzerhandbuch.md` (F-3, plan-deklarierter Aufschub),
  [`SPEC-027`](../../spec/pflichtenheft.md)/`LH-FA-SST-009.a` (gebündelt, wie
  bei C#/Kotlin an derselben Stelle praktiziert — Änderungshistorie-Zeile
  2026-09-22, konsistent mit [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Folgepflicht 2 „sobald das Folge-Release real existiert"),
  `README.md`-Wurzel Zeile 27 (bleibt für das **veröffentlichte** 0.1.0-Package
  wahr, Nachzug beim Version-Bumb gebündelt — kein Befund in diesem Slice),
  offene NATS-Slice-Planzeile (F-5).
- **`AGENTS.md` §3.11 (host-lokale Pfade)** — der ganze Diff (Code, Kommentare,
  Skripte, Plan, README-Zeile) enthält keinen host-lokalen absoluten Pfad.
- **`AGENTS.md` §3.12 Instanz A (Zahlen)** — `29 Unit-Tests` (F-1: eine der
  beiden Zahlen driftet), `change_id=804-1` (gemessen, Lauf genannt:
  `make test-sdk-python-integration EXIT=0`), `EXIT=0`-Belege beider Targets
  tragen Lauf+Ursprung; keine Zahl im Diff ohne Ursprung außer F-1.
- **`AGENTS.md` §3.5 (Accepted-ADR-Immutabilität)** — keine ADR-Datei im
  Diff; [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  bleibt unangetastet; der Umsetzungs-Schnitt (eigenständiges Skript,
  geteilte Compose-Umgebung, kein Wegwerf-Client) liegt innerhalb der
  vom Slice-Plan deklarierten Architektur-Entscheidung und der
  [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  -Delegation an den umsetzenden Zug (Folgepflicht 1).

---

## Verdikt

**Fixrunde am Implementer nötig** (2 HIGH: F-1, F-2). Die DoD-Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" bleibt
deshalb **offen** — sie wird regulär bei Schritt 21 des
Implementer-Workflows nachgezogen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde greift nicht, weil eine Fixrunde
existiert). F-3/F-4 (MEDIUM) gehen in dieselbe Rückgabe; F-5/F-6 sind
ohne Implementer-Rückgabe direkt behebbar.

**Nicht geprüft (Grenze dieses Laufs):** der reale `make gates`- und
`make test-sdk-python-integration`-Lauf (DoD-Substanz) — Verifier-Aufgabe;
der reale Bau der Stub-Module aus der `.proto` (braucht Netz) — der
Sensor-Beleg im DoD bleibt dessen Behauptung, nicht meine Messung.