# Slice slice-007: Bootstrap-Verdrahtung — das Binary trägt die Pipeline

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — wellenloser Korrektur-Zug (Bootstrap-Verdrahtung als Folgelielerung der slice-006-Grenze; einzeln lieferbar).

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md), [`ADR-0026`](../../../../docs/plan/adr/README.md), [`ADR-0037`](../../../../docs/plan/adr/README.md), [`ADR-0044`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`SPEC-015`](../../../../spec/pflichtenheft.md) (Deployment-Form), [`ARC-007`](../../../../spec/architecture.md) (Bootstrap)
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912 (Implementer-Rolle, Agent-Lauf).
**Autor:** pt9912. **Datum:** 2026-09-09.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das `cmd/pg-change-feed`-Binary verdrahtet die reale Pipeline (Store, Stream, Service, ACK je [`ADR-0026`](../../../../docs/plan/adr/README.md)) — der Feed-Container trägt dann den CDC-Runtime-Lauf statt des `--version`-Smokes; der MVP-Integrationstest (slice-006) fährt danach das verdrahtete System im Container.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Config-Format-Festlegung (TOML/YAML/Env) — der Bootstrap-Vertrag
  braucht eine Minimal-Verdrahtung (DSN, Publication, Slot-Name);
  die vollständige Konfigurationsschicht folgt in einem späteren Slice.
- CLI-Adapter und SQL-Funktionen — [`ADR-0018`](../../../../docs/plan/adr/README.md)/[`ADR-0019`](../../../../docs/plan/adr/README.md)-Rest, folgen nach dem MVP.
- Observability-Adapter mit voller Metrik-Abdeckung —
  [`LH-QA-OPS-003`](../../../../spec/lastenheft.md)-Umfang folgt nach der
  Verdrahtungs-Basis.

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

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `cmd/pg-change-feed` verdrahtet die reale Pipeline (Store,
      Stream, Service, ACK) je [`ADR-0026`](../../../../docs/plan/adr/README.md) — der
      Feed-Container fährt CDC-Runtime statt `--version`-Smoke. *Beleg:
      verify-slice-007.md (E2E-Probe: 2 INSERTs → 2 cdc.change-Zeilen,
      Binary als einziger Schreiber).*
- [x] Der MVP-Integrationstest (slice-006) fährt gegen den
      verdrahteten Feed-Container: INSERT/UPDATE/DELETE landen im
      Store (Ende-zu-Ende durch das Binary). *Beleg:
      `make test-integration` grün am verdrahteten Pfad
      (verify-slice-007.md, eigene E2E-Probe).*
- [x] `make gates` grün. *Beleg: vier Gates inkl.
      commit-traceability am HEAD (verify-slice-007.md).*
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      *Beleg: review-slice-007.md committet (`e971ed4`).*
- [ ] Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. (*§7*)
- [x] Reconciliation-Register — **entfällt**: Repos ohne
      Brownfield-Bootstrap haben die Datei nicht., **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (siehe §7).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `cmd/pg-change-feed/main.go` | update | Verdrahtung: Bootstrap konstruiert Store/Stream/ACK/Service je [`ADR-0026`](../../../../docs/plan/adr/README.md) |
| `internal/bootstrap/*.go` | update | Verdrahtungs-Funktionen (kein Test-Layer mehr allein) |
| `compose.yaml` | update | Feed-Container fährt CDC-Runtime (ENV für DSN/Publication) |
| `tools/harness/run-integration-tests.sh` | update | *Plan-Nachzug (Review F-1):* der Runner trägt die Aktivierung (Feed-Tabellen, Publication, Bindungs-Zeilen) **vor** dem Container-Start und prüft den Stream-Vertrag zweigeteilt (Slot **und** `State.Running`) — Wächter-Lücke aus der eigenen Probe belegt |
| `test/integration/mvp_test.go` | update (neuschreiben) | *Plan-Nachzug (Review F-2):* der Test konsumiert über den Store-Adapter-Lese-Pfad (`PostgresChangeStoreAdapter.ReadChanges`) — „kein anderer Schreiber als das Binary" geprüft; die Go-Baum-Verdrahtung des slice-006-Stands ist entfernt |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): welle-2 schließt (alle drei Slices in
`done/`) — der Verdrahtungs-Zug trägt den Verdrahtungslücken-Ausgang
(Review F-6, slice-006). Kein anderes Slice in `in-progress/`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wächst der
  Verdrahtungs-Scope über die drei Liefer-Punkte hinaus (z. B. vollständige
  Konfigurationsschicht) → aufteilen.
- `in-progress` → `open` (blockiert — Carveout?): der Feed-Container
  kann die CDC-Runtime nicht fahren (Image-/Compose-Faktoren).

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD (§2) vollständig abgehakt · `make test-integration` grün am
verdrahteten Pfad · Review-Report unter `docs/reviews/` · Closure-Notiz
mit Lerneintrag.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Verdrahtungs-Lücke (Review F-6, slice-006): der Feed-Container ist
  Image-Vertrag-Smoke — **Ausgang:** eingetreten; Träger ist dieser
  Slice (das Binary verdrahtet, der Container fährt CDC-Runtime).
- Fehlerbehandlungs-Grenze (Review F-2): der erste Production-Pfad endet
  auf jeden Adapter-Fehler mit Ausgang 1; die [`SPEC-008`](../../../../spec/pflichtenheft.md)-Aktion für
  `transient` (Retry/Backoff) trägt kein Element, und die Adapter-Grenze
  „kontrollierte Fortsetzung beim Aufrufer" wird vom Aufrufer mit
  Prozess-Ende beantwortet — **Ausgang: weiter offen** → Retry/Backoff
  folgt mit der Konfigurationsschicht (späterer Slice); die Grenze trägt
  der Kommentar in `wiring.go` und der `restart: "no"`-Vertrags-Zeile in
  `compose.yaml` (Fix-Runde).
- MVP-Integrationstest-Claim hängt an der Go-Baum-Ebene — **Ausgang:**
  eingetreten (Grenze benannt); der Test fährt nach diesem Slice das
  verdrahtete System am Container.

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

- **Was hat funktioniert:** die Verdrahtung trennt Verbindungsaufbau
  (`BindCapture`) von Port-Verdrahtung (ADR-0007 Option C) — der
  ACK-Adapter braucht die Verbindung erst nach deren Aufbau; der
  zweigeteilte Wächter (Slot **und** `State.Running`) wurde am
  Fehlmodus-Probe live rot gesehen (Publication fehlt → Exit 1, Slot
  blieb); der MVP-Integrationstest fährt danach das verdrahtete System
  (E2E-Probe: 2 INSERTs → 2 cdc.change-Zeilen, Binary als einziger
  Schreiber — verify-slice-007.md).
- **Was ging anders als geplant:** der Review-F-7-Schiedsspruch
  bestätigte den ADR-0044-Vertrag (Digest lauf-gebunden) — diesmal
  wechselte der Digest wirklich (Binary trägt die Verdrahtung), und der
  Beleg wurde am HEAD committet (`789b76e`). Die commit-traceability-
  Klasse färbte beim ersten F-3-Commit **selbst rot** (ADR-0045-Sensor
  greift) — die Klasse ist jetzt maschinell getragen (ADR-0045).
- **Steering-Loop-Eintrag:** *F-2-Fehlerbehandlungs-Grenze → Ausgang
  weiter offen (Retry/Backoff folgt mit der Konfigurationsschicht);
  die Grenze trägt der `wiring.go`-Kommentar und die Container-
  Vertrags-Zeile · seit slice-007.* *(Der Review-INFO-Kanal F-10
  (Signal-Ende-Vertrag) bleibt ohne automatisierten Träger —
  Binary-Lauf-Test-Aufbau wächst über den Nachzug hinaus; Restrisiko
  benannt.)*
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/adapter-fehler-ausgang/`
  neu angelegt (Beleg `evidence/slice-007.md`, 1×).
- **Folge-Slices:** keiner — wellenloser Zug, abgeschlossen vor der
  Welle-2-Closure.
- **Risiken aus §6:** Verdrahtungs-Lücke → **eingetreten** (Träger:
  dieser Slice); Fehlerbehandlungs-Grenze → **weiter offen** (Retry/
  Backoff mit der Konfigurationsschicht).
- **Drei Paarungen:** im Wellen-Betrieb an die Welle-2-Closure delegiert.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührte Sub-Area:
Bootstrap-Verdrahtung (`cmd/**`, `internal/bootstrap/**`; [`ARC-007`](../../../../spec/architecture.md)).
Achsen: (1) Konventionen-Dichte — Composition-Root-Regeln in
[`ADR-0026`](../../../../docs/plan/adr/README.md), Deployment-Form in [`SPEC-015`](../../../../spec/pflichtenheft.md);
(2) Phase-Reife — Phase 5 (alle Adapter und Ports real committet,
Verdrahtung entsteht); (3) Evidenz-Risiko niedrig (GF). Schwelle ≥ 2
erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** <Register durchgegangen;
je berührter Sub-Area der Treffer mit Zähler-Stand — oder "keine Treffer">

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

*Reiner GF-Hinweis genügt (siehe oben); kein Sub-Area-Block.*
