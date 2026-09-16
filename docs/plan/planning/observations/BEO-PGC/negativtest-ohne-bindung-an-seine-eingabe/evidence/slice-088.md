# Beleg: slice-088

Vorgang: `slice-088` — Coverage-Tail „Reine Übersetzung" (Cluster B der
`welle-20`).

Fund: **Der Vorgang, der diese Klasse beheben sollte, ist an ihr gescheitert —
und hat sie dadurch an sich selbst gefunden.** Der Slice besteht aus Zusagen,
die an ihrer Eingabeseite hängen müssen; der Reviewer hat an **zwei** der neuen
Tests nachgewiesen, dass sie es nicht taten:

- `TestConsumeRelationOnUnboundTableStaysNoop` und
  `TestConsumeRelationWithoutRegisteredVersionStaysNoop`
  (`internal/adapters/driving/replication/mapper/mapper_test.go`) behaupteten
  **Fehler-Abwesenheit gegen einen Stub, der seine Argumente ignoriert**.
- Gemessen: die Mutationen M4 (`if !activated { return nil }` entfernt) und M5
  (`if !found { return nil }`) ließen `go test ./internal/... ./cmd/...` **je
  Exit 0**.

**Das Bemerkenswerte:** der Implementer hatte **fünf** eigene Mutationen
gefahren und alle rot gesehen — die zwei löchrigen waren die, die er für
selbstverständlich hielt. Die Fixrunde hat sie über einen Stub gebunden, der die
**empfangenen Kennungen protokolliert**, und der Verifier hat M4/M5
**eigenständig reproduziert** (Exit 1, je genau der betroffene Test) sowie M6/M7
als Bindungssonden bestätigt.

Der Fall ist damit ein **vierter Vorgang** derselben Klasse, und er zeigt ihre
scharfe Kante: sie trifft nicht die Zusagen, die man mutiert, sondern die, die
man übersieht.

Quelle: `docs/reviews/review-slice-088.md` (F-1, mit eigener Messung) ·
`docs/reviews/review-slice-088-delta.md` (Reproduktion M4–M7) ·
`docs/reviews/verify-slice-088.md` (eigene Reproduktion) ·
`internal/adapters/driving/replication/mapper/mapper_test.go`.
