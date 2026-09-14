# BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Betreiber-Dokumentation gegenüber der gewachsenen Oberfläche, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Das Benutzerhandbuch (`docs/user/benutzerhandbuch.md`)
beschreibt die Betreiber-Oberfläche des Werkzeugs — Umgebungsvariablen des
Feed-Containers (§5), administrative Aufgaben (§4). Wächst diese Oberfläche
durch einen Slice, ohne dass das Handbuch im selben Zug mitgezogen wird,
entsteht eine Lücke, die **kein Sensor** fängt: `docs-check` prüft
Referenzen, nicht Vollständigkeit; die verkörperte
Handbuch-Versionshistorie-Regel (`BEO-PGC/handbuch-versionshistorie-uebersprungen`)
greift erst, wenn das Handbuch **angefasst** wird — wer es gar nicht anfasst,
löst sie nicht aus.

Drei unabhängige Vorgänge, in denen eine neue Betreiber-Oberfläche entstand
und das Handbuch unberührt blieb: `slice-059` (HTTP/JSON-API: `CDC_HTTP_ADDR`,
`CDC_API_TOKEN_READER`/`ADMIN`), `slice-066` (Spaltenausschluss:
`cdc.exclude_column`/`cdc.include_column`), `slice-069` (gRPC-Live-Stream:
`CDC_GRPC_ADDR`). Beim dritten Vorgang war der Aufschub ausdrücklich benannt
(„Operator-Doku gehört zur Erreichbarkeit in `slice-071`/`slice-072`") —
das benennt die Lücke, beseitigt sie aber nicht.

Deklaration: Planner-Koordinator, 2026-09-14, auf Nachfrage des
Auftraggebers („Ist das Handbuch noch aktuell?").
