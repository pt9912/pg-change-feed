# Welle 11: E2E-Abdeckung — Verwaltung & Observability

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-11-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-12.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Der Compose-Integrationstest prüft CDC-Status/-Liste
([`LH-FA-CFG-003`](../../../spec/lastenheft.md)/`004`) und
Administration/Observability
([`LH-FA-ADM-002`](../../../spec/lastenheft.md)…`005`) bislang nur
teilweise oder ausschließlich white-box über interne Go-Use-Cases
(`status.NewGetStatusService` u. a.), obwohl die externen SQL-Sichten
(`cdc.active_tables`, `cdc.consumer_status`, `cdc.process_heartbeat`)
dafür bereits real existieren und ausgeliefert werden — dieselbe Lücke
wie zuvor bei `cdc.changes` (`welle-9`, `BEO-PGC/lese-doppelquelle`).
Diese Welle liefert die fehlende Black-Box-Abdeckung gegen die bereits
existierenden Views, **ohne** neue Fähigkeiten zu bauen: die während der
Eröffnungs-Recherche gefundene echte Fähigkeitslücke (keine SQL-
Administrationsfunktionen, keine CDC-Deaktivierung, keine CLI-
Diagnose — registriert als `BEO-PGC/verwaltung-keine-sql-administration`)
ist ausdrücklich **nicht** Gegenstand dieser Welle, sondern einer eigenen,
bereits vorgemerkten Feature-Welle.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass CDC-Status/-Liste und die vier Observability-Signale
zusammen real über die extern ausgelieferten SQL-Sichten lesbar sind und
mit dem internen Go-Lesepfad übereinstimmen — analog zu `welle-9`s
Closure-Trigger für `cdc.changes`.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `welle-10` liegt in `done/`.
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make gates` grün.
- Ein neuer Black-Box-Testfall liest `cdc.active_tables` real über SQL und
  belegt `LH-FA-CFG-003`/`004`s Happy-Path- und Boundary-Kriterien
  (aktivierte Tabelle meldet „aktiviert", nie aktivierte Tabelle meldet
  „nicht aktiviert"/leere Liste).
- Mindestens ein neuer Black-Box-Testfall je Observability-Signal
  (`LH-FA-ADM-002`…`005`) liest über die jeweils zuständige externe
  Schnittstelle (`cdc.process_heartbeat`, `cdc.consumer_status`) — soweit
  nicht bereits real black-box belegt (Fehlerzustand/`LH-FA-ADM-003`
  bereits über `slice-033`s `awaitHeartbeatErrorClass`-Muster gedeckt,
  CDC-Abstand/`LH-FA-ADM-004` bereits über den `cdc_capture_lag`-
  Lasttest-Beleg gedeckt — dieser Slice konsolidiert/verifiziert, statt
  neu zu erfinden, wo bereits real etwas existiert).
- Closure-Notiz in `welle-11-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-034 | Black-Box-Status/Liste gegen `cdc.active_tables` | [`LH-FA-CFG-003`](../../../spec/lastenheft.md)/`004` |
| slice-035 | Black-Box-Observability-Konsolidierung (Betriebsstatus, Fehlerzustand, CDC-Abstand, Verarbeitungsrückstand) | [`LH-FA-ADM-002`](../../../spec/lastenheft.md)…`005` |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keine andere Welle.
- Intern: `slice-034` und `slice-035` sind unabhängig voneinander
  (unterschiedliche Testfälle im selben Compose-Integrationstest); keine
  erzwungene Reihenfolge.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **SQL-Administrationsfunktionen** (`LH-FA-ADM-001`: `cdc.enable_table(...)`,
  `cdc.disable_table(...)`) — real fehlende Fähigkeit, kein
  Testabdeckungs-Defizit; eigene, bereits vorgemerkte Feature-Welle
  „Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose".
- **CDC-Deaktivierung mit Live-Zugriffsweg** (`LH-FA-CFG-002`) — real
  fehlende Fähigkeit (`disable.NewDisableTableService` nirgends
  verdrahtet), dieselbe Feature-Welle wie oben.
- **CLI-Diagnose-Befehl** (`LH-FA-SST-003`) — real fehlende Fähigkeit
  (nur `register-consumer`/`acknowledge-consumer` existieren), dieselbe
  Feature-Welle wie oben.
- **`BEO-PGC/rollen-test-abdeckungsluecken`** — inhaltlich nicht berührt
  (Fork-Recherche 2026-09-12 widerlegt den ursprünglichen Roadmap-Bezug);
  bleibt eigenständig im Register.
- **Retention, Performance, Schema-Evolution-E2E** — bereits als eigene,
  andere Wellen vorgemerkt bzw. bereits abgeschlossen.

## 7. Closure-Notiz

<!--
BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md §Verwendung,
Schritt 5) und darf deshalb nichts Tragendes halten.

- Erst nach Welle-Abschluss fuellen; nur die Nummer, nicht die volle Welle-ID.
- Ziel-Form der Ergebnis-Notiz: `welle-results.template.md` — Schwester-Vorlage
  im Template-Verzeichnis, kein Artefakt deines Repos. Sie ist von der
  Ruheort-Regel ausgenommen und faellt mit diesem Kommentar ohnehin weg.
-->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [welle-11-results.md](welle-11-results.md)
Zähler: [../observations/](../observations/)
