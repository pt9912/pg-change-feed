package messagecode_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

var classes = []messagecode.Class{
	messagecode.ClassTransient, messagecode.ClassConfiguration, messagecode.ClassPermission,
	messagecode.ClassSchema, messagecode.ClassStorage, messagecode.ClassReplication,
	messagecode.ClassInternal,
}

// TestTableFormat trägt die Form jedes Codes der Tabelle: die ERE aus der
// Festlegung, Schwere `E` oder `W` (`I` ist reserviert), ein bekannter Status.
// Färbt rot, sobald ein Eintrag eine andere Form trägt.
func TestTableFormat(t *testing.T) {
	form := regexp.MustCompile(`^PCF-[EWI][0-9]{4}$`)
	if len(messagecode.Table) == 0 {
		t.Fatal("Tabelle ist leer")
	}
	for _, e := range messagecode.Table {
		if !form.MatchString(string(e.Code)) || !messagecode.Valid(e.Code) {
			t.Errorf("%q: Form verletzt", e.Code)
		}
		if e.Code[4] != 'E' && e.Code[4] != 'W' {
			t.Errorf("%q: Schwere %q, erwartet E oder W", e.Code, e.Code[4])
		}
		if e.Status != messagecode.StatusActive && e.Status != messagecode.StatusWithdrawn {
			t.Errorf("%q: unbekannter Status %q", e.Code, e.Status)
		}
	}
}

// TestTableClassMatchesFirstDigit trägt: die Klasse einer Tabellenzeile ist
// die der ersten Ziffer. Eine Zeile mit falscher Klasse färbt den Test rot.
func TestTableClassMatchesFirstDigit(t *testing.T) {
	for _, e := range messagecode.Table {
		if e.Code[4] == 'W' {
			continue
		}
		want, ok := messagecode.DigitClass(e.Code)
		if !ok {
			t.Errorf("%q: erste Ziffer außerhalb 1 bis 8", e.Code)
			continue
		}
		if e.Class != want {
			t.Errorf("%q: Klasse %q, erste Ziffer nennt %q", e.Code, e.Class, want)
		}
	}
}

// TestWarningCodesCarryTheirArea trägt: jeder Warncode der Tabelle trägt keine
// Fehlerklasse und einen Bereich 1 bis 5 in der ersten Ziffer, und jede
// benannte Konstante liegt im Bereich ihres Namens (Bereich 1 Erfassung,
// 2 Backfill, 3 Retention und Speicher, 4 Verwaltung). Färbt rot, sobald ein
// Warncode in einen anderen Bereich oder mit einer Klasse eingetragen wird.
func TestWarningCodesCarryTheirArea(t *testing.T) {
	named := map[int][]messagecode.Code{
		1: {
			messagecode.WarnNotifyFailed, messagecode.WarnStreamPublish, messagecode.WarnStreamNameSkipped,
			messagecode.WarnStreamTargetName, messagecode.WarnChangeNotEncodable, messagecode.WarnStreamCycleRetry,
		},
		2: {
			messagecode.WarnBackfillCleanup, messagecode.WarnBackfillQueueRead, messagecode.WarnBackfillRunNotAdvance,
			messagecode.WarnBackfillStateNotKept, messagecode.WarnBackfillRunFailed, messagecode.WarnBackfillRunInterrupt,
		},
		3: {messagecode.WarnWALMeasureFailed, messagecode.WarnWALOverWarn, messagecode.WarnRetentionFailed},
		4: {
			messagecode.WarnAdminWakeDisturbed, messagecode.WarnAdminRequestsRead, messagecode.WarnAdminPreflightExpiry,
			messagecode.WarnAdminRequestFailed, messagecode.WarnAdminRequestRejected, messagecode.WarnAdminOutcomeNotKept,
			messagecode.WarnAdminRowWithoutID,
		},
	}
	inTable := 0
	for _, e := range messagecode.Table {
		if e.Code[4] != 'W' {
			continue
		}
		inTable++
		if _, ok := messagecode.Area(e.Code); !ok {
			t.Errorf("%q: erste Ziffer außerhalb 1 bis 5", e.Code)
		}
		if e.Class != messagecode.ClassNone {
			t.Errorf("%q: Warnung trägt die Klasse %q", e.Code, e.Class)
		}
	}
	namedCount := 0
	for area, codes := range named {
		for _, code := range codes {
			namedCount++
			if got, ok := messagecode.Area(code); !ok || got != area {
				t.Errorf("%q: Bereich %d, erwartet %d", code, got, area)
			}
			if _, ok := messagecode.Lookup(code); !ok {
				t.Errorf("%q fehlt in der Tabelle", code)
			}
		}
	}
	if inTable != namedCount {
		t.Errorf("Tabelle führt %d Warncodes, die Konstanten nennen %d", inTable, namedCount)
	}
	if _, ok := messagecode.Area(messagecode.DecodeUnreadable); ok {
		t.Error("Area liest einen Fehlercode")
	}
	if _, ok := messagecode.Area("PCF-W9001"); ok {
		t.Error("Area liest den reservierten Bereich 9")
	}
}

// TestTableHasNoDuplicates trägt: kein Code steht zweimal in der Tabelle.
func TestTableHasNoDuplicates(t *testing.T) {
	seen := map[messagecode.Code]bool{}
	for _, e := range messagecode.Table {
		if seen[e.Code] {
			t.Errorf("%q steht mehrfach in der Tabelle", e.Code)
		}
		seen[e.Code] = true
	}
}

// TestEveryClassHasItsFallback trägt: jede der sieben Klassen führt ihren
// Rückfall `…000` in der Tabelle, und die Klassen sind die des Domänentyps.
// Fehlt ein Rückfall in der Tabelle, färbt der Test rot.
func TestEveryClassHasItsFallback(t *testing.T) {
	for i, class := range classes {
		if _, err := model.NewErrorClass(string(class)); err != nil {
			t.Errorf("Klasse %q ist keine Fehlerklasse des Domänentyps: %v", class, err)
		}
		fallback := messagecode.Fallback(class)
		want := messagecode.Code(fmt.Sprintf("PCF-E%d000", i+1))
		if fallback != want {
			t.Errorf("Rückfall der Klasse %q = %q, erwartet %q", class, fallback, want)
		}
		entry, ok := messagecode.Lookup(fallback)
		if !ok {
			t.Errorf("Rückfall %q der Klasse %q fehlt in der Tabelle", fallback, class)
			continue
		}
		if entry.Class != class || entry.Status != messagecode.StatusActive {
			t.Errorf("Rückfall %q: Klasse %q, Status %q", fallback, entry.Class, entry.Status)
		}
	}
	if _, ok := messagecode.Lookup(messagecode.RejectedFallback); !ok {
		t.Errorf("Rückfall der Ablehnungen %q fehlt in der Tabelle", messagecode.RejectedFallback)
	}
}

// TestClassesMatchDomainClasses trägt: die Klassen der Fehlercodes sind genau
// die sieben des Domänentyps, jede Fehlerklasse trägt mindestens den Rückfall.
func TestClassesMatchDomainClasses(t *testing.T) {
	domain := []model.ErrorClass{
		model.ErrorClassTransient, model.ErrorClassConfiguration, model.ErrorClassPermission,
		model.ErrorClassSchema, model.ErrorClassStorage, model.ErrorClassReplication,
		model.ErrorClassInternal,
	}
	for i, class := range classes {
		if string(domain[i]) != string(class) {
			t.Errorf("Klasse %d: Domäne %q, Meldungscodes %q", i+1, domain[i], class)
		}
	}
}

// TestNewBindsCodeToClass trägt: der Fehlerwert eines Codes beginnt mit dem
// Kopf `Fehlerklasse <klasse> [<code>]: `; ein Code außerhalb der Tabelle und
// eine Ablehnung sind kein Fehlerwert.
func TestNewBindsCodeToClass(t *testing.T) {
	err := messagecode.New(messagecode.DecodeUnreadable, "Ursache")
	want := "Fehlerklasse schema [PCF-E4001]: Ursache"
	if err.Error() != want {
		t.Errorf("Text = %q, erwartet %q", err.Error(), want)
	}
	wrapped := fmt.Errorf("Zusatz: %w", err)
	code, ok := messagecode.From(wrapped)
	if !ok || code != messagecode.DecodeUnreadable {
		t.Errorf("From = %q, %v", code, ok)
	}
	if got := messagecode.WithoutHead(wrapped); got != "Zusatz: Ursache" {
		t.Errorf("WithoutHead = %q", got)
	}
	if _, ok := messagecode.From(errors.New("ohne Code")); ok {
		t.Error("From liefert einen Code für einen Fehler ohne Code")
	}
	for _, bad := range []messagecode.Code{"PCF-E9999", messagecode.RejectedFallback} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("New(%q) ohne Panik", bad)
				}
			}()
			_ = messagecode.New(bad, "x")
		}()
	}
}

// TestCodesReadsTheWholeChain trägt: `Codes` liefert die Codes aller
// klassifizierten Fehlerwerte der Kette in der Reihenfolge des
// Tiefendurchlaufs, auch hinter `errors.Join` und mehreren `%w`; ein Fehler
// ohne Code liefert keinen. Färbt rot, sobald der Durchlauf nur dem ersten
// `Unwrap`-Ziel folgt.
func TestCodesReadsTheWholeChain(t *testing.T) {
	a := messagecode.New(messagecode.DecodeUnreadable, "a")
	b := messagecode.New(messagecode.ChangeStoreFailed, "b")
	cases := []struct {
		name string
		err  error
		want []messagecode.Code
	}{
		{"ohne Code", errors.New("x"), nil},
		{"nil", nil, nil},
		{"einer", fmt.Errorf("Zusatz: %w", a), []messagecode.Code{messagecode.DecodeUnreadable}},
		{"zwei in Kettenfolge", fmt.Errorf("%w: %w", a, b), []messagecode.Code{messagecode.DecodeUnreadable, messagecode.ChangeStoreFailed}},
		{"Join", errors.Join(b, a), []messagecode.Code{messagecode.ChangeStoreFailed, messagecode.DecodeUnreadable}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := messagecode.Codes(tc.err)
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Fatalf("Codes = %v, erwartet %v", got, tc.want)
			}
		})
	}
}

// TestInternalSentinelCodesStayInternal trägt: die Codes `…001` bis `…010` der
// Klasse 7 tragen die Klasse `internal` (erste Ziffer 7).
func TestInternalSentinelCodesStayInternal(t *testing.T) {
	for _, code := range []messagecode.Code{
		messagecode.NotifyFailed, messagecode.RequestQueueFailed, messagecode.DiagnosticsReadFailed,
		messagecode.BackfillStoreFault, messagecode.SchemaStoreFault, messagecode.SnapshotPermissionFault,
		messagecode.SnapshotConfigurationFault, messagecode.SnapshotSourceFault, messagecode.SnapshotSlotFault,
		messagecode.SnapshotReadFault,
	} {
		if got := messagecode.ClassOf(code); got != messagecode.ClassInternal {
			t.Errorf("%q: Klasse %q, erwartet internal", code, got)
		}
	}
}

// TestMessageForms trägt die Textformen von Run und Antrag.
func TestMessageForms(t *testing.T) {
	if got := messagecode.RunMessage(messagecode.SnapshotReadFailed, "Ursache"); got != "storage [PCF-E5008]: Ursache" {
		t.Errorf("RunMessage = %q", got)
	}
	if got := messagecode.RejectionMessage(messagecode.RejectedSourceEmpty, "Quelle ist leer: a"); got != "abgelehnt [PCF-E8002]: Quelle ist leer: a" {
		t.Errorf("RejectionMessage = %q", got)
	}
	if !strings.HasPrefix(messagecode.Head(messagecode.ReplicationStreamFailed), "Fehlerklasse replication [PCF-E6001]") {
		t.Errorf("Head = %q", messagecode.Head(messagecode.ReplicationStreamFailed))
	}
}
