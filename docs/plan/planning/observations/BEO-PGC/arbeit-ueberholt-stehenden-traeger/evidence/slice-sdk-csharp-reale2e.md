# Beleg: slice-sdk-csharp-reale2e

Vorgang: `slice-sdk-csharp-reale2e` liefert die vier Realserver-Phasen
des C#-SDK-Packages `PgChangeFeed.Client` (gRPC, SSE, NATS-Vollinhalt,
HTTP) samt Abdeckungs-Träger `docs/user/sdk-e2e-abdeckung.md` und
`.d-check.yml`-`trace.coverage`-Eintrag (Label `SDK-E2E`). Der
§3.13-Suchlauf des Plans (committetes Feld) trug vier Träger-Zeilen —
alle vier einzeln bestätigt (Verifikation §3) — und war an den
**deklarierten Mustersatz** gebunden: `grep` nach den geplanten Trägern
(README-Werkzeug-Zeile, `sdks/csharp/README.md` §Status,
`docs/user/benutzerhandbuch.md`, `spec/pflichtenheft.md`).

Fund: der Suchlauf verpasste drei Fundstellen, deren Übereinstimmung mit
dem Suchraum nicht am Mustersatz hängt, sondern an der **Schreib-
Verbreiterung der Bewegung** — Haupt-Review F-1 (MEDIUM, vom Reviewer
gefunden), ein eigener, breiterer Lauf (`grep` nach
`trace.coverage`/„abdeckung", beide Stände `fce7af10`/HEAD):

- (a) `harness/README.md` — die `make doc-trace`-Zeile zählt die
  `trace.coverage`-Coverage-Dimensionen wörtlich auf („liest zusätzlich
  `e2e-abdeckung.md`, `bench-abdeckung.md` und `ci-matrix-abdeckung.md`");
  der neue `.d-check.yml`-Eintrag überholt genau diese Aufzählung — der
  Träger war vor dem Diff exakt richtig und ist durch die eigene Arbeit
  falsch geworden (dieselbe Fund-Klasse wie `slice-091`: die stehende
  Aufzählung, die die Arbeit falsch macht, ohne sie anzufassen).
- (b) `harness/sensors/docs-check.md` §Grenze — nennt nur
  `docs/user/e2e-abdeckung.md` als kuratierte Coverage-Dimension
  (vorgelagert bereits für Bench/CI-Matrix überholt; der neue Eintrag
  verbreitert die Lücke — dieselbe Grep-Form findet sie).
- (c) `welle-sdk-reale2e.md` §6 — behauptet über den hier erzeugten
  Träger, er „trägt `Datei:Zeile`-Orte (Muster `e2e-abdeckung.md`)";
  der reale Träger trägt Testklassen-/Runner-Verweis-Orte **ohne**
  Zeilenanker — der Mustersatz beschrieb die erwartete Form, nicht die
  gebaute.

Fixrunde `86892bdb` zog alle drei (die README-Zeile mit Neu-Messung der
RTM-Zahlen — gemessen 2026-09-23: 79 Anforderungen, 2 Waisen —, die
Sensor-Doku auf vier Dimensionen, die Welle-Form am Ort) und ergänzte
das Suchlauf-Feld um die drei Zeilen (Verifikation §3 bestätigt alle
sieben Feld-Zeilen gegen beide Stände).

Einordnung: dieselbe Lücken-Struktur wie bei
`slice-sdk-python-nats-stream-client-flaeche` („der Suchlauf folgt
§3.13s Wortlaut, nicht dem Plan-Mustersatz") — diesmal mit einer
schärferen Wurzel: der Suchraum eines Suchlaufs folgt der
**Schreib-Verbreiterung der Bewegung**, nicht dem Mustersatz. Jede
Verdrahtungs-Änderung, die eine aufgezählte Menge verbreitert (hier:
jeder `.d-check.yml`-`trace.coverage`-Eintrag), überholt alle
Prosa-Zeilen, die dieselbe Aufzählung tragen — README-Werkzeug-Zeile,
Sensor-Doku, Welle-Plan —, auch wenn sie keinem geplanten Träger
zugeordnet waren. Geschärfte Lehre (Closure-Notiz §7 des Plans): der
Suchraum des §3.13-Suchlaufs folgt der Schreib-Verbreiterung der
Bewegung, nicht dem deklarierten Mustersatz. Ausgang bleibt verkörpert
(`AGENTS.md` §3.13), kein neuer Schwellen-Übertritt — die Schärfung ist
eine Anwendungs-Schärfung der verkörperten Regel; ob sie einen eigenen
Satz im Träger trägt, prüft der Lese-Schritt der Welle-Closure.

Quelle: Review-Report `slice-sdk-csharp-reale2e` (F-1, MEDIUM),
Verifikations-Report desselben Slices (§1, §3), Commits `ecff7370`,
`11e42f7b`, `86892bdb`, `7bede9f2`.