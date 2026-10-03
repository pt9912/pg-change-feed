# Slice meldungscodes-http-grpc-fehlerkoerper: Feld `code` im HTTP-Fehlerkörper und `ErrorInfo` bei gRPC (Teil 4 von 4)

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
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Festlegung 3 HTTP/gRPC,
9 T4), [`ADR-0057`](../../adr/0057-http-grpc-api.md) (HTTP/gRPC-API),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).
Verdikt: [`architect-verdict-meldungscodes-statt-interner-kennungen`](../../../reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md) (HTTP-Fehler-Antwortform,
heute `{"error": "<Klartext>"}`; additives Feld `code`) und
[`SPEC-031`](../../../../spec/pflichtenheft.md) (gRPC-Fehler: `ErrorInfo`) — Liefer-Punkte
dieses Slice.

**Reihenfolge:** nach [T2 `meldungscodes-registry-fehlerkopf`](../done/slice-meldungscodes-registry-fehlerkopf.md)
(Tabelle, besonders die `E8…`-Ablehnungs-Codes); unabhängig von
[T3](slice-meldungscodes-warnungen-heartbeat-diagnose.md).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `ba60c7bc`, 2026-10-02; Befehle im Feld §3).**
Der HTTP-Fehlerkörper hat genau ein Feld: `Error string \`json:"error"\`` in
`internal/adapters/driving/http/middleware.go` (1 Treffer). Die gRPC-Fehler entstehen an
7 Stellen (`administration.go` 3, `interceptor.go` 3, `server.go` 1) über `status.Error`/
`status.New`. Die SDKs (C#, Kotlin; Python mitgemessen vom Implementer) modellieren den
Fehlerkörper als `ErrorResponse`.

**Ziel:** Der HTTP-Fehlerkörper trägt additiv `{"error": "…", "code": "<code>"}`; gRPC-Fehler
tragen ein Status-Detail `google.rpc.ErrorInfo` mit `reason = <code>` und
`domain = pg-change-feed` ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 3). Der Code kommt aus der Tabelle von T2: Ablehnungen von Aufrufer-Eingaben
(HTTP `400`/`404`, gRPC `InvalidArgument`/`NotFound`) tragen Codes des Bereichs `E8`, klassifizierte
Fehler den Code ihrer Einzelursache oder den Rückfall der Klasse.

**Ausdrücklich NICHT in diesem Slice:**

- **Tabelle, Fehlerkopf im Text, Katalog-Gerüst** — T2; der Katalog wird hier nur um neu
  vergebene Codes ergänzt, das Gate `meldungscodes-check` bleibt grün.
- **Warn-Codes, Heartbeat-Spalte, Diagnose** — T3.
- **SDK-Änderungen, SDK-Release, Server-Release** — Freigaben des Auftraggebers
  ([`AGENTS.md`](../../../../AGENTS.md) und Projektabsprache: jedes Release braucht neue
  Freigabe); ob ein SDK `code` modelliert, ist Folgearbeit. Das additive Feld ist für tolerante
  Leser unschädlich — **zu messen** (Liefer-Punkt 0).
- **Numerischer Ausgang, HTTP-Statuscodes, gRPC-Status-Codes** — unverändert.
- **SSE-Stream-Ereignisse** — kein Fehlerkörper; nur wenn der Implementer einen Fehlerpfad
  mit dem Körper `{"error":…}` findet, meldet er ihn (kein stilles Mitändern).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(0) Vorab-Messung (Startbedingung, vor Code).** Tolerieren die SDKs ein zusätzliches
      JSON-Feld `code` im Fehlerkörper? Die ADR nennt die Toleranz *hergeleitet*, nicht
      nachgemessen. Gemessen wird an den Modellen und Tests der SDKs (`ErrorResponse` in C#
      und Kotlin, Python-HTTP-Client): Deserialisierung mit unbekanntem Feld (Strict-Modus?),
      Tests, die den Körper byte-gleich vergleichen. *Zu belegen durch:* Befehl und Ergebnis
      im Bericht; zeigt die Messung einen brechenden Leser, geht der Slice mit der Frage an
      den Auftraggeber zurück (Auftraggeber-Entscheidung, kein stilles Anpassen).
- [ ] **(A) HTTP.** `code` im Fehlerkörper aller Fehlerpfade der HTTP-API (Tabellen-Konstanten,
      kein Literal); Spec-Nachzug `SPEC-018`; Handbuch (Fehlerantwort-Beispiel = echte Ausgabe,
      Katalog um neu vergebene Codes). *Zu belegen durch:* Unit-Tests je Fehlerpfad (Happy:
      Code im Körper; Boundary: Ablehnung `E8…` bei `400`/`404`; Negative: Auth-Fehler
      `401`/`403` — Code oder bewusst keiner, im Bericht begründet), `make test`,
      `make test-integration` (HTTP-API-Rundlauf: ein abgelehnter Aufruf gegen den laufenden
      Feed-Container zeigt den Code, Körper im Bericht **gemessen**).
- [ ] **(B) gRPC.** `ErrorInfo` mit `reason`/`domain` an allen Fehlerstellen von
      `administration.go`, `interceptor.go`, `server.go`; Spec-Nachzug `SPEC-031`; Unit-Tests
      (Details auslesen, `reason` gleich Tabellen-Code); Beispiel-/Wegwerf-Clients, die Status
      lesen, bleiben lauffähig. *Zu belegen durch:* `make test`, gRPC-Rundlauf in
      `make test-integration` (Stream-Öffnung ohne Token → `Unauthenticated` mit `ErrorInfo`),
      `make examples-csharp`/`make examples-kotlin`/`make example-run-go`-Bau soweit sie
      gRPC-Fehler auswerten.
- [ ] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-http-grpc-fehlerkoerper.md`
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
| `internal/adapters/driving/http/middleware.go` (Fehlerkörper-Typ), `errors.go`, `registerconsumer.go` und weitere Fehlerpfade | update | Feld `code`, Tabellen-Konstanten. |
| `internal/adapters/driving/grpc/administration.go`, `interceptor.go`, `server.go` | update | `ErrorInfo` an den 7 Fehlerstellen. |
| `internal/domain/` (Tabelle), Handbuch-Katalog | update | neu vergebene `E8…`-Codes der API-Ablehnungen (Tabelle bleibt Quelle). |
| `spec/pflichtenheft.md` (`SPEC-018`, `SPEC-031`) | update | Spec-Nachzug. |
| `docs/user/benutzerhandbuch.md`, `docs/user/` (Fehlerantwort-Beispiele) | update | Beispiele = echte Ausgabe, Änderungshistorie in Betreibersicht ohne Kennung. Rebase auf `main`, nur eigene Abschnitte committen. |
| `*_test.go` der Adapter, `tools/harness/httpclient`, `tools/harness/grpcclient`, Läufer unter `tools/harness/` | update | Tests und Erwartungen. |
| `sdks/` | nur messen | Liefer-Punkt 0; keine Änderung ohne Freigabe. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der HTTP-Fehlerkörper hat nur
`error`, gRPC trägt keine Details“; Parent ist `ba60c7bc`; die `diff`-Zeilen und die Befunde
trägt der Implementer ein):**

```suchlauf
ba60c7bc 1 -n -F 'json:"error"' -- internal/adapters/driving/http
ba60c7bc 2 -n -F '{"error": "<Klartext>"}' -- spec/pflichtenheft.md
ba60c7bc 7 -n -E 'status\.(Error|Errorf|New)\(' -- internal/adapters/driving/grpc ':!*_test.go'
ba60c7bc 0 -n -F 'ErrorInfo' -- internal proto spec docs/user
```

| Träger | Messung am Parent (`ba60c7bc`, 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| HTTP-Fehlerkörper-Typ | Zeile 1: 1 | Feld `code` ergänzt |
| Spec-Beschreibung HTTP | Zeile 2: 2 (`spec/pflichtenheft.md` Zeile 699 bei `SPEC-018` und Zeile 1048, eine Wiederholung der Form, die auf `SPEC-018` verweist) | `SPEC-018` um `code` erweitert, die Wiederholung an Zeile 1048 mitgezogen oder auf `SPEC-018` verwiesen |
| gRPC-Fehlerstellen | Zeile 3: 7 | jede Stelle trägt `ErrorInfo` |
| `ErrorInfo` heute | Zeile 4: 0 | Nichtgefunden; Diff: nur die genannten Träger |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt, [T2](../done/slice-meldungscodes-registry-fehlerkopf.md) liegt in `done/`,
und die **Vorab-Messung (Liefer-Punkt 0) liegt vor** und zeigt keinen brechenden SDK-Leser
(sonst Auftraggeber-Frage).

**Rückführung:** `in-progress` → `open` (blockiert), wenn Liefer-Punkt 0 einen brechenden
Leser zeigt; `in-progress` → `next` (zu groß): Schnitt HTTP gegen gRPC, die ADR bleibt
unverändert.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration` grün nach `make image`, Beispiel-Bau soweit betroffen,
Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **SDKs lesen den Fehlerkörper.** `ErrorResponse` in den SDKs deserialisiert den Körper; ein
  unbekanntes Feld ist für tolerante Leser unschädlich (*hergeleitet*, nicht nachgemessen).
  Gegenmittel: Liefer-Punkt 0 vor Code. — **Ausgang:** (bei Closure)
- **Code an der falschen Stelle.** Zwei Fehlerpfade könnten denselben Fehler mit verschiedenem
  Code melden. Gegenmittel: nur Tabellen-Konstanten, Test je Pfad, Reviewer liest die Zuordnung. —
  **Ausgang:** (bei Closure)
- **Auth-Fehler ohne Klasse.** `401`/`403`/`Unauthenticated` sind weder Ablehnung einer
  Eingabe noch klassifizierter Fehler; ob sie einen Code tragen, ist im Bericht zu begründen
  (Vorschlag an den Architect, falls die Tabelle keine Stelle hat — keine stille Erweiterung
  der ADR). — **Ausgang:** (bei Closure)
- **Läufer-Erwartungen** an Fehlerkörper werden erst im `make test-integration`-Lauf rot, der
  nicht in `make gates` liegt. Gegenmittel: Suchlauf, vollständiger Lauf vor Closure. —
  **Ausgang:** (bei Closure)
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
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird die Default-Sub-Area `*` (Kürzel
`PGC`, Modus Greenfield, [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); eine feinere Aufteilung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register `../observations/BEO-PGC/` ist
nicht Inhalt dieses Plans; der Implementer sichtet es vor dem Start (nicht gemessen, kein
Treffer behauptet).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
