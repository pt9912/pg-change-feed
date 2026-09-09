# Review-Report: slice-001 Implementer-Diff — 2026-09-09

**Review-Art:** Diff-Review (Implementer-Range `9402bef..HEAD`, 4 Commits) —
*wogegen*: Slice-Plan §1/§3 (Plan-Treue oder still erweitert), ADR-Bezüge
([`ADR-0036`](../plan/adr)/0039/0040, [`ADR-0038`](../plan/adr) als Superseded), Hard Rules (`AGENTS.md` §3.1
Docker-only, §3.7 Kommentar-Klassen), Konventionen (MR-000 ID-Schema, MR-001),
Konfig-Konsistenz (Makefile/a-check/README/Beleg-Datei). Keine DoD-Prüfung —
das ist der Verifier (Modul 11).

**Gegenstand:** `88e1517` (go.mod/go.sum/main.go) · `01e981e` (a-check-
Aktivierung) · `f932022` (image-hash-Rezept-Fix) · `ca8f0a7` (harness/README) —
Basis `9402bef` (slice-001 `next` → `in-progress`).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft: vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst: die referenzierte
Vorlage `docs/reviews/review-report.template.md` existiert nicht — dieser
Report folgt der Form des Eröffnungs-Reports
(`review-lastenheft-pflichtenheft.md`), siehe F-8.

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Diff `git log -p 9402bef..HEAD` (Dockerfile, Makefile, cmd/pg-change-feed/
  main.go, go.mod, go.sum, harness/README.md, harness/image-hash.txt)
- `docs/plan/planning/in-progress/slice-001-bootstrap.md` (§1–§3, §6)
- `docs/plan/adr/README.md` · [`ADR-0036`](../plan/adr) (Proposed) · [`ADR-0039`](../plan/adr) (Accepted,
  Supersedes [`ADR-0038`](../plan/adr)) · [`ADR-0038`](../plan/adr) (Superseded) · [`ADR-0040`](../plan/adr) (nicht berührt)
- `spec/lastenheft.md` ([`LH-QA-OPS-001`](../../spec/lastenheft.md), [`LH-QA-POR-003`](../../spec/lastenheft.md) — in Commit-Messages
  genannt), `spec/architecture.md` ([`ARC-007`](../../spec/architecture.md))
- `AGENTS.md` §3 (Hard Rules), §4; `harness/README.md`;
  `harness/conventions.md` (MR-000/MR-001)
- `Makefile`, `a-check.mk`, `.a-check.yml`, `harness/mk/enforce.mk`
  (record-gates/GATE_CHECKS-Verkabelung), `Dockerfile`, Vorbestand-Blame
  (Dockerfile:17 → `390a7b5`, vor der Range)

---

## Findings

### F-1 — a-check hängt nicht an GATE_CHECKS: Hochschalt-Trigger von ADR-0036 wird so nicht vollzogen

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0036`](../plan/adr) (Re-Evaluierungs-Trigger: Hochschalten zu Accepted,
  „sobald das Import-Linting-Gate **im Gate-Lauf** existiert und grün läuft")
  gegen Slice-Plan Kopf/Bezug („dieser Slice vollzieht dessen
  Hochschalt-Trigger")
- `pfad`: `Makefile:16-21` (Kommentar + `include a-check.mk` ohne
  `GATE_CHECKS`-Aufnahme), `a-check.mk:15-17`,
  `harness/README.md:129` (Werkzeuge-Zeile, „nicht im `make gates`-Bündel")
- `befund`: Das Ziel ist aktiv und als eigener Aufruf lauffähig, hängt aber
  nicht an `GATE_CHECKS` — `make gates` (`record-gates: $(GATE_CHECKS)`,
  `harness/mk/enforce.mk`) führt es nicht und der Gate-Nachweis-Stempel
  belegt es nicht. Damit ist der [`ADR-0036`](../plan/adr)-Hochschalt-Trigger im Wortlaut
  („im Gate-Lauf") nicht erfüllt, obwohl der Slice-Plan behauptet, der Slice
  vollziehe ihn. Die Abweichung ist dokumentiert (Makefile-Kommentar,
  README-Zeile), nicht still — aber dokumentierte Abweichung ist keine
  Folge-ADR. Bewertung: bewusst gewählte Verkabelung, Planner/Architect-Frage;
  die Widerspruchsklasse (Proposed-ADR-Text vs. Implementierung) ist
  real, die Eskalation zur Sequenz steht erst an, wenn die Plan-Claim
  („vollzieht Hochschalt-Trigger") gegen die Verkabelung verteidigt wird.
- `verifizierbar`: ja — `GATE_CHECKS`-Akkumulation gegen die
  Include-Fragmente ist mechanisch prüfbar (a-check.mk trägt keinen
  `GATE_CHECKS +=`-Zweig)
- `klasse`: Gate außerhalb des Gate-Bündels ohne Folge-ADR

### F-2 — ADR-0036 nennt depguard; aktiviert ist a-check — ohne Folge-ADR

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0036`](../plan/adr) §Fitness Function (depguard, „statisches
  Import-Linting") · [`ADR-0039`](../plan/adr) §Konsequenzen/Fitness Function („depguard-Gate
  ist Folgepflicht")
- `pfad`: `docs/plan/adr/0036-architekturpruefung-ci.md:57-59, 68-70` gegen
  `a-check.mk:9` (gepinntes Release-Image `ghcr.io/pt9912/a-check`) und
  `.a-check.yml` (pfadbasierte Schichten-Edges)
- `befund`: Die aktivierte Maschinenprüfung ist a-check (pfadbasierte
  Hexagon-Edges, externes Release-Image), nicht das in [`ADR-0036`](../plan/adr)/[`ADR-0039`](../plan/adr)
  benannte depguard-Import-Linting. Die Divergenz ist Vorbestand (Vorbereitung
  `174ae63`, vor der Range), wird aber mit `01e981e` in-range operationalisiert
  und in der Commit-Message auf [`ADR-0036`](../plan/adr) bezogen — ohne dass ein Folge-ADR das
  Werkzeug nachschärft. Der Hochschalt-Trigger löst gegen den ADR-Text
  („Import-Linting-Gate") ohne Lesart-Entscheidung nicht.
- `verifizierbar`: nein — Urteil (Werkzeug-Äquivalenz), kein Gate
- `klasse`: ADR-Werkzeug-Nennung ersetzt ohne Folge-ADR

### F-3 — Plan-Defekt: Diff geht über §3 des Slice-Plans hinaus, Plan ohne Nachzug

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan `slice-001-bootstrap.md` §3 (Plan-Tabelle) · Modul 5
  („Wer später mitnimmt …, hat den Plan geändert, nicht nur ergänzt")
- `pfad`: `docs/plan/planning/in-progress/slice-001-bootstrap.md:109-115`
  gegen `Makefile:25` (image-hash-Extraktion), `Dockerfile:7-9`
  (Kommentar-Erneuerung), `harness/image-hash.txt` (neu)
- `befund`: Drei Berührungen fehlen in §3: die Umstellung der
  image-hash-Extraktion auf `containerimage.digest` (der alte Rezept-Ausdruck
  `grep -o 'sha256:…' | head -1` griff die erste sha256 der metadata-JSON —
  nicht notwendig den Image-Digest), die Dockerfile-Kommentar-Erneuerung und
  die Beleg-Datei selbst. Bewertung: Inhaltlich ist der Nachzug gedeckt —
  die Beleg-Datei ist in §1 Ziel genannt, der Rezept-Fix korrigiert die
  Erzeugung genau dieses Belegs ([`LH-QA-OPS-001`](../../spec/lastenheft.md)), und der Dockerfile-Kommentar
  beschrieb eine Bedingung, die dieser Slice selbst aufgelöst hat (stilles
  Stehenlassen wäre ein §3.7-Defekt gewesen). Der Defekt ist die Form: §3
  wurde nicht ergänzt, Plan und Diff driften — der Plan ist die
  Review-Prüfgrundlage (Modul 1), die Abweichung ist erst über den
  Commit-Text rekonstruierbar. Richtiger Ort: §3-Nachzug bzw. §7
  („Was ging anders als geplant") bei der Closure, nicht erst dieser Report.
- `verifizierbar`: ja — Datei-Menge des Diffs gegen die §3-Tabelle
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug

### F-4 — Dockerfile baut auf ADR-0038 (Superseded) als Begründungs-Anker

- `kategorie`: LOW
- `quelle`: ADR-Index ([`ADR-0038`](../plan/adr) Status Superseded, → [`ADR-0039`](../plan/adr)) ·
  Traceability-Klasse (Auftrag: jede [`ADR-0038`](../plan/adr)-Referenz ist ein Befund)
- `pfad`: `Dockerfile:17` („CGO aus ([`ADR-0038`](../plan/adr))")
- `befund`: Die Zeile referenziert die Superseded-ADR als direkte Begründung
  eines lebenden Build-Flags. Einstufung LOW, nicht HIGH (Klasse
  Traceability/ID-Schema), mit Begründung der Abweichung von der
  Auftrags-Lesart: (1) Vorbestand — Zeile entstand in `390a7b5`, vor der
  Review-Range; dieser Diff berührt die Datei, aber nicht diese Zeile;
  (2) [`ADR-0039`](../plan/adr) §Kontext erklärt ausdrücklich, dass „CGO-Ziel … aus [`ADR-0038`](../plan/adr)
  unverändert bestehen" bleibt — der Inhalt ist nicht verloren, nur der Anker
  zeigt auf die tote Datei, und ein folgender Leser übersieht die
  Nachfolge-Kette. Korrektur als Doku-Nachzug (Zeiger [`ADR-0038`](../plan/adr) →
  [`ADR-0039`](../plan/adr) §Kontext), kein Slice-Blocker.
- `verifizierbar`: ja — grep über den Bestand gegen den ADR-Index-Status
- `klasse`: Verweis auf superseded ADR

### F-5 — Transienter Beleg-Zwischenstand `image-hash.raw` ohne Ignorier-Schutz

- `kategorie`: LOW
- `quelle`: Maintainability (Beleg-Form, Modul 14 — `image-hash.txt` ist der
  Beleg, das `.raw`-Zwischenprodukt nicht)
- `pfad`: `Makefile:25` · fehlendes `.gitignore` (Repo-Wurzel)
- `befund`: Die Rezept-Kette erzeugt `harness/image-hash.raw` und entfernt sie
  erst nach zwei `&&`-Gliedern. Scheitert der Build oder ein grep-Glied,
  bleibt die Datei untracked im Baum; das Repo trägt kein `.gitignore`, die
  Fehlzustands-Datei ist damit committbar und würde als zweite Quelle neben
  `image-hash.txt` liegen.
- `verifizierbar`: ja — `git status` nach absichtlich fehlgeschlagenem
  `make image`
- `klasse`: Transienter Beleg ohne Vernachlässigbarkeits-Schutz

### F-6 — Fehlende Datei-Ende-Zeilen in neuen Quelldateien

- `kategorie`: LOW
- `quelle`: Maintainability (Stil, ohne semantische Auswirkung)
- `pfad`: `cmd/pg-change-feed/main.go:23`, `go.mod:3`
- `befund`: Beide neuen Dateien enden ohne Zeilenumbruch (git
  `\ No newline at end of file`); nachfolgende Anfügungen erzeugen
  irrelevante Diff-Rauschen in `go.mod`.
- `verifizierbar`: ja
- `klasse`: Fehlender Datei-Abschluss

### F-7 — Beleg-Datei `image-hash.txt` ohne Sensor (benannte Grenze)

- `kategorie`: INFO
- `quelle`: Implementer-Risiko (b) des Handoffs · harness/README §Sensors
  (kein Target gegen `image-hash.txt` eingetragen)
- `pfad`: `harness/image-hash.txt` · `Makefile:25`
- `befund`: Der Digest in `image-hash.txt` wird von keinem Gate gegen einen
  Rebuild verglichen; eine abweichende/gefälschte Beleg-Zeile ist nur manuell
  sichtbar. Als benannte Grenze aufgenommen, keine erwartete Aktion in diesem
  Slice: ob ein Abgleich-Sensor (Rebuild-Vergleich) den Aufwand trägt,
  entscheidet der Planner.
- `verifizierbar`: nein — Grenze, kein Defekt
- `klasse`: Beleg-Datei ohne Gate-Deckung

### F-8 — Report-Gerüst referenziert, aber nicht vorhanden

- `kategorie`: INFO
- `quelle`: `.harness/skills/reviewer.md` §Output-Schema („Report-Gerüst:
  `docs/reviews/review-report.template.md`")
- `pfad`: `docs/reviews/` (nur `.gitkeep` + Eröffnungs-Report)
- `befund`: Die referenzierte Vorlage existiert nicht; dieser Lauf folgt der
  Form des Eröffnungs-Reports. Zuständigkeit: Architect/Implementer des
  Gerüsts, kein Befund gegen den geprüften Diff.
- `verifizierbar`: ja — Dateiexistenz prüfbar
- `klasse`: Referenziertes Artefakt fehlt

### F-9 — image-cve-Aktivierungsfenster mit diesem Slice eingetreten

- `kategorie`: INFO
- `quelle`: `harness/README.md:136-137` („Nicht behauptet (geplant): make
  image-cve … ab dem ersten `make image`-Lauf")
- `pfad`: `harness/README.md:136-137`
- `befund`: Mit diesem Slice lief der erste `make image`-Lauf (Beleg-Datei
  liegt vor); die in der Zeile genannte Bedingung ist eingetreten, das Target
  bleibt unter „Nicht behauptet". Keine Aktion im Slice-Plan — Hinweis an den
  Planner, die Zeile zu präzisieren oder das Target zu addieren.
- `verifizierbar`: nein
- `klasse`: Aktivierungsbedingung eingetreten, Zeile unverändert

---

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff** — kein ADR-Verstoß
  gegen eine Accepted-ADR ([`ADR-0039`](../plan/adr) eingehalten: Baum startet unter
  `cmd/pg-change-feed/`, Modul-Pfad `github.com/pt9912/pg-change-feed`
  konsistent; [`ADR-0040`](../plan/adr) korrekt unberührt — kein Clock-Bezug im Bootstrap),
  keine Gate-Suppression, kein Sicherheits- oder Korrektheits-Finding im
  kritischen Pfad, keine Norm nur im Template-Kommentar, kein
  Chronik-tragendes Zustandsfeld, kein Docker-only-Verstoß (Build/Templates
  laufen im Container; `a-check` via `--network none` + ro-Mount, gepinnter
  Digest), kein Zwei-Quellen-Drift über den geprüften Umfang hinaus
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code (§3.7)** —
  `main.go` (Paket-Kopf, `version`-Konstante), Makefile-Zeilen 16–21/23–25,
  Dockerfile-Zeilen 7–9: durchgängig Zustand/Kopplung/Abgrenzung im Indikativ,
  keine verworfene Alternative, kein abwesender Text; die entfernten
  Bedingungs-Kommentare („aktivieren mit dem ersten src/-Slice") beschrieben
  eine Bedingung, die dieser Slice erfüllt — Ersetzungs-Fassung steht korrekt
- geprüft, ohne Befund: **Traceability der vier Commits** — jeder trägt
  mindestens eine `LH-*`- oder `ADR-*`-Kennung (`ADR-0039`/`LH-QA-POR-003` ·
  `ADR-0036` ×2 · `LH-QA-OPS-001`); alle genannten IDs existieren; keine
  nicht in MR-000 deklarierten Präfixe; keine `ADR-0038`-Referenz im
  Commit-Text oder im in-range Diff (F-4 betrifft Vorbestand-Zeile)
- geprüft, ohne Befund: **Spec-Stratum** — kein Erweiterungs-Befund; der
  Diff berührt kein Spec-Dokument; Struktur-IDs (`ARC-007`) nur im Slice-Plan,
  nicht in Commits (`AGENTS.md` §5 eingehalten)
- geprüft, ohne Befund: **`.a-check.yml`/`a-check.mk` gegen den entstehenden
  Baum** (unverändert, laut Plan §3 geprüft) — `composition_root` deckt
  `cmd/**`, Schichten-Globs zeigen auf noch nicht existierende `internal/**`-
  Pfade (erwartbar in diesem Slice); netzlos und read-only wie in der
  Werkzeuge-Zeile behauptet; das Deckungsverhalten des ersten Laufs ist
  Verifier-Gegenstand (DoD-Punkt 3), nicht Reviewer
- geprüft, ohne Befund: **harness/README-Konsistenz** — Werkzeuge-Zeile
  `make a-check` und Entfernung aus „Nicht behauptet" decken sich mit dem
  Makefile-Include (Ziel existiert, Zeile existiert, keine zweite Quelle für
  den Aktivierungs-Zustand); `make image`-Zeile trägt die [`ADR-0039`](../plan/adr)-Bindung
  laut Plan-DoD
- geprüft, ohne Befund: **Beleg-Datei-Form** — `harness/image-hash.txt`
  getrackt, einzeilig, Inhalt in `sha256:<64 hex>`-Form = `containerimage.digest`
  nach dem gefixten Rezept; „Beleg, kein Wiederholungs-Schlüssel" gemäß
  Dockerfile-Kopf gewahrt
- geprüft, ohne Befund: **Dockerfile-Build-Vertrag gegen neuen Baum** —
  `COPY go.mod go.sum` mit leerem `go.sum` und `go mod download/verify`
  verifiziert die leere Modulliste; Build-Pfad `./cmd/pg-change-feed` existiert;
  CGO_ENABLED=0 konsistent mit dem (über [`ADR-0039`](../plan/adr) weitergelten) CGO-Ziel

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Gate außerhalb des Gate-Bündels ohne
Folge-ADR · ADR-Werkzeug-Nennung ersetzt ohne Folge-ADR · Plan-Erweiterung
ohne Plan-Nachzug · Verweis auf superseded ADR · Transienter Beleg ohne
Vernachlässigbarkeits-Schutz · Fehlender Datei-Abschluss · Beleg-Datei ohne
Gate-Deckung · Referenziertes Artefakt fehlt · Aktivierungsbedingung
eingetreten

## Verdikt

**Merge-blockierend:** nein — der Diff ist inhaltlich schlüssig, Docker-only
und kommentar-regelkonform; alle MEDIUMs sind Plan-/Verkabelungs-Fragen, keine
Code-Defekte.

**Blockierend für Closure:** ja, in zwei Punkten, bevor der Slice nach
`done/` geht: (1) F-1/F-2 — die Slice-Plan-Behauptung, der Slice vollziehe den
[`ADR-0036`](../plan/adr)-Hochschalt-Trigger, hält in dieser Verkabelung nicht; der Planner
entscheidet GATE_CHECKS-Aufnahme oder Plan-/ADR-Korrektur, und der Architect
entscheidet die depguard-vs-a-check-Frage als Folge-ADR **vor** dem
Hochschalten von [`ADR-0036`](../plan/adr) auf Accepted. (2) F-3 — §3 des Slice-Plans bekommt
den Nachzug (Dockerfile, image-Rezept, Beleg-Datei), damit Plan und Diff vor
der Closure wieder deckungsgleich sind.

**Übergabe:** Findings an den Implementer (F-3: Plan-Nachzug §3 bzw. §7
„Was ging anders"; F-4–F-6: Richtung Doku-Nachzug/`.gitignore`/Abschluss,
kein Code-Input vom Reviewer). F-1/F-2 gehen an Planner/Architect als
Verkabelungs-/ADR-Frage. Die Finding-Klassen gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler. DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11) — insbesondere `make a-check` grün und
`make gates` grün als beobachtbare Belege.