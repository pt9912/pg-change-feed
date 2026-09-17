# Slice slice-101: NATS-Client in C# und Kotlin

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf Nutzerentscheidung, ein einzelner Slice
ohne Closure-Bedingung jenseits seiner DoD (Modul 6 §Wann Arbeit eine Welle
braucht). Dieses Repo führt derzeit **keine** offene Welle
(`docs/plan/planning/in-progress/roadmap.md` §Offene Wellen ist leer).

**Bezug:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
(Festlegung 1 — Zellen `csharp`×`nats` und `kotlin`×`nats`, „eine gepinnte
öffentliche Client-Bibliothek je Sprache — gemessen existent"; Festlegung 4 —
`NATS.Net` (NuGet) und `io.nats:jnats` (Maven Central), Registry-Existenz
gemessen 2026-09-17; §Slice-Schnitt-Empfehlung, Zeile 4: „**NATS-Client in C#
und Kotlin** … Abhängigkeit 1, 2") · [`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md)
und [`ADR-0056`](../../adr/0056-nats-tabellen-granulares-subjekt.md) (das Wecksignal und
sein Subjekt-Schema, dessen Form der Client anspricht) ·
[`ADR-0079`](../../adr/0079-nats-beispielclient-vierter-examples-client.md)
(Form-Vorbild `examples/nats-client` in Go, samt zweiseitigem Ablauf:
lauschen, dann über `GET /changes` holen) · `BEO-PGC/handbuch-nicht-
nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**, verkörpert) und
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (**3×**, verkörpert) —
beide treffen den Handbuch-Nachzug dieses Slice.

**Berührte Spec-Stellen:** `LH-FA-SST-007` (das NATS-Wecksignal) — dieser
Slice **zeigt** es in zwei weiteren Sprachen, er ändert es nicht.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `examples/csharp/nats-client/` und `examples/kotlin/nats-client/`
entstehen — beide lauschen auf das tabellen-granulare Subjekt-Schema
`cdc.changes.<source_id>.<schema>.<table>`, geben das leere Payload-Ereignis
aus und holen anschließend die tatsächliche Änderung über `GET /changes`
(zweiseitiger Ablauf, Form-Vorbild `examples/nats-client` in Go). Jede
Sprache bekommt eine gepinnte, öffentliche Client-Bibliothek — C#: `NATS.Net`
(NuGet), Kotlin/JVM: `io.nats:jnats` (Maven Central); Existenz beider
2026-09-17 gemessen ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
Festlegung 4, konkrete Version/Digest-Pin ist Sache dieses Zuges). Beide
Handbuch-Zeilen (der vierte Zugriffs-Abschnitt, der einen `**Beispiele:**`-
Block bekommt) werden nachgezogen.

**Warum beide Sprachen in einem Slice:** Beide Programme teilen dieselbe
Handbuch-Form und denselben zweiseitigen Ablauf; die einzige neue Zutat je
Sprache ist eine gepinnte Bibliothek — kein Codegenerator, kein neuer
Bau-Kontext. Zwei Programme, zwei gepinnte Abhängigkeiten, zwei
Handbuch-Zeilen bleiben bei drei Liefer-Punkten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der gRPC-Client in C#/Kotlin.** Folge-Slices `slice-102` (C#) und
  `slice-103` (Kotlin) übernehmen ihn — er trägt den benannten Zusatzkontext
  und den Protobuf-/gRPC-Generator, zwei Form-Fragen, die dieser Slice nicht
  öffnet.
- **Der Go-NATS-Client (`examples/nats-client`).** Er existiert bereits
  (`slice-083`) und ist das **Form-Vorbild**; ihn umzubauen wäre Arbeit an
  Bestand ohne Adresse.
- **Eine Änderung am Subjekt-Schema oder der Zustellsemantik des
  Wecksignals.** Der Client **benutzt** `cdc.changes.<source_id>.<schema>.
  <table>`; verlangt er eine Vertragsänderung, ist das eine Spec-Änderung,
  kein Beispiel-Umbau (`SPEC-023`, `SPEC-017`).
- **Eine Zustandsmaschine (Reconnect, Deduplizierung, Rückstand-Tracking).**
  [`SPEC-023`](../../../../spec/pflichtenheft.md) schließt das für **jeden**
  Beispiel-Client aus — Vorbild, kein Belegträger.
- **`.a-check.yml`.** Unverändert aus denselben Gründen wie in den
  vorangegangenen Sprach-Slices: keine C#-/Kotlin-Schicht in diesem Repo.
  Bestand bleibt bewusst stehen.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1 — der C#-NATS-Client samt Pin.** `NATS.Net` gepinnt im
      Projekt-/Paket-Manifest von `examples/csharp/`, mit Leser-Zeile
      („wofür sie da ist"); `examples/csharp/nats-client/` lauscht real,
      gibt das Ereignis aus, holt die Änderung über `GET /changes`; netzlos
      prüfbare Teile (Subjekt-Aufbau) sind getestet und laufen über
      `examples-csharp`.
- [x] **LP2 — der Kotlin-NATS-Client samt Pin.** `io.nats:jnats` gepinnt im
      Build-Manifest von `examples/kotlin/`, mit Leser-Zeile;
      `examples/kotlin/nats-client/` — dieselbe Zusage, über
      `examples-kotlin`.
- [x] **LP3 — die zwei Handbuch-Zeilen.** `docs/user/benutzerhandbuch.md` §4
      „Zugriff über das NATS-Wecksignal": der bestehende `**Beispiele:**`-
      Block (Go seit `slice-095`) bekommt zwei weitere Zeilen (C#, Kotlin),
      samt Änderungshistorie-Zeile.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      Bericht: `docs/reviews/review-slice-101.md` (1 HIGH, 1 MEDIUM), über
      den Konflikt-Pfad (Modul 8) an
      `docs/reviews/architect-verdict-slice-101-jnats-bouncycastle.md`
      übergeben — beide zusammen sind die vollständige Review-Übergabe
      dieses Slice.
- [x] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-101.md`
      liegt vor (Modul 11, frischer Kontext) — DoD-konform, mit zwei
      benannten, nicht-blockierenden Einschränkungen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) — *entfällt: Greenfield-
      Bootstrap (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = GF),
      keine Datei vorhanden.*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im
      Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/csharp/<Paket-Manifest>` | update | Pin von `NATS.Net`, mit Leser-Zeile. |
| `examples/csharp/nats-client/**` | neu | NATS-Client, zweiseitiger Ablauf, Form-Vorbild `examples/nats-client` (Go). |
| `examples/kotlin/<Build-Manifest>` | update | Pin von `io.nats:jnats`, mit Leser-Zeile. |
| `examples/kotlin/nats-client/**` | neu | NATS-Client, zweiseitiger Ablauf. |
| `docs/user/benutzerhandbuch.md` | update | §4 „Zugriff über das NATS-Wecksignal": zwei weitere Zeilen (C#, Kotlin) im `**Beispiele:**`-Block; Änderungshistorie-Zeile. |
| `examples/csharp/Dockerfile` | update | **Plan-Nachzug:** drittes Programm in der gemeinsamen `build`-Stufe (restore/build/test/publish) plus neue Stufe `runtime-nats` — dieselbe Form wie `runtime-sse`, notwendig, damit der Pin überhaupt einen Bau-Beleg hat. |
| `examples/kotlin/settings.gradle.kts` | update | **Plan-Nachzug:** `include("nats-client")` — ohne das Modul ist der Pin im `build.gradle.kts` des neuen Moduls für Gradle unsichtbar. |
| `examples/kotlin/Dockerfile` | update | **Plan-Nachzug:** drittes Modul in der gemeinsamen `build`-Stufe (test/installDist) plus neue Stufe `runtime-nats`, dieselbe Form wie `runtime-sse`. |
| `harness/mk/examples.mk` | update | **Plan-Nachzug:** dritter `docker build --target runtime-nats`-Aufruf je Sprachziel — ohne ihn baut `make examples-csharp`/`make examples-kotlin` das dritte Programm zwar mit, aber taggt kein eigenes Image dafür. |
| `harness/README.md` | update | **Plan-Nachzug:** Sensors-Zeilen `make examples-csharp`/`make examples-kotlin` von „zwei Images“ auf „drei Images“ gezogen (§3.13-Suchlauf-Fund, siehe §7). |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-098` **und** `slice-099` liegen in
`done/` — die C#- und Kotlin-Sprachwurzel samt Werkzeugkette und `make`-Ziel
existieren, ohne sie ist kein Bau-Kontext für den NATS-Client vorhanden. Ohne
Rückfrage feststellbar (Verzeichnis-Position der beiden Vorgänger-Slices).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn eine der
  beiden gepinnten Client-Bibliotheken transitiv eine zweite, gepinnte
  Abhängigkeit erzwingt, die eine eigene Bewertung braucht (Lizenz,
  Sicherheits-Historie) — dann ist die Größenannahme („eine Bibliothek je
  Sprache") falsch.
- `in-progress` → `open` (blockiert — Carveout?): wenn `NATS.Net` oder
  `io.nats:jnats` zum Zeitpunkt des Baus nicht mehr auflösbar ist (Registry-
  Digest/Paket zurückgezogen) oder eine der beiden Bibliotheken eine
  nicht-öffentliche Quelle voraussetzt
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — beide Clients
lauschen real auf das Subjekt-Schema und holen die Änderung über `GET
/changes`, beide sind über ihr jeweiliges Sprachziel kompiliert und getestet,
die Handbuch-Zeilen tragen beide Sprachen — **und** `make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegender Kandidat: ob die am 2026-09-17 gemessenen Bibliotheks-
Kandidaten (`NATS.Net`, `io.nats:jnats`) beim tatsächlichen Bau ohne
Zusatz-Abhängigkeit auskommen — der Lauf trägt den realen Pin nach.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Registry-/Digest-Pin einer der beiden Bibliotheken ist zum Zeitpunkt
  des Baus nicht mehr auflösbar** (Paket zurückgezogen, Version gelöscht). —
  **Ausgang: entfallen.** Beide `--no-cache`-Bauproben des Verifiers lösten
  real gegen NuGet (`NATS.Net`) bzw. Maven Central (`io.nats:jnats`) auf
  (`docs/reviews/verify-slice-101.md` #4/#6, #11/#12) — beide Pakete sind
  zusätzlich real die neuesten stabilen Versionen.
- **Eine der beiden Werkzeugketten verlangt eine nicht-öffentliche Quelle**
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). — **Ausgang: entfallen.** Beide Bauten
  liefen ausschließlich gegen öffentliche Registries (NuGet, Maven
  Central), ohne Zugangsdaten (`docs/reviews/verify-slice-101.md` #4, #6,
  §5).
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (ein
  drittes Bau-Ziel je Sprache im selben Workflow) — §Re-Evaluierungs-
  Trigger 4, `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×,
  offen), `BEO-PGC/github-actions-unverifizierbar-lokal` (7×, verkörpert in
  `AGENTS.md` §3.10). — **Ausgang: weiter offen**, ohne neue Evidenz in
  einem der beiden zitierten Register-Einträge. Der Verifier hat gemessen
  (`docs/reviews/verify-slice-101.md` #14): `git diff 03337c9..HEAD --stat
  -- .github/workflows/` ist **leer** — dieser Slice ändert
  `.github/workflows/examples.yml` nicht strukturell, `AGENTS.md` §3.10
  wird nicht neu ausgelöst (wie schon bei `slice-100`). Die oben zitierte
  Zahl „7×" ist bereits die vom Verifier korrigierte — der Slice-Kopf/§8
  dieses Plans zitierte ursprünglich „5×", real bereits 7× zum
  Niederschrift-Zeitpunkt (`docs/reviews/verify-slice-101.md` §5); Korrektur
  siehe fünfter Risiko-Punkt unten und §8. `nicht-blockierender-workflow-
  alarmmuedigkeit` bleibt bei 1×.
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen; sie ist der einzige Teil dieses Slice, den **kein**
  Kompilat erzwingt. — **Ausgang: entfallen.** Reviewer (Negativbefund) und
  Verifier (#17, direkte Lektüre) bestätigen unabhängig den
  `**Beispiele:**`-Block mit drei Zeilen (Go/C#/Kotlin) und die
  Änderungshistorie-Zeile 1.23. Beide Registereinträge
  (`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  `handbuch-versionshistorie-uebersprungen`) sind **eingelöst**, nicht
  verletzt — keine neue Evidenzdatei, beide bleiben bei 3×.
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, verkörpert in `AGENTS.md`
  §3.13 — Suchlauf-Pflicht). — **Ausgang: eingetreten, direkt behoben — und
  ein zweiter, davon unterschiedener Fund derselben Closure.** (a) Der
  Implementer-eigene §3.13-Suchlauf fand eine stale Aussage in
  `harness/README.md` §Sensors („zwei Images" → „drei Images" je
  Sprachziel, ausgelöst durch das dritte Runtime-Image je Sprache) und
  korrigierte sie im selben Commit (`79dbd5d`). Reviewer und Verifier
  bestätigen die Korrektur als vollständig
  (`docs/reviews/review-slice-101.md` Negativbefund ·
  `docs/reviews/verify-slice-101.md` #18). Neue Evidenzdatei
  `evidence/slice-101.md` bei `BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
  Zähler **7× → 8×** — kein neuer Schwellen-Übertritt, die Regel
  (`AGENTS.md` §3.13) steht bereits seit `welle-20`. (b) Ein zweiter, davon
  unabhängiger Fund gehört **nicht** zu diesem Register-Eintrag, sondern zu
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`: Der Slice-Kopf/§8
  dieses Plans zitierten `github-actions-unverifizierbar-lokal` mit „5×" —
  real bereits 7× zum Zeitpunkt der Plan-Niederschrift, vom Verifier
  gefunden (`docs/reviews/verify-slice-101.md` §5). Anders als bei (a)
  überholt hier nicht die Arbeit dieses Slice einen fremden, bei
  Niederschrift noch korrekten Träger — die Zahl war bereits **bei ihrer
  eigenen Niederschrift** veraltet, keine Drift durch diesen Slice selbst
  (dieselbe Unterscheidung, die der dortige Register-Eintrag selbst
  zieht). Neue Evidenzdatei `evidence/slice-101.md` bei
  `zahl-in-traeger-driftet-gegen-die-messung`, Zähler **9× → 10×** (bereits
  verkörpert, kein neuer Schwellen-Übertritt). Beide Fundstellen im Plan
  (dieser Risiko-Punkt und §8) sind im Rahmen dieser Closure auf „7×"
  korrigiert.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Die Rollen-Sequenz aus Modul 8 hat einen realen
  Prozessfehler aufgefangen, ohne dass am Ende etwas inhaltlich Falsches im
  Repo landete: Reviewer eskalierte F-1 (HIGH) korrekt über den
  Konflikt-Pfad statt als stille Fixrunde am Implementer
  (`docs/reviews/review-slice-101.md`); der Architect prüfte in eigenem
  Kontext eigenständig nach — mit größerer Tiefe als der Review selbst
  (drei zusätzliche Advisories über NVD gefunden,
  `docs/reviews/architect-verdict-slice-101-jnats-bouncycastle.md`); der
  Verifier hat die Architect-Zahlen ein drittes Mal, unabhängig
  reproduziert (`docs/reviews/verify-slice-101.md` #7–#9) und bestätigt.
  Beide Sprachen bauen und testen real — der Verifier fuhr zusätzlich zwei
  eigene, cache-lose `docker build --no-cache`-Bauproben, die alle drei
  Programme je Sprache real neu bauen und testen.
- **Was ging anders als geplant:** Der Implementer erkannte den in §4
  vorab deterministisch benannten Rückführungs-Trigger (transitive
  `bcprov-lts8on`-Abhängigkeit über `io.nats:jnats`), bewertete ihn aber
  selbst — statt den Konflikt-Pfad auszulösen — und belegte dabei nur die
  Lizenz-Hälfte vollständig, die Sicherheits-Historie-Hälfte nur als
  unbelegte Reputationsaussage (Review F-1 HIGH, F-2 MEDIUM). Der Architect
  hat den Trigger nachträglich vollständig geprüft (fünf Advisories über
  vier unabhängige Quellen, alle vor oder exakt bei der gepinnten Version
  `2.73.12.1` behoben — die Version **ist** sogar der Fix-Commit für die
  jüngste, `CVE-2026-15997`) und „Fortsetzen ohne Rückführung" verfügt: Die
  geforderte Bewertung trägt inhaltlich, der Prozessfehler bleibt trotzdem
  eigenständig bestehen und wird nicht durch das gute Ergebnis geheilt.
- **Steering-Loop-Eintrag:** `BEO-PGC/vorab-bedingung-nach-umsetzung-
  geprueft` erreicht mit diesem Slice **2×** (weiterhin unter der
  3×-Schwelle — kein Ausgang fällig, bleibt `offen`) — zweite, unabhängige
  Instanz derselben Klasse wie `slice-073`: eine §4-Bedingung, deren
  Klärung vor die Umsetzung gehört, wird erst nach dem Schreiben des
  Artefakts ausgewertet. Anders als bei `slice-073` deckte die
  nachträgliche Prüfung hier keine Diskrepanz auf, sondern bestätigte das
  bereits (unvollständig belegt) behauptete Ergebnis — die Klasse selbst
  ist trotzdem identisch. Zusätzlich benannt, noch nicht verkörpert: Der
  Architect hat Modul 8 §Konflikt-Pfad — dem Wortlaut nach für
  ADR-Konflikte geschrieben — hier erstmals explizit auf einen
  Plan-Trigger-Konflikt übertragen und die Übertragung begründet; der
  Verifier hat das als plausible, aber nicht zwingende
  Auslegungsentscheidung benannt (`docs/reviews/verify-slice-101.md` §3.1).
  Beides wird hier festgehalten, damit ein dritter Fall — auf
  `vorab-bedingung-nach-umsetzung-geprueft`, oder eine zweite Anwendung der
  Konflikt-Pfad-Übertragung — nicht erneut bei null herleiten muss.
- **Beobachtungs-Register (`../observations/`):** Drei Einträge
  fortgeschrieben: `arbeit-ueberholt-stehenden-traeger` jetzt **8×**
  (`evidence/slice-101.md`, `harness/README.md`-Zahlenwort-Fund, bereits
  verkörpert, kein neuer Schwellen-Übertritt); `vorab-bedingung-nach-
  umsetzung-geprueft` jetzt **2×** (`evidence/slice-101.md`, Rollen-Verstoß
  aus Review F-1, weiterhin `offen`, unter der Schwelle);
  `zahl-in-traeger-driftet-gegen-die-messung` jetzt **10×**
  (`evidence/slice-101.md`, die im eigenen Plan-Kopf/§8 veraltet zitierte
  „5×" für `github-actions-unverifizierbar-lokal`, bereits verkörpert,
  kein neuer Schwellen-Übertritt). Zwei Einträge eingelöst, keine neue
  Evidenzdatei: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (bleibt 3×) und `handbuch-versionshistorie-uebersprungen` (bleibt 3×).
  Zwei Einträge bewusst **nicht** fortgeschrieben, obwohl im
  Slice-Kopf/§8 zitiert: `github-actions-unverifizierbar-lokal` (bleibt
  7×, real bereits korrekt gezählt, nur der Plan zitierte veraltet) und
  `nicht-blockierender-workflow-alarmmuedigkeit` (bleibt 1×) — Begründung
  in §6, Risiko 3: dieser Slice ändert `.github/workflows/examples.yml`
  nicht strukturell, `AGENTS.md` §3.10 wird nicht neu ausgelöst.
- **Folge-Slices:** keine neuen. `slice-102`/`-103` existierten bereits vor
  dieser Closure in `open/`.
- **Risiken aus §6:** fünf Zeilen, fünf Ausgänge — drei **entfallen**
  (Registry-/Digest-Pin gemessen auflösbar; keine nicht-öffentliche
  Quelle, da beide Bauten ausschließlich öffentliche Registries nutzten;
  Handbuch-Nachzug gemessen vollständig), eines **weiter offen**
  (nicht-blockierender Workflow — ohne neue Registerevidenz, da dieser
  Slice den Workflow nicht strukturell ändert) und eines **eingetreten,
  direkt behoben, mit einem zweiten, davon unterschiedenen Fund**
  (`harness/README.md`-Zahlenwort-Fund →
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, jetzt 8×; eigene, im
  Plan-Kopf/§8 veraltet zitierte Zahl →
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, jetzt 10×). Details
  je Zeile in §6.
- **Drei Paarungen:** Anker-Paarung entfällt (keine neue Verkörperung durch
  diesen Slice — beide fortgeschriebenen, bereits verkörperten
  Register-Einträge waren schon vor diesem Slice verkörpert;
  `vorab-bedingung-nach-umsetzung-geprueft` bleibt unter der Schwelle).
  Folge-Slice-Paarung entfällt (keine neuen Folge-Slices; `slice-102`/`-103`
  existierten bereits vor dieser Closure in `open/`). Register-Paarung
  **grün** — alle sieben zitierten Beobachtungs-Verzeichnisse
  (`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  `handbuch-versionshistorie-uebersprungen`,
  `nicht-blockierender-workflow-alarmmuedigkeit`,
  `github-actions-unverifizierbar-lokal`,
  `arbeit-ueberholt-stehenden-traeger`,
  `vorab-bedingung-nach-umsetzung-geprueft`,
  `zahl-in-traeger-driftet-gegen-die-messung`) existieren mit nicht leerem
  `evidence/`. Letztes DoD-Häkchen wird gegen den mv-Commit bestätigt
  (Commit 3).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `examples/csharp/**`,
`examples/kotlin/**` und `docs/user/benutzerhandbuch.md` — die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Keine zu
grobe Sub-Area zu differenzieren.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Vier
Treffer: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide als LP3 in die DoD gezogen; `nicht-blockierender-
workflow-alarmmuedigkeit` (**1×**, offen) und `github-actions-
unverifizierbar-lokal` (**7×** — bei Niederschrift dieses Plans fälschlich
als „5×" zitiert, von `docs/reviews/verify-slice-101.md` §5 korrigiert und
bei dieser Closure nachgezogen, siehe §6 letzter Risiko-Punkt sowie
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-101.md`;
verkörpert in `AGENTS.md` §3.10) — beide treffen den weiter wachsenden
Workflow, als Risiko in §6. Kein Treffer zu NATS.Net/jnats/NuGet/Maven
selbst (gemessen: `grep -rli "nats\.net\|jnats\|nuget\|maven"
docs/plan/planning/observations/` → kein Fund).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**.
Der Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

### Sub-Area: `*` (Default, PGC)

- **Modus:** GF
- **Konventionen-Dichte:** `harness/conventions.md` Modus-Deklaration setzt
  GF für das gesamte Repo (Doc führt, Code folgt).
- **Phase-Reife:** Phase 5 (etabliert) — beide Sprach-Wurzeln existieren
  bereits (`slice-098`, `slice-099`); dieser Slice fügt ein drittes Programm
  je Sprache samt gepinnter Bibliothek hinzu.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — GF, Doc führt; kein Bestand, der
  von der neuen Form abweichen könnte.
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bootstrap).
