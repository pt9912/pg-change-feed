# Welle 12: Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-12-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-13.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`welle-11`s Eröffnungs-Recherche fand real, dass
[`LH-FA-ADM-001`](../../../spec/lastenheft.md) (Lastenheft, Rang 1) SQL-
Funktionen für Aktivierung/Deaktivierung/Status/Consumer-Verwaltung
verlangt, die real nicht existieren; dass
[`LH-FA-CFG-002`](../../../spec/lastenheft.md) (Deaktivierung) keinen
Live-Zugriffsweg hat; und dass
[`LH-FA-SST-003`](../../../spec/lastenheft.md) (CLI) keinen Status-/
Diagnose-Befehl hat. Eine zweite Fork-Recherche fand zusätzlich, dass
[`ADR-0046`](../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
die Frage, wie eine schreibende SQL-Funktion einen Go-Inbound-Port
erreichen soll, bewusst offen gelassen hatte, und dass
`Assembler.tables` ohne Live-Reload-Mechanismus ist. Ein Architect-Verdikt
([`ADR-0050`](../adr/0050-sql-administration-antragsqueue-und-live-reload.md))
hat beide Fragen entschieden: eine Antrags-Queue mit `LISTEN`/`NOTIFY`,
verarbeitet von einer neuen Administrations-Goroutine im laufenden
Capture-Prozess, die den bestehenden Inbound Port aufruft und die
laufende `Assembler`-Bindung live nachträgt; plus ein Boot-Wechsel von
`Assembler.tables` ausschließlich aus `CDC_TABLES` auf `cdc.source_table`
via `TableActivationPort.List`.

Diese Welle setzt `ADR-0050` real um und liefert zusätzlich, unabhängig
davon, die CLI-Diagnose (`LH-FA-SST-003`).

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass ein Administrator eine Tabelle real über
`SELECT cdc.enable_table(...)` aktivieren kann und der **laufende**
Erfassungspfad diese Tabelle danach tatsächlich erfasst, ohne dass der
Prozess neu startet — das ist erst die Summe aller drei Slices.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `welle-11` liegt in `done/`.
- `ADR-0050` liegt `Accepted` vor.
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make gates` grün.
- Ein realer Testlauf gegen den Compose-Stack belegt: `SELECT
  cdc.enable_table(...)` real ausgeführt → die Administrations-Goroutine
  verarbeitet den Antrag → eine danach an der neu aktivierten Tabelle
  ausgeführte Änderung wird vom **laufenden** Prozess real erfasst, ohne
  Neustart.
- Derselbe Nachweis für `cdc.disable_table(...)` (Erfassung endet für die
  Tabelle, ohne den gesamten Prozess zu beenden).
- Ein neuer CLI-Diagnose-Befehl liefert real die in `LH-FA-ADM-002`…`005`
  genannten Signale.
- `BEO-PGC/verwaltung-keine-sql-administration` erreicht Ausgang
  *verkörpert*.
- Closure-Notiz in `welle-12-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-036 | Antrags-Queue und schreibende SQL-Funktionen (`cdc.enable_table`, `cdc.disable_table`) | [`ADR-0050`](../adr/0050-sql-administration-antragsqueue-und-live-reload.md) |
| slice-037 | Administrations-Goroutine, Assembler-Live-Reload und Boot-Wechsel auf `TableActivationPort.List` | [`ADR-0050`](../adr/0050-sql-administration-antragsqueue-und-live-reload.md) |
| slice-038 | CLI-Diagnose-Befehl | [`LH-FA-SST-003`](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: „Retention-Löschausführung" (in der Roadmap direkt
  nachfolgend eingereiht — kein inhaltlicher Zusammenhang, reine
  Sequenz-Priorisierung).
- Wird blockiert von: keine andere Welle.
- Intern: `slice-036` → `slice-037` (die Goroutine verarbeitet Anträge,
  die `slice-036`s SQL-Funktionen erst erzeugen können). `slice-038` ist
  unabhängig von beiden (reiner Lesezugriff, kein Reload-Bezug laut
  `ADR-0050`s Folgepflicht) und kann parallel oder zuerst laufen.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Consumer-Verwaltung über SQL** (`LH-FA-ADM-001`s „Consumer-
  Verwaltung"-Beispiel, `SELECT * FROM cdc.consumers`) — `cdc.consumer_status`
  deckt die Lese-Hälfte bereits ab (`welle-11`); eine schreibende
  Consumer-Registrierung/-ACK über SQL (analog zur hier gebauten
  Antrags-Queue) ist ein eigener, hier ausgeschlossener Vorgang, da
  `LH-FA-ADM-001`s Kern-Antrieb (`LH-FA-CFG-001`/`002`) bereits die
  Welle füllt.
- **Replica-Identity-Prüfung real umsetzen** — `ADR-0050`s Kontext
  benennt sie als „separater, hier nicht zu schließender Befund"; diese
  Welle baut die Antrags-Queue-Architektur so, dass sie später ergänzt
  werden kann, führt sie aber nicht selbst ein.
- **`spec/architecture.md`s Sequenzdiagramm-Korrektur** (`ADR-0050`s
  Folgepflicht: SQL- und CLI-Pfad für `LH-FA-CFG-001.a` trennen) — läuft
  als Teil von `slice-037` (dort entsteht der asynchrone SQL-Pfad real),
  kein eigener Slice.
- **Entfernte-Spalten-Deaktivierungs-Semantik, Retention, Performance,
  NATS** — bereits als eigene, andere Wellen vorgemerkt bzw. abgeschlossen.

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

Ergebnis: <Zeiger auf `welle-12-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
