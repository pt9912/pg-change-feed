Zustand: **verkörpert** — Ausgang zugewiesen beim Lese-Schritt der
`welle-d-check`-Closure (der Architect-Verdikt zum welle-d-check-Lese-Schritt,
§2): Kein
neuer HIGH-Punkt — der bestehende Punkt „Beleg trägt seinen Satz nicht" in
`.harness/skills/reviewer.md` ist um die Verweis-Form („einen Verweis auf
eine Stelle eines anderen Dokuments — eine Abschnittsnummer, eine
ADR-Festlegung, eine Slice-/Welle-Kennung") und eine explizite
Gegenprobe-Pflicht am Original erweitert (statt einer Zusammenfassung/
Berichts-Kopfzeile), Herkunfts-Anker `seit welle-d-check`. `AGENTS.md` §3.12
**Instanz B** trägt das Prinzip weiterhin auf der Schreiber-Seite; die
Erweiterung schließt die Leser-seitige Lücke beim Reviewer.

Ein vierter Beleg (`welle-d-check-verkoerperung`) traf **denselben Commit**,
der diese Verkörperung schrieb — bestätigt den Wert der Regel unmittelbar
(siehe evidence-Datei).

Zähler (Datei-Anzahl unter `evidence/`, real ausgezählt): **9×**
(evidence/slice-antragsqueue-lesefehler-failed.md — der Verweis liegt im eigenen Plan: eine Zeile
des nummerierten Suchlauf-Feldes (Review F-2, HIGH) und die Nummer eines Findings des zugehörigen
Reviews (Verifikation V-1, LOW); beide vor dem Merge gefunden, Ausgang unverändert **verkörpert**;
evidence/slice-backfill-run-store.md — Kopplungs-Kommentar nennt eine Testdatei, die keinen
Schema-Neuaufbau ausführt, Review F-3; evidence/slice-090.md, evidence/slice-102.md,
evidence/slice-d-check-tracked-modul.md,
evidence/welle-d-check-verkoerperung.md,
evidence/slice-release-hub-description.md — zwei Fundstellen, F-1/F-2,
evidence/slice-release-doku-releasing.md,
evidence/slice-sdk-python-nats-stream-client-flaeche.md)
— **verkörpert**, weitere Belege zählen weiter, ohne die Verkörperung
erneut auszulösen. Der Beleg
`evidence/slice-sdk-python-nats-stream-client-flaeche.md`
(`slice-sdk-python-nats-stream-client-flaeche`, F-3 HIGH): der neue
`__init__.py`-Docstring zitierte `ADR-0110` Festlegung 3 als
entscheidende Schicht für „v2 desselben Packages" — die Festlegung
entscheidet ausdrücklich nichts (Delegation an den Folge-Zug); die
tatsächliche Entscheidungsschicht ist der Welle-Plan §6. Der Verweis
wurde übernommen, ohne das Original aufzuschlagen; Fixrunde 1 re-anchorte
auf die Entscheidungslage, das Fixrunden-Re-Review maß beide Anker im
Original nach.
Das
Erstauftreten fiel in `slice-090` beim Übergang **in einen stehenden Träger**
auf: aus der Kopfzeile eines Review-Reports wurde eine Abschnittsnummer
übernommen und in `harness/sensors/generated-sync.md` gesetzt. Der zweite
Beleg (`slice-102`) trifft eine andere Form derselben Klasse: kein Zahlen-,
sondern ein Slice-Kennungs-Zitat — der Plan verwies auf „`slice-095`s §1" für
eine Aussage, die tatsächlich in `slice-097` §1 steht (Reviewer F-3, vom
Verifier unabhängig reproduziert). Der dritte Beleg
(`slice-d-check-tracked-modul`) trifft eine dritte Form: kein Zahlen- oder
Kennungs-Zitat, sondern eine **Entscheidungslage** (zwei ADRs als
Präzedenzmuster, die die gegenteilige Aussage tragen) — der Reviewer hat dort
selbst den korrekten Präzedenzfall recherchiert (F-4), das Architect-Verdikt
hat ihn eigenständig bestätigt. Der Lese-Schritt (Ausgang zuweisen) gehört der
Closure von `welle-d-check` (Modul 6/8), da dieser dritte Beleg-Vorgang einer
offenen Welle angehört, nicht dem wellenlosen Pfad.

**Nicht zu verwechseln** mit `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(dort driftet ein **Wert** gegen die Messung) und mit
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (dort trägt ein **Befehl**
seinen Satz nicht). Hier ist der Gegenstand der **Verweis**: er wird nicht
nachgemessen, sondern aufgeschlagen.
