# Architect-Review slice-013 — Verdikt zu F-1 (Scope-Reduktion und Übergangs-Reihenfolge)

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-10.
**Eingang:** [`docs/reviews/review-slice-013.md`](../../reviews/review-slice-013.md)
F-1 (HIGH, Rollen-Bezug) · Slice-Plan
[`docs/plan/planning/in-progress/slice-013-fehlerzustaende-cdc-abstand.md`](../planning/in-progress/slice-013-fehlerzustaende-cdc-abstand.md)
§1/§2/§4/§6 · Commit-Sequenz `fcf442d..HEAD` (`git log --oneline`, s.
Beleg-Anker) · Präzedenzfall
[`architect-review-slice-011.md`](architect-review-slice-011.md) ·
[`done/welle-2-results.md`](../planning/done/welle-2-results.md) §Steering-
Loop-Einträge, Eintrag „Zwei gemischte Commits" (zweiter Präzedenzfall, für
Frage 2 einschlägiger als slice-011) ·
[`ADR-0023`](0023-fehlerklassifikation.md) (Accepted, permanent) ·
[`ADR-0032`](0032-postgresql-adapterdetail.md) (Accepted, permanent) ·
[`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md) (Accepted,
Supersedes ADR-0018) · Baseline-Regelwerk
`modul-05-planning-harness.md` §Lifecycle als State Machine ·
`modul-08-agentenrollen.md` §Kernidee, §Konflikt-Pfad als Rollen-Sequenz.
**Ausgang:** Übergabe-Artefakt an Planner — **kein Folge-ADR nötig, kein
Carveout**; Verdikt zu zwei getrennten Fragen (Inhalt der Scope-Reduktion /
Prozess der Übergangs-Reihenfolge) und eine geschärfte Prozess-Regel als
Textvorschlag (Umsetzung ist Planner-Sache).
**Harte Regel eingehalten:** keine Accepted-ADR in-place geändert; keine
Folge-ADR-Datei angelegt; kein ADR-Index-Eintrag (dieses Dokument ist kein
ADR); der Slice-Plan wird von diesem Lauf **nicht** editiert — Closure-Notiz
und Regel-Umsetzung sind Planner-Sache im nächsten Zug (Modul 8
§Konflikt-Pfad, Sequenz).

---

## Frage 1 — Ist die Scope-Reduktion selbst sachlich richtig?

**Kurzform:** Ja, in allen drei geprüften Teilfragen. Die Reduktion trägt
inhaltlich; der Reviewer-Negativbefund dazu ist bestätigt, nicht nur
nachvollzogen.

### Prüfung: Vier-Schichten-Einschätzung für das reale `cdc_capture_lag`

Eigenständig am Code geprüft, nicht aus der Plan-Prosa übernommen:

- **Replication-Decoder** (`internal/adapters/driving/replication/decode/decode.go:45-48`):
  `Commit struct { CommitLSN uint64; EndLSN uint64 }` — das Feld, das
  `pglogrepl.CommitMessage` bereits liefert (`CommitTime time.Time`), wird
  beim Dekodieren (Zeilen 134-138) **nicht** übernommen. Der Plan-Befund
  „keine neue Wire-Verbindung nötig, aber Zeitstempel-Durchreichen fehlt"
  ist damit am Code verifiziert, nicht behauptet.
- **Domain** (`internal/domain/model/transaction.go`): `ChangeTransaction`
  trägt `commitPosition SourcePosition` (`CommitPosition()`/`Commit(...)`),
  keinen Zeitstempel — `SourcePosition` (`position.go`) ist rein
  LSN-basiert. Ein reales `cdc_capture_lag` bräuchte ein zusätzliches
  Zeit-Feld an genau dieser Typ-Grenze.
- **Application/Ports** (`internal/application/port/inbound/capture.go`,
  `internal/application/port/outbound/changestore.go`): `CaptureCommand`
  trägt `Transaction`, `ChangeStorePort` erwartet die Domain-Typen ohne
  Zeit-Parameter — der Transport-Vertrag müsste erweitert werden.
- **Store-Adapter** (`internal/adapters/driven/postgresstorage/store.go:95-99`,
  `queries/queries.go:184`): `InsertTransaction` wird mit
  `TransactionID, SourceID, CommitPosition` aufgerufen — kein Zeit-Parameter;
  der Query-Kommentar hält explizit fest „der Aufrufer übergibt keine Uhr"
  (die DB setzt `committed_at` selbst über `current_timestamp`/`DEFAULT
  now()`).

Vier unabhängige Schichten, vier unabhängige Änderungspunkte — die
Plan-Einschätzung ist **technisch zutreffend**, nicht nur plausibel
formuliert. Der Reviewer-Negativbefund („die Rückführung war real
begründet, keine vorgeschobene Vereinfachung") wird hiermit von
unabhängiger Seite bestätigt.

### Prüfung: Kohärenz des reduzierten Scope

Der reduzierte Scope liefert zwei voneinander unabhängige, beide real
gelieferte und getestete Liefer-Punkte (§2 DoD): Fehlerzustands-Sichtbarkeit
(`ErrorClass`, Heartbeat-`Fault`, Composition-Root-Klassifikation) und die
`cdc_capture_lag`-Näherung (Persistenz-Zeit-Proxy, als Grenze markiert). Keiner
der beiden wartet auf das ausgeschlossene reale `cdc_capture_lag` — die
Näherung ist als eigenständiger, ehrlich begrenzter Wert nutzbar (ein
Monitoring-System sieht eine Pipeline-Frische-Kennzahl, keine kaputte oder
platzhalterhafte). Kein Zombie-Slice: Beide Teile sind für sich
liefer- und nutzwert-tragend, das reale Maß ist eine echte Erweiterung, keine
Bedingung für die Nützlichkeit des jetzt Gelieferten.

### Prüfung: Folge-Slice-Ausschluss (§1) und Register-Eintrag

`docs/plan/planning/observations/BEO-PGC/cdc-capture-lag-real/` existiert,
`state.md` trägt den Ausgang „weiter offen → Folge-Slice" mit
`evidence/slice-013.md` (1×) — Inhalt und Sub-Area-Zuordnung
(`Observability/Replication-Adapter`) decken sich mit §1/§6 des Plans. Der
Ausschluss zitiert den Registerpfad korrekt (keine Neu-Formulierung, keine
zweite Kennung). Eine Sache bleibt offen und ist **kein** Defekt dieses
Plan-Standes, sondern ein Closure-Vorbehalt: §6 trägt den Ausgang
„eingetreten (Träger: Folge-Slice … Kennung folgt bei der nächsten
Eröffnung)" — Modul 5 verlangt für „eingetreten" beim **Übergang nach
`done/`** eine Folge-Slice-**ID**, nicht nur die Klasse. Solange der Plan in
`in-progress/` liegt, ist „Kennung folgt" zulässig; vor dem `git mv` nach
`done/` muss entweder die Slice-ID existieren oder der Ausgang ändert sich.
Das ist ein Hinweis an den Planner für die Closure, kein Finding gegen den
jetzigen Stand.

**Ergebnis Frage 1:** Scope-Reduktion inhaltlich korrekt, reduzierter Scope
kohärent und eigenständig wertvoll, Ausschluss und Register-Eintrag sauber
gepaart.

---

## Frage 2 — Reicht ein Steering-Loop-Lerneintrag, oder braucht es einen Carveout?

**Kurzform:** Ein dokumentierter Steering-Loop-Lerneintrag (plus geschärfte
Regel) ist das richtige Werkzeug. Ein Carveout ist hier ein **Werkzeug-
Fehlgriff** — nicht die zu milde, sondern die falsche Kategorie.

### Rekonstruktion der Zeitlinie (unabhängig gegen `git log` geprüft)

```
15:15:44  426e1ed  chore(plan): next/ — Rückführung (reiner Move)
15:15:56  0588958  docs(plan): Rückführung dokumentiert, Plan-Nachzug
15:16:03  e71f1a0  feat(domain): ErrorClass
15:16:11  91c8c0b  feat(port,adapter): Heartbeat Fehlerzustand
15:16:18  79dbe3f  feat(schema): error_class + cdc_capture_lag-Näherung
15:16:25  6b05247  feat(bootstrap): Fehlerzustand klassifiziert/gemeldet
15:17:49  1f8ddf3  chore(image): Digest-Beleg
——— slice-013 liegt bis hier ohne Unterbrechung in next/ ———
15:20:42  467b78a  docs(planning): in-progress (Planner-Korrektur)
15:21:04  8d67f5   docs(planning): Scope reduziert (§1/§2/§5/§6)
15:22:04  7e72c67  docs(planning): Register-Eintrag
```

Fünf Produktions-Commits (inkl. Image-Digest) liegen real und vollständig
zwischen dem Rückführungs-Move und dem Rückkehr-Move — bestätigt, nicht nur
aus dem Review übernommen. Die Reihenfolge ist damit tatsächlich invertiert
gegenüber Modul 5: der Anspruchs-Commit (`git mv` nach `in-progress/`)
kommt **nach**, nicht **vor** der Arbeit.

### Werkzeug-Wahl-Prüfung (Modul 7 §Werkzeug-Wahl bei Diskrepanz)

Der Trichter aus Modul 7 ist für **andauernde Diskrepanzen** gebaut (Code
weicht von Doku ab, ein Gate ist rot, eine Regel ist gesenkt) — mit einem
Geltungsbereich, der sich in der Zukunft auflösen lässt. Hier passt keine
der drei Kategorien:

- **Carveout scheidet aus.** Ein Carveout trägt zwingend ein „betroffenes
  Gate" (Pflicht-Header-Feld) und einen Auflösungs-Trigger, der ein *noch
  bestehendes* rotes Gate in einen grünen Zustand überführt. Hier ist kein
  Gate rot — der Reviewer selbst hält fest: „kein Gate erzwingt
  Transition-Reihenfolge". Die Verzeichnis-Position ist inzwischen korrekt
  (`in-progress/`, seit `467b78a`); es gibt nichts mehr „aufzulösen". Ein
  Carveout für ein bereits abgeschlossenes, nicht wiederholbares
  Vergangenheits-Ereignis hätte kein Gate, das er referenziert, und keinen
  Zustand, den er in der Zukunft ändert — er wäre Form ohne Funktion.
- **BF-Sub-Area-Markierung scheidet aus.** Kein Cluster im selben
  Geltungsbereich (die Finding-Klasse steht im Report als „1. Auftreten"),
  kein systemisches Code-vor-Doku-Muster — eine Sub-Area-Modus-Frage stellt
  sich hier nicht.
- **ADR (permanent) scheidet aus.** Es gibt keine Architekturentscheidung,
  die abgelöst oder bestätigt werden müsste — die verletzte Regel ist eine
  Prozess-Regel des Baseline-Regelwerks (Modul 5), keine `ADR-*`. Der
  Konflikt-Pfad aus Modul 8 (drei ADR-Verdikte) ist damit nicht wörtlich
  anwendbar; seine **Mechanik** — unabhängiges Architect-Verdikt als
  Übergabe-Artefakt statt Selbstbestätigung derselben Kette — ist es sehr
  wohl, und genau das liefert dieses Dokument.

### Der einschlägige Präzedenzfall ist nicht slice-011, sondern welle-2

`done/welle-2-results.md` §Steering-Loop-Einträge führt „Zwei gemischte
Commits" (`13abe88`, `e3e36f2`): ein struktur-gleicher Fall — ein
Rollen-/Reihenfolge-Verstoß, der bereits geschehen und nicht mehr rückgängig
zu machen war (Implementer-Content vermischt mit Review-/Plan-Artefakten,
eine Review-Basis verdrängt). Der Ausgang dort war **weder** Carveout
**noch** ADR: Offenlegung in den Closure-Notizen von slice-004/005 plus eine
benannte Konsequenz („Rollen-Läufe strikt sequenzieren"). Kein Gate wurde
gebaut, kein Carveout eröffnet — die Disziplinierung lief über die
verkörperte Regel, nicht über einen Ausnahme-Mechanismus. Dieser Fall hier
hat dieselbe Form (einmaliges, bereits abgeschlossenes Reihenfolge-Ereignis,
keine laufende Ausnahme) und verdient denselben Träger.

### Verdikt

**Kein Carveout.** Ein **dokumentierter Steering-Loop-Lerneintrag** in der
Closure-Notiz (§7) — die Abweichung offen benannt, mit Zeitlinie wie oben,
keine verschleierte Historie — **plus eine geschärfte Prozess-Regel** ist der
richtige, verhältnismäßige Träger. Die Übergangs-Reihenfolge-Verletzung
selbst bekommt den Ausgang **eingetreten, mit Konsequenz** (keiner der drei
Risiko-Ausgänge aus Modul 5 passt wörtlich, weil es kein vorab benanntes
Risiko in §6 war — sondern ein Finding aus dem Review; dessen Ausgang läuft
über die Finding-Klassen-Route in §7, nicht über §6).

**Zusätzlich mit diesem Dokument geschlossen:** Die zweite Hälfte von F-1 —
„Rollenwechsel ohne Übergabe-Artefakt" — ist mit diesem unabhängigen
Architect-Verdikt bedient. Anders als die „Planner-Korrektur" (`467b78a`,
`8d67f52`, derselbe `Claude-Session`-Trailer wie die Implementer-Commits,
fünf Minuten Abstand) ist dieses Dokument ein eigenständiger Rollen-Zug mit
eigenem Eingabe-Kontext (Review-Report + Plan + Commit-Historie + Code,
unabhängig von der Implementer-/Planner-Selbsteinschätzung geprüft) — die
Form, die der Reviewer für den Closure-Ausgang verlangt hat.

### Geschärfte Prozess-Regel (Textvorschlag — Umsetzung ist Planner-Sache)

> Eine Implementer-Rückführung (`in-progress → next`), die im selben Lauf
> trotzdem Code liefert (Plan-Nachzug mit fortgesetzter Arbeit statt echtem
> Slice-Stopp), braucht **vor** dem Rückkehr-Commit nach `in-progress/`
> einen unabhängigen Architect-Verdikt-Zug — nicht eine Planner-Bestätigung
> im selben Kontextfenster/derselben Session wie die Implementer-
> Entscheidung. Bis das Verdikt als eigenständiges Artefakt vorliegt, bleibt
> der Plan in `next/` liegen, auch wenn der unabhängige Teil bereits real
> gemergt ist; der Rückkehr-Commit trägt im Betreff oder Body einen Verweis
> auf das Architect-Verdikt-Artefakt. Zusätzlich: Der Rückkehr-Move-Commit
> landet vor jedem weiteren Produktions-Commit dieses fortgesetzten Teils —
> Code-Auslieferung wartet auf den sichtbaren Anspruchs-Commit, nicht
> umgekehrt (Modul 5 §Lifecycle als State Machine gilt für **jeden**
> `next → in-progress`-Übergang, auch den zweiten innerhalb desselben
> Slice-Zugs).

Empfohlener Zielort für die Verkörperung (Planner-Entscheidung, nicht
Architect-Festlegung): `.claude/commands/implement-slice.md`
(Implementer-Workflow-Schritt, analog zum bestehenden Plan-Nachzug-Schritt
14) plus ein Verweis in `AGENTS.md` §6, falls der Planner das für
repo-weit sichtbar genug hält.

---

## Disposition zu F-1 (Gesamt)

| Teil-Finding | Ausgang |
|---|---|
| Vier-Schichten-Einschätzung / Kohärenz des Scope | bestätigt — kein Korrektur-Bedarf am Plan |
| Folge-Slice-Ausschluss + Register-Eintrag | bestätigt — mit Closure-Vorbehalt „Kennung vor `done/`" (Hinweis an Planner) |
| Übergangs-Reihenfolge verletzt | eingetreten, mit Konsequenz — Steering-Loop-Lerneintrag + geschärfte Regel, **kein Carveout** |
| Rollenwechsel ohne Übergabe-Artefakt | mit diesem Dokument bedient — unabhängiger Architect-Zug liegt jetzt vor |

F-1 ist damit **kein Closure-Blocker mehr**, sofern die Planner-Closure-Notiz
(§7) den Steering-Loop-Eintrag tatsächlich schreibt (Zeitlinie + Konsequenz,
Zielort für die Regel-Verkörperung benannt oder als „geplant" mit
Folge-Slice/Folge-Kennung markiert) — Modul 5 verlangt für „done/" einen
Lerneintrag, nicht bloß dessen Ankündigung.

---

## Beleg-Anker (Kurzfassung)

| Aussage | Beleg |
|---|---|
| Vier-Schichten-Einschätzung technisch zutreffend | `decode.go:45-48,134-138`, `transaction.go`, `position.go`, `capture.go`, `changestore.go`, `store.go:95-99`, `queries.go:184` |
| Produktions-Commits landen vor dem Rückkehr-Move | `git log --format='%H %ad %s' fcf442d..HEAD` (Zeitlinie oben) |
| Kein Gate erzwingt Transition-Reihenfolge → Carveout ohne Objekt | `docs/reviews/review-slice-013.md` F-1 `verifizierbar: nein` |
| Präzedenz für „Disclosure statt Carveout/ADR" bei abgeschlossenem Reihenfolge-Verstoß | `docs/plan/planning/done/welle-2-results.md` §Steering-Loop-Einträge, „Zwei gemischte Commits" |
| Register-Eintrag korrekt gepaart | `docs/plan/planning/observations/BEO-PGC/cdc-capture-lag-real/{observation,state}.md` |
