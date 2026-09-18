# ADR-0101: Zugangsdaten-Klasse der Konfigurationsdatei — sieben Schlüssel

**Status:** Accepted — Supersedes [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md),
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) und
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) — jede
nur in ihren **Zahl-Aussagen** über die Zugangsdaten-Klasse, namentlich:

- aus [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) deren §Status
  (die bestätigende Klausel „der Diskriminator (Festlegung 1 samt ihren
  **sechs** Schlüsseln)"), deren §Schärft („das Pflichtenheft führt die
  Klasse bereits mit ihren **sechs** Schlüsseln"), deren §Entscheidung (der
  Kopfsatz, der Ersatztext der Festlegung 2 und die Zahl der Festlegung 3),
  deren §Konsequenzen (beide Aussagen über die gelieferten **sechs**), beide
  Zeilen ihrer §Fitness Function und der „Sonst permanent"-Satz ihres
  §Re-Evaluierungs-Trigger;
- aus [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
  deren §Entscheidung Festlegung 1 (die Klassen-Klausel des Ersatztextes der
  Träger-Paarung), deren §Konsequenzen (dieselbe Reichweite in der
  Zahlenfolge), deren §Fitness Function Zeile 1 (die **aktive** Zeile) und
  deren §Kontext (2) — allein ihr Schlusssatz („sechs von sechs"); die
  Messzeilen darüber bleiben stehen, sie tragen ihren Lauf;
- aus [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
  deren §Entscheidung Festlegung 1 (die Klassen-Aufzählung „Diese Klasse
  umfasst heute: …"), der Kopfsatz des zweiten Bullets ihrer Festlegung 2
  („Weiter env-exklusiv: `nats_url` und die zwei Token-Schlüssel") und die
  Klassen-Zeile ihrer §Fitness Function (die Regel „**Zugangsdaten-Klasse:**"
  über einer Sechser-Enumeration).

Alles Übrige dieser drei Dokumente bleibt **hiermit bestätigt** und wird nicht
wiederholt: der Diskriminator („kann dieses Feld Zugangsdaten tragen?"), die
neun zulässigen Datei-Felder, die Durchleitung der fünf
Oberflächen-Variablen, die Fehlerform, die Paarung als Review-Prüfpflicht ohne
Sensor ([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)),
§Verglichene Alternativen, §Konsequenzen im Übrigen, §Geschichte und die
Messtabellen der §Kontext-Abschnitte.

**Datum:** 2026-09-18

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf den Review
zu `slice-nats-drittstream-core` <!-- d-check:status-provenance -->, Befund
F-4; anderer Kontext als der Implementer-Lauf, der die Meldung schrieb, und als
der Reviewer-Lauf, der den Befund erhob — Modul 8 §Rollen-Regeln: „Architect
schreibt")

**Bezug:** [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) (die
Zahl-Aussagen superseded — das Dokument, dessen Titel die Zahl trägt) ·
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) (die
Zahl-Aussagen superseded — der heute **tragende** Träger der Klasse; sein
§Fitness Function Zeile 1 ist die aktive Zeile) ·
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) (die
Klassen-Aufzählung und die Klassen-Zeile ihrer §Fitness Function superseded —
der Ursprung der Klasse) ·
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)
(unberührt: kein Sensor, Wächter Review/Verifier) ·
[`ADR-0100`](0100-nats-dritter-vollinhalts-zustellweg.md) (der Anlass des
Zuwachses: `CDC_NATS_STREAM_TOKEN` tritt als env-exklusiver,
zugangsdaten-tragender Schlüssel in die Verdrahtung) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
(nimmt §Entscheidung, §Konsequenzen und die Fitness-Function-Regeln von der
Zitat-Korrektur aus — deshalb Folge-ADR statt in-place) ·
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(das gebaute Muster: eine falsche Zahl in einem `Accepted`-Text wird namentlich
superseded) · `AGENTS.md` §3.5 · §3.12 Instanz A · §3.13 · Review zu
`slice-nats-drittstream-core` <!-- d-check:status-provenance --> (Befund F-4
samt eigener Messung) · [`SPEC-016`](../../../spec/pflichtenheft.md)
(Klassen-Satz und Feldtabelle) · `docs/user/benutzerhandbuch.md` §5 ·
`internal/bootstrap/config_file.go` (`forbiddenFileCredentialKeys`) ·
`internal/bootstrap/config_file_internal_test.go` (die Iteration über die
Klasse)

**Schärft:** — (Korrektur-ADR ohne Spec-Stratum: das Pflichtenheft führt die
Klasse bereits mit ihren **sieben** Schlüsseln und ist nicht Gegenstand dieser
ADR, wie [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) und
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass — der Zuwachs ist geliefert, die Entscheidungen zählen noch
sechs.** Der NATS-Vollinhalts-Stream
([`ADR-0100`](0100-nats-dritter-vollinhalts-zustellweg.md)) bringt
`CDC_NATS_STREAM_TOKEN` in die Verdrahtung. Der Diskriminator der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) greift
(„kann dieses Feld Zugangsdaten tragen?" — ja, ein Verbindungs-Token), und der
umsetzende Slice hat die drei **lebenden** Träger der Klasse nachgezogen:
[`SPEC-016`](../../../spec/pflichtenheft.md), `docs/user/benutzerhandbuch.md` §5
und `forbiddenFileCredentialKeys` im Code. Die drei `Accepted`-ADRs, die die
Klasse führen, blieben stehen — und was an ihnen **gilt**, sind keine
Aufzeichnungen: die §Fitness Function Zeile 1 der
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) ist die
aktive Fassung der Träger-Paarung, die §Fitness Function der
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) behauptet die
Bindung „jeder der sechs Klassenschlüssel" (die Iteration läuft heute über
sieben), und die Klassen-Zeile der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-§Fitness
Function führt die Sechser-Enumeration als Prüfregel.

**(2) Die eigene Messung dieses Zugs** (2026-09-18, HEAD `0c490da`; die
Klassenschlüssel als sortierte Menge aus `forbiddenFileCredentialKeys`
gelesen, die Klassensätze der drei Träger gelesen):

| Träger | Gemessene Klasse |
|---|---|
| `internal/bootstrap/config_file.go`, `forbiddenFileCredentialKeys` | **7** Einträge: `capture_dsn`, `admin_dsn`, `reader_dsn`, `api_token_reader`, `api_token_admin`, `nats_url`, `nats_stream_token` |
| [`SPEC-016`](../../../spec/pflichtenheft.md) (Klassen-Satz und env-exklusive Liste) | **7** — „die drei DSN-Schlüssel …, die drei Token-Schlüssel … und `nats_url`" |
| `docs/user/benutzerhandbuch.md` §5 | **7** — dieselben Namen, mit der Form-Begründung |
| `internal/bootstrap/config_file_internal_test.go` (Iteration) | **7** — „jeder ihrer sieben Schlüssel" |
| [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) §Fitness Function, Klassen-Zeile | **6** — die Enumeration `capture_dsn`/`admin_dsn`/`reader_dsn`/`api_token_reader`/`api_token_admin`/`nats_url`, `nats_stream_token` fehlt |

Die Rechnung trägt keine Sechs: 3 + 3 + 1 = **7**. Die zulässige Feldmenge des
`fileConfig` (neun `yaml`-Tags) und die sieben Feld-Schlüssel des Handbuchs
§5.2 sind von diesem Zuwachs **nicht** berührt — `nats_stream_token` ist ein
env-exklusiver Klassen-Schlüssel, kein Datei-Feld.

**(3) Was die drei `Accepted`-Dokumente tragen — und was in ihnen
Aufzeichnung ist.** Zu unterscheiden sind die **Zahl-Aussagen** — die Klasse
*gilt* sechs — von den **Aufzeichnungen**: einer §Geschichte-Zeile, einer
Zeile der §Verglichenen Alternativen, einer datierten Messzeile im §Kontext
(die ihren Lauf nennt, `AGENTS.md` §3.12 Instanz A) und dem **Titel** eines
Dokuments. Die Zahl-Aussagen stehen in §Entscheidung, §Konsequenzen, §Fitness
Function und im „Sonst permanent"-Satz des §Re-Evaluierungs-Trigger; die
Aufzeichnungen tragen den Stand ihres Laufs und werden nicht nachgezogen (für
die datierten Historie-Zeilen in
[`SPEC-016`](../../../spec/pflichtenheft.md) und
`docs/user/benutzerhandbuch.md` §5 gilt dasselbe — geprüft, nicht
Gegenstand).

**(4) Die Wirkung — eine falsche Zahl an der Zeile, die das Review liest.** Die
Paarung trägt **kein** Sensor
([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)); ihr
Wächter sind Review und Verifier, und beide lesen die §Fitness
Function-Zeilen. Nimmt der Verifier die aktive Zeile beim Wort, zählt er in
`forbiddenFileCredentialKeys` sieben und in der Zeile sechs — und kann ihr
nicht entnehmen, ob der vierte Token-Schlüssel eine Drift des Codes oder die
Reichweite der Aussage ist. Genau das ist die Klasse `AGENTS.md` §3.12
Instanz A: eine als Beleg gelesene Zahl trägt ihren Ursprung.

**(5) Warum das eine Architect-Frage ist.** Die Defekte liegen in
`Accepted`-Texten, nicht im gelieferten Stand: die drei lebenden Träger stimmen
mit dem Code überein (gemessen, Kontext (2)). `AGENTS.md` §3.5 schließt die
In-place-Korrektur aus, und
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
nimmt §Entscheidung, §Konsequenzen und die **Fitness-Function-Regeln**
ausdrücklich von der Zitat-Korrektur aus — also Folge-ADR, dasselbe Werkzeug,
mit dem die [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) die
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) und die
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) die
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) beerbt
haben. Kein Reviewer→Implementer-Pfeil: es gibt keine Fixrunde am Code.

## Entscheidung

Wir wählen: **Die Zugangsdaten-Klasse der Konfigurationsdatei umfasst sieben
Schlüssel — die drei DSN-Schlüssel, die drei Token-Schlüssel und `nats_url`.
Die Zahl-Aussagen der drei `Accepted`-Dokumente, die sie als sechs führen,
werden auf die Sieben gezogen; kein Code-, Spec-, Handbuch- oder
Sensor-Eingriff.**

Vier Festlegungen:

1. **Die geltende Klasse** — die drei DSN-Schlüssel (`capture_dsn`,
   `admin_dsn`, `reader_dsn`), die drei Token-Schlüssel (`api_token_reader`,
   `api_token_admin`, `nats_stream_token`) und `nats_url`. Die Aufzählung
   ersetzt die Zahl nicht, sie trägt sie: die Zahl ist an der Stelle
   nachzählbar, an der sie gelesen wird (`AGENTS.md` §3.12 Instanz A).
   Die drei lebenden Träger führen sie vollständig (gemessen, Kontext (2)).
2. **Die Zahl-Aussagen der drei `Accepted`-Dokumente werden ersetzt** — an
   den im §Status namentlich genannten Stellen gilt ab jetzt derselbe Satz:
   *die Zugangsdaten-Klasse umfasst sieben Schlüssel; die drei DSN-Schlüssel,
   die drei Token-Schlüssel und `nats_url`; `SPEC-016`, das Handbuch §5 und
   `forbiddenFileCredentialKeys` nennen sie vollständig.* Unberührt bleiben
   ihre §Verglichene-Alternativen-Zeilen, ihre §Geschichte-Zeilen, die
   datierten Messtabellen der §Kontext-Abschnitte und der Titel der
   [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) samt seiner
   Zelle im Index — der **Name** eines Dokuments ist keine Aussage, und die
   Supersession trägt die Korrektur, nicht eine Titel-Umschrift.
3. **Die zwei Größen bleiben getrennt lesbar.** *Fünf* meint die
   durchgereichten Oberflächen-**Variablen** der
   [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
   Festlegung 3 (`CDC_NATS_URL`, `CDC_HTTP_ADDR`, `CDC_GRPC_ADDR`,
   `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`), *sieben* die
   Zugangsdaten-**Schlüssel**. Beide Zahlen sind wahr; sie zählen
   verschiedene Mengen.
4. **Der Ort der geltenden Zahl.** Sie steht ab jetzt in §Entscheidung und
   §Fitness Function **dieser** ADR; die drei Träger tragen die Namen, die
   Paarung bleibt ein Urteil **von Hand** — kein Werkzeug rotet, wenn sie in
   einem Lauf unterbleibt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: die Sechs bleibt in allen drei `Accepted`-Dokumenten stehen | kein Eingriff an `Accepted`-ADRs; der gelieferte Stand ist richtig (gemessen, Kontext (2)) | die **aktive** §Fitness-Function-Zeile der [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) und beide Zeilen der [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) nennen eine Zahl, die ihr eigener Wächter (Review/Verifier) beim Nachzählen widerlegt — ohne dass die Zeile ihm sagt, ob er eine Code-Drift oder eine überholte Aussage sieht (`AGENTS.md` §3.12 Instanz A) |
| **B — die Zahl-Aussagen der drei Dokumente namentlich und klauselweise superseden (gewählt)** | die geltende Zahl ist an der Stelle richtig, an der sie gelesen wird; kein Code-, Spec-, Handbuch- oder Sensor-Eingriff; die Aufzeichnungen (§Geschichte, §Verglichene Alternativen, datierte Messzeilen, Titel) bleiben unangetastet — die Korrektur trifft genau die Klauseln, die eine Aussage über heute treffen | vier Dokumente ([`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md), [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md), [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md), diese) müssen an der Klassen-Frage zusammengelesen werden; die `#Geschichte`- und `#Verglichene Alternativen`-Zeilen der Vorgänger tragen weiterhin die Sechs (als Aufzeichnung, nicht als Aussage) |
| C — die Zahl überall streichen und auf die Träger zeigen („die Klasse aus [`SPEC-016`](../../../spec/pflichtenheft.md)") | **eine** Stelle je Dokument weniger zu halten; die kleinste Drift-Fläche | die zählbare Aussage verschwindet: wer die Klasse erweitert, sieht an der Stelle keine Zahl und damit keinen Hinweis, dass die Liste **vollständig** war (so schon die abgelehnte Option C der [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md)); die §Fitness Function-Zeile ist der Ort, an dem der Verifier nachzählt, nicht der Ort eines Zeigers |
| D — die Klasse maschinell tragen (Go-Test vergleicht die drei Träger) | die Drift wäre automatisch rot | genau die Klasse, die [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) §Verglichene Alternativen C verworfen hat: die Handbuch-Hälfte ist Prosa, ein Parser darauf ist eine Formpflicht auf Prosa, und der Code würde zum Schiedsrichter über eine Rang-2-Spec-Tabelle; kein Anlass, die dort benannte Grenze neu aufzumachen |

**Fazit:** B. A lässt eine falsche Zahl an der Zeile stehen, die der Verifier
als Beleg liest; C kostet die zählbare Aussage und damit die Prüfbarkeit der
Vollständigkeit; D macht eine bereits abgewogene und verworfene Grenze neu
auf, ohne dass dieser Befund sie verlangt.

## Konsequenzen

- Positiv: Die aktive Träger-Paarung
  ([`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
  §Fitness Function Zeile 1), die Bindungs-Zeile je Schlüssel
  ([`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) §Fitness
  Function Zeile 2) und die Klassen-Zeile der
  [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-§Fitness
  Function nennen die Zahl, die ihr Wächter nachzählt: sieben.
- Positiv: **kein Code-, Spec-, Handbuch- oder Sensor-Eingriff.** Die drei
  lebenden Träger stimmen mit dem Code überein (gemessen, Kontext (2)); zu
  alt war die Zahl in den Entscheidungstexten, nicht ein Träger.
- Positiv: Die Korrektur folgt dem gebauten Muster für eine enge
  Klausel-Korrektur an `Accepted`-Dokumenten (`ADR-0048`, `ADR-0062`,
  [`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md),
  [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md),
  [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md),
  [`ADR-0093`](0093-digest-korrektur-adr-0087-kotlin-basis-image.md)); kein
  zweites Werkzeug, kein neuer Träger.
- Negativ: Vier `Accepted`-Dokumente tragen jetzt je einen Teil der
  Klassen-Aussage; wer die Klasse auditiert, liest sie zusammen.
- Negativ mit benannter Grenze: Die Paarung bleibt ein **Urteil** — **kein**
  Sensor rotet, wenn sie in einem Lauf unterbleibt (so schon
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)
  §Konsequenzen und
  [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
  §Konsequenzen). Diese ADR zieht die Zahl richtig; sie macht sie nicht
  erzwungen.
- Hinweis (Index): Die `ADR-0091`- und die `ADR-0092`-Zeile tragen den
  Vorwärts-Zeiger `; → ADR-0101`; er ersetzt dort den `, teilw.`-Vermerk
  (Muster `ADR-0089` ← `ADR-0092`). Die `ADR-0088`-Zeile trägt **keinen**
  Zeiger — die `structure`-Deckelung der `Titel`-Zelle bei 80 Zeichen lässt
  ihn nicht zu (dieselbe Lage wie bei [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md)
  §Konsequenzen). Den Teil-Charakter trägt ab jetzt die Zeile **dieser** ADR
  (`Supers. ADR-0088/0091/0092, teilw.`).
- Folgepflicht (Planner-Zug): die §6-Meldung des laufenden Slice nennt den
  **heute** tragenden Träger ([`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md),
  nicht die supersedierte
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)) und
  diese ADR als die eingelöste Remediation — im selben Commit nachgezogen;
  darüber hinaus **keine** Folgepflicht: kein DoD eines Slice lehnt sich an
  eine der ersetzten Zahlen an.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (kein Sensor) | **Die Zugangsdaten-Klasse trägt sieben Schlüssel** — `forbiddenFileCredentialKeys` führt sieben Einträge, [`SPEC-016`](../../../spec/pflichtenheft.md) und `docs/user/benutzerhandbuch.md` §5 nennen dieselben sieben — **von Hand** nachzuzählen. Der Wächter ist der **Review** und der **Verifier** ([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md), [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md): die Träger-Paarung trägt kein Werkzeug) | — |
| `go test ./internal/bootstrap/...` (im gepinnten Container, netzlos) | **Die Bindung je Schlüssel** trägt die **zweite** Zeile der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-§Fitness Function in ihrer hier ersetzten Enumeration: jeder der sieben Klassenschlüssel in der Datei endet in `ErrConfiguration` mit einer eigenen, den Grund nennenden Zeile (`config_file_internal_test.go` iteriert über die sieben) | `make test` (kein Gate) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Die drei benannten Trigger der
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) bleiben in Kraft und
werden hier nicht wiederholt — ihr Trigger 1 deckt den Zuwachs eines
zugangsdaten-tragenden Feldes bereits ab; **einer kommt hinzu** — er zählt in
deren Nummerierung weiter:

4. **Ein Träger außerhalb der drei genannten spricht die Klasse aus** — ein
   weiteres Dokument, ein Kommentar oder ein Sensor nennt sie namentlich. Dann
   ist die Zahl hier gegen den neuen Trägerstand nachzumessen, bevor diese ADR
   zitiert wird; der Suchlauf dafür ist ein `grep` über die Träger, nicht über
   den eigenen Diff (`AGENTS.md` §3.13).

Sonst permanent: die Zugangsdaten-Klasse der Konfigurationsdatei umfasst
**sieben** Schlüssel — die drei DSN-Schlüssel, die drei Token-Schlüssel und
`nats_url`; die fünf Oberflächen-**Variablen** der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
Festlegung 3 bleiben davon unberührt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — Anlass: Review zu `slice-nats-drittstream-core` <!-- d-check:status-provenance -->, F-4: die §3.13-Meldung des Slice nennt als Wächter die durch [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) supersedierte [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) und schlägt `Supersedes ADR-0091` vor, was die falsche Zahl in [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md)s aktiver §Fitness-Function-Zeile stehen ließe. Eigene Messung an HEAD `0c490da`: **sieben** in `forbiddenFileCredentialKeys`, in [`SPEC-016`](../../../spec/pflichtenheft.md), im Handbuch §5 und in der Test-Iteration. Unabhängiger Architect-Zug zieht die Zahl-Aussagen der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md), [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) und [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) auf sieben; die drei lebenden Träger bleiben unverändert | Review zu `slice-nats-drittstream-core` <!-- d-check:status-provenance --> F-4 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
