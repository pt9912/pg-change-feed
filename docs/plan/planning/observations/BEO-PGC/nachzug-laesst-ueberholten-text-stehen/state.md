Stand: **offen** (5×, Schwelle erreicht mit `slice-backfill-spec-nachzug` —
Ausgang noch **nicht** zugewiesen). Gelesen wird der Eintrag im Lese-Schritt
der Closure von `welle-backfill-bestand` (Modul 6): die Regelschärfungs-Frage
— ob und wie `.harness/skills/reviewer.md` oder ein Schritt der
Implementer-Selbstprüfung die Klasse fängt — ist eine Architect-Entscheidung
(Modul 4/8), die dieser Lese-Schritt als Steering-Loop-Eintrag weiterträgt.

Zähler (abgeleitet): 5× (evidence/slice-sdk-python-projektgeruest.md,
evidence/slice-sdk-kotlin-pack-werkzeug.md,
evidence/slice-backfill-spec-nachzug.md,
evidence/slice-backfill-change-origin.md,
evidence/slice-backfill-sql-administration.md). Der fünfte Beleg
(`slice-backfill-sql-administration`, F-6, F-7, V-1) trifft **Zählwörter und Aufzählungen**
in Kommentaren, Schema-Beschreibungen und einem Skript, die der Suchlauf über den
Symbolnamen nicht fand („vier Antragsarten", „die vier Views", „sechs" Fremdobjekte,
„die beiden Tabellen-Antragsarten"); ein Verifier-Suchlauf nach den Zählwörtern
(`vier`, `beiden`) fand den Rest, den die Fixrunde übersah. Die Form „Suchlauf nach dem
Zählwort neben dem Symbolnamen" ist Kandidat für die Regelschärfung dieses Lese-Schritts
(Architect-Entscheidung, nicht getroffen). Der zweite Beleg trifft denselben
Mechanismus an einem anderen Träger-Typ (`spec/pflichtenheft.md` Fließtext
statt Slice-Plan-Dokument) und mit Ursprung und Vorkommen getrennt (Ursprung:
`slice-sdk-python-pack-werkzeug`; Vorkommen: gefunden beim Review von
`slice-sdk-kotlin-pack-werkzeug`). Der dritte Beleg trifft ihn an einer
Tabellenzeile des Pflichtenhefts (`SPEC-022`: die neue Zeile „Position und
`limit`" nimmt die Fortsetzungs-Regel der Zeile „Reihenfolge" desselben
Abschnitts zurück, ohne dass diese den Vorbehalt trägt) — Ursprung und
Vorkommen fallen dort im selben Zug zusammen; das Review fand ihn, die
Fixrunde behob ihn. Der vierte Beleg (`slice-backfill-change-origin`) trifft ihn
in einem Kommentarblock des `Makefile` und der Meldung, die er beschreibt (der
erste Absatz sagte „nur bekannte Fremdobjekte", der ergänzte Block ließ
`--allow-destructive` auch neben der View-Signatur zu) — ein vierter
Träger-Typ, dieselbe Zeit-Form: der Nachzug erzeugt den Widerspruch im selben
Zug. Der Ausgang bleibt zugewiesen beim Lese-Schritt der Closure von
`welle-backfill-bestand`.
