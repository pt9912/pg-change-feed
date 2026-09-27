**Vorgang:** adr-0129-capture-quellseite-keepalive-test-an-beiden-pins (Architect-Zug und Review der ADR
[`ADR-0129`](../../../../../adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md), 2026-09-27)

**Fund:** zwei Stellen, keine mit Wirkung auf eine Repo-Datei.

- **Architect, Grenze eingetreten.** Ein Host-`python3 --version` ohne Repo-Pfad im Befehlsstring passierte den
  Guard: der in `MR-003` benannte Rand (Host-`python`/`perl` auf einem Pfad ohne Repo-Namen geht durch). Er hatte
  keine Wirkung auf eine Repo-Datei und verstieß gegen `AGENTS.md` §3.1 (Host-Interpreter). Angabe des
  Auftraggebers, **übernommen**; nicht am Guard gemessen. Nachgemessen am 2026-09-27 (Stand `cea198fb`): die
  Hook-Eingabe `python3 --version` am Guard ergibt Exit 0 ohne Ausgabe.
- **Reviewer, Guard-Treffer.** Ein `python3 --version` im selben Vorgang wurde geblockt, ohne Wirkung auf eine
  Repo-Datei. Angabe des Auftraggebers, **übernommen**. Der Befehlsstring ist im Auftrag nicht wiedergegeben; nach
  der Regel von `MR-003` ist ein Block bei `python3 --version` nur mit einem Repo-Pfad im Befehlsstring erklärbar
  (**hergeleitet**, nicht gemessen).

**Form (Ausprägung):** **Regelgrenze**, kein neuer Fehlgriff der Klasse: dieselbe Handlung (ein Host-`python3` zur
Versionsabfrage) erreicht je nach Befehlsstring den Pass (kein Repo-Name) oder den Block (Repo-Name) — die zwei Seiten
des Randes, den die Nutzer-Entscheidung „Weg 3“ (`state.md`) mit der unbedingten Sperre von `python`/`python3` am Kopf
schließt. Ein Beleg für das erste Neubewertungs-Kriterium der Kopf-Liste (Host-Interpreter-Aufruf ohne Repo-Pfad **und**
Wirkung auf eine Repo-Datei) ist keine der beiden Stellen.

Quelle: Angabe des Auftraggebers im Auftrag an den Planner zu `slice-harness-guard-blocked-python` (2026-09-27),
übernommen; kein Report als Anker gelesen.
