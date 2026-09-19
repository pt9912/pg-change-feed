# Verifikationsbericht: slice-sdk-python-projektgeruest — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-python-projektgeruest.md` §2)
und die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff
als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-python-projektgeruest.md`](review-slice-sdk-python-projektgeruest.md)
samt Fixrunden-Nachprüfung) und **nicht** gegen realen Bedarf (Validator,
hier nicht ausgelöst).

**Gegenstand:** Diff-Range `f3e10429..HEAD`, Slice
`slice-sdk-python-projektgeruest`, Welle
`welle-sdk-python-lh-fa-sst-009`. `HEAD` =
`2ed7518daab329d86475688738c0223f65dc4f27`, Arbeitsbaum sauber
(`git status --short` leer, vor und nach dieser Verifikation).

**Frischer Kontext:** Slice-Plan (§1–§8 vollständig), Review-Report
inkl. Fixrunden-Nachprüfung (vollständig), `ADR-0107` (vollständig, alle
fünf Festlegungen und §6 „Was diese ADR nicht ändert"). Jede DoD-Zeile
einzeln gegen den realen Datei-/Lauf-Zustand geprüft, nichts aus Plan,
Commit-Message oder Review-Report ungeprüft übernommen: eigene Lektüre
von `pyproject.toml`/`Dockerfile`/`README.md`/`__init__.py`/`options.py`/
`test_options.py`/`.gitignore` im Volltext, eigener `grep` gegen private
Imports, eigener `docker build --no-cache`-Lauf (**vierte** unabhängige
Bau-Bestätigung nach Implementer, Reviewer-Erstlauf und
Reviewer-Fixrunden-Nachprüfung), eigener ungepipter `make gates`-Lauf mit
direkter Exit-Code-Prüfung (`AGENTS.md` §3.9), eigene `git diff`-Läufe
gegen die von `ADR-0107` §6 und der DoD als „unberührt"/„entfällt"
behaupteten Träger.

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 `pyproject.toml`-Metadaten und Import-/Abhängigkeitsgrenze

Eigene Lektüre `sdks/python/pgchangefeed/pyproject.toml`: `name =
"pgchangefeed"` ✓, `version = "0.1.0"` (PEP 440, `ADR-0107` Festlegung 4)
✓, `description` (nennt korrekt HTTP-API) ✓, `authors = [{ name =
"pt9912" }]` ✓, `license = "MIT"` — gegen `LICENSE` im Repo-Root
gegengelesen, tatsächlich MIT, Copyright pt9912 2026 ✓,
`[project.urls]` (`Homepage`/`Repository`, beide auf
`https://github.com/pt9912/pg-change-feed`) ✓, `readme = "README.md"`
✓ (lokal, mit dem im Plan-Nachzug §3 dokumentierten
`DistutilsOptionError`-Grund für die Abweichung von einem relativen
Ebenen-Wechsel).

`dependencies = ["httpx>=0.27"]` ist die einzige Laufzeit-
Fremdabhängigkeit; `[project.optional-dependencies] test =
["pytest>=8"]` ist eine Test-Extra, keine Laufzeitabhängigkeit. Kein
weiterer `dependencies`-Block im Volltext.

Import-Grenze eigenständig geprüft (nicht nur den Review-Befund
zitiert):

```
$ grep -rn "internal/\|cmd/pg-change-feed" sdks/python/
(kein Treffer, Exit 1)
```

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 `sdks/python/Dockerfile` — vierte unabhängige Bau-Bestätigung

Eigene Lektüre: digest-gepinnte `python:3.13-slim`-Basis
(`sha256:8d9d0b8bcf6506481eae4907c18f5e3e7902e629f5f6d684f9e7c32e85e3ddf0`,
Manifest-Index-Digest), `AS build` ohne Folgestufe, kein Runtime-Server-
Start, kein `python -m build`-Aufruf — deckungsgleich mit der
DoD-Formulierung ("kein Runtime-Server-Start nötig").

Eigener realer Bau, unabhängig von Implementer- und beiden
Reviewer-Läufen:

```
$ docker build --no-cache -f sdks/python/Dockerfile -t \
    pgcf-python-verifier-check sdks/python
...
#13 [9/9] RUN pytest
tests/test_options.py ...                                    [100%]
3 passed in 0.01s
EXIT_CODE=0
$ docker rmi pgcf-python-verifier-check
```

Exit-Code `0`, 3/3 Tests grün — identisch zur Implementer-Behauptung und
zu beiden Reviewer-Läufen (Erstlauf und Fixrunden-Nachprüfung). Keine
Divergenz über die vier unabhängigen Bau-Läufe. **Ergebnis: Checkbox
berechtigt auf `[x]`.**

### 1.3 `sdks/python/README.md`

Eigene Lektüre: durchgängig Englisch, nennt Zweck, Status (`0.x.y`,
`ADR-0107`), Installationsweg (`pip install pgchangefeed`), verweist für
den vollen Kontext auf das Repo-Root-`README.md` und
`spec/pflichtenheft.md` (`SPEC-018`) statt Draht-Details selbst zu
beschreiben — kein Duplikat der kanonischen Draht-Doku. Alle vier Links
laufen über absolute GitHub-Blob-URLs
(`https://github.com/pt9912/pg-change-feed/blob/main/...` bzw.
`/tree/main/...`), kein relativer Pfad:

```
$ grep -n '](\.\.' sdks/python/README.md; grep -n '](/[^/]' sdks/python/README.md
(kein Treffer)
```

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 „`make gates` grün" — siehe eigenständiger Abschnitt 2 unten

### 1.5 „Review durchgeführt, Report liegt vor"

Eigene Lektüre `review-slice-sdk-python-projektgeruest.md` (Erstlauf +
Fixrunden-Nachprüfung): Erstlauf fand 1 HIGH (F-1, Slice-Chronik im
`__init__.py`-Docstring); Fixrunden-Nachprüfung (Commit `70aa64ce`)
bestätigt F-1 als **behoben** — eigene Lektüre des `git diff
7281a633..HEAD`, beide gerügten Slice-Namen sind entfernt, verbleibende
Kennung ausschließlich `ADR-0107 Festlegung 1/3/4`. Gesamt-Verdikt
„Freigegeben", 0 offenes HIGH, 0 MEDIUM. Deckt sich mit dem
DoD-Zeilentext im Slice-Plan. **Ergebnis: Checkbox berechtigt auf
`[x]`.**

### 1.6 „Doku-Update für `harness/README.md` entfällt"

Eigener, vom Auftrag vorgegebener Diff-Lauf:

```
$ git diff f3e10429..HEAD --stat -- harness/README.md Makefile harness/mk/
(leer)
```

Kein neues `make`-Target entstand in diesem Slice — die Behauptung ist
korrekt. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.7 „Reconciliation entfällt"

```
$ find docs/plan/planning -iname "reconciliation*"
(kein Treffer)
```

Es existiert tatsächlich keine `reconciliation.md` in diesem Repo — die
Begründung „kein Brownfield-Bootstrap" trifft strukturell zu.
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.8 Verbleibende `[ ]`-Checkboxen (Closure-Notiz, Beobachtungs-Register, §6-Risiken, drei Paarungen)

Eigene Lektüre §2/§7 des Plans: Closure-Notiz trägt ausschließlich den
Platzhalter `<…>`, Beobachtungs-Register-Checkbox ist `[ ]`, jedes der
drei §6-Risiken trägt **keinen** erledigten Ausgang — alle drei sind im
Plan-Text selbst bereits explizit als „weiter offen" bezeichnet (nicht
etwa fälschlich als „eingetreten"/„entfallen" markiert), die
Drei-Paarungen-Checkbox ist `[ ]` mit dem korrekten Verweis „Prüfung
läuft regelkonform bei" der Welle-Closure. Der Slice-Plan liegt
weiterhin unter `docs/plan/planning/in-progress/`. Keine dieser vier
Zeilen wird im Plan fälschlich als erledigt behauptet — der `[ ]`-Zustand
ist korrekt und konsistent mit dem `in-progress`-Stand. **Ergebnis: kein
Befund — der offene Zustand ist zutreffend, nicht fälschlich verschwiegen
oder vorzeitig abgehakt.**

## 2. `make gates` — eigenständig, ungepiped ausgeführt

```
$ git log -1 --format=%H
2ed7518daab329d86475688738c0223f65dc4f27
$ git status --short
(leer)
$ make gates > /tmp/claude-501-verify-gates.log 2>&1; ec=$?; echo "EXIT_CODE=$ec"
EXIT_CODE=0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK —       54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 822 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Alle sechs Gates real gefahren und grün, Exit-Code `0` direkt (ungepiped)
geprüft, nicht durch eine Pipe/einen Wrapper hindurch. `docs-check` zeigt
`822` geprüfte Dateien — identisch mit dem Wert der
Fixrunden-Nachprüfung des Reviewers zuzüglich dieser einen neuen
Verifikationsbericht-Datei, die als eigene Datei erst nach dem
Review-Commit entsteht (keine Diskrepanz, plausibel erklärt). **Ergebnis:
Checkbox 1.4 berechtigt auf `[x]`.**

## 3. Spec-/ADR-Konformität zusätzlich zur DoD

### 3.1 `ADR-0107` Festlegung 1/3 — kein Vorgriff, kein privater Import

Eigene vollständige Lektüre `options.py`: Die Klasse trägt ausschließlich
zwei Felder (`address: str`, `api_token: str`) und eine validierende
`__post_init__` (leere Werte werfen `ValueError`). Kein HTTP-Aufruf, kein
`httpx.Client`-Feld, keine Methode, die eine Endpunkt-Operation ausführt.
Der Modul-Docstring benennt selbst explizit, dass keine Vorwegnahme von
Endpunkt-Methoden stattfindet, und zieht die Analogie zum bereits
geprüften C#-Formvorbild. `test_options.py` bindet real an die
`__post_init__`-Validierung (zwei Negativtests: leere Adresse, leerer
Token), nicht nur an den Rückgabewert. **Ergebnis: Festlegung 1
eingehalten, kein Vorgriff.**

### 3.2 `__init__.py` — eigenständige Bestätigung, kein Zitat des Review-Berichts

Eigene vollständige Lektüre `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py`:

```
"""PG Change Feed Python client library.

This release exposes only the shared connection configuration
(`ClientOptions`) that every wire surface needs regardless of transport --
address and bearer token (ADR-0107 Festlegung 1/3/4). The HTTP API client
surface itself (the SPEC-018 capabilities) is added by a follow-up release.
"""
```

Kein Slice-Name (`slice-sdk-python-projektgeruest`,
`slice-sdk-python-http-client-flaeche`) steht mehr in dieser Datei — die
einzige verbleibende Kennung ist `ADR-0107 Festlegung 1/3/4`. Die
verbleibende Vorher/Nachher-Sprache („is added by a follow-up release")
nennt keinen konkreten künftigen Slice mehr, formgleich zum zitierten
Formvorbild `PgChangeFeedClientOptions.cs`. **Ergebnis: F-1 real behoben,
unabhängig vom Review-Bericht selbst bestätigt.**

### 3.3 `ADR-0107` §6 — was diese ADR nicht ändert

```
$ git diff f3e10429..HEAD --stat -- .a-check.yml harness/README.md \
    spec/architecture.md docs/user/version.md sdks/csharp/
(leer)
```

Alle fünf explizit geprüften, von `ADR-0107` §6 als „bleibt unberührt"
benannten Träger sind in diesem Diff tatsächlich unangetastet. Zusätzlich
geprüft und ebenfalls leer: `spec/` (gesamt), `AGENTS.md`,
`harness/conventions.md`, `docs/plan/adr/`, `.github/` (bereits im
Review-Bericht belegt, hier stichprobenhaft für `spec/architecture.md`
und `docs/user/version.md` gegengemessen). **Ergebnis: §6 eingehalten.**

## 4. §6-Risiken — Zwischenstand

Alle drei §6-Risiken (Konfigurations-Skelett ohne echten gemeinsamen
Nenner; fehlender Python-Referenz-Client; Python-Mindestversion
`>=3.11`) tragen im Plan selbst explizit den Text „weiter offen" — kein
Risiko ist fälschlich als erledigt/entfallen markiert. Das ist konsistent
mit dem `in-progress`-Zustand und keine DoD-Verletzung: §2 listet den
Risiko-Ausgang ausdrücklich als eigene, noch offene Checkbox (Abschnitt
1.8 oben). Kein Befund.

## 5. Bereinigung

Test-Image `pgcf-python-verifier-check` (`docker rmi`) wurde nach der
Prüfung entfernt; `docker images | grep pgcf` danach leer. `git status
--short` vor und nach allen eigenen Läufen leer.

---

## Verdikt

**DoD erfüllt** — für den aktuellen `in-progress`-Stand des Slice. Alle
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht
aus Plan, Commit-Message oder Review-Report übernommen:

1. Jede `[x]`-Checkbox in §2 ist gegen reale Dateien/Läufe geprüft und
   berechtigt gesetzt (`pyproject.toml`-Metadaten Feld für Feld,
   Import-/Abhängigkeitsgrenze per `grep`, Dockerfile real neu gebaut —
   vierte unabhängige Bau-Bestätigung mit identischem Ergebnis zu
   Implementer und beiden Reviewer-Läufen, README.md, `make gates`,
   Review-Report-Inhalt inkl. Fixrunden-Nachprüfung, die beiden
   „entfällt"-Begründungen real per `git diff`/`find` nachgeprüft).
2. Die vier verbleibenden `[ ]`-Checkboxen (Closure-Notiz,
   Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen) stehen
   korrekt offen — keine davon wird im Plan fälschlich als erledigt
   dargestellt; alle drei §6-Risiken sind ausdrücklich als „weiter offen"
   benannt.
3. `make gates` lief eigenständig, ungepiped, mit Exit-Code `0`; alle
   sechs Gates einzeln im Log bestätigt.
4. `ADR-0107` Festlegung 1 (kein Vorgriff auf Endpunkt-Methoden),
   Festlegung 3 (keine private Import-Grenze verletzt) und §6 (fünf
   benannte Träger unberührt) sind real eingehalten.
5. Das einzige HIGH-Finding des Erstreview-Laufs (F-1, Slice-Chronik in
   `__init__.py`) ist real behoben — eigenständig durch vollständige
   Lektüre der Datei bestätigt, nicht nur den Review-Bericht zitiert.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`)
bleiben die im Plan selbst bereits als offen geführten Closure-Pflichten
zu erfüllen (§7 Closure-Notiz, Beobachtungs-Register-Entscheid,
Risiko-Ausgänge in §6, die drei Paarungen bei der
`welle-sdk-python-lh-fa-sst-009`-Closure).
