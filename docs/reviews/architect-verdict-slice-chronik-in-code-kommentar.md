# Architect-Verdikt: Chronik-Sprache in Code-Kommentaren — Sensor vs. geschärfte Prosa

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar`
erreicht 3× (dem Beleg zum Review-Report-Commit zu `slice-041`, dem Beleg
zum Review-Report-Commit zur Fixrunde von `slice-041`,
dem Beleg zum Review-Report-Commit zu `slice-044`) — Lese-Schritt des wellenlosen Betriebs
(Modul 6 „Träger im Repo ohne Wellen"), ausgelöst direkt aus der
`slice-044`-Closure, ohne hostende Welle. Planner → Architect-Zug.
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** [AGENTS.md](../../AGENTS.md) §3.7 (Hard Rule „Ein Kommentar
beschreibt, was da ist"), [ADR-0014](../plan/adr/0014-retention-domain-policy.md),
[LH-FA-RET-002](../../spec/lastenheft.md) (thematisch nächste Kennung — die
Beobachtung wurde durch die Retention-Slices `slice-041`/`slice-044`
ausgelöst, ist selbst aber keine Retention-Entscheidung),
`.claude/commands/implement-slice.md` Schritt 20 (bestehende, bereits
gescheiterte Verkörperung), `tools/harness/commit-traceability.sh`
(Sensor-Muster-Referenz, geprüft und verworfen), `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar`.

---

## Frage

Reicht eine geschärfte Prosa-Instruktion (Implementer-Workflow), oder
braucht die Beobachtung einen echten mechanischen Sensor (`.go`-/SQL-Diff
gegen ein Chronik-Sprachmuster, als neues `make`-Target)?

## Verdikt

**Kein mechanischer Sensor. Geschärfte, diff-skopierte Selbstprüf-Instruktion**
in `.claude/commands/implement-slice.md` Schritt 20 — Begründung unten. Der
bereits existierende Prosa-Schritt 20 war nicht falsch, nur unvollständig:
Er nennt die Probe (Ist-Zustand vs. Chronik), verlangt aber keine
**Enumeration** der Kandidatenstellen. Das dritte Vorkommen
(Review zu `slice-044` F-1) trat ein, obwohl derselbe Implementer im
selben Commit bereits zwei andere Chronik-Stellen selbst korrigiert
hatte — ein visueller Scan hat eine von mehreren Stellen übersehen. Das
ist eine Enumerations-Lücke, keine Verständnis-Lücke, und dafür ist die
richtige Antwort ein deterministischer Kandidatenlauf vor der Übergabe,
nicht (zwangsläufig) ein Gate.

## Geprüft: ist ein grep-Sensor praktikabel? (Frage 1+2 der Anfrage)

**Empirischer Befund gegen die Annahme „`slice-\d+` kommt in `.go`
nie legitim vor":** Ein Repo-weiter Suchlauf
(`grep -rnE '\b(slice|welle)-[0-9]+' --include='*.go' --include='*.sql'`)
liefert über 30 Treffer in bereits gemergtem, reviewtem Code —
durchweg Godoc-Kommentare direkt über `Test*`-Funktionen, die
begründen, **warum ein bestimmter Testfall existiert**: Beispiele
`internal/bootstrap/acknowledge_test.go:174`
(„TestAcknowledgeConsumerReportsInvalidPosition trägt einen der zwei
externen Domänenfehler-Pfade aus dem Review zu `slice-022` F-1"),
`internal/bootstrap/diagnose_test.go:152`
(„… siehe Plan-Nachzug slice-038 §3, Risiko 1 aus §6"),
`internal/adapters/driving/replication/receive/stream_test.go:532/544/559/625`
(vier Vorkommen aus dem Review zu `slice-025` F-1),
`internal/adapters/driven/postgresstorage/administrationrequest.go:138/214`
(Review zu `slice-037` F-3). Diese Zitierform — „dieser Testfall deckt
Finding X aus Review Y ab" — ist **etablierte, zulässige
Regressionstest-Provenienz**, keine Chronik über Produktionsverhalten;
sie ist über 15+ Dateien hinweg konsistent und wird mit jedem künftigen
Slice, der einen Review-Fund als Regressionstest verankert, **weiter
neu produziert** (kein Alt-Bestand, den man einmalig ausnimmt).

**Strukturell ununterscheidbar von den drei tatsächlichen Verstößen:**
Alle drei gezählten Funde (`config_file.go`s `mergeTables`-Kommentar,
`config_file.go:115`s `ConfigFromEnvAndFile`-Doc-Kommentar,
`retention_internal_test.go:14–16`s Dateikopf-Kommentar) sind **ebenfalls**
Godoc-Kommentare direkt über einer Funktions-/Typ-/Dateideklaration — exakt
dieselbe strukturelle Position wie die oben zitierten, akzeptierten
Testfall-Provenienz-Kommentare. Eine Heuristik über *Position*
(Kommentar unmittelbar vor `func`/`type`) oder *Dateityp* (`_test.go`
vs. Produktionsdatei) trennt die beiden Klassen **nicht**: Der dritte
Verstoß selbst steht in einer `_test.go`-Datei, an derselben Godoc-Position
wie die zulässigen Beispiele. Der einzige tragende Unterschied ist
**semantisch** — ist das Subjekt des Satzes der *Testfall* („dieser Test
existiert wegen …", zulässig) oder der *Produktionscode-Pfad* („der Code
verhält sich so wegen …", Chronik, unzulässig) —, und das ist
Prosa-Verständnis, kein Zeichenketten-Muster.

**Die implizite Variante (Vorher/Nachher-Sprache ohne Ziffer, z. B.
„zeigt identisches Verhalten wie vor diesem Slice" aus
dem Review-Report zur Fixrunde von `slice-041`) ist noch schlechter geeignet:** Das Wort
„Slice" allein kollidiert mit Gos eingebautem Datentyp
(`[]byte`-„Slice", „Byte-Slice" etc.) — in einer Go-Codebasis ein
hochfrequentes, harmloses Wort. Ein Muster ohne Ziffernbindung
(`\bslice\b`) wäre auf jeder Zeile mit einer echten Go-Slice-Diskussion
ein Fehlalarm; ein Muster mit Ziffernbindung (`slice-\d+`) erkennt diese
Variante gar nicht erst (sie trägt keine Ziffer).

**Fazit zu Frage 1/2:** Ein textmuster-basierter Sensor über den
gesamten `.go`/SQL-Bestand ist **nicht praktikabel ohne unvertretbare
Fehlalarme** — er würde entweder (a) die etablierte,
fortlaufend produzierte Testfall-Provenienz-Konvention flächendeckend
melden (und damit entweder ignoriert oder durch eine ständig wachsende
Ausnahmeliste neutralisiert, die exakt dasselbe menschliche Urteil
dupliziert, das er ersetzen sollte), oder (b) eine
Satz-Subjekt-Erkennung brauchen, die weit über ein „kleines,
dependency-freies Skript" hinausgeht. Diese Beobachtungsklasse gehört
zur zweiten Menge aus `harness/README.md`/Modul 6 — *Urteil bleibt
Urteil*, hier ohne die dort sonst übliche urteilsfreie Formprüfung
(anders als z. B. bei den drei Risiko-Ausgängen: dort ist die Form
geschlossen, hier nicht).

## Geprüft: was mechanisiert werden kann, ohne die Klassifikation zu erzwingen (Frage 3)

Kein Klassifikations-Sensor — aber die **Enumeration** der
Kandidatenstellen lässt sich mechanisieren, ohne die
Provenienz-vs.-Chronik-Unterscheidung selbst zu automatisieren: ein
**diff-skopierter** (nicht repo-weiter) Kandidatenlauf über exakt die in
diesem Lauf geänderten `.go`-/`tools/schema/*.sql`-Dateien, mit demselben
permissiven Muster (`slice-[0-9]+`, `welle-[0-9]+`,
`vor diesem [Ss]lice`, `nach diesem [Ss]lice`, `seit diesem [Ss]lice`).
Diff-Skopierung ist tragend, nicht optional: repo-weit träfe derselbe
Lauf auf denselben etablierten Testfall-Bestand und würde beim ersten
Einsatz Dutzende Alt-Treffer melden, die niemand beheben will/muss.

Diese Enumeration behebt gezielt den beobachteten Fehlermodus — eine von
mehreren korrigierten Stellen wurde **übersehen**, nicht falsch
beurteilt —, ohne eine Klassifikations-Aufgabe zu automatisieren, die
mechanisch nicht trägt. Sie bleibt **Selbstprüfung im Implementer-Workflow
(Schritt 20)**, kein `make`-Target: Ein `make`-Target, das advisory oder
gar gate-artig liefe, würde denselben Fund-Strom in jedem Lauf erzeugen,
der eine neue Testfall-Provenienz-Zeile schreibt — advisory Rauschen ohne
Handlungsaufforderung ist so wertlos wie gar keins, und ein Gate wäre
ohne die semantische Trennung schlicht falsch (es würde legitime,
gewollte Kommentare blockieren).

## Was das für den Implementer bedeutet

`.claude/commands/implement-slice.md` Schritt 20 bekommt einen
diff-skopierten Kandidatenlauf-Befehl plus die explizite
Unterscheidungsprobe (Testfall-Provenienz vs. Produktionsverhalten-Chronik)
als Pflicht-Ergänzung — Umsetzung in diesem Zug, kein Folge-Slice nötig
(kleine Doku-Änderung an einer Prozess-Datei, kein Produktionscode).

Weder Produktionscode noch eine ADR-Datei wurden im Rahmen dieses
Verdikts geändert.
