# BEO-PGC/nats-notify-validierung-koennte-bestand-ablehnen

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Aktivierungs-/Notify-Kette im gesamten Repo, keine eigene Sub-Area im Sinn
der Modus-Deklaration).

Die Beobachtung: `natsnotify.Notify` (`ADR-0056`) lehnt Schema-/
Tabellennamen mit NATS-reservierten Zeichen (`.`, `*`, `>`) oder
Whitespace vor dem ersten Publish-Versuch ab (Fehlerklasse `transient`).
Diese Prüfung existiert erst seit `slice-058` und trifft ausschließlich
den Notify-Pfad — eine heute (ohne `CDC_NATS_URL`) unauffällig
funktionierende Aktivierung mit einem ungewöhnlichen Schema- oder
Tabellennamen (z. B. ein per Anführungszeichen quotiertes
PostgreSQL-Bezeichner mit Punkt) würde sichtbar erst dann auffällig,
wenn der Betreiber `CDC_NATS_URL` erstmals aktiviert.

## Benannt, nicht gezählt

Kein Vorkommen ohne abgeschlossenen Vorgang bislang.
