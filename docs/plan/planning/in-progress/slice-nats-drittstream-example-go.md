# Slice nats-drittstream-example-go: Go-Beispiel-Client `nats-stream-client`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-nats-drittstream.

**Bezug:** [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md) (aktiv,
`Accepted`, Teilfrage 6), [`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Startform-Präzedenzfall, `SURFACE=`-Muster), [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(Beispiel-Client-Vertrag).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/pflichtenheft.md)
(Beispiel-Client-Bestand/Namenskonvention — neuer vierter Zugriffsweg).

**Verantwortlich:** —.

**Autor:** Planner-Rolleninhaber (Modul 8). **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar.

**Ziel:** `make example-run-go SURFACE=nats-stream` baut und startet real
einen Go-Client, der mit `CDC_NATS_STREAM_TOKEN` verbindet, `cdc.stream.>`
abonniert und eine empfangene Change lesbar ausgibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **C#/Kotlin-Beispiel-Clients** — eigener, unabhängiger Folge-Slice
  `slice-nats-drittstream-example-csharp-kotlin` (beide hängen nur an
  `slice-nats-drittstream-core`, nicht aneinander, `welle-nats-drittstream`
  §5).
- **Demo-Umgebung (`examples/compose.yaml`, `examples/.env`)** —
  `slice-nats-drittstream-core` trägt die additive NATS-Auth-Erweiterung der
  Demo-Umgebung (geteilte Infrastruktur, nicht sprachspezifisch — dieselbe
  Server-Konfiguration, die auch die Root-`compose.yaml`-Testumgebung
  bekommt); dieser Slice startet ausschließlich gegen eine bereits laufende
  Umgebung (`make example-demo-up`), legt keine eigene an.
- **Server-seitige Änderungen** (`natsstream.Publisher`, Bootstrap-
  Verdrahtung) — vollständig `slice-nats-drittstream-core`, dieser Slice
  konsumiert nur den bereits fertigen dritten Zustellweg.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] `examples/nats-stream-client` (Go) existiert, verbindet mit
      `CDC_NATS_STREAM_TOKEN`, abonniert `cdc.stream.>` und gibt eine
      empfangene Change lesbar aus; `make example-run-go SURFACE=nats-stream`
      baut und startet ihn real gegen die Demo-Umgebung
      ([`LH-FA-SST-008`](../../../../spec/lastenheft.md)) — real ausgeführt
      gegen eine frisch hochgefahrene Demo-Umgebung (`make example-demo-down`
      dann `-up`, weil der zuvor laufende Feed-Container älter war als die
      `CDC_NATS_STREAM_TOKEN`-Verdrahtung von `slice-nats-drittstream-core`):
      `nats-stream-client: change_id=804-1 table=public.orders
      operation=INSERT new_image={"id":"2","customer":"nats-stream-smoke-test","amount":"42.50"}`.
- [x] `harness/mk/examples.mk`s drei `$(filter …)`-Prüfungen für
      `example-run-go` (und die zugehörige `$(error …)`-Meldung) tragen
      `nats-stream` als vierten zulässigen Wert (`ADR-0100` Teilfrage 6) —
      **nur** die `example-run-go`-Zeile in diesem Slice, die beiden anderen
      (`example-run-csharp`/`-kotlin`) trägt der Folge-Slice.
- [x] Doku-Nachzug: `examples/README.md` bekommt den vierten Zugriffs-
      Abschnitt, `docs/user/benutzerhandbuch.md`s `**Beispiele:**`-Block für
      den neuen Weg trägt die Go-Zeile.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `examples/README.md`, neuer Handbuch-Abschnitt „Zugriff
      über den NATS-Vollinhalts-Stream" (siehe §2 dritter Liefer-Punkt).
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
| `examples/nats-stream-client/main.go`, `format.go`, `format_test.go` | neu | verbindet mit `CDC_NATS_STREAM_TOKEN`, abonniert `cdc.stream.>` (Wurzel-Wildcard, keine Quelle/Schema/Tabelle-Filterung nötig — anders als `examples/nats-client`, weil der Vollinhalts-Stream keinen zweiten HTTP-Holschritt braucht), gibt jede empfangene Change (10-Feld-JSON, `SPEC-024`-Schema) in einer Endlosschleife lesbar aus — Vorbild `examples/grpc-client`s `main.go`/`format.go`-Aufteilung, Nachrichtenschema/-dekodierung wie `internal/adapters/driven/natsstream` |
| `examples/Dockerfile` | update | **Plan-Nachzug** (im ursprünglichen Plan übersehen, gegen den Ist-Zustand nachgezogen vor dem Sensor-Lauf): vierter Build-Aufruf (`-o /out/nats-stream-client`) und neue Stufe `runtime-nats-stream` — ohne sie bräche `make example-run-go SURFACE=nats-stream` am `docker build --target runtime-nats-stream` ab |
| `harness/mk/examples.mk` | update | die `example-run-go`-Ziel-Zeile (`$(filter $(SURFACE),http sse grpc nats)` → `… nats nats-stream`, `$(error …)`-Text) trägt `nats-stream` |
| `examples/README.md` | update | neue Tabellenzeile in `## Go` |
| `docs/user/benutzerhandbuch.md` | update | **Plan-Korrektur**: der §4-Abschnitt „Zugriff über den NATS-Vollinhalts-Stream" existiert bereits (`slice-nats-drittstream-core`, mit einem Platzhalter-Absatz „Beispiele: folgen…"); dieser Slice ersetzt den Platzhalter durch den `**Beispiele:**`-Block (zunächst nur Go-Zeile), keine neue Sektion. Neue Versionshistorie-Zeile |

**Ansatz:** Der neue Handbuch-Abschnitt entsteht als eigener Abschnitt nach
dem Vorbild „Zugriff über das NATS-Wecksignal" — Betreiber-Doku zu Aktivierung
(`CDC_NATS_STREAM_TOKEN`), Subjekt-Form und Auth-Verhalten, **bevor** der
`**Beispiele:**`-Block folgt (derselbe Aufbau wie die drei bestehenden
Streaming-Abschnitte). Die Aggregat-Aussage „vier Zugriffsarten × drei
Sprachen" in `examples/README.md`s Intro und im Handbuch-Header wird **nicht**
in diesem Slice auf „fünf" gehoben — erst wahr, sobald alle drei Sprachen
existieren (Folge-Slice trägt das, §7-Hinweis).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-nats-drittstream-core` liegt in
`done/` — ohne einen realen `natsstream.Publisher` gäbe es nichts, gegen das
dieser Client laufen könnte. WIP-Limit 1 gegenüber
`slice-nats-drittstream-example-csharp-kotlin` gilt nicht zwingend (beide
sind voneinander unabhängig, `welle-nats-drittstream` §5) — wird aber
empfohlen, wenn ein einzelner Implementer-Rolleninhaber beide Slices trägt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): unwahrscheinlich —
  der Client ist strukturell identisch zu `examples/nats-client`, nur mit
  Token und anderem Subjekt; sollte der neue Handbuch-Abschnitt einen
  eigenständigen, größeren Konzept-Absatz brauchen, wird nur die Doku-DoD
  in einen Folge-Slice verschoben.
- `in-progress` → `open` (blockiert — Carveout?): `slice-nats-drittstream-core`
  liegt noch nicht in `done/`, obwohl dieser Slice bereits gestartet wurde
  (Prozessfehler, nicht Carveout-würdig) — zurück nach `open/`, bis der
  Vorgänger schließt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

- DoD (§2) vollständig, `make gates` grün, Review-Report vorliegt.
- `make example-run-go SURFACE=nats-stream` real ausgeführt (Docker
  verfügbar) und ein empfangenes Event gezeigt.
- Closure-Notiz mit Steering-Loop-Lerneintrag (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **`slice-nats-drittstream-core` liefert einen anderen Aktivierungspfad als
  hier angenommen** (z. B. der Config-Loader-Umbau verschiebt die
  `SPEC-016`-Feldliste in eine andere Datei) — dieser Slice läse dann gegen
  eine veraltete Annahme. — **Ausgang:** <bei Closure ausfüllen>.
- **Der neue Handbuch-Abschnitt kollidiert mit der bestehenden
  Abschnitts-Reihenfolge** (vier Streaming-/Zugriffs-Abschnitte, neuer
  fünfter) — Ankertext-Kollision mit `docs-check`s `ids`/`anchors`-Modulen
  ist unwahrscheinlich (neuer, eindeutiger Anker), aber ungeprüft. —
  **Ausgang:** <bei Closure ausfüllen>.

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

*(bei Closure zu füllen — Hinweis für den Implementer: die Aggregat-Aussage
„vier Zugriffsarten × drei Sprachen" bleibt bewusst unverändert, bis
`slice-nats-drittstream-example-csharp-kotlin` schließt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** einzige berührte Sub-Area ist `*`
(Default) — `examples/**` ist Beispiel-/Doku-Fläche ohne eigene
Sub-Area-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`handbuch`, `beispiel`, `nats`, `start`) — Treffer:
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert, `seit slice-077`, hier in §2 als DoD-Punkt
adressiert). Keine weiteren Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.

### Sub-Area: `*` (Default)

- **Modus:** GF — neues Beispiel-Verzeichnis nach bestehendem Vorbild
  (`examples/nats-client`), kein Bestandscode.
