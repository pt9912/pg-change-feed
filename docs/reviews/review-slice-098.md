# Review-Report: slice-098 — 2026-09-17

**Review-Art:** Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `git diff 9b6b4fe..HEAD` (Commits `9cc65d0`, `c4d6a11`,
`c87790b`, `7744a68`) — C#-Sprachwurzel und HTTP-Client, erste Zelle der
Beispiel-Client-Matrix.

**Skill:** `.harness/skills/reviewer.md` @ Accepted (Schärfung 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-098-csharp-sprachwurzel-http-client.md`
  (vollständig, insb. §1 Abgrenzung, §2 DoD)
- `ADR-0090` (Umfang volle Matrix, Festlegung 4 Abhängigkeits-Pins,
  Festlegung 5 nicht-blockierender Workflow, Festlegung 7 Handbuch-Form,
  §Slice-Schnitt-Empfehlung Zeile 1)
- `ADR-0087` (Festlegung 2/3 Sprach-Wurzel/Bau-Kontext, Festlegung 4
  Werkzeug-ohne-Gate, bestätigt für diesen Slice)
- `AGENTS.md` §3.1 (Docker-only), §3.6 (Pin-Hebung = bewusster Commit),
  §3.8 (Action-Pinning), §3.9 (Exit-Code-Disziplin), §3.10 (Workflow braucht
  realen Post-Push-Lauf), §3.12/§3.13 (Beleg-Herkunft, Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema; Modus-Deklaration `*`/PGC = GF)
- `LH-FA-SST-006` (HTTP-/JSON-API — berührte Spec-Stelle laut Slice-Kopf)

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings. Zwei INFO-Hinweise.

### F-1 — `SPEC-023` „Sprachen und Umfang" bleibt hinter dem Accepted-`ADR-0090` zurück

- `kategorie`: INFO
- `quelle`: `ADR-0090` §7 Festlegung 7 / §Konsequenzen „Folgepflicht (Spec-Zug)"
- `pfad`: `spec/pflichtenheft.md:460`
- `befund`: Die Zeile *Sprachen und Umfang* trägt weiterhin den `ADR-0087`-Wortlaut
  („C# und Kotlin die **HTTP-Familie** … weitere Oberflächen **perspektivisch**"),
  obwohl `ADR-0090` (Accepted, 2026-09-17, vor Beginn dieses Slice) genau diese
  Beschränkung als aufzulösen benennt und die Aktualisierung als eigene
  Folgepflicht führt. Die Lücke stammt nicht aus diesem Diff — `ADR-0090`s
  eigener Annahme-Commit (`72d026b`) hat `spec/pflichtenheft.md` nicht
  angefasst, und dieser Slice-Diff tut es ebenfalls nicht (`git diff
  9b6b4fe..HEAD -- spec/pflichtenheft.md` ist leer). Der Slice-Plan grenzt
  `SPEC-023` weder als Liefer-Punkt noch als §1-Ausschluss ein; „Berührte
  Spec-Stellen" nennt nur `LH-FA-SST-006`. Da diese Zeile außerhalb des
  Diffs liegt, bleibt sie nach Skill-Konvention INFO (Leser ist die Messung,
  nicht die Form) — Planner-Aufmerksamkeit empfohlen, damit die Folgepflicht
  nicht zwischen den sieben Matrix-Slices verloren geht.
- `verifizierbar`: ja — `grep -n "Sprachen und Umfang" spec/pflichtenheft.md`
- `klasse`: „ADR-Folgepflicht ohne zugewiesenen Träger-Slice"

### F-2 — Escaping-Divergenz zum Go-Formvorbild bei Leerzeichen in Query-Parametern

- `kategorie`: INFO
- `quelle`: Maintainability (Form-Vorbild-Treue)
- `pfad`: `examples/csharp/http-client/TablesUrlBuilder.cs:16` vs.
  `examples/http-client/tables.go`
- `befund`: `Uri.EscapeDataString` kodiert ein Leerzeichen als `%20`
  (Test `BuildEscapesReservedCharacters` erwartet
  `quelle%20%26%20test`), während Gos `url.QueryEscape` im Formvorbild
  `+` verwendet (`quelle+%26+test`). Beide Formen sind gültige
  Query-Kodierungen und werden von Standard-Query-Parsern gleich
  dekodiert — kein Funktionsfehler, nur eine beobachtbare Abweichung vom
  wörtlichen Vorbild.
- `verifizierbar`: ja — Testausgabe beider Sprachen vergleichen
- `klasse`: „Formvorbild-Divergenz ohne Funktionsfehler"

## Negativbefunde

- geprüft, ohne Befund: `examples/csharp/Dockerfile` (Digest-Pinning beider
  Basen — SDK `sha256:60a2b2230a...`, Runtime `sha256:e6541e52ae...`, je 64
  Hex-Zeichen, `FROM ...@sha256:...`-Form; SDK-Digest deckungsgleich mit dem
  in `ADR-0087` Festlegung 3 gemessenen Kandidaten)
- geprüft, ohne Befund: `examples/csharp/Directory.Packages.props` (drei
  exakte Versionen `18.10.1`/`2.9.3`/`4.0.0`, keine Ranges/Wildcards;
  `.csproj`-Dateien referenzieren ohne eigene Version, zentral verwaltet)
- geprüft, ohne Befund: `examples/csharp/http-client/**` Produktionscode
  (`Cli.cs`, `Config.cs`, `Program.cs`, `TablesUrlBuilder.cs`,
  `TablesClient.cs`) — Kommentare tragen ADR-/LH-Bezüge und Form-Vorbild-Zeiger,
  keine Chronik-Sprache
- geprüft, ohne Befund: `examples/csharp/http-client/HttpClient.Tests/**` —
  Tests decken Env-Defaults, Flag-Override, unbekanntes Flag, fehlender
  Flag-Wert, URL-Aufbau samt Escaping und den Fehlerpfad ohne erreichbaren
  Host; mehr Fälle abgedeckt als das Go-Formvorbild (`tables_test.go` prüft
  nur den URL-Aufbau)
- geprüft, ohne Befund: `harness/mk/examples.mk` und `Makefile`-Einbindung
  (`include harness/mk/*.mk`, `examples-csharp` erscheint nirgends in
  `GATE_CHECKS`)
- geprüft, ohne Befund: `.github/workflows/examples.yml` — einzige
  `uses:`-Zeile SHA-gepinnt mit Tag-Kommentar (identischer SHA wie
  `ci.yml`/`e2e.yml`), `pull_request`/`push`-Trigger ohne Required-Status-Check
  im Dateiinhalt, kein Job-Name, der eine Verwechslung mit `ci.yml`
  nahelegt, ruft ausschließlich `make examples-csharp` auf
- geprüft, ohne Befund: `.a-check.yml` — `git diff 9b6b4fe..HEAD -- .a-check.yml`
  ist leer, wie in Plan §1 zugesagt
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version 1.19→1.20,
  neue Änderungshistorie-Zeile 1.20, `**Beispiel:**` → `**Beispiele:**`-Block
  mit je einer Zeile Go/C#, entspricht `ADR-0090` Festlegung 7 (ein Block je
  Oberfläche, eine Zeile je Sprache statt eines Absatzes je Programm)
- geprüft, ohne Befund: `harness/README.md` §Werkzeuge — neue Zeile für
  `make examples-csharp`, Bindung auf `ADR-0087`/`ADR-0090`, Anker
  `· seit slice-098`, `kein Gate` explizit genannt
- geprüft, ohne Befund: `examples/csharp/.gitignore` — Plan-Nachzug (§3 der
  Slice-Plan-Datei) korrekt begründet und dokumentiert, Umfang passt zum
  bestehenden Ziel (schützt nur gegen versehentlichen lokalen `dotnet build`)
- geprüft, ohne Befund: alle vier Commit-Messages tragen `ADR-0087`/
  `ADR-0090` im Body, kein `SPEC-*`/`ARC-*`-Präfix im Betreff
  (`AGENTS.md` §5, Traceability-Regel)
- geprüft, ohne Befund: repo-weiter `§3.13`-Gegencheck — kein Fund für
  „nur Go"/„kein C#-Client" außerhalb der ADR-Kontextabschnitte selbst
  (`grep -rln "examples/http-client" --include="*.md" .` zeigt nur erwartete
  Treffer; `docs/user/benutzerhandbuch.md` ist bereits nachgezogen)
- geprüft, ohne Befund: `docs/plan/planning/observations/BEO-PGC/` — alle
  vier im Slice-Kopf/§8 zitierten Verzeichnisse
  (`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  `handbuch-versionshistorie-uebersprungen`,
  `nicht-blockierender-workflow-alarmmuedigkeit`,
  `github-actions-unverifizierbar-lokal`) existieren real

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** ADR-Folgepflicht ohne zugewiesenen
Träger-Slice · Formvorbild-Divergenz ohne Funktionsfehler

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW; die zwei INFO-Hinweise
lösen keine Fixrunde am Implementer aus (F-1 liegt außerhalb des Diffs und
ist Planner-Angelegenheit einer künftigen Zug-Zuordnung; F-2 ist ein
Formhinweis ohne Funktionsfehler).

**Übergabe:** Da keine Fixrunde folgt, zieht dieser Report nach
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde die Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan
selbst nach — im selben Commit, der diesen Report anlegt. Die Finding-Klassen
gehen in die Slice-Closure §7 und von dort in den Steering-Loop-Zähler
(F-1: Planner prüft bei nächster Gelegenheit, welcher der sieben
Matrix-Slices `SPEC-023` §Sprachen-und-Umfang nachzieht). Dieser Report ist
Lauf-Beleg und ersetzt keine Verifikation (Modul 11, separater Kontext).
