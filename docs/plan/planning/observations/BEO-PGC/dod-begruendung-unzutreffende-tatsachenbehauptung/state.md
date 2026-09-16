Zustand: **offen — 4× erreicht, Ausgang noch nicht zugewiesen.** Der
Lese-Schritt gehört der **laufenden Welle-Closure** (`welle-20`), nicht der
Slice-Closure: Modul 6 lässt den Ausgang ab 3× der laufenden Welle-Closure
zuweisen („Bei 3× wandert der Eintrag in die Steering-Loop-Einträge der
laufenden Welle-Closure"). Bis dahin ist `offen` der zulässige, vorübergehende
Stand — wie beim Schwester-Eintrag
`BEO-PGC/generierte-artefakte-ohne-sync-sensor`.

**Ein Träger ist inzwischen benannt, aber nicht gebaut:** der Architect-Zug zu
`ADR-0078` hat für diese Wurzel eine Regel formuliert — *jeder Zahlenwert in
einem `Accepted`-Dokument trägt seinen Ursprung* (gemessen / übernommen /
**abgeleitet**); ein abgeleiteter Wert darf nie als gemessen erscheinen. Als
Träger schlägt er den Kopf von `docs/plan/adr/README.md` vor oder `AGENTS.md`
§3.7. **Das Benennen ist nicht das Bauen** — die Verkörperung ist Planner-Arbeit
und steht aus; deshalb ist der Ausgang hier *nicht* vorwegnehmend als
`verkörpert` gesetzt.

Zähler (abgeleitet): **4×** (evidence/slice-036.md, evidence/slice-082.md,
evidence/slice-081.md, evidence/slice-083.md) — **Schwelle erreicht**, Ausgang
beim Lese-Schritt der `welle-20`-Closure. Der vierte Beleg ist ein weiterer
Vorgang derselben Klasse und kein neuer Handlungsbedarf: `slice-083` berief sich
auf **nicht existierende** Nachbar-Clients („wie die drei anderen"), und der
Verfasser war der Planner. Das Erstvorkommen (`slice-036`) wurde seinerzeit **ohne
Kennung** notiert: das Review nannte das Label, legte aber kein Verzeichnis an.
Der Eintrag entstand mit dem zweiten Auftreten und zitiert das erste über
**dasselbe Label** — die Zuordnung ist belegt, nicht abgeleitet. Zwei Funde **im
selben** Vorgang (`slice-082`: Review F-1 und Verifikation V-1; `slice-081`:
Review F-1 und die Berichtigung in §1/§2) sind je *eine* Gelegenheit — der
Zähler misst Wiederholung über Vorgänge, nicht die Zahl der Funde.
