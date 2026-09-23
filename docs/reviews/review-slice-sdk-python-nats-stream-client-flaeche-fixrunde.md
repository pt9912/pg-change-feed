# Review-Report: slice-sdk-python-nats-stream-client-flaeche — Fixrunde — 2026-09-23

**Review-Art:** Code (Fixrunden-Re-Review) — der Fix-Diff gegen die elf
Findings der Haupt-Review (`docs/reviews/review-slice-sdk-python-nats-stream-client-flaeche.md`,
F-1…F-11: 3 HIGH + 3 MEDIUM + 4 LOW + 1 INFO), gegen den Slice-Plan und gegen
die Hard Rules. Je Befund eine eigene Messung — die Behauptungen des Fix-Zugs
(samt des neu committeten Suchlauf-Felds) wurden nachgemessen, nicht
übernommen ([`AGENTS.md`](../../AGENTS.md) §3.12/§3.13). Kein DoD-Abgleich
als solcher — Verifier-Aufgabe (Modul 11); Zahlen behauptende Stellen sind
trotzdem nachgemessen.

**Gegenstand:** Diff-Range `0f8cc4f2..HEAD` — Commits `b4d1d352` („fix(sdk-python):
Fixrunde review F-1 bis F-11 (NATS-Flaeche)") und `264498e5` („docs(reviews):
id-unlinked-Befunde im NATS-Review-Report behoben"); 11 Dateien,
530 Insertions / 24 Deletions. Slice `slice-sdk-python-nats-stream-client-flaeche`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
**Modell:** glm-5.3-flash (Claude-Agent-SDK-Subagent) · **Datum:** 2026-09-23.

**Eingangs-Kontext:**

- Haupt-Review-Report (F-1…F-11) am HEAD:
  `docs/reviews/review-slice-sdk-python-nats-stream-client-flaeche.md`
- `Plan` `slice-sdk-python-nats-stream-client-flaeche`
  (§3 Nachzug, das neue §3.13-Suchlauf-Feld, §6 Risiken, §7 Closure-Notiz)
- [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  (Accepted) Festlegung 3 (Zeilen 203–210) und §Konsequenzen Folgepflicht 2;
  Welle-Plan `../plan/planning/welle-sdk-python-vollabdeckung.md` §6
- [`AGENTS.md`](../../AGENTS.md) §3.2 (Suppression-Verbot), §3.7
  (Kommentar-Disziplin), §3.9 (Exit-Code-Disziplin), §3.11 (Hostpfade),
  §3.12 (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- Fremd-Quelltext als Messung der Reject-Bindung: `nats-py` v2.11.0 und
  master (aktuell 2.16.0) — `nats/aio/client.py`
  (`_process_connect_init`, `_process_err`) und `nats/errors.py`
  (netzlos-nicht-prüfbar; gegen die veröffentlichten Quellen gelesen)
- Schwester-Präzedenz: C#-Fix-Zug `d86d1965` (zog `README.de.md` für seine
  Hälfte derselben Distribution-Zeile mit — Pflegehistorie der Datei)

**Eigenständig durchgeführte Prüfungen (je eigene Messung):**

- **Träger-Stände an beiden Enden gemessen** (`0f8cc4f2` UND `HEAD`,
  `git grep`/`grep`): `README.de.md:27` trägt an beiden Ständen identisch
  „die HTTP-API zusätzlich als offizielle Python-Client-Bibliothek"; der
  Fix-Zug ändert die Datei nicht (sie fehlt im Diff-Stat). Die
  `harness/README.md`-Zeile `test-sdk-python-integration` (Zeile 160) wurde
  nur im Phasen-Anteil gezogen (dritte Phase vorhanden), der
  Reject-Satz nennt unverändert nur zwei Formen, und die Endklause
  „erweitert sich Slice für Slice … (NATS-Vollinhalt)" steht unverändert.
- **F-6-Bindung gegen den nats-py-Quelltext nachgemessen** (v2.11.0 und
  master 2.16.0): der nicht-verbose Connect-Init-Pfad raised an der
  `-ERR`-Antwort `errors.Error("nats: " + err_msg.rstrip('\r\n'))` — der
  Server-Rohwortlaut des NATS-Servers (`Authorization Violation`) bleibt im
  String (`nats: 'Authorization Violation'`); Typ
  (`nats.errors.Error`, Basisklasse) und `match=` tragen beide. Die Klasse
  `nats.errors.AuthorizationError` (fester String `nats: authorization
  failed`) fliegt nur im Post-Connect-`_process_err`-Pfad — die Rezeptur der
  Haupt-Review („z. B. `nats.errors.AuthorizationError`") wäre hier die
  falsche Bindung gewesen; die Fixrunde wählte die korrekte Form.
- **§3.13-Suchlauf gegen den Baum wiederholt** (breiter als der Feldkopf):
  außer den in R-2 genannten Reststellen trägt der Baum an HEAD keine
  weitere stale Form der bewegten Eigenschaft („HTTP-API-only-Python" /
  `0.1.0`-Artefaktnamen / „no gRPC/SSE/NATS surface exists yet" /
  „NATS still follows") außerhalb historischer Records und der
  als historisch deklarierten Welle-Plan-Trigger-Zeile (Zeile 56, verifiziert).
- **Zahlen nachgemessen:** 48 Unit-Tests (`grep -c "^def test_"`:
  20 http + 3 options + 8 grpc + 10 sse + 7 nats = 48); `pyproject.toml`
  `version = "0.2.0"`; `sdks/python/dist/` trägt das reale
  `0.2.0`-Artefakt-Paar (20481/24639 Bytes, 2026-09-23 07:47) neben dem
  älteren 0.1.0-Paar (06:41); Handbuch `Version: 1.43` (Kopf Zeile 3) mit
  Änderungshistorie-Zeile 1280 (datiert 2026-09-23).
- **EOF/Suppression/Form:** `tail -c 1` je neuer Testdatei = `\n`;
  `grep -rn "noqa" sdks/` = 0 Treffer; `git diff`-Neuzüge ohne host-lokale
  Pfade.
- **Anker von F-3-Fix aufgeschlagen:** [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 3-Überschrift („Struktur … bleibt Sache des umsetzenden Zuges",
  Zeilen 203–210) und Welle-Plan §6 (Zeile 137, „diese Welle wählt v2
  desselben Packages") tragen die neue Docstring-Form beide.
- **`make gates` ungefiltert, Exit-Code direkt (kein Pipe/Wrapper): `0`**
  (Log in diesem Lauf; u. a. `a-check: gesamt: 0 Befund(e)`,
  `generated-sync: OK`).

---

## Findings

Die Nummern R-1…R-4 sind die des **Fixrunden-Reviews**; die F-Nummern der
Haupt-Review werden als F-<n> zitiert.

### R-1 — Committetes Suchlauf-Feld behauptet zwei vollzogene Nachzüge; der Baum widerlegt beide

- `kategorie`: HIGH
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.13 · §3.12 Instanz B ·
  Reviewer-Skill HIGH „Beleg trägt seinen Satz nicht"
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-nats-stream-client-flaeche.md:148-149`
- `befund`: Die Zeile 149 („Wurzel-`README.md` + `README.de.md`
  (Distribution-Zeile) … | gezogen auf dieselben vier Zustellwege") ist für
  `README.de.md` falsch — die Datei trägt an HEAD unverändert die
  „HTTP-API zusätzlich als offizielle Python-Client-Bibliothek"-Form
  (Zeile 27) und fehlt im Fix-Diff. Die Zeile 148 nennt als Befund
  „Zwei-Flächen-/Zwei-Reject-Form" und als Behandlung „gezogen auf drei
  Flächen" — die gefundene Zwei-Reject-Form ist an HEAD unbehandelt
  (`harness/README.md:160` nennt unverändert nur gRPC `Unauthenticated`
  bzw. HTTP `401`). Das committete Feld ist genau der Träger, den F-5 der
  Haupt-Review für den wahrheitsfähigen Suchlauf-Befund eingefordert hat;
  die Probe der Klasse (die genannten Fundstellen aufschlagen) widerlegt
  die Behandlungs-Angaben.
- `verifizierbar`: ja — `grep -n "offizielle Python-Client-Bibliothek"
  README.de.md` gegen Plan-Zeile 149; `grep -n "Unauthenticated" harness/README.md`
  gegen Plan-Zeile 148 (beide in diesem Lauf ausgeführt).
- `klasse`: „Beleg trägt seinen Satz nicht" (Behandlungs-Angabe widerlegt
  durch die genannte Fundstelle)

### R-2 — Träger-Nachzug F-4 nur teilweise vollzogen: deutsche Zwilling-Zeile und Reject-Dritform bleiben stehen

- `kategorie`: MEDIUM
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.13
- `pfad`: `README.de.md:27` · `harness/README.md:160`
- `befund`: Drei Reststellen desselben Nachzugs: (a) `README.de.md:27` — die
  deutsche Zwilling-Zeile blieb stehen, während `README.md:27` im selben Fix
  gezogen wurde; der C#-Schwester-Zug `d86d1965` zog dieselbe Zeile für seine
  Hälfte (Datei-Pflegehistorie), (b) `harness/README.md:160` — der „Der Lauf
  belegt:"-Satz zählt weiter nur zwei Reject-Formen; die NATS-Ablehnung
  (Server-Rohwortlaut `Authorization Violation`, gemessen am Test-Bindungswort
  aus F-6) bleibt ungenannt, (c) dieselbe Zelle endet weiter mit „erweitert
  sich Slice für Slice um die weiteren Flächen dieser Welle
  (NATS-Vollinhalt)" — seit ihre Phasenliste die dritte Fläche trägt und
  dieser Zug der letzte Flächen-Zug der Welle ist, widerspricht die
  Endklause ihrer eigenen Zelle.
- `verifizierbar`: ja — `grep` je Stelle; beide Stände gemessen
  (`README.de.md` an `0f8cc4f2` und HEAD identisch).
- `klasse`: „Arbeit überholt stehenden Träger" (Residual von F-4(b)/(c) der
  Haupt-Review)

### R-3 — Queue-Annotation behält den Sentinel-Typ nach Entzug der einzigen None-Verzweigung

- `kategorie`: LOW
- `quelle`: Maintainability (Halbsstand des INFO-F-11)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/nats_stream_client.py:75`
- `befund`: Die F-11-Fixrunde entfernte die `if payload is None: return`-Verzweigung
  (korrekt — der Sentinel hatte keinen Produzenten), ließ aber die Annotation
  `queue.Queue[bytes | None]` stehen: der Sentinel existiert jetzt nur noch
  im Typ, weder Produzent noch Konsument lesen ihn. Form-Abweichung ohne
  semantische Auswirkung.
- `verifizierbar`: nein — Lesebefund.
- `klasse`: „Unerreichbare Verzweigung" (Typhälfte, Zweig entzogen)

### R-4 — Suchlauf-Feld-Kopf deklariert den Raum schmaler, als die eigenen Zeilen dokumentieren

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-nats-stream-client-flaeche.md:140`
- `befund`: Der Feldkopf deklariert „Raum `spec/`, `docs/`, `sdks/python/`",
  während die Zeilen 5–7 des Felds Fundstellen in `harness/README.md` und
  den Wurzel-READMEs dokumentieren — der tatsächlich dokumentierte Suchlauf
  war breiter als die Kopf-Deklaration; die schmale §2-Mustersatz-Form
  (Plan-Zeile 92) wurde nicht nachgezogen. Kosmetisch; die Zeilen sind der
  Beleg für den breiteren Lauf.
- `verifizierbar`: nein — Lesebefund.
- `klasse`: „Deklarations-Raum enger als dokumentierter Lauf"

---

## Negativbefunde

- geprüft, ohne Befund: **Haupt-Review F-1** — gelöst: [`spec/pflichtenheft.md`](../../spec/pflichtenheft.md)
  Zeile 625 trägt „aktuell `0.2.0`" und die Vier-Verträge-Liste in der
  Zeilen-Form der Schwester-Zeilen (`SPEC-026`:624, `SPEC-028`:626);
  gegen [`pyproject.toml`](../../sdks/python/pgchangefeed/pyproject.toml)
  (`version = "0.2.0"`) gemessen.
- geprüft, ohne Befund: **F-2 (Suppression)** — gelöst: `grep -rn "noqa"
  sdks/` = 0 Treffer. Randnotiz: Der Fix entfernte die Kommentar-Zeile
  hinter dem Token mit (die Rezeptur der Haupt-Review sah sie „bleiben");
  die Zusage-Form lebt im Methoden-Docstring (Zeilen 70–73, „not a swallowed
  empty stream") und im Unit-Test — Rollen-Spielraum, kein Verstoß (§3.7
  verlangt keinen Kommentar).
- geprüft, ohne Befund: **F-3 (Zitierung)** — gelöst und wahr: die neue
  Docstring-Form („Welle-Plan §6: v2 desselben Packages — [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 3 delegiert … an den umsetzenden Zug", `__init__.py`
  Zeilen 14–16) trägt an beiden Belegen (Messung oben).
- geprüft, ohne Befund: **F-6 (Reject-Ursache)** — gelöst und gegen die
  nats-py-Quellen (v2.11.0, master 2.16.0) nachgemessen (Messung oben);
  der reale `make test-sdk-python-integration`-Lauf bleibt der
  Verifier-Nachweis (unverändert).
- geprüft, ohne Befund: **F-7, F-8, F-9, F-10, F-11** — gelöst: vakuum-
  Nachklammer entfernt (`rejection` bleibt genutzt, keine tote Zuweisung);
  `exceptions.py`-Docstring zählt die `SPEC-024`-NATS-Form und
  „HTTP/SSE/NATS contracts"; Plan-Zeile 138 trägt „§7
  Historie-Chronik-Zeile"; beide neue Testdateien enden auf `\n`;
  None-Verzweigung entzogen (Halbsstand: R-3).
- geprüft, ohne Befund: **F-4/F-5 zu 3,5 von 5 Stellen bzw. als Feld
  vorhanden** — die übrigen Suchlauf-Feld-Zeilen sind einzeln gegen den
  Baum verifiziert: `spec/pflichtenheft.md` §1 (Zeilen 184–192), §6 (625),
  §7 (Zeile 669), `harness/README.md` sdk-pack-Zeile (159),
  `sdks/python/README.md` Intro (Zeile 5), `options.py` (Zeile 7),
  Welle-Plan Trigger-Zeile 56 (belassen, `0.1.0` — historisch korrekt,
  verifiziert). Rest: R-1/R-2.
- geprüft, ohne Befund: **Sync-Wrapper-Exit (Plan §6 Risiko 1) am Ort
  benannt** — der Modul-Docstring (`nats_stream_client.py` Zeilen 23–28)
  trägt die Behauptung wörtlich („the connect and subscription run in a
  dedicated thread with its own event loop, the generator consumes a queue
  … the asyncio boundary stays inside this module"); der Plan-§6-Ausgang
  („Diskrepanz explizit dokumentiert statt verschwiegen") ist erfüllt.
- geprüft, ohne Befund: **Neue Zahlen des Fix-Diffs (§3.12 Instanz A)** —
  48 Tests (eigene Messung, s. o.); `0.2.0` (pyproject, spec-Zeile 625,
  harness/README.md:159 mit Herkunfts-Klammer, reales Artefakt-Paar
  gemessen); 1.43 (Handbuch-Kopf + Historie-Zeile, datiert — nicht im
  Fix-Diff, unverändert konsistent). Alle tragen Ursprung und halten der
  Messung.
- geprüft, ohne Befund: **Chronik (§3.7)** — keine neue Chronik-Form im
  Fix-Diff; „(Endstand nach der Version-Hebung `0.2.0`)" ist ein
  harness-eigener Herkunfts-Anker (disambiguiert das neben dem 0.2.0-Paar
  weiterhin liegende 0.1.0-Paar in `sdks/python/dist/`, beide gemessen),
  keine Vorher/Nachher-Erzählung.
- geprüft, ohne Befund: **Host-lokale Pfade (§3.11)** — 0 Treffer über alle
  Fix-Diff-Neuzüge (inkl. des committeten Haupt-Review-Reports).
- geprüft, ohne Befund: **`264498e5`** — rein mechanische ID-Verlinkung im
  Haupt-Review-Report (27/27 Zeilen, Form `[`SPEC-*`](../../spec/…)`), keine
  inhaltliche Änderung; Records-Form (keine §Geschichte) gewahrt;
  Betreff-Traceability ([`LH-FA-SST-009`](../../spec/lastenheft.md),
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md);
  kein `SPEC-*` im Betreff) erfüllt.
- geprüft, ohne Befund: **Runner-/Make-Wiring** — unverändert; der F-6-Fix
  berührt nur die Testdatei, der `reject_marker`-Grep des Runners bleibt
  kompatibel; keine neue Verkabelung nötig.
- geprüft, ohne Befund: **`releasing.md`** — PyPI-Ist-Stand `0.1.0`
  historisch wahr (Tag-Push Out-of-Scope, unverändert).
- geprüft, ohne Befund: **`make gates`** — ungepiped, Exit-Code direkt
  geprüft: `0` (in diesem Review-Lauf real ausgeführt).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Beleg trägt seinen Satz nicht"
(Behandlungs-Angabe widerlegt) · „Arbeit überholt stehenden Träger"
(Residual) · „Unerreichbare Verzweigung" (Typhälfte) · „Deklarations-Raum
enger als dokumentierter Lauf"

## Verdikt

**Merge-blockierend:** ja — aber eng begrenzt: 9 der 11 Findings der
Haupt-Review sind gelöst, und ihre Lösungen sind, wo sie gegen die Quelle
nachgemessen werden konnten, korrekt — inklusive der F-6-Bindung gegen den
nats-py-Quelltext, wo die Fixrunde die bessere Form wählte als die Rezeptur
des Haupt-Reviews. Offen ist der Rest von F-4/F-5: drei Träger-Stellen
(R-2) und — der schärfere Teil — das neu committete Suchlauf-Feld behauptet
für genau diese Stellen vollzogene Nachzüge (R-1). Die Restfläche: eine
deutsche README-Zeile, ein Halbsatz plus Endklause in `harness/README.md`,
zwei Plan-Tabellenzeilen, ein Typ-Token. Keine der Stellen berührt
Code-Korrektheit oder Sicherheit.

**Übergabe:** Rückgabe-Pfeil an den Implementer für eine zweite, kleine
Fixrunde — derselbe Mechanismus wie die erste; kein Rollen-Widerspruch,
keine Konflikt-Pfad-Sequenz über den Architect nötig. Erwarteter Umfang:

1. `README.de.md:27` ziehen (Muster `README.md:27`; C#-Präzedenz
   `d86d1965`).
2. `harness/README.md` Zeile 160: NATS-Ablehnungsform ergänzen (die
   Test-Bindung misst `nats: 'Authorization Violation'` am Wire) und die
   Endklause „erweitert sich Slice für Slice …" zurückbauen.
3. Plan-Zeilen 148/149 auf den realen Ist-Stand setzen — entweder die
   Nachzüge aus 1./2. vollziehen und die Zeilen dann wahr schreiben, oder
   die unbehandelten Hälften als offen melden.
4. R-3 (ein Typ-Token) und R-4 nach Ermessen der Fixrunde.

Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den
Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; er ersetzt keine
Verifikation gegen die volle DoD — der reale
`make test-sdk-python-integration`-Lauf gegen die NATS-Fläche bleibt der
Nachweis des Verifiers (Modul 11).

**DoD-Checkbox-Nachzug:** entfällt — dieses Verdikt führt zu einer
Fixrunde; die DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" wird regulär bei Schritt 21 des Implementer-Workflows
nachgezogen.