# BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Zusammenarbeit
der Rollen-Subagenten mit dem Schreib-Werkzeug, keine eigene Sub-Area im Sinn
der Modus-Deklaration).

Die Beobachtung: Ein Reviewer-Subagent versuchte, einen Entwurf seines
Review-Reports in das Scratchpad-Verzeichnis zu schreiben; das `Write`-Werkzeug
lehnte den Aufruf mit der Meldung „Subagents should return findings as text,
not write report files" ab. Der Subagent meldete daraufhin, der Zielpfad
`docs/reviews/…` sei gesperrt. Der Zielpfad war nicht gesperrt: ein zweiter
`Write`-Versuch, diesmal am Zielpfad, ging durch, der Report liegt committet
vor. Die Ablehnung galt dem einen Aufruf (Report-artige Datei am
Scratchpad-Entwurfspfad); der Subagent übertrug sie auf einen Pfad, den dieser
Aufruf nicht betraf.

**Warum das zählt:** Die Rollen-Kette trägt ihr Übergabe-Artefakt in einer Datei
unter `docs/reviews/`. Eine falsche Sperr-Meldung lässt den Auftraggeber das
Artefakt an anderer Stelle oder gar nicht erwarten und kostet eine Runde; die
Meldung war nicht der Werkzeugbefund, sondern dessen Verallgemeinerung. Kein
Gate liest, wie ein Subagent eine Werkzeug-Ablehnung berichtet — der Wächter ist
die Genauigkeit der Meldung: welcher Aufruf, welcher Pfad, welcher
Werkzeug-Wortlaut.

Ursache (im Transkript des Laufs vom Auftraggeber geprüft, hier übernommen):
Die Werkzeug-Meldung bezog sich auf den konkreten Aufruf; der Subagent hat den
Zielpfad nicht getrennt geprüft, bevor er ihn als gesperrt meldete. Eine
Schuldzuweisung folgt daraus nicht; die Beobachtung hält den Befund und den
Abstand zwischen Werkzeugmeldung und Berichtsaussage fest.

Deklaration: `slice-backfill-row-image-gemeinsam`, real aufgetreten beim Review
dieses Slice (Beleg: Commit `0d5c3f1a`).
