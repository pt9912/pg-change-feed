# Beleg: slice-094

Vorgang: `slice-094` — Coverage Cluster A (Prozess-Rand).

Fund: **Ein geschätzter Wert wurde auf dem Weg durch drei Rollen zu einer
Grenze** — und war um vier Statements zu niedrig.

[ADR-0082](../../../../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Kontext (4a) notiert für `cmd/pg-change-feed`
**„≈45 netzlos erreichbar"**. Die vier Differenz-Statements sind die **Aufrufe** der vier
dienstgebundenen Sondermodi (`main.go:44` `Healthcheck`, `:62`
`RegisterConsumer`, `:86` `AcknowledgeConsumer`, `:101` `Diagnose`).

Der Weg, Station für Station:

1. **[../../../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md]** (Accepted): „≈45 netzlos erreichbar" — eine **Schätzung**, mit
   Tilde, auf dem Stand von zwei Läufen. Der Ursprung ist dort korrekt benannt.
2. **Slice-Plan `slice-094` §1** (Planner): „Die restlichen **≈10** sind
   **nicht** erreichbar und werden **benannt**, nicht kaschiert." — Die
   Schätzung ist zur **Decke** geworden; im §2 (LP3) trägt sie den **Grund** für
   die vier Stellen.
3. **`harness/sensors/coverage-gate.md`** (Implementer): „die vier offenen sind
   je **ein** Aufruf eines dienstgebundenen Sondermodus und **liegen außerhalb
   des netzlosen Tiers**." — Die Decke ist jetzt eine **Grenze** in einem
   stehenden Träger.
4. **Review** zu `slice-094` (F-1, HIGH): gemessen — ein Test, der dieselben
   vier Modi mit **vollständigem** ENV fährt, deckt sie **alle vier**: `cmd`
   **49 von 49**, Gesamt **1581/1903 = 83,08 %**, netzlos, EC 0.

**Niemand hat den Wert geändert; nur seine Verbindlichkeit.** Die Tilde fiel
stückweise weg, und jede Station hielt den Wert für ein Zitat aus der vorigen.
Das ist der Kern: **es gab keine Messung, an der der Wert hätte driften
können** — deshalb ist es nicht `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`.
Und jeder Träger trug seinen Satz — er zitiert nur nicht als Zitat.

**Was geholfen hat, war eine Handlung, kein Werkzeug:** der Reviewer hat an der
**letzten Station** gemessen — dort, wo der Wert zur Grenze wurde. Eine
Testdatei, vier Statements.

Quelle: Review zu `slice-094` (F-1, HIGH) ·
`docs/plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md`
§Kontext (4a) · `docs/plan/planning/done/welle-20/slice-094-coverage-cluster-a.md` §1/§2
(der Plan ist nach der Closure in `done/`; die Fassung bei der Planung trug die
Decke in `in-progress/`) ·
`harness/sensors/coverage-gate.md` §Grenze Punkt 1 (berichtigt in `8292766`).
