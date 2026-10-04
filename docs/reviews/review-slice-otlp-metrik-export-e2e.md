# Review-Report: slice-otlp-metrik-export-e2e — 2026-10-04

**Review-Art:** Code (gegen Plan, [`ADR-0153`](../plan/adr/0153-otlp-einheit-consumer-lag-byte.md), [`ADR-0154`](../plan/adr/0154-spec-luecken-otlp-tls-token-ist-zustand.md) und die unten genannten ADRs, Hard Rules)
**Gegenstand:** Diff 13df935f..HEAD (Commits 8f277c8a, 36317637, 11a0461b, 5952fb27, a9ad1d58)
**Skill:** `.harness/skills/reviewer.md` · **Modell:** Sonnet 5.5

**Eingangs-Kontext:**
- Plan [`slice-otlp-metrik-export-e2e`](../plan/planning/in-progress/slice-otlp-metrik-export-e2e.md)
- [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md)
- [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
- [`LH-FA-SST-010`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.12, §3.13

Der Report-Text stammt vom Reviewer-Lauf; der Reviewer darf keine Report-Dateien anlegen, die Datei hat der Planner aus der Rückgabe angelegt.

## Eigene Proben

Alle Proben liefen an einer `git archive`-Kopie im Scratchpad.

| Probe | Ergebnis |
|---|---|
| `make suchlauf-nachmessen PLAN=<Plan>` | Exit 0, 12 Zeilen stimmen (67/12/59/2/1/1) |
| `make doc-trace` | Exit 0, `83 Anforderung(en), 0 Waise(n).` |
| `make fmt-check` | Exit 0 |
| `make handbuch-public-doc-check` | Exit 0 |
| `make kommentar-kennungen DIFF=13df935f` | Exit 0, ohne Kandidat |
| Abdeckungs-Zählung | 59 Bash-Zeilen plus 21 Go-Zeilen ergeben die 80 Tabellenzeilen |
| Go-Tests `otlpcheck`, `otlpexport` | grün |
| Mutation Soll-Einheit `cdc_consumer_lag` auf `1` im Werkzeug | `TestCheckRejectsEachDeviation` rot |
| Mutation `req.Header.Add` des Exporters entfernt | beide Header-Tests rot |

Den Runner `make test-integration` hat der Reviewer nicht neu gefahren; Aussagen über Collector-Verhalten sind aus den gedruckten Zeilen im Plan **übernommen**.

## Findings

### F-1 — Text per `cat >>` an eine Repo-Datei angehängt
- Kategorie: HIGH · Quelle: Hard Rule Docker-only (`AGENTS.md` §3.1, Umleitungs-Verbot)
- Pfad: `tools/harness/otlpcheck/check_test.go`
- Befund: Der Implementer meldet einen Anhang per `cat >>` statt Edit/Write. Der Guard liest Umleitungen nicht, am Diff ist die Schreibweise nicht ablesbar (Meldung **übernommen**). Der Inhalt ist unauffällig (Test grün, Mutation rot).
- Verifizierbar: nein.

### F-2 — „Gegenlesung des Slots“ trägt der Beleg nicht
- Kategorie: HIGH · Quelle: Skill-HIGH „Beleg trägt seinen Satz nicht“, `AGENTS.md` §3.12 Instanz B
- Pfade: `tools/harness/run-integration-tests.sh` (`abdeckung_declare`-Text), `docs/user/e2e-abdeckung.md`, `harness/README.md` (Zeile `make test-integration`)
- Befund: Drei Träger sagen, `cdc_wal_retention_bytes` werde gegen die „Gegenlesung des Slots“ gehalten; der Code liest die Messreihe aus dem Prozess-Log, nicht `pg_replication_slots`.
- Verifizierbar: ja (`git grep -n "Gegenlesung des Slots"`).

### F-3 — Mutationstabelle: gequotetes Rot zeigt gleiche Werte
- Kategorie: LOW · Pfad: Plan, Mutationstabelle
- Befund: Das gesehene Rot „Export 1, Sicht 1“ trägt die Verschiebung nicht lesbar.

### F-4 — Handbuch: „bleibt die Verbindung aus“ unklar
- Kategorie: LOW · Pfad: `docs/user/benutzerhandbuch.md`, Absatz zum `https`-Empfänger
- Befund: Die Messung zeigt, dass der TLS-Aufbau scheitert und `PCF-W6001` erscheint; der Satz lässt offen, ob versucht wird.

### F-5 — Docker-Hub-Pull ohne Cache-Pfad; Post-Push-Lauf offen
- Kategorie: INFO · Pfad: `tools/harness/run-integration-tests.sh` (Pull des Collector-Images)
- Befund: `docker pull` bricht auch bei gecachtem Image ab, wenn die Registry nicht antwortet; die Meldung nennt das Limit. Der reale `e2e.yml`-Lauf steht aus (`AGENTS.md` §3.10).

### F-6 — Mutationsgrenzen und Hilfsparser
- Kategorie: INFO
- Nicht erprobt, als hergeleitet gekennzeichnet: https-Gegenprobe, 1-MiB-Untergrenze, 16-s-Fenster. `ReadInts` akzeptiert Zeilen wie `12abc` (folgenlos, Eingabe vom Runner).

## Geprüft, ohne Befund

Handbuch-Wahrheit gegen die gedruckten Lauf-Zeilen, Versionshistorie 1.98, Nachbarabschnitte, Suchlauf-Zahlen und Lokatoren, Pin-Form (ein voller Digest an genau einer Stelle, keiner in `docs/`), Runner-Robustheit (Cleanup-Trap, Zeitfenster, Collector-Rechte), fünf von zehn Mutationsbelegen an der Kopie rot, Kommentar-Regeln, Traceability aller fünf Commits.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Verdikt:** Fixrunde am Implementer: F-2 (drei Träger auf „Messreihe des Prozesses“), F-3, F-4; F-1 als Verfahrensverstoß im Beobachtungs-Register festhalten (Inhalt geprüft).
