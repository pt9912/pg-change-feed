# Architect-Verdikt: DoD-Checkbox „Review durchgeführt" ohne Fixrunde — Instruktion vs. Sensor

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde`
erreicht mit `slice-047` real 3× (`evidence/slice-045.md`,
`evidence/slice-046.md`, Finding VF-4 in `docs/reviews/verify-slice-047.md`,
drittes eigenständiges Auftreten) — Lese-Schritt des wellenlosen Betriebs
(Modul 6 „Träger im Repo ohne Wellen"), ausgelöst direkt aus der
`slice-047`-Closure, ohne hostende Welle. Planner → Architect-Zug.
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** `LH-FA-SST-003` (thematisch nächste Kennung — die Beobachtung
wurde durch `slice-047` zum 3×-Übertritt gebracht, ist selbst aber keine
Retention-/Diagnose-Entscheidung), `docs/reviews/verify-slice-047.md`
Finding VF-4, `docs/plan/planning/observations/BEO-PGC/dod-checkbox-nachzug`
(Präzedenzfall, verkörpert seit welle-5, `architect-review-welle-5.md` Zug 2),
`docs/plan/planning/observations/BEO-PGC/dod-checkbox-nachzug-architect-pfad`
(verwandte, aber eigenständige Beobachtung — 1×, unter der Schwelle, hier
nicht Gegenstand), `.claude/commands/implement-slice.md` Schritt 18/21
(bestehende, hier strukturell nicht greifende Verkörperung),
`.harness/skills/reviewer.md` (Zielort dieses Verdikts).

---

## Frage

Dieselbe Frage wie beim Präzedenzfall `BEO-PGC/dod-checkbox-nachzug`
(dort: Implementer-Workflow), hier auf eine strukturell andere Ursache
angewandt: Reicht eine geschärfte Prosa-Instruktion — diesmal für die
**Reviewer**-Rolle statt die Implementer-Rolle —, oder braucht die
Beobachtung einen mechanischen Sensor?

## Warum die bestehende Verkörperung nicht greift

`BEO-PGC/dod-checkbox-nachzug` (3× bei `slice-015/016/017`, verkörpert seit
welle-5) hat den Nachzug an zwei Implementer-Schritten verankert: Schritt 18
(im selben Lauf, vor Handoff an den Reviewer) und Schritt 21 (nach
Reviewer-Findings, im Rahmen einer Fixrunde). Für die DoD-Zeile „Review
durchgeführt" greift **keiner** der beiden strukturell:

- **Schritt 18** kann die Checkbox nicht setzen — das Review findet
  *nach* dem Implementer-Handoff statt; der Implementer weiß im eigenen
  Lauf noch nicht, ob und wie das Review ausgeht.
- **Schritt 21** greift nur, wenn eine **Fixrunde** läuft — ein zweiter
  Implementer-Lauf, ausgelöst durch Reviewer-Findings, die eine
  Rückgabe verlangen. Bei einem sauberen oder nicht-fixrundenpflichtigen
  Review (0 HIGH, wie bei `slice-045`/`046`, oder MEDIUM/LOW ohne
  Rückgabe-Pfeil an den Implementer, wie bei `slice-047`) gibt es diesen
  zweiten Lauf nicht. Bei `slice-047` wurden die beiden Findings (F-1/F-2)
  sogar real behoben — aber direkt vom Planner (`fb6173d`), nicht über
  eine Implementer-Fixrunde; der Träger, an den Schritt 21 hängt, fehlte
  damit komplett.

Der Verifier findet die Lücke deshalb bei jedem derartigen Review erneut,
unabhängig vom Implementer-Pfad — strukturell identisch zum
`slice-chronik`-Präzedenzfall: eine bestehende Instruktion war nicht
falsch, nur an der falschen Rolle/dem falschen Zeitpunkt verankert.

## Verdikt

**Kein mechanischer Sensor. Geschärfte Instruktion — diesmal beim
Reviewer, nicht beim Implementer.** `.harness/skills/reviewer.md` bekommt
einen neuen Pflichtschritt: Kommt der Reviewer im eigenen Verdikt zu dem
Schluss, dass keine Fixrunde am Implementer nötig ist, zieht er die
DoD-Checkbox „Review durchgeführt" im selben Commit nach, der seinen
Report anlegt. Zielort und genauer Wortlaut: siehe Abschnitt „Umsetzung"
unten — die Änderung ist in diesem Zug bereits vorgenommen.

**Warum der Reviewer der richtige Träger ist:** Er ist die einzige Rolle,
die zum **richtigen Zeitpunkt** — beim Schreiben des eigenen Verdikts,
*bevor* irgendjemand sonst wieder auf den Slice-Plan schaut — bereits
weiß, ob eine Fixrunde kommt oder nicht. Das ist dieselbe Eigenschaft, die
beim Implementer für Schritt 18/21 trägt (er weiß im eigenen Lauf, was er
gerade tut), nur auf die Rolle übertragen, die hier tatsächlich das
auslösende Wissen hat.

## Geprüft: ist ein Existenz-Sensor praktikabel?

Ein Sensor der Form „jeder Slice in `done/`, zu dem ein
`docs/reviews/review-slice-<NNN>.md` existiert, muss die Checkbox `[x]`
tragen" wirkt auf den ersten Blick mechanisch robust — die Formulierung
der DoD-Zeile ist über **alle** 34 geprüften `done/`-Slices identisch
(`- [x]`/`- [ ] Review durchgeführt, Report unter \`docs/reviews/\` liegt
vor`, wörtlich aus der Vorlage übernommen, keine Abweichung gefunden). Das
unterscheidet diesen Fall vom `slice-chronik`-Präzedenzfall, wo schon die
Textform uneinheitlich war. Trotzdem trägt der Sensor aus drei Gründen
nicht:

1. **Offene Rollen-Menge auf der Erzeuger-Seite.** Ein Glob auf
   `review-slice-<NNN>.md` sieht nur den Reviewer-Pfad. Die separat
   geführte, verwandte Beobachtung
   `BEO-PGC/dod-checkbox-nachzug-architect-pfad` (1×, `slice-024`) zeigt
   real einen zweiten Pfad: Läuft die Review-Substanz über die
   Architect-Rolle (`docs/reviews/architect-review-slice-<NNN>.md`), gilt
   dieselbe Bedingung, aber mit einem anderen Dateinamensmuster. Ein
   Sensor, der beide Muster kennen müsste, hängt an derselben Fragilität
   wie beim Chronik-Sensor: Eine dritte Rolle (Validator, künftig denkbar)
   würde ihn erneut umgehen, ohne dass der Sensor selbst je „falsch" wird
   — er prüft nur, was er kennt.
2. **Der entscheidende Zustand ist nicht Existenz, sondern Verdikt.**
   Der Sensor müsste zwei Zustände unterscheiden, die beide real und
   legitim sind: (a) Review-Report existiert, keine Fixrunde nötig →
   Checkbox muss `[x]` sein; (b) Review-Report existiert, Fixrunde läuft
   noch → Checkbox **korrekt** `[ ]`, bis die Fixrunde abgeschlossen ist.
   Reine Datei-Existenz unterscheidet diese beiden Fälle nicht — dafür
   müsste der Sensor den **Verdikt-Text** des Reports parsen („keine
   Fixrunde … nötig" o. ä.), und dieser Wortlaut ist in den bisher
   geprüften Reports uneinheitlich formuliert (`review-slice-045.md`:
   „keine Findings"; `review-slice-047.md`: „keine Fixrunde am Code
   nötig"; andere Reports mit HIGH-Findings formulieren die
   Rückgabe an den Implementer wieder anders). Ein Sensor, der auf einem
   Prosa-Satz im Report keilt, ist dieselbe Art Klassifikations-Aufgabe,
   die schon beim Chronik-Verdikt als nicht mechanisierbar eingestuft
   wurde — hier sogar mit dem zusätzlichen Risiko, dass eine zufällig
   passende Formulierung im *falschen* Fall grün schaltet.
3. **Ein Sensor griffe strukturell zu spät.** Selbst wenn (1) und (2)
   gelöst wären: Ein Sensor gegen `done/` prüft **nach** dem `git mv` —
   die Hard Rule (`AGENTS.md` §3.3) verlangt aber, dass Inhaltsänderungen
   (Checkbox) *vor* dem reinen `git mv` passieren. Ein Sensor an dieser
   Stelle würde die Lücke erst melden, wenn sie bereits ein
   Ordnungsverstoß ist, statt sie vorher zu verhindern — er würde den
   Verifier-Fund (VF-1-Klasse) mechanisieren, aber nichts an der Ursache
   ändern: Ohne einen Träger *vor* dem `git mv`, der die Checkbox
   überhaupt setzen kann, bliebe der Sensor dauerhaft rot oder bräuchte
   seinerseits eine Nacharbeits-Routine — der Reviewer-Schritt ist genau
   dieser fehlende Träger.

**Fazit:** Der Sensor würde entweder (a) unvollständig bleiben, weil er
nur einen von mehreren Erzeuger-Pfaden kennt, oder (b) eine
Prosa-Klassifikation brauchen, die dieselbe Fragilität wie beim
Chronik-Fall hat, oder (c) zu spät greifen, um die Ursache zu beheben.
Die Instruktion beim Reviewer behebt die Ursache direkt, am Ort des
Wissens, ohne diese drei Probleme zu erben.

## Umsetzung

`.harness/skills/reviewer.md` bekommt den neuen Abschnitt „DoD-Checkbox-
Nachzug ohne Fixrunde" nach §Output-Schema — Umsetzung in diesem Zug,
kein Folge-Slice nötig (kleine Doku-Änderung an einer Prozess-Datei, kein
Produktionscode). Grenze: Der Reviewer rührt ausschließlich diese eine
Checkbox an — keine andere DoD-Zeile, keine Verifikations-Substanz (das
bleibt Verifier-Aufgabe, Modul 11); er bestätigt nur die Tatsache, dass
sein eigener, bereits abgeschlossener Arbeitsschritt stattgefunden hat.
Das ist kein Rollen-Übergriff auf „Verifikation gegen DoD" — der Reviewer
urteilt nicht über DoD-Konformität insgesamt, er trägt einen Fakt über
seinen eigenen Lauf nach.

**Herkunfts-Anker:** `seit slice-047`.

**Anmerkung zur eigenen Verantwortung dieses Verdikts:** Dieses Dokument
und die begleitende `.harness/skills/reviewer.md`-Änderung sind die
einzigen Artefakte, die der Architect-Zug hier ändert. `state.md` des
Registereintrags und die vierte Evidence-Datei
(`evidence/slice-047.md`) bleiben Planner-Arbeit (Modul 8: Lese-Schritt
und Zähler-Zuweisung liegen beim Planner) — dieses Verdikt liefert nur die
Entscheidung, auf die sich der Planner-Zug stützt.

Weder Produktionscode noch eine ADR-Datei wurden im Rahmen dieses
Verdikts geändert.
