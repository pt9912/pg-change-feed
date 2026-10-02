# Validierungsbeleg: welle-routing — 2026-10-02

**Rolle:** Validator (Modul 8) — „Bauen wir das Richtige?“ gegen den realen Bedarf, aus der Sicht
eines Betreibers/Integrators, der das Handbuch liest und die Beispiele nachvollzieht. Die grünen
Verifikations-Reports sind Eingabe, nicht Ergebnis.

**Gegenstand:** Welle [`welle-routing`](../plan/planning/done/welle-routing.md) mit ihrer Closure-Notiz
[`welle-routing-results`](../plan/planning/done/welle-routing-results.md) (Abschnitt „Validator-Feststellung
(Modul 8)“: Lauf nicht gelaufen, Frist vor dem Release-Entscheid). Bedarf:
[`LH-FA-CFG-008`](../../spec/lastenheft.md) in Lastenheft 0.14.0 (Happy Path, Boundary, Negative,
Out-of-Scope), Stakeholder „Datenbankadministratoren / DevOps-Teams (Betreiber)“ und „ETL-/Integrations-Entwickler
(Consumer)“ ([`spec/lastenheft.md`](../../spec/lastenheft.md) §2). Entscheidungen:
[`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
[`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md),
[`ADR-0141`](../plan/adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md). Belege der Eingabe:
[`verifikation-slice-routing-e2e`](verifikation-slice-routing-e2e.md),
[`verifikation-slice-routing-sdk-realserver-e2e`](verifikation-slice-routing-sdk-realserver-e2e.md),
[`docs/user/e2e-abdeckung.md`](../user/e2e-abdeckung.md), [`docs/user/sdk-e2e-abdeckung.md`](../user/sdk-e2e-abdeckung.md).

Dieser Lauf ändert weder Code noch Plan noch Spec; er schreibt nur diesen Beleg. Ursprung der Zahlen und
Ausgaben: **gemessen** = in diesem Lauf gefahren und unten gedruckt; **übernommen** = aus E2E-Läufen und
Reports der Eingabe, nicht nachgefahren; **hergeleitet** = aus Handbuch oder Entscheidung gelesen.

## 0. Verdikt

**Bedarfsgerecht: ja, mit Einschränkungen (ja, nicht „teils“).** Ein Betreiber kommt mit dem Handbuch allein —
ohne Quellcode-Kenntnis — von der Regel bis zur Auswahl nach Ziel; das Handbuch-Beispiel reproduziert sich
Zeile für Zeile, die Fehlertexte und die Abhilfe der Nicht-Anwendbarkeit stimmen mit der gesehenen Ausgabe
überein. Die Einschränkungen sind Bedarfskanten (Abschnitt 3), keine verfehlten Akzeptanzkriterien.

**Release-Blocker: nein** (Begründung Abschnitt 5).

## 1. Die drei Akzeptanzkriterien in der Lesart „Zustellen = unter dem Ziel auswählbar“

Aufbau (**gemessen**): Wegwerf-Umgebung im Scratchpad, PostgreSQL 18 (`postgres:18-alpine`, Digest der
Wurzel-Compose-Datei), NATS, Feed-Container `ghcr.io/pt9912/pg-change-feed:dev`; vorab
`git diff d8371372 HEAD --stat -- internal cmd proto gen` druckt nichts (Image-Stand = Code-Stand). Schema-Rollout
über `tools/schema/apply-rollout.sh` (Exit 0), Feed als Superuser, `op_admin` mit `cdc_admin`-Mitgliedschaft,
`op_reader` mit `cdc_reader`-Mitgliedschaft; Tabelle `public.orders(id, name, region)`, die Zeilen 1 bis 3
**nach** dem Feed-Start und **vor** den Regeln eingefügt, wie im Handbuch-Beispiel.

| Kriterium | Probe | Gesehen | Urteil |
|---|---|---|---|
| Happy Path | zwei Regeln nach dem Handbuch-Beispiel (`eu_orders` order 10 mit `when`, `rest` order 100 ohne `when`), dann Zeilen 4 bis 6 | beide Anträge `applied`; `cdc.changes` druckt `1,2,3` mit leerem Ziel, `4 → eu`, `5 → sonstige`, `6 (region NULL) → sonstige`; `WHERE route_target='eu'` liefert `4`; `route_target IS NULL` zählt `3`; `GET /changes?source=meine-quelle&target=eu` liefert genau die Change `824-1`, `target=sonstige` genau `824-2` und `824-3`; SSE `GET /changes/stream?target=eu` liefert beim Einfügen von (7, eu) und (8, us) genau das Event zu 7 | erfüllt, identisch zur Tabelle im Handbuch |
| Boundary | Fehlversuche gegen das Handbuch (R1 bis R6, Formprüfungen) | `order bereits vergeben: public.orders.10`, `Zielname ist ungültig: public.orders.A`, `Regelname nicht geführt: public.orders.nope`, `unbekannter Schlüssel in rule_spec: foo` — alle `failed`, Regelstand unverändert; zwei treffende Regeln löst die kleinere `order` (Zeile 4 trägt `eu`, nicht `sonstige`) | erfüllt (Auflösung durch Reihenfolge und statischen Ausschluss, nicht still) |
| Negative | Tabelle `public.rt2` per `cdc.enable_table` aktiviert, Regel `rt_eu` auf `region` (`applied`), dann `ALTER TABLE … DROP COLUMN region`, erste Change eingefügt | Prozess endet mit Exit 1, Container `exited`/`unhealthy`; Log: `Fehlerklasse schema: Routing-Regel auf die Änderung nicht anwendbar: Regel "rt_eu" an public.rt2: Spalte der Routing-Regel fehlt in den Spalten der Änderung: region`; `cdc.heartbeat.error_class = schema`; keine Change von `public.rt2` persistiert (Zählung `0`) | erfüllt, kein Standardziel, kein stilles Verwerfen |

Abhilfe (Handbuch Ursache 1, **gemessen**): `cdc.remove_route` bei stehendem Prozess → Antrag `pending`;
`docker start` → Container `healthy`, Antrag `applied`; die zuvor nicht bestätigte Change erscheint in
`cdc.changes` als `{"id": "1", "note": "a"}` mit leerem Ziel; eine weitere Zeile wird erfasst (Zählung `2`),
`cdc.heartbeat.error_class` ist wieder leer. Der Wortlaut des Handbuchs (Fehlertext, Reihenfolge der Abhilfe,
Ergebnis) deckt sich vollständig mit dem Gesehenen.

Gedruckte Auszüge:

```text
 id | region | route_target
 1  | eu     |
 2  | us     |
 3  |        |
 4  | eu     | eu
 5  | us     | sonstige
 6  |        | sonstige
```

```text
pg-change-feed: Fehlerklasse schema: Routing-Regel auf die Änderung nicht anwendbar: Regel "rt_eu" an public.rt2: Spalte der Routing-Regel fehlt in den Spalten der Änderung: region
```

Nicht nachgefahren (**übernommen** aus dem E2E-Lauf, den beiden Verifikations-Reports und den SDK-Tiers):
gRPC-Stream, NATS-Zusatz-Subjekt, `ReadChanges`, Backfill-Label, Neustart-Ableitung, PostgreSQL 17, die drei
SDKs gegen den Realserver. Die Abdeckungstabellen weisen `LH-FA-CFG-008` für `E2E` und `SDK-E2E` aus
(Zeilen zur Routing-Phase in [`e2e-abdeckung.md`](../user/e2e-abdeckung.md) und je Sprache in
[`sdk-e2e-abdeckung.md`](../user/sdk-e2e-abdeckung.md)); ich habe diese Zeilen gelesen, nicht erneut gelaufen.

## 2. Verständlichkeit für den Betreiber (Fehlertexte, Abhilfe)

- **Gut:** Fehlertexte nennen die Adresse (`schema.tabelle.name`); der Status-Poll über
  `cdc.administration_request` ist im Handbuch als Schritt 3 vorgesehen und nötig, weil `set_route` auch bei
  späterem `failed` eine Kennung zurückgibt. Das Beispiel ist ohne Vorwissen ausführbar; die Rollenvoraussetzung
  (`cdc_admin` für `set_route`, `cdc_reader` zum Lesen) ist ausgesprochen.
- **Gesehen, nicht schlimm:** Bei mehreren Fehlern meldet der Antrag nur den ersten in der Spec-Reihenfolge
  (eine Regel mit Tippfehler in der Spalte **und** doppelter `order` meldete `order bereits vergeben`, nicht
  `Spalte existiert nicht`). Das Handbuch nennt „Reihenfolge“ nach der Spec; der Betreiber korrigiert in zwei
  Runden.
- **Lücke in der Anleitung:** Die Abhilfe für Ursache 2 (Bedingungsspalte an einer Tabelle mit bekannter
  Spaltenform entfernt) ist im Handbuch ehrlich als „nicht beschrieben“ ausgewiesen (F-4).

## 3. Befunde und Lücken des Bedarfs

Schweregrade: MEDIUM = Betreiber stößt real an, Abhilfe nur mit Zusatzwissen; LOW = Unbequemlichkeit oder
dokumentierte Grenze; INFO = Einordnung.

| # | Befund | Schwere | Einordnung |
|---|---|---|---|
| F-1 | **„Zustellen“ heißt Auswählbarkeit, kein aktives Schieben an eine Senke.** Wer „unterschiedliche Zustellziele“ als „Webhook/Queue je Ziel“ liest, bekommt das nicht. Lastenheft 0.14.0 hält die Lesart ausdrücklich fest; [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) hat die aktive Senke (Option D) verworfen. Die Klarstellung entstand, nachdem der Entwurf stand (Historienzeile 0.14.0). Die Results-Notiz führt die Lesart (A-3) als „offen, beim Auftraggeber“; ein Beleg der Bestätigung der Lesart (im Unterschied zur Freigabe des Release) liegt mir nicht vor | MEDIUM (Lesart) | **bewusst ausgeklammert** (Out-of-Scope-Zeile des Kriteriums: Modell der Zustellziele ist Architekturfrage). Empfehlung: Auftraggeber bestätigt die Lesart ausdrücklich; die Release-Beschreibung nennt „Auswahl nach Ziel, keine aktive Zustellung“ |
| F-2 | **Regel ersetzen ist nicht atomar** („erst entfernen, dann neu setzen“) und es gibt kein Umetikettieren. Zwischen Entfernen und Setzen erfasste Changes tragen das Ziel des dann geltenden Regelstands (Abschlussregel oder leer) und behalten es. Das Handbuch nennt „kein Umetikettieren“, verbindet es aber nicht mit dem Ersetzen (nicht gefahren, **hergeleitet**) | LOW | Bedarfslücke, die zum Konzept gehört (Ziel ist zum Erfassungszeitpunkt fest, [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)); ein Satz im Handbuch am Ersetzen-Abschnitt würde sie ehrlich machen. Der Altbestand-Weg ist ein neuer Backfill-Run (dokumentiert) |
| F-3 | **`set_route`/`remove_route` gibt es nur über SQL**, nicht über die gRPC-Verwaltungs-API oder HTTP (`proto/cdc/administration/v1/administration.proto` trägt nur den Lesefilter `target`). `EnableTable` ist über die API erreichbar | INFO | Konsistent mit den Transformationsregeln und mit der Erwartung der DBAs („SQL-basierte Administration“, §2 des Lastenhefts); kein Kriterium von [`LH-FA-CFG-008`](../../spec/lastenheft.md) verlangt eine API. Bewusst, keine Lücke |
| F-4 | **Abhilfe der Ursache 2 fehlt** (Spalte an Tabelle mit bekannter Spaltenform entfernt; das Entfernen der Regel genügt nach Handbuch nicht). Ein Betreiber, der eine geroutete Spalte entfernt, findet im Handbuch den Befund, aber keinen Weg | MEDIUM | Die Lücke gehört der inkompatiblen Schemaänderung ([`LH-FA-SCH-004`](../../spec/lastenheft.md)), ist vorbestehend und im Handbuch offen benannt. Für diese Welle kein Blocker; als Betriebsdoku-Folge wertvoll (Vorbeugung: Regel **vor** dem `DROP COLUMN` entfernen) |
| F-5 | **Die Nichtanwendbarkeit hält die gesamte Quelle an**, nicht nur die Tabelle: nach dem Ausfall von `public.rt2` erfasste auch `public.orders` nichts mehr, bis zum Neustart (Container `exited`, `restart: "no"`). Beobachtet, im Handbuch beschrieben („Erfassung der gesamten Quelle“) | LOW | Folge der Fail-closed-Entscheidung („kein stilles Standardziel“, Negative-Kriterium); Betreiber braucht eine Überwachung/einen Neustart-Vertrag (Health-Check meldet `unhealthy`: gesehen) |
| F-6 | **Ein Tippfehler im Ziel beim Lesen liefert leer und keinen Fehler** (`target=Eu`: `{"changes":[]}`; am gRPC-Stream dokumentiert). Der Leser kann „Tippfehler“ nicht von „keine Change“ unterscheiden | LOW | Bewusste Entscheidung ([`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)); Abhilfe des Lesers: `route_target` per SQL prüfen oder ungefiltert lesen |
| F-7 | **Das Ziel steht in keiner Antwort und keiner Stream-Nachricht** der Roh-Wege, nur in `cdc.changes.route_target`; ein Leser über HTTP/gRPC/SSE muss dem Filter vertrauen. Am NATS-Weg ist das Zusatz-Subjekt Fire-and-forget (kein Replay), der SSE-Client des SDK trägt `target` als einzigen Filter | INFO | **bewusst** (Nachrichtenschema bleibt unverändert, [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)); Fire-and-forget gilt für alle Live-Wege dieses Systems |
| F-8 | **Eine Change hat höchstens ein Ziel**; ein Fan-out (dieselbe Change in zwei Zielen) ist nicht möglich | INFO | **bewusst** (Leitplanke aus [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md), Boundary „Auflösung bei Mehrdeutigkeit“); wer zwei Ziele braucht, liest ungefiltert oder vergibt ein Sammelziel |
| F-9 | **Prozess-Drift:** Die Frist „Validator vor dem Release-Entscheid“ wurde verpasst — der Auftraggeber hat die Freigabe vor diesem Lauf erteilt. Der Release selbst ist laut Auftrag noch nicht ausgeführt; dieser Lauf liegt damit noch vor dem Tag, aber nach der Entscheidung | LOW (Prozess) | Gemeldet, nicht durchgewunken. Kein Einfluss auf das Ergebnis; für künftige Wellen: Validator-Auftrag zusammen mit der Closure, nicht danach |

## 4. Grenzen dieser Validierung

- **Kein Last-/Dauer-Test:** Die Lesekosten der Regelstände im Backfill-Run und die Auswertungslast je Change
  sind laut Results-Notiz nicht gemessen; ich habe sie nicht gemessen.
- **Eine Version, ein Pfad:** gefahren an PostgreSQL 18 und am Pfad `GET /changes`, SQL, SSE; PostgreSQL 17,
  gRPC, NATS, Backfill und die SDKs sind **übernommen**.
- **Rolle des Feeds:** Der Feed lief als Superuser (wie die Wurzel-Compose-Umgebung); die Least-Privilege-Login-Identitäten
  des Betriebs für den Feed selbst habe ich nicht aufgebaut. `set_route` und die Lesezugriffe liefen unter
  Rollen mit `cdc_admin`- bzw. `cdc_reader`-Mitgliedschaft.
- **Beteiligung am Arbeitsbaum:** Während meines Laufs änderten sich im Arbeitsbaum Dateien, die nicht von mir
  stammen (`.claude/commands/*`, `AGENTS.md`, `.harness/skills/closure-note-reviewer.md`,
  `tools/harness/run-schema-rollout-guard-test.sh`, `tools/schema/plan.yaml`, `tools/schema/down.sql`; Letztere
  tragen das Ziel `cdc-schema-rollout-guard-test-pg`, also einen parallelen Lauf des Guard-Tests). Ich habe sie
  weder angefasst noch committet.

## 5. Release-Reife aus Bedarfssicht

**Blockiert ein Befund den Release (Server v0.5.0, SDKs mit neuem optionalem Parameter `target`)? Nein.**

- Alle drei Akzeptanzkriterien sind am laufenden System aus Betreibersicht erfüllt, einschließlich der
  Abhilfe der Nichtanwendbarkeit; Handbuch und gesehene Ausgabe stimmen überein.
- Kein Befund ist eine verfehlte Zusage: F-1 und F-3 bis F-8 sind in Lastenheft/ADR als Abgrenzung angelegt
  oder im Handbuch benannt; F-2, F-4 und F-6 sind Betriebsdoku-Verfeinerungen; F-9 ist Prozess.
- **Vor dem Tag empfohlen (nicht blockierend):** (a) Auftraggeber bestätigt die Lesart aus F-1 ausdrücklich;
  (b) Release-Beschreibung nennt „Auswahl nach Ziel, keine aktive Zustellung“ und „Ziel nur zum Erfassungszeitpunkt,
  kein Umetikettieren“; (c) die C#-/Kotlin-SDK-Release-Notiz nennt, dass `target` als neuer letzter Parameter
  quellkompatibel, aber nicht binärkompatibel ist (die Results-Notiz hält das fest, der Bedarf eines bestehenden
  C#-Konsumenten ist „neu übersetzen“).
- **Folge-Kandidaten** (nach dem Release): Hinweis am Ersetzen-Abschnitt (F-2), Betriebsanleitung zur
  inkompatiblen Schemaänderung (F-4).

## 6. Aufräumen

Wegwerf-Umgebung abgebaut mit `docker compose … down -v` (Container und Netz `valroute-net` entfernt;
`docker ps -a` zeigt keinen `valroute`-Container). Dieser Beleg ist die einzige Datei, die dieser Lauf im
Repo anlegt.
