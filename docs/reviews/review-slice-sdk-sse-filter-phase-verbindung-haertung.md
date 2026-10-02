# Review-Report: slice-sdk-sse-filter-phase-verbindung-haertung — 2026-10-02

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `sdk-sse-filter-phase-verbindung-haertung`, Diff `git diff e900e5c3 HEAD` (Commits `26e6a500`, `d0ed33d0`):
12 Dateien, +291/−90 (Fixture `tools/harness/lib-sdk-filter-fixture.sh`, drei Runner, drei Filter-Testdateien, zwei `PhaseEnvironment`-Hilfen,
`harness/README.md`, `docs/user/sdk-e2e-abdeckung.md`, Plan).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1`.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Ablage und Arbeitsweise:** Der Reviewer-Lauf hat diesen Report mit dem Write-Werkzeug geschrieben (kein Edit-Werkzeug im Lauf). Mutationen liefen in
Kopien im Scratchpad (`git archive HEAD` nach `rvw2`, davon eine zweite Kopie `rvw2B`; Mutation per `sed … > Datei im Scratchpad` und `cp` in die Kopie,
kein `sed -i`); das Tier-Image der Kopien trug den eigenen Tag `pg-change-feed-mutation:rvw-pya` bzw. `rvw-pyb` über `SDK_PYTHON_INTEGRATION_IMAGE`,
beide danach mit `docker rmi` (ohne `-f`) entfernt, `docker images | grep -c mutation` gibt 0. `git status --short` im echten Repo war nach allen Läufen leer.
`make image` wurde nicht aufgerufen (kein Server-Code im Diff; `:dev` lag als `ghcr.io/pt9912/pg-change-feed:dev` geladen vor); keine Verweigerung der
Berechtigungsschicht im Lauf (`AGENTS.md` §3.15 nicht ausgelöst). Gewählte Sprache der Package-Mutation: **Python** (das Plan-Feld verlangt „eine Sprache“;
C# und Kotlin trägt der Implementer, übernommen).

**Eingangs-Kontext:**

- Slice-Plan `sdk-sse-filter-phase-verbindung-haertung` (§1, §2 DoD, §3 Plan und Suchlauf-Feld, §6 Risiken)
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Festlegung 2)
- [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.7, §3.9, §3.12, §3.13, §3.15), [`harness/conventions.md`](../../harness/conventions.md)
- Vorheriger Report am selben Modul: `review-slice-sdk-sse-client-schema-table-filter-realserver`

**Eigene Messungen** (Exit-Codes je als eigener Schritt ausgewertet):

- `make fmt-check` Exit 0 („323 Go-Dateien geprüft, alle formatiert“; deckt nur Go).
- `make kommentar-kennungen DIFF=e900e5c3` Exit 0, kein Kandidat.
- `make suchlauf-nachmessen PLAN=<Plan>` Exit 0, „22 Zeilen stimmen“ (Stände `91e46048` und `diff`).
- `make docs-check` Exit 0 vor Anlage dieses Reports; `make sdk-public-doc-check` Exit 0; `make test` Exit 0.
- `git diff --stat e900e5c3 HEAD -- sdks`: fünf Dateien, ausschließlich `PhaseEnvironment.cs`/`.kt` und die drei Filter-Testdateien (Integrationstest-Pfade); kein
  Produktivcode, keine Version. Keine interne Kennung in den neuen Zeilen unter `sdks/`.
- Tier-Läufe unmutiert (Kotlin und Python je ein Lauf Exit 0; C# erster Lauf Exit 2 an der NATS-Phase, siehe F-1, zweiter Lauf Exit 0). Gedruckte Zeilen:
  - Python: `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`
  - Kotlin: `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`
  - C# (zweiter Lauf): `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`, SEEN nach 331 ms/326 ms, Versuch 1.
- `docs/user/sdk-e2e-abdeckung.md`: genau drei geänderte Zeilen (`git diff e900e5c3 HEAD --stat`: 6 Zeilen = 3 +/3 −), je die Filter-Zeile eines Abschnitts.
- Mutationen (Python, eigener Image-Tag, je einzeln):
  - Package: `sse_client.py` sendet `table` nicht (`("table", None)`), `RUN pytest` der Unit-Stufe an der Kopie zu `RUN true`: **rot**, Exit 2, gedruckt
    `FILTER_RESULT f1=4 f1_foreign=2 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15` (erwartet `f1_foreign=2`, gesehen).
  - Zeilenzahl der zweiten Gruppe falsch erwartet (Fixture Zeile 257: F1 `!= "2"` statt `!= "1"`): **rot**, Exit 2, „F1 empfing 1 Change(s) der zweiten Gruppe“.
  - Eine Mutation, die nur die `SEEN_SECOND`-Bedingung im Test entfernt, habe ich nicht gefahren; siehe F-3 (hergeleitet).

---

## Findings

### F-1 — Runner-Ablehnungs-Schleife: `grep -q` hinter `docker logs` unter `pipefail` verfehlt eine vorhandene Zeile (gilt für alle drei Runner, nicht im Diff)

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Flake im Tier-Gate); `AGENTS.md` §3.9 (Exit-Code einer Pipe) sinngemäß
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh:246` (ebenso Kotlin `:244`, Python `:280`); `set -euo pipefail` in `:54`
- `befund`: Im ersten C#-Tier-Lauf endete die NATS-Phase mit „der Ablehnungs-Beleg blieb aus (REJECTED token-rejected fehlt)“, obwohl die Zeile
  `REJECTED token-rejected: …` im ausgegebenen Container-Log stand (gemessen am Fehlertext des Laufs); der zweite Lauf war grün. Ursache **hergeleitet** und
  teilweise gemessen: der Test ist nach der `RECEIVED`-Zeile in etwa 25 ms fertig, die Schleife prüft `docker logs … | grep -qF` und danach erst den Container-Zustand;
  `grep -q` beendet sich beim ersten Treffer, `docker logs` schreibt die restlichen Zeilen des Stack-Traces (20 Zeilen) in die geschlossene Pipe, und unter
  `pipefail` wird die Pipeline dadurch falsch (SIGPIPE). Der anschließende `docker inspect` meldet „nicht laufend“ und beendet die Schleife ohne zweiten Versuch und ohne
  letzte Prüfung. Gemessen: ein beendeter Container mit einer Treffer-Zeile und 20 Folgezeilen, `docker logs | grep -qF` unter `pipefail`, 300 Aufrufe → 3 Fehlschläge (1 %).
  Der Implementer meldete dasselbe Symptom als „Race `run_phase` prüft zuerst `docker logs`, dann den Zustand“; der Zustandsvergleich allein erklärt den
  Fehlschlag nicht, solange `grep` die Zeile findet — der Treffer-Ausfall muss hinzukommen. Weitere Grep-Pipelines derselben Form (`RECEIVED`, `READY`) stehen in allen drei Runnern.
- `verifizierbar`: ja — `make test-sdk-csharp-integration` wiederholt fahren; die 300er-Schleife ist reproduzierbar.
- `klasse`: „Pipe-Grep unter pipefail verwirft einen Treffer“

### F-2 — Härtungsbeweis Arm A: das Falsch-Grün ist nur über die Zählung belegt, nicht über gedruckte `RECEIVED_F1`-Zeilen

- `kategorie`: LOW
- `quelle`: Maintainability; Plan §2 Liefer-Punkt 3 („die Probe gilt nur, wenn die `RECEIVED_F1`-Zeilen von Arm A zeigen …“)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-sse-filter-phase-verbindung-haertung.md` (Beleg-Tabelle, Zeile „Arm A“)
- `befund`: Der Plan verlangt als Bedingung der Probe die `RECEIVED_F1`-Liste; im grünen Arm-A-Lauf druckt die Fixture sie nicht (nur `FILTER_RESULT`), der Beleg
  ist `f1=1` plus die bestandene Gegenlesung (die eine Zeile ist `public.feed_e2e_sdkfilt_a`). Der Schluss trägt logisch: unter der Package-Mutation sieht ein
  verbundener F1 beide `public`-Tabellen der ersten Gruppe, `f1=1` mit Tabelle A heißt, dass B fehlt (**hergeleitet**, nicht aus einer Zeile abgelesen); der
  Arm-B-Lauf (`f1=3 f1_foreign=1`) ist dazu konsistent (A1, B2, A2). Die Indirektheit trägt, eine Fixture-Option oder ein Lauf, der die Zeilen immer druckt, würde sie zu einer
  Ablesung machen; die Bedingung des Plans ist damit wörtlich nicht erfüllt.
- `verifizierbar`: nein — Arm A ist ein einmaliger Lauf an einer Kopie des Parent; ich habe ihn nicht nachgefahren (übernommen).
- `klasse`: „Beleg trägt seinen Satz nur indirekt“

### F-3 — Die `SEEN_SECOND`-Wartebedingung ist innerhalb des 15-s-Ruhefensters nicht eigenständig falsifizierbar

- `kategorie`: INFO
- `quelle`: Maintainability; `AGENTS.md` §3.12 Instanz B
- `pfad`: `sdks/python/pgchangefeed/integration/test_sse_filter_realserver.py` (Schleife vor `SEEN_SECOND`), ebenso C# und Kotlin
- `befund`: Die Zusage trägt die zweite Gruppe (Package-Mutation rot, Zählung rot, siehe Messungen). Würde die `SEEN_SECOND`-Schleife entfallen, begänne das Ruhefenster
  unmittelbar nach `SEEN`; die Zeilen der zweiten Gruppe treffen laut gemessenem Wert (`SEEN_SECOND` nach 324–379 ms) lange vor dem Ende von 15 s ein, die Mutation
  bliebe grün (**hergeleitet**, nicht gefahren). Die Bedingung senkt das Flake-Risiko bei langsamer Zustellung und garantiert die Reihenfolge der Gegenlesung, ist aber keine
  Eingabeseite, an der ein Test rot wird. Die Beschreibung im Plan („Ruhefenster beginnt erst, wenn …“) ist wahr; ein Leser könnte ihr eine Prüfwirkung zuschreiben.
- `verifizierbar`: nein
- `klasse`: „Wartebedingung ohne eigene Falsifikation“

### F-4 — Routing-Phase hat denselben Aufbau (Dreiergruppe nach `READY`, nur `SEEN` vor dem Ruhefenster) — vom Implementer gemeldet, nicht untersucht

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Meldung statt stiller Mitänderung); `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (Nachbar)
- `pfad`: `tools/harness/lib-sdk-route-fixture.sh` (Phasenablauf), Plan §3 Suchlauf-Feld Zeile 2
- `befund`: Die Meldung steht im Plan (`Betroffenheit nicht untersucht (hergeleitet), Entscheid beim Auftraggeber`); sie nennt keine Adresse eines Folge-Slice. Ob der
  Verbindungs-Fehler der Filter-Phase (Dreiergruppe vor stehender Verbindung aller Clients) dort ein Falsch-Grün erzeugt, ist offen.
- `verifizierbar`: nein
- `klasse`: „Aufschub ohne Adresse“

### F-5 — Kotlin: eine Wartebedingung über 190 Zeichen in einer Zeile; gemischte Referenz auf den ersten Sentinel

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/SseFilterRealserverTest.kt:90`
- `befund`: Die `while`-Bedingung der ersten Gruppe nutzt `first` und im selben Ausdruck `PhaseEnvironment.sentinel` (gleicher Wert); zwei weitere Zeilen der Datei
  überschreiten 140 Zeichen. `make fmt-check` deckt Kotlin nicht (nur Go). Keine semantische Auswirkung.
- `verifizierbar`: nein
- `klasse`: „Formatierung ohne Sensor“

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/lib-sdk-filter-fixture.sh` — Sentinel `<Sentinel>Second` (`PGCHANGEFEED_FILTER_SENTINEL_SECOND`); ID-Basis `id_base + 100` (Filter-Basis 1000 →
  1101–1103) gegen die Versuchs-IDs `id_base + attempt*3 + 1..3` (1004–1018) und gegen die Basen der anderen Phasen (600–900, andere Tabellen) kollisionsfrei; beide
  Gegenlesungen nehmen beide Sentinels; die Zeilenzahl-Assertion zählt über den Sentinel der zweiten Gruppe (F1 genau 1, F2 genau 1, ohne Filter genau 3, je Tabelle eine),
  ein Versuch > 1 verfälscht sie nicht, weil die Wiederholung nur die erste Gruppe (anderer Sentinel, eigene IDs) betrifft; fremde Nachzügler der ersten Gruppe zählt `f1_foreign`.
- geprüft, ohne Befund: Fristen — Test 90 s je Wartebedingung gegen Runner-Schleife 450 × 0,2 s plus Aufruf-Overhead (länger als 90 s, der Test scheitert zuerst und der Runner meldet
  „kein SEEN_SECOND“); Gesamtfrist `quiet + 120` s.
- geprüft, ohne Befund: die drei Testklassen — die Bedingung der zweiten Gruppe verlangt F1 Tabelle A, F2 zweites Schema, den Client ohne Filter alle drei, je mit dem Sentinel der
  zweiten Gruppe; `FILTER_RESULT` zählt vor den Assertions beide Gruppen (`OwnAny`/`ownAny`/`_own_any`); Kotlin `filterSentinelSecond` folgt der Namenskonvention der Nachbarn
  (`filterSchemaA`, `filterTableA`). Übersetzbarkeit: die C#- und Kotlin-Integrationsprojekte wurden im Tier-Lauf gebaut und liefen; Python importiert die neue Variable
  im Tier-Lauf.
- geprüft, ohne Befund: Runner-Änderungen (Kopf-/Phasenkommentare, Abdeckungs-Zeile, Schlusszeile) und der Satz in `harness/README.md` (drei Zeilen, nennt zweite Gruppe,
  `SEEN_SECOND`, Zahlen genau 1/1/3) gegen die Fixture: wahr.
- geprüft, ohne Befund: `AGENTS.md` §3.1 (kein Host-Werkzeug am Repo im Diff), §3.2 (keine Suppression), §3.7 (Kommentare im Indikativ, höchstens eine Kennung je Block;
  `make kommentar-kennungen DIFF=e900e5c3` ohne Kandidat), §3.12/§3.13 (Suchlauf-Feld: 22 Zeilen stimmen, Befunde je Zeile eingetragen), Commits tragen `LH-FA-SST-009`/`ADR-0133`, kein
  `SPEC-*` im Betreff.
- geprüft, ohne Befund: `sdks/` — kein Produktivcode, keine Version, keine interne Kennung in neuen Zeilen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Pipe-Grep unter pipefail verwirft einen Treffer · Beleg trägt seinen Satz nur indirekt · Wartebedingung ohne eigene Falsifikation · Aufschub ohne Adresse · Formatierung ohne Sensor

## Verdikt

**Merge-blockierend:** nein. Der Diff selbst trägt die Zusage (Package-Mutation und falsche Zählung rot mit gedruckter Zeile, drei unmutierte Tier-Läufe grün mit
`f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6`). F-1 liegt außerhalb des Diffs (Bestand der drei Runner, nicht Gegenstand dieses Slice) und bleibt ein eigener Auftrag an den
Planner; er blockiert die Änderung nicht, färbt aber jeden C#-Tier-Lauf gelegentlich rot. Weil F-1 ein offenes MEDIUM ist, habe ich die DoD-Zeile „Review durchgeführt“ im Plan
**nicht** auf `[x]` gezogen (Regel „kein offenes HIGH/MEDIUM“); der Planner entscheidet, ob F-1 in diesen Slice zurückgereicht oder als Folge-Slice benannt wird.

**Übergabe:** F-1 an den Planner (Folge-Slice oder Rückgabe), F-2 bis F-5 zur Kenntnis. Architect-Fragen: (1) Soll die Routing-Phase (F-4) einen eigenen Slice mit gleicher Härtung
bekommen oder erst nach einer Messung entschieden werden? (2) Soll die Fixture eine Option erhalten, die `RECEIVED_*`-Zeilen immer druckt (F-2), oder genügt der indirekte
Beleg? Der Report ersetzt keine Verifikation (Modul 11).
