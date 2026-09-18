# Review-Report: slice-beispiele-start-architect-entscheidung — 2026-09-18

**Review-Art:** Architect-Entscheidung — Konsistenzprüfung einer neuen
Supersedes-ADR gegen `ADR-0076`/`ADR-0087`/`ADR-0090` (Modul 10 §Drei
Review-Arten), nicht gegen DoD (Verifier-Aufgabe).

**Gegenstand:** Commit `e491981` — neue
[`ADR-0098`](../plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Supersedes [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
in der Startform-Klausel von Festlegung 1), Start-Make-Target-Form über
Go/C#/Kotlin, Umgebungsdatei-Kontrakt einer neuen Demo-Umgebung; Slice
`slice-beispiele-start-architect-entscheidung`, Welle
`welle-beispiele-start-ueber-make`.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-beispiele-start-architect-entscheidung.md` (Plan)
- `docs/plan/planning/welle-beispiele-start-ueber-make.md` (Wellen-Kontext)
- `ADR-0076`, `ADR-0087`, `ADR-0090`, `ADR-0085` (zitierte Präzedenzfälle/Bestand)
- `AGENTS.md` §3.5, §3.9, §3.12 (Hard Rules)
- `docs/plan/adr/README.md` (Index-Nachzug)
- Drei abhängige `open/`-Slices (`slice-beispiele-compose-bootstrap.md`,
  `slice-beispiele-go-dockerfile-start.md`,
  `slice-beispiele-csharp-kotlin-start-target.md`)
- Eigene Messung: `<Dockerfile>.dockerignore`-Mechanismus real in einem
  Wegwerf-Baum reproduziert (Baseline-Bau ohne Zusatzdatei vs. mit
  `examples.Dockerfile.dockerignore` daneben vs. Kontroll-Bau `-f Dockerfile .`);
  `git show e491981 --stat` gegen §3.5-Verstoß geprüft; `ADR-0076` Festlegung 5
  und `ADR-0087` Festlegung 5 wörtlich nachgelesen; `compose.yaml`-Zitate
  zeilengenau nachgeschlagen; `make gates` selbst ausgeführt.

---

## Findings

### F-1 — Netzwerk-Ordering-Lücke, benannt aber nicht geschlossen

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md:322-351` (Festlegung 4)
- `befund`: Die ADR weist die Anlage des Netzwerks `cdc-examples` explizit
  `slice-beispiele-compose-bootstrap` zu (kein unadressierter Gap), benennt
  aber nicht, was `make example-run-go SURFACE=…` tut, wenn dieses Netzwerk
  noch nicht existiert (Nutzer startet ein Sprache-Target vor der
  Demo-Compose-Umgebung) — kein Fallback, keine Fehlermeldung spezifiziert.
- `verifizierbar`: nein — kein Gate; erst mit der Implementierung prüfbar.
- `klasse`: „Netzwerk-Ordering-Lücke, benannt aber nicht geschlossen"

## Negativbefunde

- geprüft, ohne Befund: `<Dockerfile>.dockerignore`-Mechanismus — real in
  einem Wegwerf-Baum reproduziert (Docker v29.8.0): kein halluzinierter
  Docker-Mechanismus, das ADR-Verhalten ist wie behauptet.
- geprüft, ohne Befund: `git show e491981 --stat` — `0076-*.md` erscheint
  **nicht** in der Dateiliste; die `Accepted`-ADR wurde nicht angefasst, kein
  `AGENTS.md` §3.5-Verstoß. `examples/**`, `Makefile`, `harness/mk/**`,
  `Dockerfile`, `.dockerignore` unverändert — reiner Entscheidungs-Zug ohne
  Umsetzung.
- geprüft, ohne Befund: `ADR-0076` Festlegung 5 wörtlich zitiert und exakt
  übereinstimmend („kein zweites Modul, kein Ausschluss aus `go test ./...`").
- geprüft, ohne Befund: `ADR-0087` Festlegung 5 (`CDC_CONFIG_FILE`-Verbot)
  wörtlich geprüft — die in `ADR-0098` gezogene Unterscheidung
  (Shell-/Compose-Ebene vs. programminterner Parser) widerspricht der Klausel
  nicht, sauber begründet.
- geprüft, ohne Befund: `ADR-0085` Merkmal (ii) („genau eine Datei") — die
  Feststellung, die dortige Ausnahme-Klasse passe hier nicht (Verzeichnisse,
  nicht eine Datei), ist korrekt.
- geprüft, ohne Befund: zwei `compose.yaml`-Zitate (`CDC_API_TOKEN_READER`
  Zeile 108, `networks.cdc-test` Zeilen 145-149) — beide wortgleich.
- geprüft, ohne Befund: `SPEC-023` „Laufzeit und Startform" — bereits
  sprachneutral, keine Änderung nötig, ADR-Fund korrekt.
- geprüft, ohne Befund: drei abhängige `open/`-Slices — Arbeitsnamen
  konsistent auf die festen Namen aus `ADR-0098` gezogen
  (`example-run-go/-csharp/-kotlin`, `examples/.env`, `cdc-examples`),
  Out-of-Scope sauber getrennt, jeder unabhängig startbar.
- geprüft, ohne Befund: ADR-Index (`docs/plan/adr/README.md`) — Zeile für
  `ADR-0098` korrekt, Titel-Zellenlänge im Rahmen.
- geprüft, ohne Befund: DoD des Slice-Plans — LP1-LP3, `make gates`,
  Doku-Update-Begründung, Closure-Notiz-Entwurf, Risiko-Ausgänge (§6)
  konsistent mit dem committeten Ergebnis.
- geprüft, ohne Befund: `make gates` — real ausgeführt, Exit-Code direkt
  (ungepiped) geprüft: Exit 0.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Netzwerk-Ordering-Lücke, benannt aber nicht
geschlossen

## INFO

- **INFO-1** — Zwei zuvor an den Architect-Zug gemeldete
  „Oversized-Title-Cell"-Befunde des ADR-Index sind im aktuellen Stand nicht
  mehr auffindbar (weder als offene Zelle noch als dokumentierte Korrektur) —
  der aktuelle Zustand der Tabelle ist unauffällig, bestätigt aber nur den
  Ist-Stand, nicht explizit die Behebung der zuvor gemeldeten Instanz.
- **INFO-2** — `examples/.env` existiert im Repo noch nicht (korrekt an
  `slice-beispiele-compose-bootstrap` delegiert). Die in `ADR-0098`
  Festlegung 3 zugesagte Eigenschaft „trägt keine echten Zugangsdaten" ist
  eine **Zusage**, keine geprüfte Tatsache (`AGENTS.md` §3.12) — im ADR-Text
  korrekt als Review-Prüfpflicht für den umsetzenden Zug re-adressiert.
- **INFO-3** — Die Docker-Probe der ADR lief mit einer kleineren,
  synthetischen `.dockerignore` in einem Wegwerf-Baum (nicht der realen
  Repo-`.dockerignore`) — von der ADR selbst korrekt als synthetische
  Reproduktion deklariert, keine Drift zur realen Datei.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Das eine LOW-Finding (F-1)
ist ein benannter, an die Implementierung delegierter Punkt (Fehlerverhalten
bei fehlendem `cdc-examples`-Netzwerk), keine Rückkante an den Architect-Zug.

Die DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/` liegt
vor" ist mit diesem Report erfüllt — keine Fixrunde nötig
(Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde: bei isolierten
LOW/INFO-Findings genügt Annahme oder Begründung; F-1 wird als offener
Implementierungs-Hinweis an `slice-beispiele-csharp-kotlin-start-target.md`
und `slice-beispiele-go-dockerfile-start.md` weitergegeben, nicht an den
Architect-Zug zurückgespielt).

Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11), sofern dieser Slice eine eigene Verifikation
vorsieht.
