# Verifier-Report: slice-024 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code inkl. Plan-Nachzug), §6 (Risiko-Ausgang,
bleibt Planner-Entscheidung), §8 (Sub-Area-Prüfung), sowie
Entscheidungs-Konformität gegen [ADR-0049](../plan/adr/0049-replication-fehlerklassen-schwellen.md) (Accepted) und Immutabilität von
`ADR-0023` (Accepted, Hard Rule 3.5). Nicht geprüft: Diff gegen Plan/Hard
Rules im Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe, bereits
erledigt, siehe [`review-slice-024.md`](review-slice-024.md), 0 HIGH/0
MEDIUM/0 LOW/1 INFO, kein Merge-Blocker), realer Bedarf (Validator — hier
nicht einschlägig, kein MVP-Grenz-Slice).

**Gegenstand:** drei inhaltstragende Commits im Slice-Fenster: `741f845`
(`ADR-0049` neu), `52e207d` (`spec/pflichtenheft.md` `SPEC-008`/`SPEC-013`
geschärft), `988ab9a` (Slice-Plan-Nachzug: DoD-Häkchen, Plan-Nachzug-Tabelle).
Zusätzlich `2db8408`/`1f64ee3` (Review-Report + Docs-Check-Fix darauf) sowie
drei reine Lifecycle-Move-/Ruhe-Marker-Commits (`37b3904`, `d9de434`,
`8da1ebc`, `open→next→in-progress`).

**Grundsatz:** Keine Behauptung wurde übernommen — jeder Beleg unten wurde
in diesem Lauf selbst gelesen oder ausgeführt (Dateien im Volltext, `git
show`/`git log --follow`, `make gates` selbst gestartet).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-024-adr-fehlerklassen-schwellen-praezisierung.md`)
- `docs/plan/adr/0049-replication-fehlerklassen-schwellen.md` (Accepted,
  Volltext)
- `docs/plan/adr/0023-fehlerklassifikation.md` (Accepted, Volltext, auf
  Widerspruch/Immutabilität geprüft)
- `docs/plan/adr/README.md` (Index, Zeile [`ADR-0049`](../plan/adr))
- `spec/pflichtenheft.md` §3 (`SPEC-013`, Zeile 224), §4 (`SPEC-008`, Zeile
  244), §5 (Metriken-Prosa, Zeilen 270–272) im Volltext gelesen
- `docs/reviews/review-slice-024.md` (0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO)
- `docs/plan/planning/welle-7.md` (§Welle-Ziel)
- `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/`,
  `.../walsender-wirksamkeit/`, `.../rollen-test-abdeckungsluecken/`,
  `.../dod-checkbox-nachzug/`, `.../plan-nachzug/` (Register-Volltext)
- `docs/plan/planning/open/slice-025-metrik-wal-rueckstand.md`,
  `.../slice-026-schwellen-ueberwachung-kontrollierte-fortsetzung.md`
  (Existenz-Beleg)
- `AGENTS.md` §3.5 (ADR-Immutabilität), `harness/README.md` §Sensors
- `git log`/`git show`/`git diff` über den vollen Commit-Verlauf
  (`37b3904^..1f64ee3`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 232 Datei(en), 0 Befund(e)` (Standardlauf, alle Module inkl. `immutable`/`matrix`/`versions`/`structure`) · `d-check --range HEAD~5..HEAD` (Modul `commits`): `232 Datei(en), 0 Befund(e)` · `commit-traceability.sh "HEAD~5..HEAD"`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: `0 Befund(e)` | **0** |
| `git log --follow -- docs/plan/adr/0049-…md` | genau 1 Commit (`741f845`) — kein nachträglicher Inhaltsedit | — |
| `git log --follow -- docs/plan/adr/0023-fehlerklassifikation.md` | 2 Commits, beide vor `slice-024` (`e6c42a3` Anlage, `c477758` Sensor-Härtung); kein Commit im Slice-Fenster berührt die Datei | — |
| `git show 52e207d -- spec/pflichtenheft.md` | zeigt exakt zwei Zeilenänderungen: `SPEC-013`-Zeile (§3) und `SPEC-008`/`replication`-Zeile (§4) | — |
| `grep -n 'ADR-[0-9]' spec/pflichtenheft.md` | kein Treffer | **1** (leer = bestanden) |
| `grep -l '^\*\*Bezug:\*\*' docs/plan/adr/*.md \| xargs grep -c '#'` | 0 Treffer mit `#anker` in `Bezug`/`Schärft` über alle 49 ADR-Dateien | — |
| `grep -c '^\| ADR-' docs/plan/adr/README.md` | 49 Zeilen (48 Bestand + `ADR-0049`) | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | ADR `Accepted`: entscheidet Sentinel-Trennung (a) und WAL-Rückstand-Schwellen (b) | **bestätigt** | `docs/plan/adr/0049-….md` Zeile 3 `Status: Accepted`; §Entscheidung trifft wörtlich beide Festlegungen: (a) `mapper.Err*` = Stream-Ordnungs-Verletzung hart, `receive.ErrReplication`/`outbound.ErrReplication` = Transport-/Verbindungsstörung Schwellen-Kandidat; (b) Warn 100 MiB / Fehler 1 GiB für `cdc_wal_retention_bytes`. `docs/plan/adr/README.md` trägt die Indexzeile korrekt |
| 2 | `SPEC-008`-Zeile `replication` präzisiert: beide Fehlerklassen-Unterarten getrennt benannt | **bestätigt** | `spec/pflichtenheft.md:244` trägt wörtlich „**Stream-Ordnungsverletzung** … oder **Transport-/Verbindungsstörung**" mit getrennter Aktionsspalte (harter Abbruch vs. Schwellen-Überwachung/kontrollierte Fortsetzung); durch `git show 52e207d` als tatsächliche Diff-Zeile verifiziert, nicht nur behauptet |
| 3 | `SPEC-013` um `cdc_wal_retention_bytes`-Eintrag ergänzt, im bestehenden §3-Eintrag statt neuer §5-Zeile | **bestätigt** | `spec/pflichtenheft.md:224`: `SPEC-013`-Zeile trägt jetzt `CDC_THRESHOLDS` (vorher `CDC_LAG_THRESHOLDS`) mit „WAL-Rückstand Warn > 100 MiB · Fehler > 1 GiB" zusätzlich zum bestehenden Capture-Lag-Wert. Eigene Prüfung der Gliederung: `SPEC-013` steht tatsächlich in §3 *Defaults und Konstanten* (Zeile 212), nicht §5; §5-Prosa (Zeile 270–272) verwies bereits **vor** diesem Slice auf `SPEC-013` als Initialwert-Quelle „für WAL-Rückstand und Capture-Lag" — die im Plan-Nachzug behauptete Spec-interne Inkonsistenz ist damit durch eigene Lektüre des Vorzustands (aus dem Diff-Kontext) bestätigt, keine neue §5-Zeile wurde angelegt |
| 4 | `make gates` grün | **bestätigt** | eigener Lauf in dieser Sitzung, alle vier inneren Gates 0 Befunde (siehe Sensor-Tabelle) |
| 5 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **materiell erfüllt, Checkbox nicht nachgezogen — siehe V-1** | `docs/reviews/review-slice-024.md` existiert real (Commits `2db8408`, `1f64ee3`), 0 HIGH/0 MEDIUM/0 LOW/1 INFO, „kein Merge-Blocker" — die Bedingung ist materiell erfüllt, aber die Checkbox in §2 steht weiterhin auf `[ ]` |
| 6 | Kein Doku-Update jenseits von `spec/pflichtenheft.md`/ADR-Index nötig | **bestätigt** | `git diff --stat 37b3904^..1f64ee3` zeigt nur: neue ADR-Datei, ADR-Index, Slice-Plan selbst, `spec/pflichtenheft.md`, Review-Report, plus die lifecycle-eigene `roadmap.md`-Ruhe-Marker-Zeile (Housekeeping, kein inhaltliches Doku-Update). Kein `docs/user/*`, kein `harness/README.md`, kein `AGENTS.md` berührt — konsistent mit „kein Betreiber-sichtbarer Vertrag geändert" |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin vollständig den Platzhalter-Vorlagentext — Planner-Arbeit, noch nicht fällig vor `git mv` |
| 8 | Reconciliation-Register fortgeschrieben, falls einschlägig | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht; Repo durchgehend Greenfield (`harness/conventions.md` §Modus-Deklaration) |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `BEO-PGC/spec008-replication-luecke/state.md` trägt unverändert „weiter offen", 1× (`evidence/slice-020.md`) — Plan-Nachzug kündigt genau das an (Planner-Arbeit bei Closure); Datei real unverändert seit vor diesem Slice bestätigt |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Risiken tragen weiterhin den Platzhalter `<bei Closure einzutragen>` — siehe Abschnitt „§6 Risiken" unten |
| 11 | Drei Paarungen getragen | **korrekt offen** | Slice ist wellenlos in der internen Zählung nicht zutreffend — Slice gehört zu `welle-7` (Kopf-Feld „Welle: welle-7"); die Paarungen laufen laut DoD-Zeile bei dieser Closure selbst (Repo-Fußnote „ohne Wellen-Betrieb hier geprüft" bezieht sich auf den Slice, nicht auf `welle-7`s eigene spätere Welle-Closure) — noch nicht fällig, da Closure aussteht |

**Zwischenstand: 6/11 Kriterien materiell erfüllt und in diesem Lauf selbst
nachgeprüft, 1 Item korrekt entfallen (8), 4 Items regulär offen als
Planner-Closure-Arbeit (7, 9, 10, 11). Punkt 5 trägt zusätzlich einen
eigenen Befund (V-1): materiell erfüllt, aber Checkbox nicht nachgezogen.**

## Plan-vs-Code-Diff (gegen Plan §1/§3 inkl. Plan-Nachzug)

`git diff --stat 37b3904^..1f64ee3` zeigt genau: neue ADR-Datei, ADR-Index,
Slice-Plan-Datei, `spec/pflichtenheft.md`, Review-Report — keine
unangekündigte Datei, kein Produktionscode berührt (§1 Out-of-Scope „keine
Implementierung" gewahrt).

Vier im Plan-Nachzug (§3) dokumentierte Abweichungen vom Ursprungsplan,
einzeln geprüft:

1. **ADR-Nummer `0049` statt Platzhalter `00NN`** — plausibel und trivial:
   `docs/plan/adr/README.md` zeigt `ADR-0048` als letzten Bestandseintrag
   vor `ADR-0049`; `0049` ist tatsächlich die nächste freie Nummer. Kein
   Konfliktpotential, reine Mechanik.
2. **`SPEC-013` in §3 statt §5** — plausibel und durch eigene Lektüre der
   Gliederung bestätigt (siehe DoD-Punkt 3 oben): Die Behauptung, §5 habe
   den Anspruch bereits vor diesem Slice getragen, ist wörtlich im
   Diff-Kontext sichtbar (unveränderte Zeilen 270–272) und stimmt.
3. **Beobachtungs-Register-Datei nicht geändert** — plausibel: Der Plan
   selbst benennt dies als „Planner-Arbeit bei der Closure, nicht Teil
   dieses Architect-Laufs"; eigene Prüfung bestätigt, die Datei ist
   tatsächlich unverändert seit vor `slice-024`.
4. **`Bezug:`/`Schärft:`-Felder ohne `#anker`-Suffix** — plausibel und
   durch eigene Stichprobe über alle 49 ADR-Dateien bestätigt: kein
   einziges `Bezug:`/`Schärft:`-Feld in diesem Repo führt einen
   `#anker`-Suffix. Die Begründung „Bestandskonsistenz, keine aktivierte
   Anker-Reifestufe" ist tragfähig — ein Abweichen hätte hier eine
   Insel-Konvention erzeugt, die kein anderes ADR-Dokument teilt.

Alle vier Abweichungen sind nachvollziehbar begründet und stimmen mit dem
tatsächlich gelieferten Diff überein — keine unbegründete oder unplausible
Abweichung gefunden.

## ADR-Konformität

- **`ADR-0049` entscheidet beide DoD-verlangten Fragen wörtlich** (siehe
  DoD-Punkt 1) — im Volltext gelesen, nicht nur die Kopfzeile.
- **`ADR-0023` bleibt unverändert seit `Accepted`:** `git log --follow --
  docs/plan/adr/0023-fehlerklassifikation.md` zeigt zwei Commits, beide
  vor dem Slice-Fenster (`e6c42a3` Anlage, `c477758` Sensor-Härtung außerhalb
  des ADR-Inhalts) — kein Commit in `37b3904^..1f64ee3` berührt die Datei.
  Inhaltlich: `ADR-0023` legt nur die sieben Klassennamen fest (`transient`
  … `internal`) und trifft **keine** Aussage über Sub-Verhalten innerhalb
  einer Klasse — eigene Lektüre bestätigt, `ADR-0023` §Entscheidung nennt
  ausschließlich die Klassennamen, keine Sentinel-Zuordnung, keine
  Schwellenwerte. `ADR-0049` ändert weder einen Klassennamen noch die
  `replication`-Zugehörigkeit eines Sentinels; sie präzisiert nur die
  Handlungskonsequenz **innerhalb** der bereits bestehenden Klasse
  `replication`. Kein `Supersedes ADR-0023`-Bedarf, keiner behauptet, Hard
  Rule 3.5 gewahrt — eigenständig bestätigt, nicht vom Reviewer-Befund
  übernommen.
- **Re-Evaluierungs-Trigger korrekt differenziert:** (a) `permanent`
  (Adapterstruktur-abhängig, nicht Zahlenwert-abhängig), (b) explizit
  **nicht** permanent, mit konkretem Auslöser (Ende-zu-Ende-Test aus
  `welle-7` §3 oder Produktionsbetrieb widerlegt die Werte → Folge-ADR
  `Supersedes ADR-0049` nur für (b)) — im ADR-Volltext (§Re-Evaluierungs-
  Trigger) gelesen und bestätigt.
- **Folge-Slices existieren real:** `slice-025-metrik-wal-rueckstand.md`
  und `slice-026-schwellen-ueberwachung-kontrollierte-fortsetzung.md`
  liegen tatsächlich in `docs/plan/planning/open/` (per `ls` bestätigt).

## Referenz-Richtungs-Regel

`grep -n 'ADR-[0-9]' spec/pflichtenheft.md` liefert **keinen** Treffer im
gesamten Dokument — die Sperre `{from: spec, to: adr, allow: false}`
(`.d-check.yml`) ist nicht verletzt, durch eigenen Lauf bestätigt, nicht nur
den Review-Negativbefund übernommen.

## §6 Risiken — Ausgang bei diesem Verifikationsstand

Beide im Slice-Plan §6 genannten Risiken tragen **weiterhin den
unausgefüllten Platzhalter** `<bei Closure einzutragen>`:

1. Byte-Schwellenwert als Stil-Analogie ohne Produktionsdaten.
2. Sentinel-Trennung könnte einen dritten Fall brauchen.

Beide sind **nicht** durch einen formalen §6-Ausgang aufgelöst. Inhaltlich
sind sie jedoch in `ADR-0049` selbst bereits **materiell** adressiert (nicht
verwechseln mit dem formalen §6-Ausgang, der eine separate Planner-Pflicht
bleibt):

- Risiko 1 ist in `ADR-0049` §Verglichene Alternativen (b) explizit als
  „nicht aus Produktionsmessungen abgeleitet, reine Stil-Analogie" benannt
  und trägt einen korrekten, nicht-permanenten Re-Evaluierungs-Trigger
  (Ende-zu-Ende-Test aus `welle-7` §3 widerlegt ggf. die Werte → Folge-ADR).
  Der plausible spätere Ausgang wäre „weiter offen → Register" oder
  „eingetreten → Folge-ADR", je nachdem was `slice-026`s Test zeigt — aber
  das ist Planner-Urteil, nicht dieser Lauf.
- Risiko 2 ist in `ADR-0049` §Konsequenzen (Negativ-Bullet: „ein neuer
  `replication`-Sentinel … braucht dieselbe Einordnungsentscheidung erneut")
  als strukturelle Konsequenz benannt, aber nicht als abgeschlossen erklärt.

**Für die Closure bleiben beide Risiken unbeantwortet** im Sinn der
formalen DoD-Pflicht (Punkt 10) — Modul 5 verlangt, dass kein Slice nach
`done/` geht, während ein Risiko ohne Ausgang dasteht. Das ist zum
jetzigen Zeitpunkt (`in-progress`, kein `git mv` erfolgt) kein Verstoß,
sondern der erwartete Zwischenstand.

## Eigene Befunde

### V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen, obwohl die Bedingung materiell erfüllt ist

- `kategorie`: LOW (kein Blocker für die Verifikation selbst — die
  zugrunde liegende Tatsache ist real geprüft und positiv; relevant für
  die Belastbarkeit der DoD-Checkboxen als Zustands-Beleg)
- `pfad`: `docs/plan/planning/in-progress/slice-024-….md:117-119`
  (Checkbox weiterhin `[ ]`) vs. `docs/reviews/review-slice-024.md`
  (Report existiert, 0 HIGH/0 MEDIUM/0 LOW/1 INFO, „kein Merge-Blocker")
- `befund`: Der Review ist real abgeschlossen und ohne Nacharbeitsbedarf
  (kein HIGH/MEDIUM/LOW). Die DoD-Zeile „Review durchgeführt, Report unter
  `docs/reviews/` liegt vor … kein Self-Review" ist damit materiell erfüllt,
  die Checkbox in §2 wurde dafür aber nicht auf `[x]` gesetzt — weder im
  Review-Commit (`2db8408`) noch im Nachfolge-Commit (`1f64ee3`). Das ist
  dieselbe Finding-Klasse, die im Beobachtungs-Register bereits als
  **verkörpert** geführt wird: `BEO-PGC/dod-checkbox-nachzug` (Zähler 3×,
  Ausgang *verkörpert* seit `welle-5`, Zielort
  `.claude/commands/implement-slice.md` — dort Schritt 21 „Fixrunden-
  Checkbox-Nachzug": „Löst eine Fixrunde nach Reviewer-Findings einen
  bislang offenen DoD-Punkt auf … wird die zugehörige Checkbox im
  Fixrunden-Commit mitgesetzt". Hier gab es zwar keine Fixrunde im engeren
  Sinn (0 HIGH/MEDIUM/LOW, kein Nacharbeitsbedarf), aber genau dieselbe
  Konsequenz träfe zu: Der DoD-Punkt ist jetzt wahr und sollte es in der
  Datei auch sein.
- `verifizierbar`: ja — `grep -n '^\- \[' docs/plan/planning/in-progress/slice-024-….md` zeigt die Checkbox weiterhin `[ ]`; `git log` zeigt keinen Commit nach `2db8408`, der die Zeile ändert, außer der reinen Docs-Check-Korrektur `1f64ee3` (nur `review-slice-024.md` selbst).
- **Für die Closure:** Kein Blocker — der Planner sollte die Checkbox vor
  `git mv` nach `done/` auf `[x]` setzen (materiell gedeckt); optional ein
  weiterer `evidence/`-Eintrag für `BEO-PGC/dod-checkbox-nachzug`, da das
  Muster (Review clean, aber Checkbox nicht nachgezogen) eine im Register
  noch nicht ausdrücklich erfasste Variante des bereits embodied Musters ist
  — Urteil bleibt beim Planner, ob dies als neues Auftreten zählt oder
  als dieselbe bereits abgedeckte Klasse gilt.

## Negativbefunde

- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde in
  allen vier Gates (Standardlauf und `--range HEAD~5..HEAD`).
- geprüft, ohne Befund: **ADR-Immutabilität `ADR-0023`** — `git log
  --follow` zeigt keine inhaltliche Änderung im Slice-Fenster; Volltext-
  Vergleich bestätigt, `ADR-0049` widerspricht `ADR-0023` inhaltlich nicht
  (siehe §ADR-Konformität oben).
- geprüft, ohne Befund: **Referenz-Richtung `spec/pflichtenheft.md` →
  ADR** — kein `ADR-\d{4}`-Token im gesamten Dokument, eigener `grep`-Lauf.
- geprüft, ohne Befund: **`SPEC-008`/`SPEC-013`-Präzisierung, Wortlaut** —
  beide Zeilen im Volltext gelesen (`spec/pflichtenheft.md:224,244`),
  Inhalt deckt sich exakt mit der DoD-Behauptung.
- geprüft, ohne Befund: **§1-Abgrenzung gewahrt** — kein Implementierungs-
  Code berührt, keine Verhaltensänderung an den Mapper-Sentinels, keine
  Lastenheft-Änderung; `git diff --stat` bestätigt.
- geprüft, ohne Befund: **Plan-Nachzug, alle vier Abweichungen** — einzeln
  gegen den tatsächlichen Diff/Dateibestand geprüft, alle plausibel (siehe
  §Plan-vs-Code-Diff).
- geprüft, ohne Befund: **Folge-Slice-Existenz** — `slice-025`/`slice-026`
  existieren real in `docs/plan/planning/open/`.
- geprüft, ohne Befund: **ADR-Index-Eintrag** — `docs/plan/adr/README.md`
  trägt die korrekte Zeile für `ADR-0049`, Nummer korrekt als nächste
  freie nach `0048` vergeben.
- geprüft, ohne Befund: **Beobachtungs-Register unverändert wie
  angekündigt** — `BEO-PGC/spec008-replication-luecke/state.md` trägt
  weiterhin `weiter offen`, 1× (`evidence/slice-020.md`); keine
  unangekündigte Änderung.
- geprüft, ohne Befund: **Traceability** — alle drei inhaltstragenden
  Commit-Betreffs nennen `ADR-0049`; kein `SPEC-*`/`ARC-*` im Betreff;
  `make commit-traceability` (Teil von `make gates`) lief grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 6/11 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft, 1 Item korrekt entfallen (Reconciliation-Register,
Greenfield), 4 Items regulär offen als Planner-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register, §6-Risiko-Ausgänge, drei
Paarungen). Kein DoD-Defekt im Sinn eines unbelegten „bestätigt"-Punkts —
der einzige eigene Befund (V-1, LOW) präzisiert eine Checkbox-Diskrepanz,
die die materielle Erfüllung des zugrunde liegenden DoD-Punkts nicht
widerlegt.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit einem offenen
Klein-Befund (1 LOW, V-1; non-blocking für diesen Slice, aber empfohlene
Korrektur vor `git mv` nach `done/`).** `ADR-0049` entscheidet real und
im Volltext beide DoD-verlangten Fragen (Sentinel-Trennung, WAL-Rückstand-
Schwellen); `SPEC-008`/`SPEC-013` sind im tatsächlichen Diff präzisiert,
nicht nur behauptet; `ADR-0023` ist nachweislich unverändert und inhaltlich
nicht widersprochen; die Referenz-Richtungs-Regel ist gewahrt; `make gates`
ist in diesem Lauf selbst grün gelaufen; alle vier Plan-Nachzug-Abweichungen
sind einzeln geprüft und plausibel.

**Plan-vs-Code-Diff:** vollständige Deckung — keine unangekündigte Datei,
kein Produktionscode berührt, §1-Abgrenzung gewahrt.

**Closure-Bereitschaft: noch nicht — vier reguläre Planner-Schritte
stehen aus**, keiner davon ein Defekt an bereits gelieferter Substanz:

1. Closure-Notiz §7 schreiben (Was hat funktioniert / anders als geplant /
   Steering-Loop-Eintrag / Beobachtungs-Register / Folge-Slices / §6-Risiken
   / Drei Paarungen).
2. Beide §6-Risiken auf einen der drei Ausgänge (eingetreten / entfallen /
   weiter offen) setzen — beide sind in `ADR-0049` bereits materiell
   adressiert (siehe §6-Abschnitt oben), der **formale** Ausgang fehlt noch.
3. Beobachtungs-Register fortschreiben — mindestens die Frage entscheiden,
   ob dieser Slice `BEO-PGC/spec008-replication-luecke` einen neuen
   `evidence/`-Eintrag hinzufügt oder unverändert bei „weiter offen" bleibt
   (der Slice selbst löst die Beobachtung noch nicht auf — Umsetzung folgt
   erst mit `slice-025`/`slice-026`).
4. Drei Paarungen (Anker · Folge-Slice · Register) — laufen bei dieser
   Closure, da der Slice zu `welle-7` gehört, deren eigene Welle-Closure
   erst später ansteht; die Slice-Closure selbst prüft sie hier trotzdem
   nicht doppelt, da `welle-7` noch offen ist (Modul 6 §Wann Arbeit eine
   Welle braucht — Paarungen der einzelnen Slice-Closure sind für
   wellenlose Slices vorgesehen; für `welle-7`s Slices verschiebt sich diese
   Pflicht auf die Welle-Closure. **Zur Klärung an den Planner:** DoD-Punkt
   11 des Slice-Plans benennt beide Fälle nebeneinander („hier geprüft" vs.
   „von der nächsten Welle-Closure") — bei einem wellen-zugehörigen Slice
   wie diesem ist die zweite Variante einschlägig, die Checkbox sollte das
   beim Nachzug entsprechend vermerken statt offen zu bleiben).

Zusätzlich empfohlen: V-1 (Checkbox-Nachzug für den bereits abgeschlossenen
Review) vor `git mv` nach `done/` beheben.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert.

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: 54 Dateien OK; `d-check` Standardlauf: 232
Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 232
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits; `a-check`: 0
Befunde). `git status` am Ende dieses Laufs sauber (keine Arbeitsverzeichnis-
Änderung durch die Verifikation selbst).
