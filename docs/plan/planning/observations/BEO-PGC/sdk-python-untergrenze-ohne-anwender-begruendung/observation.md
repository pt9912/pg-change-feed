# BEO-PGC/sdk-python-untergrenze-ohne-anwender-begruendung

**Sub-Area:** `sdks/python/` (`*`/`PGC` Greenfield — das Python-Package
`pgchangefeed`).

Die Beobachtung: Das Package verlangt Python 3.14 oder neuer
(`requires-python = ">=3.14"` in `sdks/python/pgchangefeed/pyproject.toml`); die
README nennt die Untergrenze, aber keinen Grund. Ein Anwender mit Python 3.12 oder
3.13 kann das Package nicht installieren, ohne dass ein Dokument sagt, ob die
Untergrenze eine bewusste Entscheidung für Anwender ist oder aus der Wahl des
Basis-Images für den Bau folgt. Die einzige Stelle in `docs/plan/adr` und `spec`,
die 3.14 nennt, ist die Begründung des Basis-Images `python:3.14-slim` in
`ADR-0108` (`grep -rln 'requires-python\|Python 3\.14\|python:3\.14' docs/plan/adr spec`
druckt diese eine Datei, gemessen 2026-09-25); zur Untergrenze für Anwender sagt sie
nichts. Keine Anwender-Aktion in der README möglich: die Aussage „Python 3.14 or
newer“ ist wahr.

**Warum das zählt:** Die Untergrenze bestimmt die erreichbare Zielgruppe der
Bibliothek (`LH-FA-SST-009`); eine unbegründete Untergrenze ist eine ungeprüfte
Entscheidung.

Deklaration: `slice-sdk-readme-nutzerdoku`, Review F-12 (INFO), Plan §6 letzter Punkt.
