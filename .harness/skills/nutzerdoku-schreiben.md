# Skill — Nutzerdokumentation schreiben — PG Change Feed

* Status: Accepted
* Bezug: [`ADR-0143`](../../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
  (Gate und Skill) · `AGENTS.md` §3.7 (Ist-Zustand statt Chronik) und §3.12
  (Ursprung einer Zahl) · [`docs/user/benutzerhandbuch-standard.md`](../../docs/user/benutzerhandbuch-standard.md)
  (Form des Handbuchs)
* Gilt für: wer ein Dokument unter `docs/user/` schreibt oder ändert, das sich an
  Betreiber und Integratoren richtet (`benutzerhandbuch.md`,
  `benutzerhandbuch-standard.md`, `version.md`); kein eigenes Make-Target —
  getragen vom Implementer (`.claude/commands/implement-slice.md` Schritt 17) und
  vom Reviewer (`.harness/skills/reviewer.md`, HIGH-Punkt Handbuch-Versionshistorie).
  Nicht Gegenstand: die Maintainer-Doku unter `docs/maintainer/` und die vier von
  Runnern geschriebenen `*-abdeckung.md`.

Ein Skill, keine Rolle: er ergänzt keine Hard Rule in `AGENTS.md`, er verweist auf
sie. Das Gate `make handbuch-public-doc-check` ist das Fangnetz (Vertrag:
[`harness/sensors/handbuch-public-doc-check.md`](../../harness/sensors/handbuch-public-doc-check.md)).

## Regeln beim Schreiben

1. **Betreiber- und Integratorsicht.** Der Text sagt, was der Nutzer tut und was
   er dabei sieht, nicht wie die Software gebaut ist. Entwicklersprache
   (Pakete, Interna, Entscheidungswege) gehört in Plan und Code.
2. **Ist-Zustand statt Chronik.** Ein Satz beschreibt, was gilt. „Früher …
   jetzt …“, „neu seit …“ im Fließtext und die verworfene Alternative stehen
   nicht im Handbuch; die Vorgängerfassung hält `git`, die Änderung je Version
   die Änderungshistorie.
3. **Keine internen Kennungen.** Weder Anforderungs-, Spezifikations-,
   Architektur- noch Entscheidungskennungen, weder Slice- noch Welle- noch
   Register-Namen. Das gilt auch für Befehlsbeispiele und Codeblöcke. Wer den
   Anlass einer Änderung wissen will, findet ihn in Commit und Plan; das Handbuch
   braucht ihn nicht.
4. **Keine Links in den Planungs- und Review-Bestand.** Verweise führen auf
   andere Nutzerdokumente, nicht nach `docs/plan/` oder `docs/reviews/`.
5. **Zahlen tragen ihren Ursprung als Wort.** Eine Messung heißt „gemessen“,
   eine fremde Zahl „übernommen“, eine Rechnung „abgeleitet“
   (`AGENTS.md` §3.12). Den Lauf nennt der Text als Befehl und gedruckte Zeile
   (`make <Ziel>`), nie als Berichts- oder Slice-Kennung.
6. **Neue Betreiber-Oberfläche zieht das Handbuch mit** (Umgebungsvariable,
   `cdc.*`-Funktion, Endpunkt): im selben Diff in den Abschnitt der Umgebungsvariablen
   bzw. der Aufgaben, oder der Aufschub trägt eine Adresse
   (`.claude/commands/implement-slice.md` Schritt 17).

## Version und Änderungshistorie

Ändert ein Diff `docs/user/benutzerhandbuch.md` inhaltlich (neuer Abschnitt, neue
Umgebungsvariable, geänderte Beschreibung), zieht derselbe Diff den
`Version:`-Kopf hoch **und** ergänzt je Version genau eine Zeile in
`### Änderungshistorie`. Die Zeile steht in Betreibersicht — was ändert sich für
den Betreiber — und ohne interne Kennung; sie ist keine Chronik des Vorgehens.
Eine reine Korrektur der Version oder der Historie braucht keine neue Zeile.

Kandidatenlauf vor der Übergabe (erste, nicht tragende Linie; die tragende ist
der unabhängige Reviewer):

```text
git diff --name-only <Basis> -- docs/user/benutzerhandbuch.md
git diff <Basis> -- docs/user/benutzerhandbuch.md | grep -E '^\+Version:|^\+\| [0-9]+\.[0-9]+ \|'
```

Trifft der erste Befehl, müssen beide Muster des zweiten treffen; fehlt eines,
wird Version oder Historie vor dem Handoff nachgetragen.

## Prüfen

- `make handbuch-public-doc-check` läuft ohne Treffer (Exit 0, ungefiltert
  ausgewertet, `AGENTS.md` §3.9).
- Die Grenze des Gates: es liest Kennungen, nicht Sinn. Chronik-Sprache ohne
  Kennung, Entwicklersicht und eine Historie-Zeile ohne Betreibernutzen bleiben
  grün. Diese Hälfte prüft der Autor beim Schreiben anhand der Regeln oben und
  der Reviewer beim Lesen des Diffs.
- Eine neue Datei unter `docs/user/` färbt das Gate rot, bis sie im Skript als
  geprüft oder ausgenommen geführt ist; die Wahl (Nutzerdoku oder Erzeugnis) ist
  eine Entscheidung, keine Formsache ([`ADR-0143`](../../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Festlegung 9).
