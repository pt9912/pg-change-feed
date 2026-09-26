**Vorgang:** slice-sdk-kotlin-cloudsmith (Review F-5, INFO)

**Fund:** Die Probe in der Stufe `pack` von `sdks/kotlin/Dockerfile` bindet die zwei Aufgabennamen `publishMavenPublicationToGitHubPackagesRepository` und `publishMavenPublicationToCloudsmithRepository` am Literal im Dockerfile; die zwei Aufruf-Literale im Workflow (`.github/workflows/sdk-kotlin-release.yml`) sind davon nicht gebunden — das Dockerfile sieht `.github/` nicht (Bau-Kontext `sdks/kotlin`). Ein Auseinanderlaufen färbt keinen Sensor, ein falscher Name endet im Tag-Lauf mit „Task not found“. Ebenso trägt die Prüfung der Upload-URL ein `grep -F`, das eine Kommentarzeile erfüllen würde, und die Namensprüfung hängt am Ausgabeformat von `--dry-run`. Gefunden hat es der Reviewer; als benannte Grenze in Plan §6 geführt. Im realen Lauf 36201941235 trugen beide Stellen dieselben Namen, kein `Task … not found` im Log.

Quelle: `docs/reviews/review-slice-sdk-kotlin-cloudsmith.md` (F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-kotlin-cloudsmith.md` (§7 Zeile F-5). <!-- d-check:status-provenance -->
