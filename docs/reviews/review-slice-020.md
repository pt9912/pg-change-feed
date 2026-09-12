# Review-Report: slice-020 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-020`) + `ADR-0023`
(Fehlerklassifikation, Accepted) + `SPEC-008` (Referenz-Tabelle der sieben
Fehlerklassen) + `AGENTS.md` §3 (Hard Rules, insbesondere §3.7).

**Gegenstand:** Commits `bbdd5f4` (Handbuch-Tabelle auf sieben Klassen
vervollständigt, DoD-Checkboxen nachgezogen) und `051096e`
(Planner-Korrektur: Entfernung „· seit slice-020" aus der Änderungshistorie)
— beide auf `docs/user/benutzerhandbuch.md` §6/§9 bzw. dem Slice-Plan selbst.

**Skill:** `.harness/skills/reviewer.md` @ HEAD zum Review-Zeitpunkt
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-020-handbuch-fehlerklassen-vollstaendig.md`
  §1–§8, vollständig
- `spec/pflichtenheft.md` `SPEC-008` (Referenz-Tabelle der sieben Fehlerklassen)
- `docs/plan/adr/0023-fehlerklassifikation.md` (Accepted)
- `internal/domain/model/errorstate.go` (die sieben `ErrorClass`-Konstanten)
- `internal/bootstrap/wiring.go` (`classifyRunError`, real konstruierte Klassen)
- `LH-FA-ADM-003` (`spec/lastenheft.md`)
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.3, §3.7)
- `harness/conventions.md` (MR-000/MR-001)
- Vorgänger-Finding-Klasse `BEO-PGC/dod-checkbox-nachzug` (verkörpert seit
  welle-5, `.claude/commands/implement-slice.md` Schritt 18) —
  Präzedenzfälle `slice-019`, `slice-021` für die Annotations-Praxis bei
  „entfällt"-Checkboxen

---

## Findings

### F-1 — `replication`-Zeile der Handbuch-Tabelle beschreibt eine andere Aktion als `SPEC-008`

- `kategorie`: MEDIUM
- `quelle`: `SPEC-008` (Maintainability — unklare Fehlerbehandlung am Rand
  des Spec-Bereichs)
- `pfad`: `docs/user/benutzerhandbuch.md:337` (Zeile unverändert von diesem
  Diff, aber Teil der jetzt „vollständigen" Sieben-Klassen-Tabelle, die §1
  dieses Slice als Ziel benennt)
- `befund`: `SPEC-008` beschreibt die Aktion für `replication` als
  „Überwachung über Schwellen (§5, WAL-Rückstand); kontrollierte
  Fortsetzung" — also kein Prozessabbruch. Das Handbuch dagegen benennt
  für dieselbe Klasse „Sichtbarer Fehler" und ordnet sie über den
  Sammelsatz direkt nach der Tabelle („Ein Fehler jeder Klasse beendet den
  Container-Prozess mit Ausgang 1") demselben Fatal-Pfad zu wie alle
  anderen sechs Klassen. Der Code bestätigt das Handbuch
  (`classifyRunError` mappt `receive.ErrReplication` u. a. auf
  `ErrorClassReplication`, die über `reportFault`/Prozessabbruch behandelt
  wird) — die Diskrepanz liegt in `SPEC-008` selbst, nicht im Code oder im
  Handbuch. Diese Zeile war bereits vor diesem Diff so formuliert (nur
  Positionswechsel in der Tabelle durch die Neusortierung nach
  `SPEC-008`-Reihenfolge), gehört aber zur Lieferung, die §1 dieses Slice
  als „bildet die dortige geschlossene Menge der sieben Fehlerklassen ab"
  bezeichnet — die Abbildung ist für genau diese eine Zeile nicht deckungsgleich.
- `verifizierbar`: nein — kein Gate vergleicht Pflichtenheft-Aktionsspalten
  sinngemäß gegen Betreiberdoku-Prosa; nur Leseverständnis.
- `klasse`: „Spec-Doku-Divergenz: Fehlerklassen-Aktion (replication)"

### F-2 — DoD-Checkbox „Reconciliation-Register … entfällt" ohne die etablierte Audit-Annotation gesetzt

- `kategorie`: LOW
- `quelle`: Maintainability / embodied Konvention
  (`BEO-PGC/dod-checkbox-nachzug`, verkörpert in
  `.claude/commands/implement-slice.md` Schritt 18)
- `pfad`: `docs/plan/planning/in-progress/slice-020-handbuch-fehlerklassen-vollstaendig.md:88`
- `befund`: Der Punkt wurde von `[ ]` auf `[x]` gesetzt, ohne die in den
  beiden unmittelbaren Vorgänger-Closures (`slice-019`, `slice-021`) für
  denselben „entfällt"-Fall verwendete Inline-Annotation
  („Geprüft: Datei existiert nicht (GF-Repo) — entfällt"). Der
  zugrundeliegende Befund selbst ist korrekt — `docs/plan/planning/reconciliation.md`
  existiert in diesem Repo nicht (reines GF-Repo, `harness/conventions.md`
  Modus-Deklaration `*` = Greenfield) —, aber ohne die Annotation ist das
  nur durch eigene Dateisystem-Prüfung nachvollziehbar, nicht aus dem
  Text selbst.
- `verifizierbar`: ja — `find docs/plan/planning -iname reconciliation.md`
  bestätigt die Abwesenheit der Datei; die Checkbox-Aussage ist damit in
  der Sache korrekt.
- `klasse`: „DoD-Checkbox ohne Audit-Annotation bei entfällt-Fall"

### F-3 — Erklärender Absatz zu `transient`/`permission` ist ohne Sensor gegen künftige Adapter-Änderungen ungeschützt

- `kategorie`: INFO
- `quelle`: Maintainability — Hinweis ohne erwartete Aktion
- `pfad`: `docs/user/benutzerhandbuch.md:340-344`
- `befund`: Der neue Absatz behauptet, `transient`/`permission` würden
  „aktuell" von keinem Adapter konstruiert. Dieser Fakt-Stand ist heute
  korrekt (per `grep` bestätigt), aber nichts im Repo würde eine künftige
  Adapter-Änderung, die eine der beiden Klassen produziert, automatisch
  gegen diesen Handbuch-Absatz spiegeln. Der Slice benennt genau das
  bereits selbst als Out-of-Scope („ein Sensor, der Handbuch-Inhalt gegen
  den Fehlerklassen-Code prüft — anderer Vorgang"); dieses Finding
  bestätigt nur, dass die Lücke real und heute noch offen ist.
- `verifizierbar`: nein — kein Gate vorhanden, Slice deklariert das bewusst.
- `klasse`: „Doku-Code-Drift-Risiko ohne Sensor (Fehlerklassen-Konstruktion)"

## Negativbefunde

- geprüft, ohne Befund: Vollständigkeit — alle sieben Klassen aus
  `ADR-0023`/`SPEC-008` (`transient`, `configuration`, `permission`,
  `schema`, `storage`, `replication`, `internal`) sind in
  `docs/user/benutzerhandbuch.md` §6 gelistet, in der `SPEC-008`-Reihenfolge.
- geprüft, ohne Befund: `internal`-Fallback-Behauptung — `classifyRunError`
  (`internal/bootstrap/wiring.go:400-421`) hat exakt einen `default:`-Zweig,
  der `model.ErrorClassInternal` zurückgibt; das Handbuch benennt das korrekt.
- geprüft, ohne Befund: `transient`/`permission` „von keinem Adapter aktuell
  konstruiert" — `grep -rn "ErrorClassTransient\|ErrorClassPermission" --include="*.go" .`
  findet nur Vorkommen in `errorstate.go` (Konstanten-Deklaration),
  `errorstate_test.go` (Aufzählungstest) und `heartbeat_test.go` (Test ruft
  den reinen Persistenz-Durchreicher `PostgresHeartbeatAdapter.Fault`
  direkt mit einer vorgegebenen Klasse auf — `Fault` klassifiziert nicht,
  es speichert nur die übergebene Klasse). Kein Adapter/keine Verdrahtung
  entscheidet real auf `transient` oder `permission`.
- geprüft, ohne Befund: Bedingungs-/Aktionsspalten für `configuration`,
  `schema`, `storage` — sinngemäß konsistent mit `SPEC-008` (unverändert
  von diesem Diff).
- geprüft, ohne Befund: Hard Rule 3.7 / Chronik-Sprache — außer der von
  `051096e` entfernten „· seit slice-020" keine weiteren Slice-Referenzen
  oder Prozess-Chronik im geänderten Text; `grep` über das gesamte
  Dokument findet keine weitere `slice-`-Referenz. Das einzige verbleibende
  „jetzt" im Dokument (Zeile 157) ist vorbestehender, von diesem Diff nicht
  berührter Text.
- geprüft, ohne Befund: §1-Abgrenzung — `git show bbdd5f4/051096e --stat`
  bestätigt: nur `docs/user/benutzerhandbuch.md` und die Slice-Plan-Datei
  geändert; kein Code-Diff an `classifyRunError`/`ErrorClass`, kein neuer
  Sensor, kein Makefile-Target berührt.
- geprüft, ohne Befund: AGENTS.md §3.3 (`git mv` + Inhaltsänderung) — beide
  Commits sind reine Inhaltsänderungen ohne Datei-Verschiebung; nicht
  einschlägig.
- geprüft, ohne Befund: Traceability — beide Commit-Betreffs nennen
  `LH-FA-ADM-003` (und `bbdd5f4` zusätzlich `ADR-0023`), kein `SPEC-*`/`ARC-*`
  im Betreff.
- geprüft, ohne Befund: `make gates` — selbst ausgeführt, grün
  (`baseline-verify`: 54 Dateien OK; `docs-check`: 221 Dateien, 0 Befunde;
  `commit-traceability`: OK über die letzten 5 Commits; `a-check`: 0 Befunde).
- geprüft, ohne Befund: Docker-only / Suppression / ADR-Immutabilität —
  vom Diff nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Spec-Doku-Divergenz: Fehlerklassen-Aktion
(replication) · DoD-Checkbox ohne Audit-Annotation bei entfällt-Fall ·
Doku-Code-Drift-Risiko ohne Sensor (Fehlerklassen-Konstruktion)

## Verdikt

**Merge-blockierend:** nein — kein HIGH-Finding. F-1 ist ein pre-existing
Spec-Doku-Widerspruch, den dieser Diff weder eingeführt noch verschlechtert
hat und den §1 dieses Slice nicht in seinem Ausschluss-Katalog nennt (das
Ziel war „ergänzen", nicht „alle Zeilen gegen `SPEC-008` audit-korrigieren");
er gehört als Beobachtung oder Folge-Slice getragen, nicht als
Merge-Blocker dieses reinen Ergänzungs-Diffs. F-2 ist eine
Konsistenz-Abweichung von einer erst zweimal beobachteten Annotations-Praxis,
ohne dass die zugrundeliegende Checkbox-Aussage falsch wäre.

**Übergabe:** Findings gehen an den Implementer (pt9912-Rolle) zur
Aufnahme in die Closure-Notiz §7 dieses Slice bzw. als Risiko-Ausgang;
F-1 sollte vor Closure entweder als „weiter offen" ins
Beobachtungs-Register wandern (Berührung der Sub-Area `*`/`PGC`) oder mit
einem Folge-Slice adressiert werden, da es ein Ausgang für §6 dieses Slice
werden könnte, sobald der Implementer die drei Risiken-Ausgänge bei
Closure zuweist. Dieser Report ist Lauf-Beleg; DoD-/Spec-Konformität prüft
der Verifier separat (Modul 11).
