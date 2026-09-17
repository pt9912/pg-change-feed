# Review-Report: slice-104 — 2026-09-17

**Review-Art:** Code — geprüft gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), nicht gegen DoD (das ist Verifier-Aufgabe).

**Gegenstand:** Diff `7d66196..HEAD` (Commits `338cfe0`, `59d5b53`),
Slice `docs/plan/planning/in-progress/slice-104-proto-generate-mount-los.md`.

**Skill:** `.harness/skills/reviewer.md` @ `7e598bb` (Repo-Stand des Laufs)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-104-proto-generate-mount-los.md` (§1
  Abgrenzung, §2 LP1–LP3, §6 Risiken)
- `ADR-0060` (grpc-streaming-mechanismus) — Bezug des Slice
- `ADR-0084` (sync-gate-fuer-generierte-artefakte) — Kontext-Bezug
- `AGENTS.md` §3.9 (Pipe-Maskierung), §3.1 (Docker-only), §3.7 (Kommentar-Klassen)
- `docs/plan/planning/observations/BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/`
  (observation.md, state.md, beide Evidence-Dateien)
- Diff selbst (`git diff 7d66196..HEAD`)

---

## Findings

### F-1 — ADR-0060-Zitat mehrdeutig platziert (Dockerfile-Kommentar)

- `kategorie`: LOW
- `quelle`: Maintainability (§3.7 Kommentar-Klassen, angrenzend)
- `pfad`: `Dockerfile:40-41`
- `befund`: Der Kommentar über der neuen Stufe `proto-export` schreibt
  „vormals `docker run -v` in den Bind-Mount des Arbeitsbaums, siehe
  `ADR-0060`" — die Zitat-Platzierung direkt hinter der Bind-Mount-Aussage
  legt nahe, `ADR-0060` beschreibe oder verlange den Bind-Mount-Mechanismus.
  Tatsächlich enthält `ADR-0060` keine Erwähnung von Bind-Mount/`docker run
  -v` (`grep -in "bind.mount\|docker run.*-v\|mount"
  docs/plan/adr/0060-grpc-streaming-mechanismus.md` — leer); der
  Bind-Mount war laut dem Slice-Kopf selbst eine Umsetzungsentscheidung von
  `slice-069`, keine `§Entscheidung` von `ADR-0060`. Der Slice-Plan trifft
  diese Unterscheidung im eigenen `Bezug`-Feld sauber, der Code-Kommentar
  übernimmt die Präzision nicht.
- `verifizierbar`: ja — `grep` gegen `ADR-0060` bestätigt die Abwesenheit
  jeder Bind-Mount-Erwähnung.
- `klasse`: „ADR-Zitat-Platzierung mehrdeutig"

---

## Negativbefunde

- **Pipe-Maskierung (§3.9), Punkt 1+2 des Auftrags** — geprüft, kein Befund.
  `grep -n -A10 "^proto-generate:" Makefile` zeigt eine reine Umleitung
  (`>`), kein Pipe-Zeichen; kein `-`-Präfix vor den relevanten Zeilen.
  Eigener Build (`docker build --target proto-export -t test-proto-export
  .`) erfolgreich; `make proto-generate` real ausgeführt — extrahierte
  Dateien gehören `db:db` (aufrufender Nutzer, kein `root`), Inhalt
  byte-identisch zum vorherigen Stand (`diff` gegen Kopie vor dem Lauf:
  keine Abweichung). Fehlschlag real simuliert (`docker run --rm
  --network none <nicht-existentes-image>` sowohl direkt als auch über ein
  isoliertes `make`-Testrezept mit identischem Rezeptmuster): `docker run`
  liefert Exit 125, `make` bricht an dieser Zeile ab
  („`make: *** [test.mk:6: t] Fehler 125`"), die nachfolgende
  `tar -xf`-Zeile läuft nachweislich nicht — die leere Zwischendatei
  bleibt liegen, exakt wie im `.gitignore`-Kommentar dokumentiert.
- **`make generated-sync`** — geprüft, kein Befund. Selbst ausgeführt,
  Exit 0, „byte-gleich" bestätigt, unabhängig vom geänderten
  `proto-generate`-Mechanismus.
- **`.dockerignore`-Kommentar zu `!proto/` (Abgrenzung zu `ADR-0085`)** —
  geprüft, kein Befund über den Klassifikations-Punkt hinaus (siehe unten).
  Formulierung ist präzise: „gelesen **und** gebaut, ein ganzes Verzeichnis
  statt einer Datei" trifft empirisch zu (`COPY proto/ proto/` gefolgt von
  `RUN protoc`), reale Gegenprobe ohne `!proto/`-Zeile bestätigt einen
  echten Build-Fehler (`COPY failed: ... "/proto": not found`).
- **`gen/cdc/stream/v1/**` in `.gitignore`** — geprüft, kein Befund;
  `git check-ignore` bestätigt, die Dateien bleiben ungeignored.
- **`make image`-Digest** — geprüft, kein Befund. `runtime` hängt in der
  `FROM`-Kette (`Dockerfile:19,35,50,76,102,108`) nachweislich nicht von
  `proto`/`proto-export` ab; realer `make image`-Lauf liefert denselben
  Digest wie `harness/image-hash.txt` (`sha256:50348a0a…`), kein
  Arbeitsbaum-Diff auf der Datei.
- **Fünf Träger-Fundstellen (LP3)** — geprüft, kein Befund; eigener,
  breiterer Suchlauf (`grep -rln "proto-generate" …`) findet keine
  sechste lebende Fundstelle. Ein Treffer in `docs/plan/adr/0085-…md`
  („Bind-Mount des Arbeitsbaums … Muster `make proto-generate`") und in
  `docs/plan/adr/0076-…md` sind unberührt korrekt: Ersterer steht in
  einer historischen `§Verglichene Alternativen`-Zeile einer `Accepted`-ADR
  (immutabel, `AGENTS.md` §3.5 — kein Nachzugs-Kandidat), Letztere nennen
  `make proto-generate` nur namentlich ohne Mechanismus-Detail.
- **Out-of-Scope-Disziplin (`generated-sync.sh`)** — geprüft, kein Befund.
  `git diff 7d66196..HEAD -- tools/harness/generated-sync.sh` ändert
  ausschließlich Kommentarzeilen; die Skript-Logik (`RUN_USER=`-Zeile
  selbst, Docker-Aufrufe, Vergleichslogik) ist byte-identisch.
- **Commit-Traceability** — geprüft, kein Befund. Beide Commits
  (`338cfe0`, `59d5b53`) nennen `ADR-0060` im Betreff; `make
  commit-traceability` (Teil von `make gates`) läuft grün.
- **§3.13-Gegencheck (bewegte Eigenschaften in fremden Trägern)** — geprüft,
  kein Befund über LP3 hinaus. Der einzige weitere Kandidat
  (`docs/plan/adr/0085-…md:98`, Kontext-Zähler „170 Kontext-Einträge")
  bleibt unverändert korrekt: Er ist eine historische, gemessene Momentaufnahme
  aus `slice-093` in einer `Accepted`-ADR, kein live nachzuziehender
  Zähler — die Zeile trägt ihren Lauf bereits als Beleg
  (`Lauf slice-093`) und bleibt damit außerhalb des Nachzugs.
- **`make gates`** — geprüft, kein Befund. Eigener Lauf, Exit 0
  (`d-check` 849 Dateien/0 Befunde, `commit-traceability` OK,
  `generated-sync` OK, `coverage-gate` 83.40 % ≥ 80 %, `a-check` 0 Befunde).

---

## Beobachtungsklassen-Frage (Auftragspunkt 4, ausführlich)

Der Implementer vermutet, der `!proto/`-Fund sei die **dritte Instanz**
von `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (bisher 2×,
`slice-049`, `slice-093`). Eigene Prüfung von `observation.md`/`state.md`:

Die beiden bestehenden Belege beschreiben spezifisch **ein Skript unter
`tools/`**, das eine neue Docker-Stufe braucht und an `.dockerignore`
scheitert — plus (in `slice-049`) eine zusätzliche, davon unabhängige
Alpine/`bash`-Facette. `state.md`s eigene Handlungsanweisung benennt den
Trigger für künftige Wiederholung explizit als „eine weitere Docker-Stage,
die **ein Skript unter `tools/`** braucht". Der Titel selbst trägt
„…-tooling" als Teil seiner **Pfad-Kennung** — und die Baseline-Regel
(Modul 6 §Das Beobachtungs-Register) ist hier ausdrücklich strikt: „Ein
Prosa-Name taugt nicht — er darf umformuliert werden, ein Pfad nicht."

Der aktuelle Fund unterscheidet sich in zwei Punkten von der bisherigen
Klasse: (a) Gegenstand ist die `.proto`-**Quelle** (ein Verzeichnis mit
Quelldateien für eine neue Erzeugungsstufe), kein Skript unter `tools/`;
(b) die zweite, bisher an jedem Beleg mitlaufende Facette (Alpine-Basis
ohne `bash`) tritt hier gar nicht auf. Der gemeinsame **Mechanismus** ist
zwar real derselbe (`.dockerignore`s Allow-Listen-Default-Deny bricht bei
jeder neuen Docker-Stufe, die einen bisher nicht gelisteten Pfad braucht —
eigene Gegenprobe: ohne `!proto/` scheitert `COPY proto/ proto/` real mit
„not found"), aber die **benannte Klasse**, wie sie in `observation.md`
und `state.md` tatsächlich geschrieben steht, ist enger gefasst als der
allgemeine Mechanismus: Sie ist auf „Skript unter `tools/`" verengt, nicht
auf „beliebiger, bisher ungelisteter Pfad, den eine neue Stufe braucht".

**Eigenes Urteil:** Die Vermutung des Implementers ist **nicht** einfach
zu übernehmen. Es handelt sich strukturell um denselben Fehler-Mechanismus,
aber **nicht** um dieselbe, im Register benannte Klasse — die dort
dokumentierte Identität ist über den Pfad `.../blockiert-tooling` und den
`tools/`-Skript-Bezug beider Evidence-Dateien enger gezogen, als der neue
Fund trägt. Ihn unter diesem Pfad als 3. Beleg einzutragen würde den
Zähler auf einen Titel heben, dessen zwei tragende Belege eine andere,
engere Sache belegen, und die dann verkörperte Regel („bei einem Skript
unter `tools/` vorab prüfen") träfe den neuen Fall (Quellverzeichnis für
eine Erzeugnisstufe) nicht mehr präzise ab. Da diese Einschätzung kein
isoliertes LOW/INFO ist, sondern eine Zähler-Frage, die die
Register-Integrität berührt (Modul 6 „Mensch urteilt, Maschine prüft
Deckung" — das Urteil *ist das dieselbe Beobachtung?* fällt beim
Schreiben), wird sie als eigenes Finding geführt statt in den
Negativbefunden zu verschwinden:

### F-2 — Beobachtungs-Register: falsche Klassenzuordnung droht

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Das
  Beobachtungs-Register (Maintainability)
- `pfad`: `docs/plan/planning/observations/BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/`
  (noch keine Datei in diesem Diff geändert — betrifft die anstehende
  Closure)
- `befund`: Der `!proto/`-Fund (neues Docker-Stufen-Erfordernis gegen
  `.dockerignore`) ist derselbe **Mechanismus**, aber nicht dieselbe
  **benannte Klasse** wie die zwei bestehenden, auf „Skript unter
  `tools/`" verengten Belege dieses Registereintrags; eine Verbuchung als
  3. Instanz würde den Zähler unter einem zu engen Titel erhöhen.
- `verifizierbar`: nein — Urteilsfrage (Modul 6: „Mensch urteilt"), kein
  Gate-Lauf entscheidet sie.
- `klasse`: „Beobachtungs-Klassifikation zu eng/zu weit gefasst"

Dieser Fund ist **nicht** merge-blockierend für den Diff selbst (die
Register-Eintragung passiert erst bei Closure, die Datei ist im
vorliegenden Diff unverändert) — er ist eine Vorab-Warnung an die
Closure: entweder ein **neuer** `BEO-PGC/`-Eintrag für den generalisierten
Mechanismus („Docker-Stufe braucht bisher ungelisteten `.dockerignore`-Pfad")
oder eine bewusste, im Eintrag selbst dokumentierte Erweiterung des
bestehenden Titels/Textes — beides ist besser als eine stille Verbuchung
unter einem zu engen Namen.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** ADR-Zitat-Platzierung mehrdeutig ·
Beobachtungs-Klassifikation zu eng/zu weit gefasst

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, ein MEDIUM (Register-Urteilsfrage
ohne Auswirkung auf den vorliegenden Diff selbst, relevant erst bei
Closure) und ein LOW (Kommentar-Präzision). Beide Findings gehen als
Hinweis an den Implementer bzw. an die Closure-Entscheidung, keine
Fixrunde am Diff selbst nötig.

Die DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt
vor" wird im selben Commit, der diesen Report anlegt, im Slice-Plan
nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde — keine
Reviewer→Implementer-Rückkante in diesem Lauf).

Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität (inkl. der
tatsächlichen Risiko-Ausgänge in §6/§7 des Slice-Plans) prüft der
Verifier separat (Modul 11).
