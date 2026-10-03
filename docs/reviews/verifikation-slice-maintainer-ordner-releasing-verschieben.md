# Verifikations-Report: slice-maintainer-ordner-releasing-verschieben — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD- und Entscheidungs-Konformität
plus Plan-vs-Code-Diff, frischer Kontext, nach Implementierung und Review.

**Gegenstand:** `git diff de24629d HEAD` — Commits `b26460a3` (reiner Move), `13b6919b`
(Verweise, Index, Gate), `36b311b4` (Zitat-Korrektur
[`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)), `2d78f702`
(Reviews, Register), `0e7bbb6c` (Suchlauf-Messung im Plan), `28cdf10c` (Review-Report).
Der Planner-Commit `ca3f04fe` (SDK-Plan) liegt im Bereich und ist fremd, ausgeschlossen.

**Eingangs-Kontext:** Slice-Plan `slice-maintainer-ordner-releasing-verschieben`
(in-progress), [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md),
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md),
[`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md),
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
[`review-slice-maintainer-ordner-releasing-verschieben.md`](review-slice-maintainer-ordner-releasing-verschieben.md)
(0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO), `AGENTS.md` §3.1, §3.3, §3.5, §3.6, §3.7, §3.9, §3.10,
§3.12, §3.13.

---

## 1. Eigene Sensor-Belege (dieser Lauf, Exit direkt, nie durch eine Pipe, `AGENTS.md` §3.9)

| Sensor | Exit | Gedruckte Zeile |
|---|---|---|
| `make gates` | 0 | Ausgabe in Datei gesichert; Schluss `gesamt: 0 Befund(e)` (a-check), `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks` |
| `make docs-check` | 0 | `d-check: 1615 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports; nach Anlage Lauf mit diesem Report: siehe Schluss der Übergabe) |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make test-handbuch-public-doc-check` | 0 | `alle 41 Fälle bestanden` (Exit zusätzlich ungepiped gelesen) |
| `make sdk-public-doc-check` | 0 | `keine interne Kennung unter sdks` |
| `make doc-immutable RANGE=de24629d..HEAD` | 0 | `1615 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=de24629d..HEAD` | 0 | `1615 Datei(en) geprüft, 0 Befund(e)` |
| `make commit-traceability` | 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make fmt-check` | 0 | `334 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=de24629d` | 0 | keine Kandidaten |
| `make suchlauf-nachmessen PLAN=<Slice-Plan>` | 0 | `12 Zeilen stimmen` (u. a. `diff 22`, `diff 41`, `diff 0`) |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |

Echtrepo nach allen Läufen: `git status --short` leer (vor Anlage dieses Reports).

## 2. DoD Zeile für Zeile

| DoD | Befund | Beleg (eigener Lauf) |
|---|---|---|
| (1) Verschiebung rein | **erfüllt** | `git diff --stat -M de24629d..b26460a3`: `docs/{user => maintainer}/releasing.md | 0`, 1 Datei, 0 Einfügungen/Löschungen. `git log --follow --oneline docs/maintainer/releasing.md | head -3`: `13b6919b`, `b26460a3`, `603413c5` (History über den Move hinaus verfolgt). Eigener Commit, getrennt vom Inhalt (§3.3). |
| (2) Verweise nachgezogen | **erfüllt**, Rest siehe V-1, V-2 | README.md und README.de.md auf `docs/maintainer/releasing.md`; [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) drei Linkziele; zwei Review-Dateien Linkziel; `state.md` des Registers (lebende Adresse); `docs-check` Exit 0. Suchlauf `diff 22 -F docs/user/releasing.md` stimmt; Aufschlüsselung der Treffer siehe §4. |
| (3) Handbuch-Gate bereinigt | **erfüllt**, Rest V-2 | `excluded` führt nur noch die vier `*-abdeckung.md`; Tabellentest `all_files` ohne `releasing.md`, Löschfall auf `bench-abdeckung.md` umgestellt; Sensor-Vertrag Absatz 3 und Reichweite-Tabelle sowie die Gate-Zeile in `harness/README.md` an das Skript angeglichen. Gate Exit 0, Tabellentest Exit 0, geprüft bleiben dieselben drei Dateien (kein Lockern, §3.6). |
| `make gates`, `docs-check`, `doc-immutable`, `suchlauf-nachmessen` | **erfüllt** | §1 |
| Zitat-Korrektur [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) | **erfüllt mit V-3, V-4** | genau **eine** neue §Geschichte-Zeile in [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md); Commit `36b311b4` nennt `ADR-0073`; [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) und [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) ohne Diff (`git diff --stat`: nicht in der Liste). |
| Review durchgeführt | **erfüllt** | Report vorhanden; `[x]` im Plan gesetzt (geprüft: Existenz des Reports). |
| Doku-Update | **erfüllt** | `harness/README.md` (Rang 6, Gate-Zeile), README-Links. |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | **offen, Closure** | bewusst noch nicht geschrieben (Plan §7 trägt die Vorlage); §7 und §8 liefern die Vorschläge. Keine DoD-Häkchen von mir. |

## 3. Entscheidungs-Konformität

- **Auftraggeber-Entscheidungen.** `git diff de24629d HEAD -- AGENTS.md` = genau **eine** Zeile (Z. 66),
  Wortlaut „6. `docs/user/` — Operations, Quality, Version; `docs/maintainer/` — Releasing.", gleich der
  Freigabe. `harness/README.md` Rang-6-Zeile gleicher Inhalt in Tabellenform. `releasing.md` ungeteilt
  verschoben. `docs/user/version.md` unverändert, `.github/**` ohne Diff.
- **`releasing.md`-Inhalt.** `git diff -M de24629d HEAD` auf die Datei (97 % Similarity): Titel, `Version: 1.13`,
  Stand, §1-Zielgruppe, der Link `../user/version.md` (löst auf, `docs-check` Exit 0) und die Historienzeile 1.13.
  Sonst unverändert. Quelle der Wahrheit der Version bleibt `docs/user/version.md` (§2 trägt den Klartext weiter,
  wahr). Beide Ordner liegen eine Ebene unter `docs/`, die übrigen relativen Links behalten ihre Tiefe.
- **[`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Festlegung 2**
  („dauerhaft außerhalb") räumlich bestätigt, nicht entschieden und nicht gelockert: das Gate prüft dieselben
  drei Dateien und nimmt dieselben vier Erzeugnisse aus; ein `excluded`-Eintrag für eine nicht mehr vorhandene
  Datei wäre Exit 2.
- **[`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md).** Die drei Korrekturen sind Linkziele
  (Z. 48, 139, 304); Aussagegehalt von Entscheidung und Konsequenzen unverändert. Siehe V-3 zur Lage der Zeilen.

## 4. Plan-vs-Code-Diff (Plan §3 gegen `git diff --stat de24629d HEAD`)

17 Dateien. Alle Plan-Zeilen haben ihr Gegenstück: Move, `docs/maintainer/README.md` (neu), README.md,
README.de.md, [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md), zwei Review-Dateien, `state.md`,
vier Gate-Dateien, `.harness/skills/nutzerdoku-schreiben.md`, `harness/README.md`, `AGENTS.md`. Nicht im Plan
und nicht Teil des Slice: der fremde Planner-Commit `ca3f04fe` (`open/slice-sdk-meldungscodes-in-fehlertypen.md`)
und der Review-Report. Keine Änderung an `.github`, `Makefile`, `harness/mk`, `.d-check.yml`, `docs/user/`.

**Aktiver Baum, `git grep -n "docs/user/releasing"`:** ohne Records und Register nur
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Z. 212 (Klartext, beschreibt das damalige Vorbild),
[`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) (drei Linktexte, ein Klartext),
[`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Z. 87 (Klartext),
der Plan selbst, `done/` (Records), Register (Klartext, `state.md` nachgezogen) und der fremde SDK-Plan
(`open/slice-sdk-meldungscodes-in-fehlertypen.md` Z. 70, 386: gemeldet, Frist Closure, Planner zieht nach).
`git grep -n "releasing" -- .github Makefile harness/mk tools .claude .harness/skills .d-check.yml`: keine
Treffer außer den Absichten (Gate-Skript, Tabellentest; `nutzerdoku-schreiben.md` nennt die Maintainer-Doku
ohne Dateinamen). Kein Workflow, kein Make-Ziel liest `releasing.md`. `git grep -n releasing -- docs/user`: leer.

## 5. Mutationsproben (einzeln, auf Kopien aus `git archive HEAD` im Scratchpad; Mutation per `awk` nach stdout in die Kopie, kein `sed -i`; Echtrepo unberührt)

| # | Mutation | Erwartung | Beobachtet |
|---|---|---|---|
| 1 | `releasing.md` wieder in `excluded` | Exit 2 | Exit 2, `genannte Datei fehlt: docs/user/releasing.md` |
| 2 | `bench-abdeckung.md` aus `excluded` | Tabellentest rot | Exit 1, mehrere Fälle `FEHLER` (u. a. `fehlende ausgenommen-genannte Datei -> Exit 0, erwartet 2`; Gate-Meldung `unklassifizierte Datei: docs/user/bench-abdeckung.md`) |
| 3 | Kennung `ADR-0051 slice-x` an `docs/maintainer/releasing.md` | Exit 0 (außerhalb) | Exit 0 |
| 4 | neues `docs/user/neu.md` mit Kennung | Exit 2 (unklassifiziert) | Exit 2, `unklassifizierte Datei: docs/user/neu.md` |

Das Gate hält seine Grenzen: die Maintainer-Doku liegt außerhalb, eine neue Datei unter `docs/user/` erzwingt
die Klassifikation. Gemessen an Kopien; die Verallgemeinerung auf andere Listenänderungen ist *hergeleitet*.

## 6. Befunde (V-Zählung dieses Reports)

- **V-1 (LOW, = Review F-1, Entscheidung des Hauptlaufs bewertet).** Drei Links in
  [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) tragen den sichtbaren Text
  `docs/user/releasing.md` bei neuem Ziel. Bewertung: der Hauptlauf-Beschluss (sichtbarer Linktext gehört
  nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) ins Verweisgerüst, in der Closure
  innerhalb derselben §Geschichte-Zeile mitkorrigiert) ist **tragfähig**: der Linktext nennt keinen anderen
  Referenten, er bezeichnet dieselbe Datei an ihrem neuen Ort; Aussage und Referent bleiben gleich. Voraussetzung:
  die Zeile nennt dann beide Korrektur-Arten (Linkziel, Linktext).
- **V-2 (LOW, neu).** Der Kopfkommentar von `tools/harness/handbuch-public-doc-check.sh` (Z. 10) sagt
  „`excluded` nennt Maintainer-Doku und Runner-Erzeugnisse". Die Liste führt nur noch die vier Erzeugnisse; der
  Satz beschreibt einen Zustand, der nicht mehr da ist (`AGENTS.md` §3.7, Ist-Zustand). `kommentar-kennungen`
  fängt das nicht (prüft die Kennungsform). Der Review hat es nicht gemeldet. Korrektur ist ein Halbsatz
  („`excluded` nennt Runner-Erzeugnisse"); gehört in die Closure oder eine kleine Nacharbeit des Implementers.
- **V-3 (LOW, neu, Frage an den Architect).** Die Korrekturen in
  [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) Z. 139 (§Entscheidung 3) und Z. 304
  (§Konsequenzen) sitzen in Abschnitten, die [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1 als unberührbar listet; Z. 48 liegt im §Kontext. Der Wortlaut von `AGENTS.md` §3.5 nennt
  „Linkziele" ausdrücklich zulässig, der Referent ist unverändert. Das ist die Spanne, die
  `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` trägt (Kurzform gegen Abschnitte-Liste; engere Lesart
  für neue Fälle laut Architect: Folge-ADR). Der Review hat sie nicht benannt. Nicht blockierend (der Plan sah die
  Korrektur vorab vor, der Linkziel-Fall ist im Wortlaut von §3.5 gedeckt), aber ein **zweites Vorkommen** der
  Beobachtung (siehe §7).
- **V-4 (INFO).** Die §Geschichte-Zeile in
  [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) nennt `Commit b26460a3` (den Move), der
  korrigierende Commit ist `36b311b4`. Der eigene Hash lässt sich nicht im selben Commit nennen; die Zeile bezeichnet
  den Auslöser, nicht den Träger. Bei der V-1-Nachkorrektur in der Closure kann die Zeile auf `36b311b4`
  (und den Folge-Commit) zeigen. Kein Verstoß gegen §3.5 (Beleg ist eine Commit-Kennung), nur ungenau.
- **Bestätigt aus dem Review:** F-2 (Plan-Prosa, Zahlen stimmen, 12 Suchlauf-Zeilen), F-3 (SDK-Plan, Hinweis
  an den Planner), F-4 ([`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
  Klartext, Record). Eigene Nachmessung ohne Abweichung.

## 7. Register-Fortschreibung (für die Closure, ich ändere nichts)

- `BEO-PGC/plattform-verhalten-nur-vom-betreiber-pruefbar`: die Adresse in `state.md` ist bereits auf
  `docs/maintainer/releasing.md` nachgezogen (geprüft, Zählerzeile unverändert 1×); keine weitere Fortschreibung.
- `BEO-PGC/zitat-nennt-die-falsche-stelle`: V-1 (Review F-1) ist **kein** Vorkommen. Die Klasse trägt den Verweis,
  der eine falsche oder nicht aufgeschlagene Stelle nennt (Abschnittsnummer, Kennung, Entscheidungslage). Hier nennt
  das Linkziel die **richtige** Datei, nur der sichtbare Text blieb stehen: unvollständige Korrektur, keine falsche
  Stelle. Kein Eintrag dort.
- `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`: V-3 ist ein Vorkommen (Zitat-Korrektur in
  §Entscheidung/§Konsequenzen einer `Accepted` ADR); Vorschlag **2×**, Beleg-Datei
  `evidence/slice-maintainer-ordner-releasing-verschieben.md` (Zählerzeile laut `state.md` derzeit 1×,
  *übernommen*). Ob für „sichtbarer Linktext gegen Linkziel" (V-1) ein neuer Eintrag lohnt: es wäre ein Vorkommen
  1× einer neuen Klasse; Vorschlag: nur im Closure-Feld des Slice vermerken, nicht eigenständig anlegen, solange
  es das erste ist.
- Lerneintrag für die Closure, den die Messung trug: der Suchlauf trennt Link-Form (bricht `docs-check`) von
  Klartext (bleibt Record). Die Messung am Parent zählte 3 Review-Links, tatsächlich 2 Linkziel-Änderungen (F-2);
  der Lerneintrag nennt diese Zahl als Lehre.

## 8. §6-Ausgang-Vorschläge (Slice-Plan §6, Entscheidung beim Planner)

| Risiko | Vorschlag | Beleg |
|---|---|---|
| Externe Rückwärts-Links | **weiter offen** (bewusst in Kauf genommen) | kein Beleg für externe Links im Repo; GitHub leitet nicht weiter |
| Docker-Hub-Beschreibung | **weiter offen** bis zum nächsten Release oder `workflow_dispatch` | `git diff --name-only de24629d HEAD -- .github` leer; `hub-description.yml` unverändert; `AGENTS.md` §3.10 nicht ausgelöst |
| Klartext-Pfade in Records und `Accepted` ADRs | **weiter offen** als Nebenwirkung der Immutabilität | §4 |
| Gate-Ausnahme räumlich bestätigt | **entfallen** | Mutationen 1 bis 4; dieselben drei Dateien geprüft |
| Reihenfolge Move/Gate | **entfallen** | Endstand gated, Gate Exit 0; `commit-traceability` prüft Messages |

Der Plan §6 dokumentiert die beiden externen Punkte wahr (beide als „weiter offen" mit Begründung, §3.10 ausdrücklich
nicht ausgelöst). Der Stand dieses Slice ist nicht gepusht; die Aussage „kein Workflow berührt" ist am Diff
**gemessen**, nicht an einem CI-Lauf.

## 9. Verdikt

**Bestanden**, mit Bedingungen für die Closure (keine Rückgabe an den Implementer zwingend):

1. V-1 (= Review F-1): sichtbaren Linktext in [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)
   im selben Zug mitkorrigieren, innerhalb derselben §Geschichte-Zeile (beide Korrektur-Arten nennen); dabei V-4 (Hash
   `36b311b4`) mitnehmen.
2. V-2: Kopfkommentar des Gate-Skripts anpassen (Halbsatz), Implementer-Nacharbeit oder Closure.
3. V-3: Architect-Notiz zur Reichweite (Linkziel in unberührbaren Abschnitten) bzw. Zähler 2× im Register; nicht
   blockierend.
4. Closure-Pflichten (Notiz, Register, §6-Ausgänge, drei Paarungen) und der Nachzug der zwei Klartext-Nennungen im
   fremden SDK-Plan (`open/slice-sdk-meldungscodes-in-fehlertypen.md`) liegen beim Planner.

Offene Punkte für den Planner: Punkte 1 bis 4, dazu die zwei bewusst offenen externen Risiken.
