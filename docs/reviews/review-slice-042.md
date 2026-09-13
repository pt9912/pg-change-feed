# Review-Report: slice-042 — 2026-09-13

**Review-Art:** Doku — geprüft gegen Plan (`slice-042`, §1/§2/§6/§8) und
`AGENTS.md` §3.7 (Kommentar-/Prosa-Disziplin, hier auf Markdown-Prosa
angewandt) sowie die reale Implementierung aus `slice-037`
(Baseline-Regelwerk `modul-10-review-harness.md`, Modul 10 §Drei
Review-Arten: Reviewer prüft gegen Plan/ADR, Maintainability — nicht
gegen DoD).

**Gegenstand:** Commit `294d165d76ae6a2e31131a65ae5f9d204e2bc9f1`
(`docs(user): SQL-Administration im Benutzerhandbuch nachdokumentiert
(LH-FA-ADM-001, LH-FA-CFG-002, ADR-0050)`) — geändert:
`docs/user/benutzerhandbuch.md` (zwei neue §4-Abschnitte „Tabelle live
aktivieren"/„Tabelle deaktivieren", Version 1.6 → 1.7, neue
Changelog-Zeile),
`docs/plan/planning/in-progress/slice-042-handbuch-sql-administration.md`
(DoD-Häkchen 1–4/6).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-042-handbuch-sql-administration.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6
  Risiken, §8)
- `docs/plan/planning/done/slice-037-administrations-goroutine-live-reload.md`
  (vollständig) — die reale Implementierungsbeschreibung, gegen die der
  neue Prosa-Text sachlich korrekt sein muss (Antrags-Fluss, `LISTEN`/
  `NOTIFY` + Fallback-Poll, keine Sofortwirkung)
- `tools/schema/nacharbeit-administration.sql` (vollständig) — exakte
  Signaturen/Parameter-Reihenfolge von `cdc.enable_table`/
  `cdc.disable_table`
- `internal/application/usecase/enable/service.go`,
  `internal/application/usecase/disable/service.go` — geprüft, ob
  `EnableTableUseCase`/`DisableTableUseCase` (von der
  Administrations-Goroutine aufgerufen, unverändert seit `slice-037`)
  weitere Vorbedingungen tragen, die der neue Doku-Text nennen müsste
- `AGENTS.md` §3.7 (Kommentar-/Prosa-Disziplin) — auf die neue
  Markdown-Prosa angewandt, nicht nur auf Go-Kommentare
- `docs/user/benutzerhandbuch.md` vollständig gelesen (nicht nur der
  Diff-Ausschnitt) — insbesondere die Parallel-Sektionen „Tabelle
  aktivieren", „Consumer registrieren", „Zugriff und Rollen" (Anker
  `#zugriff-und-rollen`) und die bestehende Änderungshistorie (Zeilen
  1.0–1.6) als Formatvorbild
- `git show 294d165` vollständig
- `docs/reviews/review-slice-039.md` (Format-Vorlage)

---

## Findings

### F-1 — Voraussetzungen der neuen Abschnitte unvollständig gegenüber der analogen „Tabelle aktivieren"-Sektion

- `kategorie`: MEDIUM
- `quelle`: Maintainability (sachliche Vollständigkeit gegenüber der
  Parallel-Sektion; `internal/application/usecase/enable/service.go`
  Zeile 26–31: „die Existenz der physischen Tabelle ist die Vorbedingung
  (`LH-FA-CFG-001` Negative: expliziter Fehlerpfad)")
- `pfad`: `docs/user/benutzerhandbuch.md:181-231` (Abschnitte „Tabelle
  live aktivieren"/„Tabelle deaktivieren")
- `befund`: Die bestehende Sektion „Tabelle aktivieren" nennt als
  Voraussetzung explizit „die physische Tabelle existiert" **und**
  „`REPLICA IDENTITY` ist wie benötigt gesetzt". Beide neuen Abschnitte
  rufen laut `slice-037` denselben, unveränderten
  `EnableTableUseCase`/`DisableTableUseCase` auf (`slice-037` §1: „ruft
  den bestehenden `EnableTableUseCase` unverändert auf, ändert seine
  interne Prüftiefe nicht") — die Tabellenexistenz-Prüfung (mit
  explizitem `failed`-Fehlerpfad) und die `REPLICA IDENTITY`-Bedingung
  gelten für den SQL-Weg identisch, stehen aber in der
  „Voraussetzung"-Zeile der neuen Abschnitte nicht. Ein Betreiber, der
  `cdc.enable_table` auf eine nicht existierende Tabelle oder ohne
  passende `REPLICA IDENTITY` aufruft, findet die Fehlerquelle nur über
  den `error_message`-Text im Antrags-Datensatz, nicht über die
  Dokumentation selbst.
- `verifizierbar`: ja — Code-Vergleich (`enable/service.go`,
  `disable/service.go`) gegen die Voraussetzungs-Zeilen beider
  Abschnitte; Vergleich mit der Parallel-Sektion „Tabelle aktivieren"
- `klasse`: „Voraussetzungs-Liste unvollständig ggü. Parallel-Sektion"
  (erstes Auftreten dieser Klasse)

### F-2 — Formatierungsinkonsistenz der Status-Poll-Anweisung zwischen den beiden neuen, laut DoD „analogen" Abschnitten

- `kategorie`: LOW
- `quelle`: Maintainability / Plan-Konformität (`slice-042` §2, DoD-Item
  2: „Neuer Abschnitt „Tabelle deaktivieren" (**analoges Format**, ...)")
- `pfad`: `docs/user/benutzerhandbuch.md:198-202` (Enable) vs.
  `docs/user/benutzerhandbuch.md:225-227` (Disable)
- `befund`: Im Enable-Abschnitt steht die Status-Poll-Abfrage
  (`SELECT status, error_message FROM cdc.administration_request ...`)
  in einem eigenen ```sql```-Codeblock unter „Ergebnis". Im
  Disable-Abschnitt ist dieselbe Abfrage stattdessen als Inline-Code in
  einen Fließtext-Satz gequetscht, der zudem zwischen dem
  „Vorgehen"-Codeblock und „Ergebnis" schwebt, statt einem der beiden
  Label eindeutig zugeordnet zu sein. Die vom DoD geforderte
  Formatgleichheit ist an dieser Stelle nicht gegeben.
- `verifizierbar`: ja (visueller Diff-Vergleich der beiden Abschnitte)
- `klasse`: „Format-Abweichung zwischen als analog deklarierten
  Abschnitten"

### F-3 — Änderungshistorie-Zeile 1.7 bricht das bisherige Zitierschema durch eingeschobene Prosa

- `kategorie`: LOW
- `quelle`: Maintainability (Formatkonsistenz der Änderungshistorie-
  Tabelle; **kein** Verstoß gegen `AGENTS.md` §3.7 — die
  Änderungshistorie ist die etablierte Ausnahme für Slice-Zitate)
- `pfad`: `docs/user/benutzerhandbuch.md:614`
- `befund`: Die Zeilen 1.1–1.6 zitieren Slices ausnahmslos als reine,
  kommagetrennte Kennungs-Liste in der Klammer (z. B. „`slice-038`",
  „`slice-041`"). Zeile 1.7 schiebt stattdessen Prosa in die Klammer ein
  („slice-036/slice-037, nachgetragen mit slice-042") — inhaltlich
  nachvollziehbar und nicht falsch, aber ein Formatbruch gegenüber dem
  bisher durchgehaltenen Muster der Tabelle.
- `verifizierbar`: ja (Zeilenvergleich 1.1–1.7 in der Änderungshistorie)
- `klasse`: „Changelog-Zitierschema-Formatbruch"

## Negativbefunde

- geprüft, ohne Befund: **Sachliche Korrektheit des Antrags-Flusses.**
  Beide neuen Abschnitte beschreiben den Aufruf korrekt als asynchron
  (Status `pending` → `applied`/`failed`, Tabelle „damit noch nicht
  aktiv"/Erfassung endet erst „nach `applied`"); kein „sofort
  aktiv/inaktiv"-Framing. Deckt sich mit `slice-037`s realem
  Live-Reload-Beleg (Goroutine verarbeitet den Antrag, `LISTEN`/`NOTIFY`
  mit Fallback-Poll, kein Neustart).
- geprüft, ohne Befund: **Funktionssignaturen/Parameter-Reihenfolge.**
  `SELECT cdc.enable_table('<source_id>', '<schema>', '<tabelle>')` und
  die spiegelbildliche `disable_table`-Zeile stimmen exakt in Reihenfolge
  und Anzahl der Parameter mit `cdc.enable_table(p_source_id text,
  p_schema_name text, p_table_name text)` /
  `cdc.disable_table(p_source_id text, p_schema_name text, p_table_name
  text)` aus `tools/schema/nacharbeit-administration.sql` überein.
- geprüft, ohne Befund: **`cdc_admin`-Voraussetzung korrekt benannt.**
  Beide neuen Abschnitte nennen „`cdc_admin`-Mitgliedschaft, verbunden
  über `CDC_ADMIN_DSN`" mit Verweis auf `[Zugriff und
  Rollen](#zugriff-und-rollen)" — wortgleiches Muster wie die bestehende
  „Consumer registrieren"-Sektion; der Anker löst korrekt auf
  (Überschrift „### Zugriff und Rollen" existiert, wird bereits an
  anderer Stelle im Dokument mit demselben Anker referenziert).
- geprüft, ohne Befund: **Keine Chronik-Sprache/Slice-Referenz in der
  operativen Prosa außerhalb der Änderungshistorie-Tabelle.** `grep -n
  "slice-\|welle-\|Welle"` über die neuen Abschnittszeilen
  (`docs/user/benutzerhandbuch.md:181-231`) liefert keinen Treffer;
  weder verworfene Alternativen („ohne X wäre …") noch abwesenter Text
  („früher stand hier …") noch abgebrochene Sätze.
- geprüft, ohne Befund: **Positionierung der neuen Abschnitte.** Direkt
  nach „Tabelle aktivieren" und vor „Aktivierte Tabellen auflisten"
  eingefügt — konsistent mit dem bestehenden Aufbau (Registrieren →
  Aktivierungswege → Auflisten → Consumer-Verwaltung).
- geprüft, ohne Befund: **§1 Out-of-Scope-Disziplin korrekt.** Alle drei
  Ausschlusspunkte tragen eine Begründung und ordnen sich den vier
  Klassen zu: „Code-/Verhaltensänderung" (Bestand bleibt stehen,
  `slice-036`/`037` bereits `done/`), „Consumer-Verwaltung über SQL"
  (verweist korrekt auf `welle-12` §6 als bereits bestehenden
  Ausschluss), „Pflichtenheft/Architektur-Änderungen" (Bestand bleibt
  stehen, bereits vollständig).
- geprüft, ohne Befund: **Traceability.** Commit-Betreff nennt
  `LH-FA-ADM-001`, `LH-FA-CFG-002`, `ADR-0050` — mindestens eine
  `LH-*`/`ADR-*`-Kennung erfüllt, keine `SPEC-*`/`ARC-*`-Struktur-ID im
  Betreff.
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung.** Einzige berührte
  Sub-Area `*`/`PGC`, korrekt GF; Beobachtungs-Register-Sichtung
  dokumentiert (kein spezifischer Treffer für Benutzerdokumentation).
- geprüft, ohne Befund: **DoD-Häkchen im Diff plausibel.** Nur die in
  diesem Commit tatsächlich erledigten Punkte (1–4, 6) sind gesetzt;
  Review-Häkchen, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgang
  (§6) und die drei Paarungen bleiben korrekt offen — kein
  Voreil-Häkchen für noch ausstehende Rollenwechsel.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Voraussetzungs-Liste unvollständig
ggü. Parallel-Sektion" (1×, erstes Auftreten) · „Format-Abweichung
zwischen als analog deklarierten Abschnitten" (1×, erstes Auftreten) ·
„Changelog-Zitierschema-Formatbruch" (1×, erstes Auftreten) — kein
Steering-Loop-Eintrag fällig, keine Klasse erreicht 3×.

## Verdikt

**Merge-blockierend:** nein im Sinne eines Rollback — der Commit ist
bereits gepusht, das einzige MEDIUM-Finding (F-1) trägt **keinen
Rollen-Widerspruch**: Es ist eine sachliche Vollständigkeits-Beobachtung
gegen die eigene, in `slice-037` real belegte Implementierung, keine
bestrittene Implementer-Aussage. Die Architect-Sequenz aus Modul 8 greift
hier nicht (kein HIGH, kein Widerspruch, keine dritte Wiederholung
derselben Klasse). Erwartete Reaktion: Der Implementer ergänzt die
fehlenden Voraussetzungen (Tabellenexistenz, `REPLICA IDENTITY`) in
beiden neuen Abschnitten und vereinheitlicht das Format der
Status-Poll-Anweisung, bevor der Slice nach `done/` geht — oder begründet
im Plan-Nachzug, warum die Lücke bewusst hingenommen wird.

**Zur zentralen Prüffrage dieses Laufs (sachliche Korrektheit gegen den
realen asynchronen Antrags-Fluss):** Bestätigt — kein irreführendes
„sofort aktiv"-Framing, Funktionssignaturen exakt, `cdc_admin`-
Voraussetzung korrekt benannt, keine Chronik-Sprache außerhalb der
etablierten Änderungshistorie-Ausnahme.

**Übergabe:** Ein MEDIUM-Finding (Voraussetzungs-Vollständigkeit), zwei
LOW-Hinweise (Format). Der Implementer entscheidet über Annahme oder
Begründung; keine Architect-Sequenz erforderlich. Dieser Report ist ein
Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen — die
Summary-Zeile speist bei Bedarf den Closure-Eintrag (Modul 5). Er ersetzt
keine Verifikation — DoD-Konformität (inkl. der übrigen offenen
DoD-Punkte: Review-Häkchen selbst, Closure-Notiz, Register, Risiko-
Ausgang, drei Paarungen) prüft der Verifier separat.
