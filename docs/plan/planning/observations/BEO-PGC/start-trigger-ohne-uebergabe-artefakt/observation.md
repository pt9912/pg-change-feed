# BEO-PGC/start-trigger-ohne-uebergabe-artefakt

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Lifecycle-Disziplin des Slice-Übergangs `next → in-progress`, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Slice-Plan macht einen Rollenwechsel (Planner→Architect,
Modul 8) explizit zur **Start-Bedingung** für den Übergang `next →
in-progress` — der Architect soll vor Implementierungsbeginn eine offene
Frage klären (hier: ob eine Bündel-Aufnahme eine eigene ADR braucht). Der
tatsächliche Übergangs-Commit trägt aber **kein** Übergabe-Artefakt für
diesen Zug und keinen Verweis auf eines — nur eine WIP-Limit-Notiz. Der
Implementer zieht die Schlussfolgerung anschließend **selbst**, in seiner
eigenen Rolle, ohne dass ein Architect-Kontext sie geprüft hätte (Modul 8:
„kein Rollenwechsel ohne Artefakt" — hier fehlt der Wechsel ganz, nicht nur
sein Artefakt).

**Unterschied zu benachbarten Klassen:** `BEO-PGC/architect-verdikt-rollen-scope-luecke`
betrifft ein **bereits existierendes** Architect-Verdikt, das eine Frage
unbehandelt lässt — hier findet der Architect-Zug **gar nicht statt**. Anders
als bei `BEO-PGC/dod-checkbox-nachzug-architect-pfad` (DoD-Häkchen bleibt
trotz sauberer Arbeit über die Architect-Rolle offen) ist hier die Arbeit
selbst — die Trigger-Bedingung des Plans — nicht eingelöst, nicht nur ihr
DoD-Beleg.

Deklaration: `slice-d-check-tracked-modul` (Review-Finding F-3, MEDIUM,
Review zu `slice-d-check-tracked-modul`): Der Slice-Plan §4
Trigger macht den Architect-Zug explizit zur Start-Bedingung; der reale
`next → in-progress`-Commit (`262bcda`, „WIP-Limit frei. Erster Slice der
Welle welle-d-check.") trägt keines. Nachträglich geheilt durch einen
Architect-Zug **nach** Implementierungsbeginn
(der Architect-Verdikt zur ADR-Frage von slice-d-check-tracked-modul) —
die Sequenz-Verletzung selbst bleibt als Record im Slice-Plan stehen, sie
wird durch den Nachtrag nicht rückwirkend geheilt (Verdikt §7).
