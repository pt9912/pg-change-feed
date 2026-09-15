# Review-Report: slice-075 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-075`, Commit `6be714d`. Der Slice-Diff ist genau dieser
Commit (13 Pfade, `git show --stat 6be714d`); sein unmittelbarer Elter ist
`bede265` (Register-Zug), die zwischen `f32d26b` und `6be714d` liegenden Züge
(`bede265`, `7ef760b`, `f90ba96`, `b37cd9f` — Verweisform und Register) sind
Rollen-Arbeit anderer Kontexte und **nicht** Teil des Gegenstands.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14 — die Regel `Neue Betreiber-Oberfläche ohne Handbuch-Zug` liegt
**vor** diesem Implementer-Lauf).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-075-dauerhafter-ausschlussstand-wiedereinspielung.md`
  vollständig (§1–§8), einschließlich der fünf §3-Nachzüge aus `6be714d`
- [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  vollständig — §Entscheidung Festlegungen 1–4, §Verglichene Alternativen
  (Option E), §Konsequenzen, §Fitness Function, §Re-Evaluierungs-Trigger
- [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) §Teilfragen 1–5
  (von `ADR-0065` nur in der Dauerhaftigkeits-Aussage superseded)
- [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  (Antrags-Queue als einziger Schreibpfad administrativer Zustände),
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt),
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Image-Digest),
  [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (Paketstruktur)
- `docs/reviews/architect-verdict-spaltenausschluss-dauerhaftigkeit.md`
  (Befundlage, aus der der Slice entstand) und `docs/reviews/review-slice-067.md`
  (F-1/F-2 — die Findings, die dorthin gingen)
- Vorgänger in `done/`: `slice-066`, `slice-067`, `slice-068`
- [`LH-FA-CFG-005`](../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../spec/lastenheft.md)
  (Lastenheft), [`SPEC-019`](../../spec/pflichtenheft.md) (Pflichtenheft),
  [`ARC-004`](../../spec/architecture.md) (Sicht)
- `AGENTS.md` §3 (Hard Rules), `harness/conventions.md` (MR-000/MR-001)

---

## Findings

### F-1 — Der neue Beleg-Baustein der Werkzeuge-Zeile trägt als einziger keinen Herkunfts-Vermerk

- `kategorie`: LOW
- `quelle`: Maintainability — Konsistenz der Werkzeuge-Zeile
  `harness/README.md` · `AGENTS.md` §3.7 (die Herkunft in **ein** auflösbares
  Feld) · Vorbild der neun vorangehenden Beleg-Bausteine derselben Zeile
- `pfad`: `harness/README.md:132`
- `befund`: Alle neun vorangehenden „Zusätzlich ein …"-Abschnitte dieser Zeile
  enden mit `· seit slice-037`/`-038`/`-051`/`-061`/`-071`/`-062`/`-063`/`-068`/
  `-072`; der neue Abschnitt („Zusätzlich ein Spaltenausschluss-Neustart-Beleg
  (`ADR-0065`, trägt `LH-QA-SEC-004`)") endet mit „… vor der Container-Ende-Grenze
  und vor dem Upgrade-Tausch" und trägt keinen Slice-Vermerk. Die Herkunft ist
  über die im selben Abschnitt zitierte `ADR-0065` auflösbar — die Zeile bleibt
  damit lesbar; es fehlt die Gleichförmigkeit, nicht die Information.
- `verifizierbar`: nein — `make docs-check` prüft Referenzen, nicht diesen
  Vermerk; die Einheitlichkeit ist grep-prüfbar
  (`grep -o '· seit slice-[0-9]*' harness/README.md`).
- `klasse`: „Herkunfts-Vermerk am Beleg-Baustein der Werkzeuge-Zeile fehlt"

### F-2 — Ein Lesefehler des Ausschlussstandes im Aktivierungs-Zweig lässt eine im Datensatz aktivierte, im Assembler aber ungebundene Tabelle zurück

- `kategorie`: INFO
- `quelle`: [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  Festlegung 1/3 · `SPEC-008` (Start-/Antrags-Fehlerklassen)
- `pfad`: `internal/bootstrap/wiring.go:1171-1179`;
  `internal/bootstrap/administration_internal_test.go:840-880`
  (`TestProcessAdministrationRequestsMarksFailedWhenExclusionReadFails`)
- `befund`: Der Aktivierungs-Zweig liest den abgeleiteten Stand **nach** dem
  erfolgreichen `EnableTableUseCase`; scheitert dieser Lese-Zug, endet der
  Antrag `failed`, während die Bindungs-Zeile und die Publication real
  geschrieben sind — die Tabelle ist damit „aktiviert, aber nicht erfasst", bis
  ein erneuter `enable_table`-Antrag oder ein Neustart greift. Die Form ist
  vorbestehend (dieselbe Lage entsteht beim `Registered`-/`CurrentVersion`-
  Fehlschlag), der neue Pfad fügt ihr einen weiteren Auslöser hinzu; die
  Richtung ist **fail-closed** (lieber keine Bindung als eine ohne den
  geführten Ausschluss) und durch den genannten Test gepinnt. Kein
  Handlungsbedarf am Diff; benannt, weil der Auslöser neu ist.
- `verifizierbar`: ja — `make test` (der genannte Test) ·
  `make test-integration` (realer Antragsweg)
- `klasse`: „Teil-Aktivierung bei Fehler nach dem Enable-Schritt"

### F-3 — Die starke Nicht-Vakuum-Aussage zum Neustart-Beleg hängt an einem Lauf, dessen Beleg nirgends abgelegt ist

- `kategorie`: INFO
- `quelle`: [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
  (Digest als Lauf-Beleg) · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-08-agentenrollen.md` §Die neun Übergaben (kein Pfeil ohne
  benennbares Artefakt)
- `pfad`: `docs/plan/planning/in-progress/slice-075-dauerhafter-ausschlussstand-wiedereinspielung.md`
  §3 (Nachzug `harness/image-hash.txt`, Klammer „— siehe Bericht")
- `befund`: Der §3-Nachzug schreibt, der erste `make test-integration`-Lauf sei
  gegen das Vorstands-Image gelaufen und habe den Neustart-Beleg rot gefärbt —
  „siehe Bericht", ohne Adresse auf ein Repo-Artefakt; der Lauf ist
  reproduzierbar nicht mehr, weil `make image` das `:dev`-Image danach ersetzt
  hat. Die Aussage trägt dennoch: der Elter-Stand `6be714d^` liest in
  `activatedTableBindings` **keinen** Ausschlussstand
  (`git show 6be714d^:internal/bootstrap/wiring.go`) und baut den Assembler bei
  jedem Neustart neu auf — auf dem alten Image muss die Zusicherung
  „ausgeschlossener Schlüssel fehlt" fehlschlagen. Die Nicht-Vakuum-Hälfte ist
  damit unabhängig vom verlorenen Lauf belegt; nicht belegt ist nur die
  Fehlfarbe selbst.
- `verifizierbar`: nein — der Lauf ist nicht mehr rekonstruierbar; die
  Nicht-Vakuum-Hälfte ist per Code-Vergleich `6be714d^` ↔ `6be714d` prüfbar
- `klasse`: „Lauf-Beleg ohne abgelegtes Artefakt (Behauptung trägt unabhängig)"

---

## Urteile zu den gestellten Prüfpunkten

**1. Trägt die Umsetzung `ADR-0065`s Entscheidung? — Ja.** Geprüft am Code, drei
Hälften:

- **Keine zweite Quelle, kein neues Schema-Objekt, kein zweiter Schreibpfad.**
  `git show --stat 6be714d` nennt kein `tools/schema/**`; `git show 6be714d --
  tools/schema/` ist leer. Die einzigen Schreib-Zugriffe auf `cdc.*` bleiben die
  vier SQL-Funktionen (`tools/schema/nacharbeit-administration.sql`), unberührt.
  Der neue Zugriff ist ein reines `SELECT` (`SelectAppliedColumnRequests`,
  `internal/adapters/driven/postgresstorage/queries/queries.go:312-319`); kein
  `INSERT`/`UPDATE`/`DELETE` ist im Diff. Die `Assembler`-Bindung bleibt, was
  `ADR-0065` §Entscheidung sagt: Laufzeit-Cache.
- **Ein Mechanismus für beide Auslöser.** Beide Aufrufer rufen dieselbe
  Fähigkeit: der Prozessstart (`internal/bootstrap/wiring.go:351`, je Quelle)
  und der Aktivierungs-Zweig (`internal/bootstrap/wiring.go:1171`); beide
  übergeben den abgeleiteten Stand in dasselbe Feld
  (`wiring.go:367`, `wiring.go:1178`).
- **Fähigkeit am bestehenden Zuschnitt.** Die Methode sitzt am vorhandenen
  `ColumnExclusionPort` (`internal/application/port/outbound/columnexclusion.go:37`),
  kein zweiter Port — genau die in `ADR-0065` §Konsequenzen verlangte Form.

**2. Der In-Prozess-Schutz aus Festlegung 4 ist unangetastet.** `mapper.go` kommt
im Diff nicht vor (`git show --stat 6be714d`): `AddBinding`-Merge
(`mapper.go:406-413`) und `setSchemaVersion` (`mapper.go:422-431`) tragen
unveränderten Quelltext, `ExcludeColumn`/`IncludeColumn` ebenso. Nachgeprüft über
`git diff 6be714d^ 6be714d -- internal/adapters/driving/replication/mapper/`
— leer.

**3. Die Ableitungslogik trägt — und der Zweitschlüssel ist deterministisch.**
`ORDER BY requested_at, administration_request_id`
(`queries.go:318`). `administration_request_id` ist Primärschlüssel
(`tools/schema/schema.yaml:228`), also ist das Paar `(requested_at, id)` total
geordnet und die Auswertungsreihenfolge innerhalb einer Datenbank eindeutig —
der Fall „zwei Anträge mit identischem `requested_at`" hat genau eine
Reihenfolge. Die Auswertung
(`internal/adapters/driven/postgresstorage/tableactivation.go:107-142`) trägt
`exclude_column` dedupliziert ein, `include_column` heraus; wird der letzte
geführte Name genommen, `delete(excluded, qualified)` — **kein** Map-Eintrag,
keine leere Liste (Zeile 129-133). Für eine Tabelle ganz ohne Ausschluss ist
`excluded[qualified]` der Nullwert `nil` und geht als `nil` in die Bindung.

**4. Der reale Neustart-Beleg steht an der richtigen Stelle und bindet die
Zusage.** `tools/harness/run-integration-tests.sh:440-510` liegt **vor** der
Container-Ende-Grenze (`TestE2ESchemaChangeDropColumn`/
`-IncompatibleTypeChange`, ab Zeile 2213) und **vor** dem Upgrade-Container-Tausch
(`:COMPOSE up -d --force-recreate`, Zeile 2157) — die Ordnungs-Aussage der
`harness/README.md`-Zeile stimmt. Der Beleg bindet drei Aussagen an derselben
Zeile: der ausgeschlossene Schlüssel fehlt (`jsonb_exists … = f`), die **nicht**
ausgeschlossene Spalte steht darin (`name = 'ColumnAfterRestart'` — die
Kontrolle gegen ein leeres Row Image), und der ausgeschlossene Wert steht
nirgends im persistierten Change (`LH-QA-SEC-004`). Die Erfassung selbst ist
vorgeschaltet belegt (Poll auf `count(*) = 1`), der Beleg ist also nicht durch
eine ausgefallene Erfassung grün. Der Neustart ist real (`docker restart`,
Prozess-Neustart) und trifft eine Tabelle ohne `CDC_TABLES`-Eintrag — ihr
Bindungs-Stand kommt aus `cdc.source_table`, ihr Ausschlussstand ausschließlich
aus dem geprüften Weg.

**5. Die vier gesetzten Mutationen — selbst gesetzt, selbst rot gesehen.** Ich
habe **vier** gesetzt (je eine für die vier Zusagen), jede einzeln, jede mit
anschließender Rücknahme und Blatt-Prüfung:

| # | Mutation | Ort | Ergebnis |
|---|---|---|---|
| 1 | Zweitschlüssel aus dem `ORDER BY` | `queries.go:318` | `make test-store` Exit 2 — `TestTableActivationExcludedColumnsDerivesAppliedColumnRequests`: „Ausschlussstand … = [], wollen [secret]" |
| 2 | `include_column` löscht den Eintrag nicht mehr | `tableactivation.go:129-133` | `make test-store` Exit 2, derselbe Test: „… wollen keinen Eintrag (include_column hat den letzten Namen genommen)" |
| 3 | Startpfad trägt den Stand nicht in die Bindung | `wiring.go:367` | `make test` Exit 2 — `TestActivatedTableBindingsCarriesExcludedColumns`: „ExcludedColumns = [], wollen [secret]" |
| 4 | Aktivierungs-Zweig trägt den Stand nicht in die Bindung | `wiring.go:1178` | `make test` Exit 2 — `TestProcessAdministrationRequestsDisableEnableCycleRestoresExclusion`: „Row Image nach dem disable/enable-Zyklus = {\"id\":\"1\",\"secret\":\"geheim\"}, wollen ohne den ausgeschlossenen Schlüssel secret" |

Rücknahme: `git status --porcelain` nach jeder Rücknahme leer (nur das vom
Sensor-Lauf geschriebene `tools/schema/plan.yaml`, s. u.); Blatt-Identität
geprüft — `git hash-object` = `git rev-parse HEAD:<pfad>` für
`queries.go` (`d4d9dcd`), `tableactivation.go` (`6196ca5`), `wiring.go`
(`3380a76`). **Alle vier Zusagen des Implementers tragen:** die Mutationen
1/2 färben den Store-Test, 3/4 die beiden Bootstrap-Tests; die Paarung
Mutation ↔ Test ist in jeder Richtung dieselbe, die der Plan §2 nennt.

**6. Der reale Fehlschlag und die Reihenfolge `make image` → `make
test-integration`.** Nachgeprüft: `harness/image-hash.txt` trägt
`sha256:9b7937d1…`, und genau diese ID trägt das lokale Image
`ghcr.io/pt9912/pg-change-feed:dev` (`docker images --no-trunc`), gebaut
**02:55:42**; der Commit `6be714d` liegt bei **02:57:38**. Das `:dev`-Image ist
also zwei Minuten vor dem Commit aus demselben Baum gebaut worden, und der
Digest im Diff ist der dieses Baus — die Reihenfolge ist damit am Artefakt
belegt, nicht nur behauptet (Grenze: ein Digest ist laut `ADR-0044` ein
Lauf-Beleg, kein Inhalts-Fingerabdruck; die Zeitfolge trägt hier als
Korroboration, nicht als Beweis). Zur **Nicht-Vakuum-Frage**: siehe F-3 — sie
trägt, gestützt auf den Vergleich `6be714d^` ↔ `6be714d`.

**7. `spec/pflichtenheft.md` — die Abweichung trägt, kein Finding.** Drei
Gründe: (a) Der Text entscheidet inhaltlich nichts — er schreibt die vier
Festlegungen von `ADR-0065` aus (dauerhaft vermerkt · `applied`-Zeilen als
einzige Herkunft · `requested_at` mit Zweitschlüssel · `applied` statt `failed`
ohne laufende Bindung); `ADR-0065` führt den `SPEC-019`-Fließtext in
`Schärft:` selbst als Zielort, die Änderung ist also **präzisieren**, nicht
erweitern. (b) Die Rollen-Attribution der ADR („Planner-/Architect-Zug") ist
keine Hard Rule, sondern die Adresse einer Folgepflicht; der Slice-Plan ist der
Planungs-Träger und ordnet sie seinem DoD zu, und der Implementer benennt die
Abweichung in §3 ausdrücklich statt still. (c) §2 DoD-Punkt 6 und §7 Historie
sind beide da (`spec/pflichtenheft.md:339-352`, `:491`).
**Zur zweiten Hälfte — braucht die neue Lesefähigkeit einen eigenen
`SPEC-*`-Eintrag? Nein.** `ADR-0065` Festlegung 2 bindet die Form der
Lesefähigkeit ausdrücklich an die Bedingung „— soweit sie einen öffentlichen
Vertrag bilden — des Pflichtenhefts"; die Fähigkeit ist eine interne
Go-Schnittstelle am bestehenden Port, ohne Betreiber-Oberfläche. Die Präzedenz
`SPEC-021` trägt hier nicht: dort war die Form ein **externer** Vertrag (neuer
Endpunkt, Drahtformat, Status-Code). Die extern sichtbare Hälfte dieser
Entscheidung — was `applied` bedeutet — ist genau der ergänzte Fließtext.

**8. Handbuch — bestätigt, kein neuer Zug.** Der Diff führt keine
Betreiber-Oberfläche ein: kein `CDC_*`-Zusatz im Container-Vertrag
(`compose.yaml` nicht im Diff), keine neue `cdc.*`-SQL-Funktion
(`tools/schema/nacharbeit-administration.sql` nicht im Diff), keine Horch-Adresse,
kein Endpunkt. Die seit `slice-077` verkörperte HIGH-Regel
(`.harness/skills/reviewer.md` §Klassifikation) ist damit **nicht** ausgelöst —
sie hängt am Diff, der die Oberfläche einführt. Ebenso ist die neue maschinelle
Verweisform-Regel (`.d-check.yml:165-175`) nicht berührt: der Diff fügt keinen
Link auf einen Slice-Plan mit Lifecycle-Verzeichnis hinzu, und `make docs-check`
läuft grün (582 Dateien, 0 Befunde). Die **vorbestehende** Lücke — `cdc.exclude_column`/
`cdc.include_column` kommen im Handbuch nicht vor — bleibt bestehen und ist
**kein** Finding dieses Laufs: sie ist der Beleg `evidence/slice-066.md` des
bereits bei 3× verkörperten Eintrags
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` mit
Remediation-Träger `slice-077`; hier ein zweiter Beleg wäre doppelte Zählung
derselben Klasse für eine Oberfläche, die der Diff nicht anfasst.

**9. Kommentar-Disziplin (`AGENTS.md` §3.7) — ohne Befund.** Grep über alle
`+`-Zeilen des Diffs auf `slice-`/`welle-` in Go-Dateien: kein Treffer in einem
Produktionscode-Kommentar (die Slice-Nennungen des Diff stehen in
Test-Godocs, Plan-Doku und Skript-Kommentaren, nicht über einem
Produktionscode-Pfad). Kein Kommentar beschreibt die verworfene Alternative im
Konjunktiv und keiner einen abwesenden Text: die neuen Blöcke nennen Zusage
(„ein Stand ohne Ausschluss ist ein fehlender Eintrag, keine leere Liste"),
Kopplung (`requested_at` als Transaktionszeit — deshalb der Zweitschlüssel) und
Herkunft (`ADR-0065`, `ADR-0059`, `ADR-0050`, `SPEC-008`).

**10. §2 DoD-Häkchen — korrekt getrennt.** Gesetzt sind die vier
Liefer-Punkte plus `make gates` plus das Doku-Update (Punkte 1–4, 6); offen
sind Review (5), Closure-Notiz (7), Reconciliation (8, mit Entfall-Begründung),
Beobachtungs-Register (9), Risiko-Ausgänge (10) und die drei Paarungen (11) —
genau die Planner-Arbeit nach diesem Report. Die sechs gesetzten Häkchen halten
der Prüfung stand: der Store-Test, die beiden Bootstrap-Tests und die drei
Bootstrap-Testnamen existieren wie benannt, `make test`/`make test-store`/
`make gates` habe ich selbst mit Exit 0 gesehen, die E2E-Zeile
*Spaltenausschluss-Neustart-Beleg* steht im Skript.

---

## Negativbefunde

- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/` (neue
  Abfrage + Auswertung; kein Schreibpfad, kein Schema-Objekt)
- geprüft, ohne Befund: `internal/adapters/driving/replication/mapper/` (nicht
  im Diff — der In-Prozess-Schutz aus `ADR-0065` Festlegung 4 ist unangetastet)
- geprüft, ohne Befund: `internal/bootstrap/` (beide Aufrufer, beide
  Fehlerpfade; die neuen Tests binden die Zusagen — vier eigene Mutationen)
- geprüft, ohne Befund: `internal/application/port/outbound/` und die beiden
  Use-Case-Fakes (Fähigkeit am bestehenden Zuschnitt, kein zweiter Port)
- geprüft, ohne Befund: `spec/pflichtenheft.md` (kein Spec-Stratum-Verstoß —
  präzisiert, erweitert nicht)
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (Position des
  Belegs, drei Zusicherungen, kein geteilter Zählraum mit anderen Abschnitten —
  alle übrigen Belege zählen je Tabelle **und** je `id`)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` (keine neue
  Betreiber-Oberfläche im Diff; die vorbestehende Lücke trägt `slice-077`, s.
  Urteil 8)
- geprüft, ohne Befund: `harness/image-hash.txt` und `.d-check.yml`
  (Referenzform-Regel eingehalten, `make docs-check` grün)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Herkunfts-Vermerk am Beleg-Baustein der
Werkzeuge-Zeile fehlt" · „Teil-Aktivierung bei Fehler nach dem Enable-Schritt"
· „Lauf-Beleg ohne abgelegtes Artefakt (Behauptung trägt unabhängig)"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Kein Finding berührt die
Richtung des Slice: der Stand wird real aus den `applied`-Zeilen abgeleitet,
ohne zweite Quelle und ohne neues Schema-Objekt (Urteil 1), der In-Prozess-Schutz
bleibt unangetastet (Urteil 2), der Zweitschlüssel ist total geordnet (Urteil 3),
der reale Neustart-Beleg steht an der richtigen Stelle und bindet drei Aussagen
(Urteil 4). Vier eigene Mutationen sind einzeln rot gesehen worden (Urteil 5);
die vier Zusagen des Implementers tragen an genau den Tests, die er nennt.

**Übergabe:** Keine Rückgabe an den Implementer — keine Fixrunde. F-1 ist ein
Ein-Zeilen-Nachzug im `harness/README.md` (Herkunfts-Vermerk), F-3 ist eine
Text-Klarstellung in §3 des Slice-Plans („siehe Bericht" → Adresse); beide sind
textlich, isoliert und ohne Verhaltenswirkung, für sie wird die Rückkante
Review → Implementer nicht geöffnet (Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-08-agentenrollen.md`: „bei isolierten LOW/INFO-Findings ist die
Sequenz Overkill"). F-2 ist eine benannte Grenze ohne Handlungsbedarf am Diff.
Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den
Zähler.

**DoD-Nachzug:** Da keine Fixrunde folgt, zieht dieser Report die Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 des Slice-Plans
selbst auf `[x]` nach, im selben Commit wie dieser Report
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Nur diese
eine Zeile; Closure-Notiz, Risiko-Ausgänge, Beobachtungs-Register und die drei
Paarungen bleiben Planner-Arbeit. Der Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11). Der reale
`make test-integration`-Lauf (DoD-Punkt 3) wurde von diesem Review **nicht**
nachgefahren; die Sensor-Exit-Codes dieses Laufs stehen unten.

---

## Ausgeführte Sensoren (Exit-Codes dieses Laufs)

| Aufruf | Exit-Code | Beleg |
|---|---|---|
| `make gates` | 0 | `docs-check` 582 Dateien / 0 Befunde · `commit-traceability` OK · `a-check` 0 Befunde · `coverage-gate` 49,30 % ≥ 35 % |
| `make test` | 0 | alle Pakete `ok`, `-race` im gepinnten Toolchain-Container |
| `make test-store` | 0 | u. a. `postgresstorage` 5,3 s `ok` (reale PostgreSQL) |
| vier Mutationen (eigene) | 2 / 2 / 2 / 2 | s. Urteil 5; danach zurückgenommen und Blatt-geprüft |

Hinweis: `make test-store` schreibt `tools/schema/plan.yaml` (Rollout-Report)
neu — der Schreibvorgang ist Lauf-Artefakt des Sensors und nicht Teil dieses
Reports.
