# Beleg: slice-084

Vorgang: `slice-084` — die Naht in `postgresack`.

Fund: Der Slice-Plan §3(b) nennt als Beleg für den Satz „die Änderung ist auf
**ein** Paket isoliert … `git diff --name-only` listet ausschließlich Dateien
unter `internal/adapters/driven/postgresack/`" den Befehl
`git diff --name-only fb6adf6..4035ee7` — **ohne Pathspec**. Auf dem sauberen
Baum listet dieser Befehl **fünf** Pfade: den Plan, `harness/image-hash.txt` und
die drei Paketdateien. Der Satz ist **wahr** (die drei Paketdateien sind die
Änderung), seine **Stütze** trägt ihn nicht — die **Range** war nachgetragen,
der **Pathspec** fehlte.

**Der erste Korrekturversuch hat es nicht behoben:** `a3cb2c4` ergänzte die
Range, ließ aber den Pathspec weg — dieselbe Klasse, eine Runde später. Erst
`-- internal/ ':!internal/adapters/driven/postgresack/'` (die Gegenrichtung,
leer) trägt die Aussage.

Gefunden hat es der Verifier im Verifikationsbericht zu `slice-084` als
**V-1** — ausdrücklich als Belegform mit Rest, nicht als Substanzfehler.

Quelle: Verifikationsbericht zu `slice-084` (V-1) ·
`docs/plan/planning/done/altbestand/slice-084-postgresack-naht.md` §3(b) (berichtigt in
`4387523`, dem Closure-Commit des Vorgangs) ·
`git diff --name-only fb6adf6..4035ee7`.
