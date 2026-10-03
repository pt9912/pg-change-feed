# Review-Report: slice-pin-stale-alle-digest-pins — 2026-10-03

**Review-Art:** Code — geprüft gegen Plan, [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Festlegung 1 bis 6) und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `pin-stale-alle-digest-pins`, Diff-Range `3eb60de6..HEAD` (`HEAD` = `90bb3f88`, ein
Implementer-Commit, 10 Dateien, +493/−28): `tools/harness/pin-stale-all.sh`, `tools/harness/lib-pin-compare.sh`,
`tools/harness/pin-stale.sh` (Refactor), `tools/harness/run-pin-stale-all-tests.sh`, `harness/sensors/pin-stale-all.md`,
`Makefile`, `.github/workflows/upstream-drift.yml`, `harness/README.md`, `docs/maintainer/releasing.md`, Slice-Plan.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03.

**Ablage:** Mutationen liefen an Kopien im Scratchpad (Änderung per `awk … > Kopie`, Prüfling per `PROG=<Kopie>`; nie
`sed -i`, keine Umleitung auf eine Repo-Datei). Es wurde kein Image gebaut. Kein verweigerter Aufruf.

**Eingangs-Kontext:** Slice-Plan `pin-stale-alle-digest-pins` (§1 bis §6, Implementer-Beleg in §3);
[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md);
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7;
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md); `AGENTS.md` (§3.1, §3.6, §3.7, §3.8, §3.9, §3.10,
§3.12, §3.13); `harness/conventions.md`.

**Eigene Messungen** (Exit-Codes direkt gesichert):

- `make test-pin-stale-all` Exit 0, gedruckt „alle 34 Prüfungen bestanden“. `make gates` Exit 0. `make test` Exit 0.
  `make fmt-check` Exit 0. `make kommentar-kennungen DIFF=3eb60de6` Exit 0, kein Kandidat.
- `make pin-stale-all` Exit 2 (über `make`), Schlusszeile `pin-stale-all: 15 Referenzen — 11 OK, 0 DRIFT, 4
  UNBESTIMMT`. Die vier `UNBESTIMMT` (`nats:2-alpine`, `eclipse-temurin:21-jdk`, `:21-jre`, `python:3.14-slim`) sind
  Docker-Hub-Abruflimit: `docker buildx imagetools inspect nats:2-alpine …` druckt `429 Too Many Requests`. Ein Lauf
  mit `UNBESTIMMT` ohne `DRIFT` ist kein Beleg für „kein Drift“ (Vertrag Grenze 5). Die PostgreSQL-17-Referenz
  (`.github/workflows/e2e.yml:80`) steht als `OK` mit Index-Digest.
- Refactor: `make pin-stale-acheck`, `-dmigrate`, `-dcheck`, `-baseline`, `-actions` Exit 0 mit `OK`-Zeilen in der
  Form der Vorgänger; `make pin-stale-race`, `-pgtest` und `make image-stale` Exit 2 durch dasselbe Abruflimit
  (`UNBESTIMMT`-Zeile, Form unverändert). Der Diff von `pin-stale.sh` ersetzt nur den Block Abfrage/Vergleich durch
  `pin_compare_digest`; die drei Zeilenformen und die Rückgaben 0/1/2 sind im Text von `lib-pin-compare.sh` gleich dem
  Parent.
- Mengen: `git grep -ohE '<Muster der ADR>' -- . ':!docs' ':!.harness' | sort -u | wc -l` = 15; `git grep -lE
  '@sha256:[0-9a-f]{64}' -- . ':!docs' ':!.harness' | wc -l` = 20; 48 Trefferzeilen. In `tools/harness/pin-stale-all.sh`,
  `run-pin-stale-all-tests.sh`, `lib-pin-compare.sh`, `harness/sensors/pin-stale-all.md` und `harness/README.md` kein
  vollständiger Digest (`git grep -l` leer): die Klausel hält, jede der 15 Referenzen ist ein echter Pin.
- `git grep -n -i "neun-achsen|P1.P9|neun achsen|P1-P9"` im lebenden Baum: Treffer nur in `ADR-0051` (`Accepted`,
  unberührbar), `ADR-0146` (zitiert die Altform), `docs/plan/planning/done/**` (Records) und im Suchlauf-Feld des
  eigenen Plans; kein lebender Träger mit dem alten Umfang.
- `git diff 3eb60de6 HEAD -- .github`: ein hinzugefügter Schritt (`if: always()`, `run: make pin-stale-all`), Kopf-
  Kommentar und Job-`name:` angepasst, keine neue `uses:`-Zeile, `timeout-minutes: 15` unverändert. `GATE_CHECKS` und
  `make gates` unberührt (`git diff -- Makefile`: nur `.PHONY`, zwei Ziele, Kommentar).
- **Mutationen** (Prüfling `PROG=<Kopie>`, Instanz `run-pin-stale-all-tests.sh`):
  (A) Exit 1 bei `DRIFT` auf 2: rot, Exit 1, drei Prüfungen (Digest weicht ab, Einzelplattform-Digest, `DRIFT` neben
  `UNBESTIMMT`) — vom Implementer als „hergeleitet“ geführt, hier gesehen.
  (B) Deduplizierung entfernt: rot, Exit 1, Prüfung „ein Registry-Aufruf je verschiedene Referenz“ (2 Aufrufe statt 1) —
  ebenfalls „hergeleitet“ im Plan, hier gesehen.
  (C) Zeitlimit entfernt (`limit=()` in `lib-pin-compare.sh`): **grün**, Exit 0, 34 Prüfungen (siehe F-1).

---

## Findings

### F-1 — Das Zeitlimit `PIN_COMPARE_TIMEOUT` ist zugesagt, aber an keiner Stelle des Tabellentests gebunden

- `kategorie`: MEDIUM
- `quelle`: Maintainability; `AGENTS.md` §3.12 (Zusage ohne Bindung); Reviewer-Skill „fehlende Negativtests bei neuem
  öffentlichem Vertrag“
- `pfad`: `tools/harness/lib-pin-compare.sh:17`, `harness/sensors/pin-stale-all.md` (Tabellen „Ausgabe“ und
  „Overrides“), `tools/harness/run-pin-stale-all-tests.sh` (kein Fall)
- `befund`: Der Vertrag sagt „Zeitlimit von `PIN_COMPARE_TIMEOUT`, Default 60 s, abgelaufen heißt `UNBESTIMMT`“; die
  Mutation, die den `timeout`-Präfix entfernt, lässt alle 34 Prüfungen grün. Der Implementer benennt die Lücke im Plan
  („hat keinen Tabellenfall“), der Sensor-Vertrag nennt sie nicht unter den Grenzen.
- `verifizierbar`: ja — `make test-pin-stale-all` mit der Mutation (C) bleibt grün.
- `klasse`: Zusage ohne Test an ihrer Eingabeseite

### F-2 — Grenze 3 des Sensor-Vertrags beschreibt die Wirkung einer Registry-Adresse mit Port ungenau

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B; [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Festlegung 1 („jede Referenz der Form“)
- `pfad`: `harness/sensors/pin-stale-all.md` §Grenze Punkt 3
- `befund`: Der Text sagt, das Muster treffe „erst ab dem Segment nach dem Port“; der Treffer beginnt am Port selbst
  (`5000/<name>:<tag>@…`, weil `:` vor dem Namen das Muster trennt), das Abfrageziel ist dann `5000/<name>:<tag>` und
  nicht ein Teil der Adresse. Am Stand trägt keine der 15 Referenzen einen Port (Lauf oben), die Grenze ist
  hergeleitet, nicht an einer Instanz gefahren.
- `verifizierbar`: nein — ohne Fixture mit Port-Adresse nicht an einem Lauf zu belegen.
- `klasse`: Grenzbeschreibung ungenauer als die Messung

### F-3 — `PIN_COMPARE_TIMEOUT` steht nur im Implementer-Beleg, nicht in ADR oder Plan-Umfang; Worst-Case entspricht dem Job-Limit

- `kategorie`: INFO
- `quelle`: [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 4;
  Slice-Plan §6
- `pfad`: `tools/harness/pin-stale-all.sh:16`, Plan §3 „Implementer-Beleg“
- `befund`: Die ADR schweigt zum Zeitlimit; sie widerspricht ihm nicht (Festlegung 4 verlangt `UNBESTIMMT` je Referenz,
  ein abgelaufener Aufruf ist genau das, Exit-Folge `UNBESTIMMT` → 2, nie `DRIFT`). Der Plan-Nachzug steht als
  Implementer-Beleg und wird dort als Antwort auf das Laufzeit-Risiko benannt, ein ADR-Nachzug ist nicht nötig. Bei 15
  Referenzen und 60 s je Aufruf wäre der Worst-Case 15 min und damit gleich `timeout-minutes: 15` des Jobs; gemessen sind
  28,8 s (Implementer, übernommen) bei erreichbarer Registry.
- `verifizierbar`: nein.
- `klasse`: Nachzug nur im Bericht des Implementers

### F-4 — Am Stand führt das Docker-Hub-Abruflimit zu `UNBESTIMMT` bei 4 von 15 Referenzen

- `kategorie`: INFO
- `quelle`: [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Frage (d) und
  Re-Evaluierungs-Trigger
- `pfad`: `harness/sensors/pin-stale-all.md` §Grenze Punkt 5
- `befund`: Mein Lauf endete Exit 2 durch 429, nicht durch Drift; zehn der 15 Referenzen liegen auf Docker Hub, P1 bis P9
  fragen dieselbe Registry im selben nächtlichen Lauf. Ob der gehostete Runner dasselbe Limit trifft, ist nicht
  gemessen; der Vertrag nennt die Grenze und hält fest, dass Exit 2 kein Beleg für „kein Drift“ ist.
- `verifizierbar`: nein (Zustand der Registry); Post-Push-Lauf nach `AGENTS.md` §3.10 trägt den Beleg.
- `klasse`: Registry-Abruflimit bei Sammel-Sensor

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/pin-stale-all.sh` gegen [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 1 (Muster wörtlich, `git grep`,
  Pathspec `':!docs' ':!.harness'`, Dedup nach Wert, keine Fundort-Liste), 2 (Index-Digest über `lib-pin-compare.sh`), 3
  (Tag-Ableitung am Namensteil nach dem letzten `/`, sonst `:latest`; im realen Lauf `aquasec/trivy:latest`,
  `d-migrate:latest`, `a-check:latest`), 4 (`UNBESTIMMT` je Referenz, Exit 1/2/0, leerer Gegenstand und `git grep`-Fehler
  Exit 2), 5 (ein Target, ein Schritt, `if: always()`), 6 (kein Tag-Wechsel)
- geprüft, ohne Befund: Refactor `lib-pin-compare.sh`/`pin-stale.sh` (Ausgabeformen, Rückgaben, Verhalten ohne
  `PIN_COMPARE_TIMEOUT` unbegrenzt wie vorher)
- geprüft, ohne Befund: Tabellentest (je Zweig der Festlegungen 1 bis 4 ein Fall mit Meldungstext; zwei weitere als
  „hergeleitet“ geführte Zusagen rot gesehen, siehe Mutationen A und B)
- geprüft, ohne Befund: Klausel „kein Digest-Literal im Skript, Test, Vertrag, Doku“; Zahl 15/20 nachgemessen
- geprüft, ohne Befund: Verdrahtung (`Makefile`, `harness/README.md`-Zeilen in der Form der Nachbarn,
  `upstream-drift.yml`, `docs/maintainer/releasing.md` §5 samt Version 1.14 und Historienzeile); kein `docs/user/`-Zug
  nötig (keine Betreiber-Oberfläche des Feed-Containers)
- geprüft, ohne Befund: §3.1 (Host-Werkzeuge nur `bash`, `git`, `docker`, `timeout`; kein `curl`/`jq`/`python`, keine
  Umleitung auf Repo-Dateien), §3.6 (keine Gate-Lockerung, keine neue Gate-Aufnahme), §3.7 (`make kommentar-kennungen
  DIFF=3eb60de6` Exit 0; Kommentare im Indikativ, ein Anker), §3.8 (keine neue `uses:`-Zeile), §3.13 (kein lebender
  Träger mit „Neun-Achsen“/„P1–P9“), SDK-Versionen unberührt
- Offen für den Verifier: §3.10 — der realer Post-Push-Lauf des zehnten Schritts von `upstream-drift.yml` ist nicht
  Gegenstand dieses Reviews.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 2 |

## Verdikt

Nicht merge-blockierend. Die Implementierung folgt [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) an jeder geprüften Festlegung; F-1 ist eine fehlende
Bindung einer benannten Zusage (Tabellenfall mit hängendem Stub-`docker` und kleinem `PIN_COMPARE_TIMEOUT` schließt sie,
das Urteil darüber liegt beim Implementer), F-2 ist eine Wortlaut-Berichtigung im Vertrag. Weil F-1 ein MEDIUM mit
Rückgabe-Pfeil an den Implementer ist, bleibt die DoD-Zeile „Review durchgeführt“ im Plan offen und wird bei der
Fixrunde nachgezogen (Reviewer-Skill „DoD-Checkbox-Nachzug ohne Fixrunde“).
