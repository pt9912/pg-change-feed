# Verifikationsbericht: slice-077 — 2026-09-15

**Rolle:** Verifier (Baseline-Regelwerk `modul-11-verifikation.md`) — „Bauen wir
es richtig?" gegen Plan (`slice-077` §1–§8) und die bindenden Entscheidungen
(`ADR-0057`, `ADR-0059`, `ADR-0060`, `ADR-0061`, `ADR-0065`, `ADR-0066`;
`AGENTS.md` §3.7, §3.9). **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe;
`review-slice-077` und `review-slice-077-delta` als Kontext gelesen) und **nicht**
gegen realen Bedarf (Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan, die vier
Spec-Stellen (`SPEC-018`/`019`/`020`/`021`), die Entscheidungen, den Review-Erst-
und Delta-Report und die berührten Artefakte. **Alle** Sensoren und Messungen
wurden in eigener Sitzung ausgeführt; kein Beleg des Implementers, Reviewers oder
Planners wurde übernommen. Exit-Codes je in eigenem, ungepiptem Schritt
(`AGENTS.md` §3.9).

**Messbedingung — der Arbeitsbaum hat sich während dieses Laufs bewegt.**
Verifiziert ist `HEAD = bdad248`. Während des Laufs lag zuerst ein **nicht
committeter** Planner-Zug im Baum (Register-Belege, N-1/N-2-Nachzüge); er wurde
mitten im Lauf als `bdad248` committet und ist damit **Teil** des geprüften
Standes. Der Baum ist danach sauber (`git status --porcelain` leer). Das
**Handbuch-Artefakt ist von `bdad248` nicht berührt** — sein Stand ist der aus
`1cf5675`.

**Gegenstand:** `slice-077`, Range `afaf5e4..bdad248`. Sechs Commits:
`10a1221` (Handbuch), `d7db513` (Review-Erstlauf), `3637e1f` (Planner-Plan-
Nachzug), `1cf5675` (Fixrunde), `8b4925e` (Delta-Review), `bdad248`
(Closure-Vorbereitung: Register-Belege, N-1/N-2). Geliefert werden zwei Dateien
(`docs/user/benutzerhandbuch.md` inhaltlich, `internal/bootstrap/wiring.go`
kommentar-only).

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt (§2) | Verdikt | Beleg (eigene Prüfung am Artefakt) |
|---|---|---|---|
| 1 | §5 trägt HTTP-Gruppe + `CDC_GRPC_ADDR`, je mit Aktivierungs-/No-Op-Semantik, gegen `wiring.go` | **erfüllt** | Vier Zeilen real vorhanden (Handbuch `:715`–`:718`). Semantik **selbst** am Code geprüft: `CDC_HTTP_ADDR` leer → kein Server/kein Pool (`wiring.go:697`-`726` nur im `if`-Zweig); gesetzt → Goroutine + `log.Error`, kein `Run`-Beitrag (`:719`-`726`); `CDC_API_TOKEN_READER`/`_ADMIN` leer → Klasse trifft nie (`middleware.go:34`-`45`), unbekanntes Token → `withToken` `roleNone` → `401` (`:78`-`84`), Admin deckt Reader implizit (`roleAdmin < roleReader`-Rang, `:85`); `CDC_GRPC_ADDR` leer → kein Listener, gesetzt → Goroutine + `log.Error` (`wiring.go:739`-`754`). **Semantik richtig, nicht nur belegt** (Einschränkung: O-1). |
| 2 | Der Deklarationsort trägt die Semantik; dieselbe Behauptung nicht an anderer Stelle | **erfüllt** | `wiring.go:113`-`115` sagt jetzt: offener Listener in eigener Goroutine, Startfehler über `log.Error`, **nicht** im `Run`-Ergebnis. Deckt `:747`-`752` (Goroutine + `log.Error`) und `:826` (einziger Nicht-Fehler-`return`). Der alte Satz ist restlos ersetzt: repo-weiter `grep` über die Aussageklasse findet „Ergebnis von `Run`" nur noch an `:113`-`115` (richtig) und `:581`-`582` (Nachbarblock, richtig). Kein Verhaltens-Change im Diff (nur Kommentarzeilen). |
| 3 | §4 trägt „Spalte vom Ausschluss konfigurieren" samt **dauerhaftem** Träger (`ADR-0065`) | **erfüllt** | Handbuch `:249`-`304`. Ist-Zustand, keine überholte Grenze: „Die `applied`-Zeilen … sind die einzige Herkunft des Standes" / „Prozessstart **und** laufende Aktivierung" / „Neustart verliert ihn nicht" / „Bereinigung verlöre den Stand" — Satz für Satz `SPEC-019` (`:339`-`352`). Im Code: `wiring.go:347`-`372` (`activatedTableBindings` samt `ExcludedColumns`) und `:1172`-`1180` (Aktivierungs-Zweig liest `columnExclusion`). `cdc.exclude_column(...)`-Signatur deckt `tools/schema/nacharbeit-administration.sql:99`. |
| 4 | §4 trägt die **drei** Zugriffswege mit Erreichbarkeit, Authentifizierung, Zustellsemantik, Nachvollziehbarkeit | **erfüllt** | HTTP `:582`-`624`, gRPC `:626`-`670`, SSE `:672`-`698` — je alle vier Hälften; Nachvollziehbarkeit in gRPC `:667`-`670` und SSE `:696`-`698`. **Zustellsemantik** deckt `SPEC-020` (`:371`) und `SPEC-021` (`:389`-`390`) wörtlich: keine Garantie, **begrenzte** Empfangswarteschlange je Abonnent, Überlauf **verworfen**, kein Replay, Erzeuger hält nie an — `broadcaster.go:36` (`queueCapacity = 64`), `:96`-`117` (`select`/`default`, Drop-Newest). **Zahl der Nachrichtenfelder: zehn** — Handbuch-Tabelle `:643`-`654` (10 Zeilen) deckt `proto/…/changestream.proto:17`-`28` (genau 10 Felder) und die `SPEC-020`-Zeile *Nachricht `Change`* (`:369`). SSE-Schema `sse.go:32`-`43` trägt dieselben 10 JSON-Schlüssel. |
| 5 | Versionshistorie fortgeschrieben | **erfüllt** | `Version: 1.14` (`:3`), `Stand: 2026-09-15` (`:5`), zwei neue Zeilen 1.13/1.14 (`:889`-`890`) mit den inhaltlichen Änderungen; Fixrunde zog 1.14 nach. |
| 6 | `make gates` grün | **nicht erfüllt am geprüften HEAD** — Ursache außerhalb des Slice (s. §2) | Eigener Lauf: **Exit 2**; **einziger** Befund `3637e1f commit-untraceable` (Planner-Commit, kein `LH-*`/`ADR-*` im Betreff). Alles Slice-Tragende ist im selben Lauf grün. Kein V-Befund gegen das Artefakt, aber das Häkchen ist am HEAD nicht wörtlich wahr (s. O-6). |
| 7 | Review durchgeführt, Report liegt vor | **erfüllt** | `review-slice-077.md` (1 HIGH) und `review-slice-077-delta.md` (Fixrunde `1cf5675`, 0 HIGH) existieren; die Delta-Report-Zeile „DoD-Nachzug" belegt den Rollenwechsel. Rollen-Trennung ist aus Artefakten allein nicht beweisbar; beide Reports tragen **eigene** Läufe (Mutationsreihen, Grep-Proben), konsistent mit getrennten Kontexten. |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich Platzhalter (gelesen). Planner-Arbeit. |
| 9 | Reconciliation-Register — **entfällt** | **erfüllt (Entfall trägt)** | `docs/plan/planning/reconciliation.md` existiert real **nicht**. |
| 10 | Beobachtungs-Register fortgeschrieben | **erfüllt** (mit einer offenen Zusatz-Zusage, s. §7) | `bdad248` legt an: `BEO-PGC/commit-traceability-kein-vorab-hook/evidence/slice-077.md` (+ `state.md` auf 4×) und den neuen Eintrag `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (2×, unter der Schwelle). **Nicht** angelegt: der von §8 zugesagte Beleg in `BEO-PGC/aufschub-adresse-verfaellt` (s. §5 und §7). |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Einträge tragen wörtlich `<bei Closure>`; keines vorzeitig geschlossen. |
| 12 | Die drei Paarungen | **korrekt offen** | Die Roadmap führt real keine offene Welle; die Prüfung trägt damit korrekt die Slice-Closure selbst. |

**Ergebnis §1:** Die **sieben** gesetzten Zeilen (1–5, 7, 9) sind real erfüllt —
jede in dieser Sitzung am Code/Artefakt nachgeprüft. Zeile 10 ist über die zwei
neuen Belege erfüllt. Zeile **6** ist am HEAD **nicht** wörtlich wahr, aber
**ausschließlich** durch den Planner-Commit `3637e1f`, nicht durch das Artefakt.
Die vier Planner-Posten (8, 11, 12 und die Register-Zusatz-Zusage) sind korrekt
offen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` (`HEAD = bdad248`) | **2** | `baseline-verify v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — 49.30 % erfüllt Schwelle 40 %` · `d-check 641 Dateien / 0 Befunde` · **Abbruch** bei `commit-traceability` |
| Befundliste des roten Laufs | — | **genau ein** Befund: `3637e1f commit-untraceable` — die Kennung fehlt im Betreff von „docs(planning): slice-077 Nachzug — Kommentar-Berichtigung als Liefer-Punkt 1 (Review F-1)" |
| `make a-check` (einzeln, wegen des Abbruchs) | **0** | `gesamt: 0 Befund(e)` |
| `make doc-commits RANGE=afaf5e4..HEAD` | **2** | 0 Befunde außer `3637e1f` — jeder andere Commit des Fensters ist kennungstragend |
| `make doc-immutable RANGE=afaf5e4..HEAD` | **0** | 0 Befunde — keine `Accepted`-ADR überschrieben |

**Die Vorgabe ist bestätigt: der rote Befund ist ausschließlich `3637e1f`.** Die
drei weiteren von dieser Rolle geforderten Sensoren (`baseline-verify`,
`docs-check`, `coverage-gate`) sind im `gates`-Lauf vor dem Abbruch durchgelaufen
und grün. **`a-check` läuft wegen des Abbruchs an der Traceability nicht mit** —
einzeln nachgeholt, Exit 0.

## 3. Plan-vs-Artefakt, beide Richtungen

**Plan → Artefakt (jede Behauptung am Handbuch geprüft):**

| Plan-Behauptung | Befund |
|---|---|
| §1/§2 Liefer-Punkt: §5 um HTTP-Gruppe + `CDC_GRPC_ADDR` | **trägt** — vier Zeilen, Semantik gegen `wiring.go`/`middleware.go` geprüft (§1/1) |
| §2 Liefer-Punkt: §4 „Spalte vom Ausschluss konfigurieren" samt dauerhaftem Träger | **trägt** — `:249`-`304`, deckungsgleich `SPEC-019` |
| §2 Liefer-Punkt: §4 drei Netzwerk-Zugriffswege | **trägt** — alle drei Abschnitte mit den vier geforderten Hälften |
| §2 Liefer-Punkt: Versionshistorie | **trägt** — 1.14 + Zeile |
| §3 Zeile 1: `docs/user/benutzerhandbuch.md` update | **trägt** — im Diff `afaf5e4..bdad248` |
| §3 Zeile 2: `internal/bootstrap/wiring.go` update, **kein Verhaltens-Change** | **trägt** — der Diff in `wiring.go` ändert ausschließlich Kommentarzeilen (kein Hunk mit Anweisungszeile) |
| §1 Ausschluss „Änderungen an `SPEC-018`/`SPEC-020`" — der Slice ändert keine Spec-Stelle | **trägt** — `git diff --name-only afaf5e4..bdad248 -- spec/` ist **leer** |
| §8 Treffer „`BEO-PGC/aufschub-adresse-verfaellt` (1×, offen) — Beleg bei Closure: `evidence/slice-077.md`" | **nicht getragen** — die Datei existiert real nicht (Register führt dort nur `slice-074.md`, `slice-076.md`; Zähler 2×). Offene Closure-Obliegenheit (§7/2). |

**Artefakt → Plan (Zustände im Handbuch, die der Plan nicht ausspricht):**

- **Die Zehn-Feld-Tabelle (§4 gRPC `:643`-`654`).** Der Plan nennt sie **nicht**
  (DoD verlangt „Zustellsemantik"), sie kam über Review F-3. **Geprüft und
  richtig**: namentlich und typlich deckungsgleich mit `changestream.proto` und
  `SPEC-020`; „bei `INSERT` leer"/„bei `DELETE` leer" entspricht dem Null-Fall.
  **Kein V-Befund.**
- **Die Ausnahme „Aufbewahrung" in der HTTP-Rahmen-Aussage (§4 `:584`-`590`).**
  Der Plan nennt sie nicht (kam über Review F-2). **Geprüft und richtig**: die
  CLI kennt real **fünf** Modi ohne Retention (`cmd/pg-change-feed/main.go:22`,
  `:26`, `:46`, `:64`, `:88`), `tools/schema/` führt genau **vier** `cdc.*`-
  Funktionen (`enable`/`disable`/`exclude`/`include`) und **keine**
  Retention-Auslösung. **Kein V-Befund.**
- **Die `503`-Zeile (§4 SSE `:686`-`688`).** Der Plan nennt sie nicht. Sie
  spiegelt `SPEC-021` (`:391`) und `sse.go:89`-`92` **wörtlich** — **richtig,
  aber aus dem Container-Vertrag unerreichbar** (bei gesetztem `CDC_HTTP_ADDR`
  ist der `Broadcaster` verdrahtet). **Kein V-Befund** (dieselbe Beobachtung wie
  Review F-4; benannt, nicht falsch).
- **„Eine Filterung nach Tabelle ist nicht Teil dieser Version"** (`:656`) und
  **der Ausschluss-Betriebs-Hinweis** (`:299`-`304`, „ohne laufende Erfassung →
  `applied`") — beide nicht im Plan, beide deckungsgleich mit `SPEC-020`
  (`:368`) bzw. `SPEC-019` (`:349`-`351`). **Kein V-Befund.**

## 4. Entscheidungs- und Spec-Konformität

1. **`spec/**` unberührt** — `git diff --name-only afaf5e4..bdad248 -- spec/`
   leer. §1-Abgrenzung hält am Diff.
2. **Handbuch konsistent mit `SPEC-018`**: Endpunkt-Tabelle `:606`-`616` — neun
   Zeilen, Methode/Pfad/Rechtsklasse deckungsgleich mit der `SPEC-018`-Tabelle
   (`spec/pflichtenheft.md:298`-`308`) und `http/server.go:79`-`98`;
   `401`/`403`-Regel (`:596`-`597`, `:618`-`619`) stimmt mit `withToken`
   (`middleware.go:78`-`91`); Fehlerform `{"error": "<Klartext>"}` deckt
   `SPEC-018:290`-`296`.
3. **`SPEC-019`**: Dauerhaftigkeits-Absatz `:289`-`297` und Betriebs-Hinweis
   `:299`-`304` spiegeln `:339`-`352` und `ADR-0065` §Entscheidung 1–3; im Code
   `wiring.go:347`-`372` und `:1172`-`1180`.
4. **`SPEC-020`/`SPEC-021`**: Zustellsemantik, Nachrichtenfelder (zehn),
   Aktivierung, `Unauthenticated`, `Last-Event-ID`, `text/event-stream`,
   `event: change`, `null`-Row-Image — je wörtlich gedeckt.
5. **`ADR-0057` Teilfrage 3** (Token-Klassen orthogonal zum DB-Rollenmodell):
   Handbuch `:602`-`604` sagt genau das. `ADR-0066`: „begrenzte
   Empfangswarteschlange, Überlauf verworfen" — gedeckt.

**Urteil zur aufgelösten Spannung (Review F-5).** Die Changelog-Wendung des
Lastenhefts 0.7.0 („`LH-FA-CFG-005` … von dauerhaftem Ausschluss auf aktive
Anforderung umgestellt") meint den **Scope der Anforderung**, **nicht** die
Dauerhaftigkeit eines Ausschlussstandes. Am Belegstand `893fc95:spec/lastenheft.md`
nachgeprüft: dort stand „### `LH-FA-CFG-005` — Spaltenauswahl **(perspektivisch)**"
mit der Out-of-Scope-Klausel „Kein Bestandteil des MVP; eine nachträgliche
Ergänzung ohne Neuanforderung ist nicht vorgesehen"; die heutige Fassung trägt
den Titel ohne Zusatz und eine Out-of-Scope-Klausel über ADR/`SPEC-*`. Die
Wendung „dauerhafter Ausschluss" adressiert also die *dauerhafte Ausschließung
der Anforderung* (perspektivisch/out-of-scope), nicht den *Spaltenausschluss*.
Das Wort ist deutsch doppeldeutig und hat die Prüf-Hypothese des Plans (§6
Risiko 2) erzeugt. **Die Begründung des Reviewers trägt; die Handbuch-Fassung ist
unter Source Precedence richtig.**

## 5. Betreiber-Doku-Hygiene (`AGENTS.md` §3.7, §3.11)

**Ohne Befund.** Die **hinzugefügten** Zeilen tragen keine `slice-`/`welle-`
Kennung, keinen host-lokalen Pfad und keine Chronik-Sprache (`grep` über die
Diff-`+`-Zeilen: kein Treffer). Der §4-Dauerhaftigkeits-Absatz schreibt den
Ist-Zustand in der Gegenwart („Ein Neustart … verliert den Ausschluss deshalb
**nicht**"), nicht die abgelöste Grenze. **Anmerkung ohne Wirkung:** die
*vorbestehenden* Historie-Zeilen 1.5–1.12 (`:881`-`888`) führen `slice-`-
Kennungen — nicht Gegenstand dieses Slice, aber die neuen Zeilen 1.13/1.14
weichen damit vom Muster der Nachbarzeilen ab (zugunsten der Hygiene-Regel).

## 6. Beobachtungen ohne DoD-Wirkung (kein V-Befund)

- **O-1 — `401` als zusammenfassende Antwort für zwei Träger.** Die §5-Zeile
  `CDC_API_TOKEN_READER` (`:716`) nennt „die HTTP- und gRPC-API" und sagt dann
  „endet `401`". Für den gRPC-Träger ist der reale Ausgang der Status
  `Unauthenticated` (`interceptor.go:91`-`98`, `SPEC-020:374`), nicht HTTP `401`.
  Der §4-gRPC-Abschnitt (`:632`-`636`) sagt das **richtig**; die §5-Zusammenfassung
  ist für die gRPC-Hälfte lose. Keine DoD-Zeile verletzt (verlangt sind die vier
  Zeilen mit Aktivierungs-/No-Op-Semantik — die tragen).
- **O-2 — „dieselbe Deaktivierung" beim Admin-Token.** `:717` sagt für leeres
  `CDC_API_TOKEN_ADMIN` „dieselbe Deaktivierung wie bei `CDC_API_TOKEN_READER`".
  Wirkungsgleich für den Aufrufer, aber nicht identisch: leeres Admin-Token lässt
  administrative Endpunkte mit einem bekannten Reader-Token `403` antworten, statt
  zu fehlen. Unschärfe ohne Falschaussage (so auch Review F-1-Umfeld).
- **O-3 — `503` unerreichbar aus dem Container-Vertrag.** S. §3; deckungsgleich
  mit Review F-4 (INFO).
- **O-4 — Feldtabelle übernimmt die Werteform des zweiten Trägers nicht.** Die
  SSE-Sektion verweist für die zehn Felder auf die gRPC-Tabelle (dort `bytes`);
  im Event sind die Row Images eingebettete JSON-Werte bzw. `null`
  (`sse.go:38`-`39`, `SPEC-021:388`). Die **Namen** sind identisch — nur die
  Typ-Spalte gilt für den zweiten Träger nicht. Deckungsgleich mit Delta-Review
  N-3 (INFO).
- **O-5 — Zählung der Liefer-Punkte.** §2 nennt (mit dem review-erzwungenen
  Deklarationsort) vier Inhalts-Punkte gegen die eigene Regel „≤ 3" — der
  Reviewer hat das benannt und begründet (review-erzwungene Korrektur, kein
  selbst gewählter Umfang). Keine DoD-Zeile.
- **O-6 — das Häkchen „`make gates` grün" am HEAD.** Zeile 6 ist gesetzt, der
  Gate-Lauf endet aber Exit 2 (nur `3637e1f`). Das ist **kein** Befund des Slice —
  es ist die Closure-Buchhaltung: entweder verlässt `3637e1f` das
  `HEAD~5..HEAD`-Fenster, oder die Zeile wird beim Closure-Zug qualifiziert
  („grün bis auf den benannten Planner-Commit").

## 7. Offene Closure-Obliegenheiten (benannt, nicht ausgeführt)

**Bereits erledigt (in `bdad248`, nach den Review-Reports):** N-1 (Wortlaut in
§1: „geltendes Verhalten des Codes … eine ADR trägt es **nicht**") · N-2 (Lücke in
`harness/sensors/docs-check.md` §Grenze, Punkt 9) · die zwei Register-Belege
(`BEO-PGC/commit-traceability-kein-vorab-hook` 4×; `BEO-PGC/kommentar-behauptet-
nicht-getragenen-fehlerpfad` neu, 2×).

**Offen — Planner-Arbeit:**

1. **Zwei §6-Risiko-Ausgänge** — Risiko 1 (Aktivierungs-Semantik) und Risiko 2
   (Ausschlussstand als dauerhaft statt prozessgebunden): je genau **ein** Ausgang.
   Beide sind durch diesen Bericht belegt — Risiko 1 durch die Code-Prüfung in
   §1/1, Risiko 2 durch §1/3 (Ist-Zustand, `SPEC-019`/`ADR-0065`) — also
   **entfallen** mit Beleg, nicht „weiter offen" ohne Grund.
2. **Register-Beleg `BEO-PGC/aufschub-adresse-verfaellt`** — plan §8 sagt
   wörtlich `evidence/slice-077.md`, Zähler dann **3×** (`slice-074`, `slice-076`,
   `slice-077`), Ausgang im Lese-Schritt. Die Datei fehlt am HEAD (real 2×). Ohne
   sie bleibt der Eintrag unter der Schwelle und der Lese-Schritt hat nichts
   zuzuweisen.
3. **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
   Sensor · benannte Spec-Lücke) und der Beobachtungs-Register-Zeile (zitiert die
   Kennungen, legt keinen Zähler von Hand).
4. **`git mv` nach `done/`** — erst Inhalt (Häkchen, Closure-Notiz), dann der
   reine Move (`AGENTS.md` §3.3).
5. **Drei Paarungen** (Anker · Folge-Slice · Register) — von der Slice-Closure
   selbst getragen (die Roadmap führt keine offene Welle); die Register-Hälfte
   prüft genau den unter 2. genannten Beleg.
6. **O-6 auflösen** — das Häkchen „`make gates` grün" am Closure-Stand
   qualifizieren oder das Fenster `3637e1f` verlassen lassen.

## Verdikt

**DoD-Konformität:** **bestätigt** — die sieben gesetzten Inhalts-Zeilen (1–5, 7,
9) sind am Artefakt selbst nachgeprüft und real erfüllt; die Register-Zeile ist
über die zwei neuen Belege erfüllt. Zeile 6 (`make gates` grün) ist am HEAD nicht
wörtlich wahr, aber **ausschließlich** durch den Planner-Commit `3637e1f` — kein
Befund des Slice (eigener Lauf belegt: kein weiterer Befund). Alle fünf inneren
Gates außer der Traceability sind grün (`a-check` einzeln Exit 0).

**Entscheidungs-/Spec-Konformität:** **bestätigt** — `spec/**` unberührt, das
Handbuch konsistent zu `SPEC-018`/`019`/`020`/`021`, zu `ADR-0057`/`0060`/`0061`/
`0065`/`0066`. Zur aufgelösten Spannung (Review F-5): **die Wendung meint den
Scope**; die Handbuch-Fassung trägt unter Source Precedence.

**Plan-vs-Artefakt:** in der Hauptrichtung deckungsgleich; in der Gegenrichtung
fünf ungenannte, aber sämtlich **richtige** Zusätze ohne V-Befund. Die einzige
**nicht getragene** Plan-Zusage ist der §8-Beleg in
`BEO-PGC/aufschub-adresse-verfaellt` — eine offene Closure-Obliegenheit, keine
DoD-Verletzung.

**Kein V-Befund.** Sechs Beobachtungen ohne DoD-Wirkung (O-1…O-6) benannt.

Der `git mv` nach `done/`, die Closure-Notiz, die zwei Risiko-Ausgänge, der
Register-Beleg und die drei Paarungen bleiben Planner-Arbeit.
