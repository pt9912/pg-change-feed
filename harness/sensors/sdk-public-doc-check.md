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
| 1 | mindestens ein Treffer; je Treffer eine Zeile `datei:zeile: text` auf stderr, darunter die Sammelzeile `sdk-public-doc-check: interne Kennung in den SDK-Dateien (siehe oben)` |

Über `make` kommt der Exit-Code des Skripts (1) als der Make-eigene Exit
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

## Sperren

- **Keine über `bash`+`git` hinaus.** Reines `find`/`grep` — kein Docker,
  kein Netz, kein Toolchain-Image; das Skript steigt über
  `git rev-parse --show-toplevel` auf die Baum-Wurzel. Fehlt `git` (kein
  Repo), scheitert der Lauf an `set -euo pipefail`, bevor geprüft wird.

## Bindung

[`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
(Aufnahme in `GATE_CHECKS`; Festlegung 3 — dieser Vertrag; Festlegung 4 —
Reichweite am Arbeitsbaum) · [`LH-FA-SST-009`](../../spec/lastenheft.md)
(offizielle Client-Bibliotheken — ihre Kommentare, Docstrings, KDoc,
XML-Doku und Fehlertexte erreichen Anwender über die Pakete) ·
`tools/harness/sdk-public-doc-check.sh` ·
`harness/mk/sdk.mk` · seit slice-sdk-public-doc-check-gate
(Wächter selbst seit slice-sdk-readme-nutzerdoku).
