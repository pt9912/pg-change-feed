# ADR-0141: Backfill-Run — Lesefehler eines Regelstands endet mit der Klasse der Ursache, `configuration` bleibt dem Wechsel des Standes vorbehalten (Schärft ADR-0139)

**Status:** Accepted

**Datum:** 2026-10-01

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch die Architect-Frage
F-N1 im Re-Review `review-slice-routing-backfill-pfad-fixrunde-1` und den Befund V-2
der Verifikation `verifikation-slice-routing-backfill-pfad`)

**Bezug:** [`LH-FA-CFG-008`](../../../spec/lastenheft.md) (Routing),
[`LH-FA-CFG-005`](../../../spec/lastenheft.md) (Ausschluss),
[`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Transformationen),
[`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill, Fail-closed),
[`LH-QA-SEC-004`](../../../spec/lastenheft.md),
[ADR-0139](0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
(Festlegung 1 — in **einem** Satz berichtigt, siehe Entscheidung),
[ADR-0111](0111-backfill-bestand-snapshot-bulk-copy.md) (Teilfrage 4),
[ADR-0113](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[ADR-0117](0117-backfill-run-fehlerklasse-schema.md) (Festlegung 5),
[ADR-0118](0118-backfill-umschreiben-im-snapshot-fenster.md),
[ADR-0023](0023-fehlerklassifikation.md)

**Schärft:** [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) (Absatz
„Fail-closed vor dem Commit“, Satz „Ein nicht lesbarer Stand gilt als
Abweichung“), [`SPEC-008`](../../../spec/pflichtenheft.md) (Zeilen
`configuration` und `storage`, Lesart im Backfill-Run)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0139` Festlegung 1 sagt für den Routing-Regelstand im Run, „derselbe
Mechanismus wie für Ausschluss- und Transformationsstand“, und: „Eine
Abweichung — auch ein nicht lesbarer Stand — rollt den Run zurück: `failed`,
Klasse `configuration`“. Die Spec trägt das als Satz „Ein nicht lesbarer Stand
gilt als Abweichung“ im Absatz, der den Routing-Regelstand beschreibt.
Review und Verifikation von `slice-routing-backfill-pfad` fanden, dass die zwei
Aussagen am Bestand nicht zusammenpassen.

### Befunde (Stand HEAD `8909ed4d`)

Gelesen am Quelltext, **nicht gefahren**:

- `internal/application/usecase/backfill/service.go`: `routingRules` wickelt den
  Lesefehler zusammen mit `ErrRoutingStateChanged`; `classifyError` bildet diesen
  auf `configuration` ab. Die Lesungen von Ausschluss- (`excludedColumns`) und
  Transformationsstand (`transformationRules`) geben den Lesefehler unverändert
  zurück; der Run endet mit der Klasse der Ursache (`storage` für
  `outbound.ErrStorage`).
- Bestands-Tests binden das Verhalten: `TestExecuteExclusionReadFailure`
  (`service_test.go`) und `TestExecuteRuleReadFailure` (`transformation_test.go`)
  erwarten Präfix `storage: ` samt Ursachentext an jeder Lesestelle;
  `TestExecuteRoutingReadFailureEndsRun` (`routing_test.go`) erwartet
  `configuration: `.
- `classifyError` prüft `configuration` vor `transient`/`replication`/`storage`.
  Ein gewickeltes `ErrRoutingStateChanged` schlägt deshalb jede Ursachenklasse
  außer `permission`: ein vorübergehend nicht erreichbarer Speicher beim Lesen
  des Routing-Stands endet als `configuration`.
- `ADR-0111` Teilfrage 4 verlangt für „jede Abweichung“ Rollback und Grund im
  `error_message`, nennt aber keine Klasse für einen Lesefehler.
  `ADR-0117` Festlegung 5 vergibt `configuration` für den **Wechsel** des
  Standes („der Zustand wechselt“). Das Lastenheft nennt für den Lesefehler
  keine Klasse (Suchlauf über `spec/lastenheft.md` nach „Abweichung“,
  „fail-closed“, „nicht lesbar“: kein Treffer) — der Wortlaut für
  `configuration` liegt allein in `ADR-0139` und dem Spec-Satz, die beide am
  1. Oktober 2026 entstanden sind.
- `SPEC-008` (Zeilen gelesen): `configuration` = „ungültige/falsch gesetzte
  Konfiguration“, Aktion „sichtbarer Fehler; kein Start und keine Fortsetzung im
  falschen Stand“; `storage` = „Persistenzfehler“; `transient` = Erneut versuchen
  mit begrenztem Backoff. Im Run gibt es keinen automatischen Neustart
  (`interrupted` und `failed` starten nur durch einen neuen Antrag); die Klasse
  steuert dort die **Diagnose** des Betreibers, nicht eine Maschinen-Aktion.

## Entscheidung

Wir wählen **Option C: ein nicht lesbarer Regelstand ist im Run fail-closed
(Rollback, `failed`, kein Commit, Ursache im Fehlertext), seine Klasse ist die
Klasse der Ursache — für alle drei Regelstände gleich (Ausschluss,
Transformationen, Routing). `configuration` gilt im Run dem Wechsel des Standes
(`ErrExclusionStateChanged`, `ErrTransformationStateChanged`,
`ErrRoutingStateChanged`) und der fehlenden Bindung, nicht einem Lesefehler.**

- **Was von `ADR-0139` Festlegung 1 gilt.** Alles außer einem Halbsatz: Stand
  nach dem Öffnen des Snapshots, Mengenvergleich je Block und vor dem Commit,
  Reihenfolge `schema` vor `configuration`, Abhilfe, Grenze. „Derselbe
  Mechanismus“ trifft jetzt zu — er ist dieser: Lesen, vergleichen, bei Lesefehler
  oder Abweichung abbrechen.
- **Was berichtigt wird.** Die Wendung „auch ein nicht lesbarer Stand“ in
  Verbindung mit „Klasse `configuration`“ in `ADR-0139` Festlegung 1 (Punkt
  „Vergleich“) gilt nicht mehr; an ihre Stelle tritt diese Festlegung. Weil
  `ADR-0139` `Accepted` und unberührbar ist (`AGENTS.md` §3.5), steht die
  Berichtigung hier und nicht dort; der Rest von `ADR-0139` bleibt unverändert in
  Kraft. Wer `ADR-0139` Festlegung 1 liest, liest sie mit dieser ADR.
- **Begründung.**
  1. *Ehrlichkeit der Klasse (`ADR-0023`).* Ein Speicherfehler beim Lesen der
     `applied`-Zeilen ist ein Persistenz- oder Transportfehler, keine „ungültige
     oder falsch gesetzte Konfiguration“. `configuration` würde den Betreiber zur
     Konfiguration schicken statt zur Datenbank.
  2. *Keine Information verlieren.* Die Klasse der Ursache trägt die
     Unterscheidung `transient` (vorübergehend, wiederholbar) gegen `storage`
     und `permission` (Rolle auf `cdc.administration_request` fehlt). Die
     Abbildung auf `configuration` löscht sie (Befund oben,
     `classifyError`-Reihenfolge).
  3. *Fail-closed ist unabhängig von der Klasse.* Die Sicherheitszusage
     (`LH-QA-SEC-004`: kein ausgeschlossener Wert wird serialisiert, kein
     Mischstand wird committet) hängt am Abbruch vor dem Commit, nicht an der
     Klassenbezeichnung; beide Klassen sind laut `SPEC-008` „sichtbarer Fehler,
     kein stiller Retry“.
  4. *Kleinste Änderung am Bestand.* Zwei der drei Stände tun es bereits, ihre
     Tests binden es; ein Satz in der Spec und eine Anweisung im Routing-Pfad
     ziehen nach. Option A hätte zwei Lesepfade und zwei Testgruppen auf eine
     Klasse umgestellt, die ihrer Bedeutung widerspricht.
  5. *Kein Lastenheft-Widerspruch.* Das Lastenheft nennt keine Klasse für den
     Lesefehler (Befund oben); die Berichtigung berührt keinen Vertragstext.
- **Preis.** Die Spec-Wendung „gilt als Abweichung“ wird enger gelesen: die
  Abweichung ist der Abbruch, nicht die Klasse. Wer in einem Run einen Lesefehler
  des Routing-Stands sieht, liest `storage`/`transient`/`permission` mit dem
  Ursachentext, nicht `configuration`.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — „nicht lesbar = `configuration`“ für alle drei Stände | ein einheitlicher Satz, wörtlich wie `ADR-0139`/Spec | schickt Betreiber bei Speicherausfall zur Konfiguration; löscht `transient`/`permission`/`storage`-Unterscheidung für drei Lesepfade; ändert zwei Bestands-Lesepfade, zwei Testgruppen (`TestExecuteExclusionReadFailure`, `TestExecuteRuleReadFailure`, je 4 bzw. 5 Aufrufstellen), `classifyError`-Kommentar, Handbuch-Zeile `storage`/`configuration` in §6 „Fehlerklassen“ und ändert eine Aussage, die `ADR-0111` Teilfrage 4 offen ließ |
| B — Routing ist die Ausnahme (`configuration`), Bestand bleibt | keine Änderung an Ausschluss/Transformation; Spec-Wortlaut erfüllt | keine Sachbegründung, warum der Lesefehler des Routing-Standes anders zu lesen wäre als der gleichartige des Ausschlussstands; dauerhaft zwei Abbildungen für dieselbe Ursache, „derselbe Mechanismus“ bliebe falsch; Betreiber lernt zwei Lesarten |
| **C — Klasse der Ursache für alle drei, `configuration` für den Wechsel (gewählt)** | Klasse bleibt ehrlich; Bestand und seine Tests unverändert; Fail-closed unberührt; ein Mechanismus für alle | berichtigt einen Halbsatz einer `Accepted` ADR und einen Spec-Satz; Routing-Pfad und sein Test werden zurückgenommen (Fixrunde `1487be88` im Kern) |

## Konsequenzen

- Positiv: ein Mechanismus, eine Lesart für alle drei Regelstände; die Klasse
  eines Run-Fehlers sagt, was ausgefallen ist.
- Negativ: die Fixrunde `1487be88` wird im Kern zurückgenommen; die Aussage von
  `ADR-0139` ist an einer Stelle durch diese ADR überlagert — Leser müssen beide
  lesen.
- Folgepflicht (Stand, nicht Chronik; jeweils eine Zeile, kein eigener Slice nötig):
  - `internal/application/usecase/backfill/service.go`: `routingRules` gibt den
    Lesefehler mit der Ursache zurück (ohne `ErrRoutingStateChanged`); Godoc von
    `routingRules` und der Absatz in `classifyError` ziehen nach („nicht
    lesbarer Routing-Regelstand“ fällt dort weg). `internal/application/usecase/backfill/routing_test.go`:
    `TestExecuteRoutingReadFailureEndsRun` erwartet Präfix `storage: ` statt
    `configuration: ` (Eingabe, Lesestellen und Abbruchpunkt unverändert).
    Träger: die offene Fixrunde bzw. Closure von `slice-routing-backfill-pfad`;
    das ist eine Rücknahme im selben Slice, kein neuer Slice.
  - `spec/pflichtenheft.md`, `LH-FA-CAP-009.a`, Absatz „Fail-closed vor dem
    Commit“: der Satz „Ein nicht lesbarer Stand gilt als Abweichung.“ lautet
    sinngemäß: „Ein nicht lesbarer Stand beendet den Run (`failed`, kein
    Commit) mit der Klasse der Ursache.“ Zeilen etwa 257–258 (Stand HEAD
    `8909ed4d`). Träger: der Planner/Spec-Nachzug am nächsten Spec-Commit, vor
    der Closure von `slice-routing-backfill-pfad`; die Änderung an der
    Spec-Datei ist nicht Sache dieser Rolle in diesem Zug.
  - `harness/README.md` und `docs/user/benutzerhandbuch.md`: kein Nachzug, die
    Fehlerklassen-Tabelle (§6 „Fehlerklassen“, `configuration` und `storage`)
    behauptet nichts über Lesefehler von Regelständen (Zeilen 1874, 1877
    gelesen). `slice-routing-betriebsdoku` und `slice-routing-e2e` sind nicht
    betroffen, soweit sie keinen Lesefehler-Fall mit Klasse `configuration`
    behaupten — *zu prüfen* durch den jeweiligen Slice-Suchlauf (`configuration`
    zusammen mit „nicht lesbar“), nicht hier gemessen.
  - Plan `slice-routing-backfill-pfad` §6: die Zeile „vom Wortlaut gedeckte
    Abweichung“ entfällt; V-2 ist damit geschlossen.

## Fitness Function (falls maschinell prüfbar)

Alle Zeilen sind **erwartet**: in dieser ADR wurde nichts gefahren und keine
Mutation erprobt; „der Implementer fährt sie“ ist eine Erwartung, keine
Erprobung. Die Aussage „alle drei Stände“ ist am Quelltext gelesen
(`routingRules`, `excludedColumns`, `transformationRules` und ihre Tests), nicht
gefahren.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (erwartet) | `TestExecuteRoutingReadFailureEndsRun`: Lesefehler (`outbound.ErrStorage`) an jeder Lesestelle des Routing-Stands endet `failed` mit Präfix `storage: ` und Ursachentext, ohne Commit | `make test` |
| Go-Test (bestehend, nicht neu gefahren) | `TestExecuteExclusionReadFailure`, `TestExecuteRuleReadFailure`: unverändert `storage: ` | `make test` |
| Go-Test (erwartet) | ein Wechsel des Routing-Standes zwischen zwei Lesungen endet weiter `configuration` (`TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration`) | `make test` |

## Re-Evaluierungs-Trigger

Ein Run-Fehlerbild zeigt, dass Betreiber einen Lesefehler als Konfigurationsproblem
brauchen (zum Beispiel ein Lesefehler, der nur durch fehlende Konfiguration
entstehen kann), oder `SPEC-008` erhält eine Klasse für „Zustand nicht
prüfbar“.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-01 | Accepted (Kurz-ADR auf Auftrag des Hauptlaufs; Schärft ADR-0139) | `review-slice-routing-backfill-pfad-fixrunde-1` F-N1, `verifikation-slice-routing-backfill-pfad` V-2 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen entstehen als neue ADR mit `Supersedes ADR-0141`.
