# Review-Report: slice-097 — Umzug der Vertragsfläche (Protobuf-Stub an öffentlichen Pfad) · 2026-09-17

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen** (Maintainability). DoD-/
Spec-Konformität, §6-Ausgänge, Closure-Notiz und die drei Paarungen sind **nicht** Gegenstand
dieses Reports (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff 4f39cd5..068f26e` — drei Commits:

- `1c1d937` — feat(grpc): erzeugte Vertragsfläche an öffentlichen Pfad umgezogen
- `57d2566` — docs(coverage): Coverage-Fläche und Bau-Kontext um `gen/...` erweitert
- `068f26e` — chore(image): Digest nach Bau-Kontext-Änderung (`gen/` hinzu) neu erfasst

17 Dateien, davon 3 Umbenennungen (`R098`/`R100`/`R099`,
`internal/adapters/driving/grpc/streamv1/*` → `gen/cdc/stream/v1/*`).

**Skill:** `.harness/skills/reviewer.md` (Accepted) · **Datum:** 2026-09-17.

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-097-umzug-vertragsflaeche.md` (§1–§8, vollständig)
- `docs/reviews/architect-verdict-slice-097-coverage-gegenstand.md` (Commit `1eb3feb`)
- `ADR-0076` (Umzugs-Folgepflicht, §3 Schnitt-Empfehlung 1), `ADR-0068` (die zurückgenommene
  `tooling → adapters`-Kante), `ADR-0071`/`ADR-0082` (Coverage-Messgegenstand-Grundsatz),
  `ADR-0044` (Image-Beleg-Semantik), `ADR-0085` (Bau-Kontext-Ausnahme, Festlegung 3)
- `AGENTS.md` §3.3 (git mv + Inhaltsänderung), §3.12 (Herkunft von Aussagen), §3.13 (Arbeit
  überholt stehenden Träger), §5 (Traceability), `harness/conventions.md` (MR-000 ID-Schema)

**Mess-Grundlage:** eigene Läufe im Arbeitsbaum am Stand `068f26e` (netzlos, gepinnte
Toolchain/Images): `make a-check` (Exit 0, „gesamt: 0 Befund(e)"), `make generated-sync`
(Exit 0, byte-gleich), `make docs-check` (Exit 0, 799 Dateien, 0 Befunde),
`make commit-traceability` (Exit 0, „Betreffs ohne Struktur-ID"), `make coverage-gate`
(Exit 0, gedruckt „83.40%" — deckt sich mit dem im Diff behaupteten Band 83,3–83,4 %). Alle
Exit-Codes direkt und ungepiped gelesen (`AGENTS.md` §3.9). `make image` nicht erneut gefahren
(Zeit-/Ressourcen-Grund) — der Digest-Wechsel in `068f26e` ist plausibel, aber nicht durch einen
eigenen Nachbau verifiziert (siehe Negativbefund).

---

## Findings

### F-1 — Commit `1c1d937` verbindet reine Umbenennung und Inhaltsänderung in einem Commit

- `kategorie`: **HIGH**
- `quelle`: `AGENTS.md` §3.3 „git mv + Inhaltsänderung = zwei Commits" (Hard Rule)
- `pfad`: Commit `1c1d937` — betrifft u. a. `gen/cdc/stream/v1/changestream.pb.go`,
  `gen/cdc/stream/v1/changestream_test.go`, `.a-check.yml`,
  `internal/adapters/driving/grpc/{server,server_test,interceptor_test}.go`,
  `tools/harness/grpcclient/main.go`
- `befund`: `git diff --name-status -M` weist die drei generierten Dateien als Umbenennung aus
  (`R098`, `R100`, `R099`), aber zwei der drei tragen im selben Commit zusätzlich eine
  Inhaltsänderung — `changestream.pb.go` (die im Raw-Descriptor eingebettete
  `go_package`-Zeichenkette ändert sich von `.../internal/adapters/driving/grpc/streamv1;streamv1`
  auf `.../gen/cdc/stream/v1;streamv1`) und `changestream_test.go` (Import-Pfad geändert). Zusätzlich
  ändert derselbe Commit `.a-check.yml` (neue Gruppe `contract`, neue Kanten, Rücknahme der
  `tooling→adapters`-Kante) und vier weitere Importstellen. §3.3 verlangt für den Regelfall „Datei
  verschoben und Inhalt umgeschrieben" **zwei** Commits — Move-Commit rein, danach die
  Inhaltsänderung. Hier liegen Move und Inhaltsänderung (samt Folgeänderungen) in einem Commit.
  Dass Git die Rename-Detection trotzdem traf (98–100 % Similarity), mindert die Regelverletzung
  nicht — die Begründung der Regel (Rename-Detection unter der 50 %-Schwelle) ist der Grund für
  die Regel, nicht ihre Bedingung.
- `verifizierbar`: nein — kein Gate prüft Commit-Struktur; nur durch `git show --stat`/`git diff -M`
  je Commit sichtbar.
- `klasse`: „git-mv-und-Inhalt-in-einem-Commit"

## Negativbefunde

- geprüft, ohne Befund: `proto/cdc/stream/v1/changestream.proto` — genau eine Zeile geändert
  (`option go_package`), kein Nachrichten-/Feld-Inhalt berührt (§1 Out-of-Scope eingehalten).
- geprüft, ohne Befund: die vier Importstellen (`server.go`, `server_test.go`,
  `interceptor_test.go`, `tools/harness/grpcclient/main.go`) — jeweils minimale
  Ein-Zeilen-Änderung, korrekt auf `gen/cdc/stream/v1` umgestellt; `make test`/`make a-check`
  real grün.
- geprüft, ohne Befund: `.a-check.yml` — neue Gruppe `contract` (`gen/**`) mit Kanten
  `adapters→contract` und `tooling→contract`; die Rücknahme der `tooling→adapters`-Kante ist
  im Kommentar explizit als „ZURÜCKGENOMMEN" mit Begründung benannt, nicht kommentarlos
  gelöscht; der alte `tooling`-Gruppenkommentar (verwies auf den alten Pfad) ist korrigiert.
  `make a-check` real Exit 0, „0 Befund(e)".
- geprüft, ohne Befund: alter Baum `internal/adapters/driving/grpc/streamv1/` vollständig
  entfernt, keine Duplikate; `changestream_test.go` zog mit dem generierten Baum nach
  `gen/cdc/stream/v1/` mit (Package `streamv1_test`, Import korrekt) — Bedingung für den im
  Architect-Verdikt genannten Nenner 1936.
- geprüft, ohne Befund: `Dockerfile` Stufe `coverage` — `./gen/...` in der `go list`-Zeile
  ergänzt, Filter-Regex unverändert (`(^|/)(postgresstorage|postgresack|replication/receive)$`);
  `make coverage-gate` real 83.40 % ≥ 80 %, Exit 0.
- geprüft, ohne Befund: `.dockerignore` — `!gen/` als eigene Zeile neben `!cmd/`/`!internal/`,
  ohne `ADR-0085`-Zitation (wäre eine Fehl-Zitation gewesen, siehe Architect-Verdikt §6.B) —
  korrekt vermieden.
- geprüft, ohne Befund: `harness/sensors/coverage-gate.md`, `harness/mk/coverage.mk`,
  `harness/README.md`, `AGENTS.md` §4, `harness/sensors/generated-sync.md` — alle konsistent auf
  Pfad `gen/cdc/stream/v1` und Fläche `./internal/...+./cmd/...+./gen/...` nachgezogen; die Zahl
  **1936** steht in `harness/sensors/coverage-gate.md` als datierte Messung („Lauf `slice-097`",
  mit Verweis auf den Architect-Verdikt), die vorherige Zahl 1903 bleibt daneben ausdrücklich als
  datierte Messung des `slice-085`-Laufs stehen, nicht als überschriebener Ist-Stand
  (`AGENTS.md` §3.12 Instanz A eingehalten).
- geprüft, ohne Befund: `harness/image-hash.txt` — Digest-Commit `068f26e` vorhanden, mit
  Begründung (Bau-Kontext-Änderung durch `.dockerignore`, `ADR-0044` Punkt 3 zitiert).
- geprüft, ohne Befund: Out-of-Scope-Disziplin (§1 des Plans) — kein gRPC-Client gebaut, kein
  zweites Go-Modul, kein neues Gate; Diff bleibt auf Umzug + seine Träger beschränkt.
- geprüft, ohne Befund: Commit-Traceability — alle drei Commit-Messages tragen `ADR-*`-Bezüge
  (`ADR-0076`/`ADR-0068`, `ADR-0071`/`ADR-0076`, `ADR-0044`/`ADR-0076`), keine `SPEC-*`/`ARC-*`-
  Kennung im Betreff; `make commit-traceability` real Exit 0.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md`, `docs/user/e2e-abdeckung.md`,
  `harness/sensors/a-check.md`, `harness/sensors/db-adapter-coverage.md`,
  `tools/harness/db-coverage.sh`, `.github/**` — keine Erwähnung des alten oder neuen Pfads,
  kein Nachzug fällig.
- geprüft, ohne Befund: ADR-Immutabilität — `ADR-0076`/`ADR-0082`/`ADR-0090`/`ADR-0068` unverändert
  im Diff; sie nennen den alten Pfad an mehreren Stellen weiter, das ist als historischer Beleg
  zulässig (keine Accepted-ADR wurde inhaltlich überschrieben).
- geprüft, ohne Befund: `make image` nicht selbst nachgebaut (Ressourcengrund) — der behauptete
  Digest-Wechsel ist plausibel (Bau-Kontext geändert), aber ungeprüft; siehe Grenze oben.

## Beobachtung außerhalb des Diffs (kein Fund im engeren Sinn, aber §3.13-relevant)

`docs/plan/planning/next/slice-095-beispiel-clients-drei.md:149–152` behauptet als Präsens-Fakt:
„die öffentliche Vertragsfläche (`gen/cdc/stream/v1`) existiert nicht (`no required module
provides package …`, Exit 1)". Das war zum Zeitpunkt der `in-progress→next`-Rückführung von
`slice-095` wahr und ist es seit diesem Diff nicht mehr — `gen/cdc/stream/v1` existiert jetzt
real. Die Datei liegt in `next/` (lebender Träger, kein `done/`-Record, keine ADR) und ist damit
kein Fall der zulässigen ADR-Ausnahme aus Punkt 8 der Prüf-Anweisung. Das Plan-eigene §6 des
Slice benennt genau diese Risikoklasse bereits ausdrücklich („Ein Träger wird überholt, den
dieser Slice nicht anfasst … `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, 3×") mit noch offenem
Ausgang — ich werte dies deshalb nicht als eigenständiges HIGH-Finding gegen den vorliegenden
Diff (der Such-/Melde-Schritt ist laut Plan der Closure vorbehalten), sondern hänge den
konkreten Fund hier an, damit er bei der Closure nicht neu gesucht werden muss.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** git-mv-und-Inhalt-in-einem-Commit

## Verdikt

**Merge-blockierend:** nein — die eine HIGH betrifft die Commit-**Struktur** (Historie/
Nachvollziehbarkeit), nicht Korrektheit, Sicherheit oder eine ADR-Verletzung im Code selbst; Git
hat die Umbenennung trotzdem korrekt erkannt (98–100 % Similarity), sodass `git log --follow`
nicht bricht. Für **diesen** Diff braucht es deshalb **keine Fixrunde** (die Historie ist bereits
gepusht bzw. wird es gleich sein; ein nachträgliches Umschreiben der Commits wäre ein größerer
Eingriff als der Nutzen rechtfertigt) — die Regel wird für künftige Umzugs-/Rename-Slices als
Steering-Loop-Kandidat vorgemerkt (Klassen-Zeile oben), nicht rückwirkend am Repo korrigiert.

**Übergabe:** Diese eine Finding-Klasse geht in die Slice-Closure §7 und von dort in den Zähler
des Beobachtungs-Registers. Die §3.13-Beobachtung zu `slice-095` (oben) gehört in den
Such-/Melde-Schritt der Closure (§6-Risiko-Ausgang „Ein Träger wird überholt …"). Da hier keine
Fixrunde am Implementer folgt, zieht dieser Report selbst die DoD-Checkbox „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" im Slice-Plan nach (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde). Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg
nicht wieder gelesen; DoD-/Spec-Konformität prüft der Verifier separat.
