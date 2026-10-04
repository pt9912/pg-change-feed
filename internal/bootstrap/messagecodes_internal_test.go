package bootstrap

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestSentinelsCarryTheirCodeAndClass bindet jeden klassifizierten Sentinel
// der Produktion an seinen Meldungscode, die Klasse des Codes an die Klasse
// der Verdrahtung (`classifyRunError`) und den Kopf des Textes an beide. Rot
// färbende Mutation: den Code eines Sentinels in seiner `messagecode.New`-Zeile
// ändern — der Fall dieses Sentinels trägt den falschen Code, bei einem Code
// anderer Klasse zusätzlich die falsche Klasse.
func TestSentinelsCarryTheirCodeAndClass(t *testing.T) {
	cases := []struct {
		name     string
		sentinel error
		code     messagecode.Code
	}{
		{"Wiederholung erschöpft", ErrTransientExhausted, messagecode.RetryExhausted},
		{"Wecksignal", outbound.ErrNotify, messagecode.NotifyFailed},
		{"Snapshot-Quelle", outbound.ErrSnapshotTransient, messagecode.SnapshotSourceFault},
		{"Verdrahtung", ErrConfiguration, messagecode.WiringPrecondition},
		{"API-Token-Liste", ErrAPITokenList, messagecode.APITokenListInvalid},
		{"TLS-Paar unvollständig", ErrTLSPairIncomplete, messagecode.TLSPairIncomplete},
		{"TLS-Paar nicht ladbar", ErrTLSPairUnusable, messagecode.TLSPairUnusable},
		{"Replication-Konfiguration", receive.ErrConfiguration, messagecode.ReplicationConfiguration},
		{"Aktivierungs-Bezeichner", postgresstorage.ErrActivationConfiguration, messagecode.ActivationIdentifier},
		{"Snapshot-Konfiguration", outbound.ErrSnapshotConfiguration, messagecode.SnapshotConfigurationFault},
		{"Replication-Berechtigung", receive.ErrPermission, messagecode.ReplicationPermission},
		{"Snapshot-Berechtigung", outbound.ErrSnapshotPermission, messagecode.SnapshotPermissionFault},
		{"Dekodierung", decode.ErrSchema, messagecode.DecodeUnreadable},
		{"TRUNCATE", mapper.ErrTruncateUnsupported, messagecode.TruncateUnsupported},
		{"Relation-Änderung", mapper.ErrIncompatibleSchemaChange, messagecode.SchemaChangeIncompatible},
		{"Transformationsregel", mapper.ErrTransformationNotApplicable, messagecode.TransformationInapplicable},
		{"Routing-Regel", mapper.ErrRoutingNotApplicable, messagecode.RoutingInapplicable},
		{"ChangeStore", outbound.ErrStorage, messagecode.ChangeStoreFailed},
		{"Antrags-Queue", outbound.ErrAdministrationStorage, messagecode.RequestQueueFailed},
		{"Backfill-Speicher", outbound.ErrBackfillStorage, messagecode.BackfillStoreFault},
		{"Consumer-Stand", outbound.ErrConsumerStateStorage, messagecode.ConsumerStateFailed},
		{"Diagnose-Views", outbound.ErrDiagnosticsStorage, messagecode.DiagnosticsReadFailed},
		{"Heartbeat", outbound.ErrHeartbeatStorage, messagecode.HeartbeatStoreFailed},
		{"Schema Store", outbound.ErrSchemaStoreStorage, messagecode.SchemaStoreFault},
		{"Snapshot lesen", outbound.ErrSnapshotStorage, messagecode.SnapshotReadFault},
		{"Replication-Stream", receive.ErrReplication, messagecode.ReplicationStreamFailed},
		{"Bestätigung", outbound.ErrReplication, messagecode.ReplicationAckFailed},
		{"BEGIN fehlt", mapper.ErrChangeWithoutBegin, messagecode.StreamOrderViolated},
		{"COMMIT ohne BEGIN", mapper.ErrCommitWithoutBegin, messagecode.StreamOrderViolated},
		{"BEGIN doppelt", mapper.ErrBeginWithoutCommit, messagecode.StreamOrderViolated},
		{"Snapshot-Slot", outbound.ErrSnapshotReplication, messagecode.SnapshotSlotFault},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, ok := messagecode.From(tc.sentinel)
			if !ok || code != tc.code {
				t.Fatalf("Code = %q (%v), erwartet %q", code, ok, tc.code)
			}
			class := messagecode.ClassOf(tc.code)
			wrapped := fmt.Errorf("%w: technischer Grund", tc.sentinel)
			if got := classifyRunError(wrapped); string(got) != string(class) {
				t.Fatalf("classifyRunError = %q, Klasse des Codes %q", got, class)
			}
			if _, err := model.NewErrorClass(string(class)); err != nil {
				t.Fatalf("Klasse %q ist keine Fehlerklasse: %v", class, err)
			}
			if want := messagecode.Head(tc.code); len(tc.sentinel.Error()) < len(want) || tc.sentinel.Error()[:len(want)] != want {
				t.Fatalf("Text %q beginnt nicht mit %q", tc.sentinel.Error(), want)
			}
		})
	}
}

// TestAdministrationFailureTextHasACodeForEveryFailure trägt den
// `error_message` eines `failed` Antrags: ein klassifizierter Fehler trägt
// seinen Kopf, die Ablehnung einer Aufrufer-Eingabe `abgelehnt [<code>]: …`
// mit dem Code ihres Grundes, jeder andere Fehler den Kopf der Klasse
// `internal`. Rot färbende Mutation: im `switch` von `rejectionCode` den Code
// eines Zweigs ändern — der Fall dieses Grundes trägt den falschen Code.
func TestAdministrationFailureTextHasACodeForEveryFailure(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"klassifiziert", fmt.Errorf("%w: Grund", outbound.ErrAdministrationStorage), "Fehlerklasse internal [PCF-E7002]: Persistenzfehler an der Antrags-Queue: Grund"},
		{"unerwartet", stderrors.New("überraschung"), "Fehlerklasse internal [PCF-E7000]: überraschung"},
		{"Regelname", domainerrors.ErrInvalidRuleName, "abgelehnt [PCF-E8010]: Regelname ist ungültig"},
		{"Regelform", domainerrors.ErrInvalidRuleSpec, "abgelehnt [PCF-E8011]: rule_spec ist ungültig"},
		{"unbekannter Regeltyp", domainerrors.ErrUnknownTransformationKind, "abgelehnt [PCF-E8011]: unbekannter Regeltyp"},
		{"Schlüssel", domainerrors.ErrUnknownRuleSpecKey, "abgelehnt [PCF-E8011]: unbekannter Schlüssel in rule_spec"},
		{"Transformationsform", domainerrors.ErrInvalidTransformation, "abgelehnt [PCF-E8011]: ungültige Transformationsregel"},
		{"Ziel gleich Quellspalte", domainerrors.ErrTransformationTargetIsColumn, "abgelehnt [PCF-E8011]: Zielname gleicht der Quellspalte"},
		{"Routing-Regel", domainerrors.ErrInvalidRoute, "abgelehnt [PCF-E8012]: ungültige Routing-Regel"},
		{"Zielname", domainerrors.ErrInvalidRouteTarget, "abgelehnt [PCF-E8012]: Zielname ist ungültig"},
		{"K1", domainerrors.ErrRuleNameTaken, "abgelehnt [PCF-E8020]: Regelname bereits vergeben"},
		{"K2", domainerrors.ErrColumnHasRule, "abgelehnt [PCF-E8021]: Spalte trägt bereits eine Regel"},
		{"K3 Regel", domainerrors.ErrTargetCollidesWithRule, "abgelehnt [PCF-E8022]: Zielname kollidiert mit einer anderen Regel"},
		{"K3 Spalte", domainerrors.ErrTargetCollidesWithColumn, "abgelehnt [PCF-E8022]: Zielname kollidiert mit einer Spalte der Tabelle"},
		{"K4 Regel", domainerrors.ErrRuleNotKept, "abgelehnt [PCF-E8023]: Regelname nicht geführt"},
		{"K4 Spalte", inbound.ErrSourceColumnMissing, "abgelehnt [PCF-E8024]: Spalte existiert nicht an der Quelle"},
		{"Tabelle fehlt", inbound.ErrSourceTableMissing, "abgelehnt [PCF-E8025]: Tabelle existiert nicht an der Quelle"},
		{"R2", domainerrors.ErrRouteOrderTaken, "abgelehnt [PCF-E8030]: order bereits vergeben"},
		{"R5", domainerrors.ErrRouteConditionTaken, "abgelehnt [PCF-E8031]: Bedingung bereits vergeben"},
		{"R4 vergeben", domainerrors.ErrRouteWithoutWhenTaken, "abgelehnt [PCF-E8032]: Regel ohne when bereits vorhanden"},
		{"R4 Ende", domainerrors.ErrRouteWithoutWhenNotLast, "abgelehnt [PCF-E8032]: Regel ohne when trägt nicht die höchste order"},
		{"R4 dahinter", domainerrors.ErrRouteBehindWithoutWhen, "abgelehnt [PCF-E8032]: order liegt hinter der Regel ohne when"},
		{"R3", domainerrors.ErrRoutingColumnExcluded, "abgelehnt [PCF-E8033]: Spalte ist ausgeschlossen"},
		{"R3 Gegenrichtung", domainerrors.ErrColumnHasRouteCondition, "abgelehnt [PCF-E8034]: Spalte trägt eine Routing-Bedingung"},
		{"Tabelle nicht aktiviert", domainerrors.ErrTableNotActivated, "abgelehnt [PCF-E8040]: Tabelle nicht aktiviert oder nicht in der Publication"},
		{"aktiver Run", domainerrors.ErrBackfillRunActive, "abgelehnt [PCF-E8041]: Für die Tabelle besteht bereits ein aktiver Backfill-Run"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := administrationFailureText(tc.err); got != tc.want {
				t.Fatalf("administrationFailureText = %q, erwartet %q", got, tc.want)
			}
		})
	}
}
