# `make sdk-public-doc-check` — keine interne Kennung in den Dateien unter `sdks/`

## Vertrag

Wird dieses Target rot, trägt eine Datei unter `sdks/` eine **interne
Kennung**: `SPEC-<NNN>`, `ADR-<NNNN>`, `ARC-<NNN>`, `LH-FA-…`/`LH-QA-…`
oder einen Slice-/Welle-Namen (`tools/harness/sdk-public-doc-check.sh`,
Muster `\b(SPEC|ADR|ARC)-[0-9]+|\bLH-(FA|QA)-[A-Z]{3}-[0-9]+|\b(slice|welle)-[a-z0-9]`).

Gegenstand ist der **Baum**: alle Textdateien unter dem Wurzelverzeichnis
(Standard `sdks`) — Quellen, Tests, README, Build-Dateien, Dockerfiles. Das
Skript liest den Baum, nicht den git-Index: das Gate läuft an jedem
Arbeitsbaum-Zustand, auch an einem mit unversionierten Dateien unter
`sdks/` (Festlegung 4 von [`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
Aufnahme in `GATE_CHECKS`).

Kommentare, Docstrings, KDoc, XML-Doku und Fehlertexte der SDKs erreichen
Anwender über die Pakete (Docstrings im Wheel, XML-Dokumentationsdatei im
`.nupkg`, Quellen im Sources-Jar und in der `sdist`), deshalb tragen sie
keine Kennung der Spezifikations-, Entscheidungs- und Anforderungsdokumente
([`LH-FA-SST-009`](../../spec/lastenheft.md)).

**Ausnahmen (Prune-Liste):** Bau-Ausgaben und erzeugter Code — `obj`, `bin`,
`build`, `dist`, `.gradle`, `__pycache__`, `.pytest_cache`, `*.egg-info`,
`grpc_gen`. Die zur Bauzeit erzeugten Python-Stubs liegen hinter dieser
Liste; ihren Inhalt gegen das installierte Paket prüft
`sdks/python/pgchangefeed/tests/test_public_text.py` im `sdk-pack-python`-Bau.

**Nicht Gegenstand:** ob die SDK-Quellen kompilieren und die Tests grün
sind (`make test-sdk-*` bzw. die `sdk-pack-*`-Bauten) und ob die Kennungen
in `spec/`, `docs/` und `harness/` korrekt getragen werden — dort sind sie
die gewollte Form (`AGENTS.md` §3.7, `make kommentar-kennungen`).

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Treffer; der Erfolgstext nennt das geprüfte Wurzelverzeichnis (`sdk-public-doc-check: keine interne Kennung unter sdks`) |
| 1 | mindestens ein Treffer; je Treffer eine Zeile `datei:zeile:text` auf stderr, darunter die Sammelzeile `sdk-public-doc-check: interne Kennung in den SDK-Dateien (siehe oben)` |
| 2 | Lesefehler (fail-closed): eine Datei oder ein Verzeichnis unter dem Wurzelverzeichnis ist nicht lesbar; stderr nennt die Datei (`sdk-public-doc-check: Lesefehler: <datei> (…)`) bzw. das Durchsuchen (`sdk-public-doc-check: Lesefehler beim Durchsuchen von <wurzel> (…)`). Ein Lesefehler endet nie mit Exit 0 |

`grep` Exit 1 (kein Treffer) ist Erfolg, `grep` Exit ≥ 2 und ein Fehler von
`find` sind Exit 2. Eine Datei mit NUL-Byte wird als Text gelesen (`grep -a`)
und liefert ihre Treffer. Die Lesefehler-Semantik ist eine Verschärfung des
Gates gegenüber [`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
keine Senkung; sie folgt dem Vertrag des Schwester-Gates
[`handbuch-public-doc-check`](handbuch-public-doc-check.md)
([`ADR-0143`](../../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)).

Über `make` kommt jeder Exit-Code ≠ 0 des Skripts als der Make-eigene Exit
`2` an — wie bei jedem gescheiterten Rezept.

## Overrides

| Variable | Wirkung |
|---|---|
| erstes Argument | abweichendes Wurzelverzeichnis (Default `sdks`) — Aufrufform für den Tabellentest `tools/harness/run-sdk-public-doc-check-tests.sh`; das Gate ruft das Skript ohne Argument |

## Grenze — was das Grün nicht abdeckt

1. **Die Wächter-Logik ist im Gate-Lauf nicht gegen ihren Tabellentest
   gesichert.** Der Gegenstand des Gates ist der Baum, nicht der Wächter;
   `make test-sdk-public-doc-check` bleibt Werkzeug
   ([`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
   Teilfrage 2, akzeptiertes Negativ). Der Ort, an dem eine Änderung am
   Wächter gegen diesen Vertrag gelesen wird, ist der Review-Diff.
2. **Eine neue Klasse erzeugter Verzeichnisse unter `sdks/` macht den Lauf
   falsch rot.** Die Prune-Liste trägt die Bau-Ausgaben-Klasse; ein
   künftig neu erzeugtes Verzeichnis außerhalb der Liste meldet das Gate
   als Treffer. Die Richtung des Fehlers ist die sichere (rot statt grün);
   die Erweiterung der Prune-Liste ist eine Umfangs-Änderung am Wächter
   (Festlegung 4 von
   [`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md))
   und zieht diese Sensor-Datei nach.
3. **Die Muster erkennen nur die gelisteten Kennungsformen.** Eine interne
   Referenz, die keine der Musterformen trägt (Freitext-Verweis auf ein
   Entscheidungsdokument ohne `ADR-<NNNN>`-Schreibweise), bleibt grün —
   das Muster- und Prüfumfangs-Verhältnis ist Gegenstand von
   `AGENTS.md` §3.7-Review, nicht des Gates.

4. **Die Lesefehler-Fälle des Tabellentests laufen nur bei uid ≠ 0.**
   `chmod 000` hindert root nicht; unter uid 0 überspringt
   `tools/harness/run-sdk-public-doc-check-tests.sh` die Fälle „nicht lesbare
   Datei“ und „nicht lesbares Unterverzeichnis“ und meldet das Überspringen.
   Das Gate selbst ist davon nicht betroffen: es liest unter root alles und hat
   dort keinen Lesefehler zu melden.

## Sperren

- **Keine über `bash`+`git` hinaus.** Reines `find`/`grep` — kein Docker,
  kein Netz, kein Toolchain-Image; das Skript steigt über
  `git rev-parse --show-toplevel` auf die Baum-Wurzel. Fehlt `git` (kein
  Repo), scheitert der Lauf an `set -euo pipefail`, bevor geprüft wird.
  Netzlos und schnell.

## Vorstufe der Paket-Bauten

Die drei Ziele `make sdk-pack-csharp`, `make sdk-pack-kotlin` und
`make sdk-pack-python` (`make sdk-pack-*`) tragen `sdk-public-doc-check` als
Voraussetzung (`harness/mk/sdk.mk`); eine Kennung unter `sdks/` bricht den
Paket-Bau ab, bevor er beginnt.

## Tabellentest

`make test-sdk-public-doc-check` fährt den Tabellentest
`tools/harness/run-sdk-public-doc-check-tests.sh` (netzlos): eine saubere Datei,
je Kennungsart ein Treffer (`SPEC-…`, `ADR-…`, `ARC-…`, `LH-…`, Slice- und
Welle-Name), die Ausnahmen `obj`, `dist` und `grpc_gen`, ein Wort mit Bindestrich
ohne Kennung, eine Datei mit NUL-Byte und die beiden Lesefehler-Fälle (Grenze 4).
Der Tabellentest bleibt Werkzeug — der Gegenstand des Gates ist der Baum, nicht
die Wächter-Logik (Grenze 1).

## Bindung

[`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
(Aufnahme in `GATE_CHECKS`; Festlegung 3 — dieser Vertrag; Festlegung 4 —
Reichweite am Arbeitsbaum) · [`LH-FA-SST-009`](../../spec/lastenheft.md)
(offizielle Client-Bibliotheken — ihre Kommentare, Docstrings, KDoc,
XML-Doku und Fehlertexte erreichen Anwender über die Pakete) ·
`tools/harness/sdk-public-doc-check.sh` ·
`harness/mk/sdk.mk` · seit slice-sdk-public-doc-check-gate
(Wächter selbst seit slice-sdk-readme-nutzerdoku).
