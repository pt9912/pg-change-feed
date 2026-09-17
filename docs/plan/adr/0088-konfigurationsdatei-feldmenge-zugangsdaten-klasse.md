# ADR-0088: Konfigurationsdatei — Feldmenge, Zugangsdaten-Klasse und die Env-Durchleitung (`SPEC-016`)

**Status:** Accepted — Supersedes [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md)
in **einer** Klausel: der Aufzählung der zulässigen Datei-Felder in deren
§Entscheidung Festlegung 6 („Nur die übrigen, nicht credential-tragenden
Felder wandern in die Datei: `source_id`, `publication`, `slot`, `tables` …,
`log_level` sowie die `WALRetentionWarnBytes`/`WALRetentionErrorBytes`-Overrides").
Diese Aufzählung wird durch die **Klassen-Regel** ersetzt, die derselbe Satz
schon formuliert; die Aufzählung steht danach vollständig da (Festlegung 2).
Alles Übrige der [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md)
bleibt **hiermit bestätigt** und wird nicht wiederholt: das strikte Decoding
(Festlegung 1), die Additivität mit Feld-für-Feld-Vorrang der
Umgebungsvariable (Festlegung 2), die Zugriffsform über `CDC_CONFIG_FILE`
ohne CLI-Flag (Festlegung 5), die Fehlerklasse `configuration` und der
Sentinel (Festlegung 7) — und das **Prinzip** der DSN-Ausschlüsse selbst
(Festlegung 6, erster Satz: „Secrets bleiben env-var-exklusiv").

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf
Nutzerhinweis vom 2026-09-17 — „Wir haben viele neue Env-Vars eingeführt, aber
das Konfigfile `SPEC-016` nicht nachgezogen"; anderer Kontext als die
Implementer-Läufe, die `CDC_HTTP_ADDR`/`CDC_GRPC_ADDR`/`CDC_NATS_URL` und die
zwei Token-Klassen einzeln in die Verdrahtung gebracht haben — Modul 8
§Rollen-Regeln: „Architect schreibt")

**Bezug:** [`SPEC-016`](../../../spec/pflichtenheft.md) (die Feldform der Datei —
von dieser ADR **geschärft**), [`LH-QA-SEC-001`](../../../spec/lastenheft.md)/
[`LH-QA-SEC-002`](../../../spec/lastenheft.md) (Least-Privilege,
Zugangsdaten-Trennung — die tragende Zusage der Ausschlussklasse),
[`LH-FA-SST-006`](../../../spec/lastenheft.md)/[`LH-FA-SST-007`](../../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../../spec/lastenheft.md)
(die drei Oberflächen, deren Adressen zur Entscheidung stehen),
[`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) (in einer Klausel
abgelöst, sonst bestätigt), [`ADR-0057`](0057-http-grpc-api.md) (HTTP-Adapter
und die zwei Token-Klassen), [`ADR-0055`](0055-nats-change-notification-wecksignal.md)
(`CDC_NATS_URL`), [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
(`CDC_GRPC_ADDR`), [`SPEC-008`](../../../spec/pflichtenheft.md)
(Fehlerklasse `configuration`), [`AGENTS.md`](../../../AGENTS.md)
§3.5/§3.12, `internal/bootstrap/wiring.go` (`ConfigFromEnv`, `Config`),
`internal/bootstrap/config_file.go` (`fileConfig`, `forbiddenFileDSNKeys`,
`mergeConfig`), `internal/bootstrap/config_file_internal_test.go`,
`internal/bootstrap/config_file_rest_internal_test.go`,
`compose.yaml`, `docs/user/benutzerhandbuch.md` (§5.1 Env-Variablen,
§5.2 Konfigurationsdatei)

**Schärft:** [`SPEC-016`](../../../spec/pflichtenheft.md) — die ADR macht die
Feldmenge, die Zugangsdaten-Klasse und die Durchleitung der
Umgebungsvariablen für die Konfigurationsdatei verbindlich. Das
Pflichtenheft trägt die Festlegung (Rang 2), diese ADR die Entscheidung und
ihre Begründung.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass.** Der Nutzerhinweis vom 2026-09-17 lautet: *„Wir haben viele
neue Env-Vars eingeführt, aber das Konfigfile `SPEC-016` nicht nachgezogen."*
`compose.yaml` führt heute **13** `CDC_*`-Variablen; `SPEC-016` beschreibt
**sieben** Datei-Felder.

**(2) Der Ist-Stand — drei Träger, eine übereinstimmende Lücke.** Gemessen in
diesem Zug an allen drei Stellen:

| Träger | Stand |
|---|---|
| `docs/user/benutzerhandbuch.md` §5.1 | **14** `CDC_*`-Variablen (inkl. `CDC_CONFIG_FILE`, der drei Oberflächen-Adressen und beider Token) |
| `docs/user/benutzerhandbuch.md` §5.2 | **7** Datei-Felder; die Secret-Klausel nennt **nur** `capture_dsn`/`admin_dsn`/`reader_dsn` |
| `spec/pflichtenheft.md` `SPEC-016` | **7** Datei-Felder; dieselbe Drei-Schlüssel-Klausel |
| `internal/bootstrap/config_file.go` | `fileConfig` deklariert **genau diese 7** `yaml`-Felder; `forbiddenFileDSNKeys` führt **3** Schlüssel |

Es ist also **keine Drift** zwischen Spec, Handbuch und Code — die drei sagen
dasselbe und lassen dasselbe weg. Damit ist der Nachzug **keine
Prosa-Korrektur**, sondern eine Erweiterung mit Code-Anteil.

**(3) Der zweite, härtere Defekt — gemessen, nicht vermutet.** Die
Haupt-Klausel von `SPEC-016` lautet: „additiv zu den Umgebungsvariablen, mit
Umgebungsvariable-schlägt-Datei-Feld-für-Feld-Precedence". Für fünf
Variablen ist sie heute **falsch**:

- `ConfigFromEnv` liest sie — `wiring.go:276-280`:
  `cfg.NatsURL = getenv(envNatsURL)`, `cfg.HTTPAddr`, `cfg.APITokenReader`,
  `cfg.APITokenAdmin`, `cfg.GRPCAddr`.
- `mergeConfig` (`config_file.go:148-193`) — der Pfad, den
  `ConfigFromEnvAndFile` bei gesetztem `CDC_CONFIG_FILE` nimmt, und den
  `cmd/pg-change-feed/main.go` für den Betriebslauf aufruft — **liest keine
  dieser fünf**. Er setzt die fünf `Config`-Felder nicht, sie bleiben leer.

Die Folge ist **still**: bei gesetztem `CDC_CONFIG_FILE` bleiben HTTP-API,
gRPC-Stream und NATS-Wecksignal **deaktiviert** und beide Token-Klassen
**leer** — obwohl die Variablen gesetzt sind. Der Prozess startet dabei
gesund. Das ist nicht „Feld fehlt in der Tabelle": die **durchgereichte**
Env-Herkunft fehlt. Kein Test nennt die fünf (Suche über die zwei
`config_file`-Tests: 0 Treffer) — die Lücke ist **unentdeckt**, nicht
entschieden.

**(4) Die Ausschlussklausel ist unvollständig — und ihre Begründung fehlt an
der entscheidenden Stelle.** Die zwei Token-Schlüssel werden heute
abgewiesen, aber nur durch `KnownFields(true)` als *unbekannter* Schlüssel —
ohne den Grund zu nennen. Die drei DSN-Schlüssel dagegen haben eine eigene,
den Secret-Grund nennende Fehlerzeile (`forbiddenFileDSNKeys`). Der
Unterschied ist nicht kosmetisch: die heutige Meldung sagt „unbekannt" statt
„unzulässig, weil Zugangsdaten" — und der Satz in `SPEC-016` liest sich
dadurch so, als seien die Token-Felder **bloß noch nicht eingetragen**, statt
**ausgeschlossen**.

**(5) Der Diskriminator ist bereits entschieden — nur nicht angewandt.**
[`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6
formuliert die Regel als **Klasse** („Nur die übrigen, **nicht
credential-tragenden** Felder wandern in die Datei") und zählt sie dann als
Liste auf. Diese ADR wendet die Klasse an und macht die Liste vollständig;
sie erfindet keine Regel.

**(6) Warum die drei Adressen nicht dieselbe Klasse sind — der Diskriminator
im Einzelnen.** `CDC_HTTP_ADDR` und `CDC_GRPC_ADDR` sind **Horch-Adressen der
eigenen Oberflächen** in der Form `host:port` (`internal/adapters/driving/http`,
`…/grpc`); sie können **keine** Zugangsdaten tragen. `CDC_NATS_URL` ist eine
**URL-Form**: sie kann Benutzer und Passwort einbetten
(`nats://user:pass@host:4222`) — und die Abwesenheit von Zugangsdaten in
einem *konkreten Wert* ist keine Eigenschaft des *Feldes*. Genau das ist der
Unterschied, an dem die beiden unten getrennt entschieden werden.

**(7) Was die Beispiele damit zu tun haben — und was nicht.** Die
Beispiel-Clients beziehen Adresse und Token aus Umgebungsvariablen und Flags;
ob sie **auch** die Konfigurationsdatei lesen dürfen, ist eine eigene
Entscheidung. Sie fällt in [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
Festlegung 5 und lautet **nein** — die Datei ist der Deployment-Eingang des
Feeds, nicht der eines Integrator-Programms. Hier ist nur die Grenze zu
nennen, damit die zwei Aussagen nicht auseinanderlaufen.

## Entscheidung

Wir wählen: **der Diskriminator ist „kann dieses Feld Zugangsdaten tragen?";
`http_addr` und `grpc_addr` werden Datei-Felder, `nats_url` und die zwei
Token-Schlüssel bleiben env-exklusiv (mit eigener, den Grund nennenden
Fehlerzeile); und die fünf optionalen Oberflächen-Variablen werden unter
gesetztem `CDC_CONFIG_FILE` durchgereicht.**

Fünf Festlegungen:

### 1 — Der Diskriminator: ein Feld, das Zugangsdaten tragen kann, ist env-exklusiv

- **Die Regel:** Ein Feld darf in der Konfigurationsdatei vorkommen, wenn es
  **nicht** zur Klasse der zugangsdaten-tragenden Felder gehört. Diese Klasse
  umfasst heute: die drei DSN-Schlüssel (`capture_dsn`, `admin_dsn`,
  `reader_dsn`), die zwei Token-Schlüssel (`api_token_reader`,
  `api_token_admin`) und `nats_url`.
- **Die Begründung ist die von [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md)
  Festlegung 6, unverändert:** eine Konfigurationsdatei ist für **andere
  Aufbewahrungs- und Verteilwege** bestimmt als eine Umgebungsvariable; ein
  zweiter Zugangsdaten-Pfad neben den etablierten Secret-Mechanismen der
  Orchestrierung ist von [`LH-QA-SEC-001`](../../../spec/lastenheft.md)/
  [`LH-QA-SEC-002`](../../../spec/lastenheft.md) nicht vorgesehen.
- **Der Unterschied zwischen den zwei Adressen und der URL ist kein Gefühl,
  sondern die Form:** `host:port` kann keine Zugangsdaten tragen; eine
  URL-Form kann es (Kontext (6)). Deshalb ist `nats_url` **nicht** dabei —
  nicht weil der heutige Wert keine Zugangsdaten trägt, sondern weil das
  Feld die Form dafür hat.

### 2 — Die Feldmenge: zwei neue Datei-Felder, und die Liste ist danach vollständig

- **Neu zulässig:** `http_addr` (entspricht `CDC_HTTP_ADDR`) und `grpc_addr`
  (entspricht `CDC_GRPC_ADDR`). Beide sind optional, additiv und dem
  Feld-für-Feld-Vorrang der Umgebungsvariable unterworfen — genau wie die
  sieben bestehenden Felder. Ein leeres Feld ist ein leeres Feld: die
  Oberfläche bleibt deaktiviert, wenn weder Datei noch Variable eine Adresse
  tragen.
- **Weiter env-exklusiv:** `nats_url` und die zwei Token-Schlüssel
  (Festlegung 1). Ein Treffer in der Datei ist ein Fehler, keine stille
  Auslassung.
- **Die vollständige Feldmenge danach:** `source_id`, `publication`, `slot`,
  `tables`, `log_level`, `http_addr`, `grpc_addr`,
  `wal_retention_warn_bytes`, `wal_retention_error_bytes` — **neun**
  zulässige Felder; dazu die **fünf** env-exklusiven Zugangsdaten-Schlüssel.
  Diese Aufzählung ist die Ersetzung der Liste aus
  [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6
  (§Status).
- **Kein Feld für „alles Übrige".** Das strikte Decoding bleibt: jeder andere
  Schlüssel bricht mit Exit-Charakter der Fehlerklasse `configuration` ab.

### 3 — Die Durchleitung: „kein Datei-Feld" heißt nicht „Env wird ignoriert"

- **Die fünf optionalen Oberflächen-Variablen behalten ihre Env-Herkunft
  unter gesetztem `CDC_CONFIG_FILE`.** `NatsURL`, `HTTPAddr`,
  `APITokenReader`, `APITokenAdmin` und `GRPCAddr` werden auf dem
  Datei-Pfad genauso aus der Umgebung gelesen wie auf dem Env-only-Pfad.
  Das ist die Reparatur des in Kontext (3) gemessenen Defekts und die
  Herstellung der Klausel, die `SPEC-016` seit ihrer ersten Fassung trägt.
- **Für die zwei neuen Felder heißt das die Precedence:** Datei-Wert als
  Basis, eine gesetzte Umgebungsvariable überschreibt ihn; eine **leere**
  Umgebungsvariable lässt den Datei-Wert stehen — dieselbe `overrideString`-
  Form wie bei `source_id`/`publication`/`slot`, kein zweiter Mechanismus.
- **Für die env-exklusiven Schlüssel heißt es:** sie kommen ausschließlich
  aus der Umgebung, auf beiden Pfaden; die Datei kann sie nicht setzen und
  nicht überschreiben.
- **Es ist eine Verhaltensänderung, und sie ist die beabsichtigte:** Ein
  Deployment, das heute `CDC_CONFIG_FILE` setzt und die drei
  Oberflächen-Variablen führt, betreibt seine Oberflächen heute **still
  deaktiviert**; nach der Reparatur sind sie aktiv. Das ist die Herstellung
  der zugesagten Semantik, kein neues Feature — und sie ist bei jedem
  solchen Deployment sichtbar zu prüfen (Konsequenzen). Der
  Compose-Betrieb ist nicht betroffen: `compose.yaml` setzt `CDC_CONFIG_FILE`
  nicht (gemessen).

### 4 — Die Fehlerform: die Zugangsdaten-Klasse nennt ihren Grund

- **Ein Zugangsdaten-Schlüssel in der Datei endet in `ErrConfiguration` mit
  einer eigenen, den Grund nennenden Fehlerzeile** — dieselbe Form, die die
  drei DSN-Schlüssel heute tragen, und **eine** Prüfung für die ganze Klasse
  (keine zweite Liste, kein zweiter Fehlertext-Bauplatz). Die Meldung nennt
  die Klasse (Zugangsdaten bleiben env-var-exklusiv) und den Schlüssel.
- **Warum nicht beim generischen „unbekannter Schlüssel" bleiben:** die
  beiden Zustände sind verschiedene Aussagen — *unbekannt* (Tippfehler,
  Formatfehler) und *unzulässig* (richtiger Name, falscher Ort). Die heutige
  Meldung macht aus dem zweiten den ersten und verschweigt damit genau die
  Zusage, die
  [`LH-QA-SEC-001`](../../../spec/lastenheft.md)/[`LH-QA-SEC-002`](../../../spec/lastenheft.md)
  trägt.
- **Kein neuer Fehlertyp, keine neue Klasse:** `ErrConfiguration` bleibt der
  Sentinel (`SPEC-008`, Fehlerklasse `configuration`), unverändert wie
  [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 7.

### 5 — Was diese ADR nicht ändert

- **Kein zweiter Precedence-Pfad, kein `--config`-CLI-Flag, kein neues
  Format.** `CDC_CONFIG_FILE` bleibt der einzige Zugriffsweg
  ([`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 5).
- **Keine Zulassung von Zugangsdaten in der Datei** — die Klasse wird
  **enger** gefasst, nicht weiter.
- **Keine Änderung an den Env-only-Vorgaben des Betriebs:**
  `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN` bleiben Start- und
  Rollen-Vorbedingungen ([`ADR-0047`](0047-rollenspezifische-dsn-verdrahtung.md));
  die drei Oberflächen-Variablen bleiben optional,
  `CDC_NATS_URL` bleibt die einzige, deren **gesetzter** Wert eine
  Start-Vorbedingung ist ([`ADR-0055`](0055-nats-change-notification-wecksignal.md)).
- **Keine Spec-Semantik außerhalb `SPEC-016`s.** Das Handbuch §5.2 und die
  Spec ziehen nach; die ADR ist kein Ersatz für sie.
- **Kein ADR-Verweis im Spec-Text.** Die Referenz-Richtung verbietet
  Spec → ADR mechanisch; der Grund des Nachzugs steht deshalb hier.
- **Die Beispiel-Clients lesen die Datei nicht** — die Grenze aus
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 5 bleibt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Die drei Oberflächen-Adressen

| Option | Pro | Contra |
|---|---|---|
| A1 — keine Datei-Felder; die drei Adressen bleiben env-exklusiv, die Lücke wird als benannte geführt | kein Code, kein Schema, kein Handbuch-Aufwand; die Datei bleibt genau so groß, wie sie ist | die Nutzeranweisung („nicht nachgezogen") bleibt unbeantwortet; die Asymmetrie ist nicht begründbar: `http_addr`/`grpc_addr` sind `host:port` und tragen **keine** Zugangsdaten — sie aus der Datei herauszuhalten hätte keinen Grund, der die DSN-Aussage überlebt; der Operator pflegt weiter zwei Herkünfte für eine Verdrahtung |
| A2 — **alle drei** werden Datei-Felder | einheitlich, „die Datei ersetzt die Env-Liste" wäre wahr | `nats_url` ist eine URL-Form und kann Zugangsdaten einbetten — die Zulassung öffnete den zweiten Zugangsdaten-Pfad, den [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6 gerade schließt; sie wäre nur mit einer Zusatzvalidierung (userinfo verboten) zu halten gewesen — eine neue Prüfung für ein Feld, dessen env-Herkunft ohnehin trägt |
| **A3 — `http_addr` und `grpc_addr` werden Datei-Felder; `nats_url` bleibt env-exklusiv (gewählt)** | die zwei Felder, die keine Zugangsdaten tragen können, wandern in die Datei; das eine, das sie tragen kann, bleibt draußen; die Regel ist **eine** („kann das Feld Zugangsdaten tragen?") und beide Entscheidungen fallen aus ihr | eine Asymmetrie, die erklärt werden muss (Kontext (6)) — sie ist begründet, nicht still; ein Operator, der die NATS-URL in der Datei erwartet, findet sie nicht und liest den Grund in `SPEC-016` |
| A4 — alle drei plus ein Verbot eingebetteter Zugangsdaten im `nats_url`-Wert | vollständig einheitliche Feldmenge; die Lücke bleibt trotzdem zu | eine neue, wert-abhängige Validierung in einem sonst rein form-geprüften Laden; die Regel hieße dann „das Feld ist zulässig, bestimmte Werte nicht" — schwerer zu lesen als „das Feld ist env-exklusiv", und für eine **optionale** Variable ohne Aufwandsvorteil |

### B — Die zwei Token-Schlüssel

| Option | Pro | Contra |
|---|---|---|
| B1 — nichts tun: sie bleiben „unbekannte Schlüssel" | kein Code; die Datei nimmt sie ohnehin nicht an | die Meldung sagt „unbekannt" statt „unzulässig"; `SPEC-016`s Satz liest sich, als fehle der Eintrag nur; die tragende Zusage (Least-Privilege) steht an der Stelle nicht, an der sie wirkt — genau die `AGENTS.md` §3.12-Klasse (eine Aussage, die als Beleg gelesen wird) |
| **B2 — sie in die Zugangsdaten-Klasse aufnehmen, mit eigener, den Grund nennenden Fehlerzeile (gewählt)** | eine Prüfung, eine Liste, ein Fehlertext für die ganze Klasse; die Meldung unterscheidet *unzulässig* von *unbekannt*; `SPEC-016` und Handbuch können den Ausschluss benennen statt ihn zu verschweigen | die Liste wird generalisiert (Code-Anteil) und braucht ihre Tests; die Klausel in [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6 muss zusammengelesen werden (§Status) |
| B3 — die Token **zulassen** (die Datei darf sie tragen) | ein Feld weniger, das env-exklusiv ist; symmetrisch zu `http_addr` | verletzt [`LH-QA-SEC-001`](../../../spec/lastenheft.md)/[`LH-QA-SEC-002`](../../../spec/lastenheft.md): ein Token ist ein Secret, und die Datei ist für andere Aufbewahrungs-/Verteilwege bestimmt — der Grund, der die DSN-Aussage trägt, trägt hier **genauso**, also gibt es keinen Grund, der sie überlebte |

### C — Die Durchleitung der fünf Variablen

| Option | Pro | Contra |
|---|---|---|
| C1 — nichts tun | kein Code | die Haupt-Klausel von `SPEC-016` bleibt für fünf Variablen **falsch** (Kontext (3)); die Oberflächen eines Deployments mit Datei bleiben still aus, und niemand kann das der Konfiguration ansehen — der schlechteste Fall dieser ADR |
| C2 — nur die zwei neuen Datei-Felder einführen, die Durchleitung unverändert lassen | kleinerer Eingriff; die Feldmenge wächst wie angefragt | der Defekt bleibt: die zwei neuen Felder würden gelesen, die **drei bestehenden** Oberflächen-Variablen weiterhin nicht — eine Konfiguration, die die Adressen aus der Datei liest, aber `CDC_NATS_URL` und beide Token aus der Umgebung **verliert**, ist schlechter als der Ist-Stand, weil sie mehr verspricht |
| **C3 — die fünf Variablen auf dem Datei-Pfad durchreichen und die zwei Datei-Felder ergänzen (gewählt)** | stellt genau die Klausel her, die `SPEC-016` trägt; keine neue Mechanik (`overrideString`, dieselbe Feld-für-Feld-Form); die zwei neuen Felder sind nur die additive Hälfte derselben Änderung | eine reale Verhaltensänderung an Deployments, die die still deaktivierten Oberflächen heute so betreiben (Konsequenzen); Tests müssen die Durchleitung **je Feld** festhalten, sonst driftet sie beim nächsten Zuwachs wieder |

**Fazit:** A3, B2, C3. A1 lässt den begründbaren Teil der Nutzeranweisung
offen; A2 öffnet einen Zugangsdaten-Pfad; A4 tauscht eine lesbare Regel gegen
eine wert-abhängige Prüfung; B1 lässt die tragende Zusage unausgesprochen;
B3 bricht sie; C1 lässt falschen Text stehen; C2 wäre eine halbe Reparatur,
die mehr verspricht als der Ist-Stand hält.

## Konsequenzen

- Positiv: Die Haupt-Klausel von `SPEC-016` („additiv … Env schlägt Datei
  Feld für Feld") ist danach **für alle** Variablen wahr — der stille Ausfall
  der drei Oberflächen und beider Token-Klassen unter `CDC_CONFIG_FILE` ist
  behoben.
- Positiv: Die Ausschlussklasse ist **anwendbar** statt aufgezählt: eine neue
  `CDC_*`-Variable wird nach **einer** Frage eingeordnet („kann das Feld
  Zugangsdaten tragen?"), und der Fehlerfall nennt seinen Grund, statt
  „unbekannt" zu sagen.
- Positiv: Kein neues Format, kein neuer Sentinel, kein zweiter
  Precedence-Pfad, keine neue Fehlerklasse; `ADR-0052`s Mechanik bleibt
  wörtlich in Kraft.
- Negativ: [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) und diese
  ADR müssen in **einer** Klausel zusammengelesen werden — das Muster der
  übrigen engen Klausel-Korrekturen.
- Negativ **mit benannter Wirkung**: Es ist eine **Verhaltensänderung**. Ein
  Deployment mit `CDC_CONFIG_FILE`, das die drei Oberflächen-Variablen führt,
  betreibt diese Oberflächen heute still deaktiviert und nach dem Zug aktiv.
  Das ist die Herstellung der zugesagten Semantik — aber es ist ein
  beobachtbarer Unterschied beim Start, und er gehört in die Closure-Notiz des
  umsetzenden Zuges, nicht in eine stille Nebenwirkung. Der Compose-Betrieb
  ist nicht betroffen (`compose.yaml` setzt `CDC_CONFIG_FILE` nicht,
  gemessen).
- Negativ mit Grenze: Die Klassen-Regel ist eine **Form**-Prüfung; sie sagt
  nichts über Werte. Ein Feld der erlaubten Menge kann einen schlecht
  gewählten Wert tragen — das prüft weiterhin das Laden je Feld
  (Pflichtangaben, Zahlformen) und im Übrigen der Betrieb.
- Folgepflicht (Code-Zug): `fileConfig` um `http_addr`/`grpc_addr` erweitern,
  `mergeConfig` um den Feld-für-Feld-Vorrang dieser zwei **und** um das
  Durchreichen der fünf Umgebungsvariablen; die Zugangsdaten-Prüfung auf die
  Klasse generalisieren (drei DSN- plus zwei Token-Schlüssel plus
  `nats_url`) samt Fehlerzeile. Tests: **je Feld** eine Durchleitung mit
  gesetzter Datei **und** gesetzter Umgebungsvariable (die Klasse fällt sonst
  beim nächsten Zuwachs wieder auseinander), plus die zwei neuen Felder in
  beiden Richtungen des Vorrangs, plus je ein Fehlerfall der
  Zugangsdaten-Klasse.
- Folgepflicht (Träger-Zug): `spec/pflichtenheft.md` `SPEC-016` — Feldtabelle,
  Ausschlussklausel (Klassen-Form) und Änderungshistorie;
  `docs/user/benutzerhandbuch.md` §5.2 und seine Änderungshistorie. Der
  Nachzug von Handbuch und Spec **im selben Zug wie der Code**
  (verkörperte Regel).
- Folgepflicht (Planner-Zug): der Vorgang ist ein **Slice** (er braucht Code,
  Spec, Handbuch und Tests) — die ADR entscheidet, der Slice liefert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test ./internal/bootstrap/...` (im gepinnten Container, netzlos) | **Durchleitung:** mit gesetzter Konfigurationsdatei **und** gesetzten `CDC_NATS_URL`/`CDC_HTTP_ADDR`/`CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` trägt das Ergebnis alle fünf Felder; **Precedence:** eine gesetzte Umgebungsvariable schlägt den Datei-Wert feldweise (`http_addr`/`grpc_addr`), eine leere lässt ihn stehen | `make test` (kein Gate) |
| `go test ./internal/bootstrap/...` | **Zugangsdaten-Klasse:** ein Datei-Dokument mit `capture_dsn`/`admin_dsn`/`reader_dsn`/`api_token_reader`/`api_token_admin`/`nats_url` endet in `ErrConfiguration` mit einer Begründung, die die Klasse nennt — nicht als generischer unbekannter Schlüssel | `make test` (kein Gate) |
| `make docs-check` | **Träger-Paarung:** die Feldmenge in `SPEC-016`, im Handbuch §5.2 und die Schlüssel-Prüfung im Code nennen dieselben Namen (Form-Prüfung des Nachzugs) | `make gates` (im Bündel) |
| Review-Prüfpflicht (nicht maschinell) | Ob ein neues Feld richtig eingeordnet wurde (Diskriminator), und ob die Verhaltensänderung im betroffenen Deployment sichtbar gemacht ist | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Drei benannte Trigger, sonst permanent:

1. **Eine neue `CDC_*`-Variable kommt in die Verdrahtung.** Dann ist der
   Diskriminator (Festlegung 1) anzuwenden und das Feld **entweder** in die
   erlaubte Menge **oder** in die Zugangsdaten-Klasse zu stellen — beides
   zusammen mit dem Träger-Nachzug. Kein Feld bleibt unentschieden liegen.
2. **Ein zugangsdaten-tragendes Feld wird in der Datei gebraucht** (etwa weil
   ein Secret-Mechanismus der Orchestrierung wegfällt oder ein Deployment
   keinen Zugang zu ihm hat). Dann ist die **Klassen-Entscheidung** neu zu
   treffen (mit eigenem Schutzkonzept, z. B. Dateirechte/Volume-Mount) — die
   Liste wird **nicht** stillschweigend um einen Schlüssel erweitert.
3. **`CDC_CONFIG_FILE` verliert seine Rolle als Deployment-Eingang** (etwa
   weil ein Secret-Store die gesamte Verdrahtung trägt oder ein
   Installations-Weg eine andere Form verlangt). Dann ist `ADR-0052`s
   Gesamtentscheidung zu prüfen, nicht diese Klausel.

Sonst permanent: die Datei trägt, was keine Zugangsdaten tragen kann; was
Zugangsdaten tragen kann, bleibt env-exklusiv; und jede gesetzte
Umgebungsvariable wirkt unabhängig davon, ob eine Datei geladen wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: Nutzerhinweis vom 2026-09-17, das Konfigfile sei bei den gewachsenen Umgebungsvariablen nicht nachgezogen. Gemessen in diesem Zug: die übereinstimmende Lücke über `SPEC-016`, Handbuch §5.2 und `fileConfig` (7 Felder, 3 ausgeschlossene Schlüssel) **und** ein zweiter Defekt — `mergeConfig` liest die fünf Oberflächen-Variablen nicht, die `ConfigFromEnv` liest, sodass bei gesetztem `CDC_CONFIG_FILE` HTTP-API, gRPC-Stream, NATS-Wecksignal und beide Token-Klassen still ausfallen. Entscheidet: Diskriminator „kann das Feld Zugangsdaten tragen?", `http_addr`/`grpc_addr` als Datei-Felder, `nats_url` und die zwei Token-Schlüssel env-exklusiv mit eigener Fehlerzeile, Durchleitung der fünf Variablen; supersedes die Feld-Aufzählung in [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6, alles Übrige dort bestätigt | [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6 · `internal/bootstrap/wiring.go` (`ConfigFromEnv`, `Config`) · `internal/bootstrap/config_file.go` (`fileConfig`, `forbiddenFileDSNKeys`, `mergeConfig`) · `spec/pflichtenheft.md` `SPEC-016` · `docs/user/benutzerhandbuch.md` §5.1/§5.2 · `compose.yaml` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0088` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
