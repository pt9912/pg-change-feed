Zustand: **verkörpert** — Ausgang: **verkörpert** → `.harness/skills/reviewer.md`,
neuer HIGH-Punkt **„Beleg trägt seinen Satz nicht“** (wer einen Beleg nennt, **fährt**
ihn: den Befehl ausführen, die Adresse auflösen, die Mutation setzen, die Zählung
nachfahren) · seit welle-20. Der **Lese-Schritt der `welle-20`-Closure** hat die Regel
**geschrieben**: vier Formen in fünf Vorgängen (Befehl, Testkommentar, Adresse,
Assertion) — alle fünf Fundstellen lagen **im** Diff, der Reviewer ist damit ein
gültiger Leser.

Zähler (Datei-Anzahl unter `evidence/`, real ausgezählt): **9×**
(evidence/slice-084.md, evidence/slice-085.md, evidence/slice-091.md,
evidence/slice-093.md, evidence/slice-094.md,
evidence/slice-release-version-und-workflow.md,
evidence/slice-sdk-kotlin-pack-werkzeug.md,
evidence/slice-sdk-python-nats-stream-client-flaeche.md,
evidence/slice-backfill-change-origin.md) —
**Schwelle erreicht**. Der neunte Beleg
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
