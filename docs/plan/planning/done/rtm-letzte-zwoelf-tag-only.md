# Slice rtm-letzte-zwoelf-tag-only: die letzten zwölf RTM-Waisen — Tag-Nachtrag an bestehenden E2E-Belegen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — kein Closure-Kriterium jenseits der eigenen DoD.

**Bezug:** [`LH-FA-CON-001`](../../../../spec/lastenheft.md),
[`LH-FA-CON-002`](../../../../spec/lastenheft.md),
[`LH-FA-CON-003`](../../../../spec/lastenheft.md),
[`LH-FA-CON-004`](../../../../spec/lastenheft.md),
[`LH-FA-CON-005`](../../../../spec/lastenheft.md),
[`LH-FA-DAT-002`](../../../../spec/lastenheft.md),
[`LH-FA-DAT-003`](../../../../spec/lastenheft.md),
[`LH-FA-DAT-005`](../../../../spec/lastenheft.md),
[`LH-FA-REA-001`](../../../../spec/lastenheft.md),
[`LH-FA-RET-001`](../../../../spec/lastenheft.md),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md),
[`LH-QA-REL-004`](../../../../spec/lastenheft.md).

**Berührte Spec-Stellen:** keine — alle zwölf Anforderungen sind bereits
real durch bestehende E2E-Phasen/-Testfälle erfüllt; dieser Slice trägt
nur den fehlenden Kennungs-Tag nach.

**Verantwortlich:** — (wellenlos, direkt umgesetzt, siehe Autor).

**Autor:** Implementer-Agent, direkt beauftragt durch den Auftraggeber im
Anschluss an `rtm-reste-sst-cfg-por` ("Ja bitte weitermachen, bis alles
abgedeckt ist"). Kein separater Priorisierungs-Schritt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Nach `bench-schwellen-per-001-002-003` und `rtm-reste-sst-cfg-por`
verbleiben zwölf RTM-Waisen — laut `harness/README.md` §`make doc-trace`
bereits als „real bereits über bestehende E2E-Phasen mitbelegt, nur ohne
eigenen Kennungs-Tag am jeweiligen Testfall" klassifiziert. Dieser Slice
verifiziert das je Anforderung gegen den tatsächlichen Testcode (nicht nur
gegen die Klassifikations-Behauptung) und trägt den Tag an genau der
Stelle nach, die die Anforderung wirklich erfüllt:

1. `LH-FA-CON-001`/`003`/`004`/`005` (Registrierung, Positions-
   Persistierung, Bestätigung, Fortsetzung nach Neustart) — alle vier an
   der bestehenden „Black-Box-CLI-Rundlauf"-Phase
   (`tools/harness/run-integration-tests.sh`): `register-consumer`,
   `cdc.consumer_position`-Rücklesung nach Bestätigung, zwei
   `acknowledge-consumer`-Aufrufe, realer `docker restart` mit
   Fortsetzen ab der zurückgelesenen Position.
2. `LH-FA-CON-002` (Unabhängige Consumer) — an
   `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer`
   (`test/integration/integration_test.go`): zwei separat registrierte
   und bestätigte Consumer, deren Positionen sich nachweislich nicht
   gegenseitig beeinflussen (`cdc.retention_blockers`-Beleg).
3. `LH-FA-DAT-002`/`003`/`005` (Quelltabelle identifizierbar,
   Operationstyp, relevante Datenwerte) — an `TestE2ECaptureFlow`: bereits
   bestehende `SourceTableID`-, `Operation`- und Row-Image-Prüfungen.
4. `LH-FA-REA-001` (Lesebereich zwischen zwei Positionen) — an der
   bestehenden „HTTP-API-Rundlauf"-Phase: `GET /changes` mit einem
   echten `[from, to)`-Bereich um eine einzelne, eigens eingefügte Zeile.
5. `LH-FA-RET-001` (Persistente Speicherung) — an der bestehenden
   „Upgrade-Sicherheits-Rundlauf"-Phase: `count(*)`-Beleg, dass eine vor
   einem realen Container-Tausch erfasste Änderung danach identisch über
   `cdc.changes` lesbar bleibt.
6. `LH-QA-REL-003` (Sichtbarer Unzuverlässigkeitszustand) — an der
   bestehenden „CLI-Diagnose-Beleg (Fehlerzustand)"-Phase: exakt der in
   der Messmethode geforderte Fehlerinjektionstest über `LH-FA-ADM-003`.
7. `LH-QA-REL-004` (Idempotente Consumer-Verarbeitung) — an
   `TestE2ECaptureFlow`s bestehendem Wiederlese-Block (bereits
   `LH-FA-REA-005` zitierend): zweimaliges Lesen desselben Bereichs
   liefert dieselben Changes.

**Ausdrücklich NICHT in diesem Slice:**

- Neue Testfälle oder neue Prüfzeilen — jede der zwölf Anforderungen wird
  ausschließlich durch bereits bestehenden, bereits grün laufenden Code
  belegt. Wo eine Anforderung eine Boundary/Negative-Akzeptanzkriterium
  trägt, das der bestehende Testfall nicht abdeckt (z. B. `LH-FA-DAT-005`s
  Negative-Kriterium „Wert nicht lieferbar"), bleibt das unbelegt — echte
  inhaltliche Lücke, kein Tag-Nachtrag-Gegenstand (dasselbe Muster wie bei
  `LH-FA-SST-001`/`005`/`LH-QA-PER-004` in `rtm-reste-sst-cfg-por`).

## 2. Definition of Done

- [x] `LH-FA-CON-001`/`002`/`003`/`004`/`005` zeigen in `make doc-trace` `ok`.
- [x] `LH-FA-DAT-002`/`003`/`005` zeigen in `make doc-trace` `ok`.
- [x] `LH-FA-REA-001` zeigt in `make doc-trace` `ok`.
- [x] `LH-FA-RET-001` zeigt in `make doc-trace` `ok`.
- [x] `LH-QA-REL-003`/`004` zeigen in `make doc-trace` `ok`.
- [x] `make test-integration` grün, regeneriert `docs/user/e2e-abdeckung.md`
      mit allen zwölf neuen Tags.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review — 1 MEDIUM (F-1,
      Kommentarpräzisierung) in derselben Fixrunde behoben, 1 INFO (F-2,
      unter der 3×-Schärfungsschwelle) übernommen, kein offenes HIGH/MEDIUM.
- [x] Doku-Update: `harness/README.md` (`make doc-trace`-Zeile mit der
      real gemessenen finalen Waisenzahl — 0, real gemessen).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | `TestE2ECaptureFlow`s Doc-Kommentar trägt zusätzlich `LH-FA-DAT-002`/`003`/`005`/`LH-QA-REL-004`; `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer`s Doc-Kommentar trägt zusätzlich `LH-FA-CON-002`. |
| `tools/harness/run-integration-tests.sh` | update | „Black-Box-CLI-Rundlauf"-`abdeckung_declare` trägt zusätzlich `LH-FA-CON-001`/`003`/`004`/`005`; „HTTP-API-Rundlauf" trägt zusätzlich `LH-FA-REA-001`; „CLI-Diagnose-Beleg (Fehlerzustand)" trägt zusätzlich `LH-QA-REL-003`; „Upgrade-Sicherheits-Rundlauf" trägt zusätzlich `LH-FA-RET-001`. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | von `make test-integration` neu geschrieben. |
| `harness/README.md` | update | `make doc-trace`-Zeile mit der finalen real gemessenen Waisenzahl. |

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — direkter Auftrag, kein
externer Trigger.

**Rückführungen:**

- `in-progress` → `next` (zu groß): entfällt — Umfang bereits real
  umgesetzt (alle zwölf Tags gesetzt).
- `in-progress` → `open` (blockiert): entfällt aus demselben Grund.

## 5. Closure-Trigger

DoD vollständig + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Reiner Tag-Nachtrag ohne neue Prüfzeile trägt keinen zusätzlichen
  Regressionsschutz über das hinaus, was der jeweilige Testfall bereits
  vor diesem Slice belegte — dieselbe bewusste Einordnung wie bei den
  drei Tag-only-Punkten aus `rtm-reste-sst-cfg-por`. **Ausgang:**
  entfallen als Risiko, das ist die Art dieses Slice.
- Mehrere Boundary-/Negative-Akzeptanzkriterien der zwölf Anforderungen
  (z. B. `LH-FA-CON-001`s „erneute Registrierung ist idempotent oder
  definierter Fehler", `LH-FA-DAT-005`s „nicht lieferbarer Wert
  erkennbar") bleiben ohne eigenen Testbeleg — der Tag deckt die
  Anforderung über ihren Happy Path, nicht über jedes Akzeptanzkriterium.
  **Ausgang:** weiter offen, benannt als unsanierte Bestandslücke (wie
  bereits in `harness/README.md` §`make doc-trace` vor diesem Slice
  beschrieben) — Sanierung wäre ein eigener inhaltlicher Vorgang mit
  neuen Testfällen, kein Tag-Nachtrag.
- Review-Fund F-1 (MEDIUM,
  `docs/reviews/review-slice-rtm-letzte-zwoelf-tag-only.md`
  <!-- d-check:status-provenance -->): der `LH-FA-CON-002`-Tag an
  `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer` belegt nur
  einen Teilbereich — dass `aheadConsumer`s spätere Bestätigung
  `behindConsumer`s gespeicherte Position nicht überschreibt —, nicht die
  wörtliche Happy-Path-/Boundary-Formulierung aus `spec/lastenheft.md`
  (gleichzeitiges Lesen desselben Bereichs durch beide Consumer).
  **Ausgang:** eingetreten, in derselben Fixrunde behoben — der
  Kommentar an der Testfunktion benennt jetzt präzise den geprüften
  Teilbereich statt der vollen Anforderung.

## 7. Closure-Notiz

- **Was hat funktioniert:** Jeder der zwölf Tags wurde vor dem Setzen
  gegen den tatsächlichen Testcode gehalten, nicht nur gegen die
  Vorab-Klassifikation aus `harness/README.md` — das hat echte
  Substanz-Unterschiede sichtbar gemacht (z. B. `LH-FA-RET-001`s Beleg ist
  stärker als der Anforderungstext verlangt, `LH-FA-CON-002`s Beleg war
  ursprünglich schwächer als der Kommentar behauptete). Reiner
  Tag-Nachtrag ohne neuen Testcode hielt den Slice trotz zwölf
  Anforderungen klein und risikoarm.
- **Was ging anders als geplant:** Der Reviewer fand 1 MEDIUM (F-1: der
  `LH-FA-CON-002`-Kommentar überzeichnete den Testumfang — der Test zeigt
  nur, dass eine spätere Bestätigung die bereits gespeicherte Position
  eines anderen Consumers nicht überschreibt, nicht die volle
  Happy-Path-/Boundary-Formulierung aus dem Lastenheft), in derselben
  Fixrunde präzisiert (`d1136723`). Der Verifier fand danach beiläufig
  (nicht Teil des beauftragten Prüfumfangs) denselben Zeilen-Lokator-Drift
  wie beim vorangegangenen Slice: die Kommentarverlängerung verschob die
  Testfunktion um 3 Zeilen, ohne dass `docs/user/e2e-abdeckung.md` neu
  erzeugt wurde. Der Drift war zum Fundzeitpunkt bereits inzidentell durch
  einen parallelen Fixrunden-Commit des Nachbar-Slices behoben — real
  durch einen erneuten `make test-integration`-Lauf bestätigt
  ("unverändert").
- **Steering-Loop-Eintrag:** kein neuer Sensor — der Zeilen-Lokator-Drift
  ist ein weiterer, unmittelbar benachbarter Beleg derselben bereits
  verkörperten `AGENTS.md` §3.13-Regel (zweiter Fall in derselben
  Arbeitssitzung, siehe Beobachtungs-Register); bestätigt dieselbe Nuance
  wie beim Vorgänger-Slice: eine Kommentarverlängerung ist eine
  Nachbaränderung im Sinne der Regel, auch wenn sie keine Zahl im
  eigentlichen Sinne bewegt.
- **Beobachtungs-Register (`../observations/`):** neuer Beleg
  `evidence/slice-rtm-letzte-zwoelf-tag-only.md` unter der bestehenden
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger/`.
- **Folge-Slices:** keine — `make doc-trace` zeigt nach diesem Slice
  0 Waisen, alle 76 Anforderungen sind gedeckt.
- **Risiken aus §6:** drei Risiken, drei Ausgänge — (1) Tag-only ohne
  neuen Regressionsschutz → **entfallen**, das ist die Art dieses Slice;
  (2) Boundary-/Negative-Akzeptanzkriterien einzelner Anforderungen ohne
  eigenen Testbeleg → **weiter offen**, unsanierte Bestandslücke,
  Sanierung wäre ein eigener inhaltlicher Vorgang; (3) Review-Fund F-1
  (`LH-FA-CON-002`-Kommentar überzeichnete den Testumfang) →
  **eingetreten**, in derselben Fixrunde präzisiert.
- **Drei Paarungen:** Anker — kein `liegt in`-Feld (kein neuer
  Sensor/keine neue Regel verkörpert). Folge-Slice — keiner benannt,
  nichts zu prüfen. Register —
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger/` existiert mit
  nicht-leerem `evidence/` (der neue Beleg oben).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „E2E-Test-Harness"
(`tools/harness/run-integration-tests.sh`, `test/integration/`) — GF
(Modus-Deklaration `harness/conventions.md`, Default `PGC`/Greenfield).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
kein Treffer zu diesem konkreten Gegenstand.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
