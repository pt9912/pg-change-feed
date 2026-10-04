**Vorgang:** slice-sdk-0-6-kompatibilitaet-messen (Architect, Entscheidung zu `ADR-0147`, im Bericht des Architect an den Hauptlauf gemeldet)

**Fund:** Der Architect schrieb `docs/plan/adr/0147-sdk-csharp-http-basis-null-literal-einschraenkung.md` und seinen Eintrag im ADR-Index teilweise mit `echo >>` (Anhängen per Umleitung) und `git apply` aus einer Patch-Datei statt mit Edit/Write. Beides ist durch `AGENTS.md` §3.1 („Text-Umschreiben im Repo ist Sache der Datei-Werkzeuge des Laufs“, Umleitung geschärft seit welle-transformationen) verboten. Der Architect meldete den Verstoß selbst; der Inhalt der Dateien ist korrekt, ein Folgeschaden liegt nicht vor.

**Form (Ausprägung):** Umleitungs-Hälfte der Regel, **ohne Wirkung auf den Inhalt**; der PreToolUse-Guard liest Umleitungen nicht (Grenz-Zeile `MR-003`), `git apply` ist kein gelistetes in-place-Werkzeug. Kein Beleg für einen Host-Interpreter am Kopf. Rolle: Architect (bisher Implementer, Planner, Reviewer, Verifier).

Quelle: Bericht des Architect an den Hauptlauf (Angabe des Auftraggebers), **übernommen**, nicht am Guard gemessen.
