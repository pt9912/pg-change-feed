# BEO-PGC/plattform-verhalten-nur-vom-betreiber-pruefbar

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhalten
externer Plattformen, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Release-Weg trägt Eigenschaften einer externen Plattform,
die weder der Agent noch ein Sensor im Repo prüfen kann, weil sie Konto, Web-App
oder einen zerstörenden Versuch verlangen: die Anzeige der POM-Beschreibung auf
der Paketseite, die Usage- und Kontingent-Anzeige des Repositories, das Verhalten
bei einem zweiten Upload derselben Version. Der Slice liefert den Weg und den
realen Lauf; diese Eigenschaften bleiben als **benannte Grenze** ohne Beleg. Sie
gehören dem Betreiber, der sie in seiner Web-App einsehen oder im ersten realen
Fall beobachten kann; ein Slice-Auftrag daraus wäre ein Auftrag ohne Ende.

**Abgrenzung zu `BEO-PGC/github-actions-unverifizierbar-lokal`:** Dort fehlt der
reale Lauf, und §3.10 verlangt ihn vor der Closure (verkörpert); hier lief er,
und die Grenze liegt hinter ihm.

Deklaration: Planner-Agent (Closure-Rolle), 2026-09-26, auf Grundlage der
Risiko-Ausgänge von `slice-sdk-kotlin-cloudsmith` (Kontingent, Doppel-Upload) und
der Verifikation (§10 V-8).
