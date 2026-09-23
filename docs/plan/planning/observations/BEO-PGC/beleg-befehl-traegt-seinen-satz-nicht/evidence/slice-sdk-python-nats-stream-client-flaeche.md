# Beleg: slice-sdk-python-nats-stream-client-flaeche

Vorgang: `slice-sdk-python-nats-stream-client-flaeche` liefert die
NATS-Vollinhalts-Client-Fläche (`SPEC-024`) samt Realserver-Phase im
Runner (`make test-sdk-python-integration`).

Fund, zwei Formen derselben Klasse in derselben Kette:

- **Haupt-Review F-6 (MEDIUM) — die Beleg-Assertion trägt ihren Satz
  nicht:** der Realserver-Rejection-Test druckte den Marker „REJECTED
  token-rejected: …" (vom Runner als `reject_marker` gegenprüft), aber
  die Assertion dahinter hielt `pytest.raises(Exception)` — ein
  nicht-authentifizierungsbedingter Verbindungsfehler hätte denselben
  Marker und denselben Runner-Satz erzeugt; die Ursachen-Aussage des
  Belegs trug seine Messung nicht. Lösung Fixrunde 1: Bindung auf
  `nats.errors.Error` + `match="Authorization Violation"` — gegen die
  nats-py-Quellen (v2.11.0 und master 2.16.0,
  `_process_connect_init`-Pfad) nachgemessen, wo der
  Server-Rohwortlaut real im String bleibt; die
  `AuthorizationError`-Rezeptur der Haupt-Review wäre die falsche Bindung
  gewesen (sie fliegt nur im Post-Connect-`_process_err`-Pfad). Der
  s3e→s3f-Wechsel ist der lebende Mutations-Beweis: die
  Bindungsänderung kippte den Lauf real rot→grün gegen denselben Server.
- **Fixrunden-Re-Review R-1 (HIGH):** das committete §3.13-Suchlauf-Feld
  des Plans behauptete zwei vollzogene Nachzüge, die der Baum widerlegte
  (`README.de.md:27` unverändert, `harness/README.md`-Reject-Satz
  unverändert) — der Beleg trägt Sätze, die die genannten Fundstellen
  nicht stützen. Der Vorgang zählt zugleich bei
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (23. Beleg, Träger-Kette).
  Lösung Fixrunde 2: die Reststellen real gezogen, die Feld-Zeilen auf
  den realen Ist-Stand gestellt.

Quelle: Haupt-Review `slice-sdk-python-nats-stream-client-flaeche` (F-6),
Fixrunden-Report desselben Slices (R-1 + Negativbefund F-6),
Verifikationsbericht desselben Slices (§2.3), Commits `0f8cc4f2`,
`b4d1d352`, `a3e9d64e`.