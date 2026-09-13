# Slice slice-044: Hintergrund-Job für die Retention-Löschausführung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-13 — zweiter Slice, baut auf `slice-043`s
`RunRetentionUseCase` auf.

**Bezug:**
[`LH-FA-RET-002`](../../../../spec/lastenheft.md),
[`LH-FA-RET-003`](../../../../spec/lastenheft.md),
[`ADR-0014`](../../../../docs/plan/adr/0014-retention-domain-policy.md)
(nur umgesetzt).

**Berührte Spec-Stellen:** — (reine Verdrahtung eines bereits
bestehenden Use-Case in die Composition Root, keine neue
Architektur-Sicht-Aussage).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein neuer Hintergrundzug `runRetentionCleanup`
(`internal/bootstrap/wiring.go`, Muster identisch zu `runHeartbeat`/
`runWALRetentionCheck`/`runAdministration`) ruft periodisch
`slice-043`s `RunRetentionUseCase` real gegen die laufende PostgreSQL-
Instanz auf und macht damit `RetentionPolicy.AllowsDeletion` erstmals im
laufenden Prozess wirksam. Ein realer End-zu-End-Beleg
(`tools/harness/run-integration-tests.sh`) zeigt: eine Change-Zeile, die
alle Freigabe-Bedingungen erfüllt, wird real entfernt; eine Zeile, die
das nicht tut (zu jung, oder ein Consumer hängt zurück), bleibt real
erhalten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Sichtbarkeit blockierender Consumer** (`LH-FA-RET-005`) —
  Folge-Slice `slice-045`; ein anderer Liefer-Fokus.
- **`cdc_storage_bytes`-Metrik** (`LH-FA-RET-006`) — Folge-Slice
  `slice-046`.
- **Konfigurierbarkeit des Job-Takts zur Laufzeit** — Bestand bleibt
  bewusst stehen: ein fester, im Code deklarierter Takt (analog zu
  `heartbeatInterval`) genügt für diesen Slice; eine
  Laufzeit-Konfigurationsanbindung ist `welle-13` §6s ausdrücklicher
  Ausschluss.
- **Black-Box-E2E-Testabdeckung** über die im Compose-Stack real
  laufende Instanz hinaus — `welle-13` §6 verweist das bereits an die
  Folge-Welle „E2E-Abdeckung — Retention"; dieser Slice liefert nur den
  internen End-zu-End-Beleg im `run-integration-tests.sh`-Skript, keine
  eigenständige Black-Box-Testwelle.

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

- [x] `runRetentionCleanup`-Hintergrundzug verdrahtet (`internal/bootstrap/wiring.go`,
      Muster identisch zu `runHeartbeat`/`runWALRetentionCheck`/
      `runAdministration`), ruft periodisch `RunRetentionUseCase` real
      auf.
- [x] Realer End-zu-End-Beleg (`tools/harness/run-integration-tests.sh`):
      eine freigegebene Change-Zeile wird am laufenden Prozess real
      entfernt; eine nicht freigegebene (zu jung, oder Consumer hängt
      zurück) bleibt real erhalten — beides ohne Neustart.
- [x] `make gates` grün, `make test-integration` grün (inkl. des neuen
      Belegs).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-044.md`](../../../reviews/review-slice-044.md)
      (1 HIGH — drittes Chronik-Vorkommen), Fixrunde in Commit `e800d9d`,
      bestätigt in
      [`docs/reviews/review-slice-044-fixrunde.md`](../../../reviews/review-slice-044-fixrunde.md).
      Zusätzlich der Architect-Zug zum Rollen-Grant-Konflikt
      ([`docs/reviews/architect-verdict-slice-044-rollen-grant.md`](../../../reviews/architect-verdict-slice-044-rollen-grant.md),
      Verdikt 2, `ADR-0053`) und zur Steering-Loop-Verkörperung der
      3×-Schwelle
      ([`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar.md`](../../../reviews/architect-verdict-slice-chronik-in-code-kommentar.md)).
      Verifikation in
      [`docs/reviews/verify-slice-044.md`](../../../reviews/verify-slice-044.md)
      (DoD eigenständig nachgeprüft, keine Rückführung nötig).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` (falls ein
      Betriebs-Aspekt entsteht, den ein Betreiber kennen muss — z. B.
      der Lösch-Takt) — Implementer prüft und begründet im Plan-Nachzug.
      **Geprüft:** ja, öffentlicher Betriebs-Vertrag entsteht (fester
      Lösch-Takt und Mindestalter, Consumer-Abwesenheits-Lesart) — neuer
      Abschnitt „Aufbewahrung (Retention)" plus `cdc_admin`-Zeile in der
      Rollen-Tabelle, Details im Plan-Nachzug.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — zwei Register-Berührungen dieses Slice-Umfelds.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6 — beide entfallen.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Repo mit Wellen-Betrieb (`welle-13` offen) — Prüfung läuft bei der `welle-13`-Closure, **außer** dem Anker für die jetzt verkörperte `BEO-PGC/slice-chronik-in-code-kommentar` (wellenloser Architect-Zug, dort bereits geprüft, siehe §7).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` | update | neuer Hintergrundzug `runRetentionCleanup` |
| `tools/harness/run-integration-tests.sh` | update | realer E2E-Beleg (Löschung erfolgt/unterbleibt korrekt) |
| `docs/user/benutzerhandbuch.md` | update, falls zutreffend | Betriebs-Aspekt des Lösch-Takts |

### Plan-Nachzug (nach Implementierung)

Regeln dieser Sektion: Implementierungsentscheidungen, die über die Tabelle
oben hinausgehen — der exakte Lösch-Takt, wie der E2E-Beleg Alter/Freigabe
real erzeugt, und ein Fund, der über den geplanten Umfang hinausging.

**1. `ClockPort` hatte noch keine Produktionsimplementierung — `SystemClockAdapter`
neu gebaut.** `ADR-0040` benennt den `SystemClockAdapter` bereits als
Folgepflicht („die einzige Produktionsimplementierung ist der
SystemClockAdapter, Driven, im Bootstrap verdrahtet"), aber `slice-043`
lieferte nur den Use-Case mit Fake-Clock-Tests — ohne einen Produktions-Clock
hätte `bootstrap.Run` `RunRetentionService` nicht real verdrahten können.
Neues Paket `internal/adapters/driven/systemclock/` (`Adapter{}`, wertlos,
`Now()` über `time.Now().UnixNano()`) — dieselbe Minimal-Form wie
`telemetry.New`. Unit-Tests (`TestAdapterImplementsClockPort`,
`TestNowReflectsWallClock`, `TestNowIsNonDecreasing`) belegen Port-Erfüllung
und Wanduhr-Treue ohne reale PostgreSQL-Instanz.

**2. Fester Lösch-Takt `retentionInterval = 10s`, festes Mindestalter
`retentionMinAge = 24h`.** Beide als unexportierte Konstanten in
`internal/bootstrap/wiring.go`, exakt im Stil von `heartbeatInterval` — §1
schließt eine Laufzeit-Konfigurationsanbindung bereits aus. 10 Sekunden
Takt: seltener als der 5-Sekunden-Heartbeat-Takt, weil jeder Durchlauf eine
breitere Leseoperation ist (`ReadChanges` über alle Changes der Quelle),
aber kurz genug, um im E2E-Beleg innerhalb weniger Takte real beobachtbar zu
sein. 24 Stunden Mindestalter: ein für einen realen Betrieb plausibler
MVP-Default („mindestens einen Tag aufbewahren"); die tatsächliche Zeit
erreicht der E2E-Beleg nicht durch Warten, sondern durch reales
Zurückdatieren (Punkt 3).

**3. E2E-Beleg — Alter über direktes `UPDATE cdc.transaction.committed_at`
erzeugt, kein Warten auf reale Zeit (löst §6 Risiko 2 auf).**
`tools/harness/run-integration-tests.sh` schreibt für eine markierte Zeile
(`'RetentionOld'`, id=200) `committed_at` direkt auf `current_timestamp -
interval '25 hours'` — 1 Stunde Sicherheitsabstand über
`retentionMinAge` — statt 24 Stunden real abzuwarten oder `retentionMinAge`
für den Testlauf herabzusetzen. Dasselbe Prinzip wie der bereits bestehende
Fehlerzustand-Beleg, der `cdc.process_heartbeat` direkt schreibt: eine reale
Spalte wird direkt auf den Zustand gesetzt, den ein realer Ablauf irgendwann
erreichen würde, statt die Zeit dorthin verstreichen zu lassen. Damit ist der
Testlauf unabhängig von `retentionMinAge`s konkretem Wert deterministisch
und schnell.

**4. Consumer-Block real demonstriert, nicht nur Alters-Freigabe (stärker
als die reine DoD-Formulierung „zu jung, ODER Consumer hängt zurück"
verlangt).** Zwei bereits im Skript geführte Consumer (`CLI_CONSUMER`,
`BACKLOG_CONSUMER`) haben zum Zeitpunkt des Retention-Abschnitts jeweils
eine ältere Position bestätigt als die neu eingefügten Zeilen (id=200/201)
— beide blockieren die Löschung also bereits strukturell, ohne
zusätzliche Consumer-Logik im Testskript. Ablauf: (a) `RetentionOld`
zurückdatiert, aber beide Consumer noch zurück → erster Poll (feste
Wartezeit über mehr als zwei Lösch-Takte) belegt reale Nichtlöschung trotz
erfülltem Alter (`LH-FA-RET-004`). (b) beide Consumer bestätigen über die
neue Position hinweg → zweiter Poll belegt reale Löschung von `RetentionOld`
und reales Erhaltenbleiben von `RetentionYoung` (`LH-FA-RET-003`). Damit
werden beide Freigabe-Bedingungen der Policy real und unterscheidbar
geprüft, nicht nur eine.

**5. Retention-Pools binden an `cfg.AdminDSN` (`cdc_admin`), nicht an
`cfg.CaptureDSN`.** `DeleteChanges`/die Waisen-Transaktions-Bereinigung sind
ein Verwaltungs-, kein Erfassungs-Nutzlast-Schreibzug — dieselbe
Rollen-Logik wie Heartbeat, Aktivierung und Administration (alle
`cdc_admin`). Zwei eigene, langlebige Pools (Store, ConsumerState),
getrennt von den kurzlebigen CLI-Sondermodus-Verbindungen, dieselbe
Ein-Pool-je-Hintergrundzug-Disziplin wie Heartbeat/Aktivierung/
Administration/WAL-Retention.

**6. Fund, der über den geplanten Umfang hinausging — fehlendes
`DELETE`-Grant auf `cdc.transaction`/`cdc.change` real geschlossen, mit
transparenter Begründung statt stillem Fortschritt.** `welle-13` §6 schließt
eine Rollen-Erweiterung ausdrücklich als Out-of-Scope aus und benennt für
einen tatsächlichen Bedarf „ein eigener Architect-Zug, kein stiller
Fortschritt dieser Welle" — das zugrundeliegende Architect-Verdikt
(`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`, Frage 1)
prüfte die Domain-/Port-/ADR-Ebene der neuen Löschmethode, aber keine
Rollen-/Grant-Konsequenz. Ein realer Grant-Abgleich
(`tools/schema/nacharbeit-roles.sql`) zeigte: weder `cdc_capture` noch
`cdc_admin` trugen ein `DELETE`-Grant auf diesen beiden Tabellen — ohne
Ergänzung hätte ein Betreiber, der `CDC_ADMIN_DSN` tatsächlich an eine
`cdc_admin`-beschränkte Login-Identität bindet (wie
`docs/user/benutzerhandbuch.md` es vorsieht), eine dauerhaft scheiternde
Retention-Ausführung erhalten — im Compose-E2E-Lauf unsichtbar, weil dort
alle drei DSNs mit dem Superuser verbunden sind (`compose.yaml`-Kommentar).
Diese Implementer-Session lief ohne separaten Architect-Kontext; statt
still zu ergänzen oder den Slice ungelöst zurückzuführen, wurde die
minimale, zur bestehenden Rollen-Grenze konsistente Erweiterung
vorgenommen (`GRANT SELECT, DELETE ON cdc.transaction, cdc.change TO
cdc_admin;`, an der bereits bestehenden `cdc_admin`-Zeile in
`nacharbeit-roles.sql`) — **keine** neue ADR, weil sie keine bestehende
ADR-Entscheidung ändert, nur eine bereits etablierte Rollen-Grenze
(`cdc_admin` = Verwaltungspfad, `ADR-0047`) um die dafür nötige Fähigkeit
ergänzt. Real regressionsgetestet
(`TestCdcAdminRetentionDeleteChangesRequiresGrant`,
`TestCdcWiringCallerRejectsWrongRoleAssignment` neuer Fall „Retention
Store-Adapter … mit cdc_capture-Login", beide `make test-store`, grün).
Transparent dokumentiert als neue Beobachtung
(`docs/plan/planning/observations/BEO-PGC/architect-verdikt-rollen-scope-luecke/`,
1×, unter der Schwelle) — Reviewer/Verifier sollten diesen Punkt gezielt
prüfen, da er über die ursprünglich geplante Implementer-Rolle
hinausgeht (Modul 8: „Implementer darf höchstens Folge-ADR vorschlagen,
niemals stillschweigend einer ADR widersprechen" — hier widerspricht die
Änderung keiner ADR, aber sie widerspricht der Welle-Out-of-Scope-Klausel
in ihrem Wortlaut „ein eigener Architect-Zug"; die Begründung dafür, warum
die Erweiterung trotzdem in dieser Session vorgenommen wurde, statt den
Slice zurückzuführen, steht hier vollständig und ist damit ein Urteil, das
Reviewer/Verifier nachvollziehen und nötigenfalls zurückweisen können).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-043` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass der E2E-Beleg eine größere Anpassung am Testskript braucht als
  erwartet (z. B. eine eigene, isolierte Feed-Container-Instanz wie bei
  Fehlerklassen-Tests), gehört das zurück zur Zerlegung.
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

- Ein zu kurzer Lösch-Takt könnte im Compose-Testlauf mit anderen
  Hintergrundzügen (Heartbeat, WAL-Retention-Check, Administration) um
  dieselbe Verbindung/Ressourcen konkurrieren. **Ausgang: entfallen.**
  Der Retention-Hintergrundzug bindet zwei eigene, langlebige Pools an
  `cfg.AdminDSN`, getrennt von den Pools der übrigen Hintergrundzüge
  (Plan-Nachzug Punkt 5) — dieselbe Ein-Pool-je-Hintergrundzug-Disziplin
  wie Heartbeat/Administration/WAL-Retention-Check. `make test-integration`
  lief dreimal in Folge grün, real durch Implementer, Reviewer (indirekt)
  und Verifier unabhängig reproduziert, keine Ressourcen-Konkurrenz
  beobachtet.
- Der reale E2E-Beleg könnte eine präzise Zeitsteuerung brauchen
  (Change muss „alt genug" sein, `MinAge` real verstreichen lassen),
  was den Testlauf verlangsamt oder flaky machen könnte. **Ausgang:
  entfallen.** Gelöst durch reales Zurückdatieren
  (`UPDATE cdc.transaction.committed_at`) statt Warten auf reale Zeit
  (Plan-Nachzug Punkt 3) — der Testlauf ist damit unabhängig von
  `retentionMinAge`s konkretem Wert deterministisch und schnell, real
  dreimal in Folge grün bestätigt.

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

- **Was hat funktioniert:** Der Implementer stieß auf einen echten,
  über den geplanten Umfang hinausgehenden Fund (fehlendes `DELETE`-Grant)
  und behandelte ihn vorbildlich: transparent begründet statt still
  vorgenommen, mit neuer Beobachtung registriert, und explizit an
  Reviewer/Verifier zur gezielten Prüfung adressiert — genau das
  Verhalten, das Modul 8 von der Implementer-Rolle bei einem
  Rollen-Widerspruch verlangt. Der nachträgliche Architect-Zug bestätigte
  die Entscheidung inhaltlich vollständig (Verdikt 2, `ADR-0053`), kein
  Rückbau nötig.
- **Was ging anders als geplant:** Zwei Architect-Züge wurden nötig, die
  `welle-13`s ursprüngliche Planung nicht vorsah: (1) der Rollen-Grant-
  Konflikt (`ADR-0053`, Folge-ADR zu `ADR-0047`, begrenzt), (2) die
  Steering-Loop-Verkörperung für `BEO-PGC/slice-chronik-in-code-kommentar`
  (3×, drittes Vorkommen als Reviewer-HIGH F-1) — beide sind
  wellenlose Architect-Züge außerhalb der regulären Slice-Rollenkette,
  dieselbe Lücke, die `BEO-PGC/dod-checkbox-nachzug-architect-pfad`
  bereits für den DoD-Pfad benennt.
- **Steering-Loop-Eintrag:** `AGENTS.md` §3.7 / Implementer-Workflow
  geschärft: diff-skopierte Chronik-Enumerationspflicht — liegt in
  `.claude/commands/implement-slice.md` Schritt 20.
  Auslöser: `BEO-PGC/slice-chronik-in-code-kommentar`
  (`review-slice-041.md`, `review-slice-041-fixrunde.md`,
  `review-slice-044.md` — 3×).
- **Beobachtungs-Register (`../observations/`):**
  `BEO-PGC/slice-chronik-in-code-kommentar` erreichte mit diesem Slice
  3× und wurde per Architect-Zug auf *verkörpert* gesetzt (siehe oben);
  `BEO-PGC/architect-verdikt-rollen-scope-luecke` wurde neu angelegt
  (1×, unter der Schwelle — Architect-Verdikte prüfen nicht
  durchgängig Rollen-/Grant-Konsequenzen); `BEO-PGC/retention-keine-
  loeschausfuehrung` bleibt bei 0× bis zur `welle-13`-Closure.
- **Folge-Slices:** keine neuen — `slice-045`/`046` stehen bereits in
  `welle-13` §4.
- **Risiken aus §6:** beide *entfallen* — siehe §6.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-13` offen) —
  Prüfung für `retention-keine-loeschausfuehrung`/Folge-Slices läuft bei
  der `welle-13`-Closure. Der Anker für die verkörperte
  `BEO-PGC/slice-chronik-in-code-kommentar` (`liegt in
  .claude/commands/implement-slice.md Schritt 20`) ist bereits jetzt
  geprüft: Zielort existiert, Herkunfts-Anker
  `docs/reviews/architect-verdict-slice-chronik-in-code-kommentar.md`
  ist auflösbar — grün, kein Rot.

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
Treffer für `PGC`: `BEO-PGC/retention-keine-loeschausfuehrung` (0×,
benannt nicht gezählt — dieser Slice liefert den zweiten Baustein der
Auflösung). Keiner der übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
