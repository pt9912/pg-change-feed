# `make pin-stale-all` — jede Digest-Pin-Referenz des Baums gegen die Registry (P10)

## Vertrag

Wird dieses Target rot, driftet mindestens eine Referenz der Form
`<image>[:<tag>]@sha256:<digest>` in einer getrackten Datei: die Registry
trägt unter dem gepinnten Tag (bzw. `:latest` bei einem Pin ohne Tag) einen
anderen Index-Digest als der Baum (`tools/harness/pin-stale-all.sh`,
[`ADR-0146`](../../docs/plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)).
Das Target ist **advisory**: es braucht Netz, steht nicht in `GATE_CHECKS` und
färbt `make gates` nicht; sein Ort ist der nächtliche Lauf von
`.github/workflows/upstream-drift.yml` (zehnter Schritt, `if: always()`).

**Gegenstand** ist die Menge der verschiedenen Referenzen, die
`git grep -nHoE` mit dem Muster
`[a-z0-9][a-z0-9./_-]*(:[A-Za-z0-9._-]+)?@sha256:[0-9a-f]{64}` in den
getrackten Dateien außerhalb von `docs/` und `.harness/` findet. Eine Datei,
ein YAML-Matrixwert, ein Shell-Default und ein `compose.yaml`-Default sind
gleichwertige Fundorte; es gibt keine Liste von Fundorten und keine
Ausnahmeliste. Dieselbe Referenz an mehreren Fundorten ist ein Eintrag (erster
Fundort, `+N weitere`) und ein Registry-Aufruf.

**Vergleich.** Der Maßstab ist der Index-Digest des Ziels
(`docker buildx imagetools inspect <ziel> --format '{{.Manifest.Digest}}'`);
die Vergleichslogik liegt in `tools/harness/lib-pin-compare.sh` und ist
dieselbe, die `tools/harness/pin-stale.sh` (P3 bis P6) benutzt. Das Ziel ist
die Referenz ohne Digest; trägt der Namensteil nach dem letzten `/` keinen
Tag, ist das Ziel `<image>:latest`. Ein Einzelplattform-Digest ist keine
zulässige Pin-Form dieses Repos und meldet `DRIFT`.

## Ausgabe und Ausgänge

Je Referenz eine Zeile; am Ende die Zusammenfassung
`pin-stale-all: <N> Referenzen — <a> OK, <b> DRIFT, <c> UNBESTIMMT`.

| Zeilenform | Bedeutung |
|---|---|
| `OK          <fundort> (<ziel>) == <digest>` | Registry und Baum gleich |
| `DRIFT       <fundort> (<ziel>): gepinnt <digest>, aktuell <digest>` | beide Digests stehen in der Zeile |
| `UNBESTIMMT  <fundort> — <ziel>: Registry nicht erreichbar` | Abfrage scheitert (Netz, Abruflimit, Zeitlimit von `PIN_COMPARE_TIMEOUT`, Default 60 s je Aufruf); die übrigen Referenzen laufen weiter |

| Exit | Bedeutung |
|---|---|
| 0 | jede Referenz `OK` |
| 1 | mindestens ein `DRIFT` (auch neben `UNBESTIMMT`) |
| 2 | `UNBESTIMMT` ohne `DRIFT`; leerer Gegenstand (`UNBESTIMMT  keine Digest-Pin-Referenz im Baum gefunden — leerer Gegenstand`); `git grep` mit Fehler |

Leer ist nicht bestanden. Über `make` kommt jeder Exit ≠ 0 des Skripts als der
Make-eigene Exit `2` an.

## Overrides

| Variable | Wirkung |
|---|---|
| `PIN_COMPARE_TIMEOUT` | Sekunden je Registry-Aufruf (Default 60 im Skript); abgelaufen heißt `UNBESTIMMT` |
| `PROG` (nur `tools/harness/run-pin-stale-all-tests.sh`) | Prüfling des Tabellentests, für Mutationsläufe an Kopien; `lib-pin-compare.sh` liegt neben dem Prüfling |

## Grenze — was das Grün nicht abdeckt

1. **Form, nicht Sinn.** Der Sensor meldet Veränderung, nicht Gefahr: ob ein
   Pin gehoben werden soll (Major-Wechsel, Kompatibilität), entscheidet ein
   Mensch. Er prüft nicht, ob mehrere Kopien desselben Pins untereinander
   gleich sind, nur dass jede gegen die Registry stimmt.
2. **Ein Tag-Wechsel ist nicht Gegenstand.** Der Sensor meldet Digest-Drift am
   gepinnten Tag, nicht, dass ein neuerer Major existiert.
3. **Was das Muster nicht trifft, sieht der Sensor nicht.** Ein Digest in einer
   zur Laufzeit zusammengesetzten Variablen oder ein Pin ohne `@sha256:` ist
   außerhalb; eine Registry-Adresse mit Port vor dem Namen trifft das Muster
   erst ab dem Segment nach dem Port. Der Treffer beginnt am Port selbst
   (`5000/<name>:<tag>@…`), das Abfrageziel ist dann die Referenz ohne Host
   und damit falsch — an einer Stub-Instanz im Tabellentest gefahren (Treffer
   `5000/…`, Aufruf ohne Host); wie eine reale Registry ein solches Ziel
   beantwortet, ist *hergeleitet*, an keiner Instanz gefahren.
4. **Eine referenz-ähnliche Zeichenkette ist ein Pin.** Ein Fixture oder
   Beispiel mit vollständigem Digest in einer getrackten Datei außerhalb von
   `docs/` und `.harness/` meldet der Sensor dauerhaft als `DRIFT` oder
   `UNBESTIMMT`. Deshalb tragen das Skript, sein Tabellentest und dieser
   Vertrag keine vollständige Referenz als Literal; die Fixture-Digests
   entstehen im Tabellentest zur Laufzeit. Eine Ausnahmeliste gibt es nicht
   (ein solcher Fund ist ein Befund für den Planner).
5. **Das Abruflimit der Registry** (Docker Hub, anonym, HTTP 429) macht
   Referenzen `UNBESTIMMT`, nicht `DRIFT`; ein Lauf mit `UNBESTIMMT` ohne
   `DRIFT` endet mit Exit 2 und ist kein Beleg für „kein Drift“, auch kein
   für Drift. Der Exit 2 färbt den Schritt im nächtlichen Lauf rot (advisory,
   blockiert nichts). Gemessen am Arbeitsplatz-Host am 2026-10-03:
   `make pin-stale-all` endete mit Exit 2, gedruckte Zeile `pin-stale-all: 15
   Referenzen — 6 OK, 0 DRIFT, 9 UNBESTIMMT`, die Registry antwortete mit `429
   Too Many Requests`. Wie oft der gehostete Runner dieses Limit trifft, ist
   nicht gemessen; das zeigt der erste `workflow_dispatch`-Lauf.
6. **Zeitlimit und Gesamtdauer.** `PIN_COMPARE_TIMEOUT` begrenzt jeden
   Registry-Aufruf (Default 60 s im Skript; `tools/harness/pin-stale.sh` setzt
   keinen Wert, dort ist der Aufruf unbegrenzt); ein Ablauf ist `UNBESTIMMT`,
   belegt im Tabellentest (Stub schläft länger als das Limit). Die Obergrenze
   des Schritts ist je Referenz das Limit mal die Zahl der Referenzen (15 × 60
   s = 15 min bei 15 Referenzen, *abgeleitet*) und liegt damit **nicht** unter
   `timeout-minutes: 15` des Jobs, der außerdem neun weitere Schritte trägt: ein
   Netzausfall aller Referenzen ließe den Job am Job-Limit enden, nicht am
   Skript. Gemessen: `real 0m28,845s` bei 15 Referenzen im Normalfall
   (Implementer-Beleg im Slice-Plan), `real 0m17,491s` am 2026-10-03 bei 9
   Antworten mit `429`.
7. **Überschneidungen** mit den neun benannten Achsen
   ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md)
   Entscheidung 7) melden denselben Drift doppelt im selben Lauf.
8. **Der Tabellentest ist Werkzeug, kein Gate**
   ([`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
   Teilfrage 2): `make test-pin-stale-all` läuft nicht in `make gates`; der
   Ort, an dem eine Änderung am Wächter gegen diesen Vertrag gelesen wird, ist
   der Review-Diff.

## Sperren

- **Host-Werkzeuge:** `bash`, `git`, `docker` (nur die Registry-Abfrage
  `buildx imagetools inspect`), `timeout` (coreutils). Kein `curl`, kein `jq`,
  kein Schreibzugriff auf den Baum. Braucht Netz; der Tabellentest ist netzlos
  (Stub-`docker`).

## Tabellentest

`make test-pin-stale-all` (`tools/harness/run-pin-stale-all-tests.sh`, netzlos)
fährt den Prüfling in einem Wegwerf-Repo mit Stub-`docker`, die
Fixture-Digests entstehen zur Laufzeit (Grenze 4), der Prüfling ist per
`PROG=<Datei>` übersteuerbar (§Overrides). Fälle, je mit Meldungstext: Treffer
in Workflow-YAML, Shell-Default und `compose.yaml`, Digest gleich und
abweichend, Einzelplattform-Digest als `DRIFT`, Pin mit und ohne Tag, `docs/`
und `.harness/` ausgenommen, dieselbe Referenz an zwei Fundorten als ein
Eintrag, Registry-Ausfall ohne Abbruch der übrigen, `DRIFT` neben
`UNBESTIMMT`, leerer Gegenstand.

## Bindung

[`ADR-0146`](../../docs/plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Festlegung 1 bis 6) · `tools/harness/pin-stale-all.sh` ·
`tools/harness/lib-pin-compare.sh` · `tools/harness/run-pin-stale-all-tests.sh`
· `.github/workflows/upstream-drift.yml` · seit slice-pin-stale-alle-digest-pins.
