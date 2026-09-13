**Vorgang:** slice-057
**Fund:** Der Verifier führte `make test-integration` dreimal aus; der
erste Lauf brach real mit Exit-Code 2 in der Retention-Lebenszyklus-
Timing-Zusicherung (`LH-FA-RET-004`) ab, die beiden folgenden Läufe
(identischer Code-Stand) liefen sauber durch. Diff-Inspektion schließt
eine Ursache in diesem Slice (reines Renaming) aus — erstes benanntes
Vorkommen dieser Flake-Klasse in diesem Repo.
