# BEO-PGC/arbeit-ueberholt-stehenden-traeger

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis
zwischen einer Änderung und den Trägern, die ihren Gegenstand **beschreiben**,
keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine Arbeit **bewegt eine gemessene Eigenschaft** eines
Gegenstands — und macht damit Träger falsch, die diese Eigenschaft
**beschreiben**, ohne sie selbst anzufassen. Die Träger stehen **nicht im Diff**
und werden von **keinem Sensor** gelesen: `d-check` prüft Links, `coverage-gate`
prüft eine Zahl gegen eine Schwelle, `generated-sync` prüft Bytes. Falsch werden
sie trotzdem, und zwar in dem Moment, in dem die Arbeit landet.

**Der Unterschied zu den benachbarten Klassen.** `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
beschreibt einen Träger, der **schon** driftet; hier driftet er **durch diese
Arbeit**, und der Drift ist die Folge einer korrekten Änderung, nicht einer
Nachlässigkeit. `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` beschreibt einen
**Beleg**, der seinen Satz nicht trägt; hier trägt der Satz seinen Beleg
**nicht mehr**. Und §3.12 (Herkunft von Aussagen) hilft nicht: die Aussage trug
ihren Ursprung korrekt — sie war zum Zeitpunkt ihrer Niederschrift wahr.

Belegt an einem abgeschlossenen Vorgang:

- **`slice-091`** (Delta-Review `review-slice-091-delta` D-1): Der Slice gab dem
  Paket `internal/adapters/driving/grpc/streamv1` eine eigene Testdatei. Damit
  wurde eine **stehende Liste** in `harness/sensors/coverage-gate.md` an **drei**
  Stellen falsch, die niemand im Diff hatte: die Überschrift („**Fünf** Pakete
  des Gegenstands [ohne] Testdatei"), der Schlusssatz (der Lauf weise die beiden
  Pakete mit Statements als `coverage: 0.0% of statements` aus) und die
  Gruppierung selbst. **Beide Nachbarsätze waren vorher wahr** — gemessen am
  Parent-Stand (`git archive f90c3f4^`, `go list`): dort waren es genau die fünf
  genannten Pakete mit `Test=0 XTest=0`. Die Liste **war** mechanisch; die
  Arbeit hat sie es nicht mehr sein lassen.

**Und die Behebung ist der zweite Teil der Beobachtung.** Der erste
Reparaturversuch hat die Gruppe mit `go list -f '{{len .TestGoFiles}}'`
begründet — gemessen trifft diese Zählung **23** der 31 Pakete des Gegenstands
und ist damit **keine** Gruppierungsregel; die Reparatur hat eine falsche
Stützung an die Stelle der überholten Aussage gesetzt. Gefunden hat das erst der
**Delta-Review**, nicht das Review und nicht die Verifikation — beide hatten
keinen Anlass, eine Datei zu öffnen, die nicht im Diff lag.

**Warum das zählt:** Eine Änderung wirkt über ihren Diff hinaus, und die
Reichweite ist nicht die der Datei-Liste, sondern die des **Gegenstands**, den
sie bewegt. Wer eine Eigenschaft ändert, die anderswo beschrieben steht, hat
einen Leser-Pfad eröffnet, den kein Gate geht — die verfügbare Falsifikation ist
wieder die **Messung**, und der billigste Wächter ist die Frage „welche Träger
beschreiben das, was ich hier gerade geändert habe?". Ein Sensor ist **nicht**
vorgeschlagen: er müsste wissen, welche Sätze von welcher Eigenschaft abhängen,
und das weiß er nicht — die Abhängigkeit steht in Prosa.
