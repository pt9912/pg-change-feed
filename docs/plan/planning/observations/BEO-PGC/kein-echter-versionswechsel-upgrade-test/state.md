Zustand: **verkörpert** — Ausgang: **verkörpert** → Phase U in
`tools/harness/run-sdk-altserver-tests.sh` (Vertrag `harness/targets/sdk-altserver.md`) ·
seit slice-upgrade-versionswechsel-alt-image. Träger der Entscheidung ist
[`ADR-0148`](../../../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)
Teil 2 (Folge des erfüllten Re-Evaluierungs-Triggers 1 von `ADR-0064`): ein anderer
Server-Build (das veröffentlichte Image 0.5.0) erfasst Zeilen, ein `--force-recreate` auf das
`:dev`-Image lässt den Datenstand der Quelle über `cdc.changes` gleich lesbar und setzt die
Erfassung fort — gemessen in zwei Läufen mit Exit 0, drei Mutationen mit gesehenem Rot
(evidence/slice-upgrade-versionswechsel-alt-image.md).

**Begrenzung des Ausgangs.** „Datenstand“ heißt hier Zeilenzahl plus `md5` über alle
`cdc.changes` der Quelle, ein Alt-Stand, ein Lauf-Paar (0.5.0 → `:dev`), eine Tabelle. Ein
Schemawechsel über Versionen ist ungemessen (akzeptiertes Negativ, `ADR-0148` Option C); der
Ausgang lautet nicht „der Versionswechsel ist vollständig belegt“. **Wiederöffnung:** ein
beobachteter Betriebs-Upgrade-Fehler, der am konstanten Schema vorbeiläuft.

Zähler (abgeleitet): 3× (evidence/slice-063.md,
evidence/slice-sdk-0-6-kompatibilitaet-messen.md,
evidence/slice-upgrade-versionswechsel-alt-image.md).
