# Slice nats-drittstream-example-csharp-kotlin: C#/Kotlin-Beispiel-Clients `nats-stream-client`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-nats-drittstream.

**Bezug:** [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md) (aktiv,
`Accepted`, Teilfrage 6), [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
(volle Matrix — kein Sprach-Pilot), [`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)
(Sprachwurzel-Bau-Form), [`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Startform).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/pflichtenheft.md)
(Beispiel-Client-Bestand/Namenskonvention).

**Verantwortlich:** —.

**Autor:** Planner-Rolleninhaber (Modul 8). **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten.

**Ziel:** `make example-run-csharp SURFACE=nats-stream` und
`make example-run-kotlin SURFACE=nats-stream` bauen (`make examples-csharp`/
`make examples-kotlin`) und starten real je einen Client, der mit
`CDC_NATS_STREAM_TOKEN` verbindet, `cdc.stream.>` abonniert und eine
empfangene Change lesbar ausgibt — die volle Drei-Sprachen-Matrix für den
vierten Zugriffsweg ist damit vollständig ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Go-Beispiel-Client** — eigener, unabhängiger `slice-nats-drittstream-example-go`
  (beide hängen nur an `slice-nats-drittstream-core`, `welle-nats-drittstream`
  §5).
- **Demo-Umgebung (`examples/compose.yaml`, `examples/.env`)** —
  `slice-nats-drittstream-core` trägt die additive NATS-Auth-Erweiterung
  (geteilte Infrastruktur); dieser Slice startet ausschließlich gegen eine
  bereits laufende Umgebung.
- **Server-seitige Änderungen** — vollständig `slice-nats-drittstream-core`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] `examples/csharp/nats-stream-client` und
      `examples/kotlin/nats-stream-client` existieren, verbinden mit
      `CDC_NATS_STREAM_TOKEN`, abonnieren `cdc.stream.>` und geben eine
      empfangene Change lesbar aus; `make examples-csharp`/
      `make examples-kotlin` bauen+testen beide netzlos
      ([`LH-FA-SST-008`](../../../../spec/lastenheft.md)).
- [ ] `harness/mk/examples.mk`: `examples-csharp`/`examples-kotlin` bekommen
      je einen vierten `docker build --target runtime-nats-stream`-Aufruf;
      `example-run-csharp`/`example-run-kotlin`s `$(filter …)`-Prüfungen
      tragen `nats-stream` als vierten zulässigen Wert (`ADR-0100`
      Teilfrage 6).
- [ ] Doku-Nachzug: `examples/README.md` bekommt die C#-/Kotlin-Zeilen im
      neuen Zugriffs-Abschnitt, `docs/user/benutzerhandbuch.md`s
      `**Beispiele:**`-Block wird um beide Sprachen komplettiert, die
      Aggregat-Aussage „vier Zugriffsarten × drei Sprachen" wird auf „fünf …
      fünfzehn Programme" gehoben (jetzt erst wahr).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `examples/README.md` (C#-/Kotlin-Zeilen, Aggregatzahlen),
      `docs/user/benutzerhandbuch.md`s `**Beispiele:**`-Block komplettiert
      (siehe §2 dritter Liefer-Punkt).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area?

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/csharp/nats-stream-client/*.csproj`, `Program.cs`, `NatsStreamClient.Tests/` | neu | analog zu `examples/csharp/nats-client`, Token-Verbindung statt anonym |
| `examples/csharp/Dockerfile` | update | neue `COPY`/`RUN dotnet restore/build/test/publish`-Zeilen für `nats-stream-client` in der `build`-Stufe, neue Stufe `runtime-nats-stream` (Vorbild: die bestehende `runtime-nats`-Stufe, Zeile ~93) |
| `examples/kotlin/nats-stream-client/` (Gradle-Modul) | neu | analog zu `examples/kotlin/nats-client` |
| `examples/kotlin/Dockerfile` | update | neues Gradle-Modul in der `build`-Stufe, neue Stufe `runtime-nats-stream` |
| `harness/mk/examples.mk` | update | vierter `docker build --target runtime-nats-stream`-Aufruf je Sprache in `examples-csharp`/`examples-kotlin`; `nats-stream` in beiden `example-run-*`-Filtern |
| `examples/README.md` | update | neue Tabellenzeilen in `## C#` und `## Kotlin`; Intro-Zahlen „vier × drei" → „fünf × drei" |
| `docs/user/benutzerhandbuch.md` | update | C#-/Kotlin-Zeilen im von `slice-nats-drittstream-example-go` angelegten Abschnitt, neue Versionshistorie-Zeile, Header-Aggregatzahlen falls vorhanden |

**Ansatz:** Mechanische Übertragung des bestehenden `nats-client`-Musters
(beide Sprachwurzeln) auf einen neuen Client mit anderem Subjekt-Filter
(`cdc.stream.>` statt `cdc.changes.>`), anderer Umgebungsvariable
(`CDC_NATS_STREAM_TOKEN` statt anonym) und anderem Nachrichtenschema
(10-Feld-JSON statt leerer Payload) — kein neuer Bau-Mechanismus, keine
`--build-context proto=proto`-Berührung (kein Protobuf-Bezug).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-nats-drittstream-core` liegt in
`done/`. Unabhängig von `slice-nats-drittstream-example-go` (§5 der
Welle-Datei) — kann vor, nach oder gleichzeitig mit ihm laufen, sofern kein
WIP-Limit-1-Rolleninhaber beide zugleich trägt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): unwahrscheinlich —
  strukturidentisch zu `slice-101` (`nats-client-csharp-kotlin`), demselben
  Bündelungsmuster.
- `in-progress` → `open` (blockiert — Carveout?): `slice-nats-drittstream-core`
  liegt noch nicht in `done/` — zurück nach `open/`.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

- DoD (§2) vollständig, `make gates` grün, Review-Report vorliegt.
- `make example-run-csharp SURFACE=nats-stream` und
  `make example-run-kotlin SURFACE=nats-stream` real ausgeführt (Docker
  verfügbar), je ein empfangenes Event gezeigt.
- Closure-Notiz mit Steering-Loop-Lerneintrag (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Zwei Sprachwurzeln in einem Slice** (wie `slice-101`, anders als die
  getrennten `slice-102`/`slice-103` für gRPC) — falls C# und Kotlin
  unerwartet divergieren (z. B. ein Kotlin-Gradle-Problem analog zu
  `BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`), könnte der
  Slice zu groß werden. — **Ausgang:** <bei Closure ausfüllen; bei
  Bedarf Aufspaltung in zwei Folge-Slices>.
- **Aggregat-Zahlen an mehreren Stellen** (`examples/README.md` Intro,
  Handbuch-Header, falls vorhanden) — ein Nachzug-Fund könnte eine Stelle
  übersehen (`AGENTS.md` §3.13, Suchlauf-Pflicht). — **Ausgang:** <bei
  Closure ausfüllen; Suchlauf-Ergebnis (gefunden/nicht gefunden) im §7
  dieses Slice dokumentiert>.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(bei Closure zu füllen — inklusive des Suchlauf-Ergebnisses aus §6
zweitem Risiko, `AGENTS.md` §3.13.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** einzige berührte Sub-Area ist `*`
(Default) — `examples/csharp/**`, `examples/kotlin/**` sind Beispiel-/
Doku-Fläche ohne eigene Sub-Area-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`handbuch`, `beispiel`, `nats`, `dockerignore`, `zusatzkontext`) — Treffer:
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert, hier in §2 als DoD-Punkt adressiert);
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (ohne Bezug — dieser
Slice fügt keinen neuen Bau-Kontext-Pfad hinzu, nur ein neues Modul innerhalb
bereits erlaubter Sprachwurzeln). Keine weiteren Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.

### Sub-Area: `*` (Default)

- **Modus:** GF — neue Beispiel-Verzeichnisse nach bestehendem Vorbild
  (`examples/csharp/nats-client`, `examples/kotlin/nats-client`), kein
  Bestandscode.
