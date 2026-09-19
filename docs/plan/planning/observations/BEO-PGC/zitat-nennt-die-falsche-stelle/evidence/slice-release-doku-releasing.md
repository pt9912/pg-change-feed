# Beleg: slice-release-doku-releasing

Vom Reviewer gefunden (Review-Report zu `release-doku-releasing`, F-1)
und durch eigenes Aufschlagen der zitierten Stelle bestätigt: `docs/user/releasing.md`
behauptete, der Abgleich von `docs/user/version.md` gegen den Git-Tag
sei in „§3" („Einen Release auslösen") beschrieben — §3 deckt
ausschließlich die SemVer-Validierung des Tags, keine Erwähnung von
`version.md`. Der Abgleich steht tatsächlich in §4 („Was beim Release
automatisch passiert"), Punkt 2. Korrigiert (`§3` → `§4`).

Dritter Fundort dieser Klasse in drei aufeinanderfolgenden Slices
derselben Welle (`release-hub-description` F-1/F-2, jetzt
`release-doku-releasing` F-1) — jeweils ein interner Selbst-Verweis
innerhalb eines neu geschriebenen Dokuments, nicht ein Verweis auf ein
fremdes Dokument wie bei den früheren Belegen. Siebter Beleg der
bereits verkörperten Beobachtung — kein neuer Zähler-Trigger.
