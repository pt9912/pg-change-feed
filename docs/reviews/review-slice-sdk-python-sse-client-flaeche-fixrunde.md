# Review-Report: slice-sdk-python-sse-client-flaeche — Fixrunde — 2026-09-23

**Review-Art:** Code (Re-Review der Fixrunde) — gegen Haupt-Review-Report +
Hard Rules geprüft (DoD-Substanz bleibt Verifier-Aufgabe, Modul 11).

**Gegenstand:** Fixrunde zu
[`review-slice-sdk-python-sse-client-flaeche.md`](review-slice-sdk-python-sse-client-flaeche.md)
(F-1…F-5) · Diff-Range `6bbe99d9..HEAD` — Fixrunden-Commit `3c941b0b` +
Link-Nachzug `228dfc9e`.

**Skill:** `.harness/skills/reviewer.md` (geschärft 2026-09-09) ·
**Modell:** `glm-5.3-flash` (Claude-Agenten-SDK, Reviewer-Rolle) ·
**Datum:** 2026-09-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne diese
Liste ist der Lauf nicht reproduzierbar):

- Haupt-Review
  [`review-slice-sdk-python-sse-client-flaeche.md`](review-slice-sdk-python-sse-client-flaeche.md)
  (F-1…F-9) und dessen Übergabe-Satz (Suchlauf-Feld vor der Closure-Notiz)
- `docs/plan/planning/in-progress/slice-sdk-python-sse-client-flaeche.md`
  (§2 DoD, §3 Plan + Plan-Nachzug + §3.13-Suchlauf-Feld)
- [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Entscheidung Festlegung 2, §Konsequenzen Folgepflicht 1–3
- `spec/pflichtenheft.md` `SPEC-021` · `spec/lastenheft.md`
  [`LH-FA-SST-009`](../../spec/lastenheft.md)
- `AGENTS.md` (§3 Hard Rules, §3.9/§3.12/§3.13) · `harness/conventions.md`
  (MR-000) · `v6.9.0` · `regelwerk/modul-08-agentenrollen.md` §Die neun
  Übergaben

---

## Findings

### FR-1 — Die Mechanik-Umstellung der Fixrunde (Argument → Umgebungsvariable) hat ihre Ist-Zustands-Träger nicht gezogen

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug, beide Stände gemessen)
- `pfad`: `tools/harness/run-sdk-python-integration-tests.sh:22` und
  `:29` · `harness/README.md:160`; mildernd:
  `docs/plan/planning/in-progress/slice-sdk-python-sse-client-flaeche.md:124`
  und `:126` (Zug-Chronik des Implementation-Zugs, vom „Real jetzt" der
  Fixrunden-Zeile `:128` überholt — kein eigenständiger Befund)
- `befund`: Die Fixrunde stellte die Testdatei-Übergabe von
  docker run-Argument auf Umgebungsvariable um (Runner `:178`
  `-e PGCHANGEFEED_TEST_FILE=…`, Dockerfile-`CMD` liest sie) und trug die
  neue Form an drei Stellen korrekt (Dockerfile-Kommentar, Plan `:128`,
  Runner-Code) — die Ist-Zustands-Träger derselben Eigenschaft blieben auf
  der alten Form: der Runner-Skriptkopf nennt weiter „Stufe-ENTRYPOINT"
  (die Stufe trägt einen `CMD`) und „benennt ihre Testdatei explizit als
  docker run-Argument" (beide Stände gemessen: bei `6bbe99d9` wahr, seit
  `3c941b0b` falsch), und die neue harness/README-Zeile `:160` schreibt
  dieselbe Argument-Form. Der committete §3.13-Suchlauf deckt nur die
  „SSE bleibt außerhalb"-Eigenschaft, nicht die Mechanik-Umstellung —
  dieselbe Lücken-Struktur, die Haupt-Review F-2 benannte. Kein stiller
  Ausschluss-Rückweg: Fehlgebrauch nach der alten Form scheitert laut
  (real gemessen, siehe Negativbefunde).
- `verifizierbar`: ja — `grep -n "docker run-Argument"
  tools/harness/run-sdk-python-integration-tests.sh harness/README.md`
  gegen `tools/harness/run-sdk-python-integration-tests.sh:178` (beide
  Stände in diesem Lauf gemessen)
- `klasse`: „Arbeit überholt stehenden Träger" (5. Auftreten; die
  Fixrunde wiederholt die Lücken-Struktur des Haupt-Review-F-2: eine
  bewegte Eigenschaft gezogen, die zweite lief ohne Suchlauf durch)

### FR-2 — Link-Nachzug `228dfc9e` an einem committeten Review-Report — gemessen, kein Befund-Inhalt geändert

- `kategorie`: INFO
- `quelle`: Maintainability (Merkung zu Übergabe-Artefakten, Modul 8)
- `pfad`: `docs/reviews/review-slice-sdk-python-sse-client-flaeche.md`
  (drei Stellen) ·
  `docs/plan/planning/in-progress/slice-sdk-python-sse-client-flaeche.md:140`
- `befund`: Der Commit behebt `id-unlinked`-Befunde durch Link-Wicklung
  in bereits committeten Trägern. Real gemessen: der Diff ändert nur
  Link-Formen (Befund-Texte, Kategorien und Pfade unverändert) und
  deklariert den Fund in seiner Message; `make docs-check` endet nach dem
  Zug real grün (912 Dateien, 0 Befunde, Ausgang 0). Der Report-Inhalt
  (Findings, Stufen, Verdikt) ist unberührt — Merkung, kein Befund.
- `verifizierbar`: ja — `git show 228dfc9e` (6 geänderte Zeilen, nur
  Link-Formen); `make docs-check`
- `klasse`: — (Hinweis ohne erwartete Aktion)

## Negativbefunde

- geprüft, ohne Befund: **F-1 gelöst** — der Python-`**SDK:**`-Absatz steht
  im SSE-Handbuch-Abschnitt nach .NET und Kotlin (dritte Sprache, real
  gezählt), Version-Kopf auf `1.42` gezogen, Historie-Zeile `1.42` im
  selben Diff; alle Assertions des Absatzes gegen den Code gemessen:
  `PgChangeFeedSseClient.stream_changes()` existiert, liefert
  `Iterator[StreamChange]`, Bearer-Token im `Authorization`-Header,
  401 → `PgChangeFeedUnauthorizedError` über `_STATUS_TO_ERROR`
  (`http_client.py:70`), „allen zehn Feldern" deckt sich mit
  `models.py:268-278` (10 Felder) und der gRPC-Feldtabelle (per
  SSE-Verweis)
- geprüft, ohne Befund: **F-2 gelöst** — die harness/README-Zeile trägt
  die Zwei-Flächen-Form (`grpc_client`, `sse_client`, beide Reject-Marker,
  Aufschub-Klausel korrekt auf NATS verkürzt); das committete
  §3.13-Suchlauf-Feld steht im Plan; die Mechanik-Rest-Form ist FR-1
- geprüft, ohne Befund: **F-3 gelöst, real gemessen** — das Image
  `pg-change-feed:sdk-python-integration` neu aus dem HEAD-Dockerfile
  gebaut (Digest `87a9c18b…` deckt sich mit dem Vorher-Bau — der
  Fixrunden-Commit baute vor dem Commit); blanker Lauf
  `docker run --rm --network none …` endet laut mit
  `PGCHANGEFEED_TEST_FILE: Testdatei je Flaeche erforderlich`, Ausgang 2
  (nicht still 41 Unit-Tests); die ENV-Form fährt die Phase real:
  `-e PGCHANGEFEED_TEST_FILE=tests/test_sse_client.py` → „10 passed",
  Ausgang 0, netzlos; `-e PGCHANGEFEED_TEST_FILE=integration/test_sse_realserver.py`
  → pytest lädt real die Realserver-Datei (collected 0/1 error ohne
  Stack) — je Phase fährt der `CMD` genau die benannte Datei. Der Guard
  feuert bei unset **und** leer (`:?`); der Runner setzt nie leer;
  Container-Ausgang = pytest-Ausgang (`sh -c`, letztes Kommando) — der
  ExitCode-Poll des Runners bleibt wirksam
- geprüft, ohne Befund: **F-4 gelöst** — der Docstring nennt die
  `SPEC-021`-SSE-Form; sein Kopf-Satz („does not match the documented
  stream schema") deckt alle drei Wurf-Formen des SSE-Clients (invalid
  JSON, nicht-Objekt, fehlendes Feld, `sse_client.py:113/120/129`); die
  Aufzählung nennt zwei von drei — Grenz-Merkung ohne Aktion
- geprüft, ohne Befund: **F-5 gelöst** — `sse_client.py` endet real mit
  Zeilenumbuch (`\n` gemessen); das Dockerfile behält die Bestandsform
  (Altbestand, konsistent mit der F-5-Begründung des Haupt-Reviews)
- geprüft, ohne Befund: **§3.13-Suchlauf-Tabelle gegen beide Stände
  gemessen** — Parent `b43251ee` trug die „SSE bleibt
  außerhalb"-Satzform in `__init__.py:8`, `options.py:6` und
  `sdks/python/README.md:9`; HEAD trägt sie in keinem der drei (NATS-Rest
  bleibt korrekt als Folge-Slice-Zusage); die „bleibt
  gebündelt"-Zeilen (`spec/pflichtenheft.md`, Wurzel-README,
  `docs/user/version.md`) sind im Fix-Diff real unberührt (diff-stat)
- geprüft, ohne Befund: **§3.12 Instanz A über die NEUEN Zahlen des
  Fix-Diffs** — `1.42` (Kopf, Historie, Plan-Vermerk) = committeter
  Versionsstand, gemessen; „allen zehn Feldern" = `models.py` (10 Felder)
  und der gRPC-Feldtabelle; „dritte Sprache" = drei `**SDK:**`-Absätze im
  SSE-Abschnitt, gezählt; Testzahlen re-gemessen: 20/3/8/10 = 41, davon
  10 in `tests/test_sse_client.py` (real „10 passed" im Image) — die
  Zahlen des committeten Haupt-Reports wurden nicht übernommen, sondern
  nachgemessen
- geprüft, ohne Befund: `AGENTS.md` §3.7 (Kommentar-Klassen) über die neu
  geschriebenen Kommentare des Fix-Diffs — der Dockerfile-Integration-
  Kommentar trägt Zusage/Abgrenzung im Indikativ („scheitert laut … statt
  still die Unit-Tests zu fahren"), kein Chronik-Satz, kein abwesender
  Text, keine `//nolint`-Form (§3.2); der Runner-Skriptkopf-Drift ist
  FR-1/§3.13, keine §3.7-Klassen-Verletzung
- geprüft, ohne Befund: `AGENTS.md` §3.11 (host-lokale Pfade) — der
  Sensor real grün (siehe unten); §3.5 (Accepted-ADR-Immutabilität) —
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  unberührt; Handbuch-Versionshistorie-Muster (Skill-HIGH) — beide Muster
  (Version-Kopf, Historie-Zeile) im selben Diff
- geprüft, ohne Befund: Commit-Message-Traceability — `3c941b0b` und
  `228dfc9e` nennen je [`LH-FA-SST-009`](../../spec/lastenheft.md) +
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md),
  kein `SPEC-*`/`ARC-*` im Betreff
- geprüft, ohne Befund: `make docs-check` real gefahren — Ausgang 0,
  912 Dateien, 0 Befunde (Exit-Code ungepiped direkt geprüft, `AGENTS.md`
  §3.9)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Arbeit überholt stehenden Träger"
(FR-1, 5. Auftreten)

## Verdikt

**Merge-blockierend:** ja — FR-1 (MEDIUM) braucht eine kurze zweite
Fixrunde am Implementer (Umformulierung der drei Ist-Zustands-Stellen auf
die ENV-Form; der Fix ist reiner Text-Nachzug, der Code ist gemessen
korrekt). Kein offenes HIGH.

**Übergabe:** FR-1 (MEDIUM) geht als Rückkante an den Implementer; im
denselben Zug wird die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" nachgezogen (Schritt 21, Fixrunden-Checkbox-
Nachzug — in `3c941b0b` wurde sie noch nicht gesetzt, zu Recht: das
„kein offenes HIGH" war zum Fix-Zeitpunkt noch nicht durch dieses
Re-Review bestätigt). FR-2 (INFO) ist eine Merkung ohne Aktion. Vor der
nächsten Closure-Notiz: der §3.13-Suchlauf um die zweite bewegte
Eigenschaft (Mechanik der Testdatei-Übergabe, beide Stände) erweitern —
FR-1 ist der Befund, den dieser Suchlauf gemeldet hätte. Dieser Report
ist Lauf-Beleg; die DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).
**Addendum (nach Fixrunde 2, Commit `beeddc3c`):** FR-1 gelöst — die drei
Phrase-Stellen tragen die ENV-Form (Runner-Kopf: „Stufe-CMD" und
„Umgebungsvariable `PGCHANGEFEED_TEST_FILE`"; `harness/README.md` Zeile 160:
ENV-Form samt Guard-Beschreibung); der §3.13-Suchlauf um die zweite bewegte
Eigenschaft erweitert, mit dem gemessenen Vorher-Stand `3c941b0b`.
Endstand: kein offenes HIGH/MEDIUM — die Substanz bestätigt die Verifikation
(`verifikation-slice-sdk-python-sse-client-flaeche.md`) durch eigenes
Nachmessen.
