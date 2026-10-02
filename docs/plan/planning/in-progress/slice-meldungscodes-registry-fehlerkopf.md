# Slice meldungscodes-registry-fehlerkopf: Code-Tabelle, Fehlerkopf `Fehlerklasse <klasse> [<code>]`, Handbuch-Katalog und Gate `meldungscodes-check` (Teil 2 von 4)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md) (Betriebsfähigkeit),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (Fehlerklassen),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Festlegungen 1–5, 8, 9 T2),
[`ADR-0023`](../../adr/0023-fehlerklassifikation.md) und
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (Klassen, unverändert),
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(Handbuch-Gate; Abgrenzung in Festlegung 7),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).
Verdikt: [`architect-verdict-meldungscodes-statt-interner-kennungen`](../../../reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md) (neuer Absatz
„Meldungscode“; der Satz „der Fehlertext beginnt mit der Klasse“ wird zum Kopf nach
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 3),
[`SPEC-029`](../../../../spec/pflichtenheft.md) (`error_message` des Backfill-Runs),
[`SPEC-019`](../../../../spec/pflichtenheft.md) (`error_message` des Antrags) — Liefer-Punkt
dieses Slice, nicht des Auftraggebers.

**Reihenfolge:** nach [T1 `meldungscodes-kennungsfreie-ausgaben`](../done/slice-meldungscodes-kennungsfreie-ausgaben.md)
(sein Gate hält die Ausgaben bei 0, bevor dieser Slice die Fehlertext-Literale anfasst); vor
[T3](../open/slice-meldungscodes-warnungen-heartbeat-diagnose.md) und
[T4](../open/slice-meldungscodes-http-grpc-fehlerkoerper.md), die an der Tabelle dieses Slice hängen.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `ba60c7bc`, 2026-10-02; Befehle im Feld §3).**
Die Quelle der Fehlertexte trägt ein einheitliches Präfix `Fehlerklasse <klasse>: …`:
**28** Literale in Produktions-Go. `classifyRunError` (`internal/bootstrap/wiring.go`)
bildet Sentinels über `errors.Is` auf die sieben Klassen ab; der Backfill-Dienst entfernt den
Kopf per Textersetzung (`internal/application/usecase/backfill/service.go`,
`strings.Replace(cause.Error(), "Fehlerklasse "+string(class)+": ", "", 1)`) — diese
Textchirurgie ersetzt der Slice durch den Code des Fehlerwerts. Das Handbuch nennt das Wort
`Fehlerklasse` in 36 Zeilen. Kein `PCF-` im Baum (0 Treffer).

**Ziel:** Es gibt eine Code-Tabelle als Quelle der Wahrheit (Paket `internal/domain`),
jeder klassifizierte Fehler trägt den Code seiner Einzelursache oder den Rückfall seiner
Klasse, die Wege dieses Slice zeigen den Kopf `Fehlerklasse <klasse> [<code>]: …`, das
Handbuch trägt den Katalog, und das Gate `meldungscodes-check` hält Quelltext, Tabelle und
Katalog gleich ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegungen 1–5).

**Ausdrücklich NICHT in diesem Slice:**

- **Warnungs-Codes (`W`), Heartbeat-Spalte `error_code`, Diagnose** — T3.
- **`code` im HTTP-Fehlerkörper, gRPC `ErrorInfo`** — T4. Die Ablehnungs-Codes `E8…`
  entstehen in der Tabelle hier (der Antrag trägt sie in `error_message`), ihre Netz-Wege in T4.
- **Änderung der Fehlerklassen-Semantik, Exit-Codes** — Klassen und Ausgang 1 bleiben
  ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2, 3).
- **Code als Metrik-Label** — ausgeschlossen (Festlegung 2).
- **SDK-Änderungen, Release, Tag** — Freigaben des Auftraggebers; der Slice ändert `sdks/`
  nicht (Server-Release-Frage: Empfehlung 3 des Verdikts, Auftraggeber-Entscheidung).
- **Lastenheft** — keine Änderung nötig; eine Anforderung wäre Auftraggeber-Entscheidung.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(A) Tabelle, Kopf und Wege.** Die Tabelle in `internal/domain` (Code, Schwere, Klasse,
      Status) enthält die sieben Rückfall-Codes `PCF-E1000`…`PCF-E7000` und die Einzelursachen
      der Sentinels (Granularität: unterschiedliche Maßnahme des Betreibers,
      [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2) sowie
      die Ablehnungs-Codes `E8…` des Antrags; die 28 `Fehlerklasse`-Literale verwenden die
      Konstanten der Tabelle (kein zweites Literal); der Kopf `Fehlerklasse <klasse> [<code>]: …`
      steht in Fehlerwert, Prozess-Ende-Zeile und `error`-Attribut; `error_message` von Run
      (`<klasse> [<code>]: …`) und Antrag (`abgelehnt [<code>]: …`); die Textchirurgie im
      Backfill-Dienst entfällt; Rollout-Fehlerzeilen als `FEHLER [<code>]: …`; Spec-Nachzug
      (`SPEC-008` Absatz, Satz zum Kopf, `SPEC-029`, `SPEC-019`). *Zu belegen durch:* Go-Test in
      `make test` (Format gegen ERE `PCF-[EWI][0-9]{4}`, erste Ziffer gegen Klasse, keine Dopplung,
      Rückfall je Klasse; ein Code mit falscher Klasse färbt rot — in der ADR nur *hergeleitet*,
      der Implementer fährt die Probe und berichtet Stelle und Farbe); Suchlauf §3;
      `make test-integration` und `make test-store` grün nach `make image`.
- [ ] **(B) Katalog und Gate.** Handbuch-Katalog (kennungsfrei: Code, Klasse bzw. Bereich,
      Bedeutung, Maßnahme; zurückgezogene Codes mit Vermerk; der Satz „Der Text ist nicht
      Vertrag, der Code ist es“; Änderungshistorie in Betreibersicht); Gate `meldungscodes-check`
      in `GATE_CHECKS` (Mengengleichheit Quelltext/Tabelle/Katalog, nur `bash`/`git`/`grep`),
      Heimat `harness/mk/doc-gate.mk`, Sensor-Vertrag `harness/sensors/meldungscodes-check.md`,
      Tabellentest `make test-meldungscodes-check` (Werkzeug ohne Gate-Aufnahme), Zeilen in
      `harness/README.md` §Sensors. *Zu belegen durch:* Tabellentest (Code im Quelltext ohne
      Tabelle → Exit 1, Tabellen-Code ohne Katalog → Exit 1, zurückgezogener Code bleibt in
      beiden Mengen, saubere Menge → Exit 0) und `make handbuch-public-doc-check` grün **mit**
      Katalog (Gate-Lauf; in der ADR nur aus den Mustern *hergeleitet*, Festlegung 7);
      `make gates` Exit 0.
- [ ] **(C) Textzusage und Bestand geprüft.** Der Implementer liest die Handbuch-Abschnitte
      Diagnose, Fehlerklassen und Exit-Codes auf eine Zusage des Fehlertexts
      ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Konsequenzen: im Handbuch
      nicht nachgemessen); findet er eine, geht der Slice mit dieser Frage zurück an den
      Auftraggeber. Er misst außerdem, ob `sdks/` Servertexte der Form `Fehlerklasse …` in Tests
      erwartet (Messung am Parent: 0 Zeilen mit `Fehlerklasse` unter `sdks/`). *Zu belegen
      durch:* Befehle und Ergebnis im Bericht (§3.12: gemessen, nicht übernommen).
- [ ] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review; der Reviewer liest den Katalog
      gegen den Code (akzeptiertes Negativ der ADR: das Gate vergleicht Mengen, nicht Sinn).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-registry-fehlerkopf.md`
      endet mit Exit 0.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten, nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/` (neue Tabelle), Go-Test dazu | neu | Quelle der Wahrheit, Registry-Test (A). |
| 28 `Fehlerklasse`-Literale (u. a. `mapper.go`, `receive.go`, `decode.go`, `tableactivation.go`, `internal/application/port/outbound/*.go`, `internal/bootstrap/wiring.go`) | update | Kopf mit Code aus der Tabelle. |
| `internal/bootstrap/wiring.go` (`classifyRunError`), `internal/application/usecase/backfill/service.go` | update | Code des Fehlerwerts statt Textersetzung; `error_message` des Runs. |
| Antragsverarbeitung (`error_message` einer abgelehnten Antrags-Zeile) | update | `abgelehnt [<code>]: …`, Codes aus dem Bereich `E8`. |
| `tools/schema/rollout.sh`, `tools/schema/rolloutguard/` | update | `FEHLER [<code>]: …`. |
| `spec/pflichtenheft.md` (`SPEC-008`, `SPEC-029`, `SPEC-019`) | update | Spec-Nachzug (Liefer-Punkt dieses Slice). |
| `docs/user/benutzerhandbuch.md` | update | Katalog, Satz zur Stabilität, Änderungshistorie. Rebase auf `main`, nur eigene Abschnitte committen. |
| Gate-Skript unter `tools/harness/`, Tabellentest-Läufer, `harness/sensors/meldungscodes-check.md`, `harness/README.md` §Sensors, `harness/mk/doc-gate.mk`, `GATE_CHECKS` | neu / update | (B). |
| `*_test.go` mit Erwartung `Fehlerklasse <klasse>:` (u. a. `internal/application/usecase/backfill/service_test.go`), Läufer unter `tools/harness/` | update | Erwartungen an den Kopf. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Fehlertext beginnt mit
`Fehlerklasse <klasse>: `“; Parent ist `ba60c7bc`; die `diff`-Zeilen und die Befunde trägt der
Implementer ein):**

```suchlauf
ba60c7bc 28 -n -E '"Fehlerklasse [a-z]+: ' -- '*.go' ':!*_test.go'
ba60c7bc 1 -n -F 'strings.Replace(cause.Error(), "Fehlerklasse "' -- internal/application/usecase/backfill/service.go
ba60c7bc 1 -n -F 'func classifyRunError' -- internal/bootstrap/wiring.go
ba60c7bc 1 -n -F 'der Fehlertext beginnt mit der Klasse' -- spec/pflichtenheft.md
ba60c7bc 36 -n -F 'Fehlerklasse' -- docs/user/benutzerhandbuch.md
ba60c7bc 2 -n -E 'Fehlerklasse [a-z]+:' -- internal/application/usecase/backfill/service_test.go
ba60c7bc 0 -n -F 'Fehlerklasse' -- sdks
ba60c7bc 0 -n -F 'PCF-' -- internal tools sdks docs/user spec cmd examples
```

| Träger | Messung am Parent (`ba60c7bc`, 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| `Fehlerklasse`-Literale in Produktions-Go | Zeile 1: 28 (die Tabellenzeile des Vorläufer-Plans gehört jetzt hierher, nicht zu T1) | Diff: Kopf mit Code; Soll je Literal die Konstante der Tabelle |
| Textchirurgie, Klassifikation | Zeilen 2, 3: je 1 | Zeile 2 entfällt (Soll 0) |
| Spec-Satz zum Kopf | Zeile 4: 1 | auf den Kopf nach Festlegung 3 umgeschrieben |
| Handbuch | Zeile 5: 36 | Implementer sichtet jede Zeile auf Textzusage und Katalog-Bezug |
| Tests | Zeile 6: 2 | Erwartungen mitgezogen |
| SDKs | Zeile 7: 0 | Nichtgefunden belegt; ändert sich in diesem Slice nicht |
| Code-Präfix im Baum | Zeile 8: 0 | Diff: nur Tabelle, Katalog, Quelltext-Konstanten |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt, [T1](../done/slice-meldungscodes-kennungsfreie-ausgaben.md) liegt in
`done/`, **und das Code-Präfix `PCF-` ist vom Auftraggeber bestätigt oder geändert**
([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Empfehlung 1 des
Verdikts: bisher nicht bestätigt). Wird es geändert, entsteht vor dem Start eine Folge-ADR
mit `Supersedes ADR-0144`, die Form und Muster nachzieht (Festlegung 1, 5, 7); nach diesem
Slice ist das Präfix nur noch per Folge-ADR änderbar. Der Hauptlauf klärt die Frage vor dem Start.

**Rückführung:** `in-progress` → `open` (blockiert), wenn der Implementer eine Zusage des
Fehlertexts im Handbuch findet (Liefer-Punkt C) oder ein SDK-Test den Servertext erwartet
(Auftraggeber-Entscheidung; die Änderung geht nie als Anpassung einer Erwartung durch).
`in-progress` → `next` (zu groß), wenn Tabelle und Wege zusammen nicht in einen Diff passen:
Schnitt entlang `application`/`adapters` (Einzelursachen-Codes) gegen Katalog/Gate, die ADR
bleibt unverändert.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration`, `make test-store` grün nach `make image`, Registry-
und Gate-Probe in den Bericht, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag.
Der sichtbar geänderte Wortlaut des Fehlerkopfs ist eine Betreiber-Änderung: ob sie ein
Server-Release mit eigener Freigabe braucht, entscheidet der Auftraggeber (Verdikt, Empfehlung 3).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **Präfix unbestätigt.** `PCF-` ist Empfehlung der ADR, nicht Auftraggeber-Entscheidung;
  nach diesem Slice ist es nur per Folge-ADR änderbar. Gegenmittel: §4-Startbedingung. —
  **Ausgang:** (bei Closure)
- **Ausgaben-Stabilität für Betreiber.** Der Kopf ändert sich zu
  `Fehlerklasse <klasse> [<code>]: …`; ein `grep` auf `Fehlerklasse schema:` bricht. Ausgang der
  ADR: keine Übergangsregel, Text nicht Vertrag; das Handbuch wird hier gemessen (Liefer-Punkt C). —
  **Ausgang:** (bei Closure)
- **SDKs und Clients reichen Servertexte durch.** Ein SDK zeigt künftig den Code im Text mit;
  SDK-Tests, die den Wortlaut prüfen, wären betroffen. Gemessen am Parent: 0 Zeilen
  `Fehlerklasse` unter `sdks/`; Test-Erwartungen anderer Fehlertexte misst der Implementer. —
  **Ausgang:** (bei Closure)
- **Granularität der Einzelursachen.** Zu feine Codes erzeugen Pflege ohne Nutzen, zu grobe
  verwischen die Maßnahme; Regel der ADR: Maßnahme des Betreibers. Gegenmittel: Rückfall
  `…000` garantiert, dass kein klassifizierter Fehler ohne Code bleibt; der Reviewer liest die
  Zuordnung. — **Ausgang:** (bei Closure)
- **Handbuch-Gate mit Katalog.** Dass `make handbuch-public-doc-check` mit Codes grün bleibt,
  ist aus den Mustern hergeleitet, nicht gelaufen. — **Ausgang:** (bei Closure)
- **Kollision mit parallelen Arbeiten am Handbuch.** Rebase auf `main`, nur eigene Abschnitte
  committen. — **Ausgang:** (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** T3 `meldungscodes-warnungen-heartbeat-diagnose`, T4 `meldungscodes-http-grpc-fehlerkoerper`.
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird die Default-Sub-Area `*` (Kürzel
`PGC`, Modus Greenfield, [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); eine feinere Aufteilung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register `../observations/BEO-PGC/` ist
nicht Inhalt dieses Plans; der Implementer sichtet es vor dem Start auf Treffer zu
„Fehlertext“ und „Kennung in Ausgabe“ (nicht gemessen, kein Treffer behauptet).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
