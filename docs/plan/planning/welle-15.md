# Welle 15: NATS-Change-Notification

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-15-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-13.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`pg-change-feed` erfüllt [LH-FA-SST-007](../../../spec/lastenheft.md)
(NATS-Change-Notification) aktuell **gar nicht** — komplett grüne Wiese,
kein NATS-Bezug irgendwo im Repo. [ADR-0055](../adr/0055-nats-change-notification-wecksignal.md)
hat vorab entschieden: Core NATS als reines, unpersistiertes Wecksignal
(leerer Payload, Subjekt `cdc.changes.<source_id>`), `transient`-
Fehlerklasse, best-effort NACH `ACK Source` in `CaptureService.Capture()`,
vollständig optional über `CDC_NATS_URL`.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass `LH-FA-SST-007` vollständig erfüllt ist — erst die Summe aus
Port/Adapter (Grundgerüst), Compose-Verdrahtung (Happy Path), Boundary-
Beleg (nicht verbundener Consumer verliert nichts) und Negative-Beleg
(Reconnect-Nachholen über den bestehenden Zugriffsweg) belegt alle drei
Akzeptanzkriterien gemeinsam.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `slice-051` liegt in `done/`.
- [ADR-0055](../adr/0055-nats-change-notification-wecksignal.md) liegt vor
  (Accepted).
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make gates` grün (unverändert — NATS ist kein Gate).
- Ein realer `make test-integration`-Lauf mit gesetztem `CDC_NATS_URL`
  belegt: ein abonnierter Consumer erhält real ein Wecksignal bei einer
  neuen Change (Happy Path).
- Ein realer Beleg zeigt: eine Change, die auftritt, während kein
  Consumer verbunden ist, bleibt über `cdc.changes` vollständig lesbar
  (Boundary).
- Ein realer Beleg zeigt: nach einem simulierten NATS-Verbindungsabbruch
  und Wiederverbindung holt ein Consumer ausschließlich über den
  bestehenden Zugriffsweg nach, nicht über NATS selbst (Negative).
- Closure-Notiz in `welle-15-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-052 | `ChangeNotificationPort` und `natsnotify`-Adapter | [ADR-0055](../adr/0055-nats-change-notification-wecksignal.md) |
| slice-053 | Compose-Verdrahtung und Happy-Path-Beleg | [LH-FA-SST-007](../../../spec/lastenheft.md) |
| slice-054 | Boundary-Beleg — nicht verbundener Consumer | [LH-FA-SST-007](../../../spec/lastenheft.md) |
| slice-055 | Negative-Beleg — Reconnect-Nachholen | [LH-FA-SST-007](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle (letzte Zeile der Roadmap).
- Wird blockiert von: keine andere Welle.
- Intern: `slice-052` (Port/Adapter/`CaptureService`-Erweiterung) muss vor
  `slice-053` (Compose-Verdrahtung, braucht den fertigen Adapter) laufen.
  `slice-053` muss vor `slice-054`/`slice-055` laufen (beide brauchen die
  laufende NATS-Verdrahtung aus `slice-053`, um sie gezielt zu stören).
  `slice-054` und `slice-055` sind voneinander unabhängig und können in
  beliebiger Reihenfolge laufen.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **JetStream/Zustellgarantien** — `ADR-0055` entscheidet ausdrücklich
  gegen JetStream; eine künftige Anforderung, die NATS selbst
  Nachvollziehbarkeit tragen lässt, braucht laut ADR-Re-Evaluierungs-
  Trigger eine eigene Folge-ADR, nicht diese Welle.
- **Change-Inhalt im NATS-Payload** — `ADR-0055` legt einen leeren
  Payload fest; eine Erweiterung um Positions-/Change-Daten wäre eine
  Architektur-Änderung, kein Umsetzungs-Detail dieser Welle.
- **Consumer-seitige NATS-Client-Bibliothek/-Beispiel** — diese Welle
  liefert die produzierende Seite (Feed-Container publiziert); ein
  Referenz-Consumer, der NATS abonniert, ist Betreiber-/Nutzer-Sache,
  nicht Repo-Gegenstand (analog zu `register-consumer`/`acknowledge-consumer`,
  die auch nur CLI-Werkzeuge sind, keine Referenz-Anwendung).
- **`LH-FA-SST-007`-Boundary/Negative testseitig ohne echten NATS-Ausfall**
  — die Belege in `slice-054`/`055` müssen real gegen einen laufenden
  oder gezielt gestörten NATS-Server erfolgen, nicht gegen eine
  Mock-Simulation; ein Mock-basierter Unit-Test allein reicht für die
  Welle-Closure nicht (ergänzend zulässig, aber nicht ausreichend).

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

Ergebnis: <Zeiger auf `welle-<NN>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
