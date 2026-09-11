# Slice slice-015: d-migrate-1.3.0-Retirement

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reaktive Wartungsarbeit (Pin veraltet, Fix-Release
liegt vor), keine Closure-Bedingung über die eigene DoD hinaus, siehe
Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht
(Modul 6).

**Bezug:** [`ADR-0043`](../../../../docs/plan/adr/README.md) (Re-Evaluierungs-Trigger: Ausweichform entfällt, wenn d-migrate die Operation ausdrücken kann)

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-11.

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

**Ziel:** d-migrate-Pin auf v1.3.0 heben (Digest verifiziert:
`sha256:d8dc38cfe3315ab4a1df9f366b262700f402a531b62abd27eb836d3bc4c1e5cc`)
und die berichtete Nacharbeit-Ausweichform zurückbauen, wo der Fix-Release
sie tatsächlich löst — Changelog nennt „Post-compare drift eliminated"
(Server-Form-zu-Server-Form-Vergleich); dieser generelle Mechanismus ist
es, der ggf. die `raw-sql-text-drift`-Klasse aus
`BEO-PGC/d-migrate-nacharbeit` (2×: CHECK-Ausdruck slice-006, gespeicherte
Views slice-010) adressiert — ein eigenständiges „Raw SQL Sandbox
Mode"-Feature existiert nicht (real geprüft, siehe DoD-Punkt 2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Rollen-DDL (`nacharbeit-roles.sql`) zurückbauen — anderer Vorgang: die
  slice-011-Planung stellte bereits fest, dass Rollen/GRANTs keinen
  `raw-sql-text-drift`-Fall auslösten (`make test-store` blieb grün ohne
  die Ausweichform-Ursache); ob d-migrate 1.3.0 ROLE/GRANT-Objekte
  überhaupt ins neutrale Modell aufnimmt, ist eine andere Frage als die
  hier adressierte Drift-Behebung und wird nur geprüft, falls der
  Implementer beim Ist-Zustand-Abgleich einen Bezug findet.
- Provenance-Overlay-Workflow (`--provenance-output`, `raw-text-provenance`)
  produktiv einführen — Bestand bleibt bewusst stehen: der reguläre
  Post-compare-Vergleich (Server-Form-zu-Server-Form) genügt allein, um
  die Drift-Klasse zu testen; das Overlay ist eine
  Convergence-Optimierung für Folge-Läufe, kein Blocker-Fix, und bräuchte
  eine eigene Bewertung des Rollout-Workflows.

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

- [x] d-migrate-Pin auf v1.3.0 gehoben (`Makefile` `D_MIGRATE_IMAGE`),
      `make schema-validate` grün am neuen Pin. Beleg: `make schema-validate`
      grün (Implementer-Lauf), Digest gegen die Registry verifiziert
      (`docker buildx imagetools inspect ghcr.io/pt9912/d-migrate:1.3.0`).
- [x] CHECK-Ausdruck `chk_change_operation` und die drei bestehenden
      SQL-Views (`active_tables`, `consumer_status`, `changes`) real
      gegen den neuen Pin (generischer Post-compare-Vergleich, kein
      separates Sandbox-Feature) getestet — je Fall entweder ins
      deklarative `tools/schema/schema.yaml` überführt (Ausweichform
      zurückgebaut) oder mit dokumentiertem Befund, warum nicht. Beleg: real
      gegen einen frischen Rollout (Testcontainer, außerhalb des Baums)
      getestet — CHECK konvergiert (Exit 0), alle drei Views einzeln und
      kombiniert weiterhin Post-execute-Drift (Exit 5); **kein separates
      „Raw SQL Sandbox Mode"-Flag existiert** (erschöpfend geprüft:
      `--help` auf `schema migrate`/`generate`/`validate`/`compare`/`reverse`
      und Top-Level — kein `sandbox`-Begriff). Die dem Slice zugrunde
      liegende Changelog-Prämisse eines separaten Sandbox-Modus trifft nicht
      zu; der einzig verwandte Mechanismus ist `--provenance-output`/
      `--migration-overlay` (explizit Out-of-Scope, §1) — und der kann den
      *ersten* View-Rollout ohnehin nicht lösen, weil er nur nach einem
      bereits driftfreien Lauf schreibt (real geprüft: kein Overlay-File bei
      gescheitertem Erstlauf). Details im Bericht an den Reviewer/Planner.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`review-slice-015.md`](../../../../docs/reviews/review-slice-015.md)
      (`0f00b52`), F-1/F-2 disponiert (`ceaf464`), 0 HIGH.
- [x] Doku-Update falls öffentlicher Vertrag berührt — geprüft: der
      d-migrate-Digest ist nirgends außer im `Makefile` dokumentiert,
      Item entfällt, sofern kein Sensor-Vertrag in `harness/README.md`
      sich ändert (z. B. wenn eine `nacharbeit-*.sql`-Datei entfällt und
      die [`ADR-0043`](../../../../docs/plan/adr/README.md)-Bindungszeile
      das erwähnt). Geprüft: die `schema-rollout`-Bindungszeile in
      `harness/README.md` nennt keine `nacharbeit-*.sql`-Dateinamen — Wegfall
      von `nacharbeit-operation-check.sql` berührt sie nicht, Item entfällt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Geprüft: `docs/plan/planning/reconciliation.md` existiert nicht (Repo ist GF, `harness/conventions.md` §Modus-Deklaration) — Item entfällt.
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
| `Makefile` | update | `D_MIGRATE_IMAGE`-Digest auf v1.3.0 gehoben ([`ADR-0043`](../../../../docs/plan/adr/README.md)) |
| `tools/schema/schema.yaml` | update | `chk_change_operation` deklarativ in `change.constraints` ergänzt (real gegen frischen Rollout konvergent, kein Drift) — die Views-Grenze bleibt bestehen und ist im Kommentarblock aktualisiert (1.3.0 real getestet, Exit 5 unverändert) |
| `tools/schema/nacharbeit-operation-check.sql` | gelöscht | Plan-Nachzug: Bedingung eingetreten — der CHECK-Ausdruck konvergiert mit d-migrate 1.3.0 deklarativ, real gegen einen frischen Rollout getestet (Exit 0, keine Drift) |
| `tools/schema/nacharbeit-views.sql` | update (Kommentar), Datei bleibt | Plan-Nachzug: Bedingung NICHT eingetreten — die drei Views konvergieren mit d-migrate 1.3.0 weiterhin nicht deklarativ (real getestet: Post-execute-Vergleich Exit 5, `pg_get_viewdef`-Katalogform weicht von der Autorenform ab, sowohl einzeln als auch kombiniert mit dem CHECK-Ausdruck reproduziert); Kommentar aktualisiert (Pin-Version, Bezug zum jetzt gelösten CHECK-Fall entfernt) |
| `Makefile` (`schema-rollout`-Target) | update | psql-Nacharbeit-Schritt für `nacharbeit-operation-check.sql` entfällt (Datei gelöscht); der Views-Schritt bleibt; erklärender Kommentarblock über dem Target aktualisiert |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | update | Plan-Nachzug: nicht ursprünglich gelistet — Pflicht-Report und Rollback-Artefakt des `make test-integration`-Laufs ([`ADR-0043`](../../../../docs/plan/adr/README.md), `schema migrate --execute`), der den neuen Pin real belegt (Closure-Trigger); Nebenprodukt des Sensor-Laufs, kein separat verfasster Inhalt |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): kein anderes Slice in `in-progress/`
(WIP-Limit 1) — wellenlos, kein Welle-Trigger nötig.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Falls der
  Sandbox-Modus ein grundlegend anderes Schema-YAML-Konstrukt verlangt
  (z. B. eine neue Overlay-Datei, die selbst mehrschichtig getestet
  werden müsste) statt der erwarteten direkten CHECK-/View-Ergänzung —
  Rückzug mit Zerlegung (Pin-Bump allein zuerst, Retirement als
  Folge-Slice).
- `in-progress` → `open` (blockiert — Carveout?): Falls v1.3.0 selbst
  einen neuen, bisher unbekannten Blocker gegen das bestehende
  `schema.yaml` wirft (Regressions-Risiko eines Major-Feature-Releases)
  — Blocker, Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss ohne
  offenes HIGH-Finding; `make schema-rollout` (oder `make
  test-integration`, das die Kette real fährt) grün am neuen Pin,
  Retirement-Umfang real belegt (nicht nur behauptet).
- Lerneintrag §7: geschärfte Regel oder benannte Spec-Lücke — ohne ihn
  bleibt der Slice abgelegt, nicht fertig.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Sandbox-Modus löst möglicherweise nur EINEN der beiden Fälle (CHECK
  vs. Views), nicht beide — der Retirement-Umfang wäre dann partiell.
  Wird bei Closure bewertet; ein Teil-Retirement ist ein legitimes
  Ergebnis, kein Scheitern.
- Der Digest-Bump selbst könnte einen Regressions-Fund gegen den
  bestehenden `schema.yaml`-Bestand auslösen (Major-Feature-Release,
  „View Portability Knowledge Shift" laut Changelog verschiebt
  Dialekt-Wissen) — wird bei Closure bewertet.

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührte Sub-Area
(Schemamigration/`tools/schema/`) ist ein Segment des GF-Baums der
Modus-Deklaration (`PGC` Greenfield, Doc führt) — erfüllt die Schwelle
≥ 2 von 3 Achsen. Nicht zu grob.

**Vorgelagert — offene Beobachtungen sichten (einziger Leser, da dieser
Slice wellenlos läuft):** Register gelesen (zehn Einträge zum Zeitpunkt
der Planung): `d-migrate-nacharbeit` 2× (unmittelbar betroffen — dieser
Slice ist der erwartete Auflösungs-Träger) — Beleg bei Closure prüfen.
Übrige neun ohne Bezug: `a-check-null-abdeckung` (verkörpert, seit
welle-1), `adapter-fehler-ausgang` (1×), `cdc-capture-lag-real` (1×),
`health-endpoint-heartbeat` (2×, eingetreten), `lese-doppelquelle` (2×),
`plan-nachzug` (verkörpert), `plan-vorlagen-defekt` (verkörpert),
`rollen-verdrahtung` (1×), `walsender-wirksamkeit` (1×). Kein Eintrag
erreicht mit diesem Slice 3× — keine Lücke.

**Modus-Begründungsblock — Umfang.** Reiner GF-Hinweis genügt (siehe oben);
kein Sub-Area-Block.
