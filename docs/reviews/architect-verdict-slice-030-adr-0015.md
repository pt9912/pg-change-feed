# Architect-Verdikt: slice-030 F-1 gegen `ADR-0015`

**Rolle:** Architect (Modul 8 §Konflikt-Pfad als Rollen-Sequenz)
**Anlass:** `docs/reviews/review-slice-030.md`, Finding F-1 (HIGH) — behaupteter
Verstoß gegen [`ADR-0015`](../plan/adr/0015-schema-evolution.md) (Accepted,
`permanent`)
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-12
**Bezug:** [`LH-FA-SCH-004`](../../spec/lastenheft.md),
[`LH-FA-SCH-005`](../../spec/lastenheft.md),
[`SPEC-004`](../../spec/pflichtenheft.md),
[`ADR-0015`](../plan/adr/0015-schema-evolution.md),
`docs/plan/planning/in-progress/slice-030-black-box-e2e-schema-aenderungen.md`

---

## Verdikt

**Verdikt 1 — ADR gilt; kein Plan hat sie fälschlich als umgesetzt
behauptet.** Die in `ADR-0015` (Option C) beschlossene und als Folgepflicht
benannte Fähigkeit — `TableSchema`-/`SchemaVersion`-Modelle je Change,
`SchemaStorePort` als Outbound Port, Fehlerklasse `schema` für nicht sicher
interpretierbare Schemaänderungen — wurde **nie gebaut**. Das ist keine
Diskrepanz zwischen Plan und Code, die ein früherer Slice erzeugt hätte: kein
Slice und keine Welle-Ergebnisnotiz in `docs/plan/planning/done/` behauptet
je, `SchemaStorePort` oder die dynamische Re-Versionierung sei geliefert. Die
Folgepflicht wurde schlicht nie als eigene Arbeit eingeplant — `welle-9`
(die `slice-030` trägt) war explizit als **E2E-Abdeckung bestehender
Fähigkeiten** geschnitten, nicht als Fähigkeits-Lieferung; slice-030 §1
schließt „Änderung des Schema-Änderungs-Erkennungsmechanismus selbst"
korrekt und bewusst aus. Der Fund ist damit real und die HIGH-Einstufung des
Reviewers **sachlich korrekt** — aber die Auflösung ist eine fehlende
Umsetzung nachzuliefern, nicht die ADR zu korrigieren oder zurückzunehmen.

Verdikt 2 (Folge-ADR, `supersedes`) und Verdikt 3 (legitime, aber
undokumentierte Lockerung) wurden geprüft und **verworfen** — Begründung
unten.

## Begründung

### Ist `ADR-0015`s Kontext heute noch gültig?

Ja, ungebrochen:

- **Lastenheft/Pflichtenheft unverändert.** `LH-FA-SCH-004` (inkompatible
  Typänderungen erkennbar melden, keine stille Fehlinterpretation) und
  `LH-FA-SCH-005` (Changes einer Schema-Version zuordenbar,
  unterscheidbar) stehen im aktuellen `spec/lastenheft.md` wortgleich zum
  ADR-Kontext. `SPEC-004` (`spec/pflichtenheft.md`) fordert weiterhin ein
  technologieunabhängiges `TableSchema`/`SchemaVersion`-Modell.
- **Architektur-Sicht führt die Fähigkeit weiterhin als vorgesehen.**
  `spec/architecture.md` zeigt im Sequenzdiagramm zu `LH-FA-CFG-001.a`
  explizit `SchemaStorePort (ARC-004)` zwischen `EnableTableUseCase` und
  `PostgresMetadataAdapter` — die Sicht behauptet die Fähigkeit nicht als
  vorhanden, aber sie ist als Ziel-Architektur weiterhin eingezeichnet, nicht
  gestrichen.
- **Drei spätere ADRs bauen auf der Port-Existenz auf**, ohne sie in Frage zu
  stellen: [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md)
  („vorgesehen: … SchemaStorePort …"),
  [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)
  (Paketstruktur nennt `SchemaStorePort` unter `outbound/`),
  [`ADR-0040`](../plan/adr/0040-clockport.md) (zählt `SchemaStorePort` als
  bestehenden Outbound-Port-Kandidaten auf). Keine dieser ADRs deutet an,
  dass die Fähigkeit verworfen wurde — sie setzen ihre künftige Existenz
  voraus.
- **Der `permanent`-Re-Evaluierungs-Trigger trägt weiterhin**: die
  Core-Unabhängigkeit (Abhängigkeitsregel, Architektur-Sicht §2) ist keine
  Randbedingung, die sich geändert hätte. Kein Ereignis im Repo entspricht
  einer beobachtbaren Bedingung, die eine Neubewertung auslösen würde.

Damit entfällt Verdikt 2: Es gibt keinen sachlichen Grund, Option A heute für
architektonisch richtiger zu halten als bei Verabschiedung von `ADR-0015`. Die
in `ADR-0015` gegen Option A vorgebrachten Contras (PostgreSQL-Details sickern
in den Core; keine stabile historische Interpretation) sind exakt die
Symptome, die der heutige Code zeigt — das bestätigt `ADR-0015`s Abwägung,
statt sie zu widerlegen.

### Gibt es eine legitime, aber undokumentierte Lockerung?

Nein. Die vier Code-Kommentare, die der Reviewer als „Metadata-Pfad"-Zeiger
identifiziert hat (`mapper.go:53`, `wiring.go:331`, `receive.go:62`,
`verwaltung.go:26`), sind selbst **Vorwärtsverweise auf die von `ADR-0015`
geforderte Fähigkeit** — sie behaupten nicht, dass ein Ersatzmechanismus
bewusst gewählt wurde, sondern benennen die Stelle, an der die noch fehlende
Übersetzung greifen soll. Das ist keine dokumentierte Abweichung von Option
C, sondern ihre (bislang unerfüllte) Ankündigung. Keine `MR-*`-Adaption, kein
Carveout und keine Konventions-Notiz im Repo erwähnt einen bewussten Verzicht
auf `SchemaStorePort` oder dynamische Re-Versionierung. Verdikt 3 entfällt.

### War die HIGH-Einstufung des Reviewers korrekt?

Ja. Der Reviewer prüft explizit gegen ADR (Modul 1 Kernidee, Modul 8
Kernidee), und die Kategorie „ADR-Verstoß" trifft strukturell zu: Der
aktuell ausgelieferte Code verhält sich bei Relation-Metadata-Änderungen wie
das in `ADR-0015` **explizit verworfene** Option A (Rohdaten-Durchreichen ohne
stabile historische Interpretation), nicht wie das beschlossene Option C.
Zwei Lastenheft-Akzeptanzkriterien (`LH-FA-SCH-004` Negative-Fall,
`LH-FA-SCH-005` Boundary) sind strukturell unerfüllbar, nicht nur zufällig
ungetestet — die neuen Black-Box-Tests aus slice-030 beweisen das erstmals
empirisch gegen den echten Compose-Stack. Das ist mehr als „Spec-Rückstand":
Es ist der reale Nachweis, dass eine `Accepted`/`permanent` ADR-Folgepflicht
seit ihrer Verabschiedung nie eingelöst wurde.

## Umsetzungsskizze (grob — kein vollständiger Implementierungsplan)

ADR-konforme Mindest-Umsetzung von Option C, entlang der bestehenden
Hexagon-Schichten:

1. **Domänen-Modell** (`internal/domain/model/`): `TableSchema`
   (Spaltenmenge inkl. Typ-/OID-Information je Version) neben dem
   bestehenden `SchemaVersionID`; eine Registrierung, die für eine
   `(SourceTableID, SchemaVersionID)`-Kombination das zugehörige
   `TableSchema` hält.
2. **Outbound Port** (`internal/application/port/outbound/`):
   `SchemaStorePort` — Fähigkeiten mindestens: aktuelle Schema-Version einer
   Tabelle lesen, neue Version registrieren (bei erkannter Änderung),
   `TableSchema` zu einer Version lesen (für historische Interpretation).
   Das ist exakt die in `ADR-0015` benannte Folgepflicht.
3. **Adapter** (`internal/adapters/driven/…` o. ä., Analogie zu
   `PostgresMetadataAdapter` aus der Architektur-Sicht): Persistenz über eine
   neue Tabelle im neutralen Schema (`tools/schema/schema.yaml`, z. B.
   `cdc.table_schema`), ausgerollt über d-migrate wie die bestehenden
   Objekte (`ADR-0043`).
4. **Decoder** (`internal/adapters/driving/replication/decode/`):
   `tupleValues`/die Relation-Dekodierung muss die Spalten-`Oid` mitführen,
   nicht nur den Text-Wert — Voraussetzung, um eine Typänderung überhaupt
   erkennen zu können.
5. **Mapper** (`internal/adapters/driving/replication/mapper/mapper.go`):
   `Assembler.Consume` behandelt `*decode.Relation` nicht mehr im
   `default`-Zweig, sondern vergleicht die eingehende Relation gegen das
   zuletzt bekannte `TableSchema` der Tabelle. Unverändert →
   nichts zu tun. Kompatible Erweiterung (z. B. neue Spalte) → neue
   `SchemaVersionID` über `SchemaStorePort` registrieren, `TableBinding`
   aktualisieren (schließt `LH-FA-SCH-005`). Nicht sicher interpretierbare
   Änderung (z. B. Typwechsel, den der Decoder nicht verlustfrei versteht) →
   sichtbarer Fehler der Klasse `schema` (`SPEC-008`), keine stille
   Fortsetzung (schließt `LH-FA-SCH-004` Negative-Fall).
6. **Wiring** (`internal/bootstrap/wiring.go`): `SchemaStorePort`-Adapter
   verdrahten; die heutige statische `Version: 1`-Verdrahtung bleibt für die
   Erstaktivierung korrekt, spätere Versionen kommen aus dem Store.
7. **Tests**: Die beiden in slice-030 bereits geschriebenen Black-Box-Tests
   (`TestMVPSchemaChangeAddColumn`, `TestMVPSchemaChangeIncompatibleTypeChange`)
   dokumentieren die Ziel-Grenze bereits korrekt und werden nach der
   Umsetzung ohne Konzeptänderung grün — sie sind der natürliche
   Abnahmetest für den letzten Slice dieser Arbeit, nicht neu zu erfinden.

## Größenordnungs-Einschätzung

**Eigene Feature-Welle, kein Einzel-Slice.** Der Umfang berührt mindestens
vier Schichten (Domäne, neuer Outbound-Port, neuer Adapter samt
DB-Schema-Migration, Decoder/Mapper der Adapter-Schicht) und liegt damit weit
über der Slice-Grenze aus Modul 5 (≤ 3 Liefer-Punkte, höchstens zwei
Schichten). Das ist derselbe Befund, den die Roadmap bereits für
„Retention-Löschausführung" (`LH-FA-RET-002`…`006`, Größe **L**) getroffen
hat — dort wurde reine Feature-Arbeit korrekt als eigene Welle abgetrennt,
statt sie in eine E2E-Abdeckungs-Welle zu pressen. Empfehlung in derselben
Form:

- Neue Feature-Welle, z. B. „Schema-Evolution-Nachlieferung (`ADR-0015`)",
  Größe **L**, geschnitten (Planner-Aufgabe) in mindestens drei Slices
  entlang der Skizze oben (Persistenz-Fähigkeit ohne Live-Verdrahtung →
  dynamische Re-Versionierung im Consume-Pfad → Typ-Auswertung/Fehlerklasse
  `schema`).
- Optional eine nachfolgende schlanke „E2E-Abdeckung — Schema-Evolution"-Welle
  (Größe **S**, Analogie zu „E2E-Abdeckung — Retention" nach
  „Retention-Löschausführung"), die die bereits geschriebenen Black-Box-Tests
  aus slice-030 real grün macht, statt neue zu erfinden.

**Nicht Gegenstand dieses Verdikts:** wie `slice-030` selbst mit seiner
Closure umgeht (Carveout vs. Beobachtungs-Register-Eintrag vs. Verweis auf
die neue Welle als Folge-Slice-Adresse) — das bleibt Planner-Entscheidung
nach Modul 5 §Offene Risiken werden bei Closure aufgelöst. Dieses Verdikt
beantwortet ausschließlich die ADR-Frage: `ADR-0015` gilt unverändert fort,
der heutige Zustand ist eine bislang nicht eingelöste Folgepflicht, keine
ADR-Verletzung, die eine Korrektur der ADR verlangt.

Weder Produktionscode noch die ADR-Datei selbst wurden im Rahmen dieses
Verdikts geändert.
