# `make baseline-verify` — verifiziert die vendored Baseline netzlos

## Vertrag

Wird dieses Target rot, stimmt die Arbeitskopie der vendored Baseline
(`.harness/baseline/<tag>/`) nicht mit ihrer `SHA256SUMS` überein — geändert,
gelöscht oder zusätzlich eingelegt —, oder mehr als ein `<tag>`-Verzeichnis
existiert. Der adoptierte Baseline-Stand (`harness/conventions.md` §Baseline)
ist dann als Quelle nicht mehr unverändert reproduzierbar.

## Grenze — was das Grün nicht abdeckt

1. **Herkunft, nicht nur Bewegung** — `SHA256SUMS` ist selbst erzeugt und
   belegt nur, dass der Baum sich seit dem Bootstrap nicht bewegt hat, nicht
   seine Übereinstimmung mit dem Release-Asset; die hängt am Asset-Hash, der
   beim Bootstrap vor dem Entpacken geprüft wurde. Permanent (die Baseline
   ist committet, ein Nach-Abgleich läuft gegen die Registry nicht).
2. **Nur der Bestand, nicht die Anwendbarkeit** — die Prüfung sagt nichts
   über die verkörperte Form (Briefing, Konventionen, README); deren
   Konformität zur Baseline lebt in Review und Drift-Audit. Permanent.
3. **Upstream-Drift des Baseline-Stands** — ob es einen neueren
   Baseline-Tag gibt, meldet dieser Sensor nicht (siehe `harness/README.md`
   §Nicht behauptet). Permanent, bis ein Freshness-Audit (Modul 2) läuft.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`sha256sum -c` über die `SHA256SUMS` plus Bestandsabgleich; der Lauf meldet
die Zahl als „54 Dateien“ (Stand der Adoption) — sie sagt etwas über den
Baseline-Ausschnitt, nicht über das Repo.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Bestand == `SHA256SUMS`, genau ein `<tag>`-Verzeichnis |
| 1 | Abweichung (geändert/gelöscht/zusätzlich) oder doppelte `<tag>`-Verzeichnisse |

Kein Netz; das Skript läuft mit bash + coreutils ohne Fremd-Laufzeit.

## Sperren

- `<tag>`-Verzeichnis fehlt — kaputter Checkout (die Baseline ist committet,
  kein Fetch-Fall) → Checkout reparieren.

## Bindung

`harness/conventions.md` §Baseline (Adoptions-Erklärung, MR-000) — der
adoptierte Stand `v6.5.0` ist die Referenz gegen die diese Prüfung gilt.
