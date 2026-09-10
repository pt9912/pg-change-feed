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

- [ ] `cmd/pg-change-feed` verdrahtet die reale Pipeline (Store,
      Stream, Service, ACK) je [`ADR-0026`](../../../../docs/plan/adr/README.md) — der
      Feed-Container fährt CDC-Runtime statt `--version`-Smoke.
- [ ] Der MVP-Integrationstest (slice-006) fährt gegen den
      verdrahteten Feed-Container: INSERT/UPDATE/DELETE landen im
      Store (Ende-zu-Ende durch das Binary).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
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

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): <Bedingung>
- `in-progress` → `open` (blockiert — Carveout?): <Bedingung>

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

<…>

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

- **Was hat funktioniert:** (fällig bei der Implementierung)
- **Was ging anders als geplant:** (fällig bei der Implementierung)
- **Steering-Loop-Eintrag:** (fällig bei der Implementierung)
- **Beobachtungs-Register (`../observations/`):** (fällig bei der Implementierung)
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

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

**Vorgelagert — Sub-Area-Wahl prüfen:** <je berührter Sub-Area: erfüllt sie
die Schwelle ≥ 2 von 3 Achsen? zu grobe vorher ausdifferenzieren>

**Vorgelagert — offene Beobachtungen sichten:** <Register durchgegangen;
je berührter Sub-Area der Treffer mit Zähler-Stand — oder "keine Treffer">

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

### Sub-Area: <Name>

- **Modus:** GF | BF | Hybrid
- **Konventionen-Dichte:** <Beleg aus `harness/conventions.md`,
  Adaptions-Block oder Code>
- **Phase-Reife:** Phase 0–5 <Begründung gegen die Phase × Modus-Matrix>
- **Evidenz-/Diskrepanz-Risiko:** <bei BF/Hybrid: was kann die
  Inventur sichtbar machen? bei GF: meist niedrig>
- **Reconciliation-Aufwand:** <Slice-Schätzung;
  Graduation-/Folge-Slice-Trigger>
