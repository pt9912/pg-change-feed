# Review-Report: slice-sdk-csharp-nats-stream-client-flaeche (Fixrunde) — 2026-09-22

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Frisches, unabhängiges Review der **Fixrunde** — nicht nur
eine Bestätigung des Implementer-Berichts; alle drei Findings des
Vorgänger-Reports wurden eigenständig neu nachgemessen, nicht aus dem
Fixrunden-Commit übernommen.

**Gegenstand:** Fixrunden-Commit `8b6cef7c` (Eltern: `8ad9f6d0`, der
Review-Report-Commit; Groß-Eltern: `d86d1965`, der ursprüngliche
Feature-Commit), Slice `slice-sdk-csharp-nats-stream-client-flaeche`.
`git diff --stat 8ad9f6d0..8b6cef7c` zeigt genau die vier vom Implementer
genannten Dateien, 34 Insertions/14 Deletions.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um weitere HIGH-Klassen
ergänzt — u. a. „Slice-/Wellen-Chronik in Produktionscode-Kommentar",
„Zahl im Träger … driftend").
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-22.

**Vorgeschichte:** Ein früherer Reviewer-Lauf gegen `d86d1965` fand 2×
HIGH + 1× MEDIUM
(`docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md`):
F-1 (Kommentar-Chronik in `PgChangeFeed.Client.csproj`), F-2 (DoD
behauptete „sechzehn neue Tests", real gemessen 21 Testfälle/15
Testmethoden), F-3 (Beobachtungs-Beleg behauptete „neun Dateien", real
zehn). Ein früherer Versuch, dieses Fixrunden-Review durchzuführen, scheiterte
ohne jeden Commit an einem Rate-Limit-Fehler — dieser Lauf startet ohne
Vorarbeit komplett neu.

**Eingangs-Kontext:**

- `harness/README.md`, `AGENTS.md` (§3.7, §3.12, §3.13), `harness/conventions.md`
  (MR-000, MR-002)
- `.harness/skills/reviewer.md` (vollständig)
- `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md`
  (Vorgänger-Report, F-1/F-2/F-3 im Detail)
- `docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md`
  (vollständig gelesen, §2 DoD, §7 Closure-Notiz)
- `docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/{state.md,evidence/slice-sdk-csharp-nats-stream-client-flaeche.md}`
- `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` (komplett neu
  gelesen)
- `git show 8b6cef7c` / `git diff 8ad9f6d0..8b6cef7c` (vollständig gelesen)

**Eigenständig durchgeführte Prüfungen (nicht nur den Implementer-Bericht
übernommen):**

- **F-1 (Kommentar-Chronik) eigenständig neu gelesen:** Der Versionskommentar
  in `PgChangeFeed.Client.csproj:12-20` trägt keine
  „startete …"/„wurde … gehoben"-Erzählung mehr; er nennt den Zustand
  (`` `0.2.0` (ADR-0106 Festlegung 3, …) ``) und den Herkunfts-Anker
  „`(· seit slice-sdk-csharp-nats-stream-client-flaeche)`" — exakt die im
  Vorgänger-Report vorgeschlagene, aber nicht vorgeschriebene Form. Der
  benachbarte, bereits im Vorgänger-Review unbeanstandete
  `NATS.Net`-Kommentar (Zeilen 61-63) bleibt unverändert. Kein
  Vorher/Nachher-Bruch mehr im Diff.
- **F-2 (Testzahl) eigene, DRITTE unabhängige Nachmessung** — nicht nur den
  Implementer-Bericht übernommen: über `git worktree add --detach <scratch>
  4b93def4` (Elternstand vor der NATS-Fläche, nicht-destruktiv, kein
  Checkout im Haupt-Arbeitsbaum) und einen zweiten Worktree gegen
  `8b6cef7c` (Fixrunden-Stand) je ein isolierter, ungecachter
  `docker build --no-cache -f sdks/csharp/Dockerfile --target build
  --build-context proto=proto sdks/csharp`-Lauf gefahren: Elternstand
  `Passed: 49, Total: 49`, Fixrunden-Stand `Passed: 70, Total: 70` — Delta
  **21**, deckt sich mit der korrigierten Plan-Zeile und mit der Messung
  des Vorgänger-Reviews. Zusätzlich eigene Attribut-Zählung
  (`grep -rn "\[Fact\]\|\[Theory\]"` über `Nats/`): **15** Treffer über
  fünf Dateien (`ChangeMessageSchemaTests.cs`,
  `PgChangeFeedNatsStreamClientAuthBoundaryTests.cs`, `SubjectTests.cs`,
  `PgChangeFeedNatsStreamClientTests.cs` tragen Attribute;
  `FakeNatsClient.cs` keine eigene Testmethode) — deckt sich exakt mit dem
  jetzt im Plan stehenden „15 neue Testmethoden … in fünf Dateien … 21
  tatsächlich laufende neue Testfälle". Beide Worktrees per
  `git worktree remove --force` wieder entfernt, Hauptarbeitsbaum
  durchgehend unverändert (`git status --short` vor und nach der Prüfung
  leer).
- **F-3 (Dateizahl) selbst nachgezählt:** Die im Beobachtungs-Beleg
  (`evidence/slice-sdk-csharp-nats-stream-client-flaeche.md:16-40`)
  aufgezählten Fundstellen einzeln durchgegangen und gegen
  `git diff --stat 4b93def4..d86d1965` gehalten: `spec/pflichtenheft.md`,
  `sdks/csharp/README.md`, `PgChangeFeed.Client.csproj`,
  `sdks/csharp/Dockerfile`, `harness/mk/sdk.mk`,
  `tools/harness/sdk-pack-csharp.sh`, `harness/README.md`, `README.md`,
  `README.de.md`, `Sse/Models/Change.cs` — genau **zehn** distinkte
  Dateien, keine mehr, keine weniger. Die jetzt im Plan und im
  Beobachtungs-Beleg stehende Zahl „zehn" ist korrekt.
- **Vollständige Neulektüre des gesamten Fixrunden-Diffs**
  (`8ad9f6d0..8b6cef7c`, alle vier Dateien Zeile für Zeile) auf weitere,
  bisher übersehene Verstöße gegen `AGENTS.md` §3.7/§3.12 — keiner
  gefunden (siehe Negativbefunde). Eine dabei aufgefallene Detailfrage
  (die Zahl „acht Fundstellen" in derselben DoD-Zeile/demselben
  Beobachtungs-Beleg, die die Fixrunde nicht anfasste) wurde eigens
  nachgerechnet (Details unten, INFO — kein sauberer mechanischer
  Widerspruch, siehe F-Kandidat unten).
- **Backtick-Parität** aller vier geänderten Dateien selbst nachgezählt:
  `slice-sdk-csharp-nats-stream-client-flaeche.md` (Plan) 480, Evidenz-Datei
  124, `state.md` 192 — alle drei gerade/paarig;
  `PgChangeFeed.Client.csproj` 10 Backticks (gerade/paarig, kein Markdown,
  aber derselbe Zitierstil).
- **Umfangs-Prüfung der Fixrunde:** `git log --oneline d86d1965..8b6cef7c`
  zeigt zwei Commits (`8ad9f6d0` Review-Report, `8b6cef7c` Fixrunde) — die
  Aufgabenstellung nahm an, `d86d1965..8b6cef7c` zeige nur die vier vom
  Implementer genannten Dateien; real zeigt dieser Zwei-Commit-Range
  zusätzlich `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md`
  (aus `8ad9f6d0`, dem legitimen, separaten Review-Report-Commit) — kein
  Scope-Verstoß der Fixrunde selbst. Der eigentliche Fixrunden-Commit
  isoliert (`git diff --stat 8ad9f6d0..8b6cef7c`) zeigt exakt die vier
  genannten Dateien, `git diff` desselben Bereichs zeigt an
  `PgChangeFeed.Client.csproj` ausschließlich eine Kommentar-Umformulierung
  — `<Version>0.2.0</Version>` selbst unverändert, keine funktionale
  Code-Änderung an der NATS-Client-Implementierung.
- `make gates` ungefiltert laufen lassen (eigener Lauf, kein Pipe/Wrapper),
  Exit-Code direkt geprüft: **0** — u. a. `generated-sync: OK`,
  `a-check: gesamt: 0 Befund(e)`.

---

## Findings

Keine offenen HIGH/MEDIUM/LOW-Findings. Ein Beobachtungspunkt unterhalb der
Findings-Schwelle:

### INFO-1 — „acht Fundstellen" bleibt eine unscharfe, nicht sauber
mechanisch nachrechenbare Zahl (von der Fixrunde nicht berührt)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Herkunft von Aussagen) — Grenzfall, kein
  bestätigter Verstoß
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md:296`,
  `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md:66`,
  `state.md` (Passage zum 20. Beleg)
- `befund`: Dieselbe Textstelle, die F-3 korrigierte (Dateizahl neun→zehn),
  behauptet zusätzlich „acht real behobene Fundstellen". Die dort selbst
  aufgezählte Liste enthält beim wörtlichen Durchzählen der
  komma-getrennten Einträge neun Positionen, der begleitende
  Beobachtungs-Beleg enthält sieben Aufzählungspunkte (`grep -c "^- "`
  liefert 7) — keine der beiden mechanischen Zählweisen ergibt direkt
  acht. Eine plausible Rekonstruktion (Bündel `Dockerfile`+`sdk.mk` als
  ein per `grep` gefundenes Paar, `sdk-pack-csharp.sh` separat als das
  „per Lesen, nicht per grep" gefundene dritte Element derselben Gruppe,
  passend zu „zwei davon außerhalb des … `grep`-Musters gefunden") kommt
  auf acht — die Zahl ist damit nicht widerlegt, aber auch nicht mit
  derselben Klarheit verifizierbar wie die Dateizahl (`git diff --stat`)
  oder die Testzahl (`dotnet test`). Diese Fixrunde hat diesen Teil der
  Zeile nicht angefasst (nur „neun"→„zehn" korrigiert), er war auch nicht
  Gegenstand von F-3.
- `verifizierbar`: nein — keine eindeutige mechanische Zählregel für
  „Fundstelle" (im Unterschied zu „Datei") ist im Text selbst definiert.
- `klasse`: „Zahl im Träger — Zählweise nicht eindeutig definiert"
  (Nachbar-Klasse zu „Zahl im Träger driftet gegen die Messung", aber ohne
  bestätigten Widerspruch)

## Negativbefunde

- geprüft, ohne Befund: F-1 — `PgChangeFeed.Client.csproj`s
  Versionskommentar ist jetzt eine reine Zustandsaussage mit korrektem
  Herkunfts-Anker, keine Vorher/Nachher-Erzählung mehr.
- geprüft, ohne Befund: F-2 — eigene, dritte unabhängige Messung
  (`git worktree` gegen `4b93def4` und `8b6cef7c`, isolierter,
  ungecachter Docker-Bau) bestätigt exakt 49/70 Tests (Delta 21) und 15
  Testmethoden über fünf Dateien; die Plan-Zeile trägt jetzt die
  zutreffende Zahl.
- geprüft, ohne Befund: F-3 — eigenes Nachzählen der in der Evidenz-Datei
  genannten Pfade bestätigt exakt zehn distinkte Dateien.
- geprüft, ohne Befund: Umfang der Fixrunde — der isolierte
  Fixrunden-Commit (`8ad9f6d0..8b6cef7c`) ändert ausschließlich die vier
  vom Implementer genannten Dateien, keine funktionale Änderung an der
  NATS-Client-Implementierung, keine Berührung von
  `docs/user/benutzerhandbuch.md` oder anderen Trägern.
- geprüft, ohne Befund: Backtick-Parität aller vier geänderten Dateien —
  alle paarig.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) über den
  gesamten Fixrunden-Diff hinweg außerhalb von F-1 — keine neue Chronik,
  kein Konjunktiv über eine verworfene Alternative, kein abgebrochener
  Satz eingeführt.
- geprüft, ohne Befund: DoD-Checkbox „Review durchgeführt" — vom
  Implementer selbst im Fixrunden-Commit korrekt auf `[x]` gesetzt, mit
  Verweis auf den Vorgänger-Report und einer Zusammenfassung der drei
  behobenen Findings; kein Reviewer-Nachzug nötig (Skill
  §DoD-Checkbox-Nachzug ohne Fixrunde greift hier nicht — es gab eine
  Fixrunde).
- geprüft, ohne Befund: Traceability — Commit-Betreff `8b6cef7c` nennt
  `LH-FA-SST-009`/`ADR-0106`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Zahl im Träger — Zählweise nicht
eindeutig definiert" (INFO, kein bestätigter Widerspruch — nicht
steering-loop-zählend, da unterhalb der Findings-Schwelle)

## Verdikt

**Merge-blockierend:** nein — alle drei Findings des Vorgänger-Reports
(F-1 HIGH, F-2 HIGH, F-3 MEDIUM) sind real behoben und durch eigenständige
Nachmessung dieses Laufs bestätigt, nicht nur aus dem Fixrunden-Commit
übernommen. Der einzige verbleibende Punkt (INFO-1) ist keine bestätigte
Drift — er bleibt unterhalb der Findings-Schwelle, weil keine der beiden
naheliegenden mechanischen Zählweisen die Zahl eindeutig widerlegt.

**Übergabe:** Keine Fixrunde nötig. Die DoD-Zeile „Review durchgeführt"
ist bereits im geprüften Commit auf `[x]` gesetzt (regulär durch den
Implementer bei Schritt 21 seines eigenen Fixrunden-Laufs, nicht durch
den Skill-Mechanismus „DoD-Checkbox-Nachzug ohne Fixrunde" — dieser greift
nur, wenn *keine* Fixrunde stattfand). Dieser Report ist ein Lauf-Beleg;
er ersetzt keine Verifikation gegen die volle DoD — das bleibt
Verifier-Aufgabe (Modul 11).
