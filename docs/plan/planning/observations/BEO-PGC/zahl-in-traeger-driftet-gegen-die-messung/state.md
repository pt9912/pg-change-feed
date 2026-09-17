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

Zähler (abgeleitet): **9×** (evidence/slice-081.md, evidence/slice-084.md,
evidence/slice-085.md, evidence/slice-088.md, evidence/slice-089.md,
evidence/slice-090.md, evidence/slice-092.md, evidence/slice-096.md,
evidence/slice-099.md) — **Schwelle erreicht**. Der siebte Beleg trifft die
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
