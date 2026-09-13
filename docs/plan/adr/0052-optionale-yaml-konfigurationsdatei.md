# ADR-0052: Optionale YAML-Konfigurationsdatei ergänzt Umgebungsvariablen — DSNs bleiben env-var-exklusiv

**Status:** Accepted

**Datum:** 2026-09-13

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-SST-003`](../../../spec/lastenheft.md), [`LH-FA-SST-002`](../../../spec/lastenheft.md), [`LH-QA-SEC-001`](../../../spec/lastenheft.md), [`LH-QA-SEC-002`](../../../spec/lastenheft.md), [ADR-0026](0026-composition-root.md), [ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md)

**Schärft:** [`ARC-007`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`internal/bootstrap/wiring.go` liest die gesamte Verdrahtungs-Vorbedingung
ausschließlich über Umgebungsvariablen (`ConfigFromEnv`, `getenv(...)`):
`CDC_CAPTURE_DSN`, `CDC_ADMIN_DSN`, `CDC_READER_DSN`, `CDC_SOURCE_ID`,
`CDC_PUBLICATION`, `CDC_SLOT`, `CDC_TABLES`, `CDC_LOG_LEVEL`. Es existiert
kein YAML-/TOML-Parser, kein `--config`-Flag und keine Konfigurationsdatei
im Repo — der Datei-Kommentar von `wiring.go` benennt das offen: „Eine
vollständige Konfigurationsschicht mit Format-Wahl ist nicht Teil dieses
Verdrahtungsstands." `CDC_TABLES` trägt bereits heute mehrere
Tabellen-Aktivierungen als eine Komma-/Doppelpunkt-getrennte Zeichenkette
(`schema.table=tabelle-id:schema-version-id`) in einer einzigen
Umgebungsvariable — ein Format, das mit wachsender Tabellenzahl unleserlich
wird und keine Kommentare oder Gruppierung erlaubt.

[`LH-FA-SST-003`](../../../spec/lastenheft.md) fordert eine CLI, die
„Installation, Diagnose und Administration" unterstützt. Eine Installation
mit mehreren aktivierten Tabellen, benannten Quellen und Betriebs-Defaults
ist über eine einzelne, flache Umgebungsvariable je Feld strukturell
unhandlich; es gibt aber weder eine Lastenheft- noch eine
Pflichtenheft-Festlegung, *wie* Installation strukturiert Konfiguration
trägt — diese Lücke ist bislang unbenannt, keine Welle, kein Slice, keine
ADR und kein Beobachtungs-Register-Eintrag adressieren sie. Es existiert
kein Konflikt mit einer bestehenden `Accepted`-ADR; diese Entscheidung
eröffnet ein neues Feld, sie widerspricht keiner bestehenden.

`go.sum` trägt `gopkg.in/yaml.v3` bereits als transitive (indirekte)
Abhängigkeit; `go.mod` selbst führt keinen YAML-/TOML-/JSON-Parser als
direkte Abhängigkeit. Das Repo nutzt YAML bereits an zwei Stellen als
etabliertes Konfigurationsformat: `tools/schema/schema.yaml` (neutrales
Schemamodell, `ADR-0043`) und `compose.yaml` (Laufumgebung). Es gibt keine
entsprechende TOML- oder JSON-Konvention im Repo.

Die drei DSN-Umgebungsvariablen tragen Zugangsdaten (Benutzer, Passwort,
Host) — [`LH-QA-SEC-001`](../../../spec/lastenheft.md) und
[`LH-QA-SEC-002`](../../../spec/lastenheft.md) verlangen
Least-Privilege-Vergabe und getrennte Berechtigbarkeit administrativer und
lesender Zugriffe; `spec/pflichtenheft.md` §4 hält explizit fest: „Keine
Credentials in Logs." Eine Konfigurationsdatei ist ein Artefakt, das
typischerweise ins Repository, in ein ConfigMap-Objekt oder in ein
Backup-Verzeichnis wandert — Kanäle, die für Zugangsdaten nicht vorgesehen
sind (anders als Umgebungsvariablen, die über Secret-Mechanismen der
jeweiligen Orchestrierung — Docker Secrets, Kubernetes Secrets als
projizierte Env-Vars — injiziert werden können).

## Entscheidung

Wir wählen **C — optionale YAML-Konfigurationsdatei, additiv zu den
bestehenden Umgebungsvariablen, mit Env-Var-Vorrang je Feld und
env-var-exklusiven DSNs**:

1. **Format: YAML.** Ein neuer, direkter `gopkg.in/yaml.v3`-Import ersetzt
   die heute nur transitive Abhängigkeit. Die Datei wird **strikt**
   dekodiert (`yaml.Decoder` mit `KnownFields(true)`) — ein unbekannter
   Schlüssel bricht mit einem `ErrConfiguration`-Fehler ab, dieselbe
   Fail-Closed-Disziplin, die `.a-check.yml` für sich selbst dokumentiert
   (Kommentarzeile „unbekannte Schlüssel brechen … ab").
2. **Precedence: Umgebungsvariable schlägt Datei, Feld für Feld.** Ist die
   Datei angegeben, liefert sie die Basis-Werte; jede gesetzte
   Umgebungsvariable überschreibt das gleichnamige Feld einzeln — kein
   Alles-oder-nichts. Das folgt der 12-Factor-Konvention (Umgebung ist die
   im Deployment am spätesten und am einfachsten änderbare Schicht) und
   dem bestehenden Betriebsmuster von `compose.yaml`/Kubernetes: ein
   Operator überschreibt einzelne Werte, ohne die eingebundene Datei
   anzufassen.
3. **Abwärtskompatibilität: reine Env-Var-Deployments bleiben unverändert
   lauffähig.** Die Datei ist vollständig optional; ist sie nicht
   angegeben, verhält sich die Verdrahtung exakt wie heute
   (`ConfigFromEnv` unverändert als Env-only-Pfad nutzbar). `compose.yaml`
   braucht keine Anpassung, um weiter zu funktionieren.
4. **Parser-Ort: `internal/bootstrap`.** Der Lade- und Merge-Vorgang ist
   Verdrahtungs-Vorbedingung, keine Domänen- oder Adapter-Fähigkeit — er
   gehört zur Composition Root (`ADR-0026`, [`ARC-007`](../../../spec/architecture.md)),
   analog zum bestehenden `ConfigFromEnv`/`parseTables`. `.a-check.yml`
   erlaubt der Composition Root (`composition_root: ["internal/bootstrap/**", …]`)
   uneingeschränkte Importe; kein Hexagon-Schichten-Edge beschränkt den
   neuen `yaml.v3`-Import.
5. **Docker-only-Konsequenz:** eine neue, optionale Umgebungsvariable
   `CDC_CONFIG_FILE` trägt den Pfad zur Datei **im Container**. Ist sie
   nicht gesetzt, entfällt jeder Dateizugriff — der heutige Env-only-Pfad
   bleibt Default. Ist sie gesetzt, aber die Datei unter diesem Pfad nicht
   lesbar, ist das ein `ErrConfiguration`-Fehler (explizit angegeben,
   also erwartet vorhanden) — kein stiller Fallback auf Env-only. Die
   Datei gelangt über einen Docker-/Compose-Volume-Mount in den Container
   (read-only, analog zum bestehenden `./tools/schema/compose-init`-Mount
   in `compose.yaml`); ein `--config`-CLI-Flag ist nicht Teil dieser
   Entscheidung — `cmd/pg-change-feed/main.go` parst Argumente heute ohne
   das `flag`-Paket (positionale Sondermodi), ein CLI-Flag für den
   Konfigurationspfad ist mit der künftigen `LH-FA-SST-003`-Installations-
   Unterstützung nachrüstbar, ohne diese Entscheidung zu revidieren
   (Folgepflicht).
6. **Secrets bleiben env-var-exklusiv — die wichtigste Einzelentscheidung:**
   `CDC_CAPTURE_DSN`, `CDC_ADMIN_DSN` und `CDC_READER_DSN` dürfen in der
   Konfigurationsdatei **nicht** vorkommen. Trägt die Datei einen dieser
   Schlüssel, ist das ein `ErrConfiguration`-Fehler beim Laden — nicht
   stillschweigend ignoriert und nicht stillschweigend übernommen. Nur
   die übrigen, nicht credential-tragenden Felder wandern in die Datei:
   `source_id`, `publication`, `slot`, `tables` (die
   Tabellen-Aktivierungen, strukturiert als YAML-Mapping statt der
   heutigen Komma-/Doppelpunkt-Zeichenkette), `log_level` sowie die
   `WALRetentionWarnBytes`/`WALRetentionErrorBytes`-Overrides. Begründung:
   Eine Konfigurationsdatei ist für andere Aufbewahrungs- und
   Verteilwege bestimmt als eine Umgebungsvariable (siehe Kontext) — DSNs
   dort zuzulassen würde einen zweiten Zugangsdaten-Pfad neben den
   etablierten Secret-Mechanismen öffnen, den [`LH-QA-SEC-001`](../../../spec/lastenheft.md)/[`LH-QA-SEC-002`](../../../spec/lastenheft.md)
   nicht vorsehen.
7. **Validierungsfehler:** Der bestehende `ErrConfiguration`-Sentinel
   (`SPEC-008`-Fehlerklasse `configuration`, `%w`-Wrapping) trägt auch die
   neuen Fehlerpfade — ungültiges YAML, unbekannter Schlüssel, ein
   DSN-Feld in der Datei, eine fehlende Pflichtangabe nach dem Merge. Kein
   neuer Fehlertyp, keine zweite Fehlerklasse.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun (ausschließlich Umgebungsvariablen, heutiger Stand) | kein Aufwand, kein neues Format, kein neuer Fehlerpfad | `CDC_TABLES` bleibt für mehrere Tabellen strukturell unhandlich (eine Zeichenkette, kein Kommentar, keine Gruppierung); `LH-FA-SST-003`s „Installation" bekommt kein strukturiertes Trägerformat; die Lücke bliebe weiterhin unbenannt statt entschieden |
| B — JSON- oder TOML-Konfigurationsdatei statt YAML | JSON: kein neuer Parser nötig (`encoding/json` Stdlib); TOML: einfachere Grammatik als YAML, weniger Fallstricke (kein „Norway problem") | JSON erlaubt keine Kommentare — für eine Betriebs-/Installationsdatei mit erklärungsbedürftigen Feldern (Tabellen-Bindungen, Schwellenwerte) ein echter Nachteil; TOML hat im Repo keine Präzedenz und keine bestehende Abhängigkeit, wäre ein zusätzliches drittes Format neben den bereits etablierten YAML-Stellen (`tools/schema/schema.yaml`, `compose.yaml`) ohne Konsistenzgewinn |
| C — vollwertiges Konfigurations-Framework (z. B. Viper: Multi-Format, automatisches Env-Binding, Live-Reload) | deckt viele Formate gleichzeitig ab, Env-Binding „geschenkt" | schwergewichtige neue Abhängigkeit in einem Repo, das bislang nur zwei direkte Abhängigkeiten führt (`pgx`, `pglogrepl`); automatisches Env-Binding verdeckt genau die explizite Feld-für-Feld-Precedence-Regel, die diese Entscheidung braucht; Live-Reload einer Konfigurationsdatei ist ein eigenständiges, hier nicht angefragtes Feature mit eigenem Betriebsrisiko |
| **D — optionale YAML-Datei, additiv, Env-Var-Vorrang je Feld, DSNs env-exklusiv (gewählt)** | konsistent mit den zwei bestehenden YAML-Stellen im Repo; `yaml.v3` ist bereits transitive Abhängigkeit (kein Fremdkörper); volle Abwärtskompatibilität für bestehende Env-only-Deployments; Secrets bleiben auf dem etablierten Injektionsweg (Env-Var/Secret-Mechanismus der Orchestrierung); Feld-für-Feld-Precedence ist explizit und ohne Framework-Magie nachvollziehbar | zwei Konfigurationsquellen zu pflegen und zu dokumentieren statt einer; ein neuer, expliziter Fehlerpfad für „DSN in der Datei" muss implementiert und getestet werden; die Struktur der Tabellen-Aktivierung wandert auf zwei Formen (weiterhin `CDC_TABLES` als Zeichenkette für den Env-only-Pfad, zusätzlich YAML-Mapping in der Datei) |

## Konsequenzen

- Positiv: `LH-FA-SST-003`s Installations-Unterstützung bekommt ein
  strukturiertes, kommentierbares Trägerformat für mehrere
  Tabellen-Aktivierungen und Betriebs-Defaults, ohne die heutige
  Env-Var-Schnittstelle zu brechen.
- Positiv: Bestehende Deployments (`compose.yaml`, jede heutige
  Env-only-Installation) bleiben ohne jede Änderung lauffähig — die Datei
  ist rein additiv, kein Breaking Change, kein eigener Migrationspfad
  nötig.
- Positiv: Die Secret-Trennung (DSNs ausschließlich über Env-Var/
  Secret-Mechanismus) bleibt entlang derselben Grenze wie heute; die
  Konfigurationsdatei vergrößert die Angriffsfläche für versehentlich
  committete oder falsch abgelegte Zugangsdaten nicht.
- Negativ: Zwei Konfigurationsquellen (Datei, Umgebungsvariablen) mit
  einer expliziten Merge-Regel erhöhen die Zahl der zu testenden
  Kombinationen gegenüber dem heutigen Einzelquellen-Zustand.
- Negativ: Die Tabellen-Aktivierung existiert nach dieser Entscheidung in
  zwei Repräsentationen — der heutigen `CDC_TABLES`-Zeichenkette und dem
  neuen YAML-Mapping — bis (falls je gewünscht) eine spätere Entscheidung
  eine der beiden Formen zurückbaut; das ist hier nicht entschieden.
- Folgepflicht: Ein Slice implementiert `internal/bootstrap`s
  Datei-Loader, das strikte YAML-Decoding, die Feld-für-Feld-Merge-Logik
  gegenüber `ConfigFromEnv`, die DSN-Ablehnung und die zugehörigen
  Tests (Whitebox, `package bootstrap`, analog zu den bestehenden
  Config-Tests).
- Folgepflicht: `spec/pflichtenheft.md` erhält die konkrete Feldform der
  Datei (Schlüsselnamen, YAML-Struktur der Tabellen-Aktivierung,
  `CDC_CONFIG_FILE`) als `SPEC-<NNN>`/`LH-*.a`-Verfeinerung — das ist
  Gegenstand des umsetzenden Slices, nicht dieser ADR.
- Folgepflicht: `compose.yaml` bzw. eine künftige Betriebsdokumentation
  benennt die Volume-Mount-Konvention für die Datei, sobald der Slice
  sie nutzt; bis dahin bleibt `compose.yaml` unverändert env-var-basiert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/bootstrap`, Whitebox) | Env-Var überschreibt Datei-Feld für Feld; eine Datei mit `capture_dsn`/`admin_dsn`/`reader_dsn` (oder gleichwertigen DSN-Schlüsseln) liefert `ErrConfiguration`; ein unbekannter YAML-Schlüssel liefert `ErrConfiguration` (striktes Decoding) | `make test` (bestehendes Ziel; die konkreten Testfälle sind Umsetzungsdetail des folgenden Slices) |
| `.a-check.yml` | unverändert — der neue `yaml.v3`-Import liegt vollständig innerhalb `composition_root`, kein Hexagon-Schichten-Edge betroffen | `make a-check` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Wird `LH-FA-SST-003`s CLI um einen `--config`-Pfad-Parameter oder einen
eigenständigen `install`-Unterbefehl erweitert, der die Datei schreibt statt
nur zu lesen: Re-Evaluierung, ob `CDC_CONFIG_FILE` als alleiniger
Zugriffsweg ausreicht oder ein CLI-Flag zusätzlich nötig wird — das ändert
diese Entscheidung nicht grundsätzlich, ergänzt sie aber um einen zweiten
Zugriffsweg auf denselben Dateipfad. Sonst permanent — die
Format-/Precedence-/Secret-Trennung dieser ADR gilt unabhängig vom
Zugriffsweg auf die Datei.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-13 | Accepted — Architect-Entscheidung zu einer bislang unbenannten Lücke (kein Slice, keine Welle, kein Beobachtungs-Register-Eintrag); benannt durch direkte Nutzerfrage vor dem Schneiden eines umsetzenden Slices | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0052` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
