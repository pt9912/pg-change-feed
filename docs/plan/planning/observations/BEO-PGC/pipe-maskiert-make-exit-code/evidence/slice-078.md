# Beleg: slice-078

Vorgang: `slice-078` — der Push des Delta-Review-Reports durch den Planner.

Fund: `make gates` und `git push` standen in **einem** Shell-Aufruf. Der
Gate-Lauf endete rot (Exit 2, direkt gelesen), der Push lief trotzdem — er war
nicht auf den Exit-Code konditioniert. Damit lag kurzzeitig ein roter Stand auf
dem Hauptzweig; die Ursache war ein `id-unlinked` im Report, das der Planner in
einem Folgeschritt reparierte (`f86fa99`).

Klasse dieses Vorkommens: die **geschärfte** §3.9-Hälfte („Prüfung und
Folgehandlung sind zwei Schritte, nicht einer") — nicht das Maskieren durch eine
Pipe, sondern die fehlende Konditionierung über einen Werkzeug-Aufruf hinweg.
Die Regel ist längst verkörpert; dieser Beleg zeigt, dass ihr Wächter
**Disziplin** bleibt: kein Sensor fängt eine unterlassene Konditionierung.

Quelle: Sitzungsverlauf · `AGENTS.md` §3.9 (geschärfte Fassung) · Commit
`f86fa99`.
