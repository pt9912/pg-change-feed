# Architect-Verdikt — Coverage-Messgegenstand beim Umzug der Vertragsfläche

**Datum:** 2026-09-17 · **Stand:** `7e598bb` · **Rolle:** Architect (eigener Kontext; die Zahlen
dieses Zugs sind selbst gemessen, nicht übernommen) · **Anlass:** `slice-097` LP3.

## 1. Verdikt: **in** — der erzeugte Code bleibt im Messgegenstand

Der Umzug bewegt den **Träger**, nicht den **Gegenstand**. Die Stufe `coverage` nimmt den
erzeugten Code am neuen Ort wieder auf: **`./gen/...` in der `go list`-Zeile**, mit
**unveränderter** Filterregel.

**Begründung aus der Eigenschaft, nicht aus der Zahl.** `ADR-0071` Punkt 1: *„Die tragende Regel
ist die **Eigenschaft**, nicht die Liste: Ein Paket, dessen Testlauf einen externen Dienst
voraussetzt, ist nicht Gegenstand dieses Gates."* Der erzeugte Vertrags-Code erfüllt die
**Ausschluss**-Eigenschaft nicht — sein Testlauf ist netzlos (gemessen: `--network none`,
Exit 0). Die Pfadangabe `./internal/... ./cmd/...` ist der **Träger** dieser Eigenschaft, und
Träger wandern mit dem Code.

**Die Beleg-Kette (alles `Accepted`):**
- `ADR-0071` Punkt 1 — die Eigenschaft als tragende Regel.
- `ADR-0076` §Konsequenzen, Folgepflicht „Umzugs-Slice": *„Der Messgegenstand (`ADR-0071`)
  **bleibt derselbe**, indem `./gen/...` in die `go list`-Zeile aufgenommen wird — sonst ändert
  sich die Zahl, ohne dass jemand den Gegenstand geändert hätte."* Die **spezifische** Klausel
  für genau diesen Fall.
- `ADR-0090` erklärt **kein** Supersedes gegenüber `ADR-0076`; nach der Reichweiten-Regel aus
  `ADR-0070` bleibt jene Klausel geltend.
- `ADR-0082` §Kontext (5) führt die **≈10** ungedeckten Proto-Runtime-Interna **innerhalb** des
  Gegenstands. Heute gemessen: 86 Statements, 76 gedeckt — **genau 10** ungedeckt.

**Der entscheidende Punkt: die Wahl „in" ändert am Gegenstand nichts.** Derselbe Paketsatz,
derselbe Nenner, dieselbe Quote. Die **Änderung wäre „out"**. „Nicht entscheiden" ist deshalb
**keine neutrale Option**: ohne die Wiederaufnahme schrumpft die Messfläche **still**.

## 2. Die gemessene Folge (beide Bandenden)

Netzlos, gepinnte Toolchain, Dedup über die Block-Position (dieselbe Zählbasis wie
`harness/sensors/coverage-gate.md` §Zählbasis; validiert, weil die gedruckte `total:`-Zeile mit
der Dedup-Rechnung übereinstimmt: 1613/1936 = 83,32 %, gedruckt 83,3 %).

| Arm | Nenner | gedeckt | Quote | gedruckt |
|---|---|---|---|---|
| **in** (Paket bleibt) | **1936** | 1613–1614 | 83,32–83,37 % | 83.3–83.4 % |
| **out** (Paket draußen) | **1850** (−86) | 1537–1538 | 83,08–83,14 % | 83.1 % |
| in, **ohne eigenes Testbinary** | 1936 | 1599 | 82,59 % | 82.6 % |

**Band: 1 Statement = 0,05 pp** (sieben Läufe des in-Arms, elf Messungen des out-Arms).
Abstand zur Endstufe 80 %: **in ≥ 3,3 pp, out ≥ 3,1 pp** — die Entscheidung ist **nicht** durch
Gate-Druck getrieben.

**Die dritte Zeile ist eine Bedingung, keine Variante.** Sie hält, weil die Testdatei
(`changestream_test.go`, `package streamv1_test`) **mit dem Baum umzieht**. Bleibt sie zurück,
hängt die Zahl an fremden Testbinaries (61 von 86 statt 76 von 86; 82,6 % statt 83,4 %).
**Der Umzug muss sie mitnehmen.**

## 3. Form: **keine neue ADR** — und warum

Der Fall ist entscheidungs-gedeckt; eine neue ADR wäre eine **zweite Quelle für eine bereits
getroffene Festlegung**. `ADR-0071` Punkt 1 trägt die Eigenschaft, `ADR-0076` §Konsequenzen
trägt genau diesen Fall wörtlich. §3.6 ist unberührt (keine Schwellen-, keine Mechanik- und —
gemessen — **keine Gegenstands-Änderung**: 1936 → 1936). Die **ablösungsbedürftige** Richtung
wäre „out".

**Verdikt-Typ (Modul 8, Konflikt-Pfad, Fall 1): die Entscheidung gilt und der Slice-Plan hat
falsch behauptet.** Das Übergabe-Artefakt ist deshalb ein **Plan-Diff** (Planner), nicht eine ADR.

## 4. Der Plan-Diff für LP3

> **LP3 — die Coverage-Fläche und ihre Träger sind nachgezogen.** Der erzeugte Code verlässt die
> **Paketliste** (`go list ./internal/... ./cmd/...`); er verlässt den **Messgegenstand nicht**:
> die tragende Regel ist die **Eigenschaft**, nicht die Liste (`ADR-0071` Punkt 1), und
> `ADR-0076` §Konsequenzen schreibt die Aufnahme vor. Die Stufe `coverage` führt deshalb
> **`./gen/...`** in der `go list`-Zeile, mit unveränderter Filterregel. Der Nenner **bleibt**
> beim Wert des `slice-097`-Laufs — gemessen **1936** (`out` wäre 1850) —, die Quote bei
> 83,3–83,4 % (Band 1 Statement). Nachzuziehen sind die Träger aus §6; die Zahl und ihr Lauf
> stehen im Bericht.

## 5. Die Zahl 1903 — Prüfung gegen die Messung

**Kein Träger führt 1903 als Ist-Stand.** Stehender Träger ist allein
`harness/sensors/coverage-gate.md` §Zählbasis; dort steht sie **mit ihrem Lauf** („Lauf
`slice-085`") und ausdrücklich als „**kein** Dauerwert". Auch `ADR-0082` §Kontext (2) führt sie
datiert. Alle übrigen Fundstellen sind Records. → **Sie halten** als datierte Messungen.

**Der Zustand ist trotzdem weiter als die Träger:** der Nenner ist heute **1936**, nicht 1903
(Differenz **+33**). Die einzige Produktionscode-Bewegung seit der `welle-20`-Closure ist
`slice-096` (`internal/bootstrap/config_file.go`, `wiring.go`) — Zuschreibung **abgeleitet** aus
`git diff 1999cfc..HEAD` über `internal/**`+`cmd/**` ohne Testdateien und ohne das ausgenommene
`receive`; kein eigener Coverage-Lauf an beiden Ständen. Kein Träger hat sie aufgenommen.

**Konsequenz:** LP3 darf **1903 nicht** als „der Nenner" verwenden. Der Wert gehört mit dem
`slice-097`-Lauf: **1936** (in) bzw. **1850** (out).

## 6. Träger — was mitzieht

**A. Träger des Coverage-Gegenstands (Entscheidung „in"):**
- `Dockerfile`, Stufe `coverage` — `./gen/...` **aufnehmen**, Filter unverändert; der
  Kommentarblock über der Stufe ist nachzuziehen. `--no-cache-filter`, `-covermode`,
  `-coverpkg`-Form bleiben unberührt.
- `harness/mk/coverage.mk` — Kommentarkopf nennt die Fläche; nachziehen. `THRESHOLD ?= 80` bleibt.
- `harness/sensors/coverage-gate.md` — §Vertrag (Pfadausdruck), §Grenze Punkt 1 (der
  `streamv1`-Absatz mit dem **alten Pfad**), §Zählbasis (Wert nach §5).
- `harness/README.md` §Sensors und `AGENTS.md` §4, je Zeile `make coverage-gate` — beide nennen
  den Pfadausdruck.

**B. Der Bau-Kontext — der Träger, den der Plan nicht nennt:**
- **`.dockerignore` — `!gen/` muss hinzu.** Ohne die Zeile ist `gen/` nicht im Kontext der Stufen
  `coverage` und `build`; `COPY . .` legt es nicht nach `/src`, und die Importkette ist nicht
  auflösbar. **Klassen-Grenze, die benannt gehört:** `!gen/` ist **keine** Ausnahme nach
  `ADR-0085` Festlegung 3 — die verlangt „gelesen, nicht gebaut" und „**genau eine Datei**".
  `gen/` ist ein **Bau-Eingang derselben Art wie `!internal/`/`!cmd/`**. Ein Kommentar, der hier
  `ADR-0085` zitiert, wäre eine **Fehl-Zitation**.
- `harness/image-hash.txt` — `.dockerignore` ist eine Build-Kontext-Datei: der Zug läuft
  `make image` vor seiner Closure (`ADR-0044` Punkt 3). Weicht der Digest ab, ist der
  Digest-Commit Teil des Slice. Dass er abweicht, ist **erwartet**, nicht gemessen.

**C. Träger des Umzugs** (vom Plan benannt): `proto/cdc/stream/v1/changestream.proto`
(`option go_package`), `gen/cdc/stream/v1/**` neu / der alte Baum entfällt (**inklusive
`changestream_test.go`** — Bedingung der Zahl, §2), die vier Importstellen, `.a-check.yml`.

- **Im Plan nicht genannt:** `harness/sensors/generated-sync.md` §Vertrag nennt den alten Pfad
  **im Präsens** — nach dem Umzug falsch. §3.13-Träger.

**D. Geprüft, kein Nachzug:** `harness/sensors/a-check.md` (nennt keinen Pfad) · `Makefile`
(`proto-generate`: `--go_opt=module=` bleibt) · `tools/harness/generated-sync.sh` (leitet den
Pfad aus `go_package`/Modul ab) · `tools/harness/run-store-tests.sh` (`go list ./...` folgt dem
Baum) · `docs/user/benutzerhandbuch.md` (kein Code-Pfad, keine Coverage-Zahl) ·
`docs/user/e2e-abdeckung.md`.
- `go list ./gen/...` ist als Träger tragfähig: **gemessen** — ein Baum mit einem Go-Paket unter
  `gen/cdc/stream/v1` und einem Nicht-Go-Verzeichnis `gen/csharp/` liefert genau das eine Paket,
  Exit 0.

**E. Fremde Träger — gemeldet, nicht mitgeändert:** `slice-095` §4 (Record eines vergangenen
Stands); die Register-Einträge `generierte-artefakte-ohne-sync-sensor` und
`arbeit-ueberholt-stehenden-traeger`; `ADR-0076`/`ADR-0082`/`ADR-0090` (immutable, datierte
Belege).

## 7. Nachbarfrage: DB-Adapter-Coverage — **nicht berührt**, belegt

`DB_COVERAGE_PKGS` in `tools/harness/db-coverage.sh` führt genau die drei Pakete, und
`harness/sensors/db-adapter-coverage.md` §Gegenstand nennt dieselben drei; `streamv1`/`gen/`
erscheint in keinem von beiden. Der Store-Lauf führt die Tests des erzeugten Pakets zwar mit
(`go list ./...`), aber ohne `-coverprofile` und außerhalb des DB-`-coverpkg`. **Keine zweite,
getrennte Folge** — benannt, nicht still.

## 8. Was anders kam als erwartet

1. **Der Plan behauptet Offenheit, wo die Entscheidungslage geschlossen ist** — Verdikt-Fall 1.
2. **`ADR-0090` enthält einen übergabe-artigen Satz** („Der Umzug-Slice bewegt ihn"), der als
   Aussage über die *Richtung* gelesen werden kann. Kein Supersedes; der Plan-Text soll diese
   Lesart nicht härten.
3. **Der Nenner ist nicht 1903, sondern 1936** (+33 aus `slice-096`), und das steht in keinem
   Träger.
4. **Der Plan nennt den Bau-Kontext nicht.** Ohne `!gen/` bricht `make image`/`make coverage-gate`.
5. **Die Zahl hängt an der mitziehenden Testdatei** — eine Bedingung, die bisher nirgends stand.
6. **`generated-sync.md` nennt den alten Pfad im Präsens** — ein §3.13-Träger, den der Plan nicht
   listet.
7. `coverage.mk` führt **keine** Nenner-Zahl, und `docs/user/` führt **keine** Coverage-Zahl —
   beide „geprüft, kein Nachzug".

## 9. Gates

`make gates` am Stand `7e598bb`: **Exit 0**, direkt und ungepiped gelesen. Im Log:
`generated-sync: OK — byte-gleich`, `a-check: gesamt: 0 Befund(e)`; der Coverage-Lauf dieses Zugs
lief separat über die Stufe `coverage` mit `THRESHOLD=80`, Exit 0, gedruckt `83.3%`.

## 10. Was nicht getan wurde

Kein Produktionscode, keine Paketliste, kein `.dockerignore`, kein Träger geändert; **keine** ADR
geschrieben (§3); kein Commit, keine Datei angelegt. Messungen in Wegwerf-Verzeichnissen unter
`/tmp`; Arbeitsbaum unverändert.
