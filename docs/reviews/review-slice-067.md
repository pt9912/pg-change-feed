# Review-Report: slice-067 — 2026-09-14

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Modul 10
§Drei Review-Arten); DoD-/Spec-Konformität ist Verifier-Aufgabe und nicht
Gegenstand dieses Reports.

**Gegenstand:** `slice-067`
(`docs/plan/planning/in-progress/slice-067-assembler-filterung-live-reload.md`),
Diff `541ee01..45c619b` (Elter-Commit `541ee01` ist ein reiner
`next→in-progress`-Move und trägt keinen Inhalt). Zwei Commits: `1943d86`
(Filterung + Schema-Bump-Erhalt + Live-Reload-Verdrahtung), `45c619b`
(Image-Digest). Beide liegen bereits auf `main` (`git branch --contains`),
kein offener PR-Zustand.

**Skill:** `.harness/skills/reviewer.md` @ `45c619b` (Stand zum
Review-Zeitpunkt, unverändert seit der letzten Schärfung 2026-09-13).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-14.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-067` (vollständig, inkl. der drei Plan-Nachzüge und §6)
- [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (vollständig,
  insbesondere Teilfrage 3, Teilfrage 4, Teilfrage 5, §Bestätigung —
  Live-Reload-Konsistenz, §Konsequenzen, §Fitness Function)
- [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  (Antrags-Queue, `tablesMu`-Synchronisation), [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
  (Image-Digest als Lauf-Beleg)
- `AGENTS.md` §3.3/§3.7/§3.9, §5 (Traceability), §6 (Rollenwechsel)
- `harness/conventions.md` (`MR-000` ID-Schema)
- [`LH-FA-CFG-005`](../../spec/lastenheft.md) (Happy Path / Boundary /
  Negative), [`LH-QA-SEC-004`](../../spec/lastenheft.md),
  [`LH-FA-DAT-005`](../../spec/lastenheft.md),
  [`LH-FA-SCH-003`](../../spec/lastenheft.md)/`LH-FA-SCH-004.a`
- `welle-18` §3/§4/§6 (Slice-Zuschnitt, Wellen-Out-of-Scope)
- Vorgänger `slice-066` (in `done/`) samt dem Review zu `slice-066`;
  vorherige Findings am gleichen Modul: das Review zu `slice-060`,
  das Review zu `slice-062`, das Review zu `slice-065`
- `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar`
  (6 Belege, `verkörpert`) — Anlass für den Kommentar-Grep

**Eigene Sensor-Läufe (Exit-Code jeweils ungepiped in einem eigenen Schritt
ermittelt, `AGENTS.md` §3.9):**

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify, d-check (522 Dateien, 0 Befunde), commit-traceability (5 Commits), a-check (0 Befunde), coverage-gate (45.90 % über Schwelle 35 %) |
| `make test` | **0** | vollständige Suite im Race-Container, alle Pakete `ok` |
| `make test-store` | **0** | reale PostgreSQL (`internal/adapters/driven/postgresstorage` 4.1s, `internal/bootstrap` 0.8s — die realen Tests laufen, nicht `skip`) |
| `make test-replication` | **0** | Replication-Stream-Tests gegen reale PostgreSQL (`receive` 4.4s, `bootstrap` 20.7s) |

**Eigene Mutationsläufe (unabhängige Nachprüfung der beiden Erhalt-Punkte und
des Konvergenz-Tests).** Vier Mutationen, jede in einem eigenen Lauf gegen das
gepinnte Toolchain-Image (netzlos, `go test -race`), danach Rücknahme per
`git checkout` — der Arbeitsbaum ist nach der Prüfung leer (`git status
--porcelain`):

| # | Mutation | Erwartung | Ergebnis |
|---|---|---|---|
| A | beide Erhalt-Punkte zurückgedreht (`AddBinding`-Merge entfernt **und** `observeRelation` wieder über `AddBinding` mit vollem `TableBinding`) | `TestConsumeExcludedColumnSurvivesSchemaBump` rot | **rot** (Exit 1, genau dieser Test) |
| B | nur der Merge in `AddBinding` entfernt (`setSchemaVersion` bleibt) | `TestAddBindingKeepsExclusionState` rot, der Bump-Test grün | **rot** (Exit 1, nur dieser Test) |
| C | nur `observeRelation` wieder über `AddBinding` (Merge bleibt) | beide grün | **grün** (Exit 0) |
| F | Ausschlussliste **vor** `classifyRelationColumns` herausgefiltert (bekannte Spaltenform und eingehende Relation) | Konvergenz-Test rot | **rot** (Exit 1) |
| G | Ausschlussliste auf einen Fremdnamen (`other` statt `secret`) gesetzt | Konvergenz-Test weiter grün | **grün** (Exit 0) |

A/B/C bestätigen die Implementer-Angabe: beide Erhalt-Punkte sind einzeln
tragend, keiner ist Dekoration — C allein bleibt grün, weil der Merge in
`AddBinding` den Pfad dann schon hält. F/G kalibrieren F-5 unten.

---

## Findings

### F-1 — Der Ausschlussstand hat keinen dauerhaften Träger: er geht über einen Prozess-Neustart **und** über einen `disable`/`enable`-Zyklus verloren

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  §Teilfrage 1 (Option D, Pro-Spalte: „löst die Neustart-Einschränkung
  strukturell (derselbe Live-Reload-Pfad wie Tabellen-Aktivierung)") ·
  §Bestätigung — Live-Reload-Konsistenz („ohne Prozess-Neustart wirksam,
  exakt wie eine Tabellen-Aktivierung heute schon") · §Konsequenzen (positiv:
  „ein ausgeschlossener Wert erreicht nie die Persistenzschicht") ·
  [`LH-QA-SEC-004`](../../spec/lastenheft.md) (Messmethode über
  `LH-FA-CFG-005`: „ausgeschlossene Spaltenwerte erscheinen nicht in den
  Changes") · Modul 5 §Offene Risiken werden bei Closure aufgelöst
- `pfad`: `internal/bootstrap/wiring.go:321-339` (`activatedTableBindings`
  baut `TableBinding` ohne `ExcludedColumns`), `:1058` (Enable-Zweig →
  `AddBinding`), `:1069` (Disable-Zweig → `RemoveBinding`),
  `internal/adapters/driving/replication/mapper/mapper.go:406-411`
  (`AddBinding`-Merge), `:509-513` (`RemoveBinding` löscht den Eintrag samt
  Ausschlussstand); Plan `slice-067` §6 (drittes, nachgetragenes Risiko)
- `befund`: Der Ausschlussstand lebt ausschließlich in der laufenden
  `Assembler`-Bindung; kein Startpfad liest ihn wieder ein (`applied`-Zeilen
  in `cdc.administration_request` wertet `activatedTableBindings` nicht aus),
  ein Neustart erfasst eine zuvor ausgeschlossene Spalte also wieder — das
  hat der Implementer selbst benannt. **Nicht benannt** ist der zweite, ohne
  jeden Neustart erreichbare Auslöser: ein `cdc.disable_table(t)` löscht die
  Bindung samt Ausschlussstand (`RemoveBinding`), ein darauf folgendes
  `cdc.enable_table(t)` legt sie über `AddBinding` neu an — der Merge greift
  dann ins Leere, weil kein Vorgänger mehr da ist; der `exclude_column`-Antrag
  bleibt `applied`, die Spalte wird wieder erfasst. Beide Auslöser liegen in
  Zeilen, die dieser Diff anfasst (`RemoveBinding`/`AddBinding`), der zweite
  ist also nicht „außerhalb des Slice". **Kein ADR-Verstoß** (Urteil unten):
  der ADR-Wortlaut sagt Live-Reload zu, nicht Dauerhaftigkeit — die einzige
  Stelle, die eine Parität nahelegt („exakt wie eine Tabellen-Aktivierung"),
  ist der Satz, an dem die beiden real auseinanderfallen. Die Entscheidung
  „zulässige Grenze oder Lücke" trägt kein Träger: sie steht weder in
  `welle-18` §6 (Out-of-Scope) noch in `ADR-0059`, und §6 des Slice nennt nur
  den Neustart.
- `verifizierbar`: nein — kein Gate deckt Dauerhaftigkeit oder einen
  `disable`/`enable`-Zyklus ab; belegt durch eigenen Code-Gang über die drei
  genannten Aufrufe (ein Unit-Test über die öffentlichen `ExcludeColumn`,
  `RemoveBinding`, `AddBinding` würde den zweiten Auslöser in fünf Zeilen
  zeigen)
- `klasse`: „Ausschlussstand ohne dauerhaften Träger"

### F-2 — `exclude_column` gegen eine nicht gebundene Tabelle wird `applied`, ohne je wirken zu können

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §Klassifikation (MEDIUM-Bullet
  „unklare Fehlerbehandlung am Rand des Spec-Bereichs") ·
  [`LH-FA-CFG-005`](../../spec/lastenheft.md) (Happy-Path-Prämisse „Given CDC
  ist für `t` aktiviert" und Negative „dann folgt ein expliziter
  Fehlerpfad") · [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  §Teilfrage 5 (gewählter Rückkanal ist der Antrags-Status — „dieselbe
  Sichtbarkeit, die jeder andere Verarbeitungsfehler dieser Queue bereits
  trägt") · §Teilfrage 1 Schlussabsatz („ein Spaltenausschluss, der von Anfang
  an gelten soll, wird **nach** der Erstaktivierung … beantragt")
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:441-449`
  (`ExcludeColumn`, No-op-Zweig; Doc „derselbe idempotente Vertrag wie
  `RemoveBinding`") · `internal/bootstrap/wiring.go:1073-1082` (Erfolg →
  `MarkApplied`, unabhängig von der Wirkung) ·
  `internal/adapters/driving/replication/mapper/mapper_test.go:714-735`
  (`TestExcludeColumnOnUnboundTableStaysNoop`) ·
  `internal/bootstrap/administration_internal_test.go:349-390`
  (`TestProcessAdministrationRequestsExcludeColumnAppliesWithoutBinding`)
- `befund`: Trifft ein `exclude_column`-Antrag eine in der laufenden Erfassung
  nicht gebundene Tabelle, endet er als `applied` ohne jede Wirkung — und er
  kann auch später keine bekommen: eine nachfolgende Aktivierung derselben
  Tabelle läuft über `AddBinding`, das bei fehlender Vorbindung nichts erbt;
  der Antrag wird nicht erneut gelesen. Der angeführte idempotente Vertrag von
  `RemoveBinding` trägt hier nicht: dort ist der Zielzustand bereits
  hergestellt, hier ist er unerreichbar. `LH-FA-CFG-005`s Negative-Kriterium
  ist **nicht** verletzt (es verlangt den Fehlerpfad für eine nicht
  existierende Spalte; die Spalte existiert und `ColumnExists` bestätigt das)
  — der Fall liegt am Rand des Spec-Bereichs, den erst `ADR-0059` §Teilfrage 5
  berührt und für diese Konstellation nicht entscheidet. Der gewählte
  Rückkanal (Antrags-Status) meldet in genau diesem Fall Erfolg; die
  Konstellation, die als Fehler sichtbar sein müsste, ist die einzige, in der
  er schweigt. Der Plan-Nachzug schreibt den No-op als Absicht fest, ein Test
  pinnt ihn.
- `verifizierbar`: ja — `make test`/`make test-store`; der No-op ist heute
  sogar als Zusage getestet (`TestExcludeColumnOnUnboundTableStaysNoop`,
  `TestProcessAdministrationRequestsExcludeColumnAppliesWithoutBinding`), ein
  abweichendes Verhalten machte beide rot
- `klasse`: „Stiller Erfolg für einen nie wirksamen Antrag"

### F-3 — Kopplungs-Kommentar über dem `WithoutBinding`-Test nicht auf die neue Verdrahtung nachgezogen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klasse „Kopplung" beschreibt, was da
  ist) · Präzedenzfall Review zu `slice-066` F-5 (dieselbe Klasse)
- `pfad`: `internal/bootstrap/administration_internal_test.go:342-347`
- `befund`: Der Test-Doc sagt weiterhin „die beiden Spalten-Antragsarten tragen
  keine `Assembler`-Bindung nach" — seit diesem Diff tragen sie sehr wohl
  etwas nach (`wiring.go:1080`/`:1091` rufen `ExcludeColumn`/`IncludeColumn`
  nach erfolgreichem Use Case), nur eben keine **neue** Bindung. In der
  jetzigen Formulierung widerspricht der Satz der Zeile, die der Diff zwei
  Bildschirme darunter neu geschrieben hat; gemeint ist ausschließlich „legt
  keine Bindung an".
- `verifizierbar`: nein — kein Gate prüft Kommentar-Inhalte; die Aussage des
  Tests selbst steht in `:385-388`
- `klasse`: „Kopplungs-Kommentar neben geändertem Artefakt nicht nachgezogen"

### F-4 — Der Schema-Bump-Pfad schreibt eine Versions-Kennung aus einem Schnappschuss auf eine inzwischen getauschte Bindung

- `kategorie`: INFO
- `quelle`: Maintainability · `ADR-0050` (`tablesMu`, zwei Goroutinen)
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:386`
  (`setSchemaVersion`) und `:422-430` (die Methode)
- `befund`: `observeRelation` liest die Bindung als Schnappschuss
  (`lookupBinding`), registriert danach die neue Version über den Store
  (Netz-/DB-Zugriff, kein kurzes Fenster) und ruft erst dann
  `setSchemaVersion`. Legt die Administrations-Goroutine in dieses Fenster
  eine Deaktivierung und eine neue Aktivierung derselben Tabelle, schreibt
  `setSchemaVersion` die aus dem alten `TableID` abgeleitete Versions-Kennung
  auf die neue Bindung — geprüft wird nur, **ob** eine Bindung da ist, nicht
  **welche**. Das Fenster ist vorbestehend (die vorige Form schrieb das ganze
  `TableBinding` des Schnappschusses und dazu dessen `TableID`); der Diff
  verkleinert es, schließt es aber nicht. Hinweis ohne erwartete Aktion an
  diesem Slice.
- `verifizierbar`: nein — kein Gate; nur durch Code-Gang oder einen
  gezielten Interleaving-Test zu zeigen
- `klasse`: „Stale-Snapshot zwischen Lese- und Schreibzugriff (vorbestehend)"

### F-5 — Die Sensitivität des Konvergenz-Tests ist einseitig: er prüft „kein Sonderfall", nicht „Erhalt"

- `kategorie`: INFO
- `quelle`: [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  §Teilfrage 4 („kein Sonderfall, reine Schichtung reicht")
- `pfad`: `internal/adapters/driving/replication/mapper/mapper_test.go:629-652`
- `befund`: Eigene Mutationsläufe F/G: der Test wird rot, sobald jemand den
  Ausschluss **in** den Spaltenvergleich hineinzieht (F), bleibt aber grün,
  wenn die Ausschlussliste auf einen Fremdnamen zeigt (G) — er hängt also
  nicht am Ausschlussstand, sondern an der Abwesenheit einer Sonderbehandlung.
  Genau das verlangt `ADR-0059` §Teilfrage 4, und der reale Spaltenverlust
  gehört zur Laufzeit-Ebene von `slice-068`; der Hinweis richtet sich an
  niemanden, der in diesem Slice etwas ändern müsste. Die generische
  `relationOther`-Strecke trägt zusätzlich `TestConsumeRelationOtherChangeReportsSchemaError`
  (`:463`, Typänderung).
- `verifizierbar`: ja — eigene Mutationsläufe (oben dokumentiert)
- `klasse`: „Test-Sensitivität einseitig (Nicht-Sonderfall statt Erhalt)"

---

## Urteil zu den zwei vom Implementer gemeldeten offenen Befunden

**1. Ausschlussstand überlebt keinen Prozess-Neustart — gegen den
ADR-Wortlaut entschieden: keine Zusage der Dauerhaftigkeit.** `ADR-0059` sagt
den Live-Reload zu („ohne Prozess-Neustart wirksam"), und die Option B/C-Contra
benennt die „Neustart-Einschränkung" als „wirkt nur beim Prozessstart" —
gesucht war also der Zug auf einen *laufenden* Prozess, nicht die Persistenz
des Zustands. Der Satz „löst die Neustart-Einschränkung strukturell (derselbe
Live-Reload-Pfad wie Tabellen-Aktivierung)" trägt die Dauerhaftigkeit **nicht**
ausdrücklich: er benennt den *Pfad*, und die Parenthese vergleicht mit einer
Fähigkeit, die dauerhaft ist — genau dort geht der Vergleich in die Irre.
**Ergebnis:** F-1 ist kein ADR-Verstoß und damit kein HIGH gegen `slice-067`.
Es bleibt aber mehr als die im §6 benannte Grenze (der `disable`/`enable`-
Zyklus ist nicht Neustart-abhängig und liegt in den von diesem Diff geänderten
Zeilen), und die Entscheidung darüber ist keine Risk-Ausgang-Frage des
Planners allein: `ADR-0059` ist `Accepted` und immutabel, eine präzisierende
Aussage zur Dauerhaftigkeit wäre eine Folge-ADR des Architects. Der
Risk-Ausgang in §6 darf deshalb nicht auf „weiter offen" laufen, ohne dass die
Entscheidung vorher zugewiesen ist.

**2. Stiller No-op für eine ungebundene Tabelle — der Vertrag trägt nicht.**
`LH-FA-CFG-005`s Negative-Kriterium ist formal nicht verletzt (es verlangt den
Fehlerpfad nur für eine nicht existierende Spalte), die Analogie zu
`RemoveBinding` ist aber sachlich falsch: dort ist der Zielzustand schon da,
hier ist er unerreichbar, und der einzige in `ADR-0059` §Teilfrage 5 gewählte
Rückkanal — der Antrags-Status — meldet Erfolg. Die Anforderung schweigt, die
ADR entscheidet den Fall nicht; damit ist es ein MEDIUM am Spec-Rand und keine
Implementer-Nachlässigkeit: der umsetzende Lauf folgt dem von ihm selbst
dokumentierten Plan-Nachzug, und ob der Antrag künftig `failed` lautet oder die
Grenze akzeptiert wird, ist eine Entscheidung, nicht eine Korrektur (Modul 8:
Plan-/ADR-Korrektur vor Neu-Implementierung).

## Negativbefunde

- **geprüft, ohne Befund: Punkt 4 (Abwesenheits-Kodierung, am Code).**
  `mapper.go:532` überspringt einen ausgeschlossenen Spaltennamen über
  dieselbe `continue`-Verzweigung wie `values[i] == nil` (NULL/TOAST) — kein
  Platzhalter, kein eigener Schlüssel, kein Sonderwert.
  `grep -rn "rowImage\|ExcludedColumns"` über die Produktionspfade zeigt
  genau **einen** Auswertungsort der Liste (`rowImage`), und `NewImage`/
  `OldImage` entstehen nur dort (`internal/adapters/driven/postgresstorage/mapper/mapper.go:80-81`
  reicht sie unverändert weiter) — kein zweiter Serialisierungspfad, über den
  ein ausgeschlossener Wert doch noch die Persistenzschicht erreichen könnte
  (`ADR-0059` §Teilfrage 3, Option D).
- **geprüft, ohne Befund: Punkt 3, erster Erhalt-Punkt (`AddBinding`-Merge).**
  `mapper.go:406-411` übernimmt ausschließlich `existing.ExcludedColumns`;
  `TableID` und `SchemaVersion` kommen aus dem Parameter. Sachlich richtig:
  der Enable-Zweig setzt die neuen Kennungen und lässt den Ausschluss stehen.
  Belegt durch Mutationslauf B (ohne den Merge wird
  `TestAddBindingKeepsExclusionState` rot).
- **geprüft, ohne Befund: Punkt 3, zweiter Erhalt-Punkt (`setSchemaVersion`).**
  `mapper.go:422-430` hebt nur `SchemaVersion`, ein nicht (mehr) getragener
  Eintrag bleibt ohne Wirkung — kein Wiederbeleben. Die Verzweigung ist nur
  über eine konkurrierende `RemoveBinding` erreichbar (der Aufrufer
  `observeRelation` prüft selbst auf `activated`); ein Wiederbeleben wäre die
  falsche Semantik, weil der Schnappschuss einen veralteten `TableID` trägt
  (siehe F-4).
- **geprüft, ohne Befund: die Unveränderlichkeits-Zusage des Plan-Nachzugs.**
  Jeder Schreibzugriff auf die Liste läuft über `appendExcluded`/
  `removeExcluded` (neu aufgebaut, `mapper.go:468-494`) oder reicht eine
  unveränderte geteilte Liste weiter; `IncludeColumn` baut auch ohne Treffer
  neu auf. `TestAssemblerColumnExclusionIsRaceFree` unter `go test -race` ist
  die passende Fitness Function (`ADR-0059` §Fitness Function) und läuft im
  `make test`-Lauf grün.
- **geprüft, ohne Befund: Punkt 5 (Konvergenz-Test).** Der Test übt
  `relationOther` über die echte `observeRelation`-Strecke aus, mit
  `store.registrations == 0` als zusätzlicher Zusage; die Strecke ist
  dieselbe wie bei einer Typänderung (`TestConsumeRelationOtherChangeReportsSchemaError`).
  Die Grenze seiner Aussagekraft steht als F-5.
- **geprüft, ohne Befund: Punkt 6 (die zwei „leeren Zeiger"-Funde) — im
  Scope, und die Zusagen sind echt.** `internal/bootstrap/administration_internal_test.go`
  und `administration_endtoend_test.go` stehen beide in der §3-Tabelle des
  Plans (Nachtrag im selben Commit, dieselbe Praxis wie `slice-066`), und
  beide Dateien testen `applyAdministrationRequest` — die von diesem Slice
  geänderte Funktion — in ihrem eigenen Paket. Kein Scope-Creep.
  `deps.assembler` ist produktiv gesetzt (`wiring.go:594`,
  `stream.Assembler()`), der nil-Zeiger war eine reine Test-Verdrahtungs-Lücke.
  Beide neuen Zusagen sind real ausgeübt: die E2E-Zusage
  (`administration_endtoend_test.go:127`, `:153`) geht über die echte
  SQL-Funktion → `ListPending` → echten `ColumnExclusionPort` → Assembler →
  Row Image, die Whitebox-Zusage
  (`TestProcessAdministrationRequestsExcludeColumnFiltersAssemblerRowImage`)
  prüft Ausschluss **und** Wieder-Einschluss über denselben Weg. Kein
  Widerspruch zu §1 des Plans (E2E am laufenden Container bleibt `slice-068`).
- **geprüft, ohne Befund: Punkt 7 (Kommentar-Disziplin, per `grep`).**
  Muster über die **hinzugefügten** Zeilen aller geänderten `.go`-Dateien:
  `slice-[0-9]+|welle-[0-9]+|bislang|zuvor|früher|jetzt|nicht mehr|vorher|ehemals|erstmals|statt bisher|nunmehr|hätte|wäre|würde`.
  Genau fünf Treffer, keiner in einem Produktionscode-Kommentar als Chronik:
  `mapper.go:95` („ohne Synchronisation wäre … eine Data Race" — unveränderte
  Begründung des Mutex, Bestandsformulierung), `mapper.go:400-402`
  („für eine bislang nicht aktivierte Tabelle" — Satzsubjekt ist der
  Aufruf-Zustand der Tabelle, nicht die Historie des Code-Pfads),
  `mapper_test.go:629`/`:643` und `administration_internal_test.go:429`/`:445`
  (Testfall-Szenario und Mutations-Hypothese; dieselbe Form ist im
  Review-Bericht zur Kommentar-Bereinigung im Umfeld von `slice-018` §4 als unauffällig bewertet). **Kein** Treffer
  der Klasse `BEO-PGC/slice-chronik-in-code-kommentar` in diesem Diff — die
  bei `slice-066` sechste real aufgetretene Klasse wiederholt sich hier nicht.
- **geprüft, ohne Befund: Punkt 8 (Scope-Fidelity).**
  `git diff 541ee01 45c619b --stat` berührt sieben Dateien: die zwei geänderten
  Produktionsdateien (`mapper.go`, `wiring.go`), die drei Testdateien, den
  eigenen Slice-Plan und `harness/image-hash.txt`. **Nicht** darunter:
  `test/integration/`, `tools/harness/run-integration-tests.sh`, `compose.yaml`,
  `.github/workflows/`, `docs/user/`, `spec/*` oder ein fremdes
  Planungsdokument. `git show --stat` je Commit bestätigt die Aufteilung
  (Image-Digest in eigenem Commit, `AGENTS.md` §3.3-Geist).
- **geprüft, ohne Befund: `harness/image-hash.txt`.** Der Digest wurde neu
  gestempelt, weil `mapper.go` Build-Kontext ist. Nach `ADR-0044` ist der Wert
  ein Lauf-Beleg und ausdrücklich kein Inhalts-Fingerabdruck; ein eigener
  `make image`-Lauf (oder ein Digest-Vergleich über Läufe) belegt deshalb
  nichts und wurde bewusst nicht als „Nachmessung" geführt.
- **geprüft, ohne Befund: Traceability/ID-Schema.** `1943d86` und `45c619b`
  tragen je `LH-FA-CFG-005` (der Feature-Commit zusätzlich `ADR-0059`) im
  Betreff, kein `SPEC-*`/`ARC-*` im Betreff; `make commit-traceability` grün
  (Exit 0). Die neue Kennung folgt `MR-000`.
- **geprüft, ohne Befund: DoD-Häkchen (`1943d86`).** Gesetzt sind genau die
  umsetzungs- und gate-tragenden Zeilen (die drei Liefer-Punkte, `make gates`,
  Doku-Update); Review, Closure-Notiz, Reconciliation-/Beobachtungs-Register,
  Risiko-Ausgänge und die drei Paarungen bleiben offen — keine überschrittene
  Checkbox.
- **geprüft, ohne Befund: `make gates`, `make test`, `make test-store`,
  `make test-replication`** — vier eigene Läufe, vier Mal Exit 0 (Tabelle
  oben).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Ausschlussstand ohne dauerhaften Träger ·
Stiller Erfolg für einen nie wirksamen Antrag · Kopplungs-Kommentar neben
geändertem Artefakt nicht nachgezogen · Stale-Snapshot zwischen Lese- und
Schreibzugriff (vorbestehend) · Test-Sensitivität einseitig
(Nicht-Sonderfall statt Erhalt)

## Verdikt

**Merge-blockierend:** nein — begründete Abweichung vom Regelfall: die beiden
Commits liegen bereits auf `main` (`git branch --contains 45c619b`), es gibt
keinen offenen PR, und keines der beiden MEDIUM-Findings ist eine Korrektur am
Diff, sondern je eine offene **Entscheidung** (Dauerhaftigkeit des
Ausschlussstandes; Fehlerpfad am Spec-Rand). Sie blockieren damit nicht den
Merge, sondern die **Closure**: `slice-067` darf nicht nach `done/`, solange
F-1/F-2 keinen Ausgang tragen, und der Risk-Ausgang in §6 darf für F-1 nicht
„weiter offen" lauten, ohne dass die Entscheidung zuvor zugewiesen ist
(Modul 5 §Offene Risiken werden bei Closure aufgelöst). Die Richtung des Slice
ist nicht bestritten: die Filterung sitzt an genau einem Ort und nutzt die
bestehende Abwesenheits-Kodierung, beide Erhalt-Punkte sind durch eigene
Mutationen als tragend bestätigt, die Konvergenz ist ein reiner
Schichtungs-Effekt, die Verdrahtung ist mit echten Zusagen belegt, und vier
Sensoren laufen Exit 0.

**Übergabe:** keine Fixrunde am Implementer — kein Reviewer→Implementer-Pfeil.
F-1 und F-2 gehen als Closure-Nachzug an Planner → Architect (Übergabe-Artefakt
ist dieser Report; `ADR-0059` ist `Accepted`, eine präzisierende Aussage zur
Dauerhaftigkeit oder zum Fehlerpfad wäre eine Folge-ADR). F-3 ist ein
Ein-Satz-Nachzug im Test-Kommentar und geht mit denselben Weg; F-4/F-5 sind
Hinweise ohne erwartete Aktion (F-5 betrifft ohnehin `slice-068`s Ebene).
Da keine Rückgabe an den Implementer erfolgt, ist die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im selben Commit, der
diesen Report anlegt, auf `[x]` gezogen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde). Dieser Report ist Lauf-Beleg und wird
über Läufe hinweg nicht erneut gelesen; Verifikation gegen DoD/Spec bleibt
Aufgabe des Verifiers (Modul 11; anderer Eingabe-Kontext).
