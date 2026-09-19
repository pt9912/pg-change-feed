# Verifikationsbericht: slice-sdk-csharp-pack-werkzeug — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-csharp-pack-werkzeug.md` §2)
und die `ADR-0106`-Konformität, in frischem Kontext. **Nicht** gegen den
Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-csharp-pack-werkzeug.md`](review-slice-sdk-csharp-pack-werkzeug.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** Diff-Range `89f62131..HEAD` (fünf Commits): `4a7b0d62`
(open→next, reiner Move), `bb23e3d8` (Verantwortlich gesetzt),
`8438ed95` (next→in-progress, reiner Move), `dcd05acc` (Implementer-Zug:
`harness/mk/sdk.mk`, `tools/harness/sdk-pack-csharp.sh`,
`sdks/csharp/Dockerfile`-Erweiterung, `sdks/csharp/.gitignore`,
`spec/pflichtenheft.md`-Träger-Nachzug, `harness/README.md`-Werkzeugzeile),
`69d2dc00` (Review-Report + DoD-Checkbox-Nachzug ohne Fixrunde).

**Frischer Kontext, eigene Läufe statt Behauptungs-Übernahme:** Slice-Plan
(vollständig), Review-Report (vollständig, 0 HIGH/MEDIUM/LOW, 1 INFO),
`ADR-0106` (Accepted, Volltext inkl. Festlegung 4/5, Konsequenzen
Folgepflicht), `spec/pflichtenheft.md` §1/§6/§7 (Volltext der berührten
Stellen), `harness/README.md` §Werkzeuge (neue Zeile). Eigene Bau-Läufe,
eigene Mutations-Probe (andere Datei/Assertion als der Reviewer), eigener
`make gates`-Lauf, eigene `grep`-Kollisionsprüfungen, eigene
`doc-commits`/`doc-immutable`-Läufe über den exakten Slice-Bereich.

---

## 1. DoD-Vertrag (§2) — jede Zeile einzeln geprüft

### 1.1 `make sdk-pack-csharp` existiert, Docker-only, kein Gate

- `harness/mk/sdk.mk` gelesen: Fragment trägt **kein** `GATE_CHECKS +=`
  (im Gegensatz zu `baseline.mk`, `doc-gate.mk` (2×), `generated-sync.mk`,
  `coverage.mk`, `a-check.mk` — alle sechs tragen es real, eigener `grep -n
  "GATE_CHECKS" Makefile harness/mk/*.mk a-check.mk` bestätigt). `Makefile`
  Zeile 322/323 (`gates: record-gates`, `record-gates: $(GATE_CHECKS)`)
  zeigt: `GATE_CHECKS` ist ausschließlich `baseline-verify`, `docs-check`,
  `commit-traceability`, `coverage-gate`, `generated-sync`, `a-check` —
  exakt die sechs in `harness/README.md` §Sensors genannten Gates.
- **Eigener Bau-Lauf 1 (sauberer Stand):** `rm -rf sdks/csharp/dist &&
  make sdk-pack-csharp` → `EXIT=0`. `sdks/csharp/dist/PgChangeFeed.Client.0.1.0.nupkg`
  existiert danach real, `ls -la` zeigt **29947 Bytes** — identisch mit
  Implementer- und Reviewer-Angabe, kein Drift (`AGENTS.md` §3.12 Instanz
  A). `file sdks/csharp/dist/PgChangeFeed.Client.0.1.0.nupkg` →
  „Zip archive data, at least v2.0 to extract, compression method=deflate"
  — bestätigt ZIP-Form.
- **Eigener `git status --porcelain`** nach dem Lauf: Artefakt erscheint
  **nicht** als untracked — `sdks/csharp/.gitignore`s `dist/`-Zeile trägt
  real (eigene Lektüre der Datei bestätigt beide Zeilen: `*.nupkg` unter
  Host-Bau-Schutz, `dist/` unter Export-Zielort).
- **Eigene, vom Reviewer unabhängige Mutations-Probe (SECHSTE
  Bau-Bestätigung, andere Mutation als der Reviewer):** Der Reviewer
  mutierte `PgChangeFeedClientOptionsTests.cs` (Token-Assertion). Diese
  Sitzung mutierte stattdessen
  `Http/PgChangeFeedHttpClientTableTests.cs:24`
  (`Assert.Equal("/tables/enable", …)` → erwarteter Pfad
  `"/tables/enable-MUTATED-verifier"`, während der reale Request weiterhin
  `/tables/enable` sendet). `rm -rf sdks/csharp/dist && make
  sdk-pack-csharp` real erneut ausgeführt: `dotnet test` schlägt sichtbar
  fehl (`Failed! - Failed: 1, Passed: 31, ..., Total: 32`), `docker build`
  bricht an der `dotnet test`-Stufe mit Exit 1 ab, `make` meldet
  `Error 1`/Gesamt-`EXIT=2`. Kein `sdks/csharp/dist/`-Verzeichnis entsteht
  (`mkdir -p` im Skript liegt nach dem `docker build`-Aufruf, wird bei
  dessen Fehlschlag nie erreicht — durch `ls sdks/csharp/dist/` real als
  „No such file or directory" bestätigt). Mutation danach exakt
  zurückgesetzt (Backup-Datei zurückkopiert), `git diff --stat`/`git status
  --porcelain` beide leer. **Ergebnis: unabhängig von der Reviewer-Probe
  real bestätigt — ein roter Test verhindert das Artefakt strukturell,
  nicht nur im vom Reviewer geprüften Einzelfall.**
- **Eigener Bau-Lauf 2 (Wiederherstellung):** `make sdk-pack-csharp`
  erneut → `EXIT=0`, `.nupkg` wieder vorhanden (29947 Bytes, unverändert).

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 `make gates` führt `sdk-pack-csharp` NICHT mit aus

Eigener, ungepiped `make gates`-Lauf (Log vollständig eingesehen):
`grep -c "sdk-pack" <log>` → `0`. Alle sechs Gates real grün:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 807 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check`-Modul `commits`: 807 Dateien, 0 Befunde; `commit-traceability.sh`: „OK — 5 Commit(s) in HEAD~5..HEAD, Betreffs ohne Struktur-ID" |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.70% erfüllt Schwelle 80%` |
| `generated-sync` | „OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)" |
| `a-check` | `gesamt: 0 Befund(e)` |

**Ergebnis: Checkbox berechtigt auf `[x]`; die strukturelle Trennung
(`sdk.mk` ohne `GATE_CHECKS`-Eintrag) trägt auch real im ausgeführten
Lauf.**

### 1.3 `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a` — Nachzug-Satz präzise

Eigene Lektüre (Zeile 183–187):

> „Für C#/NuGet ist die Frage beantwortet: `PgChangeFeed.Client`
> (`SPEC-026`) deckt HTTP-API und gRPC-Stream, real Docker-only
> paketierbar (`make sdk-pack-csharp`) und real geprüft
> (`PgChangeFeed.Client.0.1.0.nupkg`). Eine zweite Sprache oder ein
> zweiter Vertriebsweg bleibt offen — diese Kennung bleibt ihre Adresse."

Erfüllt exakt die geforderte Form: C#/NuGet nicht mehr offen, zweite
Sprache/Vertriebsweg bleibt offen, die Kennung `LH-FA-SST-009.a` selbst
steht als Überschrift unverändert (`### LH-FA-SST-009.a — Sprachmatrix und
Vertriebsweg offen`) — **nicht gestrichen**, wie gefordert. Der
ursprüngliche Absatz („… ist offen, ADR-pflichtig …") bleibt als
Fließtext-Aussage über die generelle Frage stehen; der Nachzug-Satz ergänzt
den bereits erledigten Teilfall, ohne die Kennungsadresse zu entwerten.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 `spec/pflichtenheft.md` §6 — `SPEC-026`-Zeile, kollisionsfrei

Eigene Lektüre §6, Zeile 611: `SPEC-026` | `PgChangeFeed.Client`
NuGet-Package … | SemVer 2.0, `0.x.y` (aktuell `0.1.0`) |
`sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` als
Metadaten-Quelle — Form konsistent mit den Nachbarzeilen (`SPEC-010`,
`SPEC-017`, `SPEC-020`, `SPEC-023`, `SPEC-024`).

Eigener `grep -c "SPEC-026" spec/*.md`: `spec/architecture.md:0`,
`spec/lastenheft.md:0`, `spec/pflichtenheft.md:3` (§1-Erwähnung,
§6-Zeile, §7-Historie) — genau eine Definitionsstelle, keine Kollision mit
einer anderweitig vergebenen Nummer.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.5 `spec/pflichtenheft.md` §7 — Historie-Zeilen, keine ADR-/Slice-Verweise

Eigene Lektüre der §7-Kopfregel (Zeile 617–619): „Regeln dieser Sektion:
**kein ADR- und kein Slice-Verweis.** Die Decken-Regel gilt für alle drei
Spec-Straten, auch hier — welche ADR eine Festlegung schärft, deklariert
die ADR aufwärts in ihrem `Schärft:`-Feld." Die beiden neuen Zeilen
(646/647) — „`SPEC-026` ergänzt: …" und „`LH-FA-SST-009.a` nachgezogen:
… (`make sdk-pack-csharp`, `PgChangeFeed.Client.0.1.0.nupkg`)" — enthalten
keine `ADR-*`- oder `slice-*`-Zeichenkette (eigener `grep -n "ADR-\|slice-"`
gegen beide Zeilen: kein Treffer). Die Kopfregel wird nicht nur zitiert,
sondern real eingehalten.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.6 `harness/README.md` §Werkzeuge — reale Zeile, Form konsistent

Eigene Lektüre der neuen `make sdk-pack-csharp`-Zeile: Bau-Kontext,
Stufen (`build`→`pack`→`pack-export`), Testlauf-vor-Artefakt-Zusage
(„ein roter Test bricht den `docker build` … ab, bevor `pack` je erreicht
wird — real geprüft"), `--build-context proto=proto`-Pflichtargument,
Export-Mechanismus (`tar`-Stream, `docker run --rm --network none … |
tar -x`), `.gitignore`-Hinweis, „Braucht Netz … Werkzeug statt Gate" —
dieselbe Form wie die `make examples-csharp`-Zeile (Bau-Kontext, Digest-
Pinnung, Netz-Begründung, `kein Gate, [ADR] … · seit slice-<name>"-Endung).
Die Bindung ist korrekt auf Festlegung 4 gesetzt (dieselbe ADR wie in
§2 dieses Berichts geprüft). Das referenzierte Target existiert real
(§1.1) — `AGENTS.md` §4 eingehalten.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.7 `sdks/csharp/.gitignore` trägt `dist/`

Bereits unter 1.1 real geprüft: nach dem eigenen Pack-Lauf zeigt `git
status --porcelain` das `.nupkg` **nicht** an.

### 1.8 `make gates` grün

Bereits unter 1.2 real geprüft: `EXIT=0`, ungepiped ermittelt.

### 1.9 Reconciliation entfällt

`harness/conventions.md`/Repo-Struktur bestätigt: keine
Reconciliation-Datei in diesem Repo — Entfall korrekt begründet, keine
eigene Prüfhandlung nötig.

### 1.10 Verbleibende `[ ]`-Checkboxen (Doku-Update entfällt, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen)

Der Slice-Plan liegt weiterhin unter `in-progress/`, §7 „Closure-Notiz"
trägt weiterhin Platzhalter (`<…>`). Für einen noch nicht geschlossenen
Slice ist das korrekt — Closure-Pflichten, keine Liefer-Punkte; ihr
`[ ]`-Zustand widerspricht dem übrigen DoD-Bild nicht. Diese Prüfung ist
laut Auftrag nicht Teil des Verifier-Umfangs für diesen Lauf
(„Was du NICHT tust: … keine Closure-Notiz").

## 2. ADR-0106-Konformität

### 2.1 §5 „Was diese ADR nicht ändert" — unberührte Pfade real bestätigt

Eigener Befehl:

```
git diff 89f62131..HEAD --stat -- docs/user/version.md .a-check.yml \
  examples/csharp/ sdks/csharp/PgChangeFeed.Client/Http/ \
  sdks/csharp/PgChangeFeed.Client/Grpc/
```

Ergebnis: **keine Ausgabe** — alle fünf genannten Pfade (Server-Version,
a-check-Konfiguration, Beispiel-Clients, beide Fläche-Verzeichnisse HTTP/
gRPC) sind seit dem letzten Closure-Commit (`89f62131`) unverändert.
Deckt sich mit `ADR-0106` §5: „Kein Produktionscode …", „`examples/csharp/**`
bleibt unverändert", „`docs/user/version.md` bleibt unberührt",
„`.a-check.yml` bleibt unberührt".

Vollständiger Diff-Stat (`git diff 89f62131..HEAD --stat`) zeigt exakt
acht geänderte Dateien: Slice-Plan, Review-Report, `harness/README.md`
(+1 Zeile), `harness/mk/sdk.mk` (neu), `sdks/csharp/.gitignore` (neu),
`sdks/csharp/Dockerfile` (+32/-14), `spec/pflichtenheft.md` (+9),
`tools/harness/sdk-pack-csharp.sh` (neu) — 403 Insertions, 14 Deletions
insgesamt, kein Treffer außerhalb des im Plan §3 deklarierten Umfangs.

### 2.2 Kein neues Gate

Bereits unter 1.1/1.2 real bestätigt: `GATE_CHECKS` (Root-`Makefile` +
alle Fragmente) enthält `sdk-pack-csharp` nicht.

### 2.3 Reales `.nupkg` als Smoke-Beleg — Größen-Konsistenz

29947 Bytes selbst gemessen (`ls -la`), deckt sich mit der im
Implementer-Commit und im Review-Report genannten Zahl — kein Drift
(`AGENTS.md` §3.12 Instanz A).

## 3. Commit-Traceability über den exakten Slice-Bereich

```
$ make doc-commits RANGE=89f62131..HEAD
d-check: 807 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=89f62131..HEAD
d-check: 807 Datei(en) geprüft, 0 Befund(e)
```

Beide Läufe grün, `EXIT=0`. `git log --oneline 89f62131..HEAD` eigenständig
gelesen: alle fünf Commit-Betreffe tragen `LH-FA-SST-009` und `ADR-0106`,
kein `SPEC-*`/`ARC-*`-Struktur-ID im Betreff.

## 4. Was diese Sitzung NICHT erneut geprüft hat

- Die externe Netz-Aussage im Review-Report zu `make docs-check`s
  Vollständigkeit über 806 vs. 807 Dateien (Reviewer maß 806 auf seinem
  Zwischenstand vor dem eigenen Report-Commit, diese Sitzung misst 807 auf
  `HEAD` nach dem Report-Commit `69d2dc00`) — kein Widerspruch, beide Werte
  sind Messungen ihres jeweiligen Standes (`AGENTS.md` §3.12 eingehalten,
  keiner wird als „der" Ist-Stand ausgegeben).
- Die Pipe-Disziplin-Probe des Reviewers (`docker run … nonexistent-image
  | tar -x`, `PIPE_EXIT=125`) wurde nicht erneut gefahren — die eigene
  Mutations-Probe dieser Sitzung deckt denselben Mechanismus (Exit ≠ 0
  bricht vor `mkdir`/`tar` ab) bereits strukturell ab.

## Verdikt

**DoD erfüllt** (für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten in §2 sind bewusst noch offen, siehe §1.10). Alle
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen:

1. `make sdk-pack-csharp` existiert, ist Docker-only, kein Gate — SECHSTE
   unabhängige Bau-Bestätigung (nach Implementer und Reviewer), inklusive
   einer von beiden Vorläufern unabhängigen, eigenen Mutations-Probe
   (anderer Test, andere Assertion) mit identischem Ergebnis: roter Test
   verhindert das Artefakt strukturell.
2. `make gates` führt `sdk-pack-csharp` real nicht mit aus
   (`grep -c "sdk-pack"` gegen den eigenen Gate-Log: `0`), alle sechs
   Gates grün.
3. `spec/pflichtenheft.md` §1/§6/§7 tragen die geforderten Nachzüge exakt
   in der verlangten Form — Kennung `LH-FA-SST-009.a` bleibt bestehen,
   `SPEC-026` kollisionsfrei, keine ADR-/Slice-Verweise in §7.
4. `harness/README.md` §Werkzeuge trägt die reale, formkonsistente Zeile.
5. `sdks/csharp/.gitignore` verhindert real, dass das Artefakt getrackt
   wird.
6. `ADR-0106` §5 (unberührte Pfade) real bestätigt: `docs/user/version.md`,
   `.a-check.yml`, `examples/csharp/**`, beide Fläche-Verzeichnisse
   unverändert seit `89f62131`.
7. Commit-Traceability und Immutabilität über den exakten Slice-Bereich
   (`89f62131..HEAD`) beide grün.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die im Plan selbst bereits als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz mit Steering-Loop-Lerneintrag,
Beobachtungs-Register-Prüfung, formaler Risiko-Ausgangs-Häkchen-Nachzug in
§2, die drei Paarungen bei der nächsten Welle-Closure
`welle-sdk-csharp-lh-fa-sst-009`).
