**Vorgang:** slice-sdk-readme-nutzerdoku (Review F-3 mit Umfangserweiterung, Verifikation §3 Zeile F-3, §5 und §4 G/H)

**Fund:** Die README-Neufassung nahm die Metadaten-Felder und die README in den Umfang, die öffentlichen API-Kommentare nicht. Der Reviewer fand die Kennungen dort: 8 Python-Dateien mit 93 Zeilen, dazu 20 Laufzeit-Fehlertexte mit `SPEC-`-Kennung (Python 8, C# 6, Kotlin 6), im Wheel als `help()` und IDE-Tooltip sichtbar. Der Nutzer erweiterte den Umfang auf alle öffentlichen API-Kommentare; die Fixrunde bereinigte `sdks` (499 Zeilen in 98 Dateien → 0), die `.proto`-Kommentare (die in den Stubs aller Sprachen stehen) und ergänzte die Wächter (`test_public_text.py`, `make sdk-public-doc-check`). Der Verifier fuhr den Artefakt-Scan über Wheel, sdist, `.nupkg`, Jar und Sources-Jar (0 Treffer) und zwei Wächter-Mutationen (Stub-Pfad über die `.proto`, C#-Datei), beide rot.

Quelle: `docs/reviews/review-slice-sdk-readme-nutzerdoku.md` (F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-readme-nutzerdoku.md` (§3 Zeile F-3, §4, §5). <!-- d-check:status-provenance -->
