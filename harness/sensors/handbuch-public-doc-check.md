# `make handbuch-public-doc-check` — keine interne Kennung in den Nutzerdokumenten unter `docs/user/`

## Vertrag

Wird dieses Target rot, trägt eines der geprüften Nutzerdokumente eine
**interne Kennung** oder einen **Link nach `docs/plan/` bzw. `docs/reviews/`**
(`tools/harness/handbuch-public-doc-check.sh`). Die Nutzerdokumentation ist die
Betreiber- und Integrator-Sicht ([`LH-QA-OPS-001`](../../spec/lastenheft.md));
Kennungen der Spezifikations-, Entscheidungs- und Anforderungsdokumente sowie
Slice-, Welle- und Register-Namen gehören in Plan und Commit. Das Muster `P`
unten fängt `LH-FA-…`, `LH-QA-…` (und `LH-RB-…`), `ADR-…`, `SPEC-…`, `ARC-…`,
`BEO-…`, `MR-…`, `CO-…` und Slice-/Welle-Namen, das Muster `L` die Links.

**Muster** (ERE, zeilenweise `grep -InE`, kein `grep -P`):

```text
P = \b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]
L = docs/(reviews|plan)/|\.\./(reviews|plan)/
```

`MR-`/`CO-` verlangen drei Ziffern; `slice-`/`welle-` verlangen am linken Rand
Zeilenanfang oder ein Zeichen außerhalb von Buchstabe, Ziffer, `_` und `-`
(ein zusammengesetztes Wort wie `byte-slice-x` trifft nicht).

**Reichweite — zwei benannte Listen plus Vollständigkeitsprüfung**
([`ADR-0143`](../../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
Festlegung 1):

| Liste | Dateien unter `docs/user/` | Grund |
|---|---|---|
| geprüft | `benutzerhandbuch.md`, `benutzerhandbuch-standard.md`, `version.md` | Nutzerdokumentation |
| ausgenommen | `bench-abdeckung.md`, `ci-matrix-abdeckung.md`, `e2e-abdeckung.md`, `sdk-e2e-abdeckung.md` | von Runnern geschriebene Erzeugnisse; die Kennung je Zeile ist ihr Inhalt und die Eingabe von `make doc-trace` |

Jede `*.md` unter `docs/user/`, die in keiner Liste steht, und jede genannte,
nicht vorhandene Datei endet mit Exit 2: eine neue Datei erzwingt die
Klassifikation (geprüft oder ausgenommen). Nicht-`.md`-Dateien sind nicht
Gegenstand. **Keine Ausnahme** für Befehlsbeispiele oder Fenced-Blöcke — das
Skript liest Zeilen.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Treffer; der Erfolgstext nennt die Zahl der geprüften Dateien (`handbuch-public-doc-check: keine interne Kennung in 3 Nutzerdokumenten unter docs/user/`) |
| 1 | mindestens ein Treffer; je Treffer eine Zeile `docs/user/<datei>:<zeile>:<text>` auf stderr, darunter die Sammelzeile `handbuch-public-doc-check: interne Kennung oder Link nach docs/plan/ bzw. docs/reviews/ in den Nutzerdokumenten (siehe oben)` |
| 2 | Klassifikationsfehler: unklassifizierte `*.md` (`unklassifizierte Datei: docs/user/<datei>`) oder genannte, fehlende Datei (`genannte Datei fehlt: docs/user/<datei>`) auf stderr; es wird nicht gegen die Muster geprüft. Ebenso **Lesefehler** (fail-closed, nie „sauber“): eine nicht lesbare geprüfte Datei (`Lesefehler: docs/user/<datei> (<Meldung von grep>)`, `grep`-Exit ≥ 2) oder ein nicht lesbares Verzeichnis unter `docs/user/` (`Lesefehler beim Durchsuchen von docs/user/ (<Meldung von find>)`). `grep`-Exit 1 (kein Treffer) bleibt Erfolg; eine Datei mit NUL-Byte wird als Text gelesen (`grep -a`), ihre Treffer bleiben sichtbar |

Über `make` kommt jeder Exit ≠ 0 des Skripts als Make-eigener Exit `2` an.

## Overrides

| Variable | Wirkung |
|---|---|
| erstes Argument | abweichende Wurzel (Default: die Repo-Wurzel; das Skript liest `<Wurzel>/docs/user/`) — Aufrufform für den Tabellentest `tools/harness/run-handbuch-public-doc-check-tests.sh`; das Gate ruft das Skript ohne Argument |

## Grenze — was das Grün nicht abdeckt

1. **Das Gate liest Kennungen, nicht Sinn.** Chronik-Sprache ohne Kennung
   („früher … jetzt …“), Entwickler- statt Betreibersicht und eine
   Änderungshistorie-Zeile ohne Betreibernutzen bleiben grün; diese Hälfte
   tragen der Skill `.harness/skills/nutzerdoku-schreiben.md` und der Reviewer.
2. **Groß-/Kleinschreibung zählt.** Eine groß geschriebene Namensform
   (`Slice-1`) und Kennungsformen außerhalb der Muster bleiben grün.
3. **Die vier Erzeugnisse sind ausgenommen, die Maintainer-Doku liegt außerhalb.**
   Eine Kennung dort fängt nur der Reviewer; die Maintainer-Doku unter
   `docs/maintainer/` (Release-Prozess) liegt außerhalb von `docs/user/` und
   ist nie Gegenstand des Gates. Verschiebt jemand eine Datei von „geprüft“ nach
   „ausgenommen“, ist das eine Milderung und braucht eine eigene ADR
   ([`ADR-0143`](../../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
   Festlegung 9); eine **neue** Datei in einer Liste zu führen, ist der
   erwartete Weg.
4. **Die Wächter-Logik ist im Gate-Lauf nicht gegen ihren Tabellentest
   gesichert.** `make test-handbuch-public-doc-check` bleibt Werkzeug; der Ort,
   an dem eine Änderung am Wächter gegen diesen Vertrag gelesen wird, ist der
   Review-Diff. Die Lesefehler-Fälle des Tabellentests (`chmod 000`) entfallen
   unter root, der Rechte nicht bindet; der Lauf druckt das dann.
5. **Muster-Ränder.** `ADR 0143` (Leerzeichen statt Bindestrich) trifft nicht;
   `--slice-x` (Bindestriche davor) trifft nicht; eine GitHub-URL mit
   `docs/plan/` trifft als Link-Muster (Falsch-Positiv, gewollt: auch eine
   absolute URL auf Plan-Dokumente gehört nicht ins Handbuch).

## Sperren

- **Host-Werkzeuge:** `bash`, `git`, `grep`, `find`, `sort`, `sed`, `head`,
  `mktemp` (`AGENTS.md` §3.1, POSIX-/coreutils-Basis). Netzlos und schnell
  (reines `grep`), kein Docker;
  das Skript steigt über `git rev-parse --show-toplevel` auf die Wurzel und
  legt ein Temp-Verzeichnis an, das es beim Ende entfernt. Der Wächter liest
  den Arbeitsbaum, nicht den Index.

## Tabellentest

`make test-handbuch-public-doc-check` fährt den Tabellentest gegen
`tools/harness/handbuch-public-doc-check.sh`
(`tools/harness/run-handbuch-public-doc-check-tests.sh`, netzlos): je
Kennungsart und je Link-Klasse ein Treffer (auch im Fenced-Block, in
`version.md` und in einer Datei mit NUL-Byte), Wortrand-Fälle ohne Treffer,
saubere Wurzel, je ein Fall für die vier ausgenommenen Dateien der Liste oben,
eine Nicht-`.md`-Datei ohne Befund; mit Exit 2 eine unklassifizierte Datei (auch
im Unterverzeichnis), eine fehlende geprüfte und eine fehlende ausgenommene Datei
sowie, außer unter root, eine nicht lesbare Datei und ein nicht lesbares
Unterverzeichnis, je mit Meldungstext. Die Zahl vier ist gemessen (2026-10-05):
die Liste `excluded=` in `tools/harness/handbuch-public-doc-check.sh` trägt vier
Namen, der Tabellentest je einen Fall „… ausgenommen“ für jede
(`grep -c '^expect 0 ".*ausgenommen"' tools/harness/run-handbuch-public-doc-check-tests.sh`
druckt `4`). Die Gesamtzahl der Fälle druckt der Lauf selbst.
Der Tabellentest bleibt Werkzeug — der Gegenstand des Gates ist der Bestand der
Nutzerdokumente, nicht die Wächter-Logik
([`ADR-0143`](../../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
Festlegung 5; Grenze 4).

## Bindung

[`ADR-0143`](../../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(Aufnahme in `GATE_CHECKS`, Muster, Reichweite) ·
[`LH-QA-OPS-001`](../../spec/lastenheft.md) (Dokumentation für Betreiber) ·
`tools/harness/handbuch-public-doc-check.sh` · `harness/mk/doc-gate.mk` · seit
slice-handbuch-public-doc-check-gate-und-skill.
