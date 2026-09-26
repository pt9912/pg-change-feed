# Slice transformationen-betriebsdoku: Betriebsdokumentation und SDK-Beleg — Transformationsregeln im Benutzerhandbuch, Row-Image-Schlüssel in den SDK-Modellen als undurchsichtig belegt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md),
[`LH-FA-ADM-001`](../../../../spec/lastenheft.md),
[`LH-FA-SST-009`](../../../../spec/lastenheft.md) (die drei SDK-Packages),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 6 (SDK-/Doku-Beleg) und Teilfrage 8 (keine Änderung des
Nachrichtenschemas).

**Berührte Spec-Stellen:** [`SPEC-026`](../../../../spec/pflichtenheft.md),
[`SPEC-027`](../../../../spec/pflichtenheft.md),
[`SPEC-028`](../../../../spec/pflichtenheft.md) (die drei SDK-Packages),
[`SPEC-018`](../../../../spec/pflichtenheft.md),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md) (Drahtverträge) — gelesen,
nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Zwei Belege, beide erst jetzt wahr formulierbar: (1) die drei
SDK-Packages setzen keine bestimmte Menge von Row-Image-Schlüsseln voraus — zu
belegen, nicht als geprüft behauptet (§3.12 Instanz B), erwartet wird keine
Code-Änderung
([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 8); (2) das Benutzerhandbuch beschreibt die Betreiber-Oberfläche der
Transformationsregeln, nachdem `e2e-wirkung` und `e2e-abhilfe` die Wirkung
belegt haben. Das ist der **aufgeschobene Gegenstand** von `antragsweg-schema`
(Adresse, siehe dessen §1) und die Adresse aller Slices, die Handbuch-Anteile
an dieses Slice gemeldet haben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine SDK-Code-Änderung oder ein SDK-Release** — erwartet keine; findet der
  Suchlauf einen Zugriff auf feste Schlüssel im Produktivcode eines SDK, ist
  das ein Plan-Nachzug (Rückführung §4) samt Versionshebung, kein stiller
  Zusatz. Ein Tag-Push bleibt Betreiber-Handlung außerhalb der Welle.
- **Ein SDK-Realserver-E2E mit aktiver Regel** — die drei SDK-Tiers fahren
  gegen einen Server ohne Regeln; ihre Test-Daten tragen eigene Tabellen und
  Schlüssel (`name` als Sentinel), keine SDK-Voraussetzung.
- **Aussagen, die kein Beleg trägt** — das Handbuch nennt nur, was
  `e2e-wirkung` und `e2e-abhilfe` real belegt haben; die Aussage „die Abhilfe
  wirkt“ steht erst mit dem Beleg.
- **Routing** ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) — keine
  Handbuch-Aussage; das Routing hat keine Entscheidung.
- **Ein allgemeiner Recovery-Weg für Schema-Fehler** — das Handbuch beschreibt
  die Abhilfe für die Nichtanwendbarkeit einer Transformationsregel, nicht mehr
  (`BEO-PGC/kein-admin-weg-schema-fehler-recovery`, offen, 1×).

## 2. Definition of Done

- [ ] SDK-Beleg: die Modelle der Row Images in den drei Packages behandeln
      `old_image`/`new_image` undurchsichtig (Anker am Start gelesen:
      `sdks/csharp/PgChangeFeed.Client/…/Change.cs`,
      `sdks/kotlin/pgchangefeed-kotlin/…/Change.kt`,
      `sdks/python/pgchangefeed/src/pgchangefeed/models.py`), und **kein
      Produktivcode** unter `sdks/*/` greift auf einen festen Bild-Schlüssel zu
      (Suchlauf, Befund und Nichtbefund im Bericht); die Test-Schlüssel der
      SDK-Tests sind Fixture-Daten. *Zu belegen durch:* der Suchlauf in §3 und
      Lesen der Modelle; ändert sich kein SDK-Code, laufen die SDK-Tiers
      unverändert — kein `make sdk-*` nötig, die Nicht-Ausführung wird im
      Bericht begründet.
- [ ] Betriebsdokumentation im Benutzerhandbuch
      ([`docs/user/benutzerhandbuch.md`](../../../user/benutzerhandbuch.md)):
      ein neuer Abschnitt in §4 (Aufgaben) „Transformationsregel konfigurieren“
      mit dem Aufruf von `cdc.set_transformation`/`cdc.remove_transformation`
      (Beispiele `rename_column`, `map_value`; Aufrufform der Regelform: der
      Parameter `rule_spec` ist `json`, ein Literal oder ein `::json`-Wert wird
      angenommen, ein `::jsonb`-Wert nicht — „function cdc.set_transformation(…,
      jsonb) does not exist“, gemessen im Review
      `review-slice-transformationen-antragsweg-schema`,
      [`ADR-0125`](../../adr/0125-transformationen-parametertyp-regelform-json.md)
      Folgepflicht 2; Annahmemenge nach
      [`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md)
      Festlegung 1: SQL-`NULL`, JSON-`null`, ein Wert ohne Objekt und ein
      doppelter Schlüssel werden angenommen und in Go geprüft, `\u0000`, eine
      Zahl außerhalb des Zahlbereichs von `numeric` und ein Syntaxfehler enden
      beim Aufruf ohne Antrags-Zeile), der Rolle `cdc_admin`, dem
      Status `pending`/`applied`/`failed` und den Fehlertexten von K1–K4, der
      Wirkung („ab `applied` für künftige Changes, nicht rückwirkend; die
      Rohform wird nicht gespeichert“, der Informationsverlust bei
      `map_value`), dem Verhältnis zum Spaltenausschluss (der Ausschluss gilt
      zuerst), der Dauerhaftigkeit (überlebt Neustart und
      Deaktivierung/Aktivierung; Grenze aus `antragsweg-usecase`: eine vermerkte Regel, die
      ein älterer Binärstand nicht lesen kann — etwa eine `map_value`-Regel unter einem Stand
      vor `map-value` —, hält Prozessstart und jeden Regel-Antrag der Quelle an, hergeleitet
      aus dem Quelltext, nicht erprobt), der Nichtanwendbarkeit samt Abhilfe-Prozedur
      (`cdc.remove_transformation` beantragen, Prozess starten — so, wie
      `e2e-abhilfe` sie belegt hat), der Reihenfolge der Aufrufe (Aufrufe einer
      Transaktion werden in Aufrufreihenfolge verarbeitet; Anträge auf dieselbe
      Regel oder Spalte nicht aus überlappenden Transaktionen absetzen —
      [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)
      Folgepflicht 4 und Festlegung 3 Punkt 1; die Aussage gilt für alle sieben
      Funktionen und gehört in den Transformations-Abschnitt in §4) und dem
      Backfill-Bezug; dazu die Zeile
      `schema` in §6 Fehlerklassen, die Rollen-Beschreibung in §2, das Glossar
      und die Änderungshistorie. *Zu belegen durch:* Review des Abschnitts
      gegen die Belege von `e2e-wirkung`/`e2e-abhilfe` und `make docs-check`.
- [ ] Jede Zahl und jede Wirkungs-Aussage trägt ihren Ursprung
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12): gemessen (Lauf), übernommen
      (Beleg-Anker) oder abgeleitet; eine Aussage, die nur die ADR trägt, steht
      als Zusage bzw. „erwartet“. Die Änderungshistorie trägt eine Zeile
      (`BEO-PGC/handbuch-versionshistorie-uebersprungen`, verkörpert, 3×). *Zu
      belegen durch:* Review; die Suchläufe in §3 belegen, dass keine
      Handbuch-Stelle mit einer Aufzählung der Antragsarten oder der Klasse
      `schema` übersehen ist.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: der Slice **ist** das Doku-Update (Handbuch,
      Änderungshistorie); `harness/README.md` ist unberührt — die Sensor-Zeilen
      tragen `e2e-wirkung` und `e2e-abhilfe`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | neuer Abschnitt in §4, §2 Rollen, §6 Fehlerklassen (Zeile `schema`), Glossar, Änderungshistorie. |
| `sdks/**` (Modelle, Tests) | lesen | Beleg-Suchlauf; Änderung nur nach Plan-Nachzug. |
| SDK-READMEs (`sdks/csharp/README.md`, `sdks/python/README.md`, `sdks/kotlin/pgchangefeed-kotlin/README.md`) | prüfen | Aussagen über Row Images (jede README trägt einen Abschnitt zum Change-Objekt mit `old_image`/`new_image`, Suchlauf am Start); die README ist die Paketbeschreibung auf PyPI und NuGet und erscheint dort erst mit einer neuen Package-Version, eine Änderung wird im Bericht benannt. Jede Textdatei unter `sdks/` trägt keine interne Kennung (`SPEC-`/`ADR-`/`ARC-`/`LH-FA-`/`LH-QA-`, Slice-/Welle-Name): `make sdk-public-doc-check` färbt sonst jedes `make sdk-pack-*`. |
| `docs/user/benutzerhandbuch-standard.md` | prüfen | Standard-Form des Handbuchs — der neue Abschnitt folgt ihr. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Menge der
Antragsarten“, „die Sätze über Row-Image-Schlüssel und Fehlerklasse `schema`“,
„die Aussagen der SDK-Modelle über Row Images“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Zugriffe auf feste Bild-Schlüssel im Produktivcode der SDKs | `grep -rn 'old_data\|new_data\|oldData\|newData\|old_image\|new_image' sdks --include=*.py --include=*.cs --include=*.kt`, dann Fundstellen in `src/main`/`src/pgchangefeed`/`PgChangeFeed.Client/` (ohne Tests) lesen | *(Implementer trägt ein)* | Nur Durchreichen (opake Typen `JsonElement?`, `Any | None`) ist zulässig; ein Zugriff auf einen benannten Schlüssel ist ein Befund |
| Aufzählungen der Antragsarten im Handbuch | `grep -rn 'exclude_column\|enable_table' docs/user` | *(Implementer trägt ein)* | jede Aufzählung trägt die zwei weiteren Arten oder begründet, warum nicht |
| Sätze über die Fehlerklasse `schema` | `grep -rn 'schema' docs/user/benutzerhandbuch.md` | *(Implementer trägt ein)* | Fehlerklassen-Tabelle (§6), „Container startet nicht“ und „Neustart nach einem Fehler“ nennen die Nichtanwendbarkeit und ihre Abhilfe, ohne die Recovery zu verallgemeinern |
| Beschreibung des SQL-Lesezugriffs `cdc.changes` und `GET /changes` im Handbuch | `grep -rn 'Row Image\|row_image\|new_data' docs/user` | *(Implementer trägt ein)* | Sätze über die Bildform nennen den Regelstand („die Schlüsselmenge folgt dem Regelstand zum Erfassungszeitpunkt“) |
| Übernahme aus Nachbar-Dokumenten | Lesen des `harness/README.md`-Zeile `make test-integration` | *(Implementer trägt ein)* | Wirkungs-Aussagen des Handbuchs stimmen mit den dort genannten Belegen überein |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-e2e-abhilfe`
in `done/` liegt (die Wirkung und die Abhilfe sind belegt, bevor das Handbuch
sie beschreibt) und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1).
Letzter Slice der Welle.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  Doku-Zug mit einem Such-Beleg; sprengte er den Umfang, wäre der SDK-Beleg
  (erster Liefer-Punkt) der abtrennbare Teil.
- `in-progress` → `open` (blockiert): falls der SDK-Suchlauf einen Zugriff auf
  feste Bild-Schlüssel im Produktivcode findet (Plan-Nachzug: SDK-Änderung samt
  Versionshebung, eine Architect-Frage nach
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 8 wird geführt) oder falls ein Beleg von
  `e2e-wirkung`/`e2e-abhilfe` eine Handbuch-Aussage nicht trägt (dann steht die
  Aussage nicht im Handbuch).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Das Handbuch behauptet mehr, als die Belege tragen**
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B), etwa „die Abhilfe
  wirkt“ vor dem realen Beleg. *Erwartet, zu belegen durch:* Review liest jede
  Wirkungs-Aussage gegen den Beleg-Anker; der Start-Trigger stellt den Beleg
  voran. **Ausgang:** *(bei Closure)*
- **Das Handbuch übersieht eine Stelle**, die die Antragsarten oder die Klasse
  `schema` aufzählt (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, offen,
  2×; `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  verkörpert, 3×). *Erwartet, zu belegen durch:* die Suchläufe in §3.
  **Ausgang:** *(bei Closure)*
- **Der SDK-Suchlauf findet einen Schlüssel-Zugriff**, den die Erwartung
  („keine Code-Änderung“) ausschließt. *Erwartet, zu belegen durch:* der
  Suchlauf mit Befund und Nichtbefund im Bericht. **Ausgang:** *(bei Closure)*
- **Die Abhilfe-Prozedur wird zu einer allgemeinen Recovery** verallgemeinert
  (`BEO-PGC/kein-admin-weg-schema-fehler-recovery`, offen, 1×). *Erwartet, zu
  belegen durch:* Review des Abschnitts auf die benannte Ursache. **Ausgang:**
  *(bei Closure)*
- **Zahlen und Fristen im Handbuch ohne Ursprung** (etwa eine Wartezeit nach
  dem Neustart) — jede Zahl trägt Lauf und Ursprung oder entfällt. *Erwartet,
  zu belegen durch:* Review. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); das Handbuch und die SDK-Verzeichnisse (`sdks/*/`, nur
gelesen) sind keine Sub-Areas dieses Slice — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` und
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, je 3×,
einschlägig — Kern dieses Slice),
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (offen, 2×, einschlägig —
dieser Slice ist die Adresse von `antragsweg-schema`; sein §2 deckt den
aufgeschobenen Gegenstand), `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
(offen, 2×), `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×,
Abgrenzung §1), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×) und
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 6×,
Risiko §6), `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert,
3×, gesichtet — SDK-READMEs werden nur gelesen),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
