# Slice slice-046: cdc_storage_bytes-Metrik

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-13 — vierter Slice, unabhängig von den übrigen drei
(reine View-Nacharbeit, analog zu `cdc.metrics`).

**Bezug:** [`LH-FA-RET-006`](../../../../spec/lastenheft.md).

**Berührte Spec-Stellen:** — (reine SQL-View-Ergänzung, keine neue
Architektur-Sicht-Aussage).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `cdc.metrics` (`tools/schema/nacharbeit-observability.sql`)
bekommt eine neue Zeile `cdc_storage_bytes` über `pg_relation_size('cdc.change')`
(oder eine geeignete Aggregation mehrerer `cdc`-Tabellen — Implementer
entscheidet und begründet im Plan-Nachzug). Der Architect-Verdikt
([`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`](../../../../docs/reviews/architect-verdict-retention-loeschausfuehrung.md))
bestätigt vorab: das etablierte View-Owner-Muster (View läuft mit den
Rechten ihres Eigentümers, `GRANT SELECT` nur auf die fertige View an
`cdc_reader`) trägt das bereits — keine Rollen-Erweiterung nötig, da
`pg_relation_size()` eine reguläre, `PUBLIC`-ausführbare Funktion ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Rollen-/Grant-Änderung an `cdc_reader`** — der Architect-Verdikt
  schließt das ausdrücklich aus; die neue Metrik-Zeile kommt über das
  bestehende View-`GRANT`, keine neue Berechtigung.
- **Weitere, laut `nacharbeit-observability.sql` bereits als „nicht
  abgedeckt" benannte Metriken** (`cdc_changes_pending`,
  `cdc_errors_total`, `cdc_wal_retention_bytes` als `cdc.metrics`-Zeile
  statt eigenem Log) — anderer Vorgang, hier nicht angefragt;
  `cdc_wal_retention_bytes` existiert bereits als eigene, strukturierte
  Log-Ausgabe (`slice-020`/`ADR-0049`), nicht als `cdc.metrics`-Zeile.
- **Löschausführung/Sichtbarkeit blockierender Consumer** — `slice-043`/
  `044`/`045`; dieser Slice liefert ausschließlich die Speicher-Metrik.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `cdc.metrics` trägt eine neue `cdc_storage_bytes`-Zeile, real
      gegen PostgreSQL getestet (ein numerischer Wert > 0 nach dem
      Einfügen von Testdaten).
      `TestMetricsViewCarriesStorageBytes`
      (`internal/adapters/driven/postgresstorage/roles_test.go`), real
      grün über `make test-store`, siehe Plan-Nachzug Punkt 1.
- [x] `LH-FA-RET-006` real erfüllt: ein Integrationstest liest die neue
      Metrik über `cdc.metrics` (analog zum bestehenden
      `cdc_capture_lag`-Testmuster).
      `TestMVPMetricsCarriesStorageBytes`
      (`test/integration/integration_test.go`), real grün über
      `make test-integration`, siehe Plan-Nachzug Punkt 2.
- [x] Bestätigt: `cdc_reader`s Grant-Fläche bleibt unverändert (kein
      neuer direkter Grant außerhalb des View-`GRANT SELECT`). Siehe
      Plan-Nachzug Punkt 3 — `TestCdcReaderRoleReadsViewsNotBaseTables`
      unverändert grün, kein neuer Grant in
      `tools/schema/nacharbeit-roles.sql`.
- [x] `make gates` grün, `make test-integration` grün. Beide real
      ausgeführt (Ausgaben im Implementer-Bericht); zusätzlich
      `make test-store` real grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` §„Metriken lesen"
      nennt die neue Zeile.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — keine Beobachtung angefallen.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6 — beide entfallen.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Repo mit Wellen-Betrieb (`welle-13` offen) — Prüfung läuft bei der `welle-13`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/nacharbeit-observability.sql` | update | neue `cdc_storage_bytes`-Zeile in `cdc.metrics` |
| `test/integration/integration_test.go` | update | Testfall für die neue Metrik |
| `docs/user/benutzerhandbuch.md` | update | neue Metrik dokumentiert |

### Plan-Nachzug (nach Implementierung)

Regeln dieser Sektion: Implementierungsentscheidungen, die über die Tabelle
oben hinausgehen — Aggregation, Testansatz (zwei Testebenen statt einer) und
die reale Bestätigung der unveränderten `cdc_reader`-Grant-Fläche.

**1. Aggregation: einzelne Tabelle `cdc.change`, keine Summe mehrerer
`cdc`-Tabellen.** `pg_relation_size('cdc.change')::numeric` als zusätzlicher
`UNION ALL`-Zweig der bestehenden View, exakt wie im Ziel dieses Slice-Plans
und im Architect-Verdikt
([`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`](../../../reviews/architect-verdict-retention-loeschausfuehrung.md),
Frage 2) besprochen. Keine Summe über `cdc.transaction`,
`cdc.source_table`, `cdc.schema_version`, `cdc.consumer`,
`cdc.consumer_position` gebildet: `cdc.change` trägt die Row Images als
`jsonb` (`LH-FA-CAP-008`) und wächst mit jedem erfassten Change; die übrigen
fünf Tabellen tragen ausschließlich Referenz-/Katalogdaten (Quelle,
Tabellen-/Schema-Katalog, Transaktions-Kopf, Consumer-Zustand) mit fester
oder linear zur Zahl der Quelltabellen/Consumer wachsender Zeilenzahl, nicht
zum Erfassungsvolumen. Eine Summe hätte keinen zusätzlichen
Informationsgewinn für den in `LH-FA-RET-006` benannten Zweck
(„Kontrolle des Datenwachstums") geliefert, aber die Bedeutung der Kennzahl
verwässert (ein Wachstum von `cdc.change` wäre neben dem konstanten Anteil
der übrigen Tabellen schwerer erkennbar) und die Abfrage unnötig verbreitert.

**2. Zwei Testebenen, kein Widerspruch zur Slice-Größe.** Die DoD zählt einen
Liefer-Punkt „real gegen PostgreSQL getestet" — dieser Punkt ist über zwei
Testfälle auf unterschiedlichen Testebenen belegt, nicht über zwei
Liefer-Punkte:
`TestMetricsViewCarriesStorageBytes`
(`internal/adapters/driven/postgresstorage/roles_test.go`, `make test-store`)
fügt eine vollständige Change-Zeile über direkte SQL-Inserts ein (Transaktion,
Tabellen-/Schema-Referenz, Change) und prüft den numerischen Wert isoliert
gegen eine frisch ausgerollte Instanz — dieselbe Testebene und dasselbe Muster
wie das bestehende `TestMetricsViewCarriesConsumerLag` für `cdc_consumer_lag`.
`TestMVPMetricsCarriesStorageBytes`
(`test/integration/integration_test.go`, `make test-integration`) liest
denselben Wert am verdrahteten Feed-Container nach einer realen
CDC-Erfassung — dieselbe externe SQL-Lesezugriffsweg-Disziplin wie der
bestehende `cdc_capture_lag`-Lasttest-Beleg
(`tools/harness/run-integration-tests.sh`). Die neue Testfunktion wurde in
das bestehende `-run`-Muster in `run-integration-tests.sh` aufgenommen
(`BEO-PGC/test-runner-stiller-ausschluss`: eine Testfunktion, die in keinem
`-run`-Muster auftaucht, liefe unter `make test-integration` dauerhaft und
stillschweigend nie) — real bestätigt: `go test ./test/integration/...` mit
gesetztem `-run` zeigt `TestMVPMetricsCarriesStorageBytes` explizit in der
Ausgabe.

**3. `cdc_reader`-Grant-Fläche real unverändert bestätigt.** Kein neuer
Eintrag in `tools/schema/nacharbeit-roles.sql` — das bestehende
`GRANT SELECT ON cdc.metrics TO cdc_reader;`
(`tools/schema/nacharbeit-observability.sql`) trägt die neue Zeile bereits
mit, weil sie ein zusätzlicher `UNION ALL`-Zweig derselben View ist, kein
neues Objekt. Real bestätigt statt nur angenommen:
`TestCdcReaderRoleReadsViewsNotBaseTables`
(`internal/adapters/driven/postgresstorage/roles_test.go`) bleibt
unverändert und lief real grün gegen die um `cdc_storage_bytes` erweiterte
View — derselbe Testfall, der bereits vor diesem Slice `cdc_reader`s
`SELECT`-Zugriff auf `cdc.metrics` insgesamt belegt (Zeilenzahl > 0). Eine
isolierte Prüfung nur der neuen Zeile war nicht nötig: PostgreSQLs
View-Owner-Semantik (Definer ohne `security_invoker`) kennt keine
Zeilen-/Spalten-granulare Rechteprüfung innerhalb einer View — der Zugriff
gilt für die View als Ganzes oder gar nicht. Der Architect-Verdikt
(Frage 2) hatte das bereits vorab hergeleitet
(`pg_relation_size()` ist eine reguläre, für `PUBLIC` ausführbare
Systemfunktion ohne eigenes Privileg auf `cdc.change`); die reale Prüfung
deckte keine bislang unbekannte PostgreSQL-Versions-/
Berechtigungs-Eigenheit auf (§6, Risiko 2 — Ausgang: entfallen).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-13` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei —
unabhängig von `slice-043`/`044`/`045`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei einer einzelnen Metrik-Zeile — falls doch, wäre das ein
  Zeichen für eine unerwartet komplexe Aggregation (mehrere Tabellen,
  Indizes).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- `pg_relation_size('cdc.change')` könnte nach einer realen Löschung
  (`slice-043`/`044`) durch PostgreSQLs Tabellen-Bloat (gelöschte
  Tupel, kein automatisches `VACUUM FULL`) einen irreführend hohen Wert
  zeigen, der den tatsächlich freigegebenen Platz nicht widerspiegelt.
  **Ausgang: entfallen.** `LH-FA-RET-006`s Akzeptanzkriterium verlangt
  genau das Gegenteil einer „Bereinigung" des Werts: „Given wachsende
  CDC-Daten, when der Verbrauch beobachtet wird, then ist er über die
  Metriken ablesbar" (Happy Path) und „Given die Retention kann ein
  Wachstum nicht begrenzen … then ist er erkennbar" (Boundary,
  `spec/lastenheft.md`). Bloat aus verzögerter, durch einen
  zurückhängenden Consumer blockierter Löschung (`LH-FA-RET-004`) ist
  reales, noch nicht freigegebenes physisches Datenwachstum — genau der
  Zustand, den dieser Boundary-Fall sichtbar verlangt. Ein Wert, der Bloat
  herausrechnete, wäre die irreführende Variante, nicht die hier gebaute.
- Der Architect-Verdikt zum View-Owner-Muster wurde vor der
  tatsächlichen Implementierung getroffen — eine reale Prüfung könnte
  eine bisher unbekannte PostgreSQL-Versions-/Berechtigungs-Eigenheit
  aufdecken. **Ausgang: entfallen.** Real durch
  `TestCdcReaderRoleReadsViewsNotBaseTables`
  (`internal/adapters/driven/postgresstorage/roles_test.go`) bestätigt,
  unverändert grün gegen die um `cdc_storage_bytes` erweiterte View — kein
  neuer Grant nötig, keine unbekannte Eigenheit aufgetreten
  (Plan-Nachzug Punkt 3).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Die einzelne `UNION ALL`-Zeile über
  `pg_relation_size('cdc.change')` hielt den Slice tatsächlich auf einen
  Liefer-Punkt für die View selbst, ohne Rollen-/Grant-Änderung — der
  Architect-Verdikt (Frage 2) traf real zu: `TestCdcReaderRoleReadsViewsNotBaseTables`
  blieb unverändert grün, kein neuer Eintrag in
  `tools/schema/nacharbeit-roles.sql` nötig. Die zwei Testebenen
  (`make test-store` für den isolierten numerischen Beleg,
  `make test-integration` für den End-zu-Ende-Beleg über den verdrahteten
  Feed-Container) ließen sich beide direkt aus bereits etablierten Mustern
  ableiten (`TestMetricsViewCarriesConsumerLag`,
  `cdc_capture_lag`-Lasttest-Beleg) — kein neuer Testansatz nötig.
- **Was ging anders als geplant:** Nichts Wesentliches — die Implementierung
  folgte dem im Architect-Verdikt vorgezeichneten Weg ohne Abweichung. Eine
  neue Testfunktion in `test/integration/integration_test.go` musste
  zusätzlich in das bestehende `-run`-Filtermuster in
  `tools/harness/run-integration-tests.sh` aufgenommen werden, sonst liefe
  sie unter `make test-integration` nie (`BEO-PGC/test-runner-stiller-ausschluss`,
  bereits bekannte, unter der Schwelle liegende Beobachtung) — real geprüft,
  keine neue Instanz dieser Klasse.
- **Steering-Loop-Eintrag:** *(kein Eintrag verkörpert — der Normalfall.)*
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen. `BEO-PGC/retention-keine-loeschausfuehrung` bleibt bei 0×
  (analog zu `slice-043`/`044`/`045`) — mit diesem Slice sind alle vier in
  der Beobachtung benannten Lücken (Löschausführung, Hintergrundjob,
  Sichtbarkeit blockierender Consumer, `cdc_storage_bytes`-Metrik)
  geliefert; der Ausgang selbst bleibt der `welle-13`-Closure vorbehalten
  (Lese-Schritt, Modul 6). `BEO-PGC/test-runner-stiller-ausschluss` bleibt
  bei 1× (kein neues Auftreten, siehe oben).
- **Folge-Slices:** keine neuen — `welle-13` trägt keine weiteren Slices
  über `slice-046` hinaus.
- **Risiken aus §6:** beide *entfallen* — siehe §6.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-13` offen) —
  Prüfung läuft bei der `welle-13`-Closure. Kein `liegt in`-Feld in diesem
  Slice (nichts verkörpert), also kein Anker-Paarungs-Gegenstand aus diesem
  Slice selbst.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC`: `BEO-PGC/retention-keine-loeschausfuehrung` (0×,
benannt nicht gezählt). Keiner erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
