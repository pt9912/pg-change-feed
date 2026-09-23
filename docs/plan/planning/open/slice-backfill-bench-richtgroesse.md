# Slice backfill-bench-richtgroesse: Bench der Kopierdauer je Tabellengröße — Warn-Richtgröße aus der Messung, Warnung in `diagnose`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Out-of-Scope: Parallelisierung/Durchsatz über sehr
große Tabellen ist Ausbaustufe — hier die Messung, ab wann eine Tabelle dafür
in Frage kommt), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 3 (große Tabellen: Warnung, keine
Ablehnung; Richtgröße zu messen), Teilfrage 4 (Ein-Transaktions-Form),
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b) (Benchmark-Infrastruktur), [`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md) (Benchmark-Schwellen).

**Berührte Spec-Stellen:** [`SPEC-025`](../../../../spec/pflichtenheft.md) (Bench-Schwellen — gelesen; **kein** neuer
Wert in §3 ohne schärfende ADR, siehe §6), [`SPEC-029`](../../../../spec/pflichtenheft.md) (Run-Zustand mit
`estimated_rows`) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Aus einer Messung entsteht die **Warn-Richtgröße** für große
Tabellen — und aus ihr eine Warnung, nie eine Ablehnung.

- **Bench.** Ein viertes eigenständiges Bench-Skript im Muster von
  `tools/bench-source-impact.sh`, `tools/bench-scaling.sh` und
  `tools/bench-batch-vs-single.sh` (Arbeitsname `tools/bench-backfill.sh`, dieselbe
  Umgebung über `tools/bench-lib.sh`): Kopierdauer eines Runs je Tabellengröße
  (Antrag bis `completed`), an einer Stufenfolge von Tabellengrößen (die Stufen
  legt der Slice fest und nennt sie im Bericht); gemessen wird auch, wie nah die
  **geschätzte** Zeilenzahl (`pg_class.reltuples`) der tatsächlichen kommt — in
  einer frisch befüllten, in einer analysierten und in einer nie analysierten
  Tabelle. `make bench` fährt das Skript mit; es ist eine **Messung** mit
  gedruckter Zahl, kein Pass/Fail gegen eine Schwelle (die Richtgröße ist ihr
  Ergebnis, nicht ihre Vorbedingung).
- **Richtgröße und Warnung.** Aus der Messung und dem im Start-Trigger
  festgelegten Warn-Kriterium folgt eine Zeilenzahl-Richtgröße. Sie steht als
  benannte Konstante im Code (Anker: der gemessene Lauf) und im Handbuch
  (Abschnitt „Grenzwerte“, mit ihrem Ursprung); `diagnose` gibt eine Warnzeile
  aus, wenn die geschätzte Zeilenzahl eines Runs die Richtgröße übersteigt;
  eine unbekannte Schätzung warnt nicht, sagt aber „unbekannt". Die
  Administrations-Goroutine protokolliert dieselbe Warnung beim Antrag. **Keine
  Ablehnung** ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Pass/Fail-Schwellenwert im Bench** — [`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md) regelt die
  Schwellen der drei bestehenden Bench-Kennungen; die Richtgröße ist eine
  gemessene Orientierung, keine Anforderung. Ein Schwellenwert im Bench wäre eine
  neue Gate-Aussage ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Ein Eintrag der Richtgröße in `spec/pflichtenheft.md` §3** — die Regel dieses
  Abschnitts verlangt, dass eine ADR den Wert in ihrem `Schärft:`-Feld
  deklariert; [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) tut das nicht (`Accepted`, unveränderlich). Soll die
  Richtgröße Vertrag werden, braucht sie eine schärfende ADR (Auftraggeber-Frage,
  siehe Bericht der Welle-Eröffnung).
- **Checkpoint, Block-Transaktionen, Parallelisierung** — der
  Re-Evaluierungs-Trigger für sie ist „Kopierdauer über der Betriebs-Toleranz"
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)); diese Messung liefert dafür Zahlen, entscheidet ihn nicht.
- **Eine Warnung am Status-View** — die View kennt die Konstante nicht; sie trägt
  die Rohwerte (Schätzung, Laufzeit).
- **Aussagen über andere Hardware** — die Zahl gilt für den gemessenen Host und
  nennt ihn.

## 2. Definition of Done

- [ ] Bench: `tools/bench-backfill.sh` (Arbeitsname) läuft in `make bench` mit
      und druckt die Kopierdauer je Tabellengröße und die Abweichung der
      Schätzung von der tatsächlichen Zeilenzahl; jede gedruckte Zahl trägt
      Größe und Lauf. *Zu belegen durch:* ein realer `make bench`-Lauf (Exit
      ungefiltert gesichert), Ausgabe im Bericht; ein gemessener Lauf, nicht
      abgeleitet.
- [ ] Richtgröße und Warnung: die Richtgröße folgt aus der Messung und dem
      festgelegten Warn-Kriterium (Herleitung im Bericht, jede Zahl mit
      Ursprung: gemessen · übernommen · abgeleitet); sie steht als benannte
      Konstante mit Anker im Code und im Handbuch-Abschnitt „Grenzwerte“;
      `diagnose` warnt oberhalb der Richtgröße, schweigt darunter und bei
      „unbekannt", und keine Stelle lehnt einen Antrag wegen der Größe ab.
      *Zu belegen durch:* Unit-Tests der Warn-Auswertung an der Grenze (genau
      auf, eins darüber, unbekannt) mit je einer Mutation (`make test`).
- [ ] Träger: die `make bench`-Beschreibung ist auf vier Skripte gezogen
      (Makefile-Kommentar und -Hilfetext, `harness/README.md` §Sensors,
      ggf. `docs/user/bench-abdeckung.md`); das Handbuch trägt die Änderungshistorie.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: Handbuch „Grenzwerte“ mit Richtgröße und Ursprung, Änderungshistorie eine Zeile.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/bench-backfill.sh` (Arbeitsname) | neu | die Messung im Muster der drei bestehenden Skripte; nutzt `tools/bench-lib.sh`. |
| `Makefile` (Target `bench`) | update | ruft das vierte Skript; Hilfetext und Kommentar zählen neu. |
| `internal/bootstrap/wiring.go` (`Diagnose`, Antrags-Zweig) | update | Warn-Auswertung an einer Stelle (Arbeitsname `backfillWarnRows`), Warnzeile in `diagnose`, Log-Warnung beim Antrag. |
| `internal/bootstrap/diagnose_test.go` u. a. | update | Grenzfälle der Warnung. |
| `docs/user/benutzerhandbuch.md` | update | Richtgröße im Abschnitt „Grenzwerte“ mit Ursprung; Diagnose-Beispiel; Änderungshistorie. |
| `harness/README.md` §Sensors | update | Zeile `make bench` (drei → vier Skripte) — aus dem realen Lauf geschrieben. |
| `docs/user/bench-abdeckung.md` | prüfen | die Datei bindet `LH-QA-PER-001`…`003` an durchgesetzte Schwellen; eine Zeile ohne Schwelle für `LH-FA-CAP-009` passt nicht zu dieser Form — Entscheidung im Slice, im Bericht begründet. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Zahl und der Inhalt der Bench-Skripte hinter `make bench`"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Zeile `make bench` in `harness/README.md` („drei eigenständige Bench-Skripte") | `grep -rn 'drei' harness/README.md docs/user/bench-abdeckung.md Makefile` (Fundstellen nach Bench-Bezug lesen) | *(Implementer trägt ein)* | nachziehen |
| Hilfetext und Kommentar des Targets `bench` | Lesen von `Makefile` | *(Implementer trägt ein)* | nachziehen |
| Bench-Abdeckungs-Träger | Lesen von `docs/user/bench-abdeckung.md` und dem Generator in `tools/bench-batch-vs-single.sh` | *(Implementer trägt ein)* | Entscheidung siehe §3, Zeile `bench-abdeckung.md` |
| Handbuch-Grenzwerte | Lesen von §Grenzwerte in `docs/user/benutzerhandbuch.md` | *(Implementer trägt ein)* | Richtgröße ergänzen; die Replication-Verbindungs-Zahl ist Sache von `sql-administration` |
| Build-Kontext | Lesen von `.dockerignore` und `Dockerfile` | *(Implementer trägt ein)* | Bench-Skripte laufen auf dem Host; keine Docker-Stufe kopiert `tools/bench-*.sh` — sonst Freigabe nach [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md) |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `e2e` in `done/` liegt, kein
anderer Slice in `in-progress/` liegt **und** das **Warn-Kriterium festgelegt
ist**: eine Dauer (in Sekunden) der Kopie, ab der eine Tabelle als „groß" gilt,
oder ein anderes benanntes Maß — [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt „die Betriebs-Toleranz",
beziffert sie aber nicht. Die Festlegung trifft der Auftraggeber oder der
Architect (Verdikt unter `docs/reviews/`); ohne sie stünde ein selbst gewählter
Wert als Grenze im Code (`BEO-PGC/geschaetzter-wert-als-grenze`). Das ist eine
**Vorab-Bedingung** und steht am Start, nicht als Rückführung
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Bench und
  Warn-Auswertung nicht in einem Review tragen — der abtrennbare Teil ist die
  Warn-Auswertung in `diagnose`.
- `in-progress` → `open` (blockiert): falls die Schätzung (`reltuples`) so weit
  von der tatsächlichen Zeilenzahl abweicht, dass eine Warnung an ihr nichts
  trägt (dann Architect-Frage: statt der Schätzung eine andere Größe, etwa ein
  Seitenzähler der Tabelle).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer `make bench`-Lauf + `make test`
grün + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Das Warn-Kriterium fehlt** — [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) beziffert die „Betriebs-Toleranz"
  nicht; der Start-Trigger trägt die Vorab-Bedingung. *Erwartet, zu belegen
  durch:* das Verdikt. **Ausgang:** *(bei Closure)*
- **Die Richtgröße wird zur Grenze** (`BEO-PGC/geschaetzter-wert-als-grenze`,
  offen, 1×; dieser Slice ist der Ort, an dem der Fehler passiert): eine
  gemessene Zahl eines Hosts steht im Code und im Handbuch als allgemeine
  Grenze. *Erwartet, zu belegen durch:* jede Nennung trägt „gemessen auf …,
  Lauf …" und das Wort „Richtgröße"; keine Stelle lehnt ab. **Ausgang:** *(bei
  Closure)*
- **Die Schätzung ist so alt wie das letzte `ANALYZE`.** Eine frisch geladene,
  nie analysierte Tabelle trägt „unbekannt" (erwartet), die Warnung schweigt dann
  gerade bei den Tabellen, für die sie gedacht ist. *Erwartet, zu belegen durch:*
  die Schätz-Abweichungs-Messung; das Handbuch nennt den Zusammenhang.
  **Ausgang:** *(bei Closure)*
- **Die Messung hängt am Host** (Docker-Umgebung, Plattengeschwindigkeit). Das
  Handbuch nennt die Zahl als Orientierung mit ihrem Lauf; eine zweite
  Umgebung ist nicht Teil. **Ausgang:** *(bei Closure)*
- **Laufzeit von `make bench`** wächst; der Default-Modus fährt verkürzte Stufen
  (Muster `bench-scaling.sh`), ein Volllauf ist Option. *Erwartet, zu belegen
  durch:* Laufzeit im Bericht. **Ausgang:** *(bei Closure)*
- **Ein Wert in `spec/pflichtenheft.md` §3 ohne schärfende ADR** wäre gegen die
  Regel des Abschnitts; der Slice führt keinen. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); das Bench-Werkzeug unter `tools/` und die Diagnose im
Bootstrap sind keine eigene Sub-Area — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/geschaetzter-wert-als-grenze` (offen, 1×, einschlägig — Start-Trigger und
Risiko §6), `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×,
Start-Trigger), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×, jede Zahl mit Ursprung), `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(verkörpert, 8×, die `make bench`-Zeile beschreibt nur, was der Lauf getan hat),
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (offen, 2×) und
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×) —
Suchlauf-Zeile Build-Kontext, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, 26×, Suchlauf §3), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×, Warn-Tests mit Mutation).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
