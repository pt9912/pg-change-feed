# Reviewer-Skill — PG Change Feed

* Status: Accepted
* Bezug: `AGENTS.md` §3 Hard Rules und §5 Dokumentations-Regeln ·
  `harness/conventions.md` (MR-000 ID-Schema, MR-001 Pflichtenheft-Umbenennung) ·
  <!-- d-check:ignore (Kurs-/ADR-Referenzen; Anker gelten im Ziel-Repo) -->
* Gilt für: kein eigenes Make-Target — getragen vom Reviewer-Agenten
  (`.claude/agents/reviewer.md`); Reports nach `docs/reviews/`

## Kontext-Eingang (Pflicht)

Was der Reviewer *immer* mitbringt, bevor er den Diff liest:

- Diff des PR
- `spec/lastenheft.md` (für referenzierte `LH-*`-IDs)
- `spec/pflichtenheft.md` (für Verfeinerungen `LH-*.<a>` und Festlegungen
  `SPEC-*`)
- ADRs, deren ID im PR oder in der Commit-Message vorkommt
- `AGENTS.md` §"Hard Rules"
- `harness/conventions.md` (MR-000/MR-001 — ID-Schema und Datei-Namen)
- vorherige Findings am gleichen Modul (letzte ~5 PRs)

Ohne diesen Block sieht der Reviewer den Code, aber nicht *die Verträge, gegen
die er prüft*.

## Klassifikation

Jeder Anker HIGH/MEDIUM/LOW hat eine *konkrete* Liste — nicht generisch. INFO ist
bewusst kurz (Ergänzungs-Kanal, nicht Hauptkanal).

**HIGH** — eines der folgenden:
- ADR-Verstoß (Layer, Tool, Hard Rule)
- Sicherheits-Anti-Pattern (Injection, fehlende Auth-Prüfung)
- Korrektheitsfehler im *kritischen* Pfad (Persist-before-ACK-Verletzung,
  Retention löscht benötigte Changes, Lesen verändert Positionen)
- Suppression eines Gates (`#noqa`, `//nolint`, `[SuppressMessage]`) ohne ADR
- **Norm nur im Template-Kommentar** — eine Regel steht im `<!-- -->`-Block
  eines `.template.md` und nirgends sonst. Sie ist beim Adopter weg, sobald er
  die Kommentare entfernt. Kein Gate fängt das (siehe Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Template-Schichtung)
- **Kommentar trägt keine der Kommentar-Klassen** — ein Kommentar in Code, Config
  oder Skript beschreibt die verworfene Alternative („Ohne X wäre …"), einen
  abwesenden Text („früher stand hier …") oder bricht mitten im Satz ab, weil
  eine Teilersetzung den Rest stehen ließ. Kein Gate fängt das (siehe
  Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein Kommentar trägt).
  *Skopus:* der Punkt gilt für Code, Konfiguration und Skripte **einschließlich**
  Tests und Runner (`tools/harness/*.sh`); ein Kommentar, der ein Vorher/Nachher
  andeutet („schließt die beiden zuvor fehlenden …“) oder eine verworfene
  Alternative im Konjunktiv nennt („würde diese verzögern“), gehört hierher,
  nicht unter INFO; der Punkt „Slice-/Wellen-Chronik“ bleibt auf
  Produktionscode-Pfade begrenzt. *Zusage:* ein Kommentar der Klasse **Zusage**,
  der ein Verhalten zusichert (Fehlerpfad, Ausgang, Rückgabe), das der Code an
  dieser Stelle nicht trägt; Probe ist das Nachfahren des zugesagten Pfads im
  Code. Liegt das Verhalten in einem anderen Slice oder Paket, trägt der
  Kommentar einen **Rang-Zeiger** darauf. Herkunft:
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (3×),
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (3×) · seit
  welle-backfill-bestand
- **Slice-/Wellen-Chronik in Produktionscode-Kommentar** — ein Godoc- oder
  Inline-Kommentar über einem **Produktionscode**-Pfad (Funktion, Typ, Datei —
  nicht ein `Test*`-Godoc) begründet eine Aussage mit einer Slice-/
  Wellen-Nummer (`slice-<NNN>`, `welle-<NN>`) oder impliziter
  Vorher/Nachher-Sprache, statt mit `ADR-*`/`LH-*` oder dem Herkunfts-Anker
  `· seit slice-<NNN>` (`AGENTS.md` §3.7). Abgrenzung zur **zulässigen**
  Testfall-Provenienz („`TestXyz` trägt/deckt … aus `review-slice-NNN.md`
  F-x“ — Subjekt ist der Test, nicht der Produktionscode-Pfad, siehe
  dem Architect-Verdikt zur Slice-Chronik in Code-Kommentaren):
  Probe ist das **Satzsubjekt** — die Funktion/der Code-Pfad (Chronik,
  unzulässig) oder der Testfall (Provenienz, zulässig). Kein Gate fängt das
  (repo-weiter Textmuster-Sensor geprüft und verworfen, s.o.). Erstes
  benanntes Auftreten als eigener HIGH-Punkt: Review F-1
  (Review zu `slice-052`, 2026-09-13) — vierter gezählter Beleg von
  `BEO-PGC/slice-chronik-in-code-kommentar`, der Architect-Verdikt-Nachtrag
  zur Slice-Chronik in Code-Kommentaren (4. Auftreten).
- **Handbuch-Versionshistorie nicht fortgeschrieben** — ein Diff ändert
  `docs/user/benutzerhandbuch.md` inhaltlich (neuer Abschnitt, neue
  Umgebungsvariable, geänderte Beschreibung), ohne im selben Diff den
  `Version:`-Kopf hochzuzählen **und** eine neue Zeile in
  `### Änderungshistorie` zu ergänzen. Dies ist die tragende
  Verteidigungslinie: Die Implementer-Selbstprüfung
  (`.claude/commands/implement-slice.md` Schritt 17) läuft im selben
  Kontext, der die Doku-Änderung geschrieben hat, und hat die Klasse real
  dreimal übersehen (`slice-045`, `-046`, `-053` — jedes Mal erst bei
  einem späteren Slice bemerkt), bevor sie geschärft wurde — dieselbe
  Struktur wie beim Chronik-Fall oben. Erstes benanntes Auftreten als
  eigener HIGH-Punkt: der Architect-Verdikt zur übersprungenen
  Handbuch-Versionshistorie
  (3× `BEO-PGC/handbuch-versionshistorie-uebersprungen`) · seit slice-053.
- **Neue Betreiber-Oberfläche ohne Handbuch-Zug** — ein Diff führt eine neue
  Betreiber-Oberfläche ein (eine `CDC_*`-Umgebungsvariable des
  Feed-Containers, eine administrative `cdc.*`-SQL-Funktion, eine
  Horch-Adresse oder einen Endpunkt), ohne dass derselbe Diff
  `docs/user/benutzerhandbuch.md` inhaltlich mitzieht (§5
  „Umgebungsvariablen des Feed-Containers", §4 „Aufgaben") **und** ohne einen
  benannten Aufschub mit Adresse (Folge-Slice-ID). Die
  Handbuch-Versionshistorie-Regel oben greift erst, wenn das Handbuch
  angefasst wird — wer es gar nicht anfasst, löst sie nicht aus; hier ist
  der Diff selbst die Fundstelle. Dies ist die tragende Verteidigungslinie:
  Die Implementer-Selbstprüfung (`.claude/commands/implement-slice.md`
  Schritt 17) läuft im selben Kontext, der die Oberfläche eingeführt hat —
  dieselbe Struktur wie beim Versionshistorie-Fall oben. Kein Gate fängt die
  Klasse: `docs-check` prüft Referenzen, nicht Vollständigkeit, und ein
  Sensor gegen die ENV-Variablen im Code bräuchte eine Semantik-Entscheidung,
  welche Variablen „Betreiber-Oberfläche" sind (geprüft und verworfen).
  Herkunft: `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (3×, `slice-059`/`-066`/`-069`) · seit slice-077.
- **Zustandsfeld trägt Chronik** — eine `Stand`-/`Status`-Zelle (Roadmap,
  Beobachtungs-Register, Meilenstein) erzählt, wie der Zustand entstand, statt
  Zustand und Beleg als Anker zu nennen; oder ein Drift-Log protokolliert
  Schließungen und erreichte Meilensteine. Kein Gate fängt das (siehe
  Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein Kommentar trägt,
  *Dieselbe Regel für Zustandsfelder*)
- **Zahl im Träger ohne Ursprung — oder gegen die Messung driftend** — ein
  Doku-Träger (Sensor-Doku, ADR, Slice-Plan, README-Tabelle, Bericht) nennt
  eine Zahl über den Gegenstand, ohne ihren **Ursprung** zu tragen (gemessen ·
  übernommen · abgeleitet) und, wo sie eine Messung ist, ohne den **Lauf**;
  oder ein übernommener Wert driftet gegen die eigene Messung. Kein Gate fängt
  das: es gibt **keinen** Sensor, und der Verzicht ist entschieden — eine
  Formpflicht auf Prosa erzeugte Pflichterfüllung (`AGENTS.md` §3.12, siehe
  [`ADR-0083`](../../docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
  §Die benannte Grenze). Die Probe ist das **Nachmessen**, nicht das Lesen der
  Form. Träger **außerhalb** des Diffs haben nur einen Leser — die Messung —
  und bleiben INFO. Herkunft:
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (4×,
  `slice-081`/`-084`/`-085`/`-088` — alle vier vom Reviewer durch eigenes
  Nachmessen gefunden) · seit slice-089.
- **Beleg trägt seinen Satz nicht** — ein Träger nennt einen **Beleg** als
  Stütze einer Aussage — einen Befehl, eine Abfrage, einen Pfad, eine Adresse,
  eine Mutationsangabe, eine Assertion oder einen Verweis auf eine Stelle
  eines anderen Dokuments (eine Abschnittsnummer, eine ADR-Festlegung, eine
  Slice-/Welle-Kennung) —, und der genannte Beleg trägt die Aussage **nicht**:
  er misst etwas anderes, zählt etwas anderes oder liefert ein anderes
  Ergebnis als das behauptete. *Wer einen Beleg nennt, fährt ihn: den
  genannten Befehl ausführen, die genannte Adresse auflösen, die genannte
  Mutation setzen, die genannte Zählung nachfahren — oder schlägt die
  genannte Stelle im Original auf, nicht in einer Zusammenfassung (die
  Kopfzeile eines Berichts, die eigene Erinnerung). Die Aussage darf dabei
  wahr sein — geprüft wird die Stütze, nicht der Satz.* Die Probe ist
  mechanisch und billig und wird genau deshalb übersehen: der Satz ist
  plausibel, erst der ausgeführte oder aufgeschlagene Beleg zeigt etwas
  anderes. Abgrenzung zu **„Zahl im Träger ohne Ursprung"** (dort trägt die
  **Aussage** nicht, hier trägt sie und ihre Stütze nicht) und zu **„Zusage
  ohne Bindung an ihre Eingabeseite"** (dort ist die Zusage nicht rot zu
  färben, hier ist der genannte Weg zur Prüfung der falsche). Kein Gate fängt
  das: ein Beleg-Befehl läuft netzlos und liefert ein Ergebnis, eine
  zitierte Stelle lässt sich aufschlagen — ob das Ergebnis oder die Stelle
  **den Satz stützt**, ist eine Lese-Handlung, kein mechanischer Vergleich.
  Herkunft: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (5×,
  `slice-084`/`-085`/`-091`/`-093`/`-094`; die Fundstellen sind ein `git diff`
  ohne Pathspec, ein `go list` ohne das zweite Test-Datei-Feld, ein
  Testkommentar, eine Adresse auf ein Artefakt, das es nicht gibt, und eine
  Assertion, die auch aus einem anderen Pfad hält) · seit welle-20;
  `BEO-PGC/zitat-nennt-die-falsche-stelle` (3×,
  `slice-090`/`slice-102`/`slice-d-check-tracked-modul` — dieselbe Fehlerklasse,
  über drei Vorkommen uneinheitlich als MEDIUM/INFO/HIGH klassifiziert, bevor
  sie hier als Nachbar-Form explizit gefasst wurde) · seit welle-d-check.
- **Zusage ohne Bindung an ihre Eingabeseite — „grün ohne Aussage"** — eine
  Zusage (Test, Negativtest, Filter-/Limit-Prüfung) ist **vorhanden** und läuft
  grün, kann aber an ihrer **Eingabeseite** nicht rot werden: mutiert wurde nur
  die Ausgabeseite (der Fake, der Rückgabewert), nicht der **Eingabewert**.
  *Eine Zusage ist nur dann gebunden, wenn der Test an ihrer Eingabeseite rot
  werden kann: mutiere den Eingabewert, nicht nur die Ausgabeseite. Wer nur den
  Fake oder den Rückgabewert mutiert, prüft den Fake — die Aussage bleibt grün,
  egal was der Adapter mit der Eingabe tut. Fehlt die Mutation der Eingabeseite,
  ist die Zusage grün ohne Aussage: ein Befund, kein Formfehler.* Abgrenzung zur
  MEDIUM-Klasse „fehlende Negativtests bei neuem öffentlichem Vertrag": dort
  **fehlt** die Abdeckung, hier steht eine vorhandene Zusage ohne Bindung — die
  Kategorie entscheidet über die Fixrunde. Kein Gate fängt das: ein
  Mutations-Harness gibt es in diesem Repo nicht, die Prüfung **ist** die
  Mutation. Herkunft: `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (4×,
  `slice-083`/`-086`/`-087`/`-088`; in drei der vier Fälle fand der Reviewer die
  Klasse durch Mutieren der Eingabeseite), der Architect-Verdikt zur Zusage
  ohne Bindung an ihre Eingabeseite
  · seit slice-089. Die Träger-Seite derselben Regel steht in
  `.claude/commands/implement-slice.md` Schritt 19.
- **Form-Vorbild-Kopie trägt ein sprachgebrochenes Wortfragment weiter** —
  ein Form-Teil stammt aus der deutschen Plan-Prosa (frisch formuliert oder
  aus einem Vorgänger-Plan übernommen) oder aus einem Formvorbild (ein
  SDK-Sprachpaket aus seinem Vorgänger-Paket, ein gespiegelter Runner-/
  Dockerfile-/Make-Target-Kopf) und trägt ein Wortfragment, das die Sprache
  seines Trägers bricht: ein unübersetztes deutsches Fachwort mitten im
  englischen Satz (README-Formulierung, Klassen-Doc-Kommentar) oder ein
  gebrochenes Hybrid-Fragment in deutscher Plan-Prosa. *Wer einen solchen
  Form-Teil schreibt oder liest, sichtet ihn einzeln auf Sprachreinheit —
  beim Kopieren mit dem Formvorbild selbst in der Hand; die wortgleiche
  Übereinstimmung mit dem Vorbild ist kein Beleg für die Sprachreinheit,
  sondern genau der Weg, auf dem das Fragment weiterwandert. Die
  §8-Deklaration („dieser Slice kopiert kein Form-Vorbild") sichtet ihre
  eigene Prosa mit.* Abgrenzung: die Sprache des Trägers entscheidet, nicht
  das Wort — ein deutsches Fachwort in einem deutschsprachigen Kommentar
  (Runner-Köpfe in `tools/harness/`, Build-Datei-Kommentare der SDK-Bäume)
  ist an seinem Ort und kein Befund dieser Klasse; die orthografische
  Form-Familie ihrer byte-identischen Kopien („FlaecheN") bleibt die
  Nachbar-Form der jeweiligen Reviews, kein Zähler dieser Klasse. Kein Gate
  fängt das: es gibt keinen Sensor für Prosa-Orthografie, und eine
  Formpflicht auf Prosa erzeugt Pflichterfüllung (dieselbe Grenze wie beim
  Zahl-im-Träger-Punkt); die verfügbare Falsifikation ist das Lesen des
  Form-Teils gegen sein Vorbild. Herkunft:
  `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (3×,
  `slice-sdk-csharp-projektgeruest`/`-kotlin-projektgeruest`/
  `-python-http-reale2e` — je LOW befundet; die beiden Code-Kommentar-
  Fragmente „unstrittige" stehen bis heute im Bestand:
  `grep -rni unstrittige sdks/` → zwei Treffer,
  `PgChangeFeedClientOptions.cs:7` und `PgChangeFeedClientOptions.kt:9`,
  gemessen beim Architect-Verdikt zur Form-Vorbild-Kopie — die als LOW an
  den „nächsten Slice, der dieselbe Datei berührt" weitergereichte
  Korrektur hat sie nie gezogen) und
  `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (3×, derselbe
  Mechanismus an den README-Sätzen aller drei Sprachpakete, README-Treffer
  heute gezogen), der Architect-Verdikt zur Form-Vorbild-Kopie und der
  Sprachreinheit je Form-Teil · seit slice-sdk-python-http-reale2e.
- **Traceability-/ID-Schema-Verstoß** — Commit oder PR nennt keine
  `LH-*`- oder `ADR-*`-Kennung; oder eine Kennung nutzt ein Präfix, das MR-000
  nicht deklariert (`LH-FA/QA-<BEREICH>-<NNN>`, `SPEC-<NNN>`, `ARC-<NNN>`,
  `ADR-<NNNN>`, `CO-<NNN>`, `slice-<NNN>`, `MR-<NNN>`, `BEO-<KUERZEL>/<slug>`,
  `RC-<NNN>`). Erstes Auftreten dieser Klasse: Review F-2 (2026-09-09,
  `PH-*`/`TST-*` im Traceability-Beispiel).
- **Spec-Stratum-Verstoß** — das Technik-Stratum erweitert, wo es nur
  präzisieren darf („präzisieren ja, erweitern nie“): das Pflichtenheft führt
  eine neue bindende Anforderung ein, statt eine bestehende `LH-*`-ID zu
  schärfen; oder das Lastenheft wird geändert, ohne dass die Änderung in
  einem eigenen Commit vor dem umsetzenden Slice liegt. Erstes Auftreten
  dieser Klasse: Review F-3 (2026-09-09, Zurückstellung bindender
  Anforderungen im Sammelabschnitt).
- **Docker-only-Verstoß** — ein Build-/Test-/Betriebs-Skript installiert ein
  lokales Toolchain (venv, SDK, Paketmanager) statt über `make` zu laufen
  (`AGENTS.md` §3.1).
- **Zwei-Quellen-Drift** — derselbe Zustand wird in zwei Dateien geführt,
  ohne dass der Gewinner deklariert ist (z. B. Anforderungstext doppelt in
  Lastenheft und Pflichtenheft; Zustand in Verzeichnis *und* Status-Feld).

**MEDIUM** — eines der folgenden:
- unklare Fehlerbehandlung am Rand des Spec-Bereichs
- fehlende Negativtests bei neuem öffentlichem Vertrag
- Wiederholung eines Musters, das schon zweimal LOW war
- **Delegation ohne Entsprechung** — das Lastenheft delegiert eine Festlegung
  an das Pflichtenheft (`LH-QA-PER-002`, `LH-QA-PER-004`, `LH-QA-POR-001`),
  dort steht sie nicht, auch nicht als offene Festlegung mit Schlusspunkt.
  Erstes Auftreten: Review F-4 (2026-09-09).
- **GWT-Pfad unvollständig** — ein Akzeptanzkriterium fehlt ganz (`—` ohne
  die `—`-Regel aus Lastenheft §3), bricht das Given/When/Then-Muster oder
  verfehlt die drei Pfade Happy/Boundary/Negative. Erstes Auftreten:
  Review F-5 (2026-09-09).
- **MVP-Kennzeichnung inkonsistent** — `MVP: ja`-Marker und die
  MVP-Abnahme-Mapping-Tabelle (Lastenheft §1) weichen voneinander ab.
  Erstes Auftreten: Review F-1 (2026-09-09).
- **Nachzug widerspricht dem Nachbarn im selben Träger** — ein Diff ergänzt einen
  Absatz, eine Tabellenzeile oder einen Kommentarblock, der eine Aussage ersetzt
  oder einschränkt, ohne dass der Gegen-Absatz **desselben** Dokuments, Blocks
  oder Abschnitts angepasst oder auf den neuen verwiesen wird (zwei Aussagen,
  keine verweist auf die andere). Probe: den Kontext um jede hinzugefügte Zeile
  lesen (`git diff -U20`), nicht nur die Zeilen. Liegt der Träger **außerhalb**
  des Diffs, ist es der Träger-Nachzug von `AGENTS.md` §3.13 (INFO/LOW, Meldung
  an den Planner). Kein Gate fängt das: ob zwei Aussagen sich widersprechen, ist
  eine Lese-Handlung. Herkunft: `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
  (8×) · seit welle-backfill-bestand.

**LOW** — stilistisch unschön ohne semantische Auswirkung, einmalige Tippfehler,
unbenutzte Imports.

**INFO** — Hinweis ohne erwartete Aktion (z. B. „diese Stelle hat ein passendes
ArchUnit-Pendant, das du nicht kennst“; „Metrik deckt die Anforderung nur
indirekt ab“ — erstes Auftreten: Review F-9, 2026-09-09).

> **Pflicht beim Ausfüllen (Modul 10 §Übungen):** Die HIGH-Liste muss mindestens
> *zwei* repo-spezifische Regeln nennen, die ein generischer Skill nicht abdeckt.
> Ist der Skill ohne sie, ist er noch nicht scharf genug — dann kommt bei einem
> Lauf auf einem realen Diff keines deiner Repo-HIGHs zur Anwendung.
> **Erfüllt seit 2026-09-09:** vier repo-spezifische HIGH-Regeln
> (Traceability/ID-Schema, Spec-Stratum, Docker-only, Zwei-Quellen-Drift).

## Was dieser Skill NICHT macht

- Keine Lösungsvorschläge („schreib das so“) — Reviewer kategorisiert,
  Implementer entscheidet.
- Kein Refactoring-Vorschlag, der über den Diff hinausgeht.
- Keine Verifikation gegen DoD — das ist Verifier-Aufgabe (Modul 11).
- Keine Validation gegen reale Bedürfnisse — das ist Validator-Aufgabe.

Wenn etwas auffällt, das in diese Kategorien gehört: ein INFO-Finding mit Verweis
auf die zuständige Rolle.

## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: ADR-ID, `LH-*`-ID, Hard-Rule-Name oder „Maintainability“
- `pfad`: Datei:Zeile
- `befund`: 1–2 Sätze, beobachtbar, ohne Lösungsvorschlag
- `verifizierbar`: ja/nein — gibt es einen Gate-Lauf, der es bestätigen würde?
- `klasse`: stabile Kurz-Bezeichnung des Fehlermusters, z. B. „Delegation ohne
  Entsprechung“ — speist den Steering-Loop-Zähler (siehe §Pflege)

Zusätzlich am Ende: eine Zeile „geprüft, ohne Befund“ pro betrachtetem
Verzeichnis (Negativbefund-Zeile — sonst ist „keine Findings“ nicht von „nicht
geprüft“ unterscheidbar). Report-Gerüst für den ganzen Lauf:
`docs/reviews/review-report.template.md`, ein Report pro Lauf, Folgeläufe als
neue Datei statt Überschreibung.

## DoD-Checkbox-Nachzug ohne Fixrunde

Kommt dein eigenes **Verdikt** zu dem Schluss, dass **keine Fixrunde am
Implementer** nötig ist — 0 HIGH, oder alle HIGH/MEDIUM/LOW-Findings werden
ohne einen Reviewer→Implementer-Rückgabe-Pfeil weitergereicht (z. B. direkt
vom Planner behoben) —, zieh die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im betroffenen Slice-Plan **selbst** auf `[x]`
nach, mit Verweis auf den eigenen Report-Pfad, **im selben Commit**, der den
Report anlegt.

**Warum hier und nicht beim Implementer:** Der reguläre Nachzug-Mechanismus
(`.claude/commands/implement-slice.md` Schritt 18/21,
`BEO-PGC/dod-checkbox-nachzug`) hängt an einem *zweiten Implementer-Lauf* —
Schritt 18 kommt zu früh (vor dem Review), Schritt 21 greift nur bei einer
echten Fixrunde. Bleibt die Fixrunde aus, gibt es keinen Implementer-Lauf
mehr, an den sich der Nachzug hängen könnte. Du bist die einzige Rolle, die
im richtigen Moment — beim Schreiben deines eigenen Verdikts — bereits weiß,
ob eine Fixrunde kommt oder nicht.

**Grenze:** Nur diese eine Checkbox. Kein anderer DoD-Punkt, keine
Verifikations-Substanz (das bleibt Verifier-Aufgabe, Modul 11) — du
bestätigst ausschließlich die Tatsache, dass dein eigener, abgeschlossener
Arbeitsschritt stattgefunden hat. Braucht der Slice eine Fixrunde, bleibt die
Checkbox offen; sie wird dann regulär bei Schritt 21 des
Implementer-Workflows nachgezogen.

Herkunft: `BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde` (3×,
`slice-045`/`slice-046`/`slice-047`), der Architect-Verdikt zur
DoD-Checkbox „Review durchgeführt" ohne Fixrunde
· seit slice-047.

## Pflege (Steering-Loop)

Bei dreimaligem Auftreten desselben Findings:

- ist die Kategorie noch richtig? → Klassifikation schärfen
- gibt es einen ADR/`AGENTS.md`-Eintrag, der das verhindert hätte?
  → Folge-ADR oder `AGENTS.md`-Update
- gibt es eine Fitness Function, die das prüfen würde? → Modul 13, Gate hinzufügen

Diese Skill-Datei wird **nicht** überschrieben, sondern versioniert
(ADR-Hard-Rule, Modul 4). Erste Schärfung 2026-09-09 aus dem
Eröffnungs-Review zu Lastenheft und Pflichtenheft.