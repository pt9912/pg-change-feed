# Welle 11 — E2E-Abdeckung — Verwaltung & Observability — Closure-Notiz

**Welle:** welle-11
**Abschluss:** 2026-09-13
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-CFG-003`](../../../../spec/lastenheft.md)/`004` (`slice-034`):
  ein neuer Testfall liest `cdc.active_tables` real über SQL für eine
  aktivierte und eine nie aktivierte Tabelle und hält beide Lesewege
  (SQL-Sicht, Go-Use-Case) gegeneinander.
- [`LH-FA-ADM-002`](../../../../spec/lastenheft.md)/`005` (`slice-035`):
  ein neuer Testfall liest `cdc.heartbeat` real über SQL und belegt eine
  frische, fehlerfreie Lebenszeichen-Zeile im laufenden Betrieb
  (Happy Path); ein neuer Abschnitt im Runner-Skript belegt real, dass
  `cdc.consumer_status` einen Verarbeitungsrückstand zeigt, solange ein
  Consumer nicht auf die aktuelle Position bestätigt hat, und `0` danach.
- Zusammen mit bereits vorhandener Abdeckung
  (`LH-FA-ADM-003`/`004` über `slice-033`s Fehlerklassen-Test bzw. den
  `cdc_capture_lag`-Lasttest-Beleg) sind jetzt alle vier
  Observability-Signale (`LH-FA-ADM-002`…`005`) und CDC-Status/-Liste
  (`LH-FA-CFG-003`/`004`) real black-box über externe SQL-Sichten belegt.
- Ein bedeutender, ungeplanter Fund während der Eröffnungs-Recherche:
  [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (SQL-Administration)
  verlangt explizit SQL-Funktionen für Aktivierung/Deaktivierung/Status/
  Consumer-Verwaltung, die real nicht existieren;
  [`LH-FA-CFG-002`](../../../../spec/lastenheft.md) (Deaktivierung) hat
  keinen Live-Zugriffsweg; [`LH-FA-SST-003`](../../../../spec/lastenheft.md)
  (CLI) hat keinen Status-/Diagnose-Befehl. Registriert als
  `BEO-PGC/verwaltung-keine-sql-administration`, adressiert in einer
  eigenen, bereits vorgemerkten Feature-Welle.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Die Eröffnungs-Recherche (Fork) vor dem Schneiden der Slices verhinderte
  eine falsch gescopte Welle: Ohne sie wäre versucht worden, Testabdeckung
  für eine Fähigkeit (SQL-Administration, CDC-Deaktivierung) zu bauen, die
  real gar nicht existiert — dieselbe Disziplin wie zuvor bei Retention
  und Schema-Evolution.
- Beide Implementer-Läufe fanden reale Präzisierungsbedarfe (`slice-034`s
  Boundary-Kriterien-Zuordnung, `slice-035`s `consumer_status`-
  Semantikabweichung) über echte Testläufe und dokumentierten sie
  transparent im Plan-Nachzug, statt sie zu verschweigen oder die
  ursprüngliche Formulierung stur zu erzwingen.
- Beide Slices erinnerten sich proaktiv an `BEO-PGC/test-runner-stiller-ausschluss`
  (aus `slice-033`s Closure) und nahmen ihre neuen Testfälle korrekt ins
  Runner-Skript-Muster auf.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- Die Welle wurde bei Eröffnung enger gescopt als ursprünglich in der
  Roadmap vorgemerkt: Der ursprüngliche Eintrag verwies fälschlich auf
  `BEO-PGC/rollen-test-abdeckungsluecken` (widerlegt) und vermischte
  reine Testabdeckung mit real fehlenden Fähigkeiten
  (SQL-Administration, Deaktivierung, CLI-Diagnose). Konsequenz: neue
  Feature-Welle „Verwaltungsfunktionen — SQL-Administration &
  CLI-Diagnose" vorgemerkt, `welle-11` selbst auf reine, real
  existierende Lesezugriffswege begrenzt.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreicht in dieser Welle 3× — der Normalfall.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Neu
angelegt in dieser Welle: `verwaltung-keine-sql-administration` (0×,
benannt nicht gezählt — adressiert durch die vorgemerkte Feature-Welle
„Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose"). Unter der
Schwelle, unverändert: `adapter-fehler-ausgang` (1×),
`dod-checkbox-nachzug-architect-pfad` (1×),
`retention-keine-loeschausfuehrung` (0×, benannt nicht gezählt),
`rollen-test-abdeckungsluecken` (2×, bestätigt inhaltlich nicht berührt
von dieser Welle), `schema-rollout-fremdobjekte` (1×),
`test-isolation-geteilter-zustand` (1×), `test-runner-stiller-ausschluss`
(1×), `walsender-wirksamkeit` (1×). Bereits verkörpert, unverändert:
`a-check-null-abdeckung` (3×), `d-migrate-nacharbeit` (4×),
`dod-checkbox-nachzug` (3×), `lese-doppelquelle` (3×, `seit slice-029`),
`plan-nachzug` (2×), `plan-vorlagen-defekt` (3×),
`schema-evolution-nicht-dynamisch` (2×, `seit slice-033`). Bereits
eingetreten, unverändert: `cdc-capture-lag-real` (2×),
`health-endpoint-heartbeat` (2×), `rollen-verdrahtung` (4×),
`spec008-replication-luecke` (2×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — die neu vorgemerkte Feature-Welle „Verwaltungsfunktionen —
SQL-Administration & CLI-Diagnose" ist bislang nur eine Vorschau-Zeile in
der Roadmap (*Nächste Wellen*), noch nicht als Slice geschnitten.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Beide Slices (`slice-034`, `slice-035`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure).
- Alle in `welle-11` §3 genannten Closure-Kriterien real erfüllt:
  `cdc.active_tables`-Beleg (`slice-034`), alle vier
  Observability-Signale (`LH-FA-ADM-002`…`005`) real belegt
  (`slice-033`, `slice-035`, Lasttest-Beleg) — dreifach reproduziert von
  Implementer, Reviewer und Verifier unabhängig voneinander bei beiden
  Slices dieser Welle.
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig". Kein Carveout in dieser Welle. Kein
  bootstrap-aware Gate berührt. Kein aktives ADR geändert oder berührt.
- Drei Paarungen (Anker · Folge-Slice · Register): Anker — kein
  Steering-Loop-Eintrag mit `liegt in` in dieser Welle, nichts zu prüfen.
  Folge-Slice — keiner genannt, nichts zu prüfen. Register — die in
  dieser Welle neu zitierte `BEO-PGC/verwaltung-keine-sql-administration`
  existiert als Verzeichnis (0 Belege, das ist bereits im eigenen
  `observation.md` unter „Benannt, nicht gezählt" transparent
  ausgewiesen, keine stille Lücke).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; beide Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
