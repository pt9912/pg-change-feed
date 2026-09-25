Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.12 Instanz A und
`.harness/skills/reviewer.md` (HIGH-Punkt „Zahl im Träger ohne Ursprung — oder gegen
die Messung driftend“) · seit slice-089. Der **Lese-Schritt der `welle-20`-Closure**
hat den Ausgang zugewiesen: der Träger trägt die **drei** Vorkommen nach der
Verkörperung (`slice-089`, `-090`, `-092`) — alle drei von **Lesern** gefunden,
keines von einem Sensor.

**Die Form ist entschieden — und sie ist der Ausgangs-Kandidat.** Der
`slice-085`-Vorgang hat sie formuliert und in zwei Dokumenten angewandt:

> Jede Zahl eines Doku-Trägers trägt ihren **Ursprung**; ist sie eine Messung,
> trägt sie zusätzlich den **Zeitpunkt** (Lauf/Stand). Der **Nenner** hängt am
> Code-Stand, die **gedeckte Zahl** am Lauf.

Das ist `ADR-0078`s „gemessen / übernommen / **abgeleitet**" **plus die
`wann`-Hälfte** — die Ergänzung, die dieser Vorgang erzwungen hat: der
Implementer musste den Ursprung alter Zahlen per `git log -S` rekonstruieren,
weil ihr Träger bei der Geburt uneindeutig datiert war. Sein Satz dafür: *„der
Schreiber setzt den Zeitpunkt, der Leser kann ihn nicht erraten."*

**Sie ist noch nicht vollständig durchgesetzt** (Verifikation V-1:
`coverage-gate.md` §Grenze Punkt 1 und zwei §Ausgabe-Abschnitte führen bewegliche
Werte ohne Zeitpunkt). **Genau das ist der Gegenstand des Ausgangs:** die Regel
**verkörpern** (Träger-Kandidat: `AGENTS.md` §3.7 als Geschwister-Ort, oder der
Kopf des ADR-Index) **und** die verbleibenden Träger darauf prüfen.

Zähler (abgeleitet): **20×** (evidence/slice-backfill-slot-leerlauf-bestaetigung.md —
ein Lauf ohne auflösbaren Träger im Handbuch, eine Trefferzahl im Suchlauf-Feld und ein
Zählwort im selben Feld, das den Selbstverweis der Plan-Datei mitzählt, Review F-4, F-5 und
Verifikation V-1; evidence/slice-backfill-bench-richtgroesse.md — Zahlen
mit nicht auflösbarem Ursprung, eine Zwischenstands-Zahl und zwei Zählungen im
Suchlauf-Feld, Review F-2, F-9 bis F-11 und Verifikation V-2, V-5, V-6;
evidence/slice-backfill-e2e.md — eine Trefferzahl im
Suchlauf-Feld des Plans, Review F-3; evidence/slice-backfill-run-store.md — Nenner der zwei
Coverage-Messungen in den Sensor-Dokumenten, Verifikation V-1, und eine Dateizahl im
Suchlauf-Feld, V-2; evidence/slice-081.md, evidence/slice-084.md,
evidence/slice-085.md, evidence/slice-088.md, evidence/slice-089.md,
evidence/slice-090.md, evidence/slice-092.md, evidence/slice-096.md,
evidence/slice-099.md, evidence/slice-101.md, evidence/slice-102.md,
evidence/slice-103.md, evidence/slice-sdk-python-grpc-client-flaeche.md,
evidence/slice-backfill-row-image-gemeinsam.md,
evidence/slice-backfill-change-origin.md,
evidence/slice-backfill-snapshot-reader.md) —
**Schwelle erreicht** (bereits verkörpert, kein neuer Schwellen-Übertritt; der
Lese-Schritt der Closure von `welle-backfill-bestand` liest den Eintrag mit, 20×).
Der zwanzigste Beleg (slice-backfill-slot-leerlauf-bestaetigung) trifft einen Handbuch-Lauf
ohne auflösbaren Träger und **Zahlen im Suchlauf-Feld eines Plans**, dessen Suchraum die
Plan-Datei selbst enthält: ihre Treffer sind ein Selbstverweis und driften mit jedem Edit
(Ausschluss der Plan-Datei aus dem Suchraum und Ausweis „ohne Plan-Datei“ neben der
Zahl); gefunden haben alle der Reviewer und der Verifier durch Nachmessen an beiden Ständen.
Kein Schwellen-Übertritt (bereits verkörpert).
Der neunzehnte Beleg (slice-backfill-bench-richtgroesse) trifft die Klasse an einem
**Mess-Slice**: Zahlen eines Laufs standen als „gemessen“, obwohl die gedruckte Zeile
nicht im Repository liegt (Lauf `20260925T012459Z`, 91,7 %); eine Speicher-„Spitze“ war
kleiner als ein Wert derselben Messung (Probe 20 s nach dem Run); zwei Trefferzahlen
des Suchlauf-Felds nannten „je 1“ statt 2; und eine Träger-Zeile
(`make bench` in `harness/README.md`) führte eine Kalibrierungs-Zahl ohne Host, während
der Messhost des Slice bei 87,5 % bis 95,8 % liegt. Die Fixrunde und die Closure setzen
je Zahl den Ursprung (gemessen · übernommen · abgeleitet) und den Host; gefunden haben
alle der Reviewer und der Verifier durch Nachmessen. Kein Schwellen-Übertritt (bereits
verkörpert).
Der achtzehnte Beleg (slice-backfill-e2e) trifft wieder eine **Suchlauf-Feld**-Zahl:
die Trefferzahl einer Zeile stand für den Stand eines früheren Fixrunden-Commits
(127 Zeilen in 59 Dateien) statt für den Diff-Stand (130 in 60, die Fixrunde legt
eine ADR mit zwei Treffern an); der Reviewer fand sie durch Nachmessen an beiden
Ständen, die Fixrunde setzte die gemessene Zahl.
Der sechzehnte Beleg (slice-backfill-snapshot-reader) trifft einen
Skript-Kommentar („ZWEIMAL" gegen gemessen dreimal), eine lauf-gebunden
streuende gedeckte Zahl (Plan-Zahlen 1691/1731 gegen gemessen 1689 bis 1733) und
zwei stand-relative Zellen des Suchlauf-Felds; die Zellen tragen seither ihre
Commit-Kennung, gefunden hat alle vier der nachmessende Leser. Der fünfzehnte Beleg (slice-backfill-change-origin) trifft wieder ein
**Suchlauf-Feld**: eine Trefferzahl (18) stand als gemessen, real 16 an vier
Ständen, und eine zweite (24 „unverändert") verschob sich durch die Meldungszeile
des Plans selbst auf 25; beide vom Verifier durch Nachmessen gefunden, der Zug
der Closure setzt die gemessene Zahl mit ihren Ständen und kennzeichnet die
abgeleitete. Der vierzehnte Beleg (slice-backfill-row-image-gemeinsam) trifft die Klasse an
einem **Suchlauf-Feld** — dem Träger, der selbst Zahlen über zwei Stände
führt: eine Summe (17) stand als gemessen, real 16, und der Parent-Bezug nannte
`HEAD` statt des Parent-Commits; vom Reviewer durch eigenes Nachmessen
gefunden, in der Fixrunde gezogen. Der dreizehnte Beleg (slice-sdk-python-grpc-client-flaeche) führt die Klasse
dreifach in derselben Korrektur-Kette (Haupt-Review F-1, Verifikation V-1,
Re-Review F-5 — 5/29 → 6/29 → 1.39/1.40): jede Nachzug-Zeile zog die Zahl aus
dem Vorgängerstand des Trägers statt aus der Messung, die sie belegt; die
Lektion „Messung zuerst, dann Text" trägt die Closure-Notiz des Belegs
(`done/slice-sdk-python-grpc-client-flaeche.md` §7). Der zwölfte Beleg (`slice-103`) trifft exakt dieselbe Fundstelle wie die
beiden vorigen (`slice-101`, `slice-102`) — alle drei Pläne entstanden in
rascher Folge und übernahmen dieselbe, bei ihrer jeweiligen Niederschrift
bereits veraltete Zahl „5×" für `BEO-PGC/github-actions-unverifizierbar-lokal`
(real 7×), jedes Mal vom Verifier durch eigene Neuzählung gefunden, nicht
übernommen. Der elfte Beleg (`slice-102`) trifft exakt dieselbe Fundstelle
wie der zehnte (`slice-101`) — beide Pläne entstanden praktisch zeitgleich
und übernahmen dieselbe, bei ihrer jeweiligen Niederschrift bereits
veraltete Zahl „5×" für `BEO-PGC/github-actions-unverifizierbar-lokal`
(real 7×), beide vom Verifier durch eigene Neuzählung gefunden, nicht
übernommen. Der zehnte Beleg (`slice-101`) trifft dieselbe Form wie `slice-099`: kein Drift
durch fortschreitende Arbeit, sondern eine im Plan-Kopf/§6/§8 selbst
zitierte Zahl, die bereits bei ihrer eigenen Niederschrift veraltet war.
Der siebte Beleg trifft die
Klasse **im Korrektur-Vorgang selbst**: Glob-Fehler → vier Zahlen für einen
Gegenstand → Zähl-Wort über vier Größen → invertiertes Herkunfts-Etikett,
vier Runden an Sätzen über Zahlen. Die Fundstellen je Vorgang liegen **im
selben** Vorgang und sind damit je *eine* Gelegenheit — der Zähler misst
Wiederholung über Vorgänge, nicht die Zahl der Funde. `slice-084` trug
**vier** driftende Werte in zwei Sensor-Dokumenten; zwei davon stammten aus
**anderen** Vorgängen und wurden mitgezogen, weil sie sonst in derselben Datei
gegen die eigene Messung stünden. `slice-085` hat die Form dann **an sich selbst
gebrochen** (der Nenner, den die Vorgänger-Fixrunde zum Zustand erklärt hatte,
blieb stehen) — der Grund, warum die Durchsetzung ein Träger und kein Vorsatz
sein muss. `slice-090` hat die Klasse dann **in ihrer eigenen Behebung**
getroffen: der Satz, der die falsche Zeilenangabe (`review-slice-090` F-2)
ersetzte, trug eine neue falsche Abstands-Angabe („eine Zeile voraus" gegen
gemessen 0/1/2/3). `slice-096` traf **zwei** falsche Zahlen in zwei am selben
Tag entstandenen ADRs, in einer Runde gefunden. `slice-099` erweitert die
Reichweite der Klasse: Ein Digest-Tabellenwert in einer bereits `Accepted`-ADR
(`ADR-0087`) war **seit ihrer Annahme nie** korrekt (63 statt 64 Hex-Zeichen)
— kein Drift durch spätere Arbeit, sondern ein von Anfang an fehlerhafter,
nie gegen die reale Registry nachgeprüfter Wert; der real gebaute Code war
davon nie betroffen. Aufgelöst über die engräumige Folge-ADR `ADR-0093`
(`Supersedes ADR-0087` für die eine Tabellenzelle) — kein neuer Sensor,
`ADR-0093`s eigene Fitness-Function-Tabelle hält das ausdrücklich fest.

Ein Sensor ist **nicht** vorgeschlagen: verlangte er, dass jede Zahl ihren
Ursprung trägt, wäre er eine Formpflicht auf Prosa und erzeugte
Pflichterfüllung — er hätte genau die Klasse, die er prüft. Verfügbare
Falsifikation bleibt die **Messung** selbst.

**Nicht zu verwechseln** mit dem benachbarten Eintrag
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (2×): dort sagt ein
Kommentar einen **Fehlerpfad** zu, den der Code nicht trägt (die Aussage kehrt
das Verhalten um); hier nennt ein Träger eine **Zahl**, die gegen die Messung
driftet. Der Review zu `slice-081` hatte seinen F-1 jenem Eintrag zugeordnet —
das ist nach Prüfung der beiden Vorgänger-Findings (`review-slice-070` F-2,
`review-slice-077` F-1, beide Fehlerpfad-Fälle) **nicht** dieselbe Beobachtung.
