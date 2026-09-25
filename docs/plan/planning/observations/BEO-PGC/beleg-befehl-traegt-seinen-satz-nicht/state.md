Deckel bei 13× (seit welle-backfill-bestand): weitere Auftreten, die vor dem Merge
vom Reviewer oder Verifier gefunden werden, Schwere ≤ LOW haben und einen bekannten
Träger-Typ treffen, bekommen keine `evidence/`-Datei, sondern stehen mit Finding-Kennung
in der Closure-Notiz des Slice (`../../README.md`, Deckel für verkörperte Einträge
ab 10×). Ausgang unverändert **verkörpert**; der Suchlauf-Befehl in der Tabellenzelle
ist mit `AGENTS.md` §3.13 §Suchform (Codeblock) und dem Werkzeug in
[`slice-harness-suchlauf-nachmessen`](../../../open/slice-harness-suchlauf-nachmessen.md)
adressiert, der Rest ist eine Lese-Handlung des Reviewers
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§3.4).

Zustand: **verkörpert** — Ausgang: **verkörpert** → `.harness/skills/reviewer.md`,
neuer HIGH-Punkt **„Beleg trägt seinen Satz nicht“** (wer einen Beleg nennt, **fährt**
ihn: den Befehl ausführen, die Adresse auflösen, die Mutation setzen, die Zählung
nachfahren) · seit welle-20. Der **Lese-Schritt der `welle-20`-Closure** hat die Regel
**geschrieben**: vier Formen in fünf Vorgängen (Befehl, Testkommentar, Adresse,
Assertion) — alle fünf Fundstellen lagen **im** Diff, der Reviewer ist damit ein
gültiger Leser.

Zähler (Datei-Anzahl unter `evidence/`, real ausgezählt): **13×**
(evidence/slice-backfill-sdk-origin.md, evidence/slice-084.md, evidence/slice-085.md, evidence/slice-091.md,
evidence/slice-093.md, evidence/slice-094.md,
evidence/slice-release-version-und-workflow.md,
evidence/slice-sdk-kotlin-pack-werkzeug.md,
evidence/slice-sdk-python-nats-stream-client-flaeche.md,
evidence/slice-backfill-change-origin.md,
evidence/slice-backfill-snapshot-reader.md,
evidence/slice-backfill-run-usecase.md,
evidence/slice-backfill-sql-administration.md) —
**Schwelle erreicht**; bereits verkörpert, der Lese-Schritt der Closure von
`welle-backfill-bestand` liest den Eintrag mit (13×). Der dreizehnte Beleg
(`evidence/slice-backfill-sdk-origin.md`, F-2) trifft die Form **Befehl**: ein
`git log -S`-Befund über Vorgänger-Commits belegte das Gegenteil des Satzes, den er
tragen sollte (die Commits hoben, der Satz sagte „hob nicht“). Der zwölfte Beleg
(`evidence/slice-backfill-sql-administration.md`, F-1 und F-3) trifft die Form **Befehl**
an zwei Trägern: ein im Repo genannter Guard-Lauf, der die zugesagten Rechte, die Funktion
und die CHECK-Menge nicht führte (die Messung lag als Wegwerf-Skript außerhalb des Repos;
die Grant-Mutation ließ den Lauf grün), und ein Suchbefehl, dessen gedruckte Zahlen zu einem
Befehl mit weniger Mustern gehörten. Ausgang bleibt **verkörpert**; die Bindung im Repo
(Lauf 5 prüft im Skript, Grant-Mutation rot) ist der Beleg der Behebung. Der elfte Beleg
(`evidence/slice-backfill-run-usecase.md`, F-6, V-2 und die Nachmessung der Closure)
trifft die Form **Befehl** in einer neuen Ausprägung: **Tabellen-Escape und Filter**. Ein
Suchbefehl in einer Tabellenzelle trug `\|` (als `-E`-Alternation und als Pipe zu `wc -l`)
und lieferte aus der Zelle kopiert 0 Treffer; ein Ausschluss-Filter (`grep -v -E` über die
`-n`-Ausgabe) schloss nach Zeileninhalt statt nach Pfad aus und ließ die Behauptung „nur
dieser Plan" falsch stehen. Ausgang bleibt **verkörpert** (der Reviewer fährt den Befehl);
die Form „ein Befehl in einer Tabellenzelle steht ohne Pipe-Zeichen, mehrere Muster als
`-e`-Argumente, Ausschlüsse als Pfad-Pathspec" trägt das Suchlauf-Feld von
`slice-backfill-run-usecase`. Der zehnte Beleg
(`evidence/slice-backfill-snapshot-reader.md`, F-3 und F-5) trifft die Form
**Befehl**: ein grüner `make coverage-gate` als Beleg für drei namentliche
Stellen, von denen er nur eine bindet, und ein unveränderter Nenner als Beleg
für „keine netzlos prüfbare Logik im Paket"; Ausgang bleibt **verkörpert**. Der neunte Beleg
(`evidence/slice-backfill-change-origin.md`, F-2) trifft die Form **Assertion**
an einem Guard-Test: der Lauf, der die Bekannt-Liste der Schema-Rollout-Wache
belegen sollte, endete auch mit entfernter Prüfung grün, weil der Abbruch aus
d-migrate kam; gefunden hat es der Reviewer durch die Mutation der
Eingabeseite, der Ausgang bleibt **verkörpert**. Der Beleg
`evidence/slice-sdk-python-nats-stream-client-flaeche.md`
(`slice-sdk-python-nats-stream-client-flaeche`, F-6 + R-1): die
Realserver-Rejection-Assertion hielt `pytest.raises(Exception)`, während
der gedruckte Marker die Ursache „token-rejected" behauptete — ein
nicht-authentifizierungsbedingter Verbindungsfehler erzeugte denselben
Beleg-Satz; Fixrunde 1 band `nats.errors.Error` +
`match="Authorization Violation"` (gegen die nats-py-Quellen
nachgemessen, s3e→s3f rot→grün am selben Server). R-1 traf dieselbe
Klasse über das committete §3.13-Suchlauf-Feld des Plans, dessen
Behandlungs-Angaben der Baum widerlegte; der Vorgang zählt zugleich bei
`BEO-PGC/arbeit-ueberholt-stehenden-traeger`. Ausgang bleibt
**verkörpert**. Der fünfte Beleg ist der erste, in dem der Beleg eine
**Assertion** ist: die Ausgabe-Hälfte von vier Fällen prüfte eine Zeichenkette,
die auch aus einem anderen Pfad kam. Der
vierte Beleg trifft die Klasse an einer **Adresse**: ein Testkopf verwies auf
Exit-Codes „im Lauf-Bericht", und dieses Artefakt existiert im Repo nicht.
Dazu ein Fall, in dem die **Wirkung** wahr war und nur die **Begründung**
falsch (eine Mutation, die über einen anderen Pfad färbt); den Ausgang weist der
**Lese-Schritt der `welle-20`-Closure** zu (Modul 6), nicht die Slice-Closure.
Die vier Belege sind vier abgeschlossene Vorgänge
(`verify-slice-084` V-1, `verify-slice-085` V-3, `review-slice-091`/`-delta`,
`review-slice-093`/`-delta`) —
**und vier verschiedene Gegenstände des Belegs**: dort ein `git diff` ohne
Pathspec, hier ein `go list` ohne das zweite Test-Datei-Feld, und mit
`slice-091` ein **Testkommentar**, der eine Mutation als rot färbend nennt, die
es nicht ist, und mit `slice-093` eine **Adresse** (ein Verweis auf einen
„Lauf-Bericht", den es nicht gibt). **Ursprung und Vorkommen sind getrennt:** der zweite
Satz stammt aus `slice-079` (`65aead2`), gefunden wurde er in `slice-085`; der
Zähler folgt den **Vorkommen**.


**Nicht zu verwechseln** mit `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(dort driftet der **Wert** gegen die Messung, die Aussage ist falsch) und mit
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (dort trägt die
**Behauptung** nicht). Hier trägt die Aussage, und ihr **Beleg** trägt sie nicht.
