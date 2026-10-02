# Review-Report: slice-sdk-public-doc-check-lesefehler-fail-closed — 2026-10-03

**Review-Art:** Code — Diff gegen Plan, [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
(`Accepted`, unberührt), Vorbild [ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
und Hard Rules. DoD-Vollständigkeit bleibt Verifier-Aufgabe.

**Gegenstand:** Diff `8e38389a..8f6430a0` (Implementer-Commit `8f6430a0`, vier Dateien:
`tools/harness/sdk-public-doc-check.sh`, `tools/harness/run-sdk-public-doc-check-tests.sh`,
`harness/sensors/sdk-public-doc-check.md`, der Slice-Plan).

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-public-doc-check-lesefehler-fail-closed.md`
- [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md), [ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
- `AGENTS.md` §3.1, §3.2, §3.6, §3.7, §3.9, §3.12, §3.13, §3.15; `tools/harness/handbuch-public-doc-check.sh` (Vorbild)
- Entscheidung des Hauptlaufs: keine Folge-ADR (Verschärfung, Festlegung im Sensor-Vertrag)

---

## Findings

### F-1 — Dateiauswahl ändert sich für Binärdateien (`-I` zu `-a`)

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: tools/harness/sdk-public-doc-check.sh:53
- `befund`: Alt prüfte `grep -I` Binärdateien nicht; neu liest `grep -a` sie als Text. Im Baum liegt eine Binärdatei (`sdks/kotlin/pgchangefeed-kotlin/gradle/wrapper/gradle-wrapper.jar`), die aktuell keinen Treffer liefert; eine künftige Binärdatei mit Musterbytes meldete einen Treffer. Die Änderung ist im Plan (§3) und im Sensor-Vertrag benannt und durch einen Tabellenfall gebunden.
- `verifizierbar`: ja — `make test-sdk-public-doc-check` (NUL-Byte-Fall), `make sdk-public-doc-check` Exit 0
- `klasse`: bewusste Verhaltensverschärfung, benannt

### F-2 — Nicht existierende Wurzel (Exit 2) ohne Tabellenfall

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: tools/harness/run-sdk-public-doc-check-tests.sh
- `befund`: Ein nicht existierendes Wurzelverzeichnis endet im Skript mit Exit 2 und Meldung (von Hand gefahren; alt: `|| true` schluckte den find-Fehler, Ausgabe „keine Kennung“, Exit 0). Der Fall steht nicht im Tabellentest; der `find`-Zweig ist über das nicht lesbare Unterverzeichnis gebunden (Mutation rot).
- `verifizierbar`: ja — Handlauf `bash tools/harness/sdk-public-doc-check.sh <fehlender Pfad>`
- `klasse`: Randfall ohne Tabellenfall

### F-3 — Ausgabeform: Dokumentation zog nach, Altform war ungenau

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: harness/sensors/sdk-public-doc-check.md:39
- `befund`: Die Form `datei:zeile:text` (mit `-H`) ist jetzt auch für eine einzelne Datei garantiert (alt: ein `xargs`-Aufruf mit einer Datei druckte ohne Dateinamen; die Altdoku nannte `datei:zeile: text`). Skriptkopf und Vertrag nennen die neue Form; kein Aufrufer parst die Ausgabe (`git grep` über den ganzen Baum: nur Make-Vorgängerkanten, Tests, Dokumentation).
- `verifizierbar`: ja
- `klasse`: Nachzug korrekt

## Geprüfte Schwerpunkte

- **Skript.** Muster, Prune-Liste (obj, bin, build, dist, .gradle, __pycache__, .pytest_cache, *.egg-info, grpc_gen), Dateiauswahl und Reihenfolge (find-Reihenfolge) gegen `8e38389a` unverändert. grep-Exit 1 ist Erfolg, ≥ 2 Exit 2 mit Dateiname, find-Fehler Exit 2, Treffer Exit 1. `find -print0` mit `read -d ''` und `grep … -- "$f"`: Dateiname mit Leerzeichen und Zeilenumbruch sowie mit führendem `-` von Hand gefahren, korrekt. Keine Pipe hinter weiterschreibendem Erzeuger; `set -euo pipefail`, `|| frc=$?`/`|| grc=$?` vermeiden den `set -e`-Abbruch; Trap `rm -rf "${scratch:?}"` auf EXIT. `make sdk-public-doc-check` am echten Baum ohne spürbare Dauer.
- **Tabellentest.** 19 Fälle, Exit und Meldungstext je Klasse gebunden; das root-Überspringen wird gemeldet (Skript und Vertrag Grenze 4). Mutationen einzeln auf einer Kopie (`git archive` im Scratchpad, ohne `sed -i`), alle rot: Lesefehler-Zweig deaktiviert; find-Zweig deaktiviert; `-a` zu `-I`; `-H` entfernt (Exit-1-Meldungsbindung); Sammelzeile geändert; `exit 1` zu `exit 0`; `-ge 2` zu `-ge 1`; Prune `grpc_gen` entfernt; Prune `dist` entfernt. Echtrepo danach unverändert.
- **Sensor-Vertrag.** Exit 2 in der Ausgangstabelle, Make-Satz, Grenze 4 (root), Verweis auf das Schwester-Gate; `harness/README.md` zählt Ausgänge nicht auf, unverändert korrekt.
- **Weitere Aufrufer.** Keine Aussage zu Exit-Codes oder Ausgabeform im Baum, die jetzt falsch ist. [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) nennt „elf Fälle“ des Tabellentests als Zeitpunktangabe (`Accepted`, unberührbar).
- **Diff-Umfang.** Kein SDK-Code, keine Version, kein ADR, kein `.github/`, kein `AGENTS.md` im Diff.
- **§3.7.** Skriptkommentare tragen Kopplung/Abgrenzung im Indikativ; `make kommentar-kennungen DIFF=8e38389a` Exit 0.
- **Hard Rules.** §3.1 (Tests schreiben nur ins Temp-Verzeichnis, Aufräumen mit `chmod -R u+rwx` vor `rm -rf "${tmp:?}"`), §3.2 (keine Suppression), §3.6 (Verschärfung, keine Senkung; Folge-ADR nicht nötig), §3.12/§3.13 (Suchlauf-Feld mit Gefundenem und Nichtgefundenem; Nachmessen Exit 0).

## Läufe (selbst gefahren)

| Lauf | Ergebnis |
|---|---|
| `make test-sdk-public-doc-check` | Exit 0, „alle 19 Fälle bestanden“ (uid 1000, Lesefehler-Fälle gelaufen) |
| `make sdk-public-doc-check` | Exit 0, „keine interne Kennung unter sdks“ |
| `make gates` (ungepiped, in Datei) | Exit 0 |
| `make test` | Exit 0 |
| `make fmt-check` | Exit 0 |
| `make kommentar-kennungen DIFF=8e38389a` | Exit 0 |
| `make suchlauf-nachmessen PLAN=<Plan>` | Exit 0 |
| `make docs-check` | Exit 0 (vor Anlage dieses Reports) |
| `make commit-traceability` | Exit 0 |
| `make sdk-pack-python` | nicht gefahren (Netz für Paketbezug nötig; Vorgänger-Kante im Make-Graph nicht geändert) |

## Negativbefunde

- geprüft, ohne Befund: tools/harness (Skript, Tabellentest)
- geprüft, ohne Befund: harness/sensors/sdk-public-doc-check.md, harness/README.md
- geprüft, ohne Befund: sdks/, docs/plan/adr, .github, AGENTS.md (kein Diff)
- geprüft, ohne Befund: Aufrufer des Gates (harness/mk, Makefile, Workflows)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 3 |

**Verdikt:** keine Fixrunde nötig, nicht merge-blockierend. Architect-Fragen: keine.
