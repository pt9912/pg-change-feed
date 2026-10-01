# Beleg: slice-routing-betriebsdoku

Vorgang: `slice-routing-betriebsdoku` — Review F-1 (LOW) und Verifikation V-1 (LOW),
Klasse „Zahl im Träger ohne tragfähigen Ursprung".

Quelle: `docs/reviews/review-slice-routing-betriebsdoku.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-betriebsdoku.md` (V-1). <!-- d-check:status-provenance -->

Fund: Das Benutzerhandbuch führte die Kosten-Spanne der zweiten NATS-Veröffentlichung
(0,92 bis 1,30) mit einem Ursprung, den ein Betreiber nicht auflösen kann: „gemessen und
abgeleitet in diesem Handbuch-Zug", „übernommen aus den Berichten der Umsetzung" und
„(Review, Fixrunde, Verifikation)" ohne Link; die gedruckte Zeile des Handbuch-Zugs
(Verhältnis 1,02) stand in keinem committeten Report. Die Zahlen waren wahr (der
Nachmessungs-Lauf des Verifiers ergab 1,20, innerhalb der Spanne); beanstandet war die Form.
Gleiche Ursache im Review F-1 für die Aussagen „gemessen im Lauf von `make test-integration`"
ohne Lauf-Anker.

Gezogen: Das Handbuch nennt für die Spanne nur Werte, deren Bericht verlinkt und committet
ist (Review-Bericht und Verifikations-Bericht des Quell-Slice), und sagt, dass weitere Läufe
nicht einzeln genannt werden; die E2E-Aussagen tragen „übernommen" samt Verweis auf den
Verifikations-Report.

Einordnung: Ausgang bleibt **verkörpert** (`AGENTS.md` §3.12 Instanz A, Reviewer-Skill);
neue Form (der Ursprung ist vorhanden, der Leser kann ihn nicht auflösen), deshalb eine
Datei trotz Deckel. Geschärfte Anwendung: in einem Betreiber-Handbuch ist ein Ursprung nur so
viel wert, wie der Leser ihn auflösen kann; ein Prozess-Vokabular („in diesem Zug") ist kein
Anker. Verlinke den committeten Report oder nenne keine Einzelzahl.
