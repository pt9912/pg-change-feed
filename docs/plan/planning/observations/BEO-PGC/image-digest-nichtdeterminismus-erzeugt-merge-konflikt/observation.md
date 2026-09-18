# Image-Digest als Lauf-Beleg erzeugt inhaltsloses Merge-Rauschen

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft `harness/image-hash.txt`
als Werkzeug-Ausgang, keine eigene Sub-Area im Sinn der Modus-Deklaration).

`harness/image-hash.txt` trägt den OCI-Image-Digest des letzten `make
image`-Laufs als reinen Lauf-Beleg, ausdrücklich **kein**
Inhalts-Fingerabdruck ([`ADR-0044`](../../../../adr/0044-image-beleg-semantik.md)
§Entscheidung Punkt 1, dort bereits real belegt: bit-identisches Binary,
unterschiedlicher Digest zwischen zwei Builds derselben Commit-Historie).
Weil der Digest builder-/lauf-gebunden ist, erzeugt jeder Zug, der `make
image` erneut laufen lässt, potenziell einen Diff auf genau dieser einen
Zeile — unabhängig davon, ob sich der tatsächliche Build-Inhalt geändert
hat. Zwei parallele Branches, die beide `make image` laufen lassen,
kollidieren beim Merge auf dieser Zeile, ohne dass einer von beiden
inhaltlich etwas Konfliktträchtiges beigetragen hätte — der Auftraggeber
hat das real als Merge-Störung benannt.

Ein Kandidat für die bessere Lösung wurde in diesem Gespräch mitgeführt
(unentschieden, Architect-Zuständigkeit): den Datei-Inhalt von
„OCI-Image-Digest" auf „sha256 des aus dem Image extrahierten Binaries"
umstellen (`ADR-0044`s bereits verglichene, aber verworfene Option B) —
deterministisch bei unverändertem `cmd/`/`internal/`/`gen/`/`proto/`-Inhalt,
also ohne Diff/Konflikt für Züge, die diese Pfade nicht berühren.
`ADR-0044`s Gegenbegründung für Option B („der Digest ist die
Registry-Adresse des exportierten Images") trägt praktisch nicht mehr real
geprüft: `make image` baut ausschließlich lokal (`docker buildx build
--load`), kein `docker push` existiert irgendwo im Repo — der Digest
adressiert also aktuell nichts tatsächlich Auflösbares in einer Registry.

## Benannt, nicht gezählt

- **`slice-nats-drittstream-example-go`** (zum Zeitpunkt dieser Beobachtung
  noch nicht in `done/`, kein abgeschlossener Vorgang — daher hier benannt,
  nicht als Evidence-Datei gezählt): Ein `make image`-Lauf während der
  Implementierung (ausgelöst durch eine — sich als falsch herausstellende —
  Annahme, `examples/**` läge im Build-Kontext des Root-`Dockerfile`) änderte
  den Digest, obwohl kein von diesem Slice berührter Pfad in der
  `.dockerignore`-Whitelist liegt (`cmd/`, `internal/`, `gen/`, `proto/`,
  `go.mod`, `go.sum`, zwei `ADR-0085`-Ausnahmen). Der fälschlich committete
  Digest-Wechsel wurde im selben Zug wieder zurückgenommen (Review-Fund F-1
  zu diesem Slice). Der
  Auftraggeber benannte in diesem Zusammenhang das allgemeinere Ärgernis:
  ein Digest-Diff auf dieser Datei stört beim Mergen, unabhängig davon, ob
  er berechtigt ist.
