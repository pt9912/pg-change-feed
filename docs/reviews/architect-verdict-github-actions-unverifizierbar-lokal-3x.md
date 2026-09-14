# Architect-Verdikt: GitHub-Actions-Workflows lassen sich nicht lokal verifizieren — 3×

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal`
erreicht mit `slice-064` real 3× (`evidence/slice-039.md`,
`evidence/slice-056.md`, `evidence/slice-064.md`) — Lese-Schritt der
laufenden `welle-17`-Closure (Modul 6 §Wellen-Closure-Prozedur, Schritt 3;
Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b Verkörperung,
Planner → Architect → Planner-Zug).
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-14
**Bezug:** `ADR-0058` (Entscheidung 4, Ursprung der `slice-064`-Matrix),
`LH-QA-POR-001` (PostgreSQL-Major-Versionen, fachlicher Anlass der
Matrix — die Beobachtung selbst ist reine Ausführungs-/
Verifikationsdisziplin, keine fachliche Anforderung, dieselbe Wahl wie
bei den zitierten Präzedenzfällen), `ADR-0051` (CI/CD-Pipeline über
GitHub Actions), `AGENTS.md` §3.1 (Docker-only), §3.9 (verwandtes
Verifikations-/Sequenzierungsmuster), Modul 5 §Offene Risiken werden bei
Closure aufgelöst, Modul 6 §Das Beobachtungs-Register,
[`docs/reviews/architect-verdict-report-nackte-id-ohne-link.md`](architect-verdict-report-nackte-id-ohne-link.md)
und
[`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar.md`](architect-verdict-slice-chronik-in-code-kommentar.md)
(Präzedenzfälle für „trägt ein bestehender Mechanismus schon, oder
braucht es eine neue Verkörperung"),
`docs/plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal/`.

---

## Frage

Reicht der bestehende Umgang (jeder betroffene Slice trägt das Risiko
korrekt als *weiter offen* im §6-Abschnitt, der reale Beleg wird vor bzw.
bei Closure nachgetragen — dreimal empirisch beobachtet, dreimal
funktioniert), oder braucht diese Beobachtung — anders als beim
Nackte-Kennung-Fall — eine Verkörperung, weil hier (anders als dort)
**kein** mechanischer Sensor existiert, der den Risikopfad strukturell
schließen könnte?

## Unterschied zum Nackte-Kennung-Präzedenzfall: kein Sensor möglich, aber auch keine reine Interpretationsfrage

Im Nackte-Kennung-Fall trug `gestrichen`, weil ein deterministischer
Sensor (`make docs-check`s `ids`-Prüfung) den realen Risikopfad
(unentdeckter oder dauerhafter Verstoß) bereits vollständig schloss —
der Sensor kann jede nackte Kennung ohne Interpretationsspielraum
finden, und der einzige Weg, seinen roten Befund zu ignorieren
(Exit-Code-Maskierung), war bereits über `AGENTS.md` §3.9 verschlossen.

Hier ist die Lage doppelt anders:

1. **Kein Sensor ist möglich, aus einem strukturellen — nicht
   behebbaren — Grund.** GitHub Actions ist ein externer, gehosteter
   Dienst; ob ein Workflow auf dem echten Runner grün läuft, ist per
   Definition erst nach einem realen Lauf entscheidbar. Ein
   Docker-only/netzloser Sensor kann diese Frage nicht beantworten —
   nicht weil er noch nicht gebaut wurde, sondern weil sein
   Geltungsbereich (netzlos, `AGENTS.md` §3.1) den Gegenstand
   ausschließt. Das unterscheidet die Beobachtung auch vom
   Chronik-Fall: Dort war ein Sensor *denkbar*, aber wegen
   Interpretationsspielraum verworfen; hier ist er *nicht denkbar*.
2. **Es ist trotzdem keine Interpretationsfrage.** „Lief der Workflow
   auf dem echten Runner grün?" ist eine harte, binäre, extern
   beobachtbare Tatsache (`gh run view`/`gh run list`) — nur eben nicht
   lokal vor dem Push feststellbar. Es gibt also, anders als beim
   Chronik-Fall, keine Grauzone, die ein zweiter menschlicher/agentischer
   Blick besser auflösen könnte als ein erster; es gibt nur einen
   Zeitpunkt (vor vs. nach dem realen Lauf), zu dem die Antwort
   überhaupt existiert.

Diese Kombination — Sensor strukturell unmöglich, Tatsache aber hart und
extern beobachtbar — ist in den bisherigen Präzedenzfällen nicht
vorgekommen und braucht eine eigene Analyse.

## Empirischer Befund: 3/3 korrekt als „weiter offen“ geführt, 3/3 real (positiv) bestätigt — aber jedes Mal neu hergeleitet

| Vorgang | §6-Risiko-Führung | Realer Beleg | Wer bestätigte |
|---|---|---|---|
| `slice-039` | korrekt *weiter offen*, an `BEO-PGC/github-actions-unverifizierbar-lokal` verwiesen | zum Zeitpunkt der Closure noch nicht möglich (erster Workflow überhaupt — Henne-Ei) | — (Sache des Nutzers nach Merge, §1 explizit als Out-of-Scope benannt) |
| `slice-056` | korrekt *weiter offen* | vier unabhängige grüne Läufe | Reviewer **und** Verifier, unabhängig per `gh run view`/`gh run list` |
| `slice-064` | korrekt *weiter offen* | beide Matrix-Legs `completed`/`success`, Run `34822131377` | Planner-Koordinator **und** Verifier, unabhängig |

Kein einziger der drei Fälle zeigt einen Fehlschlag der Disziplin: In
allen drei Fällen wurde das Risiko korrekt nicht vorzeitig als
*entfallen* verbucht, und in den beiden Fällen, in denen zum
Closure-Zeitpunkt bereits ein realer Push möglich war (`slice-056`,
`slice-064`), wurde der reale Lauf tatsächlich unabhängig geprüft, bevor
die jeweilige Closure-Notiz geschrieben wurde. Das ist der Regelfall, für
den das Harness-Design sorgen soll — und er hat dreimal gehalten.

**Aber:** Nichts in `AGENTS.md`, `harness/README.md` oder einer Sensor-
Definition sagte bisher *explizit*, dass genau dieses Vorgehen verlangt
ist. Jeder der drei Fälle beruhte auf derselben, jedes Mal neu
angewendeten Schlussfolgerung — nicht auf einer bereits geltenden,
zitierbaren Regel. Das ist der Unterschied zum Nackte-Kennung-Fall: Dort
gab es bereits ein zweites Artefakt *und* einen Sensor, der es prüft;
hier gibt es kein Artefakt, das ein Sensor prüfen könnte, sondern nur
eine wiederholt richtig angewendete, aber nirgends kodifizierte
Verhaltensregel.

## Warum weder `gestrichen` noch `geplant` trägt

**Gegen `gestrichen`:** Der reale Risikopfad — ein Slice/eine Welle wird
als abgeschlossen behandelt, während ein neuer oder strukturell
geänderter Workflow tatsächlich rot laufen würde, unentdeckt oder erst
spät bemerkt — ist **nicht** durch einen bestehenden Mechanismus
strukturell geschlossen. Es gibt keinen Sensor (kann es strukturell
nicht geben, s.o.) und keine bereits existierende Hard Rule, die dieses
Verhalten erzwingt; es gibt nur eine bislang dreimal erfolgreiche, aber
rein diskretionäre Praxis. Anders als beim Nackte-Kennung-Fall wäre
„gestrichen“ hier eine Wette auf die vierte Wiederholung derselben
Diskretion, keine Feststellung einer bereits geschlossenen Lücke.

**Gegen `geplant`:** Es gibt keine ausstehende Entscheidung, die erst
noch getroffen werden müsste, und keinen sinnvollen Folge-Slice, der
etwas *baut* — die einzig sinnvolle Handlung ist, die bereits dreifach
bewährte Praxis in eine explizite, zitierbare Regel zu übersetzen.
`geplant` würde eine Kennung für etwas verlangen, das keine
Bau-Arbeit ist, sondern eine Prosa-Feststellung, die sich in diesem Zug
direkt schreiben lässt.

## Verdikt: verkörpert — als Hard Rule, kein Sensor, mit klarer Struktur-Begründung

**Ausgang: verkörpert.** Neue Hard Rule `AGENTS.md` §3.10 „Ein neuer
oder strukturell geänderter GitHub-Actions-Workflow gilt erst nach einem
realen, grünen Post-Push-Lauf als abgeschlossen“, mit Herkunfts-Anker
`seit welle-17`. Die Regel:

- benennt die strukturelle Ursache (externer Runner, nicht
  Docker-only/netzlos simulierbar, §3.1) explizit, statt sie implizit zu
  lassen — genau die „Tatsachenfeststellung über die Natur des
  Werkzeugs“, die künftige Slice-/Wellen-Planungen sonst jedes Mal neu
  herleiten müssten;
- macht **kein** Sensor-Versprechen, das nicht einlösbar wäre — die
  Begründung nennt ausdrücklich, warum keiner folgt;
- bindet die Pflicht an einen bereits bestehenden, bekannten Ort im
  Slice-Lifecycle (das §6-Risiko und seinen Ausgang, Modul 5), statt
  einen neuen Mechanismus zu erfinden — die Regel formalisiert nur, was
  in allen drei Belegen bereits richtig geschah.

Diese Verkörperung schließt keine neue Lücke im Sinn „ein Fehler ist
passiert, der jetzt verhindert wird“ — alle drei Belege fielen positiv
aus. Sie schließt eine andere, ebenfalls reale Lücke: dass diese
Erkenntnis bislang nirgends kodifiziert war und jede künftige
Workflow-Änderung sie erneut aus erster Prinzip herleiten müsste, ohne
Garantie, dass ein künftiger, weniger sorgfältiger Lauf dieselbe
Disziplin zeigt wie die drei gezählten. Das ist dieselbe Logik, die den
3×-Chronik-Fall trug (geschärfte Instruktion statt Sensor, weil kein
Sensor möglich ist) — nur mit einer noch härteren strukturellen
Begründung, warum kein Sensor je möglich sein wird.

## Was dieses Verdikt NICHT tut

- Kein neuer oder geschärfter Sensor — strukturell nicht möglich
  (externer, gehosteter Dienst; Docker-only/netzlos-Geltungsbereich
  schließt den Gegenstand aus).
- Keine Ergänzung in `.harness/skills/reviewer.md` — die Prüfung ist
  keine Interpretationsfrage, die ein zweiter Blick besser löst; sie ist
  eine Sequenzierungs-/Closure-Disziplin, die bei der Rolle sitzt, die
  Risiken abschließt (Planner/Verifier), nicht beim Reviewer.
- Kein neues ADR — reine Ausführungs-/Verifikationsdisziplin, keine
  Architektur- oder Vertragsentscheidung; `ADR-0051`/`ADR-0058` bleiben
  unverändert und werden nur als thematischer Kontext referenziert.
- Keine Duplizierung in `harness/README.md` — anders als beim
  Pipe-/Exit-Code-Fall (dort zusätzlich ein Kurzverweis in beiden
  Workflow-Abschnitten) ist diese Regel nicht Teil des generischen
  8-Schritt-Workflows, sondern eine GitHub-Actions-spezifische
  Ergänzung neben der bereits vorhandenen GitHub-Actions-spezifischen
  Hard Rule §3.8 (Action-Pinning) — derselbe Ort, dieselbe
  Geltungsbereichs-Logik.
- Der Register-Ausgang (`state.md`) wird von diesem Architect-Zug direkt
  mitgeführt, als Lese-Schritt der laufenden `welle-17`-Closure (Modul 8
  §Rollen-Sequenz für eine Welle, Schritt 3b) — die eigentliche
  `welle-17`-Closure (Schritte 3c–6) bleibt Sache des anschließenden
  Planner-Zugs.

Weder Produktionscode noch eine ADR-Datei noch eine Skill-Datei wurden im
Rahmen dieses Verdikts geändert. Geändert wurde ausschließlich
`AGENTS.md` (neue Hard Rule §3.10) sowie, als Teil desselben Zugs, der
Registereintrag `docs/plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal/state.md`.
