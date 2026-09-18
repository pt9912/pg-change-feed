# ADR-0102: Zugangsdaten-Klasse — Supersede-Liste vervollständigt

**Status:** Accepted — Supersedes [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)
in **einer** Klausel: deren **§Status-Supersede-Erklärung** — der Aufzählung
der ersetzten Stellen **und** der bestätigenden Klausel „Alles Übrige dieser
drei Dokumente bleibt **hiermit bestätigt**". Die Aufzählung nennt nicht jeden
Satz, der die Klasse noch als *sechs* führt; die bestätigende Klausel hätte die
übrigen als geltend bestätigt. An ihre Stelle tritt die **vollständige**
Aufzählung unten, geschlossen durch die Regel in §Entscheidung Festlegung 2.

Alles Übrige der [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)
bleibt bestätigt und wird nicht wiederholt: die geltende Klasse mit ihren
**sieben** Schlüsseln, die Trennung der **fünf** Oberflächen-Variablen von den
sieben Zugangsdaten-Schlüsseln, der Ort der geltenden Zahl, §Verglichene
Alternativen, §Konsequenzen, §Fitness Function, §Re-Evaluierungs-Trigger und
§Geschichte.

**Zahl-Aussagen, die ersetzt werden** — kraft Regel (§Entscheidung
Festlegung 2) und namentlich an diesen verifizierten Stellen:

- aus [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md):
  deren §Entscheidung-Kopfsatz („`nats_url` und die **zwei** Token-Schlüssel
  bleiben env-exklusiv"), deren §Konsequenzen-Bullet „Folgepflicht (Code-Zug)"
  („drei DSN- plus **zwei** Token-Schlüssel plus `nats_url`") und deren
  §Kontext (4) („Die **zwei** Token-Schlüssel werden heute abgewiesen");
- aus [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md): deren
  §Entscheidung Festlegung 1 („`nats_url` und die **zwei** Token-Schlüssel
  env-exklusiv"), deren §Kontext (1) („er führt **sechs** Schlüssel"),
  §Kontext (3) („fünf *Variablen*, **sechs** *Schlüssel*" und „Gemeinsam sind
  drei: `nats_url` und die **zwei** Token-Schlüssel") und §Kontext (4)
  („deckt dann fünf statt **sechs** ab");
- aus [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md):
  deren §Entscheidung Festlegung 1 („`nats_url` und die **zwei**
  Token-Schlüssel env-exklusiv") und deren §Kontext (2) („`SPEC-016` führt
  neun zulässige Feld-Schlüssel und **sechs** env-exklusive");
- die **bestätigenden** Klauseln der §Status-Abschnitte, soweit sie eine der
  ersetzten Stellen mitbestätigen — namentlich die Bestätigung der
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)-Festlegung 1
  in der §Status der
  [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md).

**Aufzeichnungen, die stehen bleiben** — sie tragen den Stand ihres Laufs bzw.
der damaligen Abwägung: der **Titel** der
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) samt seiner
Index-Zelle, die §Geschichte-Zeilen der vier Dokumente, die datierten
Messzeilen und Nachweis-Tabellen der §Kontext-Abschnitte und die
§Verglichene-Alternativen-Zeilen
([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) Option C,
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) Optionen B und D,
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) Option B).

**Datum:** 2026-09-18

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf den
Re-Review zu `slice-nats-drittstream-core` <!-- d-check:status-provenance -->,
Befunde F-4 und F-12; anderer Kontext als der Architect-Lauf, der die
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) schrieb, und als
der Re-Review-Lauf, der die Befunde erhob — Modul 8 §Rollen-Regeln:
„Architect schreibt")

**Bezug:** [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)
(§Status-Supersede-Erklärung superseded; alles Übrige bestätigt) ·
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) ·
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) ·
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) ·
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) (die vier
Dokumente, deren Zahl-Aussagen über die Klasse ersetzt werden) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
(§Status ist keine Zitat-Korrektur-Fläche; §Entscheidung und §Konsequenzen
sind von ihr ausgenommen) ·
[`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md)
(das gebaute Muster: auch eine einzelne Zeile wird per Folge-ADR korrigiert) ·
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) §Status
(„eine Klausel und die drei Stellen, die deren Aussage wörtlich wiederholen" —
das Muster für eine Klausel samt ihren Wiederholungen) · `AGENTS.md` §3.5 ·
§3.12 · §3.13 · Re-Review zu `slice-nats-drittstream-core` <!-- d-check:status-provenance -->
(F-4, F-12 samt eigener Messung) · [`SPEC-016`](../../../spec/pflichtenheft.md)

**Schärft:** — (Korrektur-ADR ohne Spec-Stratum: das Pflichtenheft führt die
Klasse bereits mit ihren sieben Schlüsseln und ist nicht Gegenstand dieser ADR,
wie [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) und
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass — die Supersede-Erklärung trägt ihren Umfang nicht.** Die
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) hat die
Zahl-Aussagen der
[`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) und
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
namentlich und klauselweise auf die **sieben** gezogen — die Substanz ist
richtig und bleibt. Ihre §Status-Aufzählung ist jedoch unvollständig: sie
nennt nicht jeden Satz der drei Dokumente, der die Klasse noch als *sechs*
führt oder eines ihrer Elemente „die zwei Token-Schlüssel" nennt. Der
Re-Review zu `slice-nats-drittstream-core` <!-- d-check:status-provenance -->
weist das als F-4 (HIGH) aus, und die bestätigende Klausel der
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) verschärft es: sie
bestätigt mit „Alles Übrige" genau die nicht genannten Stellen als geltend.

**(2) Die eigene Messung dieses Zugs** (2026-09-18, HEAD `2fc3499`; `grep` über
die vier `Accepted`-Dokumente nach „sechs Schlüssel" / „zwei Token-Schlüssel" /
„er führt sechs" und Lesen jeder Fundstelle im Zusammenhang):

| Dokument | Stelle (Sektion, Zeile am HEAD) | Trägt heute | Nach dieser ADR |
|---|---|---|---|
| [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) | §Entscheidung-Kopfsatz (:136–137); §Konsequenzen „Folgepflicht (Code-Zug)" (:314); §Kontext (4) (:100) | „zwei Token-Schlüssel"/„sechs" | **ersetzt** |
| [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) | §Entscheidung Festlegung 1 (:124); §Kontext (1)/(3)/(4) (:62, :80, :88, :97) | „zwei Token-Schlüssel"/„sechs" | **ersetzt** |
| [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) | §Entscheidung Festlegung 1 (:145); §Kontext (2) (:73) | „zwei Token-Schlüssel"/„sechs" | **ersetzt** |
| [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) | §Entscheidung Festlegung 1, §Konsequenzen, §Fitness Function Zeile 1, §Kontext (2) Schlusssatz | „sechs" | bereits durch [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) ersetzt; §Status-Bestätigung der [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md)-Festlegung 1 hier mitgezogen |
| die vier | Titel, §Geschichte, datierte Messzeilen der §Kontext-Abschnitte, §Verglichene Alternativen | „sechs" | **Aufzeichnung**, bleibt stehen |

Die vier Dokumente tragen die alte Zahl, nicht zwei — die §6-Meldung des
laufenden Slice nennt sie als „**zwei** Dokumente" (Re-Review F-12) und wird im
selben Commit nachgezogen. Der **heute tragende** Träger der Klasse ist die
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md); die
[`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) ist durch
sie teilweise supersediert, trägt aber in ihren nicht ersetzten Stellen die
alte Zahl weiter.

**(3) Warum eine Aufzählung die Klasse nicht schließt.** Die
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) hat die Klasse als
**Liste** von Stellen gefasst und alles Übrige pauschal bestätigt. Eine Liste
kann nur nennen, was ihr Verfasser gesehen hat; was sie nicht nennt, bestätigt
die pauschale Klausel — die Lücke wird also nicht nur übersehen, sie wird
**geltend gemacht**. Genau diese Form ist der Defekt, den F-4 misst: die
Nachbar-Präzedenz `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` führt
dieselbe Klasse. Die Korrektur muss deshalb nicht nur die fehlenden Stellen
nennen, sondern die **Klasse** schließen — durch eine Regel über einen
definierten Abschnitts-Scope (Festlegung 2), deren namentliche Aufzählung nur
die verifizierte Instanz ist.

**(4) Die In-place-Form ist ausgeschlossen — und die Lebenslage der
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) ändert das nicht.**
Der Fehler sitzt in deren §Status. `AGENTS.md` §3.5 erklärt §Status (neben
§Entscheidung, §Konsequenzen, §Verglichene Alternativen und der
`Supersedes`-Kette) für **unberührbar**; die **Zitat-Korrektur** der
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
deckt nur das Zitat- und Verweisgerüst bei **unverändertem Referenten**. Hier
änderte sich der Referent: **welche** Klauseln als ersetzt gelten, ist der
Gegenstand. Damit ist die Erweiterung der Aufzählung eine inhaltliche Änderung
und keine Zitat-Korrektur. Die Lebenslage der
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) — wenige Minuten
vor diesem Zug committet (`05f9d52`), allein vom Index referenziert, nie als
eigenes Artefakt übergeben oder reviewt — öffnet keinen anderen Weg: §3.5
knüpft die Immutabilität an den **Status `Accepted`**, den die
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) in ihrer §Status
selbst erklärt und der auf `main` gelandet und indexiert ist; eine
Review-Pflicht als Vorbedingung kennt §3.5 nicht. Dieselbe Lage hatten die
[`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) und die
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) — beide
waren Architect-Antworten auf einen Befund, beide wurden für ihren eigenen
Defekt per Folge-ADR korrigiert. Das gebaute Muster für eine einzelne Zeile
ist [`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md);
für eine Klausel samt ihren Wiederholungen
[`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) §Status.

**(5) Warum das eine Architect-Frage ist.** Die Defekte liegen in
`Accepted`-Texten, nicht im gelieferten Stand: die drei lebenden Träger stimmen
mit dem Code überein (gemessen, [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) §Kontext (2)). Es gibt keine
Fixrunde am Code, also auch keinen Reviewer→Implementer-Pfeil (Modul 8
§Konflikt-Pfad, Verdikt 1).

## Entscheidung

Wir wählen: **Die geltende Klasse bleibt sieben Schlüssel
([`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) §Entscheidung
Festlegung 1). Neu gefasst wird allein die Reichweite der Supersession: der
Ersatztext der [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)
gilt ab jetzt an jeder Zahl-Aussage über die Zugangsdaten-Klasse der vier
Dokumente, die keine Aufzeichnung ist; kein Code-, Spec-, Handbuch- oder
Sensor-Eingriff.**

Vier Festlegungen:

1. **Die geltende Klasse** — die drei DSN-Schlüssel (`capture_dsn`,
   `admin_dsn`, `reader_dsn`), die drei Token-Schlüssel (`api_token_reader`,
   `api_token_admin`, `nats_stream_token`) und `nats_url`. Die drei lebenden
   Träger führen sie vollständig (gemessen, [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) §Kontext (2)); diese ADR
   ändert keine Zahl des geltenden Standes, sondern die Reichweite seiner
   Supersession.
2. **Der Ersatztext gilt kraft Regel, nicht nur kraft Aufzählung.** Der in
   [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) §Entscheidung
   Festlegung 2 formulierte Satz gilt an **jeder** Stelle der vier Dokumente,
   die die Zugangsdaten-Klasse heute als *sechs* führt oder eines ihrer
   Elemente „die zwei Token-Schlüssel" nennt — einschließlich der
   §Status-Bestätigungen, die eine solche Stelle mitbestätigen —, ausgenommen
   allein die in Festlegung 3 genannten **Aufzeichnungen**. Die namentliche
   Aufzählung im §Status ist die **verifizierte Instanz** dieser Regel, nicht
   ihre Grenze: eine weitere Stelle derselben Form ist ersetzt, auch wenn
   diese ADR sie nicht einzeln nennt. Das ist der Unterschied zur
   [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)-Aufzählung, an
   der der Befund F-4 hing — und zugleich die Bedingung, unter der die
   bestätigende Klausel im §Status dieser ADR stehen darf.
3. **Aufzeichnung bleibt Aufzeichnung.** Nicht nachgezogen werden: der
   **Titel** der
   [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) samt seiner
   Index-Zelle (der Name eines Dokuments ist keine Aussage), die
   §Geschichte-Zeilen, die **datierten** Messzeilen und Nachweis-Tabellen der
   §Kontext-Abschnitte (sie nennen ihren Lauf) und die
   §Verglichene-Alternativen-Zeilen (sie tragen die damalige Abwägung). Sie
   sind kein Beleg über heute (`AGENTS.md` §3.12 Instanz A); wer die Klasse
   auditiert, liest sie als Historie, nicht als geltend.
4. **Der Ort der geltenden Zahl.** Sie steht weiterhin in §Entscheidung und
   §Fitness Function der
   [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md); der Verweis
   ihrer §Entscheidung Festlegung 2 auf „die im §Status namentlich genannten
   Stellen" liest ab jetzt gegen die Aufzählung **dieser** ADR. Die Paarung
   bleibt ein Urteil **von Hand** — kein Werkzeug rotet, wenn sie in einem Lauf
   unterbleibt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: die Aufzählung der [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) bleibt, wie sie steht | kein neuer Träger; die Substanz (sieben Schlüssel) ist richtig | die nicht genannten Stellen führen die Klasse weiter als *sechs*, und die bestätigende Klausel macht sie geltend; der Verifier nimmt eine Zahl, die sein eigenes Nachzählen widerlegt, ohne zu wissen, ob er eine Code-Drift oder eine überholte Aussage sieht (`AGENTS.md` §3.12 Instanz A) — F-4 bleibt offen |
| B — die Aufzählung in der [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) in-place erweitern | kleinster Eingriff in den Text | §3.5 erklärt §Status für unberührbar, und die Zitat-Korrektur der [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1 deckt nur das Verweisgerüst bei unverändertem Referenten — hier ändert sich der Referent (welche Klauseln ersetzt sind); das Muster [`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md) zeigt dieselbe Antwort für eine einzelne Zeile |
| C — **die Aufzählung durch eine Regel über einen definierten Abschnitts-Scope ersetzen, die fehlenden Stellen namentlich als Instanz nennen (gewählt)** | die Klasse ist **geschlossen**: jede Zahl-Aussage derselben Form ist ersetzt, auch die noch unbemerkte; die Aufzeichnungen bleiben benannt stehen; kein Code-, Spec-, Handbuch- oder Sensor-Eingriff; die Substanz der [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) bleibt unangetastet | eine **Regel** statt einer endlichen Liste: wer eine Stelle einordnet, muss entscheiden, ob sie Zahl-Aussage oder Aufzeichnung ist — die Grenze ist in Festlegung 2/3 benannt, aber sie bleibt ein Urteil (Konsequenzen) |
| D — nur die fehlenden Stellen namentlich an die Liste hängen, ohne Regel | kleinste Supersession; die konkreten Fundstellen sind benannt | wiederholt genau die Klasse, an der F-4 hing: die Liste kann nur nennen, was dieser Zug gesehen hat, und die bestätigende Klausel macht den Rest geltend; der nächste Zuwachs erzeugt dieselbe Lücke neu |

**Fazit:** C. A lässt eine falsche Zahl an Stellen stehen, die ihre eigene
Geltung behaupten; B verbietet §3.5; D wiederholt die Aufzählungs-Lücke, die
den Befund erzeugt hat.

## Konsequenzen

- Positiv: Die Zahl-Aussagen über die Zugangsdaten-Klasse sind in den vier
  Dokumenten vollständig auf die **sieben** gezogen; die bestätigende Klausel
  pauschalisiert nicht mehr, sondern benennt die Aufzeichnungen als
  Aufzeichnungen.
- Positiv: **kein Code-, Spec-, Handbuch- oder Sensor-Eingriff.** Die drei
  lebenden Träger stimmen mit dem Code überein (gemessen, [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)
  §Kontext (2)); falsch war allein die Reichweite der Supersession.
- Positiv: Die Korrektur folgt dem gebauten Muster für eine Klausel samt ihren
  Wiederholungen ([`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
  §Status) und für eine einzelne Zeile
  ([`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md));
  kein zweites Werkzeug, kein neuer Träger.
- Negativ: Die Klassen-Aussage ist jetzt über **sechs** `Accepted`-Dokumente
  verteilt — [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md),
  [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md),
  [`ADR-0092`](0092-feldmengen-paarung-reichweite-der-drei-traeger.md),
  [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) und diese; wer die
  Klasse auditiert, liest sie zusammen.
- Negativ mit benannter Grenze: Die Abgrenzung Zahl-Aussage ↔ Aufzeichnung ist
  eine **Regel**, kein Zähler. Ob ein Satz eine Aussage über heute oder eine
  Aufzeichnung ist, entscheidet ein Leser von Hand; **kein** Sensor rotet, wenn
  er falsch einordnet (so schon
  [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) und
  [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) §Konsequenzen).
  Die benannte Grenze ist zugleich der Re-Evaluierungs-Trigger 5.
- Hinweis (Index): Die [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md)-Zeile des ADR-Index behält ihren Wortlaut
  `(Supers. ADR-0088/0091/0092, teilw.)` — ihn nennt deren §Konsequenzen
  ausdrücklich, und die `structure`-Deckelung der `Titel`-Zelle bei 80 Zeichen
  lässt einen Vorwärts-Zeiger `; → ADR-0102` nur zu, wenn genau dieser
  Wortlaut geändert würde (dieselbe Lage wie bei
  [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md), die
  aus demselben Grund keinen Zeiger trägt). Den Teil-Charakter und die
  Nachfolge trägt ab jetzt die Zeile **dieser** ADR
  (`Supers. ADR-0101, teilw.`).
- Folgepflicht (Planner-Zug): die §6-Meldung des laufenden Slice
  (`slice-nats-drittstream-core` <!-- d-check:status-provenance -->) zählt die
  Träger der alten Zahl auf **vier** (F-12) und nennt diese ADR als die
  eingelöste Vervollständigung — im selben Commit nachgezogen; darüber hinaus
  **keine** Folgepflicht: kein DoD eines Slice lehnt sich an eine der ersetzten
  Zahlen an.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (kein Sensor) | **Die Zugangsdaten-Klasse trägt sieben Schlüssel, und keine Zahl-Aussage der vier Dokumente führt sie als sechs** — `forbiddenFileCredentialKeys` führt sieben Einträge, [`SPEC-016`](../../../spec/pflichtenheft.md) und `docs/user/benutzerhandbuch.md` §5 nennen dieselben sieben; die vier `Accepted`-Dokumente nennen sie an keiner Stelle mehr als sechs — **von Hand** nachzuzählen. Der Wächter ist der **Review** und der **Verifier** ([`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md): die Träger-Paarung trägt kein Werkzeug) | — |
| `go test ./internal/bootstrap/...` (im gepinnten Container, netzlos) | **Die Bindung je Schlüssel** trägt unverändert die **zweite** Zeile der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)-§Fitness Function in der durch [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) ersetzten Enumeration: jeder der sieben Klassenschlüssel in der Datei endet in `ErrConfiguration` mit einer eigenen, den Grund nennenden Zeile (`config_file_internal_test.go` iteriert über die sieben) | `make test` (kein Gate) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Die vier benannten Trigger der
[`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) bleiben in Kraft
und werden hier nicht wiederholt; **einer kommt hinzu** — er zählt in deren
Nummerierung weiter:

5. **Eine Zahl-Aussage über die Klasse ist zwischen Aussage und Aufzeichnung
   strittig** — ein Review oder ein Verifier ordnet eine Stelle anders ein als
   Festlegung 2/3. Dann ist die Grenze hier nachzuschärfen, bevor diese ADR
   zitiert wird; der Suchlauf dafür ist ein `grep` über die Träger nach der
   alten Zahl, nicht über den eigenen Diff (`AGENTS.md` §3.13).

Sonst permanent: die Zugangsdaten-Klasse der Konfigurationsdatei umfasst
**sieben** Schlüssel — die drei DSN-Schlüssel, die drei Token-Schlüssel und
`nats_url`; und keine Zahl-Aussage der vier Dokumente führt sie als sechs.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — Anlass: Re-Review zu `slice-nats-drittstream-core` <!-- d-check:status-provenance -->, F-4 (HIGH) und F-12 (HIGH). Eigene Messung an HEAD `2fc3499` (`grep` über die vier `Accepted`-Dokumente nach „sechs Schlüssel"/„zwei Token-Schlüssel"/„er führt sechs"): die §Status-Aufzählung der [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) nennt den §Entscheidung-Kopfsatz und den §Konsequenzen-Bullet der [`ADR-0088`](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md), deren §Kontext (4), die §Entscheidung Festlegung 1 und die §Kontext (1)/(3)/(4) der [`ADR-0091`](0091-zugangsdaten-klasse-sechs-schluessel.md) sowie die §Entscheidung Festlegung 1 und die §Kontext (2) der [`ADR-0089`](0089-feldmengen-paarung-kein-sensor-review-waechter.md) nicht; die bestätigende Klausel hätte sie geltend bestätigt. Unabhängiger Architect-Zug ersetzt die Aufzählung durch die verifizierte Liste samt der Regel, die die Klasse schließt; die Substanz der [`ADR-0101`](0101-zugangsdaten-klasse-sieben-schluessel.md) und die drei lebenden Träger bleiben unverändert | Re-Review zu `slice-nats-drittstream-core` <!-- d-check:status-provenance --> F-4/F-12 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
