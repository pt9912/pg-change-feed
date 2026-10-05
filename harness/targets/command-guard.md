# `make test-command-guard` — Tabellentest des PreToolUse-Guards

Ausführliche Fassung der Index-Zeilen aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

## `make test-command-guard`

Tabellentest gegen den PreToolUse-Guard `.claude/hooks/pretooluse-command-guard.sh` (`tools/harness/run-command-guard-tests.sh`): je Fall ein Kommandostring als Hook-JSON, erwartet ein Block der Klasse `pkg`/`inplace`/`interp` oder der Pass-Fall — Bestandsregeln (Paketmanager, Präfixe, `bash -c`-Rekursion, Tiefe, fail-closed), die Kopf-Erkennung hinter Wrapper-Optionen und Schlüsselwörtern, die Anführungszeichen-Lesung (`tools/harness/mask-quotes.awk`), die in-place Formen von `sed`/`perl`/`awk`/`gawk` je Position und je Mitglied der Zeichenklassen, ihre Nicht-Treffer, Host-`python3.<N>`/`perl` auf Repo-Pfaden, Host-`python`/`python3` am Kopf unbedingt (Fragment `tools/harness/blocked/python`, je Position, Pass-Fälle, benannte Falsch-Positiv-Ränder und benannte Grenzen der nicht gelisteten Namen), benannte Falsch-Positiv-Ränder und benannte Grenzen; der Guard läuft in einem Wegwerf-Repo im Temp-Verzeichnis, Prüfling per `GUARD=<Datei>`, Maskierer per `MASKER=<Datei>`, Fragmente per `BLOCKED_DIR=<Verzeichnis>` übersteuerbar (Mutationsläufe an Kopien); Host-Werkzeuge `bash`, `awk`, `mktemp`; netzlos

**Bindung:** kein Gate (ein Wächter verhindert eine Handlung, er prüft kein Ergebnis), [`MR-003`](../conventions/MR-003-guard-inplace-textwerkzeug.md), [`MR-004`](../conventions/MR-004-guard-host-python-am-kopf.md) · seit slice-harness-guard-inplace-textwerkzeug, erweitert seit slice-harness-guard-blocked-python
