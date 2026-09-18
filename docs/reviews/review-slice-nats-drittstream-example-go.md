# Review — `slice-nats-drittstream-example-go`

**Review-Art:** Code und Doku gegen Plan, Entscheidung und Hard Rules
(`.harness/skills/reviewer.md`, Modul 10).

**Gegenstand:** `git diff ee6b2ea0..96e54007` (5 Commits: 2× reine `git
mv`/Ruhe-Marker, 1× feat, 1× chore/Digest, 1× DoD-Beleg-Nachzug) gegen den
Slice-Plan zu `slice-nats-drittstream-example-go` (bei diesem Lauf unter
`in-progress/`) und [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(`Accepted`).

**Datum:** 2026-09-18.

**Selbst ausgeführt (Exit-Code jeweils direkt und ungepiped, `AGENTS.md`
§3.9):**

- `make gates` — Exit 0 (d-check 720 Dateien, 0 Befunde · coverage-gate
  82,80 % ≥ 80 % · commit-traceability OK · a-check 0 Befunde ·
  generated-sync OK).
- `docker build -f Dockerfile --target build . && docker run --rm <img> ls
  /src` — zeigt kein `examples/` im Build-Kontext des Root-`Dockerfile`
  (Beleg zu F-1).
- `docker run golang:1.27-alpine sh -c "go build ./examples/... && go vet
  ./examples/... && go test ./examples/nats-stream-client/..."` — Exit 0;
  `gofmt -l` leer.
- `docker build -f examples/Dockerfile --target runtime-nats-stream .` —
  Exit 0, reale Probe der neuen Runtime-Stufe.

**Verdikt (erster Durchlauf):** 1 HIGH, 1 MEDIUM offen → Fixrunde nötig.
**Verdikt (nach Fixrunde, siehe Nachtrag 1):** alle vier Findings adressiert,
F-2 auf LOW herabgestuft (funktional gelöst, ein benannter, nicht
blockierender Rest). **Freigegeben für Handoff an die Verifier-Rolle.**

---

## HIGH

### F-1 — Digest-Commit behauptet eine Build-Kontext-Kausalität, die das eigene `.dockerignore` widerlegt

- **klasse:** „Beleg trägt seinen Satz nicht" (Build-Kontext-Kausalbehauptung
  widerlegt durch `.dockerignore`-Whitelist)
- **pfad:** `harness/image-hash.txt` (Commit `7ea44409`)
- **befund:** Commit `7ea44409` begründet den Digest-Wechsel mit
  „examples/-Aenderungen liegen im Build-Kontext des Root-Dockerfiles
  (COPY . .)". Real geprüft: das Root-`.dockerignore` schließt per
  `*`-Whitelist alles außer `cmd/`, `internal/`, `gen/`, `proto/`, `go.mod`,
  `go.sum`, `tools/coverage-gate.sh`, `tools/schema/nacharbeit-roles.sql`
  aus — `examples/**` gehört nicht dazu. Ein realer `docker build -f
  Dockerfile --target build . && docker run --rm <img> ls /src` zeigt: kein
  `examples/` im Build-Kontext. Keiner der 5 Commits dieses Diffs berührt
  einen Pfad aus der Whitelist. Der beobachtete Digest-Wechsel ist der in
  `ADR-0044` selbst dokumentierte lauf-/builder-gebundene
  Nicht-Determinismus, nicht examples/-Inhalt.
- **verifizierbar:** ja — siehe „Selbst ausgeführt" oben.

## MEDIUM

### F-2 — Zitat-Korrektur an `done/`-Record ohne ADR-0073-Zitat in der Commit-Message

- **klasse:** [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)-Belegform bei Zitat-Korrektur an Record fehlt
- **pfad:** `docs/plan/planning/done/slice-nats-drittstream-core.md`
  §„Folge-Slices" (Commit `106bba7d`)
- **befund:** Der Fix selbst ist korrekt in Scope — ein gebrochener
  Markdown-Link (`../open/slice-….md`, gebrochen durch den `git mv` nach
  `in-progress/`) auf zwei Folge-Slices wird durch eine
  Inline-Code-Kennungs-Zitierung ersetzt; Aussage und Referent bleiben
  unverändert (`ADR-0073` Klasse 1, korrekt angewendet). Die
  Commit-Message nennt `ADR-0073` aber nicht, obwohl Punkt 3(a) genau das
  als Beleg-Form fordert: „Die korrigierende Commit-Message nennt
  `ADR-0073`". Records (`done/`) tragen keine §Geschichte-Zeile — bei ihnen
  ist die Commit-Message selbst der Beleg (Punkt 3(b)); ohne die Nennung
  ist die Änderung an einem abgeschlossenen Record aus dem Commit allein
  nicht als legitimierte Zitat-Korrektur erkennbar.
- **verifizierbar:** ja — Commit-Message von `106bba7d` enthält
  „[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)"
  nicht.

## LOW

### F-3 — `main.go`-Docstring zitierte die durch ADR-0098 abgelöste Startform (4. Fundstelle)

- **klasse:** Docstring zitiert abgelöste Startform (Vorbild-Drift)
- **pfad:** `examples/nats-stream-client/main.go:6`
- **befund:** Alle vier Go-Beispiele (`grpc-client`, `sse-client`,
  `nats-client`, jetzt auch `nats-stream-client`) trugen im
  Package-Doc-Kommentar noch „Startform ist `go run
  ./examples/<name>`" — die durch `ADR-0098` abgelöste Form. Der neue
  Client übernahm das treu vom Vorbild `grpc-client`, ohne die bereits
  dreifach bestehende Inkonsistenz zu korrigieren.
- **verifizierbar:** ja — `grep -n "Startform ist" examples/*/main.go`.

### F-4 — `.a-check.yml`-Kommentar zur `examples`-Gruppe zählte den vierten Zugriffsweg nicht mit auf

- **klasse:** Beschreibender Kommentar nicht um neue Instanz ergänzt
- **pfad:** `.a-check.yml:23-36`
- **befund:** Der beschreibende Kommentar zur `examples`-Layer-Gruppe zählt
  die Beispiel-Programme namentlich auf, nannte
  `examples/nats-stream-client` nicht, obwohl strukturell identisch (nats.go,
  keine domain/ports/app/adapters-Importe). Keine Aussage falsch, nur
  unvollständig; `a-check` selbst braucht keine Änderung.
- **verifizierbar:** nein — reine Lese-Prüfung.

## INFO

- **I-1** — `SPEC-023` „Sprachen und Umfang" bleibt bei „vier"
  Zugriffs-Oberflächen (`spec/pflichtenheft.md:463`) — mit dem fünften Weg
  (bislang nur Go) im Geiste veraltet. Weder dieser Slice-Plan noch der
  Folge-Slice-Plan noch `ADR-0100`s Folgepflicht-Liste nennen eine
  `SPEC-023`-Aktualisierung als Liefer-Punkt — Planner-Zuständigkeit für
  einen benannten Abschluss.
- **I-2** — `examples/README.md` §Abgrenzung nennt `natssub` (welle-15) und
  `natsstreamsub` (`slice-nats-drittstream-core`) weiterhin nicht —
  vorbestehende Lücke, von diesem Diff nicht verursacht (dieser Slice fasst
  nichts unter `tools/harness/**` an), korrekt nicht mitgezogen.
- **I-3** — Aggregat-Aussage „vier Zugriffsarten × drei Sprachen"
  (`examples/README.md:3-4`, Handbuch-Kopf) bleibt während der
  Übergangszeit bewusst unverändert stehen — explizit im Plan als
  Aufschub mit benannter Adresse (`slice-nats-drittstream-example-csharp-kotlin`)
  deklariert.

---

## Negativbefunde (geprüft, ohne Befund)

- Hard Rule 3.3 (git mv + Inhalt = zwei Commits) — alle 5 Commits einzeln
  per `git show --stat` geprüft; `4e713a58` ist reiner `git mv` (0/0), jede
  Inhaltsänderung liegt in einem eigenen Commit.
- `ADR-0100`-Konformität in `main.go`/`format.go` — Subjekt-Namensraum
  (`cdc.stream.>`, Root-Wildcard, Teilfrage 2), Fire-and-Forget/Kein-Replay-
  Framing (Teilfrage 3), Token-Auth (`nats.Token`, Teilfrage 4) korrekt
  wiedergegeben.
- `SPEC-023` Import-Grenze — kein `internal/`-Import in
  `examples/nats-stream-client/**`.
- Nachrichtenschema-Konsistenz `format.go`/`format_test.go` gegen
  `internal/adapters/driven/natsstream/publisher.go` — die Testbehauptung
  „fehlendes Row Image → JSON `null`" stimmt mit `rowImage()`s realem
  Verhalten überein.
- `examples/Dockerfile` neue Stufe `runtime-nats-stream` — real gebaut.
- `harness/mk/examples.mk` — `$(filter …)`/`$(error …)`-Erweiterung um
  `nats-stream`, **nur** die `example-run-go`-Zeile berührt.
- DoD-Beleg-Nachzug (`96e54007`) — der eingetragene Smoke-Test-Output
  (`change_id=804-1`) deckt sich mit dem DoD-Punkt.
- `docs/user/benutzerhandbuch.md` Versionshistorie — Version 1.28→1.29 und
  neue Zeile im selben Diff wie die inhaltliche Änderung.
- Reconciliation-Register/Beobachtungs-Register/§6-Risiken-Ausgänge/„drei
  Paarungen" — DoD-Punkte korrekt noch `[ ]`, Closure noch nicht erfolgt.
- Traceability — jeder der 5 Commits nennt mindestens eine `LH-*`/`ADR-*`-
  Kennung im Betreff.

---

## Nachtrag 1 — Fixrunde und Re-Check (Commits `b67a1a63`, `bce81ca6`, `20562f44`)

**F-1 — RESOLVED.** `harness/image-hash.txt` per `b67a1a63` auf den
Vor-Slice-Wert (`sha256:2c0178a7…`) zurückgesetzt; die neue Commit-Message
trägt die korrekte Kausal-Erklärung ([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md), kein Build-Kontext-Bezug).
Unabhängig nachgeprüft: Revert korrekt und vollständig, `.dockerignore`-
Whitelist bestätigt keinen `examples/**`-Pfad.

**F-2 — funktional resolved, ein benannter, nicht blockierender Rest.**
`bce81ca6` ist ein reiner leerer Commit (`git diff-tree` bestätigt: keine
Datei geändert), dessen Message `ADR-0073` explizit nennt und `106bba7d`
als korrigierten Commit benennt. Buchstäblich erfüllt das Punkt 3(a) nicht
ganz (3(a) verlangt die Nennung in „der korrigierenden Commit-Message"
selbst — `106bba7d`s eigene Message bleibt unverändert, kein Amend/Rebase,
wie explizit gewollt); funktional ist die Lücke geschlossen: `git log
--grep=ADR-0073` findet `bce81ca6`, das den Referenten namentlich trägt —
ein Auditor findet die Zuordnung ohne Rätselraten. Gegeben die Randbedingung
„kein Amend/Rebase ohne explizite Anweisung" ist ein Folge-Commit die
einzige nicht-destruktive Reparaturform; `ADR-0073` selbst sieht diesen
Nachtrags-Fall nicht ausdrücklich vor — das ist eine Lücke in der ADR, keine
im Fix. **Downgrade MEDIUM → LOW**, nicht blockierend. Falls das Repo
künftig eine kanonische Antwort will, ob ein Nachtrags-Commit 3(a) generell
erfüllt, wäre das ein eigener, kleiner `ADR-0073`-Klärungspunkt.

**F-3/F-4 — RESOLVED, korrekt gescoped.** `20562f44` korrigiert
ausschließlich `examples/nats-stream-client/main.go`s eigene Docstring-Zeile
(jetzt `make example-run-go SURFACE=nats-stream`); die drei Geschwister-
Vorkommen derselben Drift (`grpc-client`, `http-client`, `nats-client`,
`sse-client`) blieben bewusst unverändert. `.a-check.yml`s
`examples`-Gruppenkommentar zählt jetzt `nats-stream-client` mit auf, im
selben Stil wie die bestehenden Einträge.

**Hard Rule 3.3** — geprüft per `git show --stat --find-renames` für alle
drei neuen Commits: keiner enthält eine Umbenennung, kein Verstoß.

**`make gates`** — selbst erneut ausgeführt, Exit 0 (Coverage 82,70 %,
d-check 0 Befunde, a-check 0 Befunde, generated-sync OK,
commit-traceability OK).

**Gesamtverdikt:** F-1, F-3, F-4 vollständig geschlossen; F-2 in der Sache
geschlossen mit einem benannten, nicht-blockierenden Rest (LOW statt
MEDIUM). Keine neuen HIGH/MEDIUM-Funde in den drei nachgeprüften Commits,
kein neuer §3.3-Verstoß. **Freigegeben für Handoff an die Verifier-Rolle** —
keine weitere Fixrunde am Implementer nötig.
