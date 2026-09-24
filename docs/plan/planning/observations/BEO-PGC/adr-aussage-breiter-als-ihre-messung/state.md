Zustand: offen — Ausgang: **weiter offen** (Schwelle 3× erreicht; die Schärfungsfrage gehört dem
Lese-Schritt der Closure von `welle-backfill-bestand`).

**Adresse der Lücke:** die Entscheidung über eine Schärfung von `ADR-0114`
Entscheidung 3 liegt beim Architect (Notiz an ihn: Verifikation
`verifikation-slice-backfill-change-origin` V-3; §Entscheidung ist nach
`AGENTS.md` §3.5 und `ADR-0073` unberührbar, eine Änderung des Wortlauts wäre
eine neue ADR mit `Supersedes`). Diese Beobachtung ordnet keinen ADR-Auftrag an.

**Bis zur Entscheidung trägt die Grenze**, an drei Stellen (gemessen mit
`git grep -n -i -E 'Eigene Rechte|ACL-Verlust|Rechteliste' HEAD -- docs/user harness docs/plan/planning`):
Handbuch `docs/user/benutzerhandbuch.md` §4 „Schema aktualisieren" (Aufzählungspunkt
„Eigene Rechte"), Target-Dokument `harness/targets/schema-rollout.md` §Grenze
Punkt 3 („ACL-Verlust") und der Plan von `slice-backfill-change-origin`
(§3, Zeile zu `ADR-0114` F-5). Die Aussage der ADR bleibt für jede Rolle wahr,
die `tools/schema/nacharbeit-roles.sql` vergibt; die Lücke betrifft Rechte, die
ein Betreiber außerhalb des Repos setzt.

**Zweiter Beleg (`slice-backfill-snapshot-reader`):** die Zusage „byte-gleich"
in `ADR-0111` Teilfrage 1 und 2 (`col::text`) war an acht Typen belegt; gemessen
trifft der Cast 79 von 86 Typ-Spalten, der Architect schloss die Lücke mit
`ADR-0115` (`Supersedes` für drei Stellen), der Typ-Satz steht als Daten mit
Strukturtest im Paritätstest. Ausgangs-Kandidat bei 3× (Architect-Entscheidung,
nicht getroffen): eine ADR, die eine Aussage über **alle Werte einer Menge**
(„byte-gleich", „für jeden Typ") trifft, nennt die Menge, an der sie geprüft ist.

**Dritter Beleg (`slice-backfill-run-store`):** `ADR-0113` (`Accepted`) Festlegung 1
setzt die Annahme als `cdc_admin` voraus und nennt die bestehende Schleife „unverändert
lauffähig"; ihr Beleg deckt die Rechte auf `cdc.backfill_run`, das Recht auf
`cdc.administration_request` war nie gesetzt (`ADR-0047` führt die Antrags-Queue in ihrer
Rollen- und Aufrufer-Tabelle nicht). **Schwelle 3× erreicht mit diesem Slice**; der Ausgang
gehört dem Lese-Schritt der Closure von `welle-backfill-bestand`, nicht entschieden.
**Adresse der Lücke:** Architect (Plan `slice-backfill-run-store` §6: ob `ADR-0047` oder
`ADR-0113` eine Ergänzung braucht — §Entscheidung ist unberührbar, eine Änderung wäre eine
neue ADR mit `Supersedes` — und ob die Grants-Tabelle von `SPEC-029` den Vermerk auf
`cdc.administration_request` nennt). Der Betrieb ist über den Grant in
`tools/schema/nacharbeit-roles.sql` getragen. Verwandt, nicht doppelt gezählt:
`BEO-PGC/architect-verdikt-rollen-scope-luecke` (1×).

**Vierter Beleg (`slice-backfill-e2e`):** zwei Aussagen von `ADR-0118` (`Accepted`)
tragen eine andere Reichweite als die Messung — die Sperr-Warteschlange nennt nur
Leser („hergeleitet, nicht gemessen“), gemessen (PostgreSQL 18) warten auch
Schreiber und die Publication-Abfrage der Administration; der Satzteil
„E2E-belegt“ zu `RENAME COLUMN` hat im Repo keinen Beleg (nur `DROP COLUMN`).
Beide sind mit der neuen `ADR-0119` (`Supersedes ADR-0118`, teilweise) auf die
gemessene Reichweite gesetzt; das Handbuch trug die Messung bereits. Dieser Beleg
ist der Ausgang der zwei Aussagen, nicht der drei früheren (deren Adresse bleibt
der Architect).

**Stand 4×.** Der Lese-Schritt der Closure von
`welle-backfill-bestand` liest den Eintrag mit; die Frage einer Regel („eine ADR,
die eine Aussage über eine Menge trifft, nennt die Menge, an der sie geprüft ist“)
bleibt dort, nicht entschieden.

Gelesen wird der Eintrag im Sichtungs-Schritt der Slice-Planung. Zähler
(abgeleitet): **4×** (evidence/slice-backfill-change-origin.md,
evidence/slice-backfill-snapshot-reader.md,
evidence/slice-backfill-run-store.md,
evidence/slice-backfill-e2e.md).
