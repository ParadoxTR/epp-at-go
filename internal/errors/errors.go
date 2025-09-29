package errors

import (
	"fmt"
)

type EPPError struct {
	Code    string
	Message string
	Details string
}

func (e *EPPError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("EPP error %s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("EPP error %s: %s", e.Code, e.Message)
}

func IsEPPError(err error) bool {
	_, ok := err.(*EPPError)
	return ok
}

func NewEPPError(code, message, details string) *EPPError {
	return &EPPError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

const (
	CodeSuccess              = "1000" // Command completed successfully
	CodeSuccessActionPending = "1001" // Command completed successfully; action pending
	CodeSuccessNoMessages    = "1300" // Command completed successfully; no messages
	CodeSuccessAckToDequeue  = "1301" // Command completed successfully; ack to dequeue
	CodeSuccessEndingSession = "1500" // Command completed successfully; ending session

	CodeUnknownCommand            = "2000" // Unknown command
	CodeCommandSyntaxError        = "2001" // Command syntax error
	CodeCommandUseError           = "2002" // Command use error
	CodeRequiredParameterMissing  = "2003" // Required parameter missing
	CodeParameterValueRangeError  = "2004" // Parameter value range error
	CodeParameterValueSyntaxError = "2005" // Parameter value syntax error

	CodeUnimplementedProtocolVersion = "2100" // Unimplemented protocol version
	CodeUnimplementedCommand         = "2101" // Unimplemented command
	CodeUnimplementedOption          = "2102" // Unimplemented option
	CodeUnimplementedExtension       = "2103" // Unimplemented extension

	CodeBillingFailure               = "2104" // Billing failure
	CodeObjectNotEligibleForRenewal  = "2105" // Object is not eligible for renewal
	CodeObjectNotEligibleForTransfer = "2106" // Object is not eligible for transfer

	CodeAuthenticationError      = "2200" // Authentication error
	CodeAuthorizationError       = "2201" // Authorization error
	CodeInvalidAuthorizationInfo = "2202" // Invalid authorization information

	CodeObjectPendingTransfer         = "2300" // Object pending transfer
	CodeObjectNotPendingTransfer      = "2301" // Object not pending transfer
	CodeObjectExists                  = "2302" // Object exists
	CodeObjectDoesNotExist            = "2303" // Object does not exist
	CodeObjectStatusProhibitsOp       = "2304" // Object status prohibits operation
	CodeObjectAssociationProhibitsOp  = "2305" // Object association prohibits operation
	CodeParameterValuePolicyError     = "2306" // Parameter value policy error
	CodeUnimplementedObjectService    = "2307" // Unimplemented object service
	CodeDataManagementPolicyViolation = "2308" // Data management policy violation

	CodeCommandFailed                    = "2400" // Command failed
	CodeCommandFailedServerClosing       = "2500" // Command failed; server closing connection
	CodeAuthenticationErrorServerClosing = "2501" // Authentication error; server closing connection
	CodeSessionLimitExceeded             = "2502" // Session limit exceeded; server closing connection
)

func IsSuccessCode(code string) bool {
	switch code {
	case CodeSuccess, CodeSuccessActionPending, CodeSuccessNoMessages,
		CodeSuccessAckToDequeue, CodeSuccessEndingSession:
		return true
	default:
		return false
	}
}
