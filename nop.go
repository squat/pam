package pam

var _ ServiceModule = (*Nop)(nil)

// Nop is a ServiceModule that instructs PAM to ignore it.
type Nop struct{}

// Authenticate is ignored.
func (n *Nop) Authenticate(handle Handle, flags int, args []string) error {
	return ErrorIgnore
}

// SetCredentials is ignored.
func (n *Nop) SetCredentials(handle Handle, flags int, args []string) error {
	return ErrorIgnore
}

// AccountManagement is ignored.
func (n *Nop) AccountManagement(handle Handle, flags int, args []string) error {
	return ErrorIgnore
}

// OpenSession is ignored.
func (n *Nop) OpenSession(handle Handle, flags int, args []string) error {
	return ErrorIgnore
}

// CloseSession is ignored.
func (n *Nop) CloseSession(handle Handle, flags int, args []string) error {
	return ErrorIgnore
}

// ChangeAuthToken is ignored.
func (n *Nop) ChangeAuthToken(handle Handle, flags int, args []string) error {
	return ErrorIgnore
}
