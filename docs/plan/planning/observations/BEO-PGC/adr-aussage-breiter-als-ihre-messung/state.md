Zustand: offen — Ausgang: **weiter offen** (unter der 3×-Schärfungsschwelle).

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

Gelesen wird der Eintrag im Sichtungs-Schritt der Slice-Planung. Zähler
(abgeleitet): 1× (evidence/slice-backfill-change-origin.md).
