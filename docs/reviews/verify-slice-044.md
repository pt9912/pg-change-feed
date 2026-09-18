# Verifikationsbericht: slice-044 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-044` §1/§2 DoD/§3 Plan-Nachzug/§4/§6/§8), `welle-13` §6 und
`ADR-0014`/`ADR-0047`/`ADR-0048`/`ADR-0053`, nicht gegen Diff (Reviewer-
Aufgabe, bereits abgeschlossen: das Review zu `slice-044` und der
Review-Report zur Fixrunde von `slice-044`, beide vollständig gelesen) und nicht gegen realen Bedarf
(Validator, hier nicht ausgelöst — kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Diese Prüfung liest den vollständigen, aktuellen
Slice-Plan, `ADR-0053`, den Architect-Verdikt
zu den Rollen-Grants in `slice-044`, `welle-13.md` §6, beide
Review-Reports und den tatsächlichen Code selbst — keine Behauptung aus
einem Bericht wird ungeprüft übernommen; jeder unten genannte Sensor-Lauf
wurde in dieser Sitzung **selbst** ausgeführt, nicht aus den Reports
zitiert.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-044-retention-hintergrundjob.md`
zum Stand `HEAD = b2255d5`. Commits (chronologisch): `f8949a4`
(`next→in-progress`), `33e2d35` (Implementierung), `824e001`
(Grant-Erweiterung `cdc_admin` DELETE), `48a34a2` (realer E2E-Beleg),
`fd28a61` (Benutzerhandbuch), `5a4c71e` (Image-Digest), `2a4ff05`
(DoD/Plan-Nachzug), `d59f670` (Kommentar-Korrektur), `3a2f457`
(Architect-Verdikt Rollen-Grant, Folge-`ADR-0053`), `69a230e` (Review:
1 HIGH — F-1 Slice-Chronik-Kommentar), `e2f4a5c` (Register 3×), `e800d9d`
(Fixrunde F-1 behoben), `bd78dc6` (Architect-Verdikt Steering-Loop-
Verkörperung), `ecf2deb`/`b2255d5` (Review-Bestätigung, Link-Fix). Sequenz
selbst geprüft (`git log --oneline f8949a4^..b2255d5`): Move →
Implementierung → Fund + Grant-Fix → E2E-Beleg → Doku → Plan-Nachzug →
Architect-Zug (Rollen-Konflikt) → Review (1 HIGH) → Architect-Zug
(3×-Verkörperung) → Fixrunde → Review-Bestätigung — keine Rolle springt
rückwärts ohne Übergabe-Artefakt.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `runRetentionCleanup`-Hintergrundzug verdrahtet, Muster identisch zu `runHeartbeat`/`runWALRetentionCheck`/`runAdministration`, ruft periodisch `RunRetentionUseCase` real auf | **erfüllt** | `internal/bootstrap/wiring.go` selbst gelesen: eigener Pool je Aufgabe (`retentionStore`, `retentionConsumerState`, beide `cfg.AdminDSN`), eigener `context.WithCancel` (`retentionCtx`/`stopRetention`), eigene `sync.WaitGroup` (`retentionDone`), Goroutine-Start und Shutdown-Reihenfolge (`stopRetention(); retentionDone.Wait()`) exakt an derselben Stelle wie die drei Vorbild-Züge. `runRetentionCleanup` (Zeile 774ff.) ruft `useCase.Run` in einer `for`-Schleife mit `interval`-Ticker auf, Best-effort-Fehlerbehandlung (Log + `continue`) spiegelt `runWALRetentionCheck`. `retentionInterval = 10s`, `retentionMinAge = 24h` als unexportierte Konstanten im Stil von `heartbeatInterval` — §1 schließt Laufzeit-Konfiguration bereits aus |
| 2 | Realer E2E-Beleg (`run-integration-tests.sh`): freigegebene Zeile real entfernt, nicht freigegebene (zu jung ODER Consumer hängt zurück) bleibt real erhalten, beides ohne Neustart | **erfüllt, dreifach selbst reproduziert** | `tools/harness/run-integration-tests.sh` gelesen: `RetentionOld` (id=200) real auf `committed_at = now() - 25h` zurückdatiert (direkter `UPDATE`, kein Warten — dieselbe Technik wie der bestehende Heartbeat-Fehlerzustand-Beleg), `RetentionYoung` (id=201) bleibt real jung. Beide Consumer (`CLI_CONSUMER`, `BACKLOG_CONSUMER`) blockieren strukturell zunächst, ein erster Poll belegt reale Nichtlöschung trotz erfülltem Alter (`LH-FA-RET-004`), nach Bestätigung beider Consumer belegt ein zweiter Poll reale Löschung von `RetentionOld` und reales Erhaltenbleiben von `RetentionYoung` (`LH-FA-RET-003`). **Selbst dreimal in Folge ausgeführt** (`make test-integration`, je ca. 1m17s–1m26s): alle drei Läufe enden mit identischer Beleg-Zeile „`RetentionOld` … blieb erhalten, solange ein Consumer zurückhing … und wurde nach Freigabe … real entfernt; `RetentionYoung` … blieb durchgehend erhalten“ und mit `PASS`/`ok`; der Feed-Container läuft über den gesamten Testlauf durch (kein `docker restart` im Skript, `docker exec` läuft gegen denselben laufenden Container wie beim Black-Box-CLI-Rundlauf davor) |
| 3 | `make gates` grün, `make test-integration` grün (inkl. des neuen Belegs) | **erfüllt, selbst reproduziert** | `make gates` selbst ausgeführt gegen `HEAD = b2255d5`: `baseline-verify` (v6.5.0, 54 Dateien) OK, `d-check` Struktur (356 Dateien, 0 Befunde), `d-check` Commits (`HEAD~5..HEAD`, 0 Befunde), `commit-traceability.sh` (5 Commits, Betreffs ohne Struktur-ID) OK, `a-check` (0 Befunde). `make test-integration` dreimal in Folge grün (siehe Punkt 2). Zusätzlich `make test -race` selbst ausgeführt: alle Pakete `ok`, insbesondere `internal/bootstrap` (1.041s) und das neue `internal/adapters/driven/systemclock` (1.016s) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **inhaltlich erfüllt, Formular-Diskrepanz — siehe Finding V-1** | Beide Reports vollständig gelesen: das Review zu `slice-044` (1 HIGH, F-1 Slice-Chronik-Kommentar, drittes Auftreten) und der Review-Report zur Fixrunde von `slice-044` (F-1 bestätigt behoben, Architect-Zug `bd78dc6` geprüft und für nachvollziehbar befunden, 0 HIGH/MEDIUM/LOW/INFO im Endstand). Die Bedingung ist damit tatsächlich erfüllt. Checkbox in §2 (Zeile 92) steht aber weiterhin auf `- [ ]` — kein Commit von `2a4ff05` bis `b2255d5` hat sie auf `[x]` gesetzt (`git log -p f8949a4..HEAD` gegen die Plan-Datei geprüft: nur die Punkte 1/2/3/6 wurden je auf `[x]` gesetzt) |
| 5 | Doku-Update `docs/user/benutzerhandbuch.md`, falls Betriebs-Aspekt entsteht | **erfüllt** | Abschnitt „Aufbewahrung (Retention)“ (Zeile 331ff.) selbst gelesen: Takt (10s) und Mindestalter (24h) stimmen mit `wiring.go`s Konstanten überein, Consumer-Abwesenheits-Lesart korrekt wiedergegeben; `cdc_admin`-Zeile der Rollen-Tabelle (Zeile 76) nennt „Retention-Löschausführung“ als neuen Verwaltungszugriff |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit (Modul 8); §7 ist noch die Bedienhinweis-Vorlage. Kein Verifikations-Gegenstand dieser Prüfung |
| 7 | Reconciliation-Register fortgeschrieben, falls einschlägig | **entfällt strukturell** | Kein `docs/plan/planning/reconciliation.md` — `harness/conventions.md` §Modus-Deklaration führt ausschließlich `*`/`PGC` im Modus Greenfield, kein Brownfield-Bootstrap |
| 8 | Beobachtungs-Register fortgeschrieben | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit. Zur Einordnung selbst geprüft: `BEO-PGC/architect-verdikt-rollen-scope-luecke` existiert bereits (1×, `evidence/slice-044.md`, `state.md` bestätigt „offen — kein Träger bislang … unter der 3×-Schwelle“); `BEO-PGC/slice-chronik-in-code-kommentar` bereits mit Ausgang *verkörpert* versehen (siehe unten, wellenloser Architect-Zug, nicht dieser Slice-Closure zugeschrieben) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit; beide Risiken in §6 tragen noch `<bei Closure einzutragen>`. Plan-Nachzug §3 Punkt 3 liefert die inhaltliche Grundlage für Risiko 2 (Zeitsteuerung über reales `UPDATE` statt Warten gelöst) |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit, wellenlos hier bzw. bei `welle-13`-Closure fällig, ohnehin erst **nach** dem `git mv` nach `done/` sinnvoll prüfbar |

## 2. Finding V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen

- **Klasse:** Verifier-only — für Tests und Review unsichtbar, weil beide
  Rollen ihre eigene Arbeit erledigt haben; nur ein Blick auf den
  *Formular-Zustand nach* der vollständigen Review-Sequenz (Erstlauf +
  Fixrunde) deckt die Lücke auf. Dieselbe Klasse wie `V-1` im
  Verifikationsbericht zu `slice-043` — bereits ein zweites Auftreten in dieser Repo-
  Historie (Steering-Loop-Signal: 1× notieren · 2× Symptom, noch keine
  Lücke).
- **Befund:** DoD-Punkt 4 in `slice-044-retention-hintergrundjob.md:92`
  steht auf `- [ ]`, obwohl die Bedingung — Review durchgeführt, Report
  liegt vor — seit `69a230e` faktisch erfüllt und seit `ecf2deb` zusätzlich
  mit 0 offenen Findings abgeschlossen ist.
- **Einordnung:** kein inhaltlicher Mangel an Review-Substanz — beide
  Durchläufe sind, wie oben belegt, real und reproduzierbar erfüllt. Es
  ist eine Diskrepanz zwischen dem DoD-Formular und der tatsächlichen
  Sachlage.
- **Erwartete Korrektur:** Checkbox auf `[x]` setzen, in einem eigenen,
  kleinen Commit vor dem `git mv` nach `done/` (Inhalt vor Move, Modul 5
  §git mv + Inhaltsänderung). Kein Rollback, keine Rückführung.

## 3. Rollen-Konflikt Grant-Erweiterung — eigenständig geprüft

- **`welle-13` §6:** Der nachgetragene Absatz ist ehrliche Nachdokumentation
  — er behauptet nicht rückwirkend, der Architect-Zug sei vorab gelaufen,
  sondern benennt „real eingetreten während `slice-044`“ und verweist auf
  den nachgetragenen Architect-Zug. Selbst gelesen, keine Beschönigung
  gefunden.
- **`ADR-0053`:** vollständig gelesen (Kontext, Entscheidung, drei
  verglichene Alternativen inkl. „nichts tun“, Konsequenzen inkl. der
  negativen Blast-Radius-Konsequenz, Re-Evaluierungs-Trigger, Geschichte).
  `Supersedes ADR-0047` ist explizit auf die Rollen-Zuordnungstabelle
  begrenzt — dieselbe Form wie `ADR-0048`.
- **`ADR-0047` unverändert (Hard Rule 3.5):** `git diff f8949a4 HEAD --
  docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` selbst
  ausgeführt — leer. `git log --oneline --all -- 0047-*.md` zeigt nur die
  Erstanlage und einen vor `slice-044` liegenden, rein pfadbezogenen
  Link-Fix (`f8bbfb5`, geprüft: ändert nur einen relativen Link, keine
  inhaltliche Aussage) — kein Commit dieses Slice berührt die Datei.
- **ADR-Index:** `docs/plan/adr/README.md` selbst geprüft — `ADR-0053`
  eingetragen, `ADR-0047`s Zeile trägt den `→ ADR-0053`-Zeiger.
- **Architect-Verdikt (zu den Rollen-Grants in `slice-044`):**
  vollständig gelesen. Verdikt 2 (Folge-ADR) wird nachvollziehbar
  begründet und von Verdikt 1/3 sauber abgegrenzt (Verdikt 1 trägt nicht,
  weil `welle-13` §6 die Bedingung korrekt offen formuliert hatte, keine
  falsche Behauptung vorlag; Verdikt 3 trägt nicht, weil der
  strukturell schwächere `ADR-0048`-Präzedenzfall bereits eine Folge-ADR
  verlangte). Kein Rückbau von Commit `824e001` — inhaltlich korrekt,
  nur der formale Architect-Zug fehlte.
- **Grant-Scope real geprüft:** `tools/schema/nacharbeit-roles.sql` selbst
  gelesen — `GRANT SELECT, DELETE ON cdc.transaction, cdc.change TO
  cdc_admin;` betrifft ausschließlich diese Tabellen und ausschließlich
  `cdc_admin`; `cdc_capture` bleibt auf `INSERT` beschränkt
  (unverändertes `GRANT INSERT ON cdc.transaction, cdc.change TO
  cdc_capture;` in derselben Datei). Regressionstests
  (`TestCdcAdminRetentionDeleteChangesRequiresGrant`,
  `TestCdcWiringCallerRejectsWrongRoleAssignment`, neuer Fall) real im
  Code gefunden (`internal/bootstrap/roles_wiring_test.go:236,267,432,443`)
  und über `make test-integration`s Rollen-DSN-Verifikationsabschnitt
  indirekt mitbestätigt (Compose-Lauf zeigt reale Rollen-Grenzen-
  Durchsetzung an anderer Stelle desselben Musters).

**Verdikt zu diesem Punkt:** Der Rollen-Konflikt ist vollständig und
ordnungsgemäß über den Architect-Zug (Verdikt 2, `ADR-0053`) aufgelöst.
Keine offene Frage, keine stille Lockerung, keine ADR-Immutabilitäts-
Verletzung.

## 4. Steering-Loop-Verkörperung `BEO-PGC/slice-chronik-in-code-kommentar` (3×)

- **Register (`state.md`):** selbst gelesen — Zustand *verkörpert*, Ausgang
  *verkörpert*, drei Belege ausgewiesen
  (dem Beleg zum Review-Report-Commit zu `slice-041`, dem Beleg zum
  Review-Report-Commit zur Fixrunde von `slice-041`,
  dem Beleg zum Review-Report-Commit zu `slice-044`) — `ls .../evidence/` bestätigt exakt
  diese drei Dateien, kein viertes oder fehlendes Element. Zielort
  benannt (`implement-slice.md` Schritt 20), Herkunfts-Anker (der
  Architect-Zug selbst, wellenlos) korrekt als eine der drei zulässigen
  Ausgangs-Formen (Modul 6 §Das Beobachtungs-Register) erkennbar — keine
  vierte, erfundene Form.
- **`.claude/commands/implement-slice.md` Schritt 20:** selbst gelesen
  (Zeilen 147–172) — die Enumerations-Pflicht ist konkret und
  ausführbar verankert: ein diff-skopierter `grep`-Kandidatenlauf gegen
  die in diesem Lauf geänderten `.go`-/`tools/schema/*.sql`-Dateien plus
  eine explizite Testfall-Provenienz-vs.-Chronik-Probe je Treffer. Grenze
  (kein Sensor/Gate, bleibt Selbstprüf-Disziplin) ist ausdrücklich benannt,
  nicht verschwiegen.
- **Eigener, unabhängiger `grep`-Durchlauf** über alle seit `f8949a4`
  geänderten `.go`-Dateien
  (`internal/adapters/driven/systemclock/{systemclock,systemclock_test}.go`,
  `internal/bootstrap/{retention_internal_test,roles_wiring_test,wiring}.go`)
  gegen das Muster `slice-[0-9]+|vorher|nachher|jetzt|früher|neu
  hinzugefügt|wurde (entfernt|geändert|umgestellt)|Ohne (dieses|diese|
  diesen)`: **ein** Treffer, `wiring.go:643` (`` `slice-026` Fixrunde,
  Review F-1 ``). Selbst per `git log -L 643,643:internal/bootstrap/
  wiring.go` und `git diff f8949a4 HEAD -- internal/bootstrap/wiring.go`
  nachvollzogen: letzter ändernder Commit dieser Zeile ist `29a49ea`
  (2026-09-12, vor dem ersten `slice-044`-Commit) — die Zeile ist nicht
  Teil des `f8949a4..HEAD`-Diffs. Kein weiterer, nicht-historischer
  Treffer gefunden — deckt sich mit dem eigenen Befund beider
  Review-Reports.

**Verdikt zu diesem Punkt:** Die Verkörperung ist korrekt im Register und
im Implementer-Briefing verankert; eigener Grep-Durchlauf bestätigt den
Endzustand ohne weiteren Chronik-Fund.

## 5. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

- **Ausschluss 1 (Sichtbarkeit blockierender Consumer, `slice-045`):**
  `grep -rn "blockierend" internal --include="*.go"` trifft nur
  Dokumentationskommentare, die den Ausschluss selbst benennen (`retention.go`
  Inbound-Port, `retention_test.go`, `heartbeat_internal_test.go`), keinen
  Code, der Blockier-Sichtbarkeit implementiert. Eingehalten.
- **Ausschluss 2 (`cdc_storage_bytes`, `slice-046`):** `grep -rn
  "cdc_storage_bytes" internal tools/schema` trifft nur die bereits
  vorbestehende „nicht abgedeckt“-Zeile in
  `nacharbeit-observability.sql`. Eingehalten.
- **Ausschluss 3 (Laufzeit-Konfigurierbarkeit des Job-Takts):**
  `git diff f8949a4 HEAD -- internal/domain/model/retention.go` ist leer
  — `RetentionPolicy` unverändert; Takt/Mindestalter bleiben feste
  Konstanten in `wiring.go`. Eingehalten.
- **Ausschluss 4 (Black-Box-E2E über den Compose-Stack hinaus):** Der
  neue Beleg läuft ausschließlich im bestehenden
  `run-integration-tests.sh`-Skript gegen den Compose-Stack, keine neue
  Testwelle oder externe Testinfrastruktur hinzugefügt. Eingehalten.

## 6. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Wie im Auftrag benannt und durch §Träger im Repo ohne Wellen-Betrieb
(Modul 6/8) gedeckt: **Closure-Notiz** (DoD 6), **Beobachtungs-Register-
Fortschreibung für diesen Slice** (DoD 8), **Risiko-Ausgänge §6** (DoD 9),
**die drei Paarungen** (DoD 10) — alle vier zugehörigen §2-Häkchen sind
**korrekt unbeansprucht**, kein Mangel, sondern der vorgesehene Zustand vor
dem nächsten Rollenwechsel an den Planner. Das Reconciliation-Register-Item
(DoD 7) entfällt strukturell (Greenfield-Repo). Auch nicht Gegenstand:
Validierung gegen realen Bedarf (kein MVP-Meilenstein-Slice, kein
Validator-Zug ausgelöst).

## Verdikt

**DoD-Konformität: bestätigt**, mit einer benannten, nicht
merge-blockierenden Formular-Diskrepanz (V-1: Checkbox 4 „Review
durchgeführt“ nicht nachgezogen — inhaltlich längst erfüllt). Alle fünf
substanziellen Implementer-DoD-Punkte (1–5) sind durch eigene, unabhängige
Reproduktion gedeckt: `runRetentionCleanup` folgt real dem etablierten
Hintergrundzug-Muster, der E2E-Beleg zeigt reale Löschung/Nicht-Löschung
ohne Neustart in **drei** unabhängigen `make test-integration`-Läufen,
`make gates` und `make test -race` sind selbst grün gelaufen, das
Benutzerhandbuch ist sachlich korrekt fortgeschrieben.

**Rollen-Konflikt (`ADR-0053`):** vollständig und ordnungsgemäß über den
Architect-Zug aufgelöst. `ADR-0047` selbst bleibt unverändert (Hard Rule
3.5, `git diff` leer), `welle-13` §6 dokumentiert den Konflikt ehrlich als
real eingetreten statt ihn zu verschweigen, kein stiller Fortschritt.

**Steering-Loop-Verkörperung (`BEO-PGC/slice-chronik-in-code-kommentar`,
3×):** korrekt im Register (`state.md`, Ausgang *verkörpert*, drei
Belege) und in `.claude/commands/implement-slice.md` Schritt 20 verankert;
eigener Grep-Durchlauf über die gesamte seit `f8949a4` geänderte
`.go`-Dateimenge findet keinen weiteren, nicht-historischen Chronik-Fund.

**Die vier Planner-Closure-Punkte** (Closure-Notiz, Beobachtungs-Register-
Fortschreibung für diesen Slice, Risiken-Ausgänge, drei Paarungen) sind
**korrekt unbeansprucht** — bewusst nicht Teil dieser Prüfung.

**Keine Rückführung nötig.** Weder `in-progress→next` (der Slice ist nicht
zu groß — drei Liefer-Punkte aus §2 sauber erfüllt, der Grant-Fund wurde
über den Architect-Zug transparent nachgetragen statt eine vierte Schicht
zu öffnen) noch `in-progress→open` (kein Blocker). Vor dem `git mv` nach
`done/` ist lediglich die Checkbox-Korrektur aus V-1 fällig — ein
Ein-Zeilen-Commit, kein Zerlegungs- oder Blocker-Fall.

**Übergabe an Planner:** Dieser Bericht bestätigt DoD-Konformität für die
Closure-Entscheidung, mit dem Hinweis V-1 zur Nachbesserung vor dem
`git mv`. Kein Validator-Zug ausgelöst — `slice-044` ist kein
MVP-Meilenstein-Slice im Sinn von Modul 8.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
