**Vorgang:** slice-063 (Blocker-Report)
**Fund:** Der Planner-Koordinator committete und pushte
`docs/reviews/blocker-slice-063.md` mit drei nackten `ADR-*`-Kennungen
ohne Link. Anders als bei den drei vorherigen Belegen dieser Beobachtung
wurde der Fehlschlag diesmal NICHT vor dem Push gefangen: `make gates`
lief ungepiped, der Exit-Code wurde korrekt separat erfasst (`EXIT=2`,
sichtbar im eigenen Tool-Output), aber die nachfolgende `git push`-Aktion
lief im selben Arbeitsschritt, ohne dass der bereits sichtbare rote
Exit-Code die Aktion tatsächlich blockierte — der Push geschah, bevor der
Fehler bemerkt wurde. Das widerspricht der Prämisse des
`gestrichen`-Architect-Verdikts dieser Beobachtung
(`architect-verdict-report-nackte-id-ohne-link.md`): „solange er
weiterhin vor jeder folgenreichen Konsequenz vom Sensor gefangen wird" —
hier war die Konsequenz (Push auf `main`) bereits eingetreten, als der
Fund bemerkt wurde. Viertes Auftreten, nach `gestrichen` bei 3×; die
zugrunde liegende Fehlerklasse (Exit-Code korrekt erfasst, aber nicht
tatsächlich als Gate genutzt) unterscheidet sich von den drei
vorherigen Belegen (die betrafen entweder Pipe-Maskierung oder
Selbstkorrektur vor Commit) und rechtfertigt eine erneute Prüfung des
gestrichenen Ausgangs.
