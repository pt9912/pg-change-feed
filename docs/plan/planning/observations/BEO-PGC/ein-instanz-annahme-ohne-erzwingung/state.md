Zustand: offen — Ausgang: **weiter offen**, adressiert. Die Annahme steht im Godoc von
`NewBackfillAdmission` und im Plan von `slice-backfill-run-store` (§6); Adresse ist der
Re-Evaluierungs-Trigger in `ADR-0113` §Re-Evaluierungs-Trigger („Mehr als eine Instanz
je Quelle nimmt Anträge an"). Der sequenzielle Fall trägt der DoD „Verarbeitung" in
`slice-backfill-sql-administration` (zwei aufeinanderfolgende `backfill`-Anträge derselben
Tabelle, der zweite endet `failed`; belegt im Login-Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins`,
Verifikation §2 Zeile 2); der nebenläufige Fall bleibt ungetestet, eine Sperre
oder eine Unique-Kante auf aktive Runs ist Sache der Re-Evaluierung, nicht dieses Eintrags.
Der zweite Beleg (`slice-backfill-sql-administration`, F-10) betrifft die andere Hälfte der
Annahme, „je Quelle nimmt eine Instanz Anträge an": für die Antragsart `backfill` bindet seit
der Fixrunde der Quellvergleich in `processAdministrationRequests` sie am Code (Test mit zwei
Quellen, Mutation rot); für die vier übrigen Antragsarten liest die Queue weiter ohne
Quellbezug (Bestand aus `ADR-0050`), und zwei Instanzen derselben Quelle bleiben ohne Schutz.
Zähler (abgeleitet, gemessen mit `ls evidence | wc -l`): **3×**
(evidence/slice-backfill-run-store.md,
evidence/slice-backfill-sql-administration.md,
evidence/slice-start-vorlauf-grenze.md). Der dritte Beleg
(`slice-start-vorlauf-grenze`) trifft dieselbe Annahme aus einer neuen
Richtung: die Verschiebung von `START_REPLICATION` nach `Stream.Run`
(`ADR-0128`) verschiebt den Zeitpunkt, an dem eine zweite Instanz derselben
Quelle am Slot scheitert, nach hinten (Vorlauf plus Goroutine-Start statt
sofort beim Verbindungsaufbau); am realen Container beim Debuggen jenes
Slice trat einmalig ein verwandter, aber nicht identischer Fall auf
(„context canceled“ statt SQLSTATE 55006 bei zwei überlappenden
Vorläufen) — kein committeter Test.

**Architect-Verdikt (Lese-Schritt der Closure von `welle-transformationen`, frischer
Kontext): Trigger NICHT ausgelöst — Ausgang bleibt `weiter offen`.**

Geprüft: alle drei Belege (`evidence/slice-backfill-run-store.md`,
`evidence/slice-backfill-sql-administration.md`, `evidence/slice-start-vorlauf-grenze.md`),
`ADR-0113` §Re-Evaluierungs-Trigger im Original, `ADR-0128` §Kontext/§Gemessen im Original.

Der Trigger-Wortlaut („Mehr als eine Instanz je Quelle nimmt Anträge an, oder ein zweiter
Worker kommt (Parallelisierung, `ADR-0111`-Ausbaustufe)") beschreibt eine **beabsichtigte
Betriebsänderung**: ein Betreiber schaltet dauerhaft mehr als eine annehmende Instanz je Quelle
scharf, oder ein zweiter Worker wird eingeführt — beides eine Skalierungs-Entscheidung mit Folgen
für Aufnahme und Annahme-Prüfung, die neu zu entscheiden sind (DB-seitige Beanspruchung, Kanal).
Der dritte Beleg trägt das nicht:

1. **Anderer Mechanismus.** Der dritte Beleg betrifft nicht die Annahme eines
   `backfill`-Antrags (der Gegenstand dieses Eintrags — keine Sperre, keine Unique-Kante auf
   `cdc.backfill_run`), sondern den Start des Replikationsstroms selbst. Dort erzwingt PostgreSQL
   die Ein-Instanz-Eigenschaft bereits **strukturell**: ein logischer Replication Slot lässt sich
   nicht von zwei Sitzungen gleichzeitig aktiv halten, `START_REPLICATION` der zweiten scheitert
   an SQLSTATE 55006. `ADR-0128` verschiebt nur den **Zeitpunkt**, an dem dieser bereits
   bestehende Wächter greift (nach dem Vorlauf statt beim Verbindungsaufbau) — er entfernt ihn
   nicht und schafft keine neue Lücke, in der zwei Instanzen dauerhaft nebeneinander Anträge
   annehmen könnten. Der Gegenstand dieses Eintrags (`NewBackfillAdmission`, `cdc.backfill_run`
   ohne Unique-Kante) bleibt von `ADR-0128` unberührt.
2. **Anderes Ereignis.** Der beobachtete reale Fall ist ein `docker restart` **während** eines
   noch laufenden Vorlaufs — ein transientes Überlappungsfenster beim Prozess-Neustart derselben
   Quelle, kein dauerhafter Zwei-Instanzen-Betrieb. Der Ausgang war zudem unauffällig: die
   auslaufende Instanz brach ihren Antrag mit „context canceled" ab (ihr Prozess-Kontext wurde
   durch den Neustart beendet), *bevor* sie mit der neuen Instanz um denselben Slot konkurrierte;
   die neue Instanz nahm den Antrag danach erneut auf. Das ist ein geordneter, sich selbst
   auflösender Übergang — kein Beleg für zwei gleichzeitig aktiv annehmende Instanzen.
3. **Schwächere Beleg-Form.** Der dritte Beleg selbst benennt seine Grenze: „kein committeter
   Test, kein Teil der Belege der DoD" — eine einmalige Beobachtung beim Debuggen, keine
   reproduzierte oder dauerhaft geführte Zusicherung (`AGENTS.md` §3.12 Instanz B: ein Befund ohne
   Beleg-Anker bleibt eine Beobachtung, keine geprüfte Tatsache).

**Verdikt:** Der dritte Beleg ist eine benannte, verwandte Beobachtung zum selben Ober-Thema
(„was passiert, wenn die Ein-Instanz-Annahme kurzzeitig nicht hält") — er löst aber nicht den in
`ADR-0113` benannten Trigger aus, weil er weder die Annahme-Transaktion des Backfills betrifft
noch einen dauerhaften Mehr-Instanzen-Betrieb zeigt, sondern eine transiente, durch PostgreSQLs
eigene Slot-Exklusivität bereits gedeckte Neustart-Überlappung, einmalig beobachtet, ungetestet.
Eine neue oder geänderte ADR ist damit **nicht** angezeigt; dieser Eintrag bleibt **weiter offen**
mit unverändertem Adressaten (`ADR-0113` §Re-Evaluierungs-Trigger) für den ursprünglichen
Gegenstand (Backfill-Annahme-Transaktion). Der dritte Beleg bleibt im Register stehen, weil er
real beobachtet und potenziell für einen künftigen, eigenständigen Eintrag zur
Vorlauf-Überlappung relevant ist — er zählt aber nicht als drittes Auftreten *dieses* Triggers im
engen Sinn des Wortlauts; da der Zähler bereits vor diesem Zug 3× stand und der Lese-Schritt genau
diese Frage klären sollte, bleibt die Zahl unverändert und der Eintrag unter Beobachtung, ohne
dass ein neuer Steering-Loop-Zyklus (erneutes 3×) nötig wäre, um ihn erneut vorzulegen — ein
vierter, hinreichend verschiedener Beleg (z. B. ein committeter Test mit echtem
Zwei-Instanzen-Wettlauf um dieselbe `queued`-Zeile) würde den Trigger dagegen auslösen.
