# Slice slice-027: Black-Box-E2E-Test über die bestehende CLI

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-8`](welle-8.md) — der repo-weite Nachweis, dass ein
echter externer Aufrufer das System nur über seine reale Schnittstelle
bedienen kann, geht über die DoD dieses einzelnen Slice hinaus (welle-8 §3).

**Bezug:** [`LH-QA-POR-003`](../../../../spec/lastenheft.md)
(reproduzierbare Testumgebung) — dieser Slice erweitert die Testabdeckung
für einen bestehenden Vertrag, ändert ihn nicht. Kein aktives ADR wird
geändert; `ADR-0030` (Testpyramide) trägt die Unterscheidung
Integrationstest vs. E2E-Test, die dieser Slice umsetzt, ohne sie zu ändern.

**Berührte Spec-Stellen:** — (reine Testinfrastruktur, keine
Verhaltensänderung an einer Spec-Stelle).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Compose-Integrationstest
([`tools/harness/run-integration-tests.sh`](../../../../tools/harness/run-integration-tests.sh))
bekommt einen echten Black-Box-Testfall: `register-consumer` und
`acknowledge-consumer` werden ausschließlich als externer Prozess
(`docker exec <feed-container> pg-change-feed register-consumer …`) gegen
den laufenden, containerisierten Produktions-Binary aufgerufen — kein
Import interner Go-Pakete (`bootstrap.RegisterConsumer` o. ä.). Der Testfall
belegt real den vollen Consumer-Rundlauf: registrieren → Changes schreiben
→ lesen → bestätigen → Feed-Container neu starten (simulierter Neustart) →
Fortsetzen ab der bestätigten Position — ausschließlich über extern
ausgelieferte Schnittstellen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein neuer Lese-Zugriffsweg über die CLI** — Bestand bleibt bewusst
  stehen: Lesen läuft weiterhin über den bestehenden externen Lesezugriffsweg
  (SQL-Sicht `cdc.changes`, `LH-FA-REA-002`), da die CLI heute keinen
  Lese-Unterbefehl hat. Das ist eine ehrliche Grenze dieses Tests, keine
  versteckte Lücke: Nur `register-consumer`/`acknowledge-consumer` sind
  aktuell CLI-Unterbefehle.
- **Rollen-spezifische DSN-Trennung im Compose-Integrationstest verifizieren**
  — Folge-Slice `slice-028` übernimmt das; dieser Slice liefert nur das
  Black-Box-Aufrufmuster (CLI-Subprozess), auf dem `slice-028` aufbaut.
- **WAL-Rückstand-Schwellen-Verhalten zusätzlich gegen den Compose-Stack
  prüfen** — Bestand bleibt bewusst stehen: bereits real gegen PostgreSQL
  bewiesen (`slice-026`, welle-8 §6 Out-of-Scope).
- **Umbenennung von `make test-integration`/`mvp_test.go`** — Folge-Slice
  `slice-028` übernimmt das, zusammen mit den übrigen Compose-
  Integrationstest-Erweiterungen.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `LH-QA-POR-003` erfüllt: neuer Black-Box-Testfall ruft
      `register-consumer`/`acknowledge-consumer` real als externen Prozess
      (`docker exec` gegen den laufenden Feed-Container) auf — kein
      Go-Paket-Import interner Anwendungslogik in diesem Testfall. Beleg:
      neuer Abschnitt „Black-Box-CLI-Rundlauf" in
      [`tools/harness/run-integration-tests.sh`](../../../../tools/harness/run-integration-tests.sh)
      (`exec_feed()` ruft ausschließlich `docker exec … /pg-change-feed …`).
- [x] Voller Rundlauf real bewiesen: registrieren (extern) → Change
      schreiben (Quelltabelle) → lesen (bestehender SQL-Lesezugriffsweg) →
      bestätigen (extern) → Feed-Container-Neustart simulieren (analog zum
      bestehenden `welle6_endtoend_test.go`-Muster, aber mit dem
      *containerisierten* Prozess statt einem in-process-Aufruf) →
      Fortsetzen ab der bestätigten Position real gezeigt. Beleg: drei
      Folge-Läufe von `make test-integration`, je mit der Ausgabezeile
      „Black-Box-CLI-Rundlauf belegt — … Fortsetzen nach simuliertem
      Neustart …"; die zentrale Assertion (kein Wiederholen der bereits
      bestätigten Änderung nach dem Neustart) wurde einmal mit einer
      absichtlich mutierten Grenze (`>=` statt `>`) rot gesehen (Ausgabe
      „gelesene id-Folge '95,96', wollen '96'"), danach zurückgesetzt.
- [x] `make gates` grün, `make test-integration` dreimal in Folge grün
      (real gegen den Compose-Stack).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-027.md`](../../../reviews/review-slice-027.md)
      (0 HIGH/MEDIUM/LOW/INFO, kein Merge-Blocker), unabhängig bestätigt
      durch [`docs/reviews/verify-slice-027.md`](../../../reviews/verify-slice-027.md)
      (eigene Rot-Grün-Gegenprobe, 1× LOW V-1, siehe §7).
- [x] Doku-Update für `harness/README.md` (Sensors-Tabelle,
      `make test-integration`-Zeile) falls sich der geprüfte Umfang
      erkennbar ändert. Entscheidung: durchgeführt — der geprüfte Umfang
      ändert sich für einen Leser erkennbar (bislang rein interne
      Store-/Go-Assertions, jetzt zusätzlich ein echter externer
      CLI-Rundlauf über einen simulierten Container-Neustart).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register — **entfällt**: Repos ohne
      Brownfield-Bootstrap haben die Datei nicht (`docs/plan/planning/reconciliation.md`
      existiert in diesem Repo nicht, Sub-Area `*`/`PGC` ist Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      Keine Beobachtung angefallen.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Repo **mit** Wellen-Betrieb (`welle-8` offen) — verschoben auf die welle-8-Closure.

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | neuer Abschnitt „Black-Box-CLI-Rundlauf": `register-consumer`/`acknowledge-consumer` als `docker exec`-Aufrufe gegen `$FEED_CONTAINER` (Helfer-Funktion `exec_feed()`), Vorbedingungs-/Ergebnis-Prüfung per SQL gegen `cdc.consumer`/`cdc.consumer_position`/`cdc.changes`, `docker restart` für den simulierten Neustart — analog zum bestehenden Lasttest-Beleg-Abschnitt (Bash + Exit-Code-Prüfung) |
| `harness/README.md` | update | `make test-integration`-Zeile ergänzt: der geprüfte Umfang trägt jetzt sichtbar den externen CLI-Rundlauf, nicht nur interne Store-/Go-Assertions |

**Plan-Nachzug (Abweichungen vom ursprünglichen §3, seit dem ersten Implementer-Lauf):**

- **`test/integration/mvp_test.go` / neue Go-Datei — nicht realisiert.** Die
  Bash-Variante in `run-integration-tests.sh` reicht: Sie ruft die CLI
  bereits als externen Prozess auf, prüft Vor-/Nachbedingungen über
  dieselben `docker exec psql`-Aufrufe wie der Rest des Skripts und braucht
  keinen zusätzlichen Go-Testprozess. Ein Go-Test hätte selbst wieder
  `os/exec` gegen `docker exec` aufrufen müssen — mehr Indirektion ohne
  zusätzlichen Beleg.
- **`compose.yaml` — nicht geändert.** `docker restart` (SIGTERM + Neustart
  desselben Containers) braucht keine zusätzliche Compose-Konfiguration;
  Slot und Publication bleiben auf der PostgreSQL-Seite unberührt, ein
  eigener Neustart-Mechanismus im Container-Vertrag ist nicht nötig.
- **`harness/README.md` ergänzt** (war im ursprünglichen §3 nicht gelistet,
  aber im Plan-Ziel unter §2 DoD vorgesehen) — siehe Zeile oben.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  der simulierte Container-Neustart eine grundlegende Änderung an
  `compose.yaml`/der Test-Infrastruktur braucht (mehr als eine Schicht),
  gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
dreimal in Folge grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein simulierter Container-Neustart (`docker restart`/`docker stop && up`)
  könnte länger dauern oder sich anders verhalten als ein echter Prozess-
  Neustart und den Test flaky machen. **Ausgang: entfallen** — real
  widerlegt: drei Folge-Läufe von `make test-integration` waren grün, das
  Health-Gate (60×1s-Warteschleife auf `healthy`) trägt die Wartezeit
  strukturell ab, kein Flake beobachtet.
- Die Lese-Verifikation bleibt technisch „intern" (SQL-Query direkt gegen
  `cdc.changes`, kein CLI-Lese-Unterbefehl vorhanden) — der Test ist damit
  kein reiner Black-Box-Test, sondern black-box für Schreiben/Bestätigen und
  white-box fürs Lesen. **Ausgang: entfallen** — wie in §1 vorweggenommen,
  eine bewusste, benannte Grenze; Reviewer und Verifier bestätigten
  unabhängig, dass die *Aktionen* (Schreiben, Bestätigen) real extern
  laufen, nur die Lese-Verifikation nicht.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Der Bash-Ansatz (Erweiterung von
  `run-integration-tests.sh` statt eines neuen Go-Tests) hielt die
  Indirektion niedrig — `exec_feed()` ruft die CLI genauso auf, wie ein
  echter Betreiber es täte, ohne einen zusätzlichen Testprozess-Layer.
  Die eigene Rot-Grün-Verifikation zog sich wie in den vorigen Slices durch
  alle drei Rollen (Implementer, Reviewer, Verifier je einmal unabhängig).
- **Was ging anders als geplant:** Kein neuer Go-Test/keine `compose.yaml`-
  Änderung nötig (Plan-Nachzug §3) — die Bash-Erweiterung reichte für den
  vollen Rundlauf inkl. echtem Container-Neustart.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× und
  wird verkörpert — der Normalfall. Die vom Verifier notierte
  DoD-Checkbox-Lücke (V-1) ist strukturell dieselbe wie bei
  `slice-024`/`025`/`026`: Die Review-Zeile kann erst nach dem Review
  getickt werden.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen.
- **Folge-Slices:** `slice-028` (Compose-Integrationstest-Nachzug —
  Rollen-DSN-Trennung, MVP-Umbenennung) — bereits eine Datei in `open/`
  (welle-8 §4).
- **Risiken aus §6:** beide *entfallen* (Neustart real dreimal flake-frei
  reproduziert; White-Box-Lesegrenze bereits in §1 bewusst benannt und von
  Reviewer/Verifier bestätigt).
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-8` offen) —
  verschoben auf die welle-8-Closure (Modul 8 §Rollen-Sequenz für eine
  Welle).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC`: `BEO-PGC/rollen-test-abdeckungsluecken` (1×, nicht
einschlägig für diesen Slice — betrifft `slice-028`),
`BEO-PGC/test-isolation-geteilter-zustand` (1×, nicht einschlägig — anderer
Test-Layer, siehe welle-8 §6 Out-of-Scope). Keiner der Treffer erreicht mit
diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
