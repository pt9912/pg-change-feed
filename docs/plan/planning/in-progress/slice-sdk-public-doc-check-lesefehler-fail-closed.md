# Slice sdk-public-doc-check-lesefehler-fail-closed: Das SDK-Gate meldet einen Lesefehler als Exit 2 statt „keine interne Kennung“

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (offizielle
Client-Bibliotheken; ihre Texte erreichen Anwender über die Pakete, Haupt-Bezug),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (`Accepted`,
unberührbar; das Gate, das hier verschärft wird),
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(`Accepted`; das Schwester-Gate `handbuch-public-doc-check` ist fail-closed bei
Lesefehlern, die ADR nennt es nicht). Anlass: die Verifikation von
`handbuch-public-doc-check-gate-und-skill` ([Report](../../../reviews/verifikation-slice-handbuch-public-doc-check-gate-und-skill.md)
§7 Punkt 2) maß das SDK-Gate als fail-open.

**Berührte Spec-Stellen:** keine.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, auf Auftrag des Auftraggebers (Folge-Slice der Closure von
`handbuch-public-doc-check-gate-und-skill`). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (vom Verifier gemessen, 2026-10-02):** In einer Kopie des `sdks`-Baums
liegt eine `README.md` mit einer internen Kennung und `chmod 000`. `tools/harness/sdk-public-doc-check.sh`
druckt die `grep`-Meldung „Keine Berechtigung“ auf stderr, meldet aber
`sdk-public-doc-check: keine interne Kennung unter …` und endet mit **Exit 0**. Ursache
(am Stand `6a8c9c51` gelesen): `xargs -0 -r grep … || true` verschluckt den `grep`-Exit 2, und
`find` läuft in einer Pipe, deren Fehler nicht ausgewertet wird. Das Gate ist Teil von
`make gates` und Vorstufe der drei `make sdk-pack-*`-Ziele; ein Lesefehler lässt dort eine
Kennung unbemerkt in ein Paket gelangen.

**Ziel:** Ein Lesefehler (nicht lesbare Datei, nicht lesbares Verzeichnis) endet mit Exit 2
und einer Meldung, nie mit „keine interne Kennung“. Muster ist das Schwester-Skript
`tools/harness/handbuch-public-doc-check.sh`: Scratch-Verzeichnis mit Trap, `find` vor der
Schleife in eine Datei, `grep -a` mit Auswertung (Exit 1 = kein Treffer = Erfolg, Exit ≥ 2 =
Lesefehler = Exit 2). Der Tabellentest `make test-sdk-public-doc-check` bekommt die
Lesefehler-Fälle (nicht lesbare Datei, nicht lesbares Unterverzeichnis, Datei mit NUL-Byte
und Kennung; nur bei uid ≠ 0, wie im Handbuch-Tabellentest) und bindet den Meldungstext. Der
Sensor-Vertrag `harness/sensors/sdk-public-doc-check.md` trägt den neuen Ausgang (Exit 2) in
der Ausgangstabelle.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung von Muster oder Prune-Liste** — die Verschärfung betrifft allein die
  Lesefehler-Semantik; Muster und Umfang sind Festlegungen von `ADR-0134`.
- **Eine Änderung am Handbuch-Gate** — es ist bereits fail-closed (Verifikation, sieben
  Mutationen rot).
- **Eine Aufnahme des Tabellentests in `make gates`** — `ADR-0134` Teilfrage 2 hält den
  Tabellentest bewusst als Werkzeug.
- **Eine Änderung an `.github/workflows/`** — das Gate läuft unverändert mit `make gates`.

## 2. Definition of Done

- [ ] **Liefer-Punkt 1 — Das Skript.** `tools/harness/sdk-public-doc-check.sh` endet bei
      einem Lesefehler mit Exit 2 und einer Meldung auf stderr; Exit 0 und 1 behalten ihre
      Bedeutung. *Zu belegen durch:* `git grep -n -F '|| true' -- tools/harness/sdk-public-doc-check.sh`
      ist leer; `make sdk-public-doc-check` Exit 0 am echten Baum (gedruckte Zeile im Bericht);
      Mutationsprobe auf einer Kopie im Scratchpad (Rücknahme der Lesefehler-Behandlung
      auf `|| true`): der Tabellentest wird rot — Stelle, Instanz und Farbe im Bericht
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
- [ ] **Liefer-Punkt 2 — Der Tabellentest.** `tools/harness/run-sdk-public-doc-check-tests.sh`
      deckt nicht lesbare Datei, nicht lesbares Unterverzeichnis und NUL-Byte-Datei, jeweils
      mit erwartetem Exit und gebundenem Meldungstext; die Lesefehler-Fälle laufen nur bei
      uid ≠ 0 und melden das Überspringen. *Zu belegen durch:* `make test-sdk-public-doc-check`
      Exit 0 mit gedruckter Schlusszeile; je Fall eine Mutation am Skript (Kopie), die den Fall rot färbt.
- [ ] **Liefer-Punkt 3 — Der Sensor-Vertrag.** `harness/sensors/sdk-public-doc-check.md`
      nennt Exit 2 (Lesefehler, fail-closed) in der Ausgangstabelle und die Grenze der
      uid-0-Fälle des Tabellentests. *Zu belegen durch:* `git diff` der Datei;
      `git grep -n -F Lesefehler -- harness/sensors/sdk-public-doc-check.md` hat mindestens einen Treffer.
- [ ] **Nur diese Pfade.** Der Diff berührt `docs/plan/adr/0134-*`, `.github/` und `AGENTS.md` nicht.
      *Zu belegen durch:* `git diff --name-only <Basis> -- docs/plan/adr .github AGENTS.md` ist leer.
- [ ] `make gates` grün — Exit-Code ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report `docs/reviews/review-slice-sdk-public-doc-check-lesefehler-fail-closed.md` liegt vor (`.harness/skills/reviewer.md`),
      kein offenes HIGH/MEDIUM — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow
      ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review.
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes, beide
      Stände gemessen; `make suchlauf-nachmessen PLAN=<diese Datei>` endet mit Exit 0.
- [ ] Doku-Update: Sensor-Vertrag (Liefer-Punkt 3); `harness/README.md` §Sensors, Zeile
      `make sdk-public-doc-check`, nennt den Ausgang, falls sie Ausgänge aufzählt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert (Nachbar: `BEO-PGC/intern-kennungen-in-ausgelieferten-texten`).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

**Umfang:** S — Schätzung, nicht gemessen: ein Skript (rund 20 Zeilen mehr), drei bis vier
Tabellenfälle, ein Abschnitt im Sensor-Vertrag.

**Voraussetzung:** Die Architect-Frage in §4 ist beantwortet (ob eine Folge-ADR nötig ist).
Das Handbuch-Gate und sein Tabellentest liegen vor und sind das Muster.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/sdk-public-doc-check.sh` | update | Lesefehler → Exit 2 nach dem Muster von `handbuch-public-doc-check.sh` |
| `tools/harness/run-sdk-public-doc-check-tests.sh` | update | Lesefehler-Fälle (nur uid ≠ 0), Meldungsbindung; Aufräumen mit `chmod -R u+rwx` vor `rm -rf` |
| `harness/sensors/sdk-public-doc-check.md` | update | Ausgang Exit 2, Grenze der uid-0-Fälle |
| `harness/README.md` | prüfen | Zeile `make sdk-public-doc-check` zählt Ausgänge nicht auf (am Start gelesen: unverändert) |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „das SDK-Gate ist fail-open bei
Lesefehlern“; Symbolnamen `|| true`, `Lesefehler`; Hedge `fail-`. Stand `6a8c9c51`, gemessen
2026-10-02; die `diff`-Zeilen und die Befunde trägt der Implementer nach):**

```suchlauf
6a8c9c51 1 -n -F '|| true' -- tools/harness/sdk-public-doc-check.sh
6a8c9c51 0 -n -F '|| true' -- tools/harness/handbuch-public-doc-check.sh
6a8c9c51 0 -n -F 'Lesefehler' -- tools/harness/sdk-public-doc-check.sh harness/sensors/sdk-public-doc-check.md tools/harness/run-sdk-public-doc-check-tests.sh
6a8c9c51 4 -n -F 'Lesefehler' -- tools/harness/handbuch-public-doc-check.sh
6a8c9c51 2 -n -F 'Lesefehler' -- harness/sensors/handbuch-public-doc-check.md
6a8c9c51 0 -n -F 'chmod' -- tools/harness/run-sdk-public-doc-check-tests.sh
6a8c9c51 1 -n -E 'fail-closed|fail-open' -- docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md
```

Die `diff`-Zeilen (Arbeitsbaum nach der Umsetzung, gemessen 2026-10-03; die Plan-Datei ist
aus dem Suchraum ausgeschlossen):

```suchlauf
diff 0 -n -F '|| true' -- tools/harness/sdk-public-doc-check.sh
diff 0 -n -F '|| true' -- tools/harness/handbuch-public-doc-check.sh
diff 14 -n -F 'Lesefehler' -- tools/harness/sdk-public-doc-check.sh harness/sensors/sdk-public-doc-check.md tools/harness/run-sdk-public-doc-check-tests.sh
diff 4 -n -F 'Lesefehler' -- tools/harness/handbuch-public-doc-check.sh
diff 2 -n -F 'Lesefehler' -- harness/sensors/handbuch-public-doc-check.md
diff 6 -n -F 'chmod' -- tools/harness/run-sdk-public-doc-check-tests.sh
diff 1 -n -E 'fail-closed|fail-open' -- docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md
```

**Befund am Diff (Implementer, gemessen):** Gefunden: der Fallback `|| true` im SDK-Skript ist
entfernt (1 → 0); `Lesefehler` steht in Skript (4), Sensor-Vertrag (4) und Tabellentest (6);
`chmod` im Tabellentest 0 → 6. Nicht gefunden: kein weiterer Träger der Aussage „SDK-Gate
fail-open/Lesefehler“ außerhalb dieser Dateien; `harness/README.md` zählt die Ausgänge des
Gates nicht auf (keine Änderung); `ADR-0134` unverändert (`fail-closed` weiter 1 Zeile, die
Prune-Liste); das Handbuch-Gate unverändert. Die Verschärfung ist ohne Folge-ADR festgelegt
(Entscheidung des Hauptlaufs auf Delegation des Auftraggebers zur Architect-Frage in §4: Verschärfung,
keine Senkung nach `AGENTS.md` §3.6; Vertrag im Sensor-Vertrag).

| Träger | Messung am Stand `6a8c9c51` | Behandlung (Befund am Diff trägt der Implementer ein) |
|---|---|---|
| Fallback `true` hinter dem `grep` im SDK-Skript (Zeile 1) | 1 Trefferzeile (der verschluckte `grep`-Exit) | Soll am Diff 0 (Zeile ergänzt der Implementer) |
| `Lesefehler` im SDK-Skript, Tabellentest und Sensor-Vertrag (Zeile 3) | 0 | Soll am Diff mindestens je ein Treffer in Skript und Sensor-Vertrag; der Implementer trägt die Zahl nach |
| [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) `fail-closed` (Zeile 7) | 1 Trefferzeile; sie betrifft die Prune-Liste (Festlegung 4), nicht Lesefehler | unberührbar (§3.5); Soll am Diff weiter 1 |
| Handbuch-Gate (Zeilen 2, 4, 5) | Fallback 0, `Lesefehler` 4 und 2 | Muster, bleibt unverändert |

## 4. Trigger

**Start** (`next` → `in-progress`): Die Architect-Frage unten ist beantwortet; kein anderer
Slice liegt in `in-progress/` (WIP-Limit 1).

**Architect-Frage (nicht entschieden, zur Beantwortung vor dem Start):** `ADR-0134` ist
`Accepted` und nach [`AGENTS.md`](../../../../AGENTS.md) §3.5 unberührbar. Die
Lesefehler-Semantik (Exit 2 statt stillem Grün) ist eine **Verschärfung**: das Gate wird
strenger, nicht lockerer; §3.6 betrifft nur Senkungen. Zu klären: (a) Reicht die Festlegung im
Sensor-Vertrag `harness/sensors/sdk-public-doc-check.md` (wie beim Handbuch-Gate, wo
die Verifikation „Abweichung mit Deckung im Zweck, ohne Wortlaut in der ADR“ nannte), oder
(b) braucht es eine Folge-ADR (`Schärft: ADR-0134`), weil ein neuer Exit-Code ein Vertragsmerkmal
ist und `ADR-0143` Exit 2 nur für Klassifikationsfehler nennt? Die Frage 1 des
[Review-Reports](../../../reviews/review-slice-handbuch-public-doc-check-gate-und-skill.md)
(Exit 2 für Lesefehler im Handbuch-Gate) ist ebenfalls unbeantwortet; eine einzige
Folge-ADR könnte beide Gates decken.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Folge-ADR wird nötig und
  ihr Entwurf ist Teil des Slice — dann trennt sich der Architect-Zug ab.
- `in-progress` → `open` (blockiert): die Architect-Antwort verlangt, die Prune-Liste oder
  das Muster mitzuändern.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
Mutationsprobe rot gesehen, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

- **Der Tabellentest läuft unter uid 0 nicht.** `chmod 000` hindert root nicht; die
  Lesefehler-Fälle werden dort übersprungen (wie im Handbuch-Tabellentest). Läuft die CI als
  root-Container, sind die Fälle dort nicht gedeckt. — **Ausgang:** (bei Closure)
- **Falsch rot durch ein nicht lesbares Verzeichnis, das früher unbemerkt blieb.** Ein
  Verzeichnis unter `sdks/`, das der Nutzer nicht lesen darf (Bau-Artefakte mit fremdem
  Eigentümer, etwa von einem Container-Lauf), färbt den Lauf künftig rot. Die Richtung ist
  die sichere; der Befund am echten Baum ist am Start zu messen. — **Ausgang:** (bei Closure)
- **Architect-Antwort bewegt den Umfang.** Verlangt die Antwort eine Folge-ADR, ist sie
  Voraussetzung und kein Teil dieses Slice (§4). — **Ausgang:** (bei Closure)

## 7. Closure-Notiz

*(Wird bei der Closure vom Planner geschrieben. Ursprung der Angaben nach
[`AGENTS.md`](../../../../AGENTS.md) §3.12: gemessen · übernommen · hergeleitet.)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag (Lerneintrag):** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den Pfaden
`tools/harness/` und `harness/sensors/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Nachbar
`intern-kennungen-in-ausgelieferten-texten` (das Gate ist dort als Träger der Bereinigung
genannt); kein weiterer Treffer gelesen.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
