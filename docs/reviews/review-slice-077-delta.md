# Review-Report: slice-077 — Delta-Nachlauf zur Fixrunde — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** Fixrunde `1cf5675` (drei Pfade: `internal/bootstrap/wiring.go`
Kommentar, `docs/user/benutzerhandbuch.md` §4/Version, Slice-Plan §2/§3) —
als **frisches Artefakt** gelesen, nicht als Bestätigung des Erstlaufs.
Zusätzlich gelesen, ohne ihn zu bewerten: der Planner-Nachzug `3637e1f`
(Plan-Pfad) und der Erstlauf-Report `docs/reviews/review-slice-077.md`.

**Gate-Status dieses Laufs — rot, aber nicht wegen dieses Slice.**
`make gates` endet mit Exit 2; der **einzige** Befund ist
`3637e1f:1  commit-untraceable` (Betreff „… Nachzug — Kommentar-Berichtigung
als Liefer-Punkt 1 (Review F-1)" trägt keine `LH-*`/`ADR-*`-Kennung). Der
Fehler liegt beim Planner, ist entschieden (kein Umschreiben, aus dem
`HEAD~5..HEAD`-Fenster laufen lassen) und wird hier **nicht** als Finding
geführt. Alles, was diesen Slice trägt, ist im selben Lauf grün: docs-check
635 Dateien/0 Befunde, `coverage-gate` 49,30 % ≥ 40 %, baseline-verify OK;
`a-check` läuft wegen des Abbruchs an der Traceability nicht mehr mit und
wurde einzeln nachgeholt (Exit 0, 0 Befunde).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd`.
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Erstlauf-Report `docs/reviews/review-slice-077.md` (F-1 HIGH, F-2/F-3 MEDIUM,
  F-4/F-5 INFO)
- Slice-Plan `docs/plan/planning/in-progress/slice-077-handbuch-betreiber-stand.md`
  §1/§2/§3 im Stand nach `3637e1f` und `1cf5675`
- `spec/pflichtenheft.md` `SPEC-018`/`SPEC-019`/`SPEC-020`/`SPEC-021`,
  `spec/lastenheft.md` `LH-FA-CFG-005`/`SST-006`/`SST-008`
- [`ADR-0057`](../plan/adr/0057-http-grpc-api.md),
  [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) §Teilfrage 6,
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md),
  [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
- `proto/cdc/stream/v1/changestream.proto`,
  `internal/adapters/driving/grpc/server.go`,
  `internal/adapters/driving/http/sse.go`, `internal/bootstrap/wiring.go`,
  `cmd/pg-change-feed/main.go`,
  `tools/schema/nacharbeit-administration.sql`, `tools/schema/nacharbeit-roles.sql`
- `AGENTS.md` §3.7/§3.9/§3.11, `harness/sensors/docs-check.md` §Grenze,
  `.d-check.yml` (`modules`, `anchors`)

---

## Verdikt je Erstlauf-Finding

| Finding | Verdikt | Eigener Beleg (nicht aus der Commit-Message übernommen) |
|---|---|---|
| **F-1 (HIGH)** | **behoben** | `wiring.go:113-115` nennt jetzt „öffnet den Listener in eigener Goroutine (`Run` unten)" und „geht **nicht** in das Ergebnis von `Run` ein"; Verhalten unverändert (`:747-752` Goroutine + `log.Error`, `:825` Rückgabewert) — die drei Teilfragen unten |
| **F-2 (MEDIUM)** | **behoben** | `handbuch:584-590` nennt die Ausnahme namentlich („kein CLI-/SQL-Zugriffsweg … einziger manueller Auslöser"); eigener Nachweis der Belege unten |
| **F-3 (MEDIUM)** | **behoben** | `handbuch:641-654` trägt die Zehn-Feld-Tabelle; sie stimmt namentlich **und** typlich mit `changestream.proto:17-28` überein (unten Zeile für Zeile) |
| F-4 (INFO) | unverändert gültig, kein Nachzug erwartet | `handbuch:686-688` (die `503`-Zeile steht weiter, jetzt einzeilig eingebettet) |
| F-5 (INFO) | unverändert gültig | `spec/lastenheft.md:1245` unberührt (Rang-1-Wortlaut, nicht im Zug) |

## Findings dieses Laufs

### N-1 — Der Plan nennt das Listener-Verhalten „die geltende Entscheidung"; kein Artefakt trägt sie

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-04-adrs.md`
  (§ADR = der Träger einer Entscheidung) · `AGENTS.md` §5 („Neue oder geänderte
  Anforderungen brauchen einen Beleg") — Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-077-handbuch-betreiber-stand.md:57`
  (aus `3637e1f`) gegen `docs/plan/adr/0060-grpc-streaming-mechanismus.md:205-212`,
  `spec/pflichtenheft.md:375,391`
- `befund`: Der neue §1-Ausschluss sagt, dass ein gescheitertes Binden nur
  geloggt wird und den Lauf nicht beendet, sei „die geltende Entscheidung".
  Gefunden wurde kein Artefakt, das sie trägt: `ADR-0060` Teilfrage 6
  entscheidet die **Aktivierung** (additive Variable, No-Op bei leerem Wert),
  die Aktivierungszeilen von `SPEC-018`/`SPEC-020`/`SPEC-021` ebenfalls; die
  Fehlersemantik des Listeners steht im Code (`wiring.go:747-752`) und — seit
  diesem Slice — im Handbuch. Die Aussage des Ausschlusses selbst trägt (der
  Fixrunde-Diff ändert in `wiring.go` ausschließlich Kommentarzeilen, kein
  Verhalten); lose ist die Zuordnung „Entscheidung".
- `verifizierbar`: nein — Zuordnungs-/Wortlautfrage im Zeitdokument; die
  Abwesenheit ist per Grep über `spec/`, `docs/plan/adr/`, `harness/` belegt
- `klasse`: „Plan nennt ein Verhalten ‚geltende Entscheidung', ohne dass ein Artefakt sie trägt"

### N-2 — Benannte Gate-Reichweiten-Lücke: ein Link, dessen Text oder Ziel über einen Zeilenumbruch geht, ist für `links` **und** `anchors` unsichtbar

- `kategorie`: INFO
- `quelle`: `harness/sensors/docs-check.md` §Vertrag („lokaler Link oder
  Heading-Anker ins Leere … `target-missing`, `anchor-missing`") gegen
  §Grenze (dort **nicht** benannt) · `AGENTS.md` §3.7 (benannte Lücke statt
  stiller) — Nebenbefund des Implementers, von mir nachgemessen und **erweitert**
- `pfad`: `.d-check.yml:8` (`modules: [links, anchors, …]`) ·
  `docs/user/benutzerhandbuch.md:186`, `:219`, `:527`, `:531`, `:532`
  (die fünf Vorkommen)
- `befund`: Sechs eigene Mutationen am `anchors`/`links`-Modul (gepinntes
  Image `sha256:18e9cd85…`, Tabelle unten): derselbe kaputte Anker ist
  **rot**, wenn der Link einzeilig steht (`anchor-missing`, Exit 1), und
  **grün**, wenn der Linktext über einen Zeilenumbruch geht (Exit 0); dasselbe
  gilt für einen Umbruch zwischen `](` und dem Ziel und — über den
  Implementer-Befund hinaus — für ein **fehlendes Linkziel**
  (`target-missing`: einzeilig Exit 1, zweizeilig Exit 0). Die Lücke betrifft
  damit nicht nur Anker, sondern beide Referenzprüfungen, und beide
  Umbruchstellen. Betroffen sind repo-weit 9 mehrzeilige Links in 5 Dateien,
  davon fünf im Handbuch (`git blame`: alle aus 2026-09-13, **keine** aus
  diesem Slice); alle fünf sind inhaltlich in Ordnung — die vier Zielanker
  existieren (`#zugriff-und-rollen`, `#tabelle-aktivieren`,
  `#metriken-lesen`, `#blockierende-consumer-erkennen`), geprüft von Hand,
  weil kein Gate sie prüft. Kein Finding an einem der Links; die Lücke selbst
  gehört als Punkt in `harness/sensors/docs-check.md` §Grenze und von dort in
  die Closure — der Wächter dieser Form ist das Review, kein Gate.
- `verifizierbar`: ja — die sechs Mutationsläufe unten (Exit-Codes je Aufruf);
  ein Dauer-Sensor wäre die Nennung der Form in §Grenze, kein neues Gate
- `klasse`: „Gate-Reichweite enger als der Vertrag — mehrzeilige Links ungeprüft"

### N-3 — Die gemeinsame Feldtabelle typt die Row Images als `bytes`; im SSE-Event sind sie eingebettete JSON-Werte

- `kategorie`: INFO
- `quelle`: `SPEC-021` (Zeile *Event-Daten*: „Die Row Images stehen als
  eingebettete JSON-Werte; ein fehlendes Bild ist `null`") ·
  `SPEC-020` (Zeile *Nachricht*, dort `bytes` für den gRPC-Träger) ·
  `internal/adapters/driving/http/sse.go:32-53`
- `pfad`: `docs/user/benutzerhandbuch.md:643-654` (Tabelle) ·
  `:682-686` (SSE-Verweis)
- `befund`: Der SSE-Abschnitt verweist für die Feldliste jetzt auf den
  gRPC-Abschnitt — die zehn **Namen** sind in beiden Trägern identisch
  (`sse.go:33-42` JSON-Tags vs. `changestream.proto:18-27`), die **Wertform**
  der beiden Row Images nicht: dort `bytes`, im Event ein eingebettetes
  JSON-Objekt bzw. `null`. Die Zielstelle führt nur „`bytes`"; die
  SSE-Sektion trägt die Null-Hälfte („ein fehlendes Row Image ist `null`"),
  nicht die Werteform. Kein Falschsatz — „dieselben zehn Felder" gilt für die
  Namen —, aber ein Leser der SSE-Sektion, der die Typ-Spalte mitliest, kann
  die Werteform falsch ableiten. Ohne Rückgabe-Pfeil; geht in die Closure.
- `verifizierbar`: nein — Doku-Konsistenz
- `klasse`: „Übernommene Feldtabelle trägt die Werteform des zweiten Trägers nicht"

## Negativbefunde

- geprüft, ohne Befund: F-1 (a) **Wortlaut trifft das Verhalten** — „öffnet den
  Listener in eigener Goroutine (`Run` unten)" deckt `wiring.go:738-753`,
  „über `log.Error` gemeldet" deckt `:749-751`, „geht nicht in das Ergebnis von
  `Run` ein" deckt `:825` (einziger Nicht-Fehler-`return`); der Kommentar trägt
  eine Zusage-Klasse im Sinne `AGENTS.md` §3.7, keine Chronik, keine verworfene
  Alternative
- geprüft, ohne Befund: F-1 (b) **die Datei ist frei von der Behauptung** —
  `grep` über `internal/bootstrap/wiring.go` findet „Ergebnis von `Run`"
  nur noch an `:114-115` (neu, richtig) und `:581-582` (Nachbarblock, richtig);
  die frühere Fassung ist restlos ersetzt
- geprüft, ohne Befund: F-1 (c) **genau eine Fundstelle** — repo-weiter
  Grep über die Aussageklasse (`Startfehler` × `Run`, „Ergebnis von `Run`")
  findet außerhalb der beiden richtigen Kommentare nur zwei **historische**
  Nennungen der alten Fassung: `docs/plan/planning/done/slice-070-grpc-capture-integration.md:186`
  (der dortige Plan-Record der Korrektur) und das Zitat im Erstlauf-Report.
  Kein Produktionspfad, kein Sensor, kein ADR trägt die Behauptung
- geprüft, ohne Befund: F-2 **die Belege tragen** — `cmd/pg-change-feed/main.go`
  kennt fünf Modi (`--version`, `--healthcheck`, `register-consumer`,
  `acknowledge-consumer`, `diagnose`), keinen Retention-Modus;
  `tools/schema/nacharbeit-administration.sql` definiert vier `cdc.*`-Funktionen
  (`enable_table`, `disable_table`, `exclude_column`, `include_column`), keine
  Retention-Auslösung; für die acht übrigen Zeilen der Endpunkt-Tabelle ist ein
  CLI- oder SQL-Weg vorhanden (`nacharbeit-roles.sql:66`: `SELECT, INSERT,
  UPDATE, DELETE` für `cdc_admin` auf `cdc.consumer`/`cdc.consumer_position`;
  `:105`: `SELECT` für `cdc_reader` auf den Lese-Views) — „Eine Ausnahme" ist
  damit genau eine, und die genannte; der Verweis
  `#aufbewahrung-retention` löst auf (`:388`)
- geprüft, ohne Befund: F-3 **Feldtabelle Zeile für Zeile** — zehn Zeilen,
  `change_id`/`transaction_id`/`source_table_id`/`schema`/`table`
  (string), `sequence` (int64), `operation` (string, drei Werte),
  `old_image`/`new_image` (bytes), `schema_version` (string) decken
  `proto/cdc/stream/v1/changestream.proto:17-28` **und** die
  `SPEC-020`-Zeile *Nachricht* **und** `grpc/server.go:142-155`
  (`toProtoChange`, alle zehn Felder gesetzt) sowie `http/sse.go:33-42`
  (dieselben zehn JSON-Schlüssel); keine Spalte fehlt, keine ist erfunden,
  keine ist umbenannt; „bei `INSERT` leer" / „bei `DELETE` leer" entspricht
  dem Null-Fall aus `handbuch:385-386`
- geprüft, ohne Befund: Plan §2, neuer Punkt 2 („Der Deklarationsort trägt die
  beschriebene Semantik …) — beide Hälften sind an `1cf5675` real erfüllt
  (Belege oben), das Häkchen ist berechtigt; die Klasse ist mit
  `review-slice-070` F-2 benannt und der Vorgang damit zählbar
- geprüft, ohne Befund: Plan §3, zweite Zeile (`internal/bootstrap/wiring.go`,
  update, „kein Verhaltens-Change") — deckt den Diff exakt: `1cf5675` ändert
  in der Datei fünf Kommentarzeilen, keine Anweisung (kein `git diff`-Hunk mit
  Code-Zeile)
- geprüft, ohne Befund: Handbuch-Versionshistorie der Fixrunde — `Version:
  1.13` → `1.14` und eine neue Zeile 1.14, die genau die zwei inhaltlichen
  Änderungen nennt; die Hochregel ist nicht ausgelöst
- geprüft, ohne Befund: der Fix hat **keine** neue mehrzeilige Verweisform
  eingeführt (die zwei neuen Links `#aufbewahrung-retention` und
  `#zugriff-über-den-grpc-change-stream` stehen einzeilig und sind damit
  geprüft — docs-check 0 Befunde)
- geprüft, ohne Befund: `3637e1f` selbst — Inhalt (drei Absätze: §1-Ausschluss,
  §2-Liefer-Punkt, dessen Zählung) ist mit dem Fix konsistent; **kein** Finding
  an diesem Commit, nur der benannte Traceability-Fehler des Betreibers

## Eigene Messungen

Alle Läufe read-only, ungepiped, Exit-Code je Aufruf einzeln (`AGENTS.md` §3.9).
Die Mutationsläufe liefen in einer Kopie unter `/tmp` (außerhalb des
Arbeitsbaums, `tar`-Auszug ohne `.git`); das gepinnte Modul-Image
`sha256:18e9cd85…` aus `d-check.mk`.

| Messung | Ergebnis | Exit |
|---|---|---|
| `make gates` (Arbeitsbaum) | baseline-verify 54 Dateien OK · coverage 49,30 % ≥ 40 % · docs-check 635/0 — Abbruch bei `commit-traceability` | **2** |
| Befundliste des roten Laufs | **genau ein** Befund: `3637e1f  commit-untraceable` | — |
| `make a-check` (einzeln, wegen des Abbruchs) | 0 Befunde | **0** |
| Mutation A — einzeiliger Link, Anker mutiert | `handbuch:590  anchor-missing` | 1 |
| Mutation B — derselbe Anker, Linktext zweizeilig | 0 Befunde | **0** |
| Mutation C1 — vier bestehende mehrzeilige Links, Anker mutiert | Befund nur an den **einzeiligen** Vorkommen (`:185`, `:252`, `:318`) | 1 |
| Mutation C1′ — nur Anker, die ausschließlich in mehrzeiligen Links stehen (`#metriken-lesen`, `#blockierende-consumer-erkennen`) | 0 Befunde | **0** |
| Mutation C2 — Umbruch zwischen `](` und Ziel, Anker mutiert | 0 Befunde | **0** |
| Mutation D1 — einzeiliger Link, Ziel existiert nicht | `handbuch:590  target-missing` | 1 |
| Mutation D2 — derselbe Link, Text zweizeilig | 0 Befunde | **0** |
| `git blame` der fünf mehrzeiligen Links | alle aus 2026-09-13 (`9cae8f9`, `b47bbfb`, `8a0ccee`) — keine aus `slice-077` | 0 |
| `git status --porcelain` des Arbeitsbaums | leer (keine Probe im Baum) | — |

Die Paarung A/B und C1′/C1 schließt aus, dass der Unterschied an der Auswahl
der Stelle liegt: **derselbe** Anker, **dieselbe** Datei, nur die Zeilenform des
Links unterscheidet rot von grün. D1/D2 zeigen, dass die Lücke auch die
`target-missing`-Prüfung trifft.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Plan nennt ein Verhalten ‚geltende
Entscheidung', ohne dass ein Artefakt sie trägt" · „Gate-Reichweite enger als
der Vertrag — mehrzeilige Links ungeprüft" (verwandt: `BEO-PGC/regel-weiter-als-ihr-sensor`)
· „Übernommene Feldtabelle trägt die Werteform des zweiten Trägers nicht"

## Verdikt

**Merge-blockierend:** **nein** — 0 HIGH; F-1 ist mit eigenem Beleg
geschlossen, F-2/F-3 tragen. N-1 ist ein Wortlaut-Punkt im Zeitdokument
(**Planner**, kein Implementer-Pfeil), N-2 eine benannte Gate-Reichweiten-Lücke
(Träger: §Grenze in `harness/sensors/docs-check.md` und die Closure), N-3 ohne
Rückkante.

**Übergabe:** keine Rückkante Reviewer → Implementer; kein Konflikt-Pfad (kein
Befund bestreitet eine Regel, kein HIGH mit Widerspruch). N-1/N-2/N-3 gehen in
die Slice-Closure §7 und von dort als Klassen in den Zähler; N-2 ist die
einzige davon mit einem benannten Träger außerhalb des Slice.

**Zur Zählung des neuen §2-Punkts:** die Plan-Regel nennt „≤ 3 Liefer-Punkte";
der Plan stand mit vier Inhalts-Punkten bereits vor der Fixrunde so, und der
neue Punkt ist eine **review-erzwungene** Korrektur, keine selbst gewählte
Umfangs-Erweiterung — ein Neu-Schnitt ist damit nicht angezeigt. Benannt, damit
die Zahl nicht unbemerkt wächst.

**DoD-Nachzug:** **ja** — die Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" in §2 des Slice-Plans ist auf `[x]` nachgezogen, mit
Verweis auf Erstlauf und diesen Delta-Lauf (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde). **Nur diese eine Zeile**;
§6-Risiko-Ausgänge, Beobachtungs-Register, Closure-Notiz und die drei Paarungen
bleiben Planner-Arbeit.

**Was ausdrücklich trägt** (nicht Gegenstand einer Fixrunde): der geänderte
Kommentar in `wiring.go` samt seiner drei Hälften, die Rahmen-Aussage der
HTTP-§4 samt ihren beiden Belegen, die Zehn-Feld-Tabelle gegen Proto, `SPEC-020`
und beide Adapter, die Versionshistorie der Fixrunde und der Plan-Nachzug in
§3. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11).
