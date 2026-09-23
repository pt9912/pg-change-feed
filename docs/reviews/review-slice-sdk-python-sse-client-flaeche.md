# Review-Report: slice-sdk-python-sse-client-flaeche — 2026-09-23

**Review-Art:** Code — gegen Slice-Plan + ADR + Hard Rules geprüft
(DoD-Substanz bleibt Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-sdk-python-sse-client-flaeche` · Diff-Range
`b43251ee..6bbe99d9` (Implementations-Commit `6bbe99d9`; davor nur reine
Move-Commits des Slice-Plans).

**Skill:** `.harness/skills/reviewer.md` (geschärft 2026-09-09) ·
**Modell:** `glm-5.3-flash` (Claude-Agenten-SDK, Reviewer-Rolle) ·
**Datum:** 2026-09-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne diese
Liste ist der Lauf nicht reproduzierbar):

- `docs/plan/planning/in-progress/slice-sdk-python-sse-client-flaeche.md`
  (§2 DoD, §3 Plan + Plan-Nachzug, §6 Risiken, §8 Beobachtungs-Sichtung)
- [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Entscheidung Festlegung 1/2, §Konsequenzen Folgepflicht 1–3
- `spec/pflichtenheft.md` `SPEC-021` (Draht-Vertrag) · `SPEC-022`
  (Abgrenzung zwölf gegen zehn Felder) · `SPEC-027`
- `internal/adapters/driving/http/sse.go` (Server-Gegenprobe, nur gelesen)
- Vorgänger-Slice `slice-sdk-python-grpc-client-flaeche` + dessen
  Haupt-Review (F-1…F-6) und Re-Review (F-1…F-4) — Lernklassen
- `AGENTS.md` (§3 Hard Rules) · `harness/conventions.md` (MR-000) ·
  `v6.9.0` · `regelwerk/modul-08-agentenrollen.md` §Die neun Übergaben

---

## Findings

### F-1 — Handbuch: SSE-`**SDK:**`-Absatz für Python fehlt; der im Diff geschriebene DoD-Vermerk widerspricht seinem eigenen Kästchen

- `kategorie`: MEDIUM
- `quelle`: [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Konsequenzen Folgepflicht 3 („SDK-Hinweis je neu gedeckter Oberfläche,
  im selben Zug wie das jeweilige Client-Programm") · `AGENTS.md` §3.13
- `pfad`: `docs/user/benutzerhandbuch.md:870-871` (nach dem
  Kotlin-`**SDK:**`-Absatz der SSE-Sektion, vor „Zugriff über das
  NATS-Wecksignal"); Abgrenzungs-Quellen:
  `docs/plan/planning/in-progress/slice-sdk-python-sse-client-flaeche.md`
  §2 (Doku-Update-Zeile samt Vermerk) und §1 („Ausdrücklich NICHT")
- `befund`: Der SSE-Abschnitt des Handbuchs trägt bei HEAD
  `**SDK:**`-Absätze für .NET (1.38) und Kotlin/JVM (1.39), aber keinen
  Python-Absatz — die Fläche dieses Slices ist im Handbuch-Hinweis-Muster
  seiner drei Vorgänger (C# 1.38, Kotlin 1.39, Python-gRPC 1.41 in der
  Fixrunde des Vorgänger-Slices) die einzige ohne eigenen Zug. Der
  in **diesem Diff** geschriebene DoD-Vermerk behauptet „der
  SSE-`**SDK:**`-Absatz gehört in denselben Zug wie die Fläche", während
  der Kästchen-Kopf desselben Punkts „bewusst nicht in diesem Slice"
  sagt — der Träger widerspricht sich im selben Absatz, und der
  benannte Zustand fehlt. Der plan-seitige Aufschub auf den letzten
  Flächen-Slice (§1) deckt nur den NATS-Handbuch-Teil und
  `spec/pflichtenheft.md`; für SSE widersprechen sich Aufschub und Vermerk.
- `verifizierbar`: ja — `grep -n "PgChangeFeedSseClient"
  docs/user/benutzerhandbuch.md` (kein Treffer im SSE-Abschnitt;
  `pgchangefeed`-Treffer nur in HTTP-/gRPC-Absätzen); Änderungshistorie
  endet bei 1.41
- `klasse`: „Träger-Nachzug gebündelt statt im Zug" (Familie
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`;
  drittes Auftreten der Handbuch-Muster-Abweichung — Kotlin 1.39,
  gRPC-Python Haupt-Review F-3 MEDIUM, gRPC-Python Re-Review F-2 LOW)

### F-2 — harness/README.md: die `make test-sdk-python-integration`-Zeile trägt die Ein-Flächen-Form weiter

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug, beide Stände gemessen)
- `pfad`: `harness/README.md:160` gegen
  `tools/harness/run-sdk-python-integration-tests.sh` (`run_surface_phase`,
  zwei Phasen je Fläche)
- `befund`: Die Werkzeug-Zeile beschreibt den Lauf weiter in der
  Ein-Flächen-Form — „der Prüfling ist das SDK selbst
  (`pgchangefeed.grpc_client`, kein Wegwerf-Duplikat-Client)" und „ein
  Stream-Öffnungsversuch ohne Token endet mit gRPC-Status
  `Unauthenticated`" —, während der Runner seit diesem Diff je Fläche
  eine Phase mit eigener Testdatei, eigenem Sentinel/ID-Bereich
  (gRPC 300ff., SSE 310ff.) und eigener Reject-Marker-Form
  (`Unauthenticated` bzw. `401`) fährt; die Aufschub-Klausel der Zeile
  („erweitert sich Slice für Slice um die weiteren Flächen dieser Welle
  (SSE, NATS-Vollinhalt)") ist mit diesem Slice verbraucht — die
  Erweiterung ist hier gelandet, die Zeile nennt sie nicht. Der
  Plan-Nachzug (§3) listet `__init__.py`/`options.py`/
  `sdks/python/README.md`, aber nicht diese Zeile; ein committetes
  Suchlauf-Feld für Gefundenes/Nichtgefundenes trägt der Plan noch nicht
  (§3.13 verlangt es vor der Closure-Notiz).
- `verifizierbar`: nein — Lesen (kein Gate; `docs-check` prüft Referenzen,
  nicht Vollständigkeit)
- `klasse`: „Arbeit überholt stehenden Träger" (4. Auftreten der
  Beobachtungsklasse; die Vorgänger-Zeile selbst wurde im Vorgänger-Slice
  mit derselben Begründung im Zug nachgezogen)

### F-3 — ENTRYPOINT-Default-Flip: blanke Aufrufform der `integration`-Stufe läuft grün über die 41 Unit-Tests statt über die Integrationstests

- `kategorie`: LOW
- `quelle`: Maintainability; verwandt mit der offenen Klasse
  `BEO-PGC/test-runner-stiller-ausschluss` (hier als Image-Default statt
  als Runner-Aufruf)
- `pfad`: `sdks/python/Dockerfile:138` gegen
  `sdks/python/pgchangefeed/pyproject.toml:52` (`testpaths = ["tests"]`)
- `befund`: Die alte Fassung startete ohne Argumente `pytest integration`;
  die neue (`ENTRYPOINT ["python", "-u", "-m", "pytest"]`) sammelt ohne
  Argumente `testpaths = ["tests"]` — real gemessen: `docker run --rm
  --network none pg-change-feed:sdk-python-integration` endet mit
  „41 passed", Ausgang 0, ohne dass ein Integrationstest je lief. Der
  Runner nennt je Phase seine Testdatei explizit (der eigentliche
  stiller-Ausschluss-Mechanismus ist sauber), aber die blanke
  Image-Aufrufform — derselbe Aufruf, der vorher die Realserver-Tests
  fuhr — liefert jetzt ein grünes Ergebnis ohne einen einzigen der
  belegten Rundläufe.
- `verifizierbar`: ja — der blanke `docker run`-Aufruf des
  Integration-Images (oben ausgeführt; Ausgang 0, 41 Unit-Tests)
- `klasse`: „Stiller Ausschluss der Testdatei im Default-Aufruf"

### F-4 — `PgChangeFeedMalformedResponseError`-Docstring nennt weiter nur SPEC-018/SPEC-022 als Formen, die ihn tragen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Deklarations-Nachzug — dieselbe Träger-Klasse
  wie Verifikations-Auflage V-3 des Vorgänger-Slices)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py:51-57`
  gegen `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py:115-133`
- `befund`: Der Docstring sagt, die Klasse trage Verletzungen der
  „[`SPEC-018`](../../spec/pflichtenheft.md)/[`SPEC-022`](../../spec/pflichtenheft.md) response shape"; seit diesem Slice wirft auch der
  SSE-Client ihn (Status 200, drei Befund-Formen gegen `SPEC-021`). Kein
  falscher Satz über eine Stelle, die er jetzt trägt — aber die
  Klassen-Doku nennt zwei von inzwischen drei Straten nicht.
- `verifizierbar`: nein — Lesen
- `klasse`: „Deklarations-Nachzug fehlt"

### F-5 — `sse_client.py` endet ohne Zeilenumbuch am Dateiende

- `kategorie`: LOW
- `quelle`: Maintainability (Stil)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py:133`
- `befund`: Die neue Datei trägt `\ No newline at end of file`; der
  Altbestand (Dockerfile, Runner-Skript) trägt dieselbe Form, im neuen
  Python-Modul ist sie eine neue Datei-Eigenschaft ohne semantische
  Auswirkung.
- `verifizierbar`: nein — `tail -c 1` auf die Datei
- `klasse`: „Fehlender Zeilenumbuch am Dateiende"

### F-6 — Der `\r`-Strip im Frame-Parser ist mit der realen httpx-Form erreichbar-leer

- `kategorie`: INFO
- `quelle`: Maintainability (Grenze-Merkung, kein Handlungsbedarf)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py:94`
- `befund`: `httpx` 0.28.1 (real im Integration-Image gelesen) implementiert
  `iter_lines` über `LineDecoder` mit `splitlines`-Semantik — CRLF ist
  bereits dekodiert, bevor eine Zeile beim Parser ankommt; der
  `\r`-Strip ist für einen anderen (naiven `\n`-Splitter) Zeilen-Iterator
  gebaut. Real gemessen: CRLF, CR/NL-Split über Chunk-Grenzen,
  Mid-Line-Split und LF liefern je das Event — der Strip bleibt harmlose
  Verteidigung, kein Kommentar behauptet das Gegenteil.
- `verifizierbar`: ja — die Container-Probe (MockTransport + vier
  Chunk-Formen) im Report-Verlauf ausgeführt
- `klasse`: — (Hinweis ohne erwartete Aktion)

### F-7 — `_build_error`-Wiederverwendung: Kopplung verhaltensseitig deklariert, Import-Site trägt sie nicht

- `kategorie`: INFO
- `quelle`: Maintainability (Kopplungs-Deklaration, `AGENTS.md` §3.7)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py:51`
- `befund`: `sse_client` importiert das private `_build_error` aus
  `http_client` (neu in diesem Slice — die gRPC-Fläche nutzt es nicht).
  Die Verhaltenskopplung ist im Modul-Docstring deklariert („the same
  typed status-code error as the HTTP surface"); ein Umbenennen schlägt
  beim Modul-Import laut fehl, nicht erst an einer Stelle. Kein
  Handlungsbedarf — Merkung für die Klassen-Dokumentation des Packages.
- `verifizierbar`: nein — Lesen
- `klasse`: — (Hinweis ohne erwartete Aktion)

### F-8 — SSE-Realserver-Test prüft die Frist erst nach Empfang eines Events

- `kategorie`: INFO
- `quelle`: Maintainability (Grenze-Merkung)
- `pfad`: `sdks/python/pgchangefeed/integration/test_sse_realserver.py:49-57`
- `befund`: Der gRPC-Test prüft die 90-s-Frist vor jedem `next()`, der
  SSE-Test erst nach Empfang eines Events — ein stummer Stream endet am
  `httpx`-Read-Timeout (95 s, höher als die Frist) mit einem untypisierten
  httpx-Fehler statt der Frist-Assertion. Der Ausgang bleibt rot sichtbar
  (der Runner pollt selbst und bricht mit eigener Meldung ab), die
  Frist-Assertion trägt nur den Fall „Events kommen, aber keiner passt".
- `verifizierbar`: nein — Lesen
- `klasse`: — (Hinweis ohne erwartete Aktion)

### F-9 — Server-seitige `sse.go`-Kommentare zitieren SPEC-018 für den SSE-Wire, kanonisch ist SPEC-021

- `kategorie`: INFO
- `quelle`: Maintainability (Träger außerhalb des Diffs — nur die Messung
  liest ihn)
- `pfad`: `internal/adapters/driving/http/sse.go:12-15` (und zwei weitere
  Stellen desselben Musters)
- `befund`: Die Server-Kommentare nennen für Endpunkt, Event-Typ und
  Nachrichtenschema „`SPEC-018`, Abschnitt `GET /changes/stream`" — das
  Pflichtenheft trägt die SSE-Festlegung als `SPEC-021` (eigener Eintrag
  „statt einer Erweiterung von `SPEC-018`"). Der Träger liegt außerhalb
  dieses Diffs; hier nur zur Meldung, kein Nachzug-Pfeil aus diesem Lauf.
- `verifizierbar`: nein — Lesen
- `klasse`: „Zitat nennt die falsche Stelle" (Familie
  `BEO-PGC/zitat-nennt-die-falsche-stelle`, außerhalb des Diffs — INFO statt
  HIGH)

## Negativbefunde

- geprüft, ohne Befund: `spec/pflichtenheft.md` `SPEC-021` (Endpunkt,
  Event-Form, zehn Felder) gegen `sse_client.py`/`models.py` — Frame-Form,
  Auth-Rechtsklasse (`reader`/`admin`, Bearer-Header), Feldsatz und
  `null`-Bild-Boundary decken den Draht; `SPEC-027`/`LH-FA-SST-009.a`-Deferral
  konsistent mit [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Folgepflicht 2 (der Träger beschreibt das
  veröffentlichte 0.1.0-Release, der Version-Bump ist gebündelt — wie im
  Vorgänger-Review praktiziert)
- geprüft, ohne Befund: `internal/adapters/driving/http/sse.go`
  (Gegenprobe, nur gelesen) — Frame-Form (`event: change\ndata: <json>\n\n`),
  503 ohne Broadcaster, 401 vor jedem Event; der Zitat-Befund ist F-9/INFO
- geprüft, ohne Befund: `sdks/python/pgchangefeed/tests/test_sse_client.py`
  — 10 Tests, Zusage-Bindung an der Eingabeseite geprüft (Header-Form am
  Handler-Assert des Fakes; 401/503 an der Status-Prüfung im Stream-Kontext;
  Chunk-Boundary real gebunden — Probe im Integration-Image mit vier
  Chunk-Formen gefahren, alle grün); 41 Unit-Tests real grün im
  Integration-Image (netzlos, Ausgang 0)
- geprüft, ohne Befund: Import-Grenze (`ADR-0110` Festlegung 4) — `grep`
  über `sdks/python/` nach Importen aus `internal/**`/`cmd/**`/`gen/**`:
  kein Treffer; `grpc_gen` ist eigenes Unterpackage des Packages
- geprüft, ohne Befund:
  `tools/harness/run-sdk-python-integration-tests.sh` §3.9-Disziplin —
  `set -euo pipefail` am Kopf, `exit 1` je Phase propagiert über die
  `$(...)`-Aufrufe der Phasen (Subshell-Ausgang wird mit `set -e` wirksam),
  `|| true`-Bewachung der `change_id`-Extraktion mit nachgelagerter
  Leer-Prüfung, ungepipedes `bash`-Target in `harness/mk/sdk.mk`; die
  RECEIVED-Marker-Form passt zu beiden Testdateien (beide drucken
  `change_id/table/operation/new_image` in derselben Feldreihenfolge)
- geprüft, ohne Befund: Plan-vs-Code (§3-Tabelle + Plan-Nachzug) — alle
  neun Planzeilen sind im Diff vertreten, keine still gestrichen; die
  drei Abweichungen (Pfad `integration/` statt `tests/integration/`,
  flexibler ENTRYPOINT, `StreamChange` + Runner-Refactor) sind als
  Plan-Nachzug deklariert
- geprüft, ohne Befund: `AGENTS.md` §3.12 Zahlen — `grep -c "^def test_"`
  über `tests/*.py`: 20/3/8/10, Summe 41, im Plan exakt so behauptet und
  als „aus derselben Messung gezogen" deklariert; reale Sammlung im Image:
  „41 tests collected"; die `change_id=804-1`/`807-1`-Werte tragen ihren
  Lauf (`make test-sdk-python-integration` EXIT=0) und sind als
  Lauf-Messung nicht netzlos reproduzierbar — Mechanik (Marker,
  SQL-Gegenprüfung, Marker-Asserts) im Diff geprüft
- geprüft, ohne Befund: Träger-Nachzug an den drei Plan-Nachzug-Stellen,
  beide Stände gemessen (Parent `b43251ee` trug die „SSE bleibt
  außerhalb"-Form in `sdks/python/README.md:9`, `__init__.py:5-7` und
  `options.py`-Docstring; HEAD in keinem der drei) — `spec/pflichtenheft.md`
  `SPEC-027`-Text bleibt für das veröffentlichte 0.1.0-Release wahr
  (Deferral deklariert, Wurzel-README-Form wie im Vorgänger-Review
  praktiziert)
- geprüft, ohne Befund: Kommentar-Klassen (`AGENTS.md` §3.7) über Code,
  Dockerfile und Runner — Zusage/Kopplung/Abgrenzung, keine Chronik, kein
  abwesender Text, keine `//nolint`-Form (§3.2); Deutsch/Englisch-Mix in
  Docstrings ist die offene Beobachtung
  `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (Architect-Entscheidung,
  kein Handlungsbedarf) — kein neuer Befund
- geprüft, ohne Befund: `AGENTS.md` §3.11 (host-lokale Pfade) — der ganze
  Diff nennt keinen host-lokalen Absolutpfad; §3.5 (Accepted-ADR-Immutabilität)
  — [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) unberührt
- geprüft, ohne Befund: Commit-Message-Traceability (`6bbe99d9` nennt
  `LH-FA-SST-009` + `ADR-0110`, kein `SPEC-*`/`ARC-*` im Betreff)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** „Träger-Nachzug gebündelt statt im Zug"
(F-1, 3. Auftreten der Handbuch-Muster-Abweichung) · „Arbeit überholt
stehenden Träger" (F-2, 4. Auftreten) · „Stiller Ausschluss der
Testdatei im Default-Aufruf" (F-3, verwandt mit
`BEO-PGC/test-runner-stiller-ausschluss`) · „Deklarations-Nachzug fehlt"
(F-4) · „Fehlender Zeilenumbuch am Dateiende" (F-5)

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 brauchen eine Fixrunde am
Implementer (beide sind Inhalt des Slice-eigenen Nachzugs: der
Handbuch-`**SDK:**`-Absatz je nach dem eigenen DoD-Vermerk und der
C#/Kotlin/gRPC-Präzedenz in diesem Slice, die `harness/README.md`-Zeile im
selben Zug). Die LOW/INFO-Findings sind ohne Pfeil zur Fixrunde behebbar
oder brauchen keine Aktion (F-6–F-9).

**Übergabe:** Findings F-1/F-2 (MEDIUM) gehen als Rückkante an den
Implementer; F-3–F-5 (LOW) in dieselbe Rückgabe, ohne eigenen
Blockier-Pfeil; F-6–F-9 (INFO) sind Merkungen ohne Aktion. Die
DoD-Zeile „Review durchgeführt" bleibt offen — der Nachzug läuft regulär
bei der Fixrunde (Schritt 21 des Implementer-Workflows), der
Selbst-Nachzug ohne Fixrunde greift hier nicht. Vor der Closure-Notiz:
der §3.13-Suchlauf-Befund (Gefundenes/Nichtgefundenes) in einem
committeten Feld des Slice-Plans festhalten — F-2 ist der Befund, den
der Suchlauf hätte melden müssen. Dieser Report ist Lauf-Beleg; die
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).