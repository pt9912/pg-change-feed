# Review-Report (Fixrunde): slice-081 — 2026-09-15

**Review-Art:** Code-Review gegen Plan + Entscheidungen (Modul 10), **Fixrunde** —
nur die Korrekturen am Verdikt aus
`docs/reviews/review-slice-081.md`; **kein** Vollauf, **nicht** gegen die DoD.

**Gegenstand:** die Korrektur-Commits auf `main`:
`c5b88a9` · `8167510` · `ba5f31c` (erste Runde) und
`80ab5d5` · `a4b4e19` · `de7ba0b` (zweite Runde), Stand `de7ba0b`.
Nicht Gegenstand: der ursprüngliche Slice-Diff (dort geprüft) und die
ADR-Commits `fad0da5`/`938ff90`/`4951c73`/`69d12f4` (eigene Vorgänge).

**Skill:** `.harness/skills/reviewer.md` @ `slice-081`-Stand · **Modell:** deepseek-v4.1-flash · **Datum:** 2026-09-15

**Eingangs-Kontext:** `slice-081` §1–§3 · `ADR-0071` Punkt 5 · `ADR-0078` ·
`AGENTS.md` §3.6/§3.7 · Beobachtungs-Register
`BEO-PGC/plan-vorlagen-defekt` (verkörpert, `seit slice-010`),
`BEO-PGC/vorlagenrest-in-closure-notiz`,
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` · der eigene Vorlauf-Report.

---

## Stand der Findings aus `review-slice-081`

| Finding | Kategorie | Stand in dieser Runde |
|---|---|---|
| F-1 Zählbasis-Kommentar `610` statt gemessen `472` (`tools/harness/db-coverage.sh`) | HIGH | **behoben** (`8167510`) — eigene Zählung stimmt zeichengenau (Datei-Split `132·93·87·58·52·27·17·6 = 472`); Kommentarblock sonst clean (die zweite Zahl „`replication/receive` auf 0 von 155" ist die Aussage über den *fehlerhaften* Merge, keine aktuelle Messung; 155/23 am frischen gemergten Profil bestätigt) |
| F-2 zugesagte „minimale `Rows`" nicht verdrahtet | MEDIUM | **behoben in beiden Hälften** — Plan richtig berichtigt (`ba5f31c`, `de7ba0b`), Code nachgezogen (`a4b4e19`) |
| F-3 Folge-Vorgänge ohne auflösende Adresse | MEDIUM | **behoben** (`c5b88a9`, `80ab5d5`) — `open/slice-084-postgresack-naht.md` · `open/slice-085-receive-naht.md`, in §3 als Adressen genannt; Symbolzahlen **nachgemessen**: `postgresack` 6, `receive` 18 |
| F-4 Schwankungs-Bandbreite überzeichnet | LOW | **behoben** (`8167510`) — `2/1854 = 0,108 pp` → „±0,11 pp" ✓; `0,70 · 1854 = 1297,8`, `1308 − 1297,8 = 10,2` → „≈11 Statements" ✓ |
| F-5 Beleg-Kommandos lösen nicht auf | LOW | **behoben** (`80ab5d5`) — jetzt `4c7ea8a~1..8e9fe4f` |
| F-6 stale Zahlen in lebenden Trägern außerhalb des Diffs | INFO | verabredet ausgelagert (Herkunfts-Regel aus `ADR-0078`), kein Slice-081-Gegenstand |

`F-5` hatte in seiner ersten Fassung (`252962b`) **nicht** getragen: die Kennung
war eine Rebase-Waise — `git merge-base --is-ancestor 252962b HEAD` leer, in
keinem Ref (`git branch -a --contains` leer), `git rev-list --all` kennt sie
nicht; sie lebte nur im Reflog (`HEAD@{25}: rebase (pick)`). Das war eine
Verschlechterung gegenüber `main..HEAD`, weil die Kennung in einem frischen Klon
gar nicht existiert. Die jetzige Kennung ist geprüft: Vorfahr von `HEAD`, in
`rev-list --all`, und beide Bereiche liefern denselben Dateisatz.

---

## Findings dieser Runde

### F-1 — §3 der Slice-081 zitiert die **alte** Wortlaut-Fassung von §1/§2 als deren geltende

- `kategorie`: LOW
- `quelle`: Maintainability (Plan-interner Querverweis driftet gegen die
  korrigierte Sektion)
- `pfad`: `docs/plan/planning/in-progress/slice-081-executor-naht.md:222-224`
  gegen `:44-50` (§1) und `:146-152` (§2)
- `befund`: §3 sagt „die in §1/§2 genannte Schnittstelle (`Query`/`Exec` plus
  minimales `Rows`)". §1 nennt seit `ba5f31c` `Query`/`Exec` plus `QueryRow`;
  §2 Liefer-Punkt 1 nennt `Query`/`Exec`/`QueryRow` — „minimales `Rows`" steht
  dort nicht mehr. Der Querverweis schreibt §1/§2 damit eine Fassung zu, die
  nur noch die **berichtigte** ist. Die Argumentation der Zeile (eine
  Query/Exec-Schnittstelle drückt `pglogrepl`-Aufrufe nicht aus) bleibt richtig;
  falsch ist allein die zitierte Form. Dieselbe Familie eine Sektion höher:
  `:53` sagt „die deklarierte `sqlexec.Rows` fordert **4**" im Präsens, während
  `a4b4e19` die Deklaration entfernt hat.
- `verifizierbar`: ja — Wortlaut-Abgleich der beiden Sektionen (kein Gate)

---

## Negativbefunde

- geprüft, ohne Befund: **F-2, Code-Hälfte (`a4b4e19`)** — `sqlexec.Rows` samt
  Kommentar ist entfernt, `seam.go` trägt keine tote Deklaration mehr; `var _
  DB = (*pgxpool.Pool)(nil)` steht unverändert; kein Verweis auf
  `sqlexec.Rows` bleibt im Code (`git grep`). Die Test-Zusicherung hält jetzt
  gegen den entschiedenen Zeilentyp (`var _ pgx.Rows = (*fakeRows)(nil)`), und
  der Fake kompiliert gegen dieselbe Fläche, die `fakeExecutor.Query` liefert.
- geprüft, ohne Befund: **F-2, Plan-Hälfte (`ba5f31c`, `de7ba0b`)** — §1 führt
  die Berichtigung mit benannter, verworfener Alternative und benannter Grenze
  („der Fake muss deshalb 10 Methoden erfüllen statt vier"); §2 Liefer-Punkt 1
  sagt „Die Zeilen-Schnittstelle bleibt `pgx.Rows`"; die §3-Träger-Tabelle nennt
  jetzt „`Query`/`Exec`/`QueryRow`; Zeilen als **`pgx.Rows`**" und „Fakes
  (`pgx.Rows`/`pgx.Row`/`Executor`)".
- geprüft, ohne Befund: **§2-Vorlagen-Rest** — `- [ ] Doku-Update für
  <Schnittstelle X> …` ist in **beiden** neuen Adressen durch ein konkretes
  Kriterium ersetzt (Transfer-Nachweis in den zwei Sensor-Docs, ohne neue
  Schwellen-ADR). Kein `<Schnittstelle X>` mehr in `open/`; die verbleibenden
  `<…>` liegen in §6/§7 und sind laut `BEO-PGC/plan-vorlagen-defekt` in `open/`
  zu Recht da.
- geprüft, ohne Befund: **Adapter-Richtung** — `slice-085` sagt jetzt an allen
  vier Stellen `driving` (§Bezug, §Berührte-Spec-Stellen, §Ziel, §Träger-Tabelle).
- geprüft, ohne Befund: **die compilerseitige Zusicherung prüft wirklich** —
  eigene Mutation (M1): `fakeRows.TypeMap` entfernt ⇒ `go test` **Exit 1**,
  `*fakeRows does not implement pgx.Rows (missing method TypeMap)` an zwei
  Stellen (`fakeExecutor.Query`-Rückgabe und die Zusicherung). Eigene Mutation
  (M2): `rows.Err()`-Pfad in `ReadChanges` entfernt ⇒ `go test` **Exit 1**
  (`TestReadChangesClassifiesIterationFailure`) — der Fake prüft nach der
  Änderung unverändert.
- geprüft, ohne Befund: **`make gates`** — eigener Lauf auf `de7ba0b`, Exit **0**
  (d-check 680 Dateien / 0 Befunde · commit-traceability OK · a-check 0 ·
  coverage-gate grün). Alle Mutationen zurückgenommen, Arbeitsbaum sauber.
- geprüft, ohne Befund: **Nebenwirkungen der Code-Änderung** — der
  Produktionspfad ist unberührt (nur ein Test und eine unbenutzte Deklaration);
  `make test-store`/`make test-replication` bleiben außerhalb dieser Änderung
  grün (im Vorlauf gemessen, `73.38 % (477 von 650)`).

## Eigene Messungen (Exit-Codes ungepiped)

| Lauf / Prüfung | Ergebnis |
|---|---|
| `make gates` auf `de7ba0b` | **Exit 0** (d-check 680/0, commit-traceability OK, a-check 0) |
| `git merge-base --is-ancestor 4c7ea8a HEAD` · `git rev-list --all \| grep -c 4c7ea8a` | Vorfahr · 1 (erreichbar) |
| `diff` der Bereiche `4c7ea8a~1..8e9fe4f` und `252962b~1..8e9fe4f` (`-- internal/`) | identisch |
| `252962b`: `--contains` / `rev-list --all` / Reflog | leer / nicht enthalten / nur `rebase (pick)` → Waise |
| M1 `fakeRows.TypeMap` entfernt | **Exit 1**, `missing method TypeMap` (2 Stellen) |
| M2 `rows.Err()`-Pfad entfernt | **Exit 1**, `TestReadChangesClassifiesIterationFailure` |
| Symbolzahlen `postgresack` / `receive` (`pgconn.`/`pglogrepl.`, distinct) | 6 / 18 — wie in `slice-084`/`slice-085` angegeben |
| Stale-Zahl-Scan (`610`, `75,25`, `1171`, `2467`, `69,70`) über lebende Träger | nur noch in Übergangs-/Zitat-Form (`610 → 472` in §3 und `ADR-0078`) und in Lauf-Beleg-Reports; keine geltende Zahl driftet |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Plan-interner Querverweis driftet gegen die
korrigierte Sektion" (F-1, LOW, mit dem Präsens-Rest in `:53` als zweitem Pfad
derselben Klasse).

## Verdikt

**Merge-blockierend: nein.** Alle vier benannten Punkte aus
`review-slice-081` tragen; die verbleibende LOW liegt im **Plan-Text** der
Planner-Rolle (Sektion §3/§1 desselben Dokuments) und braucht **keinen**
Rückgabe-Pfeil an den Implementer — sie geht mit dem Planner-Zug der Closure.

**DoD-Häkchen:** das Kästchen „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" ist mit dem Anlegen dieses Reports auf `[x]` nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde: 0 HIGH, kein
Implementer-Rückgabe-Pfeil). Die Verifikation bleibt davon unberührt
(Modul 11, eigener Kontext).
