# Verifikationsbericht: slice-sdk-csharp-nats-stream-client-flaeche — 2026-09-22

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md`
§2/§3/§6/§7), in frischem Kontext. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-csharp-nats-stream-client-flaeche.md`](review-slice-sdk-csharp-nats-stream-client-flaeche.md)
und der
[Fixrunden-Nachprüfung](review-slice-sdk-csharp-nats-stream-client-flaeche-fixrunde.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** `4b93def4..c18f3545` (`d86d1965` Implementer-Inhalt,
`8ad9f6d0` Review-Report, `8b6cef7c` Fixrunde, `c18f3545`
Fixrunden-Nachprüfung), Slice `slice-sdk-csharp-nats-stream-client-flaeche`,
**letzter** Flächen-Slice der Welle `welle-sdk-csharp-vollabdeckung`,
`LH-FA-SST-009`, `LH-FA-SST-008`, `ADR-0106`, `SPEC-024`.

**Eingangs-Kontext (eigen gelesen, nicht aus Bericht übernommen):**
`harness/README.md`, `AGENTS.md`, `harness/conventions.md`, Slice-Plan §1–§8
vollständig, `ADR-0106` vollständig (Kontext, Entscheidung, Alternativen,
Konsequenzen, Geschichte), `spec/pflichtenheft.md` §`SPEC-024`,
`LH-FA-SST-009.a`, `SPEC-026`-Zeile und Änderungshistorie im Volltext,
beide Review-Reports (Erstlauf + Fixrunde) vollständig, der
Beobachtungs-Beleg `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md`,
alle fünf neuen `Nats/**`-Quelldateien und die fünf neuen Testdateien,
`PgChangeFeed.Client.csproj`, `Directory.Packages.props`,
`docs/user/benutzerhandbuch.md`-Diff.

---

## 1. DoD-Vertrag (§2) — jede Zeile einzeln geprüft

### 1.1 `LH-FA-SST-009` erfüllt — NATS-Vollinhalts-Fläche, zehn `SPEC-024`-Felder, Testzahl

**Feld-für-Feld-Gegenprobe** (eigene Instanz, unabhängig vom Review):
`spec/pflichtenheft.md` §`SPEC-024` (Zeile 537, Nachrichteninhalt-Zeile)
nennt zehn Felder — `change_id`, `transaction_id`, `source_table_id`,
`sequence`, `operation`, `old_image`, `new_image`, `schema_version`,
`schema`, `table`. `sdks/csharp/PgChangeFeed.Client/Nats/Models/Change.cs`
trägt exakt dieselben zehn `JsonPropertyName`-Attribute, keins fehlt, keins
zusätzlich — bestätigt.

**Testzahl — vierte unabhängige Nachmessung.** Über zwei nicht-destruktive
`git worktree`-Checkouts (`git worktree add --detach`, kein Checkout im
Hauptbaum) je ein isolierter, ungecachter
`docker build --no-cache -f sdks/csharp/Dockerfile --target build
--build-context proto=proto sdks/csharp`-Lauf:

| Stand | Commit | Ergebnis |
|---|---|---|
| Elternstand | `4b93def4` | `Passed: 49, Total: 49` |
| Fixrunden-Stand (= aktueller `HEAD`-Inhalt) | `8b6cef7c` | `Passed: 70, Total: 70` |

Delta **21** — deckt sich exakt mit den drei vorigen Messungen
(Implementer-Plan-Nachzug, Erstreview, Fixrunden-Review). Zusätzlich eigene
Attribut-Zählung: `grep -rn "\[Fact\]\|\[Theory\]"` über
`PgChangeFeed.Client.Tests/Nats/` liefert **15** Treffer über vier der fünf
Dateien (`FakeNatsClient.cs` trägt keine eigene Testmethode) — deckt sich
mit der jetzt im Plan stehenden Formulierung „15 neue Testmethoden … 21
tatsächlich laufende neue Testfälle". Beide Worktrees per
`git worktree remove --force` entfernt, Hauptarbeitsbaum vor und nach der
Prüfung unverändert (`git status --short` leer).

**Ergebnis: berechtigt auf `[x]`.** DoD-Punkt 1 trägt jetzt die korrekt
gemessene Zahl, kein Rest-Drift (F-2 des Erstreviews ist real geschlossen).

### 1.2 `NATS.Net`-Pin real gesetzt und aktuell

`sdks/csharp/Directory.Packages.props` trägt
`<PackageVersion Include="NATS.Net" Version="3.2.0" />`. Eigene
Nachmessung: `curl https://api.nuget.org/v3-flatcontainer/nats.net/index.json`
— jüngster Eintrag in der Versionsliste ist weiterhin `3.2.0`, kein
neuerer Kandidat. **Ergebnis: berechtigt auf `[x]`.**

### 1.3 `<Version>0.2.0</Version>` real gesetzt, Kommentar-Chronik geprüft (vierte Bestätigung)

`PgChangeFeed.Client.csproj:22-23` trägt `<Version>0.2.0</Version>`
(vorher `0.1.0`, additive Erweiterung, kein Breaking Change). Der
begleitende Kommentar (Zeilen 12-20) wurde eigenständig gelesen: er nennt
den **Zustand** — `` „Version `0.2.0` (ADR-0106 Festlegung 3, unabhängig von
docs/user/version.md) — additive, rückwärtskompatible Erweiterung um die
SSE- und NATS-Vollinhalts-Client-Fläche neben der bestehenden HTTP-/
gRPC-Fläche (· seit slice-sdk-csharp-nats-stream-client-flaeche), kein
Breaking Change …" `` — **keine** „startete …/wurde … gehoben"-Erzählung mehr.
Dies ist die **vierte** unabhängige Bestätigung (Implementer-Fix,
Fixrunden-Review, diese Verifikation), dass F-1 des Erstreviews real und
dauerhaft behoben ist. **Ergebnis: berechtigt auf `[x]`.**

### 1.4 Reales `.nupkg` mit allen vier Client-Flächen

**Eigener, vierter unabhängiger Lauf** von `bash tools/harness/sdk-pack-csharp.sh`
(nach `rm -f sdks/csharp/dist/*.nupkg`): erzeugt
`sdks/csharp/dist/PgChangeFeed.Client.0.2.0.nupkg`, 37451 Bytes — exakte
Byte-Übereinstimmung mit der im Plan behaupteten Größe.
`unzip -l` bestätigt `lib/net10.0/PgChangeFeed.Client.dll`. Über die reine
Existenzprüfung hinaus: die extrahierte DLL wurde per `strings` auf alle
vier Namensräume geprüft — `PgChangeFeed.Client.Http.*`,
`PgChangeFeed.Client.Grpc.*`, `PgChangeFeed.Client.Sse.*` und
`PgChangeFeed.Client.Nats.*` (inkl. `PgChangeFeedNatsStreamClient`,
`Nats.Models.Change`) sind alle in derselben Assembly vorhanden — kein
Namensraum fehlt. **Ergebnis: berechtigt auf `[x]`.**

### 1.5 Träger-Nachzug `spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md`/`README.md`/`README.de.md`/`sdks/csharp/README.md`

`spec/pflichtenheft.md` Zeile 183-187 (`LH-FA-SST-009.a`) und Zeile 622
(`SPEC-026`-Tabellenzeile) sowie die Änderungshistorie-Zeile vom
2026-09-22 (Zeile 665) bestätigen den Nachzug: „`PgChangeFeed.Client` deckt
HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt … Version auf `0.2.0`
gehoben". `docs/user/benutzerhandbuch.md` Version 1.38, Changelog-Zeile
1.38 vorhanden, beide neuen `**SDK:**`-Absätze (SSE §"Zugriff über
Server-Sent-Events", NATS §"Zugriff über den NATS-Vollinhalts-Stream")
gelesen — nennen korrekt `PgChangeFeedSseClient`/
`PgChangeFeedNatsStreamClient`, `StreamChangesAsync`, „alle zehn Felder",
`BuildSubject`/`BuildSourceSubject`. `README.md`/`README.de.md`
§Distribution-Zeile nennt „HTTP API, gRPC stream, SSE stream, and NATS
full-content stream" bzw. deutsches Äquivalent. `sdks/csharp/README.md`
Zeile 5/9 nennt alle vier Flächen inkl. NATS-Vollinhalts-Stream.

**Eigener Suchlauf mit anderen Formulierungsvarianten** als die drei
vorigen Läufe (Implementer-Suchlauf, Reviewer-Suchlauf 1, Reviewer-Suchlauf
2 der Fixrunde):

```
grep -rn "nur HTTP\|only HTTP\|zwei Client-Flächen\|two client surfaces\|
two surfaces\|0\.1\.0\.nupkg\|PgChangeFeed\.Client\.0\.1\.0\|
deckt HTTP-API und gRPC\|HTTP and gRPC only" \
  --include="*.md" --include="*.cs" --include="*.sh" --include="*.mk" \
  --include="*.props" --include="*.csproj" --include="Dockerfile" \
  spec/ docs/ harness/ sdks/csharp/ README.md README.de.md tools/
```

Ergebnis: alle verbleibenden Treffer auf „0.1.0"/„nur HTTP-API und
gRPC-Stream" liegen in `Accepted`-ADRs (unberührbar, `AGENTS.md` §3.5),
`docs/plan/planning/done/**`/`docs/reviews/**` (historische Records),
`docs/plan/planning/open/slice-sdk-{kotlin,python}-nats-stream-client-flaeche.md`
(künftige, noch nicht umgesetzte Slice-Pläne — der dort zitierte Satz
bezieht sich auf die jeweils **andere** Sprache und ist dort zum
gegenwärtigen Zeitpunkt weiterhin korrekt), oder auf `pgchangefeed-kotlin`
(deckt real weiterhin nur HTTP-API und gRPC-Stream — zutreffend). **Keine
vierte, von allen drei vorigen Suchläufen übersehene Stelle gefunden.**
Ergebnis deckungsgleich mit Reviewer-Suchlauf 2 (Fixrunde), aber mit
eigenen, neuen Suchbegriffen erzielt. **Ergebnis: berechtigt auf `[x]`.**

### 1.6 Kein Import aus `internal/**`/`cmd/**`/`gen/**`

Eigener Lauf: `grep -rn "internal/\|cmd/\|gen/" sdks/csharp/PgChangeFeed.Client/Nats/`
→ kein Treffer, Exit `1`. **Ergebnis: berechtigt auf `[x]`.**

### 1.7 `make gates` grün

Siehe §3 unten. **Ergebnis: berechtigt auf `[x]`.**

### 1.8 Review durchgeführt, kein offenes HIGH

Erstreview: 2× HIGH (F-1 Kommentar-Chronik, F-2 Testzahl-Drift) + 1×
MEDIUM (F-3 Dateizahl-Drift). Fixrunden-Nachprüfung: alle drei Findings
durch eigenständige Nachmessung des Reviewers real bestätigt behoben, 0
HIGH/MEDIUM/LOW, 1 INFO. Eigene Prüfung bestätigt dasselbe Bild (siehe
1.1/1.3 oben und §2 unten für die INFO-Frage). **Ergebnis: berechtigt auf
`[x]`.**

### 1.9 Reconciliation, Beobachtungs-Register, Closure-Notiz, §6-Risiko-Ausgänge, drei Paarungen

Reconciliation „entfällt" — zutreffend, keine Reconciliation-Datei in
diesem Repo. Beobachtungs-Register: neuer Beleg
`evidence/slice-sdk-csharp-nats-stream-client-flaeche.md` zu
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Zähler 20×, real in
`state.md` nachgezählt: elf benannte + fünf nachträglich ergänzte + der
achtzehnte/neunzehnte/zwanzigste Beleg = 20 Dateien unter `evidence/`,
plausibel). §6-Risiken siehe §4 unten. Drei Paarungen: korrekt auf die
Welle-Closure verwiesen (Slice liegt weiterhin nicht in `done/`).

**Kein DoD-Blocker in dieser Gruppe.**

---

## 2. „Acht Fundstellen" — eigene, unabhängige Entscheidung zur vom Fixrunden-Reviewer offen gelassenen INFO-Frage

**Befund nachvollzogen:** `state.md` und der Beobachtungs-Beleg behaupten
„acht real behobene Fundstellen … über zehn Dateien". Die Dateizahl
(„zehn") ist über `git diff --stat 4b93def4..d86d1965` mechanisch
eindeutig nachprüfbar und wurde bereits dreifach bestätigt (Reviewer F-3,
Fixrunden-Review, diese Verifikation — siehe §1.5). Die Zahl „acht"
dagegen zählt eine Einheit („Fundstelle"), die im Text selbst **nicht**
mit einer eindeutigen Zählregel definiert ist — anders als „Datei", die
trivial auf `git diff --stat`-Pfade abbildet.

**Eigener Nachzähl-Versuch, unabhängig vom Fixrunden-Reviewer:**

- `grep -c "^- "` über die Evidenz-Datei: **7** Bullets (selbst
  nachvollzogen, siehe unten).
- Wörtliche Aufzählung der im Text benannten Einzeldateien/-Aussagen
  (jede referenzierte Datei einzeln gezählt, auch wenn mehrere in einem
  Bullet gebündelt sind): `spec/pflichtenheft.md`, `sdks/csharp/README.md`,
  `PgChangeFeed.Client.csproj`, `sdks/csharp/Dockerfile`,
  `harness/mk/sdk.mk`, `tools/harness/sdk-pack-csharp.sh`,
  `harness/README.md`, `README.md`, `README.de.md`, `Sse/Models/Change.cs`
  — **10** Einzeldateien (deckt sich mit „zehn Dateien").
- Der Text selbst nennt `tools/harness/sdk-pack-csharp.sh` explizit als
  „dritten Fundort" und `Sse/Models/Change.cs` als „vierten Fundort" —
  beides Fundorte, die **außerhalb** des vorgeschriebenen `grep`-Musters
  lagen. Das impliziert eine interne Zählung, die nicht mit den Bullets
  (7) und nicht mit den Einzeldateien (10) übereinstimmt, sondern
  offenbar Bündel wie „`Dockerfile`+`sdk.mk` als ein per `grep` gefundenes
  Paar" von den beiden separat benannten „dritten"/„vierten" Fundorten
  unterscheidet — eine Rekonstruktion, die auf 8 kommt, aber nicht aus dem
  Text selbst als geschlossene Rechenregel ableitbar ist.

**Eigene Entscheidung:** Ich schließe mich der Einordnung des
Fixrunden-Reviewers an (INFO, kein bestätigter Zahlen-Drift) — mit einer
eigenen Begründung, nicht nur durch Übernahme: Der entscheidende
Unterschied zu F-2 (Testzahl) und F-3 (Dateizahl) ist, dass für beide
**eine einzige, im Repo bereits etablierte, mechanisch eindeutige
Zähleinheit** existiert (`dotnet test`-Ausgabe; `git diff --stat`-Pfade) —
gegen die die behauptete Zahl **falsifizierbar** war und sich als falsch
erwies. Für „Fundstelle" gibt es **keine** einzige mechanisch eindeutige
Zähleinheit im Text: „ein Bullet", „eine Datei" und „eine Aussage
innerhalb einer Datei" sind alle drei plausible, aber unterschiedliche
Lesarten, und keine ist im Text als die verbindliche deklariert. Eine Zahl
ohne definierte Zähleinheit ist nicht dasselbe wie eine Zahl, die gegen
eine definierte Messung driftet (`AGENTS.md` §3.12 „Zahl im Träger …
driftend" setzt eine vergleichbare Messgröße voraus) — hier fehlt die
Messgröße selbst, nicht nur ihre korrekte Anwendung.

**Verdikt zu diesem Punkt: kein eigener Finding-Punkt, keine
Korrektur-Pflicht vor Closure.** Empfehlung an den Planner (nicht
blockierend): Bei der nächsten inhaltlichen Berührung dieser Zeile (z. B.
bei der Welle-Closure oder einem künftigen Beleg derselben
`BEO-PGC/arbeit-ueberholt-stehenden-traeger`-Familie) die Formulierung
„acht real behobene Fundstellen" entweder mit einer expliziten
Zähldefinition zu versehen oder ersatzlos zu streichen und ausschließlich
die mechanisch eindeutige Zahl „zehn Dateien" (samt der Liste selbst) zu
tragen — die Liste selbst ist bereits vollständig und korrekt, nur die
zusammenfassende Zahl „acht" ist die vermeidbare Unschärfe.

---

## 3. `make gates` real, ungepiped, in dieser Sitzung ausgeführt

```
$ git log -1 --format=%H
c18f3545...
$ make gates > /tmp/verifier-gates.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 896 Datei(en) geprüft, 0 Befund(e)` |
| `a-check` | `gesamt: 0 Befund(e)` |
| `commit-traceability` | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` |

**Ergebnis: DoD-Checkbox „`make gates` grün" berechtigt auf `[x]`.**

---

## 4. `git diff af5e59a5..HEAD` — Scope-Sauberkeit (gesamter Slice-Umfang, af5e59a5 = SSE-Closure-Commit, Elternstand)

```
$ git diff --name-status af5e59a5..HEAD
M   README.de.md
M   README.md
A   docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md
A   docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-sdk-csharp-nats-stream-client-flaeche.md
M   docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/state.md
D   docs/plan/planning/open/slice-sdk-csharp-nats-stream-client-flaeche.md
A   docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche-fixrunde.md
A   docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md
M   docs/user/benutzerhandbuch.md
M   harness/README.md
M   harness/mk/sdk.mk
M   sdks/csharp/Directory.Packages.props
M   sdks/csharp/Dockerfile
A   sdks/csharp/PgChangeFeed.Client/Nats/Models/Change.cs
A   sdks/csharp/PgChangeFeed.Client/Nats/PgChangeFeedNatsMalformedMessageException.cs
A   sdks/csharp/PgChangeFeed.Client/Nats/PgChangeFeedNatsStreamClient.cs
A   sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Nats/ChangeMessageSchemaTests.cs
A   sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Nats/FakeNatsClient.cs
A   sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Nats/PgChangeFeedNatsStreamClientAuthBoundaryTests.cs
A   sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Nats/PgChangeFeedNatsStreamClientTests.cs
A   sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Nats/SubjectTests.cs
M   sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj
M   sdks/csharp/PgChangeFeed.Client/Sse/Models/Change.cs
M   sdks/csharp/README.md
M   spec/pflichtenheft.md
M   tools/harness/sdk-pack-csharp.sh
```

26 Dateien, jede einzelne einem der DoD-/Plan-Punkte zuordenbar: neue
`Nats/**`-Produktionsklassen (3) und Testdateien (5), Paket-/Projekt-Meta
(`Directory.Packages.props`, `.csproj`), Träger-Nachzug (`spec/pflichtenheft.md`,
`sdks/csharp/README.md`, `harness/README.md`, `README.md`/`README.de.md`,
`Dockerfile`, `harness/mk/sdk.mk`, `tools/harness/sdk-pack-csharp.sh`,
`Sse/Models/Change.cs`-Kommentar), Handbuch-Update, Observation-Register
(state.md + neue Evidenz-Datei), zwei Review-Reports, der Slice-Plan selbst
(Move `open` → `in-progress`, sichtbar als D+A relativ zu `af5e59a5`).
**Kein unerwarteter Treffer, keine Berührung von `internal/**`, `cmd/**`,
`gen/**`, `spec/architecture.md`, `.a-check.yml` oder `sdks/python/**`/
`sdks/kotlin/**`.** Scope ist sauber.

---

## 5. Backtick-Parität aller in diesem Slice geänderten Markdown-Dateien (eigenständig nachgezählt)

| Datei | Backtick-Anzahl | Parität |
|---|---|---|
| `README.de.md` | 62 | gerade |
| `README.md` | 62 | gerade |
| `docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md` | 480 | gerade |
| `docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-sdk-csharp-nats-stream-client-flaeche.md` | 124 | gerade |
| `docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/state.md` | 192 | gerade |
| `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche-fixrunde.md` | 212 | gerade |
| `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md` | 494 | gerade |
| `docs/user/benutzerhandbuch.md` | 1848 | gerade |
| `harness/README.md` | 1804 | gerade |
| `sdks/csharp/README.md` | 54 | gerade |
| `spec/pflichtenheft.md` | 1470 | gerade |

Alle elf Dateien: gerade Anzahl, paarig. Deckt sich mit den bereits von
beiden Reviewer-Läufen gemeldeten Werten für die vier bzw. drei von ihnen
geprüften Dateien (`slice`-Plan 464→480 nach Fixrunde, `harness/README.md`
1804, `spec/pflichtenheft.md` 1470 unverändert seit dem Erstreview).

---

## 6. §6-Risiken — Ausgänge bestätigt

1. „Ein gefakter NATS-Verbindungsfehler-Test könnte das reale
   Token-Ablehnungsverhalten des NATS-Servers nicht exakt nachbilden." —
   **Ausgang: weiter offen**, korrekt so geführt. Kein realer NATS-Server-
   Rundlauf-Beleg in diesem Slice; ein solcher bleibt strukturell
   `make test-integration`s bestehendem `tools/harness/natsstreamsub`-Muster
   vorbehalten (dieselbe Teststrategie wie HTTP/gRPC/SSE zuvor). Zu Recht
   nicht als erledigt markiert.
2. „Der Träger-Nachzug in `spec/pflichtenheft.md` könnte eine Stelle
   übersehen." — **Ausgang: weiter offen** für `spec/pflichtenheft.md`
   selbst (kein Fund in dieser Verifikation, siehe §1.5); für den
   breiteren Suchraum korrekt als „eingetreten und behoben" geführt (zwei
   Fundstellen außerhalb des `grep`-Musters, real behoben). Beide
   Teil-Ausgänge sind zutreffend differenziert, nicht pauschal
   „erledigt".
3. „Version-Hebung ohne SemVer-Minor-Konvention" — **Ausgang: entfallen**,
   zutreffend begründet (additive Erweiterung ist im SemVer-Vokabular
   unmissverständlich ein Minor-Bump).

**Docker-Bau-Kopplungs-Beobachtung** (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`):
eigenständig im Register nachgelesen — Zähler bei **3×** (zwei
`examples/**`-Vorgänge + `slice-sdk-csharp-sse-client-flaeche`), Ausgang
korrekt „noch nicht zugewiesen, Architect-Entscheidung bei Welle-Closure
verwiesen". Dieser Slice hat **keinen** vierten Beleg erzeugt — zutreffend
so berichtet, kein isolierter Bau-Versuch ohne `--build-context proto=proto`
in diesem Zug unternommen.

**Kein DoD-Blocker in dieser Gruppe.**

---

## Verdikt

**DoD-konform: ja.** Alle neun geprüften DoD-Zeilen aus §2 sind für den
aktuellen `in-progress`-Stand berechtigt auf `[x]`. Alle beauftragten
Prüfpunkte wurden real und unabhängig nachgemessen, keiner aus Bericht
oder Commit-Message unbesehen übernommen:

1. Alle zehn `SPEC-024`-Nachrichtenfelder real gegen `Change.cs` gehalten
   — keine Abweichung.
2. Testzahl **vierte** unabhängige Messung (`git worktree`, isolierter,
   ungecachter Docker-Bau gegen Parent `4b93def4` und HEAD-Inhalt
   `8b6cef7c`): 49→70, Delta 21, deckt sich mit allen drei vorigen
   Messungen; 15 Testmethoden über vier von fünf Dateien selbst gezählt.
3. `NATS.Net` 3.2.0 real gegen `api.nuget.org` erneut bestätigt, weiterhin
   jüngste stabile Version.
4. `<Version>0.2.0</Version>` real gesetzt; Kommentar-Chronik-Fix (F-1)
   **vierte** unabhängige Bestätigung: reine Zustandsaussage mit
   Herkunfts-Anker, keine „startete/wurde gehoben"-Erzählung.
5. **Vierter** unabhängiger `bash tools/harness/sdk-pack-csharp.sh`-Lauf:
   `PgChangeFeed.Client.0.2.0.nupkg`, 37451 Bytes, alle vier
   Client-Namensräume (`Http`, `Grpc`, `Sse`, `Nats`) real in der
   extrahierten DLL nachgewiesen (`strings`-Prüfung, über die reine
   Existenzprüfung hinaus).
6. Träger-Nachzug: eigener Suchlauf mit dritten, neuen Formulierungs-
   varianten — keine vierte, übersehene Stelle gefunden.
7. Kein Import aus `internal/**`/`cmd/**`/`gen/**` in `Nats/**`.
8. `make gates` eigenständig, ungepiped, `EXIT=0` — alle sechs Gates grün
   (baseline-verify 54 Dateien OK, docs-check 896/0, a-check 0,
   commit-traceability OK, coverage-gate 82.80% ≥ 80%, generated-sync OK).
9. `git diff --name-status af5e59a5..HEAD` (voller Slice-Umfang, 26
   Dateien) — Scope sauber, kein unerwarteter Treffer.
10. Backtick-Parität aller elf in diesem Slice geänderten Markdown-Dateien
    — alle gerade/paarig.
11. §6-Risiken: alle drei Ausgänge zutreffend differenziert geführt, kein
    Risiko fälschlich als erledigt markiert; die
    Docker-Bau-Kopplungs-Beobachtung bleibt korrekt bei 3× ohne vierten
    Beleg aus diesem Slice.

**Eigene Entscheidung zur „acht Fundstellen"-Frage (Auftrag Punkt 2):**
**keine bestätigte Zahlen-Drift, kein eigener Finding-Punkt.** Anders als
die bereits korrigierten Zahlen „21 Testfälle" und „zehn Dateien" hat
„Fundstelle" keine im Text definierte, mechanisch eindeutige Zähleinheit
— die Diskrepanz ist eine vermeidbare Formulierungs-Unschärfe, keine
Messung gegen eine falsche Zahl. Empfehlung (nicht blockierend) an den
Planner: bei nächster inhaltlicher Berührung dieser Zeile entweder die
Zähleinheit explizit machen oder die Zahl „acht" zugunsten der bereits
mechanisch verifizierten „zehn Dateien" streichen.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die im Plan selbst bereits als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz liegt bereits ausgefüllt vor; die drei
Paarungen sind explizit an die Closure von `welle-sdk-csharp-vollabdeckung`
verwiesen — diese Welle kann nach diesem letzten Flächen-Slice geschlossen
werden). Keine Fixes durchgeführt, kein `git mv` nach `done/`, keine
Wellen-Closure — das bleibt Planner-/Architect-Aufgabe.
