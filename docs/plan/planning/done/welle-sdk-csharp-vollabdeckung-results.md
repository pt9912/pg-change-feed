# Welle welle-sdk-csharp-vollabdeckung — Closure-Notiz

**Welle:** welle-sdk-csharp-vollabdeckung
**Abschluss:** 2026-09-22
**Verantwortlich:** pt9912

## Was wurde geliefert?

[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md) §Entscheidung
Festlegung 1 letzter Absatz eingelöst — `PgChangeFeed.Client` erreicht die
volle Vier-Wege-Parität, die `examples/csharp/{http,grpc,sse,nats-stream}-client`
bereits real belegen, in zwei Slices:

- **`slice-sdk-csharp-sse-client-flaeche`**: öffentliche SSE-Client-Fläche
  (`SPEC-021`) — `PgChangeFeedSseClient`, eigener Frame-Parser für
  `text/event-stream`, alle zehn `SPEC-021`-Nachrichtenfelder.
- **`slice-sdk-csharp-nats-stream-client-flaeche`**: öffentliche
  NATS-Vollinhalts-Client-Fläche (`SPEC-024`) — `PgChangeFeedNatsStreamClient`,
  Subjekt-Formatierer (`BuildSubject`/`BuildSourceSubject`), alle zehn
  `SPEC-024`-Nachrichtenfelder; bündelt zusätzlich die Version-Hebung
  (`0.1.0` → `0.2.0`, additive, rückwärtskompatible Erweiterung) und den
  Träger-Nachzug (`AGENTS.md` §3.13) in `spec/pflichtenheft.md`
  (`LH-FA-SST-009.a`, `SPEC-026`) und `docs/user/benutzerhandbuch.md` für
  **beide** neu gelieferten Flächen.

Ein real neu gebautes `sdks/csharp/dist/PgChangeFeed.Client.0.2.0.nupkg`
(37451 Bytes, eigenständiger, frischer `bash tools/harness/sdk-pack-csharp.sh`-Lauf
zu dieser Closure) trägt alle vier Client-Flächen (`Http`, `Grpc`, `Sse`,
`Nats`-Namensräume) in einer einzigen Assembly — `LH-FA-SST-009` gilt für
C# damit als vollständig für die Vier-Wege-Matrix erfüllt, kein isolierter
Flächen-Slice allein hätte das belegt (Welle-Ziel §1).

`make gates` grün auf dem Endstand (897 Dateien, 0 `docs-check`-Befunde,
Coverage 82,80 % ≥ 80 %, `a-check`/`generated-sync`/`commit-traceability`/
`baseline-verify` je ohne Befund). Kein realer `sdk-csharp-v0.2.0`-Tag-Push
(`git tag -l "sdk-csharp-v*"` zeigt weiterhin nur `sdk-csharp-v0.1.0`) —
diese Welle liefert bewusst nur das paketierbare Ergebnis, wie im
Closure-Trigger (§3 der Welle-Datei) vorab festgelegt.

## Was hat funktioniert?

- Beide Slices liefen mit vollständiger, mehrfach unabhängiger
  Verifikation: `slice-sdk-csharp-sse-client-flaeche` glatt durch (0
  HIGH/MEDIUM, 3 LOW, 1 INFO, keine Fixrunde, Verifier bestätigt „DoD
  erfüllt"), `slice-sdk-csharp-nats-stream-client-flaeche` mit einer
  gezielten Fixrunde (2 HIGH + 1 MEDIUM → real behoben → 0 HIGH/MEDIUM/LOW,
  1 INFO → Verifier bestätigt „DoD-konform: ja") — je mit eigenständig
  nachgefahrenen, nicht übernommenen Docker-Build-Läufen (`git worktree`
  gegen den jeweiligen Elternstand, isolierter, ungecachter Bau) für jede
  behauptete Testzahl.
- Die Design-Lehre aus der SSE-Fixrunde trug real in den Folge-Slice: Der
  SSE-Review benannte (F-2) eine wörtliche Kopie des HTTP-Fehler-Mappings
  als Wartbarkeitsrisiko und warnte vorab, ein zweites Auftreten desselben
  Zuschnitts im Folge-Slice wäre ein Kandidat für einen neuen
  Beobachtungs-Eintrag. Der NATS-Slice vermied das bewusst — eine
  eigenständige `PgChangeFeedNatsMalformedMessageException` statt einer
  Wiederverwendung der HTTP-Exception-Hierarchie, mit einer echten
  strukturellen Begründung (kein `StatusCode`-Äquivalent bei NATS) statt
  einer bloßen Kopie — vom NATS-Reviewer eigenständig gegen die
  SSE-Vorgeschichte gehalten und als begründete, konsistente
  Design-Abweichung bestätigt. Kein zweites Auftreten dieser Klasse.
- Das etablierte 3-Commit-Move-Muster (`git mv` · Inhalt · `git mv`) hielt
  `AGENTS.md` §3.3 über beide Slices sauber; Backtick-Paritäts-Checks vor
  jedem Commit liefen ohne Nachbesserung.
- Vier unabhängige Messungen derselben Testzahl (Implementer, Reviewer,
  Fixrunden-Reviewer, Verifier) beim NATS-Slice konvergierten exakt auf
  dasselbe Ergebnis (Delta 21 Testfälle) — ein reales Beispiel dafür, dass
  wiederholtes, unabhängiges Nachmessen (statt Übernahme aus dem
  Bericht der vorigen Rolle) Zahlen-Drift zuverlässig fängt, bevor er in
  `done/` landet.

## Was ging anders als geplant?

- Der NATS-Slice-Review fand zwei reale Zahlen-Drifts in der DoD-Zeile und
  im Beobachtungs-Beleg (F-2 „sechzehn" statt real 21 Testfälle, F-3
  „neun" statt real zehn Dateien) — beide in einer gezielten Fixrunde
  behoben, von einem frischen Fixrunden-Review eigenständig nachgemessen
  bestätigt (kein Rollen-Konflikt, keine Konflikt-Pfad-Eskalation nötig,
  Modul 8). Eine dritte, verwandte Unschärfe („acht Fundstellen", keine im
  Text definierte Zähleinheit) blieb als nicht-blockierendes INFO bzw.
  als Verifier-Empfehlung offen — der Planner hat sie bei der
  Slice-Closure umgesetzt: die Formulierung streicht die unscharfe
  Summenzahl ersatzlos und trägt ausschließlich die bereits dreifach
  unabhängig verifizierte Zahl „zehn Dateien"
  (`docs/plan/planning/done/slice-sdk-csharp-nats-stream-client-flaeche.md`
  §7, `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md`,
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger/state.md`).
- `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` erreichte real
  während dieser Welle die 3×-Schwelle (`slice-sdk-csharp-sse-client-flaeche`
  als dritter Beleg, nach zwei `examples/**`-Vorgängen in C#/Kotlin) — der
  NATS-Slice erzeugte bewusst keinen vierten Beleg. Siehe
  „Steering-Loop-Einträge" unten für den Lese-Schritt.
- Der Träger-Nachzug (`AGENTS.md` §3.13) im NATS-Slice betraf real zwei
  Dateien außerhalb des im DoD-Wortlaut vorgeschriebenen `grep`-Musters
  (`tools/harness/sdk-pack-csharp.sh`, `Sse/Models/Change.cs`s
  „three surfaces"-Kommentar) — gefunden durch aufmerksames Lesen
  benachbarter Dateien, nicht durch den Suchlauf-Befehl selbst; derselbe
  Grenzfall, den `AGENTS.md` §3.13 §Grenze bereits für Zahlen/
  Prosa-Umformulierungen benennt, hier erstmals für einen Kommentar-Text
  belegt (zwanzigster Gesamtbeleg der Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, weiterhin verkörpert seit
  `AGENTS.md` §3.13, kein neuer Verkörperungsbedarf).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier stehen die Beobachtungen, die im
Register 3× erreicht haben (Lese-Schritt dieser Closure).

- **`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`** — erreichte
  während dieser Welle real die 3×-Schwelle
  (`evidence/slice-102.md`, `evidence/slice-103.md`,
  `evidence/slice-sdk-csharp-sse-client-flaeche.md`). Die ersten zwei
  Vorgänge treffen dieselbe strukturelle Ursache in `examples/**` (eine
  gemeinsame `build`-Stufe koppelt den benannten `--build-context
  proto=proto`-Zusatzkontext an alle Programme einer Sprache statt nur an
  `grpc`); der dritte Vorgang trifft denselben Symptom-Kern in einem
  SDK-Package-Baum (`sdks/csharp/**`), aber mit einer strukturell anderen,
  bewussten Ursache (ein einziges `.csproj`/Package statt vermeidbar
  geteilter Docker-Stufen über eigenständige Programme). Sein eigenes
  `state.md` benannte die Closure dieser Welle ausdrücklich als „nächste
  reguläre Lesegelegenheit" — **diese Lesung findet hiermit statt.** Die
  eigentliche Regelschärfung (ob/wie `.harness/skills/reviewer.md` oder
  ein Dockerfile-Struktur-Muster — separate Bau-Stufen je Zielfläche statt
  einer gemeinsamen — die Klasse künftig fängt, und ob die dritte,
  ursachenmäßig andersartige Instanz dieselbe Regel oder eine eigene
  braucht) bleibt bewusst **außerhalb** dieser Planner-geführten
  Wellen-Closure: Es ist eine Architect-Entscheidung (Modul 4/8), kein
  Planner-Zug. `state.md` bleibt entsprechend `offen`, mit dem Vermerk,
  dass die Lesung stattgefunden hat und der nächste Schritt bei einer
  künftigen Architect-Sichtung liegt.
- **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** (verkörpert seit
  `AGENTS.md` §3.13): zwei neue Belege in dieser Welle
  (`evidence/slice-sdk-csharp-sse-client-flaeche.md`,
  `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md`, Zähler jetzt
  real 20×) — bereits verkörpert, keine erneute Verkörperung nötig. Der
  zwanzigste Beleg trägt eine Erweiterung des bereits benannten Grenzfalls
  (§Grenze: Prosa-/Zahlen-Umformulierungen trifft `grep` nicht
  zuverlässig) — hier erstmals für einen Kommentar-Text statt für eine
  Zahl oder einen Symbolnamen; kein neuer Regelbedarf, die bestehende
  Formulierung deckt den Fall bereits.
- **`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`**
  (bereits verkörpert): geprüft, kein neuer Handlungsbedarf — beide
  Flächen-Slices trugen den Handbuch-Hinweis bereits als eigenen DoD-Punkt
  (SSE + NATS-Vollinhalt, `docs/user/benutzerhandbuch.md` 1.36 → 1.38).
- **`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`** (bereits 3×,
  verkörpert/gelesen bei `welle-sdk-kotlin-lh-fa-sst-009`): geprüft, kein
  neuer Beleg in dieser Welle — Zähler bleibt unverändert bei 3×, keine
  weitere Lesung nötig.

**Unter der Schwelle, unverändert, kein Verkörperungsbedarf:**

- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` —
  weiterhin 1×, kein zweites Auftreten (Publish-Mechanismus unverändert in
  dieser Welle, kein realer Tag-Push).
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als
  `AGENTS.md` §3.10) — unverändert, kein neuer Beleg, da kein realer
  Tag-Push stattfand.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/). Ein Eintrag
erreichte in dieser Welle die 3×-Schwelle
(`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`) — der
Lese-Schritt dazu steht oben unter *Steering-Loop-Einträge*; die
eigentliche Regelschärfung bleibt eine offene Übergabe an eine künftige
Architect-Sichtung, kein Teil dieser Wellen-Closure.

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine belegte
Feststellung.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` enthält
  ausschließlich `.gitkeep`, kein Carveout referenziert `ADR-0106` oder
  eine der beiden Slices dieser Welle.
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate`
  steht unverändert bei 80 %-Endstufen-Schwelle (real gemessen 82,80 %)
  und wurde von keinem der beiden Slices berührt (reiner C#-Baum,
  außerhalb der netzlos prüfbaren Go-Fläche `./internal/...`+`./cmd/...`+
  `./gen/...`); `.a-check.yml` `languages: go` liest `sdks/csharp/**`
  strukturell nicht (real bestätigt: `a-check` meldet 0 Befunde auf dem
  Endstand).
- **Entscheidung/ADR —
  [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
  §Re-Evaluierungs-Trigger, durch diese Welle **nicht** ausgelöst:**
  - Trigger 1 (zweite Sprache/zweiter Vertriebsweg verlangt) — nicht
    eingetreten, out-of-scope §6 der Welle-Datei.
  - Trigger 2 (real auf NuGet.org veröffentlicht + Nutzungsdaten legen
    andere Priorisierung nahe) — **nicht** eingetreten: Diese Welle liefert
    SSE/NATS-Vollinhalt real als v2 desselben Packages, aber genau das war
    bereits durch `ADR-0106` §Entscheidung Festlegung 1 letzter Absatz
    ausdrücklich **ohne neue ADR** erlaubt („entweder als v2 desselben
    Packages oder als eigenständiges Package … Entscheidung eines
    Folge-Zuges, kein Vorgriff hier") — die Wahl „v2 desselben Packages"
    ist damit die Einlösung einer bereits vorgesehenen Delegation, keine
    neue, re-evaluierungspflichtige Entscheidung. `git tag -l
    "sdk-csharp-v*"` zeigt weiterhin nur `sdk-csharp-v0.1.0` — kein realer
    `0.2.0`-Veröffentlichungs-Versuch, also auch keine Nutzungsdaten, die
    Trigger 2 auslösen könnten.
  - Trigger 3 (gRPC-Abhängigkeits-Footprint erweist sich als echtes
    Consumer-Problem) — nicht eingetreten, kein Bezug zu dieser Welle
    (betrifft eine mögliche künftige Aufspaltung von HTTP/gRPC, nicht die
    hier gelieferten SSE-/NATS-Flächen).
  - Trigger 4 (`docs/user/version.md`s Server-SemVer erreicht `1.0.0`) —
    nicht eingetreten.

  `ADR-0106` bleibt `Accepted` und inhaltlich unberührt durch diese Welle.

## Nebenbefunde (außerhalb des Scopes dieser Welle, hier nur gemeldet)

- `gh secret list` zeigt zum Zeitpunkt dieser Closure `DOCKERHUB_TOKEN`,
  `DOCKERHUB_USERNAME`, `NUGET_API_KEY`, `PYPI_API_TOKEN` — `NUGET_API_KEY`
  existiert bereits seit der ersten csharp-Welle; unverändert durch diese
  Welle, hier nur zur Vollständigkeit vermerkt (ändert nichts an `ADR-0106`s
  eigenem Re-Evaluierungs-Trigger 2, der oben separat geführt wird).

## Folge-Slices

Keine. `ADR-0106` ist mit dieser Welle für die volle C#-Vier-Wege-Matrix
vollständig umgesetzt. Was ansteht, sind keine weiteren Slices dieser
Welle, sondern eigenständige, künftige Entscheidungen:

- Ein realer `sdk-csharp-v0.2.0`-Tag-Push — eine irreversible, extern
  sichtbare Aktion (öffentlicher NuGet.org-Push), die nur nach expliziter,
  gesonderter Rückfrage beim Auftraggeber läuft (`AGENTS.md` §3.10, wie
  bei der ersten csharp-Welle).
- Eine Aufspaltung in ein eigenständiges Folge-Package (statt v2 desselben
  Packages) — bleibt laut `ADR-0106` Festlegung 1 eine künftige,
  eigenständige Entscheidung, falls Nutzungsdaten dafür sprechen; diese
  Welle rollt das nicht neu auf.
- Eine vierte SDK-Sprache oder ein zweiter C#-Vertriebsweg — bleibt laut
  `ADR-0106` Re-Evaluierungs-Trigger 1 offen für eine künftige, separate
  ADR.
- Eine künftige Architect-Sichtung von
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (3×, siehe
  Steering-Loop-Einträge oben) — kein Slice dieser Welle, eine
  eigenständige Regelschärfungs-Entscheidung.

## Verifikation

- `docs/reviews/review-slice-sdk-csharp-sse-client-flaeche.md` (0
  HIGH/MEDIUM, 3 LOW, 1 INFO, keine Fixrunde) +
  `docs/reviews/verifikation-slice-sdk-csharp-sse-client-flaeche.md` (DoD
  erfüllt).
- `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md` (2
  HIGH, 1 MEDIUM) + Fixrunde (Commit `8b6cef7c`) +
  `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche-fixrunde.md`
  (0 HIGH/MEDIUM/LOW, 1 INFO) +
  `docs/reviews/verifikation-slice-sdk-csharp-nats-stream-client-flaeche.md`
  (DoD-konform: ja).
- `make gates`: grün auf dem Endstand (897 Dateien, 0 Befunde; Coverage
  82,80 % ≥ 80 %; `a-check`/`generated-sync`/`commit-traceability`/
  `baseline-verify` je ohne Befund), ungepiped geprüft (`AGENTS.md` §3.9)
  nach jedem Commit dieser Closure sowie erneut vor dieser Welle-Closure.
- Reales `PgChangeFeed.Client.0.2.0.nupkg`-Artefakt — frisch für diese
  Wellen-Closure eigenständig neu gebaut (`bash
  tools/harness/sdk-pack-csharp.sh`, Exit 0, nach `rm -f
  sdks/csharp/dist/*.nupkg`), 37451 Bytes, `unzip -l` bestätigt
  `lib/net10.0/PgChangeFeed.Client.dll`.
- `git tag -l "sdk-csharp-v*"`: zeigt ausschließlich `sdk-csharp-v0.1.0` —
  bestätigt, dass kein realer `0.2.0`-Veröffentlichungsversuch stattfand
  (Welle-Datei §3 vierter Punkt eingehalten).
- `docs/plan/carveouts/`: nur `.gitkeep` — 0 offene Carveouts dieser Welle.

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente (kein
`archiv`-Ziel in `Makefile`/`harness/mk/*.mk`, real geprüft per `grep`,
dieselbe Feststellung wie bei den vorigen SDK-Wellen) — die Bedingung für
Schritt 4 der Closure-Prozedur ist nicht eingetreten, keine Handarbeit als
Ersatz.
