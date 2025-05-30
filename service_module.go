package pam

// ServiceModule is an implementation of a PAM service module.
type ServiceModule interface {
	// Authentication APIs

	// Authenticate performs the task of authenticating the user.
	// See pam_sm_authenticate(3).
	Authenticate(handle Handle, flags int, args []string) error

	// SetCredentials performs the task of altering the credentials
	// of the user with respect to the corresponding authorization scheme.
	// Generally, an authentication module may have access to more
	// information about a user than their authentication token.
	// This function is used to make such information available to the application.
	// It should only be called after the user has been
	// authenticated but before a session has been established.
	// See pam_sm_setcred(3).
	SetCredentials(handle Handle, flags int, args []string) error

	// Account Management APIs

	// AccountManagement performs the task of establishing whether
	// the user is permitted to gain access at this time.
	// It should be understood that the user has previously
	// been validated by an authentication module.
	// See pam_sm_acct_mgmt(3).
	AccountManagement(handle Handle, flags int, args []string) error

	// Session Management APIs

	// OpenSession is called to commence a session.
	// See pam_sm_open_session(3).
	OpenSession(handle Handle, flags int, args []string) error

	// CloseSession is called to terminate a session.
	// See pam_sm_close_session(3).
	CloseSession(handle Handle, flags int, args []string) error

	// Password Management APIs

	// ChangeAuthToken is called to (re-)set the authentication token of the user.
	// See pam_sm_chauthtok(3).
	ChangeAuthToken(handle Handle, flags int, args []string) error
}
