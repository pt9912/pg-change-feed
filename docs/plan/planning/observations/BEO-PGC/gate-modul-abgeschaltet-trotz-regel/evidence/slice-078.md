# Beleg: slice-078

Vorgang: `slice-078` — die Entfernung der host-lokalen Pfade und die Aktivierung
des `hostpaths`-Moduls.

Fund: Die Regel „keine host-lokalen absoluten Pfade in diesem Repo" bestand als
Absicht — sie stand in Entscheidungen als Vorbild-Zitat und war Gegenstand eines
Architect-Verdikts —, aber **ihr Sensor war abgeschaltet**: `hostpaths` fehlte in
der `modules`-Liste der `.d-check.yml`. Folge: **42 Vorkommen in 15 Dateien**
sammelten sich über viele Slices hinweg, ohne dass etwas sie hätte melden können;
sichtbar wurden sie erst, als der Auftraggeber das Modul einforderte. Belegt mit
dem Modul-Aufruf über den Basisstand (31 gescannte Fundstellen) und einem
repo-weiten Grep (42 Vorkommen in 15 Dateien).

Quelle: `docs/reviews/architect-verdict-slice-078-konfliktpfad.md` ·
`docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md` ·
`docs/reviews/verify-slice-078.md` (eigene Zählung, 31/42).
