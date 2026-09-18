Zustand: **verkörpert** — Ausgang: **verkörpert** → neue Hard Rule
`AGENTS.md` §3.10 („Ein neuer oder strukturell geänderter
GitHub-Actions-Workflow gilt erst nach einem realen, grünen
Post-Push-Lauf als abgeschlossen“) · seit welle-17. Der Architect-Verdikt
zu github-actions-unverifizierbar-lokal (3×)
(Diagnose: anders als beim Nackte-Kennung-Präzedenzfall ist hier kein
Sensor möglich — GitHub Actions ist ein externer, gehosteter Dienst, den
ein Docker-only/netzloser Sensor strukturell nicht erreichen kann —, aber
auch keine reine Interpretationsfrage wie beim Chronik-Fall; die drei
Belege lösten sich real jeweils positiv auf, aber jedes Mal durch neu
hergeleitete, nirgends kodifizierte Disziplin statt einer bereits
geltenden Regel. Deshalb weder `gestrichen` — kein bestehender
Mechanismus schließt den Risikopfad strukturell — noch `geplant` — keine
ausstehende Bau-Entscheidung, nur eine sofort schreibbare Prosa-Regel).
Lese-Schritt lief als Teil der laufenden `welle-17`-Closure (Modul 8
§Rollen-Sequenz für eine Welle, Schritt 3b, Planner → Architect →
Planner-Zug).

Zähler (abgeleitet): **7×** (evidence/slice-039.md, evidence/slice-056.md,
evidence/slice-064.md, evidence/slice-082.md, evidence/slice-090.md,
evidence/slice-098.md, evidence/slice-099.md) — Schwelle erreicht, Ausgang
zugewiesen. Der fünfte Beleg fällt **außerhalb** der Regel an, die sie
verkörpert: `slice-090` ändert am Workflow nur einen Schrittnamen und
Kommentarzeilen, `AGENTS.md` §3.10 ist dem Buchstaben nach **nicht**
ausgelöst — sein Grund aber trifft zu, weil `ci.yml` mit `make gates` jetzt
einen Schritt fährt, dessen kalter Layer-Cache auf dem Runner lokal nicht
prüfbar ist. Der Beleg hält den Unterschied fest, statt die Regel zu dehnen.
Der sechste (`slice-098`) und siebte (`slice-099`) Beleg sind wieder der
Regelfall: je ein strukturell neuer bzw. erweiterter Workflow
(`.github/workflows/examples.yml`, C#- bzw. Kotlin-Job), `AGENTS.md` §3.10
dem Buchstaben nach ausgelöst, korrekt an allen drei Trägern als offen
geführt — kein neuer Schwellen-Übertritt, die Regel steht bereits.

(Hinweis, vom Verifier bei `slice-065` gefunden,
Verifikationsbericht zu `slice-065`: Die vorige Kopfzeile dieser Datei
trug „weiter offen“ als Ausgangsbezeichnung — das ist die
Ausgangsbezeichnung für ein Slice-§6-Risiko (Modul 5), nicht einer der
drei gültigen Register-Ausgänge ab 3× (verkörpert/geplant/gestrichen,
Modul 6). Mit diesem Zug korrigiert.)
