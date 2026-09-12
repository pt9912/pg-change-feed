# Welle 7: Replication-Schwellen-Überwachung

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-7-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-12.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`SPEC-008`](../../../spec/pflichtenheft.md) verlangt für die Fehlerklasse
`replication`: „Überwachung über Schwellen (§5, WAL-Rückstand); kontrollierte
Fortsetzung." Real ist davon nichts umgesetzt
(`BEO-PGC/spec008-replication-luecke`): Jeder als `replication`
klassifizierte Fehler beendet den Prozess sofort, die Metrik
`cdc_wal_retention_bytes` existiert nicht, und
[`SPEC-013`](../../../spec/pflichtenheft.md) definiert Schwellenwerte nur für
`cdc_capture_lag`, nicht für WAL-Rückstand.

Beim Nachvollziehen des Codes zeigt sich zusätzlich: Die heutige
`replication`-Klasse vermischt zwei verschiedene Fehlerarten —
Stream-Ordnungs-Verletzungen (`mapper.ErrChangeWithoutBegin` u. ä., echte
Dateintegritäts-Korruption, muss hart abbrechen) und
Transport-/Verbindungsstörungen (`receive.ErrReplication`,
`outbound.ErrReplication`). Diese Welle baut die WAL-Rückstand-Überwachung
als **eigenständigen, proaktiven Health-Check** (periodische Messung gegen
`pg_replication_slots`), unabhängig vom reaktiven Fehlerpfad: unterhalb der
Warnschwelle läuft der Capture-Betrieb unverändert weiter (die geforderte
„kontrollierte Fortsetzung"), oberhalb der Fehlerschwelle eskaliert die
Überwachung selbst zu einem `replication`-Fehler (Risiko eines WAL-Verlusts).
Stream-Ordnungs-Verletzungen bleiben von dieser Welle unberührt — sie lösen
weiterhin sofort aus, unabhängig vom WAL-Rückstand.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass ein künstlich erzeugter WAL-Rückstand unterhalb der
Warnschwelle den Capture-Betrieb tatsächlich unverändert fortsetzen lässt
und derselbe Rückstand oberhalb der Fehlerschwelle real zum Abbruch führt.
Das zeigt erst ein Ende-zu-Ende-Test, der beide Seiten der Schwelle real
durchläuft.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `welle-6` liegt in `done/`.
- `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/` steht
  auf `weiter offen` (Zähler 1×) — der Auftrag, die Lücke tatsächlich zu
  schließen, kommt von außerhalb der Welle-Mechanik (Nutzeranweisung, nicht
  Register-Schwelle 3×).
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices (`slice-024`, `slice-025`, `slice-026`) liegen in `done/`.
- `make gates` grün.
- Ein Ende-zu-Ende-Test (Teil der DoD von `slice-026`) belegt real, mit
  künstlich erzeugtem WAL-Rückstand: unterhalb der Warnschwelle setzt der
  Capture-Betrieb unverändert fort; oberhalb der Fehlerschwelle bricht er
  kontrolliert mit `replication`-Klassifikation ab. Das ist das *Mehr* dieser
  Welle — kein Einzel-Slice-DoD deckt beide Seiten der Schwelle im selben Lauf.
- Closure-Notiz in `welle-7-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-024 | ADR — Fehlerklassen-Trennung und Schwellen-Präzisierung (`SPEC-008`/`SPEC-013`) | [`SPEC-008`](../../../spec/pflichtenheft.md), [`SPEC-013`](../../../spec/pflichtenheft.md) |
| slice-025 | Metrik `cdc_wal_retention_bytes` | [`SPEC-009`](../../../spec/pflichtenheft.md) |
| slice-026 | Schwellen-Überwachung mit kontrollierter Fortsetzung im Capture-Pfad | [`SPEC-008`](../../../spec/pflichtenheft.md), [`SPEC-013`](../../../spec/pflichtenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keine andere Welle.
- Intern (Slice-Reihenfolge innerhalb dieser Welle, kein Welle-zu-Welle-Bezug):
  `slice-025` und `slice-026` setzen die ADR aus `slice-024` voraus (sie legt
  fest, welche Fehler-Sentinels als Transport-/Verbindungsstörung gelten und
  welche Schwellenwerte gelten). `slice-026` setzt zusätzlich die Metrik aus
  `slice-025` voraus (die Schwellen-Überwachung liest `cdc_wal_retention_bytes`).
  Reihenfolge: `slice-024` → `slice-025` → `slice-026`.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Stream-Ordnungs-Verletzungen bleiben hart abbrechend.** `mapper.ErrChangeWithoutBegin`,
  `ErrCommitWithoutBegin`, `ErrBeginWithoutCommit` (echte Dateintegritäts-Korruption)
  ändern sich durch diese Welle nicht — Bestand bleibt bewusst stehen. Eine
  Lockerung dieses Pfads wäre eine eigenständige, riskante Entscheidung und
  keine Nebenwirkung der Schwellen-Überwachung.
- **`BEO-PGC/walsender-wirksamkeit` (Publication-Entzug-Timing) bleibt unberührt.**
  Ein verwandtes, aber eigenständiges Replication-Stream-Thema — es wäre ein
  anderer Vorgang mit eigenem Trigger, kein Teil der WAL-Rückstand-Schwelle.
- **Kein generisches Alerting-/Notification-System.** Diese Welle liefert die
  Metrik und die Schwellen-Reaktion im Prozess selbst (Log, kontrollierte
  Fortsetzung/Abbruch) — Weiterleitung an ein externes Alerting-System (E-Mail,
  Pager, Webhook) wäre ein anderer Vorgang (Integrationsarbeit, keine
  Fehlerklassen-Präzisierung).
- **`LH-FA-SST-006`/`LH-FA-SST-007`** (HTTP/gRPC-API bzw. NATS-Change-Notification,
  beide neu im Lastenheft) bleiben unberührt — Schicht-Abgrenzung: Diese Welle
  arbeitet im Capture-Pfad/Observability, nicht am Consumer-Zugriffsweg.

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

Ergebnis: [`welle-7-results.md`](welle-7-results.md), Geschwister im Ruheort `done/`.
Zähler: [`../observations/`](../observations/)`BEO-PGC/`, eine Ebene über dem Ruheort.
