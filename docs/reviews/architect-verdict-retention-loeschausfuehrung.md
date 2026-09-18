# Architect-Verdikt: Retention-Löschausführung — Entscheidungsgrundlage vor dem Schneiden

**Rolle:** Architect (Modul 8)
**Anlass:** Anfrage des Planners vor Eröffnung/Schnitt der Welle
„Retention-Löschausführung" (`docs/plan/planning/in-progress/roadmap.md`,
Trigger `welle-12` liegt in `done/`), Größe **L**
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** [`LH-FA-RET-002`](../../spec/lastenheft.md)…[`006`](../../spec/lastenheft.md),
[`LH-QA-OPS-003`](../../spec/lastenheft.md),
[`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md),
[`ADR-0011`](../plan/adr/0011-persist-before-ack.md),
[`ADR-0012`](../plan/adr/0012-at-least-once.md),
[`ADR-0014`](../plan/adr/0014-retention-domain-policy.md),
[`ADR-0029`](../plan/adr/0029-domain-invarianten.md),
[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
`docs/plan/planning/observations/BEO-PGC/retention-keine-loeschausfuehrung`

---

## Verdikt

**Beide Fragen: bestehende Entscheidungen/Muster tragen bereits — keine neue
ADR.** Dieselbe Konstellation wie beim `slice-030`/`ADR-0015`-Präzedenzfall
(der Architect-Verdikt zu `slice-030`/`ADR-0015`): Eine
`Accepted`/`permanent` ADR benennt die Fähigkeit bereits als
**Folgepflicht**, die schlicht noch nicht umgesetzt wurde. Es ist keine
Architektur-*Entscheidung* zu treffen, sondern eine bereits getroffene
Entscheidung *einzulösen*.

## Frage 1 — Löschausführung (neue `ChangeStorePort`-Delete-Methode + Use-Case)

**Trägt bereits, ohne Folge-ADR — [`ADR-0014`](../plan/adr/0014-retention-domain-policy.md) benennt die fehlende Umsetzung explizit als Folgepflicht.**

- **`ADR-0014`** entscheidet bereits die Rollenteilung: „Der Domain Core
  entscheidet, was sicher gelöscht werden darf; der Storage-Adapter führt
  die physische Löschung aus" — und nennt als Folgepflicht wörtlich:
  „RunRetentionUseCase unter den Inbound Use Cases; blockierende Consumer
  sind sichtbar zu melden (`LH-FA-RET-005`)." Re-Evaluierungs-Trigger:
  `permanent` — die Trennung von Löschentscheidung und physischer Ausführung
  ist an keine Bedingung geknüpft, die sich geändert hätte.
- **`ADR-0009`** (Option C, Accepted) hat den `ChangeStorePort` bewusst als
  wachsenden Fähigkeits-Port entworfen — die dort vermerkte Contra „Port-
  Vertrag muss Lesen und Schreiben gleichermaßen gut abdecken und wird damit
  breiter" ist exakt die hier anstehende Erweiterung um eine Delete-Fähigkeit,
  keine neue Grundsatzfrage. `re-Evaluierungs-Trigger: permanent`.
- **Kein Zielkonflikt mit Persist-before-ACK/At-Least-Once**, geprüft am
  Code: `RetentionPolicy.AllowsDeletion` (`internal/domain/model/retention.go`)
  verweigert die Freigabe, solange irgendeine übergebene Consumer-Position
  nicht bestätigt ist oder vor der Change-Position liegt
  (`!position.Acknowledged() || position.Position.Before(changePosition)`).
  Das ist bereits `ADR-0029` Invariante 5 („Retention löscht keine benötigten
  Changes") als Code-Konstruktion. Eine Löschausführung, die diese Policy als
  einzige Freigabe-Instanz respektiert, kann eine Zeile physisch immer erst
  löschen, *nachdem* alle betrachteten Consumer sie bereits bestätigt haben —
  das liegt zeitlich hinter dem Persist-before-ACK-Pfad (`ADR-0011`: Persist
  → COMMIT Store → ACK Source) und hinter der At-Least-Once-Zusage
  (`ADR-0012`), nicht in Konkurrenz dazu. Es gibt keine Sequenz, in der die
  physische Löschung eine dieser beiden Garantien unterlaufen könnte, solange
  der Use-Case ausschließlich über `AllowsDeletion` freigibt (keine
  Bypass-Löschung am Domain Core vorbei) — genau das ist der Kern von
  `ADR-0014`s Entscheidung und bereits geschrieben.
- **`ARC-002`** (`spec/architecture.md`) führt „Retention" bereits als
  Application-Use-Case-Kategorie neben Capture/Consumer/Konfiguration — die
  Sicht behauptet die Fähigkeit als vorgesehen, nicht als geliefert (dieselbe
  Lesart wie bei `ADR-0015`/`SchemaStorePort` im Präzedenzfall).
- **`LH-FA-RET-005`** (Erkennbarkeit blockierender Consumer) ist bereits zu
  wesentlichen Teilen durch bestehende Infrastruktur getragen: die View
  `cdc.consumer_status` (`tools/schema/schema.yaml`) projiziert je Consumer
  die Differenz aus bestätigter und letzter Commit-Position — genau die
  Information, die einen blockierenden Consumer einzeln erkennbar macht.
  Der Use-Case muss diese Sichtbarkeit für den Retention-Kontext nicht neu
  erfinden, nur (ggf.) im Löschlauf selbst mitführen/loggen.

**Fazit F1:** Eine neue `ChangeStorePort`-Delete-Methode und ein
`RunRetentionUseCase`, der `AllowsDeletion` real aufruft, ist reine
Umsetzung einer bereits `Accepted`/`permanent` stehenden Entscheidung.
Keine der geprüften ADRs (`0009`, `0011`, `0012`, `0014`, `0029`) steht der
Umsetzung entgegen; keine widerspricht sich; kein Zielkonflikt zwischen
Löschausführung und At-Least-Once-Garantie besteht, weil die Domain-Policy
diesen Konflikt bereits strukturell ausschließt.

## Frage 2 — `cdc_storage_bytes`-Metrik (Rollen-Erweiterung nötig?)

**Trägt bereits über das etablierte View-Owner-Muster — keine
Rollen-Erweiterung für `cdc_reader` nötig, keine neue ADR.**

- **Das Muster ist bereits zweimal etabliert und dokumentiert:**
  `cdc.metrics` (`tools/schema/nacharbeit-observability.sql`) und
  `cdc.heartbeat` (`tools/schema/nacharbeit-heartbeat.sql`) tragen exakt
  dieselbe Aussage im Kommentar: PostgreSQL-Views laufen mit den Rechten
  ihres Eigentümers (Definer-Semantik ohne `security_invoker`); `GRANT SELECT
  ... TO cdc_reader` auf die View allein trägt den Lesezugriff, ohne dass
  `cdc_reader` je einen Grant auf eine zugrundeliegende Basistabelle bekommt
  (`nacharbeit-roles.sql`, Kommentar zu `cdc_reader`). Eine dritte View nach
  demselben Muster ist keine neue Entscheidung, sondern eine weitere
  Instanz derselben.
- **`pg_relation_size()` selbst verlangt keine Sonderrechte:** Es ist eine
  reguläre, für `PUBLIC` ausführbare Systemfunktion (kein
  `SECURITY DEFINER` nötig, kein `REVOKE EXECUTE FROM PUBLIC` in
  PostgreSQL) und liefert reine Metadaten (Relationsgröße), keinen
  Zeileninhalt — sie prüft kein `SELECT`-Privileg auf der Zieltabelle. Die
  einzige Voraussetzung ist, dass die aufrufende Rolle die Relation über
  ihren qualifizierten Namen auflösen kann (`USAGE` auf dem Schema); das ist
  für `cdc_reader` bereits erteilt (`GRANT USAGE ON SCHEMA cdc TO ...
  cdc_reader;`, `nacharbeit-roles.sql`).
- **Der eigentliche Least-Privilege-Punkt, den der bisherige Kommentar in
  `nacharbeit-observability.sql` zu Recht benennt, ist nicht
  `pg_relation_size()` selbst, sondern ein *direkter* Grant an `cdc_reader`,
  der ihm erlauben würde, `pg_relation_size()` auf einer **beliebigen**
  Relation der Datenbank zu parametrisieren** — das wäre tatsächlich eine
  Ausweitung der Angriffs-/Informationsfläche. Genau das vermeidet das
  View-Owner-Muster strukturell: Die View hardcodet `'cdc.change'` in ihrer
  `SELECT`-Klausel; `cdc_reader` bekommt nur `SELECT` auf die fertige View
  und kann den Parameter nicht selbst wählen. Die Sorge aus dem
  Ursprungskommentar bezog sich der Sache nach auf einen direkten,
  parametrisierbaren Zugriff — nicht auf eine gekapselte View, wie sie
  `cdc.metrics`/`cdc.heartbeat` bereits vorführen.
- **`ADR-0046`** (Accepted, supersedes `ADR-0018`) entscheidet bereits
  genau diese Klasse: „Reine Lese-SQL-Views dürfen Driven-Adapter-Tabellen
  direkt lesen, sofern sie ausschließlich Projektion/Join über bereits
  persistierte, bereits validierte Daten liefern und keine Domänenregel
  auswerten." Eine `cdc_storage_bytes`-View, die `pg_relation_size('cdc.change')`
  projiziert, trifft keine Domänenentscheidung (keine Schwellenwert-Logik,
  keine Autorisierungsprüfung) — sie ist Projektion im Sinne dieser ADR,
  Kategorie C, wie `cdc.metrics` es für die übrigen `SPEC-009`-Zeilen bereits
  ist.
- **Rollen-Erweiterung als Präzedenz, nicht als neue Frage:**
  `nacharbeit-roles.sql` zeigt bereits mehrfach denselben Nacharbeitsschritt
  („Lückenschließung") für neu hinzugekommene Objekte (Heartbeat-Grants,
  Schema-Version-Grants) — ohne dass dafür je eine ADR nötig war. Ein
  zusätzlicher, sich selbst grantender View-Block (wie `cdc.metrics`/
  `cdc.heartbeat`) reiht sich in dieselbe etablierte Form ein.

**Fazit F2:** Die `cdc_storage_bytes`-Metrik lässt sich als vierte
View im etablierten Muster liefern (View-Eigentümer ruft
`pg_relation_size('cdc.change')` auf, `GRANT SELECT` auf die View an
`cdc_reader`), ohne `cdc_reader` einen direkten Systemkatalog- oder
Funktionsaufruf-Zugriff zu geben. Der in `nacharbeit-observability.sql`
vermerkte Least-Privilege-Vorbehalt betraf einen direkten, parametrisierbaren
Zugriff — nicht das hier anwendbare, bereits zweifach bewährte
View-Owner-Muster. Keine neue ADR nötig; `ADR-0046` trägt die Einordnung
bereits.

## Was das für den Planner bedeutet

Beide Fähigkeiten sind **Umsetzung**, keine **Entscheidung** — der
Slice-Plan referenziert `ADR-0009`, `ADR-0011`, `ADR-0012`, `ADR-0014`,
`ADR-0029` und `ADR-0046` als bereits geltende Constraints (Implementer
liest sie als Constraint, Modul 8 §Rollen-Regeln), ohne dass eine dieser
ADRs geändert oder eine neue geschrieben wird. Die Größenordnung **L** aus
der Roadmap bleibt unverändert plausibel (mehrere Schichten: Domain-Port-
Erweiterung, neuer Inbound-Use-Case, Adapter-Implementierung der Löschung,
SQL-Nacharbeit für die Metrik) — das ist ausschließlich Schnitt-Aufgabe des
Planners (Modul 5), nicht Gegenstand dieses Verdikts.

Weder Produktionscode noch eine ADR-Datei wurden im Rahmen dieses Verdikts
geändert.
