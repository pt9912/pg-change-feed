# Architect-Review slice-015 — Verdikt zu Trigger-Audit (ADR-0043) und Register-Verkörperung (BEO-PGC/d-migrate-nacharbeit, 3×)

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-11.
**Eingang:** Slice-Plan
[`docs/plan/planning/done/slice-015-d-migrate-1.3.0-retirement.md`](../planning/done/slice-015-d-migrate-1.3.0-retirement.md)
§1–§8 · [`docs/reviews/review-slice-015.md`](../../reviews/review-slice-015.md)
(0 HIGH, F-1 MEDIUM/F-2 LOW disponiert) ·
[`docs/reviews/verify-slice-015.md`](../../reviews/verify-slice-015.md)
(DoD 5/10 real geprüft, kein Defekt, V-1/V-2 offen für Planner) ·
[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (Accepted, permanent,
Re-Evaluierungs-Trigger) ·
`docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`
(`observation.md`, `state.md`, `evidence/{slice-006,slice-010,slice-015}.md`)
· Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register ·
`modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle (Tabelle „ohne
Wellen-Betrieb": Trigger-Audit und Verkörperung bleiben, der
Verifikations-Zug entfällt — dieser Slice ist wellenlos).

**Ausgang:** Zwei unabhängige Architect-Züge, beide mit Verdikt:

1. **Trigger-Audit (ADR-0043):** Der Re-Evaluierungs-Trigger feuert
   **nicht**. ADR-0043 bleibt `permanent`, unverändert, kein Folge-ADR
   fällig — bestätigt den Reviewer- und Verifier-Negativbefund unabhängig.
2. **Register-Verkörperung (BEO-PGC/d-migrate-nacharbeit, 3×):**
   Ausgang **verkörpert** — Regeltext und Zielort unten, Herkunftsanker
   `seit slice-015`. `geplant` scheidet mangels Kennung aus, `gestrichen`
   scheidet mangels Wegfall-Begründung aus (Modul 6 zwingt zu einem der
   drei Ausgänge bei Closure — Details in Zug 2).

**Zusätzlich:** Beide vorgeschlagenen §6-Risiko-Ausgänge werden bestätigt
(Risiko 1 „eingetreten", Risiko 2 „entfallen") — siehe §Risiko-Bestätigung.

**Harte Regel eingehalten:** ADR-0043 (`Accepted`) wird von diesem Lauf
**nicht** inhaltlich geändert — der Trigger-Audit bestätigt sie nur; die in
Zug 2 vorgeschlagene Regel liegt **außerhalb** der ADR-Aussage selbst (eine
operative Test-Kadenz-Regel, keine Änderung an „welches Werkzeug, welche
Bedingung löst Re-Evaluierung aus"), deshalb kein Folge-ADR, kein
`supersedes`. Kein ADR-Index-Eintrag (dieses Dokument ist kein ADR). Weder
der Slice-Plan noch `harness/README.md` noch das Beobachtungs-Register
werden von diesem Lauf editiert — das Schreiben der verkörperten Regel
(„liegt in …") ist Planner-Arbeit im Closure-Zug (Modul 8, Schritt 3b: der
Architect liefert Zielort und Regeltext, der Planner schreibt sie ein).

---

## Zug 1 — Trigger-Audit (ADR-0043)

**Trigger-Text (wörtlich):** „d-migrate kann eine benötigte Operation
nicht ausdrücken — der Lauf blockiert (Exit 8) oder meldet einen Blocker,
**und** es gibt keine Ausweichform (Overlay, `--rename-*`, berichtete
manuelle Nacharbeit) — dann ist der Werkzeugwechsel als Folge-ADR mit
`supersedes` zu prüfen. Sonst `permanent`."

Der Trigger ist eine **Konjunktion**: beide Teilbedingungen müssen
gleichzeitig gelten — *nicht ausdrückbar* **und** *keine Ausweichform*.
Fehlt eine der beiden, feuert der Trigger nicht.

**Stand nach diesem Slice, Fall für Fall:**

| Fall | Ausdrückbar? | Ausweichform vorhanden? | Trigger-Bedingung erfüllt? |
|---|---|---|---|
| `chk_change_operation` (CHECK) | **ja, ab 1.3.0** (real reproduziert: `--dry-run` gegen migrierte DB schlägt keine Operation vor, deklarativ in `schema.yaml`) | entfällt (Ausweichform `nacharbeit-operation-check.sql` bereits zurückgebaut) | **nein** — Fall ist gelöst, nicht mehr Trigger-relevant |
| Drei Views (`active_tables`, `consumer_status`, `changes`) | **nein, weiterhin nicht** (Exit 5, real reproduziert von Implementer **und** unabhängig vom Verifier, zwei getrennte Testcontainer-Instanzen) | **ja** — `nacharbeit-views.sql`, berichtete manuelle Nacharbeit im `schema-rollout`-Target, in Betrieb | **nein** — die zweite Teilbedingung fehlt |

Für keinen der beiden Fälle sind **beide** Teilbedingungen gleichzeitig
erfüllt. Der Trigger feuert nicht.

**Prüfung der Ausweichform-Qualität** (damit „Ausweichform vorhanden" nicht
zur Formalie wird): `nacharbeit-views.sql` ist keine Behauptung, sondern
ein real laufender psql-Nacharbeitsschritt im `schema-rollout`-Target, vom
Implementer und unabhängig vom Verifier gegen frische Rollouts bestätigt
(`make test-integration` grün, alle drei Views angelegt). Das ist exakt die
im Trigger-Text genannte Kategorie „berichtete manuelle Nacharbeit" — keine
Grauzone.

**Zusatzbefund aus `evidence/slice-015.md`, ohne Wirkung auf dieses
Verdikt:** Die isolierte Reproduktion (Repro-Paket an d-migrate übergeben)
zeigt eine plausible Erklärung für die Views-Drift — Post-Compare prüft
laut `spec/cli-spec.md:940-951` eigentlich keine vom Plan selbst
angelegten Objekte, tut es aber bei `CREATE VIEW` doch — und d-migrate hat
einen eigenen internen Bug bestätigt, Fix-Arbeit begonnen, kein
Release-Termin. Das ändert an der Trigger-Bewertung **nichts**: Die
Ausweichform besteht unabhängig davon, ob die Ursache ein bekannter,
bereits in Arbeit befindlicher Bug ist oder eine grundsätzliche
Werkzeug-Grenze. Es ist aber ein starkes zusätzliches Argument **gegen**
einen verfrühten Werkzeugwechsel: Die Grenze ist kein dauerhaftes
Charakteristikum von d-migrate, sondern ein von den Autoren selbst
anerkannter, in Bearbeitung befindlicher Defekt.

**Verdikt Zug 1:** Bestätigt — Reviewer-Negativbefund und
Verifier-Negativbefund sind beide unabhängig zutreffend. ADR-0043 bleibt
`Accepted`, `permanent`, unverändert. Kein Folge-ADR.

---

## Zug 2 — Verkörperung (Beobachtungs-Register 3×)

### Warum überhaupt ein Ausgang fällig ist

`evidence/slice-015.md` ist die dritte Beleg-Datei unter
`BEO-PGC/d-migrate-nacharbeit/evidence/` (nach `slice-006.md`,
`slice-010.md`) — der abgeleitete Zähler steht bei 3×. Modul 6 lässt hier
keinen Spielraum: „Nicht zulässig ist ein Eintrag, der eine Closure ohne
Ausgang übersteht." Da `evidence/slice-015.md` mit **diesem** Slice
entstanden ist, ist der Ausgang eine Bedingung für die Closure **dieses**
Slice, nicht einer späteren.

### Die drei Ausgänge, geprüft der Reihe nach

- **Gestrichen — scheidet aus.** Die Beobachtung kann nicht als „tritt
  nicht mehr auf" gestrichen werden: Die Views-Ausweichform besteht
  unverändert, der d-migrate-Bug ist real und ungelöst (siehe Zug 1). Eine
  Streichung ohne tragende Begründung wäre stilles Vergessen genau des
  Musters, das die Schwelle markieren soll.
- **Geplant — scheidet aus.** Modul 6: „Geplant ist ein Ausgang **mit
  Kennung**, kein Vorsatz." Es gibt aktuell keine Slice- oder Welle-Kennung,
  die „die Regel später schreibt" — der Fix-Termin bei d-migrate ist
  unbekannt („Arbeit begonnen", kein Release-Datum), und die Regel, um die
  es hier geht, ist ohnehin nicht an den externen Fix gebunden (siehe
  nächster Absatz). Ein „geplant" ohne Kennung wäre genau der laut Modul 6
  unzulässige Vorsatz.
- **Verkörpert — trägt.** Das ist keine Verlegenheitslösung, sondern folgt
  aus einer wichtigen Trennung, die die Aufgabenstellung selbst schon
  andeutet: **Zwei verschiedene Dinge** hängen an dieser Beobachtung, und
  nur eines davon liegt außerhalb unserer Kontrolle.
  1. *Die technische Auflösung der Views-Drift* (d-migrate behebt seinen
     Post-Compare-Bug) — das ist extern, terminlich nicht kontrollierbar,
     und **dafür** gäbe es in der Tat keinen sauberen der drei Ausgänge
     (kein „geplant" ohne Kennung, kein „verkörpert" ohne stehende Regel,
     kein „gestrichen" ohne Wegfall). Aber diese Frage muss hier **nicht**
     entschieden werden — der Registereintrag verlangt einen Ausgang für
     die *Beobachtung* (das wiederkehrende Muster), nicht für das
     zugrunde liegende Fremd-Ticket.
  2. *Wie wir mit dem bekannten, unveränderten Zustand umgehen* — das ist
     vollständig in unserer Kontrolle, ab sofort wirksam, ohne auf
     irgendetwas Externes zu warten. Genau das ist der Gegenstand der
     Regel unten, und genau das kann „stehen", ohne dass sich am Views-Bug
     selbst etwas ändert.

  Der dreimal wiederholte Kern des Musters (slice-006, slice-010,
  slice-015) ist nicht „die Views-Drift existiert" allein — das war schon
  nach dem zweiten Auftreten klar und stand unter der Schwelle. Der neue,
  mit slice-015 hinzugekommene Aspekt ist: Ein Pin-Bump-Slice testet den
  bekannten, unveränderten Fall routinemäßig erneut real gegen eine
  Testcontainer-Datenbank — ein teurer Lauf, der bei unverändertem
  Ausweichform-Zustand keine neue Information liefert, solange d-migrate
  selbst nichts Neues meldet. Das ist der Punkt, an dem eine Regel jetzt
  etwas verändert.

### Regeltext (final, zur Übernahme durch den Planner)

> Ein d-migrate-Pin-Bump testet die bekannte Views-Drift
> (`raw-sql-text-drift` bei `CREATE VIEW`, `BEO-PGC/d-migrate-nacharbeit`)
> nicht routinemäßig real gegen einen frischen Testcontainer-Rollout nach.
> Ein erneuter Real-Test dieser Klasse ist erst fällig, wenn d-migrate ein
> explizites Fix-Signal für CREATE-VIEW-Post-Compare-Drift gibt
> (Changelog-Eintrag oder direkte Team-Rückmeldung, die genau diesen Fall
> benennt) — nicht bei jedem Pin-Bump an sich. Der reguläre
> Pin-Bump-Testlauf (`make schema-validate`, `make gates`, `make
> test-integration`) bleibt davon unberührt und läuft wie bisher; entfällt
> soll nur der zusätzliche, gezielte Reproduktionslauf gegen die bereits
> bekannte, unveränderte Views-Grenze. Tritt das Fix-Signal ein, ist der
> Real-Test wieder Pflicht — Ergebnis (gelöst oder weiterhin Exit 5) landet
> als neue `evidence/slice-<NNN>.md` in `BEO-PGC/d-migrate-nacharbeit/`,
> unabhängig vom Ausgang.

### Zielort

`harness/README.md`, §Sensors, Zeile `make schema-rollout`
(Bindung-Spalte, aktuell `kein Gate,
[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md)`) — dieselbe Zeile
trägt bereits den ADR-0043-Bezug und ist damit der Ort, an dem ein Leser
die Testcadence für den d-migrate-Rollout ohnehin nachschlägt. Ergänzung
nach dem Muster der bestehenden `commit-traceability`-Zeile (die dort
bereits einen Trigger-Satz plus `seit slice-006` führt):

> kein Gate, [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md);
> Views-Ausweichform-Retest nur bei explizitem d-migrate-Fix-Signal für
> CREATE-VIEW-Post-Compare-Drift, nicht bei jedem Pin-Bump — seit slice-015
> (`BEO-PGC/d-migrate-nacharbeit`, 3×)

**Herkunftsanker:** `seit slice-015` (`BEO-PGC/d-migrate-nacharbeit`, 3×:
`evidence/slice-006.md`, `evidence/slice-010.md`, `evidence/slice-015.md`).

### Hinweis für die Register-`state.md` (Planner-Arbeit, nicht Teil dieses Verdikts)

Der Registereintrag selbst bleibt nach Modul 6 „mit Vermerk stehen" — sein
`Zustand` wechselt von `offen` auf `verkörpert`, **obwohl** die
zugrundeliegende Views-Drift technisch unverändert real und offen bleibt.
Das ist kein Widerspruch: „verkörpert" bewertet, ob aus der dreifachen
Wiederholung eine stehende Regel wurde — nicht, ob der externe Bug
behoben ist. Der Eintrag bleibt aktiv und nimmt weitere `evidence/`-Dateien
auf, sobald der genannte Fix-Signal-Trigger einen neuen Real-Test auslöst.

**Verdikt Zug 2:** `verkörpert`. Regeltext und Zielort wie oben. Die
Ausführung (Eintrag in `harness/README.md`, `state.md`-Update) ist
Planner-Arbeit im Closure-Zug dieses Slice.

---

## Risiko-Bestätigung (§6 des Slice-Plans)

Beide vom Implementer/Verifier vorgeschlagenen Ausgänge werden bestätigt —
unabhängig anhand der vorliegenden Belege geprüft, nicht aus dem
Verifier-Vorschlag übernommen:

- **Risiko 1** („Der Sandbox-Modus löst möglicherweise nur EINEN der
  beiden Fälle, nicht beide") → **eingetreten**. Exakter Deckungsgleich:
  CHECK real gelöst (Exit 0, deklarativ), Views real weiterhin offen
  (Exit 5) — von Implementer und Verifier unabhängig reproduziert. Ein
  Teil-Retirement ist laut Risikoformulierung selbst „ein legitimes
  Ergebnis, kein Scheitern" — dem schließt sich dieses Verdikt an. Kein
  Carveout, kein rotes Gate: Der Closure-Trigger des Slice-Plans (§5)
  verlangt kein vollständiges Retirement, sondern „real belegt (nicht nur
  behauptet)" — das ist für beide Teilergebnisse der Fall.
- **Risiko 2** („Der Digest-Bump könnte einen Regressions-Fund gegen den
  bestehenden `schema.yaml`-Bestand auslösen") → **entfallen**, mit
  Begründung: Kein Regressions-Fund in irgendeinem real gefahrenen Sensor
  — weder `make gates` noch `make test-integration` (Implementer- **und**
  unabhängiger Verifier-Lauf) noch in den gezielten Zusatz-Proben des
  Verifiers (u. a. testweise Deklaration der drei Views im `views:`-Knoten,
  um eine mögliche Nebenwirkung zu provozieren — Ergebnis war der bekannte,
  erwartete Exit-5-Fall, keine neue Drift an bestehenden Objekten). Die im
  Risiko benannte Changelog-Prämisse („View Portability Knowledge Shift")
  hat sich nicht in einer Regression gegen den bestehenden Bestand
  niedergeschlagen — mehrfach und unabhängig real geprüft, nicht nur
  einmal beobachtet.

Beide Ausgänge sind bei Closure so in §7/§6 des Slice-Plans zu übernehmen.

---

## Disposition (Gesamt)

| Frage | Ergebnis |
|---|---|
| Feuert der ADR-0043-Re-Evaluierungs-Trigger? | Nein — Konjunktion nicht erfüllt (CHECK: gelöst; Views: Ausweichform vorhanden) |
| ADR-0043-Status | unverändert, `Accepted`, `permanent` |
| Folge-ADR nötig? | Nein |
| Register-Ausgang `BEO-PGC/d-migrate-nacharbeit` bei 3× | **verkörpert** |
| Regeltext | siehe §Regeltext oben |
| Zielort | `harness/README.md` §Sensors, `make schema-rollout`-Zeile |
| Herkunftsanker | `seit slice-015` |
| Risiko 1 (§6) | bestätigt: **eingetreten** |
| Risiko 2 (§6) | bestätigt: **entfallen**, mit Begründung |
| Carveout nötig? | Nein — kein rotes Gate, keine Ausnahme, kein Trigger-Bezug |

---

## Beleg-Anker (Kurzfassung)

| Aussage | Beleg |
|---|---|
| CHECK jetzt ausdrückbar, deklarativ konvergent | `verify-slice-015.md` Sensor-Tabelle Zeile „Eigener frischer Rollout … `--dry-run`" |
| Views weiterhin nicht ausdrückbar, Ausweichform aktiv und real getestet | `verify-slice-015.md` Sensor-Tabelle Zeile „Eigener Rollout mit … Views … testweise deklariert"; `tools/schema/nacharbeit-views.sql`; `make test-integration` (beide Läufe) |
| Trigger-Wortlaut (Konjunktion) | [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) §Re-Evaluierungs-Trigger |
| d-migrate-Bug bestätigt, Fix in Arbeit, kein Termin | `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/evidence/slice-015.md` |
| Zähler bei 3× | `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/evidence/{slice-006,slice-010,slice-015}.md` |
| Register-Regel „geplant nur mit Kennung" | Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register |
| Zielort-Präzedenz (Sensors-Zeile trägt Trigger + `seit slice-NNN`) | `harness/README.md` §Sensors, `make commit-traceability`-Zeile |
| Risiko 1 Deckung | Commit-Botschaft `3cb0c8e` („Retirement damit teilweise"); `verify-slice-015.md` DoD-Punkt 9 |
| Risiko 2 Deckung | `verify-slice-015.md` Sensor-Tabelle (`make gates`, `make test-integration`, Zusatzproben), Negativbefunde |
