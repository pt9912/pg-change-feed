# Verifikationsbericht: slice-087 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-087` §2) und die im Slice referenzierten Entscheidungen
([`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) §Fitness
Function (die Zeile, die dieser Slice einlöst) und §Testabdeckung,
[`ADR-0057`](../plan/adr/0057-http-grpc-api.md) (der bestehende HTTP-Rundlauf),
[`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(ein Client, begrenzte Import-Berechtigung),
[`ADR-0030`](../plan/adr/0030-testpyramide.md) (E2E-Tier, kein neues Gate),
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Digest-Beleg);
`AGENTS.md` §3.1, §3.9, §3.10, §3.11). **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe; `review-slice-087.md` wurde als Kontext gelesen, **nicht**
als Beleg übernommen) und **nicht** gegen realen Bedarf (Validator, hier
nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf hat den Slice-Plan in der Fassung von
`HEAD` gelesen, dazu die fünf Entscheidungen, den Review-Report und die
berührten Artefakte. **Alle** Zahlen dieses Berichts stammen aus eigenen, in
dieser Sitzung gefahrenen Läufen; kein Beleg des Implementers, des Reviewers
oder des Planners wurde übernommen. Exit-Codes je **ungepiped** und in
eigenem Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und Folgehandlung
getrennt beauftragt. Mutationen wurden aus einer `cp`-Kopie
byte-gleich zurückgenommen; der Arbeitsbaum ist am Ende dieses Laufs
unverändert (`git status --porcelain` leer).

**Gegenstand:** `slice-087`, geprüfter Stand `HEAD` = `aede3f4`
(„Review slice-087"), Zweig `main`. Der Slice liegt weiterhin in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt. Der
Implementer-Commit ist `93ac92e`, der Review-Commit `aede3f4`.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make test-integration` (unmutiert; Ausgabe in Logdatei, Exit direkt gelesen) | **0** | `READ changes=1 table=feed_e2e_full schema=public change_id=966-1 operation=INSERT commit_position=30901520 new_image={"id":"285","name":"HttpChangesReadE2ESentinel"}`; Beleg-Zeile nennt „GET /changes real per HTTP mit reader-Token (die eigens eingefügte Zeile id=285/HttpChangesReadE2ESentinel, Bereich [30901520,30901521), change_id=966-1 gegen cdc.changes gehalten)". `git status --porcelain` danach **leer** |
| 2 | `make gates` (ungepiped, Exit direkt gelesen) | **0** | baseline-verify `v6.5.0` OK (54 Dateien) · d-check **698** Dateien / **0** Befunde · commit-traceability OK (5 Commits, Betreffe ohne Struktur-ID) · a-check **0** Befunde · coverage-gate **OK — 71.90 %** erfüllt Schwelle 70 % |
| 3 | `make doc-commits RANGE=997a726..HEAD` | **0** | „698 Datei(en) geprüft, 0 Befund(e)" — jede Message der Range mit `LH-*`/`ADR-*`, keine `SPEC-*`/`ARC-*` im Betreff |
| 4 | `make doc-immutable RANGE=997a726..HEAD` | **0** | „698 Datei(en) geprüft, 0 Befund(e)" — kein Accepted-/Record-Artefakt inhaltlich überschrieben |
| 5 | **Eigene Gegenprobe** — Client druckt in der `READ`-Zeile eine **fabrizierte** `change_id` (`erfunden-…`) | **2** | `HTTP-API-Rundlauf — die über GET /changes gelesene Änderung (change_id=erfunden-966-1, HttpChangesReadE2ESentinel) ist nicht real über cdc.changes lesbar (count=0)` — die **`change_id`-Bindung gegen `cdc.changes`** (`run-integration-tests.sh:2043-2053`) ist real und fängt eine erfundene Client-Ausgabe |
| 6 | **Wiederholung der DoD-Gegenprobe** — Client ruft `/changes-gegenprobe` statt `/changes` (falscher Pfad) | **2** | `httpclient: GET /changes (reader) fehlgeschlagen: status 404 (erwartet 200): 404 page not found` → `endete mit Ausgang 1` — deckt die Implementer-Angabe (Liefer-Punkt 2, K3) **wörtlich** |
| 7 | Tabellen-Byte-Vergleich nach vollem Lauf | — | `docs/user/e2e-abdeckung.md` sha256 vor = nach (`e6fde2c3…`), `git status --porcelain` leer — die committete Fassung ist **byte-gleich** zu dem, was ein voller Lauf schreibt |
| 8 | Zählprobe `docs/user/e2e-abdeckung.md` | — | **13** Go-Zeilen (`test/integration/integration_test.go:*`) + **24** Bash-Zeilen (`run-integration-tests.sh:*`) = **37** Zeilen — deckt die DoD-Angabe |
| 9 | `git diff --name-only 997a726..HEAD` | — | `roadmap.md` · `slice-087-…md` · `review-slice-087.md` · `docs/user/e2e-abdeckung.md` · `harness/README.md` · `tools/harness/httpclient/main.go` · `tools/harness/run-integration-tests.sh` — **keine** Workflow-, `spec/`-, `tools/schema/`-, `.a-check.yml`- oder `Makefile`-Datei |
| 10 | Umgebung nach den Läufen | — | keine `cdc-*`-Container, kein `cdc-*`-Netz zurückgeblieben; alle drei Läufe nahmen ihre Mutationen zurück |

**Anker-Spotprobe (zu Messung 6/7):** Die `abdeckung_declare`-Deklaration
steht in `run-integration-tests.sh:1948`; ihr Anker
(`GET /changes real per HTTP mit reader-Token (die eigens eingefügte Zeile`)
löst auf die Echo-Zeile `2061` auf — genau der Wert in der Tabelle. Dass der
volle Lauf die Tabelle **nicht** korrigiert (Messung 7), ist der stärkere
Beleg: die committeten `Ort`-Werte sind die, die der Generator selbst
berechnet.

**Beobachtung zum Coverage-Wert:** Die DoD-Zeile (und der Review) nennen
`coverage-gate 72,00 %`; mein Lauf maß **71.90 %** — derselbe Commit, beide
über der Schwelle 70 %, Gate in beiden Fällen grün. Der Wert ist also
run-gebunden und um ±0,1 % beweglich; die DoD-**Aussage** („grün") trägt, die
exakte Zahl ist keine stabile Größe (kein Befund, benannt).

---

## 2. DoD-Konformität, Kriterium für Kriterium

Gelesen ist der Text in der Fassung `aede3f4` — sechs Häkchen vom
Implementer (Liefer-Punkte 1–3), das Review-Häkchen vom Reviewer.

| # | DoD-Zeile (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| LP1-K1 | `httpclient` ruft zusätzlich `GET /changes` — mit `source`/optional `schema`/`table`/`from`/`to`/`limit` — und prüft die Antwort **inhaltlich**, nicht nur den Status | **erfüllt** | `readChanges` (`httpclient/main.go:85-146`) baut die Abfrage aus allen sechs Werten (leerer Wert lässt den Parameter weg, `main.go:86-102`) und prüft: nicht leere Liste (`:114`), je Eintrag nicht leere `change_id` (`:121`), bekannte Operation (`:130-134`), `commit_position ≥ 1` (`:135`), Filtertreue (`:124-129`), nicht absteigende Position (`:138`). Im Lauf real ausgeführt (Messung 1: `changes=1 … operation=INSERT`). **Grenze:** die Filtertreue-Prüfung kann in dieser Aufruf-Form nicht feuern — siehe F-1/§6. |
| LP1-K2 | Er bleibt **ein** Client (`ADR-0068`): kein zweites Programm, keine zweite Import-Berechtigung, kein neuer `.a-check.yml`-Eintrag | **erfüllt** | Derselbe `tools/harness/httpclient`; der neue Import ist `net/url` (Standardbibliothek, `main.go:18`) — kein `internal/`-Import. `git diff --name-only` (Messung 9) enthält **keine** `.a-check.yml`; `make a-check` 0 Befunde (Messung 2). Der Bestand `tooling: ["tools/harness/**"]` deckt die Datei unverändert. |
| LP2-K1 | Die HTTP-Rundlauf-Phase ruft ihn real gegen den **laufenden** Feed-Container und wertet aus; bei Scheitern endet die Phase **rot** | **erfüllt** | `run-integration-tests.sh:2001-2014`: `docker run --network "$NETWORK" … go run ./tools/harness/httpclient …` gegen den Compose-Netz-Alias; `http_status=$?` wird ausgewertet und bei `≠ 0` mit eigener Zeile rot. Real belegt: unmutiert grün (Messung 1), mutiert rot (Messungen 5/6, beide Exit 2). |
| LP2-K2 | `abdeckung_declare` trägt den erweiterten Nachweis — Anker und Kurzbeschreibung nennen das Lesen | **erfüllt** | Deklaration `run-integration-tests.sh:1948` nennt „liest zusätzlich Changes über `GET /changes` mit dem reader-Token; … der gelesene Change gegen cdc.changes"; die erzeugte Zeile `docs/user/e2e-abdeckung.md:50` trägt dieselbe Aussage und den Anker-Ort `:2061`. |
| LP2-K3 | **Eine rote Gegenprobe:** der Beleg wird rot gesehen, wenn der Endpunkt nicht antwortet (etwa falscher Pfad) — der Beweis, dass er den **Endpunkt** prüft, nicht sich selbst | **erfüllt** | Eigener Nachlauf (Messung 6): falscher Pfad → `404`, Exit **2**. Die Implementer-Angabe ist damit **wörtlich** reproduziert (nicht nur übernommen). Ergänzend bindet meine eigene Gegenprobe (Messung 5) die `change_id` unabhängig gegen `cdc.changes`. |
| LP3-K1 | Die Abdeckungstabelle ist durch ein **volles** `make test-integration` real neu erzeugt, trägt den Nachweis; die Zeilen-Drift ist mitgezogen | **erfüllt** | Messung 7: nach meinem vollen Lauf ist die Tabelle **byte-gleich** (`git status` leer) — die committete Fassung ist genau das Generator-Erzeugnis. Die HTTP-Zeile (`:50`) zeigt auf `:2061` und nennt das Lesen; die fünf verschobenen Nachbarzeilen tragen die mitgezogene Drift (`git show 93ac92e`). |
| LP3-K2 | `harness/README.md` nennt den Beleg in der Aufzählung des `make test-integration`-Ziels | **erfüllt** | Die Ziel-Zeile trägt den Lese-Beleg des Wegwerf-Clients samt `· seit slice-087` (Diff in `93ac92e`); der Satz deckt sich mit dem gemessenen Phasen-Ausgang. |
| LP3-K3 | `make gates` grün (Exit direkt, ungepiped) | **erfüllt** | Messung 2: **Exit 0**, fünf innere Gates grün (Zahl 71.90 % statt 72,00 %, beide ≥ 70 % — s. o.). |
| — | Review durchgeführt, Report liegt vor | **erfüllt (vom Reviewer gesetzt)** | `review-slice-087.md` liegt vor; das Häkchen kam im Review-Commit `aede3f4` (`git show aede3f4` ändert nur diese Zeile + legt den Report an). 0 HIGH / 0 MEDIUM, kein Rückgabe-Pfeil. Inhaltlich **nicht** nachgeprüft (Auftrag). |
| — | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **offen (Planner)** | §7 trägt noch Platzhalter; der Slice ist `in-progress`. |
| — | Reconciliation-Register `../reconciliation.md` | **entfällt (trägt)** | `docs/plan/planning/reconciliation.md` existiert real **nicht** — die §2-Zeile sieht genau diesen Entfall vor. |
| — | Beobachtungs-Register fortgeschrieben, **kein** Zähler gesetzt | **offen (Planner)** | Der Lauf steht aus; die §8-Sichtung nennt zwei Einträge (beide real vorhanden, `observations/BEO-PGC/`). |
| — | Jedes Risiko aus §6 trägt einen Ausgang | **offen (Planner)** | Die drei §6-Einträge tragen noch `<bei Closure>`. |
| — | Die drei Paarungen | **korrekt offen, hier nicht zuständig** | Das Repo **hat** Wellen (`welle-20` offen); die Paarungen trägt die nächste Wellen-Closure (Modul 6, Schritt 3c) — die §2-Zeile sagt das. |

**Ergebnis §2:** Alle **drei Liefer-Punkte tragen vollständig** (LP1-K1/K2,
LP2-K1/K2/K3, LP3-K1/K2/K3 bestätigt). Die nicht abgehakten Posten sind die
Closure-Pflichten des **Planners** (§7, Register, Risiko-Ausgänge), kein
Code-/Vertragsdefekt; der Slice ist bewusst noch `in-progress`.

---

## 3. Die eingelöste Fitness-Function-Zeile — prüft der Beleg den Endpunkt?

**Ja — und die Frage ist mit dem eigenen Lauf entschieden, nicht aus dem Text
geschlossen.** `ADR-0081` §Fitness Function führt die Zeile
`make test-integration` / „realer `GET /changes` gegen den laufenden
Feed-Container über den bestehenden HTTP-Wegwerf-Client". Bis `slice-086`
hatte sie keinen Träger (`grep '"/changes' tools/ test/` fand nur
`/changes/stream`). Jetzt:

- **Der Endpunkt wird erreicht, über das Netz** (nicht `docker exec`, kein
  Mock): `docker run --network "$NETWORK" … go run ./tools/harness/httpclient`
  gegen `$HTTP_BASE_URL` (`run-integration-tests.sh:2001-2007`).
- **Die Antwort wird inhaltlich geprüft**, nicht ihr Status: die `READ`-Zeile
  meines Laufs trägt die reale Server-Antwort
  (`new_image={"id":"285","name":"HttpChangesReadE2ESentinel"}`, `operation=INSERT`,
  `commit_position=30901520`) — genau der eingelegte Sentinel.
- **Der Beleg prüft den Endpunkt, nicht den Prozess:** zwei unabhängige
  Gegenproben — der falsche Pfad (`404`, Messung 6) und die fabrizierte
  `change_id` (Messung 5). Die zweite ist die stärkere: sie zeigt, dass die
  `READ`-Zeile **nicht** selbstgenügsam ist, sondern über die unabhängige
  SQL-Bindung an den Store (`run-integration-tests.sh:2043-2053`) hängt.

Damit ist die Fitness-Function-Zeile **eingelöst**: der Beleg existiert, läuft
real grün, wird rot gesehen und bindet die gelesene Änderung gegen den
Lesezugriffsweg `cdc.changes`.

---

## 4. Entscheidungs-Konformität

| ADR | Prüfung | Verdikt |
|---|---|---|
| [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) **§Fitness Function**, Zeile `make test-integration` (realer `GET /changes` gegen den laufenden Container über den Wegwerf-Client) | **eingelöst** — Messung 1 (grün), Messungen 5/6 (rot gesehen), Messung 3/§3 (Endpunkt + `cdc.changes`-Bindung) |
| `ADR-0081` **§Testabdeckung**, letzter Aufzählungspunkt („Erweiterung des E2E-Rundlaufs … über den bestehenden HTTP-Wegwerf-Client") | **eingelöst** — derselbe Wegwerf-Client, derselbe Rundlauf, eine zusätzliche Phase; die ADR führt die Zeile als „Erwartung, keine abschließende Festlegung" | **konform** |
| `ADR-0081` §Re-Evaluierungs-Trigger | nicht ausgelöst: kein Filter jenseits der Achsen, kein Default-Limit, keine zweite Antwortform — die Phase nutzt genau die bestehende Grammatik | **konform** |
| [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) (der bestehende HTTP-Rundlauf, den dieser Slice erweitert) | Die Phase ist **derselbe** HTTP-Rundlauf (`abdeckung_declare "HTTP-API-Rundlauf"` bleibt eine Zeile) und **erweitert** ihn um das Lesen; die zwei Token-Klassen und die Adapter-Platzierung unberührt | **konform** |
| [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md) (ein Client, begrenzte Import-Berechtigung) | Kein zweiter Client; der neue Import `net/url` ist Standardbibliothek unter dem Glob `tooling: ["tools/harness/**"]`; `.a-check.yml` **unverändert**; `make a-check` 0 Befunde | **konform** |
| [`ADR-0030`](../plan/adr/0030-testpyramide.md) (E2E-Tier, kein neues Gate) | `make test-integration` bleibt **kein** Gate; keine `.github/workflows/**`- und keine `Makefile`-Änderung; der Slice fügt dem E2E-Tier eine Phase hinzu, keinem Gate | **konform** |
| [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Digest-Beleg) | **nicht berührt** — `.dockerignore` schließt `tools/**` aus dem Build-Kontext aus (`*` plus nur `cmd/`, `internal/`, `go.mod`, `go.sum`, `tools/coverage-gate.sh`); der Diff ändert ausschließlich `tools/**` und Doku, damit den Build-Kontext **nicht**; kein Digest-Commit nötig | **konform** |

---

## 5. Plan-vs-Code-Diff

**Zuschnitt (§3):** Der Slice nennt drei Träger — `tools/harness/httpclient/main.go`
(update), `tools/harness/run-integration-tests.sh` (update) und
`docs/user/e2e-abdeckung.md` · `harness/README.md` (Erzeugnis/Beleg-Aufzählung).
Der Diff (Messung 9) berührt **genau** diese vier plus den Slice-Plan selbst
(und, außerhalb dieses Slice, die Roadmap-Marker-Commits der Lifecycle-Übergänge).

**„Nicht in dieser Liste" trägt:** keine Datei unter `internal/**`, `spec/**`,
`tools/schema/**`, `.a-check.yml` (Messung 9) — der Endpunkt ist mit `slice-086`
abgenommen, dieser Slice belegt ihn nur.

**Kein stiller Umfangszuwachs:** Die berührten Doku-Dateien tragen den
geänderten Zustand (keine neue Zusage); kein neuer Endpunkt, kein neues
`CDC_*`-Vertragsfeld; `docs/user/benutzerhandbuch.md` unberührt (keine neue
Betreiber-Oberfläche); kein `//nolint`; d-check `hostpaths`/`links`/`structure`
grün (Messung 2).

**Die `set +e`-Form ist konform, kein `§3.9`-Verstoß:** Das Fenster
(`run-integration-tests.sh:2000-2010`) umschließt genau die Zuweisung; die
unbedingte Prüfung `if [ "$http_status" -ne 0 ]` (`:2011`) folgt und meldet
rot. `§3.9` regelt den **Rolleneingriff** auf `make`-Kommandos, nicht `set -e`
in einem Skript — die zitierte Rang-Grenze ist weit, aber nicht falsch
(konform mit dem Review, kein Befund).

---

## 6. Urteil zu den Review-INFOs F-1 und F-2 (und F-3)

**F-1 — „Beleg ohne Bindung an seine Eingabeseite" (INFO).** Die Phase ruft
`GET /changes` in einem Bereich der Breite **1** auf; der Endpunkt kann dort
nur die eine Sentinel-Change liefern, deren `schema`/`table` ohnehin die
gefilterten Werte sind. Der Client **prüft** die Filtertreue
(`main.go:124-129`), aber die Prüfung **kann in dieser Aufruf-Form nicht
feuern** — Mutation C des Reviews (Filterachse entfällt) blieb grün.

- **Berührt es die DoD-Erfüllung?** **Nein.** LP1-K1 verlangt, dass der Client
  die Antwort *inhaltlich* prüft — das tut er (fünf Prüfungen, real
  ausgeführt). Kein DoD-Kriterium verlangt, dass die **Filterwirkung selbst**
  im E2E belegt wird. F-1 ist eine **Grenze der Aussagekraft** des Belegs:
  die DoD-**Evidenz**-Formulierung „die zum Filter passende
  Klartext-Identität" und der Phasen-Satz „mit Filter gelesen" sind stärker
  als das, was der Beleg in dieser Breite binden kann.
- **Wohin gehört es?** In die **Closure §7** als Finding-Klasse („Beleg ohne
  Bindung an seine Eingabeseite") und von dort in den Zähler des
  **Beobachtungs-Registers**. Die Klasse ist die **positive Schwester** von
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (1×, `slice-086`) — ob
  der Lese-Schritt die bestehende Kennung zitiert oder eine neue anlegt,
  entscheidet die Closure, nicht dieser Bericht. **Nicht closure-blockierend.**

**F-2 — „Assertion, die die Aufruf-Form nicht ausüben kann" (INFO).** Die
Ordnungsprüfung vergleicht `commit_position` **zwischen** Einträgen, die
Limit-Angabe wirkt erst oberhalb der Trefferzahl; die Aufruf-Form liefert
`changes=1` bei `limit=100`. Beide Prüfungen können hier nicht feuern; der
Doc-Kommentar nennt nur den ersten der drei Sortierschlüssel.

- **Berührt es die DoD-Erfüllung?** **Nein.** Die DoD nennt dieselben
  Prüfungen als *vorhanden* (LP1-K1), nicht als *ausgeübt*; die Ordnung trägt
  die Unit-Ebene. F-2 ist ebenfalls eine Grenze der Aussagekraft.
- **Wohin gehört es?** Ebenfalls in die **Closure §7** (Finding-Klasse) und
  den Register-Zähler; eine Korrektur des Doc-Kommentars ist Formkosmetik
  (`AGENTS.md` §3.7), kein Vertragsdefekt. **Nicht closure-blockierend.**

**F-3 (LOW) — Anker bricht mitten in der Phrase ab.** Der Anker löst auf
(§1, Messung 6/7) — Lesbarkeit, keine Wirkung. Auch hier: Closure/Register,
nicht der Verifier. Kein DoD-Bezug.

**Zusammengefasst:** F-1/F-2/F-3 sind **Grenzen der Aussagekraft** und
Formulierungen, keine Verhaltens- oder Vertragsfehler; sie berühren **kein**
DoD-Kriterium und gehören in die Closure §7 / das Beobachtungs-Register. Sie
verlangen **keine** Fixrunde — konform mit dem Review.

---

## 7. `AGENTS.md` §3.10

**§3.10 greift nicht.** `git diff --name-only 997a726..HEAD` (Messung 9)
enthält **keine** Datei unter `.github/workflows/`; die drei Slice-Commits
(`93ac92e` Implementer, `aede3f4` Review, plus die Planungs-Commits) berühren
keinen Workflow. Der Slice führt **keinen** neuen Workflow und **keine**
strukturelle Änderung an einem bestehenden ein — die Hard Rule ist nicht
ausgelöst (bestätigt, nicht angenommen).

---

## 8. Was ich nicht prüfen konnte

- **Der Post-Push-Lauf auf dem GitHub-Runner** (`AGENTS.md` §3.10): für
  diesen Slice **nicht einschlägig** — kein Workflow ist berührt (§7). Die
  strukturelle Grenze bleibt gleichwohl bestehen (ein Docker-only-Sensor kann
  sie nicht leisten), nur hat sie hier kein Objekt.
- **Den Inhalt des Review-Reports** habe ich **nicht** nachgeprüft (Auftrag);
  ich habe ihn als Kontext gelesen und nur seinen **Stand** (Häkchen gesetzt,
  kein Rückgabe-Pfeil) bestätigt.
- **Die §6-Risiko-Ausgänge, die Closure-Notiz, das Beobachtungs-Register und
  die drei Paarungen** sind **Planner**-Pflichten beim Übergang nach `done/`;
  sie sind hier bewusst noch offen (der Slice ist `in-progress`).
- **Die vollständige Ordnungs-/Limit-Semantik des Endpunkts** über die
  Aufruf-Form hinaus (F-2): trägt die Unit-Ebene, nicht dieser E2E-Beleg.

---

## 9. Verdikt

**Der Slice ist closure-fähig.** Die drei Liefer-Punkte und ihr Vertrag
tragen vollständig; es fehlen allein die **Planner**-Closure-Pflichten
(§7-Notiz, Register, Risiko-Ausgänge), die naturgemäß erst beim Übergang nach
`done/` entstehen.

Getragen ist alles, was der Slice zusagt:

- **Die Fitness-Function-Zeile ist eingelöst:** realer `GET /changes` gegen
  den laufenden Feed-Container über den bestehenden Wegwerf-Client, grün
  gemessen (Exit 0) und **zweifach rot gesehen** (falscher Pfad `404`;
  fabrizierte `change_id`);
- **der Beleg prüft den Endpunkt, nicht sich selbst:** die `READ`-Zeile
  hängt über die unabhängige `cdc.changes`-Bindung am Store (eigene
  Gegenprobe, Messung 5) — nicht über den `make test-integration`-Prozess;
- **ein Client bleibt ein Client** (`ADR-0068`), `internal/**`/`spec/**`/
  `tools/schema/**`/`.a-check.yml` unberührt;
- **die Abdeckungstabelle ist das echte Generator-Erzeugnis:** ein voller
  Lauf schreibt **nichts** (byte-gleich, Messung 7);
- **`make gates` grün** (Exit 0); `doc-commits`/`doc-immutable` über die
  Slice-Range grün;
- **`§3.10` nicht ausgelöst** — kein Workflow berührt.

**Was offen bleibt:** (1) die Planner-Closure-Posten; (2) F-1/F-2/F-3 als
Finding-Klassen der Closure §7 / des Registers — **nicht
closure-blockierend**, keine Fixrunde.

---

**Übergabe:** Dieser Bericht geht an den **Planner** (Closure-Entscheidung;
F-1/F-2/F-3 als Finding-Klassen in §7 und den Register-Zähler) und als
Entlastung an den **Implementer** (der Beleg ist nachweislich tragend und
zweifach rot gesehen); an den **Validator** nur, falls dieser Slice als
MVP-Slice validiert wird — der reale Bedarf liegt hier repo-extern nicht vor.
Der Bericht ist ein **Lauf-Beleg** (dieser Stand, diese Läufe, dieses Verdikt)
und ersetzt weder das Review noch die Planner-Closure.
