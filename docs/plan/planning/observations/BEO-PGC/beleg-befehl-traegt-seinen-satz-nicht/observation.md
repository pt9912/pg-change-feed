# BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form, in der
ein Beleg-Befehl seine Aussage trägt, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Träger **nennt einen Beleg** — einen Befehl, eine Abfrage,
einen Pfad — als Stütze einer Aussage, und der genannte Beleg **trägt diese
Aussage nicht**: Er misst etwas anderes, zählt etwas anderes oder liefert ein
anderes Ergebnis als das behauptete. Die Aussage kann dabei **wahr** sein; falsch
ist die **Stütze**. Das macht den Fall schwer zu sehen: Wer den Satz liest,
findet ihn plausibel; wer den Befehl ausführt, findet etwas anderes — und nur
der zweite Blick deckt es auf.

Belegt an zwei abgeschlossenen Vorgängen — **zwei verschiedene Befehle**:

- **`slice-084`** (Verifikation `verify-slice-084` V-1): der Slice-Plan §3(b)
  nennt als Beleg für „die Änderung ist auf **ein** Paket isoliert" den Befehl
  `git diff --name-only fb6adf6..4035ee7` — **ohne Pathspec**. Auf dem sauberen
  Baum listet er **fünf** Pfade (Plan, `image-hash.txt` und die drei
  Paketdateien), während der Satz „ausschließlich Dateien unter
  `postgresack/`" behauptet. Die **Range** war nachgetragen, der **Pathspec**
  fehlte; die Ergänzung
  `git diff --name-only fb6adf6..4035ee7 -- internal/ ':!internal/adapters/driven/postgresack/'`
  trägt die Aussage (die Gegenrichtung ist leer).
- **`slice-085`** (Verifikation `verify-slice-085` V-3):
  `harness/sensors/coverage-gate.md` §Grenze Punkt 1 zitiert
  `go list -f '{{len .TestGoFiles}}'` als Beleg für „fünf Pakete ohne
  Testdatei". Mit dieser Formel liefert der Lauf **25** — die externen
  Testpakete (`XTestGoFiles`) zählt sie nicht mit; erst beide Felder ergeben
  genau die fünf genannten.

**Ursprung und Vorkommen sind zwei Dinge.** Der zweite Beleg ist
**Alt-Bestand**: Den Satz eingeführt hat `slice-079` (`65aead2`) — **gefunden**
wurde er erst in `verify-slice-085`. `slice-079` ist damit der **Ursprung**, kein
Beleg-Vorgang; der Zähler zählt die Vorgänge, in denen der Fund **auftrat**.

**Warum das zählt:** Ein Beleg ist die Prüf-Form einer Aussage. Trägt er sie
nicht, prüft ein späterer Leser gegen etwas anderes als das Behauptete — und die
Aussage selbst bleibt dabei unangetastet, sodass niemand einen Anlass sieht,
nachzusehen. Die verwandte Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
trifft den **Wert**; hier trifft es den **Befehl**, der ihn stützen soll.
