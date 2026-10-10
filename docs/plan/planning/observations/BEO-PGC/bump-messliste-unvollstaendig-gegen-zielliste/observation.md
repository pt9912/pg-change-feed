# BEO-PGC/bump-messliste-unvollstaendig-gegen-zielliste

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Liste der Ziele, die ein
Pin-Bump alt gegen neu misst, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Der Plan eines Pin-Bumps zählt die Einzel-Ziele auf, die er alt gegen neu misst,
und leitet die Liste aus dem Changelog und aus den Zielen ab, die ihm einfallen, nicht aus der
Zielliste der Datei, die den Pin trägt. Ziele, die der Pin ebenfalls bewegt, bleiben ungemessen,
obwohl der Plan „dazu die Einzel-Ziele“ sagt, als wäre die Aufzählung erschöpfend.

**Erster Vorgang:** `slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1` (Review F-4, INFO). Liefer-Punkt
1 listete `doc-immutable`, `doc-planning`, `doc-targets` und die Gates; `doc-trace`, `doc-complete`
(laut `d-check.mk` ein Vollständigkeits-Gate), `doc-doctor` und `doc-repair` fehlten. Der Reviewer
maß alle vier alt gegen neu: byte-gleich, je Exit 0. Kein Schaden; die Liste war nicht gegen
`d-check.mk` abgeleitet.

**Abgrenzung:** `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` ist die Regel, dass eine
Erweiterung ihr Artefakt braucht; hier fehlt nicht das Urteil, sondern ein Gegenstand der Messung.
