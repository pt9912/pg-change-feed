# Beleg: slice-sdk-kotlin-sse-client-flaeche

Vorgang: `slice-sdk-kotlin-sse-client-flaeche` (erster Slice der Welle
`welle-sdk-kotlin-vollabdeckung`).

Fund: Der Slice-Plan behauptete durchgängig (§3, §6, §7, §8) der Zähler von
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` stehe bei „2×,
unter der Schwelle". Das war zum Zeitpunkt des Slice-Starts (`next →
in-progress`, Commit `baf6e8bd`, 2026-09-22 06:06:49) bereits **falsch**:
Die parallele Geschwister-Welle `welle-sdk-csharp-vollabdeckung` hatte den
Zähler über `slice-sdk-csharp-sse-client-flaeche` (Commit `d8cc10a1`,
2026-09-22 00:10:21 — vor dem Start dieses Slice) bereits auf real **3×
(Schwelle erreicht)** gehoben. Der Welle-Plan
(`welle-sdk-kotlin-vollabdeckung.md`) trägt diese Korrektur bereits seit
Commit `3438f52b` (06:05:38 — eine Minute vor diesem Slice-Start); die
Kopie im Slice-Plan selbst wurde bei dessen eigener Niederschrift nicht
gegen den zu diesem Zeitpunkt bereits korrigierten Welle-Plan
nachgeprüft, sondern aus einem älteren Entwurfsstand übernommen.

Gefunden hat es nicht der Reviewer (dessen Prüfung „keine Änderung an
`observations/`" war für den Diff korrekt, prüfte aber nicht den
absoluten Registerstand zum eigenen Lesezeitpunkt), sondern der
**Verifier** — genau die Klasse Fund, für die diese Rolle existiert
(`docs/reviews/verifikation-slice-sdk-kotlin-sse-client-flaeche.md` <!-- d-check:status-provenance --> §6).
Die **Handlung** des Slice war unberührt richtig (kein vierter Beleg
erzeugt, `--build-context proto=proto` von Anfang an gesetzt); nur die
**Prosa-Begründung** war eine ungeprüfte, aus einem Nachbardokument in
einem älteren Stand übernommene Tatsachenbehauptung. Korrigiert bei der
Planner-Closure dieses Slice (§3/§6/§7/§8 auf den realen 3×-Stand
gehoben).

**Geprüft und verworfen:** ein eigener, neuer Beobachtungs-Eintrag für die
Unterklasse „eine Behauptung wird durch eine **parallele** Welle stale,
nicht durch die eigene Arbeit" (zur Unterscheidung von
`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, der die *eigene* Arbeit als
überholendes Ereignis voraussetzt). Der Mechanismus ist hier identisch mit
der bestehenden Definition dieser Klasse — ungeprüfte Übernahme statt
frischer Prüfung gegen den Gegenstand zum Schreibzeitpunkt; die Herkunft
des überholenden Ereignisses (eigene vs. parallele Arbeit) ändert den
Mechanismus nicht. Kein neuer Eintrag, dieser Beleg zählt hier.

Quelle:
`docs/reviews/verifikation-slice-sdk-kotlin-sse-client-flaeche.md` <!-- d-check:status-provenance --> §6
(Verdikt Punkt 6) ·
`docs/plan/planning/done/slice-sdk-kotlin-sse-client-flaeche.md` §3/§6/§7/§8
(Berichtigung) ·
`docs/plan/planning/welle-sdk-kotlin-vollabdeckung.md` §Beobachtungs-Register
(bereits korrigiert, Commit `3438f52b`).
