package pam

import "fmt"

// Error represents a PAM error.
type Error int

const (
	// ErrorOpen indicates a dlopen() failure when dynamically
	// loading a service module.
	ErrorOpen Error = iota + 1
	// ErrorSymbol indicates a symbol not found error.
	ErrorSymbol
	// ErrorService indicates an error in the service module.
	ErrorService
	// ErrorSystem indicates a system error.
	ErrorSystem
	// ErrorBuffer indicates a memory buffer error.
	ErrorBuffer
	// ErrorPermissionDenied indicates permission was denied.
	ErrorPermissionDenied
	// ErrorAuth indicates an authentication failure.
	ErrorAuth
	// ErrorCredentialsInsufficient indicates inability to access authentication data
	// due to insufficient credentials.
	ErrorCredentialsInsufficient
	// ErrorAuthInfoUnavailable indicates the underlying authentication service
	// cannot retrieve authentication information.
	ErrorAuthInfoUnavailable
	// ErrorUserUnknown indicates the user is not known to the underlying
	// authentication module.
	ErrorUserUnknown
	// ErrorMaxTries indicates an authentication service has
	// maintained a retry count which has
	// been reached. No further retries
	// should be attempted.
	ErrorMaxTries
	// ErrorNewAuthTokenRequired indicates a new authentication token is required.
	// This is normally returned if the
	// machine security policies require
	// that the password should be changed
	// because the password is NULL or it
	// has aged.
	ErrorNewAuthTokenRequired
	// ErrorAccountExpired indicates the user account has expired.
	ErrorAccountExpired
	// ErrorSession indicates inability to make/remove an entry for
	// the specified session.
	ErrorSession
	// ErrorCredentialsUnavailable indicates an underlying authentication service
	// cannot retrieve user credentials.
	ErrorCredentialsUnavailable
	// ErrorCredentialsExpired indicates user credentials are expired.
	ErrorCredentialsExpired
	// ErrorCredentials indicates a failure setting user credentials.
	ErrorCredentials
	// ErrorNoModuleData indicates no module-specific data is present.
	ErrorNoModuleData
	// ErrorConversation indicates a conversation error.
	ErrorConversation
	// ErrorAuthToken indicates an authentication token manipulation error.
	ErrorAuthToken
	// ErrorAuthTokenRecovery indicates authentication information
	// cannot be recovered.
	ErrorAuthTokenRecovery
	// ErrorAuthTokenLockBusy indicates an authentication token lock is busy.
	ErrorAuthTokenLockBusy
	// ErrorAuthTokenDisableAging indicates authentication token aging is disabled.
	ErrorAuthTokenDisableAging
	// ErrorTryAgain indicates a preliminary check by the password service.
	ErrorTryAgain
	// ErrorIgnore indicates an instruction to ignore the underlying account module
	// regardless of whether the control
	// flag is required, optional, or sufficient.
	ErrorIgnore
	// ErrorAbort indicates a critical error (?module fail now request).
	ErrorAbort
	// ErrorAuthTokenExpired indicates a user's authentication token has expired.
	ErrorAuthTokenExpired
	// ErrorModuleUnknown indicates a module is not known.
	ErrorModuleUnknown
	// ErrorBadItem indicates a bad item passed to pam_*_item().
	ErrorBadItem
	// ErrorConversationAgain indicates a conversation function is event driven
	// and data is not available yet.
	ErrorConversationAgain
	// ErrorIncomplete indicatesa an instruction to please call this function again to
	// complete authentication stack. Before
	// calling again, verify that conversation
	// is completed.
	ErrorIncomplete
)

// Error implements the error interface.
func (err Error) Error() string {
	return fmt.Sprintf("%d", int(err))
}
