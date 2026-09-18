# Beleg: slice-090

Vorgang: `slice-090` — das Sync-Gate des generierten Protobuf-Codes.

Fund: Der Slice fügt `make gates` ein neues Ziel hinzu (`generated-sync`), und
`ci.yml` fährt `make gates` auf einem **frischen Runner** — also ohne den warmen
Layer-Cache, auf dem das Gate lokal läuft. Gemessen lokal: der erste Lauf baut
die Dockerfile-Stufe `proto` real mit Netz (`apk add` und zwei `go install`; die
gedruckte Build-Zeile `#9 DONE 8.5s`), ab dem zweiten Lauf greift der Cache
(`real 0m0,941s` / `real 0m0,928s`). Ob der kalte Build auf dem Runner grün
durchläuft und innerhalb des Zeitlimits bleibt, ist **lokal nicht prüfbar**
(`AGENTS.md` §3.1) — der Beleg ist der erste reale Post-Push-Lauf.

**§3.10 ist dem Buchstaben nach nicht ausgelöst** — dieser Vorgang ändert am
Workflow nur einen Schrittnamen und zwei Kommentarzeilen, keine Stufe, keine
Matrix, keine Abhängigkeit, kein `uses:`. Sein **Grund** trifft aber zu: die
Wirksamkeit einer Änderung an dem, was der Workflow fährt, ist statisch nicht
beweisbar. Der Unterschied ist festgehalten, nicht eingeebnet: die Regel gilt
hier nicht, ihr Gegenstand schon.

**Warum das zählt:** Die Zusage des Slice ist, dass das Gate **im Gate-Lauf**
läuft. Lokal ist das belegt; auf dem Runner ist es das nicht. Wer den Slice
danach liest, soll die Grenze sehen und nicht annehmen, `make gates` sei überall
dasselbe.

Quelle: Verifikationsbericht zu `slice-090` (V-5) ·
Review zu `slice-090`, Delta-Review (Negativbefunde, `#6`–`#9 CACHED`) ·
`docs/plan/planning/done/slice-090-sync-gate-protobuf.md` §6 (viertes
Risiko) · `.github/workflows/ci.yml`.
