# Review-Report: slice-upgrade-versionswechsel-alt-image — 2026-10-04

**Review-Art:** Code — geprüft gegen Plan, [`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Teil 2,
[`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md), [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `upgrade-versionswechsel-alt-image`, Diff-Range `b965d80a..HEAD` (`HEAD` = `552919f7`, zwei
Implementer-Commits: `1ec033ff` Runner, Vertrag, `harness/mk/sdk.mk`, `harness/README.md`; `552919f7` Plan-Belege und
Register-Evidenz). Sechs Dateien, +232/−10; kein Go-, Spec-, Handbuch-, SDK-, Test-, `compose.yaml`- oder
`Makefile`-Diff, `tools/harness/run-integration-tests.sh` unberührt (gemessen: `git diff b965d80a HEAD --stat` mit den
Pathspecs `*.go spec docs/user sdks tools/harness/run-integration-tests.sh test/ compose.yaml Makefile` leer).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-04.

**Ablage / Vorgehen:** Der Reviewer-Lauf hat diesen Report mit dem Write-Werkzeug geschrieben. Die Mutation lief an einer
Scratchpad-Kopie des Runners (`sed … Datei > Kopie`, Ausgabe nach stdout; nie `sed -i`, nie eine Umleitung auf eine
Repo-Datei). Es wurde kein Image gebaut; `:dev` blieb unberührt (Image-ID vor dem Lauf `sha256:f2ab8eb9…`, vom Runner nur
gelesen). Nach den Läufen: `docker ps -a` und `docker network ls` ohne `cdc-*`, `git status --short` leer.
**Eigener Fehlgriff, vom Guard verhindert:** ein Aufruf enthielt `python3 --version` (Host-Interpreter am Kopf); der
PreToolUse-Guard blockte ihn vor der Ausführung, es lief nichts. Der Aufruf wurde gestrichen, nicht auf anderem Weg
wiederholt (`AGENTS.md` §3.15).

**Eingangs-Kontext:**

- Slice-Plan `slice-upgrade-versionswechsel-alt-image` (in-progress)
- [`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Teil 2,
  [`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md),
  [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-QA-OPS-005`](../../spec/lastenheft.md)
- `AGENTS.md` (§3.1, §3.7, §3.9, §3.12, §3.13, §3.15), `harness/conventions.md`
- Bestehender Rundlauf `tools/harness/run-integration-tests.sh` (Upgrade-Sicherheits-Rundlauf, Zeilen 5150–5245) als
  Vergleichsmaßstab
- Vorherige Findings am Modul: `review-slice-sdk-0-6-kompatibilitaet-messen`

---

## Findings

### F-1 — `make pin-stale-all` am Endstand nicht Exit 0 (Registry nicht erreichbar)

- `kategorie`: INFO
- `quelle`: [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) (P10-Klausel), Plan-DoD „P10-Klausel“
- `pfad`: `docs/plan/planning/in-progress/slice-upgrade-versionswechsel-alt-image.md` (DoD „P10-Klausel“)
- `befund`: Eigener Lauf `make pin-stale-all` (Exit direkt gelesen): Exit 2, „16 Referenzen — 7 OK, 0 DRIFT, 9
  UNBESTIMMT“, alle neun UNBESTIMMT sind Docker-Hub-/Registry-Fehler („Registry nicht erreichbar“); die Zeile
  `tools/harness/run-sdk-altserver-tests.sh:42 (ghcr.io/pt9912/pg-change-feed:0.5.0) == sha256:f99a77ff…` steht auf `OK`.
  Der Implementer meldete 13 OK / 3 UNBESTIMMT; die Differenz ist Erreichbarkeit zum Messzeitpunkt, kein DRIFT. Die
  DoD-Zeile „Exit 0“ ist an diesem Stand nicht belegbar; das ist Sache des Verifiers nach Wartezeit.
- `verifizierbar`: ja — `make pin-stale-all` bei erreichbarer Registry
- `klasse`: Netz-abhängiger Beleg nicht reproduzierbar

### F-2 — Mutationen M1 und M3 nicht nachgefahren; Aussage ist übernommen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Instanz B)
- `pfad`: `harness/targets/sdk-altserver.md` §Test (Phase-U-Tabelle); Plan §7 Mutationstabelle
- `befund`: Ich habe M2 selbst gefahren (siehe Prüfprotokoll); M1 (Änderung einer `cdc.change`-Zeile nach dem Tausch) und
  M3 (Override behält `image:`) habe ich nicht nachgefahren — für den Reviewer **übernommen** aus Plan und Vertrag. Die
  Vertragstabelle trägt je Zusage Mutation, Instanz und gesehene Farbe und kennzeichnet die Verallgemeinerung als
  *hergeleitet*; das genügt der Form.
- `verifizierbar`: ja — Mutation an einer Kopie, `bash <Kopie>`
- `klasse`: Mutation vom Reviewer nicht nachgefahren

### F-3 — Schlusszeile „Datenstand“ ohne Mengenvermerk

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Instanz A/B)
- `pfad`: `tools/harness/run-sdk-altserver-tests.sh:403` (Schlusszeile „erhält den Datenstand (U)“)
- `befund`: Die Schlusszeile sagt „erhält den Datenstand“ ohne Begrenzung. Gemessen wird Zeilenzahl und `md5` über alle
  `cdc.changes` der Quelle (`SOURCE_ID`), nicht nur `feed_e2e_full` — die Menge ist korrekt und eher weiter als
  befürchtet; die Begrenzung (Schema konstant, ein Alt-Stand, INSERT/UPDATE/DELETE auf einer Tabelle als einzige
  erzeugte Last, Übertragung *hergeleitet*) steht in Vertrag Grenze 6, nicht in der gedruckten Zeile. Kein Widerspruch,
  nur ein Hinweis zur Lesbarkeit der Zeile als Beleg.
- `verifizierbar`: nein
- `klasse`: Beleg-Zeile ohne Mengenvermerk

---

## Prüfprotokoll (je Schwerpunkt)

1. **Prüfumfang gegen [`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Teil 2 / [`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md).**
   Die ADR verlangt: Container aus 0.5.0, Zeilen erfassen, `--force-recreate` auf `:dev`, „derselbe Prüfumfang“ wie der
   Rundlauf (Datenstand vor/nach identisch lesbar, danach eingefügte Zeile erfasst). Der bestehende Rundlauf
   (`run-integration-tests.sh:5150–5245`) prüft: Container-ID neu, postgres/nats-ID unverändert, Health `healthy`,
   `count(*) = 1` der Zeile id=250, danach id=251 erfasst. Phase U (`run-sdk-altserver-tests.sh`, `phase_u`) prüft
   dieselben Punkte und mehr: Zeilenzahl **und** `md5` über `change_id`/`commit_position`/beide Bilder aller Changes der
   Quelle vor/nach (stärker als `count = 1`), Image-Referenz und Image-ID des neuen Containers gleich Ziel und ungleich
   Start, Slot da mit nicht kleinerer `confirmed_flush_lsn`, Zeilenzahl danach genau +1. Der beschlossene Umfang ist
   erfüllt und übertroffen.
2. **Trivial grün möglich?** Nein, nach Quelltext und Lauf. Jede Prüfung kann an ihrer Eingabeseite rot werden:
   Container-ID (`feed_neu = feed_alt`), Build-Wechsel (`neu_ref`, `neu_id`, `ziel_id`, `start_id`; im B3-Modus fängt
   das Überspringen die Phase ab, siehe 4.), Momentaufnahme (md5 schlägt an, wenn der Neustart Changes dupliziert oder
   ändert), Erfassung (Warteschleife auf id=9102, 120 mal 0,25 s). Eigene Mutation **M2** (Kopie des Runners,
   `docker stop "$FEED_CONTAINER"` vor dem INSERT von id=9102): Exit 1, `U ROT — die nach dem Tausch eingefügte Zeile
   (id=9102) erscheint nicht über cdc.changes, die Erfassung setzt nicht fort` (gesehen, `bash <Kopie>`). M1/M3:
   übernommen (F-2). Verallgemeinerung auf die übrigen Prüfungen: *hergeleitet*.
3. **Ziel-Image per `awk` aus `docker compose config`.** Kein Tag-/Digest-Literal; Ergebnis im Lauf
   `ghcr.io/pt9912/pg-change-feed:dev`. Fehlerpfade am Quelltext gelesen: `docker compose config` schlägt fehl → unter
   `set -euo pipefail` beendet die Zuweisung den Runner mit Exit ≠ 0 (kein stilles Grün); leeres Ergebnis → Meldung +
   Exit 1; Image nicht geladen → „make image vorher“ + Exit 1, vor jedem Start (nicht gefahren). Die `awk`-Erkennung
   hängt an der Zwei-Leerzeichen-Einrückung der normalisierten Compose-Ausgabe; das ist die ausgegebene Form von
   `docker compose config`. `awk` ohne `-i` ist im Vertrag benannt und nach `AGENTS.md` §3.1 zulässig.
4. **Vorprüfung / B3.** `:dev` fehlt → Exit 1 vor jedem Start (am Quelltext). B3 (`SDK_ALTSERVER_IMAGE` gleich Ziel):
   B1 wird rot; ohne `SDK_ALTSERVER_WEITER` bricht `befund` mit `exit 1` ab (Phase U nicht erreicht), mit `WEITER=1`
   druckt U `ÜBERSPRUNGEN` und `rot=1` hält den Gesamtausgang auf Exit 1 — kein stilles Grün im B3-Modus. B3 selbst
   habe ich nicht gefahren (Plan-Beleg des Implementers übernommen).
5. **[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)-Klausel.**
   `git diff b965d80a HEAD | grep -E '^\+.*@sha256:'`: die Treffer stehen in Plan-Prosa und Register-Evidenz
   (Messzeilen bzw. Zitat der Pin-Form), kein neues Literal in Skript, Makefile oder Vertrag (Suchlauf `@sha256:`
   außerhalb `docs`/`.harness`: 86 gleich Parent). Der Alt-Pin bleibt EIN Literal mit Tag
   (`run-sdk-altserver-tests.sh:42`). `make pin-stale-all`: siehe F-1.
6. **Mutationen.** M2 selbst gefahren (oben). Zusage · Mutation · Instanz · Rot der drei Tabellenzeilen stehen im
   Vertrag und im Plan; je eine Mutation je Zusage; Verallgemeinerung ausdrücklich *hergeleitet* (Vertrag Grenze 6).
7. **Verträge/Doku.** `harness/targets/sdk-altserver.md`: Phase-U-Zeile, Aufruf (`make image` zuerst), Host-Werkzeug
   `awk`, Exit-Tabelle, Grenze 6 konsistent mit dem Runner. `harness/mk/sdk.mk`: Kommentar und Hilfetext nennen
   `make image` vorher. `harness/README.md`: zwei Zeilen (`make test-integration` mit einem Satz,
   `make test-sdk-altserver`); Suchlauf-Zählung `sdk-altserver` 2 gleich Soll. Register-Evidenz
   `…/kein-echter-versionswechsel-upgrade-test/evidence/slice-upgrade-versionswechsel-alt-image.md`: Form
   `Vorgang`/`Fund`/`Form (Ausprägung)` wie die Nachbardateien; `state.md` unberührt (Zähler 2× bleibt, „Ausgang setzt
   die Closure“ ausdrücklich vermerkt, nicht von Hand gesetzt).
8. **Unberührtes.** Kein Go-/Spec-/Handbuch-/SDK-/Test-/`compose.yaml`-/`Makefile`-Diff;
   `docs/user/e2e-abdeckung.md` unverändert; `make test-sdk-altserver` steht nicht in `GATE_CHECKS` (`make gates` ohne
   dieses Ziel, Exit 0).
9. **Suchlauf.** `make suchlauf-nachmessen PLAN=…` Exit 0: „10 Zeilen stimmen“. Die Abweichung Soll 8 statt 5 zu
   „desselben Images“ erklärt der Plan plausibel (das Muster trifft „versionswechsel“): die acht Trefferzeilen am
   Arbeitsbaum (`git grep -n -i -E` mit dem Muster des Plans) sind fünf Bestandszeilen (`docs/user/e2e-abdeckung.md:89`,
   `run-integration-tests.sh:50/5146/5151`, `harness/README.md:147`) und drei neue (`harness/README.md:185`,
   `run-sdk-altserver-tests.sh:26` und `:312`). Die Zuordnung der fünf als Bestand ist abgeleitet aus Soll 5 am Parent;
   den Parent habe ich nicht einzeln nachgemessen.
10. **Kommentare / Hard Rules.** `make kommentar-kennungen DIFF=b965d80a`: Exit 0, kein Kandidat. Die Skriptkommentare
    der Phase U tragen Zusage/Kopplung/Grenze im Indikativ, höchstens eine Kennung, keine Chronik. Kein
    `sed -i`/Umleitung auf Repo-Dateien im Diff (§3.1). Ersatzweg nach Verweigerung (§3.15): aus dem Diff nicht
    ablesbar, kein Hinweis darauf.

**Läufe (Exit direkt gelesen):** `make gates` 0 · `make test` 0 · `make fmt-check` 0 · `make sdk-public-doc-check` 0 ·
`make kommentar-kennungen DIFF=b965d80a` 0 · `make suchlauf-nachmessen PLAN=…` 0 · `make test-sdk-altserver` 0
(gedruckt: `ALTSERVER U: Tausch d7d65e690ea5 -> 1e6180195447, Image …:0.5.0@sha256:f99a77ff… -> …:dev, Datenstand vor dem
Tausch (4 Zeilen, Prüfsumme 2e5ff42113514778710c1e905ddd4e3a) identisch lesbar, danach eingefügte Zeile erfasst
(Position 30850848), Phase U 7 s`; danach kein `cdc-*`-Container/-Netz) · `make pin-stale-all` 2 (F-1) ·
`make docs-check` nach dem Anlegen dieses Reports gesondert (Commit-Schritt).
**Nicht gelaufen:** `make test-integration` (nicht Teil des Slice), B3-Läufe, Mutationen M1/M3.

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/run-sdk-altserver-tests.sh` (Phase U, Vorprüfung, Fehlerpfade, Aufräumen)
- geprüft, ohne Befund: `harness/targets/sdk-altserver.md` (Vertrag, Exit-Tabelle, Grenzen, Mutationstabelle)
- geprüft, ohne Befund: `harness/mk/sdk.mk` (Kommentar und Hilfetext)
- geprüft, ohne Befund: `harness/README.md` (beide Sensor-Zeilen, Gate-Bündel unverändert)
- geprüft, ohne Befund: `docs/plan/planning/` Slice-Plan und Register-Evidenz (Form, Zähler, `state.md` unberührt)
- geprüft, ohne Befund: `spec/`, `docs/user/`, `sdks/`, `test/`, `compose.yaml`, `Makefile`, `*.go`,
  `tools/harness/run-integration-tests.sh` (nicht berührt, Diff leer)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Netz-abhängiger Beleg nicht reproduzierbar · Mutation vom Reviewer nicht nachgefahren ·
Beleg-Zeile ohne Mengenvermerk

## Verdikt

**Merge-blockierend:** nein — kein HIGH, MEDIUM oder LOW; drei INFO ohne erwartete Fixrunde.

**Übergabe:** Keine Fixrunde am Implementer. F-1 an den Verifier (P10-Klausel „Exit 0“ nach Wartezeit bzw. erreichbarer
Registry nachfahren; UNBESTIMMT ist nicht DRIFT). F-2 und F-3 ohne erwartete Aktion. Die DoD-Zeile „Review
durchgeführt“ im Plan habe ich auftragsgemäß **nicht** gesetzt; sie ist Nachzug von Planner oder Implementer. Dieser
Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
