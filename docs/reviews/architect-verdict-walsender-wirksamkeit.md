# Architect-Verdikt: Publication-Entzug-Wirksamkeit am laufenden Walsender (`BEO-PGC/walsender-wirksamkeit`)

**Rolle:** Architect (Modul 8)
**Anlass:** Nächste Roadmap-Zeile „Publication-Entzug-Wirksamkeit am
laufenden Stream" soll `BEO-PGC/walsender-wirksamkeit` schließen (Zähler
abgeleitet: 1×, `evidence/slice-008.md`, `welle-3`) — Klärung *vor* der
Slice-Planung, welcher Testansatz trägt, weil eine Fork-Recherche den
bestehenden Beleg als strukturell unzureichend identifiziert hat.
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** [`LH-FA-CFG-002`](../../spec/lastenheft.md) (CDC-Deaktivierung
je Tabelle, Happy Path: „werden fortan keine Änderungen an `t` mehr
erfasst"), [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md)
(Verwaltungsverträge), [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antragsqueue, Live-Reload), `internal/adapters/driving/replication/mapper/mapper.go`
(`Assembler.RemoveBinding`/`AddBinding`), `internal/adapters/driven/postgresstorage/tableactivation.go`
(`Unpublish`), `internal/bootstrap/wiring.go` (`applyAdministrationRequest`),
`tools/harness/run-integration-tests.sh` (SQL-Administration
Live-Reload-Beleg, disable-Zweig),
`docs/plan/planning/observations/BEO-PGC/walsender-wirksamkeit`

---

## Verdikt

**A — Testansatz spezifizieren**, aber nicht in der Form, die die
Fork-Recherche als Kandidat nennt. `pg_logical_slot_peek_changes` (oder
`_get_changes`) auf dem **vom laufenden Feed benutzten Slot** ist technisch
nicht ausführbar, solange der Feed-Container läuft: PostgreSQL erlaubt pro
logischem Replication-Slot genau einen aktiven Leser; ein zweiter
Lesezugriff auf denselben Slot scheitert mit `replication slot "..." is
active for PID ...`. Ein *zweiter, dedizierter* Slot auf derselben
Publication umgeht dieses Problem, verfehlt aber die eigentliche Frage: Eine
frisch aufgebaute Decoding-Session sieht den aktuellen Publication-Katalog
immer korrekt (`get_rel_sync_entry` liest ihn beim ersten Aufbau der
Rel-Sync-Cache neu ein) — sie kann nicht zeigen, ob eine **bereits laufende**
Session ohne Neuaufbau reagiert, und genau das ist die offene Frage der
Beobachtung.

Der tragfähige Testansatz braucht keinen zweiten Slot und kein
Slot-Introspektions-Kommando. Er trennt die beiden Filterebenen anders: über
den bestehenden Administrations-Mechanismus selbst.

## Warum der bestehende Beleg die Frage nicht beantwortet

Der „SQL-Administration Live-Reload-Beleg (disable)"
(`tools/harness/run-integration-tests.sh`, Zeilen ~933–986) führt
`cdc.disable_table('src-mvp', 'public', $ADMIN_TABLE)` aus. Das ist kein
reiner Publication-Vorgang: `applyAdministrationRequest`
(`internal/bootstrap/wiring.go`, Fall `AdministrationRequestDisable`) ruft
`DisableTableUseCase.Disable` (der intern `TableActivationPort.Unpublish` —
`ALTER PUBLICATION ... DROP TABLE` — ausführt) **und unmittelbar danach**
`deps.assembler.RemoveBinding(qualified)` auf. `RemoveBinding` löscht die
Tabellen-Bindung **im Go-Prozess selbst**, synchron zur SQL-Ausführung.

`Assembler.change()` (`mapper.go`) verwirft jede Änderung, deren
qualifizierter Name keine Bindung mehr trägt (`activated == false`), *bevor*
sie in `cdc.changes` ankommen könnte — unabhängig davon, ob PostgreSQL das
zugehörige WAL-Segment für die Tabelle noch einen Moment lang gesendet
hätte. Der Test prüft also die **App-seitige** Filterung; ob PostgreSQLs
Walsender für eine *bereits laufende* Session sofort oder verzögert
reagiert, bleibt dabei vollständig verdeckt — mit ODER ohne
Walsender-Verzögerung zeigt `cdc.changes` in beiden Fällen dieselbe leere
Menge, weil der App-Filter beide Fälle gleich behandelt. Der bestehende
Beleg ist damit real vorhanden, beweist aber eine andere Eigenschaft
(App-seitige Idempotenz-Filterung) als die, die `LH-FA-CFG-002` inhaltlich
braucht (Erfassung endet, gleich auf welcher Ebene sie endet — aber genau
*das* „auf welcher Ebene" ist die offene Frage von `BEO-PGC/walsender-wirksamkeit`).

## Der isolierte Testansatz

Die Trennung gelingt, indem der Test die beiden Schreibpfade **entkoppelt**,
die `applyAdministrationRequest` normalerweise zusammen ausführt — nicht
über ein zusätzliches PostgreSQL-Introspektions-Werkzeug, sondern indem er
den Administrations-Antragsweg (`cdc.disable_table`) für diesen einen Check
**bewusst umgeht** und die Publication-Änderung direkt per DDL setzt:

1. **Dedizierte Tabelle**, bereits über die reguläre `cdc.enable_table`
   -Antragskette aktiviert und gebunden (wie `$ADMIN_TABLE` im bestehenden
   Beleg) — nicht dieselbe Tabelle wie der enable/disable-Beleg
   wiederverwenden, damit beide Belege sich nicht gegenseitig stören.
2. **Direkt per `psql` als `$PG_USER`** (Superuser, ausreichend privilegiert)
   `ALTER PUBLICATION <publication> DROP TABLE public.<tabelle>;` ausführen
   — **ohne** `cdc.disable_table(...)` aufzurufen. Es entsteht keine
   Zeile in `cdc.administration_request`; die Administrations-Goroutine wird
   nie tätig; `Assembler.RemoveBinding` wird für diese Tabelle **nicht**
   aufgerufen. Geprüft (`internal/bootstrap/wiring.go`,
   `nacharbeit-administration.sql`): `RemoveBinding` hat im Repo genau eine
   Aufrufstelle, ausgelöst ausschließlich durch einen verarbeiteten
   `disable`-Antragsdatensatz; es gibt keinen Katalog-Watcher, der
   `pg_publication_rel`-Änderungen unabhängig davon beobachtet. Die
   Assembler-Bindung der Tabelle bleibt also über den gesamten Testverlauf
   aktiv (`activated == true`).
3. Eine neue Zeile in die Tabelle einfügen (`INSERT ... VALUES (<neue
   id>, ...)`).
4. **Real, aber begrenzt warten** (Muster wie Zeile 972 im bestehenden
   Skript: ein reales `sleep`, kein Poll auf ein Ereignis, das nicht
   eintreten *soll* — hier ist die Abwesenheit die zu prüfende Eigenschaft),
   dann `cdc.changes` auf die neue Zeile abfragen.
5. **Auswertung:**
   - Erscheint die Zeile **nicht** in `cdc.changes`: PostgreSQLs bereits
     laufende Decoding-Session hat die WAL-Änderung der nun
     unpublizierten Tabelle **sofort** ausgefiltert — der Assembler hätte
     sie (Bindung war ja noch aktiv) andernfalls unverändert
     durchgelassen. Das widerlegt die Verzögerungs-Annahme der Beobachtung
     empirisch für die geprüfte PostgreSQL-Version.
   - Erscheint die Zeile **doch**: Der Walsender hat trotz entzogener
     Publication-Mitgliedschaft weiter dekodiert und ausgeliefert — die
     Assembler-Bindung war der einzige Grund, warum sie sonst nicht
     ankäme. Das bestätigt die Verzögerungs-Annahme real und macht die
     App-seitige Doppel-Schutz-Architektur (Punkt darunter) zur
     tragenden, nicht nur zur zusätzlichen Absicherung.
6. **Nachlauf:** Tabelle wieder über `ALTER PUBLICATION ... ADD TABLE`
   publizieren (Zustand für nachfolgende Prüfungen im selben Lauf
   wiederherstellen) — oder eine Wegwerf-Tabelle verwenden, die im
   restlichen Lauf nicht mehr gebraucht wird.

## Warum das isoliert, was der bestehende Beleg nicht trennt

Der Assembler kennt nur einen Filter-Zustand pro qualifiziertem Namen
(`tables[qualified]` vorhanden oder nicht). Solange Schritt 2 diesen
Zustand *nicht* anfasst, ist jede Abwesenheit einer Änderung in
`cdc.changes` ausschließlich durch PostgreSQL selbst erklärbar — es gibt
keinen zweiten Mechanismus mehr, der als Alternativerklärung in Frage kommt.
Umgekehrt: Bliebe die Bindung aktiv und die Zeile käme trotzdem nicht an,
gäbe es keine App-seitige Erklärung dafür — der einzige verbleibende Pfad
ist PostgreSQLs eigene Dekodierung. Das ist der Unterschied zum
`pg_logical_slot_peek_changes`-Ansatz: Der hier vorgeschlagene Test braucht
keinen zweiten Zugriffspfad auf den WAL-Strom, weil er den vorhandenen
(App-seitigen) Pfad durch gezieltes Nicht-Anfassen der Assembler-Bindung
zum neutralen Beobachter macht.

**Die bestehende Doppel-Schutz-Architektur bleibt davon unberührt und
notwendig, unabhängig vom Testergebnis:** Sie schützt heute schon jede
Deaktivierung, die über den regulären Antragsweg läuft (der einzige Weg, den
`LH-FA-CFG-002` als Vertrag kennt); der oben beschriebene Test verlässt
diesen Weg bewusst, um eine isolierte Messung zu bekommen, und ist damit
kein Ersatz für den bestehenden Beleg, sondern eine Ergänzung, die eine
andere Eigenschaft prüft.

## Was das für den Planner bedeutet

Die offene Beobachtung bleibt bis zur realen Testausführung **weiter
offen** (`state.md` unverändert) — dieses Verdikt liefert die
Testspezifikation, nicht das Testergebnis; welcher der drei
Register-Ausgänge (verkörpert / geplant / gestrichen bzw. bei einem echten
Risiko: eingetreten mit Folge-Slice / entfallen mit Begründung / weiter
offen) trägt, entscheidet sich erst am Ergebnis dieses Tests, nicht hier.

Ein Folge-Slice (Replication-Berührung, Sub-Area-Kürzel `PGC`, Bezug
`LH-FA-CFG-002`, `ADR-0050`, diese Verdikt-Datei) liefert:

- den oben beschriebenen isolierten Test, ergänzend zum bestehenden
  `tools/harness/run-integration-tests.sh`-Beleg (nicht als dessen Ersatz),
  mit einer dedizierten Tabelle;
- abhängig vom realen Ausgang entweder (a) eine Notiz, dass PostgreSQL
  sofort filtert und die Beobachtung bei Closure als *entfallen* mit
  Testbeleg gestrichen wird, oder (b) — falls der Walsender tatsächlich
  verzögert liefert — eine benannte Behandlung (z. B. ein kurzer
  Grace-Wait vor dem Wirksamkeits-Anspruch von `cdc.disable_table`, oder
  eine explizite Dokumentation, dass die App-seitige Filterung die einzige
  verlässliche Ebene ist und `LH-FA-CFG-002` entsprechend zu lesen ist) als
  eigenen Liefer-Punkt.

Größenordnung: klein (ein Testabschnitt, eine dedizierte Tabelle, keine
Schema- oder Domänen-Änderung vorausgesetzt) — das ist Schnitt-Aufgabe des
Planners (Modul 5), nicht Gegenstand dieses Verdikts.

Weder Produktionscode noch eine ADR-Datei noch `state.md` wurden im Rahmen
dieses Verdikts geändert.
