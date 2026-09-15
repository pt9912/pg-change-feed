# BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form von
Negativtests in Adapter-Paketen, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein **Negativtest** prüft, dass ein Aufruf **abgelehnt** wird —
aber er bindet die Ablehnung nicht an den **Eingabewert**, der sie auslösen
soll. Er stellt gegen einen Fake, der die Ablehnung **unabhängig von der
Abfrage** liefert. Der Test ist damit grün, egal was der Adapter mit der Eingabe
macht; die Zusage „dieser Wert wird abgelehnt" hat keinen Träger.

**Der Unterschied zu einem Fake, der nichts prüft**, ist fein und entscheidend:
dieser Fake prüft etwas (er liefert den Fehler), und der Test sieht vollständig
aus. Was fehlt, ist die **Kette**: *Eingabe → unveränderte Weitergabe →
Ablehnung*. Nur das mittlere Glied ist die Arbeit des Adapters — und genau es
war ungetestet.

Sichtbar wird das **nicht durch Lesen, sondern durch Mutieren der
Eingabeseite**: Wer nur die Ausgabeseite mutiert (den Fake, den Rückgabewert),
sieht weiter grün. Belegt an `slice-086`: der Test `?limit=0 → 400` stand gegen
einen Fake, der `ErrNonPositiveLimit` unabhängig von der Abfrage lieferte; die
Mutation `readChangesLimit("0") → nil` ließ die **ganze Suite grün**, obwohl der
Endpunkt danach bei `?limit=0` **still unbegrenzt** liest.

**Warum das zählt:** Der Negativtest ist die Prüf-Form einer Abgrenzung. Trägt
er nicht, ist die Abgrenzung nur behauptet — und sie fällt erst auf, wenn jemand
sie verletzt.
