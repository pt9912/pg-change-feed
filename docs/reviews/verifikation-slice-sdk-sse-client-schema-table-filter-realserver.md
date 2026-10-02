# Verifikations-Report: slice-sdk-sse-client-schema-table-filter-realserver — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + ADR-Konformität +
Plan-vs-Code-Diff + Gates, in frischem Kontext. Review-Artefakt:
[`review-slice-sdk-sse-client-schema-table-filter-realserver.md`](review-slice-sdk-sse-client-schema-table-filter-realserver.md)
(0 HIGH/0 MEDIUM; F-1 LOW, F-2 LOW, F-3 bis F-5 INFO). Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md).

**Gegenstand:** Slice-Plan
[`slice-sdk-sse-client-schema-table-filter-realserver`](../plan/planning/done/slice-sdk-sse-client-schema-table-filter-realserver.md)
(wellenlos), Diff-Range `448ee4a5..HEAD` (`c6c7aaf7`): Implementer-Commit `8b0e1e52`, Review-Commit
`3c6851df`, Stilfix `c6c7aaf7` (F-1), dazu Planner-Commits ohne Code. Bezug:
[`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 1 und 4,
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2,
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md).

Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Alle Mutationen liefen an
`git archive HEAD`-Kopien im Scratchpad (je Kopie `git init`, eigener Image-Tag `pg-change-feed-mutation:vf-<name>`,
nach dem Lauf `docker rmi`); Mutanten per `sed … Datei > Scratchpad-Temp` und `cp` in die Kopie, kein `sed -i`.
`git status --short` im Echtrepo war nach den Tier-Läufen, nach den Mutationen und nach den Sensoren leer.
Keine Berechtigungs-Verweigerung im Lauf ([`AGENTS.md`](../../AGENTS.md) §3.15 nicht ausgelöst).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | gedruckte Zeile |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.13.0 OK — 54 Dateien`; `d-check: 1549 Datei(en) geprüft, 0 Befund(e)` (docs-check und `commits`); `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`; `generated-sync: OK`; `sdk-public-doc-check: keine interne Kennung unter sdks`; a-check `gesamt: 0 Befund(e)` |
| `make test` | Exit 0 | alle Pakete `ok` (kein Go-Code im Diff, Gegenprobe) |
| `make fmt-check` | Exit 0 | `fmt-check: 323 Go-Dateien geprüft, alle formatiert` |
| `make docs-check` | Exit 0 | `d-check: 1549 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make sdk-public-doc-check` | Exit 0 | `keine interne Kennung unter sdks` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-sse-client-schema-table-filter-realserver.md` | Exit 0 | `suchlauf-nachmessen: 12 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=448ee4a5` | Exit 0 | keine Ausgabe, kein Kandidat |
| `make doc-trace` | Exit 0 | `80 Anforderung(en), 0 Waise(n).` (gemessen in diesem Lauf am Stand `c6c7aaf7`) |
| `make doc-commits RANGE=448ee4a5..HEAD` | Exit 0 | `d-check: 1549 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`) |
| `make doc-immutable RANGE=448ee4a5..HEAD` | Exit 0 | `d-check: 1549 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |

Nachtrag nach Anlage dieses Reports: `make docs-check` wird im Commit-Schritt erneut gefahren (Ergebnis siehe Verdikt).

## 2. Die drei Tier-Läufe (selbst gefahren, unmutiert)

Image: `ghcr.io/pt9912/pg-change-feed:dev` war geladen, `harness/image-hash.txt` nennt
`sha256:48e5699d9b2fbf9761fa1f853d6d2faa6b0ba3e43007e81b320e8a57d466a878` (identisch mit dem Digest im Review-Report).
Ich habe `make image` nicht erneut gefahren: `git diff 448ee4a5 HEAD --stat -- internal cmd proto gen Dockerfile compose.yaml`
ist leer, der Server-Code ist durch den Slice unverändert. Dass das geladene Image aus genau diesem Baum stammt, ist **übernommen**
(Review-Report), nicht von mir nachgebaut.

| Tier | Exit | gedruckte Zeile (Beleg, nicht der Bau-Exit) |
|---|---|---|
| `make test-sdk-csharp-integration` | 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15`, `SEEN nach 337ms ab dem letzten Commit (Versuch 1)` |
| `make test-sdk-kotlin-integration` | 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15`, `SEEN nach 335ms … (Versuch 1)` |
| `make test-sdk-python-integration` | 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15`, `SEEN nach 369ms … (Versuch 1)` |

Die Gegenlesung über `cdc.changes` (`schema_name`, `table_name`, Sentinel) ist Teil des Runners und bestand in allen drei Läufen
(sonst Exit ≠ 0). Die zwölf Alt-Phasen liefen in allen drei Gesamtläufen grün (Exit 0); der gemeldete NATS-Flake (Review F-4)
trat in meinen drei Läufen nicht auf (kein Wiederholungslauf, keine Aussage über die Häufigkeit). Idempotenz der Abdeckung:
`git status --short` nach den drei Läufen leer, `docs/user/sdk-e2e-abdeckung.md` blieb unverändert.

## 3. Mutationen (auf Kopien; gesehene Farbe mit gedruckter Zeile)

| Nr. | Sprache | Stelle | Farbe | gedruckte Zeile |
|---|---|---|---|---|
| P (Package) | Python | `sse_client.py:90`: `("table", table)` aus der Query-Menge entfernt **und** `sdks/python/Dockerfile` `RUN pytest` → `RUN true` (Unit-Stufe umgangen, Diff in `vfm/py-P.diff`) | rot, Exit 2, **in der Filter-Phase** | `FILTER_RESULT f1=2 f1_foreign=1 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15` |
| P (Package) | Kotlin | `PgChangeFeedSseClient.kt:99`: `"table" to (null as String?)` **und** Dockerfile `gradlew test` → `RUN true`, `build` → `build -x test` | rot, Exit 2, in der Filter-Phase | `FILTER_RESULT f1=2 f1_foreign=1 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15` |
| M3 | Python | `test_sse_filter_realserver.py:142`: F2 mit `_SCHEMA_A` statt zweitem Schema | rot, Exit 2 | keine `FILTER_RESULT`-Zeile; Runner: „der Test meldete nach 5 Dreiergruppen kein SEEN“, pytest: `AssertionError: innerhalb der Frist weder die Change … des zweiten Schemas …` |
| M3 | C# | `SseFilterRealserverTests.cs:90`: F2 mit `FilterSchemaA` | rot, Exit 2 | Runner: „der Test meldete nach 5 Dreiergruppen kein SEEN“, `Failed: 1` |

Hinweise zur Lesart:

- **Ein erster Python-M3-Versuch war eine Nullmutation** (mein `sed`-Muster traf nicht, `vfm/py-M3.diff` leer): Exit 0 mit
  `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 …` — er zählt als zusätzlicher unmutierter Python-Lauf, nicht als Mutation.
  Die berichtigte Mutation (Diff gezeigt: eine Zeile) ist die in der Tabelle genannte.
- **M3 färbt über die positive Zusage, nicht über einen Fremdwert**: ein Client, der das falsche Schema bekommt, erreicht
  `IsOtherSchema` nie, der Test endet erst an seiner 90-s-Frist, der Runner bricht bereits nach 5 × 12 s mit „kein SEEN“ ab.
  Der Plan erwartete „Gegenlesung rot“; tatsächlich ist die Farbe die fehlende SEEN-Bedingung. Die Mutation ist damit gebunden,
  der Mechanismus ist ein anderer als der im Plan genannte.
- **Review F-5 geschlossen:** Mit umgangener Unit-Stufe trägt die Filter-Phase allein die Bindung des Packages — sowohl in Python als auch in
  Kotlin färbt sie rot mit gezähltem Fremdwert 1. Die Package-Mutation in einer Kopie ist gangbar (Plan §6 „Ob der Tier-Bau eine
  Kopie annimmt, ist nicht gemessen“ ist damit gemessen: ja, mit `git init` in der Kopie und eigenem Image-Tag).
- **Mutations-Matrix aus allen mir sichtbaren Belegen** (Review-Report + dieser Lauf): C#: M1, M2 (Review), M3 (hier); Kotlin: M1, M2 (Review),
  Package (hier); Python: M1 (Review), M3, Package (hier). **Nicht erprobt**: C# Package, Kotlin M3, Python M2. Der Plan verlangt je
  Sprache M1 bis M3 und eine Package-Mutation; diese drei Zellen sind in keinem mir vorliegenden Artefakt belegt (der Implementer-Bericht
  liegt mir nicht als Datei vor). Für die Closure-Notiz gilt: erprobt oder als **hergeleitet** gekennzeichnet
  ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B).

## 4. DoD-Abgleich, Zeile für Zeile

| DoD-Zeile | Befund |
|---|---|
| LP 1 — Hilfsdatei + C#-Tier | **erfüllt**: `tools/harness/lib-sdk-filter-fixture.sh` vorhanden (Anlage der drei Tabellen inkl. zweitem Schema, `enable_table` mit Poll auf `applied`, Phasenablauf mit Dreiergruppen B, zweites Schema, A, `SEEN`, `FILTER_RESULT` vor den Assertions im Test, Gegenlesung `cdc.changes`); Test `SseFilterRealserverTests` gelesen: F1/F2/U wie im Plan, Ruhefenster beginnt nach Empfang aller drei Gruppen am Client ohne Filter. Grüner Lauf siehe §2. Mutationsteil: M1, M2, M3 belegt (Review + hier), Package-Mutation C# **nicht erprobt** |
| LP 2 — Kotlin-Tier | **erfüllt** bis auf Mutations-Zelle Kotlin M3 (nicht erprobt); Lauf grün, Package-Mutation rot gesehen |
| LP 3 — Python-Tier | **erfüllt** bis auf Mutations-Zelle Python M2 (nicht erprobt); Lauf grün, M1 (Review), M3 und Package-Mutation rot gesehen |
| Mutationsproben je Sprache | **teilweise belegt**: neun von zwölf Zellen erprobt (§3), drei offen; kein Hindernis, nur nicht gefahren |
| Nur Test-Code und Runner | **erfüllt**: `git diff 448ee4a5 HEAD --name-only -- sdks` nennt fünf Pfade (`PhaseEnvironment.cs`, `SseFilterRealserverTests.cs`, `PhaseEnvironment.kt`, `SseFilterRealserverTest.kt`, `test_sse_filter_realserver.py`), keine Versionsdatei; `docs/user/version.md` unverändert; kein Tag; `make sdk-public-doc-check` Exit 0 |
| Abdeckung getragen | **erfüllt**: `git diff 448ee4a5 HEAD -- docs/user/sdk-e2e-abdeckung.md` = genau drei Einfügungen (je Runner eine Zeile, Kennungen `LH-FA-SST-008`/`LH-FA-SST-009`, Nachweis `SseFilterRealserverTests`/`SseFilterRealserverTest`/`test_sse_filter_realserver.py`); zweiter Lauf je Tier ändert nichts (§2); `make doc-trace` 80/0; `make docs-check` Exit 0 |
| `make gates` grün | **erfüllt** (§1, Exit 0 direkt gesichert) |
| Review durchgeführt | **erfüllt** (Report vorhanden, 0 HIGH/0 MEDIUM; Häkchen im Plan bereits gesetzt) |
| §3.13-Suchlauf | **erfüllt**: Feld in §3 trägt Gefundenes und Nichtgefundenes, `suchlauf-nachmessen` Exit 0 (12 Zeilen) |
| Doku-Update | **erfüllt**: `harness/mk/sdk.mk` (drei Hilfetexte „dreizehn Phasen … eine mit Tabellenfilter“), `harness/README.md` §Sensors (drei Zeilen), Runner-Köpfe; Handbuch unberührt |
| Closure-Notiz, Beobachtungs-Register, §6-Ausgänge, drei Paarungen | **offen** — Closure-Handlungen des Planners, noch nicht fällig vor dem Verdikt (§6 unten) |

Entscheidungs-Konformität: [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4 (SSE nutzt `schema`/`table`) und
Teilfrage 1 (Kombinatorik) sind am Realserver in den drei Packages belegt: `schema` + `table` liefert genau die Tabelle, `schema` allein das Schema,
das zweite Schema mit gleichnamiger Tabelle wird getrennt; die vom Plan als `table` allein ausgeschlossene Form bleibt Unit-Ebene.
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2 (Mechanik): Prüfling ist das Package, Gegenlesung
unabhängig über `cdc.changes`. Plan-vs-Code: keine Abweichung vom Plan-Umfang gefunden (Hilfsdatei, drei Runner, drei Tests, zwei
`PhaseEnvironment`-Hilfen plus Python-Test; kein Dockerfile-, Compose- oder Workflow-Diff; kein Produktivcode).

## 5. Bewertung Review-F-2 (Verbindung der Clients vor dem ersten Commit)

Code-Lesung: In allen drei Tests druckt `READY` unmittelbar nach dem Start der drei Konsumenten, nicht nach stehender Verbindung; der Runner
wartet nur auf `READY` (`lib-sdk-filter-fixture.sh:99–112`) und committet dann die Gruppe in der Reihenfolge B, zweites Schema, A. Das
Positiv-`SEEN` verlangt von F1 die Change der Tabelle A, von F2 die des zweiten Schemas, von U alle drei. `SEEN` wird bereits im Versuch 1
erreicht und beendet die Schleife (`:134`); das Ruhefenster danach committet nichts mehr.

Schwäche: Ein gefilterter Client, dessen Verbindung zwischen zwei Commits der ersten Gruppe zustande kommt, bleibt ohne Fremdwert, weil ihn die
fremde Change nie erreichte, und kann dennoch `SEEN` erfüllen. Ein Package-Defekt (`table` nicht auf dem Draht) bliebe dann unentdeckt.
Das Fenster ist die Zeit zwischen den drei `INSERT`s in einer `psql`-Sitzung (Millisekunden) gegenüber dem Vorlauf aus Docker-Log-Polling und
`docker exec`-Start (hunderte Millisekunden); in allen acht gesehenen Läufen (drei unmutiert, vier Mutationen, Review) wurde das Fenster nicht
getroffen, jede Mutation zählte den Fremdwert. Das ist eine Eintrittswahrscheinlichkeit, kein Beweis der Abwesenheit; die Aussage ist
**hergeleitet**, ein Lauf mit verzögerter Verbindung wurde nicht gefahren.

**Urteil:** Robustheit genügt für das Verdikt (Falsch-Grün nur bei gleichzeitig defektem Package und Treffer eines Millisekunden-Fensters;
die Positivzusage und die Gegenlesung bleiben unberührt) — **keine Nachbesserung als Bedingung**. Eine Nachbesserung ist billig und
empfohlen, sie braucht keine Test-Änderung: Nach `SEEN` committet der Runner (`lib-sdk-filter-fixture.sh`, hinter der Versuchsschleife)
**eine zweite Dreiergruppe** (neue IDs, gleicher Sentinel). Weil F1 und F2 ihre eigene Change der ersten Gruppe nachweislich empfangen haben,
sind sie bei der zweiten Gruppe sicher verbunden; deren fremde Changes liegen im Ruhefenster, und `FILTER_RESULT` zählt sie
(dann `f1=2`, `f2=2`, `unfiltered=6`; die Zeilenzahlprüfung `lines == count` bleibt gültig). Alternativ nennt die Closure-Notiz die Grenze als
akzeptiert und als hergeleitet.

Zu F-3 (Reinheit statt Vollständigkeit): richtig gelesen; `f1 ≥ 1`, `f2 ≥ 1` und die vollständige Dreiergruppe am Client ohne Filter tragen die
Zusage „genau seine Auswahl“ nur für die Reinheit. Die Vollständigkeit am Filter-Client ist durch die Ein-Gruppen-Form nicht über 1 hinaus
gemessen; mit der zweiten Gruppe wäre `f1 == 2` zusätzlich prüfbar. Nicht blockierend.

## 6. Register-Fortschreibung für die Closure (melden, nicht ändern)

Beobachtungen unter [`docs/plan/planning/observations/BEO-PGC/`](../plan/planning/observations/BEO-PGC/) (Stand dieses Laufs gelesen):

- [`integrationsprojekt-uebersetzt-nicht-unbemerkt`](../plan/planning/observations/BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt/state.md)
  (heute **1×**): Der Kotlin-Fund F-1 (fehlendes Leerzeichen) ist **kein** Beleg — es war ein Formatierungsfund, die Datei übersetzte. Ein
  Übersetzungsfehler ist in diesem Slice nicht aufgetreten; die drei Tier-Läufe sahen die Übersetzung selbst. Fortschreibung: **keine
  neue `evidence/`-Datei**; der Plan nennt den Eintrag als „weiter offen → Register“ (§6), der Ausgang dort ist „kein Auftreten“.
- [`formatierungs-drift-ohne-gate`](../plan/planning/observations/BEO-PGC/formatierungs-drift-ohne-gate/state.md) (heute **3×**, verkörpert:
  `make fmt-check`, Go-Dateien): Review-F-1 ist die Klasse „Formatierung, Nicht-Go“, die `make fmt-check` nicht deckt (Kotlin). Das ist ein
  weiteres Auftreten bei gelaufenem Schritt 18 des Implementers an einem Träger, den der Sensor nicht liest — der im Eintrag genannte Trigger
  für die Gate-Aufnahme („ein weiteres Auftreten, das der Reviewer trotz gelaufenem Schritt 18 findet“) ist **dem Wortlaut nach nur für Go
  gemeint**; ob ein Kotlin-/Nicht-Go-Fund denselben Trigger zählt, ist eine Auslegungsfrage für den Planner/Architect. Empfehlung: `evidence/`-Datei
  (4×) mit der Auslegung, kein Gate.
- [`docker-cache-ueberspringt-tests-still`](../plan/planning/observations/BEO-PGC/docker-cache-ueberspringt-tests-still/state.md) (heute 2×):
  **kein Auftreten.** Die Beleg-Form des Slice (gedruckte `FILTER_RESULT`-Zeile mit Fremdwert) hat getragen, und die Package-Mutationen
  druckten eine geänderte Zeile — ein Cache-Treffer ohne Quelländerung ist nicht eingetreten. Ergänzung für den Hinweis-Absatz: die
  Tier-Bauten mit umgangener Unit-Stufe laufen über eine eigene Dockerfile-Kopie und eigenen Image-Tag, die Mutationsläufe berühren die
  Tag-Images der echten Tier-Bauten nicht (anders als der Hinweis im Eintrag für den Vorgänger-Slice).
- [`negativtest-ohne-bindung-an-seine-eingabe`](../plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/state.md)
  (verkörpert, über dem Deckel 14×): Review-F-2 (Vorbedingung des Negativ-Belegs nicht nachgewiesen) und die Lücke der M3-Lesart sind Funde ≤ LOW
  an einem bekannten Träger-Typ vor dem Merge — nach dem Deckel kein `evidence/`, sondern **Finding-Kennung in der Closure-Notiz**.
- [`drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md):
  **kein Auftreten** in diesem Slice (ein Phasenablauf, ein Eingabesatz). Die gesehene Auffälligkeit ist eine reine Vermeidung (Hilfsdatei).
  Der Zählerstand des Eintrags ist hier nicht nachgezählt.

Weitere Closure-Pflichten: §6 des Plans braucht je Risiko einen Ausgang — Risiko „Zweites Schema“: **entfallen** (alle drei Tabellen wurden
in allen sechs Läufen aktiviert, `applied`); Risiko „Mutation am Package ohne Arbeitsbaum-Änderung“: **entfallen/gangbar** (§3); Risiko
„Negativ-Beleg über Abwesenheit“: Mutationen und Zähler tragen, F-2 als benannte, hergeleitete Grenze; „Docker-Cache“: kein Auftreten (s. o.);
„Kein Release“: entfallen. Die drei offenen Mutations-Zellen (§3) stehen in der Closure-Notiz entweder erprobt oder als hergeleitet.

## 7. Verdikt

**bestanden** — mit den Bedingungen B-1 und B-2 für die Closure.

- **B-1 (Closure, Pflicht):** Die Closure-Notiz führt die Mutations-Matrix (§3) und kennzeichnet die drei nicht erprobten Zellen (C# Package,
  Kotlin M3, Python M2) als **hergeleitet**, oder sie werden gefahren. Die Beschreibung der M3-Farbe nennt „kein SEEN“ statt „Gegenlesung rot“.
- **B-2 (Closure, Pflicht):** F-2 (Verbindung vor erstem Commit) bekommt in der Closure-Notiz eine Antwort: Nachbesserung (zweite Gruppe nach
  `SEEN`, §5) **oder** benannte, als hergeleitet gekennzeichnete Grenze. Ein Merge-Blocker ist F-2 nicht.
- Offen für den Planner (keine Bedingung des Verdikts): Zuordnung des Kotlin-Formatfunds zu `formatierungs-drift-ohne-gate` (§6), Ausgänge der
  §6-Risiken, Closure-Notiz mit Lerneintrag, Beobachtungs-Register, DoD-Häkchen (hier gesetzt: **keine**).

`make docs-check` nach Anlage dieses Reports: siehe Commit-Schritt (Exit 0 verlangt).
