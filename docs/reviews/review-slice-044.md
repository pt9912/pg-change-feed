# Review-Report: slice-044 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-044`, §1/§2/§3/§4/§6/§8),
`welle-13` (inkl. des nachgetragenen §6-Absatzes), den Architect-Verdikt
zu den Rollen-Grants in `slice-044`, `ADR-0053` und
`ADR-0040`/`ADR-0047`/`ADR-0048` sowie `AGENTS.md` §3 Hard Rules (§3.1,
§3.5, §3.7) — Rollentrennung Modul 8: diese Prüfung läuft gegen
Plan/ADR/Hard Rules (Maintainability), nicht gegen DoD (Verifier-Aufgabe).

**Gegenstand:** Commits `33e2d35` (`feat(bootstrap): Retention-
Hintergrundzug ruft RunRetentionUseCase real auf`), `824e001`
(`fix(schema): cdc_admin trägt DELETE auf cdc.transaction/change`),
`48a34a2` (`test(integration): realer Retention-Beleg im Compose-
Testlauf`), `fd28a61` (`docs(user): Betriebsabschnitt „Aufbewahrung
(Retention)"`), `5a4c71e` (`chore(image): Image-Digest neu`), `2a4ff05`
(`docs(planning): slice-044 DoD/Plan-Nachzug und neue Beobachtung`),
`d59f670` (`docs(bootstrap): Kommentar-Formulierung an Ist-Zustand-
Disziplin angepasst`), `3a2f457` (`docs(adr): Architect-Verdikt slice-044
Rollen-Grant — Folge-ADR-0053`) — neu: `internal/adapters/driven/
systemclock/{systemclock,systemclock_test}.go`, `internal/bootstrap/
retention_internal_test.go`, `docs/plan/adr/0053-retention-
loeschausfuehrung-cdc-admin-delete-grant.md`, der Architect-Verdikt
zu den Rollen-Grants in `slice-044`, zwei Beobachtungs-Register-Dateien
(`BEO-PGC/architect-verdikt-rollen-scope-luecke/`); geändert:
`internal/bootstrap/wiring.go`, `internal/bootstrap/roles_wiring_test.go`,
`tools/schema/nacharbeit-roles.sql`, `tools/harness/run-integration-
tests.sh`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/welle-13.md`,
`docs/plan/adr/README.md`, `docs/plan/planning/in-progress/slice-044-
retention-hintergrundjob.md` (DoD-Häkchen, Plan-Nachzug), `harness/image-
hash.txt`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft
2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-044-retention-hintergrundjob.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug, §4
  Trigger, §6 Risiken, §8)
- `docs/plan/planning/welle-13.md` (vollständig, inkl. des nachgetragenen
  §6-Absatzes zum Rollen-Konflikt)
- der Architect-Verdikt zu den Rollen-Grants in `slice-044` (vollständig
  — Konflikt-Pfad-Verdikt, Modul 8)
- `docs/plan/adr/0053-retention-loeschausfuehrung-cdc-admin-delete-
  grant.md` (vollständig — Kontext, Entscheidung, Verglichene
  Alternativen, Konsequenzen, Re-Evaluierungs-Trigger, Geschichte)
- `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md`,
  `0048-heartbeat-grant-korrektur-select-ergaenzung.md` (Präzedenzfall für
  Verdikt 2/Folge-ADR-Form) — geprüft, dass `ADR-0047`s Datei selbst nicht
  verändert wurde (Hard Rule 3.5, `git diff` gegen den Dateiinhalt: leer)
- `docs/plan/planning/done/slice-043-changestoreport-loeschmethode.md`
  (vollständig — Grundlage, auf der dieser Slice aufbaut)
- `AGENTS.md` §3.1 (Docker-only), §3.5 (ADR-Immutabilität), §3.7
  (Kommentar-Disziplin — `BEO-PGC/slice-chronik-in-code-kommentar`, 2× vor
  diesem Lauf laut Auftrag, besonders geprüft)
- `internal/bootstrap/wiring.go` — bestehende Hintergrundzüge
  (`runHeartbeat`, `runWALRetentionCheck`, `runAdministration`) als
  Vorbild-Muster, vollständig gelesen (Pool-Aufbau, Goroutine-Start,
  Shutdown-Reihenfolge)
- `docs/plan/adr/0040-*.md` (`ClockPort`/`SystemClockAdapter`-
  Folgepflicht) — `grep -rn "ClockPort"` bestätigt: vor `33e2d35` existierte
  keine Produktionsimplementierung, nur der Port und Fake-Test-Doubles
- Vollständiger `git show` für `33e2d35`, `824e001`, `48a34a2`, `3a2f457`
  (alle Dateien, nicht nur die Implementer-Zusammenfassung); `fd28a61`,
  `5a4c71e`, `2a4ff05`, `d59f670` ergänzend gelesen
- Review zu `slice-043` (Format-Vorlage)
- `grep -rn "slice-0[0-9]"` über alle in diesem Diff neuen/geänderten
  `.go`-Dateien (vollständiger Durchlauf, nicht stichprobenartig)

---

## Findings

### F-1 — Slice-Chronik in Go-Quellcode-Kommentar (drittes Auftreten)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Disziplin), Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- `pfad`: `internal/bootstrap/retention_internal_test.go:14–16`
  (Datei-Kopf-Kommentar über `fakeRunRetentionUseCase`)
- `befund`: Der Kommentar lautet: „Whitebox-Test (`package bootstrap`,
  nicht `bootstrap_test`): der periodische Auslöse-Zug ist ein
  unexportiertes Verdrahtungsdetail (`slice-044`, `ADR-0014`) — der reale
  Ende-zu-Ende-Beleg … liegt in `tools/harness/run-integration-
  tests.sh` …". Der Kommentar nennt explizit die Slice-Kennung
  `slice-044` neben der ADR als Beleg-Quelle einer Aussage über den
  Code — dasselbe Muster, das in diesem Repo bereits zweimal HIGH
  eingestuft wurde (Review zu `slice-041` F-1, bestätigt im
  Review-Report zur Fixrunde von `slice-041`; beide als Evidence in
  `BEO-PGC/slice-chronik-in-code-kommentar/evidence/` geführt, Zähler vor
  diesem Lauf bei 2×). Ein Slice-Plan ist kein dauerhaftes Artefakt — er
  wandert bei Welle-Closure nach `done/` und kann später archiviert
  werden (gekürzter Stub); der Verweis referenziert dann eine Adresse,
  deren Inhalt sich ändert, während der Code bestehen bleibt. `ADR-0014`
  allein wäre der stabile, kanonische Verweis gewesen — die
  Slice-Kennung trägt keine der fünf zulässigen Kommentar-Klassen (Zusage
  · Kopplung · Abgrenzung · Rang-Zeiger · Grenze), sondern Chronik.
  **Dieser Fund ist das dritte registrierte Auftreten der Klasse** — bei
  Closure erreicht `BEO-PGC/slice-chronik-in-code-kommentar` damit 3×
  (Modul 5 §Closure- und Lerneintrag-Regeln: die Summary-Zeile dieses
  Review-Reports speist den Zähler als dritte Quelle) und braucht bei der
  Slice-Closure einen Ausgang — im Regelfall *verkörpert*: eine
  mechanische Prüfung (z. B. ein `grep`-Schritt gegen
  `slice-[0-9]{3}`-Muster in geänderten `.go`-Dateien vor dem Commit,
  wie in `BEO-PGC/slice-chronik-in-code-kommentar/state.md` bereits als
  Vorschlag geführt), nicht nur ein weiterer manueller Fixrunden-Commit.
- `verifizierbar`: ja — `grep -n "slice-[0-9]" internal/bootstrap/*.go
  internal/adapters/driven/systemclock/*.go` (in dieser Sitzung
  ausgeführt, ein Treffer außerhalb bereits vorbestehender, älterer
  Vorkommen in `wiring.go:643`, `walretention_internal_test.go`,
  `welle6_endtoend_test.go` — letztere drei datieren vor der
  Registrierung/Schärfung der Regel am 2026-09-13 und sind historischer
  Bestand, kein Gegenstand dieses Diffs)
- `klasse`: „Slice-Chronik in Go-Quellcode-Kommentar“ (3×)

## Geprüft, ohne (weiteren) Befund

- **Rollen-Konflikt-Auflösung (Architect-Verdikt 2, `ADR-0053`):**
  inhaltlich überzeugend. Verdikt 1 (Plan hat ADR missverstanden) trägt
  nicht, weil `welle-13` §6 die Bedingung korrekt offen formuliert hatte
  ("träfe das … zu, wäre das ein eigener Architect-Zug") — keine falsche
  Behauptung, die zurückzuweisen wäre. Verdikt 3 (Bestätigung ohne neue
  ADR) trägt nicht, weil der Präzedenzfall `ADR-0048` für den
  strukturell schwächeren Heartbeat-Fall bereits eine Folge-ADR verlangte
  — hier fehlte nicht nur ein Grant-Text, sondern die
  Rollen-*Zuordnung* selbst vollständig, und die verworfene Alternative
  (eigene vierte Rolle) berührt `ADR-0047`s eigenen
  Re-Evaluierungs-Trigger (a). Der `welle-13` §6-Nachtrag ist ehrliche
  Nachdokumentation — er behauptet nicht rückwirkend, der Architect-Zug
  sei vorab gelaufen, sondern benennt explizit "real eingetreten während
  `slice-044`" und dass der fällige Zug nachgetragen wurde. `ADR-0053`s
  `Supersedes ADR-0047`-Beziehung ist sauber begrenzt: `git diff` gegen
  `0047-rollenspezifische-dsn-verdrahtung.md` selbst ist leer (Hard Rule
  3.5 eingehalten), der Index-Eintrag wurde um den Zeiger ergänzt, und
  `ADR-0053`s Kontext/Entscheidung benennt ausdrücklich, dass der
  Drei-DSN-Vertrag und alle übrigen Zuordnungen unverändert bleiben —
  exakt dieselbe begrenzte Form wie `ADR-0048`.
- **`runRetentionCleanup`-Hintergrundzug:** folgt dem etablierten Muster.
  Eigener Pool je Rolle/Aufgabe (`retentionStore`,
  `retentionConsumerState`, beide `cfg.AdminDSN`), eigener
  `context.WithCancel`, eigene `sync.WaitGroup`, `stopRetention()` +
  `retentionDone.Wait()` an derselben Stelle im Shutdown wie
  `stopHeartbeat`/`stopWALRetention`/`stopAdministration`. Best-effort-
  Fehlerbehandlung (Log + `continue`) spiegelt `runWALRetentionCheck`.
  `retentionInterval = 10s` (seltener als der 5s-Heartbeat-Takt,
  begründet mit der breiteren Leseoperation) und `retentionMinAge = 24h`
  sind als unexportierte Konstanten im Stil von `heartbeatInterval`
  verdrahtet, `§1` schließt eine Laufzeit-Konfiguration bereits
  ausdrücklich aus.
- **`SystemClockAdapter`:** Die Implementer-Behauptung stimmt —
  `grep -rn "ClockPort"` traf vor diesem Commit nur den Port selbst und
  Test-Doubles, keine Produktionsimplementierung; `ADR-0040` benennt den
  Adapter bereits explizit als offene Folgepflicht. Die Implementierung
  ist minimal und korrekt: zustandsloser Wert-Typ, `Now()` liest
  ausschließlich `time.Now()`, keine Testmocks oder Slice-Referenzen im
  Produktionspfad. Drei Unit-Tests belegen Port-Erfüllung, Wanduhr-Treue
  und Monotonie ohne reale PostgreSQL-Instanz.
- **E2E-Beleg (direktes `UPDATE … committed_at`):** legitimer Testbeleg,
  keine Verwässerung des DoD-Anspruchs. Das Muster ist identisch zum
  bereits etablierten Fehlerzustand-Beleg (`cdc.process_heartbeat` direkt
  geschrieben, `tools/harness/run-integration-tests.sh` Zeile ~635): eine
  reale Spalte wird auf einen Zustand gesetzt, den ein realer Ablauf
  irgendwann real erreichen würde, statt 24 Stunden real verstreichen zu
  lassen — der Testlauf bleibt damit unabhängig vom konkreten Wert von
  `retentionMinAge` deterministisch und schnell, ohne die geprüfte
  Eigenschaft (`AllowsDeletion`s Alters-Freigabe) zu verändern: Die
  Löschung selbst erfolgt weiterhin real über den laufenden
  `runRetentionCleanup`-Takt, nicht simuliert. Zusätzlich real
  demonstriert (stärker als die reine DoD-Formulierung verlangt): der
  Consumer-Block wird über eine feste, mehrere Takte überspannende
  Wartezeit (`sleep 25`) belegt, bevor beide Consumer bestätigen — beide
  Freigabe-Bedingungen der Policy sind damit unterscheidbar und real
  geprüft, nicht nur eine.
- **Grant-Erweiterung (`tools/schema/nacharbeit-roles.sql`):** korrekt
  gescopt — `GRANT SELECT, DELETE ON cdc.transaction, cdc.change TO
  cdc_admin;` betrifft ausschließlich die beiden genannten Tabellen und
  ausschließlich `cdc_admin`; `cdc_capture` bleibt unverändert auf
  `INSERT` beschränkt. Real regressionsgetestet in beide Richtungen
  (`TestCdcAdminRetentionDeleteChangesRequiresGrant`:
  `SQLSTATE 42501` ohne Grant, Erfolg mit Grant;
  `TestCdcWiringCallerRejectsWrongRoleAssignment` neuer Fall: `DELETE`
  scheitert real unter `cdc_capture`-Login).
- **Kommentar-Disziplin (`grep`-Durchlauf, sonst):** `internal/bootstrap/
  wiring.go`s neue Retention-Abschnitte sind sauber — `d59f670` hat die
  beiden einzigen anderen Chronik-Formulierungen in diesem Diff
  („Vorher/Nachher“, „jetzt“ in `roles_wiring_test.go` und
  `run-integration-tests.sh`) bereits vor diesem Review-Lauf korrigiert;
  beide Ersetzungen sind reine Ist-Zustand-Beschreibungen ohne
  inhaltliche Änderung. Kein weiterer Treffer außerhalb von F-1.
- **`docs/user/benutzerhandbuch.md` — Abschnitt „Aufbewahrung
  (Retention)“:** sachlich korrekt und vollständig. Takt (10s) und
  Mindestalter (24h) stimmen mit `wiring.go`s Konstanten überein, die
  Consumer-Abwesenheits-Lesart ist korrekt und konsistent mit
  `slice-043`s Plan-Nachzug Punkt 8 wiedergegeben (ein Consumer zählt
  erst ab seiner ersten Bestätigung als schützenswert), die
  `cdc_admin`-Zeile der Rollen-Tabelle nennt die Retention-
  Löschausführung als neuen Verwaltungszugriff.
- **`docs/plan/adr/README.md`:** Index-Zeile für `ADR-0053` ergänzt,
  `ADR-0047`s Zeile um den `→ ADR-0053`-Zeiger erweitert (Präzedenzform
  identisch zu `ADR-0048`s Eintrag).
- **Traceability:** alle acht Commit-Betreffe tragen mindestens eine
  `LH-*`- oder `ADR-*`-Kennung, keiner trägt eine `SPEC-*`/`ARC-*`-
  Struktur-ID im Betreff.
- **Docker-only (§3.1):** keine lokale Toolchain-Installation in diesem
  Diff.
- **Zwei-Quellen-Drift:** keine gefunden — `welle-13` §6 und `ADR-0053`
  tragen dieselbe Aussage konsistent, keiner widerspricht dem anderen.

## Zusammenfassung

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Slice-Chronik in Go-Quellcode-
Kommentar“ (1×, drittes Gesamt-Auftreten der Klasse — Zähler-Übertritt
3× fällig bei Closure).

## Verdikt

**Merge-blockierend:** nein im Sinne von „Commit zurückrollen“ — der
betroffene Kommentar ist eine reine Doc-Kommentar-Formulierung ohne
Verhaltensänderung, und der Implementer hat in `d59f670` bereits gezeigt,
dass er dieselbe Korrektur-Disziplin anwendet, sobald sie benannt ist.
**Aber Closure-blockierend:** F-1 ist HIGH und muss vor `git mv` nach
`done/` behoben sein (Kommentar auf Ist-Zustand + `ADR-0014`-Verweis
umformulieren, `slice-044` streichen) — dieselbe Behandlung wie beim
Review zu `slice-041` F-1. Zusätzlich: Da dies das **dritte** Auftreten
ist, verlangt Modul 6 bei der Slice-Closure einen Ausgang für
`BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert/geplant/gestrichen)
statt eines weiteren offenen Zählers — die Closure-Notiz sollte diesen
Übertritt explizit als Steering-Loop-Eintrag mit Zielort führen, nicht
nur als weiteren Fixrunden-Commit.

Kein Rollen-Widerspruch: F-1 ist ein reiner Kommentar-Formfehler ohne
Implementer-Gegenposition zu erwarten, keine Architect-Sequenz nach
Modul 8 erforderlich. Der Architect-Verdikt zum Rollen-Grant-Konflikt
selbst (separates Thema, `ADR-0053`) wird bestätigt und nicht erneut zur
Disposition gestellt.

**Übergabe:** ein HIGH-Finding — Implementer behebt den Kommentar in
`internal/bootstrap/retention_internal_test.go` in einer Fixrunde und
trägt den fälligen 3×-Ausgang für `BEO-PGC/slice-chronik-in-code-
kommentar` in der Closure-Notiz nach, bevor `slice-044` nach `done/`
wandert. Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg
nicht wieder gelesen; die Summary-Zeile speist den Closure-Eintrag
(Modul 5).
