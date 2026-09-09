Zustand: offen — Ausgang wird vom Lese-Schritt (die Welle-1-Closure)
zugewiesen; der Zähler ist die Zahl der `evidence/`-Dateien (2× nach
slice-002) und wird nie gespeichert.

Beobachtung: die Layer-Globs der `.a-check.yml` matchen null Dateien
(`internal/**` entsteht erst mit slice-002/003); `make a-check` ist
seit slice-001 grün, aber nur über `composition_root: cmd/**`
nichtleer (Verifier: verify-slice-001.md, Implementer-Risiko (b)).

Stand nach slice-002 (evidence/slice-002.md): `domain`/`ports` matchen
echten Content; `app`/`adapters` weiterhin null — teils entkräftet.