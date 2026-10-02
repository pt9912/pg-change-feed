# Slice meldungscodes-warnungen-heartbeat-diagnose: Warn-Codes als Log-Attribut, Spalte `error_code` im Heartbeat, Code in der Diagnose (Teil 3 von 4)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md) (Betriebsfähigkeit),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (Fehlerzustand erkennbar),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (Fehlerklassen),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Festlegung 1 `W`-Bereiche,
3 Abbildung, 9 T3), [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
(Heartbeat, Schwellen), [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
(Vorlauf bei View-Signatur-Änderung), [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).
Verdikt: [`architect-verdict-meldungscodes-statt-interner-kennungen`](../../../reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md).

**Berührte Spec-Stellen:** Tabelle von `cdc.process_heartbeat` in
[`spec/pflichtenheft.md`](../../../../spec/pflichtenheft.md) (additive Spalte `error_code`;
Liefer-Punkt dieses Slice), die `Diagnose`-Zeile des gRPC-/HTTP-Vertrags (additives Feld
`error_code`) — jeweils die Stelle, die die Spalte bzw. das Feld heute beschreibt.

**Reihenfolge:** nach [T2 `meldungscodes-registry-fehlerkopf`](slice-meldungscodes-registry-fehlerkopf.md)
(Tabelle, Katalog, Gate); unabhängig von [T4](slice-meldungscodes-http-grpc-fehlerkoerper.md).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `ba60c7bc`, 2026-10-02; Befehle im Feld §3).**
Nach T2 tragen klassifizierte Fehler Codes; Warnungen und der Betriebszustand noch nicht.
`error_class` kommt in Go-Produktion in 14 Zeilen vor (Treffer in `queries.go`, `http/diagnose.go`,
`gen/cdc/administration/v1/administration.pb.go`, `examples/grpc-client/diagnose.go`,
`tools/harness/grpcadminclient/main.go`), in `tools/schema` in 7 Zeilen
(`nacharbeit-heartbeat.sql`, `nacharbeit-observability.sql`, `schema.yaml`), im Handbuch in 7.
Log-Warnungen: 35 Zeilen mit `.Warn(` in Produktions-Go (nicht jede ist eine Warnung im Sinn
der ADR; der Implementer trennt Betreiber-Warnungen mit Maßnahme von Fortschritts-/Fehlerlogs).

**Ziel:** Warnungen tragen einen `W`-Code als Log-Attribut `code=<code>`
([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1: Bereiche
1 Erfassung/Replikation, 2 Backfill, 3 Retention/Speicher, 4 Verwaltung, 5 Konfiguration/Start);
`cdc.process_heartbeat` und `cdc.heartbeat` tragen die additive, nullable Spalte `error_code`
neben unveränderter `error_class` (`NULL` bei Normalbetrieb); `diagnose` (CLI, `GET /diagnose`,
RPC `Diagnose`) nennt in der Zeile „Fehlerzustand“ Klasse und Code (`schema [<code>]`) und
trägt das additive Feld `error_code`. **Alle anderen Diagnose-Zeilen tragen weder Kennung
noch Code** (Zustandsberichte, Festlegung 7 der ADR).

**Ausdrücklich NICHT in diesem Slice:**

- **Registry, Fehlerkopf, Katalog-Gerüst** — T2; dieser Slice erweitert Tabelle und Katalog
  um die `W`-Codes und hält das Gate `meldungscodes-check` grün.
- **HTTP-Fehlerkörper `code` und gRPC `ErrorInfo`** — T4.
- **Änderung von `error_class`, Metrik-Label, Schwellen** — bleiben
  ([`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)); kein Code als Label.
- **Numerischer Ausgang je Code** — Ausgang 1 bleibt.
- **SQL-Funktionen** — kein Code (0 `RAISE EXCEPTION`).
- **Release/Tag, SDK-Änderungen** — Freigaben des Auftraggebers. Die Diagnose-Antwort ist
  additiv; ob ein SDK sie modelliert, ist nicht Gegenstand.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(A) Warn-Codes.** Die Betreiber-Warnungen (Liste im Bericht, Auswahl begründet:
      Maßnahme des Betreibers, Granularität wie in der ADR) tragen `code=<code>` als Attribut;
      Tabelle und Handbuch-Katalog um die `W`-Codes erweitert; `meldungscodes-check` und der
      Registry-Test (Ziffer gegen Bereich) grün. *Zu belegen durch:* Test je Warn-Code, der das
      Attribut einer ausgelösten Warnung liest (Happy), der Registry-Test; Suchlauf §3.
- [ ] **(B) Heartbeat-Spalte und Diagnose.** Additive Spalte `error_code` in
      `cdc.process_heartbeat` und `cdc.heartbeat` samt Rollout (Vorlauf bei View-Signatur-Änderung,
      [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)) und Schreib-/Lesepfad;
      `diagnose` in CLI, HTTP und gRPC: Zeile „Fehlerzustand“ mit Klasse und Code, Feld
      `error_code`; Proto-Erzeugnis (`make proto-generate`, committet), Handbuch-Beispiele =
      echte Ausgabe; Spec-Nachzug (Heartbeat-Tabelle, Diagnose-Vertrag). *Zu belegen durch:*
      `make test-store` (Spalte, Rollen, zweiter Rollout gegen migriertes Ziel Exit 0,
      Rollout-Wache `tools/harness/run-schema-rollout-guard-test.sh`), `make generated-sync`
      grün, `make test-integration` (Fehlerzustand-Beleg: der direkt in `cdc.process_heartbeat`
      geschriebene Fehlerzustand erscheint mit Code, Normalbetrieb mit `NULL`), die gedruckten
      Zeilen eines realen `diagnose`-Laufs (**gemessen**) gegen das Handbuch.
- [ ] **(C) Kompatibilität des Lesers.** Die Spalte und das Feld sind additiv: bestehende Leser
      (`cdc.heartbeat`-Konsumenten, SDKs, Beispiel-Clients) brechen nicht. *Zu belegen durch:*
      Messung der Leser mit Befehl im Bericht (`git grep` der Spalten-/Feldlisten in `sdks/`,
      `examples/`, `tools/harness/`; Ergebnis, nicht Erwartung), `make examples-*`-Bau soweit
      ein Beispiel-Client das Feld liest.
- [ ] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-warnungen-heartbeat-diagnose.md`
      endet mit Exit 0.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten, nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/` (Tabelle), Registry-Test | update | `W`-Codes (A). |
| Log-Stellen der Betreiber-Warnungen (`internal/…`, `cmd/…`) | update | Attribut `code`. |
| `tools/schema/schema.yaml`, `tools/schema/nacharbeit-heartbeat.sql`, `tools/schema/nacharbeit-observability.sql`, Rollout-Vorlauf | update | Spalte `error_code`, Views, View-Signatur-Vorlauf. |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, Heartbeat-Port/-Adapter | update | Schreiben und Lesen der Spalte. |
| `proto/cdc/administration/v1/administration.proto`, `gen/…` | update | Feld `error_code` im Diagnose-Ergebnis; Erzeugnis committet. |
| `internal/bootstrap/wiring.go` (Diagnose-Ausgabe), `internal/adapters/driving/http/diagnose.go`, `tools/harness/grpcadminclient/main.go`, `examples/grpc-client/diagnose.go` | update | Zeile „Fehlerzustand“, Feld. |
| `spec/pflichtenheft.md` (Heartbeat-Tabelle, Diagnose-Vertrag) | update | Spec-Nachzug. |
| `docs/user/benutzerhandbuch.md` | update | Katalog der `W`-Codes, Diagnose-Beispiele, Änderungshistorie. Rebase auf `main`, nur eigene Abschnitte committen. |
| Läufer unter `tools/harness/` (Erwartungen an Diagnose-Zeile) | update | Erwartungen. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „Heartbeat/Diagnose tragen nur die
Klasse“; Parent ist `ba60c7bc`; die `diff`-Zeilen und die Befunde trägt der Implementer ein):**

```suchlauf
ba60c7bc 7 -n -F 'error_class' -- tools/schema
ba60c7bc 7 -n -F 'error_class' -- docs/user/benutzerhandbuch.md
ba60c7bc 14 -n -F 'error_class' -- '*.go' ':!*_test.go'
ba60c7bc 0 -n -F 'error_code' -- '*.go' '*.sql' '*.yaml' '*.proto' ':!*_test.go'
ba60c7bc 35 -n -E '\.Warn\(' -- internal cmd ':!*_test.go'
```

| Träger | Messung am Parent (`ba60c7bc`, 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| Schema-Träger der Spalte | Zeile 1: 7 | Spalte `error_code` daneben; `error_class` unverändert (Zeile 1 bleibt ≥ 7) |
| Handbuch | Zeile 2: 7 | Beschreibung ergänzt |
| Go-Träger von `error_class` | Zeile 3: 14 | `error_code` daneben |
| `error_code` heute | Zeile 4: 0 | Nichtgefunden; Diff: nur die genannten Träger |
| Warn-Aufrufe | Zeile 5: 35 | Implementer klassifiziert: Betreiber-Warnung (Code) / andere (kein Code) |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt, [T2](slice-meldungscodes-registry-fehlerkopf.md) liegt in `done/`.

**Rückführung:** `in-progress` → `next` (zu groß), wenn Warn-Codes und Spalte/Diagnose zusammen
nicht in einen Diff passen: Schnitt `W`-Codes gegen `error_code` (Spalte, Proto, Rollout);
`in-progress` → `open` (blockiert), wenn die Leser-Messung (C) einen brechenden Leser zeigt —
dann Auftraggeber-Frage, nicht stille Anpassung.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration`, `make test-store` grün nach `make image`, `make generated-sync`
grün, realer `diagnose`-Lauf mit Handbuch verglichen, Suchlauf-Block nachgemessen,
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **Schema-Rollout mit View-Signatur-Änderung.** Die Spalte ändert die Signatur von
  `cdc.heartbeat`; ohne Vorlauf scheitert der Rollout gegen migrierte Ziele
  ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)). Gegenmittel:
  Rollout-Wache-Lauf und zweiter Rollout gegen migriertes Ziel. — **Ausgang:** (bei Closure)
- **Proto-/Leser-Kompatibilität.** Ein additives Proto-Feld ist für Leser unschädlich
  (*hergeleitet*, nicht nachgemessen); Beispiel-Clients und Wegwerf-Clients lesen das
  Diagnose-Ergebnis. Gegenmittel: Liefer-Punkt C. — **Ausgang:** (bei Closure)
- **Zu viele oder zu wenige Warn-Codes.** Die Auswahl der 35 Warn-Stellen ist Urteil;
  Gegenmittel: Regel „Maßnahme des Betreibers“, Auswahl im Bericht, Reviewer liest sie. —
  **Ausgang:** (bei Closure)
- **Läufer-Erwartungen an die Diagnose-Zeile** werden erst im `make test-integration`-Lauf rot,
  der nicht in `make gates` liegt. Gegenmittel: Suchlauf, vollständiger Lauf vor Closure. —
  **Ausgang:** (bei Closure)
- **Kollision mit parallelen Arbeiten am Handbuch.** Rebase auf `main`, nur eigene Abschnitte
  committen. — **Ausgang:** (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird die Default-Sub-Area `*` (Kürzel
`PGC`, Modus Greenfield, [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); eine feinere Aufteilung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register `../observations/BEO-PGC/` ist
nicht Inhalt dieses Plans; der Implementer sichtet es vor dem Start (nicht gemessen, kein
Treffer behauptet).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
