# Review-Report: slice-019 — 2026-09-12

**Review-Art:** Code — geprüft gegen Slice-Plan
(`docs/plan/planning/in-progress/slice-019-cdc-capture-lag-ablösen.md`)
und `welle-5.md` (Closure-Trigger der Welle; Maintainability).

**Gegenstand:** fünf Commits auf `main`:
`9de0e1b` (`fix(schema): cdc_capture_lag_approx nach cdc_capture_lag
umbenannt`, `tools/schema/nacharbeit-observability.sql`),
`47c8909` (`docs(user): Benutzerhandbuch nennt cdc_capture_lag statt der
Näherung`, `docs/user/benutzerhandbuch.md`),
`3ee79a4` (`test(integration): Ende-zu-Ende-Lasttest-Beleg für
cdc_capture_lag`, `tools/harness/run-integration-tests.sh`),
`e7232b5` (`chore(image): Image-Hash-Beleg nach dem
Commit-Zeitstempel-Fix nachgezogen`, `harness/image-hash.txt`),
`88d0803` (`docs(planning): slice-019 Plan-Nachzug + DoD-Häkchen`,
Slice-Plan §2/§3).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-019-cdc-capture-lag-ablösen.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Nachzug, §4 Trigger, §6
  Risiken, §8 Sub-Area-Prüfung)
- `docs/plan/planning/welle-5.md` (§1 Welle-Ziel, §3 Closure-Trigger —
  dieser Slice liefert den Ende-zu-Ende-Lasttest-Beleg der Welle)
- `spec/pflichtenheft.md` `SPEC-009` (Metrik-Katalog, kanonischer Name
  `cdc_capture_lag`) und `SPEC-013` (`CDC_LAG_THRESHOLDS`)
- [`LH-FA-ADM-004`](../../spec/lastenheft.md) (messbarer CDC-Abstand)
- [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Image-Beleg-
  Semantik — Digest ist Lauf-Beleg, welle-1-Regel: Build-Kontext-Änderung
  → `make image` vor Closure)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.3 Move/Inhalt-Trennung
  — hier nicht einschlägig, kein Lifecycle-Übergang in diesen fünf
  Commits —, §3.7 Kommentar-Klassen)
- `docs/user/benutzerhandbuch-standard.md` (Stil-Maßstab für die
  Benutzerhandbuch-Korrektur)
- Commit-Traceability (`AGENTS.md` §5, `ADR-0045`)
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/
  cdc-capture-lag-real/` (Auflösungs-Träger dieses Slice, Ausgang
  `eingetreten` erst bei Closure fällig — nicht Gegenstand dieses Diffs)
- Vorherige Findings am gleichen Modul: `docs/reviews/review-slice-018.md`
  F-1 (künstliche Test-Sleep ohne Unterscheidungskraft — direkt
  einschlägiges Präzedenzmuster für Prüfpunkt 3 dieses Laufs, siehe
  Prüfung 3 unten) und `docs/reviews/verify-slice-018.md` (nennt die noch
  offene `cdc_capture_lag_approx`-Umbenennung als DoD-Rest, den dieser
  Slice abschließt)

---

## Prüfungen im Detail

**1. Diff gegen Plan — Metrik-Umbenennung vollständig.** `git grep -rn
"cdc_capture_lag_approx"` über das ganze Repo liefert nach diesem Diff nur
noch Treffer in `done/slice-013…`, `done/slice-017…`,
`done/slice-018-…md` (Plan-Zeile "Metrik-Umbenennung … steht noch aus"),
`done/welle-4-results.md`, `docs/reviews/verify-slice-018.md`,
`docs/reviews/review-comment-cleanup.md` und in `slice-019` selbst
(§1 Ziel-Satz, §6 Risiko-2 — beide beschreiben den *alten* Namen als das,
was abgelöst wird) sowie in `welle-5.md` (Welle-Ziel-Beschreibung). Alle
das sind Zeitdokumente bzw. die eigene Ziel-Beschreibung dieses Slice —
keine aktive Spec-Stelle (`SPEC-009`/`SPEC-013` nennen bereits den
kanonischen Namen, unverändert seit vor diesem Slice) und kein
Produktcode (`grep` über `*.go` ist leer). Konsistent mit Plan §2 DoD
Punkt 1.

**2. Der Lasttest-Beleg real.** `tools/harness/run-integration-tests.sh`
misst `cdc_capture_lag` zweimal per `psql -tAc "SELECT value FROM
cdc.metrics WHERE metric_name = 'cdc_capture_lag'"` — real gegen die
SQL-View, nicht gegen einen Anwendungs-internen Wert. Die
`docker pause`/`unpause`-Technik ist sauber: Der `EXIT`-Trap (bereits vor
diesem Diff registriert) wurde um einen `docker unpause … || true` am
Anfang von `cleanup()` ergänzt — das läuft *vor* `$COMPOSE down -v`, sonst
hinge das Abräumen an einem eingefrorenen Container. Eigener Testlauf
(`make test-integration`, zweimal hintereinander): beide Male grün,
Baseline 0,1297 s / 0,1301 s, verzögert 1,2175 s / 1,1742 s — reproduziert
die im Plan genannten Bereiche. Die Toleranzprüfung (`awk … d*0.5`) trennt
mit deutlichem Abstand (Baseline < 0,5 s Margin ≈ 0,37 s; verzögert > 0,5 s
Margin ≈ 0,65–0,72 s bei den beobachteten Werten) — für den Default
(`LAG_DELAY_SECONDS=1`) robust. Siehe aber F-1 zur Override-Fläche.

**3. Der `wal_sender_timeout`-Kompromiss.** Die Präzedenz-Begründung
(Kommandozeilen-`-c`-Flag in `compose.yaml` sticht `ALTER SYSTEM`, weil
`ALTER SYSTEM` nur `postgresql.auto.conf` schreibt, das unterhalb von
Kommandozeilen-GUCs rangiert) ist technisch zutreffend — eigene Prüfung
gegen `compose.yaml:27` (`-c max_replication_slots=10 -c
wal_sender_timeout=2000`) bestätigt den Wert. Die Wahl 1 s Pause bei 2 s
Timeout ist durch die beiden Messläufe oben bestätigt robust: Faktor
~9–12× zwischen Baseline und verzögerter Messung, weit über dem
Rauschen. Anders als im Präzedenzfall `review-slice-018 F-1` (künstliche
Sleep ohne Beitrag zur Unterscheidungskraft) trägt die Verzögerung hier
tatsächlich die gesamte Testaussage — ohne sie gäbe es keinen Kontrast.
Siehe F-1 zur fehlenden Absicherung, falls `LAG_DELAY_SECONDS`
überschrieben wird.

**4. Der Image-Rebuild-Fund.** `harness/image-hash.txt` trägt jetzt
`sha256:b4c49bf2…` (vorher `sha256:ed23e283…`). Eigene Prüfung: das lokal
geladene Image `ghcr.io/pt9912/pg-change-feed:dev` hat exakt diese
Image-ID (`docker images`) und ein `Created`-Datum `2026-09-12T07:55:44`
— nach beiden Commit-Zeitstempeln von `3fa7e5d` (06:26) und `f35a94d`
(06:57) desselben Tages. Die Diagnose („altes Image trug den Fix nie")
ist damit nicht nur behauptet, sondern durch Zeitstempel-Vergleich
nachvollziehbar, zusätzlich indirekt bestätigt durch den eigenen
Testlauf: Der Lasttest misst mit dem neu geladenen Image korrekt den
realen Commit-Zeitstempel-Abstand (Prüfung 2). Vorgehen entspricht
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (welle-1-Regel:
Build-Kontext-relevanter Zug → `make image` vor Closure, Digest-Commit
nur bei geändertem Digest — hier geändert, also korrekt committet).

**5. Benutzerhandbuch-Korrektur.** Der neue Satz („Abstand zwischen der
letzten Quelländerung und der CDC-Verfügbarkeit, gemessen über den
Commit-Zeitstempel aus dem WAL") ist fachlich deckungsgleich mit dem
SQL-Kommentar und `SPEC-009`, im Indikativ, ohne Konjunktiv über den
verworfenen `_approx`-Zustand. Stilistisch konsistent mit der
unmittelbar vorangehenden Passage im selben Dokument (Heartbeat-Absatz
erklärt ebenfalls in einem Klammer-Einschub, wie ein Wert zu lesen ist) —
kein Bruch mit `benutzerhandbuch-standard.md` (aufgabenbasiert, direkte
Sprache, kein "man kann").

**6. Kommentar-Stil (Hard Rule 3.7).** `nacharbeit-observability.sql`:
neuer Kommentar-Block ist präsentisch/indikativ, trägt `SPEC-009`,
`LH-FA-ADM-004`, `SPEC-013` als auflösbare Anker, keine Slice-Referenz,
kein Verweis auf einen "früheren" Text oder eine verworfene Alternative
im Konjunktiv. `run-integration-tests.sh`: der neue Kommentar-Block vor
dem Lasttest beschreibt, was die folgenden Zeilen tun und warum
(`wal_sender_timeout`-Grenze), ohne Chronik-Sprache oder Slice-ID. Beide
Kommentare bestehen die Prüfung auf die fünf zulässigen Klassen (Zusage ·
Kopplung · Abgrenzung · Rang-Zeiger · Grenze).

**7. Hard Rules allgemein.** Docker-only: Lasttest-Block läuft
ausschließlich über bereits laufende Docker-/Compose-Container
(`docker exec`, `docker pause`/`unpause`), kein lokales Toolchain-Install;
das neu genutzte `awk` ist bereits an drei anderen Stellen im Repo
(`tools/harness/working-tree-hash.sh`, `tools/harness/image-stale.sh`,
`tools/harness/extract-command.awk`) etabliertes Muster, kein neuer
Werkzeug-Bruch. Suppression-Verbot: kein `nolint`/`noqa`/
`SuppressMessage`/`d-check:ignore` im Diff. Git-mv+Inhalt-Trennung: keine
dieser fünf Commits enthält einen `git mv` — nicht einschlägig.
Traceability: alle fünf Commit-Betreffs tragen `LH-FA-ADM-004`, keine
`SPEC-*`/`ARC-*`-Kennung im Betreff; `make commit-traceability` lief in
dieser Sitzung gegen `HEAD~5..HEAD` grün. `make gates` lief in dieser
Sitzung eigenständig (nicht nur den Implementer-Bericht übernommen):
grün, 0 Befunde in allen vier inneren Gates.

**8. Was für Closure/Welle-Closure noch fehlt (nur genannt, nicht
bewertet):** Review-Schluss selbst (dieser Report), Closure-Notiz §7 des
Slice-Plans, Beobachtungs-Register-Eintrag `BEO-PGC/cdc-capture-lag-real`
(Ausgang `eingetreten`), Ausgänge für beide §6-Risiken, der `git mv` nach
`done/`, und — weil dieser Slice den Closure-Trigger von `welle-5`
liefert — im Anschluss die komplette Welle-5-Closure-Prozedur (Schritte 1
und 3–6 aus Modul 6, inklusive der drei Paarungen und der
Zeitdokument-Archivierung).

---

## Findings

### F-1 — `LAG_DELAY_SECONDS`-Override ohne Grenzprüfung gegen `wal_sender_timeout`

- `kategorie`: MEDIUM
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:156-157` (`LAG_TABLE=…`,
  `LAG_DELAY_SECONDS=${LAG_DELAY_SECONDS:-1}`) und die begründende
  Kommentar-Passage `tools/harness/run-integration-tests.sh:143-155`
- `befund`: Der Kommentar begründet den Default (1 s) explizit mit der
  `wal_sender_timeout=2000`-Grenze aus `compose.yaml` — eine längere
  Pause lasse die Quelle die Replication-Verbindung beenden und den
  Feed-Container (ohne Restart-Policy) abstürzen. Das Skript übernimmt
  `LAG_DELAY_SECONDS` aber ungeprüft aus der Umgebung; wählt ein
  künftiger Lauf einen Wert ≥ ~2 s (z. B. um eine größere Verzögerung zu
  demonstrieren), tritt genau das im Kommentar beschriebene
  Absturzszenario ein, aber das Skript meldet dafür keinen eigenen,
  erklärenden Fehler — der Lauf scheitert stattdessen an der generischen
  Meldung „verzögerte Transaktion … wurde nach dem Fortsetzen nicht
  erfasst" (Zeile 202), die die eigentliche Ursache (WAL-Sender-Timeout,
  Container-Absturz) nicht benennt.
- `verifizierbar`: ja — `LAG_DELAY_SECONDS=3 make test-integration`
  ausführen und beobachten, ob die Fehlermeldung die Timeout-Ursache
  benennt (kein automatisiertes Gate prüft das).
- `klasse`: Unclamped Env-Override kollidiert mit einer im selben Diff
  dokumentierten harten Infrastruktur-Grenze, ohne eigene Diagnose.

### F-2 — Hartkodierte Test-IDs ohne Reservierungs-Hinweis

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:162,186` (`VALUES (90,
  'LagBaseline')`, `VALUES (91, 'LagDelayed')`)
- `befund`: Die IDs `90`/`91` in `public.feed_mvp_full` sind frei gewählte
  Konstanten ohne Kommentar, der ihre Reservierung gegen die IDs anderer
  Testfälle auf derselben Tabelle (aktuell `id=1` in
  `test/integration/mvp_test.go`) festhält. Aktuell kollisionsfrei, aber
  ein künftiger Testfall auf `feed_mvp_full` mit dem Wertebereich 90–99
  würde still auf `ON_ERROR_STOP` scheitern statt auf einen erklärenden
  Reservierungs-Hinweis zu treffen.
- `verifizierbar`: nein — Stilfrage, kein Gate-Gegenstand.
- `klasse`: Unkommentierte Konstanten-Reservierung über Testfall-Grenzen
  hinweg.

### F-3 — `BEO-PGC/cdc-capture-lag-real` und §6-Risiko-Ausgänge zum Prüfzeitpunkt offen

- `kategorie`: INFO
- `quelle`: Maintainability (Slice-Plan §6, §2 DoD, §7)
- `pfad`:
  `docs/plan/planning/in-progress/slice-019-cdc-capture-lag-ablösen.md`
  §6, §7 (leer)
- `befund`: Beide Risiken aus §6 (Teilstrecken-Abdeckung,
  externe `_approx`-Konsumenten) und der Register-Eintrag
  `BEO-PGC/cdc-capture-lag-real` tragen noch keinen der drei
  Closure-Ausgänge — erwartungsgemäß, da die Datei noch in
  `in-progress/` liegt und dieser Diff kein Closure-Commit ist. Kein
  Diff-Defekt, sondern ein Hinweis für die nächste Rolle (Planner bei
  Closure).
- `verifizierbar`: nein — Planungs-/Closure-Frage, kein Gate-Gegenstand.
- `klasse`: Offenes Risiko ohne Ausgang zum Review-Zeitpunkt —
  erwartungsgemäß vor Closure (identische Klasse wie
  `review-slice-018.md` F-2).

---

## Negativbefunde

- geprüft, ohne Befund: **verbleibende `cdc_capture_lag_approx`-Treffer
  im Repo** — ausschließlich in `done/`-Zeitdokumenten, laufenden
  Review-/Verify-Reports und der eigenen Ziel-Beschreibung von
  `slice-019`/`welle-5`; kein Treffer in Produktcode, aktiver Spec oder
  Benutzerdoku.
- geprüft, ohne Befund: **`SPEC-009`/`SPEC-013`** — nennen bereits vor
  diesem Diff den kanonischen Namen `cdc_capture_lag`; `spec/
  pflichtenheft.md` unverändert, wie im Plan-DoD behauptet.
- geprüft, ohne Befund: **Go-Code** — kein `.go`-Treffer für
  `cdc_capture_lag`/`_approx` in beide Richtungen; die Metrik lebt
  ausschließlich in der SQL-View, keine Zwei-Quellen-Drift zwischen
  Anwendungscode und Schema.
- geprüft, ohne Befund: **`docker pause`/`unpause`-Aufräumen bei
  Fehlerfällen** — `cleanup()` ruft `docker unpause … || true` vor
  `$COMPOSE down -v`; ein Abbruch zwischen `docker pause` und
  `docker unpause` (z. B. durch eine fehlschlagende `INSERT`) lässt den
  Container nicht eingefroren zurück.
- geprüft, ohne Befund: **Eigener `make test-integration`-Lauf, zweimal**
  — beide grün, `cdc_capture_lag`-Werte reproduzieren die im Plan
  genannten Bereiche (Baseline ≈0,13 s, verzögert ≈1,17–1,22 s).
- geprüft, ohne Befund: **`wal_sender_timeout`-Präzedenz-Behauptung** —
  Kommandozeilen-`-c`-Flag rangiert über `ALTER SYSTEM`
  (`postgresql.auto.conf`); technisch zutreffend.
- geprüft, ohne Befund: **Image-Rebuild-Diagnose** — Image-Erstellzeit
  (`docker inspect … Created`) liegt nach beiden slice-017/018-Commit-
  Zeitstempeln; Diagnose nicht nur behauptet, sondern zeitlich
  nachvollziehbar.
- geprüft, ohne Befund: **`ADR-0044`-Konformität des Digest-Commits** —
  Digest geändert, also Commit korrekt gesetzt (welle-1-Regel: entfiele
  bei unverändertem Digest).
- geprüft, ohne Befund: **Benutzerhandbuch-Stil** — aufgabenorientiert,
  Indikativ, konsistent mit `benutzerhandbuch-standard.md` und dem
  umgebenden Dokument-Stil.
- geprüft, ohne Befund: **Kommentar-Klassen** (Hard Rule 3.7) — beide
  geänderten Kommentar-Blöcke (`nacharbeit-observability.sql`,
  `run-integration-tests.sh`) tragen Anker (`SPEC-009`, `SPEC-013`,
  `LH-FA-ADM-004`), keine Slice-Referenz, keine Konjunktiv-Rede über eine
  verworfene Alternative, kein abgebrochener Satz.
- geprüft, ohne Befund: **Docker-only** — kein lokales
  Toolchain-Install; `awk`-Nutzung ist etabliertes Repo-Muster.
- geprüft, ohne Befund: **Suppression-Verbot** — kein
  `nolint`/`noqa`/`SuppressMessage`/`d-check:ignore`-Marker im Diff.
- geprüft, ohne Befund: **Hard Rule 3.3 (git mv + Inhalt)** — keiner der
  fünf Commits verschiebt eine Datei; nicht einschlägig für diesen
  Implementer-Lauf.
- geprüft, ohne Befund: **Traceability der Commit-Messages** — alle fünf
  Betreffs tragen `LH-FA-ADM-004`, keine `SPEC-*`/`ARC-*`-Kennung; `make
  commit-traceability` lief in dieser Sitzung gegen `HEAD~5..HEAD` grün.
- geprüft, ohne Befund: **`make gates`** — in dieser Sitzung eigenständig
  ausgeführt, grün, 0 Befunde (baseline-verify, docs-check,
  commit-traceability, a-check).
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung des Slice-Plans** —
  Sub-Area bleibt `PGC`/Observability, Modus GF, kein Modus-Wechsel durch
  diesen Diff; Register-Sichtung (elf Einträge, `cdc-capture-lag-real`
  1×) konsistent mit dem tatsächlichen Registerstand.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Unclamped Env-Override kollidiert mit
dokumentierter Infrastruktur-Grenze ohne eigene Diagnose (F-1) ·
Unkommentierte Konstanten-Reservierung über Testfall-Grenzen hinweg
(F-2) · Offenes Risiko ohne Ausgang zum Review-Zeitpunkt — erwartungsgemäß
vor Closure (F-3, identische Klasse wie `review-slice-018.md` F-2 — zweites
Auftreten dieser Klasse).

## Verdikt

**Merge-blockierend:** nein — kein HIGH-Finding, kein Rollen-Widerspruch.
Das einzige MEDIUM (F-1) betrifft eine Override-Fläche eines
kein-Gate-Testskripts mit sicherem Default (1 s < 2 s Timeout, durch
zwei eigene Läufe bestätigt); der Implementer kann es akzeptieren oder
mit Begründung zurückweisen (kein Rollen-Widerspruch, keine Sequenz nach
Modul 8 nötig — dies ist zudem erst das erste Auftreten dieser
Finding-Klasse).

**Übergabe:** Findings gehen an den Implementer. Die Finding-Klassen
gehen zusätzlich in die Slice-Closure §7 und von dort in den
Beobachtungs-Register-Zähler — insbesondere F-3s Klasse ist bereits beim
zweiten Auftreten (nach `review-slice-018.md` F-2); ein drittes
Auftreten wäre ein Steering-Loop-Signal. Dieser Report ist ein
Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen. Er ersetzt
keine Verifikation — DoD-/Spec-Konformität (inkl. der offenen
Closure-Pflichten aus Prüfung 8) prüft der Verifier separat.
