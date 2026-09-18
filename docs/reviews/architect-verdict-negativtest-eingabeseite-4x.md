# Architect-Verdikt: Zusage ohne Bindung an ihre Eingabeseite — 4. Auftreten

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
— der Eintrag steht bei **4×** (Schwelle erreicht). Der Lese-Schritt gehört
der laufenden `welle-20`-Closure; er wird hier **vorgezogen**, weil der gerade
geschlossene `slice-088` allein sechs Funde dieser und der Schwesterklasse
erzeugt hat. Der Zug läuft **vor** der formalen Closure: der Zähler-Eintrag und
der Register-Ausgang durch den Planner stehen noch aus
(siehe §Was dieses Verdikt NICHT tut).
**Rolleninhaber:** pt9912 (Architect-Rolle, dieser Lauf)
**Datum:** 2026-09-16
**Bezug:** die vier Beleg-Dateien des Eintrags
(`evidence/slice-083.md`, `evidence/slice-086.md`, `evidence/slice-087.md`,
`evidence/slice-088.md`), Review zu `slice-086` (F-2, mit
eigener Mutation D), Review zu `slice-087` (F-1, Mutation
C), Review zu `slice-083` (F-2) und
Verifikationsbericht zu `slice-083`,
Review zu `slice-088` (F-1),
Delta-Review zu `slice-088` (M4–M7),
Verifikationsbericht zu `slice-088` (eigene Reproduktion) ·
`.harness/skills/reviewer.md` (Vorschlags-Träger des Eintrags),
`.claude/commands/implement-slice.md` Schritt 19 (die Mutations-Pflicht),
`.claude/agents/verifier.md` („Prüfe die Belege, nicht die Behauptung") ·
[`../plan/adr/0083-herkunft-von-aussagen-in-traegern.md`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
(Schwesterklasse, anderer Mechanismus) · Modul 8 §Rollen-Regeln und §Welche
Rolle braucht welche Artefaktklasse · Modul 6 §Das Beobachtungs-Register

---

## Frage

Trägt die Klasse „eine Aussage, die nicht an ihre Eingabeseite gebunden ist,
ist grün ohne Aussage" eine verkörperte Regel — und wenn ja, an welchem Träger?
Der Eintrag nennt als Kandidaten `.harness/skills/reviewer.md`, „wo die
Mutations-Pflicht bereits steht — aber ohne diese Richtung". Zu prüfen ist
damit zweierlei: **ob** die Klasse verkörperbar ist (statt eine benannte
Grenze zu bleiben) und **welcher** Träger sie trägt.

## Empirischer Befund

**Vier Vorgänge, vier Mal vor Merge gefangen — kein Beleg hat je `main`
erreicht.** `slice-086` (Negativtest `?limit=0 → 400`), `slice-087` (E2E-Beleg
für den Filter), `slice-083` (Kommentar samt Test über eine Eigenschaft, die
kein Input ausübt), `slice-088` (zwei Noop-Zusagen gegen einen
argument-ignorierenden Stub). Je Fall wurde eine Fixrunde gefahren und die
Bindung nachgeholt (`slice-086`: Commit `f20d2c9`; `slice-088`: die Fixrunde
über einen Stub, der die empfangenen Kennungen protokolliert, plus
`t.Cleanup`-Nachweis).

**Herkunft dieser Angaben:** übernommen aus den vier Beleg-Dateien des
Eintrags; die nennen ihrerseits Report und Commit. **Nicht** in diesem Zug
nachgemessen — die Mutationen sind Go-Testläufe (Docker-only, Kosten des vollen
Paketlaufs), und der Zug prüft die **Form** der Klasse (sind die vier Belege
derselbe Mechanismus?), nicht ihre Zahl. Die vier Belege sind im Mechanismus
deckungsgleich: *eine Aussage, die nicht an ihre Eingabeseite gebunden ist, ist
grün ohne Aussage.*

**Der schärfste Beleg ist `slice-088`, und er entscheidet den Träger.** Der
Implementer hat in diesem Vorgang **fünf** eigene Mutationen gefahren und alle
rot gesehen; die zwei löchrigen Tests waren genau die, die er für
selbstverständlich hielt (Review zu `slice-088`, F-1). Die Mutations-Pflicht
(`implement-slice.md` Schritt 19) war also **ausgeführt** und hat die Klasse
trotzdem nicht gefangen. Das ist keine Compliance-Lücke, sondern eine
**Richtungs**-Lücke: mutiert wurde die **Ausgabeseite** — der Rückgabewert, der
Fake —, nicht die Eingabe.

**In drei der vier Fälle war der Reviewer die tragende Instanz**, und in allen
dreien fand er die Klasse nicht durch Lesen, sondern durch **Mutieren der
Eingabeseite** (086 Mutation D: `readChangesLimit("0") → nil`; 087 Mutation C:
Filterachse entfällt; 088 M4/M5: die beiden `if`-Zweige entfernt, je Exit 0).
Der Verifier hat in `slice-088` M4/M5 eigenständig reproduziert und M6/M7 als
Bindungssonden bestätigt.

## Befund zum Träger-Vorschlag des Eintrags

Der Eintrag nennt `.harness/skills/reviewer.md` als Ort, „wo die
Mutations-Pflicht bereits steht". **Geprüft: dort steht sie nicht.** Der
Reviewer-Skill führt heute keine Mutations-Pflicht; sie steht in
`.claude/commands/implement-slice.md` **Schritt 19** („Zu jedem neuen oder
geänderten Wächter die rot färbende Mutation benennen") — auf der
**Implementer**-Seite. Der Vorschlag trifft damit die **Finder**-Seite (die
Rolle, die die Klasse in 3 von 4 Fällen gefunden hat), nicht die **Träger**-
Seite der Pflicht. Beide Seiten werden verkörpert; die Richtung reitet in
beide. Die Korrektur ändert den Ausgang des Eintrags nicht, nur seine Adresse.

## Diagnose

Warum die Ausgabeseite der Default ist: Ein Negativtest **sieht vollständig
aus**, weil er etwas prüft — der Fake liefert den Fehler. Das sichtbare Objekt
des Tests ist der Rückgabewert; das mittlere Glied der Kette (Eingabe →
**unveränderte Weitergabe** → Ablehnung) ist die einzige Arbeit des Adapters
und liegt unsichtbar dazwischen. Wer am Rückgabewert mutiert, prüft den Fake;
wer die Eingabe mutiert, prüft, ob die Aussage an ihr hängt.

Daraus folgt der Zuschnitt des Verdikts: Die Klasse braucht **keinen neuen
Wächter**, sondern eine **Richtung** an den zwei Stellen, an denen die Prüfung
ohnehin stattfindet. Sie hat eine korrekt ausgeführte Mutations-Pflicht
überlebt (088) und eine vorhandene Verifier-Disziplin („Prüfe die Belege,
nicht die Behauptung" — der Verifier fand die Klasse nicht als eigenes
Ergebnis, sondern reproduzierte die Mutation des Reviewers).

## Verdikt

**Verkörpert — an zwei Trägern, mit einer Regel.** Die beiden Träger sind zwei
Rollen mit **verschiedenem Eingabe-Kontext** (Autor: eigene Arbeit; Reviewer:
der Diff), die Mehrfachzuweisung ist damit nach Modul 8 §Rollen-Regeln sauber
und keine doppelte Arbeit.

**Der Wortlaut — eine Regel, beide Träger:**

> *Eine Zusage ist nur dann gebunden, wenn der Test an ihrer
> **Eingabeseite** rot werden kann: mutiere den **Eingabewert**, nicht nur die
> Ausgabeseite. Wer nur den Fake oder den Rückgabewert mutiert, prüft den
> Fake — die Aussage bleibt grün, egal was der Adapter mit der Eingabe tut.
> Fehlt die Mutation der Eingabeseite, ist die Zusage **grün ohne Aussage**:
> ein Befund, kein Formfehler.*

**Träger 1 — `.harness/skills/reviewer.md`, benannter HIGH-Unterpunkt.** Die
Hausform der übrigen repo-spezifischen HIGH-Punkte: benannte Klasse,
Abgrenzung, „kein Gate fängt das", Herkunft mit Zählerstand und einem
Herkunfts-Anker. Begründung: der Reviewer ist die tragende Instanz (3 von 4
Fällen, alle durch eigene Mutation), und eine wiederkehrende Klasse verdient
einen eigenen Anker statt einer Subsumtion unter „fehlende Negativtests bei
neuem öffentlichem Vertrag" (MEDIUM) — sie ist keine fehlende Abdeckung,
sondern eine **vorhandene Zusage ohne Bindung**.

**Träger 2 — `.claude/commands/implement-slice.md` Schritt 19, die Richtung.**
Der Schritt verlangt heute „zu jedem neuen oder geänderten Wächter die rot
färbende Mutation benennen". Er bekommt die Richtung **explizit**: welche
Mutation der **Eingabeseite** müsste diesen Test rot machen, und wurde sie
einmal gesehen? Dazu in derselben Form wie Schritt 20 („Enumerations-Pflicht
statt Erinnerung"): **je Zusage eine benannte Eingabeseiten-Mutation** — die
bloße Zahl der gefahrenen Mutationen trägt nicht, wie `slice-088` zeigt (fünf
gefahren, zwei Aussagen ungebunden).

**Kein Sensor, und keiner wird bestellt.** „Welche Änderung müsste diesen Test
rot machen?" ist die Mutation selbst; ein Mutations-Harness (Mutation Testing
über den Baum) gibt es in diesem Repo nicht, und Schritt 19 führt es
ausdrücklich als „falls vorhanden". Ein solcher Tier wäre ein eigener Vorgang
mit eigenem Aufwand (Laufzeit je Mutation × Testfälle) — nicht als Nebenprodukt
einer Klassen-Entscheidung zu beschließen. **Die verfügbare Falsifikation ist
die Mutation**, und sie hat die Klasse in allen vier Vorgängen widerlegt.

## Was dieses Verdikt NICHT tut

- Es schreibt **keine** ADR: kein Hard Rule in `AGENTS.md`, kein Gate, keine
  Schwelle. Die Verkörperung ist ein Rollen-Skill und ein Workflow-Schritt —
  die Hausform dafür ist ein Verdikt, kein Beschluss (Vergleich: der
  Architect-Verdikt-Nachtrag zur Slice-Chronik in Code-Kommentaren, 4x).
- Es legt **keinen** Eintrag in `evidence/` an und bumpt **keinen** Zähler. Der
  Ausgang wird regulär beim Lese-Schritt der `welle-20`-Closure durch den
  Planner ins Register geschrieben (Modul 6).
- Es fasst **keinen** Produktionscode und **keinen** Test an.
- Es führt **keine** neue Rolle und **keinen** neuen Skill ein (Modul 8 §Welche
  Rolle braucht welche Artefaktklasse: Skills wachsen pro Urteilstyp, nicht pro
  Rolle — der Mutations-Typ bleibt beim bestehenden Skill).

## Was das für die Rollen bedeutet

**Implementer:** Schritt 19 wird um die Richtung und die Enumerations-Pflicht
erweitert. Beobachtbarer Unterschied: nicht „ich habe mutiert", sondern eine
Liste — *Zusage · mutierte Eingabe · gesehenes Rot*. Wo die Liste leer bleibt,
steht der Grund.

**Reviewer:** ein eigener, direkt zitierbarer HIGH-Anker statt der Subsumtion
unter die Negativtest-Klasse. Beobachtbare Folge: die Klasse wird nicht mehr
als „fehlender Negativtest" (MEDIUM) geführt, sondern als **ungebundene
Zusage** — die Kategorie ist genau das, was über die Fixrunde entscheidet.

## Lese-Schritt vorweggenommen

Für den Lese-Schritt der `welle-20`-Closure ist die Antwort damit
vorweggenommen: Ausgang **`verkörpert`**, Zielorte `.harness/skills/reviewer.md`
(HIGH-Unterpunkt) und `.claude/commands/implement-slice.md` Schritt 19
(Richtung + Enumerations-Pflicht); Herkunfts-Anker `seit welle-20`, wenn die
Welle-Closure die Verkörperung selbst schreibt — sonst `seit slice-<NNN>` des
schreibenden Vorgangs. Der Planner muss keine neue Architect-Eskalation
auslösen, nur den Ausgang und den Verweis nachtragen.
