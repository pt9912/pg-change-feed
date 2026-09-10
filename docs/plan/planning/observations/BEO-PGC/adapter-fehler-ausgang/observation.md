# BEO-PGC/adapter-fehler-ausgang

**Sub-Area:** Bootstrap-Verdrahtung (Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: der erste Production-Pfad endet auf jeden Adapter-Fehler
mit Prozess-Ausgang 1 — die `SPEC-008`-Aktion für `transient` („Erneut
versuchen mit begrenztem Backoff") trägt kein Element des Pfads, und die
Adapter-Grenze „kontrollierte Fortsetzung beim Aufrufer"
([`ADR-0012`](../../../../docs/plan/adr/README.md)) wird vom Aufrufer mit
Prozess-Ende beantwortet.
