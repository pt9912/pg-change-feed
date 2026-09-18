# Beleg: slice-093

Vorgang: `slice-093` — Coverage Cluster D2 (Bootstrap-Rest und Telemetrie).

Fund: **Zwei Stellen**, beide vom Review gefunden, beide an derselben Form: ein
Träger nennt einen Beleg, und der Beleg trägt ihn nicht.

- **F-3** (Review zu `slice-093`, LOW): `internal/adapters/driven/telemetry/slog_levels_internal_test.go`
  nannte eine Mutation als rot färbend — der Test färbte **rot**, aber **nicht aus
  dem genannten Grund**: der Handler steht auf `slog.LevelWarn`, ein INFO-Record
  wird gefiltert, die Zeile wird **gar nicht geschrieben**, und der Test fällt am
  `json.Unmarshal`. Die **Wirkung** war wahr, die **Begründung** nicht. Die
  Fixrunde nennt jetzt **zwei** gemessene Mutationen mit ihrer **verschiedenen**
  Wirkung (`Warn → ErrorContext` färbt über `decoded["level"]`; `Warn →
  InfoContext` färbt am JSON) — beide vom Delta-Review wörtlich bestätigt.
- **D-3** (Delta-Review zu `slice-093`, LOW): der Testkopf verwies auf Exit-Codes „im
  **Lauf-Bericht**" — gemessen ist das **kein Artefakt dieses Repos** (repo-weit
  genau ein Treffer: die Zeile selbst). Er nennt jetzt §7 der Closure-Notiz, und
  §7 **führt** die sieben Exit-Codes (Verifikation V-4).

**Ein Vorgang, eine Zählung.** Beide fallen in dieselbe Klasse und in denselben
Vorgang — der Zähler bewegt sich **einmal**.

**Was diese Vorkommen von den früheren unterscheidet:** Die drei ersten Belege
(`slice-084`, `-085`, `-091`) trafen **Befehle** und **Mutationsangaben** in
Kommentaren; hier trifft die Klasse zusätzlich eine **Adresse** — einen Verweis,
der auf ein Artefakt zeigt, das es nicht gibt. Beide Formen haben denselben Kern:
der Leser kann nicht nachschlagen, was der Satz ihm zu prüfen gibt.

Quelle: Review zu `slice-093` (F-3) ·
Delta-Review zu `slice-093` (D-3, Negativbefunde) ·
Verifikationsbericht zu `slice-093` (V-4) ·
`internal/adapters/driven/telemetry/slog_levels_internal_test.go`,
`internal/bootstrap/roles_rollout_file_internal_test.go` (berichtigt in `e140363`).
