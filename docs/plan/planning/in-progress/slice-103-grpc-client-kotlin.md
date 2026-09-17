# Slice slice-103: gRPC-Client in Kotlin — die in `slice-102` gemessene Form übertragen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf Nutzerentscheidung, ein einzelner Slice
ohne Closure-Bedingung jenseits seiner DoD (Modul 6 §Wann Arbeit eine Welle
braucht). Dieses Repo führt derzeit **keine** offene Welle
(`docs/plan/planning/in-progress/roadmap.md` §Offene Wellen ist leer).

**Bezug:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
(Festlegung 1 — Zelle `kotlin`×`grpc`; **Festlegung 2** — der `.proto`-Weg:
derselbe benannte Zusatzkontext, jetzt auf die Kotlin-Werkzeugkette
übertragen; **Festlegung 3** — der erzeugte Stub entsteht **im Bau**, wird
**nicht** committet; Festlegung 4 — `io.grpc:grpc-kotlin-stub`,
`io.grpc:protoc-gen-grpc-kotlin`, `io.grpc:grpc-netty-shaded` (Maven Central,
Registry-Existenz gemessen 2026-09-17); §Slice-Schnitt-Empfehlung, Zeile 6:
„**gRPC-Client in Kotlin** … dieselbe, jetzt **gemessene** Form auf den
Kotlin-Bau übertragen … Abhängigkeit 2, 5 (die Form ist in 5 belegt)") ·
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) (der Server-Stream
und die `.proto`) · `slice-102` (die Form-Referenz — Zusatzkontext,
Generator-Stufe — real gemessen und hier übertragen) ·
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `BEO-PGC/handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide treffen den Handbuch-Nachzug dieses Slice.

**Berührte Spec-Stellen:** `LH-FA-SST-008` (der Live-Change-Stream, hier über
gRPC) — dieser Slice **zeigt** ihn in der letzten der zwölf Zellen, er ändert
ihn nicht.

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

**Ziel:** `examples/kotlin/grpc-client/` öffnet real
`ChangeStream/StreamChanges` und gibt jede Nachricht aus (Form-Vorbild
`examples/grpc-client` in Go, Form-Referenz `examples/csharp/grpc-client` aus
`slice-102`). Dieser Slice **überträgt** die in `slice-102` real gemessene
Form der zwei Form-Fragen auf die Kotlin-Werkzeugkette, statt sie
unabhängig neu zu erfinden:

1. **Der benannte Zusatzkontext.** `examples/kotlin/Dockerfile` bekommt
   denselben Mechanismus wie `examples/csharp/Dockerfile` in `slice-102`
   (`docker buildx build --build-context <name>=<proto-Verzeichnis> …`,
   `COPY --from=<name> …`) — der Bau-Kontext selbst bleibt das
   Sprach-Wurzelverzeichnis, unverändert. Die Fehlermeldung bei fehlendem
   Kontext ist irreführend (Docker meldet „pull access denied" statt eines
   fehlenden Bau-Kontexts, `slice-102` Review-F-2) — der erklärende
   Kommentar direkt über der `COPY`-Zeile in `examples/csharp/Dockerfile`
   ist die Milderung dafür und wird beim Übertragen auf
   `examples/kotlin/Dockerfile` mitgeführt.
2. **Der Generator.** `io.grpc:protoc-gen-grpc-kotlin` erzeugt den Stub **im
   Bau** aus der kopierten `.proto`; `io.grpc:grpc-kotlin-stub` und
   `io.grpc:grpc-netty-shaded` tragen die Laufzeit. Der Stub liegt **nicht**
   im Baum ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
   Festlegung 3).

Der Handbuch-Abschnitt „Zugriff über den gRPC-Change-Stream" bekommt die
**dritte und letzte** Zeile (Kotlin) im `**Beispiele:**`-Block — mit diesem
Slice ist die volle Matrix (vier Oberflächen × drei Sprachen, zwölf
Programme) vollständig geliefert.

**Warum dieser Slice nach `slice-102`, nicht parallel:** Die Form
(Zusatzkontext + Generator) ist mit `slice-102` real im offiziellen
Sprach-Bau gemessen; dieser Slice **überträgt** sie, statt eine zweite,
unabhängige Bauprobe für Kotlin zu führen — dieselbe Disziplin wie bei den
Sprach-Wurzeln (`slice-098`/`slice-099`), nur diesmal mit einer bereits
belegten statt einer neu zu beweisenden Form.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Go-gRPC-Client (`examples/grpc-client/`).** Er ist **kein** Teil
  dieser Matrix-Erweiterung um C#/Kotlin — er läuft als eigener Liefer-Punkt
  (LP2) von `slice-095` (aktuell `in-progress/`), gebunden an den Umzug der
  Vertragsfläche (`slice-097`, `done/`, nicht an diesen Slice). Anderer
  Vorgang, andere Sprache, anderer Träger.
- **Eine erneute Bewertung des benannten Zusatzkontexts.** `slice-102` hat
  ihn real im Sprach-Bau gemessen; dieser Slice übernimmt die Form. Zeigt
  sich beim Übertragen ein Unterschied (siehe §4 Rückführung), ist das ein
  Befund dieses Slice, keine Wiederholung der Bewertung aus `slice-102`.
- **Der Inhalt der `.proto`.** Unverändert — dieser Slice **liest** sie über
  den Zusatzkontext, wie bereits `slice-102`; kein weiterer lesender Bauweg
  ändert die Quelle.
- **Ein committeter Kotlin-Stub.** Aus demselben Grund wie in `slice-102`:
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 3
  entscheidet gegen ein committetes Erzeugnis. Bestand bleibt bewusst stehen.
- **`.a-check.yml`.** Unverändert — dieselbe Begründung wie in `slice-102`:
  die `contract`-Gruppe bleibt Go, der erzeugte Kotlin-Stub liegt nicht im
  Baum. Bestand bleibt bewusst stehen.
- **Eine Zusammenfassung/Konsolidierung der zwölf Beispiel-Programme** (z. B.
  ein Übersichts-Dokument über die volle Matrix). Wäre ein anderer Vorgang —
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) legt die Form
  des Handbuch-Nachzugs bereits fest (ein `**Beispiele:**`-Block je
  Oberfläche, keine Matrix-Übersicht als eigenes Artefakt).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **LP1 — der `.proto`-Weg im Bau, übertragen aus `slice-102`.**
      `examples/kotlin/Dockerfile` liest die `.proto` über denselben
      Mechanismus (zusätzlicher, benannter Bau-Kontext); `make
      examples-kotlin` ruft den Bau mit diesem Zusatzkontext auf. Beleg: der
      Bau erzeugt den Stub real, bricht ohne den Kontext ab.
- [ ] **LP2 — der Client.** `examples/kotlin/grpc-client/` öffnet real den
      Server-Stream `ChangeStream/StreamChanges`, gibt jede Nachricht aus,
      liest Adresse/Token aus `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER`; netzlos
      prüfbare Teile sind getestet und laufen über `examples-kotlin`.
- [ ] **LP3 — die Handbuch-Zeile.** `docs/user/benutzerhandbuch.md` §4
      „Zugriff über den gRPC-Change-Stream": der `**Beispiele:**`-Block (Go
      seit `slice-095`, C# seit `slice-102`) bekommt die **dritte** Zeile
      (Kotlin), samt Änderungshistorie-Zeile — die volle Matrix ist damit im
      Handbuch vollständig.
- [ ] `make gates` grün.
- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-103.md` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-103.md`
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
| `examples/kotlin/Dockerfile` | update | Generator-Stufe (`protoc-gen-grpc-kotlin`) und der zusätzliche, benannte Bau-Kontext für `proto/` — Form übertragen aus `examples/csharp/Dockerfile` (`slice-102`). |
| `examples/kotlin/<Build-Manifest>` | update | Pin von `io.grpc:grpc-kotlin-stub`, `io.grpc:protoc-gen-grpc-kotlin`, `io.grpc:grpc-netty-shaded`. |
| `examples/kotlin/grpc-client/**` | neu | gRPC-Server-Stream-Client, Form-Vorbild `examples/grpc-client` (Go) und `examples/csharp/grpc-client` (`slice-102`). |
| `Makefile` / `harness/mk/*.mk` | update | `examples-kotlin`-Ziel ruft den Bau mit dem benannten Zusatzkontext auf. |
| `docs/user/benutzerhandbuch.md` | update | §4 „Zugriff über den gRPC-Change-Stream": dritte Zeile (Kotlin) im `**Beispiele:**`-Block; Änderungshistorie-Zeile. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-099` **und** `slice-102` liegen in
`done/` — die Kotlin-Sprachwurzel samt Werkzeugkette existiert
(`slice-099`), und die Form des `.proto`-Wegs (Zusatzkontext, Generator) ist
im C#-Bau real gemessen (`slice-102`) und damit übertragbar. Ohne Rückfrage
feststellbar (Verzeichnis-Position beider Vorgänger-Slices).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die aus
  `slice-102` übertragene Form auf der Kotlin-/JVM-Werkzeugkette **nicht**
  unverändert trägt (z. B. weil `protoc-gen-grpc-kotlin` einen anderen
  Kontext-Zugriff braucht als `Grpc.Tools`) und eine eigenständige,
  abweichende Lösung nötig wird — dann ist die Übertragungs-Annahme falsch
  und der Slice ist neu zu bewerten.
- `in-progress` → `open` (blockiert — Carveout?): wenn eine der drei
  gRPC-Kotlin-Bibliotheken nicht mehr auflösbar ist oder die
  Kotlin-gRPC-Werkzeugkette eine nicht-öffentliche Quelle voraussetzt
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — der Bau erzeugt
den Kotlin-Stub über den übertragenen Zusatzkontext reproduzierbar, der
Client öffnet real den Server-Stream und gibt Nachrichten aus, die
Handbuch-Zeile trägt Kotlin als dritte und letzte Sprache der Zelle — **und**
`make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegender Kandidat: ob die Übertragung der `slice-102`-Form
unverändert trägt oder eine sprachspezifische Anpassung brauchte — das ist
der letzte Beleg dafür, ob „eine Form, zwei Sprachen" (statt einer
sprachneutralen `contract`-Gruppe, [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
Festlegung 6) für den benannten Zusatzkontext wirklich trägt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die aus `slice-102` übertragene Form trägt auf der Kotlin-Werkzeugkette
  nicht unverändert** (abweichender Kontext-Zugriff des Kotlin-/JVM-Bauwerks
  gegenüber dem C#-Bauwerk). — **Ausgang:** <…>
- **Eine der drei gRPC-Kotlin-Bibliotheken ist zum Zeitpunkt des Baus nicht
  mehr auflösbar** (Paket zurückgezogen, Version gelöscht). —
  **Ausgang:** <…>
- **Die Kotlin-gRPC-Werkzeugkette verlangt eine nicht-öffentliche Quelle**
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). — **Ausgang:** <…>
- **Der fremdsprachige gRPC-Bau kann seinen Stub nicht mehr aus der `.proto`
  erzeugen** — [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 2 (dieselbe Bedingung wie in `slice-102`, hier für
  die Kotlin-Route geprüft). — **Ausgang:** <…>
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** — nach
  diesem Slice deckt er **alle** vier Bau-Ziele beider Sprachen; §Re-
  Evaluierungs-Trigger 4, `BEO-PGC/nicht-blockierender-workflow-
  alarmmuedigkeit` (1×, offen), `BEO-PGC/github-actions-unverifizierbar-lokal`
  (5×, verkörpert in `AGENTS.md` §3.10). — **Ausgang:** <…>
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen; mit diesem Slice wird die Matrix im Handbuch
  vollständig, ein vergessener Nachzug wäre hier am sichtbarsten. —
  **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, verkörpert in `AGENTS.md`
  §3.13 — Suchlauf-Pflicht; insbesondere Träger, die die Matrix noch als
  „unvollständig" beschreiben, z. B. [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)s
  eigener Ist-Stand-Abschnitt oder `SPEC-023` „Sprachen und Umfang", falls bis
  dahin nicht bereits nachgezogen). — **Ausgang:** <…>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `examples/kotlin/**`,
`proto/` (nur **lesend**, über den Zusatzkontext), `Makefile`/
`harness/mk/*.mk` und `docs/user/benutzerhandbuch.md` — die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Keine zu
grobe Sub-Area zu differenzieren.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Fünf
Treffer: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide als LP3 in die DoD gezogen; `nicht-blockierender-
workflow-alarmmuedigkeit` (**1×**, offen) und `github-actions-
unverifizierbar-lokal` (**5×**, verkörpert in `AGENTS.md` §3.10) — beide
treffen den nun vollständigen Workflow, als Risiko in §6;
`arbeit-ueberholt-stehenden-traeger` (verkörpert in `AGENTS.md` §3.13) — als
Risiko in §6, weil dieser Slice die letzte Zelle liefert und mehrere Träger
(„Matrix unvollständig") mit ihm falsch werden könnten. Kein Treffer zu
`grpc-kotlin-stub`/`protoc-gen-grpc-kotlin`/`grpc-netty-shaded` selbst
(gemessen: `grep -rli "grpc-kotlin\|netty-shaded"
docs/plan/planning/observations/` → kein Fund).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**.
Der Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

### Sub-Area: `*` (Default, PGC)

- **Modus:** GF
- **Konventionen-Dichte:** `harness/conventions.md` Modus-Deklaration setzt
  GF für das gesamte Repo (Doc führt, Code folgt).
- **Phase-Reife:** Phase 5 (etabliert), sofern `slice-102` die Form real
  bestätigt — dieser Slice überträgt eine bereits gemessene Form, statt eine
  neue zu beweisen.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — GF, Doc führt; das
  Haupt-Diskrepanz-Risiko (trägt die Form auf der zweiten Sprache?) ist in
  §6 als Risiko geführt, nicht als stille Annahme.
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bootstrap).
