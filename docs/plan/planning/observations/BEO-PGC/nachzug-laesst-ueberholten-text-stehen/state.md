Stand: **offen** (8×, Schwelle erreicht mit `slice-backfill-spec-nachzug` —
Ausgang noch **nicht** zugewiesen). Gelesen wird der Eintrag im Lese-Schritt
der Closure von `welle-backfill-bestand` (Modul 6): die Regelschärfungs-Frage
— ob und wie `.harness/skills/reviewer.md` oder ein Schritt der
Implementer-Selbstprüfung die Klasse fängt — ist eine Architect-Entscheidung
(Modul 4/8), die dieser Lese-Schritt als Steering-Loop-Eintrag weiterträgt.

Zähler (abgeleitet): 8× (evidence/slice-sdk-python-projektgeruest.md,
evidence/slice-sdk-kotlin-pack-werkzeug.md,
evidence/slice-backfill-spec-nachzug.md,
evidence/slice-backfill-change-origin.md,
evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-bench-richtgroesse.md,
evidence/slice-backfill-slot-leerlauf-bestaetigung.md,
evidence/slice-backfill-sdk-origin.md). Der achte Beleg (`slice-backfill-sdk-origin`,
Review F-6) trifft einen Träger **außerhalb des Diffs**: ein Handbuch-Absatz mit einer
Zukunftsaussage über ein Package, die schon vor dem Slice überholt war; der Nachzug lief
als Planner-Zug in der Fixrunde. Der siebte Beleg
(`slice-backfill-slot-leerlauf-bestaetigung`, Verifikation V-2) trifft einen Träger in
einer **fremden Datei** (Register-Sichtung der Welle-Datei): der Implementer fand ihn im
Suchlauf und meldete ihn an den Planner, statt ihn still mitzuändern; der überholte Text
stand bis zum Nachzug der Closure im Baum. Der Ausgang bleibt beim Lese-Schritt der
Welle-Closure; die Meldung der fremden Datei ist die Regel, ihre Frist die offene Frage.
Der sechste Beleg
(`slice-backfill-bench-richtgroesse`, Review F-1) trifft **Plan und Handbuch nach
einem Architect-Verdikt**: das Verdikt beantwortete eine im Plan offen geführte Frage
und wies dem Planner den Nachzug zu; zum Stand des Reviews trugen beide Träger den
Text davor (die Fixrunde zog nach, der Verifier bestätigte) — derselbe Mechanismus
mit einem **Verdikt** statt eines Plan-Nachzugs als auslösendem Artefakt. Der fünfte Beleg
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
