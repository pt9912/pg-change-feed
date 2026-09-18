# Architect-Verdikt: Nackte Kennung ohne Link — 4. Auftreten widerlegt die Prämisse des `gestrichen`-Verdikts

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/report-nackte-id-ohne-link`
— nach dem `gestrichen`-Verdikt bei 3× (der erste Architect-Verdikt zu report-nackte-id-ohne-link)
tritt die Klasse ein viertes Mal auf: `slice-063` (das Blocker-Protokoll zu `slice-063`), real committet und gepusht mit drei
nackten `ADR-*`-Kennungen ohne Link (`evidence/slice-063-blocker.md`). Der
Fehlerpfad unterscheidet sich real von allen drei vorherigen Belegen —
weder Pipe-Maskierung (`slice-054`) noch Selbstkorrektur vor Commit
(`slice-055`) noch der bereits vom 3×-Verdikt als „funktioniert genau wie
vorgesehen" gewertete Fall (`slice-056`). Direkter Architect-Zug aus einer
Nutzeranfrage, vor jeder formalen Slice-Closure (`slice-063` liegt zum
Zeitpunkt dieses Verdikts in `open/`, zurückgeführt über einen Blocker —
siehe `docs/plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md`).
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-14
**Bezug:** [`AGENTS.md`](../../AGENTS.md) §3.9 (die hier geschärfte Regel),
Modul 8 §Kernidee und §Konflikt-Pfad als Rollen-Sequenz, Modul 6 §Das
Beobachtungs-Register,
der erste Architect-Verdikt zu report-nackte-id-ohne-link
(3×-Verdikt, hier neu bewertet),
der Architect-Verdikt-Nachtrag zur Slice-Chronik in Code-Kommentaren (4x)
(Präzedenzfall für ein 4.-Auftreten-Verdikt, dort: Status quo bestätigt +
zwei gezielte Verkörperungen; hier: Prämisse widerlegt + eine gezielte
Verkörperung),
der Architect-Verdikt dazu, dass eine Pipe den `make`-Exit-Code maskiert
(Ursprungs-Verkörperung von `AGENTS.md` §3.9, hier ergänzt, nicht
ersetzt), `docs/plan/planning/observations/BEO-PGC/report-nackte-id-ohne-link/`,
das Blocker-Protokoll zu `slice-063`.

---

## Frage

Trägt der `gestrichen`-Ausgang der 3×-Auswertung noch, nachdem ein reales
viertes Auftreten seine tragende Prämisse — „der einzige Weg, den roten
Befund zu ignorieren, ist bereits durch `AGENTS.md` §3.9 verschlossen" —
empirisch widerlegt? Oder bestätigt das vierte Auftreten den Ausgang mit
geschärfter Begründung, weil die zugrunde liegende Disziplin bereits
existiert und nur nicht befolgt wurde (Compliance-Lücke, keine
Konstruktions-Lücke — die Diagnose, die der Chronik-4×-Fall für seinen
Fall traf)?

## Was am vierten Beleg real anders ist

Die drei vorherigen Belege teilen sich in genau zwei Fehlerklassen, beide
bereits durch das `gestrichen`-Verdikt korrekt zugeordnet:

| Vorgang | Fehlerklasse | Status nach 3×-Verdikt |
|---|---|---|
| `slice-054` | Exit-Code durch eine Pipe **maskiert** (`make gates \| tail`) — der rote Wert wurde nie sichtbar | geschlossen durch `AGENTS.md` §3.9 (Pipe-Verbot) |
| `slice-055` | Selbstkorrektur vor Commit — kein Sensor-Fall | kein Fehlerpfad, nie offen |
| `slice-056` | Sensor griff korrekt **vor** dem Push, weil der Exit-Code diesmal nicht maskiert war | als Beleg für die Wirksamkeit von §3.9 gewertet |

Der vierte Beleg (`slice-063-blocker`, `evidence/slice-063-blocker.md`)
fällt in **keine** der beiden Klassen: Der Exit-Code wurde real korrekt
und ungepiped ermittelt (`EXIT=2`, sichtbar im eigenen Tool-Output — exakt
das, was §3.9 verlangt) — aber die davon abhängige Folgehandlung
(`git push`) lief im selben Arbeitsschritt-Batch, **bevor** der bereits
sichtbare rote Wert die Aktion tatsächlich blockierte. Es handelt sich
nicht um eine falsche Messung (wie bei `slice-054`), sondern um eine
korrekte Messung, deren Ergebnis die nachfolgende Handlung nicht
tatsächlich gesteuert hat.

`AGENTS.md` §3.9 adressiert in seinem heutigen Wortlaut ausschließlich die
**Mess-Ebene**: „Exit-Code … wird direkt geprüft, nie durch eine
Pipe/einen Wrapper hindurch" — die Diagnose und alle drei genannten
Ursprungsfälle (`welle-15`) drehen sich um *falsch ermittelte* Werte
(Pipe, Wrapper). Die „Richtig"-Beispielzeile enthält zwar bereits das
Muster `test $ec -eq 0 && git push` — aber als *ein* Beispiel unter
mehreren, nicht als eigenständig benannte Pflicht, dass Prüfung und
Folgehandlung **getrennte Schritte** sein müssen. Der vierte Beleg zeigt
genau die Lücke, die diese Unschärfe offen lässt: eine korrekt gemessene
rote Ampel, gefolgt von einer Handlung, die im selben Zug beauftragt
wurde, ohne dass der bereits vorliegende Messwert sie tatsächlich
verhinderte.

## Diagnose: eine andere Fehlerklasse, keine Wiederholung derselben

Der Chronik-4×-Fall hatte sein viertes Auftreten als *Compliance-Lücke,
keine Konstruktions-Lücke* diagnostiziert, weil das bestehende Muster
(Schritt 20) den Fall strukturell bereits abdeckte und schlicht nicht
ausgeführt wurde. Hier liegt der Fall anders: Das `test $ec -eq 0 &&
git push`-Muster stand zwar als *Beispiel* in §3.9, aber §3.9 selbst
benennt an keiner Stelle explizit, dass eine Gate-Prüfung und ihre
abhängige Folgehandlung **nicht im selben Werkzeug-Aufruf-Batch**
beauftragt werden dürfen. Das ist eine Lücke, die spezifisch für
agentische Ausführung zutrifft: In einem einzelnen Shell-Skript erzwingt
`&&` die Reihenfolge mechanisch; über mehrere Werkzeug-Aufrufe eines
Agenten-Laufs hinweg (ein Aufruf ermittelt den Exit-Code, ein späterer
Aufruf pusht) gibt es diese mechanische Kopplung nicht von selbst — sie
hängt von der **Disziplin des Agenten**, das Ergebnis des ersten Aufrufs
tatsächlich abzuwarten und auszuwerten, bevor der zweite beauftragt wird.
Der ursprüngliche Pipe-Fix schloss die Mess-Ebene; er sagte nichts über
die Sequenzierungs-Ebene zwischen zwei getrennten Werkzeug-Aufrufen.

Das ist **keine** Wiederholung der bereits geschlossenen Pipe-Klasse
(`slice-054`) und **kein** Beleg gegen die Wirksamkeit von §3.9 für den
Fall, den es ursprünglich adressierte (`slice-056` zeigt weiterhin, dass
die Mess-Ebene hält). Es ist eine **vierte, eigenständige** Fehlerklasse
innerhalb derselben Beobachtung — die Beobachtung selbst führt Fälle mit
nackten Kennungen im Report, unabhängig von deren jeweiliger Ursache.

## Verdikt: Ausgang wechselt von `gestrichen` zu `verkörpert`

Die Prämisse des 3×-Verdikts — „der reale Risikopfad … ist strukturell
geschlossen … der einzige Weg, den roten Befund zu ignorieren, ist bereits
verschlossen" — trifft nicht mehr uneingeschränkt zu: Es gibt einen
zweiten, real aufgetretenen Weg, einen roten Befund wirkungslos zu
lassen, den §3.9 in seinem bisherigen Wortlaut nicht ausdrücklich
schließt. Der 3×-Verdikt wird **nicht** per `supersedes`-Logik verworfen
(ADRs kennen dieses Mittel, ein Architect-Verdikt als Prosa-Dokument
nicht) — er bleibt für die Mess-Ebene weiterhin richtig. Dieses Verdikt
tritt **neben** ihn und schärft die Lücke, die der erste Verdikt nicht
sehen konnte, weil sie zum Zeitpunkt seiner Abfassung noch nicht real
aufgetreten war.

**Verkörperung:** `AGENTS.md` §3.9 wird um einen eigenen, benannten Absatz
ergänzt — Prüfung und Folgehandlung sind zwei getrennte Schritte, nicht
einer; ein Gate-Lauf und die davon abhängige Folgehandlung dürfen nicht im
selben Werkzeug-Aufruf-Batch beauftragt werden. Kein neuer Sensor: Die
Verletzung liegt wie beim ursprünglichen Pipe-Fall in der Ausführung
selbst, nicht im committeten Ergebnis — kein zweites, unabhängig
einsehbares Artefakt (Diff) hätte sie fangen können. Dieselbe Begründung,
mit der §3.9 selbst ohne Sensor auskommt, trägt auch für seine Schärfung.

## Was dieses Verdikt tut

- `AGENTS.md` §3.9 bekommt einen neuen Absatz „Prüfung und Folgehandlung
  sind zwei Schritte, nicht einer" mit Falsch-/Richtig-Beispiel, der das
  vierte Auftreten benennt (Anker `seit slice-063`, analog zur bestehenden
  `welle-15`-Referenz für die ersten drei Fälle).
- `docs/plan/planning/observations/BEO-PGC/report-nackte-id-ohne-link/state.md`
  wird auf Ausgang **verkörpert** aktualisiert, mit Zielort- und
  Herkunfts-Anker auf den neuen §3.9-Absatz — als wellenloser Lese-Schritt,
  direkt von diesem Architect-Zug mitgeführt (Modul 6 §Wann Arbeit eine
  Welle braucht, Tabelle *Träger im Repo ohne Wellen*; Modul 8 §Rollen-
  Sequenz für eine Welle, Übergabe Planner → Architect → Planner bleibt
  auch ohne Wellen-Betrieb erhalten). `slice-063` liegt zum Zeitpunkt
  dieses Verdikts in `open/`, nicht in `done/` — der reguläre
  Zähler-/Ausgang-Eintrag durch den Planner bei einer künftigen
  `slice-063`-Closure ist damit vorweggenommen, analog zum
  Chronik-4×-Verdikt bei `slice-052`.

## Was dieses Verdikt NICHT tut

- Es hebt den 3×-Verdikt nicht auf und erklärt ihn nicht für falsch — für
  die Fehlerklasse, die er analysierte (Pipe-/Wrapper-Maskierung), bleibt
  seine Analyse zutreffend, belegt durch `slice-056`.
- Kein neuer Sensor, kein neues Gate — dieselbe Begründung wie beim
  Ursprungs-§3.9: Die Verletzung liegt in der Ausführung, nicht im
  committeten Ergebnis, kein Artefakt im Working Tree träfe sie.
- Keine Änderung an `.harness/skills/reviewer.md` — dies ist keine
  Review-Findings-Frage (kein Diff, den ein Reviewer prüfen könnte,
  betrifft die Ausführungsdisziplin selbst, wie beim Ursprungsfall von
  §3.9).
- Kein neues ADR — reine Ausführungsdisziplin-Schärfung wie der
  Ursprungsfall, keine Architektur- oder Vertragsentscheidung.
- Der bestehende `gestrichen`-Vermerk in `state.md` für die ersten drei
  Belege bleibt als Historie stehen; er wird nicht rückwirkend
  umgeschrieben (`git` hält die Chronik, das Feld selbst nennt nur den
  aktuellen Zustand und seinen Beleg, `AGENTS.md` §3.7).

Weder Produktionscode noch eine ADR-Datei noch eine Skill-Datei wurden im
Rahmen dieses Verdikts geändert — ausschließlich `AGENTS.md` §3.9 und der
Registereintrag.
