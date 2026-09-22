# Welle welle-sdk-kotlin-vollabdeckung — Closure-Notiz

**Welle:** welle-sdk-kotlin-vollabdeckung
**Abschluss:** 2026-09-22
**Verantwortlich:** pt9912

## Was wurde geliefert?

[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
§Entscheidung Festlegung 1 letzter Absatz eingelöst — `pgchangefeed-kotlin`
erreicht die volle Vier-Wege-Parität, die
`examples/kotlin/{http,grpc,sse,nats-stream}-client` bereits real belegen,
in zwei Slices:

- **`slice-sdk-kotlin-sse-client-flaeche`**: öffentliche SSE-Client-Fläche
  (`SPEC-021`) — `PgChangeFeedSseClient`, eigener Frame-Parser für
  `text/event-stream`, alle zehn `SPEC-021`-Nachrichtenfelder.
- **`slice-sdk-kotlin-nats-stream-client-flaeche`**: öffentliche
  NATS-Vollinhalts-Client-Fläche (`SPEC-024`) — `PgChangeFeedNatsStreamClient`,
  Subjekt-Formatierer (`buildSubject`/`buildSourceSubject`), alle zehn
  `SPEC-024`-Nachrichtenfelder; bündelt zusätzlich die Version-Hebung
  (`0.1.0` → `0.2.0`, additive, rückwärtskompatible Erweiterung) und den
  Träger-Nachzug (`AGENTS.md` §3.13) in `spec/pflichtenheft.md`
  (`LH-FA-SST-009.a`, `SPEC-028`) und `docs/user/benutzerhandbuch.md` für
  **beide** neu gelieferten Flächen.

Ein real neu gebautes `sdks/kotlin/dist/pgchangefeed-kotlin-0.2.0.jar`
(144598 Bytes, eigenständiger, frischer `bash tools/harness/sdk-pack-kotlin.sh`-Lauf
zu dieser Closure, nach `rm -f sdks/kotlin/dist/*.jar`) trägt alle vier
Client-Flächen (`http`/`grpc`/`sse`/`nats`-Paketpfade) in einem einzigen
Jar — `LH-FA-SST-009` gilt für Kotlin damit als vollständig für die
Vier-Wege-Matrix erfüllt, kein isolierter Flächen-Slice allein hätte das
belegt (Welle-Ziel §1).

`make gates` grün auf dem Endstand (906 Dateien, 0 `docs-check`-Befunde,
Coverage 82,80 % ≥ 80 %, `a-check`/`generated-sync`/`commit-traceability`/
`baseline-verify` je ohne Befund). Kein realer `sdk-kotlin-v0.2.0`-Tag-Push
(`git tag -l "sdk-kotlin-v*"` liefert **keinen** Treffer — auch keinen
`sdk-kotlin-v0.1.0` aus einem früheren Slice) — diese Welle liefert bewusst
nur das paketierbare Ergebnis, wie im Closure-Trigger (§3 der Welle-Datei)
vorab festgelegt.

## Was hat funktioniert?

- Beide Slices liefen mit vollständiger, mehrfach unabhängiger
  Verifikation: `slice-sdk-kotlin-sse-client-flaeche` glatt durch (0
  HIGH/MEDIUM, 2 LOW, 2 INFO, keine Fixrunde, Verifier bestätigt „DoD
  erfüllt"), `slice-sdk-kotlin-nats-stream-client-flaeche` mit einer
  gezielten Fixrunde (1 HIGH + 1 MEDIUM → real behoben → 0
  HIGH/MEDIUM/LOW/INFO im Fixrunden-Review → Verifier bestätigt „DoD
  erfüllt", fünfter unabhängiger Docker-Bau dieser Welle) — je mit
  eigenständig nachgefahrenen, nicht übernommenen Docker-Build-Läufen
  (`git worktree` gegen den jeweiligen Elternstand, isolierter,
  ungecachter Bau) für jede behauptete Testzahl.
- Das bereits bei den ersten beiden Kotlin-Flächen (HTTP/gRPC) bewährte
  `SseTransport`/`GrpcStreamTransport`-Seam-Muster hat sich unverändert auf
  NATS übertragen lassen (`NatsStreamTransport`, ein
  `() -> ByteArray?`-Next-Payload-Supplier) — netzlose Tests ohne echten
  NATS-Server, exakt wie bei den beiden vorigen Flächen. Die C#-Geschwister-
  Fläche (`sdks/csharp/PgChangeFeed.Client/Nats/PgChangeFeedNatsStreamClient.cs`)
  war ein direkt übertragbares Formvorbild für Konstruktoren, Subjekt-
  Validierung und die Entscheidung, keine zweite Fehlerklassen-Hierarchie
  für Verbindungsfehler zu erfinden.
- Vier unabhängige Messungen derselben Testzahl (Implementer, Reviewer,
  Fixrunden-Reviewer, Verifier) beim NATS-Slice konvergierten exakt auf
  dasselbe Ergebnis (Delta 14 Testfälle: 4 Auth + 3 Schema + 7 Subjekt) —
  wiederholtes, unabhängiges Nachmessen fängt Zahlen-Drift zuverlässig,
  bevor sie in `done/` landet.
- Das etablierte 3-Commit-Move-Muster (`git mv` · Inhalt · `git mv`) hielt
  `AGENTS.md` §3.3 über beide Slices sauber; Backtick-Paritäts-Checks vor
  jedem Commit liefen ohne Nachbesserung.

## Was ging anders als geplant?

- Der Träger-Nachzug (`AGENTS.md` §3.13) im NATS-Slice betraf real zwei
  Stellen außerhalb des im DoD-Wortlaut vorgeschriebenen `grep`-Musters
  (`harness/README.md`s `make sdk-pack-kotlin`-Zeile,
  `harness/mk/sdk.mk`s gleichlautender Kommentarblock — beide seit der
  Kotlin-SSE-Fläche bereits falsch, aber übersehen) — der Implementer-
  eigene §3.13-Suchlauf fand und behob beide im selben Zug. Zusätzlich
  übersah dieser Suchlauf zwei Kommentarzeilen in `sdks/kotlin/Dockerfile`
  (Zeile 77/104, stale `pgchangefeed-kotlin-0.1.0.jar`) — diese Stelle fand
  erst der unabhängige Reviewer (F-2, MEDIUM), nicht der
  Implementer-Suchlauf: dieselbe Unter-Klasse „gefunden vom Reviewer,
  nicht vom Implementer-Suchlauf" wie bei mehreren früheren Belegen dieser
  Beobachtung. In einer gezielten Fixrunde real behoben, von einem
  frischen Fixrunden-Review und dem Verifier je eigenständig bestätigt.
- Derselbe NATS-Review fand einen zweiten, unabhängigen Fund derselben
  Klasse wie beim analogen C#-NATS-Sibling: `build.gradle.kts`s
  Version-Kommentar erzählte Vorher/Nachher (Arrow-Muster `0.1.0 -> 0.2.0`,
  Slice-Kennung statt `ADR-*`/`LH-*`-Anker) statt eine reine Zustandsaussage
  zu treffen (F-1, HIGH, „Slice-/Wellen-Chronik in Produktionscode-
  Kommentar") — real in der Fixrunde behoben, vierfach unabhängig
  bestätigt (Reviewer, Fixrunden-Reviewer, Verifier, Planner-Closure).
  Beide Findings dieser Runde sind bei der Planner-Closure als weitere
  Belege zweier bereits verkörperter Beobachtungen nachgetragen worden
  (`BEO-PGC/slice-chronik-in-code-kommentar` neunte, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
  einundzwanzigste Evidenzdatei) — keine neue Klasse, kein neuer
  Registereintrag.
- `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` erreichte
  bereits während der Geschwister-Welle `welle-sdk-csharp-vollabdeckung`
  die 3×-Schwelle; kein Kotlin-Slice dieser Welle erzeugte einen vierten
  Beleg — beide Kotlin-Flächen-Slices benannten `--build-context
  proto=proto`s Zwangsbreite explizit im eigenen Plan (§1/§3), statt sie
  als überraschenden Fund im Review auftreten zu lassen. Siehe
  „Steering-Loop-Einträge" unten für den (wiederholten) Lese-Schritt.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier stehen die Beobachtungen, die im
Register 3× erreicht haben (Lese-Schritt dieser Closure).

- **`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`** — bereits
  während `welle-sdk-csharp-vollabdeckung` real auf die 3×-Schwelle
  gestiegen (`evidence/slice-102.md`, `evidence/slice-103.md`,
  `evidence/slice-sdk-csharp-sse-client-flaeche.md`); jener Lese-Schritt
  ließ `state.md` bewusst `offen` (die Regelschärfungs-Frage ist eine
  Architect-Entscheidung, kein Planner-Zug) und benannte „eine künftige
  Architect-Sichtung" als nächste Gelegenheit. Diese Welle erzeugte
  **keinen** vierten Beleg — beide Kotlin-Flächen-Slices (§6 der jeweiligen
  Slice-Pläne, §1/§3 des NATS-Slice-Plans) benannten die
  `--build-context proto=proto`-Zwangsbreite von
  `sdks/kotlin/Dockerfile`s einziger `build`-Stufe von vornherein
  explizit, statt einen engen DoD-Satz zu schreiben, den die Kopplung dann
  im Nachhinein widerlegt. Zähler bleibt bei **3×**, `state.md` bleibt
  `offen` mit unverändertem Vermerk — kein weiterer Lese-Schritt nötig, die
  Architect-Entscheidung steht weiterhin aus.
- **`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`** — erreichte
  bereits während des ersten Flächen-Slices dieser Welle
  (`slice-sdk-kotlin-sse-client-flaeche`) real die 6×-Schwelle und ist
  seit `slice-089` als `AGENTS.md` §3.12 Instanz B **verkörpert** (Leser:
  Verifier und Planner). `slice-sdk-kotlin-nats-stream-client-flaeche`
  erzeugte **keinen** weiteren Beleg dieser Klasse — kein DoD-/Plan-Satz
  dieses Slices behauptete eine bei Niederschrift bereits unzutreffende
  Tatsache. Zähler bleibt bei **6×**, Ausgang bleibt **verkörpert**, kein
  neuer Handlungsbedarf.
- **`BEO-PGC/slice-chronik-in-code-kommentar`** (bereits verkörpert seit
  einem wellenlosen Architect-Zug): ein neuer Beleg in dieser Welle
  (`evidence/slice-sdk-kotlin-nats-stream-client-flaeche.md`, neunte
  Evidenzdatei) — Reviewer fing die Stelle vor Merge, bestätigt die
  bestehende Diagnose erneut, kein neuer Regelschärfungs-Anlass.
- **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** (verkörpert seit
  `AGENTS.md` §3.13): ein neuer Beleg in dieser Welle
  (`evidence/slice-sdk-kotlin-nats-stream-client-flaeche.md`,
  einundzwanzigste Evidenzdatei) — mit einer Variante gegenüber dem
  unmittelbar vorangehenden, zwanzigsten (C#-)Beleg: dort fand der
  Implementer-eigene Suchlauf die Fundstelle selbst, hier fand sie erst
  der Reviewer. Bereits verkörpert, keine erneute Verkörperung nötig.
- **`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`**
  (bereits verkörpert): geprüft, kein neuer Handlungsbedarf — beide
  Flächen-Slices trugen den Handbuch-Hinweis bereits als eigenen DoD-Punkt
  (SSE + NATS-Vollinhalt).
- **`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`** (bereits 3×,
  verkörpert/gelesen bei `welle-sdk-kotlin-lh-fa-sst-009`): geprüft, kein
  neuer Beleg in dieser Welle — Zähler bleibt unverändert bei 3×.

**Unter der Schwelle, unverändert, kein Verkörperungsbedarf:**

- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` —
  weiterhin 1×, kein zweites Auftreten (Publish-Mechanismus unverändert in
  dieser Welle, kein realer Tag-Push).
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als
  `AGENTS.md` §3.10) — unverändert, kein neuer Beleg, da kein realer
  Tag-Push stattfand.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/). Kein Eintrag
erreichte in dieser Welle **neu** die 3×-Schwelle — der Lese-Schritt oben
bestätigt für `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`
lediglich den bereits während der Geschwister-Welle erreichten Stand
(unverändert `offen`, Architect-Entscheidung steht aus) und für
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` den bereits
beim SSE-Slice dieser Welle erreichten, bereits verkörperten Stand.

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine belegte
Feststellung.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` enthält
  ausschließlich `.gitkeep`, kein Carveout referenziert `ADR-0109` oder
  eine der beiden Slices dieser Welle.
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate`
  steht unverändert bei 80 %-Endstufen-Schwelle (real gemessen 82,80 %)
  und wurde von keinem der beiden Slices berührt (reiner Kotlin-Baum,
  außerhalb der netzlos prüfbaren Go-Fläche `./internal/...`+`./cmd/...`+
  `./gen/...`); `.a-check.yml` `languages: go` liest `sdks/kotlin/**`
  strukturell nicht (real bestätigt: `a-check` meldet 0 Befunde auf dem
  Endstand).
- **Entscheidung/ADR —
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  §Re-Evaluierungs-Trigger, durch diese Welle **nicht** ausgelöst:**
  - Trigger 1 (vierte Sprache/vierter Vertriebsweg verlangt) — nicht
    eingetreten, out-of-scope §6 der Welle-Datei.
  - Trigger 2 (real auf GitHub Packages veröffentlicht + Nutzungsdaten
    legen andere Priorisierung nahe) — **nicht** eingetreten: Diese Welle
    liefert SSE/NATS-Vollinhalt real als v2 desselben Packages, aber genau
    das war bereits durch `ADR-0109` §Entscheidung Festlegung 1 letzter
    Absatz ausdrücklich **ohne neue ADR** erlaubt — die Einlösung einer
    bereits vorgesehenen Delegation, keine neue, re-evaluierungspflichtige
    Entscheidung. `git tag -l "sdk-kotlin-v*"` liefert **keinen** Treffer —
    kein realer `0.2.0`-Veröffentlichungsversuch, also auch keine
    Nutzungsdaten, die Trigger 2 auslösen könnten.
  - Trigger 3 (Gradle wird real auf ≥ 9.1.0 angehoben) — nicht
    eingetreten; diese Welle blieb bei `eclipse-temurin:21-jdk`/Gradle
    `8.14`, unverändert (out-of-scope §6 der Welle-Datei).
  - Trigger 4 (der reale GitHub-Packages-Publish-Workflow scheitert im
    ersten realen Lauf) — nicht einschlägig: kein realer Tag-Push in
    dieser Welle, `.github/workflows/sdk-kotlin-release.yml` wurde von
    keinem der beiden Slices verändert.
  - Trigger 5 (Authentifizierungspflicht beim Lesen erweist sich als
    echtes Consumer-Problem) — nicht eingetreten, keine Nachfrage
    aufgetreten.

  `ADR-0109` bleibt `Accepted` und inhaltlich unberührt durch diese Welle.

## Folge-Slices

Keine. `ADR-0109` ist mit dieser Welle für die volle Kotlin-Vier-Wege-Matrix
vollständig umgesetzt. Was ansteht, sind keine weiteren Slices dieser
Welle, sondern eigenständige, künftige Entscheidungen:

- Ein realer `sdk-kotlin-v0.2.0`-Tag-Push — eine irreversible, extern
  sichtbare Aktion (öffentlicher GitHub-Packages-Push), die nur nach
  expliziter, gesonderter Rückfrage beim Auftraggeber läuft (`AGENTS.md`
  §3.10, wie bei den vorigen SDK-Wellen).
- Eine vierte SDK-Sprache oder ein vierter Vertriebsweg — bleibt laut
  `ADR-0109` Re-Evaluierungs-Trigger 1 offen für eine künftige, separate
  ADR.
- Eine künftige Architect-Sichtung von
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (weiterhin 3×,
  siehe Steering-Loop-Einträge oben) — kein Slice dieser Welle, eine
  eigenständige Regelschärfungs-Entscheidung, unverändert von der
  Geschwister-Welle übernommen.

## Verifikation

- `docs/reviews/review-slice-sdk-kotlin-sse-client-flaeche.md` (0
  HIGH/MEDIUM, 2 LOW, 2 INFO, keine Fixrunde) +
  `docs/reviews/verifikation-slice-sdk-kotlin-sse-client-flaeche.md` (DoD
  erfüllt).
- `docs/reviews/review-slice-sdk-kotlin-nats-stream-client-flaeche.md` (1
  HIGH, 1 MEDIUM) + Fixrunde (Commit `6b279311`) +
  `docs/reviews/review-slice-sdk-kotlin-nats-stream-client-flaeche-fixrunde.md`
  (0 HIGH/MEDIUM/LOW/INFO) +
  `docs/reviews/verifikation-slice-sdk-kotlin-nats-stream-client-flaeche.md`
  (DoD erfüllt, fünfter unabhängiger Docker-Bau dieser Welle).
- `make gates`: grün auf dem Endstand (906 Dateien, 0 Befunde; Coverage
  82,80 % ≥ 80 %; `a-check`/`generated-sync`/`commit-traceability`/
  `baseline-verify` je ohne Befund), ungepiped geprüft (`AGENTS.md` §3.9)
  nach jedem Commit dieser Closure sowie erneut vor dieser Welle-Closure.
- Reales `pgchangefeed-kotlin-0.2.0.jar`-Artefakt — frisch für diese
  Wellen-Closure eigenständig neu gebaut (`bash
  tools/harness/sdk-pack-kotlin.sh`, Exit 0, nach `rm -f
  sdks/kotlin/dist/*.jar`), 144598 Bytes.
- `git tag -l "sdk-kotlin-v*"`: liefert **keinen** Treffer — bestätigt,
  dass kein realer `0.2.0`-Veröffentlichungsversuch stattfand (auch kein
  `sdk-kotlin-v0.1.0` aus einem früheren Slice; Welle-Datei §3 vierter
  Punkt eingehalten).
- `docs/plan/carveouts/`: nur `.gitkeep` — 0 offene Carveouts dieser Welle.

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente (kein
`archiv`-Ziel in `Makefile`/`harness/mk/*.mk`, real geprüft per `grep`,
dieselbe Feststellung wie bei den vorigen SDK-Wellen) — die Bedingung für
Schritt 4 der Closure-Prozedur ist nicht eingetreten, keine Handarbeit als
Ersatz.
