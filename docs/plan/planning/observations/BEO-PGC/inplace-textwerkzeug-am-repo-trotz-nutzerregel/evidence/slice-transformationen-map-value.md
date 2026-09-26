**Vorgang:** slice-transformationen-map-value (Review F-7, INFO; Verifikation V-5 und V-6; die
Closure-Sitzung des Planners) — der nächste Vorgang nach `slice-harness-guard-inplace-textwerkzeug`,
dessen Läufe unter dem PreToolUse-Guard liefen (`MR-003`, zweites Neubewertungs-Kriterium).

**Fund:** drei Stellen, keine mit Wirkung auf eine Repo-Datei außer der beabsichtigten.

- **Implementer, Anhängen per Umleitung.** Ein Lauf hängte Testtext per `cat >>` an eine
  Testdatei an (Angabe des Auftraggebers an den Review, **übernommen**; im Diff keine Spur,
  `make fmt-check` und `make kommentar-kennungen` ohne Befund). Der Weg ist weder `sed -i`
  noch `perl -pi` noch ein Host-Interpreter: der Guard liest Umleitungen nicht (Grenz-Zeile in
  `MR-003`), und `AGENTS.md` §3.1 trägt dafür kein ausdrückliches Urteil — der Verbotssatz nennt
  in-place Umschreiben, die Überschrift des Absatzes („Text-Umschreiben im Repo ist Sache der
  Datei-Werkzeuge des Laufs“) lässt eine strengere Lesung zu.
- **Verifier, Guard-Treffer.** Ein `sed -i` auf einer Scratch-Hilfsdatei wurde vom Guard
  geblockt (Klasse in-place), der Aufruf lief nicht, keine Wirkung auf eine Repo-Datei; die
  Änderung lief danach über das Write-Werkzeug (Angabe des Verifikations-Reports,
  **übernommen**).
- **Planner der Closure-Sitzung, Grenze eingetreten.** Ein Host-`python3 --version` im
  Befehlsstring `cd <Repo-Wurzel> && python3 --version` lief und druckte die Version; er hatte
  keine Wirkung auf eine Repo-Datei. Der Guard ließ ihn passen — der in `MR-003` benannte Rand
  „`cd <Repo> && python3 x.py` geht durch“. Nachgemessen in derselben Sitzung: die Hook-Eingabe
  mit diesem Befehlsstring am Guard ergibt Exit 0 ohne Ausgabe (Pass); ein
  `python3 --version` mit einem Repo-Pfad als Argument ergibt die Meldung der Klasse `interp`.
  Der Aufruf verstieß gegen `AGENTS.md` §3.1 (Host-Interpreter) und wurde in derselben Sitzung
  nicht wiederholt.

**Form (Ausprägung):** **Regelgrenze**, kein neuer Fehlgriff der Klasse: zwei der drei Stellen
liegen hinter der Grenz-Zeile des Guards (Umleitung, `cd` im selben Kommando) und sind vom Review
bzw. von der Selbstdisziplin der Rolle getragen; die dritte zeigt, dass der Guard greift, wo er
liest. Ein Beleg für das erste Neubewertungs-Kriterium der Kopf-Liste `tools/harness/blocked/go`
(Host-Interpreter-Aufruf ohne Repo-Pfad im Befehlsstring **und** Wirkung auf eine Repo-Datei) ist
keine der drei Stellen.

Quelle: `docs/reviews/review-slice-transformationen-map-value.md` (F-7) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-transformationen-map-value.md` (V-5, V-6). <!-- d-check:status-provenance -->
