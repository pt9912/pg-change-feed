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

**Verantwortlich:** — *(bis zur Priorisierung `open` → `next`)*.

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

- [ ] **LP1 — der C#-NATS-Client samt Pin.** `NATS.Net` gepinnt im
      Projekt-/Paket-Manifest von `examples/csharp/`, mit Leser-Zeile
      („wofür sie da ist"); `examples/csharp/nats-client/` lauscht real,
      gibt das Ereignis aus, holt die Änderung über `GET /changes`; netzlos
      prüfbare Teile (Subjekt-Aufbau) sind getestet und laufen über
      `examples-csharp`.
- [ ] **LP2 — der Kotlin-NATS-Client samt Pin.** `io.nats:jnats` gepinnt im
      Build-Manifest von `examples/kotlin/`, mit Leser-Zeile;
      `examples/kotlin/nats-client/` — dieselbe Zusage, über
      `examples-kotlin`.
- [ ] **LP3 — die zwei Handbuch-Zeilen.** `docs/user/benutzerhandbuch.md` §4
      „Zugriff über das NATS-Wecksignal": der bestehende `**Beispiele:**`-
      Block (Go seit `slice-095`) bekommt zwei weitere Zeilen (C#, Kotlin),
      samt Änderungshistorie-Zeile.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-101.md`
      liegt vor (Modul 11, frischer Kontext).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) — *entfällt: Greenfield-
      Bootstrap (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = GF),
      keine Datei vorhanden.*
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
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
  **Ausgang:** <…>
- **Eine der beiden Werkzeugketten verlangt eine nicht-öffentliche Quelle**
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). — **Ausgang:** <…>
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (ein
  drittes Bau-Ziel je Sprache im selben Workflow) — §Re-Evaluierungs-
  Trigger 4, `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×,
  offen), `BEO-PGC/github-actions-unverifizierbar-lokal` (5×, verkörpert in
  `AGENTS.md` §3.10). — **Ausgang:** <…>
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen; sie ist der einzige Teil dieses Slice, den **kein**
  Kompilat erzwingt. — **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, verkörpert in `AGENTS.md`
  §3.13 — Suchlauf-Pflicht). — **Ausgang:** <…>

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

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
unverifizierbar-lokal` (**5×**, verkörpert in `AGENTS.md` §3.10) — beide
treffen den weiter wachsenden Workflow, als Risiko in §6. Kein Treffer zu
NATS.Net/jnats/NuGet/Maven selbst (gemessen: `grep -rli
"nats\.net\|jnats\|nuget\|maven" docs/plan/planning/observations/` → kein
Fund).

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
