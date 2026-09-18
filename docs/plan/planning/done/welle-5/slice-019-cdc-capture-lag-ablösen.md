# Slice slice-019: `cdc_capture_lag` ablösen — Lasttest-Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-5`](welle-5.md) — dieser Slice liefert den
Closure-Trigger der Welle (Ende-zu-Ende-Lasttest-Beleg).

**Bezug:** [`LH-FA-ADM-004`](../../../../spec/lastenheft.md)

**Berührte Spec-Stellen:** [`SPEC-013`](../../../../spec/pflichtenheft.md)
(`CDC_LAG_THRESHOLDS`, Metrik `cdc_capture_lag`)

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Näherung `cdc_capture_lag_approx` durch den kanonischen
`cdc_capture_lag` ([`SPEC-013`](../../../../spec/pflichtenheft.md))
ablösen und real — per Ende-zu-Ende-Lasttest mit künstlich eingebauter
Verzögerung zwischen Quell-Commit und Verarbeitung — belegen, dass die
Metrik den tatsächlichen Abstand abbildet statt nur die
Persistenz-Zeit-Differenz von vorher.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- [`SPEC-013`](../../../../spec/pflichtenheft.md)-Schwellen (p95 ≤ 1 s, Warnschwelle > 5 s, Fehlerschwelle
  > 60 s) als aktive Alarmierung verdrahten — anderer Vorgang: [`SPEC-013`](../../../../spec/pflichtenheft.md)
  definiert die Schwellen als Messziel, nicht als Pflicht zur
  Alarm-Automatisierung; diese Welle liefert die Messung, nicht die
  Reaktion darauf.
- `cdc_capture_lag_approx` rückwirkend aus historischen Daten neu
  berechnen — Bestand bleibt bewusst stehen: bereits gespeicherte
  Transaktionen behalten ihren zum damaligen Zeitpunkt gültigen
  `committed_at`-Wert; eine Rückrechnung wäre eine Datenmigration, kein
  Observability-Wechsel.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `tools/schema/nacharbeit-observability.sql`: `cdc_capture_lag_approx`
      durch `cdc_capture_lag` ersetzt (Metrik-Name, Kommentar); die
      übrigen Metriken unverändert.
- [x] Ende-zu-Ende-Lasttest: eine Testtransaktion mit künstlich
      eingebauter Verzögerung zwischen Quell-Commit und Verarbeitung
      zeigt, dass `cdc_capture_lag` diese Verzögerung real abbildet —
      nicht nur die Persistenz-Zeit-Differenz. Real gegen die Compose-
      Umgebung belegt (`make test-integration`, `tools/harness/
      run-integration-tests.sh`): Feed-Container über die Freezer-Cgroup
      pausiert (`docker pause`), Quell-Commit während der Pause,
      Fortsetzen — gemessen Baseline (ungehindert) ~0,1–0,13 s, verzögert
      (1 s künstliche Pause) ~1,15–1,16 s, stabil über mehrere Läufe.
      Fund unterwegs: der geladene `ghcr.io/pt9912/pg-change-feed:dev`
      war seit vor slice-017/018 nicht neu gebaut — `make image` lief vor
      dem Beleg (`harness/image-hash.txt` fortgeschrieben), sonst hätte
      der Lasttest gegen die alte Persistenz-Zeit-Näherung gemessen statt
      gegen den realen Fix.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`review-slice-019.md`](../../../../docs/reviews/review-slice-019.md)
      (F-1 MEDIUM/F-2 LOW disponiert, 0 HIGH), Verifier bestätigt in
      [`verify-slice-019.md`](../../../../docs/reviews/verify-slice-019.md).
- [x] Doku-Update falls öffentlicher Vertrag berührt — geprüft:
      [`SPEC-013`](../../../../spec/pflichtenheft.md) und der
      Metriken-Katalog-Eintrag (`SPEC-009`) nennen die Metrik bereits
      unter ihrem kanonischen Namen `cdc_capture_lag` als Zielzustand,
      unabhängig vom Umsetzungsstand — kein Textinhalt behauptet dort
      etwas, das dieser Slice widerlegen würde. `spec/pflichtenheft.md`
      bleibt unverändert. Zusätzlich im ersten Lauf gefunden und
      korrigiert: `docs/user/benutzerhandbuch.md` nannte die Metrik noch
      unter dem `_approx`-Namen mit veralteter Beschreibung — aktualisiert
      auf `cdc_capture_lag` mit der realen Bedeutung.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Geprüft: Datei existiert nicht (GF-Repo) — entfällt.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — `BEO-PGC/cdc-capture-lag-real` bekommt den Auflösungs-Beleg dieser Welle (Ausgang `eingetreten`, Träger dieser Slice). Beleg: `evidence/slice-019.md`, `state.md` auf `eingetreten` gesetzt.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Dieser Slice gehört zu `welle-5` — die Paarungen prüft die **Welle-Closure**, nicht dieser Slice. Item entfällt hier bewusst.

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/nacharbeit-observability.sql` | update | Metrik-Zeile `cdc_capture_lag_approx` → `cdc_capture_lag`, Kommentar aktualisiert |
| `tools/harness/run-integration-tests.sh` | update | Lasttest-Block: Feed-Container pausieren, Quell-Commit während der Pause, `cdc_capture_lag` nach dem Fortsetzen gegen die künstliche Verzögerung prüfen |
| `docs/user/benutzerhandbuch.md` | update | Metriken-Liste nennt `cdc_capture_lag_approx` mit veralteter Beschreibung („Näherung über die Persistenzzeit") — öffentlicher Vertrag (Source Precedence Rang 6), im ersten Lauf gefunden, nicht im ursprünglichen Plan |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): kein anderes Slice in `in-progress/`
(WIP-Limit 1); slice-018 liegt in `done/` (der reale Wert steht in
`cdc.transaction.committed_at`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Falls der
  Lasttest-Aufbau (künstliche Verzögerung real reproduzierbar einbauen)
  eine größere Erweiterung der Testinfrastruktur braucht, als ein Slice
  trägt — Rückzug mit Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Falls der reale Test
  zeigt, dass `cdc_capture_lag` trotz slice-017/018 nicht die erwartete
  Verzögerung abbildet (z. B. weil eine weitere, bisher unentdeckte
  Zeitquelle dazwischenliegt) — Blocker, Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss ohne
  offenes HIGH-Finding; der Lasttest-Beleg liegt real vor — das ist
  zugleich der Closure-Trigger von `welle-5`.
- Lerneintrag §7: `BEO-PGC/cdc-capture-lag-real` löst sich auf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Lasttest könnte zeigen, dass `cdc_capture_lag` zwar die
  Persistenz-Verzögerung, aber nicht die reine Netzwerk-/
  Replication-Slot-Latenz abbildet (verschiedene Teilstrecken derselben
  Gesamtverzögerung). **Ausgang: entfallen.** Die künstliche `docker
  pause`-Verzögerung liegt zwischen Quell-Commit und Verarbeitung — vor
  der Persistenz —, und die Messung bildet genau diese Verzögerung ab
  (Faktor 8–9× gegenüber der Baseline, real reproduziert von
  Implementer, Reviewer und Verifier unabhängig). Der Beleg deckt damit
  die volle Commit-zu-Verfügbarkeit-Strecke, nicht nur einen Teilabschnitt.
- Bestehende Dashboards/Beobachter, die noch `cdc_capture_lag_approx`
  erwarten, sehen den Namen nach dieser Umbenennung nicht mehr.
  **Ausgang: entfallen.** Kein externer Konsument dieser Metrik ist in
  diesem Repo bekannt oder dokumentiert.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Der Implementer fand während der Arbeit
  selbst heraus, dass der erste Lasttest-Lauf fälschlich nahezu keine
  Verzögerung maß — statt das als Testfehler zu verwerfen, isolierte er
  die Ursache bis zu einem veralteten Container-Image (gebaut vor
  slice-017/018) und baute es neu, bevor er den Beleg als erbracht
  meldete. Genau das ist der Zweck eines Ende-zu-Ende-Belegs: er hat
  einen realen Bereitstellungs-Fehler gefangen, den keine der drei
  Einzel-Slice-DoDs allein hätte aufdecken können (`welle-5` §1, „Das
  *Mehr* gegenüber den einzelnen Slice-DoDs"). Reviewer und Verifier
  haben den Lasttest und die neue `LAG_DELAY_SECONDS`-Fehlerprüfung
  jeweils unabhängig real reproduziert.
- **Was ging anders als geplant:** Zwei Punkte lagen nicht im
  ursprünglichen §3-Plan und wurden per Plan-Nachzug ergänzt: die
  Korrektur von `docs/user/benutzerhandbuch.md` (öffentlicher Vertrag,
  Rang 6) und der Image-Rebuild (`harness/image-hash.txt`, ausgelöst
  durch den oben beschriebenen Fund).
- **Steering-Loop-Eintrag:** Mit diesem Slice wurde nichts verkörpert.
  Der Eintrag ist gezählt, nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-019.md`
  in `BEO-PGC/cdc-capture-lag-real/` ergänzt — Zähler steht bei 2×,
  Ausgang `eingetreten` (Träger `welle-5`, siehe DoD-Beleg oben).
- **Folge-Slices:** keine.
- **Risiken aus §6:** beide `entfallen` (siehe §6).
- **Drei Paarungen:** entfällt hier — dieser Slice gehört zu `welle-5`;
  die Paarungen prüft die anschließende Welle-Closure für alle drei
  Slices gemeinsam.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührte Sub-Area
(Observability/Schemamigration) ist ein Segment des GF-Baums der
Modus-Deklaration (`PGC` Greenfield, Doc führt) — erfüllt die Schwelle
≥ 2 von 3 Achsen. Nicht zu grob.

**Vorgelagert — offene Beobachtungen sichten:** Register gelesen (zwölf
Einträge — `BEO-PGC/dod-checkbox-nachzug` kam während der Closure von
slice-017 hinzu, siehe slice-018 §8): `cdc-capture-lag-real` 1× — dieser
Slice ist der vorgesehene Auflösungs-Träger (Ausgang `eingetreten`,
siehe §2 DoD). Übrige elf ohne Bezug — siehe slice-017/018 §8 für die
vollständige Liste. Kein Eintrag erreicht mit diesem Slice neu 3× —
keine Lücke.

**Modus-Begründungsblock — Umfang.** Reiner GF-Hinweis genügt (siehe oben);
kein Sub-Area-Block.
