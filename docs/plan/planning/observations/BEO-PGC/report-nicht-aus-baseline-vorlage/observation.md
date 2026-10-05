# BEO-PGC/report-nicht-aus-baseline-vorlage

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verfassen von
Review- und Verifikations-Reports durch die Rolle, die sie ablegt; keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Review- und Verifikations-Reports unter `docs/reviews/` wurden in einer
Sitzung **frei** geschrieben statt aus der Baseline-Vorlage
für Review-Reports unter `.harness/baseline/<tag>/templates/` kopiert und ausgefüllt. Anlass auf Planner-Seite: der Auftrag untersagte dem Reviewer das Schreiben der
Report-Datei, der Planner legte sie an und nahm dafür die Vorlage nicht. Sechs Reports
mussten später aus der Vorlage neu gesetzt werden (Commit `9b360010`); die Repo-Kopie der
Vorlage wurde im Zug entfernt, damit das Gerüst nur noch aus der Baseline kommt (Commit
`f5b2840a`).

**Warum das zählt:** Ein frei geschriebener Report verliert die Pflichtfelder, die die Vorlage
trägt (Negativbefund-Pflicht, Findings-Schema); kein Sensor liest die Form eines Reports, die
Abweichung zeigt sich erst beim Gegenlesen.

**Abgrenzung.** `BEO-PGC/plan-vorlagen-defekt` und `BEO-PGC/vorlagenrest-in-closure-notiz`
betreffen **gefüllte** Vorlagen mit Resten; hier wurde die Vorlage nicht genommen.
`BEO-PGC/report-nackte-id-ohne-link` betrifft den Inhalt eines Reports, nicht seine Herkunft.

Herkunft: Rückfrage des Auftraggebers; Ursprung der Angabe **übernommen**, die Commits
`9b360010` und `f5b2840a` sind am Repo gemessen.
