# Architect-Verdikt — Lesung von Exit 2 bei mehrdeutiger Stellung (`SPEC-038`, `slice-spec-festlegungen-harness-werkzeuge`, F-1)

**Datum:** 2026-10-07 · **Stand:** `941b1207` (Baum sauber) · **Rolle:** Architect (Modul 8),
Konflikt-Pfad, anderer Kontext als Implementer- und Reviewer-Lauf des Slice ·
**Rolleninhaber:** pt9912 (Architect-Agent im Auftrag) · **Anlass:**
[`review-slice-spec-festlegungen-harness-werkzeuge`](review-slice-spec-festlegungen-harness-werkzeuge.md)
F-1 (zweite Hälfte); der Implementer hat angehalten (Plan §3 „Fixrunde“), weil keine
Quell-ADR die Lesung „nicht messbar, Urteil am Diff“ trägt.

**Bezug:** [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
Entscheidung 3 und 4,
[`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
Entscheidung 2 und §Konsequenzen („fail-closed gewählt“),
[`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
Entscheidung 1 (c) und 2,
[`SPEC-038`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
[`AGENTS.md`](../../AGENTS.md) §3.5, §3.12.

---

## 0. Ergebnis in einem Blick

- **Verdikt 1 aus Modul 8: die Entscheidung gilt, die Festlegung hat falsch behauptet.**
  Lesung (a) gilt: Endet `make zitat-vergleich` wegen mehrdeutiger Stellung mit Exit 2,
  hat diese Seite **keine Einheit**, und die Korrektur **besteht nicht**
  (`ADR-0158` Entscheidung 3). Lesung (b) gilt nicht.
- **Keine Folge-ADR.** (a) folgt wörtlich aus `ADR-0158` Entscheidung 3 und
  `ADR-0162` Entscheidung 2; der Satz „Mehrdeutig heißt „nicht messbar““ in `SPEC-038`
  ist ein Fehler der Festlegung und wird berichtigt. Wer (b) will, ändert die Folge von
  Entscheidung 3 und braucht eine Folge-ADR; dafür gibt es heute keinen Fall (§3).
- **Vertrag `harness/targets/zitat-vergleich.md`:** kein Nachzug. Er verweist für die
  mehrdeutige Stellung und für die Lesung der Ausgänge nur auf `SPEC-038`.

## 1. Begründung

1. **Mehrdeutig ist eine Form von „keine Einheit“.** `SPEC-038` führt die mehrdeutige
   Stellung in der Aufzählung „Keine Einheit haben: …“ und in der Exit-Tabelle unter
   „eine Seite hat keine Einheit“. Für eine Seite ohne Einheit sagt `ADR-0158`
   Entscheidung 3: „Ist eine Einheit leer oder nicht lesbar (Anker unbekannt, Datei
   fehlt), besteht die Korrektur nicht.“ Die Klammer nennt Beispiele, keine
   abschließende Liste; eine Einheit, die das Werkzeug nicht sicher lesen kann, ist
   „nicht lesbar“.
2. **„Nicht messbar“ ist an den Referenten gebunden, nicht an das Werkzeug.**
   `ADR-0158` Entscheidung 4 und `ADR-0162` Entscheidung 1 (c) gelten, wenn die alte
   Adresse außerhalb des Repos liegt oder an keinem Stand auflöst — dann gibt es
   keinen Referenten. Bei mehrdeutiger Stellung liegt die Adresse im Repo und kann
   im Renderer auflösen; es fehlt nur die sichere Lesung durch das Werkzeug. Das ist
   eine Grenze der Nachbildung, kein fehlender Referent.
3. **`ADR-0162` Entscheidung 2 nennt den Fall fail-closed.** „Mehrdeutig mit Exit 2
   ist fail-closed und hält die Zusage „nie `cmp 0`“.“ Fail-closed heißt in dieser
   Kette (`ADR-0159` §Konsequenzen): „der Fehler ist eine Folge-ADR zu viel, nie eine
   stille Inhaltsänderung“. Lesung (b) machte aus demselben Exit 2 in den übrigen
   Abschnitten und an einem MR-Eintrag ein Bestehen nach Urteil, also fail-open.
4. **Der alte Vertrag hatte keinen Träger.** `abbe11b4:harness/targets/zitat-vergleich.md`
   band „mehrdeutig“ an „nicht messbar, Urteil am Diff“ und zitierte dafür `ADR-0159`
   Entscheidung 2. Deren Absatz gilt einer Adresse, die an keinem Stand auflöst; die
   Entscheidung ist durch `ADR-0161` abgelöst, ihr Zweig lebt in `ADR-0162`
   Entscheidung 1 (c) mit derselben Bedingung. Keine der drei ADRs nennt die
   mehrdeutige Stellung als Fall von „nicht messbar“.

## 2. Wortlaut für den Implementer

**`spec/pflichtenheft.md`, `SPEC-038`, Absatz vor „Zu `SPEC-038` — roh und
Normalisierung“.** Der Satz

> Mehrdeutig heißt „nicht messbar“.

wird ersetzt durch:

> Eine Seite in mehrdeutiger Stellung hat wie jede hier genannte Form keine
> Einheit: der Lauf endet mit Exit 2, und die Korrektur besteht nicht (Absatz
> „Stände“). Ein nicht messbarer Referent ist sie nicht; der Absatz dazu gilt
> nur für eine alte Adresse außerhalb des Repos oder ohne Stand, an dem sie
> auflöst.

**`spec/pflichtenheft.md` §8 Historie, Zeile 2026-10-07.** An das Ende der Zeile, nach
„… lässt die Korrektur nicht bestehen“, wird angehängt:

> , auch bei mehrdeutiger Stellung (keine Einheit, kein nicht messbarer Referent)

**`harness/targets/zitat-vergleich.md`:** keine Änderung.

## 3. Folge in der Praxis — gemessen

Unter (a) blockiert eine mehrdeutige Stellung am Referenten jede Zitat-Korrektur dorthin;
die Änderung geht dann den Weg einer Inhaltsänderung (Folge-ADR an einer `Accepted` ADR,
Nachfolger an einem MR-Eintrag). Ob es eine solche Stelle gibt, ist **gemessen** am Stand
`941b1207` mit `tools/harness/zitat-vergleich.sh`, je Anker derselbe Stand und dieselbe
Adresse auf beiden Seiten, `mehrdeutig` auf stderr gezählt:

| Messung | Ergebnis |
|---|---|
| jede `<a id="…"` in getrackten `.md` (`git grep -n -o '<a id="[^"]*"' -- '*.md'`, 96 Fundstellen samt Baseline) | 18 Exit 0, 78 Exit 2; davon 13 eindeutige Paare (Datei, `id`) mit `mehrdeutig`, alle in Inline-Code-Beispielen in `docs/reviews/architect-verdict-zitat-vergleich-werkzeug.md`, `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde-2.md`, `-fixrunde-3.md` und `harness/targets/zitat-vergleich.md` (`attr`, `x`, `X`, `…`, `qa` bis `qf`) |
| Links auf diese `id`s (`git grep -n -E '\]\([^)]*#(attr\|x\|X\|qa\|…\|qf\|…)\)' -- '*.md'`) | keine Zeile |
| jede Link-Adresse `](<pfad>.md#<anker>)` in getrackten `.md`, Pfad relativ zur Quelldatei aufgelöst (91 eindeutige Ziele, 84 davon als Datei am Stand vorhanden) | 84 Exit 0, keine `mehrdeutig` |

Ergebnis: Heute zeigt kein Verweis auf einen Anker in mehrdeutiger Stellung; die
mehrdeutigen `id`s sind Beispiele im Code-Span, und dort soll das Werkzeug gerade nicht
lesen. Die Strenge von (a) kostet heute nichts. **Grenze der Messung:** Verweise in
Backtick-Form ohne Link (`` `datei.md` §Abschnitt ``) und Link-Ziele mit Leerzeichen hat das
Muster nicht erfasst; die Exit-2-Fälle ohne `mehrdeutig` (andere Form, Inline-Code,
leere Einheit) sind nicht Gegenstand dieses Verdikts.

**Hergeleitet, nicht erprobt:** Tritt der Fall künftig auf und ist das Referenz-Dokument
änderbar, kann die `id`-Zeile dort in einem eigenen Commit vor der Korrektur in eine
eindeutige Stellung gebracht werden; die alte Adresse wird dann am jüngsten Stand gelesen,
an dem sie auflöst. Ob die Einheit dabei gleich bleibt, zeigt erst die Messung. Das ist
ein akzeptiertes Negativ ohne eigenen Vorgang: einmalig, heute ohne Fundstelle, und die
Messung färbt den Fall sichtbar rot.

## 4. Folgen für den Slice

- F-1 ist mit dem Wortlaut aus §2 vollständig aufgelöst; der Plan §3 „Fixrunde“ vermerkt
  dieses Verdikt als Grund.
- Keine neue ADR, kein Eintrag im ADR-Index, keine Rückführung des Slice.
