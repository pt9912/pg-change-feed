# BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Verträglichkeit
der SDK-Lesemodelle mit einem additiven Feld der `GET /changes`-Antwort, keine
eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Die Antwort von `GET /changes` trägt mit
`slice-backfill-change-origin` ein dreizehntes Feld `origin`. Dass die
Decoder der SDK-Lesemodelle (C# `System.Text.Json`, Kotlin Gson, Python `json`
mit expliziten Feldern) ein unbekanntes Feld ignorieren, ist **gelesen**, nicht
gegen den realen Server gemessen. Gemessen im Slice (`git grep` gegen `HEAD`):
strikte Dekoder-Muster im Baum 0 Treffer; der Wegwerf-Client
`tools/harness/httpclient` dekodiert mit `json.Unmarshal` in eine Struktur ohne
`origin` und ist nicht strikt. Ein realer Beleg mit den SDK-Decodern gegen einen
Server, der das Feld liefert, fehlt.

**Warum das zählt:** Ein additives Feld ist nur dann abwärtsverträglich, wenn jeder
Leser es ignoriert; ein einziger strikter Decoder bricht am Feld. Kein Gate im
Bestand liest die SDK-Decoder gegen eine Antwort mit dem neuen Feld.

Deklaration: `slice-backfill-change-origin`, Risiko §6 dritter Punkt, Ausgang
*weiter offen*.
