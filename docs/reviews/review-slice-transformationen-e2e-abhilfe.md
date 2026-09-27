# Review-Report: slice-transformationen-e2e-abhilfe — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan, `ADR-0112` (Accepted, Teilfrage 4 und Folgepflicht 5)
und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-transformationen-e2e-abhilfe` (Welle `welle-transformationen`),
Diff-Range `e6c5d087..HEAD`: Lifecycle `c674d637`/`11195110`/`e6c5d087` (Moves + Ein-Zeilen-Feld),
Implementierung `c883a5bf` (neuer Go-Test `TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass`,
neue Runner-Phase „Transformationen-Nichtanwendbarkeit und Abhilfe", `docs/user/e2e-abdeckung.md`
als Erzeugnis, `harness/README.md`-Nachzug, zwei Kommentar-Korrekturen). 5 geänderte Dateien.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09", seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Ablage:** Keine Textmutation an Repo-Dateien in diesem Lauf — die Mutationsprüfung des neuen
E2E-Tests lief ausschließlich durch **Lesen** gegen die bereits existierenden, mutationsdokumentierten
Domänen-Unit-Tests (`internal/adapters/driving/replication/mapper/transformation_test.go`), nicht
durch eine eigene Textänderung. Ein realer, voller `make test-integration`-Lauf (eigenes Docker-Log,
Scratchpad) sowie ein realer, voller `make gates`-Lauf und `make test` liefen ungefiltert gegen den
Arbeitsbaum, Exit-Code jeweils direkt gelesen (`AGENTS.md` §3.9).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-transformationen-e2e-abhilfe` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan mit
  Suchlauf-Feld, §6 Risiken)
- `ADR-0112` Teilfrage 4 (sichtbarer Fehlerzustand bei nicht anwendbarer Regel) und Folgepflicht 5
  (Abhilfe-Akzeptanzkriterium (a)–(d))
- `ADR-0023` (Fehlerklassifikation, keine achte Klasse), `ADR-0128` (Vorlauf-Frist, Grenze des
  Beleg-Umfangs dieses Slice)
- `AGENTS.md` §3.1, §3.3, §3.7, §3.9, §3.12, §3.13, §3.15
- `harness/README.md` §Sensors, Zeile `make test-integration`

---

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

- **Neuer Go-Test wörtlich gelesen** (`test/integration/integration_test.go:1235–1298`): Die
  Regel `rename_column(name → label)` wird vom Runner vorab gesetzt; der Test fügt zuerst eine
  Zeile ein (id=1, Boundary-Anker), fährt dann eine **Gegenprobe** (`ALTER TABLE … ADD COLUMN
  other text`, id=2) und erst danach die **reale Kollision** (`ALTER TABLE … ADD COLUMN label
  text`, id=3). Die Gegenprobe wird aktiv geprüft (Row Image von id=2 muss `label`/`other` tragen,
  `awaitChangesViewRows` würde bei einem fälschlich ausgelösten Fehler blockieren/timeouten) — der
  Beleg ist damit an die **Namenskollision** gebunden, nicht an irgendeine Spalten-Erweiterung.
  Nach der Kollision: Heartbeat-Fehlerklasse `schema` erwartet, `cdc.changes` darf id=3 **nicht**
  tragen (Persist-before-ACK gewahrt), und beide vorherigen Zeilen (id=1, id=2) werden erneut
  gelesen und **byte-gleich** gegen ihren ersten Lesestand gehalten (Boundary).
- **Mutationsprüfung des zugrunde liegenden Domänenmechanismus** (statt einer eigenen Text-Mutation
  am E2E-Test — Begründung siehe unten): `internal/adapters/driving/replication/mapper/transformation_test.go`
  trägt bereits `TestConsumeRuleTargetCollisionAfterCompatibleExtension` mit derselben Story wie der
  neue E2E-Test (kompatible Erweiterung um den Zielnamen passiert `observeRelation` fehlerfrei, die
  nächste Change endet mit `ErrTransformationNotApplicable`/`ErrTransformationTargetCollides`; eine
  spalten-**entfernende** Relation endet stattdessen mit `ErrIncompatibleSchemaChange` **ohne**
  `ErrTransformationNotApplicable` — exakt die Abgrenzung, die auch der E2E-Test referenziert) sowie
  `TestConsumeRuleRemedyRestoresCapture` (Abhilfe `RemoveTransformation` setzt die Erfassung fort).
  Beide Tests tragen im Kommentar eine benannte, rot färbende Mutation
  (`TestConsumeRuleColumnMissingInRelationIsNotApplicable`-Kommentarblock: „in `CheckApplicable` die
  Prüfung … entfernen — dann liefert der Assembler ein ungeändertes Bild ohne Fehler"). `make test`
  (voller Repo-Lauf, Race-Detector) bestätigt beide Pakete `ok`
  (`internal/adapters/driving/replication/mapper 1.047s`).
  **Entscheidung, keinen zweiten realen `make test-integration`-Lauf mit eigener Text-Mutation zu
  fahren:** Ein voller Compose-Lauf kostet ≈8 Minuten; ich habe bereits einen für Schwerpunkt 4
  gefahren. Die Kombination aus (a) der bereits im E2E-Test eingebauten, aktiv geprüften Gegenprobe
  und (b) den bereits bestehenden, mutationsdokumentierten Unit-Tests derselben Fehlerklasse auf der
  darunterliegenden Domänenebene ist eine hinreichende Bindung — ein zweiter voller E2E-Lauf hätte
  nur bestätigt, was auf Unit-Ebene bereits mit Mutation belegt ist. Diese Abwägung steht hier
  transparent, nicht als stille Annahme.
- **Runner-Phase wörtlich gelesen** (`tools/harness/run-integration-tests.sh:4159–4225`): Ordnung
  der Schritte real nachvollzogen — (1) Slot-Drop + `docker start` + Health-Poll (Escape aus dem von
  `TestE2ESchemaChangeIncompatibleTypeChange` hinterlassenen inkompatiblen Zustand), (2) Regel
  setzen, (3) `go test -run '^TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass$'`
  (der Container endet dabei real, weil der Test über `env.pool` direkt schreibt, nicht über
  `docker exec`), (4) `docker inspect .State.Running` = `false` geprüft (Kriterium a, Container-Lauf),
  (5) Log-Sentinel geprüft, (6) `cdc.remove_transformation` **per reinem SQL-Aufruf, während der
  Container nachweislich bereits gestoppt ist** (Kriterium b — das ist eine echte Messung, kein
  Konjunktiv: der Running-Check aus Schritt (4) liegt vor diesem SQL-Aufruf), Status `pending`
  geprüft, (7) zweiter `docker start` **ohne** Slot-Eingriff (die am Slot unbestätigte Transaktion
  id=3 bleibt erhalten), Health-Poll, `bf_await_applied`, kein zweiter `schema`-Fehler, Zeile id=3
  über `cdc.changes` gelesen und Rohform geprüft (Kriterium c/d).
- **Kriterium (c) — Grenze bestätigt, wie im Plan §1/§2 benannt:** Die Phase misst **nicht** eine
  strikte Zeitordnung „Antrag applied, bevor die erste Transaktion assembliert wird" — sie misst
  ausschließlich die **Folge** (Antrag `applied`, kein zweiter `schema`-Fehler, Transaktion erscheint
  in Rohform). Das ist im Plan (§1 vierter „NICHT"-Punkt, §2 DoD-Text zu (c) wörtlich: „(c) wird über
  seine Folge belegt … und über die Ordnungs-Tests aus `start-reihenfolge`") explizit und mit Verweis
  auf die dort bereits vorhandenen Quelltext-Ordnungstests benannt — kein stiller Gap, sondern eine
  nach `AGENTS.md` §3.12 Instanz B korrekt gekennzeichnete Grenze. **Bewertung:** ausreichend für
  diesen Slice-Zuschnitt; eine echte Zeitordnungsmessung am System würde eine eigene Instrumentierung
  brauchen (z. B. einen Log-Zeitstempel-Vergleich zwischen „Antrag applied" und „erste Assemblierung"),
  die weder im Plan noch in `ADR-0112` Folgepflicht 5 gefordert ist — die ADR verlangt die Wirkung
  (a)–(d), nicht den Beweis einer internen Ausführungsreihenfolge. Kein Finding, INFO (siehe unten).
- **Escape-Mechanismus geprüft (Schwerpunkt 3):** Zwei Neustarts mit unterschiedlicher
  Slot-Behandlung sind im Kommentar begründet (erster Neustart: Slot-Neuanlage, weil die von
  `TestE2ESchemaChangeIncompatibleTypeChange` hinterlassene Transaktion „echt inkompatibel" ist und
  ein Replay ohne Slot-Neuanlage denselben Fehler erneut auslöste; zweiter Neustart: **kein**
  Slot-Eingriff, weil genau die am Slot unbestätigte Transaktion id=3 der Beleg für Kriterium (d)
  ist). Die Phase ist der **letzte** Rundlauf des Skripts — danach folgt nur noch die
  Abdeckungstabellen-Erzeugung. Die `cleanup()`-Funktion (`trap cleanup EXIT`) fährt unbedingt
  `$COMPOSE down -v --remove-orphans`, unabhängig vom Slot-Zustand am Laufende — kein versteckter
  Seiteneffekt auf einen späteren oder parallelen Zustand, weil kein späterer Zustand mehr folgt.
- **Realer, voller `make test-integration`-Lauf** (eigener Lauf, `:dev`-Image frisch mit `make image`
  gebaut — vollständig aus dem Cache, Digest unverändert `sha256:01e6ebc7f32e…`): Exit sauber, kein
  `FAIL`, kein `panic`. Zitierte Ausgabezeilen:
  ```
  === RUN   TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass
  --- PASS: TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass (0.39s)
  PASS
  ok  	github.com/pt9912/pg-change-feed/test/integration	0.395s
  run-integration-tests: Transformationen-Nichtanwendbarkeit und Abhilfe (ADR-0112 Folgepflicht 5)
  belegt — Kollision auf feed_e2e_transform_abhilfe beendete den Erfassungspfad real
  (error_class=schema, Sentinel im Log), cdc.remove_transformation
  (677d573f-eeee-4a74-ab14-29675a65df62) während des Stillstands beantragt (pending), nach dem
  Neustart applied vor der ersten Transaktion der Tabelle, die zuvor nicht bestätigte Zeile (id=3)
  erscheint über cdc.changes in Rohform ohne zweiten schema-Fehler
  run-integration-tests: E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht
  dem Quelltext-Stand
  ```
  `git status --porcelain` nach dem Lauf: leer — keine committete Datei bewegt, `docs/user/e2e-abdeckung.md`
  byte-gleich bestätigt (durch den Runner selbst **und** durch die leere `git status`-Ausgabe).
- **Realer, voller `make gates`-Lauf** (eigener Lauf, ungefiltert, Exit direkt geprüft): Exit 0 —
  `baseline-verify` (v6.9.0, 54 Dateien), `docs-check` (1359 Dateien, 0 Befunde), `a-check` (0
  Befunde), `generated-sync` (byte-gleich), `coverage-gate` (85.40 % ≥ 80 %),
  `commit-traceability` (`HEAD~5..HEAD`, keine Struktur-ID im Betreff) — deckungsgleich mit den im
  DoD genannten Zahlen bzw. besser (Coverage 85.40 % gemessen, kein Wert im DoD dieses Slice
  genannt).
- **`make test`** (voller Repo-Lauf mit Race-Detector, eigenständig gefahren): Exit 0, alle Pakete
  `ok`, inklusive `internal/bootstrap`, beide `mapper`-Pakete und `test/integration`.
- **`make docs-check`**: eigenständig gefahren, Exit 0, 0 Befunde.
- **`make commit-traceability RANGE=c674d637~1..HEAD`**: eigenständig gefahren, Exit 0 — 4 Commits,
  jeder trägt `(ADR-0112)`, keine `SPEC-*`/`ARC-*`-Kennung im Betreff.
- **`make kommentar-kennungen DIFF=e6c5d087`**: Exit 0, keine Ausgabe — kein Kandidat. Die zwei
  Kommentar-Korrekturen (Entfernung von „, letzter" an `TestE2ESchemaChangeIncompatibleTypeChange`)
  wurden zusätzlich per `git grep -n "eigener, letzter go-test-Aufruf"` gegen den ganzen Baum
  geprüft: **0 Treffer** — die alte, jetzt falsche Formulierung ist vollständig entfernt; genau
  zwei Stellen tragen „als eigener go-test-Aufruf" (Zeile 4045, unverändert seit vorher, und die
  korrigierte Zeile 4137), keine trägt mehr „letzter" an der falschen Stelle.
- **`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-transformationen-e2e-abhilfe.md`**:
  Exit 0, „8 Zeilen stimmen" — alle vier Suchmuster-Paare (je Parent `e6c5d087` und Diff-Stand)
  exakt wie im Plan behauptet.
- **Diff-Umfang:** `git diff --name-only e6c5d087..HEAD` liefert exakt die 5 im Auftrag genannten
  Dateien (Plan-Datei, `docs/user/e2e-abdeckung.md`, `harness/README.md`,
  `test/integration/integration_test.go`, `tools/harness/run-integration-tests.sh`) — keine
  unbenannte Nebenwirkung.
- **`make doc-immutable` — Target existiert in diesem Repo nicht** (`grep -rn immutable Makefile
  harness/mk/*.mk harness/README.md` liefert 0 Treffer). Ich habe deshalb stattdessen manuell
  geprüft: `ADR-0112` selbst ist **nicht** Teil des Diffs (keine ADR-Datei in
  `git diff --name-only`), Immutability nach §3.5 ist damit trivial gewahrt — kein Fund, aber kein
  Gate-Beleg dieser Art existiert (dies ist eine Beobachtung über den Auftrag, kein Slice-Befund).
- **Fremde-Datei-Träger geprüft** (§3.13-Nachzug-Pflicht): Die vier Muster aus dem Suchlauf treffen
  ausschließlich innerhalb des eigenen Diffs bzw. bestätigen „nichts zu tun" (die beiden
  `ErrIncompatibleSchemaChange`/`ErrTransformationNotApplicable`-Symbole standen bereits am Parent
  nebeneinander, kein dritter Auslöser entsteht) — keine fremde Datei bleibt unbehandelt stehen.
- **`BEO-PGC/kein-admin-weg-schema-fehler-recovery`** (§1 Abgrenzung, Restrisiko 4 des Auftrags):
  `grep -rn` bestätigt den Eintrag existiert unter `docs/plan/planning/observations/BEO-PGC/kein-admin-weg-schema-fehler-recovery/`,
  konsistent als „offen, 1×" in mehreren Plänen (`welle-transformationen.md`,
  `slice-transformationen-betriebsdoku.md`, `slice-capture-transient-wiederholung.md`) referenziert
  und in **diesem** Slice-Plan §1 korrekt als bewusste Abgrenzung geführt („der Beleg deckt die
  Abhilfe für **diese** Ursache; die Frage … bleibt offen") — kein Finding, korrekt geführt.
- **Verweigerte Aktion / §3.15 (Schwerpunkt 9):** Laut Aufgabenstellung wurde ein `docker run … go
  test …`-Aufruf außerhalb von `make` von der Berechtigungsschicht verweigert, und der Implementer
  wich danach auf `make test` — den regulär vorgesehenen, sanktionierten Weg — aus, **ohne** einen
  Ersatzweg zum selben (unsanktionierten) Ziel zu suchen. Ich kann die Verweigerung selbst nicht aus
  Diff oder Artefakten nachvollziehen (kein Gate hinterlässt eine Signatur einer Permission-Denial);
  dieser Punkt ist **übernommen**, nicht gemessen (`AGENTS.md` §3.12/§3.15 Grenze). Auf Basis der
  Angabe: korrekt gehandhabt — `make test` ist nicht nur kein Umgehungsweg, sondern der nach §3.1
  ohnehin vorgeschriebene Weg; es liegt kein Ersatzweg-zum-selben-unsanktionierten-Ziel vor.
- **Dangling-Docker-Ressourcen:** `docker volume ls -f dangling=true -q | wc -l` vor meinen eigenen
  Docker-Läufen: **38**. Nach `make image` + `make test-integration` + `make gates` (eigene Läufe):
  **39** — ein neues dangling Volume, vermutlich ein Zwischenlayer des `coverage`-Build-Stage aus
  `make gates`. Kein `prune` verwendet.

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings in diesem Lauf.

## Negativbefunde

- geprüft, ohne Befund: **Bindung des neuen E2E-Tests an die Namenskollision.** Gegenprobe
  (`ADD COLUMN other`) aktiv geprüft, reale Kollision (`ADD COLUMN label`) korrekt ausgelöst,
  Boundary (id=1/id=2 byte-gleich nach der Kollision) bestätigt.
- geprüft, ohne Befund: **Zugrunde liegender Domänenmechanismus.** Bereits bestehende,
  mutationsdokumentierte Unit-Tests (`TestConsumeRuleTargetCollisionAfterCompatibleExtension`,
  `TestConsumeRuleRemedyRestoresCapture`) binden exakt dieselbe Story auf Domänenebene; `make test`
  bestätigt beide Pakete grün.
- geprüft, ohne Befund: **Kriterium (a)/(b) der Runner-Phase.** Container-Stopp real per
  `docker inspect` geprüft, `cdc.remove_transformation` real erst danach beantragt (echte
  Reihenfolge, nicht nur behauptet), Status `pending` real gelesen.
- geprüft, ohne Befund: **Kriterium (c)/(d) — Grenze korrekt benannt.** Die Phase misst die Folge,
  nicht die strikte Zeitordnung; Plan §1/§2 benennen diese Grenze explizit und verweisen auf die
  Ordnungstests aus `slice-transformationen-start-reihenfolge` — kein stiller Gap.
- geprüft, ohne Befund: **Escape-Mechanismus (zwei Neustarts, unterschiedliche Slot-Behandlung).**
  Kein versteckter Seiteneffekt: die Phase ist real der letzte Rundlauf, `cleanup()` räumt am Ende
  unbedingt ab.
- geprüft, ohne Befund: **Realer `make test-integration`-Lauf.** Grün, neuer Test und neue Phase
  wörtlich zitiert; `docs/user/e2e-abdeckung.md` unverändert (Runner-Meldung **und** leere
  `git status`), keine committete Datei bewegt.
- geprüft, ohne Befund: **Kommentar-Korrekturen.** Genau zwei Stellen tragen „als eigener
  go-test-Aufruf" im ganzen Baum, keine mehr die falsche „letzter"-Zuschreibung an
  `TestE2ESchemaChangeIncompatibleTypeChange`; `make kommentar-kennungen DIFF=e6c5d087` liefert 0
  Kandidaten.
- geprüft, ohne Befund: **§3.13-Suchlauf.** `make suchlauf-nachmessen` bestätigt 8/8 Zeilen exakt;
  Suchraum (ganzer Baum ohne die drei Standard-Ausnahmen) und Muster (Symbol, Zählwort/Beschreibung,
  Hedge) sind vollständig für die im Plan benannten bewegten Eigenschaften.
- geprüft, ohne Befund: **Traceability, ID-Schema, Moves.** Alle 4 Commits tragen `(ADR-0112)`,
  keine `SPEC-*`/`ARC-*` im Betreff; die drei Lifecycle-Commits sind reine `git mv`-/
  Ein-Zeilen-Commits (§3.3 eingehalten).
- geprüft, ohne Befund: **Diff-Umfang.** `git diff --name-only` liefert exakt die 5 im Auftrag
  genannten Dateien, keine unbenannte Nebenwirkung.
- geprüft, ohne Befund: **`BEO-PGC/kein-admin-weg-schema-fehler-recovery`.** Korrekt als offene,
  bewusste Abgrenzung in §1 des Slice-Plans geführt.
- geprüft, ohne Befund: **Docker-only, Suppression, Spec-Stratum, Zwei-Quellen-Drift.** Kein
  `//nolint`; kein Host-Werkzeug außerhalb der erlaubten Klasse im Diff; `spec/lastenheft.md` nicht
  berührt; kein doppelt geführter Zustand ohne erklärten Gewinner.
- geprüft, ohne Befund: **Realer `make gates`-Lauf.** Exit 0, alle sechs inneren Gates grün,
  Coverage 85.40 % ≥ 80 %.
- geprüft, ohne Befund: **`BEO-PGC/ersatzweg-nach-verweigerter-aktion` (§3.15, übernommen).** Der
  gewählte Ersatzweg (`make test`) ist der ohnehin vorgeschriebene Weg, kein Umgehen der
  Verweigerung — auf Basis der Angabe im Auftrag, nicht selbst aus Artefakten nachweisbar.

## INFO

- **`make doc-immutable` existiert in diesem Repo nicht** (im Auftrag genannt, im Makefile nicht
  vorhanden) — kein Slice-Befund, sondern eine Diskrepanz zwischen Auftragstext und Repo-Bestand;
  Immutability von `ADR-0112` ist über die triviale Beobachtung „ADR nicht im Diff" gewahrt.
- **Kriterium (c) bleibt strukturell eine Folge-Messung, keine Zeitordnungs-Messung** — vom Plan
  selbst so benannt (§1, §2) und durch `AGENTS.md` §3.12 Instanz B gedeckt formuliert. Für einen
  künftigen Slice, der eine echte Ordnungsgarantie im laufenden System verlangt (nicht nur an
  Fakes/Quelltext), wäre eine eigene Instrumentierung nötig — nicht Gegenstand dieses Slices.
  Hinweis für Planner/Architect, kein Reviewer-Fund.
- **Dangling-Docker-Volumes 38 → 39** durch meine eigenen Läufe (`make image`/`make
  test-integration`/`make gates`) — zur Kenntnis, kein `prune` verwendet.
- **Kein eigener Mutationslauf an der Bash-Runner-Phase selbst** (nur Lesen) — angesichts der
  Kombination aus aktiv geprüfter Gegenprobe im Go-Test und mutationsdokumentierten
  Domänen-Unit-Tests halte ich das für ausreichend; ein Fixrunde-Anlass wäre es nur, wenn diese
  Kombination fehlte.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** keine — sauberer Lauf ohne HIGH/MEDIUM/LOW-Befund.

## Verdikt

**Merge-blockierend: nein.** Der neue E2E-Test bindet die Nichtanwendbarkeits-Prüfung korrekt an
die Namenskollision (aktive Gegenprobe, Boundary bestätigt), die neue Runner-Phase bindet drei der
vier Buchstaben des Abhilfe-Akzeptanzkriteriums real und misst den vierten (c) transparent nur über
seine Folge — eine im Plan explizit benannte, durch `AGENTS.md` §3.12 gedeckte Grenze, kein stiller
Gap. Der Escape-Mechanismus (zwei Neustarts, unterschiedliche Slot-Behandlung) ist sauber, ohne
Seiteneffekt auf spätere Zustände, weil die Phase real der letzte Rundlauf ist. Ein realer, voller
`make test-integration`-Lauf und ein realer, voller `make gates`-Lauf sind in diesem Review
eigenständig gefahren worden und beide grün; `docs/user/e2e-abdeckung.md` ist byte-gleich
bestätigt. Traceability, Lifecycle, Suchlauf (§3.13) und Kommentar-Form (§3.7) sind je eigenständig
nachgemessen, alle konform.

**DoD-Checkbox-Nachzug:** Da keine Fixrunde nötig ist (0 HIGH, keine Findings), ziehe ich die
DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan selbst auf
`[x]` nach, mit Verweis auf diesen Report, im selben Commit, der den Report anlegt
(Skill §DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** keine Fixrunde. Die vier INFO-Punkte sind Hinweise ohne erwartete Aktion; der Hinweis
zu Kriterium (c) richtet sich an einen künftigen Slice/Architect-Zug, falls eine echte
Zeitordnungsgarantie am System jemals gefordert wird — nicht an diesen Slice.
