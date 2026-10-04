# Verifikations-Report: slice-upgrade-versionswechsel-alt-image — 2026-10-04

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). Frischer Kontext. Gegenstand:
`git diff b965d80a HEAD` (7 Dateien, 409 Zeilen hinzu): Implementer-Commit `1ec033ff`
(Runner, Vertrag, `sdk.mk`, `harness/README.md`), `552919f7` (Plan-Belege, Register-Evidenz),
Review-Commit `1c359505` ([Review](review-slice-upgrade-versionswechsel-alt-image.md): 0 HIGH/MEDIUM/LOW,
3 INFO), Review-Haken `1383d75e`. Entscheidungsgrundlage:
[`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Teil 2,
[`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md),
[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md);
Anforderung [`LH-QA-OPS-005`](../../spec/lastenheft.md).
Plan: [`slice-upgrade-versionswechsel-alt-image`](../plan/planning/done/slice-upgrade-versionswechsel-alt-image.md).

**Verdikt: bestanden** — mit den Bedingungen und offenen Punkten in §6 (keine blockierende).
Alle Läufe wurden in diesem Kontext selbst gefahren, Exit direkt (nie durch eine Pipe, §3.9), je als
eigener Hintergrund-Schritt mit gesichertem Exit. Mutationen liefen ausschließlich an Kopien des
Runners im Scratchpad (`sed … > Kopie`, kein `-i`), Aufruf `bash <Kopie>` aus der Repo-Wurzel. Der
Arbeitsbaum blieb unberührt (`git status --short` am Ende leer). Docker-Hub-Abruflimit: in keinem
Lauf aufgetreten (kein 429, keine Pull-Fehler; ob die Läufe aus dem lokalen Cache zogen, ist nicht
gemessen).

## 1. Eigene Läufe (gedruckt, Exit direkt)

| Lauf | Exit | Gedruckte Zeile (Auszug) |
|---|---|---|
| `make image` (zuerst) | 0 | `naming to ghcr.io/pt9912/pg-change-feed:dev done` |
| `make test-sdk-altserver` | 0 | `ALTSERVER U: Tausch 9be019bbfe1f -> ff6fe02469d7, Image ghcr.io/pt9912/pg-change-feed:0.5.0@sha256:f99a77ff… -> ghcr.io/pt9912/pg-change-feed:dev, Datenstand vor dem Tausch (4 Zeilen, Prüfsumme 1a71d28c02d907c4d56c440d4d2f79b4) identisch lesbar, danach eingefügte Zeile erfasst (Position 30850848), Phase U 7 s`; Schlusszeile `Altserver-Messung grün — … der Tausch auf ghcr.io/pt9912/pg-change-feed:dev erhält den Datenstand (U)`; davor `ALTSERVER B0`, `B1`, `B2 csharp/kotlin/python` je `RECEIVED code=none http=404 grpc=NOT_FOUND text=74 diag_error_code=leer (Exit 0)` |
| B3 `SDK_ALTSERVER_IMAGE=…:dev make test-sdk-altserver` | 2 (Runner 1) | `B1 ROT — der Fehlerkörper … ist nicht der eines Servers ohne Meldungscodes (Ausgang 1)`; kein Phase-U-Lauf |
| B3 mit `SDK_ALTSERVER_WEITER=1` | 2 | B1 rot, B2 csharp/kotlin/python rot (`Marker NORMAL_DONE blieb aus`), `ALTSERVER U ÜBERSPRUNGEN: die Start-Referenz ghcr.io/pt9912/pg-change-feed:dev ist die Ziel-Referenz …`, `Altserver-Messung ROT` |
| `make pin-stale-all` | **0** | `pin-stale-all: 16 Referenzen — 16 OK, 0 DRIFT, 0 UNBESTIMMT`; Zeile `OK tools/harness/run-sdk-altserver-tests.sh:42 (ghcr.io/pt9912/pg-change-feed:0.5.0) == sha256:f99a77ff…` |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`; `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`; `generated-sync: OK`; `a-check gesamt: 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s)`; `d-check: 1648 Datei(en) geprüft, 0 Befund(e)` |
| `make docs-check` | 0 | `d-check: 1648 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make test` | 0 | alle Pakete `ok` (letzte Zeile `tools/schema/rolloutguard ok`) |
| `make fmt-check` | 0 | `334 Go-Dateien geprüft, alle formatiert` |
| `make sdk-public-doc-check` | 0 | `keine interne Kennung unter sdks` |
| `make doc-immutable RANGE=b965d80a..HEAD` | 0 | `0 Befund(e)` |
| `make doc-commits RANGE=b965d80a..HEAD` | 0 | `0 Befund(e)` |
| `make commit-traceability` | 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make kommentar-kennungen DIFF=b965d80a` | 0 | keine Ausgabe, kein Kandidat |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |
| `make suchlauf-nachmessen PLAN=…` | 0 | `suchlauf-nachmessen: 10 Zeilen stimmen` (`diff`-Zeilen 2, 86, 4, 8, 1) |

Nach allen Läufen: `docker ps -a` ohne `cdc-*` (der Container `gitea` gehört nicht zu diesem
Repo und bestand vor dem Start; ein zu Beginn gesehener `jolly_nash` war am Ende verschwunden), `docker network
ls` ohne `cdc-*` (Zählung 0). Der Aufräum-`trap` trifft damit nach dem Override-Wechsel das richtige
Projekt (Risiko „Aufräum-trap“ gemessen entfallen).

**F-1 des Reviews (`make pin-stale-all` Exit 2) ist an meinem Stand geschlossen:** Exit 0, 16 OK,
0 UNBESTIMMT. Die P10-Klausel „Exit 0, Zeile zum Altserver-Runner `OK`“ ist belegt. Der Befund des
Reviews war eine Erreichbarkeits-Lage des Laufs, kein Drift dieses Slice.

## 2. Mutationen M1 bis M3 — selbst nachgefahren

Instanz: PostgreSQL und Docker-Daemon des Runners, je ein Lauf. Mutanten per `sed` nach stdout in
Kopien erzeugt, Diff gegen den Runner vor dem Lauf gelesen (je genau die beabsichtigte Zeile).

| Zusage | mutierte Eingabe (Kopie) | gesehene Farbe |
|---|---|---|
| M1 Datenstand identisch | vor der Prüfsumme nach dem Tausch `UPDATE cdc.change SET new_data = '{"m":1}'::json WHERE new_data->>'name' = 'altserver-b0'` | Exit 1, `U ROT — der Datenstand der Quelle über cdc.changes ist nach dem Tausch nicht identisch (vorher: 4 68c4027e…, nachher: 4 3263fa78…)` |
| M2 Erfassung setzt fort | `docker stop` des Feed-Containers nach dem Health-Poll des Tauschs | Exit 1, `U ROT — die nach dem Tausch eingefügte Zeile (id=9102) erscheint nicht über cdc.changes, die Erfassung setzt nicht fort` |
| M3 Tausch wechselt den Build | Override des Tauschs trägt weiter `image: <SDK_ALTSERVER_IMAGE>` | Exit 1, `U ROT — der neue Container läuft nicht den Ziel-Build (Referenz …:0.5.0@sha256:f99a77ff…, Image-ID sha256:af59013f…; Ziel …:dev, sha256:f2ab8eb9…; Start sha256:af59013f…)` |

Alle drei Rot-Farben stimmen mit der Behauptung in Plan §3 und Vertrag überein; die Prüfsummen von M1
unterscheiden sich im „nachher“-Wert vom Plan (zeit-/LSN-abhängig, die Aussage ist die Farbe, nicht der
Hash). Das schließt F-2 des Reviews (M1/M3 nicht nachgefahren): beide bestätigt. M2 ist an dieser
Stelle strukturell die Probe der Zeile „danach eingefügt“; ein Stopp *vor* dem Tausch würde an der
Container-ID bzw. Health rot, das ist nicht gefahren (*hergeleitet*).

## 3. DoD-Zeile für Zeile (Beleg gegen Behauptung)

| Zeile | Stand im Plan | Befund des Verifiers |
|---|---|---|
| Liefer-Punkt 1 (Phase U) | `[x]` | Belegt. Runner gelesen (Zeilen 103–116 Vorprüfung „make image vorher“, 284–413 Phase U): Vorprüfung vor jedem Start, Start gleich Ziel druckt `ÜBERSPRUNGEN` (B3 mit WEITER bestätigt); INSERT/UPDATE/DELETE `id=9101`; Momentaufnahme `count(*)` + `md5` über **alle** `cdc.changes` der Quelle; Tausch ohne `image:`-Zeile; geprüft Container-ID, `postgres`/`nats`, Referenz **und** Image-ID gleich Ziel und verschieden vom Start, Health, Datenstand identisch, Zeile `id=9102`, Slot mit nicht kleinerer `confirmed_flush_lsn`, Zeilenzahl n+1. Gedruckte Zeile wie gefordert (plus Dauer). |
| Liefer-Punkt 2 (Mutationen) | `[x]` | Belegt, §2. B3 rot an B1, Lauf kommt nicht bis U (Exit 2). |
| Liefer-Punkt 3 (Verträge, Verdrahtung, Register) | `[ ]` | Inhaltlich geliefert, Haken ist Planner-Arbeit (siehe §6). Vertrag, `sdk.mk`, README-Zeilen gelesen (§5), Register-Evidenz liegt (§7); `state.md` bewusst unverändert. |
| `make gates` grün | `[ ]` | Erfüllt am Stand `1383d75e`, §1 (Haken Planner, am Endstand der Closure erneut). |
| `make docs-check` Exit 0 | `[ ]` | Erfüllt vor Report; nach `git add` dieses Reports erneut, siehe §8. |
| `suchlauf-nachmessen` Exit 0 | `[ ]` | Erfüllt, 10 Zeilen. |
| `fmt-check`/`test` unberührt, kein `*.go` im Diff | `[ ]` | Erfüllt: Diff-Stat über `*.go`, `spec`, `docs/user`, `sdks`, `test`, `compose.yaml`, `Makefile`, `run-integration-tests.sh` leer (0 Zeilen); beide Läufe Exit 0. |
| `kommentar-kennungen DIFF=` | `[ ]` | Erfüllt, kein Kandidat. |
| P10-Klausel | `[ ]` | Erfüllt: `make pin-stale-all` Exit 0 (16 OK). Kein neues `@sha256:`-Literal in Skript, Makefile, Vertrag: `git diff b965d80a HEAD \| grep -E '^\+.*@sha256:'` trifft nur Plan-Prosa, Register-Evidenz und Review-Zitat (4 Zeilen, alle unter `docs/`); Suchlauf `@sha256:` `diff` Soll 86 gleich Parent. |
| `make test-sdk-altserver` Exit 0, `ALTSERVER U`-Zeile, kein `cdc-*`-Rest, Arbeitsbaum unberührt | `[ ]` | Erfüllt, §1 (`git status --short` vor und nach leer). |
| `make test-integration` nicht Teil | `[ ]` | Zusage eingehalten: `run-integration-tests.sh` und `docs/user/e2e-abdeckung.md` ohne Diff. Nicht gefahren, nicht nötig. |
| Review durchgeführt | `[x]` | Report liegt, 0 HIGH/MEDIUM/LOW. |
| Doku-Update | `[ ]` | Zusage eingehalten (nichts unter `docs/user/`, `spec/`, `sdks/`). |
| Closure-Notiz, Register, §6-Ausgänge, Paarungen | `[ ]` | Planner-Arbeit, offen. |

## 4. Zusage „Datenstand“ und die Aussage „Server 0.5.0 übersteht den Tausch“ (§3.12 Instanz B)

- **Menge der Prüfung:** Zeilenzahl und `md5` über `change_id`, `commit_position`, `old_data`, `new_data`
  aller `cdc.changes` der Quelle (`SOURCE_ID`), in `change_id`-Reihenfolge — das ist weiter als
  „nur `feed_e2e_full`“ (Review F-3 bestätigt: Menge korrekt, eher weiter). Im Lauf sind es 4 Zeilen
  (B0-Zeile plus INSERT/UPDATE/DELETE von `id=9101`, *abgeleitet*, nicht je Zeile gelesen). Der
  Vergleich ist der Nachweis „gleiche Zeilen, gleiche Bilder, gleiche Positionen“, nicht „alle
  Operationen und Tabellen“.
- **Wahr formuliert:** Plan §3 und Register-Evidenz nennen „Schema konstant, ein Alt-Stand, keine
  Zeilen während des Tauschs“ und stellen die Übertragung als *hergeleitet*; Vertrag Grenze 6
  trägt dieselbe Begrenzung. Die Aussage „Server 0.5.0 übersteht den Tausch bei konstantem Schema“ ist
  als *gemessen* (ein Lauf des Implementers, ein weiterer Lauf in diesem Report: beide `Exit 0`, eine
  Quelle, eine Tabelle) richtig begrenzt. Die Prüfsummen beider Läufe unterscheiden sich
  (`2e5ff421…` Implementer, `1a71d28c…` Verifier); der Wert hängt an LSN-Positionen, nur die Gleichheit
  vor/nach dem Tausch innerhalb eines Laufs trägt die Aussage (so auch formuliert).
- **Restlücke (INFO, benannt):** die Schlusszeile des Runners „erhält den Datenstand (U)“ trägt keinen
  Mengenvermerk (Review F-3). Die Begrenzung steht im Vertrag, nicht in der gedruckten Zeile; die
  gedruckte `ALTSERVER U:`-Zeile selbst nennt `4 Zeilen` und die Prüfsumme und ist damit
  ausreichend belegbar. Kein Befund gegen die DoD.
- **Nicht gemessen:** der Tausch bei laufendem Schreiber, Rückweg `:dev` → 0.5.0, Schemawechsel
  (akzeptiertes Negativ von [`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)
  Option C). Alle drei stehen im Plan §1 bzw. Vertrag Grenze 6.

## 5. Plan-vs-Code-Diff, Entscheidungs-Konformität

- **Umfang von [`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Teil 2:** „eine Phase mehr: Feed-Container aus dem veröffentlichten Image
  0.5.0 startet, erfasst Zeilen, wird per `--force-recreate` durch das `:dev`-Image ersetzt; derselbe
  Prüfumfang wie im bestehenden Rundlauf (identisch lesbar, Erfassung setzt fort)“. Der bestehende
  Rundlauf (`git grep -n 'force-recreate' tools/harness/run-integration-tests.sh`, Block um Zeile 5191)
  prüft Container-ID neu, `postgres`/`nats` unberührt, Health, Vorher-Zeile lesbar, Nachher-Zeile
  erfasst. Phase U deckt dies vollständig ab und ergänzt: Referenz **und** Image-ID gleich Ziel und
  verschieden vom Start (schließt „trivial grün“), Prüfsumme statt Einzelzeile, Slot-Position,
  Zeilenzahl n+1. Konform; Zuschnitt (Altserver-Runner statt Integrations-Runner) wie [`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)
  zulässt, Begründung im Plan §1.
- **[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md):** ein Pin-Literal (Zeile 42 des Runners, `OK` im Pin-Inventar), Ziel-Image aus
  `compose.yaml` gelesen (`awk` auf `docker compose config`), kein zweites Literal. Konform.
- **Unberührtes:** Diff-Stat über `*.go`, `spec`, `docs/user`, `sdks`, `test`, `compose.yaml`,
  `Makefile`, `tools/harness/run-integration-tests.sh` leer; `docs/user/e2e-abdeckung.md` unverändert;
  kein Dockerfile im Diff. Keines der Targets steht in `GATE_CHECKS` (`git grep altserver -- Makefile
  harness/mk` trifft nur `harness/mk/sdk.mk`; `GATE_CHECKS` führt es nicht).
- **Verträge:** `harness/targets/sdk-altserver.md` — Zeile U in der Schritt-Tabelle, Aufruf mit
  `make image` vorweg, Exit-Tabelle (neuer Exit-1-Grund „Ziel-Image nicht geladen“), Host-Werkzeuge um
  `awk` ergänzt (benannt), Grenze 6, Mutationstabelle mit gesehenem Rot. Wahr gegen den gelesenen
  Runner. `harness/mk/sdk.mk`: Kommentar und `##`-Text nennen Phase U und `make image` vorher.
  `harness/README.md`: die zwei Zeilen (`make test-integration`, `make test-sdk-altserver`) tragen den
  Verweis auf Phase U; Suchlauf `sdk-altserver` in der README 1 → 2 gemessen.
- **Planner-Erwartung vs. Messung:** Suchlauf-Zeile „desselben Images“ 5 → 8 statt 5 — im Plan
  begründet (das Muster trifft das Wort „versionswechsel“ in Slice-Name, Kopfkommentar und
  Meldung); ich habe die Zeilen nachgemessen (Suchlauf `OK soll=8 ist=8`), die Begründung trägt. Ein
  Träger außerhalb von `harness/`/`tools/harness/`, der den Upgrade-Tausch als „derselbe Bau“
  beschreibt und neu hinzukam, ist nicht gefunden; bestehende Träger dieser Art (z. B. die
  Beobachtung selbst) sind absichtlich Records.

## 6. Bedingungen, offene Punkte für den Planner

Keine blockierende Bedingung. Zu setzen durch den Planner bei der Closure:

1. **Haken** Liefer-Punkt 3, Gate-Zeilen (`make gates`, `docs-check`, `suchlauf-nachmessen`,
   `fmt-check`/`test`, `kommentar-kennungen`, P10-Klausel, `make test-sdk-altserver`, Doku-Update) am
   Endstand erneut belegen (Dieser Report belegt den Stand `1383d75e`).
2. **§6-Ausgänge (Vorschlag):**
   - *Alt-Build läuft mit konstantem Schema durch den Tausch:* **entfallen** (Erwartung traf zu, zwei
     Läufe Exit 0; Begrenzung: eine Quelle, eine Tabelle, ein Alt-Stand, kein Schema-Wechsel).
   - *Phase trivial grün:* **entfallen**, getragen von M3 (Image-Referenz/-ID) und M1 (Prüfsumme),
     beide in diesem Report mit gesehenem Rot nachgefahren.
   - *Docker-Hub-Abruflimit/ghcr-Pull:* **weiter offen** — kein Limit in drei bis fünf Läufen gesehen
     (nicht widerlegt, Cache-Herkunft nicht gemessen); Anker: das Register
     `kein-echter-versionswechsel-upgrade-test` ist dafür nicht der Ort, Adresse ist der
     bestehende Sensor-Vertrag Grenze 5.
   - *Laufzeit:* **entfallen** — Phase U 7 s (gedruckt, ein Lauf); der ganze Lauf von
     `make test-sdk-altserver` etwa 2,5–3 Minuten (aus Dateizeitstempeln der Läufe, nicht gestoppt).
   - *Container-/Netzkollision:* **weiter offen** — Vorprüfung deckt den Start, nicht den Tausch mitten
     im Lauf (hergeleitet aus dem Runner; nicht gefahren: ein paralleler Runner). Kein Befund.
   - *`:dev` als Vorbedingung:* **eingetreten und getragen** — Vorprüfung mit Meldung „make image
     vorher“ im Code und im Vertrag; ein veraltetes `:dev` bleibt unbenannt-gültiges Ziel (im Vertrag
     benannt). Die negative Vorprüfung (Image fehlt → Exit 1) habe ich **nicht** gefahren
     (*hergeleitet* aus Zeile 113–116 des Runners).
   - *Aufräum-`trap`:* **entfallen**, Messung `docker ps -a`/`docker network ls` ohne `cdc-*` nach
     allen fünf Läufen.
   - *`.dockerignore`/Dockerfile:* **entfallen**, kein Dockerfile im Diff.
3. **Register** — siehe §7.
4. **Zusammenhang Review F-3** (Schlusszeile ohne Mengenvermerk): kein Handlungsbedarf, hier als
   ausreichend begrenzt gewertet (§4). Wer die Zeile als Beleg zitiert, nennt die Begrenzung aus
   Vertrag Grenze 6.
5. **Vorläufige Ablage:** der Slice liegt noch unter `in-progress/`; `git mv` nach `done/` erst nach
   Inhalt der Closure-Notiz (`AGENTS.md` §3.3). Von mir keine DoD-Häkchen gesetzt.

## 7. Register `BEO-PGC/kein-echter-versionswechsel-upgrade-test`

- Evidenz-Datei `evidence/slice-upgrade-versionswechsel-alt-image.md` liegt; Dateien unter
  `evidence/` jetzt drei (`slice-063.md`, `slice-sdk-0-6-kompatibilitaet-messen.md`,
  `slice-upgrade-versionswechsel-alt-image.md`) — Zähler 3×, *abgeleitet* aus den Dateien (nicht
  geschrieben). `state.md` unverändert: „offen — Trigger erfüllt, Träger benannt“ (Ausgang erst nach
  dem gelesenen Lauf, wie im Plan zugesagt).
- **Vorschlag für die Closure** (Eintragsform: Zustandswort aus *verkörpert · geplant · gestrichen* mit
  auflösbarem Anker, Muster `a-check-null-abdeckung`): **Ausgang verkörpert** · seit
  `slice-upgrade-versionswechsel-alt-image` — Phase U in `tools/harness/run-sdk-altserver-tests.sh`
  (Vertrag `harness/targets/sdk-altserver.md`; Träger der Entscheidung
  [`ADR-0148`](../plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Teil 2):
  ein anderer Server-Build (0.5.0) erfasst Zeilen, ein `--force-recreate` auf `:dev` lässt den
  Datenstand der Quelle gleich lesbar und setzt die Erfassung fort (gemessen, zwei Läufe, Exit 0; drei
  Mutationen mit gesehenem Rot). **Begründung:** die Lücke der Beobachtung („alt“ und „neu“ sind
  dasselbe Image) ist für den Fall konstantes Schema geschlossen, und der Sensor steht im Baum.
  **Ehrlich zu vermerken:** der Schemawechsel über Versionen bleibt ungemessen (akzeptiertes Negativ,
  Option C der ADR-Teilfrage; Wiederöffnungs-Trigger: ein beobachteter Betriebs-Upgrade-Fehler, der am
  konstanten Schema vorbeiläuft), also **nicht** „der Versionswechsel ist vollständig belegt“. Wer das
  für zu weit hält, wählt *geplant* mit Anker auf Option C; ich lese *verkörpert* als wahr für den
  gemessenen Umfang, wenn der Eintrag die Begrenzung trägt.
- Steering-Loop-`liegt in`-Feld: erst wenn mit der Closure etwas verkörpert ist; Kandidaten
  `tools/harness/run-sdk-altserver-tests.sh` und `harness/targets/sdk-altserver.md`.

## 8. Hinweise zur Form

- Die Report-Datei war vor `git add` untracked; `make docs-check` wurde nach `git add` dieses Pfads
  erneut gefahren (Exit siehe Commit-Nachricht und Bericht an den Planner).
- Alle Kennungen in diesem Report sind verlinkt oder stehen in Backticks als Pfad/Name.
