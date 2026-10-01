# Review-Report: slice-routing-spec-nachzug — 2026-10-01

**Review-Art:** Plan/Design — der Diff ist reine Spec-Doku (Pflichtenheft und
Architektur-Sicht) plus Slice-Plan; geprüft gegen Plan, ADRs und `AGENTS.md`
Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `b8085839..4b6dc973`, zwei Commits `aee60747`
(Spec) und `4b6dc973` (Plan und Suchlauf-Feld); drei Dateien:
`spec/pflichtenheft.md` (+349/−), `spec/architecture.md`, Slice-Plan
`slice-routing-spec-nachzug` (Lifecycle `in-progress`). `spec/lastenheft.md`
ist im Diff nicht enthalten.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-routing-spec-nachzug` (Ziel, Umfang, Suchlauf-Feld) und
  Welle `welle-routing`; die Pläne der Folge-Slices `slice-routing-backfill-pfad`,
  `slice-routing-antragsweg`, `slice-routing-lesewege` (nur als Adressen)
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  (Accepted, Teilfragen 1–7, R1–R6, Folgepflichten, Entscheidungen 1–6) und
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
  (Accepted, Festlegungen 1 und 2, V3 offen)
- Lastenheft [`LH-FA-CFG-008`](../../spec/lastenheft.md),
  [`LH-FA-CAP-009`](../../spec/lastenheft.md) und die Pflichtenheft-Stellen
  [`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md),
  [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md),
  [`SPEC-002`](../../spec/pflichtenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md),
  [`SPEC-019`](../../spec/pflichtenheft.md),
  [`SPEC-020`](../../spec/pflichtenheft.md),
  [`SPEC-021`](../../spec/pflichtenheft.md),
  [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-024`](../../spec/pflichtenheft.md),
  [`SPEC-029`](../../spec/pflichtenheft.md),
  [`SPEC-030`](../../spec/pflichtenheft.md),
  [`SPEC-031`](../../spec/pflichtenheft.md),
  [`SPEC-032`](../../spec/pflichtenheft.md),
  [`ARC-001`](../../spec/architecture.md)
- `AGENTS.md` (Hard Rules §3.4, §3.5, §3.7, §3.9, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Vorgänger-Review gleicher Art:
  [`review-slice-transformationen-spec-nachzug`](review-slice-transformationen-spec-nachzug.md)
  (Report-Form, Klassen)

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen):

- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-spec-nachzug.md`
  → Exit 0, `suchlauf-nachmessen: 16 Zeilen stimmen` (acht Zeilen am Parent
  `30fd6cb5`, acht am Arbeitsbaum `diff`). Die Zahlen des Feldes stimmen; die
  Befund-Spalte des Plans (Zeile 2: 3 Treffer, Zeile 6: 9 Treffer) deckt sich
  mit meinem Lesen der Treffer.
- `make docs-check` → Exit 0, `d-check: 1469 Datei(en) geprüft, 0 Befund(e)`
  (ungefiltert in eine Datei gesichert, Exit gesondert gelesen, `AGENTS.md`
  §3.9), vor Anlage dieses Reports.
- `AGENTS.md` §3.4: `git diff b8085839 HEAD -- spec | grep '^+' | grep -c -E
  'ADR-0|slice-|welle-'` → `0`; ich habe den Diff von `spec/architecture.md`
  zusätzlich vollständig gelesen (nicht nur den `matrix`-Lauf): die fünf
  Hunks tragen weder Wellen-, Slice- noch ADR-Bezug.
- Zählwörter am Kopf (§3.13), eigene Stichproben außerhalb des Plan-Musters:
  `grep -rn -E 'sieben (Arten|Antragsarten|Werte|SQL)|fünf übrigen|sechs
  übrigen|zwei optionale|zehn Felder|dreizehn Felder' spec harness docs/user`
  — `request_kind` jetzt neun Werte (Tabelle neun Einträge), `column_name`
  „sieben übrigen" (9 − 2 = 7, Aufzählung der sieben Namen stimmt),
  `rule_name` „fünf übrigen" (9 − 4 = 5, stimmt), `rule_spec` „sieben übrigen"
  (9 − 2 = 7, stimmt), Architektur „neun Arten" bei neun Tabellenzeilen,
  `SPEC-019` „neun SQL-Funktionen"; `zehn Felder`/`dreizehn Felder` in der Spec
  unverändert und weiterhin wahr (das Label ist nach `SPEC-020`/`SPEC-021`/
  `SPEC-024`/`SPEC-031` kein Nachrichtenfeld). Träger außerhalb des Diffs:
  `harness/targets/schema-rollout.md:74` und `:90` („sieben Funktionen"), das
  Handbuch (`zwei optionale` an drei Stellen) — beide sind an Adressen
  gemeldet, die den Gegenstand annehmen (`git grep` der Kernbegriffe im Plan
  der Adresse: `slice-routing-antragsweg` führt `sieben SQL-Funktionen|sieben
  Funktionen` in seinem Suchlauf und erwähnt `harness/targets/schema-rollout.md`;
  `slice-routing-betriebsdoku` ist der genannte Handbuch-Träger).
- Text-Abgleich `LH-FA-CFG-008.a`, `SPEC-002`, `SPEC-008`, `SPEC-019`,
  `SPEC-020`/`-021`/`-022`/`-024`, `SPEC-031`, `SPEC-032` und die
  Architektur-Sicht mit `ADR-0137` (Teilfragen 1–6, R1–R6, Entscheidungen 1–6)
  und `ADR-0138` (Festlegungen 1 und 2) Satz für Satz; Abweichungen und
  Ergänzungen stehen in F-1 bis F-6 und in der Tabelle unter „Setzungen".
- Lifecycle und Commit-Struktur: `git log --oneline b8085839..HEAD` — zwei
  Commits, beide Messages tragen `ADR-0137`, `ADR-0138`, `LH-FA-CFG-008`,
  Betreffs ohne `SPEC-`/`ARC-`-Kennung.
- Keine Mutation gefahren: die Zusagen des Diffs sind Prosa-Aussagen; ihre
  Eingabeseite ist der Wortlaut der ADRs, geprüft durch Lesen.

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Fail-closed-Prüfung vor dem Commit nennt den Routing-Regelstand nicht; „Regelstand zum Run" ist undefiniert

- `kategorie`: MEDIUM
- `quelle`: Maintainability (unklare Fehlerbehandlung am Rand des
  Spec-Bereichs); Prüfmuster „Nachzug widerspricht dem Nachbarn im selben
  Träger"; [`LH-FA-CAP-009`](../../spec/lastenheft.md) (Fail-closed);
  [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Teilfrage 6 (Backfill)
- `pfad`: `spec/pflichtenheft.md:247-253` (Absatz „Fail-closed vor dem
  Commit", im Diff unberührt), `spec/pflichtenheft.md:262-265` (neuer Absatz
  „Ziel der Backfill-Changes"), `spec/pflichtenheft.md:1351-1354`
  (`SPEC-032`, Bildbasis „Backfill-Change wie ein `INSERT`")
- `befund`: Der Absatz „Fail-closed vor dem Commit" führt als zu prüfenden
  Stand weiter nur den Ausschlussstand und den Regelstand der Transformationen;
  der neue Absatz sagt, Backfill-Changes trügen `route_target` „des
  Regelstands zum Run", ohne zu sagen, ob dieser Stand beim Annehmen, beim
  Start, je Block oder beim Commit gilt und ob eine Änderung des
  Routing-Regelstands während des Runs den Run beendet. Die Antragsverarbeitung
  läuft nebenläufig zum Worker; ein `set_route`/`remove_route` zwischen zwei
  Blöcken ergibt ein Label aus zwei Ständen. Weder `ADR-0137` (Teilfrage 6:
  „das Label des Regelstands zum Run") noch `ADR-0138` entscheidet das; der
  Plan des Folge-Slices `slice-routing-backfill-pfad` (§1 Ziel, DoD
  „Fail-closed") setzt dagegen eine Erweiterung der Prüfung mit Klasse
  `configuration` bereits voraus — eine Zusage, die die Spec (Rang 2) nicht
  trägt. Der Implementer hat die Lücke im Plan benannt (Umfangszeile
  „eine Fail-closed-Prüfung … steht **nicht** darin"); der Spec-Text selbst
  lässt sie weder offen markiert noch geschlossen.
- `verifizierbar`: nein — kein Gate liest die Lücke; Falsifikation ist das
  Lesen von `LH-FA-CAP-009.a` mit der Frage „was tut der Run, wenn zwischen
  Block 1 und Block 2 `cdc.set_route` `applied` wird?"
- `klasse`: Regelstand-Dimension in der Fail-closed-Prüfung nicht bestimmt
- Architect-Frage: siehe A-1.

### F-2 — `target` außerhalb des Alphabets liefert auf allen Lesewegen „leer, kein Fehler" — Zusage breiter als die Entscheidung, ohne Erwartungs-Kennzeichnung

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz B (Aussage über eine Menge: „alle
  Lesewege"); [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
  Festlegung 1 (letzter Punkt)
- `pfad`: `spec/pflichtenheft.md:974` (`SPEC-020` Zeile Request),
  `spec/pflichtenheft.md:1028` (`SPEC-022` Zeile Zustellziel),
  `spec/pflichtenheft.md:1263` (`SPEC-031` Zeile `ReadChanges`); `SPEC-021`
  (Zeile Query-Parameter) übernimmt die Form über `SPEC-020`
- `befund`: `ADR-0138` führt „ein `target` außerhalb des Alphabets liefert
  eine leere Liste, kein Fehler" nur für `ReadChanges` und ausdrücklich als
  *Erwartung*, vom Slice `slice-routing-lesewege` zu belegen. Die Spec setzt
  den Satz unbedingt (ohne „Erwartung" oder „hergeleitet") für den gRPC-Stream,
  `GET /changes`, `ReadChanges` und — über die Zeile Query-Parameter — den
  SSE-Stream; für die beiden SQL-gestützten Wege ist nicht bestimmt, was ein
  `target` mit dem Zeichen U+0000 (das PostgreSQL in einem Text-Parameter
  ablehnt, hergeleitet, nicht gemessen) tut, obwohl „kein Fehler" zugesagt ist.
  Der Satz ist eine Setzung der Spec über die ADR-Festlegung hinaus
  (Alternative: `400`/`InvalidArgument` für eine Alphabet-Verletzung).
- `verifizierbar`: nein
- `klasse`: Erwartung der ADR als unbedingte Zusage auf mehr Wege ausgedehnt
- Architect-Frage: siehe A-2.

### F-3 — Grenzen der Regelform ohne Herkunft: `order` ≤ 2147483647, „ohne Bruchteil und Exponent"

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B; Maintainability (Setzung des
  Implementers)
- `pfad`: `spec/pflichtenheft.md:1332` (`SPEC-032` Tabelle, Zeile `order`),
  `spec/pflichtenheft.md:1317-1321` (Vorspann)
- `befund`: `ADR-0137` Teilfrage 3 sagt „positive ganze Zahl". Die Spec setzt
  den Wertebereich 1 bis 2147483647 und schließt die Schreibweise mit Bruchteil
  oder Exponent aus; der Vorspann kennzeichnet Grenzen als Festlegung
  „über den Wortlaut des Lastenhefts hinaus", die Obergrenze geht aber über
  den Wortlaut der ADR hinaus, und für sie ist kein Grund (Spaltentyp,
  Domänentyp) genannt. `jsonb` selbst trägt die Grenze nicht.
- `verifizierbar`: nein
- `klasse`: Setzung ohne genannte Herkunft

### F-4 — zwei Zeilen nach Teilersetzung nicht umbrochen

- `kategorie`: LOW
- `quelle`: Maintainability (dieselbe Klasse wie F-6 des Vorgänger-Reviews
  [`review-slice-transformationen-spec-nachzug`](review-slice-transformationen-spec-nachzug.md))
- `pfad`: `spec/architecture.md:418` (96 Zeichen: „Snapshots nicht anwendbar,
  endet der Run `failed` mit der Klasse `schema`, bevor die erste Zeile"),
  `spec/pflichtenheft.md:1348` (101 Zeichen: „`jsonb`, ohne eigene
  Längenvorgabe; die leere Zeichenkette ist zulässig. Die Bedingung trifft,
  wenn")
- `befund`: Beide Absätze folgen sonst dem Umbruch von etwa 80 Zeichen; die
  Ersetzung hat je eine Zeile länger stehen lassen (gemessen mit `awk
  'length($0)>…'` über die Hunks).
- `verifizierbar`: nein — `structure` prüft Zellenlängen, keine Zeilenlängen
- `klasse`: Umbruch nach Teilersetzung nicht nachgezogen (Wiederholung des
  Musters aus dem Vorgänger-Review, dort LOW)

### F-5 — Lastenheft unverändert: „Zustellung" als Abruf/Abonnement ist eine Lesart des Pflichtenhefts, die allein in der ADR verankert ist

- `kategorie`: INFO
- `quelle`: Spec-Stratum („präzisieren ja, erweitern nie");
  [`LH-FA-CFG-008`](../../spec/lastenheft.md) Happy Path;
  [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Entscheidung 1
- `pfad`: `spec/pflichtenheft.md:364-368` (`LH-FA-CFG-008.a`, „Zielmodell":
  „keine aktive Zustellung an externe Senken"), `spec/lastenheft.md:310-330`
- `befund`: Das Lastenheft sagt „wird sie dem Zustellziel zugestellt"; das
  Pflichtenheft legt „Zustellung" als Abruf oder Abonnement eines Kanals über
  die vorhandenen Wege fest und schließt aktive Senken aus. Das ist eine
  einschränkende Lesart; ihre Legitimation ist die Entscheidung des
  Auftraggebers in `ADR-0137` (Entscheidung 1), im Lastenheft findet sich
  davon nichts, und es gibt keinen Eintrag in der Lastenheft-Versionstabelle.
  Im Konfliktfall gewinnt nach dem Kopf von `spec/pflichtenheft.md` das
  Lastenheft. Die Out-of-Scope-Zeile des Lastenhefts („vollständige
  Routing-Sprache nicht gefordert") stützt die Lesart nur für die
  Ausdrucksform, nicht für die Zustellform. Verstoß gegen die Spec-Stratum-HIGH
  liegt nicht vor: das Lastenheft wurde nicht geändert, und es wird keine
  neue bindende Anforderung eingeführt.
- `verifizierbar`: nein
- `klasse`: Pflichtenheft-Lesart ohne Gegenzeichnung im Vertragsstratum
- Vertretbar im Sinne der Frage (g): ja, solange die Abnahme gegen
  `LH-FA-CFG-008` die Lesart der ADR akzeptiert; das ist eine Auftraggeber-,
  keine Reviewer-Frage (siehe A-3).

### F-6 — Sicht nennt Use-Case-Namen und die gegenseitige Prüfung von Ausschluss und Routing

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.4 (Sicht trägt Komponenten und Sequenzen);
  Maintainability
- `pfad`: `spec/architecture.md:277-278` (Zeilen `set_route`/`remove_route`
  mit `SetRouteUseCase`/`RemoveRouteUseCase`), `spec/architecture.md:295-300`
- `befund`: Die Namen entsprechen dem Muster der Bestandszeilen
  (`SetTransformationUseCase`) und dem Plan von `slice-routing-antragsweg`
  (derselbe Name an der Adresse, per `git grep` bestätigt); sie sind eine
  Setzung des Implementers, die kein ADR-Text trägt. Der Absatz zur
  gegenseitigen Prüfung von Routing-Bedingung und Spaltenausschluss wiederholt
  einen Spec-Inhalt (R3) in der Sicht — dasselbe tut der Bestandsabsatz für die
  Transformationen.
- `verifizierbar`: nein
- `klasse`: Setzung im Sicht-Stratum ohne ADR-Träger (Bestandsmuster)

---

## Setzungen des Implementers (Frage d) — Einzelbefund je Punkt

„Zu Recht" heißt: die Spec trägt die Aussage als Zusage an die Umsetzung, sie
folgt aus dem ADR-Text oder dem Bestandsmuster, und eine Gegenentscheidung
bräuchte keinen Auftraggeber.

| Setzung | Befund | Einordnung |
|---|---|---|
| `order` JSON-Ganzzahl 1..2147483647 | ADR: „positive ganze Zahl"; Obergrenze und Schreibweise sind Zusatz | F-3 (LOW) |
| `when.column` zeichengenau, nicht leer, ohne U+0000 | wortgleich zum Bestandsmuster von `SPEC-030` (`column`), Gegenstand des Katalog-Vergleichs ist R3 (ADR) | zu Recht |
| `when.equals` darf leer sein, keine Längengrenze | folgt aus „Zeichenkette" (ADR); `NULL`/abwesend ist durch die Abwesenheits-Regel getrennt abgedeckt; U+0000 in `equals` scheitert am `jsonb`-Cast wie bei den Transformationen | zu Recht |
| `rule_name`-Alphabet wie Transformationen | ADR: „eigener Namensraum", Muster `ADR-0112`; Alphabet-Gleichheit ist ein Bestandsmuster | zu Recht (INFO: Wiederverwendung von `SPEC-030` Bezeichner) |
| Fehlertext-Wortlaute und Prüfreihenfolge | ADR verlangt nur „`failed` mit Fehlertext"; Form (Klartext, Doppelpunkt, Adresse) ist die Bestandsform der Transformationen; Spec kennzeichnet die Wortlaute als Zusage (`spec/pflichtenheft.md:943`) | zu Recht |
| drei Texte für R4 | R4 der ADR („höchstens eine Regel ohne `when`, und sie trägt die höchste `order`") verlangt beide Hälften; die dritte Zeile (Regel mit `when` hinter geführter Regel ohne `when`) ist die dritte Verletzungsrichtung derselben Invariante | zu Recht |
| Antrag gegen Tabelle ohne Bindung endet `applied` | Bestandsmuster `SPEC-019` (Spalten- und Transformations-Antragsarten, Zeilen 797 und 872); [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 6 „Dauerhaftigkeit … in jedem Pfad, der eine Bindung anlegt" trägt es | zu Recht |
| `include_column` bleibt zulässig (R3 nur in Richtung `exclude_column`) | ADR R3 nennt nur `exclude_column`; folgt | zu Recht |
| `target` außerhalb des Alphabets = leere Antwort auf allen Wegen | ADR nennt es als Erwartung nur für `ReadChanges` | F-2 (MEDIUM), Frage A-2 |
| Use-Case-Namen in der Sicht | Bestandsmuster, Plan der Adresse stimmt überein | F-6 (INFO) |
| „Regelstand zum Run" ohne Definition | ADR-Wortlaut übernommen, Lücke bleibt | F-1 (MEDIUM), Frage A-1 |

## Beantwortung der Prüfauftrags-Punkte

**(a) Treue zu den ADR-Festlegungen.** Gegen den Text gelesen und treu:
Zielmodell (Kanal je Quelle, Alphabet `[a-z0-9][a-z0-9_-]{0,62}`, höchstens ein
Ziel, kein Standardziel), Regelschlüssel `target`/`order`/`when`, R1–R6 (Wortlaut
je Invariante), Auflösung „erster Treffer in aufsteigender `order`", NULL ohne
Standardziel, Abwesenheit = Nicht-Treffer, Label nicht rückwirkend (Altbestand
über Backfill-Run), Auswahl statt Ausblendung, Klasse `schema` im Erfassungspfad
und im Run, Filter-Parameter je Weg, `SPEC-031` `target = 7`, NATS-Subjekt
`cdc.route.<source_id>.<ziel>` und „Wecksignal trägt kein Ziel",
`SPEC-008`-Zeile `schema`. Abweichung: F-1 und F-2.

**(b) §3.12.** Alle in der Frage genannten Aussagen sind als hergeleitet,
nicht gemessen oder Zusage gekennzeichnet: DELETE gegen PostgreSQL 17/18
(`spec/pflichtenheft.md` zweimal), NATS-Alphabet (zweimal), NATS-Last
(„nicht gemessen"), Proto3-Kompatibilität (`SPEC-031`), V3 / Erreichbarkeit
(„Offen und nicht gemessen"), Gleichheit der Annahmemenge mit
`cdc.set_transformation` („hergeleitet"). Keine Aussage steht als belegt, die die
ADRs als hergeleitet führen. Befund dazu nur F-2 (unbedingter Satz, wo die ADR
„Erwartung" sagt).

**(c) §3.4 / Stratum und Mehrumfang.** Kein ADR-, Slice- oder Wellen-Bezug in
`spec/architecture.md` (Diff gelesen, nicht nur `matrix`). Der Mehrumfang
(Komponentenliste `ARC-001`, Backfill-Absatz, Fehlermodell-Zeile) ist
gerechtfertigt: alle drei Stellen nannten „Transformationsregel" allein als
Gegenstand der Regel-Wirkung und würden nach dem Nachzug falsch lesen (§3.13,
Träger-Nachzug). Der Plan weist den Mehrumfang aus (Umfangszeile und
Befund-Spalte). Nebenbefund: F-4 (Umbruch).

**(e) Fail-closed für den Routing-Regelstand.** Die Spec lässt die Lücke offen
(F-1); sie behauptet die Prüfung nicht, setzt aber „Regelstand zum Run" voraus.
Der Folge-Slice `slice-routing-backfill-pfad` plant die Prüfung bereits als
Fakt. Präzise Frage an den Architect: A-1.

**(f) Zählwörter.** Nachmessung grün (16 Zeilen), Stichproben ohne stehengebliebenes
Zählwort in der Spec; Träger außerhalb gemeldet und adressiert (siehe oben).
Mein Nachzählen von „Zeile 2" im Befund des Plans (3 Zeilen) stimmt.

**(g) Lastenheft unverändert.** Vertretbar (F-5, INFO): Der Diff führt keine
neue bindende Anforderung ein, alles trägt `LH-FA-CFG-008` als Verfeinerung.
Die einzige Lesart mit Vertragswirkung („Zustellung = Abruf/Abonnement") ruht
auf der Auftraggeber-Entscheidung der ADR.

## Architect-Fragen

- **A-1 (zu F-1).** Gilt die Fail-closed-Prüfung vor dem Commit eines
  Backfill-Runs (`LH-FA-CAP-009.a`) auch für den Routing-Regelstand — und wenn
  ja: mit welcher Fehlerklasse (der Plan von `slice-routing-backfill-pfad`
  nimmt `configuration` an), welcher Stand ist der Vergleichsmaßstab
  (Annahme, Start des Runs oder je Block), und was ist die Begründung dafür,
  dass das Label (kein Teil des Row Images, `ADR-0137` Teilfrage 6) diese
  Prüfung braucht, obwohl der Ausschluss-Leck-Grund (`LH-QA-SEC-004`) durch R3
  bereits gedeckt ist? Alternativen: (a) Prüfung wie bei Ausschluss- und
  Transformationsstand; (b) keine Prüfung, der Run trägt je Block den dann
  gelesenen Stand, in der Spec als Zusage benannt; (c) Stand bei Annahme
  festgeschrieben. Auch ein Verdikt „(b)" verlangt eine Zeile im Absatz
  „Ziel der Backfill-Changes".
- **A-2 (zu F-2).** Soll ein `target` außerhalb des Alphabets des Zielnamens
  auf allen Lesewegen (gRPC-Stream, SSE, `GET /changes`, `ReadChanges`) eine
  leere Antwort ohne Fehler liefern (Spec-Stand), oder nur auf `ReadChanges`
  (`ADR-0138`, als Erwartung) und auf den übrigen Wegen einen Fehler? Und:
  gilt „kein Fehler" auch für ein `target` mit U+0000?
- **A-3 (zu F-5, optional).** Reicht die Auftraggeber-Entscheidung 1 aus
  `ADR-0137` als alleinige Verankerung der Lesart „Zustellung = Abruf", oder
  braucht der Auftraggeber eine Klarstellung im Lastenheft (eigener Commit vor
  dem umsetzenden Slice, Draft-Regel des Lastenhefts)?

## Negativbefunde (geprüft, ohne Befund)

- `spec/architecture.md` (fünf Hunks): geprüft, ohne Befund außer F-4 und F-6
  (kein ADR-/Slice-/Wellen-Bezug, keine Kennung ohne Link).
- `spec/pflichtenheft.md` `SPEC-002`, `SPEC-008`, `SPEC-019` (Zählwörter,
  R1–R6, Fehlertext-Tabelle, `applied`-Absatz), `SPEC-020`/`-021`/`-022`/`-024`
  (Parameter bzw. Zusatz-Subjekt, Nachrichtenzählwörter), `SPEC-031` (Feld 7,
  Wire-Hedge), `SPEC-032` (Regelform, Auswertung, Beispiele — das Beispiel
  `order bereits vergeben: public.orders.10` entspricht der Adressregel):
  geprüft, ohne Befund außer F-1 bis F-4.
- `SPEC-029`: im Diff nicht berührt. `ADR-0138` nennt `SPEC-029` in
  `Schärft:`; die Run-Fehlerklasse steht in der Spec an `LH-FA-CAP-009.a` und
  `SPEC-008`, `SPEC-029` verweist für `error_message` generisch auf `SPEC-008`
  — kein Nachzug nötig, kein Widerspruch.
- `docs/plan/planning/in-progress/slice-routing-spec-nachzug.md`: Suchlauf-Feld
  gemessen und stimmig; kein Zustandsfeld mit Chronik; die Meldung der Träger
  außerhalb des Diffs nennt Adressen, die den Gegenstand annehmen.
- `spec/lastenheft.md`: unverändert (Diff leer); keine Änderung ohne eigenen
  Commit.
- Docker-only (`AGENTS.md` §3.1): keine Skripte oder Host-Werkzeuge im Diff.
- Traceability: beide Commit-Messages nennen `ADR-0137`, `ADR-0138`,
  `LH-FA-CFG-008`; kein Verstoß gegen das ID-Schema.

## Verdikt

- Findings: 0 HIGH, 2 MEDIUM (F-1, F-2), 2 LOW (F-3, F-4), 2 INFO (F-5, F-6).
- **Merge-blockierend: nein.** Der Diff ist reine, ungepushte Spec-Doku, keine
  Aussage steht als belegt, die die ADRs als hergeleitet führen, und die
  Stratum-Regeln (§3.4, Lastenheft unberührt) sind eingehalten. Vor dem Start
  von `slice-routing-backfill-pfad` muss A-1, vor `slice-routing-lesewege` A-2
  entschieden sein; die Antworten ziehen Zeilen in `spec/pflichtenheft.md` nach
  (Fixrunde am Spec-Text), deshalb bleibt die DoD-Zeile „Review durchgeführt"
  im Slice-Plan unberührt, bis diese Runde entschieden ist.
