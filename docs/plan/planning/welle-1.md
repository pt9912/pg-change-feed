# Welle welle-1: MVP-Grundlage — Bootstrap, Domänenkern, Capture-Persist

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<NN>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M1 (MVP-Abnahme) — diese Welle legt seine Grundlage, der
Meilenstein schließt erst mit dem MVP-Integrationstest.

**Verantwortlich:** pt9912. **Datum:** 2026-09-09.

---

## 1. Welle-Ziel

<!-- BEDIENHINWEIS: Eine Aussage, die sich an einem Lasttest oder
Akzeptanzkriterium spiegelt. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die MVP-Grundlage steht: Go-Modul und Build-Vertrag sind real baut
([`LH-QA-POR-003`](../../../spec/lastenheft.md)), der Domänenkern trägt seine
Invarianten ([`ADR-0029`](../../../docs/plan/adr/README.md)), und der Capture-Pfad
erfüllt die Persist-before-ACK-Ordnung gegen Fake Ports
([`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md)). Nach dieser Welle ist die
Grundlage für die reale PostgreSQL-Integration (Welle 2) belegt — gemessen
an den Abnahmekriterien des MVP ([Lastenheft §1](../../../spec/lastenheft.md),
MVP-Schnitt) als Fortschritt, nicht als Abschluss.

## 2. Trigger (Welle startet)

<!-- BEDIENHINWEIS: Was muss vorher passiert sein? -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `make gates` grün über den Architektur-Bestand: alle 40 ADRs
  ([ADR-0001](../../../docs/plan/adr/README.md)…0040) committet, die drei Spec-Straten committet, das
  Doc-Gate grün (Beleg: Gate-Lauf nach a-check-Vorbereitung).
- [`ADR-0039`](../../../docs/plan/adr/README.md) und
  [`ADR-0040`](../../../docs/plan/adr/README.md) `Accepted` (Paketstruktur und
  Port-Bestand stehen fest) — beobachtbar am ADR-Index.

## 3. Closure-Trigger (Welle schließt)

<!-- BEDIENHINWEIS: Aktion, nicht Termin. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/` (slice-001, slice-002,
  slice-003).
- `make gates` grün **und** `make a-check` grün (aktivierte Fassung über
  den entstehenden Baum) — repo-weite Belege, die in keinem einzelnen
  Slice-DoD stehen.
- `make image` baut das Image mit dem neuen `internal/`-Baum;
  `harness/image-hash.txt` trägt den Digest.

## 4. Slices in dieser Welle

<!-- BEDIENHINWEIS: keine Status-Spalte ergaenzen. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-001 | Go-Modul-Bootstrap und Gate-Aktivierung | [`LH-QA-POR-003`](../../../spec/lastenheft.md), [`LH-QA-OPS-001`](../../../spec/lastenheft.md) |
| slice-002 | Domänenkern — Modelle, Invarianten, ClockPort | [`LH-FA-DAT-001`](../../../spec/lastenheft.md), [`LH-FA-DAT-004`](../../../spec/lastenheft.md) |
| slice-003 | Capture-Persist-Pfad — Ports, Service, Persist-before-ACK | [`LH-QA-REL-001`](../../../spec/lastenheft.md), [`LH-QA-REL-002`](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

<!-- BEDIENHINWEIS: Falls jemand diese Welle aendert — was bricht? -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle 2 (reale PostgreSQL-Integration) — sie baut auf dem
  Domänenkern und den Ports dieser Welle auf; der Store-Adapter
  ([`ADR-0009`](../../../docs/plan/adr/README.md)) braucht die Port-Verträge.
- Wird blockiert von: nichts — der Start-Trigger ist mit dem
  Architektur-Bestand eingetreten.

## 6. Out-of-Scope für diese Welle

<!-- BEDIENHINWEIS: explizite Nicht-Inhalte, schuetzt vor Scope-Creep. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Reale PostgreSQL-Adapter (Store, ACK, Metadata) und Integrationstests —
  Welle 2; die Fak dieser Welle halten die Ordnungs-Logik prüfbar, aber
  keinen Treiber-Vertrag.
- MVP-Abnahmekriterien als Abschluss — der Meilenstein M1 schließt erst mit
  dem MVP-Integrationstest; diese Welle legt nur die Grundlage.
- Consumer-Verwaltung, Retention, SQL-/CLI-Adapter — spätere Wellen;
  [`ADR-0028`](../../../docs/plan/adr/README.md) listet die Use Cases, die hier
  nicht angetastet werden.
- `codepaths`-Aktivierung und `image-cve`-Sensor — die Pfade/Belege
  entstehen später; die Bedingungen stehen in `.d-check.yml` und
  `harness/README.md`.

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
