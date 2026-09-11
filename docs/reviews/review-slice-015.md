# Review-Report: slice-015 — 2026-09-11

**Review-Art:** Code — geprüft gegen Slice-Plan + [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) (Maintainability).

**Gegenstand:** Commit `3cb0c8e`
(`feat(schema): d-migrate-Pin auf v1.3.0, CHECK-Ausweichform zurückgebaut (ADR-0043)`),
einziger Commit dieses
Laufs — `Makefile`, `tools/schema/schema.yaml`,
`tools/schema/nacharbeit-operation-check.sql` (gelöscht),
`tools/schema/nacharbeit-views.sql`, `tools/schema/plan.yaml`,
`tools/schema/down.sql`, sowie der Slice-Plan selbst
(`docs/plan/planning/in-progress/slice-015-d-migrate-1.3.0-retirement.md`).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-11

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-015-d-migrate-1.3.0-retirement.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4–6 Trigger/Risiken, §8
  Sub-Area-Prüfung)
- [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  (Schemamigrationen mit d-migrate; Re-Evaluierungs-Trigger: kein
  Werkzeugwechsel nötig, solange eine Ausweichform existiert)
- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`
  (Register-Stand: 2× — CHECK slice-006, Views slice-010; unter der
  3×-Schwelle)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.3 Move/Inhalt-Trennung,
  §3.7 Kommentar-Klassen) · `harness/conventions.md` (MR-000, genau eine
  Sub-Area `PGC`, Greenfield)
- Commit-Traceability (`AGENTS.md` §5, `ADR-0045`)

---

## Findings

### F-1 — §1 (Ziel/Ausschluss) und der unveränderte Teil von DoD-Punkt 2 tragen weiterhin die widerlegte „Sandbox-Modus"-Prämisse

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Plan-interne Konsistenz, AGENTS.md §3.7 sinngemäß
  auf Plan-Prosa übertragen — dort ist die Regel wörtlich auf
  Code/Konfiguration/Skripte und Zustandsfelder begrenzt, der Fall hier ist
  eine dritte, nicht wörtlich gelistete Instanz derselben Sorge: ein Text
  behauptet einen Zustand, den derselbe Diff im selben Absatz widerlegt)
- `pfad`:
  `docs/plan/planning/in-progress/slice-015-d-migrate-1.3.0-retirement.md:41-44`
  (§1 Ziel: „…einen neuen „Raw SQL Sandbox Mode" mit Provenance-Overlay,
  die genau die `raw-sql-text-drift`-Klasse aus `BEO-PGC/d-migrate-nacharbeit`
  (2×…) adressieren") ·
  `docs/plan/planning/in-progress/slice-015-d-migrate-1.3.0-retirement.md:55-59`
  (§1 Ausschluss-Begründung für das Provenance-Overlay: „der Sandbox-Modus
  allein genügt, um die Drift-Klasse zu testen") ·
  `docs/plan/planning/in-progress/slice-015-d-migrate-1.3.0-retirement.md:85-93`
  (DoD-Punkt 2, erster — unveränderter — Satz: „real gegen den Sandbox-Modus
  getestet", fünf Zeilen später im selben Punkt, im **neu hinzugefügten**
  Beleg-Text: „kein separates „Raw SQL Sandbox Mode"-Flag existiert")
- `befund`: Der Diff fügt in DoD-Punkt 2 den Befund an, dass die
  Changelog-Prämisse eines eigenständigen „Raw SQL Sandbox Mode" nicht
  zutrifft (`--help` erschöpfend geprüft) — korrigiert aber weder den
  unmittelbar vorausgehenden, unveränderten Satz desselben Punkts noch §1
  (Ziel und Ausschluss-Begründung), die beide dieselbe Prämisse als Tatsache
  behaupten bzw. eine Scope-Entscheidung darauf stützen. Ein Leser, der nur
  §1 liest, hält die Overlay-Ausgrenzung weiterhin für „der Sandbox-Modus
  allein genügt" begründet — ein Mechanismus, dessen Existenz derselbe Slice
  im selben Commit widerlegt. Die Stelle wurde bearbeitet (die
  Beleg-Ergänzung liegt direkt daneben), die Korrektur der widersprechenden
  Nachbarsätze blieb aus.
- `verifizierbar`: nein — kein Gate prüft Plan-Prosa auf interne Konsistenz.
- `klasse`: Plan-Prämisse durch eigene Beleglage widerlegt, Text
  unkorrigiert stehen gelassen (1. Auftreten)

### F-2 — `schema.yaml`-Kopfkommentar behauptet weiterhin eine gemeinsame „benannte Grenze" von Views und `chk_change_operation`, die der Diff selbst auflöst

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/schema/schema.yaml:27-28` („Die drei Views tragen dieselbe
  benannte Grenze wie `chk_change_operation` unten (roher SQL-Text,
  `raw-sql-text-drift`) und leben deshalb nicht im `views:`-Knoten dieser
  Datei…") — außerhalb des in diesem Commit geänderten Hunks (dieser liegt
  bei Zeile 38 ff.), aber durch ihn stale geworden
- `befund`: Der Diff ersetzt den weiter unten stehenden Absatz „Benannte
  Grenze — der CHECK über der Operation" durch „Retired (slice-015)" —
  `chk_change_operation` ist ab diesem Commit keine „benannte Grenze" mehr,
  sondern deklarativ konvergent. Der frühere, unveränderte Absatz bei Zeile
  27-28 sagt weiterhin, die Views trügen „dieselbe benannte Grenze wie
  `chk_change_operation` unten" — ein Vergleichsziel, das es an dieser
  Stelle der Datei nicht mehr gibt. Die Aussage bleibt für sich lesbar
  (Views sind weiterhin Ausweichform), der Vergleich ist aber falsch.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Querverweise innerhalb
  einer YAML-Datei.
- `klasse`: Stale Querverweis nach Teiländerung am Zielabsatz (1. Auftreten)

---

## Negativbefunde

- geprüft, ohne Befund: **[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Konformität** — der Re-Evaluierungs-
  Trigger („d-migrate kann eine Operation nicht ausdrücken **und** es gibt
  keine Ausweichform") feuert nicht: für `chk_change_operation` entfällt die
  Ausweichform, weil die Operation jetzt ausdrückbar ist (im Sinn der ADR
  positiv aufgelöst); für die drei Views existiert weiterhin eine
  dokumentierte Ausweichform (`nacharbeit-views.sql`). Kein Werkzeugwechsel
  fällig, ADR bleibt `permanent` ohne Folge-ADR-Pflicht.
- geprüft, ohne Befund: **Pin-Änderung minimal und korrekt** —
  `Makefile:73` ändert ausschließlich den `D_MIGRATE_IMAGE`-Digest
  (64-Hex-`sha256`, formal valide); keine Nebenänderung an
  `SCHEMA_SOURCE`/`SCHEMA_TARGET`/`SCHEMA_ROLLOUT_NETWORK`.
- geprüft, ohne Befund: **`schema.yaml`-Deklarativität** — der neue
  `chk_change_operation`-Eintrag liegt sauber unter `change.constraints`
  (vor `chk_change_sequence_positive`), Form identisch zu den bestehenden
  Constraint-Einträgen (`name`/`type`/`expression`); keine
  Nebeneffekt-Änderung an anderen Tabellen/Spalten.
- geprüft, ohne Befund: **§3 Plan-Nachzug vollständig** — jede im Diff
  geänderte Datei (`Makefile`, `schema.yaml`,
  `nacharbeit-operation-check.sql`, `nacharbeit-views.sql`, `plan.yaml`,
  `down.sql`) ist in §3 der Slice-Plan-Datei mit Änderungsart und
  Begründung gelistet; `plan.yaml`/`down.sql` sind korrekt als „Plan-
  Nachzug … Nebenprodukt des Sensor-Laufs" gekennzeichnet, nicht als
  separat verfasster Inhalt.
- geprüft, ohne Befund: **Hard Rule 3.3 (git mv + Inhalt = zwei Commits)**
  — es findet kein `git mv` statt (der Slice bleibt in `in-progress/`, kein
  Lifecycle-Übergang); `nacharbeit-operation-check.sql` wird gelöscht, nicht
  verschoben. Ein einziger Commit mit Inhaltsänderung ist damit korrekt,
  die Hard Rule ist nicht einschlägig.
- geprüft, ohne Befund: **Kommentar-Klassen in `Makefile` und
  `nacharbeit-views.sql`** — die aktualisierten Blöcke sind indikativ,
  nennen den geltenden Zustand mit Rang-Zeiger (`ADR-0043`, `slice-015`,
  Beleg „real gegen einen frischen Rollout getestet"); kein Konjunktiv über
  eine verworfene Alternative, kein abgebrochener Satz. Der neue „Retired
  (slice-015)"-Block in `schema.yaml:38-46` erklärt die frühere Ausweichform
  im Dienst der aktuellen Zusage („konvergiert deklarativ … kein E5/E012
  mehr") und ist damit selbst kein Verstoß gegen §3.7 — anders als die
  stehengebliebene Nachbarstelle in F-2.
- geprüft, ohne Befund: **Docker-only** — alle Läufe (`schema validate`,
  `schema migrate`, psql-Nacharbeit) laufen über `docker run` im Makefile,
  kein lokales Toolchain-Install im Diff.
- geprüft, ohne Befund: **Suppression-Verbot** — kein `nolint`/`noqa`/
  `SuppressMessage`-Marker im Diff.
- geprüft, ohne Befund: **Traceability der Commit-Message** — trägt
  `ADR-0043` im Betreff und im Body, keine `SPEC-*`/`ARC-*`-Kennung im
  Betreff.
- geprüft, ohne Befund: **dangling Referenzen auf die gelöschte Datei** —
  `nacharbeit-operation-check.sql` kommt nach dem Diff nur noch in
  historischen Artefakten vor (`done/welle-3-results.md`,
  `done/slice-010-…`, frühere Review-/Verify-Reports, das
  Beobachtungs-Register), keine in aktiv geltender Doku
  (`harness/README.md`, `AGENTS.md`) oder im `schema-rollout`-Target
  selbst — die Bindungszeile in `harness/README.md` nennt ohnehin keine
  `nacharbeit-*.sql`-Dateinamen (DoD-Punkt „Doku-Update" korrekt als
  entfallend geprüft).
- geprüft, ohne Befund: **§6 Risiken vorab benannt** — beide Risiken
  (partielles Retirement; Regressions-Fund durch Pin-Bump) waren bereits
  vor der Implementierung im Plan benannt und decken exakt den
  eingetretenen Fall (partiell: CHECK gelöst, Views nicht); die
  Ausgangs-Zuweisung selbst ist Planner-Aufgabe bei Closure, hier nicht
  bewertet.
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung** — GF-Sub-Area `PGC`
  korrekt referenziert, Register-Sichtung mit Stand zum Planungszeitpunkt
  dokumentiert, kein Eintrag erreicht mit diesem Slice die 3×-Schwelle.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Kein HIGH-Finding.** Insbesondere: kein ADR-Verstoß, kein
Traceability-/ID-Schema-Verstoß, keine Suppression, kein Docker-only-
Verstoß, keine Zwei-Quellen-Drift zwischen zwei Dateien (F-1 ist
plan-intern, zwischen zwei Abschnitten derselben Datei — das erfüllt die
HIGH-Definition „zwei Dateien" nicht wörtlich und wird deshalb als MEDIUM
geführt, nicht als HIGH-Analogieschluss).

**Finding-Klassen dieses Laufs:** Plan-Prämisse durch eigene Beleglage
widerlegt, Text unkorrigiert stehen gelassen (1. Auftreten) · Stale
Querverweis nach Teiländerung am Zielabsatz (1. Auftreten)

## Verdikt

**Merge-blockierend:** nein — kein HIGH-Finding, kein Rollen-Widerspruch.
F-1 (MEDIUM) und F-2 (LOW) sind reguläre Findings ohne Konflikt zwischen
Rollen; der Implementer akzeptiert oder begründet sie (Modul 8 §Konflikt-
Pfad: Sequenz-Modellierung wäre hier Overkill, nicht ausgelöst). Der
Closure-Trigger des Slice-Plans (§5: „Review-Schluss ohne offenes
HIGH-Finding") ist erfüllt.

**Zur eigentlichen technischen Frage (Pin-Bump, CHECK-Retirement,
Views-Grenze):** unproblematisch. Der CHECK-Ausdruck ist sauber
deklarativ überführt, real gegen einen frischen Rollout belegt (`plan.yaml`/
`down.sql` zeigen den Constraint jetzt im generierten `CREATE TABLE
"change"`); die Views-Grenze bleibt korrekt und mit aktualisiertem Beleg
dokumentiert bestehen. Die Diskrepanz zur ursprünglichen Changelog-Prämisse
(„Sandbox Mode") ist selbst kein Code- oder ADR-Problem — sie ist eine
Planungs-Prosa-Frage (F-1), die die tatsächliche technische Deckung des
DoD-Punkts nicht schmälert, aber künftige Leser des Plans fehlleiten kann.

**Übergabe:** F-1 und F-2 gehen an den Implementer/Planner zur
redaktionellen Korrektur von §1 und der Kopfkommentare — beide sind ohne
Rollen-Widerspruch lösbar (Text angleichen an die bereits vorliegende
Beleglage, keine neue Entscheidung nötig). Die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser
Report selbst ist ein **Lauf-Beleg** und wird über Läufe hinweg nicht
wieder gelesen. Der Report ersetzt keine Verifikation — DoD-/
Spec-Konformität (inkl. ob `make gates`/`make schema-rollout` tatsächlich
grün liefen) prüft der Verifier separat (Modul 11).

---

**Gate-Beleg:** in dieser Review-Sitzung nicht erneut ausgeführt (Docker-
/Netzzugriff außerhalb des Scopes eines Maintainability-Reviews); die vom
Implementer berichteten Läufe (`make schema-validate`, `make
test-integration`, `make gates`) sind Gegenstand der Verifikation, nicht
dieses Reviews.
