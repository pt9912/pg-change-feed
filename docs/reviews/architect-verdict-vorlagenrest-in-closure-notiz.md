# Architect-Verdikt: Vorlagenrest in der §7-Closure-Notiz

**Rolle:** Architect (Modul 8)
**Anlass:** Der Lese-Schritt der `welle-19`-Closure (Modul 6
§Das Beobachtungs-Register) fand eine Klasse, die mit dieser Welle die
3×-Schwelle erreicht: `slice-069`…`slice-072` ließen den
Vorlagen-Guidance-Block der §7-Anker-Zeile stehen, obwohl mit ihnen nichts
verkörpert wurde. Der Planner-Lauf konnte den fälligen Architect-Zug nicht
selbst dispatchen (kein Agent-Werkzeug) und hat den Fall in
`BEO-PGC/plan-vorlagen-defekt/state.md` vermerkt. Dieser Lauf führt den
Rollenwechsel nach (Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b —
Verkörperung ist eine Entscheidung, keine Planung).
**Rolleninhaber:** pt9912 (dieser Lauf)
**Datum:** 2026-09-15
**Bezug:** `BEO-PGC/vorlagenrest-in-closure-notiz` (neu),
`BEO-PGC/plan-vorlagen-defekt` (verwandter, verengter Eintrag),
`done/welle-19-results.md`, Modul 6 §Das Beobachtungs-Register,
`grundlagen-traceability.md` §Herkunfts-Anker, `.d-check.yml`
(Zielort der Verkörperung), `.claude/commands/implement-slice.md`
(Zielort 1).

---

## Frage

Welchen Ausgang bekommt die Klasse — Erweiterung des verwandten Eintrags,
eigener Eintrag, `gestrichen` oder Stehenlassen (`benannt, nicht gezählt`)?
Und wenn `verkörpert`: an welchem Zielort, und trägt dort ein Sensor oder
eine Disziplin-Zeile die Prüfung?

## Verglichene Alternativen

**1. Den bestehenden Eintrag `plan-vorlagen-defekt` erweitern — verworfen.**
Die beiden sind verwandt (beide Male bleibt Vorlagen-Text stehen, den kein
Gate fängt), aber der bestehende Eintrag ist auf §2 (DoD) und den Fill der
Wellen-Eröffnung **verengt**, und seine `observation.md` ist „unveränderlich
ab Anlage" (Modul 6). Eine Erweiterung müsste diese Identität überschreiben
oder §7-Belege unter einen Text hängen, der §7 nicht nennt — beides
verfälscht die Deklaration. Das Register trennt nach Beobachtung, nicht nach
Ähnlichkeit; die beiden Zählräume bleiben getrennt.

**2. Eigener Eintrag `BEO-PGC/vorlagenrest-in-closure-notiz` — gewählt.**
Eigene Stelle (die §7-Closure-Notiz statt §2), eigener Fill-Moment (die
Slice-Closure statt die Wellen-Eröffnung), eigener Zielort — damit eine
eigene Beobachtung im Sinn der Register-Regel.

**3. `gestrichen` — verworfen.** Die Klasse ist strukturell: Die Vorlage
bleibt unverändert (sie ist vendored und `baseline-verify`-gepinnt), §7 wird
bei jedem Slice neu aus ihr kopiert, und kein Gate fängt den Rest. Sie kann
weiter auftreten — `slice-077` trägt den Block in `open/` ungefüllt.

**4. Stehenlassen (`benannt, nicht gezählt`) — verworfen.** Die Regel
„ein Vorkommen ohne abgeschlossenen Vorgang bekommt keinen Beleg" trägt die
Zurückhaltung des Planner-Laufs **vor** dem Architect-Zug, nicht danach:
Vier abgeschlossene Vorgänge sind vier Belege, die Klasse steht über der
Schwelle, und ein Eintrag, der den Ausgang der Schwelle schuldig bleibt, ist
genau der Fall, den der Lese-Schritt ausschließt.

## Befund: vier Belege, nicht drei

Die `welle-19`-Closure nannte drei Vorgänge (`slice-070`…`slice-072`, die den
**ganzen** Block stehen ließen — Teil-Zeile `— liegt in …`,
`Auslöser:`-Platzhalter und kursiver Hinweis). Dieser Zug fand einen
**vierten**: `slice-069` behielt nur den kursiven Hinweis; Anker-Zeile und
`Auslöser:`-Zeile waren entfernt. Genau diese Form ist der Grund, warum die
Anker-Paarung ihn nicht sah — sie löst allein über das Feld `liegt in` aus,
und `slice-069` trägt dieses Feld nicht mehr. Der vierte Beleg ist gezählt
(`evidence/slice-069.md`); damit stehen vier Belege, nicht drei.

Das ist kein Zähl-Fehler, sondern die Bestätigung der Diagnose: **eine
Prüfung, die an das Feld `liegt in` gebunden ist, ist blind für genau die
halb-aufgelösten Reste.** Der Rest ist ein Form-Phänomen, kein Feld-Vorkommen.

## Verdikt: Sensor als tragende Linie, Selbstprüfung als erste

Zwei Linien, keine Reviewer-HIGH:

1. **Tragend, mechanisch:** die fünfte `structure`-Regel in `.d-check.yml`
   (Dateiklasse `docs/plan/planning/done/slice-*.md`, Abschnitt
   `## 7. Closure-Notiz`) bekommt ein `forbid-pattern` auf den Ausfüll-Hinweis.
   `make docs-check` färbt einen stehengebliebenen Hinweis als
   `section-forbidden` rot. Diese Linie ist unabhängig vom schreibenden
   Kontext — sie greift an der Ablage, nicht am Lauf.
2. **Erste, nicht tragend:** die datei-skopierte Selbstprüfung im
   Planner-Closure-Schritt (`.claude/commands/implement-slice.md` Schritt 24)
   auf dem **rohen** Text, damit der Rest schon vor dem `git mv` sichtbar wird.

**Kein Reviewer-HIGH-Punkt — geprüft und verworfen.** Bei den bisherigen
3×-Klassen war die Prüfung ein Urteil (Satz-Subjekt, Diff-Inhalt) und der
unabhängige Reviewer deshalb die tragende Linie; hier ist sie mechanisch. Der
Closure-Note-Reviewer schließt eine Dopplung des Struktur-Gates ausdrücklich
aus, und das Modul 8 warnt vor einer Mehrfachzuweisung ohne **anderen**
Eingabe-Kontext. Eine zweite, inferentielle Linie auf dieselbe Zeichenkette
wäre doppelte Arbeit mit denselben blinden Flecken.

## Der Sensor und seine Grenze

Möglich **und** präzise: Der Rest ist Vorlagen-Guidance, keine Prosa; er ist
über den Hinweis-Text greifbar, der in **allen vier** Belegen mitging.

Grenze, benannt statt verschwiegen: Die `structure`-Bedingungen lesen den
**bereinigten** Abschnittstext, in dem Inline-Code-Spans geleert werden. Ein
stehengebliebener Zielort-Platzhalter (`<…>` im Backtick-Zielort der
Anker-Zeile) ist damit **unsichtbar** — eine Regel allein auf ihn wäre blind
gegen den Normalfall. Der gewählte Marker (der Hinweis-Text) trifft alle vier
Belege; die rohe `grep`-Selbstprüfung deckt die Rest-Grenze (Hinweis
entfernt, Platzhalter-Zeile geblieben) ab. Die Grenze steht in der `state.md`
des Eintrags.

Keine Bestands-Ausnahme nötig: Nach dem Entfernen des `slice-069`-Rests ist
die Regel gegen den ganzen `done/slice-*.md`-Korpus grün (gemessen: 574
Dateien, 0 Befunde). Der Alt-Bestand (`welle-1`, `welle-2`, `slice-007`)
trägt andere Reste an anderer Stelle und wird von dieser Regel nicht
berührt — er bleibt, wie die `welle-19`-Closure ihn festgestellt hat.

## Umsetzung

1. `.d-check.yml` — `forbid-pattern` auf der fünften `structure`-Regel
   ergänzt, mit Herkunfts-Anker `· seit welle-19` im Kommentar. Umgesetzt.
2. `.claude/commands/implement-slice.md` Schritt 24 — die Selbstprüfung mit
   Kandidatenlauf. Umgesetzt.
3. `docs/plan/planning/done/welle-19/slice-069-grpc-streaming-adapter-grundgeruest.md`
   — der restliche Hinweis entfernt (dieselbe Klasse wie `041dc0c`; reine
   Guidance, nichts Tragendes). Umgesetzt.
4. Register — Neuanlage `BEO-PGC/vorlagenrest-in-closure-notiz` mit vier
   Belegen, Ausgang `verkörpert`, Herkunfts-Anker `seit welle-19`; der
   Nachtrag in `BEO-PGC/plan-vorlagen-defekt/state.md` zeigt auf den neuen
   Eintrag. Umgesetzt.

## Was dieses Verdikt NICHT tut

- Kein neues ADR — Prozess-/Instruktions-Schärfung, kein Architektur- oder
  Vertragsgegenstand; Produktionscode ist nicht berührt.
- Keine Änderung an der vendorten Baseline (`.harness/baseline/v6.5.0/…`) —
  die Vorlage ist referenziert und `baseline-verify`-gepinnt.
- Keine Umschreibung geschlossener Vorgänge über die Guidance hinaus — nur
  der `slice-069`-Rest ist entfernt, nichts Tragendes.
- Kein Reviewer-HIGH-Punkt — siehe Verdikt.
