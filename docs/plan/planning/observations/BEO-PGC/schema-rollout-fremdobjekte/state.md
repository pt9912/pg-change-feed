Zustand: **geplant** → Träger `docs/plan/planning/open/schema-rollout-zentrale-idempotenz-wache.md`.

Vorheriger Zustand (Commit `7d0bf05`): `verkörpert` → `AGENTS.md` §3.14
("Ein Aufrufer von `make schema-rollout` gegen ein möglicherweise bereits
migriertes Ziel trägt seine eigene Idempotenz-Wache"). Der Architect-Verdikt
zum Lese-Schritt (F-1) fand real, dass diese Entscheidung nur zwei Optionen gegeneinander
abgewogen hatte — (a) Fremdobjekte ins neutrale Modell überführen (an
`POST_EXECUTE_DRIFT`/Exit 5 gebunden, bleibt korrekt verworfen) und (b) ein
generisches `--allow-destructive`-Handling (zu grobkörnig, bleibt korrekt
verworfen) — nicht aber eine dritte: die Idempotenz-Wache **zentral im
`schema-rollout`-Makefile-Target selbst** statt bei jedem Aufrufer einzeln.

**Die dritte Option, ernsthaft geprüft:**

- **Tractability.** Der Bericht (`tools/schema/plan.yaml`) ist
  strukturiertes JSON, kein Freitext — jede Operation trägt `objectType`
  (`TABLE`/`VIEW`/…), `kind` (`CreateView`/`DropView`/…) und `path` (die
  Objekt-Segmente). Ein Blocker gegen genau die sechs bekannten
  Fremdobjekte ist damit maschinell von einem Blocker gegen ein
  unbekanntes Objekt unterscheidbar — reale Probe gegen einen frischen,
  unblockierten Rollout-Report bestätigt die Feldform. Die Klassifikation
  selbst ist also **tractable**, anders als eine Freitext-Heuristik es
  wäre.
- **Eine einfachere Variante gibt es zusätzlich**, die der Reviewer nicht
  benannt hatte: statt Blocker zu klassifizieren, genügt derselbe
  Existenz-Check, den heute jeder Aufrufer einzeln fährt
  (`to_regclass('cdc.source_table')`), zentral **vor** dem `--execute`-
  Schritt im Makefile-Target selbst — ohne jede Objekt-Allowlist. Diese
  Variante ist am billigsten, hat aber einen realen Nachteil, den keiner
  der Aufrufer-seitigen Checks heute trägt: Sie würde `make schema-rollout`
  für **jeden** Aufrufer — auch einen künftigen, der eine echte
  inkrementelle Schema-Änderung gegen ein bereits migriertes,
  langlebiges Ziel ausrollen will — grundsätzlich zu einem stillen No-op
  machen, sobald irgendein früherer Rollout stattfand. Das ist ein
  Vertragsbruch am gemeinsamen, von allen Aufrufern geteilten Target,
  nicht nur eine Bequemlichkeit für die drei bekannten
  Wegwerf-/Bootstrap-Aufrufer, die heute existieren — real geprüft: **alle**
  aktuellen Aufrufer (`tools/harness/run-integration-tests.sh`,
  `tools/bench-lib.sh`, `tools/schema/apply-rollout.sh`) fahren ohnehin
  gegen eine frisch angelegte, leere Ziel-DB (`$COMPOSE down -v` o. ä.
  vor jedem Lauf) — real geprüft (`grep -rn "to_regclass" --include="*.sh"
  --include="Makefile" --include="*.mk" .`): nur **ein** persistenter
  Aufrufer trägt heute den Skip-Bedarf, `examples/bootstrap.sh`
  (`make example-demo-up`) — die drei anderen realen `example-run-*`-Ziele
  starten nur bereits gebaute Images gegen dieselbe, von
  `example-demo-up` verwaltete Umgebung, kein zweiter unabhängiger
  Aufrufer. Diese einfache Variante bliebe also heute folgenlos, ist aber
  als **dauerhafte** Eigenschaft eines geteilten Ziels riskanter als der
  Nutzen für den aktuell einen betroffenen Aufrufer rechtfertigt — kein
  Ein-Zeilen-Fix ohne Nebenwirkung, unabhängig von der genauen
  Aufruferzahl (ein geteiltes Target würde auch für jeden **künftigen**
  Aufrufer zum stillen No-op).
- **Die vom Reviewer skizzierte Klassifikations-Variante** (nur überspringen,
  wenn *alle* Blocker exakt die bekannte Liste sind, sonst wie bisher
  Exit 8) trägt diesen Nachteil nicht, weil sie nicht pauschal überspringt,
  sondern gezielt. Sie ist aber kein Ein-Zeilen-Fix: Sie braucht (1) einen
  JSON-Parsing-Schritt im bislang `jq`-freien Docker-only-Makefile-Fluss,
  (2) eine Klärung, ob d-migrate überhaupt erlaubt, einzelne Operationen
  gezielt zu bestätigen und den Rest eines Laufs trotzdem auszuführen
  (ungeklärt — noch nicht gegen das gepinnte Image geprüft), und (3) einen
  eigenen Negativ-Testfall (ein künstlich eingefügtes, nicht gelistetes
  destruktives Ziel muss weiterhin blockieren). Das ist ein eigener,
  recherche-gebundener Liefer-Umfang, kein Nachtrag in dieser Beobachtung.
- **Wartungsaufwand — differenzierter als die Ausgangsfrage.** Eine
  Objekt-Allowlist bräuchte zwar bei jedem neuen `nacharbeit-*.sql`-Skript
  einen Nachtrag — aber an **derselben Stelle**, an der ohnehin schon die
  Aufrufzeile des Skripts im Makefile-Target ergänzt wird (Kolokation im
  selben Commit, kein zusätzliches Erinnern an einer entfernten Stelle).
  Das ist ein echter Vorteil gegenüber dem heutigen Zustand, in dem N
  verschiedene Aufrufer-Dateien unabhängig voneinander denselben
  Existenz-Check erfinden müssen (real dreifach unabhängig geschehen). Der
  heutige Aufrufer-Check selbst braucht dagegen **keine** Allowlist-Pflege
  — er prüft nur, ob überhaupt schon migriert wurde, unabhängig vom
  Fremdobjekt-Bestand; das in der vorigen Fassung dieses Eintrags
  unterstellte "wächst mit jedem Skript" traf also nur auf den
  *Blocker*, nicht auf den *Aufrufer-seitigen Check* selbst zu — eine
  Präzisierung gegenüber der vorherigen Fassung.

**Verdikt:** Die zentrale, blocker-klassifizierende Wache ist die
plausibel bessere Lösung — geteilte statt verteilte Brüchigkeit, an genau
der Stelle behoben, die jeder neue Fremdobjekt-Zug ohnehin schon anfasst.
Sie ist aber real recherche- und implementierungsgebunden (d-migrates
Bestätigungs-Mechanismus ist ungeklärt) und damit kein Nachtrag zu dieser
Beobachtung, sondern ein eigener Slice:
`docs/plan/planning/open/schema-rollout-zentrale-idempotenz-wache.md`.
`AGENTS.md` §3.14 bleibt bis zu dessen Lieferung die geltende
Zwischenlösung — sie schützt real, auch wenn sie nicht die
architektonisch beste Endform ist, und wird bei Closure jenes Slice
angepasst oder gestrichen.

Zähler (abgeleitet): 3× (evidence/slice-016.md, evidence/slice-063.md,
evidence/slice-beispiele-compose-bootstrap.md) — Schwelle erreicht;
Ausgang jetzt `geplant` statt `verkörpert` (Träger siehe oben).
