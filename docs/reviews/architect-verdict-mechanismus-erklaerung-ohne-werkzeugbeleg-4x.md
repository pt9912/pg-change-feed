# Architect-Verdikt — Mechanismus-Erklärung ohne Werkzeugbeleg, 4. Auftreten (Modul 6, Prosa ausgeschöpft?)

**Datum:** 2026-10-10 · **Stand:** `11da9df8` (Baum sauber) · **Rolle:** Architect (Modul 8),
Lese-Schritt der wellenlosen Closure, anderer Kontext als Planner-, Implementer-, Reviewer- und
Verifier-Lauf des Slice · **Rolleninhaber:** pt9912 (Architect-Agent im Auftrag) · **Anlass:**
Closure-Notiz von `slice-spec-festlegungen-doku-gates` (`docs/plan/planning/done/`), §7
„Lese-Schritt“ und „Steering-Loop-Eintrag“: der Planner hat den Ausgang des Eintrags
`BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg` unverändert gelassen und die Frage, ob der
Modul-6-Satz zum vierten Auftreten greift, an den Architect gegeben.

**Bezug:** `docs/plan/planning/observations/BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg/`
(vier Beleg-Dateien) · [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B ·
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Entscheidung und benannte
Grenze: kein Sensor als Formpflicht auf Prosa) · `.claude/commands/plan-welle.md` Schritt 6 ·
Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3
(„Verkörpert heißt nicht zwangsläufig automatisiert“) ·
[`review-slice-spec-festlegungen-doku-gates`](review-slice-spec-festlegungen-doku-gates.md)
F-1 bis F-4 · [`verify-slice-spec-festlegungen-doku-gates`](verify-slice-spec-festlegungen-doku-gates.md)
A-1 · [`verify-slice-spec-festlegungen-doku-gates-nachpruefung`](verify-slice-spec-festlegungen-doku-gates-nachpruefung.md)
A-6.

---

## 0. Ergebnis in einem Blick

- **Der Modul-6-Satz greift.** Die Klasse ist nach ihrer Verkörperung als Prosa
  (`AGENTS.md` §3.12 Instanz B, seit slice-089) erneut aufgetreten. Das ist der Fall, für den
  Modul 6 die Prosa-Form für ausgeschöpft erklärt. Der neue Eintrag muss einen Sensor
  benennen oder begründen, warum es keinen gibt. „Ausgang unverändert“ ohne diese Begründung
  genügt nicht (§1).
- **Kein Sensor möglich, mit Begründung (§3).** Ein Werkzeug müsste entscheiden, welcher Satz
  ein Werkzeug beschreibt und ob ein Lauf ihn trägt. Das erste ist Semantik. Das zweite ist die
  Messung selbst. Eine Pflicht „jeder Satz hat eine Zeile“ wäre die Formpflicht auf Prosa, die
  `ADR-0083` als Grenze ausschließt.
- **Eine Anpassung des Trägers, keine zweite Schärfung von §3.12.** Der vierte Beleg hatte
  einen Beleg-Anker, nur an der falschen Stelle: den Quellen statt dem Werkzeug (§2).
  `plan-welle` Schritt 6 verlangte bisher genau diese Gegenprobe der Quellen. Er bekommt deshalb
  den Punkt (c) **Messung am Werkzeug** (§4). Das ist Prosa und als Prosa benannt. Sie verschiebt
  die Prüfung von der Beschreibung auf das Werkzeug, auf das die Messungen in diesem Slice schon
  gezeigt haben.
- **Kein Folge-Slice, keine ADR, kein neuer Register-Eintrag.** Der Ausgang bleibt
  **verkörpert**, jetzt an zwei Zielorten; `state.md` trägt die Begründung.

## 1. Greift der Satz? Und wie lief der Pilot-Fall?

**Wortlaut** (`modul-06-roadmap.md`, Closure Schritt 3): „Erreicht dieselbe Fehlerklasse
trotzdem ein **viertes** Mal die Schwelle, gilt die Prosa-Form als ausgeschöpft“. Er steht
direkt nach dem Satz, dass ein Eintrag beim dritten Auftreten „meist als geschärfte Instruktion
verkörpert“ wird. Das Wort „trotzdem“ verweist auf diese Verkörperung. Gemeint ist also ein
Auftreten, das die bereits stehende Prosa-Regel nicht verhindert hat.

**Dieser Eintrag:** `ls evidence | wc -l` zeigt am Stand `11da9df8` den Wert `4`. Die
Verkörperung kam mit dem dritten Beleg (slice-089). Der vierte
(`evidence/slice-spec-festlegungen-doku-gates.md`) ist das erste Auftreten nach ihr. Damit ist
der Wortlaut erfüllt, ohne dass man ihn auslegen muss.

**Wie das Repo den Satz bisher anwendet** (gelesen in `docs/reviews/`): Es gibt drei Verdikte
zum vierten Auftreten nach einer Prosa-Verkörperung. Das sind
[`architect-verdict-negativtest-eingabeseite-4x`](architect-verdict-negativtest-eingabeseite-4x.md),
[`architect-verdict-report-nackte-id-ohne-link-4x`](architect-verdict-report-nackte-id-ohne-link-4x.md)
und
[`architect-verdict-slice-chronik-in-code-kommentar-4x`](architect-verdict-slice-chronik-in-code-kommentar-4x.md).
Jedes hat die Frage nach dem Sensor ausdrücklich gestellt und beantwortet. Die Praxis wendet
den Satz also an.

**Pilot-Fall `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`** (6 Belege; die Verkörperung
kam mit dem fünften, `seit welle-transformationen`): Der sechste Beleg aus
`slice-spec-festlegungen-harness-werkzeuge` ist ebenfalls das erste Auftreten nach der
Verkörperung. Die Closure dieses Slice (`a1da6114`) hat den Ausgang „unverändert verkörpert“
eingetragen, ohne Architect-Zug. Das zugehörige `state.md` trägt aber einen Satz der
geforderten Art: „ein Sensor ist nicht vorgeschlagen — ob eine Adresse den Gegenstand trägt,
ist eine Lese-Handlung am Plan der Adresse“. **Im Ergebnis war der Satz damit angewandt**, in
Kurzform und ohne Verweis auf Modul 6, **nicht übergangen**. Als Präzedenz für „unverändert,
ohne Begründung“ taugt er nicht. Diesem Eintrag fehlt dagegen genau dieser Satz: Sein
`state.md` und §7 der Closure sagen nur „Neuer Sensor: keiner“. Das holt dieses Verdikt nach.
Am Pilot-Eintrag selbst ist nichts nachzuziehen. Seine Begründung trägt, und der Satz dort,
„ein weiteres Auftreten ist Anlass, die Frage eines Sensors erneut zu stellen“, verlangt beim
nächsten Auftreten ohnehin einen Architect-Zug.

## 2. Diagnose des vierten Belegs

Die drei früheren Belege (Review slice-079, Review slice-080, slice-089) waren Sätze **ohne**
Anker: Eine Erklärung stand da, ohne dass jemand nachgesehen hatte. Instanz B von §3.12 zielt
genau darauf. Sie verlangt einen Beleg-Anker (Befehl, Datei, Abfrage) oder die Form „erwartet“.

Beim vierten Beleg war das anders. Hinter den Sätzen von `SPEC-040` stand eine vollständige
Gegenprobe mit 76 Zeilen (Plan §3). Jeder Satz hatte damit einen Anker: eine Datei, nämlich
Vertrag, Kommentar in `.d-check.yml` oder ADR. Bei F-4 war es eine Messung aus slice-077 an
einem älteren d-check, ohne Kennzeichnung als *übernommen*. **§3.12 war der Form nach
erfüllt.** Der Anker zeigte aber auf eine **Beschreibung** des Werkzeugs, nicht auf das
Werkzeug. Laut Closure-Notiz §7 (1) waren die Quellen selbst breiter oder enger als das
Werkzeug. Die Messungen des Implementers (M0 bis M23 im Beleg zu Liefer-Punkt 1) gab es zwar,
sie waren aber nicht Satz für Satz an die Festlegung gebunden. Gefunden haben die Abweichungen
die Mutationen des Reviewers (R1 bis R13) und die 20 Läufe des Verifiers. A-6 gehört zur
Variante „gemessen, aber zu schmal“: Die Reihenfolge stützte sich auf fünf Läufe, war im 20er-Lauf
6-mal anders und ist deshalb jetzt nicht mehr festgelegt (`928de0ff`).

**Folgerung:** Eine dritte Prosa-Schärfung von §3.12 („wirklich einen Anker nennen“) hätte
diesen Fall nicht verhindert, denn einen Anker gab es. Der Fehler lag in einem bestimmten
Schritt. `plan-welle` Schritt 6 schreibt für Festlegungen die Gegenprobe der **Quellen** als
Beleg vor. Für eine Festlegung über ein Werkzeug ist das der falsche Gegenstand. Der Schritt
hat damit eine Sicherheit erzeugt, die er nicht liefern kann.

## 3. Verdikt zum Sensor: keiner möglich, und warum

Die Modul-6-Frage lautet: Kann ein mechanischer Sensor die Klasse künftig fangen? Ich habe drei
Formen geprüft. Keine fängt die Klasse.

1. **Ein Sensor, der Sätze über Werkzeuge erkennt und einen Lauf verlangt.** Ob ein Satz in
   Spec, Vertrag oder Skript-Kopf ein Werkzeugverhalten behauptet, ist eine Frage der
   Bedeutung. Ein Muster wie „meldet“ oder „prüft“ trifft die Hälfte der Prosa im Repo und
   übersieht zum Beispiel „gilt jedem geschlossenen Slice-Plan“ (F-1). Das ist dieselbe
   Unmöglichkeit, die `AGENTS.md` §3.13 und `ADR-0083` §Grenze schon für Zahlen benennen.
2. **Eine Formpflicht: Jeder Festlegungssatz trägt eine Kennung, die eine Messungszeile
   referenziert, und ein Skript zählt nach.** Das ist die Formpflicht auf Prosa, die
   `ADR-0083` ausschließt. Sie würde Pflichterfüllung erzeugen: Die Zeile würde im selben
   Kontext geschrieben, der den Satz aus den Quellen gelesen hat, und könnte eine Quelle statt
   eines Laufs nennen. Prüfen ließe sich nur, dass die Zeile da ist, nicht, dass sie trägt.
   Accepted-ADRs ändert dieses Verdikt nicht, und eine Folge-ADR dafür sehe ich nicht
   begründet.
3. **Ein Konformitätstest: die Festlegungssätze als Fixture-Fälle gegen das gepinnte
   Werkzeug, als `make`-Ziel.** Das wäre mechanisch. Er fängt aber eine andere Klasse: Ein
   Satz wird falsch, wenn das Werkzeug sich bewegt (Pin-Bump). Den ersten falschen Satz fängt
   er nicht, weil die Fälle erst entstehen, wenn jemand den Satz misst. Und die vier Belege
   kamen alle aus dem Schreiben, keiner aus dem Driften. Für Drift tragen schon die
   Re-Evaluierungs-Trigger von
   [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
   (an einer Änderung von d-check). Ein eigenes Ziel dafür wäre ein neuer Vorgang ohne
   Beleg im Register. Ich bestelle es nicht.

**Die verfügbare Falsifikation ist die Messung am Werkzeug.** In diesem Slice hat sie alle
sechs Abweichungen vor dem Merge gefunden. In allen vier Vorgängen hat kein falscher Satz
`main` erreicht. Das steht in den Beleg-Dateien und ist für diesen Zug *übernommen*, nicht
nachgemessen. Die tragende Linie arbeitet also. Was die Klasse kostet, sind Fixrunden, keine
falschen Aussagen auf `main`. Das ist ein **akzeptiertes Negativ**: Die Selbstprüfung beim
Schreiben bleibt die erste Linie und nicht die tragende.

## 4. Die Anpassung des Trägers

`.claude/commands/plan-welle.md` Schritt 6, Absatz „Festlegung im Pflichtenheft“, bekommt den
Punkt (c) **Messung am Werkzeug** (Anker `seit slice-spec-festlegungen-doku-gates`). Inhalt:
Für jeden Satz, der beschreibt, was ein Werkzeug prüft, meldet, nicht meldet oder ausgibt, gibt
es eine Zeile *Satz → Lauf* (Befehl, konstruierte Eingabe, gesehene Ausgabe, Werkzeug-Stand)
oder die Kennzeichnung *übernommen* bzw. *hergeleitet* nach §3.12. Ein Satz über eine
Reihenfolge oder eine Menge nennt die Zahl der Läufe.

**Warum das trotz Modul 6 Prosa sein darf:** Ausgeschöpft ist die Prosa-Form **dieser Regel**,
also §3.12 Instanz B. Der neue Absatz wiederholt sie nicht. Er korrigiert einen Schritt, der
den falschen Gegenstand verlangte, und macht aus der Messungstabelle eine Pflicht, die in
diesem Slice schon ohne Pflicht getragen hat (Closure-Notiz §7, Steering-Loop-Eintrag).
Prosa ist er trotzdem, und so steht er im Text.

**Reichweite auf den Bestand:** Fünf Folge-Slices der Reihe liegen in `open/`
(`slice-spec-festlegungen-code-gates`, `-commit-baseline-gates`, `-coverage-gates`,
`-kennungs-gates`, `-pruefer-hooks`). Gemessen mit `grep -ci` je Plan am Stand `11da9df8`:
„Gegenprobe“ kommt je 2- bis 4-mal vor, „Messung“ 0- oder 1-mal, „Mutation“ in keinem. Sie
beschreiben alle Werkzeuge. Hier tritt die Klasse am ehesten wieder auf. Ihre §2 nennen als
Beleg nur die Gegenprobe und die Anschluss-Frage. **Adresse:** Der Planner ergänzt beim
Übergang `open → next` je Plan im „Zu belegen durch“ von Liefer-Punkt 1 den Punkt (c). Das ist
ein Halbsatz je Plan, kein eigener Vorgang. Das Verdikt fasst die Pläne nicht an, weil
parallel an ihnen gearbeitet wird. Bis dahin fängt der Reviewer die Klasse, so wie er es
viermal getan hat.

## 5. Was dieses Verdikt nicht tut

- Es ändert **keine** Accepted-ADR, keinen Record unter `done/` und keinen Bericht unter
  `docs/reviews/` außer dieser Datei. `AGENTS.md` §3.12 bleibt im Wortlaut unverändert. Die
  Paraphrase in `state.md` („eine Aussage über einen Mechanismus nennt den Beleg-Anker“) ist
  eine zulässige Lesart von Instanz B („Tatsache über den Gegenstand“), auch wenn das Wort
  „Mechanismus“ dort nicht steht (`grep -c Mechanismus AGENTS.md`: `0`).
- Es legt **keinen** Beleg an und erhöht **keinen** Zähler.
- Es bestellt **keinen** Folge-Slice und **keinen** Sensor (§3).
- Die Notiz aus §7 der Closure (die Gegenprobe der Quellen fand keine der sechs Abweichungen)
  ist durch §2 und §4 beantwortet. Sie braucht keinen eigenen Register-Eintrag: Es ist kein
  zweiter Mechanismus, sondern die Erklärung, warum dieser Eintrag ein viertes Mal auftrat.

## 6. Nächstes Auftreten

Tritt die Klasse nach dieser Anpassung wieder auf, entscheidet die Form, wie es weitergeht.
Fehlt in einer Festlegung über ein Werkzeug die Zeile (c), ist der Träger nicht gelesen worden,
und das ist eine Frage an Reviewer und Planner. Steht die Zeile da und der Satz ist trotzdem
falsch, war der Lauf zu schmal (wie A-6). Dann ist die nächste Frage an den Architect, ob der
Konformitätstest aus §3 Punkt 3 Fälle erzwingen sollte, die nicht aus der Festlegung kommen.
