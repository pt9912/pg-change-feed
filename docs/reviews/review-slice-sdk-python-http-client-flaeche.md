# Review-Report: slice-sdk-python-http-client-flaeche — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-python-http-client-flaeche.md`),
`ADR-0107` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `0324b9e6..HEAD` (Abschluss von
`slice-sdk-python-projektgeruest` bis Implementer-Commit dieses Slice),
Slice `slice-sdk-python-http-client-flaeche`, Welle
`welle-sdk-python-lh-fa-sst-009`.
Vier Commits: `8f76564e` (open→next, reiner Move),
`225b5070` (Verantwortlich gesetzt, nur Slice-Datei),
`1d971e6b` (next→in-progress, reiner Move),
`314b51f6` (Inhalt: `PgChangeFeedHttpClient` + Modelle + Exceptions +
Tests + Handbuch-Update + README-Nachzug). `git log` bestätigt: genau
**ein** Commit trägt den Inhalt — kein verwaister/doppelter Commit aus dem
lokalen `git commit --amend`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. AGENTS.md §3.13 Träger-Nachzug, „Beleg trägt seinen Satz
nicht").
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-python-http-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken,
  §8 Sub-Area/Modus)
- `ADR-0107` (Accepted) — Festlegung 1/3, §Kontext „Was das ändert",
  §Konsequenzen Folgepflicht 2/4, Bindung 5 (Träger-Nachzug)
- `spec/pflichtenheft.md` §2 `SPEC-018` (Endpunkte, Token-Header-Form,
  Fehler-Antwortform, Feld-für-Feld-Tabelle), `SPEC-022` (`GET /changes`)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md` —
  Formvorbild und Vorgeschichte (beide F-1/F-2-Klassen dieses
  Geschwister-Reports wurden vom Implementer für diesen Slice proaktiv
  adressiert, siehe Prüfungen unten)

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- `docker build --no-cache -f sdks/python/Dockerfile sdks/python` real
  ausgeführt: Exit-Code direkt `0`. `pytest`-Ausgabe real gelesen:
  **„collected 23 items" / „tests/test_http_client.py .................... "
  (20 Punkte) / „tests/test_options.py ..." (3 Punkte) / „23 passed"** —
  die Aufteilung selbst nachgemessen, nicht die Commit-Message-Zahl
  übernommen (siehe F-1). Test-Image danach mit `docker rmi` entfernt.
- `make docs-check` real ausgeführt, Exit-Code direkt geprüft: `0`
  (`d-check: 827 Datei(en) geprüft, 0 Befund(e)`) — deckt sich exakt mit
  der im Commit behaupteten Zahl.
- `grep -c "^def test_" sdks/python/pgchangefeed/tests/test_http_client.py`
  → `20`, nicht die im Commit behauptete „23 pytest-Tests" für die dort
  benannte Kategorisierung (siehe F-1).
- Alle zehn öffentlichen Methoden von `PgChangeFeedHttpClient`
  (`register_consumer`, `acknowledge_consumer`, `get_consumer_position`,
  `remove_consumer`, `enable_table`, `disable_table`, `get_status`,
  `list_tables`, `run_retention`, `read_changes`) einzeln nachgezählt: genau
  zehn, keine fehlt, keine zusätzliche (kein Diagnose-/Health-Endpunkt).
- Jedes Feld jeder Dataclass in `models.py` gegen `spec/pflichtenheft.md`
  `SPEC-018`/`SPEC-022` (Zeile 362–372, 468–477) Feld für Feld
  gegengelesen — Feldnamen, Reihenfolge, Request-/Response-Form: **keine
  Abweichung** in allen zehn Fähigkeiten plus `Change`/`ReadChangesResponse`
  (zwölf Change-Felder exakt: `commit_position`, `change_id`,
  `transaction_id`, `source_table_id`, `schema`, `table`, `sequence`,
  `operation`, `old_image`, `new_image`, `schema_version`, `committed_at`).
  `ListTablesResponse.retained` trägt korrekt dieselbe `TableInfo`-Form wie
  `tables`, wie `SPEC-018`s Tabellenzeile verlangt.
- `PgChangeFeedHttpClient.__init__` gelesen: Bearer-Token kommt
  ausschließlich aus dem bei Konstruktion übergebenen `ClientOptions`
  (`self._options.api_token`), kein Modul-Level-/globaler Zustand.
- Zentraler `_handle`-Helfer gelesen: **beide** Klassen aus dem
  C#-Geschwister-Review real geprüft — Nicht-Erfolgs-Statuscodes
  (`400`/`401`/`403`/`404`/`500` + Fallback `PgChangeFeedUnexpectedStatusError`)
  **und** der malformte `2xx`-Erfolgspfad (ungültiges JSON via
  `except ValueError`, valides JSON ohne erwartetes Feld via
  `except (KeyError, TypeError)`) laufen durch denselben Helfer für alle
  zehn Methoden — kein methodenspezifischer Zweitpfad, der ihn umgehen
  könnte (jede `_get`/`_post`-Stelle ruft `self._handle` auf).
- `sdks/python/README.md`-Diff gelesen: Status-Absatz nennt jetzt korrekt
  „provides … a full HTTP API client surface"; „follow-up release" bezieht
  sich jetzt ausschließlich auf gRPC/SSE/NATS — die real zutreffende
  Abgrenzung (`ADR-0107` Festlegung 1). Die C#-HIGH-Klasse F-1
  (veraltete README-Statusbehauptung) ist hier tatsächlich **nicht**
  aufgetreten, mit einer Einschränkung (siehe F-1 unten: eine dritte,
  strukturell identische Stelle in `pyproject.toml` wurde vom
  Implementer-Suchlauf nicht gefunden).
- `grep -rn "internal/\|cmd/" sdks/python/` — kein Treffer (Exit 1).
- `docs/user/benutzerhandbuch.md`-Diff gelesen: `Version: 1.34` → `1.35`
  korrekt hochgezogen, neue Zeile `| 1.35 | 2026-09-19 | … |` in
  `### Änderungshistorie` vorhanden und inhaltlich zutreffend (nennt
  „dieselben zehn Fähigkeiten" — reale Methodenzahl stimmt exakt).
- `git show 314b51f6 -s --format=%s` geprüft: Betreff „feat(sdks/python):
  HTTP-API-Client-Flaeche (`LH-FA-SST-009`, `ADR-0107`)" — kein `SPEC-*`/`ARC-*`
  im Betreff (`AGENTS.md` §5).
- Neue `.py`-Dateien nach `slice-`/`welle-`-Nennungen im Produktionscode
  durchsucht (`AGENTS.md` §3.7 Chronik-Klasse) — kein Treffer in `src/`.
- Eigener, breiterer `grep -rn "follow-up\|folgt erst\|added by" sdks/python/`
  ausgeführt (bewusst weiter als der im Plan dokumentierte
  Implementer-Suchlauf `grep -rn "follow-up\|added by" sdks/python/README.md
  sdks/python/pgchangefeed/src/pgchangefeed/*.py`) — fand einen dritten
  Treffer außerhalb des Implementer-Suchlaufs (siehe F-1).

---

## Findings

### F-1 — Träger-Nachzug unvollständig: `pyproject.toml`-Kommentar behauptet weiterhin, die HTTP-Client-Fläche „folge erst" mit diesem Slice — der real bereits gelandet ist

- `kategorie`: **HIGH**
- `quelle`: `AGENTS.md` §3.13 (Hard Rule: „Eine Arbeit, die eine
  beschriebene Eigenschaft bewegt, zieht ihre Träger nach")
- `pfad`: `sdks/python/pgchangefeed/pyproject.toml:22-24` (unverändert in
  diesem Diff — nicht in `git diff 0324b9e6..HEAD --stat` gelistet)
- `befund`: Der Kommentar über `dependencies` sagt: „Einzige
  Fremdabhaengigkeit (`ADR-0107` Festlegung 1: httpx, keine weitere, keine
  schwere Fremdabhaengigkeit) -- die HTTP-Client-Flaeche selbst folgt erst
  mit slice-sdk-python-http-client-flaeche." Genau dieser Slice liefert
  real die HTTP-API-Client-Fläche (`PgChangeFeedHttpClient`, zehn
  Methoden, siehe oben) — die Aussage ist damit zum Zeitpunkt dieses Diffs
  veraltet. Der Implementer hat für diesen Slice denselben Fehlerklasse
  proaktiv gesucht (Plan-Nachzug §3: `grep -rn "follow-up\|added by"
  sdks/python/README.md sdks/python/pgchangefeed/src/pgchangefeed/*.py`)
  und dabei zwei Treffer gefunden und korrekt nachgezogen
  (`README.md`, `__init__.py`) — der Suchlauf deckte aber weder
  `pyproject.toml` (Glob `src/pgchangefeed/*.py` schließt es aus) noch die
  in `pyproject.toml` tatsächlich verwendete deutsche Formulierung „folgt
  erst" (Suchmuster nur „follow-up"/„added by") ab. Dieselbe
  Fehlerklasse, die für zwei Stellen proaktiv vermieden wurde, ist an
  einer dritten, strukturell identischen Stelle real aufgetreten.
- `verifizierbar`: ja — Datei lesen (`pyproject.toml:22-24`) gegen die
  real vorhandene Methodenzahl in `http_client.py` (10 Methoden) halten.
- `klasse`: Arbeit überholt stehenden Träger
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`-Familie — dieselbe
  Fehlerklasse wie F-1 des C#-Geschwister-Reviews, hier an einer dritten,
  vom eigenen proaktiven Suchlauf nicht erfassten Stelle)

### F-2 — Commit-Message zählt „23 pytest-Tests" für eine Kategorisierung, die nur 20 der 23 tatsächlich vorhandenen Tests trägt

- `kategorie`: MEDIUM
- `quelle`: Maintainability / verwandt „Beleg trägt seinen Satz nicht"
  (Reviewer-Skill HIGH-Klasse, hier abgeschwächt: der Fehlbeleg steht
  ausschließlich in der Commit-Message, nicht in einem persistenten
  Doku-Träger — `grep -rn "23" docs/plan/planning/in-progress/
  slice-sdk-python-http-client-flaeche.md docs/user/benutzerhandbuch.md
  sdks/python/README.md` liefert keinen Treffer)
- `pfad`: Commit-Message `314b51f6` (kein Datei:Zeile-Anker — Träger ist
  die Commit-Message selbst)
- `befund`: Die Commit-Message behauptet „23 pytest-Tests (Happy Path je
  Faehigkeit, Auth-Boundary 401/403, Statuscode-Mapping 400/404/500/
  unerwartet, malformter 2xx-Body …)". Real trägt
  `test_http_client.py` genau **20** Tests in dieser Kategorisierung
  (nachgezählt, `grep -c "^def test_"`); die „23" ergibt sich nur, wenn
  die drei bereits vor diesem Slice existierenden, unter keine der
  genannten Kategorien fallenden `ClientOptions`-Validierungstests aus
  `test_options.py` (Vorgänger-Slice `slice-sdk-python-projektgeruest`)
  mitgezählt werden. Der Docker-Testlauf zeigt real „23 passed" — die
  Gesamtzahl der Testsuite stimmt, aber die im selben Satz genannte
  Kategorisierung trägt nur 20 davon.
- `verifizierbar`: ja — `docker build --no-cache -f sdks/python/Dockerfile
  sdks/python` zeigt die Aufschlüsselung nach Datei (`test_http_client.py`
  20 Punkte, `test_options.py` 3 Punkte).
- `klasse`: Beleg trägt seinen Satz nicht (Zahlen-Klasse — Gesamtzahl
  korrekt, Kategorisierungs-Zuordnung nicht)

## Negativbefunde

- geprüft, ohne Befund: Methodenzahl gegen `SPEC-018`/`SPEC-022` — exakt
  zehn öffentliche Methoden, keine fehlt, kein Diagnose-/Health-Endpunkt
  hinzugefügt (Plan §1 „Ausdrücklich NICHT in diesem Slice" korrekt
  eingehalten).
- geprüft, ohne Befund: JSON-Feldnamen aller Dataclasses in `models.py`
  gegen `SPEC-018`/`SPEC-022` Feld für Feld — keine Abweichung, inklusive
  der zwölf `Change`-Felder und der `ListTablesResponse.retained`-Form.
- geprüft, ohne Befund: Bearer-Token-Übergabe und Abwesenheit globalen/
  Modul-Level-Zustands (`PgChangeFeedHttpClient.__init__`).
- geprüft, ohne Befund: `PgChangeFeedError`-Hierarchie konsistent über
  alle zehn Methoden, inklusive des Erfolgspfads mit malformtem
  `2xx`-Body — beide Klassen des C#-Geschwister-Findings F-2 sind hier
  bereits über den zentralen `_handle`-Helfer abgedeckt und mit zwei
  eigenen Tests belegt (`test_non_json_success_body_raises_typed_
  malformed_response_error`, `test_success_body_missing_expected_field_
  raises_typed_malformed_response_error`).
- geprüft, ohne Befund: Import-Grenze (`grep -rn "internal/\|cmd/"
  sdks/python/` — kein Treffer, weder als Import noch als Kommentar/
  Docstring-Erwähnung).
- geprüft, ohne Befund: Testabdeckung je Fähigkeit real gezählt (nicht die
  Commit-Message-Gesamtzahl unbesehen übernommen, siehe F-2) — jede der
  zehn Fähigkeiten trägt mindestens einen Happy-Path-Test, `401`/`403`
  sind als Auth-Boundary abgedeckt, `400`/`404`/`500`/undokumentierter
  Status sind einzeln getestet; Docker-Testlauf bestätigt real 23/23 grün
  (20 davon `test_http_client.py`, 3 vorbestehende `test_options.py`).
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version
  1.34→1.35 korrekt, Änderungshistorie-Zeile vorhanden und inhaltlich
  zutreffend (Methodenzahl „zehn" stimmt).
- geprüft, ohne Befund: `sdks/python/README.md` §Status — die
  C#-HIGH-Klasse F-1 (veraltete „follow-up"-Statusbehauptung für die real
  gelieferte Fläche) ist hier für README **und** `__init__.py` korrekt
  proaktiv vermieden; siehe F-1 oben für die dritte, nicht gefundene
  Stelle in `pyproject.toml`.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) im
  Produktionscode — keine Konjunktiv-Begründung über eine verworfene
  Alternative, kein abwesender Text, keine Slice-/Wellen-Chronik im
  Produktionscode-Kommentar (`grep` nach `slice-`/`welle-` in `src/`: kein
  Treffer).
- geprüft, ohne Befund: Traceability — Commit-Betreff `314b51f6` nennt
  `LH-FA-SST-009` und `ADR-0107`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: genau ein Commit trägt den Inhalt dieses Slices
  (`git log 0324b9e6..HEAD`) — der lokale `git commit --amend` zur
  Korrektur der beiden vom Implementer selbst gefundenen Fehler
  (Docstring-Erwähnung des Go-Server-Pfads, Commit-Betreff mit
  `SPEC-018`/`SPEC-022`) hat keinen verwaisten/doppelten Commit hinterlassen.
- geprüft, ohne Befund: reale Docker-Build- und Testausführung — Exit-Code
  direkt (ungepiped) geprüft, `0`; 23/23 Tests grün. Test-Image nach
  Prüfung entfernt.
- geprüft, ohne Befund: `make docs-check` real gefahren, Exit-Code direkt
  geprüft, `0` (827 Dateien, 0 Befunde — deckt sich exakt mit der
  Commit-Message-Zahl).
- geprüft, ohne Befund: host-lokale absolute Pfade im Diff — `git diff
  0324b9e6..HEAD` gegen die Wurzel-Segmente eines Entwicklerrechners
  (AGENTS.md §3.11 Präfixliste) durchsucht, kein Treffer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Arbeit überholt stehenden Träger
(`pyproject.toml`-Kommentar, `AGENTS.md` §3.13) · Beleg trägt seinen Satz
nicht (Commit-Message-Testzahl-Kategorisierung)

## Verdikt

**Merge-blockierend:** ja — F-1 ist ein Hard-Rule-Verstoß (`AGENTS.md`
§3.13) an einer Stelle, die real Teil des über `pip install pgchangefeed`
bezogenen Sdist sein wird. Eine Fixrunde am Implementer ist nötig:

1. `pyproject.toml:22-24` auf den tatsächlichen Lieferstand dieses Slice
   heben (HTTP-Client-Fläche jetzt vorhanden), analog der bereits
   korrekten Formulierung in `README.md`/`__init__.py`.
2. F-2 (MEDIUM) kann in derselben Fixrunde mitgezogen werden (Commit ist
   bereits gelandet und unveränderlich in der Historie — kein
   `--amend`, da dieser Report bereits existiert; die Korrektur betrifft
   nur zukünftige Träger, keine rückwirkende Text-Reparatur der
   Commit-Message selbst), ist aber für sich allein keine Merge-Blockade.

**Übergabe:** Rückmeldung an den Implementer-Agenten mit diesem Report.
Da eine Fixrunde nötig ist, bleibt die DoD-Checkbox „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" im Slice-Plan **offen** — sie wird
regulär bei Schritt 21 des Implementer-Workflows nach der Fixrunde
nachgezogen (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift hier
nicht, weil eine Fixrunde stattfindet). Dieser Report ersetzt keine
Verifikation gegen die DoD — das bleibt Verifier-Aufgabe (Modul 11).
