# Review-Report: schema-rollout-zentrale-idempotenz-wache — 2026-09-18

**Review-Art:** Code — geprüft gegen Slice-Plan
(`docs/plan/planning/in-progress/schema-rollout-zentrale-idempotenz-wache.md`),
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) und
`AGENTS.md` Hard Rules (nicht gegen DoD — Verifier-Aufgabe).

**Gegenstand:** Commits `b41cc381` + `588140df` (Fixrunde innerhalb desselben
noch offenen Slice), Diff-Range `252f1e7f..588140df`, als ein effektiver Diff
gelesen.

**Skill:** `.harness/skills/reviewer.md` (Stand dieses Repos, ohne separates
Versions-Tag) · **Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/schema-rollout-zentrale-idempotenz-wache.md`
  (§1, §2, §3, §6 vollständig gelesen)
- [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
- `docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/state.md`
  und `evidence/*`
- `AGENTS.md` §3.3, §3.5, §3.7, §3.9, §3.12, §3.13
- `d-migrate schema migrate --help` real gegen das im Makefile gepinnte
  Image (`D_MIGRATE_IMAGE`) ausgeführt, zur Prüfung der
  `--allow-destructive`-Semantik

---

## Findings

### F-1 — Vorher/Nachher-Chronik im Produktionscode-Kommentar (Makefile + guard.go)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klassen) i. V. m. dem
  Reviewer-Skill-HIGH-Punkt „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar"
- `pfad`: `Makefile:197` („statt beim bloßen Überspringen von --execute
  verlustig zu gehen … kam bei einem reinen Skip nicht zurück"),
  `tools/schema/rolloutguard/guard.go:35-41` (identisches Muster:
  „ein reines Überspringen von `--execute` ließ … nicht zurückkommen")
- `befund`: Beide Kommentare begründen das aktuelle Verhalten, indem sie das
  Verhalten der **verworfenen** Vorgänger-Implementierung (Skip von
  `--execute` bei bekannten Blockern, commit `b41cc381`) narrativ referenzieren,
  statt ausschließlich den geltenden Ist-Zustand indikativ zu beschreiben und
  auf [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) /
  einen `· seit <slice>`-Anker zu verweisen. Das ist exakt das
  in AGENTS.md §3.7 benannte Muster („Die Abwägung gehört in die ADR, die
  Historie in git") und die im Reviewer-Skill bereits viermal belegte Klasse
  `BEO-PGC/slice-chronik-in-code-kommentar` — hier ein fünftes/sechstes
  Auftreten, in zwei Dateien derselben Fixrunde.
- `verifizierbar`: nein (kein Gate; Textmuster-Sensor wurde laut Skill bereits
  geprüft und verworfen)
- `klasse`: Slice-Chronik/Vorher-Nachher-Sprache in Produktionscode-Kommentar

### F-2 — AGENTS.md §3.13 verletzt: gestrichenes §3.14 bleibt als lebende Analogie-Referenz in einer `Accepted`-ADR stehen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.13 (Hard Rule — „Eine Arbeit, die eine
  beschriebene Eigenschaft bewegt, zieht ihre Träger nach")
- `pfad`: [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md):200
  („… wird bewusst nicht mit dieser ADR miterledigt (siehe
  Re-Evaluierungs-Trigger), analog zu `AGENTS.md` §3.14s Umgang mit einer
  erkannten, aber vertagten saubereren Lösung.")
- `befund`: Commit `b41cc381` streicht `AGENTS.md` §3.14 ersatzlos.
  [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md)
  (`Status: Accepted`, `Datum: 2026-09-18`, committet **vor** diesem Slice —
  siehe `git log` `e4edba5a` vor `252f1e7f`) stützt eine eigene
  Design-Entscheidung ausdrücklich mit einer Analogie auf §3.14 als
  bestehendes Muster. Diese Referenz ist ein wörtlicher, grep-barer
  Symbolname (`§3.14`) — genau die Klasse, die §3.13 selbst als „zuverlässig
  treffbar" benennt (Abgrenzung zu Zeilen-Lokatoren) — und wurde vom
  Such-/Nachzugslauf dieses Slice nicht gefunden: Weder Commit-Message noch
  Slice-Plan §3/§7 erwähnen
  [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md). Da
  diese ADR `Accepted` ist (§3.5, Immutabilität), ist eine stille Inhaltsänderung ohnehin nicht der
  richtige Weg — die Regel verlangt an dieser Stelle explizit **Meldung**
  statt stiller Mitänderung oder Schweigen; beides ist nicht erfolgt.
- `verifizierbar`: nein (kein Gate — `docs-check`s Anchor-Prüfung deckt
  Fragment-Links, nicht Prosa-Erwähnungen einer Abschnittsnummer ohne Link)
- `klasse`: Träger nicht nachgezogen (§3.13) — stale Analogie-Referenz in immutabler ADR

### F-3 — Blanket-`--allow-destructive` deckt keine zwischen Precheck und `--execute` neu entstandene destruktive Operation

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Korrektheit des Wache-Designs
  ([`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md))
- `pfad`: `Makefile:208-217` (`schema-rollout`-Target: zwei getrennte
  `docker run … schema migrate`-Aufrufe gegen dasselbe `SCHEMA_TARGET`)
- `befund`: `--plan-only` und `--execute --allow-destructive` sind zwei
  unabhängige, sequenzielle `schema migrate`-Läufe, die jeweils ihren Plan
  **frisch gegen den aktuellen Zielzustand** berechnen (bestätigt über
  `schema migrate --help`: kein Flag, das einen zuvor berechneten,
  signierten Plan zur Ausführung wieder einliest — `--plan-artefact`
  ist ein reiner Ausgabe-Pfad). `decide()` prüft ausschließlich den
  `--plan-only`-Report; `--allow-destructive` ist danach ein pauschales Flag
  für den *gesamten* nachfolgenden `--execute`-Lauf. Entstünde zwischen
  beiden `docker run`-Aufrufen (durch einen konkurrierenden Schreiber auf
  demselben Ziel) eine neue, echte destruktive Operation, die nicht zu den
  sechs bekannten Fremdobjekten gehört, würde sie vom Guard nie geprüft,
  liefe aber unter dem bereits gesetzten `--allow-destructive` durch — die
  Sicherheitseigenschaft ist damit „kein ungeprüfter destruktiver Blocker
  zum Zeitpunkt des Precheck", nicht „kein ungeprüfter destruktiver Blocker,
  der tatsächlich ausgeführt wird". Weder der Slice-Plan (§1, §6) noch ein
  Kommentar benennt diese Lücke; sie ist in der Praxis eng (Fenster von
  wenigen Sekunden, typischer Aufrufer ist ein einzelner Operator/CI-Lauf),
  aber strukturell vorhanden und unbenannt.
- `verifizierbar`: nein (kein Gate; ein Nachweis bräuchte einen echten
  Race zwischen zwei `make schema-rollout`-Läufen)
- `klasse`: TOCTOU-Lücke zwischen Precheck und Execute bei pauschalem Flag

### F-4 — Fehlende Unit-Test-Abdeckung des Mischfalls ist strukturell bedingt, kein realer Gap

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/schema/rolloutguard/guard_test.go` (5 Tests, keiner mit
  einem gemischten Report aus sechs bekannten Blockern **plus** einer
  nicht-blockierenden, „echten" Operation)
- `befund`: Der durch `588140df` behobene Fehler lag ausschließlich im
  Makefile-Kontrollfluss (Skip vs. immer-`--execute`), nicht in `decide()`
  selbst — `decide()` betrachtet ausschließlich `r.Blockers` und ignoriert
  jede nicht-blockierende `Operation` bereits heute, unabhängig davon, ob
  ein solcher Mischfall im Test-Report vorkommt. Ein Unit-Test mit einer
  zusätzlichen, nicht-blockierenden Operation neben den sechs bekannten
  Blockern würde exakt dasselbe Ergebnis wie
  `TestDecideAllowsDestructiveWhenAllBlockersKnown` liefern und den
  behobenen Fehler nicht abbilden — der einzige Ort, an dem sich dieser
  Fehler überhaupt zeigen konnte, ist die Reihenfolge der `docker run`-Aufrufe
  im Makefile, was `tools/harness/run-schema-rollout-guard-test.sh` Lauf 3
  bereits real abdeckt. Kein Nachbesserungsbedarf.
- `verifizierbar`: ja (`make test`, `tools/harness/run-schema-rollout-guard-test.sh`)
- `klasse`: Test-Ebenen-Zuordnung korrekt (kein Finding)

### F-5 — `foreignObject`-Schlüssel setzt implizit Single-Schema und punktfreie Pfadsegmente voraus

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/schema/rolloutguard/guard.go:53-56` (`strings.Join(op.Path, ".")`)
  und `guard.go:16-23` (`knownForeignObjects`-Schlüssel ohne Schema-Segment)
- `befund`: Der Schlüssel `(kind, objectType, path)` disambiguiert die
  aktuell sechs bekannten Objekte korrekt (Funktions-Signaturen bzw.
  einzelne View-Namen, alle implizit im einzigen Schema `cdc`). Zwei
  Annahmen sind dabei nicht dokumentiert: (a) es gibt nur ein Schema — ein
  künftiges zweites Schema mit gleichnamigem Objekt würde kollidieren, weil
  `path` kein Schema-Segment führt; (b) `strings.Join(path, ".")` verschmilzt
  z. B. `["a.b"]` und `["a","b"]` zum identischen String. Für den heutigen
  Bestand real folgenlos (kein betroffenes Objekt), aber unadressiert.
- `verifizierbar`: nein
- `klasse`: Implizite Eindeutigkeits-Annahme im Map-Schlüssel

### F-6 — AGENTS.md §3.14-Streichung lief vor dem im Plan benannten Zeitpunkt

- `kategorie`: LOW
- `quelle`: Maintainability / Plan-Treue
- `pfad`: Slice-Plan §1 Abgrenzung, dritter Punkt ("§3.14 wird erst mit der
  Closure dieses Slice angepasst oder gestrichen, nicht vorher") vs.
  `AGENTS.md` (§3.14 bereits in Commit `b41cc381` vollständig entfernt)
- `befund`: Der Plan benennt „Closure dieses Slice" (Review + Verifier +
  Closure-Notiz + `git mv` nach `done/`) als Zeitpunkt für die
  §3.14-Streichung; real gestrichen wurde sie bereits im ersten
  Implementierungs-Commit, vor jedem Reviewer-/Verifier-Zug. Die DoD-Zeile
  selbst ist bereits mit Begründung auf `[x]` gesetzt. Funktional folgenlos,
  solange dieser Slice tatsächlich liefert (wofür die Belege real vorliegen),
  aber eine Abweichung von der eigenen, expliziten Plan-Formulierung des
  Zeitpunkts.
- `verifizierbar`: nein
- `klasse`: Zeitpunkt-Abweichung Plan vs. Ausführung

## Negativbefunde

- geprüft, ohne Befund: `tools/schema/rolloutguard/main.go`, `report.go` —
  §3.7-konform, keine Chronik-Sprache, Semantik-Wechsel `skip`→`allowDestructive`
  korrekt durchgezogen.
- geprüft, ohne Befund: `Makefile`-Kontrollfluss selbst (Exit-Code-Erfassung
  von `--plan-only`, bedingtes Setzen von `allow_destructive`, `--execute`
  läuft in jedem Zweig) — konsistent mit §1/§3 der Plan-Korrektur, kein
  AGENTS.md §3.9-Verstoß (keine Pipe/kein Wrapper zwischen `make` und
  Exit-Code).
- geprüft, ohne Befund: `.a-check.yml`-Erweiterung `tools/schema/**` unter
  `tooling` — Begründung im Kommentar ist §3.7-konform (indikativ, kein
  Vorher/Nachher), Layer-Grenze (keine Importe aus `domain`/`ports`/`app`/
  `adapters`/`contract`) real im Code bestätigt (`report.go`/`guard.go`/
  `main.go` importieren nur `encoding/json`, `fmt`, `os`, `strings`).
- geprüft, ohne Befund: `.gitignore`-Ergänzung `tools/schema/rollout-precheck.yaml`
  — Kommentar §3.7-konform, Trennung von committetem Pflicht-Report
  (`plan.yaml`) sachlich korrekt begründet.
- geprüft, ohne Befund: `harness/README.md` (§Guides-Vorgabe aus
  AGENTS.md §3.13) — die beiden betroffenen Sensors-Zeilen
  (`make schema-rollout`, `make example-demo-up`/`-down`) sind im
  Netto-Diff (`588140df`) korrekt nachgezogen; dass Commit `b41cc381`
  allein diese Zeile noch nicht trug, ist durch die Fixrunde im selben
  Slice vor jedem Review-Zug geheilt.
- geprüft, ohne Befund: Zwei-Commit-Historie (`b41cc381` + `588140df`) —
  entspricht dem Git-Safety-Protocol (neue Commits statt Amend) und ist
  eine sachgerechte Fixrunde innerhalb eines noch offenen, nicht amendierten
  Slice, vor jedem Reviewer-/Verifier-Zug; kein Normverstoß.
- geprüft, ohne Befund: `tools/harness/run-schema-rollout-guard-test.sh` —
  vier reale Läufe, Cleanup via `trap`, eigenständiges Docker-Netz, deckt
  exakt die im Slice-Plan §2 benannten Szenarien inklusive des in §1
  beschriebenen, real gefundenen Regressionsfalls (Lauf 3).
- geprüft, ohne Befund: `decide()`-Kernlogik selbst (Allowlist-Vergleich
  über `kind`/`objectType`/`path`, nicht über `id`-Hashes) — für den
  aktuellen Objektbestand korrekt und „fail closed" bei jeder Anomalie
  (leerer Report, unbekannte Reason-Klasse, fehlende Operation-ID).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Slice-Chronik/Vorher-Nachher-Sprache in
Produktionscode-Kommentar · Träger nicht nachgezogen (§3.13) · TOCTOU-Lücke
zwischen Precheck und Execute bei pauschalem Flag · Implizite
Eindeutigkeits-Annahme im Map-Schlüssel · Zeitpunkt-Abweichung Plan vs.
Ausführung

## Verdikt

**Merge-blockierend:** ja — zwei HIGH-Findings (F-1, F-2) lösen eine
Rückkante zum Implementer aus (Kommentar-Bereinigung in `Makefile` +
`guard.go`; Meldung/Behandlung der stale
[`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md)-Referenz
— Zitat-Korrektur prüfen, ob
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
greift, sonst Meldung ins Beobachtungs-Register statt stiller Änderung).
F-3 (MEDIUM) sollte vor Closure mindestens als benanntes Risiko in
Slice-Plan §6 nachgetragen werden, auch wenn keine Code-Änderung zwingend
folgt.

**Übergabe:** Findings gehen an den Implementer (Fixrunde nötig — die
DoD-Zeile „Review durchgeführt" bleibt deshalb bewusst `[ ]`, kein
Nachzug ohne Fixrunde). Die Finding-Klassen gehen in die Slice-Closure §7
und von dort in den Beobachtungs-Register-Zähler. Dieser Report ersetzt
keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
