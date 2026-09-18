# Slice slice-001: Go-Modul-Bootstrap und Gate-Aktivierung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-1.

**Bezug:** [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (reproduzierbares
Deployment, MVP), [`LH-QA-OPS-001`](../../../../spec/lastenheft.md) (containerisierter
Betrieb), [`ADR-0039`](../../../../docs/plan/adr/README.md),
[`ADR-0036`](../../../../docs/plan/adr/README.md) (dieser Slice vollzieht
dessen Hochschalt-Trigger: die Maschinenprüfung der §2-Constraints wird real)

**Berührte Spec-Stellen:** [`ARC-007`](../../../../spec/architecture.md)
(Bootstrap entsteht als Paket) · [architecture.md §1](../../../../spec/architecture.md)
(Komponenten-Übersicht)

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

**Ziel:** Ein lauffähiges Go-Modul mit minimalem `cmd/pg-change-feed`-
Einstiegspunkt, das den Build-/Deploy-Vertrag real macht: `make image` baut
das erste Image (Beleg: `harness/image-hash.txt`), und der bislang
auskommentierte `a-check`-Include wird aktiviert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Domänenmodelle und Invarianten — ein Folge-Slice übernimmt es
  (slice-002, [`ADR-0039`](../../../../docs/plan/adr/README.md)-Baum wächst mit
  Lieferung); ein leeres `internal/domain/` ohne Modell wäre ein
  Schicht-Schnitt ohne Lieferwert.
- PostgreSQL-Store- und Replication-Stream-Adapter — es wäre ein anderer
  Vorgang (reale Treiber-Integration, Umfang späterer Wellen;
  [`ADR-0010`](../../../../docs/plan/adr/README.md),
  [`ADR-0006`](../../../../docs/plan/adr/README.md)).
- Aktivierung von `codepaths` im Doc-Gate — Bestand bleibt bewusst stehen:
  mehrere Doku-Dateien nennen zukünftige Pfade (Archiv-Stubs, Lifecycle-Orte);
  die Aktivierung gehört in den Moment, in dem diese Pfade existieren.

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

- [x] `go.mod`/`go.sum` existieren; ein minimales `cmd/pg-change-feed`-Binary
      baut (`go build`, CGO aus) und antwortet auf ein `--version`-Flag mit
      beobachtbarem Beleg — Teil-Beleg zu [`LH-QA-POR-003`](../../../../spec/lastenheft.md).
      *Beleg: `docker run --rm ghcr.io/pt9912/pg-change-feed:dev --version` →
      `pg-change-feed 0.1.0-bootstrap` (verify-slice-001.md).*
- [x] `make image` baut das OCI-Image; `harness/image-hash.txt` trägt den
      Digest — Teil-Beleg zu [`LH-QA-OPS-001`](../../../../spec/lastenheft.md).
      *Beleg: Digest am HEAD `sha256:3e8a106c…` (`be405e4`), Build
      deterministisch (zwei Rebuilds, verify-slice-001.md).*
- [x] `a-check`-Include in `Makefile` aktiviert und `make a-check` grün
      (aktivierte `.a-check.yml`-Fassung). *Beleg: im Gate-Bündel seit
      [`ADR-0041`](../../../../docs/plan/adr/README.md) (a-check.mk,
      GATE_CHECKS += a-check).*
- [x] `make gates` grün. *Beleg: baseline-verify OK (54 Dateien), d-check
      73 Dateien/0 Befunde, a-check im Bündel 0 Befunde
      (verify-slice-001.md, selbst gefahrene Sensoren).*
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      *Beleg: review-slice-001.md (HIGH 0 · MEDIUM 3 · LOW 3 · INFO 3,
      alle behoben oder als Grenze benannt).*
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist (hier:
      `harness/README.md` — `make image`-Zeile auf [`ADR-0039`](../../../../docs/plan/adr/README.md)
      umziehen, `make a-check` aus *Nicht behauptet* in die Werkzeuge-Tabelle).
      *Beleg: getragen; `make a-check` ist inzwischen Gate ([ADR-0041](../../../../docs/plan/adr/README.md)).*
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. *(§7)*
- [x] Reconciliation-Register — **entfällt**: Repos ohne
      Brownfield-Bootstrap haben die Datei nicht.
- [x] Beobachtungs-Register fortgeschrieben — `BEO-PGC/a-check-null-abdeckung/`
      mit `evidence/slice-001.md` (Zähler 1×); keine weitere Beobachtung
      angefallen.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (siehe §7).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — im
      Wellen-Betrieb von der Welle-1-Closure geprüft (auch für diesen
      Slice).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `go.mod`, `go.sum` | neu | Go-Modul `github.com/pt9912/pg-change-feed`, Go 1.26 (Toolchain-Digest im Dockerfile gepinnt), keine externen Deps im Bootstrap |
| `cmd/pg-change-feed/main.go` | neu | Minimal-Einstiegspunkt (Bootstrap-Verdrahtung noch leer; `--version` als beobachtbarer Beleg) |
| `Makefile` | update | `# include a-check.mk` aktivieren (Aktivierungsbedingung dieses Slice erfüllt) |
| `a-check.mk`, `.a-check.yml` | update | Prüfen, ob die Deklaration gegen den entstehenden Baum hält (Schichten-Globs) |
| `harness/README.md` | update | `make a-check` von *Nicht behauptet* in die Sensors-/Werkzeuge-Tabelle ziehen |
| `Makefile` (image-Rezept) | update | *Plan-Nachzug (nachgezogen vor der Closure):* der erste Lauf griff beim Extrahieren des Image-Digests den ersten `sha256:`-Treffer der Metadaten-Datei — das ist der **Base-Digest**, nicht der Digest des gebauten Images. Fix: Extraktion liest den `containerimage.digest`-Key (Implementer-Commit f932022) |
| `Dockerfile` (Kommentar) | update | *Plan-Nachzug:* deps-Layer-Kommentar (leeres `go.mod`/`go.sum` bleibt verifiziert); CGO-Anker von [`ADR-0038`](../../adr) (superseded) auf [`ADR-0039`](../../adr) gezogen |
| `harness/image-hash.txt` | neu | Beleg-Artefakt (Modul 14): trägt den Digest des gebauten Images; getrackt, kein Gate darauf (benannte Grenze, s. Review F-7) |
| `a-check.mk` | update | *Plan-Nachzug (V-2):* `GATE_CHECKS += a-check` — Verkabelung ins Gate-Bündel ([`ADR-0041`](../../../../docs/plan/adr/README.md)) |
| `.gitignore` | neu | *Plan-Nachzug (V-2):* `.tmp/` + `harness/image-hash.raw` ignoriert (Review F-5) |
| `tools/harness/image-stale.sh` | update | *Plan-Nachzug (V-2):* Major-Drift-Erkennung ergänzt (golang:1.27-alpine existiert, MAJOR-DRIFT gemeldet) |
| `harness/sensors/*.md` | neu | *Plan-Nachzug (V-2):* Sensor-Dateien nach gate.template.md ([`ADR-0041`](../../adr)-Zug) |
| `docs/reviews/review-report.template.md` | neu | *Plan-Nachzug (V-2):* aus der vendored Referenz-Form wiederhergestellt (Review F-8) |
| `harness/image-hash.txt` (erneuert) | update | *Plan-Nachzug (V-1):* Beleg am HEAD erneuert (`be405e4`) — der committete Beleg trug den f932022-Digest, der Zeilenenden-Fix (F-6) änderte den Build-Kontext ohne Re-Build |
| `cmd/pg-change-feed/main.go`, `go.mod` | update | *Plan-Nachzug (V-2):* Datei-Ende-Zeilenumbrüche (Review F-6) |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): welle-1 ist eröffnet (flache Datei
`welle-1.md` liegt), `make gates` grün, kein anderes Slice in
`in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wachsen die
  DoD-Punkte über die drei Liefer-Punkte hinaus (z. B. weil der erste echte
  a-check-Lauf Architektur-Befunde erzeugt, die eigene Slices verdienen).
- `in-progress` → `open` (blockiert — Carveout?): Registry-/Digest-Probleme
  (Base-Image nicht erreichbar) oder Go-Toolchain-Faktoren, die den
  Docker-Vertrag blockieren.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD (§2) vollständig abgehakt · `make image`-Beleg (`image-hash.txt`) liegt
vor · Review-Report unter `docs/reviews/` vorhanden · Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der erste echte a-check-Lauf meldet Abdeckungs-/Auflösungs-Hinweise, die
  auf Lücken in `.a-check.yml` deuten (Globs treffen den entstehenden Baum
  nicht) — **Ausgang:** offen; behoben oder ins Register erst bei Closure
  bewertet.
- Das minimale Binary ohne Verhalten liest sich als „leerer Gate" —
  **Ausgang:** entfallen; das `--version`-Flag ist der beobachtbare Beleg,
  und `make image` + `a-check` sind die tragenden Lieferungen.

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

- **Was hat funktioniert:** die Rollen-Sequenz mit Übergabe-Artefakten
  (Architect-Review vor dem `next`-Übergang; Review/Verifier in frischem
  Kontext mit selbst gefahrenen Belegen); die Docker-only-Kette ohne
  Host-Toolchain (Belege über stdout, Arbeitsbaum read-only); der Build ist
  deterministisch (zwei Rebuilds, identischer Digest — verify-slice-001.md).
- **Was ging anders als geplant:** der image-hash-Beleg driftete vom HEAD
  (V-1) — die benannte Grenze „kein Abgleich-Sensor" (Review F-7) war
  erstmals wirksam; `a-check` war vor [`ADR-0041`](../../../../docs/plan/adr/README.md)
  außerhalb des Gate-Bündels (Review F-1) — gelöst durch die Folge-ADR;
  der image-hash-Extraktion-Defekt (erster Lauf griff den Base-Digest,
  Implementer-Commit f932022) war ein Plan-Defekt, nachgezogen in §3.
- **Steering-Loop-Eintrag:** geschärfte Regel: *Beleg-Artefakte mit
  Baum-Kopplung (`harness/image-hash.txt`) gelten am HEAD — ein Zug, der
  Build-Kontext-Dateien ändert, erneuert den Beleg vor seiner Closure
  (Re-Build + Commit), sonst trägt er einen Vorgangs-Beleg*
  — liegt in `harness/README.md` (Werkzeuge-Zeile `make image`) · seit
  slice-001. *(Zählerstand: 1×; kein `BEO`-Auslöser, Erstverkörperung.)*
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/a-check-null-abdeckung/`
  neu angelegt, Beleg `evidence/slice-001.md` — Zähler steht damit bei 1×.
- **Folge-Slices:** slice-002 (Domänenkern) und slice-003
  (Capture-Persist-Pfad) — Dateien in `open/`, welle-1.
- **Risiken aus §6:** Risiko 1 (a-check-Hinweise/Glob-Lücken) → **weiter
  offen**, trägt im Register `BEO-PGC/a-check-null-abdeckung` (bewertet,
  wenn die Layer-Globs echten Content matchen, slice-002/003); Risiko 2
  (minimales Binary als leerer Gate gelesen) → **entfallen** (das
  `--version`-Flag ist der beobachtbare Beleg; `make image` + `a-check`
  sind die tragenden Lieferungen).
- **Drei Paarungen:** im Wellen-Betrieb an die Welle-1-Closure delegiert
  (Anker · Folge-Slice · Register — geprüft bei der Closure von welle-1,
  auch für diesen Slice).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührte Sub-Area: Toolchain/Bootstrap
(`go.mod`, `cmd/`, Make-Gate-Aktivierung; [`ARC-007`](../../../../spec/architecture.md)).
Achsen: (1) Konventionen-Dichte — Build-/Deploy-Vertrag verankert (Dockerfile,
([`ADR-0039`](../../../../docs/plan/adr/README.md))-Baum, `a-check.mk`/`.a-check.yml` committet); (2) Phase-Reife —
Phase 3 (Spec + ADRs committet, Code entsteht mit diesem Slice); (3) Evidenz-
Risiko niedrig — Greenfield, kein Bestand. Schwelle ≥ 2 erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** das Register
(`docs/plan/planning/observations/`) trägt nur seine `README.md` — **keine
Treffer**; notiert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Deklaration
`harness/conventions.md`, Default-Zeile) — kein BF/Hybrid, kein Block.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

*Reiner GF-Hinweis genügt (siehe oben); kein Sub-Area-Block.*
