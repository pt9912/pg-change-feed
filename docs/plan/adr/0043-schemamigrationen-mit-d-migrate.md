# ADR-0043: Schemamigrationen mit d-migrate

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`ADR-0015`](0015-schema-evolution.md) (Schema-Evolution — diese
Entscheidung gibt der Schema-Form ihre Werkzeug-Abdeckung),
[`ADR-0032`](0032-postgresql-adapterdetail.md) (PostgreSQL bleibt
Adapterdetail)

**Schärft:** [`SPEC-001`](../../../spec/pflichtenheft.md),
[`SPEC-002`](../../../spec/pflichtenheft.md) (Schema-Form der CDC-Tabellen —
[Pflichtenheft §2](../../../spec/pflichtenheft.md)),
[`ARC-006`](../../../spec/architecture.md),
[`ARC-009`](../../../spec/architecture.md) (Driven-Adapter-Infrastruktur) —
die Herkunft der DDL des Store-Adapters. Die Domänen-Invarianten
([`ADR-0029`](0029-domain-invarianten.md), darunter Regel 1
Persist-before-ACK) und die Idempotenz des Store-Schreibpfads
([`ADR-0011`](0011-persist-before-ack.md)) bleiben Port-Kontrakte,
unberührt: Die Entscheidung ändert die Quelle der DDL, nicht das
Verhalten des Adapters. Aufwärts-Deklaration der Änderungskopplung: wer
diese ADR ändert, zieht die Schema-Form und die Adapter-Infrastruktur
nach.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Das CDC-Schema — die Tabellen aus [`SPEC-001`](../../../spec/pflichtenheft.md),
die Change-Struktur aus [`SPEC-002`](../../../spec/pflichtenheft.md) — wird
derzeit handgeschrieben als DDL geführt
(`internal/adapters/driven/postgresstorage/schema.sql`, Schema-Aufbau des
PostgreSQL-Store-Adapters) und über `ApplySchema` idempotent
(`IF NOT EXISTS`) in den Testcontainer gerollt. Diese Form trägt den
Erststand, aber drei Grenzen wachsen mit dem Schema:

- Die Idempotenz-Form sichert nur die *Erst-Anlage*; eine Änderung an
  einem bestehenden Objekt (Spalte, Constraint, Index) ist in ihr nicht
  ausdrückbar — der Stand driftet still, sobald ein Ausbau bestehende
  Objekte berührt.
- Ein Rollout ist unbeobachtet: es gibt kein Rücknahme-Artefakt und
  keinen Beleg darüber, was ein Lauf getan hat.
- Die Schema-Form existiert nur als SQL je Dialekt; ein zweiter Dialekt
  hieße ein zweites handgeschriebenes DDL-Dokument für denselben Stand.

Zugleich steht der Ausbau an: die Tabellen der Consumer- und
Betriebs-Zustände (`cdc.consumer`, `cdc.consumer_position`,
`cdc.capture_state`, [`SPEC-001`](../../../spec/pflichtenheft.md)) trägt die
DDL noch nicht — sie wächst.

Der Auftraggeber hat entschieden: **d-migrate**
(`ghcr.io/pt9912/d-migrate`, Container-Werkzeug mit neutralem
YAML-Schemamodell, das je Dialekt SQL erzeugt) wird das
Schemamigrations-Werkzeug dieses Repos. Die Herleitung steht im
Anwenderhandbuch des Werkzeugs (`docs/user/anwenderhandbuch.md`,
§2.3 Schnelldurchlauf, §3.1 Schema prüfen, §3.2 SQL erzeugen, §3.5
Ausrollen und Zurücknehmen). Die Annahme, die kippen würde: dass
d-migrate die benötigten Operationen ausdrücken kann
(Re-Evaluierungs-Trigger unten). Die Decken-Regel trennt hier die
Schichten: Diese ADR nennt die Einsatzmomente; die Planning-Dateien
tragen die Slice-Kennungen und beziehen sich von dort auf diese ADR.

## Entscheidung

Wir wählen **d-migrate als Schemamigrations-Werkzeug**:

1. **Das neutrale Schemamodell ist die Quelle.** Das Repo führt
   `tools/schema/schema.yaml` im neutralen d-migrate-Schemaformat; es
   beschreibt die CDC-Tabellen aus [`SPEC-001`](../../../spec/pflichtenheft.md)/[`SPEC-002`](../../../spec/pflichtenheft.md)
   technologieunabhängig. Die DDL ist von dort erzeugt, nicht mehr die
   Quelle selbst.
2. **SQL wird erzeugt, nicht handgeschrieben.** Vor jedem
   `generate`/`migrate` läuft `schema validate --source schema.yaml`
   als Vorlauf (netzlos); das Ziel-SQL entsteht mit
   `schema generate --source schema.yaml --target postgresql`.
3. **Rollout über `schema migrate` mit Pflicht-Report.** Der Rollout
   auf eine Datenbank läuft als
   `schema migrate --source schema.yaml --target db:… --execute --report plan.yaml --generate-rollback --rollback-output down.sql`.
   Der Report ist Pflicht und ist der Beleg des Rollouts (Operationen,
   Risiken, Rollback-Fähigkeit). Destruktive Operationen bleiben
   default blockiert (Exit 8); ihre Zulassung ist ein bewusster,
   berichteter Entschluss, kein Default. Die Rollback-Artefakte
   (`plan.yaml`, `down.sql`) werden je Rollout aufbewahrt.
4. **Ersteinsatz ist der Compose-Rollout.** Der erste Einsatz ist der
   Schema-Rollout in die Compose-Testdatenbank vor jedem E2E-Lauf; der
   Ausbau-Slice um die ACK-/Consumer-Tabellen nutzt das Werkzeug nach
   Bedarf (beide Adressen im Planning-Stratum).
5. **Benannte Grenze — der Testloader bleibt.** Der handgeschriebene
   Test-Schema-Loader des Schema-Aufbau-Slices (DDL + `ApplySchema`)
   bleibt in den Testcontainer-Tests stehen, bis der d-migrate-Rollout
   ihn ersetzt. Er ist Bestand mit benanntem Ersatz, kein Vorbild für
   neue Tabellen — neue Tabellen entgehen der DDL-Herleitung.
6. **Port-Kontrakte unberührt.** Die Domänen-Invarianten
   ([`ADR-0029`](0029-domain-invarianten.md)) und die Idempotenz des
   Store-Schreibpfads ([`ADR-0011`](0011-persist-before-ack.md)) sind
   Kontrakte am Port; die Werkzeug-Herkunft der DDL ändert nichts an
   ihrem Verhalten.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Hand-DDL + Loader (Status quo) | kein zweites Werkzeug-Image; netzlos; die Idempotenz-Form hat sich im Testcontainer bewährt | Änderungen an bestehenden Objekten sind in der Idempotenz-Form nicht ausdrückbar; kein Rollback-Artefakt, kein Rollout-Beleg; je weiterem Dialekt ein weiteres Hand-DDL für denselben Stand |
| B — d-migrate als Export-Quelle für Flyway/Liquibase (`export flyway`/`liquibase`, Anwenderhandbuch §3.11) | etablierte Migrations-Läufer mit Versionstabelle; Datei-Artefakte im Repo | zweites Werkzeug-Bild (Export-Lauf + Fremdwerkzeug-Lauf) und ein externer Lauf-Zustand (Versions-Tabelle) im CDC-Schema — für einen Soll-Stand-Abgleich ohne Migrationshistorie Zustands-Doppelbuchung |
| C — nichts tun; DDL weiter handführen | kein neuer Pin, kein neuer Workflow | die drei Kontext-Grenzen wachsen mit dem Schema; der anstehende Ausbau (Consumer-/ACK-Tabellen) wiederholt das Muster |
| **D — d-migrate direkt (gewählt)** | ein Image (Digest-Pin); neutrales Modell als eine Quelle; Pflicht-Report als Rollout-Beleg; Rollback-Artefakt je Lauf; destruktive Operationen default blockiert; `--deterministic` für diff-freundliche DDL | ein weiteres Werkzeug-Image mit Pin-Pflicht (Modul 14 — Pin-Hebung = bewusster Commit); `schema migrate --execute` braucht DB-Zugang und ist damit nicht netzlos-Gate-tauglich; die Überführung DDL → YAML ist einmalige Nacharbeit |

## Konsequenzen

- Positiv: Die Schema-Form hat eine Quelle (`schema.yaml`) statt
  doppelt geführter DDL; SQL wird abgeleitet; jeder Rollout trägt
  Pflicht-Report (Beleg) und Rollback-Artefakt (Rücknahme); der
  Schutz gegen Datenverlust ist Werkzeug-Verhalten (Exit 8), keine
  Konvention, die ein Lauf vergessen kann.
- Negativ: Ein zweites Werkzeug-Image im Baum — der Digest-Pin folgt
  der Modul-14-Disziplin, die Pin-Hebung ist ein bewusster Commit;
  die Rollout-Targets sind nicht netzlos (DB-Zugang) und hängen an
  keinem Gate.
- Folgepflicht: **Das neutrale Schema-YAML entsteht als Erstlieferung
  des ADR-Einbaus** — die Überführungs-Quelle ist die handgeschriebene
  DDL (`internal/adapters/driven/postgresstorage/schema.sql`, Stand des
  Schema-Aufbaus), nicht eine Neuerfindung. Vor dieser Überführung
  bleibt die DDL die Quelle; Abweichung zwischen DDL und YAML wäre
  zwei Quellen für dieselbe Schema-Form und ein Befund.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| d-migrate (gepinntes Image) | `schema validate` als Vorlauf vor jedem `generate`/`migrate`; der Pflicht-Report jedes `--execute`-Laufs ist der Beleg und wird aufbewahrt | `make schema-validate` (Vorlauf, netzlos, kein Gate) · `make schema-rollout` (braucht DB-Zugang, kein Gate) |

Keines der beiden Targets hängt an `GATE_CHECKS` — das Werkzeug
erzeugt Belege, es bewacht keine Schwelle; die Gates dieses Repos
bleiben unverändert.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbarer Trigger: **d-migrate kann eine benötigte Operation nicht
ausdrücken** — der Lauf blockiert (Exit 8) oder meldet einen Blocker,
und es gibt keine Ausweichform (Overlay, `--rename-*`, berichtete
manuelle Nacharbeit) — dann ist der Werkzeugwechsel als Folge-ADR mit
`supersedes` zu prüfen. Sonst `permanent` — das Werkzeug wächst mit
dem Schema.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted — Anlass: handgeschriebene DDL des Schema-Aufbaus und anstehender Ausbau um Consumer-/ACK-Tabellen; Vorgabe pt9912; Herleitung aus dem Anwenderhandbuch des Werkzeugs (§2.3, §3.1, §3.2, §3.5) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).