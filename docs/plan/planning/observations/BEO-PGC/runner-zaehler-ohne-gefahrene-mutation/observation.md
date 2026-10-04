# Ein Zähler einer Runner-Phase trägt seine Zusage nur hergeleitet, wenn die Mutation an ihm nicht gefahren ist

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Phasen von `tools/harness/run-integration-tests.sh`).

Eine Zusage der Runner-Phase („beim Ladefehler entsteht kein Replication-Slot“,
[`SPEC-034`](../../../../../../spec/pflichtenheft.md)) wird vom Unit-Test des Pakets an einer Mutation
gebunden (`Run` lädt das Paar nach `postgresstorage.New`: rot), vom Slot-Zähler des Runners aber nur an
dem Zustand „kein Slot“. Der Slot entsteht erst in `NewStream`; der Zähler färbt sich daher bei der
Mutation „Laden nach der Datenbankverbindung“ nicht und ist nur an einer Mutation gebunden, die das Laden
hinter `NewStream` zieht. Diese Mutation verlangt einen vollständigen Lauf von `make test-integration`
(rund fünfzehn Minuten, dazu `make image-mutation`); sie ist im Slice nicht gefahren, die Bindung des
Zählers steht als *hergeleitet*.

## Benannt, nicht gezählt

Die Verifikation (`docs/reviews/verifikation-slice-tls-http-grpc-server.md` §3) <!-- d-check:status-provenance -->
führt die
Runner-Mutation als *hergeleitet*; die Zusage „kein Slot bei Ladefehler“ ruht auf dem Code-Aufbau (Laden
vor `NewStream`, gelesen) und dem Unit-Test.
