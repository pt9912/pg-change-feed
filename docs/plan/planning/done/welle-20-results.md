# Welle 20 — Coverage 80 % über der netzlos prüfbaren Fläche — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-NNN` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-20
**Abschluss:** 2026-09-17
**Verantwortlich:** pt9912

## Was wurde geliefert?

<!-- BEDIENHINWEIS: Ergebnis, nicht Taetigkeit. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Das Welle-Ziel ist erreicht: 80 % über der netzlos prüfbaren Fläche.** Die
  Gate-getragene Zahl steht bei **83,10 %** (gedruckte Zeile `make gates`, Lauf
  dieser Closure) — Einstieg war 70 %. Der **Nenner** ist dabei **unverändert**
  1903 geblieben: die Welle hat ihr Ziel durch **Test-Arbeit** erreicht, nicht
  durch einen Schnitt am Gegenstand.
- **Vier Cluster, sechs Slices** — `slice-079` (Scope-Schnitt), `slice-088` (B),
  `slice-091` (C), `slice-092` (D1), `slice-093` (D2), `slice-094` (A). Alle in
  `done/`. Zusammen **~161 Statements** neu gedeckt.
- **Die Rampe steht auf ihrer Endstufe**: `THRESHOLD ?= 80` in `make coverage-gate`
  (Hochschalt-Trigger ausgeschöpft); zwei Belege, je eigener Lauf: grün bei 80
  (`83.10%`, Exit 0) und **rot** bei 85 (Exit ≠ 0).
- **§3.13 neu** (`AGENTS.md`) und ein neuer HIGH-Punkt in `.harness/skills/reviewer.md`
  („Beleg trägt seinen Satz nicht") — beide aus dem Lese-Schritt dieser Closure.
- **Fünf ADRs** aus dieser Welle: `ADR-0082` (Schnittmaß), `-0083` (Herkunft von
  Aussagen), `-0084` (Sync-Gate), `-0085` (Build-Kontext-Ausnahme), `-0086`
  (Schwere folgt der Konsequenz).

## Was hat funktioniert?

<!-- BEDIENHINWEIS: was du im naechsten Zyklus bewusst wieder so machen wuerdest. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Messen, bevor geschnitten wird.** Die vier Cluster wurden nach einer
  **eigenen Messung** geschnitten, nicht nach einem Plan; das hat den ersten
  Vorschlag der Welle („`internal/bootstrap` ist der Hebel") **widerlegt** —
  `Run` ist netzlos gar nicht prüfbar.
- **Die Mutation als Rückgrat jedes Test-Slice.** Über die Welle: **~120**
  Mutationsproben über alle Rollen; sie haben zwei Testkommentare, eine
  Assertion, eine Adresse und drei Sätze als **unwahr** entlarvt — und in
  `slice-094` eine erfundene Grenze.
- **Der Register-Auftrag.** Aus `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
  entstand in jedem folgenden Slice ein `grep`-Auftrag — er hat **vier** Stellen
  gefunden, die niemand sonst gesucht hätte.
- **Das Delta-Review als Pflicht.** Seit `slice-092` steht sie in der DoD; sie
  hat in **jedem** Slice dieser Welle mindestens einen Satz gefunden, den die
  erste Runde nicht sehen konnte.

## Was ging anders als geplant?

<!-- BEDIENHINWEIS: Beobachtungen, keine Schuldzuweisung. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Aus vier Clustern wurden sechs Slices**, und aus vier Test-Slices wurden
  **vierzehn** Korrekturrunden. Der Grund liegt nicht im Schnitt, sondern im
  Gegenstand: **jeder** Fund lag in einem **Satz über** den Mechanismus, keiner
  im Mechanismus. Die vier Slices lieferten korrekte Tests und beschrieben sie
  falsch.
- **Die Schätzung „≈52 netzlos erreichbar" für Cluster C war um fünf zu hoch**
  (real 47); die für Cluster A um vier zu niedrig (real 58). Beide Zahlen trugen
  die Welle weiter, aber die **Differenz** war als Grenze formuliert — siehe
  `BEO-PGC/geschaetzter-wert-als-grenze`.
- **Die Coverage-Zahl der Welle ist selbst nicht stabil**: ein Band von 1–2
  Statements über denselben Baum, dessen Quelle erst spät lokalisiert war
  (`runWALRetentionCheck`, dann `runAdministration`).
- **Die Welle hat ihre eigene Slice-Liste nicht nachgezogen** — sie stand seit
  Cluster C bei zwei Zeilen, während vier weitere entstanden. Gefunden hat das
  der **Nutzer**, nicht eine der vier Rollen.

## Steering-Loop-Einträge

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

**Acht Einträge** haben in dieser Welle 3× erreicht (§Lese-Schritt). Ihre
Ausgänge, alle **verkörpert**:

- **Reviewer-Skill** ergänzt: neuer HIGH-Punkt „**Beleg trägt seinen Satz nicht**"
  (wer einen Beleg nennt, **fährt** ihn) — liegt in `.harness/skills/reviewer.md`.
  Auslöser: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (`slice-084`,
  `slice-085`, `slice-091`, `slice-093`, `slice-094` — 5×).
- **Hard Rule** ergänzt: **§3.13** „Eine Arbeit, die eine beschriebene Eigenschaft
  bewegt, zieht ihre Träger nach" — liegt in `AGENTS.md §3.13`.
  Auslöser: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (`slice-091`,
  `slice-093`, `slice-094` — 3×).
- **Regel bestätigt** (kein neuer Träger): `AGENTS.md §3.12` Instanz A und der
  Reviewer-HIGH-Punkt „Zahl im Träger ohne Ursprung" — liegt in `AGENTS.md §3.12`.
  Auslöser: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (`slice-081`,
  `slice-084`, `slice-085`, `slice-088`, `slice-089`, `slice-090`, `slice-092` — 7×).
- **Regel bestätigt**: der Reviewer-HIGH-Punkt „Zusage ohne Bindung an ihre
  Eingabeseite" („mutiere den Eingabewert, nicht nur die Ausgabeseite") — liegt in
  `.harness/skills/reviewer.md`. Auslöser:
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (`slice-083`, `slice-086`,
  `slice-087`, `slice-088`, `slice-091`, `slice-092` — 6×).
- **Regel bestätigt**: `AGENTS.md §3.12` Instanz B — liegt in `AGENTS.md §3.12`.
  Auslöser: `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`
  (`slice-036`, `slice-081`, `slice-082`, `slice-083` — 4×).
- **Gate bestätigt**: `make generated-sync` — liegt in `Makefile:generated-sync`.
  Auslöser: `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (`slice-069`,
  `slice-074`, `slice-082`, `slice-086` — 4×).
- **Regel bestätigt**: `AGENTS.md §3.12` Instanz B — liegt in `AGENTS.md §3.12`.
  Auslöser: `BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg` (`slice-079`,
  `slice-080`, `slice-089` — 3×).
- **Regel bestätigt**: `AGENTS.md §3.12` Instanz A und `harness/sensors/coverage-gate.md`
  §Zählbasis — liegt in `AGENTS.md §3.12`. Auslöser:
  `BEO-PGC/test-integration-retention-timing-flake` (`slice-057`, `slice-090`,
  `slice-091` — 3×). **Nicht gestrichen** — die Beobachtung kann noch auftreten;
  ein benannter Träger ist weggefallen (`slice-093` fährt ihn deterministisch).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in `observations/BEO-PGC/` (je Eintrag ein Verzeichnis; die
Zahl seiner `evidence/`-Dateien **ist** der Zähler).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*.

## Folge-Slices

<!--
DERIVATIV: der Folge-Slice selbst ist eine Datei in `open/`; diese Liste
zeigt nur darauf. Deshalb braucht sie keinen eigenen Konsumenten — wohl
aber eine Deckung: jeder genannte Folge-Slice MUSS als Datei im
Planning-Lifecycle existieren (`open/`, `next/`, `in-progress/`, `done/` —
nicht nur `open/`, er kann bis zur Prüfung weitergewandert sein).
Folge-Slice-Paarung, geprüft am Ende von Schritt 3 der Closure-Prozedur.
Genannt ohne angelegt ist dieselbe Klasse wie ein halluziniertes Gate.
-->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

- `slice-095` (Beispiel-Clients unter `examples/` — die drei fehlenden) —
  liegt in `next/`; die **Rückführung** `in-progress → next` ist vollzogen, weil
  einer der drei (der gRPC-Client) ohne den fehlenden Umzugs-Slice aus
  [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §3 nicht lieferbar ist. Zwei der drei sind geliefert
  und committet.
- **Kein weiterer Folge-Slice ist benannt.** Der **Umzugs-Slice**
  ([`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §3, Empfehlung 1) ist eine **offene Adresse ohne
  Kennung** — er wird als eigener Schnitt angelegt, und *benannt ohne angelegt*
  wäre dasselbe wie ein halluziniertes Gate.

## Verifikation

<!--
Die Belege aus Schritt 1 der Closure-Prozedur. Keine Behauptung ohne
nachprüfbaren Anker (Hash, Lauf, Zahl).
-->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- **`make gates` grün, Exit 0** (Lauf dieser Closure): `baseline-verify` 54
  Dateien · `docs-check` 784 Dateien, 0 Befunde · `a-check` 0 Befunde ·
  `commit-traceability` OK · `generated-sync` byte-gleich · `coverage-gate` OK.
- **Coverage: 83,10 %** über der netzlos prüfbaren Fläche, Nenner **1903**
  (gedruckte Zeile des Gate-Laufs dieser Closure).
- **Grün-Beleg bei der Endstufe:** `make coverage-gate THRESHOLD=80` → **Exit 0**
  (`83.10% erfüllt Schwelle 80%`).
- **Rot-Beleg unmittelbar darüber:** `make coverage-gate THRESHOLD=85` → **Exit 2**
  (`83.10% unter Schwelle 85%`) — die Stufe prüft real.
- **Offene Carveouts: keine.**
- **Archivierung (Schritt 4): nicht eingetreten.** Dieses Repo führt das
  Werkzeug nicht, das die Zeitdokumente der Welle einsammelt (`tools/harness/`
  trägt kein Archiv-Skript, `Makefile`/`harness/mk/` kein Ziel). Die Bedingung
  ist damit nicht eingetreten — die Feststellung steht hier, statt einen
  Handlauf zu erfinden (Baseline-Regelwerk `modul-06-roadmap.md`
  §Wellen-Closure-Prozedur, Schritt 4). Der Geltungsbereich der Sensoren ist
  deshalb **nicht** zu prüfen: es wurde nichts bewegt.
