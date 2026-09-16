Zustand: **offen — 3× erreicht, Ausgang noch nicht zugewiesen.** Der
Lese-Schritt gehört der **laufenden Welle-Closure** (`welle-20`), nicht der
Slice-Closure (Modul 6). Bis dahin ist `offen` der zulässige, vorübergehende
Stand.

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

Zähler (abgeleitet): **3×** (evidence/slice-081.md, evidence/slice-084.md,
evidence/slice-085.md) — **Schwelle erreicht**. Die Fundstellen je Vorgang
liegen **im selben** Vorgang und sind damit je *eine* Gelegenheit — der Zähler
misst Wiederholung über Vorgänge, nicht die Zahl der Funde. `slice-084` trug
**vier** driftende Werte in zwei Sensor-Dokumenten; zwei davon stammten aus
**anderen** Vorgängen und wurden mitgezogen, weil sie sonst in derselben Datei
gegen die eigene Messung stünden. `slice-085` hat die Form dann **an sich selbst
gebrochen** (der Nenner, den die Vorgänger-Fixrunde zum Zustand erklärt hatte,
blieb stehen) — der Grund, warum die Durchsetzung ein Träger und kein Vorsatz
sein muss.

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
