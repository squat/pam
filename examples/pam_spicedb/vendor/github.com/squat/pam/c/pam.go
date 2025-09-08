package pamc

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/squat/pam"
)

/*
#cgo LDFLAGS: -lpam
#include <security/pam_modules.h>
#include <stdlib.h>
typedef const char cchar_t;
typedef const void cvoid_t;
*/
import "C"

var serviceModule pam.ServiceModule = &pam.Nop{}

func Register(sm pam.ServiceModule) {
	serviceModule = sm
}

func next(p **C.char) **C.char {
	return (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(p)))
}

func unpackCharPointerPointer(c C.int, v **C.char) []string {
	if c == 0 || v == nil {
		return nil
	}
	args := make([]string, 0, c)
	for p := v; *p != nil; p = next(p) {
		args = append(args, C.GoString(*p))
	}
	return args
}

type handle struct {
	pamh *C.pam_handle_t
}

func (h *handle) SetItem(itemType pam.ItemType, item unsafe.Pointer) error {
	if err := pam.Error(C.pam_set_item(h.pamh, C.int(itemType), item)); err != 0 {
		return err
	}
	return nil
}

func (h *handle) GetItem(itemType pam.ItemType) (unsafe.Pointer, error) {
	var item unsafe.Pointer
	if err := pam.Error(C.pam_get_item(h.pamh, C.int(itemType), &item)); err != 0 {
		return nil, err
	}
	return item, nil
}

func (h *handle) GetUser(prompt string) (string, error) {
	var user *C.cchar_t
	if err := pam.Error(C.pam_get_user(h.pamh, &user, C.CString(prompt))); err != 0 {
		return "", err
	}
	return C.GoString(user), nil
}

func (h *handle) PutEnv(key, value string) error {
	if err := pam.Error(C.pam_putenv(h.pamh, C.CString(fmt.Sprintf("%s=%s", key, value)))); err != 0 {
		return err
	}
	return nil
}

func (h *handle) GetEnv(key string) string {
	v := C.pam_getenv(h.pamh, C.CString(key))
	if v == nil {
		return ""
	}
	return C.GoString(v)
}

func (h *handle) GetEnvList() (map[string]string, error) {
	env := make(map[string]string)
	p := C.pam_getenvlist(h.pamh)
	if p == nil {
		return nil, pam.ErrorBuffer
	}
	for q := p; *q != nil; q = next(q) {
		parts := strings.SplitN(C.GoString(*q), "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
		C.free(unsafe.Pointer(*q))
	}
	C.free(unsafe.Pointer(p))
	return env, nil
}

//export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.cchar_t) C.int {
	args := unpackCharPointerPointer(argc, argv)
	if err := serviceModule.Authenticate(&handle{pamh}, int(flags), args); err != nil {
		if perr, ok := err.(pam.Error); ok {
			return C.int(perr)
		}
	}
	return C.PAM_SUCCESS
}

//export pam_sm_setcred
func pam_sm_setcred(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.cchar_t) C.int {
	args := unpackCharPointerPointer(argc, argv)
	if err := serviceModule.SetCredentials(&handle{pamh}, int(flags), args); err != nil {
		if perr, ok := err.(pam.Error); ok {
			return C.int(perr)
		}
	}
	return C.PAM_SUCCESS
}

//export pam_sm_acct_mgmt
func pam_sm_acct_mgmt(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.cchar_t) C.int {
	args := unpackCharPointerPointer(argc, argv)
	if err := serviceModule.AccountManagement(&handle{pamh}, int(flags), args); err != nil {
		if perr, ok := err.(pam.Error); ok {
			return C.int(perr)
		}
	}
	return C.PAM_SUCCESS
}

//export pam_sm_open_session
func pam_sm_open_session(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.cchar_t) C.int {
	args := unpackCharPointerPointer(argc, argv)
	if err := serviceModule.OpenSession(&handle{pamh}, int(flags), args); err != nil {
		if perr, ok := err.(pam.Error); ok {
			return C.int(perr)
		}
	}
	return C.PAM_SUCCESS
}

//export pam_sm_close_session
func pam_sm_close_session(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.cchar_t) C.int {
	args := unpackCharPointerPointer(argc, argv)
	if err := serviceModule.CloseSession(&handle{pamh}, int(flags), args); err != nil {
		if perr, ok := err.(pam.Error); ok {
			return C.int(perr)
		}
	}
	return C.PAM_SUCCESS
}

//export pam_sm_chauthtok
func pam_sm_chauthtok(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.cchar_t) C.int {
	args := unpackCharPointerPointer(argc, argv)
	if err := serviceModule.ChangeAuthToken(&handle{pamh}, int(flags), args); err != nil {
		if perr, ok := err.(pam.Error); ok {
			return C.int(perr)
		}
	}
	return C.PAM_SUCCESS
}
