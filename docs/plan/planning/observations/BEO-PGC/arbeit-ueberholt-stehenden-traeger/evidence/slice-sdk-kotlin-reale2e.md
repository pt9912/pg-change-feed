# Beleg: slice-sdk-kotlin-reale2e

Vorgang: `slice-sdk-kotlin-reale2e` liefert die vier Realserver-Phasen
des Kotlin-SDK-Packages `pgchangefeed-kotlin` (gRPC, SSE, NATS-Vollinhalt,
HTTP) samt Kotlin-Abschnitt im Abdeckungs-Träger
`docs/user/sdk-e2e-abdeckung.md`. Der §3.13-Suchlauf des Plans (committetes
Feld) trug sechs Träger-Zeilen — alle sechs nach der Fixrunde gegen beide
Stände bestätigt (Verifikation §3).

Fund: die Prüf-Angaben des committeten Suchlauf-Felds selbst trugen zwei
Baum-Widerlegungen — beide vom Reviewer gefunden, beide in der Fixrunde
`c6523009` gezogen, beide vom Verifier gegen beide Stände bestätigt:

- (a) F-1 (MEDIUM): die Geprüft-Zeile nannte als Träger
  `sdks/kotlin/README.md` §Status — auf dieser Ebene existiert keine
  README (`ls sdks/kotlin/`), der reale Träger ist
  `sdks/kotlin/pgchangefeed-kotlin/README.md`; die Aussage der Zeile
  („trägt keine Teststrategie-Aussage über Realserver-Läufe") hält gegen
  die reale Datei, ihre Adresse war unauflösbar (Beleg-Adresse ohne
  Artefakt).
- (b) F-6 (INFO): die Behandlungs-Spalte zur C#-Zeile in
  `harness/README.md` behauptete „gezogen", während die Zeile an beiden
  Ständen byte-identisch bleibt (nur Zeilenverschiebung durch die
  Einfügung der neuen Target-Zeile) — die wahre Behandlung ist „gemeldet,
  nicht gezogen"; die Endklause „erweitert sich Slice für Slice um die
  Kotlin- und Python-HTTP-Abschnitte" bleibt als Verlaufs-Aussage wahr und
  ist zur Hälfte verbraucht.

Die bewegte Eigenschaft selbst („die Kotlin-Flächen tragen reale
Realserver-Belege") traf keinen falsch werdenden Träger: die `done/`-
Records der Kotlin-Flächen-Slices tragen ihre „kein eigener
Realserver-Test"-Aussagen sliced-/welle-scoped und bleiben als Records
wahr, `harness/sensors/docs-check.md` listet den Träger file-level und
wird nicht überholt, `spec/pflichtenheft.md` bleibt im Range unverändert
(breiterer Suchlauf des Verifikations-Laufs über beide Stände,
Verifikation §3). Die Fundstruktur ist dieselbe wie beim
dreiundzwanzigsten Beleg: der Suchlauf-Befund selbst wurde zum Träger,
dessen Prüf-Angaben der Review am Baum nachmisst.

Einordnung: dieselbe Klasse, eine neue Prüf-Angaben-Form — hier widerlegt
der Baum nicht nur behauptete vollzogene Nachzüge, sondern auch die
Adresse einer Geprüft-Zeile; beide Angaben sind Übergabe-Anteile des
Suchlaufs, ihre Auflösbarkeit ist die Prüf-Bedingung des Felds. Ausgang
bleibt verkörpert (`AGENTS.md` §3.13), kein neuer Schwellen-Übertritt —
die Schärfung ist eine Anwendungs-Schärfung der verkörperten Regel.
Geschärfte Lehre (Closure-Notiz §7 des Plans): eine Spiegelung überträgt
nicht die Assertions-Semantik, sondern die Draht-Aussage;
Bibliotheksspezifische Null-Formen (gson `JsonNull` vs. C# `null` vs.
Python `None`) brauchen je Sprache eine eigene Bindungs-Form.

Quelle: Review-Report `slice-sdk-kotlin-reale2e` (F-1 MEDIUM, F-6 INFO),
Fixrunde `c6523009`, Verifikations-Report desselben Slices (§1, §3, §6),
Commits `060aa62f`, `c6523009`, `e69a77eb`.