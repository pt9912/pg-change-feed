# ADR-0092: Feldmengen-Paarung — die Reichweite der drei Träger

**Status:** Accepted — Supersedes [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)
in **einer** Klausel und den drei Stellen, die deren Reichweiten-Aussage
wörtlich wiederholen: ihrer §Entscheidung **Festlegung 2**, dem Ersatztext der
Träger-Paarung — dessen ersten Satz („die Feldmenge in `SPEC-016`, im Handbuch
§5.2 und die Schlüssel-Prüfung im Code nennen dieselben Namen — **von Hand**
nachzuzählen") — sowie derselben Reichweite in ihrem §Kontext (2) („das
Handbuch §5.2 dieselben"), in ihrer §Fitness Function Zeile 1 und im „Sonst
permanent"-Satz ihres §Re-Evaluierungs-Trigger. Alles Übrige der
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) bleibt
**hiermit bestätigt** und wird nicht wiederholt: die Entscheidung selbst (die
Paarung ist Review-Prüfpflicht, **kein** Sensor), ihre Festlegungen 1 und 3,
die vier Messungen ihres §Kontext zum `docs-check`-Lauf, §Verglichene
Alternativen A–E, §Konsequenzen und die zwei benannten Re-Evaluierungs-Trigger
(1 und 2).

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf den
Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance -->, Befund
V-2 (LOW); anderer Kontext als der Verifier-Lauf, der den Befund erhob, und als
der Architect-Zug, der [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)
geschrieben hat — Modul 8 §Rollen-Regeln: „Architect schreibt")

**Bezug:** [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)
(§Entscheidung Festlegung 2 — superseded samt ihren drei Wiederholungen) ·
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(§Fitness Function, dritte Zeile — die Kette dieser Korrektur: `ADR-0088` →
`ADR-0089` → hier) ·
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(das gebaute Muster: eine falsche Zahl in §Kontext eines `Accepted`-Dokuments
wird namentlich superseded) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
(nimmt die Fitness-Function-Regeln von der Zitat-Korrektur aus — deshalb
Folge-ADR) · `AGENTS.md` §3.5 · §3.12 Instanz B ·
Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance --> (Befund
V-2 samt eigener Messung) · [`SPEC-016`](../../../spec/pflichtenheft.md)
(Feldtabelle und Klassen-Satz) · `docs/user/benutzerhandbuch.md` §5.2 ·
`internal/bootstrap/config_file.go` (`fileConfig`,
`forbiddenFileCredentialKeys`) ·
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`

**Schärft:** — (Korrektur-ADR ohne Spec-Stratum: das Pflichtenheft führt die
neun zulässigen Schlüssel und die Klasse bereits; diese ADR zieht die
**Reichweite** einer Paarungs-Aussage auf das, was die drei Träger tragen, wie
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) und
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass — eine zu weite Reichweite an der Zeile, die der Verifier
liest.** Der Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance -->
weist V-2 (LOW) aus: der Ersatztext der Träger-Paarung, den
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) an die
Stelle der dritten [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-Fitness-Function-Zeile
gesetzt hat, sagt, die drei Träger nennten „dieselben Namen". Gemessen gilt
das für die **Zugangsdaten-Klasse** — für die **Feldmenge** ist es zu weit:
das Handbuch §5.2 nennt sieben der neun zulässigen Schlüssel.

**(2) Die eigene Messung dieses Zugs** (2026-09-17, HEAD `825d5fd`; die
Schlüsselmengen als sortierte Mengen aus `fileConfig`s `yaml`-Tags, der
`SPEC-016`-Tabelle und dem §5.2-Beispiel gelesen):

| Träger | zulässige Feld-Schlüssel | Zugangsdaten-Klasse |
|---|---|---|
| `spec/pflichtenheft.md` `SPEC-016` (Feldtabelle) | **9** | **6** |
| `internal/bootstrap/config_file.go` (`fileConfig`; `forbiddenFileCredentialKeys`) | **9** — mengengleich mit der Tabelle | **6** |
| `docs/user/benutzerhandbuch.md` §5.2 (Beispiel und Prosa) | **7** — `source_id`, `publication`, `slot`, `tables`, `log_level`, `http_addr`, `grpc_addr` | **6** — dieselben Namen |

Die zwei fehlenden Namen sind `wal_retention_warn_bytes` und
`wal_retention_error_bytes` — die zwei einzigen zulässigen Schlüssel **ohne**
Env-Gegenstück (`SPEC-016`s Tabelle führt bei beiden „kein Env-Gegenstück"); im Handbuch kommen
sie weder im Beispiel noch in der Prosa vor. Die **Klasse** dagegen ist in
allen drei Trägern vollständig: sechs von sechs.

**(3) Warum die Klasse trägt und die Feldmenge nicht.** Die Paarung hat zwei
Gegenstände. Die **Klasse** ist die tragende Zusage
([`LH-QA-SEC-001`](../../../spec/lastenheft.md)/[`LH-QA-SEC-002`](../../../spec/lastenheft.md) —
ein zugangsdaten-tragendes Feld bleibt env-exklusiv), und sie ist die Hälfte,
die das Handbuch für den Nutzer aussprechen **muss**; sie nennt sie
vollständig. Die **Feldmenge** ist ein Katalog, und ihr vollständiger Ort ist
die Rang-2-Tabelle in `SPEC-016`; das Handbuch §5.2 zeigt die Konfiguration,
die es erklärt, und erhebt an keiner Stelle Anspruch auf Vollständigkeit. Eine
Aussage, die beide Gegenstände zu „dieselben Namen" zusammenzieht, ist damit
für den einen wahr und für den anderen zu weit.

**(4) Die Wirkung — eine Aussage, die in ihrer Form nicht falsifizierbar
ist.** Die Zeile ist die Nachfolgerin der dritten Fitness-Function-Zeile und
ausdrücklich für den **Verifier** geschrieben („Prüfe die **Belege**, nicht die
Behauptung"); sie sagt ihm zugleich, dass **kein** Werkzeug sie trägt. Nimmt er
sie beim Wort, zählt er im Handbuch und findet sieben statt neun — und kann der
Zeile nicht entnehmen, ob die zwei fehlenden Namen eine Drift des Handbuchs
oder die Reichweite der Aussage sind. Genau das ist die Klasse `AGENTS.md`
§3.12 Instanz B: eine als Beleg gelesene Aussage nennt ihren Anker — oder sie
ist als *erwartet* formuliert.

**(5) Warum das eine Architect-Frage ist.** Der Defekt liegt im
`Accepted`-Text der
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md). `AGENTS.md`
§3.5 schließt die In-place-Korrektur aus, und
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
nimmt die **Fitness-Function-Regeln** ausdrücklich von der Zitat-Korrektur aus —
also Folge-ADR in der engen Klausel, dasselbe Werkzeug, mit dem
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) ihren
Vorgänger beerbt hat. Kein Reviewer→Implementer-Pfeil: es gibt keine Fixrunde
am Code.

## Entscheidung

Wir wählen: **Der Ersatztext der Träger-Paarung wird auf die Reichweite
gezogen, die die drei Träger tatsächlich haben — die zulässige Feldmenge führen
`SPEC-016` und der Code mengengleich, die Zugangsdaten-Klasse nennen alle drei
vollständig. Das Handbuch §5.2 wird nicht erweitert.**

Drei Festlegungen:

1. **Der Ersatztext der Träger-Paarung** lautet:

   > **Träger-Paarung:** die **zulässige Feldmenge** führen
   > [`SPEC-016`](../../../spec/pflichtenheft.md) (neun Schlüssel) und der
   > Code (`fileConfig`s `yaml`-Tags) **mengengleich**; das Handbuch §5.2
   > nennt **sieben** der neun — die zwei `wal_retention_*`-Overrides kommen
   > dort nicht vor. Die **Zugangsdaten-Klasse** — **sechs** Schlüssel —
   > nennen `SPEC-016`, das Handbuch §5.2 und `forbiddenFileCredentialKeys`
   > im Code **vollständig**. **Von Hand** nachzuzählen. **Kein Sensor trägt
   > diese Zeile**
   > ([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md),
   > Kontext (1)–(4)): `make docs-check` liest keinen Go-Quelltext, und keine
   > `structure`-Regel stellt zwei Dokumente einander gegenüber. Der Wächter
   > ist der **Review** und der **Verifier** („Prüfe die **Belege**, nicht die
   > Behauptung").

   Der Satz „Kein Sensor trägt diese Zeile" und die Benennung der zwei Wächter
   bleiben wörtlich in Kraft; ersetzt wird die Reichweiten-Aussage.
2. **Die drei Wiederholungen derselben Reichweite werden mitgezogen** — der
   §Kontext (2), die §Fitness Function Zeile 1 und der „Sonst
   permanent"-Satz des §Re-Evaluierungs-Trigger der
   [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md). An
   ihrer Stelle gilt derselbe Satz: *die zulässige Feldmenge führen `SPEC-016`
   und der Code mengengleich; die Zugangsdaten-Klasse nennen alle drei Träger
   vollständig.*
3. **Das Handbuch §5.2 wird nicht erweitert.** Es zeigt eine betriebsfähige
   Konfiguration und nennt die Klasse; einen **vollständigen** Feldkatalog
   trägt die Rang-2-Tabelle in `SPEC-016`. Die zwei `wal_retention_*`-Overrides
   dort nachzutragen, erzeugte eine **zweite** vollständige Feldliste außerhalb
   der Spec: die Paarung hätte dann zwei Kataloge zu halten statt einen, und
   genau die Drift wüchse, gegen die sie gerichtet ist (Alternativen C und D).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: „dieselben Namen" bleibt stehen | kein Eingriff an einer `Accepted`-ADR; die Entscheidung der [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) bleibt wörtlich in Kraft | die Aussage ist für die Feldmenge zu weit, und zwar an der Zeile, die den Verifier zu seiner Prüfung anhält; er prüft gegen eine Zahl, die er nicht finden kann, oder liest die zwei fehlenden Namen als Handbuch-Drift — §3.12 Instanz B verlangt für eine als Beleg gelesene Aussage ihren Anker oder die Form *erwartet* |
| **B — der Ersatztext nennt die Reichweite je Träger (gewählt)** | die Aussage wird in der Form falsifizierbar, in der sie dasteht: neun ↔ neun mengengleich, sieben im Handbuch, sechs als Klasse in allen drei; kein Code-, Handbuch-, Spec- oder Sensor-Eingriff | die Zeile trägt drei Zahlen statt einer und wird länger; jede ist an ihrer Stelle nachzählbar (`AGENTS.md` §3.12) |
| C — das Handbuch §5.2 trägt die zwei `wal_retention_*`-Felder nach; dann ist „dieselben Namen" wahr | die knappe Form bliebe wahr; der Nutzer fände die zwei Overrides im Handbuch | erzeugt eine **zweite vollständige Feldtabelle** außerhalb `SPEC-016`s — die Paarung hätte zwei Kataloge zu halten; das Beispiel wüchse um Felder **ohne** Env-Gegenstück, deren Default `SPEC-013` ist, und außerhalb von `SPEC-016`, dem Code und seinen Tests nennt heute **kein** Träger die zwei Namen (gemessen: kein Treffer in `compose.yaml`, `tools/`, `test/`, `harness/`, `docs/user/`); eigener Slice mit eigener Begründung |
| D — die Paarung auf `SPEC-016` ↔ Code verkürzen (das Handbuch aus der Zeile nehmen) | die kürzeste tragfähige Aussage; die zwei maschinennächsten Träger blieben gepaart | das Handbuch §5.2 ist der Träger, den der **Nutzer** liest, und der einzige, der die **Zugangsdaten-Klasse** für ihn ausspricht — sie aus der Zeile zu nehmen, ließe die Klasse an der Stelle ohne Paarung, an der die tragende Zusage ([`LH-QA-SEC-001`](../../../spec/lastenheft.md)/[`LH-QA-SEC-002`](../../../spec/lastenheft.md)) hängt |

**Fazit:** B. A lässt eine in ihrer Form nicht falsifizierbare Aussage an der
Zeile stehen, die der Verifier als Beleg liest; C erzeugt eine zweite
vollständige Feldliste und einen Slice, den der Befund nicht verlangt; D gibt
die Paarung der Klasse auf — der Klasse, die alle drei Träger vollständig
führen.

## Konsequenzen

- Positiv: Die Zeile nennt jetzt je Träger, was er führt — und damit für den
  Verifier, was er nachzählen muss: neun ↔ neun mengengleich, sieben im
  Handbuch, sechs als Klasse in allen drei.
- Positiv: **kein Code-, Spec-, Handbuch- oder Sensor-Eingriff.** Die drei
  Träger stimmen mit dem Code überein (gemessen, Kontext (2)); zu weit war die
  Reichweite der Aussage, nicht ein Träger.
- Positiv: Die Korrektur folgt dem gebauten Muster
  ([`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  für die supersedierte §Kontext-Zeile,
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) für die
  Kette am selben Träger); kein zweites Werkzeug.
- Negativ: [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)
  und diese ADR müssen an einer Stelle zusammengelesen werden.
- Negativ mit benannter Grenze: Die Paarung bleibt ein **Urteil** — **kein**
  Sensor rotet, wenn sie in einem Lauf unterbleibt (so schon
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)
  §Konsequenzen). Diese ADR macht die Aussage prüfbar; sie macht sie nicht
  erzwungen.
- Hinweis (Index): Die `ADR-0089`-Zeile trägt den Vorwärts-Zeiger
  `; → ADR-0092`. Er **ersetzt** dort den `, teilw.`-Vermerk (Muster
  `ADR-0087` ← `ADR-0090`); die Zelle misst am Inhalt **75** Zeichen ohne und
  **79** mit dem Zeiger und bleibt unter der 80-Zeichen-Deckelung des
  `structure`-Moduls. Den Teil-Charakter trägt ab jetzt die Zeile **dieser**
  ADR (`Supers. ADR-0089, teilw.`).
- Folgepflicht (Planner-Zug): **keine** — wie bei
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md): kein
  DoD eines Slice lehnt sich an die Zeile an.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (kein Sensor) | **Träger-Paarung** — die Nachfolge-Zeile der dritten Fitness-Function-Zeile aus [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md): die **zulässige Feldmenge** führen `SPEC-016` und der Code (`fileConfig`s `yaml`-Tags) **mengengleich** (neun); das Handbuch §5.2 nennt **sieben**; die **Zugangsdaten-Klasse** (sechs) nennen `SPEC-016`, das Handbuch §5.2 und `forbiddenFileCredentialKeys` **vollständig** — **von Hand** nachzuzählen. Der Wächter ist der **Review** und der **Verifier**; **kein** Werkzeug trägt diese Zeile | — |
| `go test ./internal/bootstrap/...` (im gepinnten Container, netzlos) | **Nicht** Gegenstand dieser Zeile: die **Code-Hälfte** der Paarung — Durchleitung und Feld-für-Feld-Vorrang je Feld, die Zugangsdaten-Klasse je Schlüssel — tragen unverändert die erste und zweite Zeile der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-Fitness-Function | `make test` (kein Gate) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Die zwei benannten Trigger der
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) bleiben in
Kraft und werden hier nicht wiederholt; **zwei kommen hinzu** — sie zählen in
deren Nummerierung weiter:

3. **Das Handbuch §5.2 wird aus anderem Anlass überarbeitet** — oder die zwei
   `wal_retention_*`-Overrides werden dort gebraucht. Dann ist die Feldmenge
   des Handbuchs **neu zu messen**, bevor diese Zeile zitiert wird; trägt es
   die zwei Namen nach, ist Festlegung 1 dieser ADR an dieser Stelle zu
   superseden.
4. **Die `SPEC-016`-Feldtabelle oder die `yaml`-Tags von `fileConfig` wachsen**
   (ein neues zulässiges Datei-Feld, Re-Evaluierungs-Trigger 1/2 der
   [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)).
   Dann ist die Neun hier gegen den neuen Trägerstand nachzumessen — die Zeile
   ist die einzige Stelle, die beide Mengen nebeneinander führt.

Sonst permanent: die zulässige Feldmenge führen `SPEC-016` und der Code
mengengleich; die Zugangsdaten-Klasse nennen alle drei Träger vollständig; die
Paarung zwischen ihnen trägt das Review.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance -->, V-2 (LOW): der Ersatztext der Träger-Paarung aus [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) sagt „dieselben Namen". Eigene Messung an HEAD `825d5fd`: das Handbuch §5.2 nennt **sieben** der neun zulässigen Schlüssel (die zwei `wal_retention_*` fehlen), die Zugangsdaten-Klasse (sechs) dagegen vollständig; `SPEC-016` und der Code sind mengengleich (neun). Unabhängiger Architect-Zug ersetzt die Reichweiten-Aussage samt ihren drei Wiederholungen; das Handbuch §5.2 bleibt unverändert | Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance --> |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
