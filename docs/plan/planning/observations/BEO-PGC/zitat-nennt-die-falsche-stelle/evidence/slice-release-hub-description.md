# Beleg: slice-release-hub-description

Zwei Fundstellen im selben Slice, beide vom Reviewer gefunden und durch
eigenes Aufschlagen der zitierten Stelle bestätigt (Review-Report zu
`release-hub-description`, Findings F-1/F-2), vom Verifier unabhängig
gegengeprüft:

- **F-1:** `.github/workflows/hub-description.yml:7` zitierte
  „`ADR-0051` §Konsistenz-Pruefung" für den Satz „Präsentation, kein
  Bestandteil der
  Distributions-Zusage" — der Wortlaut steht tatsächlich in **Entscheidung
  8**, nicht im Abschnitt „Konsistenz-Prüfung" (der behandelt `AGENTS.md`
  §3.1/§3.6). Ein `grep` gegen die volle ADR-Datei bestätigt: die
  Formulierung kommt genau einmal vor, exakt in Entscheidung 8.
- **F-2:** Die DoD-Zeile im Slice-Plan verwies auf „§1 Abgrenzung der
  Welle-Datei" für die Aussage über `DOCKERHUB_USERNAME`/
  `DOCKERHUB_TOKEN` — die Welle-Datei hat kein §1 „Abgrenzung"; die
  zutreffende Stelle ist §6 „Out-of-Scope für diese Welle". Der Text
  stammte ursprünglich vom Planner (unverändert vor diesem Commit), wurde
  aber im Implementer-Commit per DoD-Checkbox bestätigt.

Beide Aussagen selbst waren sachlich zutreffend — nur die zitierten Anker
falsch, unabhängig von der inhaltlichen Richtigkeit (Kernmerkmal dieser
Beobachtung: ein Beleg, der die Aussage nicht trägt, ist ein Befund, auch
wenn die Aussage stimmt). Beide in derselben Fixrunde korrigiert
(`5cc7fa49`), vom Verifier real gegen die Originalstellen nachgeschlagen
und bestätigt.

Fünfter und sechster Beleg der bereits verkörperten Beobachtung — kein
neuer Zähler-Trigger (bereits `verkörpert` seit `welle-d-check`).
