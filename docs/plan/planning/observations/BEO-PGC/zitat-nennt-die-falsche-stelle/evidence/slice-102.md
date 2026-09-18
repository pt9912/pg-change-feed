# Beleg: slice-102

Vorgang: `slice-102` — gRPC-Client in C# (`ADR-0090`).

Fund: Der Slice-Plan zitierte in §6 (letzter Risiko-Punkt) und §8
(„Vorgelagert — offene Beobachtungen sichten") die Aussage „ihre
Bau-Kontexte erreichen sie heute nicht" als „`slice-095`s §1". Gemessen
(`grep -n "erreichen" docs/plan/planning/done/slice-095-beispiel-clients-
drei.md docs/plan/planning/done/slice-097-umzug-vertragsflaeche.md` → Treffer
nur in `slice-097:107`): Der Satz steht wörtlich in `slice-097` §1
(„Ausdrücklich NICHT in diesem Slice", zweiter Punkt), nicht in `slice-095`.
`slice-095` (done) enthält weder das Wort „erreichen" noch „Kotlin".

Gefunden hat es der Reviewer (F-3,
Review zu `slice-102`), unabhängig reproduziert vom Verifier
(Verifikationsbericht zu `slice-102`, #11). Anders als beim Erstauftreten
(`slice-090`, Übernahme aus einer Bericht-Kopfzeile) ist der Ursprung hier
nicht eindeutig rekonstruiert — plausibel ist eine Verwechslung zweier
benachbarter, thematisch ähnlicher `done/`-Slices aus derselben Matrix-Serie
(`slice-095`, `slice-097`), beide zur Beispiel-Client-Erweiterung gehörig.

**Behoben:** Beide Fundstellen im Slice-Plan (§6, §8) sind im Rahmen dieser
Closure auf „`slice-097`" korrigiert.

Quelle: Review zu `slice-102`, F-3 ·
Verifikationsbericht zu `slice-102`, #11 ·
`docs/plan/planning/done/slice-097-umzug-vertragsflaeche.md:107` ·
`docs/plan/planning/done/slice-102-grpc-client-csharp.md` §6/§8.
