Zustand: **verkörpert** — Ausgang: **verkörpert** → geschärfte
Selbstprüf-Instruktion (kein mechanischer Sensor — geprüft und verworfen,
siehe Architect-Verdikt) als Pflicht-Ergänzung in
`.claude/commands/implement-slice.md` Schritt 20: diff-skopierter
`grep`-Kandidatenlauf gegen die in diesem Lauf geänderten
`.go`-/`tools/schema/*.sql`-Dateien, plus die explizite
Unterscheidungsprobe Testfall-Provenienz vs. Produktionsverhalten-Chronik.
Architect-Verdikt zur Slice-Chronik in Code-Kommentaren
(empirischer Befund: über 30 bereits gemergte, akzeptierte
`slice-\d+`/`welle-\d+`-Zitate in `_test.go`-Godoc-Kommentaren als
etablierte Testfall-Provenienz — strukturell ununterscheidbar von den drei
gezählten Verstößen, ein repo-weiter Textmuster-Sensor würde entweder
diese Konvention flächendeckend fehlalarmieren oder eine
Satz-Subjekt-Erkennung brauchen, die kein „kleines, dependency-freies
Skript" mehr ist). Wellenloser Architect-Zug (kein Slice/keine Welle
trägt diese Verkörperung) — der Herkunfts-Anker ist dieser Architect-Zug
selbst (obiger Verdikt-Pfad), analog zu
`BEO-PGC/architect-verdikt-ablageort-uneinheitlich`s wellenloser
Behebung direkt aus einer Nutzerfrage.

Zähler (abgeleitet): 6× (dem Beleg zum Review zu `slice-041`,
dem Beleg zum Review-Report zur Fixrunde von `slice-041`, dem Beleg zum
Review zu `slice-044`,
evidence/slice-052.md, evidence/slice-060.md, evidence/slice-066.md) —
Schwelle erreicht, Ausgang im Lese-Schritt dieses Architect-Zugs
zugewiesen (wellenlos, siehe Modul 6 „Träger im Repo ohne Wellen": der
Lese-Schritt läuft normalerweise in der Slice-Closure; hier lief er als
eigener, von der Nutzerin ausgelöster Architect-Zug außerhalb einer
Slice-Closure — dieselbe Lücke, die `BEO-PGC/dod-checkbox-nachzug-architect-pfad`
für den DoD-Pfad bereits benennt). Die historische `slice-018`-Korrektur
(siehe `observation.md`) zählt laut Präzedenzfall `BEO-PGC/plan-nachzug`
weiterhin nicht mit (vor der Registrierung).

Restrisiko, benannt statt gezählt: Die geschärfte Instruktion behebt die
beobachtete Enumerations-Lücke (eine von mehreren Stellen übersehen),
nicht die schwerere, im Verdikt ausdrücklich unmechanisierte
Klassifikations-Frage (Testfall-Provenienz vs. Chronik). Tritt die Klasse
trotz der geschärften Instruktion ein viertes Mal auf, ist das ein
Signal, dass Enumeration allein nicht trägt — neue Beobachtung oder
Zähler-Fortschreibung, Urteil beim nächsten Lese-Schritt.
Vorgezogene Antwort auf das Restrisiko (4. Beleg, `slice-052` F-1, Review
zu `slice-052` samt Fixrunde; formal
noch nicht als `evidence/slice-052.md` gezählt — das ist reguläre
Slice-Closure-Arbeit des Planners, `slice-052` liegt noch in
`in-progress/`): Architect-Verdikt-Nachtrag zur Slice-Chronik in
Code-Kommentaren (4x).
Diagnose: kein neuer Enumerations-Fall (das bestehende Pattern hätte den
Fund getroffen), sondern Bestätigung, dass Schritt 20 strukturell nur die
erste, nicht die tragende Verteidigungslinie sein kann (Modul 8 §Kernidee)
— der unabhängige Reviewer hat 4/4 Fälle vor Merge gefangen, kein
Hard-Rule-Verstoß hat je `main` erreicht. Verkörpert statt eines neuen
Sensors: eigener benannter HIGH-Punkt in `.harness/skills/reviewer.md`
(„Slice-/Wellen-Chronik in Produktionscode-Kommentar") und eine
Grenz-Klarstellung in `.claude/commands/implement-slice.md` Schritt 20.
Für den Lese-Schritt bei `slice-052`s Closure ist der Ausgang damit
vorweggenommen: **verkörpert** (erneut), Herkunfts-Anker `seit slice-052`
auf diesen Nachtrag — keine weitere Architect-Eskalation nötig, nur
Zähler und `evidence/slice-052.md` nachtragen.

Belege 4–6 (`slice-052`, `slice-060`, `slice-066`) bestätigen die
Verdikt-Diagnose: in allen drei Fällen fing der unabhängige Reviewer den
Fund vor dem Merge, kein Hard-Rule-Verstoß hat `main` erreicht. Der
Ausgang bleibt **verkörpert** — die tragende Verteidigungslinie ist der
Reviewer-HIGH-Punkt, nicht die Selbstprüfung des schreibenden Laufs. Der
`slice-066`-Beleg zeigt zusätzlich den Nutzen des datei-skopierten
Enumerationslaufs: er fand vier weitere Produktionscode-Stellen derselben
Klasse in `postgresstorage/administrationrequest.go` und
`bootstrap/wiring.go`, die über den ursprünglichen Review-Befund
hinausgingen.

Beleg 7 (`slice-schema-rollout-zentrale-idempotenz-wache`, evidence-Datei):
zwei Produktionscode-Stellen (`Makefile`, `guard.go`) — Reviewer fing beide
vor Merge, bestätigt die Diagnose erneut. **Zusätzlich eine neue Grenze
real gefunden:** eine dritte, wortidentische Instanz in `harness/README.md`
(Doku-Prosa, kein Code-Kommentar) lag außerhalb des wörtlichen Skopus des
Reviewer-HIGH-Punkts („Produktionscode-Kommentar") und wurde erst vom
Verifier gefunden, nicht vom Reviewer. Ausgang bleibt **verkörpert** für die
Produktionscode-Hälfte; die Doku-Hälfte ist ein offener Randbefund für den
nächsten Lese-Schritt (Skopus-Erweiterung des HIGH-Punkts erwägen), kein
eigener Zähler-Eintrag, da bislang nur 1× beobachtet.

Beleg 8 (`slice-generated-sync-tar-export`, evidence-Datei, Planner-Closure
2026-09-21): erstes Vorkommen der Klasse in einem **Shell-Skript**
(`tools/harness/generated-sync.sh`, Kopf-Kommentar) — Datei-Typ-Erweiterung
ggü. der bisherigen Go-/SQL-/`schema.yaml`-Fläche, keine neue Klasse.
Reviewer fing die Stelle vor Merge (F-1, HIGH), Fixrunde real geprüft
(0 HIGH danach) — Ausgang bleibt **verkörpert**, neunter Beleg derselben
Diagnose. **Zusätzlich zweiter Beleg für die in Beleg 7 benannte
Doku-Prosa-Grenze:** dieselbe Chronik-Struktur stand in
`harness/sensors/generated-sync.md`, diesmal bereits vom **Reviewer selbst**
gefunden (F-2, INFO, nicht erst vom Verifier), blieb über die Fixrunde
unverändert bewusst INFO. Die Zwei-Beleg-Schwelle für die Doku-Prosa-
Variante ist damit erreicht — noch kein eigener Registereintrag (unter der
Drei-Beleg-Schwelle, `harness/README.md` §Traceability rules Analogie),
aber ein verstärktes Signal für den nächsten Lese-Schritt.

Beleg 9 (`slice-sdk-kotlin-nats-stream-client-flaeche`, evidence-Datei,
Planner-Closure 2026-09-22): `build.gradle.kts`s Version-Kommentar trug ein
Arrow-Muster (`0.1.0 -> 0.2.0`) mit Slice-Kennung statt `ADR-*`/`LH-*`-Anker
— dieselbe Instanz wie beim analogen C#-NATS-Sibling-Slice
(`slice-sdk-csharp-nats-stream-client-flaeche`, für den keine eigene
Evidenzdatei angelegt wurde — Lücke, hier nur benannt, nicht rückwirkend
geschlossen). Reviewer fing die Stelle vor Merge (F-1, HIGH), Fixrunde real
geprüft, vierfache unabhängige Bestätigung (Reviewer, Fixrunden-Reviewer,
Verifier). Ausgang bleibt **verkörpert**, neunte Evidenzdatei in diesem
Verzeichnis (real ausgezählt).

**Benannt, nicht gezählt — außerhalb eines Slice.** Im Kommentarblock über dem
Target `schema-rollout` im `Makefile` und in den `proto`-Kommentaren standen
acht Zeilen mit Slice-Bezug als Begründung („seit slice-011/012/036/066",
„slice-015 hat das …", „slice-016 …", „seit slice-104"); der Nutzer meldete sie
nach der Verifikation von `slice-backfill-change-origin`, der Zug
`c54f873a` bereinigte sie (Kommentar auf den Ist-Vertrag gekürzt, der Vertrag
steht in `harness/targets/schema-rollout.md`). Ein Vorgang außerhalb eines
Slice-Diffs, gefunden ohne Reviewer, deshalb nach dem Präzedenzfall der
`slice-018`-Korrektur (siehe `observation.md`) **keine** `evidence/`-Datei und
kein Zähler-Beitrag. Derselbe Suchlauf gegen `HEAD` (`c54f873a`, gemessen bei der
Closure von `slice-backfill-change-origin`) findet weitere Kandidaten: mit
`grep -cE 'slice-[0-9a-z]'` **19** Zeilen in `harness/mk/*.mk` (sechs Dateien),
**32** in `tools/**/*.sh` (16 Dateien) und **57** in `harness/sensors/*.md`
(vier Dateien). Die Treffer sind Kandidaten, nicht Verstöße: Herkunfts-Anker der
Form `· seit slice-<NNN>` und Testfall-Provenienz sind nach `AGENTS.md` §3.7
zulässig, die Klassifikation je Treffer (Zustand des Codes gegen Testfall) ist
nicht geschehen. Ein Aufräum-Zug über diese Flächen ist nicht angelegt.

**Benachbart, kein Widerspruch zum Verdikt.** `make kommentar-kennungen`
(`harness/sensors/kommentar-kennungen.md`) zählt verschiedene Kennungen je Go-Kommentarblock und
liest weder ein Satz-Subjekt noch eine Slice-/Wellen-Nummer; es unterscheidet Testfall-Provenienz
nicht von Produktionsverhalten-Chronik. Der Chronik-Kandidatenlauf in
`.claude/commands/implement-slice.md` Schritt 20 bleibt der Weg dieser Klasse. Die
Kennungsdichte ist ein eigener Eintrag: `BEO-PGC/kommentar-herkunft-als-kette`.
