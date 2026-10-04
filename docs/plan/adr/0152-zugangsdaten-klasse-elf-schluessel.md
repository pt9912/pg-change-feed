# ADR-0152: Zugangsdaten-Klasse der Konfigurationsdatei — elf Schlüssel

**Status:** Accepted — Supersedes [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)
und [`ADR-0102`](0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md),
jede nur in ihren **Zahl-Aussagen und Aufzählungen der Klasse** („sieben
Schlüssel“, die Aufzählung der drei DSN-, drei Token-Schlüssel und `nats_url`),
namentlich:

- aus [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) deren
  §Entscheidung Festlegung 1 (die geltende Klasse) und Festlegung 4 (der Ort der
  geltenden Zahl), beide Zeilen ihrer §Fitness Function und der
  „Sonst permanent“-Satz ihres §Re-Evaluierungs-Trigger;
- aus [`ADR-0102`](0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md)
  deren §Entscheidung Festlegung 1 (die Bestätigung der geltenden Klasse mit
  ihren sieben Schlüsseln), beide Zeilen ihrer §Fitness Function und der
  „Sonst permanent“-Satz ihres §Re-Evaluierungs-Trigger.

Alles Übrige beider ADRs bleibt bestätigt: der Diskriminator („kann dieses Feld
Zugangsdaten tragen?“, [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)),
die Trennung der Oberflächen-**Variablen** von den Klassen-**Schlüsseln**, die
Paarung als Review-Prüfpflicht ohne Sensor
([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)),
§Verglichene Alternativen, §Konsequenzen im Übrigen, §Geschichte und die
datierten Messtabellen. Die Zahlen **sechs** und **sieben** in den Vorgängern
bleiben als Aufzeichnung ihres Stands stehen.

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8; Anlass: Review Spec 0.15.0, Befund F-4)

**Bezug:** [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) und
[`ADR-0102`](0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md) (die
Zahl-Aussagen superseded),
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md),
[`ADR-0149`](0149-otlp-metrik-export-mechanismus.md) (`otlp_endpoint`,
`otlp_headers`),
[`ADR-0150`](0150-tls-und-mehrfach-token.md) (die zwei Token-Listen),
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
(nimmt Entscheidung und Fitness-Function-Regeln von der Zitat-Korrektur aus —
deshalb Folge-ADR), [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md),
`AGENTS.md` §3.5, §3.12, §3.13,
`internal/bootstrap/config_file.go` (`forbiddenFileCredentialKeys`),
`internal/bootstrap/config_file_internal_test.go`
(`TestConfigFromFileLehntZugangsdatenAb`, die Iteration über die Klasse).

**Schärft:** [`SPEC-016`](../../../spec/pflichtenheft.md#spec-016--konfigurationsdatei-cdc_config_file)
(Klassen-Satz; das Pflichtenheft führt die Klasse bereits mit elf Schlüsseln).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Form.** Das Vorbild ist [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md):
auch dort wuchs die Klasse (sechs auf sieben), und die Zahl-Aussagen der
`Accepted`-ADRs wurden namentlich durch eine Folge-ADR ersetzt. Eine bloße
Schärfung genügt nicht, weil die Zahl in §Entscheidung und §Fitness Function
der Vorgänger steht, die `AGENTS.md` §3.5 unberührbar nennt.

**Anlass.** Das Pflichtenheft führt in `SPEC-016` die Klasse jetzt mit elf
Schlüsseln: die drei DSN-Schlüssel (`capture_dsn`, `admin_dsn`, `reader_dsn`),
die drei Token-Schlüssel (`api_token_reader`, `api_token_admin`,
`nats_stream_token`), `nats_url` und vier neue — `otlp_endpoint`, `otlp_headers`
([`ADR-0149`](0149-otlp-metrik-export-mechanismus.md): URL-Form und Header
können Zugangsdaten tragen) sowie `api_tokens_reader` und `api_tokens_admin`
([`ADR-0150`](0150-tls-und-mehrfach-token.md): Listen von Token). Beide
neuen ADRs sind `Accepted` und nennen die Klassenzahl nicht.

**Gemessen** (2026-10-04, HEAD `3d8f0a4c`, Lesung ohne Lauf): `SPEC-016` zählt
elf Namen (3 + 3 + 1 + 4); `forbiddenFileCredentialKeys` in
`internal/bootstrap/config_file.go` führt noch **7** Einträge, und
`TestConfigFromFileLehntZugangsdatenAb` iteriert über dieselben sieben. Die
Spec führt also schon elf, der Code noch sieben: das ist der Rückstand, den der
umsetzende Slice schließt, kein Widerspruch der Spec in sich.

`tls_cert_file` und `tls_key_file` gehören **nicht** zur Klasse: sie tragen
Dateipfade, keine Zugangsdaten (`SPEC-016` führt sie als zulässige Datei-Felder).

## Entscheidung

Wir wählen: **Die Zugangsdaten-Klasse der Konfigurationsdatei umfasst elf
Schlüssel; die Zahl-Aussagen von `ADR-0101` und `ADR-0102` werden darauf gezogen.**

1. **Die geltende Klasse** (elf): `capture_dsn`, `admin_dsn`, `reader_dsn`,
   `api_token_reader`, `api_token_admin`, `nats_stream_token`, `nats_url`,
   `otlp_endpoint`, `otlp_headers`, `api_tokens_reader`, `api_tokens_admin`.
   Die Aufzählung trägt die Zahl: sie ist an der Stelle nachzählbar, an der sie
   gelesen wird.
2. **Die Zahl-Aussagen der zwei Vorgänger** werden an den im §Status genannten
   Stellen durch den Satz aus Festlegung 1 ersetzt. Unberührt bleiben ihre
   §Verglichene-Alternativen-Zeilen, §Geschichte-Zeilen, datierten Messtabellen
   und Titel.
3. **Die Größen bleiben getrennt lesbar.** *Fünf* meint die durchgereichten
   Oberflächen-Variablen von `ADR-0088` Festlegung 3, *elf* die
   Zugangsdaten-Schlüssel.
4. **Der Ort der geltenden Zahl** ist ab jetzt §Entscheidung und §Fitness
   Function **dieser** ADR; die Paarung der Träger bleibt ein Urteil von Hand.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: `ADR-0101`/`ADR-0102` führen weiter sieben | kein Eingriff an `Accepted`-ADRs | die aktive Fitness-Function-Zeile nennt eine Zahl, die ihr Wächter (Review/Verifier) gegen `SPEC-016` widerlegt, ohne dass die Zeile sagt, ob Drift oder überholte Aussage vorliegt (`AGENTS.md` §3.12 Instanz A) |
| **B — Zahl-Aussagen beider ADRs namentlich superseden (gewählt)** | die geltende Zahl steht dort, wo sie gelesen wird; das Muster von `ADR-0101` | mehrere Dokumente sind zusammen zu lesen |
| C — in `ADR-0149`/`ADR-0150` nachtragen | ein Dokument weniger | beide sind `Accepted` und unberührbar (`AGENTS.md` §3.5); eine Schärfung ersetzt die Zahl-Aussage der Vorgänger nicht |
| D — die Zahl streichen, auf `SPEC-016` zeigen | kleinste Drift-Fläche | die zählbare Aussage verschwindet; so schon verworfen in `ADR-0091` |

## Konsequenzen

- Positiv: Fitness-Function-Zeile und Spec nennen dieselbe Zahl.
- Positiv: kein Spec- oder Handbuch-Eingriff durch diese ADR; `SPEC-016` führt
  elf bereits.
- Negativ: `ADR-0101`, `ADR-0102` und diese ADR sind an der Klassen-Frage
  zusammen zu lesen.
- Folgepflicht (erwartet, nicht gefahren): der umsetzende Slice für
  [`ADR-0149`](0149-otlp-metrik-export-mechanismus.md) und
  [`ADR-0150`](0150-tls-und-mehrfach-token.md) trägt die vier Schlüssel in
  `forbiddenFileCredentialKeys` ein und passt `TestConfigFromFileLehntZugangsdatenAb`
  auf die elf Schlüssel an (Liste und Kommentar „sieben“); zieht das
  Benutzerhandbuch §5 (Klassen-Satz) nach.
- Hinweis (Index): die Zeile von `ADR-0101` trägt `; → ADR-0152` (Muster
  `ADR-0091` ← `ADR-0101`); die Zeile von `ADR-0102` behält ihren Wortlaut — die
  Titel-Deckelung bei 80 Zeichen lässt den Zeiger nicht zu.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (kein Sensor) | **Die Zugangsdaten-Klasse trägt elf Schlüssel** — `forbiddenFileCredentialKeys`, `SPEC-016` und das Benutzerhandbuch §5 nennen dieselben elf, **von Hand** nachzuzählen; Wächter sind Review und Verifier (`ADR-0089`) | — |
| `go test ./internal/bootstrap/...` (gepinnter Container, netzlos) | jeder der elf Klassenschlüssel in der Datei endet in `ErrConfiguration` mit eigener, den Grund nennender Zeile; `TestConfigFromFileLehntZugangsdatenAb` iteriert über die elf — *erwartet*: der Slice passt den Test an (heute sieben), nicht gefahren | `make test` (kein Gate) |

## Re-Evaluierungs-Trigger

Ein weiteres Feld kann Zugangsdaten tragen (Diskriminator von `ADR-0088`), oder
ein Träger außerhalb von `SPEC-016`, Handbuch und Code nennt die Klasse
namentlich; dann ist die Zahl gegen den neuen Trägerstand nachzumessen
(`AGENTS.md` §3.13).

Sonst permanent: die Klasse umfasst **elf** Schlüssel, die in Festlegung 1
genannten.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Accepted — Anlass: Review Spec 0.15.0, F-4 (die Klasse wächst von sieben auf elf, `ADR-0101` nennt sieben, `ADR-0149`/`ADR-0150` ohne Supersedes) | Review Spec 0.15.0, F-4 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0152` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
