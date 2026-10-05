# Review-Report: slice-otlp-metrik-export-e2e — 2026-10-04

**Review-Art:** Code — gegen Plan, [`ADR-0153`](../plan/adr/0153-otlp-einheit-consumer-lag-byte.md), [`ADR-0154`](../plan/adr/0154-spec-luecken-otlp-tls-token-ist-zustand.md), die unten genannten ADRs und die Hard Rules

**Gegenstand:** Diff 13df935f..HEAD (Commits 8f277c8a, 36317637, 11a0461b, 5952fb27, a9ad1d58)

**Skill:** `.harness/skills/reviewer.md` @ 675246dd
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-04

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-otlp-metrik-export-e2e`
- [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md)
- [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
- [`LH-FA-SST-010`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.12, §3.13

Der Report-Text stammt vom Reviewer-Lauf; der Reviewer darf keine Report-Dateien anlegen, die Datei hat der Planner aus der Rückgabe angelegt.

**Eigene Proben** (alle an einer `git archive`-Kopie im Scratchpad):

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

---

## Findings

Jedes Finding folgt dem §Output-Schema des Reviewer-Skills.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Implementer meldet einen Anhang per `cat >>` statt Edit/Write. Der Guard liest Umleitungen nicht, am Diff ist die Schreibweise nicht ablesbar (Meldung **übernommen**). Der Inhalt ist unauffällig (Test grün, Mutation rot). | Hard Rule Docker-only (`AGENTS.md` §3.1, Umleitungs-Verbot) | `tools/harness/otlpcheck/check_test.go` · Kurzzitat nicht erhoben | nein | Docker-only-Verstoß |
| F-2 | HIGH | Drei Träger sagen, `cdc_wal_retention_bytes` werde gegen die „Gegenlesung des Slots“ gehalten; der Code liest die Messreihe aus dem Prozess-Log, nicht `pg_replication_slots`. | Skill-HIGH „Beleg trägt seinen Satz nicht“, `AGENTS.md` §3.12 Instanz B | `tools/harness/run-integration-tests.sh` (`abdeckung_declare`-Text), `docs/user/e2e-abdeckung.md`, `harness/README.md` (Zeile `make test-integration`) · „Gegenlesung des Slots“ | ja — `git grep -n "Gegenlesung des Slots"` | Beleg trägt seinen Satz nicht |
| F-3 | LOW | Das gesehene Rot „Export 1, Sicht 1“ der Mutationstabelle trägt die Verschiebung nicht lesbar. | nicht erhoben | Plan `slice-otlp-metrik-export-e2e`, Mutationstabelle · Kurzzitat „Export 1, Sicht 1“ | nicht erhoben | nicht erhoben |
| F-4 | LOW | Die Messung zeigt, dass der TLS-Aufbau scheitert und `PCF-W6001` erscheint; der Satz lässt offen, ob versucht wird. | nicht erhoben | `docs/user/benutzerhandbuch.md`, Absatz zum `https`-Empfänger · „bleibt die Verbindung aus“ | nicht erhoben | nicht erhoben |
| F-5 | INFO | `docker pull` bricht auch bei gecachtem Image ab, wenn die Registry nicht antwortet; die Meldung nennt das Limit. Der reale `e2e.yml`-Lauf steht aus (`AGENTS.md` §3.10). | `AGENTS.md` §3.10 | `tools/harness/run-integration-tests.sh` (Pull des Collector-Images) · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |
| F-6 | INFO | Nicht erprobt, als hergeleitet gekennzeichnet: https-Gegenprobe, 1-MiB-Untergrenze, 16-s-Fenster. `ReadInts` akzeptiert Zeilen wie `12abc` (folgenlos, Eingabe vom Runner). | nicht erhoben | `ReadInts` · Datei nicht erhoben | nicht erhoben | Mutationsgrenzen und Hilfsparser |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Handbuch-Wahrheit gegen die gedruckten Lauf-Zeilen | geprüft, ohne Befund |
| Versionshistorie 1.98 und Nachbarabschnitte | geprüft, ohne Befund |
| Suchlauf-Zahlen und Lokatoren | geprüft, ohne Befund |
| Pin-Form (ein voller Digest an genau einer Stelle, keiner in `docs/`) | geprüft, ohne Befund |
| Runner-Robustheit (Cleanup-Trap, Zeitfenster, Collector-Rechte) | geprüft, ohne Befund |
| Mutationsbelege (fünf von zehn an der Kopie rot) | geprüft, ohne Befund |
| Kommentar-Regeln | geprüft, ohne Befund |
| Traceability aller fünf Commits | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Docker-only-Verstoß · Beleg trägt seinen Satz nicht · Mutationsgrenzen und Hilfsparser (F-3 und F-4 ohne Klassenbezeichnung: nicht erhoben)

## Verdikt

**Merge-blockierend:** ja — Fixrunde am Implementer: F-2 (drei Träger auf „Messreihe des Prozesses“), F-3, F-4; F-1 als Verfahrensverstoß im Beobachtungs-Register festhalten (Inhalt geprüft).

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
