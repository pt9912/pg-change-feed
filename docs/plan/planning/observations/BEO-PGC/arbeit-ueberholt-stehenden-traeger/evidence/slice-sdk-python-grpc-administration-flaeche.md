# Beleg: slice-sdk-python-grpc-administration-flaeche

Vorgang: `slice-sdk-python-grpc-administration-flaeche` — der
Python-Verifikations-Report
(`docs/reviews/verifikation-slice-sdk-python-grpc-administration-flaeche.md`,
V-1, MEDIUM) führt die Klasse an: der Träger liegt diesmal in der Datei,
die der Zug selbst berührt hat.

Fund: Der Modul-Docstring von
`sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py` (Zeilen
28–29) behauptete weiterhin „The stream is fire-and-forget: it has no
replay and **cannot be filtered by table**." — derselbe Docstring
dokumentiert elf Zeilen darüber die neuen `schema`/`table`-Filter, die
derselbe Zug an dieselbe Signatur gebaut hat (`stream_changes(timeout,
schema=None, table=None)`). Der Satz stand bereits am Parent
(`fd39b68b`) und wurde von dem Zug, der genau die beschriebene
Eigenschaft (Filterbarkeit des Streams) bewegt hat, nicht nachgezogen;
er wird ab diesem Zug falsch und widerspricht dem Lieferumfang im
selben Atemzug. Beide Review-Runden (Haupt-Review, Fixrunden-Review)
hatten die Stelle nicht gesehen; gefunden hat sie der Verifier gegen
Parent und Diff (`sed -n '28,29p'` an beiden Ständen,
`git diff fd39b68b..HEAD` berührt den Absatz nicht).

Gezogen: einzeilige Berichtigung des Satzes auf „no replay of
already-delivered changes" durch `8b12198e`; der frische pytest-Lauf im
Frischbau meldet 130 passed — der Träger ist gezogen, kein Rest-Vorhaben.

Einordnung: dieselbe Klasse wie die Gründungsfälle (`slice-091`/`-093`/
`-094`) — die Arbeit war korrekt, der Träger steht nicht im Diff, kein
Sensor liest ihn; hier zusätzlich verschärft dadurch, dass der Träger
in einer Datei lag, die der Zug selbst anfasste, und dass ihn zwei
Review-Runden übersprachen. Die verfügbare Falsifikation bleibt die
Messung an beiden Ständen; der Verifier hat sie gefahren.
