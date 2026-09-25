**Vorgang:** slice-sdk-readme-nutzerdoku (Review F-4, LOW; Verifikation §3 Zeile F-4)

**Fund:** `docs/user/releasing.md` (Version 1.6) nannte den Kotlin-Weg „noch unbewiesen“ und `0.1.0` als einziges Beispiel, obwohl `sdk-kotlin-v0.2.0` sowie `sdk-csharp-v0.2.0` und `sdk-python-v0.2.0` existierten und die Publish-Läufe `success` waren; dieselbe Aussage trugen die drei SDK-Workflow-Zeilen der `harness/README.md`. Der Träger lag außerhalb des Diffs, der Slice bewegte die Eigenschaft „veröffentlichte SDK-Versionen“ nicht selbst; der Reviewer fand ihn beim Nachlesen. Behoben in der Fixrunde (`releasing.md` Version 1.7, die drei Workflow-Zeilen); der Verifier fuhr `gh run list` je Workflow und `gh api` der Kotlin-Paketversionen nach (Läufe 36117929191, 36117929298, 36117929552 je `success`, Paketversion `0.2.0`).

Quelle: `docs/reviews/review-slice-sdk-readme-nutzerdoku.md` (F-4) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-readme-nutzerdoku.md` (§3 Zeile F-4). <!-- d-check:status-provenance -->
