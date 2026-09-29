# BEO-PGC/test-methode-lauft-still-nicht

**Sub-Area:** `sdks/kotlin/` (Testbestand des Kotlin-SDK-Packages;
Sub-Area-Kürzel `PGC` aus der Modus-Deklaration, Repo-Default).

Die Beobachtung: Eine Test-Methode existiert im Quelltext, kompiliert und
zählt im Testbestand — wird aber vom Test-Framework **nie ausgeführt**,
weil ihre Signatur-Form die Entdeckungs-Regel der Plattform nicht erfüllt.
Der Lauf bleibt grün und meldet die kleinere Testfall-Zahl nicht als
Verlust: die Testresultat-XML der Klasse trug `tests="6"`, während die
Klasse sieben `@Test`-Methoden deklariert — nichts schlägt an, kein Gate
liest die Differenz.

Konkret (Kotlin/JUnit-Plattform): eine `@Test`-Methode, deren
`runBlocking`-Block als letzter Ausdruck `assertFailsWith<…> { … }`
stehen lässt, trägt den Rückgabetyp `PgChangeFeedGrpcNotFoundException`
statt `void`/`Unit` — die JUnit-Plattform entdeckt Testmethoden mit
Nicht-void-Rückgabe nicht; ein `--tests`-Filter auf genau diese Methode
bricht `:test` mit „no tests found" ab, dieselbe Probe auf eine
void-Methode derselben Klasse läuft grün.

**Abgrenzung zum Nachbar-Eintrag:** `BEO-PGC/test-runner-stiller-ausschluss`
trifft dieselbe Wirkung (geschriebener Test läuft stillschweigend nie)
über einen anderen Mechanismus — dort filtert das **eigene
Runner-Skript** (`tools/harness/run-integration-tests.sh`, `-run`-Muster)
die Testfunktion aus; hier entdeckt die **Plattform-Regel des
Test-Frameworks** die Methode wegen ihrer Signatur-Form nicht. Beide
tragen dieselbe Kernfrage („ist jeder geschriebene Test auch gelaufen?"),
aber der Träger der Ursache ist ein anderer: Skript-Wartung hier,
Signatur-Form dort. Dieser Eintrag bleibt bewusst eigenständig, statt
den älteren Eintrag über seine unveränderliche, skriptbezogene
Definition hinaus zu dehnen.

**Warum das zählt:** Der Negative-Fall (`disableTable` → `notFound`)
war testseitig geschrieben und schriftlich als Abdeckung beansprucht,
lief aber nie — die beanspruchte Abdeckung trug einen Fall weniger als
behauptet, und der grüne Lauf hat die Lücke verdeckt statt sie zu
zeigen. Die verfügbare Falsifikation ist die Zählung gegen den
Quelltext (Testfall-Zahl des Ergebnisses gegen deklarierte
`@Test`-Methoden), nicht ein roter Lauf: kein Ergebnis dieser Art wird
rot, solange niemand zählt.
