# Roadmap

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**,
keine Reihenfolge von Terminen (siehe
Baseline-Regelwerk `modul-06-roadmap.md`).
Termine werden — falls überhaupt — als Konsequenz der Wellen-Schätzung
gezeigt, nicht als Treiber.

---

## Offene Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte — *Offene Wellen* ist **derivativ**: Der
Zustand sind die flachen Welle-Dateien; woran gearbeitet wird, sagt das
`Welle:`-Feld der Slices in `in-progress/`. Ziel, Trigger und
Closure-Kriterien stehen in der Welle-Datei, nicht hier.

<!-- BEDIENHINWEIS: Zwei unabhängige Aussagen in diesem Block. Die Liste oben
folgt den Dateien (ein Zeiger je offener Welle-Datei). Trägt in-progress/
keinen Slice, kommt der Ruhe-Marker ZUSÄTZLICH dazu — nicht an Stelle der
Liste; beides zugleich ist der Normalfall direkt nach der Wellen-Eröffnung.
Zwei Aussagen, zwei Wächter: die Marker-Hälfte (Marker genau dann, wenn
in-progress/ keinen Slice trägt) und die Listen-Hälfte (Bijektion Zeiger <->
flache Welle-Dateien, in beide Richtungen; der Marker geht nicht ein). Die
Listen-Hälfte braucht einen Sensor, der das Kardinalitäts-Modell kennt — ein
Ein-Wellen-Wächter hält den Block gegen GENAU EINE Datei und meldet legitime
Zustände (mehrere offene Wellen; eine Welle eröffnet, nichts beansprucht) als
Drift. Welche Hälfte dein Sensor prüft, musst du wissen; eine ungewächterte
Hälfte ist zulässig, wenn sie benannt ist.

WORTLAUT des Markers: Baseline-Regelwerk modul-06-roadmap.md, Bullet
"Offene Wellen". Hier bewusst NICHT zitiert, und das ist eine Regel, keine
Nachlässigkeit: Ein Doku-Sensor matcht den Marker als Substring dieses
Blocks, also matcht sich jeder Regel-, Hinweis- oder Beispieltext selbst,
der den Wortlaut literal trägt — der Block meldete "Ruhe" bei beanspruchtem
Slice. Aus demselben Grund gehört der Wortlaut in keine Sektions-Regel-Zeile
oben. Wer den Sensor selbst baut: Code-Fences beim Matchen aus dem Block
nehmen, sonst schlägt ein Beispiel-Auszug durch. -->

- [welle-17.md](../welle-17.md) — E2E-Testbelege für fünf testfreie Lastenheft-Kennungen (`ADR-0058`).
- [welle-18.md](../welle-18.md) — Spaltenauswahl — Antrags-Queue-Erweiterung, Assembler-Filterung, E2E-Beleg (`LH-FA-CFG-005`, `ADR-0059`).
- [welle-19.md](../welle-19.md) — Live-Change-Streaming — gRPC-Server-Streaming und HTTP/SSE (`LH-FA-SST-008`, `ADR-0060`, `ADR-0061`).

## Nächste Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Nächste Wellen* — die geordnete
Vorschau: je Zeile Welle, Trigger als beobachtbare Bedingung, wichtigste Slices
und geschätzter Aufwand (S/M/L, kein Termin).

| Welle | Trigger | Wichtigste Slices | Geschätzter Aufwand |
|---|---|---|---|

Nichts geplant — die einzige zuvor hier geführte Zeile (`LH-FA-SST-007` —
NATS-Change-Notification) ist mit dieser Änderung als `welle-15` eröffnet
(Zeiger unter *Offene Wellen*).

## Meilensteine

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Welle ≠ Meilenstein ≠ Release.

<!--
Externe Versprechen oder interne Trigger-Punkte.
"M2: erste lauffähige Fassung" ist ein Meilenstein.
Status: "offen" oder "erreicht YYYY-MM-DD" plus Beleg als auflösbarer
Anker (z. B. Tag, Workflow-Lauf, Ergebnis-Notiz). Erreichte Meilensteine
bleiben hier in der Tabelle —
sie gehören nicht ins Drift-Log, und die Status-Zelle erzählt nicht, wie es
dazu kam (Baseline-Regelwerk grundlagen-harness-dateien.md §Was ein
Kommentar trägt, "Dieselbe Regel für Zustandsfelder").
-->

| Meilenstein | Welle(n) | Trigger | Status |
|---|---|---|---|
| M1 — MVP-Abnahme | welle-1, welle-2 | MVP-Integrationstest grün ([Lastenheft §1, MVP-Schnitt](../../../../spec/lastenheft.md)) | erreicht 2026-09-10 — Beleg: `make test-integration` grün am verdrahteten System ([welle-2-results.md](../done/welle-2-results.md), Verifikation) |

## Abhängigkeitsgraph

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Nächste Wellen* — die Abhängigkeit
steht als beobachtbare Bedingung in der `Trigger`-Spalte **und** als gerichtete
Kante hier; eine Welle, die ohne fertige Vorgängerin nicht starten kann, ist
eine Phantom-Welle.

```mermaid
flowchart LR
    W1[welle-1: MVP-Grundlage]
    W2[welle-2: reale PostgreSQL-Integration]
    W3[welle-3: CDC-Verwaltung, Lesen, Sicherheit]
    W4[welle-4: Observability-Vervollständigung]
    W5[welle-5: Realer CDC-Capture-Lag]
    W6[welle-6: Consumer-Zugriffsweg]
    W7[welle-7: Replication-Schwellen-Überwachung]
    W8[welle-8: Black-Box-E2E und Integrationstest-Nachzug]
    W9[welle-9: E2E-Abdeckung — CDC-Kernpfad]
    W9B[welle-10: Schema-Evolution-Nachlieferung ADR-0015]
    W10[welle-11: E2E-Abdeckung — Verwaltung & Observability]
    W10B[welle-12: Verwaltungsfunktionen — SQL-Administration und CLI-Diagnose]
    W11[welle-13: Retention-Löschausführung]
    S047[wellenlos: slice-047/048 Retention-CLI-Sichtbarkeit und kombinierter E2E-Rundlauf]
    W13[welle-14: Performance-Benchmarks & Test-Coverage-Gate]
    S051[wellenlos: slice-051 Walsender-Wirksamkeit isolierter Beleg]
    W15[welle-15: NATS-Change-Notification]
    W16[welle-16: HTTP/JSON-API mit Token-Authn]

    W1 --> W2 --> W3 --> W4 --> W5 --> W6 --> W7 --> W8 --> W9 --> W9B --> W10 --> W10B --> W11 --> S047 --> W13 --> S051 --> W15 --> W16
```

## Abgeschlossene Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

| Welle | Abschluss | Closure-Notiz |
|---|---|---|
| welle-1 — MVP-Grundlage | 2026-09-09 | [welle-1-results.md](../done/welle-1-results.md) |
| welle-2 — Reale PostgreSQL-Integration | 2026-09-10 | [welle-2-results.md](../done/welle-2-results.md) |
| welle-3 — CDC-Verwaltung, Lesen-Vollabdeckung, Sicherheit und Observability-Basis | 2026-09-10 | [welle-3-results.md](../done/welle-3-results.md) |
| welle-4 — Observability-Vervollständigung | 2026-09-10 | [welle-4-results.md](../done/welle-4-results.md) |
| welle-5 — Realer CDC-Capture-Lag | 2026-09-12 | [welle-5-results.md](../done/welle-5-results.md) |
| welle-6 — Consumer-Zugriffsweg | 2026-09-12 | [welle-6-results.md](../done/welle-6-results.md) |
| welle-7 — Replication-Schwellen-Überwachung | 2026-09-12 | [welle-7-results.md](../done/welle-7-results.md) |
| welle-8 — Black-Box-E2E und Integrationstest-Nachzug | 2026-09-12 | [welle-8-results.md](../done/welle-8-results.md) |
| welle-9 — E2E-Abdeckung — CDC-Kernpfad | 2026-09-12 | [welle-9-results.md](../done/welle-9-results.md) |
| welle-10 — Schema-Evolution-Nachlieferung (`ADR-0015`) | 2026-09-12 | [welle-10-results.md](../done/welle-10-results.md) |
| welle-11 — E2E-Abdeckung — Verwaltung & Observability | 2026-09-13 | [welle-11-results.md](../done/welle-11-results.md) |
| welle-12 — Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose | 2026-09-13 | [welle-12-results.md](../done/welle-12-results.md) |
| welle-13 — Retention-Löschausführung | 2026-09-13 | [welle-13-results.md](../done/welle-13-results.md) |
| welle-14 — Performance-Benchmarks & Test-Coverage-Gate | 2026-09-13 | [welle-14-results.md](../done/welle-14-results.md) |
| welle-15 — NATS-Change-Notification | 2026-09-14 | [welle-15-results.md](../done/welle-15-results.md) |
| welle-16 — HTTP/JSON-API mit Token-Authn (`LH-FA-SST-006`, `ADR-0057`) | 2026-09-14 | [welle-16-results.md](../done/welle-16-results.md) |

## Historische Trigger-Verschiebungen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Historische Trigger-Verschiebungen*
— das Drift-Log: jede Umplanung mit Datum, Änderung, Grund. Leer heißt starre
Roadmap, jede Zeile voll heißt treibende.

<!--
Wenn Wellen umgeplant wurden: Datum, Grund, neue Reihenfolge.
Steering-Loop-relevant.
NUR Umplanungen: Trigger verschoben, präzisiert oder ersetzt; Slice oder
Welle umgehängt. KEINE Schließungen (die stehen im Closure-Log oben) und
KEINE erreichten Meilensteine (Status-Spalte) — sonst führt diese Tabelle ein
zweites Closure-Log, und zwei Logs driften.
-->

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-09-12 | Neue Feature-Welle „Schema-Evolution-Nachlieferung (`ADR-0015`)" zwischen `welle-9` und „E2E-Abdeckung — Verwaltung & Observability" eingefügt; deren Trigger von „`welle-9` liegt in `done/`" auf „Vorherige Welle (Schema-Evolution-Nachlieferung) liegt in `done/`" umgehängt. | `slice-030` fand real, dass [`ADR-0015`](../../adr/0015-schema-evolution.md)s Folgepflicht (`SchemaStorePort`, dynamische Re-Versionierung) nie eingelöst wurde (`docs/reviews/review-slice-030.md` F-1 HIGH, Architect-Verdikt `docs/reviews/architect-verdict-slice-030-adr-0015.md`) — Größenordnung Feature-Welle, nicht Einzel-Slice. |
| 2026-09-12 | Neue Feature-Welle „Verwaltungsfunktionen — SQL-Administration (`LH-FA-ADM-001`)" zwischen „E2E-Abdeckung — Verwaltung & Observability" und „Retention-Löschausführung" eingefügt; „E2E-Abdeckung — Verwaltung & Observability"s Scope auf den bereits existierenden Lesezugriff (`LH-FA-CFG-003`/`004`, `LH-FA-ADM-002`…`005`) präzisiert, ihr fälschlicher Bezug zu `BEO-PGC/rollen-test-abdeckungsluecken` gestrichen; „Retention-Löschausführung"s Trigger entsprechend umgehängt. | Fork-Recherche zur Eröffnung von `welle-11` fand real, dass [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (Lastenheft, Rang 1) explizite SQL-Administrationsfunktionen (Aktivierung/Deaktivierung/Status/Consumer-Verwaltung) verlangt, die nicht existieren, und dass `LH-FA-CFG-002` (Deaktivierung) keinen Live-Zugriffsweg hat — registriert als `BEO-PGC/verwaltung-keine-sql-administration`; Größenordnung Feature-Welle, kein Testabdeckungs-Defizit. |
| 2026-09-13 | „Retention-Löschausführung" als `welle-13` eröffnet (verlässt *Nächste Wellen*, Zeiger unter *Offene Wellen*); „E2E-Abdeckung — Retention"s Trigger von „Vorherige Welle (Retention-Löschausführung) liegt in `done/`" auf „`welle-13` liegt in `done/`" präzisiert. | `welle-12` liegt in `done/`, Architect-Verdikt zur Löschausführung liegt vor (`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`) — Eröffnungs-Trigger erfüllt. |
| 2026-09-13 | „E2E-Abdeckung — Retention" aus *Nächste Wellen* entfernt — kein eigener Welle-Schnitt mehr, sondern der wellenlose `slice-047` (Diagnose-CLI-Erweiterung für Retention-Sichtbarkeit); „Performance-Benchmarks & Test-Coverage-Gate"s Trigger von „Vorherige Welle (E2E-Abdeckung Retention) liegt in `done/`" auf „`slice-047` liegt in `done/`" umgehängt. | Eröffnungs-Recherche (Fork) fand real: Die SQL-View-Black-Box-Ebene für Retention ist bereits vollständig geliefert (`slice-045`/`046` selbst, kein Go-Domain-Import); nur die `diagnose`-CLI-Sichtbarkeit fehlt real — ein einzelner Slice ohne Closure-Bedingung jenseits seiner eigenen DoD trägt laut Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht keine eigene Welle. |
| 2026-09-13 | Header-Absatz „E2E-Abdeckungsprogramm" (Nutzerentscheidung 2026-09-12, ursprünglich vier Wellen: `welle-9`, `welle-11`, `welle-13`, „E2E-Abdeckung — Retention") aus *Nächste Wellen* entfernt — er stand fälschlich über den drei völlig unabhängigen Folgezeilen (Performance-Benchmarks, Publication-Entzug, NATS), nachdem die vorherige Korrektur seine letzte zugehörige Tabellenzeile entfernt hatte, ohne den darüberstehenden Programm-Header mitzuziehen. Das Programm selbst ist damit vollständig abgeschlossen: drei seiner vier Wellen liegen in `done/` (`welle-9`, `welle-11`, `welle-13`), die vierte wurde wie oben beschrieben als wellenloser `slice-047`/`048` statt als eigene Welle realisiert. | Nutzerfrage („warum steht das in roadmap.md?") deckte real auf, dass der Header nach der vorherigen Zeilen-Entfernung verwaist stehen geblieben war. |
| 2026-09-13 | „Performance-Benchmarks & Test-Coverage-Gate" als `welle-14` eröffnet (verlässt *Nächste Wellen*, Zeiger unter *Offene Wellen*); „Publication-Entzug-Wirksamkeit"s Trigger von „Vorherige Welle (Performance-Benchmarks) liegt in `done/`" auf „`welle-14` liegt in `done/`" präzisiert. | `slice-047` liegt in `done/`, Architect-Verdikt zu Scope/Schwelle/Suppression liegt vor ([ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)) — Eröffnungs-Trigger erfüllt; zwei unabhängige Slices geschnitten (`slice-049` Coverage-Gate, `slice-050` Benchmark-Infrastruktur). |
| 2026-09-13 | „Publication-Entzug-Wirksamkeit am laufenden Stream" aus *Nächste Wellen* entfernt — kein eigener Welle-Schnitt mehr, sondern der wellenlose `slice-051` (isolierter Walsender-Timing-Beleg); „`LH-FA-SST-007` — NATS-Change-Notification"s Trigger von „Vorherige Welle (Publication-Entzug-Wirksamkeit) liegt in `done/`" auf „`slice-051` liegt in `done/`" umgehängt. | Eröffnungs-Recherche (Fork) plus Architect-Verdikt (`docs/reviews/architect-verdict-walsender-wirksamkeit.md`) fanden real: kein Mehr über einen einzelnen isolierten Testabschnitt hinaus, keine Schema-/Domänen-Änderung — trägt laut Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht keine eigene Welle. |
| 2026-09-13 | „`LH-FA-SST-007` — NATS-Change-Notification" als `welle-15` eröffnet (verlässt *Nächste Wellen*, Zeiger unter *Offene Wellen*) — vier Slices geschnitten (`slice-052` Port/Adapter, `slice-053` Compose-Verdrahtung/Happy-Path, `slice-054` Boundary-Beleg, `slice-055` Negative-Beleg). Damit ist *Nächste Wellen* leer — es steht keine weitere Welle in der Vorschau. | `slice-051` liegt in `done/`, [ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md) (Accepted) liegt vor — Eröffnungs-Trigger erfüllt; die im Lastenheft komplett unimplementierte Anforderung erfordert laut ADR mindestens vier voneinander abhängige Liefer-Punkte (Grundgerüst, Verdrahtung, Boundary-Beleg, Negative-Beleg), die einzeln DoD-fähig, aber gemeinsam erst `LH-FA-SST-007` vollständig belegen — genau das *Mehr* gegenüber den Einzel-DoDs. |
