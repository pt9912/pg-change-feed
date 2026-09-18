# Beleg: slice-096

Vorgang: `slice-096` — Konfigurationsdatei-Nachzug.

Fund: **Zwei falsche Zahlen in zwei ADRs, die am selben Tag entstanden sind** —
und beide in **einer** Runde gefunden, weil die Verifikation sie nachgezählt hat.

- **`ADR-0088` §Entscheidung Festlegung 2** sagt „dazu die **fünf** env-exklusiven
  Zugangsdaten-Schlüssel". Die Klasse hat **sechs**: die drei DSN-Schlüssel, die
  **zwei** API-Tokens und **`nats_url`**. **Fünf Träger sagen sechs** —
  Festlegung 1 desselben Dokuments, `SPEC-016`, Handbuch §5.2,
  `forbiddenFileCredentialKeys` im Code und `ADR-0089`.
- **`ADR-0089`s Ersatztext** (Fitness Function) sagt, die drei Träger nennten
  „dieselben Namen". Gemessen nennt das Handbuch §5.2 **7 von 9** zulässigen
  Datei-Feldern — die zwei `wal_retention_*` fehlen dort; die **Klasse** nennen
  alle drei vollständig.

**Die zwei Größen sind nicht dieselbe — und das ist der Kern des ersten Falls.**
`ADR-0088` trägt **14 Vorkommen von „fünf"**; **13** meinen die fünf
Oberflächen-**Variablen** (`CDC_NATS_URL`, `CDC_HTTP_ADDR`, `CDC_GRPC_ADDR`,
beide Tokens) und sind **richtig**; **genau eines** zählt die **Klasse**. Die
Klassen-Liste führt drei Einträge, die in keiner der fünf Variablen vorkommen
(die DSN-Schlüssel); die Variablen führen zwei, die nicht in der Klasse sind
(`http_addr`, `grpc_addr`); gemeinsam sind drei von sechs bzw. fünf. **Beide
Zahlen sind wahr — über verschiedene Dinge.**

Beide Korrekturen sind als Folge-ADR gelaufen (`ADR-0091`, `ADR-0092`), weil
`ADR-0073` §Entscheidung 1 §Entscheidung und Fitness-Function-Regeln von der
Zitat-Korrektur **ausnimmt**.

**Ein Vorgang, eine Zählung:** beide Funde fallen in denselben Vorgang
(`slice-096`) — der Zähler bewegt sich **einmal**, obwohl zwei ADRs betroffen
sind.

Quelle: Verifikationsbericht zu `slice-096` (V-1, V-2) ·
`docs/plan/adr/0091-zugangsdaten-klasse-sechs-schluessel.md` ·
`docs/plan/adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md` ·
`internal/bootstrap/config_file.go` (`forbiddenFileCredentialKeys`) ·
`spec/pflichtenheft.md` (`SPEC-016`).
