# Verifier-Report: slice-016 — 2026-09-11

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code), §6 (Risiko-Vorschlag, Ausgang bleibt
Planner-Entscheidung) und Entscheidungs-Konformität gegen
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
(Re-Evaluierungs-Trigger). Nicht geprüft: Diff gegen Plan/Hard Rules im
Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe, bereits erledigt,
siehe [`review-slice-016.md`](review-slice-016.md)), realer Bedarf
(Validator).

**Gegenstand:** `b86cf70` (Implementierung: Pin-Bump + Views-Retirement),
`4367809` (Kommentar-Korrektur `schema.yaml`), `e582a40` (Review-Report,
F-1 HIGH), `04590a4` (Review-Fixrunde F-1).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren, inklusive zweier eigener,
frischer Testcontainer-Rollout-Läufe gegen den gepinnten
d-migrate-1.3.1-Digest (unabhängig von Implementer- und Reviewer-Lauf,
eigene PostgreSQL-18-Instanz, eigenes Docker-Netz) und einer eigenen
Gegenprobe (Views ohne `columns:`-Signatur). Alle selbst angelegten
Container/Netze/Temp-Dateien wurden nach dem Lauf entfernt;
`tools/schema/plan.yaml`/`down.sql` auf den committeten Stand
zurückgesetzt (`git checkout --`); `git status` danach sauber.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-11

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-016-d-migrate-1.3.1-views-retirement.md`, nach
  `04590a4` — Plan-Datei selbst von diesem Commit nicht mehr berührt)
- `review-slice-016.md` (F-1 HIGH, committet `e582a40`, disponiert
  `04590a4`, kein offenes HIGH)
- [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) im
  Volltext (Re-Evaluierungs-Trigger: kein Werkzeugwechsel, solange
  d-migrate die Operation ausdrücken kann oder eine Ausweichform besteht)
- `harness/conventions.md` (MR-000, Sub-Area `PGC`, Greenfield),
  `harness/README.md` (Sensors-Tabelle, `schema-rollout`-Bindungszeile)
- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`
  (Register-Stand: 3×, `verkörpert`, keine `evidence/slice-016.md`)
- Code im Volltext: `Makefile` (Zeilen 60–120), `tools/schema/schema.yaml`,
  `tools/schema/nacharbeit-roles.sql`,
  `tools/schema/nacharbeit-observability.sql`,
  `tools/schema/nacharbeit-heartbeat.sql`, `tools/schema/plan.yaml`,
  `tools/schema/down.sql`, `tools/schema/compose-init/01-cdc-schema.sql`

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make schema-validate` | zieht den gepinnten v1.3.1-Digest, `Validation passed: 0 warning(s)`, 8 Tabellen/29 Spalten/1 Index/7 Constraints | **0** |
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 169 Datei(en), 0 Befund(e)` (voll und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` · `a-check gesamt: 0 Befund(e)` | **0** |
| `make commit-traceability RANGE=c58531f..04590a4` | 4 Commits (Slice-Übergang bis Review-Fix), alle mit `ADR-0043`-Bezug, kein Struktur-ID-Betreff | **0** |
| Eigener frischer Rollout (eigene PostgreSQL-18-Instanz, `wal_level=logical`, Compose-Init-Skript, eigenes Docker-Netz): `make schema-rollout` gegen leere DB | 8× `CreateTable`, 3× `CreateView` (`active_tables`, `changes`, `consumer_status`), `primaryBlockedReason: null`; `\dv cdc.*` zeigt alle fünf Views (drei aus `schema.yaml` + `heartbeat`/`metrics` aus den psql-Nacharbeit-Skripten) | **0** |
| Zweiter `make schema-rollout`-Lauf gegen dieselbe, jetzt migrierte DB | `d-migrate schema migrate --execute` selbst: 2× `ReplaceView` (`active_tables`, `changes`), 1× `ReplaceView` (`consumer_status`) — **kein** `VIEW_SIGNATURE_UNKNOWN`, kein Signatur-Blocker; daneben 2× `DropView` (`cdc.heartbeat`, `cdc.metrics`, außerhalb des neutralen Modells) mit `primaryBlockedReason: DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION` | **8** (durch die `DropView`-Bestätigungspflicht auf `heartbeat`/`metrics` — siehe Einordnung unten) |
| Eigene Gegenprobe: `schema.yaml`-Kopie ohne `columns:`-Knoten in den drei Views (auf separater, frisch migrierter DB), `source_dialect: postgresql` unverändert | `schema migrate --execute` (isolierter `docker run`, kein `make`-Wrapper): alle drei `ReplaceView`-Operationen `skipped: true`, `primaryBlockedReason: MANUAL_ACTION_REQUIRED`, je Operation `VIEW_SIGNATURE_UNKNOWN` im Report | **8** (erwartet — Gegenprobe reproduziert den Blocker, den `columns:` verhindert) |

**Nicht selbst neu gebaut:** `make image` (kein Gate, in diesem Diff nicht
berührt).

## DoD-kritischer Punkt (2): eigenständig, unabhängig reproduziert

Beide Teilaussagen aus DoD-Punkt 2 sind mit eigenen, frischen Läufen
bestätigt:

1. **Positiv-Fall:** Erstanlage *und* Folgelauf (`ReplaceView`) der drei
   Views laufen mit `source_dialect: postgresql` + `columns:`-Signatur
   ohne `VIEW_SIGNATURE_UNKNOWN` — die eigentliche technische Zusage des
   Slice ist erfüllt.
2. **Gegenprobe:** Ohne `columns:`-Signatur blockiert derselbe
   `ReplaceView`-Pfad reproduzierbar mit `VIEW_SIGNATURE_UNKNOWN` — die
   Kausalität (`columns:` ist der wirksame Faktor, nicht `source_dialect`
   allein oder Zufall) ist damit nicht nur behauptet, sondern gezeigt.

**Eine Präzisierung der DoD-Formulierung selbst, kein Befund gegen die
Implementierung:** DoD-Punkt 2 schreibt wörtlich „`schema migrate
--execute` konvergiert (Exit 0) … auch gegen eine bereits migrierte DB
(Folgelauf, kein `ReplaceView`-Blocker)". Der **volle**
`make schema-rollout`-Lauf (vier Schritte: `d-migrate` +
drei psql-Nacharbeit-Dateien) exitet im Folgelauf real mit **8**, nicht 0
— aber nicht wegen der drei in diesem Slice überführten Views, sondern
wegen `DropView` auf `cdc.heartbeat`/`cdc.metrics` (Rollout-Nebenprodukt
außerhalb von `schema.yaml`, unverändert seit slice-011/slice-012). Liest
man die Klammer „(Folgelauf, kein `ReplaceView`-Blocker)" als die
eigentliche Erfüllungsbedingung des zweiten Halbsatzes — wofür die
Formulierung selbst spricht, ebenso wie die Review-Interpretation —, ist
der DoD-Punkt materiell erfüllt. Liest man „Exit 0" wörtlich für **beide**
Szenarien inklusive aller vier Rollout-Schritte, ist er es nicht. Diese
Zweideutigkeit stand schon vor der Implementierung im Plan-Text und ist
kein Implementierungs-Defekt; sie gehört vor `git mv` nach `done/`
sprachlich geschärft oder im §7-Lerneintrag benannt, damit künftige
DoD-Formulierungen „Exit 0" und „kein X-Blocker" nicht in einem Satz
vermischen, wenn sie unterschiedliche Prüfebenen (Voll-Target vs.
Werkzeug-Teilschritt) meinen könnten.

## Nebenfinding-Verifikationskette — plausibel, nicht zirkulär

Die Kette Implementer-Bericht → Reviewer-Reproduktion (separater
`git worktree` auf `b86cf70^`, alter Pin, `schema.yaml` ohne
`views:`-Knoten) → „keine Regression" ist methodisch sauber: Der Reviewer
hat nicht die Implementer-Aussage zitiert, sondern **denselben Befund an
einem anderen Artefakt-Stand neu erzeugt** (unabhängige Reproduktion,
keine Übernahme). Meine eigene Reproduktion an der aktuellen `main`
bestätigt den *Mechanismus* strukturell: `cdc.heartbeat`/`cdc.metrics`
entstehen ausschließlich über die psql-Nacharbeit-Skripte
(`nacharbeit-observability.sql`, `nacharbeit-heartbeat.sql|`), nicht über
`schema.yaml`. Ein deklaratives Migrationswerkzeug, das Ist- gegen
Soll-Zustand vergleicht, liest jedes Katalog-Objekt, das im Soll-Modell
nicht deklariert ist, als Kandidat für eine destruktive Operation
(`DropView`) — dieses Verhalten hängt an der *Modell-Lücke* (Views
außerhalb von `schema.yaml`), nicht am Digest-Pin. Der Blocker ist damit
pin-unabhängig erwartbar, was die „keine Regression"-Aussage zusätzlich
zur Reviewer-Reproduktion plausibilisiert. Ich habe den Lauf nicht selbst
gegen den alten Pin `d8dc38c…` wiederholt (Reviewer hat das bereits
getan; aus Zeitgründen nicht redundant dupliziert) — auf Basis der
strukturellen Erklärung und der bereits vorliegenden unabhängigen
Reproduktion sehe ich keinen Anlass zum Zweifel.

**Weiterhin unverankert:** Wie vom Reviewer vermerkt, steht dieser reale,
reproduzierbare `heartbeat`/`metrics`-Blocker nach wie vor nicht in §1
(Ausschluss) oder §6 (Risiko) des Slice-Plans und hat keinen
Beobachtungs-Register-Eintrag. Meine eigene Reproduktion bestätigt: Er
tritt bei **jedem** zweiten `make schema-rollout`-Lauf gegen dieselbe DB
auf, unabhängig vom Pin — ein stehender, nicht nur einmaliger Zustand.
`make test-integration` trifft ihn nicht (frische Umgebung je Lauf), aber
ein Operator, der `make schema-rollout` manuell gegen eine bestehende
Instanz erneut fährt, sieht ihn. Kein DoD-Blocker dieses Slice (DoD-Punkt
2 spricht explizit von den drei Views, nicht von `heartbeat`/`metrics`),
aber ein Kandidat für einen eigenen Beobachtungs-Register-Eintrag oder
eine §6-Ergänzung bei der Closure — reine Empfehlung.

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Pin v1.3.1, `make schema-validate` grün | **bestätigt** | eigener Lauf, `Validation passed`, Digest im `Makefile` identisch mit dem real gezogenen |
| 2 | Views real getestet, `source_dialect`+`columns:` deklarativ, Exit 0 Erstanlage + kein `ReplaceView`-Blocker im Folgelauf | **materiell bestätigt** (siehe eigene Reproduktion + Gegenprobe oben); Formulierungs-Präzisierung „Exit 0" vs. „kein `ReplaceView`-Blocker" für den vollen `make schema-rollout`-Lauf empfohlen (kein Blocker) |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, Exit 0, alle vier inneren Gates |
| 4 | Review durchgeführt, Report liegt vor, kein offenes HIGH | **bestätigt** | `review-slice-016.md` liegt vor (`e582a40`), F-1 (HIGH) disponiert (`04590a4`), Report bestätigt „kein offenes HIGH" |
| 5 | Doku-Update, `harness/README.md`-Bindungszeile | **bestätigt** | Zeile nennt `active_tables`/`consumer_status`/`changes` deklarativ überführt, `nacharbeit-views.sql` zurückgebaut, „seit slice-016"; keine tote Referenz auf die gelöschte Datei mehr in aktiv geltender Doku (eigener repo-weiter `grep`, siehe unten) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin Platzhalter — Planner-Closure-Arbeit |
| 7 | Reconciliation-Register, falls Inventur-Fund | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo durchgehend GF) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | keine `evidence/slice-016.md` unter `BEO-PGC/d-migrate-nacharbeit/` — Planner-Closure-Arbeit; Register-Zustand aktuell 3× `verkörpert` (unverändert), ein vierter Beleg würde den bereits zugewiesenen Ausgang nicht ändern |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen, mit Verifier-Vorschlag** | Risiko 1 (Metadaten-Lücke jenseits `columns:`/`source_dialect`) — **entfallen**: eigene Reproduktion (Erstanlage + Folgelauf, plus Gegenprobe) zeigt, dass `source_dialect`+`columns:` hinreichend sind. Risiko 2 (Digest-Bump löst Regressions-Fund im `schema.yaml`-Bestand aus) — **entfallen**: `schema-validate`/`make gates`/eigene Rollouts zeigen keinen Regressions-Fund gegen den deklarierten Bestand; der einzige reale Blocker (`heartbeat`/`metrics`) liegt außerhalb des `schema.yaml`-Bestands und ist pin-unabhängig (siehe Nebenfinding-Abschnitt). Beides bleibt **Vorschlag**, Ausgangs-Zuweisung ist Planner-Aufgabe |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen** | wellenlos — bei dieser Closure fällig, kein `git mv` bisher |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (real, nicht nur behauptet — inklusive einer echten
Gegenprobe für den kritischen Punkt), 1 korrekt entfallen (Item 7), 4
korrekt noch offen als Planner-Closure-Arbeit (Items 6, 8, 9, 10). Keine
eigenen DoD-Blocker-Findings — die einzige Anmerkung betrifft eine
sprachliche Zweideutigkeit im DoD-Text selbst (Punkt 2), kein
Implementierungs-Defekt.**

## Plan-vs-Code-Diff (gegen Plan-§3)

```
git diff --stat b86cf70^..04590a4 -- Makefile 'tools/schema/*' harness/README.md
```

liefert neun geänderte Dateien: `Makefile`, `harness/README.md`,
`tools/schema/down.sql`, `tools/schema/nacharbeit-heartbeat.sql`,
`tools/schema/nacharbeit-observability.sql`,
`tools/schema/nacharbeit-roles.sql`, `tools/schema/nacharbeit-views.sql`
(gelöscht), `tools/schema/plan.yaml`, `tools/schema/schema.yaml`.

§3 des Slice-Plans listet sechs davon: `Makefile` (zweimal — Digest und
`schema-rollout`-Target, als zwei Zeilen geführt), `tools/schema/
schema.yaml`, `tools/schema/nacharbeit-views.sql` (löschen),
`harness/README.md`, `tools/schema/plan.yaml`/`down.sql` (als
Plan-Nachzug ergänzt in `b86cf70`).

**Deckungslücke gefunden:** Die drei Dateien `nacharbeit-roles.sql`,
`nacharbeit-observability.sql`, `nacharbeit-heartbeat.sql` — geändert im
Review-Fixrunden-Commit `04590a4` (F-1-Korrektur der toten
Kopfkommentar-Querverweise) — stehen in **keiner** Zeile von §3. Das ist
ein realer, wenn auch kleiner Plan-vs-Code-Diff-Befund: Der Slice-Plan
sagt „der Implementer-Agent erweitert die Liste in seinem ersten Lauf"
(§3-Bedienhinweis), aber die Fixrunde ist ein zweiter, review-getriebener
Lauf, und die Liste wurde dort nicht nachgezogen. Inhaltlich unstrittig
(reine Kommentar-Korrektur, vom Reviewer selbst als Übergabe an den
Implementer verlangt, siehe `review-slice-016.md` F-1 „Übergabe"), aber
formal eine Lücke — vor `git mv` nach `done/` per Plan-Nachzug-Zeile in §3
zu schließen (analog zur bereits bestehenden Plan-Nachzug-Zeile für
`plan.yaml`/`down.sql`).

**Deckt sich der Diff mit §3? Im Kern ja — sechs von sechs geplanten
Dateien vollständig gedeckt, keine unangekündigte *funktionale* Änderung.
Drei zusätzliche, rein redaktionelle Dateien aus der Review-Fixrunde sind
in §3 nicht nachgezogen (siehe oben) — kein funktionaler Diff-Defekt,
aber eine offene Vollständigkeits-Lücke im Plan-Nachzug.**

## Eigene Befunde

### VF-1 — Drei review-getriebene Dateien (`04590a4`) fehlen in Plan-§3

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-016-d-migrate-1.3.1-views-retirement.md` §3
- `befund`: siehe „Plan-vs-Code-Diff" oben — `nacharbeit-roles.sql`,
  `nacharbeit-observability.sql`, `nacharbeit-heartbeat.sql` sind real
  geänderte Dateien dieses Slice, aber nicht in §3 gelistet.
- `verifizierbar`: nein — kein Sensor gleicht §3 gegen den tatsächlichen
  Diff ab.
- **Für die Closure:** trivial per Plan-Nachzug-Zeile zu schließen.

### VF-2 — Keine der DoD-Checkboxen in §2 ist gesetzt, obwohl fünf Punkte materiell erledigt sind

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-016-d-migrate-1.3.1-views-retirement.md` §2, Zeilen 85–96 (Punkte 1–3, 5)
- `befund`: Anders als im Vorgänger-Slice (slice-015, wo der
  Implementer-Commit die eigenen Punkte 1/2/3/5 direkt auf `[x]` setzte
  und nur Punkt 4 — Review, naturgemäß erst nach dem eigenen Commit
  abschließbar — offen blieb) sind in slice-016 sämtliche zehn
  Checkboxen weiterhin `[ ]`, obwohl fünf Punkte (1–5) laut dieser
  Verifikation real erledigt sind. Dieselbe Klasse wie
  `verify-slice-015.md` V-1, hier auf mehr Punkte ausgedehnt.
- `verifizierbar`: nein.
- **Für die Closure:** trivial zu beheben (Checkboxen 1–5 setzen), gehört
  vor den `git mv` nach `done/`.

## Negativbefunde

- geprüft, ohne Befund: **[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Konformität**
  — Re-Evaluierungs-Trigger feuert nicht (d-migrate drückt die Views-Operation
  jetzt aus, kein Blocker mehr für den Vertrags-Fall dieses Slice); kein
  Folge-ADR fällig.
- geprüft, ohne Befund: **Digest-Korrektheit** — `Makefile`-Pin
  `sha256:862dfb0…` identisch mit dem real gezogenen und mehrfach
  verwendeten Image in allen eigenen Läufen.
- geprüft, ohne Befund: **Hard Rule 3.3** — kein `git mv` in diesem
  Commit-Fenster, keine Lifecycle-Übergänge zu prüfen.
- geprüft, ohne Befund: **Keine toten `nacharbeit-views.sql`-Referenzen
  mehr in aktiv geltender Doku/Code** — eigener repo-weiter `grep`:
  verbleibende Treffer liegen ausschließlich in historischen
  Artefakten (`done/`-Slices, frühere Review-/Verify-Reports, dem
  Beobachtungs-Register-Beleg `evidence/slice-010.md`/`slice-015.md`,
  `state.md` mit Vermerk „bleibt technisch bestehen" — Planner-Closure
  aktualisiert das mit dem neuen Beleg) und der `Accepted`-[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (immutabel, Kontext-Abschnitt beschreibt korrekt den historischen
  Stand zum Zeitpunkt ihrer Annahme). Keiner dieser Treffer ist ein
  aktiv geltender, gebrochener Verweis.
- geprüft, ohne Befund: **Arbeitsbaum nach allen Sensor-/Probenläufen** —
  `git status` nach diesem Verifikationslauf sauber; alle selbst
  angelegten Testcontainer/-netze/-dateien entfernt,
  `tools/schema/plan.yaml`/`down.sql` auf committeten Stand
  zurückgesetzt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 (VF-1, VF-2) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (`make schema-validate`, `make gates`,
`make commit-traceability`, zwei eigene frische `schema-rollout`-Läufe
plus eine eigene Gegenprobe ohne `columns:`), 1 Item korrekt entfallen
(Reconciliation-Register), 4 Items regulär noch offen als
Planner-Closure-Arbeit. **Kein DoD-Defekt im Sinn eines unbelegten
„bestätigt"-Punkts.** Beide eigenen Findings (VF-1, VF-2) sind
Vollständigkeits-/Konsistenz-Lücken außerhalb dessen, was die
DoD-Checkliste als materiellen Liefer-Punkt zählt.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit zwei offenen
Klein-Findings vor Closure.** Kein ADR-Verstoß, kein
Traceability-/ID-Schema-Verstoß, keine Halluzination in den
Implementer-/Reviewer-Behauptungen — insbesondere die kritische
DoD-Zusage (Views real gegen den neuen Pin konvergent, sowohl Erstanlage
als auch Folgelauf, `columns:` als wirksamer Faktor) ist **eigenständig
reproduziert und per Gegenprobe kausal bestätigt**, nicht nur übernommen.
Die vom Reviewer bereits unabhängig reproduzierte
Nebenfinding-Verifikationskette (Implementer-Bericht →
Reviewer-Reproduktion im separaten Worktree) ist plausibel und nicht
zirkulär; die strukturelle Erklärung (Blocker hängt an der Modell-Lücke
`heartbeat`/`metrics`, nicht am Pin) stützt sie zusätzlich.

**Plan-vs-Code-Diff:** deckt sich im Kern mit §3 — alle sechs
ursprünglich geplanten Dateien vollständig, kein funktionaler
Überschuss. Drei review-getriebene, rein redaktionelle Dateien
(`04590a4`) sind in §3 nicht nachgezogen (VF-1).

**Vor `git mv` nach `done/` zu klären (Planner):**

1. VF-1 — Plan-Nachzug-Zeile in §3 für die drei `nacharbeit-*.sql`-Dateien
   aus `04590a4` ergänzen.
2. VF-2 — DoD-Checkboxen 1–5 auf `[x]` setzen (materiell erledigt).
3. DoD-Punkt 2 — optional die „Exit 0"/„kein `ReplaceView`-Blocker"-Zweideutigkeit
   im §7-Lerneintrag benennen (geschärfte Plan-Formulierungs-Regel), kein
   Blocker dieses Slice.
4. §6-Risiken disponieren — Verifier-Empfehlung: beide „entfallen" (siehe
   DoD-Prüfung Punkt 9 oben) — **reine Empfehlung, kein gesetzter
   Ausgang.**
5. §7 Closure-Notiz, Beobachtungs-Register-Fortschreibung
   (`BEO-PGC/d-migrate-nacharbeit`, vierter Beleg, Ausgang bleibt
   `verkörpert`) und die drei Paarungen — reguläre Planner-Closure-Arbeit.
6. Empfehlung (kein DoD-Blocker): den `heartbeat`/`metrics`-`DropView`-Blocker
   (pin-unabhängig, bei jedem Folgelauf gegen eine bereits migrierte DB
   reproduzierbar) als eigenen Beobachtungs-Register-Eintrag oder
   §6-Ergänzung aufnehmen — er steht bislang nirgends.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert; alle Testcontainer/-netze
und Temp-Dateien dieses Laufs wurden vollständig entfernt.

---

**Gate-Beleg:** `make schema-validate` und `make gates` in diesem Lauf,
beide Exit 0 (siehe Sensor-Tabelle oben). Zwei eigene frische
`schema-rollout`-Läufe (Erstanlage + Folgelauf) plus eine eigene
Gegenprobe (ohne `columns:`) real ausgeführt, Docker/PostgreSQL 18,
netzgebunden, isolierte eigene Container/Netze, vollständig entfernt nach
Lauf-Ende.
