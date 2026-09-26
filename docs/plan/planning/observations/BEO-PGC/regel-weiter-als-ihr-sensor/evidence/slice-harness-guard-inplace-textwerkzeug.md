**Vorgang:** slice-harness-guard-inplace-textwerkzeug (Review-Funde F-4 MEDIUM und F-6 LOW; Verifikation V-2, LOW)

**Fund:** Die Regel in `AGENTS.md` §3.1 (kein in-place Text-Umschreiben am Repo) reicht weiter als ihr Sensor, der PreToolUse-Guard. Die Grenze steht in `MR-003` als Grenz-Zeile; sie war in der ersten Fassung unvollständig, und der Review fand die Lücken, die sie nicht nannte.

- **F-4:** Schlüsselwörter der Shell-Kontrollstrukturen (`for … do sed -i`, `while … do`, `if … then`, `else`, `!`, `case`-Label) wurden nicht als Kopf gelesen — die typische Form einer Massenänderung mit `sed -i` passierte den Guard und stand in keiner Grenz-Zeile und in keinem Tabellenfall. Die Fixrunde löste sie im Guard.
- **F-6:** weitere Formen, die der Guard nicht las und keine Grenz-Zeile nannte (`sed -i''`, `\sed`, `busybox`, `gsed`, die Abkürzung `--in-p`, `find … -exec sh -c '…'`, andere Interpreter); dazu lag der Wortlaut der Grenz-Zeile in Kopfkommentar, `AGENTS.md` und `MR-003` nicht gleichlautend vor. Die Fixrunde löste einen Teil im Guard und nannte den Rest in `MR-003`; die `AGENTS.md`-Liste ist als Teilmenge ausgewiesen.
- **V-2 (Verifikation):** die Quote-Lesung, die F-2 des Reviews löste, hat eine Kehrseite, die die Grenz-Zeile nicht nannte: zwei Apostrophe in zwei verschiedenen Heredoc-Zeilen maskieren die Zeilen dazwischen (Falsch-Negativ). Sie steht als benannter Rand in `MR-003` und ist im Tabellentest gebunden.

**Form (Ausprägung):** derselbe Mechanismus wie in `slice-078` und `slice-079` (eine Lücke zwischen Regel und Sensor), in einer neuen Domäne (Wächter-Skript statt Doku-Sensor) und in der Variante **„benannte Grenze unvollständig“**: nicht die Lücke, sondern ihre Benennung fehlte. Wie dort ist kein zusätzlicher Sensor beauftragt; die Benennung ist die Antwort (`MR-003` Grenz-Zeile, Kurzform im Kopfkommentar des Guards und in `AGENTS.md` §3.1), der Reviewer ihr Leser. Alle drei Funde vor dem Merge von Reviewer und Verifier.

Quelle: `docs/reviews/review-slice-harness-guard-inplace-textwerkzeug.md` (F-4, F-6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-harness-guard-inplace-textwerkzeug.md` (§7 V-2). <!-- d-check:status-provenance -->
