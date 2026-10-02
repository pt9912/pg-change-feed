# ADR-0144: Meldungscodes — nutzerseitige Kennungen für Fehler und Warnungen statt interner Kennungen in Programm-Ausgaben (Schärft ADR-0023, ADR-0049)

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-02

**Autor:** pt9912 (Architect-Rolle, Modul 8; anderer Kontext als der Planner-Lauf
des Slice `meldungscodes-statt-interner-kennungen`, dessen sieben Fragen diese ADR
beantwortet)

**Bezug:** [`LH-QA-OPS-001`](../../../spec/lastenheft.md) (Betriebsfähigkeit,
Betreiber-Sicht), [`LH-QA-REL-003`](../../../spec/lastenheft.md) (Fehlerklassen),
[`LH-FA-ADM-003`](../../../spec/lastenheft.md) (Fehlerzustand erkennbar),
[ADR-0143](0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) (das
Handbuch trägt keine interne Kennung; Abgrenzung in Festlegung 7),
[ADR-0134](0134-sdk-public-doc-check-gate-make-gates.md) (Muster des Wächters),
[ADR-0083](0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen),
`AGENTS.md` §3.5 (Accepted-ADRs sind unberührbar), §3.6 (ein neues Gate braucht
eine ADR), §3.11, §3.12; Slice-Plan `docs/plan/planning/open/slice-meldungscodes-statt-interner-kennungen.md`.

**Schärft:** [ADR-0023](0023-fehlerklassifikation.md) und
[ADR-0049](0049-replication-fehlerklassen-schwellen.md): die **sieben Klassen und
die Sentinel-Trennung bleiben unverändert**; diese ADR legt eine feinere
Kennung **unter** der Klasse fest und ändert keine Klassen-Semantik. Spec-Stelle:
[`SPEC-008`](../../../spec/pflichtenheft.md) (Fehler-Codes und Logging-Felder, §4) —
der Nachzug ist in Festlegung 8 benannt. Die zwei Accepted-ADRs werden nicht
geändert (`AGENTS.md` §3.5).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Anlass.** Der Auftraggeber hat festgelegt: Interne Kennungen (`LH-*`, `ADR-*`,
`SPEC-*`, `ARC-*`) gehören nicht in die Programm-Ausgabe; für Fehler und Warnungen
braucht der Betreiber eigene Nummern, Vorbild ist das Schwester-Repo d-migrate.

**(2) Ist-Stand der Ausgaben** (gemessen am Stand `abf71ace`, 2026-10-02, mit
`git grep`, Befehle in der Fitness Function): 14 Zeilen Go-Produktionscode
tragen eine Kennung in einem Ausgabe-Literal — zwölf in der Diagnose-Ausgabe
(`internal/bootstrap/wiring.go`), eine im Fehlertext der Konfigurationsdatei
(`internal/bootstrap/config_file.go`), eine in der Begründung der Rollout-Wache
(`tools/schema/rolloutguard/guard.go`); dazu 3 `echo`-Zeilen in `tools/schema/rollout.sh`. Alle
14 Go-Zeilen sind gelesen und echte Literale. Die Fehlertexte der Quelle tragen
bereits ein einheitliches Präfix `Fehlerklasse <klasse>: …` (28 Literale in
Produktions-Go, gemessen); die Klasse ist heute die einzige stabile
Maschinen-Kennung einer Fehlermeldung. SQL-Funktionen und Views erzeugen selbst keinen
Fehlertext: `git grep -n 'RAISE EXCEPTION' -- tools/schema internal/adapters/driven/postgresstorage`
liefert 0 Zeilen (gemessen) — Fehler dort sind native PostgreSQL-Fehler mit SQLSTATE.

**(3) Das Vorbild** (Schwester-Repo d-migrate, `spec/ledger.md`, `spec/cli-spec.md`
§4.1, gelesen): Code = Buchstabe (`E`/`W`) + drei Ziffern, der Buchstabe trägt die
Schwere, die Nummer einen Bereich; der Code ist ein eigenes Feld des Ergebnisobjekts,
die CLI formatiert `Error [E053]: …`; ein Ledger (YAML) hält je Code Status, Test
und Evidenz, ein Test prüft „kein Code ohne Test und Evidence“; ein Code wird nie neu
belegt (`E137` ist zurückgezogen, die Kennung bleibt vergeben). Zwei Grenzen der
Übertragbarkeit: (a) d-migrate ist ein Batch-Werkzeug, dieses Repo ein Dauerprozess
mit Log, Diagnose, HTTP, gRPC und SQL; (b) d-migrate läuft **im Rollout** dieses Repos
(`make schema-rollout`) und druckt seine eigenen `E`/`W`-Codes in dasselbe Terminal —
bloße `E053`-Codes dieses Projekts wären dort mit denen des Werkzeugs verwechselbar
(*hergeleitet*: die Ausgabe von d-migrate im Rollout ist nicht nachgemessen).

**(4) Randbedingung Handbuch.** [ADR-0143](0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
hält das Handbuch frei von internen Kennungen (Muster `P`). Ein Meldungscode ist eine
**nutzerseitige** Kennung und gehört in den Handbuch-Katalog; er darf nicht wie eine
interne Kennung aussehen, sonst färbt er das Gate rot.

## Entscheidung

Wir wählen **einen Meldungscode je Fehler- und Warnursache, feiner als die
Fehlerklasse, mit Projekt-Präfix und mit der Klasse in der ersten Ziffer**. Der
Code ist ein eigenes Feld, wo der Weg ein strukturiertes Feld hat, und ein fester
Kopf im Text, wo er keines hat. Quelle der Wahrheit ist eine Tabelle im Code; das
Handbuch trägt den Katalog; ein netzloses Gate hält beide gleich. Interne Kennungen
verschwinden aus allen Ausgabe-Literalen, und ein zweites netzloses Gate hält diesen
Zustand. Die Umsetzung läuft in vier Teil-Slices (Festlegung 9).

### Teilfrage 1 — Schema der Codes

| Option | Pro | Contra |
|---|---|---|
| A — d-migrate übernommen: `E`/`W` + 3 Ziffern, ohne Präfix | größte Nähe zum Vorbild | im Rollout Verwechslung mit den Codes von d-migrate (Kontext 3b); 3 Ziffern ließen nur Bereiche zu, die d-migrate bereits sprengen musste (das Ledger führt unregelmäßige Bereiche) |
| **B — Projekt-Präfix + Schwere + 4 Ziffern: `PCF-E4001` (gewählt)** | eindeutig, wo mehrere Werkzeuge in einem Log stehen; 999 Einzelursachen je Bereich ohne Umbruch der Bereiche | länger als `E053` (9 statt 4 Zeichen); ein Präfix, das zu pflegen ist |
| C — rein numerisch / UUID | kürzeste Form / kollisionsfrei | nicht lesbar, keine Schwere, keine Klasse erkennbar |

### Teilfrage 2 — Beziehung zu den sieben Fehlerklassen

| Option | Pro | Contra |
|---|---|---|
| A — zwei unabhängige Achsen (Klasse und Code nebeneinander) | einfach zu erklären | Klasse und Code können sich widersprechen; ein Test müsste die Paarung je Fund prüfen |
| **B — Code in der Klasse: erste Ziffer = Klasse, Klasse bleibt Maschinen-Achse (gewählt)** | kein Widerspruch möglich (Klasse folgt aus dem Code); `error_class`, das Metrik-Label und das Wort `Fehlerklasse` bleiben unverändert; ein Test prüft Ziffer gegen Klasse | eine Änderung der Klasse eines Fehlers ist ein neuer Code (gewollt, Festlegung 4) |
| C — Klasse nur noch aus dem Code abgeleitet, `error_class` entfällt | eine Achse | bricht `cdc.process_heartbeat.error_class`, die Metrik und jeden Betreiber, der sie liest — verstößt gegen die Zusage „Klassen unverändert“ |

### Teilfrage 3 — Abbildung je Weg

| Option | Pro | Contra |
|---|---|---|
| A — Code immer als Präfix im Text | ein Mechanismus | wo ein strukturiertes Feld existiert (JSON, gRPC-Detail, Log-Attribut, Spalte), muss der Leser den Text parsen |
| B — Code immer als eigenes Feld | maschinenfreundlich | der Prozess-Ende-Text, die Textspalte `error_message` und die Skript-Zeile haben kein Feld |
| **C — Feld, wo der Weg eines hat; sonst fester Kopf im Text (gewählt)** | folgt dem Weg; der Text bleibt für Menschen lesbar | zwei Formen (in Festlegung 3 tabelliert) |

### Teilfrage 4 — Stabilität

| Option | Pro | Contra |
|---|---|---|
| A — Text und Code stabil | maximale Betreiber-Sicherheit | jede Textkorrektur wird ein Vertragsbruch; Texte tragen Laufzeitdetails (Regelname, Spalte), die nie stabil sein können |
| **B — Code stabil, Text frei, Klasse und Ausgang 1 stabil (gewählt)** | der Betreiber hat einen Anker, der Text bleibt verbesserbar | ein `grep` auf Text bricht; Übergangsregel in den Konsequenzen |
| C — nichts zusagen | keine Pflege | der Code hätte keinen Wert; entspricht dem Zustand ohne Nummern |

### Teilfrage 5 — Dokumentation und Prüfung

| Option | Pro | Contra |
|---|---|---|
| A — Handbuch-Katalog als einzige Quelle | eine Stelle | der Code im Quelltext kann ohne Katalog-Eintrag entstehen, nichts färbt rot |
| B — zusätzlich ein YAML-Ledger wie d-migrate | maschinenlesbar mit Test- und Evidenz-Feldern | eine dritte Stelle neben Code und Handbuch; d-migrate braucht sie, weil seine Codes über Kotlin-Module und Dokumente verteilt sind, hier liegen Code und Test im selben Modul |
| **C — Tabelle im Go-Code als Quelle, Handbuch-Katalog, netzloses Gate gleicht beide ab (gewählt)** | zwei Stellen statt drei; der Code im Quelltext verweist auf die Konstante der Tabelle, ein Tippfehler ist ein Übersetzungsfehler | der Abgleich Handbuch gegen Tabelle ist ein Textvergleich, keine Semantik-Prüfung |

### Teilfrage 6 — Spec-Stelle

| Option | Pro | Contra |
|---|---|---|
| A — nur diese ADR | keine Spec-Änderung | `SPEC-008` nennt den Fehlertext-Beginn (`schema: `), die Prosa wäre falsch |
| **B — `SPEC-008` um einen Absatz erweitern (gewählt)** | der Abschnitt heißt bereits „Fehler-Codes und Logging-Felder“; keine neue Kennung | der Pflichtenheft-Nachzug ist Arbeit des Planners (Festlegung 8) |
| C — neues `SPEC-*` und eine Lastenheft-Anforderung | stärkste Verankerung | das Lastenheft ist vertraglich; eine neue Anforderung ist Auftraggeber-Entscheidung und für den Betrieb nicht nötig |

### Teilfrage 7 — Reichweite

| Option | Pro | Contra |
|---|---|---|
| **A — Fehler und Warnungen; Info reserviert (gewählt)** | folgt der Vorgabe des Auftraggebers | `I` bleibt ungenutzt |
| B — alle Ausgabezeilen mit Code | einheitlich | die Diagnose-Zeilen sind Zustandsberichte; ein Code je Zeile ist Rauschen |
| C — nur Fehler | kleiner | Warnungen (WAL-Schwelle, Wiederholung, Backfill-Größe) bleiben ohne Anker |

### Festlegungen

1. **Form und Nummernraum (Teilfrage 1).** Ein Meldungscode hat die Form
   `PCF-<S><NNNN>`, ERE `PCF-[EWI][0-9]{4}`; `S` ist die Schwere (`E` Fehler, `W` Warnung,
   `I` Information). `I` ist reserviert und derzeit nicht vergeben (Teilfrage 7). Der
   Nummernraum:

   | Schwere | erste Ziffer | Bedeutung |
   |---|---|---|
   | `E` | 1 `transient`, 2 `configuration`, 3 `permission`, 4 `schema`, 5 `storage`, 6 `replication`, 7 `internal` | die Fehlerklasse (Reihenfolge der `SPEC-008`-Tabelle) |
   | `E` | 8 | **Ablehnung einer Aufrufer-Eingabe** (abgelehnter Antrag, HTTP `400`/`404`, gRPC `InvalidArgument`/`NotFound`): trägt keine Fehlerklasse |
   | `E` | 9 | reserviert |
   | `W` | 1 Erfassung und Replikation, 2 Backfill, 3 Retention und Speicher, 4 Verwaltung (Anträge, Regeln), 5 Konfiguration und Start, 9 reserviert | der Bereich (Warnungen tragen keine Klasse) |

   Die drei übrigen Ziffern sind fortlaufend je erster Ziffer, am Ende angefügt; `000` ist
   der **Rückfall der Klasse** (`PCF-E1000` … `PCF-E7000`): jeder klassifizierte Fehler, dem
   keine Einzelursache zugeordnet ist, trägt den Rückfall seiner Klasse und nie keinen
   Code. Die Beispielformen in dieser ADR (`PCF-E4001` u. ä.) sind **nicht vergeben**; die
   Vergabe steht allein in der Tabelle (Festlegung 5).
2. **Granularität und Klasse (Teilfrage 2).** Eine Einzelursache ist eine Ursache, bei der
   sich die **Maßnahme des Betreibers** von der anderer Ursachen derselben Klasse
   unterscheidet; Stellen mit gleicher Maßnahme teilen einen Code. Die Klasse
   eines Fehlers ist die des Codes; `cdc.process_heartbeat.error_class`, das Label der
   Metrik `cdc_errors_total{class}` und das Wort `Fehlerklasse` im Text bleiben unverändert.
   **Kein Code als Metrik-Label** (der Metriken-Vertrag bleibt; ein Label je Code
   vervielfachte die Kardinalität ohne Nutzen, die Klasse genügt für Alarmierung).
3. **Abbildung je Weg (Teilfrage 3).**

   | Weg | Form |
   |---|---|
   | Fehlerwert, Prozess-Ende-Zeile, `error`-Attribut einer Log-Zeile | Kopf `Fehlerklasse <klasse> [<code>]: <Ursache>` (heute `Fehlerklasse <klasse>: <Ursache>`; das `[<code>]` wird eingefügt, Klasse und Doppelpunkt bleiben) |
   | Log-Zeile mit Warnung | eigenes Attribut `code=<code>`, der Meldungstext trägt keinen Code |
   | Prozess-Ausgang | Ausgang 1 unverändert; **kein** numerischer Ausgang je Code |
   | `error_message` einer Backfill-Run-Zeile | `<klasse> [<code>]: <Ursache>` (heute `<klasse>: <Ursache>`) |
   | `error_message` einer abgelehnten Antrags-Zeile | `abgelehnt [<code>]: <Klartext>`, Code aus dem Bereich `E8` |
   | `cdc.process_heartbeat` und `cdc.heartbeat` | additive, nullable Spalte `error_code` neben unveränderter `error_class` (`NULL` bei Normalbetrieb) |
   | `diagnose` (CLI, `GET /diagnose`, RPC `Diagnose`) | die Zeile „Fehlerzustand“ nennt Klasse und Code (`schema [<code>]`), die Antwort trägt ein additives Feld `error_code`; **alle anderen Diagnose-Zeilen tragen weder Kennung noch Code** (Zustandsberichte, Teilfrage 7) |
   | HTTP-Fehlerkörper | additives Feld `code` neben `error`: `{"error": "…", "code": "<code>"}` |
   | gRPC | Status-Detail `google.rpc.ErrorInfo` mit `reason = <code>` und `domain = pg-change-feed` |
   | SQL-Funktionen und Views | **kein Code**: sie erzeugen keinen eigenen Fehlertext (Kontext 2); native PostgreSQL-Fehler behalten SQLSTATE |
   | Rollout-Skript | `FEHLER [<code>]: <Text>` für Fehlerzeilen; die Fortschrittszeilen (`schema-rollout: …`) tragen keinen Code |

   Die SDKs unter `sdks/` ändern sich in dieser ADR nicht; sie reichen Servertexte
   durch (additive Felder sind für tolerante Leser unschädlich — *hergeleitet*, die
   Toleranz der drei Clients gegen ein unbekanntes JSON-Feld ist nicht nachgemessen).
4. **Stabilität (Teilfrage 4).** Ein Code wird **nie neu belegt**. Ein Code, dessen
   Ursache entfällt, wird **zurückgezogen** (Status `withdrawn`), nicht gelöscht; er
   bleibt in Tabelle und Katalog, mit dem Vermerk. Die Bedeutung eines Codes wird nicht
   erweitert; teilt sich eine Ursache, entstehen neue Codes, der alte bleibt. Die Klasse
   eines Codes ändert sich nie (die erste Ziffer ist Teil der Kennung): ändert sich die
   Klasse eines Fehlers, entsteht ein neuer Code. **Stabil sind** Code, Klasse, das Wort
   `Fehlerklasse` am Kopf-Anfang und Ausgang 1. **Nicht stabil ist der Text** nach
   dem Kopf; wer maschinell auswertet, wertet den Code oder die Klasse aus. Der Katalog
   nennt das in einem Satz.
5. **Quelle, Katalog und Prüfung (Teilfrage 5).** Quelle der Wahrheit ist eine Tabelle
   im Paket `internal/domain` (Code, Schwere, Klasse, Status); jede Stelle im Go-Code
   verwendet die Konstante der Tabelle, kein zweites Literal. Der Katalog steht im
   Benutzerhandbuch, kennungsfrei: je Code **Code, Klasse bzw. Bereich, Bedeutung,
   Maßnahme**, zurückgezogene Codes mit Vermerk. Zwei Prüfungen tragen das:
   (a) ein Go-Test in `make test`: Format gegen die ERE aus Festlegung 1, erste Ziffer
   gegen die Klasse, keine Dopplung, der Rückfall `…000` je Klasse vorhanden;
   (b) ein netzloses Gate `meldungscodes-check` in `GATE_CHECKS` (nur `bash`/`git`/`grep`):
   jeder Code in der ERE-Form in Produktions-Go, den Skripten unter `tools/schema/` und
   `examples/` steht in der Tabelle, und die Menge der Codes der Tabelle ist gleich der
   Menge der Codes im Katalog. Heimat `harness/mk/doc-gate.mk`, Sensor-Vertrag
   `harness/sensors/meldungscodes-check.md`, Tabellentest als Werkzeug ohne Gate-Aufnahme
   (Muster [ADR-0134](0134-sdk-public-doc-check-gate-make-gates.md) Teilfrage 2).
   Eine Milderung (weniger Prüfung, Ausnahme) ist eine ADR (`AGENTS.md` §3.6).
6. **Spec-Stelle (Teilfrage 6).** Die Festlegung gehört in `SPEC-008`; kein neues `SPEC-*`
   und keine Lastenheft-Änderung (Festlegung 8 nennt den Nachzug).
7. **Reichweite und Abgrenzung zum Handbuch-Gate (Teilfrage 7).** Codes tragen
   klassifizierte Fehler, Ablehnungen von Aufrufer-Eingaben und Warnungen. Diagnose-Zeilen,
   Fortschrittszeilen und Info-Log-Zeilen tragen keinen Code. **Der Code sieht nicht wie
   eine interne Kennung aus:** das Präfix `PCF-` kommt in keinem Muster der Gates vor.
   Gemessen (`grep -cE` mit dem Muster `P` aus [ADR-0143](0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
   Festlegung 3 an `printf`-Zeilen, 2026-10-02): `PCF-E4003`, `PCF-W2001`, `PCF-I1001` und
   `E053` je 0 Treffer; `LH-FA-X`, `ADR-0049`, `CO-123` je 1 Treffer. Der Handbuch-Katalog
   darf Codes nennen und bleibt grün; `PCF-` + Buchstabe + 4 Ziffern ist
   die einzige Form eines Codes — ein Code mit anderem Präfix wäre ein Fehler in Tabelle und Katalog.
   Das Gate [ADR-0143](0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
   bleibt unverändert. Für `docs/user/` ist diese Aussage aus den Mustern hergeleitet;
   der Lauf von `make handbuch-public-doc-check` mit einem Katalog im Handbuch gehört
   zum Teil-Slice 2.
8. **Spec-Nachzug (Aufgabe des Planners, nicht dieser ADR).** `spec/pflichtenheft.md`
   §4 `SPEC-008`: ein Absatz „Meldungscode“ (Form, Nummernraum, Stabilität, Kopf);
   der Satz „der Fehlertext beginnt mit der Klasse (`schema: `)“ in der Prosa nach der
   Tabelle wird zum Kopf nach Festlegung 3; `SPEC-029` (`error_message` des Runs),
   `SPEC-019` (`error_message` des Antrags), `SPEC-018` (HTTP-Fehlerkörper), `SPEC-031`
   (gRPC-Fehler) und die Tabelle von `cdc.process_heartbeat` (`error_code`) werden mit
   dem Teil-Slice nachgezogen, der den jeweiligen Weg ändert (Festlegung 9). Das
   Lastenheft braucht keine Änderung; eine eigene Anforderung wäre Auftraggeber-Entscheidung.
9. **Umsetzung in vier Teil-Slices, in dieser Reihenfolge.**
   - **T1 `meldungscodes-kennungsfreie-ausgaben`** (braucht keinen Code): die 14 Go-Zeilen
     und die 3 `echo`-Zeilen ohne Kennung, die Erwartungen in den Läufern nachgezogen,
     Hinweissatz und Platzhalter im Handbuch entfernt, die `description:`-Felder
     von `tools/schema/schema.yaml`: ist die Beschreibung als Datenbank-Kommentar beim
     Betreiber sichtbar (zu messen am Rollout-Ergebnis), wird sie kennungsfrei, sonst ist
     sie Quelltext-Dokumentation und bleibt; Gate `ausgabe-kennungen-check` (Festlegung 10).
   - **T2 `meldungscodes-registry-fehlerkopf`:** die Tabelle mit den sieben Rückfall-Codes
     und den Einzelursachen der Sentinels, der Kopf in Fehlerwerten, Prozess-Ende-Zeile,
     `error`-Attribut, `error_message` von Run und Antrag, `FEHLER [<code>]` im Rollout,
     der Handbuch-Katalog, das Gate `meldungscodes-check`, der Go-Test. Die
     Pflichtenheft-Stellen dieser Wege zieht der Planner im selben Slice nach.
   - **T3 `meldungscodes-warnungen-heartbeat-diagnose`:** `W`-Codes als Log-Attribut,
     Spalte `error_code` mit View und Rollout-Vorlauf, Diagnose in CLI, HTTP und gRPC.
   - **T4 `meldungscodes-http-grpc-fehlerkoerper`:** `code` im HTTP-Fehlerkörper,
     `ErrorInfo` bei gRPC.

   T1 vor T2, weil T1 den Betreiber-Teil der Vorgabe (keine Kennung in der Ausgabe) ohne
   Code-Entscheidung erfüllt und sein Gate den Ausgangszustand 0 festhält, bevor T2 die
   Literale der Fehlertexte anfasst. T3 und T4 hängen an T2 (die Tabelle) und sind
   untereinander unabhängig. Die Aufteilung erfüllt die Rückführung „zu groß“ des Plans.
   Jeder Teil-Slice ist für sich lauffähig; ein Server-Release ist keiner von ihnen.
10. **Wächter `ausgabe-kennungen-check` ist ein Gate in `make gates`, ab Ist-Stand 0.**
    Gegenstand: Ausgabe-Literale in Produktions-Go unter `internal/`, `cmd/` und `tools/schema/`
    (ohne `*_test.go`) und `echo`-/`printf`-Zeilen in den Skripten unter `tools/schema/` und
    `examples/`. Muster, ERE (kein `-P`, `AGENTS.md` §3.1), zweistufig: eine Zeile trifft,
    wenn eine Kennung `(LH|ADR|SPEC|ARC)-[A-Z0-9]` hinter einem **ungeschlossenen**
    Anführungszeichen steht (also im Literal), und die Zeile nicht mit `//` beginnt.
    Grenzen, im Sensor-Vertrag zu nennen: ein nachgestellter Kommentar, der ein Literal mit
    Kennung zitiert, trifft (laute Richtung); ein **mehrzeiliges** Raw-String-Literal wird
    nicht gelesen (akzeptiertes Negativ: alle 14 Treffer des Bestands sind einzeilig; ein
    AST-Wächter wie `tools/harness/kommentar-kennungen/` kostet ein Docker-Programm und ist
    erst bei einem realen Fund eines mehrzeiligen Literals angezeigt). Ausgenommen sind Tests,
    die Läufer unter `tools/harness/` (Entwickler-Ausgabe, nicht Betreiber-Ausgabe) und
    erzeugter Code. Die Gate-Aufnahme geschieht mit T1, wenn der Ist-Stand 0 ist; vorher
    wäre das Gate rot.

## Verglichene Alternativen

Die Optionen je Teilfrage stehen in den sieben Tabellen oben. Gesamtbild:

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Kennungen bleiben in der Ausgabe | kein Aufwand | widerspricht der Vorgabe des Auftraggebers; Kennungen ohne Bedeutung für den Betreiber |
| B — nur Kennungen entfernen, keine Codes (T1 allein) | kleinster Eingriff | der Betreiber hat weiter nur die Klasse als stabile Kennung; die Vorgabe „eigene Nummern“ bleibt offen |
| **C — T1 sofort, Codes in drei Folge-Slices (gewählt)** | der sichtbare Teil der Vorgabe ist sofort erfüllt; Codes lassen sich nach Weg ausliefern | vier Slices statt einem |
| D — alle Wege in einem Slice | ein Durchgang | Handbuch, vier Läufer, Schema-Rollout (`error_code`) und zwei Netz-Wege in einem Diff; der Plan nennt selbst die Rückführung |

## Konsequenzen

- Positiv: Der Betreiber hat je Fehler und Warnung einen stabilen, im Handbuch
  katalogisierten Anker; die Klasse bleibt, wie sie ist; ein Text darf besser werden.
- Positiv: Zwei netzlose Gates verhindern den Rückfall — Kennung im Ausgabe-Literal, und
  Code ohne Tabelle oder ohne Katalog.
- Negativ (Kompatibilität): Der Kopf `Fehlerklasse <klasse>: …` ändert sich zu
  `Fehlerklasse <klasse> [<code>]: …`, ebenso `error_message`. Wer den Text mit `grep` auf
  `Fehlerklasse schema:` auswertet, bricht; wer auf `Fehlerklasse schema` oder die Klasse
  prüft, nicht. **Übergangsregel: keine** — kein Doppelformat und kein Schalter, weil der
  Text nie als Vertrag zugesagt war (*übernommen* aus dem Plan; im Handbuch nicht nachgemessen:
  der Implementer von T2 liest die Abschnitte Diagnose, Fehlerklassen und Exit-Codes und trägt
  den Satz „Der Text ist nicht Vertrag, der Code ist es“ in den Katalog; findet er eine
  Textzusage, geht der Slice mit dieser Frage zurück an den Auftraggeber). Die
  Änderungshistorie des Handbuchs nennt die Änderung in Betreibersicht und ohne Kennung.
- Negativ: Zwei Formen der Abbildung (Feld und Kopf) und ein Präfix, das zu pflegen ist.
- Negativ (akzeptiertes Negativ): Das Gate `meldungscodes-check` vergleicht Mengen, nicht
  Sinn; ein falsch beschriebener Katalog-Eintrag bleibt grün. Diese Hälfte trägt der Reviewer.
- Folgepflicht: der Spec-Nachzug (Festlegung 8) und der Handbuch-Katalog (T2);
  `harness/README.md` §Sensors trägt je eine neue Zeile für die zwei Gates mit den Teil-Slices.
- Folgepflicht: `AGENTS.md`, `spec/lastenheft.md` und die zwei geschärften ADRs bleiben unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Wächter-Zählung am Bestand (erprobt) | `git grep -n -E '^[^"`]*("[^"]*"[^"`]*)*["`][^"`]*(LH\|ADR\|SPEC\|ARC)-[A-Z0-9]' -- '*.go' ':!*_test.go'` abzüglich Kommentarzeilen liefert **14** Zeilen, alle gelesen und echte Literale (Falsch-Positive: 0); dieselbe Zahl liefert das Muster der Plan-Messung (`git grep -n -P '^[^/]*(LH\|ADR\|SPEC\|ARC)-[A-Z0-9]'`); `echo`-Zeilen mit Kennung in `tools/schema` und `examples`: **3** (gemessen 2026-10-02, Stand `abf71ace`) | `make ausgabe-kennungen-check` (nach T1) |
| Wächter-Muster an Eingabezeilen (erprobt) | acht `printf`-Zeilen, Pipeline `grep -E <Muster> \| grep -vE '^[[:space:]]*//'`: Treffer bei `fmt.Println("x (LH-FA-ADM-002): y")`, `x := f("a/b") + "ADR-0043"`, `errors.New("http://x/y ADR-0043")` und `x := 1 // Fehler "ADR-0049" früher` (letzterer ein **Falsch-Positiv**, laute Richtung); kein Treffer bei `// siehe ADR-0049 Teilfrage 1`, `x := 1 // ADR-0049`, `// Fehlertext "…(LH-FA-ADM-002)…" früher` und der Struct-Tag-Zeile mit Kennung in Backticks (gemessen 2026-10-02) | — |
| Code-Form gegen Muster `P` (erprobt) | `PCF-E4003`, `PCF-W2001`, `PCF-I1001`, `E053` je 0 Treffer, `LH-FA-X`, `ADR-0049`, `CO-123` je 1 Treffer (`grep -cE` mit `P`, gemessen 2026-10-02) | — |
| Wächter-Mutation (hergeleitet; **nicht erprobt**) | Eine in eine Kopie von `internal/bootstrap/wiring.go` im Scratchpad eingefügte Kennung in einem `fmt.Println`-Literal färbt das Skript mit Wurzel = Kopie-Verzeichnis rot; dieselbe Kennung im Kommentar färbt nicht. Stelle, Instanz und Farbe trägt der Implementer in den Bericht; „der Implementer fährt sie“ ist eine Erwartung | `make test-ausgabe-kennungen-check` |
| Registry-Test (hergeleitet; **nicht erprobt**) | Format, erste Ziffer gegen Klasse, keine Dopplung, Rückfall `…000` je Klasse; ein Code mit falscher Klasse färbt `make test` rot — Erwartung, nicht gefahren | `make test` |
| `meldungscodes-check` (hergeleitet; **nicht erprobt**) | ein Code im Quelltext ohne Tabellen-Eintrag und ein Tabellen-Code ohne Katalog-Eintrag färben je Exit 1; zurückgezogene Codes bleiben in beiden Mengen | `make meldungscodes-check`, `make test-meldungscodes-check` |
| Gate-Verdrahtung (hergeleitet) | beide Gates stehen in `GATE_CHECKS` | `make gates` |
| Handbuch mit Katalog (hergeleitet) | `make handbuch-public-doc-check` bleibt mit einem Katalog grün (Festlegung 7) | `make handbuch-public-doc-check` |

## Re-Evaluierungs-Trigger

Ändert sich die Menge der Klassen ([ADR-0023](0023-fehlerklassifikation.md) wird
abgelöst), ist das Schema der ersten Ziffer nachzuziehen (neue ADR mit `Supersedes`).
Wird das Präfix `PCF-` mit einem anderen Werkzeug im Betrieb verwechselt oder verlangt ein
Betreiber einen maschinenlesbaren Katalog (JSON), öffnet das die Frage nach einem
Export aus der Tabelle. Findet sich ein mehrzeiliges Literal mit Kennung im Bestand,
ist der Wächter ein AST-Programm (Festlegung 10).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-02 | Accepted — Architect-Entscheidung als Vorstufe der Umsetzungs-Slices | `docs/plan/planning/open/slice-meldungscodes-statt-interner-kennungen.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0144` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
