# ADR-0091: Zugangsdaten-Klasse der Konfigurationsdatei — sechs Schlüssel

**Status:** Accepted — Supersedes [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
in **einer** Klausel: dem Schlusssatz des dritten Bullets ihrer §Entscheidung
**Festlegung 2** („… — **neun** zulässige Felder; dazu die **fünf**
env-exklusiven Zugangsdaten-Schlüssel."). Alles Übrige der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) bleibt
**hiermit bestätigt** und wird nicht wiederholt: der Diskriminator
(Festlegung 1 samt ihren sechs Schlüsseln), die zulässige Feldmenge und ihre
neun Schlüssel, die Durchleitung der fünf Oberflächen-Variablen
(Festlegung 3), die Fehlerform (Festlegung 4), die Nicht-Änderungen
(Festlegung 5), §Verglichene Alternativen A/B/C, §Konsequenzen und alle vier
Zeilen ihrer §Fitness Function.

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf den
Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance -->, Befund
V-1 (MEDIUM); anderer Kontext als der Verifier-Lauf, der den Befund erhob, als
der Reviewer-Lauf, der ihn nicht sah, und als die Implementer-Läufe, die den
Code geschrieben haben — Modul 8 §Rollen-Regeln: „Architect schreibt")

**Bezug:** [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(§Entscheidung Festlegung 2 — eine Klausel superseded; Festlegung 1 ist die
enumerierende Fassung derselben Klasse) ·
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) (dasselbe
Werkzeug am selben Träger: eine enge Klausel-Korrektur an einer
`Accepted`-ADR) ·
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(das gebaute Muster für eine falsche Zahl in einem `Accepted`-Text) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
(nimmt §Entscheidung von der Zitat-Korrektur aus — deshalb Folge-ADR statt
in-place) · `AGENTS.md` §3.5 · §3.12 Instanz A ·
`docs/reviews/verify-slice-096.md` <!-- d-check:status-provenance --> (Befund
V-1 samt eigener Messung) · [`SPEC-016`](../../../spec/pflichtenheft.md)
(Klassen-Satz und Feldtabelle) · `docs/user/benutzerhandbuch.md` §5.2 ·
`internal/bootstrap/config_file.go` (`forbiddenFileCredentialKeys`) ·
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`

**Schärft:** — (Korrektur-ADR ohne Spec-Stratum: das Pflichtenheft führt die
Klasse bereits mit ihren sechs Schlüsseln und ist nicht Gegenstand dieser ADR,
wie [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) und
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass — ein Befund am Entscheidungstext, nicht am Gegenstand.**
Der Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance --> weist
V-1 (MEDIUM) aus: die §Entscheidung Festlegung 2 der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) sagt
von ihrer Liste, sie sei „danach vollständig", und zählt die
Zugangsdaten-Klasse in ihrem Schlusssatz eine zu klein:

> … — **neun** zulässige Felder; dazu die **fünf** env-exklusiven
> Zugangsdaten-Schlüssel.

Der ausgelieferte Stand ist richtig — er führt sechs Schlüssel. Der Defekt
sitzt allein im Entscheidungstext.

**(2) Die eigene Messung dieses Zugs** (2026-09-17, HEAD `825d5fd`; die
Schlüsselliste des Codes als sortierte Menge gelesen, die Sätze der drei
Träger und der beiden ADRs gelesen):

| Träger | Gemessene Klasse |
|---|---|
| `internal/bootstrap/config_file.go`, `forbiddenFileCredentialKeys` | **6** Einträge: `capture_dsn`, `admin_dsn`, `reader_dsn`, `api_token_reader`, `api_token_admin`, `nats_url` |
| `spec/pflichtenheft.md` `SPEC-016` | **6** — „die drei DSN-Schlüssel …, die zwei Token-Schlüssel … und `nats_url`" |
| `docs/user/benutzerhandbuch.md` §5.2 | **6** — dieselben Namen, mit der Form-Begründung |
| [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) §Entscheidung Festlegung 1 | **6** — dieselbe Aufzählung |
| [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) §Kontext (2) | „**sechs** env-exklusive" |

Die ersten vier Zeilen nennen dieselben sechs Namen; die fünfte nennt die
Zahl. Die Rechnung trägt keine Fünf: 3 + 2 + 1 = **6**.

**(3) Die zwei Größen — fünf *Variablen*, sechs *Schlüssel*.** Die andere Fünf
desselben Dokuments ist richtig und bleibt es: die **fünf optionalen
Oberflächen-Variablen** der Festlegung 3 (`CDC_NATS_URL`, `CDC_HTTP_ADDR`,
`CDC_GRPC_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`). Die beiden
Mengen sind nicht dieselbe. Die Klasse führt die **drei DSN-Schlüssel**, die
zu den fünf Variablen nicht gehören (sie sind Pflicht-Variablen ohne
Datei-Gegenstück); die fünf Variablen führen **`http_addr`** und
**`grpc_addr`**, die der Klasse nicht angehören (sie sind zulässige
Datei-Felder). Gemeinsam sind drei: `nats_url` und die zwei Token-Schlüssel.
Die Zahl am falschen Ort ist die **Klassen**-Zahl.

**(4) Die Wirkung — und warum kein Sensor sie sieht.** Das Bullet, in dem der
Satz steht, nennt sich selbst die Ersetzung der Liste aus
[`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6: wer die
Klasse erweitert (Re-Evaluierungs-Trigger 1/2 der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)) oder
sie auditiert und die Zahl nimmt, nimmt eine zu kleine — und die Zusage
„vollständig" deckt dann fünf statt sechs ab. `make docs-check` ist über den
ganzen Vorgang grün (im Verifikationsbericht gemessen, V-1): die Zahl hat
keinen Sensor, sie hat Leser.

**(5) Warum das eine Architect-Frage ist.** Der Defekt liegt im
`Accepted`-Text, nicht am Code. `AGENTS.md` §3.5 schließt die
In-place-Korrektur aus, und
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
nimmt §Entscheidung ausdrücklich von der Zitat-Korrektur aus — also Folge-ADR
mit `Supersedes` in der engen Klausel, dasselbe Werkzeug, mit dem
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) die
dritte Fitness-Function-Zeile derselben ADR beerbt hat. Kein
Reviewer→Implementer-Pfeil: es gibt keine Fixrunde am Code.

## Entscheidung

Wir wählen: **Der Schlusssatz des dritten Bullets der §Entscheidung
Festlegung 2 aus
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) wird
ersetzt. Er nennt die Zugangsdaten-Klasse mit ihren sechs Schlüsseln.**

Drei Festlegungen:

1. **Die Entscheidung der
   [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
   gilt unverändert.** Der Diskriminator („kann dieses Feld Zugangsdaten
   tragen?"), die neun zulässigen Datei-Felder, `http_addr`/`grpc_addr` als
   Datei-Felder, `nats_url` und die zwei Token-Schlüssel env-exklusiv mit
   eigener, den Grund nennenden Fehlerzeile und die Durchleitung der fünf
   Oberflächen-Variablen bleiben in Kraft. Diese ADR ändert **keine**
   Festlegung, keinen Code, keine Spec-Stelle und kein Handbuch — sie
   korrigiert eine Zahl in einem Satz ihres §Entscheidung.
2. **Der Ersatztext des Schlusssatzes** lautet — das Bullet bleibt im Übrigen
   unverändert:

   > dazu die **sechs** env-exklusiven Zugangsdaten-Schlüssel — `capture_dsn`,
   > `admin_dsn`, `reader_dsn`, `api_token_reader`, `api_token_admin`,
   > `nats_url` (Festlegung 1).

   Die Aufzählung ersetzt die Zahl nicht, sie trägt sie: die Zahl ist an der
   Stelle nachzählbar, an der sie gelesen wird (`AGENTS.md` §3.12 Instanz A —
   eine Zahl trägt ihren Ursprung). Unberührt bleiben die neun zulässigen
   Felder des Bullets und der Satz, dass diese Aufzählung die Liste aus
   [`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md) Festlegung 6
   ersetzt.
3. **Die zwei Größen bleiben getrennt lesbar.** *Fünf* meint die
   durchgereichten Oberflächen-**Variablen** (Festlegung 3), *sechs* die
   Zugangsdaten-**Schlüssel** (Festlegung 1). Beide Zahlen sind wahr; sie
   zählen verschiedene Mengen (Kontext (3)).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: „fünf" bleibt stehen | kein Eingriff an einer `Accepted`-ADR; der ausgelieferte Stand ist richtig | die Klausel, die ihre Liste als **vollständig** ausgibt, zählt sie eine zu klein; wer die Klasse erweitert oder auditiert, nimmt eine falsche Zahl; kein Sensor sieht sie (gemessen, V-1) — die nächste Leserin müsste den Widerspruch erneut herleiten |
| **B — die Zahl auf sechs ziehen und die sechs Namen an der Stelle nennen (gewählt)** | die Zahl ist an der Stelle nachzählbar, an der sie gelesen wird; wer die Klasse erweitert, hat die Liste vor sich; kein Code-, Spec-, Handbuch- oder Sensor-Eingriff — die drei Träger führen sechs (gemessen) | die Klassen-Liste steht danach zweimal in derselben ADR (Festlegung 1 und Festlegung 2); beim nächsten Zuwachs sind beide Stellen nachzuziehen |
| C — die Zahl streichen und auf Festlegung 1 zeigen („dazu die env-exklusiven Zugangsdaten-Schlüssel der Festlegung 1") | **eine** Liste in der ADR — die kleinste Drift-Fläche; die Klasse bleibt vollständig benannt | die zählbare Aussage verschwindet: wer die Klasse erweitert, sieht an der Stelle keine Zahl, die ihm sagt, dass die Liste vollständig war; das Bullet daneben zählt seine neun Felder — die zwei Hälften desselben Bullets trügen verschiedene Formen |
| D — die Zahl auf sechs ziehen, ohne die Namen zu nennen | kleinster Eingriff in den Text | die Zahl bleibt eine Behauptung ohne ihre Liste an der Stelle — genau die Form, die den Befund erzeugt hat; nachzählbar ist sie nur, wenn der Leser Festlegung 1 daneben aufschlägt (`AGENTS.md` §3.12 Instanz A) |

**Fazit:** B. A lässt eine falsche Zahl an einer Stelle stehen, die ihre eigene
Vollständigkeit behauptet; bei C verschwände die zählbare Aussage, und die zwei
Hälften desselben Bullets trügen zwei Formen; D ließe die Zahl eine Behauptung
bleiben.

## Konsequenzen

- Positiv: Die Klausel, die die Klassen-Liste als vollständig ausgibt, zählt
  richtig — und die Zahl ist an ihrer Stelle nachzählbar.
- Positiv: **kein Code-, Spec- oder Handbuch-Eingriff.** `SPEC-016`, das
  Handbuch §5.2 und `forbiddenFileCredentialKeys` führen je sechs Schlüssel
  (gemessen, Kontext (2)); der ausgelieferte Stand war richtig, falsch war
  allein der Entscheidungstext.
- Positiv: Die Korrektur folgt dem gebauten Muster für eine enge
  Klausel-Korrektur (`ADR-0048`, `ADR-0062`, `ADR-0063`, `ADR-0065`,
  `ADR-0067`, `ADR-0089`); kein zweites Werkzeug, kein neuer Träger.
- Negativ: [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
  und diese ADR müssen in **einem Satz** zusammengelesen werden.
- Negativ mit benannter Grenze: Die Liste der Klasse steht danach **zweimal**
  in derselben ADR (Festlegung 1 und der ersetzte Satz). Zwei Stellen können
  gegen die Messung driften — genau die Klasse, die dieser Vorgang korrigiert;
  die Alternative C hätte die eine Stelle gehabt und die zählbare Aussage
  gekostet.
- Hinweis (Index): Die `ADR-0088`-Zeile des ADR-Index trägt **keinen**
  Vorwärts-Zeiger auf diese ADR — die `structure`-Deckelung der `Titel`-Zelle
  bei 80 Zeichen lässt ihn nicht zu (am Zellinhalt gemessen: **79** Zeichen;
  mit dem Zusatz `; → ADR-0091` **91**, mit dem im Repo beim Paar
  `ADR-0087` ← `ADR-0090` gebrauchten Tausch `, teilw.` → `; → ADR-0091`
  **83**). Die Nachfolge trägt die Zeile **dieser** ADR
  (`Supers. ADR-0088, teilw.`) — dieselbe Lage wie bei den Paaren
  `ADR-0083`/`ADR-0086` und `ADR-0057`/`ADR-0081`, die ebenfalls keinen
  Vorwärts-Zeiger führen.
- Folgepflicht (Planner-Zug): **keine** — kein DoD eines Slice lehnt sich an
  die ersetzte Zahl an; der gelieferte Stand trägt sechs (gemessen).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (kein Sensor) | **Die Zugangsdaten-Klasse trägt sechs Schlüssel** — `forbiddenFileCredentialKeys` führt sechs Einträge, `SPEC-016` und das Handbuch §5.2 nennen dieselben sechs — **von Hand** nachzuzählen. Der Wächter ist der **Review** und der **Verifier** ([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md): die Träger-Paarung trägt kein Werkzeug) | — |
| `go test ./internal/bootstrap/...` (im gepinnten Container, netzlos) | **Die Bindung je Schlüssel** trägt unverändert die **zweite** Zeile der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-§Fitness Function: jeder der sechs Klassenschlüssel in der Datei endet in `ErrConfiguration` mit einer eigenen, den Grund nennenden Zeile | `make test` (kein Gate) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Drei benannte Trigger, sonst permanent:

1. **Ein zugangsdaten-tragendes Feld tritt in die Verdrahtung** (Trigger 1/2
   der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)).
   Dann wachsen Klasse und Zahl zusammen; der Nachzug trifft **beide** Stellen
   derselben ADR (Festlegung 1 und den ersetzten Satz) samt `SPEC-016`,
   Handbuch und Code.
2. **Der Klassen-Satz eines Trägers wird umformuliert** — in `SPEC-016`, im
   Handbuch §5.2 oder im Code-Kommentar. Dann ist die Zahl hier gegen den
   neuen Trägerstand nachzumessen, bevor sie zitiert wird.
3. **Die Paarung driftet real auseinander**, und ein Review oder ein Verifier
   findet es. Dann ist zu entscheiden, ob der Befund bei der benannten Grenze
   bleibt ([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md),
   Trigger 2) oder ob er einen Sensor belegbar macht — die Beobachtung dafür
   ist der Beleg, nicht der Vorsatz.

Sonst permanent: die Zugangsdaten-Klasse der Konfigurationsdatei umfasst
**sechs** Schlüssel — die drei DSN-Schlüssel, die zwei Token-Schlüssel und
`nats_url`; die fünf Oberflächen-**Variablen** der Festlegung 3 bleiben davon
unberührt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: Verifikationsbericht zu `slice-096` <!-- d-check:status-provenance -->, V-1 (MEDIUM): die §Entscheidung Festlegung 2 der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) zählt die Zugangsdaten-Klasse als „fünf". Eigene Messung an HEAD `825d5fd`: **sechs** in `forbiddenFileCredentialKeys`, in `SPEC-016`, im Handbuch §5.2 und in `ADR-0088` Festlegung 1. Unabhängiger Architect-Zug ersetzt den Schlusssatz des Bullets (sechs, mit der Namensliste); die Fünf der Festlegung 3 meint die Oberflächen-**Variablen** und bleibt richtig | `docs/reviews/verify-slice-096.md` <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
