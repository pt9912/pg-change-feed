Zustand: offen (**3×**) — die Schwelle ist erreicht, ein Ausgang ist nicht zugewiesen: der
Lese-Schritt der Closure von `welle-transformationen` liest den Eintrag. Kein Slice ist
wegen des Eintrags fällig: die Form „benannte Grenze braucht eine Adresse, die eintreten
kann“ ist eine Planner-Konvention (`implement-slice` Schritt 17 nennt sie für den
Handbuch-Aufschub), ihr Träger ist eine Zeile in Skill oder Regelwerk, also ein
Architect-Zug des Lese-Schritts. Die zwei
Funde des Erstauftretens haben je eine Adresse; der Fund des zweiten Auftretens
(`slice-transformationen-map-value`: die Größenordnung der linearen Suche in
`lookupMappedValue`, eine benannte Grenze ohne Betreiber-Aussage) hat seine Adresse in
`slice-transformationen-betriebsdoku` §2 (Übergabe-Block) und, für die Spec-Frage nach
einer Obergrenze der Paare, in `welle-transformationen` §5 (Fragen für den nächsten
Architect-Zug, Punkt (d)); der Fund des dritten Auftretens
(`slice-transformationen-e2e-wirkung`: die Restfläche der Aussage „alle Wege dieselbe
Form“ — UPDATE, DELETE und Alt-Bild auf den Stream-Wegen) hat seine Adresse in
`welle-transformationen` §3 (Closure-Kriterium zur Restfläche der Zustellwege):

- Kommentar an der gRPC-Nachricht `Change` in `proto/cdc/stream/v1/changestream.proto`
  („mit denselben Feldern wie der Domain-Typ“; der Domain-Typ trägt `origin`, die
  Live-Nachricht nach dem Pflichtenheft nicht): **offen**, Adresse dieser Eintrag. Die
  Korrektur verlangt `make proto-generate` und ist ein Implementer-Zug; sie geht mit dem
  nächsten Slice, der die `.proto` ändert, oder als eigener kleiner Zug.
- Kommentar am Feld `allowDestructive` in `tools/schema/rolloutguard/guard.go`
  („nur bekannte Fremdobjekte blockieren“, neben der Klasse „View-Signatur“ ungenau):
  Adresse ist der Plan von `slice-transformationen-antragsweg-schema`, der dieselbe Datei
  ändert (§3-Zeile zu `guard.go`).

Zähler (abgeleitet): **3×** (evidence/slice-backfill-change-origin.md,
evidence/slice-transformationen-map-value.md,
evidence/slice-transformationen-e2e-wirkung.md).
