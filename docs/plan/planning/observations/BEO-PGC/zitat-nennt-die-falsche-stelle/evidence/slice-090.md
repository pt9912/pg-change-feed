# Beleg: slice-090

Vorgang: `slice-090` — das Sync-Gate des generierten Protobuf-Codes.

Fund: Das in diesem Zug neu angelegte `harness/sensors/generated-sync.md` zitiert
`ADR-0084` mit einer falschen Festlegungs-Nummer — „Festlegung 2 = der Befund
nennt Datei und Zeile". Das ADR führt **vier** Festlegungen (`**1.` „Der
Protobuf-Code bekommt ein Gate" mit den zwei Bedingungen als Unterpunkten, `**2.`
die E2E-Abdeckungstabelle, `**3.` `plan.yaml`/`down.sql`, `**4.`
`harness/image-hash.txt`); die gemeinte Stelle ist **Festlegung 1, zweite
Bedingung**. Das ADR zitiert sich in seinen Re-Evaluierungs-Triggern selbst
**richtig** (`:297` „das Gate aus Festlegung 2" = die E2E-Tabelle) — das Original
war die ganze Zeit korrekt.

**Der Ursprung ist eine Zusammenfassung.** Die falsche Nummer stammt aus der
**Kopfzeile** von `review-slice-090.md:22` („Festlegung 1 = Baum nicht
schreiben; Festlegung 2 = Befund nennt Datei und Zeile"), die ihrerseits aus dem
Delta-Review als übernommen benannt wird; dieselbe Kurzform steht in
`verify-slice-090.md:6`. Ein Bericht ist **Lauf-Beleg**, kein Zitat-Träger —
gelesen wurde er hier wie eine Quelle, und die Wiedergabe wanderte in einen
**stehenden** Träger.

**Gefunden hat es ein fremder Kontext.** Der Delta-Review hat das ADR
aufgeschlagen und die vier Überschriften gezählt; weder das diff-skopierte
Review der ersten Runde noch der Verifier der ersten Runde hatten einen Anlass,
die Nummer zu prüfen — sie stand in einem neuen Satz, der eine **alte** Aussage
zu beheben schien. Beide Berichte sind als Lauf-Belege **nicht** rückdatiert;
die Korrektur steht in §7 der Closure-Notiz.

Quelle: `docs/reviews/review-slice-090-delta.md` (D-1) ·
`docs/plan/adr/0084-sync-gate-fuer-generierte-artefakte.md` §Entscheidung ·
`docs/reviews/review-slice-090.md:22` (Ursprung der Nummer) ·
`harness/sensors/generated-sync.md` §Bindung (berichtigt in `81f1fff`).
