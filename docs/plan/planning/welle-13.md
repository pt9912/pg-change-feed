# Welle 13: Retention-Löschausführung

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-13-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-13.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`RetentionPolicy.AllowsDeletion` (`internal/domain/model/retention.go`) ist
reine Entscheidungslogik ohne Aufrufer — der eigene Doc-Kommentar benennt
den fehlenden „Run-Retention-Use-Case" explizit, aber er existiert nicht:
`grep -rn "AllowsDeletion"` trifft ausschließlich die eigene Testdatei.
[`ADR-0014`](../../../spec/lastenheft.md) (Retention als Domain Policy,
Accepted, `permanent`) benennt diesen Use-Case bereits als Folgepflicht —
derselbe Präzedenzfall wie `slice-030`/`ADR-0015`: die Entscheidung steht,
ihre Umsetzung fehlt. Diese Welle liefert die reale Löschausführung
([`LH-FA-RET-002`](../../../spec/lastenheft.md)…`004`), die Sichtbarkeit
blockierender Consumer ([`LH-FA-RET-005`](../../../spec/lastenheft.md))
und die Metrik `cdc_storage_bytes`
([`LH-FA-RET-006`](../../../spec/lastenheft.md)), die
`tools/schema/nacharbeit-observability.sql` bisher explizit als „nicht
abgedeckt" ausweist. Ein Architect-Verdikt
([`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`](../../reviews/architect-verdict-retention-loeschausfuehrung.md))
hat vorab bestätigt: keine neue ADR nötig — bestehende ADRs
([`ADR-0009`](../adr/0009-change-store-outbound-port.md),
[`ADR-0011`](../adr/0011-persist-before-ack.md),
[`ADR-0012`](../adr/0012-at-least-once.md),
[`ADR-0014`](../adr/0014-retention-domain-policy.md)) tragen die
Port-Erweiterung bereits, und die `cdc_storage_bytes`-Metrik kann dem
etablierten View-Owner-Muster (`cdc.metrics`, `cdc.heartbeat`) folgen,
ohne die `cdc_reader`-Rolle zu erweitern.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass ein Administrator eine reale Löschausführung anstoßen kann,
die `RetentionPolicy.AllowsDeletion` tatsächlich befragt, `cdc.change`-
Zeilen real entfernt, blockierende Consumer sichtbar macht und den
resultierenden Speicherverbrauch über `cdc_storage_bytes` misst — das ist
erst die Summe aller vier Slices.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `welle-12` liegt in `done/`.
- Architect-Verdikt zur Retention-Löschausführung liegt vor
  (`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`).
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make gates` grün.
- Ein realer Testlauf gegen den Compose-Stack belegt: eine reale
  Löschausführung entfernt `cdc.change`-Zeilen, für die
  `RetentionPolicy.AllowsDeletion` real `true` liefert (alle Consumer
  haben bestätigt, `MinAge` erreicht), und lässt Zeilen unangetastet, für
  die das nicht gilt (ein Consumer hängt zurück).
- Dieselbe Sichtbarkeit für einen blockierenden Consumer: eine SQL-Sicht
  oder CLI-Ausgabe zeigt real, welcher Consumer bis zu welcher Position
  die Löschung verhindert.
- `cdc_storage_bytes` liefert real einen numerischen Wert über
  `cdc.metrics`, ohne dass `cdc_reader`s Grant-Fläche erweitert wurde.
- `BEO-PGC/retention-keine-loeschausfuehrung` erreicht Ausgang
  *verkörpert*.
- Closure-Notiz in `welle-13-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-043 | `ChangeStorePort`-Löschmethode und `RunRetentionUseCase` | [ADR-0009](../adr/0009-change-store-outbound-port.md), [ADR-0014](../adr/0014-retention-domain-policy.md) |
| slice-044 | Hintergrund-Job/CLI-Trigger für die Löschausführung | [ADR-0014](../adr/0014-retention-domain-policy.md) |
| slice-045 | Sichtbarkeit blockierender Consumer | [LH-FA-RET-005](../../../spec/lastenheft.md) |
| slice-046 | `cdc_storage_bytes`-Metrik | [LH-FA-RET-006](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: „E2E-Abdeckung — Retention" (in der Roadmap direkt
  nachfolgend eingereiht — ohne diese Welle gäbe es nichts zu testen).
- Wird blockiert von: keine andere Welle.
- Intern: `slice-043` → `slice-044` (der Hintergrund-Job ruft den in
  `slice-043` gebauten Use-Case auf). `slice-045` und `slice-046` sind
  von `slice-043`/`044` unabhängig (reine Sichtbarkeits-/Metrik-Arbeit
  auf bestehendem bzw. neuem, aber nicht wechselseitig benötigtem
  Zustand) und können parallel oder in beliebiger Reihenfolge laufen.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Black-Box-E2E-Testabdeckung der neuen Löschausführung** — eigene,
  bereits in der Roadmap vorgemerkte Folge-Welle „E2E-Abdeckung —
  Retention"; diese Welle liefert die Fähigkeit, nicht ihre externe
  Testabdeckung (dieselbe Trennung wie bei `welle-9`/`welle-10`
  gegenüber `welle-11`).
- **`SPEC-<NNN>`-Verfeinerung des Lösch-Ausführungsmechanismus** (CLI-Name,
  Job-Takt, Batch-Größe) — entsteht als Teil des jeweils umsetzenden
  Slice (`slice-043`/`044`), nicht vorab hier; `LH-FA-RET-004.a` (Safe
  Watermark) existiert bereits und deckt die Auswahllogik ab.
- **Erweiterung der Least-Privilege-Rollen** (`cdc_reader`/`cdc_admin`) —
  der Architect-Verdikt bestätigt, dass weder die Löschausführung noch
  `cdc_storage_bytes` eine Rollen-Erweiterung braucht; träfe das während
  der Umsetzung doch zu, wäre das ein eigener Architect-Zug, kein
  stiller Fortschritt dieser Welle.
- **Konfigurierbarkeit der Retention-Policy zur Laufzeit** (z. B. über
  die neue YAML-Konfigurationsdatei aus `slice-041`) — `RetentionPolicy`
  bleibt, wie sie ist (`MinAge`, Consumer-Positionen); eine
  Laufzeit-Konfigurationsanbindung ist ein anderer Vorgang.

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
