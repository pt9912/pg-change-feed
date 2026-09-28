# Verifikations-Report: gRPC-`Administration`-Service (ADR-0130) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/ADR-Konformitätsprüfung
+ Plan-vs-Code-Diff + Gates, in frischem Kontext, nach der Fixrunde. Kein
Slice-Plan trägt diesen Zug — er lief als direkter Architect→Implementer→
Reviewer-Auftrag über [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md);
die DoD dieses Reports ist die ADR selbst (Entscheidung, Fitness Function,
Folgepflichten) plus die Findings des Code-Reviews.

**Gegenstand:** `git diff b317d529..HEAD` — sechs Commits: `9160ad32` (proto +
generierter Code), `2c400b8a` (Handler + Interceptor + Server), `23ad7424`
(E2E-Test), `d888c01e` (§3.13-Suchlauf-Nachzug), `1bbe3e06` (Review-Report),
`60d845d3` (Fixrunde: Handbuch-Zug + Kommentar-Kennung), `08ae3e84`
(`id-unlinked`-Korrektur am Review-Report selbst — unabhängig von der
Slice-DoD, aber Voraussetzung für ein grünes `make gates`).

**Eingangs-Kontext:**

- [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) (Accepted) — vollständig gelesen, alle sechs Teilfragen,
  Fitness-Function-Tabelle, Folgepflichten
- [`review-grpc-administration-server.md`](review-grpc-administration-server.md) — vollständig gelesen (1 HIGH F-1, 1 MEDIUM F-2,
  1 INFO F-3, Verdikt merge-blockierend wegen F-1)
- `git show 60d845d3` und die diff-Ausschnitte gegen `docs/user/benutzerhandbuch.md`
  und `gen/cdc/administration/v1/administration_test.go`
- `AGENTS.md` §3.7, §3.9, §3.12 — insbesondere die Nutzerdirektive „keine
  ADR-/Review-Verweise im benutzerseitigen Fließtext des Benutzerhandbuchs“

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien`; `db-package-lists-check: OK`; `d-check: 1377 Datei(en) geprüft, 0 Befund(e)` (Modul-Bündel); `--enable commits`: `1377 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `coverage-gate: OK — Coverage 80.80% erfüllt Schwelle 80%`; `generated-sync: OK — byte-gleich` (beide `.proto`-Quellen, alle vier generierten Dateien geprüft); `a-check: gesamt: 0 Befund(e)` |
| `make kommentar-kennungen DIFF=b317d529` | **Exit 0** | keine Ausgabe, kein Kandidat außerhalb von `gen/` (strukturell ausgenommen, siehe F-2-Einordnung) |
| `git status --short` | **leer** | Arbeitsbaum sauber vor und nach diesem Lauf |
| `git diff --stat b317d529..HEAD` | **19 Dateien, deckungsgleich mit ADR-Folgepflichten** | siehe §3 |

Nicht erneut gefahren: `make test-integration` (schwerer ~20+-Minuten-Lauf) —
der Reviewer hat ihn bereits **selbst und frisch** bis zum Ende
(„Lauf abgeschlossen“) durchlaufen lassen und den neuen gRPC-Administration-
Abschnitt real beobachtet (`REGISTERED`/`LISTED`/`REJECTED code=Unauthenticated`/
`REJECTED code=PermissionDenied`); kein Produktionscode-Diff seit diesem
Review-Lauf (`git diff --name-only 1bbe3e06..HEAD` zeigt nur
`docs/user/benutzerhandbuch.md`, `gen/cdc/administration/v1/administration_test.go`
und den Review-Report selbst — kein `.go`-Nicht-Test-Diff, kein
Runner-Skript-Diff) — der reale, grüne E2E-Beleg bleibt gültig,
**übernommen**, Anker: Review-Report §„Bestätigte Implementer-Behauptungen“.

## 2. F-1 (HIGH) — Handbuch-Zug: geprüft, behoben

**Behauptung der Fixrunde:** neuer §4-Abschnitt „Zugriff über die
gRPC-Verwaltungs-API“, §5-`CDC_GRPC_ADDR`-Zeile korrigiert, Versionshistorie
nachgezogen.

**Eigen geprüft (nicht nur gelesen):**

- `git show 60d845d3 -- docs/user/benutzerhandbuch.md` zeigt den vollständigen
  Diff: neuer Abschnitt zwischen den bestehenden Abschnitten „Zugriff über den
  gRPC-Change-Stream“ (Zeile 1238) und „Zugriff über Server-Sent-Events“
  (jetzt Zeile 1407, vorher 1342) — `grep -n "^### Zugriff über"` bestätigt
  die neue Zeile 1342 „Zugriff über die gRPC-Verwaltungs-API“ zwischen beiden,
  saubere Einfügung ohne Kollateralschaden an Nachbarabschnitten.
- **Alle neun RPCs vorhanden**, gegen [ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) Teilfrage 2 Zeile für Zeile
  gehalten: `RegisterConsumer`, `AcknowledgeConsumer`, `GetConsumerPosition`,
  `RemoveConsumer`, `EnableTable`, `DisableTable`, `GetTableStatus`,
  `ListTables`, `RunRetention` — keine erfunden, keine fehlend.
- **Rechtsklassen-Tabelle deckungsgleich mit Teilfrage 4**: 6× `admin`
  (`RegisterConsumer`, `AcknowledgeConsumer`, `RemoveConsumer`, `EnableTable`,
  `DisableTable`, `RunRetention`), 3× `reader` (`GetConsumerPosition`,
  `GetTableStatus`, `ListTables`) — exakte Übereinstimmung mit der ADR-Tabelle
  (`roleAdmin`/`roleReader`-Zuordnung).
- **Fehlercode-Tabelle deckungsgleich mit Teilfrage 5**: `InvalidArgument`,
  `Unauthenticated`, `PermissionDenied`, `NotFound`, `Internal` — dieselben
  fünf Zeilen wie die ADR-Tabelle (der strukturell entfallende
  JSON-Decode-Fall wird im Nutzerhandbuch korrekt nicht aufgeführt, da er dort
  keine beobachtbare Fähigkeit ist).
- **Erreichbarkeit/Auth korrekt beschrieben**: `CDC_GRPC_ADDR` aktiviert
  beide Dienste gemeinsam, Bearer-Token-Mechanik und beide Token-Klassen
  benannt, `Unauthenticated`/`PermissionDenied`-Verhalten beschrieben —
  deckt sich mit dem real vom Reviewer beobachteten E2E-Verhalten.
- **§5-`CDC_GRPC_ADDR`-Zeile korrigiert**: `git show 60d845d3` zeigt die
  Änderung von „Horch-Adresse des gRPC-**Streaming**-Servers“ auf
  „Horch-Adresse des gRPC-Servers … aktiviert gemeinsam den Change-Stream
  (`ChangeStream`) und die Verwaltungs-API (`Administration`)“ — behebt genau
  die im Finding benannte Unvollständigkeit.
- **Versionshistorie korrekt nachgezogen**: Kopf `Version: 1.70` → `1.71`
  (Zeile 3), neue Tabellenzeile `1.71 | 2026-09-28 | …` am Ende der
  Änderungshistorie-Tabelle angehängt (nicht eingeschoben — Reihenfolge
  bleibt chronologisch), Inhalt der Zeile deckt exakt die zwei
  vorgenommenen Änderungen (neuer Abschnitt, `CDC_GRPC_ADDR`-Zeile).
- **Anker/Struktur:** die neue Abschnittsüberschrift verlinkt intern per
  `#zugriff-über-die-http-json-api` und `#zugriff-über-den-grpc-change-stream`
  — beide Ziel-Überschriften existieren unverändert (`grep -n` bestätigt),
  `make gates`s `docs-check` (Modul `anchors`) lief mit 0 Befunden über den
  gesamten Baum.
- **Keine ADR-/Review-Verweise im benutzerseitigen Fließtext:** der neue
  Abschnittstext selbst (Prosa zwischen der Überschrift und der
  Versionshistorie-Tabelle) trägt **keine** `ADR-*`/Review-Kennung — geprüft
  durch vollständiges Lesen des eingefügten Texts. Die einzigen `ADR-0130`-/
  `LH-FA-SST-006`-Nennungen stehen in der **Versionshistorie-Zeile**
  (Changelog, kein Fließtext) — dieselbe bereits etablierte Konvention wie
  die Nachbarzeilen 1.68/1.69 derselben Tabelle (die dort ebenfalls
  `LH-FA-CFG-007`/`ADR-0112`/`ADR-0125` tragen). Kein neuer Verstoß gegen die
  Nutzerdirektive.

**Verdikt F-1:** vollständig und korrekt behoben — kein verbleibender Mangel.

## 3. F-2 (MEDIUM) — Kommentar-Kennung: geprüft, behoben

`git show 60d845d3 -- gen/cdc/administration/v1/administration_test.go`
zeigt den vollständigen Diff des Kopf-Kommentars:

```diff
-// Service an den Stellen, die eine Zusage des Draht-Vertrags tragen
-// (`SPEC-031`, `ADR-0130`): die Feld-Getter, die Fehler-Weitergabe des
+// Service an den Stellen, die eine Zusage des Draht-Vertrags (`ADR-0130`)
+// tragen: die Feld-Getter, die Fehler-Weitergabe des
```

Der Kopf-Kommentar-Block (Zeilen 1–14) trägt jetzt **genau eine** Kennung
(`ADR-0130`) — von Hand gezählt, nicht nur über das (für `gen/` strukturell
blinde) Werkzeug. Ein zweiter, unabhängiger Kommentarblock weiter unten in
derselben Datei (`TestAdministrationClientReichtInvokeErgebnisJeRPCWeiter`,
Zeile ~255) trägt ebenfalls genau eine Kennung (`ADR-0130`) — kein Verstoß,
da §3.7 „höchstens eine Kennung je **Block**“ verlangt, nicht je Datei. Beide
Vorkommen von `ADR-0130` in der Datei liegen in getrennten Blöcken (`grep -n
"ADR-0130" gen/cdc/administration/v1/administration_test.go` → Zeilen 2 und
255, mit weit über 200 Zeilen Abstand und einem eigenen Testfall-Kommentar
dazwischen).

**Verdikt F-2:** vollständig und korrekt behoben — kein verbleibender Mangel.

## 4. F-3 (INFO) — folgenlos, keine Aktion nötig

Der Reviewer maß 80,70 % Coverage, mein eigener Lauf nach der Fixrunde misst
80,80 % (identisch mit der ursprünglichen Implementer-Angabe) — beide Werte
liegen über der 80-%-Schwelle, keine Gate-Auswirkung. Die geringe
Schwankung ist mit dem `gen/`-Testnachzug bzw. normaler Lauf-zu-Lauf-Varianz
der Coverage-Messung vereinbar und ohne DoD-Relevanz; der Reviewer hatte dies
bereits korrekt als „folgenlos“ eingeordnet.

## 5. ADR-0130-Konformität — Fitness-Function-Tabelle

| Zeile der ADR-Tabelle | Status |
|---|---|
| Go-Unit-Test `authUnaryInterceptor` (Unauthenticated/PermissionDenied/admin-erreicht-beides) | vom Reviewer geprüft, ohne Befund (`interceptor_test.go`), `make test` in meinem `make gates`-Lauf grün mitgelaufen |
| Go-Unit-Test je Handler (ruft exakt den Inbound Use Case seines HTTP-Äquivalents) | vom Reviewer Zeile für Zeile gegen die vier HTTP-Handler-Dateien gehalten, ohne Befund |
| `.a-check` unverändert | eigener Lauf in `make gates`: `a-check: gesamt: 0 Befund(e)` |
| `make test-integration` — realer gRPC-Rundlauf je Token-Klasse + zwei Negativpfade | vom Reviewer selbst frisch bis zum Ende durchlaufen lassen (siehe §1), Beleg über `cdc.consumer`/`cdc.source_table` unabhängig bestätigt |

Alle vier Fitness-Function-Zeilen sind erprobt (nicht nur behauptet) —
keine Lücke gegenüber der ADR.

**Kompatibilität (ADR §Kompatibilität):** `proto/cdc/stream/v1/changestream.proto`
und `gen/cdc/stream/v1/*` sind über den gesamten Diff `b317d529..HEAD`
unverändert — bereits vom Reviewer per `git diff --stat` bestätigt, durch
meinen `make generated-sync`-Lauf erneut bestätigt (beide Quellen, alle vier
generierten Dateien byte-gleich zum gepinnten Generator).

**Folgepflichten der ADR, die bewusst außerhalb dieses Zugs bleiben** (kein
DoD-Bruch, ADR nennt sie ausdrücklich als „Gegenstand eines eigenen
Folge-Schritts/Slice“, nicht dieser ADR): `spec/pflichtenheft.md`
`SPEC-*`-Eintrag ([SPEC-031](../../spec/pflichtenheft.md) bereits vorbestehend, laut Review außerhalb des
Diffs), Beispiel-Clients (Go/C#/Kotlin) und SDK-Erweiterung um die neun
Administration-RPCs. Die Fixrunde benennt diese Lücke jetzt zusätzlich
**explizit im Benutzerhandbuch** selbst („Für keine der drei Sprachen … noch
kein `examples/`-Programm und kein SDK-Package“) — eine Betreiber-sichtbare,
ehrliche Aussage über den aktuellen Stand, kein Aufschub-Ersatz für eine
künftige Slice-Adresse, aber auch keine falsche Vollständigkeitsbehauptung.

## 6. Plan-vs-Code-Diff

`git diff --stat b317d529..HEAD` — 19 Dateien, deckungsgleich mit der ADR-
Entscheidung (neues `.proto`-Paket, generierter Code, Handler, Interceptor,
Server-Verdrahtung, Wiring, E2E-Test, Harness-Doku-Nachzug, Review- und
Handbuch-Fix) — kein unbenannter Nebeneffekt; `changestream.proto`/
`gen/cdc/stream/v1/*` fehlen im Diff wie von der ADR gefordert.

## 7. Verdikt

**DoD erfüllt: ja.**

- Beide merge-blockierenden bzw. zu behebenden Findings des Code-Reviews
  (F-1 HIGH, F-2 MEDIUM) sind durch Commit `60d845d3` **tatsächlich und
  vollständig** behoben — nicht nur behauptet, sondern in diesem Lauf
  eigenständig am Diff nachgeprüft (§2, §3).
- F-3 (INFO) bleibt folgenlos.
- Alle vier Fitness-Function-Zeilen der ADR sind erprobt, nicht nur
  hergeleitet (§5).
- `make gates` läuft in meinem eigenen, frischen, ungepipten Lauf mit
  Exit 0 (§1) — inklusive `generated-sync` für beide `.proto`-Quellen,
  `coverage-gate` über der Schwelle, `commit-traceability` für die letzten
  fünf Commits.
- Der neue Handbuch-Abschnitt selbst verstößt an keiner Stelle gegen
  `AGENTS.md`-Hard-Rules — insbesondere keine ADR-/Review-Verweise im
  benutzerseitigen Fließtext (nur in der bereits etablierten
  Changelog-Konvention).
- Kein unbenannter Diff-Nebeneffekt; `ChangeStream`/`changestream.proto`
  bleiben byte-identisch unverändert.

**Verbleibendes Restrisiko (kein DoD-Bruch, zur Kenntnis):** Die
ADR-Folgepflichten „Beispiel-Clients“ und „SDK-Erweiterung“ um die neun
Administration-RPCs sind noch offen und tragen keine committete
Folge-Slice-Adresse — die ADR selbst benennt sie ausdrücklich als
eigenständige Folge-Schritte außerhalb dieses Zugs, und das Benutzerhandbuch
macht die Lücke jetzt betreiberseitig sichtbar. Da weder Auftrag noch ADR
noch Review dies als Teil der DoD dieses Zugs führen, ist es kein
Verifikations-Mangel dieses Reports, sondern ein Hinweis für die nächste
Planungsrunde.

**Gates:** `make gates` — Exit 0 (eigener Lauf, ungepiped, `AGENTS.md` §3.9).
`make kommentar-kennungen DIFF=b317d529` — Exit 0.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt
keine künftige Verifikation an einem späteren Stand.
