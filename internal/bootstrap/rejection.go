package bootstrap

import (
	"errors"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// administrationFailureText bildet den `error_message` eines `failed`
// Antrags aus dem Fehler seiner Verarbeitung. Ein klassifizierter Fehler
// trägt seinen Kopf `Fehlerklasse <klasse> [<code>]: …` selbst. Die Ablehnung
// einer Aufrufer-Eingabe (`rejectionCode`) steht als
// `abgelehnt [<code>]: <Klartext>`. Jeder andere Fehler ist unerwartet und
// trägt den Kopf der Klasse `internal`; kein `failed`-Antrag bleibt ohne Code.
func administrationFailureText(err error) string {
	if _, ok := messagecode.From(err); ok {
		return err.Error()
	}
	if code, ok := rejectionCode(err); ok {
		return messagecode.RejectionMessage(code, err.Error())
	}
	return messagecode.Head(messagecode.InternalFallback) + err.Error()
}

// rejectionCode ordnet den Fehler einer abgelehnten Aufrufer-Eingabe ihrem
// Code im Bereich `E8` zu. Die Zuordnung folgt den Ablehnungsgründen der
// Antragsarten: Regelname, Regelform, Konfliktfreiheit (K1 bis K4, R1 bis R6),
// fehlende Quell-Objekte und die Vorbedingungen des Backfills.
func rejectionCode(err error) (messagecode.Code, bool) {
	switch {
	case errors.Is(err, domainerrors.ErrInvalidRuleName):
		return messagecode.RejectedRuleName, true
	case errors.Is(err, domainerrors.ErrInvalidRuleSpec),
		errors.Is(err, domainerrors.ErrUnknownTransformationKind),
		errors.Is(err, domainerrors.ErrUnknownRuleSpecKey),
		errors.Is(err, domainerrors.ErrInvalidTransformation),
		errors.Is(err, domainerrors.ErrTransformationTargetIsColumn):
		return messagecode.RejectedRuleSpec, true
	case errors.Is(err, domainerrors.ErrInvalidRoute),
		errors.Is(err, domainerrors.ErrInvalidRouteTarget):
		return messagecode.RejectedRouteForm, true
	case errors.Is(err, domainerrors.ErrRuleNameTaken):
		return messagecode.RejectedRuleNameTaken, true
	case errors.Is(err, domainerrors.ErrColumnHasRule):
		return messagecode.RejectedColumnHasRule, true
	case errors.Is(err, domainerrors.ErrTargetCollidesWithRule),
		errors.Is(err, domainerrors.ErrTargetCollidesWithColumn):
		return messagecode.RejectedTargetCollides, true
	case errors.Is(err, domainerrors.ErrRuleNotKept):
		return messagecode.RejectedRuleNotKept, true
	case errors.Is(err, inbound.ErrSourceColumnMissing):
		return messagecode.RejectedColumnMissing, true
	case errors.Is(err, inbound.ErrSourceTableMissing):
		return messagecode.RejectedTableMissing, true
	case errors.Is(err, domainerrors.ErrRouteOrderTaken):
		return messagecode.RejectedOrderTaken, true
	case errors.Is(err, domainerrors.ErrRouteConditionTaken):
		return messagecode.RejectedConditionTaken, true
	case errors.Is(err, domainerrors.ErrRouteWithoutWhenTaken),
		errors.Is(err, domainerrors.ErrRouteWithoutWhenNotLast),
		errors.Is(err, domainerrors.ErrRouteBehindWithoutWhen):
		return messagecode.RejectedWithoutWhenOrder, true
	case errors.Is(err, domainerrors.ErrRoutingColumnExcluded):
		return messagecode.RejectedConditionExcluded, true
	case errors.Is(err, domainerrors.ErrColumnHasRouteCondition):
		return messagecode.RejectedColumnHasCondition, true
	case errors.Is(err, domainerrors.ErrTableNotActivated):
		return messagecode.RejectedNotActivated, true
	case errors.Is(err, domainerrors.ErrBackfillRunActive):
		return messagecode.RejectedRunActive, true
	default:
		return messagecode.RejectedFallback, false
	}
}
