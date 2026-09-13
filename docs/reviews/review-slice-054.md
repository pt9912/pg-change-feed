# Review-Report: slice-054 — 2026-09-13

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `1cb1f06` — NATS-Boundary-Beleg (Change ohne verbundenen
Consumer) in `tools/harness/run-integration-tests.sh`. Elter `ec8b3d3`
(reiner `next→in-progress`-Move, keine Inhaltsänderung).

**Skill:** `.harness/skills/reviewer.md` @ `f07ba3d` (Stand 2026-09-13, vier
repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-054-nats-boundary-beleg.md` (§1
  Ziel/Abgrenzung, §2 DoD, §6 Risiken)
- `ADR-0055` (`docs/plan/adr/0055-nats-change-notification-wecksignal.md`) —
  Punkt 4 (Fehlerklasse/kritischer Pfad, best-effort nach ACK)
- `ADR-0056` (`docs/plan/adr/0056-nats-tabellen-granulares-subjekt.md`) —
  Vier-Token-Subjekt `cdc.changes.<source_id>.<schema>.<table>`, hier bereits
  in `NATS_SUBJECT` (Bestand aus `slice-053`) korrekt verwendet
- `LH-FA-SST-007` (Boundary-/Negative-Akzeptanzkriterium)
- `compose.yaml` (NATS-Service, Digest-Pin, `-m 8222` Monitor-Port)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7), §6 Workflow
- `harness/conventions.md` (MR-000 ID-Schema)
- Vorherige Findings am gleichen Modul: `docs/reviews/review-slice-053.md`

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings. Eine INFO-Beobachtung unten.

### F-1 — Fixer `sleep 2` als einziges Zeitfenster für den Negativ-Beleg

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:1310-1314`
- `befund`: Der Testablauf wartet nach der auslösenden Change eine feste
  Sekundenzahl, bevor er `cdc.changes`, den Feed-Container-Status und
  `/subsz` erneut abfragt. Das ist für einen Negativ-Beleg („nichts soll
  passieren") eine bewusste, im Kommentar begründete Design-Entscheidung
  (kein Poll auf ein erwartetes Ereignis möglich) und keine Lücke im
  Boundary-Beleg selbst — ein zu kurzes Fenster könnte im ungünstigsten Fall
  einen späten, fälschlich doch propagierten Fehler verpassen. Kein
  aktueller Mangel, da der reale `make test-integration`-Lauf (siehe
  Negativbefunde) den Container in genau dieser Konstellation grün zeigt.
- `verifizierbar`: ja — `make test-integration`, bereits real ausgeführt.
- `klasse`: „fixe Wartezeit als einziges Fenster für Negativ-Beleg"

## Negativbefunde

- geprüft, ohne Befund: **Mechanismus-Realität des `/subsz`-Belegs.** Das
  JSON-Format des NATS-HTTP-Monitors (`"subject": "<wert>"`, Leerzeichen
  nach dem Doppelpunkt) wurde gegen einen frisch gestarteten Container mit
  demselben Digest-Pin (`nats:2-alpine@sha256:065e8355c2…`) real reproduziert
  — sowohl der Leerlauf-Zustand (nur `$SYS.*`-Subscriptions, keine
  Überschneidung mit `cdc.changes.*`) als auch eine real geöffnete/wieder
  geschlossene Test-Subscription auf exakt
  `cdc.changes.src-mvp.public.feed_mvp_full`: Das Subjekt erscheint während
  offener Verbindung im `/subsz`-Output und verschwindet unmittelbar nach
  `docker rm -f` des Subscriber-Containers, ohne messbare Verzögerung. Der
  `grep -qF "\"subject\": \"$NATS_SUBJECT\""`-Mechanismus ist damit
  mechanismus-basiert und keine Annahme — er prüft real gegen den laufenden
  NATS-Server, nicht gegen ein angenommenes Testablauf-Verhalten
- geprüft, ohne Befund: **Reihenfolge.** Vor-Check (`nats_subs_before_boundary`)
  → `INSERT` (id=231) → `sleep 2` → `cdc.changes`-Lesebeleg → `feed_running`
  → Nach-Check (`nats_subs_after_boundary`). Deckt beide im Auftrag
  genannten Prüfpunkte in der richtigen Kausalreihenfolge ab; die
  Positionierung von `feed_running` zwischen Lesebeleg und Nach-Check statt
  ganz am Ende ist eine reine Reihenfolge-Variante ohne Auswirkung auf die
  Aussagekraft der drei Belege
- geprüft, ohne Befund: **`feed_running`-Check** — vorhanden, gleiches
  Muster wie beim bestehenden Happy-Path-Beleg (`docker inspect --format
  '{{.State.Running}}'`), mit auf `ADR-0055` Punkt 4 verweisender
  Fehlermeldung
- geprüft, ohne Befund: **ID-Kollision.** `id=231` ist repo-weit einmalig
  (Grep gegen `feed_mvp_full`/`$RETENTION_TABLE`/`$LIFECYCLE_TABLE`-Inserts
  bestätigt keine Wiederverwendung); `NATS_SUBJECT` wird unverändert aus dem
  bestehenden Happy-Path-Block übernommen (kein Redefinitions-Risiko)
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7).** Der neue
  17-zeilige Kommentar beschreibt den gegenwärtigen Mechanismus (Zusage:
  reale Server-Momentaufnahme statt Annahme; Abgrenzung: `$SYS.*` vs.
  `cdc.changes.*`) im Indikativ, ohne Slice-/Wellen-Nummer und ohne
  Bezugnahme auf einen früheren/entfernten Code-Zustand — keine der drei
  HIGH-Kommentar-Klassen (verworfene Alternative im Konjunktiv, abwesenter
  Text, abgebrochener Satz) trifft zu
- geprüft, ohne Befund: **Out-of-Scope-Einhaltung (§1 des Slice-Plans).**
  `git show --stat` bestätigt genau eine geänderte Datei
  (`tools/harness/run-integration-tests.sh`, nur Insertionen); keine
  Änderung an `ChangeNotificationPort`/`internal/adapters/driven/natsnotify/`,
  keine neue Compose-Verdrahtung, kein Reconnect-Verhalten
- geprüft, ohne Befund: **Commit-Trailer/Traceability.** `1cb1f06` trägt
  keine `Co-Authored-By:`/`Claude-Session:`-Zeilen; Betreff referenziert
  `LH-FA-SST-007` und `ADR-0055`, keine `SPEC-*`/`ARC-*`-ID im Betreff
- geprüft, ohne Befund: **Docker-only (`AGENTS.md` §3.1).** Ausschließlich
  `docker exec`/`docker run`/`docker inspect` gegen bereits laufende,
  digest-gepinnte Container; keine lokale Toolchain-Installation
- geprüft, ohne Befund: **reale Sensor-Läufe** — `make gates` (baseline-verify
  v6.5.0 OK, Coverage-Gate 40,90 % ≥ 35 %, `d-check` 427 Dateien/0 Befunde
  zweimal — voller Lauf und Commit-Range —, Commit-Traceability OK, `a-check`
  0 Befunde plus dem vorbestehenden, unveränderten Abdeckungs-Hinweis zu
  `tools/harness/natssub/main.go`) und `make test-integration` (voller
  Compose-Lauf inkl. des neuen Boundary-Abschnitts: „NATS-Boundary-Beleg
  (LH-FA-SST-007) — Change (id=231, feed_mvp_full) entstand real ohne einen
  auf … abonnierten Client …, blieb vollständig über cdc.changes lesbar, und
  der Feed-Container lief unverändert weiter"), beide real in diesem Review
  ausgeführt, Exit 0

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** fixe Wartezeit als einziges Fenster für
Negativ-Beleg

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW. Die einzige INFO ist
eine bereits im Code begründete, bewusste Design-Entscheidung ohne
erwartete Aktion.

**Übergabe:** Kein Fixrunden-Pfad zum Implementer nötig. Nach
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde wird die
DoD-Zeile „Review durchgeführt" im Slice-Plan in diesem Commit mitgezogen.
Die Finding-Klasse geht zusätzlich in die Slice-Closure §7 und von dort in
den Beobachtungs-Register-Zähler. Dieser Report selbst ist ein
**Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und wird über Läufe hinweg nicht wieder gelesen. Der Report
ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier
separat (Modul 11).
