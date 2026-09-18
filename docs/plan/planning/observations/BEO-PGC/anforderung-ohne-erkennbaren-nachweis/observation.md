# BEO-PGC/anforderung-ohne-erkennbaren-nachweis

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das
Verhältnis zwischen Lastenheft-Anforderungen und den Trägern, die ihren
Nachweis führen, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Werkzeug, das Anforderungen gegen ihre erkennbaren
Nachweis-Träger (ADR, Slice, kuratierte E2E-Coverage-Tabelle) abgleicht,
findet Anforderungen, für die **kein** dieser Träger sie zitiert — obwohl sie
im Lastenheft als Anforderung geführt werden. Das ist keine Aussage über
Vollständigkeit der Umsetzung (die Anforderung könnte real erfüllt sein), nur
über die **Auffindbarkeit** ihres Nachweises über die drei bekannten
Trägerklassen.

**Aktivierung, kein Neufund durch Inhalt:** Die sieben Waisen wurden nicht
durch inhaltliche Prüfung entdeckt, sondern durch die **Aktivierung** eines
Werkzeugs (`d-check --trace`), das vorher nicht lief. Das unterscheidet diese
Klasse von einer inhaltlichen Spec-Lücke: Die Anforderungen könnten
längst über ADR, Slice-Historie oder Code real abgedeckt sein — nur zitiert
sie kein bekannter Träger in einer für das Werkzeug erkennbaren Form.

Deklaration: `slice-d-check-trace-rtm` (`docs/plan/planning/done/…`, §1/§2 —
Plan benennt dies ausdrücklich als „Fund, der sonst nirgends steht", explizit
außerhalb des eigenen Lieferumfangs). Sieben konkrete Kennungen (Stand
2026-09-17, `make doc-trace` mit `trace.coverage` aktiv):
`LH-FA-CFG-006`, `LH-FA-CON-002`, `LH-FA-DAT-002`, `LH-FA-DAT-003`,
`LH-FA-SST-001`, `LH-FA-SST-005`, `LH-QA-REL-004`. Keine der sieben hat eine
Erwähnung in `docs/user/e2e-abdeckung.md` (eigene Gegenprobe des Verifiers,
Verifikationsbericht zu `slice-d-check-trace-rtm`, §A).

**Warum das zählt:** Ohne diese Beobachtung verschwindet der Fund im
`make doc-trace`-Output eines einzelnen Laufs — das Werkzeug ist advisory,
kein Gate, und niemand liest seinen stdout ein zweites Mal. Der
Sichtungs-Schritt künftiger Slice-Planungen (Modul 5 §Zwei Schritte vor der
Modus-Begründung) ist der einzige Ort, an dem dieser Fund wiederkehren kann,
solange er unter der 3×-Schwelle steht.
