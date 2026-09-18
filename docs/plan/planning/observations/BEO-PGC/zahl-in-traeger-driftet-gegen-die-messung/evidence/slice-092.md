# Beleg: slice-092

Vorgang: `slice-092` — Coverage Cluster D1 (Anwendungs-Kern).

Fund: **Drei Funde desselben Vorgangs**, alle an **Sätzen über** Zahlen, keiner
an einer Rechnung:

- **F-1** (Review zu `slice-092`, MEDIUM): Ein Dokument, **vier Zahlen für denselben
  Gegenstand** (10 · 11 · 12 · 13 an verschiedenen Stellen, zwei davon an einer
  Aussage, die eine andere Größe meint). Entstanden durch einen
  Berichtigungs-Zug, der einen Glob-Fehler beheben sollte — und dabei die Zahlen
  an Stellen setzte, die **andere** Aussagen tragen.
- **D-1** (Delta-Review zu `slice-092`, MEDIUM): Die Überschrift der Behebung sagte
  „**Drei** Zahlen, drei Dinge" und führte **vier**. Dritte Runde an derselben
  Stelle.
- **D-2** (ebd., MEDIUM): Der Herkunfts-Marker war **invertiert** — die vier
  gedeckten Zahlen sind **gemessen** (je mit Lauf), die **Differenz** ist nach
  `AGENTS.md` §3.12 Instanz A ausdrücklich **abgeleitet**; der Satz nannte es
  umgekehrt. Und die „Deutung" war **contra-gemessen**: die zwei Bänder liegen
  auf demselben Produktionsstand (Nenner 1903 an beiden, kein Produktcode
  dazwischen), die Verschiebung liegt an den **Tests**.

**Ein Vorgang, eine Zählung.** Alle drei fallen in dieselbe Klasse und in
denselben Vorgang — der Zähler bewegt sich **einmal**, obwohl drei Berichte sie
führen.

**Was diese Vorkommen von den früheren unterscheidet:** die früheren Belege
dieser Klasse trafen **Bestands-Zahlen** (ein Wert, der gegen eine neue Messung
driftet). Hier trifft die Klasse **den Korrektur-Vorgang selbst** — die
Berichtigung einer Zahl erzeugt die nächste falsche Zahl, und zwar in der
Reihenfolge Glob-Fehler → vier Zahlen → Zähl-Wort → invertiertes Etikett. Vier
Runden, jede an einem Satz, der eine Zahl erklärt.

Quelle: Review zu `slice-092` (F-1) ·
Delta-Review zu `slice-092` (D-1, D-2) ·
Verifikationsbericht zu `slice-092` (V-1, V-2) ·
`docs/plan/planning/done/slice-092-coverage-cluster-d1.md` §1 ·
`harness/sensors/coverage-gate.md` §Zählbasis (berichtigt in `1d0d11e`).
