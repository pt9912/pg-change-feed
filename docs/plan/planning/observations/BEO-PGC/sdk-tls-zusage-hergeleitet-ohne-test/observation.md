# BEO-PGC/sdk-tls-zusage-hergeleitet-ohne-test

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die SDK-Packages unter
`sdks/` und das Messwerkzeug unter `tools/harness`; keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Die TLS-Option der drei SDKs (Anker-Datei) trägt Zusagen in den READMEs, im
Handbuch und im Plan, die am Ende von `slice-sdk-tls-optionen` aus dem Code **hergeleitet**,
aber nicht von einem Test gebunden sind. Benannt, je mit Grund:

1. **Systemanker-Ausschluss** („mit Anker gelten die Anker des Betriebssystems nicht
   zusätzlich“): nur Python HTTP und SSE gebunden (Eingabeseite `SSL_CERT_FILE`, Mutation
   `load_default_certs()` rot). C#, Kotlin und Python gRPC nicht: der Wurzelspeicher ist dort
   prozessweit initialisiert, ein Test im selben Prozess wäre ein Scheintest oder bräuchte einen
   Kindprozess je Fall.
2. **Übergebener Client behält seine Einstellung** (Verifikation V-1, LOW): ein selbst
   übergebener `HttpClient`, Kanal oder `httpx.Client` bekommt den Anker nicht; hergeleitet aus
   der Konstruktorform, kein Test in keinem SDK. Verifizierbar: ein Fall je SDK.
3. **Ausstellerfälle nur auf Unit-Ebene** (V-3, INFO): „Anker ist das Zertifikat oder sein
   Aussteller“ ist in allen drei Sprachen im Unit-Test gebunden, am Realserver nicht
   (`tools/harness/certgen` erzeugt nur selbstsignierte Zertifikate).
4. **Anker plus nicht-TLS-Adresse wirft** (V-2, INFO): die Festlegung steht in keiner Norm; die
   drei SDKs setzen sie gleich um (Kotlin: jedes Nicht-`https`; Python gRPC: schemaloses
   `host:port` mit Anker wird TLS) und dokumentieren sie. Eine benannte Spec-Lücke, die der
   Architect schließt oder als Umsetzungsentscheidung bestätigt.
5. **Änderung des C#-Mappings** (`Internal` mit `HttpRequestException` wird
   `UnexpectedStatus`) ist von `make test-sdk-kompat` nicht gedeckt: das Ziel liest die
   Fehlertypen, nicht das Mapping; belegt nur durch den Unit-Fall `ExpectConnectionFailureAsync`.
6. **README-Überschrift `0.7.0`** (C# und Kotlin, Abschnitt Upgrading) gegen Package-Version
   `0.6.1`: beschreibt den Inhalt der nächsten Minor-Version; die Auflösung ist der Release-Zug
   (Versionen anheben vor jeder Veröffentlichung, braucht eine neue Freigabe des
   Auftraggebers).

**Warum das zählt:** Eine Zusage ohne Test ist die Klasse, die `AGENTS.md` §3.12 als
*hergeleitet* zu kennzeichnen verlangt; bleibt sie in README und Handbuch stehen, liest sie ein
Betreiber als geprüft.

Gelesen wird der Eintrag bei der Planung eines Folge-Slice an den SDK-Packages oder am
Realserver-Fixture; kein Folge-Slice ist zugewiesen.
