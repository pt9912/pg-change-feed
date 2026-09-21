# BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Accepted-ADR-Immutabilitäts-Disziplin (`AGENTS.md` §3.5), keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Implementer weicht beim Schreiben von der wörtlichen
Festlegung einer bereits `Accepted`-ADR ab und stellt die Abweichung im
Slice-Plan als bewusste, begründete Entscheidung dar — mit einer eigenen
Lesart einer einschlägigen Hard Rule — statt sie als offene Frage an
Reviewer/Architect zu kennzeichnen. `make gates` fängt das nicht: Der
Widerspruch steht in Prosa/Workflow-Konfiguration gegen ADR-Text, keine
geprüfte Struktur. Gefunden hat es der Reviewer (HIGH), nicht ein Sensor.

Deklaration: `slice-sdk-kotlin-publish-workflow` (Review-Finding F-1,
`docs/reviews/review-slice-sdk-kotlin-publish-workflow.md` <!-- d-check:status-provenance -->) — `ADR-0109`
§Entscheidung Festlegung 5 legt wörtlich Docker-only bis einschließlich
`publish` fest; der ursprüngliche Implementer-Commit ließ `./gradlew
publish` stattdessen direkt auf dem GitHub-hosted Runner laufen und
begründete das im Slice-Plan mit einer eigenen Lesart von `AGENTS.md`
§3.1. Aufgelöst über eine Fixrunde (Code-Fix, `21872c3d`), **nicht** über
eine Supersede-ADR: Der Fund betraf den *Ort* der Ausführung, keine
inhaltliche ADR-Korrektur — `ADR-0109` §Re-Evaluierungs-Trigger 4 zieht
diese Unterscheidung bereits selbst („ein reiner Workflow-Bugfix ohne
Entscheidungsänderung braucht keine neue ADR").
