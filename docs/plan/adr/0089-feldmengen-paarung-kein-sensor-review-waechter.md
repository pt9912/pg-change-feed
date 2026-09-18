# ADR-0089: Feldmengen-Paarung — kein Sensor; der Wächter ist der Review

**Status:** Accepted — Supersedes [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
in **einer** Klausel: der **dritten Zeile** ihrer §Fitness Function
(„`make docs-check` | **Träger-Paarung:** die Feldmenge in `SPEC-016`, im
Handbuch §5.2 und die Schlüssel-Prüfung im Code nennen dieselben Namen
(Form-Prüfung des Nachzugs) | `make gates` (im Bündel)"). Alles Übrige der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) bleibt
**hiermit bestätigt** und wird nicht wiederholt: §Entscheidung mit ihren fünf
Festlegungen (der Diskriminator, die Feldmenge, die Durchleitung, die
Fehlerform, die Nicht-Änderungen), §Verglichene Alternativen A/B/C, die erste,
zweite und vierte Zeile ihrer §Fitness Function und die drei
Re-Evaluierungs-Trigger.

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Reviewer-Lauf von `slice-096` <!-- d-check:status-provenance -->,
der F-1 als HIGH auswies und seine eigene Messung vorlegte, und als die
Implementer-Läufe, die den Code geschrieben haben; Modul 8 §Konflikt-Pfad,
Verdikt 1: die Fitness-Function-Zeile ist der Defekt, die Entscheidung gilt)

**Bezug:** [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(§Fitness Function, dritte Zeile — superseded; alles Übrige bestätigt) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
(nimmt „die Fitness-Function-Regeln" von der Zitat-Korrektur aus — deshalb
Folge-ADR statt in-place) ·
[`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md)
(dasselbe Werkzeug für dieselbe Klasse: eine einzelne Fitness-Function-Zeile,
die ihren Satz nicht trägt) ·
[`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md) (Sync-Gate-Form für
Erzeugnisse — die verworfene Alternative D) ·
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) (`AGENTS.md` §3.12 als
Entscheidung) · `AGENTS.md` §3.12 Instanz B (die tragende Regel) · §3.13 (die
Regel, die den fremden Träger gemeldet statt still geändert hat) · §3.6 · §4 ·
§5 · `.d-check.yml` (`modules:`) · [`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md)
§Grenze 7 · Review zu `slice-096` <!-- d-check:status-provenance -->
(Anlass, F-1) · `Makefile` (`test:` — der Mount `$(CURDIR):/src:ro`) ·
`spec/pflichtenheft.md` `SPEC-016` · `docs/user/benutzerhandbuch.md` §5.2 ·
`internal/bootstrap/config_file.go`

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie
[`ADR-0063`](0063-lh-fa-sch-003-testform-korrektur.md)/[`ADR-0064`](0064-lh-qa-ops-005-testansatz-korrektur.md)/[`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md);
trifft eine Fitness-Function-Korrektur, ändert keine Lastenheft-/
Pflichtenheft-Zusage)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass — und die eigene Messung.** Das Review zu
`slice-096` <!-- d-check:status-provenance --> weist F-1 als **HIGH** aus: die
dritte Zeile der `ADR-0088`-Fitness-Function nennt `make docs-check` als
Träger einer Paarung, die dieser Lauf nicht prüft. Die Zeile ist hier
**nachgemessen**, nicht übernommen (2026-09-17, HEAD `b261dc3`):

| Messung (2026-09-17, HEAD `b261dc3`) | Ergebnis |
|---|---|
| `make docs-check`, unmutiert | 789 Datei(en), 0 Befunde, Exit 0 |
| dieselben drei Träger in einer Arbeitskopie **außerhalb** des Repos mutiert — `grpc_addr` → `grpc_addr_x` in der Feldtabelle der `SPEC-016`, `nats_url` aus der Klassenliste des Handbuchs §5.2 entfernt, `http_addr`/`grpc_addr` aus dem §5.2-Beispiel gestrichen — `make docs-check` | 789 Datei(en), **0 Befunde**, Exit 0 |
| dieselbe mutierte Kopie mit **allen** optionalen Modulen (`citations`, `sources`, `tracked`, `targets`, `spans`, `codepaths`, `diagrams`, `pins`, `immutable`, `workflows`) | 789 Datei(en), 137 Befunde, Exit 1 — **133 × `codepath-missing`, 4 × `span-unclosed`**; **kein** Befund nennt `SPEC-016`, das Handbuch oder einen Feldnamen |
| `.d-check.yml`, `modules:` | `[links, anchors, ids, matrix, versions, structure, hostpaths]` — **kein** Modul liest Go-Quelltext |

Der genannte Beleg liefert **mit und ohne** Paarung dasselbe Ergebnis.
*Benannte Abweichung der Zahl (§3.12):* das Review nennt für seine zwei
Arbeitskopien 788 Dateien — ein anderer Stand; die Zahl hier ist an `b261dc3`
mit `make docs-check` gemessen, nicht übernommen.

**(2) Die Aussage ist wahr, ihre Stütze nicht.** Die drei Träger nennen heute
tatsächlich dieselben Namen — `SPEC-016` führt neun zulässige Feld-Schlüssel
und sechs env-exklusive, das Handbuch §5.2 dieselben, und
`forbiddenFileCredentialKeys` samt den `yaml`-Tags von `fileConfig` im Code
ebenso. Die ADR behauptet also nichts Falsches über den Gegenstand; sie nennt
für die Behauptung einen **Beleg, der sie nicht trägt**. Genau das ist die
Klasse `AGENTS.md` §3.12 Instanz B: eine als Beleg gelesene Aussage nennt
ihren Beleg-Anker — oder sie ist als *erwartet* formuliert.

**(3) Was `make docs-check` an diesen drei Trägern tatsächlich prüft.** Die
`ids`-Regel hält eine nackte `SPEC-016`-Kennung zur Linkpflicht auf
`spec/pflichtenheft.md`; `links`/`anchors` prüfen die Verweise **innerhalb**
der Dokumente; `structure` prüft Abschnitts-Invarianten und
Spalten-Mindestbreiten. Keine dieser Regeln vergleicht **zwei** Dokumente
miteinander, und keine liest Go-Quelltext — die `structure`-Regeln führen je
eine `files:`-Angabe auf **eine** Datei und adressieren darin Abschnitt,
Tabelle oder Muster.

**(4) Warum die Paarung an kein Modul fällt.** Die verfügbaren Module sind
vollständig aufgezählt (Kontext (1), vierte Zeile); die nicht im Bündel
laufenden sind darüber hinaus einzeln gefahren, ohne einen Befund zu den drei
Trägern (Kontext (1), dritte Zeile). `codepaths` prüft Pfade, keine
Feldmengen, und ist zudem aus; `citations` allein gefahren meldet nichts. Die
tragende Grenze steht **benannt** in
[`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md)
§Grenze 7: „**Kein Modul prüft einen Symbol- oder Funktionsnamen.**"

**(5) Was maschinell getragen ist — und was nicht.** Die **Code-Hälfte** der
Paarung tragen unverändert die erste und zweite Zeile der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-Fitness-Function:
Durchleitung und Feld-für-Feld-Vorrang je Feld, die Zugangsdaten-Klasse je
Schlüssel — beide über `go test ./internal/bootstrap/...` (`make test`, kein
Gate; das Review hat die Bindungen an der Eingabeseite mit eigenen Mutationen
rot gesehen). **Nicht** getragen ist die Paarung **zwischen** den drei
Trägern; sie ist die einzige der vier Zeilen, die ein Gate behauptet und
keines hat.

**(6) Die Hausform für eine benannte Grenze.** Dieses Repo führt nicht-gedeckte
Flächen regelmäßig als **benannte** Grenze mit benanntem Wächter, statt sie
mit einem Sensor zu behaupten: [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md)
§„Was diese Regel nicht hat — die benannte Grenze" („Die durchsetzenden Leser
existieren bereits und sind Rollen, kein Werkzeug: der **Reviewer** für
Instanz A, der **Verifier** für Instanz B"), `AGENTS.md` §3.11 („diese Lücke
ist benannt, nicht still; der Wächter dort ist das Review, kein Gate"),
`AGENTS.md` §3.13 („**Kein Sensor.** … Die verfügbare Falsifikation ist die
**Messung**"), und die Grenzen 8 und 9 in
[`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md).
Der Wächter hat hier **real funktioniert**: das Review zu
`slice-096` <!-- d-check:status-provenance --> hat die drei Träger von Hand
nachgezählt („`SPEC-016` und Handbuch §5.2 stimmen mit dem Code überein —
nachgezählt") und die Paarung damit geprüft, obwohl kein Sensor sie trägt.

**(7) Warum das eine Architect-Frage ist.** Der Defekt liegt im
`Accepted`-Text der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
nicht am Code: der Implementer-Diff folgt deren §Entscheidung. `AGENTS.md`
§3.5 schließt die In-place-Korrektur aus, und
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
nimmt die **Fitness-Function-Regeln** ausdrücklich von der Zitat-Korrektur aus
— eine Schärfung des Beleg-Ankers ist eine **inhaltliche** Änderung. Modul 8
§Konflikt-Pfad (Verdikt 1) führt die Korrektur als Folge-ADR mit `Supersedes`.
Kein Reviewer→Implementer-Pfeil: es gibt keine Fixrunde am Code.

## Entscheidung

Wir wählen: **Die dritte Fitness-Function-Zeile der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) wird
ersetzt. Sie nennt die Feldmengen-Paarung als Review-Prüfpflicht und nicht
mehr `make docs-check` als ihren Träger.**

Drei Festlegungen:

1. **Die Entscheidung der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
   gilt unverändert.** Der Diskriminator („kann dieses Feld Zugangsdaten
   tragen?"), `http_addr`/`grpc_addr` als Datei-Felder, `nats_url` und die
   zwei Token-Schlüssel env-exklusiv mit eigener, den Grund nennenden
   Fehlerzeile und die Durchleitung der fünf Oberflächen-Variablen auf dem
   Datei-Pfad bleiben in Kraft. Diese ADR ändert **keine** Festlegung, keinen
   Code und keine Spec-Stelle — sie korrigiert eine Zeile ihres
   Beleg-Apparats.

2. **Der Ersatztext der dritten Fitness-Function-Zeile** lautet:

   > **Träger-Paarung:** die Feldmenge in `SPEC-016`, im Handbuch §5.2 und die
   > Schlüssel-Prüfung im Code nennen dieselben Namen — **von Hand**
   > nachzuzählen. **Kein Sensor trägt diese Zeile** (Kontext (1)–(4)):
   > `make docs-check` liest keinen Go-Quelltext, und keine
   > `structure`-Regel stellt zwei Dokumente einander gegenüber. Der Wächter
   > ist der **Review** und der **Verifier** („Prüfe die **Belege**, nicht die
   > Behauptung").

   Die erste, zweite und vierte Zeile der
   [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-Fitness-Function
   bleiben unverändert in Kraft; die **Code-Hälfte** der Paarung trägt
   weiterhin die erste und zweite Zeile (Kontext (5)).

3. **Die Grenze wird an der Zeile benannt, die sie betrifft.** Die
   Review-Prüfpflicht steht **in** der Fitness-Function-Zeile und nicht nur in
   [`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md):
   der Verifier liest die Zeile als Beleg für das Nachzug-Kriterium des
   Slice, und er soll der Zeile entnehmen, dass er sie **selbst** führen muss.
   Eine Lücke, die nur im Sensor-Dokument steht, erreicht diese Leserin nicht.
   **Ein Gate, das mit und ohne die geprüfte Eigenschaft dasselbe Ergebnis
   liefert, ist keine Deckung, sondern eine Behauptung** — es ist schlechter
   als kein Gate, weil es einen Lauf ersetzt, der sonst stattgefunden hätte.

### Was diese ADR nicht ändert

- **`SPEC-016` und das Handbuch §5.2 bleiben unberührt.** Beide tragen die
  Feldmenge bereits und stimmen mit dem Code überein; der Defekt lag allein in
  der Zeile, die diese Übereinstimmung einem Gate zuschrieb.
- **Kein Sensor wird gebaut und keiner abgeschaltet.** Der Bestand der inneren
  Gates (`make gates`) bleibt unverändert;
  [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Fitness
  Function bleibt in Kraft. Der Sensor hat nicht versagt — er war für diese
  Aussage nie zuständig; falsch war ihre Zuschreibung.
- **Keine Schwellen-Änderung** (`AGENTS.md` §3.6): nichts wird gelockert, weil
  nichts an einer Schwelle hing.
- **Kein Code-Eingriff, kein Implementer-Rückweg.** Der ausgelieferte Diff von
  `slice-096` <!-- d-check:status-provenance --> folgt §Entscheidung der
  [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md); die
  Korrektur ist reiner Text an der Entscheidungslage.
- **Die Träger-Zeile des Slice-Plans braucht keinen Nachzug** (gemessen):
  `slice-096`s <!-- d-check:status-provenance --> LP2 verlangt „Je Klasse ein
  Test", LP3 den Handbuch-Nachzug — **keine** DoD-Zeile lehnt sich an
  `make docs-check` als Beleg der Paarung. Der Plan trägt die falsche Zusage
  nicht.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: die Zeile bleibt, wie sie steht | kein Eingriff an einer `Accepted`-ADR; die Aussage („die drei nennen dieselben Namen") ist heute **wahr** | die Zeile nennt ein Gate als Träger, das mit und ohne Paarung dasselbe Ergebnis liefert (gemessen, Kontext (1)); der Verifier liest sie als Beleg und prüft damit nichts; §3.12 Instanz B verlangt für eine als Beleg gelesene Aussage ihren Beleg-Anker — der nächste Lauf müsste den Widerspruch erneut herleiten |
| B — die Zeile auf das zurücknehmen, was es gibt: Review-Prüfpflicht mit **benannter** Lücke, Code-Hälfte bei Zeile 1/2 **(gewählt)** | behebt den Defekt an seiner Stelle (der ADR-Text), kein Code-Eingriff, kein neuer Sensor; die Zeile nennt jetzt, was sie trägt **und** was sie nicht trägt; die tragende Rolle (Verifier) steht dort, wo sie gelesen wird; Hausform (Kontext (6)) | es entsteht eine benannte **Urteils**-Fläche statt eines Messwerts: die Paarung kann in einem Lauf unterbleiben, und kein Sensor rotet dann (Konsequenzen) |
| C — ein Go-Test vergleicht die drei Feldmengen | die Paarung wäre für die `SPEC-016`-Tabellenhälfte maschinell getragen; `make test` mountet `$(CURDIR):/src:ro`, der Baum läge dem Test vor (gemessen, `Makefile`) | (i) er machte den Inhalt einer **Rang-2**-Spec-Tabelle vom Code-Test abhängig — in einem GF-Repo („Doc führt, Code folgt") wäre der Code der Schiedsrichter über die Spec; (ii) die Handbuch-Hälfte ist **Prosa** (ein Satz mit sechs Schlüsseln in Inline-Code) — ein Parser darauf ist eine Formpflicht auf Prosa, genau die Klasse, die §3.12 für einen Sensor verwirft; (iii) `make a-check` sähe die Kante **nicht**: es prüft Go-Importe, kein Datei-Lesen — es entstünde ein ungewächterter Rückwärts-Verweis von Code auf Spec; (iv) eigener Slice |
| D — die `SPEC-016`-Feldtabelle als **Erzeugnis** aus dem Code, mit Sync-Gate ([`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md)-Form) | ein Erzeuger, eine Quelle; die Tabelle könnte nicht mehr driften; das Muster ist im Repo gebaut (`generated-sync`, `proto-generate`, `schema-rollout`) | [`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md) setzt einen **Erzeuger** mit netzlos-deterministischem Lauf voraus — hier gibt es keinen, und die Handbuch-Hälfte bliebe unerzeugt; vor allem kehrt es den **GF-Modus** um: die Feldtabelle ist ein normativer Träger (das `Schärft:` der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) zeigt auf `SPEC-016`), aus dem Code erzeugt würde sie vom Code regiert |
| E — die Zeile ersatzlos streichen (die Code-Hälfte trägt Zeile 1/2 ohnehin) | kleinster Eingriff; keine neue Formulierung nötig | die Paarung hätte dann **keinen** benannten Träger — weder Werkzeug noch Rolle; sie wäre eine stille Annahme, und §3.12/§3.13 verlangen, dass die Lücke **benannt** ist; der Verifier hätte keinen Zeiger auf seine eigene Prüfpflicht |

**Fazit:** B. A lässt eine Aussage stehen, die einen Beleg behauptet, den es
nicht gibt; C tauscht eine Prosa-Paarung gegen eine Formpflicht und einen
ungeprüften Code→Spec-Verweis; D kehrt die Referenz-Richtung des GF-Modus
um; E macht aus einer benannten Lücke eine stille.

## Konsequenzen

- Positiv: Kein Träger dieses Repos nennt mehr ein Gate als Deckung einer
  Prüfung, die es nicht leistet. Die Zeile, die der Verifier als Beleg liest,
  beschreibt jetzt genau das, was er zu tun hat.
- Positiv: **kein Code-Eingriff, kein Implementer-Rückweg, kein neuer
  Sensor.** §Entscheidung und §Verglichene Alternativen der
  [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
  bleiben in jeder Festlegung unberührt; geändert wird eine Zeile ihres
  Beleg-Apparats.
- Positiv: Die Korrektur folgt dem gebauten Muster für dieselbe Klasse
  ([`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md),
  eine einzelne Fitness-Function-Zeile); kein zweites Werkzeug.
- Negativ mit benannter Grenze: Die Paarung ist danach ein **Urteil**, kein
  Messwert. Sie kann in einem Lauf unterbleiben, und **kein Sensor rotet
  dann** — dieselbe Grenze, die
  [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) für §3.12 und
  `AGENTS.md` §3.13 für die bewegte Eigenschaft benennen. Das ist benannt,
  nicht wegdefiniert; der Re-Evaluierungs-Trigger greift genau hier.
- Negativ: [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
  und diese ADR müssen **in einer Zeile** zusammengelesen werden (Muster der
  übrigen engen Klausel-Korrekturen des Repos: `ADR-0048`, `ADR-0062`,
  `ADR-0063`, `ADR-0065`, `ADR-0067`).
- Folgepflicht (Planner-Zug): **keine** — der Plan von
  `slice-096` <!-- d-check:status-provenance --> lehnt sich an die ersetzte
  Zeile nicht an (gemessen, „Was diese ADR nicht ändert").
- Hinweis (Closure): Die Finding-Klasse „Beleg trägt seinen Satz nicht"
  (Review-Summary zu `slice-096` <!-- d-check:status-provenance -->) geht bei
  der Slice-Closure in §7 und von dort in den Zähler (Beobachtungs-Register);
  Träger ist die Closure, nicht diese ADR.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (kein Sensor) | **Träger-Paarung** — die Nachfolge-Zeile der dritten Fitness-Function-Zeile aus [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md): die Feldmenge in `SPEC-016`, im Handbuch §5.2 und die Schlüssel-Prüfung im Code nennen dieselben Namen — **von Hand** nachzuzählen. Der Wächter ist der **Review** und der **Verifier** („Prüfe die **Belege**, nicht die Behauptung"); **kein** Werkzeug trägt diese Zeile (Kontext (1)–(4)) | — |
| `go test ./internal/bootstrap/...` (im gepinnten Container, netzlos) | **Nicht** Gegenstand dieser Zeile: die **Code-Hälfte** der Paarung — Durchleitung und Feld-für-Feld-Vorrang je Feld, die Zugangsdaten-Klasse je Schlüssel — tragen unverändert die erste und zweite Zeile der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-Fitness-Function | `make test` (kein Gate) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Zwei benannte Trigger, sonst permanent:

1. **Ein Werkzeug wird verfügbar, das die Paarung wirklich prüfen kann** —
   etwa ein `d-check`-Modul, das Go-Quelltext liest, oder ein Erzeuger für die
   `SPEC-016`-Feldtabelle (Alternativen C/D). Dann ist diese Zeile erneut zu
   superseden: die Paarung wandert in den Sensor, und die Review-Prüfpflicht
   fällt auf das, was der Sensor deckt **nicht** — die Handbuch-Prosa und die
   Einordnung eines neuen Feldes.
2. **Die Paarung driftet real auseinander**, und ein Review oder ein Verifier
   findet es. Dann ist zu entscheiden, ob der Befund bei der benannten Grenze
   bleibt (Urteil) oder ob er Alternative C/D belegbar macht — die
   Beobachtung dafür ist der Beleg, nicht der Vorsatz.

Sonst permanent: die drei Träger nennen dieselben Namen; die Code-Hälfte
tragen die ersten beiden Fitness-Function-Zeilen der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md), die
Paarung zwischen ihnen trägt das Review.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: Review zu `slice-096` <!-- d-check:status-provenance -->, F-1 (HIGH): die dritte Fitness-Function-Zeile der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) nennt `make docs-check` als Träger der Feldmengen-Paarung. Eigene Messung an HEAD `b261dc3`: der Lauf ist mit und ohne die Paarung grün (789 Dateien, 0 Befunde, Exit 0); auch mit **allen** optionalen Modulen kein Befund zu den drei Trägern. Unabhängiger Architect-Zug ersetzt die Zeile (Verdikt 1, Modul 8 §Konflikt-Pfad): die Paarung ist Review-Prüfpflicht mit benannter Lücke, die Code-Hälfte bleibt bei den Zeilen 1/2 der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) | Review zu `slice-096` <!-- d-check:status-provenance --> |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0089` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
