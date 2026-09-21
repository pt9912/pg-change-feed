# Review-Report: slice-sdk-kotlin-publish-workflow (Fixrunde) — 2026-09-21

**Review-Art:** Code — geprüft gegen Plan + `ADR-0109` + Hard Rules
(`AGENTS.md` §3, Modul 10 §Drei Review-Arten). Frisches Review der
Fixrunde nach dem HIGH-Finding F-1 des vorangehenden Laufs — eigener
Kontext, kein Self-Review, keine ungeprüfte Übernahme der
Implementer-Einschätzung (`.harness/skills/reviewer.md`, Modul 8).

**Gegenstand:** Fixrunden-Commit `21872c3de45d7bfc0c27e617019182b688aa764d`
(`fix(sdk-kotlin): gradlew publish laeuft Docker-only statt auf dem
Runner`) gegen Elternstand `546a7503` — `git show 21872c3d` /
`git diff 546a7503..21872c3d` (Fixrunden-Diff selbst, sechs Dateien);
zusätzlich im Kontext des Gesamt-Slice-Diffs `eaf6ded5..21872c3d`
gegengelesen, um Zwischencommits (`24077306` Review-Report,
`3f754976` unabhängiger Slice-Plan) korrekt aus dem Scope
auszuschließen.

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-13 (unverändert
seit dem vorangehenden Lauf zu diesem Slice; Repo-`HEAD` zum Zeitpunkt
dieses Laufs).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-21.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-sdk-kotlin-publish-workflow.md`
  (vollständig gelesen, §1–§8, inkl. der Fixrunden-Ergänzung in §2)
- `ADR-0109` §Entscheidung Festlegung 5 (`Accepted`, wörtlich
  gegengelesen: „baut/testet/paketiert (`./gradlew build`,
  `./gradlew publish`) im gepinnten `eclipse-temurin:21-jdk`-Image")
- `docs/reviews/review-slice-sdk-kotlin-publish-workflow.md` (F-1 im
  Detail, Ausgangsbefund dieser Fixrunde)
- `AGENTS.md` §3.1 (Docker-only), §3.5 (Accepted-ADR-Immutabilität), §3.7
  (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.9 (Exit-Code-Disziplin),
  §3.10 (Post-Push-Risiko), §3.12/§3.13 (Beleg-Herkunft, Träger-Nachzug)
- `harness/conventions.md` (MR-000/MR-001/MR-002)
- Vollständig gelesen: `sdks/kotlin/Dockerfile` (Endstand nach dem Fix),
  `.github/workflows/sdk-kotlin-release.yml` (Endstand),
  `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` (geänderter
  Kommentarblock), `tools/harness/sdk-pack-kotlin.sh`,
  `harness/mk/sdk.mk`, `docs/user/releasing.md` §4 (Endstand)
- Eigener Docker-Bau/-Inspektionslauf gegen den realen Endstand (siehe
  F-Liste unten) — kein realer Registry-Push, kein Tag

---

## Findings

Kein HIGH, kein MEDIUM in diesem Lauf. Ein INFO-Punkt unten.

### F-1 — Laufzeit-Wiederverwendungs-Annahme im `publish`-Stufen-Kommentar bleibt bis zum realen Post-Push-Lauf unbewiesen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.10 (Post-Push-Risiko, bereits im Slice-Plan §6
  strukturell offen geführt), Nachbarklasse zu §3.12 Instanz B
  (Tatsachenbehauptung ohne Beleg-Anker)
- `pfad`: `sdks/kotlin/Dockerfile:106-118` (Kommentar der `publish`-Stufe:
  „Zur Laufzeit findet Gradles Up-to-date-Prüfung für
  `generateProto`/`compileKotlin`/`jar` den seit dem Bau unveränderten
  Dateizustand bereits erfüllt vor und lädt nur den bereits gebauten Jar
  hoch — kein zweiter Compile-/Codegenerierungs-Lauf außerhalb von
  Docker")
- `befund`: Die Behauptung, dass Gradles Up-to-date-Prüfung beim
  `docker run … ./gradlew publish`-Aufruf zur Laufzeit tatsächlich keinen
  zweiten Compile-/Codegenerierungs-Lauf auslöst, ist technisch plausibel
  (Docker-Image-Layer sind ein eingefrorener Dateizustand, Gradles
  Incremental-Build-Modell prüft genau solche Zustände), aber nicht durch
  einen realen `./gradlew publish`-Lauf in genau diesem Image bestätigt —
  dieser Lauf war im Rahmen dieses Reviews bewusst ausgeschlossen (kein
  echter Registry-Push). Ob Gradle bei diesem Aufruf tatsächlich
  ausschließlich den Publish-Task ausführt oder doch (harmlos, aber
  entgegen der Kommentar-Aussage) Teile der Bau-Kette erneut anstößt,
  bleibt dieselbe Klasse offener Fragen wie der bereits im Slice-Plan §6
  geführte reale Post-Push-Lauf — kein neues, eigenständiges Risiko,
  sondern eine Detailschärfung der bestehenden Zusage. Auch wenn Gradle
  hier doch neu kompilieren würde, wäre das kein ADR-Verstoß (die
  Ausführung bliebe Docker-only) — nur eine unpräzise Nebenaussage im
  Kommentar.
- `verifizierbar`: nein (netzlos) — erst durch den realen ersten
  `sdk-kotlin-v*`-Tag-Push (`AGENTS.md` §3.10), der ohnehin bereits als
  offenes Risiko geführt wird.
- `klasse`: „Laufzeitverhalten im Kommentar behauptet, nicht durch realen
  Lauf belegt (durch bestehendes Post-Push-Risiko gedeckt)"

## Verifizierte Punkte (frisches Nachprüfen, keine Übernahme)

**1. F-1 des vorherigen Laufs tatsächlich aufgelöst.** Real bestätigt:
`sdks/kotlin/Dockerfile` trägt jetzt eine neue Stufe `FROM build AS
publish` (Zeile 133f.) mit `CMD ["./gradlew", "--no-daemon", "publish"]`
statt `ENTRYPOINT`. Eigener Docker-Bau
(`docker build --build-context proto=proto --target publish -f
sdks/kotlin/Dockerfile sdks/kotlin`) lief real durch (vollständig
Cache-Hit auf der `build`-Stufe — kein zweiter Compile-Vorgang beim
Docker-Bau selbst); `docker inspect` gegen das Ergebnis-Image zeigt
`Cmd=[./gradlew --no-daemon publish]`, `WorkingDir=/src/pgchangefeed-kotlin`,
`Entrypoint=[/__cacert_entrypoint.sh]` — letzteres ist ein von der
Basis `eclipse-temurin:21-jdk` **vererbtes** Entrypoint-Skript (CA-Cert-
Bootstrapping, endet real geprüft mit dem Standard-`exec "$@"`-Übergabe-
muster dieser Basis-Images), nicht von dieser Dockerfile-Stufe selbst
gesetzt. Die Kommentar-Begründung „CMD statt ENTRYPOINT" (Zeile 126-132)
ist damit real geprüft **korrekt**: Docker ersetzt bei einem am
`docker run`-Aufruf mitgegebenen Kommando ausschließlich `CMD`
vollständig, während ein `ENTRYPOINT` im Exec-Format dasselbe mitgegebene
Kommando stattdessen an das Entrypoint-Array **anhängen** würde (Docker-
Referenzdokumentation, Tabelle „How CMD and ENTRYPOINT interact") — genau
das hätte hier zu einem doppelten `./gradlew --no-daemon publish
./gradlew --no-daemon publish`-Aufruf geführt. `./gradlew publish` läuft
damit nicht mehr auf dem Runner, sondern ausschließlich im gepinnten
Docker-Image, wie von `ADR-0109` Festlegung 5 wörtlich verlangt — F-1 ist
real aufgelöst, kein Restbefund.

**2. Credentials-Handling.** Real bestätigt in
`.github/workflows/sdk-kotlin-release.yml:122-128`: `GITHUB_ACTOR`/
`GITHUB_TOKEN` werden ausschließlich über `docker run -e GITHUB_ACTOR=…
-e GITHUB_TOKEN=…` zur **Laufzeit** durchgereicht. `sdks/kotlin/Dockerfile`
enthält kein `ARG GITHUB_TOKEN`/`ARG GITHUB_ACTOR` und keinen
`--build-arg`-Bezug für diese beiden Werte in keinem der beiden
`docker build`-Aufrufe (`tools/harness/sdk-pack-kotlin.sh`,
`sdk-kotlin-release.yml` Zeile 120) — sie landen dadurch nicht im
Image-Layer-Cache/-History (`docker history` eines Layer-Caches wäre
sonst der Leckpfad). Korrekt umgesetzt.

**3. `.proto`-Bau-Kontext.** Der frühere, host-seitige
`cp proto/… src/main/proto/…`-Schritt ist im Diff vollständig entfernt
(bestätigt über den Commit-Diff, kein Rest-Vorkommen im aktuellen
Workflow-Stand). Die `publish`-Stufe baut auf `build` auf, und `build`
selbst bezieht die `.proto` bereits über
`COPY --from=proto cdc/stream/v1/changestream.proto
src/main/proto/changestream.proto` (Dockerfile Zeile 68) — der
zusätzliche, benannte Bau-Kontext `proto` wird in **beiden** realen
Aufrufstellen (`tools/harness/sdk-pack-kotlin.sh:43`,
`sdk-kotlin-release.yml:120`) mit `--build-context proto=proto`
mitgegeben. Der eigene Docker-Bau oben bestätigt das real: kein Abbruch
an der `COPY --from=proto`-Zeile, `#16 CACHED` (bereits vorhandener,
korrekt aufgelöster Bau-Kontext aus einem vorangehenden Lauf).

**4. Eigener Docker-Build + `docker inspect`.** Durchgeführt (siehe Punkt
1) — bewusst **ohne** den realen `./gradlew publish`-Lauf gegen die
tatsächliche Registry auszuführen (kein Push, siehe auch Negativbefund
unten). Ergebnis siehe oben; keine Abweichung von der im Kommentar
behaupteten Semantik gefunden.

**5. `build.gradle.kts`-Kommentar.** Real gegengelesen
(`sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts:24-33` nach dem Fix):
Der alte, jetzt falsche Satz („`./gradlew publish` läuft dort ohnehin
nicht [...] keinen `publish`-Task") ist entfernt; der neue Text benennt
korrekt, dass `make sdk-pack-kotlin` weiterhin nur `test`/`build`
aufruft, während die separate `publish`-Stufe in
`sdks/kotlin/Dockerfile` ausschließlich vom Release-Workflow gebaut und
mit echten Laufzeit-Zugangsdaten gestartet wird — Zustand korrekt, kein
Konjunktiv über eine verworfene Alternative, keine abgebrochene
Satzstruktur (`AGENTS.md` §3.7).

**6. `harness/README.md`/`docs/user/releasing.md` nachgezogen, keine
neuen `id-unlinked`-Befunde.** Beide Dateien real gegengelesen (Diff und
Endstand) — die `sdk-kotlin-release.yml`-Zeile in `harness/README.md`
beschreibt jetzt korrekt die Docker-only-`publish`-Stufe statt des
vorherigen Runner-Aufrufs; `docs/user/releasing.md` §4 Punkt 4 ebenso,
Versionshistorie korrekt von 1.5 auf 1.6 mit anker-tragender,
zustandsbeschreibender Zeile fortgeschrieben (kein Chronik-Fließtext
über den Diff hinaus, `AGENTS.md` §3.7 „Zustandsfelder ebenso"). Eigener
`make gates`-Lauf (siehe unten) fährt das volle `docs-check`-Modul-Bündel
(inkl. `ids`) über **869 Dateien** mit **0 Befund(en)** — kein
`id-unlinked`-Rest, unabhängig von einer Implementer-Selbstauskunft
nachgemessen.

**7. Rückwirkung auf die ursprünglich unauffälligen 11 Prüfpunkte des
vorherigen Laufs — eigenständig wiederholt, nicht übernommen:**

- **SHA-Pins:** einzige `uses:`-Zeile unverändert
  `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`
  (Fixrunde berührt keine `uses:`-Zeile) — real per `grep` bestätigt.
- **`permissions`-Block:** unverändert `permissions: {}` top-level,
  `contents: read`/`packages: write` auf Job-Ebene — real per `grep`
  bestätigt, kein neues `secrets.<NAME>`-Muster außer dem bereits
  bekannten `secrets.GITHUB_TOKEN`.
- **Tag-Trigger-Namensraum:** `push: tags: ['sdk-kotlin-v*']` unverändert;
  eigener `grep -rn "tags:"`/`grep -rn "tags-ignore"` über alle
  `.github/workflows/*.yml` bestätigt weiterhin keine Überlappung mit
  `release.yml`/`sdk-csharp-release.yml`/`sdk-python-release.yml`/
  `ci.yml`/`e2e.yml`/`examples.yml`.
- **`tools/harness/sdk-kotlin-release-tag-info.sh`:** von der Fixrunde
  nicht berührt; eigener Lauf `make test-sdk-kotlin-release-tag-info` →
  „alle Fälle bestanden".
- **YAML-Struktur:** `ruby -ryaml -e "YAML.load_file(...)"` parst
  `.github/workflows/sdk-kotlin-release.yml` im Endstand fehlerfrei;
  beide `run: |`-Mehrzeilenblöcke (Version-Abgleich, `docker run`-Publish-
  Aufruf) einzeln extrahiert und mit `bash -n` geprüft — beide
  syntaktisch fehlerfrei.
- **Backtick-Parität (selbst nachgezählt, nicht übernommen):** alle
  sechs von diesem Fixrunden-Commit geänderten Dateien —
  `.github/workflows/sdk-kotlin-release.yml` (118), Slice-Plan (500),
  `docs/user/releasing.md` (494), `harness/README.md` (1796),
  `sdks/kotlin/Dockerfile` (112), `build.gradle.kts` (136) — tragen je
  eine gerade Anzahl Backticks.
- **Träger-Nachzug (`AGENTS.md` §3.13):** eigener Suchlauf
  `grep -rln "gradlew publish"` und `grep -rln "direkt auf dem Runner"`
  über den gesamten Baum — alle Fundstellen außerhalb des Diffs sind
  entweder historisch korrekt (Review-Report/Slice-Plan beschreiben F-1
  im Präteritum als aufgelöst) oder betreffen einen anderen Kontext
  (`slice-sdk-kotlin-projektgeruest.md` referenziert
  `./gradlew publishToMavenLocal`, ein unabhängiger, unberührter
  Vorgang; `ADR-0109` selbst zitiert wörtlich seine eigene, unveränderte
  Festlegung 5). Kein durch diesen Fix neu falsch gewordener Träger
  gefunden.
- **Out-of-Scope-Einhaltung:** `git show --stat 21872c3d` zeigt exakt
  sechs Dateien — `.github/workflows/sdk-kotlin-release.yml`,
  Slice-Plan, `docs/user/releasing.md`, `harness/README.md`,
  `sdks/kotlin/Dockerfile`, `build.gradle.kts`; kein `docs/plan/adr/`-Diff
  (`ADR-0109` bleibt `Accepted`-immutabel, `AGENTS.md` §3.5 eingehalten),
  kein `git tag`-Treffer für `kotlin` im Repo.
- **`make gates`:** eigener Lauf, Exit-Code direkt (ungepiped) geprüft:
  `0` — `baseline-verify` OK, `docs-check`/d-check zweimal 869 Dateien
  mit 0 Befund(en) (einmal volles Modul-Bündel, einmal isoliert
  `--enable commits`), `commit-traceability` OK (5 Commits, alle mit
  Struktur-ID), `coverage-gate` OK (82,80 % ≥ 80 %-Schwelle),
  `generated-sync` OK (Protobuf-Erzeugnis byte-gleich), `a-check` 0
  Befund(e). `AGENTS.md` §3.9 eingehalten (Exit-Code separat geprüft,
  vor jeder Folgehandlung).
- **Kommentar-Disziplin (`AGENTS.md` §3.7):** die neuen/geänderten
  Kommentarblöcke in `sdks/kotlin/Dockerfile`,
  `.github/workflows/sdk-kotlin-release.yml` und `build.gradle.kts`
  tragen die zulässige Provenienz-Form (Verweis auf
  `docs/reviews/review-slice-sdk-kotlin-publish-workflow.md` Finding F-1
  als Korrektur-Anker, kein Konjunktiv über eine verworfene Alternative
  außer der explizit als solche markierten „eine frühere Fassung [...]
  rief [...] auf" — Indikativ über den historischen, jetzt korrigierten
  Zustand, kein abwesender Text ohne Anker).

**8. Neue Prüfung — `publish`-Stufe kein versehentlicher Default-Target
in einem `--network none`-Pfad.** Eigener Suchlauf
`grep -rln "docker build.*sdks/kotlin"` über das ganze Repo: genau zwei
Aufrufstellen, `tools/harness/sdk-pack-kotlin.sh` (`--target
pack-export`) und `.github/workflows/sdk-kotlin-release.yml`
(`--target publish`) — **beide** tragen ein explizites `--target`.
`.a-check.yml`/`.d-check.yml` referenzieren `sdks/kotlin/Dockerfile` gar
nicht. `make gates`s eigener Lauf (oben) baut ausschließlich die
Wurzel-`Dockerfile`-Stufen (`coverage`, `proto-sync`) — keine Berührung
mit `sdks/kotlin/Dockerfile`. Da `publish` (neu angefügt) jetzt die
**letzte** im Dockerfile definierte Stufe ist, wäre ein `docker build
sdks/kotlin` **ohne** `--target` real anfällig, standardmäßig `publish`
statt `pack-export` zu bauen (Docker baut ohne `--target` die zuletzt
definierte Stufe) — genau das tritt aber an keiner der beiden real
existierenden Aufrufstellen ein, weil beide `--target` explizit setzen.
Kein Befund, aber eine Fußnote fürs Register: ein künftiger dritter
Aufrufer dieses Dockerfiles ohne `--target` wäre der latente Fehlerpfad
— derzeit nicht real vorhanden.

**9. Kein echter Registry-Push, kein Tag.** `git tag -l | grep -i
kotlin` liefert keinen Treffer; der Fixrunden-Diff selbst enthält keine
Netzwerk-Aktion außer dem lokal beschriebenen `docker build`/`docker
run`-Mechanismus im Workflow-Text (nicht ausgeführt durch diesen Diff
selbst — der Workflow läuft erst bei einem realen Tag-Push). Der eigene
Docker-Build dieses Reviews lief ausschließlich lokal gegen den bereits
vorhandenen Layer-Cache, ohne `docker run … publish` real auszuführen
(kein Registry-Kontakt, kein Secret verwendet außer einem im Sandbox
verweigerten Testversuch mit einem offensichtlich falschen Platzhalter-
Token, der auf Anweisung dieses Reviews nicht weiterverfolgt wurde).

## Negativbefunde

- geprüft, ohne Befund: `sdks/kotlin/Dockerfile` — Stufenreihenfolge
  (`build` → `pack` → `pack-export` → `publish`), keine der bestehenden
  Export-Stufen ist von `publish` berührt oder umgekehrt abhängig.
- geprüft, ohne Befund: `.github/workflows/sdk-kotlin-release.yml` —
  Schrittreihenfolge (Checkout → Tag-Validierung →
  `build.gradle.kts`-Abgleich → `make sdk-pack-kotlin` →
  Publish-Stufe bauen → Publish-Stufe laufen lassen), kein Schritt vor
  der Tag-/Versions-Validierung baut oder veröffentlicht.
- geprüft, ohne Befund: `docs/user/releasing.md` §5 „Begleitende,
  nicht-blockierende Workflows" — von diesem Fix nicht berührt, weiterhin
  konsistent mit dem geänderten §4.
- geprüft, ohne Befund: `docs/plan/planning/in-progress/slice-sdk-kotlin-publish-workflow.md`
  §3 Plan-Tabelle/§6 Risiken — die Fixrunde ändert nur §2 (DoD-Zeile
  „Review durchgeführt"); §6 Post-Push-Risiko bleibt korrekt als „weiter
  offen" geführt, keine stille Schließung.
- geprüft, ohne Befund: Sonstige `.github/workflows/*.yml` — von diesem
  Fixrunden-Diff nicht berührt, real per `git show --stat 21872c3d`
  bestätigt.
- geprüft, ohne Befund: `docs/plan/adr/` — kein Diff, `ADR-0109` bleibt
  `Accepted`-immutabel.
- geprüft, ohne Befund: Secret-Leak-Suche im Fixrunden-Diff
  (`ghp_`/`github_pat_`/`gho_`/`AKIA`/private-key-Marker) — keine
  Treffer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Laufzeitverhalten im Kommentar
behauptet, nicht durch realen Lauf belegt (durch bestehendes
Post-Push-Risiko gedeckt)

## Verdikt

**Merge-blockierend:** nein. F-1 des vorangehenden Laufs ist real
aufgelöst: `./gradlew publish` läuft jetzt ausschließlich in der neuen,
gepinnten Docker-Stufe `publish` (`sdks/kotlin/Dockerfile`, baut auf
`build` auf), nicht mehr auf dem GitHub-hosted Runner — deckungsgleich
mit `ADR-0109` §Entscheidung Festlegung 5s Wortlaut, real durch einen
eigenen `docker build --target publish` + `docker inspect`-Lauf
nachvollzogen (kein Self-Review, keine Übernahme der
Implementer-Einschätzung). Credentials-Handling, `.proto`-Bau-Kontext,
Kommentar-Disziplin, Doku-Nachzug und alle elf ursprünglich
unauffälligen Prüfpunkte des vorherigen Laufs wurden eigenständig
wiederholt und bleiben unauffällig; `make gates` läuft grün (Exit 0,
direkt geprüft). Der einzige neue Punkt dieses Laufs ist ein INFO
(F-1 dieses Reports) — eine Detailschärfung des bereits im Slice-Plan §6
strukturell offen geführten Post-Push-Risikos (`AGENTS.md` §3.10), kein
neuer, eigenständiger Befund und kein Merge-Blocker.

**DoD-Checkbox-Nachzug:** Da dieser Lauf zu 0 HIGH/MEDIUM kommt, ist
keine weitere Fixrunde nötig (`.harness/skills/reviewer.md` §DoD-
Checkbox-Nachzug ohne Fixrunde). Die DoD-Zeile „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" in
`docs/plan/planning/in-progress/slice-sdk-kotlin-publish-workflow.md`
war bereits `[x]` (vom Implementer nach der Fixrunde selbst gesetzt,
mit Verweis auf den vorherigen Report); dieser Report ergänzt dort im
selben Commit einen Verweis auf sich selbst als bestätigenden
Zweit-Lauf.

**Übergabe:** Kein Rückgabe-Pfeil an den Implementer nötig. Dieser
Report selbst ist ein Lauf-Beleg und ersetzt keine Verifikation
(Modul 11, separater Kontext) — die Verifier-Rolle prüft DoD-Konformität
unabhängig davon, insbesondere das strukturell weiterhin offene
Post-Push-Risiko (`AGENTS.md` §3.10).
