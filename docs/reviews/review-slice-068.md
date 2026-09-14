# Review-Report: slice-068 — 2026-09-14

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Modul 10
§Drei Review-Arten); DoD-/Spec-Konformität ist Verifier-Aufgabe und nicht
Gegenstand dieses Reports.

**Gegenstand:** [`slice-068`](../plan/planning/in-progress/slice-068-e2e-spaltenausschluss.md),
Diff `d088bb0..7fa5784` (Elter-Commit `d088bb0` ist ein reiner
`next→in-progress`-Move und trägt keinen Inhalt). Ein Commit `7fa5784`,
vier Dateien (+262/−8): `tools/harness/run-integration-tests.sh`,
`harness/README.md`, `tools/schema/plan.yaml`, der Slice-Plan selbst.

**Skill:** `.harness/skills/reviewer.md` @ `7fa5784` (Stand zum
Review-Zeitpunkt, unverändert seit der letzten Schärfung 2026-09-13).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-14.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-068`](../plan/planning/in-progress/slice-068-e2e-spaltenausschluss.md)
  (vollständig, inkl. der vier Plan-Nachzüge und §6)
- [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (vollständig,
  insbesondere Teilfrage 3 Option D, Teilfrage 5, §Konsequenzen,
  §Fitness Function)
- [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  (Abgrenzung: die Dauerhaftigkeit des Ausschlussstandes ist **nicht**
  Gegenstand dieses Slice; ihre Kontext- und Alternativen-Passagen tragen für
  die Auslegungsfrage unten)
- [`LH-FA-CFG-005`](../../spec/lastenheft.md) (Happy Path / Boundary /
  Negative wörtlich), [`LH-QA-SEC-004`](../../spec/lastenheft.md),
  [`LH-FA-DAT-005`](../../spec/lastenheft.md),
  [`LH-FA-SCH-003`](../../spec/lastenheft.md)
- [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  (Antrags-Queue), [`ADR-0030`](../plan/adr/0030-testpyramide.md)
  (E2E-Tier), [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
- `AGENTS.md` §3.7/§3.9/§5, §6 (Rollenwechsel); `harness/conventions.md`
  (`MR-000` ID-Schema)
- Vorgänger `slice-066`/`slice-067` (in `done/`) samt
  `review-slice-066.md`/`review-slice-067.md`; `welle-18` §3/§6
- `docs/plan/planning/observations/BEO-PGC/` — `plan-nachzug`,
  `plan-vorlagen-defekt`, `dod-checkbox-nachzug-review-ohne-fixrunde`,
  `slice-chronik-in-code-kommentar`, `report-nackte-id-ohne-link`,
  `test-isolation-geteilter-zustand`, `test-runner-stiller-ausschluss`,
  `test-integration-retention-timing-flake`

**Eigene Sensor-Läufe** (Exit-Code jeweils ungepiped in einem eigenen Schritt
ermittelt, `AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify, docs-check (530 Dateien, 0 Befunde), commit-traceability (5 Commits), a-check (0 Befunde; Hinweis: 2 Dateien in keiner Schicht, unverändert), coverage-gate (45.80 % über Schwelle 35 %) |
| `make test-integration` | **0** | voller Compose-Stack-Lauf; beide neuen Beleg-Zeilen in der Ausgabe (`Spaltenausschluss-Rundlauf (LH-FA-CFG-005 Happy Path, ADR-0059) belegt` und `Spaltenausschluss-Negative-Beleg (LH-FA-CFG-005)`) |

**Eigene Mutationsläufe** (unabhängige Nachprüfung der beiden vom Implementer
gemeldeten roten Mutationsbefunde). Jede Mutation erzeugt einen roten
`make test-integration`-Lauf, danach Rücknahme per `git checkout` — der
Arbeitsbaum ist nach jeder Prüfung leer (`git status --porcelain`), und der
Blob-Hash der zurückgenommenen Datei stimmt mit `HEAD` überein
(`git hash-object` = `git rev-parse HEAD:<pfad>` = `cc3f399…`):

| # | Mutation | Erwartung | Ergebnis |
|---|---|---|---|
| A | Antragsart verfälscht: `cdc.exclude_column(...)` → `cdc.include_column(...)` in der Happy-Path-Anforderung (`tools/harness/run-integration-tests.sh:348-349`) | rot an der Row-Image-Zusage | **rot** (Exit 2; genau `…die nach dem Ausschluss erfasste Change (id=2, …) trägt den ausgeschlossenen Spaltenschlüssel secret weiterhin im Row Image`) |
| B | Katalog-Blick ersetzt die nicht existierende Spalte: `'nicht_vorhandene_spalte'` → `'name'` in der Negative-Anforderung (`tools/harness/run-integration-tests.sh:443`) | rot an der Negative-Zusage | **rot** (Exit 2; genau `… endete mit status=applied, wollen failed (LH-FA-CFG-005 Negative: expliziter Fehlerpfad)`) |

Beide Angaben des Implementers sind damit unabhängig bestätigt: die zwei
Zusagen des Abschnitts sind einzeln tragend, keine ist Dekoration.

---

## Findings

### F-1 — Der §1-Satz „künftige **und historische** Changes … den ausgeschlossenen Wert nicht mehr tragen" überzeichnet die zitierte Anforderung

- `kategorie`: MEDIUM
- `quelle`: [`LH-FA-CFG-005`](../../spec/lastenheft.md) Happy Path („dann sind
  **künftige** Changes von `t` ohne die für CDC relevanten Datenwerte von
  `c`") · [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  §Teilfrage 3 Option D („der sensible Wert wird **nie** serialisiert und nie
  an die Persistenzschicht übergeben") · [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  §Verglichene Alternativen, Option A („Ein bereits persistierter Wert ist
  nicht zurückholbar; genau den schließt [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  Teilfrage 3 Option D als Wirkort aus")
- `pfad`: `docs/plan/planning/in-progress/slice-068-e2e-spaltenausschluss.md:43-45`
  (§1 *Ziel*)
- `befund`: Der Satz verlangt für die **vor** dem Ausschluss erfasste Change
  einen Wertverlust, den der im selben Satz als Bezug geführte Wirkort
  (Row-Image-Konstruktion) strukturell nicht herstellen kann — die zitierte
  Anforderung selbst fordert ausschließlich „künftige Changes von `t`". Der
  Diff belegt stattdessen die *unveränderte Lesbarkeit* der davor erfassten
  Change (`tools/harness/run-integration-tests.sh:419-437`) und folgt damit
  der Anforderung; der Plan-Text ist die abweichende Stelle.
- `verifizierbar`: nein — kein Gate deckt die Aussage ab; der Widerspruch ist
  am Verhältnis Plan-Satz ↔ Anforderungs-Wortlaut entschieden, nicht am Code
- `klasse`: „Plan-Text überzeichnet die zitierte Anforderung"

### F-2 — Der Kopf nennt „alle drei Akzeptanzkriterien", §1 schließt die Boundary aus

- `kategorie`: LOW
- `quelle`: [`LH-FA-CFG-005`](../../spec/lastenheft.md) Boundary · Plan
  `slice-068` §1 (Out-of-Scope-Disziplin)
- `pfad`: `docs/plan/planning/in-progress/slice-068-e2e-spaltenausschluss.md:12-13`
  gegen `:51-56`
- `befund`: Der `Bezug:`-Kopf führt „alle drei Akzeptanzkriterien: Happy Path,
  Boundary, Negative" als Bezug dieses Slice, während §1 den Boundary-Fall
  ausdrücklich ausschließt (Delegation an `slice-067`s Konvergenz-Test). Kein
  Deckungsverlust — `TestConsumeExcludedColumnDroppedInSourceReportsSchemaError`
  (`internal/adapters/driving/replication/mapper/mapper_test.go:629`) übt den
  Pfad real aus —, aber der Kopf liest sich als Zusage an diesen Slice.
- `verifizierbar`: ja — der genannte Testfall läuft in `make test`; die
  Aussage selbst ist eine Text-Konsistenz und kein Gate-Gegenstand
- `klasse`: „Kopf-Zusage deckt die §1-Abgrenzung nicht"

### F-3 — Der E2E-Beleg deckt nur den INSERT-Pfad; der `old_data`-Pfad ist Unit-Test-Beleg

- `kategorie`: INFO
- `quelle`: [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  §Fitness Function („erscheint nie als Schlüssel im resultierenden
  `old_data`/`new_data`-JSON eines `Change`")
- `pfad`: `tools/harness/run-integration-tests.sh:395-418`
- `befund`: Die drei Zusagen der nach dem Ausschluss erfassten Change prüfen
  den Schlüssel nur in `new_data` (`jsonb_exists(new_data, …)`); der
  `old_data`-Pfad eines UPDATE tritt in diesem Rundlauf nicht auf und ist
  ausschließlich über `TestConsumeExcludedColumnAbsentFromRowImages`
  (`mapper_test.go:540`, INSERT/UPDATE/DELETE inkl. `OldImage`) belegt.
- `verifizierbar`: ja — der genannte Testfall läuft in `make test`
- `klasse`: „E2E-Tier deckt nur einen Operations-Pfad"

### F-4 — `tools/schema/plan.yaml` reist mit, ohne Zeile in der §3-Tabelle

- `kategorie`: INFO
- `quelle`: Maintainability · Plan `slice-068` §3 (Pfad-Kandidaten-Quelle für
  §8)
- `pfad`: `tools/schema/plan.yaml:5`
- `befund`: Die Datei trägt in diesem Commit eine andere Ziel-DSN als zuvor
  (`postgres://postgres:***@cdc-test-postgres:5432/cdc`), weil der
  vorgeschriebene `make test-integration`-Lauf `make schema-rollout` ausführt;
  ihr Inhalt hängt damit vom zuletzt gelaufenen Rollout ab und wechselt
  zwischen Lauf-Umgebungen. Der Plan-Nachzug `:171-177` benennt sie, die
  §3-Tabelle `:126-129` führt sie nicht.
- `verifizierbar`: ja — ein `git status` nach `make test-integration` zeigt
  die Datei genau dann als geändert, wenn die Umgebung abweicht (in diesem
  Lauf blieb der Baum leer, die committete Fassung reproduziert sich lokal)
- `klasse`: „Generiertes Lauf-Artefakt wandert ohne §3-Zeile mit"

### F-5 — Die Sensors-Zeile benennt die tragende QA-Kennung des Leck-Checks nicht

- `kategorie`: INFO
- `quelle`: [`LH-QA-SEC-004`](../../spec/lastenheft.md) (Messmethode:
  „ausgeschlossene Spaltenwerte erscheinen nicht in den Changes") ·
  `harness/README.md` §Sensors
- `pfad`: `harness/README.md:131`
- `befund`: Der neue Satz der Sensors-Zeile nennt
  [`LH-FA-CFG-005`](../../spec/lastenheft.md) und
  [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md), nicht aber
  [`LH-QA-SEC-004`](../../spec/lastenheft.md), deren Messmethode die dritte
  Zusage des Abschnitts (Wert kommt im persistierten Change nicht vor,
  `run-integration-tests.sh:408-418`) trägt; der Skript-Kommentar `:399`
  nennt sie.
- `verifizierbar`: nein — Prosa-Vollständigkeit, kein Gate-Gegenstand
- `klasse`: „Beleg nennt die zuständige QA-Kennung nicht"

## Negativbefunde

- **geprüft, ohne Befund: Punkt 1 (Lage vor der Container-Ende-Grenze und
  Weiterlaufen des Feed-Containers) — am Skript-Code, nicht am Log.** Der
  neue Abschnitt endet bei `run-integration-tests.sh:481`; die Grenze liegt bei
  `:1765`/`:1868` (`TestE2ESchemaChangeDropColumn`, danach
  `TestE2ESchemaChangeIncompatibleTypeChange` bei `:1963`). Zwischen beiden
  steht kein `docker restart`/`docker pause`/`up`-Aufruf für den Feed-Container
  — die zwei `docker inspect --format '{{.State.Running}}'`-Wächter bei `:431-435`
  und `:476-480` sind reine Beobachtungen. Der Beleg-Rundlauf ist eine
  Shell-Strecke und berührt das Go-`-run`-Muster bei `:277-279` nicht; die
  bekannte Lücke `BEO-PGC/test-runner-stiller-ausschluss` (offen, 1×) wird
  dadurch weder berührt noch verschärft.
- **geprüft, ohne Befund: Punkt 2 (Baseline trennt „Wert abwesend" von
  „Spalte war nie im Row Image").** `run-integration-tests.sh:327-346` legt vor
  dem Ausschluss eine eigene Change an und verlangt real
  `new_data->>'secret' = 'ColumnBeforeExclusionSentinel'`. Damit ist genau die
  Alternativerklärung ausgeschlossen, die den ganzen Happy-Path-Beleg sonst
  entwerten würde: eine publication-/decoder-seitige Spaltenauswahl, die
  `secret` schon vor dem Antrag nie ins Row Image ließe. Der Wechsel id=1 → id=2
  ist die einzige Variable zwischen den beiden Belegen, und Mutationslauf A
  zeigt, dass die Zusage am Ausschluss hängt und nicht am Dekodier-Pfad.
- **geprüft, ohne Befund: Punkt 3 (Kontrolle der nicht ausgeschlossenen
  Spalte).** `:396-402` verlangt `new_data->>'name' = 'ColumnAfter'` in
  derselben Change, deren Schlüssel-Abwesenheit geprüft wird; sie
  unterscheidet „Filter wirkt gezielt" von „Filter wirft alles weg" (ein
  leeres `new_data` würde hier rot). Kein Leerlauf-Assert: die `id`-Bedingung
  der `WHERE`-Klausel allein beweist die Nicht-Leerheit bereits, die
  `name`-Zusage benennt sie.
- **geprüft, ohne Befund: Punkt 4 (Negative-Beleg).** `:439-481` fragt
  `cdc.exclude_column('src-e2e','public','feed_e2e_column_exclusion','nicht_vorhandene_spalte')`
  gegen eine Tabelle mit laufender Bindung ab, pollt auf `failed` (ein
  `applied` bricht die Schleife und führt in die rote Meldung, `:463-469`)
  und prüft `error_message` auf den Text des Sentinels
  `ErrSourceColumnMissing` (`internal/application/port/inbound/verwaltung.go:26`).
  Der Text-Abgleich imitiert die etablierte Form des Skripts für
  CLI-Ausgaben (`:990-1000`) — Kopplung an eine Meldung, kein zweiter
  Zustandsträger.
- **geprüft, ohne Befund: Punkt 5 (die zwei Mutationsbefunde des
  Implementers).** Beide selbst nachgestellt, beide rot, beide exakt an der
  jeweils zugesagten Stelle (Tabelle oben), beide ohne Restspur im Baum
  (`git status --porcelain` leer, Blob-Hash identisch).
- **geprüft, ohne Befund: Punkt 6 (die §3-Nachzüge — im Scope).**
  *Platzierung* (`slice-068:131-141`): der Grund hält am Code — die neue
  Tabelle wird erst nach dem Go-Tier aktiviert (`:298-299`), kein Go-Testfall
  liest oder schreibt sie (alle `cdc.changes`-Abfragen des Skripts filtern per
  `table_name`, geprüft über den ganzen Bestand), und der Abschnitt bleibt vor
  der Grenze. *Zusatz-Beobachtungen* (`:162-169`): Baseline und Kontrolle sind
  keine vierter Liefer-Punkt, sondern die zwei Alternativerklärungs-Ausschlüsse
  desselben Happy-Path-Belegs; die DoD-Zeile bleibt bei drei Liefer-Punkten.
  *`tools/schema/plan.yaml`*: kein Liefer-Punkt, sondern ein Erzeugnis des
  vorgeschriebenen Sensor-Laufs (siehe F-4, INFO). Kein Scope-Creep.
  Keine Zeile der §3-Tabelle wurde stillschweigend erweitert — die vier
  Nachzüge stehen als eigene Absätze vor §4 (`BEO-PGC/plan-nachzug`,
  `verkörpert`).
- **geprüft, ohne Befund: Punkt 7 (Kommentar-Disziplin, per `grep`).** Muster
  über die **hinzugefügten** Kommentarzeilen des Skripts (`:280-288`,
  `:327-329`, `:395-399`, `:419-423`, `:439-441`):
  `slice-[0-9]+|welle-[0-9]+|bislang|früher|zuletzt|nicht mehr|vorher|ehemals|hätte|wäre`
  → in Kommentaren **kein** Treffer der Klasse
  `BEO-PGC/slice-chronik-in-code-kommentar`; alle fünf Blöcke nennen den
  Ist-Zustand und tragen einen `ADR-*`/`LH-*`- oder Testfall-Zeiger,
  `:286-288` eine Kopplung (Grenze + Funktionen). Der Satz
  `harness/README.md:131` trägt den Herkunfts-Anker `· seit slice-068` in der
  Form der Nachbarzeilen, kein neues Gate (Bindungsspalte unverändert
  `kein Gate, ADR-0030, LH-QA-POR-003`).
- **geprüft, ohne Befund: Punkt 8 (`harness/README.md`-Update).** Der Satz
  beschreibt genau, was das Skript tut (aktivierte, nicht in `CDC_TABLES`
  gelistete Tabelle — `compose.yaml:89` führt sie nicht —, Poll auf `applied`,
  drei Aussagen zur danach erfassten Change, unveränderte Lesbarkeit der davor
  erfassten, Negative mit Fehlertext), und steht dort, wo die DoD-Zeile ihn
  verlangt (§Sensors, Zeile `make test-integration`).
- **geprüft, ohne Befund: Traceability/ID-Schema.** `7fa5784` trägt
  [`LH-FA-CFG-005`](../../spec/lastenheft.md) und
  [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) im Betreff, kein
  `SPEC-*`/`ARC-*`; `make commit-traceability` grün (Exit 0). Die im Diff
  vergebenen Kennungen folgen `MR-000` (keine neue Kennung).
- **geprüft, ohne Befund: Zustandsfeld-Chronik und Beobachtungs-Register.**
  Der Diff ändert kein `Stand`-/`Status`-Feld und kein Drift-Log; die
  §8-Sichtung des Plans (`:274-282`) zitiert die zwei Registereinträge
  (beide 1×, weiter offen) statt sie neu zu formulieren.
- **geprüft, ohne Befund: keine ungewollte Zustands-Überschneidung.** Die
  Aktivierung von `feed_e2e_column_exclusion` wirkt nur auf diese Tabelle; die
  danach laufenden Abschnitte (Lasttest, Black-Box-CLI, Retention-Lifecycle,
  Diagnose, Upgrade-Tausch) lesen ausschließlich per `table_name` gefiltert,
  und die `min(commit_position)`-Rechnung des Retention-Abschnitts (`:663`)
  wird durch zusätzliche junge Transaktionen nicht verändert.
  `BEO-PGC/test-isolation-geteilter-zustand` und
  `BEO-PGC/test-integration-retention-timing-flake` (beide 1×, weiter offen)
  werden von diesem Diff nicht berührt.
- **geprüft, ohne Befund: DoD-Häkchen (`7fa5784`).** Gesetzt sind genau die
  umsetzungs- und gate-tragenden Zeilen (drei Liefer-Punkte, `make gates` +
  `make test-integration`, Doku-Update); Review, Closure-Notiz,
  Reconciliation-/Beobachtungs-Register, Risiko-Ausgänge und die drei
  Paarungen bleiben offen.
- **geprüft, ohne Befund: `docs/reviews/`, `docs/plan/planning/observations/`,
  `harness/`, `tools/schema/`** — kein Diff in `spec/*`, `.github/workflows/`,
  `compose.yaml`, `internal/` oder `test/integration/`; der Rollenwechsel
  bleibt in seinem Slice (Modul 8).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Plan-Text überzeichnet die zitierte
Anforderung · Kopf-Zusage deckt die §1-Abgrenzung nicht · E2E-Tier deckt nur
einen Operations-Pfad · Generiertes Lauf-Artefakt wandert ohne §3-Zeile mit ·
Beleg nennt die zuständige QA-Kennung nicht

## Verdikt

**Zur vorgelegten Auslegungs-Abweichung — sie trägt.** Maßgeblich ist der
Wortlaut der Anforderung, nicht die Begründung des Implementers:
[`LH-FA-CFG-005`](../../spec/lastenheft.md) Happy Path verlangt „**künftige**
Changes von `t` ohne die für CDC relevanten Datenwerte von `c`" — die
Bestimmung „historisch" kommt im Lastenheft nicht vor. Die zweite Fundstelle
ist der Wirkort selbst: [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
§Teilfrage 3 Option D setzt die Filterung in die Row-Image-Konstruktion und
sagt zu, der Wert werde „nie serialisiert" — was *nach* der Filterung
entsteht, ist erfasst; was davor lag, ist davon unberührt. Dass ein bereits
persistierter Wert nicht zurückholbar ist, hält
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
§Verglichene Alternativen selbst fest und benennt genau diesen Wirkort als
den, der ihn ausschließt. Die Forderung des §1-Satzes ist damit unter der
geltenden Entscheidung **unerfüllbar**, und der Diff erfüllt statt ihrer die
Anforderung: die nach dem Ausschluss erfasste Change trägt den Schlüssel nicht
mehr und den Wert nirgends, die davor erfasste bleibt unverändert lesbar —
beide Hälften real belegt. Die Abweichung ist folglich **kein Finding gegen
den Code**: zu korrigieren ist der Plan-Text (F-1, Planner-Arbeit), nicht der
Abschnitt. `ADR-0059`s eigene Schnitt-Empfehlung
(`:248-254`) formuliert die historische Hälfte mit der eingeschobenen
Klammer „für den bereits vor Slice 2 gedeckte Fall" und ist als Verweis auf
den anderweitig gedeckten Fall lesbar; der §1-Satz des Plans lässt die Klammer
weg und macht daraus eine Behauptung. Da §Konsequenzen die Empfehlung
ausdrücklich als „Empfehlung an den Planner, keine Festlegung dieser ADR"
führt, braucht die Korrektur keine Folge-ADR — sie ist ein Planer-Nachzug.

**Merge-blockierend:** nein — begründete Abweichung vom Regelfall. Der Commit
liegt bereits auf `main` (`git branch --contains 7fa5784`), es gibt keinen
offenen PR, und das einzige MEDIUM-Finding ist keine Korrektur am Diff,
sondern eine Textstelle des Plans (Rückkante Review → Plan, Modul 8). Der
Code-Teil des Diff ist in allen geprüften Punkten ohne Befund: der neue
Abschnitt liegt vollständig vor der Container-Ende-Grenze, der Feed-Container
läuft nach beiden Belegen weiter, kein Neustart, beide Zusagen sind durch
eigene Mutationsläufe als tragend bestätigt, die Baseline trennt „Wert
abwesend" von „Spalte war nie im Row Image", die Kontrolle unterscheidet
gezielten von totalem Filter, der Negative-Beleg landet real `failed` samt
Sentinel-Text, kein §3-Nachzug ist Scope-Creep, `make gates` und
`make test-integration` laufen Exit 0.

**Übergabe:** keine Fixrunde am Implementer — kein Reviewer→Implementer-Pfeil.
F-1 und F-2 gehen an den **Planner** (Plan-Text); F-3 bis F-5 sind Hinweise
ohne erwartete Aktion. Da keine Rückgabe an den Implementer erfolgt, ist die
DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor" im
selben Commit, der diesen Report anlegt, auf `[x]` gezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde;
`BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde`). Dieser Report ist
Lauf-Beleg und wird über Läufe hinweg nicht erneut gelesen; Verifikation
gegen DoD/Spec bleibt Aufgabe des Verifiers (Modul 11; anderer
Eingabe-Kontext).
