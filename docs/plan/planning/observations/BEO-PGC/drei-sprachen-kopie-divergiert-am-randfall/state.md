Zustand: gestrichen (**3×**, akzeptiertes Negativ; Ausgang am Ende). Behoben im Vorgang:
die Randwert-Regel steht in einer Formulierung an den Trägern (`spec/pflichtenheft.md` §7
Historie, `docs/user/benutzerhandbuch.md` Absatz nach den `**SDK:**`-Absätzen, die drei
`sdks/*/README.md`) und je Sprache als Test mit sechs Fällen; die Eingabeseiten-Mutationen
(Verifikation §4: C1 bis C3, P1, P2, K1, K2) färben je Sprache rot.

Zähler (abgeleitet): **2×** (evidence/slice-backfill-sdk-origin.md,
evidence/slice-sdk-readme-nutzerdoku.md). Der zweite Beleg trifft dieselbe Klasse an
einer **Doku-Aussage** statt an einem Lesemodell: die README-Zeile über das Bild bei
INSERT/DELETE kopierte den Wortlaut der Schwester-READMEs, und die Kotlin-Bibliothek
liefert dort `JsonNull` statt `null` (Review F-2, behoben, Kotlin-README nennt
`JsonNull`). Unter der Schwelle, kein Ausgang zugewiesen.

**Dritter Beleg, Schwelle erreicht (slice-routing-sdk-beispiel-target, Review F-3):** der
Zähler steht bei **3×** (evidence/slice-backfill-sdk-origin.md,
evidence/slice-sdk-readme-nutzerdoku.md, evidence/slice-routing-sdk-beispiel-target.md); der Ausgang
gehört zum Lese-Schritt der Closure von `welle-routing`. Einordnung nach der Zählregel: gezählt wird
der **Satz**, nicht das Verhalten. Das Verhalten (Python liest das leere Ziel am NATS-Stream als „kein
Ziel“, C# und Kotlin weisen es am Subjekt-Bau ab) ist eine **gewollte API-Form-Differenz** — im Plan
benannt, am Parameter-Schnitt begründet, kein Divergenz-Fund einer Bibliothek am Randwert. Der Satz der
Python-README war dagegen dem Wortlaut der Schwester-READMEs gefolgt und enger als der Code: dieselbe
Form wie der zweite Beleg (Doku-Aussage über einen Randfall, in drei Sprachen mit unterschiedlichem
Verhalten kopiert; vom Reviewer vor dem Merge gefunden). Der Eintrag ist `offen`, nicht `verkörpert`,
der Deckel gilt nicht. Gegenläufig belegt: HTTP, SSE und gRPC lasen dieselbe Eingabetabelle
(`null`/`""`/`eu`/`a&b=c`) in allen drei Sprachen gleich, weil der Plan sie als Pflicht-Eingabesatz
je Sprache nannte (der Reviewer mutierte 13 Stellen in 10 Läufen, der Verifier 10 Mutationen, je selbst gefahren) — die Gegenmaßnahme „derselbe Eingabesatz je Sprache im Plan“ ist wirksam, aber nicht für
die Sätze der README.

**Verwandt, nicht gleich:** `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
(die Kopie trägt einen Fehler des Vorbilds weiter) und
`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (zwei Quellen nennen dieselbe
Aussage verschieden).

**Gegenmaßnahme beobachtet, kein neuer Anfall (slice-sdk-regel-realserver-e2e, §6/§7):**
die drei SDK-Realserver-Tiers brauchten für ihre Regel-Testphasen dieselbe
Vorbereitung (Tabelle anlegen, `cdc.enable_table`, `cdc.set_transformation`,
Poll); statt sie unabhängig je Sprache zu implementieren, trägt eine
gemeinsame Hilfsdatei (`tools/harness/lib-sdk-rule-fixture.sh`) die
Aufrufform genau einmal — der Suchlauf des Slice fand sie ausschließlich dort
(0 Kopien in den drei Runnern). Das ist keine Auflösung der Beobachtung: die
**Produktivcode**-Entscheidung je Sprachbibliothek (wie diese eine Regel
JSON-Randwerte liest) bleibt unabhängig und divergenzfähig — hier war der
Gegenstand eine **Testvorbereitung**, kein Lesemodell. Zähler bleibt bei 2×;
kein neuer `evidence/`-Eintrag, weil kein Divergenz-Fund vorliegt.

**Vermerk, kein neuer Anfall (slice-routing-sdk-realserver-e2e, Review F-6):** die Sammler der
Routing-Phasen enden in den drei Tiers verschieden — C# beendet die Verbraucher vor der Auswertung,
Kotlin und Python werten einen Schnappschuss, während die Daemon-Threads weiterlesen. Das ist eine
Divergenz der **Testhilfe**, kein Divergenz-Fund am Randwert: eine Change nach dem Schnappschuss
liegt außerhalb des geprüften Fensters, kein Ergebnis ändert sich, Schwellen, Prüfreihenfolge und
Meldetexte sind gleich gelesen (Review „Negativbefunde“). Die Gegenmaßnahme „eine Fixture-Quelle“
hat wieder gewirkt (`tools/harness/lib-sdk-route-fixture.sh` trägt Vorbereitung und Phasenablauf
einmal). Zähler bleibt bei **3×**, keine neue Datei.

**Ausgang (Lese-Schritt der Closure von `welle-routing`, Architect-Verdikt 2026-10-02): gestrichen (akzeptiertes Negativ).** Der Zustand ist `gestrichen`. Begründung (Verdikt `architect-verdict-welle-routing-lese-schritt` §3.3): (a) alle drei Belege sind LOW und der Reviewer fand jeden vor dem Merge; das dritte Auftreten ist ein Satz der Python-README, enger als der Code, die Verhaltens-Differenz selbst (leeres Ziel am NATS-Stream) eine gewollte API-Form-Differenz. (b) Der mögliche Schaden ist eine zu enge Doku-Aussage oder ein Randwert, den der Server nicht sendet. (c) Die wirksamen Gegenmaßnahmen stehen im Code und im Plan der Slices und werden dort vom Reviewer gelesen: ein Pflicht-Eingabesatz je Sprache und eine Fixture-Quelle (`tools/harness/lib-sdk-route-fixture.sh`, `lib-sdk-rule-fixture.sh`). Eine Plan-Pflicht wäre Zeremonie bei kleinem Schaden. **Wiederaufnahme-Trigger:** ein Divergenz-Fund nach dem Merge, oder ein Fund mit Schwere ab MEDIUM. Kein Folge-Artefakt; die Alternative (Plan-Pflicht in `.claude/commands/plan-welle.md`) bleibt für den Auftraggeber verfügbar, der Architect empfiehlt sie nicht.
