# Slice slice-057: E2E-Testumgebung — Umbenennung des MVP-Namensrelikts

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reine Umbenennung ohne Verhaltensänderung, die
Closure-Bedingung ist ausschließlich die eigene DoD (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht). Orthogonal zu
`welle-15` (NATS) — reine Test-/Compose-Infrastruktur, keine CDC-Fähigkeit.

**Bezug:** [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (vollständige
Testumgebung automatisiert und reproduzierbar) — keine ADR, da reine
Namensänderung ohne Architekturentscheidung.

**Berührte Spec-Stellen:** — (reine Namensänderung, keine Spec-Stelle
berührt).

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

**Ziel:** Die E2E-/Compose-Testumgebung (`compose.yaml`,
`test/integration/integration_test.go`,
`tools/harness/run-integration-tests.sh`,
`docs/user/benutzerhandbuch.md`-Beispiele) trägt seit `welle-1`/`welle-2`
die Namensgebung `src-mvp`/`pub_pgc_mvp`/`slot_pgc_mvp`/`feed_mvp_*` — ein
Relikt aus der ursprünglichen MVP-Abnahme (`M1`), obwohl die Umgebung
längst der allgemeine Black-Box-E2E-Rundlauf für alle Wellen ist
(`ADR-0030` E2E-Tier), nicht mehr MVP-beschränkt. Dieser Slice benennt
alle betroffenen Bezeichner **verhaltensgleich** um: `src-mvp` →
`src-e2e`, `pub_pgc_mvp` → `pub_pgc_e2e`, `slot_pgc_mvp` → `slot_pgc_e2e`,
`feed_mvp_flow`/`feed_mvp_full`/`feed_mvp_idle`/`feed_mvp_missing`/
`feed_mvp_schema`/`feed_mvp_sql`/`feed_mvp_walsender` →
`feed_e2e_flow`/`feed_e2e_full`/`feed_e2e_idle`/`feed_e2e_missing`/
`feed_e2e_schema`/`feed_e2e_sql`/`feed_e2e_walsender` — inklusive der
zugehörigen Go-Bezeichner in `test/integration/integration_test.go`
(`mvpSource`/`mvpPublication`/`mvpSlot`/`newMVPEnv` →
`e2eSource`/`e2ePublication`/`e2eSlot`/`newE2EEnv`), konsistent mit dem
bereits in `slice-053` neu eingeführten Container-Namen
`cdc-e2e-natssub` (der schon die `e2e`-Konvention trägt).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Archivierte/geschlossene Dokumente** (`docs/plan/planning/done/*.md`,
  `docs/reviews/*.md`) — Historie wird nicht rückwirkend umgeschrieben
  (dieselbe Disziplin wie bei Slice-Chronik: `git` hält die Historie,
  nicht der aktuelle Bestand; ein Review-Report von `slice-007`
  beschreibt korrekt, was zum Zeitpunkt seiner Entstehung galt).
- **Der Meilenstein-Begriff „MVP" selbst** (`spec/lastenheft.md` §1
  MVP-Schnitt, `roadmap.md` M1 — MVP-Abnahme) — das ist ein fachlicher
  Meilenstein-Name, kein technisches Bezeichner-Relikt; er bleibt
  unverändert.
- **Verhaltensänderung jeder Art** — reines Renaming; kein Test darf
  eine andere Aussage treffen als vorher, nur unter neuem Namen.
- **Compose-Servicenamen/Container-Namen** (`cdc-test-postgres`,
  `cdc-test-feed`) — tragen kein `mvp`-Fragment, bleiben unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Alle in §1 genannten Bezeichner real umbenannt in `compose.yaml`,
      `test/integration/integration_test.go`,
      `tools/harness/run-integration-tests.sh`,
      `docs/user/benutzerhandbuch.md` — real per `grep -rn "mvp"` (case-
      insensitive, außer den zwei §1-Ausnahmen) auf null Treffer geprüft.
- [x] `make test-integration` läuft real vollständig durch und ist
      inhaltlich unverändert erfolgreich (derselbe Rundlauf, neue Namen).
- [x] `make gates`, `make test` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md`-Beispiele, falls sie
      die alten Namen zeigen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `compose.yaml` | update | `CDC_SOURCE_ID`/`CDC_PUBLICATION`/`CDC_SLOT`/`CDC_TABLES` umbenannt — inklusive der Bindungs-/Schema-Version-Alias-Fragmente in `CDC_TABLES` (`tbl-mvp-*`→`tbl-e2e-*`, `sv-mvp-*`→`sv-e2e-*`), die derselbe Bezeichner-Stamm sind, aber nicht wörtlich in §1 aufgeführt waren (*Plan-Nachzug*) |
| `test/integration/integration_test.go` | update | Go-Konstanten/-Funktionsnamen + Tabellen-/Kennungs-Literale umbenannt — zusätzlich zu den in §1 genannten alle `TestMVP*`-Funktionsnamen sowie der Typ `mvpEnv`→`e2eEnv` (*Plan-Nachzug*, DoD verlangt Null-Treffer für `grep -in mvp` über die ganze Datei) |
| `tools/harness/run-integration-tests.sh` | update | alle `mvp`-Literale (Subjekt, SQL, Kommentare) umbenannt — zusätzlich `feed_mvp_sql_admin`→`feed_e2e_sql_admin`, `feed_mvp_walsender_timing`→`feed_e2e_walsender_timing`, die `-run`-Regex-Liste der `TestMVP*`-Namen sowie der Freitext-Wert `'MVP-Quelle'`→`'E2E-Quelle'` (*Plan-Nachzug*) |
| `docs/user/benutzerhandbuch.md` | update | Beispiel-Wert `src-mvp`→`src-e2e` in der `diagnose`-Beispielausgabe (§4); zieht per Implementer-Workflow (`BEO-PGC/handbuch-versionshistorie-uebersprungen`) den `Version:`-Kopf (1.11→1.12) und eine neue Zeile in `### Änderungshistorie` mit (*Plan-Nachzug*) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-053` liegt in `done/` (dieselben
Dateien `compose.yaml`/`run-integration-tests.sh` dürfen nicht gleichzeitig
von zwei Implementer-Läufen bearbeitet werden — WIP-Limit 1 erzwingt das
ohnehin), `Verantwortlich:` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten — reines, mechanisches Renaming ohne Verhaltensänderung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein übersehenes `mvp`-Literal (z. B. in einer SQL-Zeichenkette innerhalb
  eines Heredocs in `run-integration-tests.sh`, wo `grep` über
  Zeilenumbrüche hinweg leichter etwas übersieht) könnte den Testlauf
  unbemerkt gegen einen inkonsistenten Namensmix laufen lassen. **Ausgang:**
  <bei Closure einzutragen>
- Falls `slice-054`/`055` (Boundary-/Negative-Belege für NATS) inzwischen
  bereits Text mit den alten `mvp`-Namen tragen (sie sind zum Zeitpunkt
  dieser Planung noch unangetastet in `open/`), müssen ihre Plan-Texte
  nachgezogen werden, bevor sie aktiviert werden — sonst driftet ihr Text
  vom tatsächlichen Namensstand. **Ausgang:** <bei Closure einzutragen>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Keine Treffer für `PGC` zu Namensgebung/Umbenennung der Testumgebung.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Entfällt — reines Renaming ohne neue Sub-Area-Berührung.
