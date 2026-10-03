package messagecode

// Die Tabelle der Meldungscodes ist die Quelle der Wahrheit: jede Stelle des
// Programms verwendet die Konstante, kein zweites Literal. Der Katalog im
// Benutzerhandbuch führt dieselbe Menge; `make meldungscodes-check` vergleicht
// beide mit dem Quelltext. Ein Code wird nie neu belegt, ein entfallener Code
// wird zurückgezogen (`StatusWithdrawn`) und bleibt in Tabelle und Katalog.
//
// Die erste Ziffer der Fehlercodes ist die Fehlerklasse (1 `transient`,
// 2 `configuration`, 3 `permission`, 4 `schema`, 5 `storage`, 6 `replication`,
// 7 `internal`); 8 ist die Ablehnung einer Aufrufer-Eingabe ohne Fehlerklasse.
// `…000` ist der Rückfall der Klasse.

// Klasse 1, `transient`.
const (
	TransientFallback         Code = "PCF-E1000"
	RetryExhausted            Code = "PCF-E1002"
	SnapshotSourceUnavailable Code = "PCF-E1003"
)

// Klasse 2, `configuration`.
const (
	ConfigurationFallback     Code = "PCF-E2000"
	WiringPrecondition        Code = "PCF-E2001"
	ReplicationConfiguration  Code = "PCF-E2002"
	ActivationIdentifier      Code = "PCF-E2003"
	SnapshotConfiguration     Code = "PCF-E2004"
	BackfillTableNotActivated Code = "PCF-E2005"
	BackfillStateChanged      Code = "PCF-E2006"
	SchemaSourceMissing       Code = "PCF-E2007"
)

// Klasse 3, `permission`.
const (
	PermissionFallback    Code = "PCF-E3000"
	ReplicationPermission Code = "PCF-E3001"
	SnapshotPermission    Code = "PCF-E3002"
)

// Klasse 4, `schema`.
const (
	SchemaFallback             Code = "PCF-E4000"
	DecodeUnreadable           Code = "PCF-E4001"
	TruncateUnsupported        Code = "PCF-E4002"
	SchemaChangeIncompatible   Code = "PCF-E4003"
	TransformationInapplicable Code = "PCF-E4004"
	RoutingInapplicable        Code = "PCF-E4005"
)

// Klasse 5, `storage`.
const (
	StorageFallback      Code = "PCF-E5000"
	ChangeStoreFailed    Code = "PCF-E5001"
	BackfillStoreFailed  Code = "PCF-E5003"
	ConsumerStateFailed  Code = "PCF-E5004"
	HeartbeatStoreFailed Code = "PCF-E5006"
	SchemaStoreFailed    Code = "PCF-E5007"
	SnapshotReadFailed   Code = "PCF-E5008"
)

// Klasse 6, `replication`.
const (
	ReplicationFallback     Code = "PCF-E6000"
	ReplicationStreamFailed Code = "PCF-E6001"
	ReplicationAckFailed    Code = "PCF-E6002"
	StreamOrderViolated     Code = "PCF-E6003"
	SnapshotSlotFailed      Code = "PCF-E6004"
)

// Klasse 7, `internal`. Die Codes `…001` bis `…010` tragen die Sentinels der
// Ports, die der Capture-Pfad nicht einer Klasse zuordnet (`classifyRunError`
// bildet sie auf `internal` ab); der Text eines Runs trägt für dieselben
// Ursachen die Codes der Klassen 1 bis 6 (`failureCode`).
const (
	InternalFallback           Code = "PCF-E7000"
	NotifyFailed               Code = "PCF-E7001"
	RequestQueueFailed         Code = "PCF-E7002"
	DiagnosticsReadFailed      Code = "PCF-E7003"
	BackfillStoreFault         Code = "PCF-E7004"
	SchemaStoreFault           Code = "PCF-E7005"
	SnapshotPermissionFault    Code = "PCF-E7006"
	SnapshotConfigurationFault Code = "PCF-E7007"
	SnapshotSourceFault        Code = "PCF-E7008"
	SnapshotSlotFault          Code = "PCF-E7009"
	SnapshotReadFault          Code = "PCF-E7010"
)

// Bereich 8, Ablehnung einer Aufrufer-Eingabe (keine Fehlerklasse).
const (
	RejectedFallback           Code = "PCF-E8000"
	RejectedIdentifierEmpty    Code = "PCF-E8001"
	RejectedSourceEmpty        Code = "PCF-E8002"
	RejectedSchemaEmpty        Code = "PCF-E8003"
	RejectedTableEmpty         Code = "PCF-E8004"
	RejectedKindUnknown        Code = "PCF-E8005"
	RejectedColumnEmpty        Code = "PCF-E8006"
	RejectedRuleName           Code = "PCF-E8010"
	RejectedRuleSpec           Code = "PCF-E8011"
	RejectedRouteForm          Code = "PCF-E8012"
	RejectedRuleNameTaken      Code = "PCF-E8020"
	RejectedColumnHasRule      Code = "PCF-E8021"
	RejectedTargetCollides     Code = "PCF-E8022"
	RejectedRuleNotKept        Code = "PCF-E8023"
	RejectedColumnMissing      Code = "PCF-E8024"
	RejectedTableMissing       Code = "PCF-E8025"
	RejectedOrderTaken         Code = "PCF-E8030"
	RejectedConditionTaken     Code = "PCF-E8031"
	RejectedWithoutWhenOrder   Code = "PCF-E8032"
	RejectedConditionExcluded  Code = "PCF-E8033"
	RejectedColumnHasCondition Code = "PCF-E8034"
	RejectedNotActivated       Code = "PCF-E8040"
	RejectedRunActive          Code = "PCF-E8041"
)

// Table führt jeden vergebenen Code genau einmal; die Klasse eines Eintrags
// ist die der ersten Ziffer (`Entry.Class`, geprüft im Test des Pakets).
var Table = []Entry{
	{TransientFallback, ClassTransient, StatusActive},
	{RetryExhausted, ClassTransient, StatusActive},
	{SnapshotSourceUnavailable, ClassTransient, StatusActive},

	{ConfigurationFallback, ClassConfiguration, StatusActive},
	{WiringPrecondition, ClassConfiguration, StatusActive},
	{ReplicationConfiguration, ClassConfiguration, StatusActive},
	{ActivationIdentifier, ClassConfiguration, StatusActive},
	{SnapshotConfiguration, ClassConfiguration, StatusActive},
	{BackfillTableNotActivated, ClassConfiguration, StatusActive},
	{BackfillStateChanged, ClassConfiguration, StatusActive},
	{SchemaSourceMissing, ClassConfiguration, StatusActive},

	{PermissionFallback, ClassPermission, StatusActive},
	{ReplicationPermission, ClassPermission, StatusActive},
	{SnapshotPermission, ClassPermission, StatusActive},

	{SchemaFallback, ClassSchema, StatusActive},
	{DecodeUnreadable, ClassSchema, StatusActive},
	{TruncateUnsupported, ClassSchema, StatusActive},
	{SchemaChangeIncompatible, ClassSchema, StatusActive},
	{TransformationInapplicable, ClassSchema, StatusActive},
	{RoutingInapplicable, ClassSchema, StatusActive},

	{StorageFallback, ClassStorage, StatusActive},
	{ChangeStoreFailed, ClassStorage, StatusActive},
	{BackfillStoreFailed, ClassStorage, StatusActive},
	{ConsumerStateFailed, ClassStorage, StatusActive},
	{HeartbeatStoreFailed, ClassStorage, StatusActive},
	{SchemaStoreFailed, ClassStorage, StatusActive},
	{SnapshotReadFailed, ClassStorage, StatusActive},

	{ReplicationFallback, ClassReplication, StatusActive},
	{ReplicationStreamFailed, ClassReplication, StatusActive},
	{ReplicationAckFailed, ClassReplication, StatusActive},
	{StreamOrderViolated, ClassReplication, StatusActive},
	{SnapshotSlotFailed, ClassReplication, StatusActive},

	{InternalFallback, ClassInternal, StatusActive},
	{NotifyFailed, ClassInternal, StatusActive},
	{RequestQueueFailed, ClassInternal, StatusActive},
	{DiagnosticsReadFailed, ClassInternal, StatusActive},
	{BackfillStoreFault, ClassInternal, StatusActive},
	{SchemaStoreFault, ClassInternal, StatusActive},
	{SnapshotPermissionFault, ClassInternal, StatusActive},
	{SnapshotConfigurationFault, ClassInternal, StatusActive},
	{SnapshotSourceFault, ClassInternal, StatusActive},
	{SnapshotSlotFault, ClassInternal, StatusActive},
	{SnapshotReadFault, ClassInternal, StatusActive},

	{RejectedFallback, ClassNone, StatusActive},
	{RejectedIdentifierEmpty, ClassNone, StatusActive},
	{RejectedSourceEmpty, ClassNone, StatusActive},
	{RejectedSchemaEmpty, ClassNone, StatusActive},
	{RejectedTableEmpty, ClassNone, StatusActive},
	{RejectedKindUnknown, ClassNone, StatusActive},
	{RejectedColumnEmpty, ClassNone, StatusActive},
	{RejectedRuleName, ClassNone, StatusActive},
	{RejectedRuleSpec, ClassNone, StatusActive},
	{RejectedRouteForm, ClassNone, StatusActive},
	{RejectedRuleNameTaken, ClassNone, StatusActive},
	{RejectedColumnHasRule, ClassNone, StatusActive},
	{RejectedTargetCollides, ClassNone, StatusActive},
	{RejectedRuleNotKept, ClassNone, StatusActive},
	{RejectedColumnMissing, ClassNone, StatusActive},
	{RejectedTableMissing, ClassNone, StatusActive},
	{RejectedOrderTaken, ClassNone, StatusActive},
	{RejectedConditionTaken, ClassNone, StatusActive},
	{RejectedWithoutWhenOrder, ClassNone, StatusActive},
	{RejectedConditionExcluded, ClassNone, StatusActive},
	{RejectedColumnHasCondition, ClassNone, StatusActive},
	{RejectedNotActivated, ClassNone, StatusActive},
	{RejectedRunActive, ClassNone, StatusActive},
}
