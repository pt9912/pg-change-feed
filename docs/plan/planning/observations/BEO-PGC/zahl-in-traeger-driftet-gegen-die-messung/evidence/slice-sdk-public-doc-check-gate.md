# Beleg: slice-sdk-public-doc-check-gate

Vorgang: `slice-sdk-public-doc-check-gate` — Erst-Review F-1 (HIGH,
`review-slice-sdk-public-doc-check-gate.md`), Klasse „Zahl im Träger
driftet gegen die Messung"; nicht identisch mit dem Vorgang
`welle-sdk-grpc-administration-flaeche` (dessen Beleg trägt die zwei
Findings des Closure-Note-Reviews an der Welle-Results-Notiz).

Fund: Das §3.13-Suchlauf-Feld des Slice-Plans behauptete für den
Arbeitsbaum nach dem Zug 130 Treffer (als 112 + 18 zerlegt) bei Parent
`f287c81b`. Nachgemessen: der genannte Stand `72293784` misst 132 — die
+18-Breakdown zählte nur die eigenen vier Dateien und überging die +2
Treffer des im selben Range liegenden Commits `865c273e`
(Closure-Note-Report, 2 Treffer); `f287c81b` war zudem nicht der Parent
der Zug-Commits. Die 130 war damit aus 112 + 18 abgeleitet und als
Messung des genannten Stands ausgewiesen — dieselbe Form wie beim
vierundzwanzigsten Beleg (Zahl stand schon beim Commit falsch, nicht erst
durch spätere Arbeit), hier mit zwei zusätzlichen Bindungs-Fehlern
(falscher Parent, übersehener Fremd-Commit im Range).

Gezogen: Fixrunde `4ca5af64` — beide Zahlen tragen ihren Stand als
Commit-Kennung (Parent `865c273e` 114, Stand `72293784` 132, je Datei
verifiziert) statt als „Arbeitsbaum"-Ableitung; das Fixrunden-Review
maß beide Stände byte-gleich nach (114/132, Delta +18 je Datei
aufgelöst). Kein Rest-Vorkommen der alten Zahlen im DoD-Feld.

Einordnung: Ausgang bleibt **verkörpert** (`AGENTS.md` §3.12 Instanz A,
Reviewer-Skill HIGH-Punkt „Zahl im Träger ohne Ursprung — oder gegen die
Messung driftend"); HIGH → Datei unabhängig vom Deckel. Anwendungs-Schärfung
für das Suchlauf-Feld: liegt im Range der Messung ein fremder Commit,
trägt jede Zahl ihren Stand als Commit-Kennung und die Zerlegung jede
Datei — abgeleitet wird nie als gemessen ausgegeben.
