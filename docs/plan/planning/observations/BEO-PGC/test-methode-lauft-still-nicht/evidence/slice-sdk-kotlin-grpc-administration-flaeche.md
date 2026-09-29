# Beleg: slice-sdk-kotlin-grpc-administration-flaeche

Vorgang: `slice-sdk-kotlin-grpc-administration-flaeche` — Erstauftreten
der Klasse, befundet als V-1 (MEDIUM) im
`docs/reviews/verifikation-slice-sdk-kotlin-grpc-administration-flaeche.md`.

Fund: Die `@Test`-Methode
`PgChangeFeedAdministrationClientTableTest.disableTable table missing at
source throws notFound` (Zeilen 76–89) lief nie — still. Ihr
`runBlocking`-Block trug als letzten Ausdruck
`assertFailsWith<PgChangeFeedGrpcNotFoundException> { … }`, damit trug
die Methode den Rückgabetyp `PgChangeFeedGrpcNotFoundException` statt
`void`, und die JUnit-Plattform entdeckt Testmethoden mit Nicht-void-
Rückgabe nicht. Die Testresultate zeigten für die Klasse `tests="6"`
(sieben `@Test`-Methoden, `skipped="0"`), ein `--tests`-Filter auf
genau diese Methode brach `:test` mit „no tests found" ab, dieselbe
Probe auf die void-Methode `enableTable table missing at source throws
notFound` lief grün (`javap` am kompilierten Klassenfile als
Signatur-Beleg). V-3 des Verifikations-Reports systemprüfte alle
Testklassen des Packages auf dieselbe Form: V-1 blieb der einzige Fall.

Wirkung: der `NOT_FOUND`-Negative-Fall von `disableTable` war ohne
laufenden Test; der Review-Nachweis F-2 („zwei notFound-Tests …
enableTable, disableTable") trug seine Hälfte nicht. Kein DoD-Bruch:
alle elf RPCs behielten je mindestens einen laufenden Happy-Pfad-Test,
und die `NOT_FOUND`-Zuordnung war zweifach anders belegt
(`ErrorMappingTest` generisch, `enableTable`-notFound).

Behoben: Ein-Zeilen-Korrektur in `3b381f5b` — Ergebnis verwerfen
(`val ex = assertFailsWith<…> { … }`) samt Statuscode-Assertion, Form
des `enableTable`-Gegenstücks; real verifiziert im frischen Docker-Bau:
Testresultat der Klasse `tests="7"`, `failures="0"`.

Warum hier ein eigenes Verzeichnis und kein Beleg bei
`BEO-PGC/test-runner-stiller-ausschluss`: dessen unveränderliche
Beobachtung definiert den Mechanismus über das eigene Runner-Skript
(`-run`-Muster); hier liegt die Ursache in der Entdeckungs-Regel des
Test-Frameworks gegen eine Signatur-Form — dieselbe Wirkung, anderer
Träger der Ursache (Abgrenzung in der `observation.md`).
