package pam

import "unsafe"

type Handle interface {
	// SetItem allows applications and PAM service modules
	// to update PAM information.
	// See pam_set_item(3).
	SetItem(ItemType, unsafe.Pointer) error

	// GetItem allows applications and PAM service modules
	// to retrieve PAM information.
	// See pam_get_item(3).
	GetItem(ItemType) (unsafe.Pointer, error)

	// GetUser returns the name of the user specified by `pam_start`
	// if this is NULL, the username is obtained via the `pam_conv`
	// mechanism.
	// See pam_get_user(3).
	GetUser(prompt string) (string, error)

	// PutEnv is used to add or change the value of PAM environment variables.
	// See pam_putenv(3).
	PutEnv(key, value string) error

	// GetEnv searches the PAM environment list for an item
	// that matches the key and returns the value.
	// See pam_getenv(3).
	GetEnv(key string) string

	// GetEnvList returns a complete copy of the PAM environment.
	// See pam_get_envlist(3).
	GetEnvList() (map[string]string, error)
}
