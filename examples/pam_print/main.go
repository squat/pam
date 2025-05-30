package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/squat/pam"
	pamc "github.com/squat/pam/c"
)

var flagSet = flag.NewFlagSet("pam_print", flag.ContinueOnError)

func init() {
	flagSet.String("output", "text", `The desired output format; one of "text" or "json"`)
	pamc.Register(new(print))
}

var _ pam.ServiceModule = (*print)(nil)

// print is a ServiceModule that prints the args passed to it.
type print struct{}

func (p *print) print(handle pam.Handle, flags int, args []string) {
	var output string
	if err := flagSet.Parse(args); err != nil {
		fmt.Printf("failed to parse arguments: %v", err)
		return
	}
	user, err := handle.GetUser("")
	if err != nil {
		fmt.Printf("failed to get user: %v\n", err)
		return
	}
	env, err := handle.GetEnvList()
	if err != nil {
		fmt.Printf("failed to get environment variables: %v\n", err)
		return
	}
	if output == "json" {
		out := struct {
			User  string            `json:"user"`
			Flags int               `json:"flags"`
			Args  []string          `json:"args"`
			Env   map[string]string `json:"env"`
		}{
			User:  user,
			Flags: flags,
			Args:  args,
			Env:   env,
		}

		b, err := json.Marshal(out)
		if err != nil {
			fmt.Printf("failed to marshal json: %v\n", err)
			return
		}
		println(string(b))
		return
	}

	fmt.Printf("User:\n")
	fmt.Printf("\t%s\n", user)
	fmt.Printf("Flags:\n\t%d\n", flags)
	fmt.Printf("Arguments:\n\t%s\n", strings.Join(args, "\n\t"))
	fmt.Printf("Environment Variables:\n")
	for k, v := range env {
		fmt.Printf("\t%s=%s\n", k, v)
	}

}

// Authenticate is ignored.
func (p *print) Authenticate(handle pam.Handle, flags int, args []string) error {
	p.print(handle, flags, args)
	return pam.ErrorIgnore
}

// SetCredentials is ignored.
func (p *print) SetCredentials(handle pam.Handle, flags int, args []string) error {
	p.print(handle, flags, args)
	return pam.ErrorIgnore
}

// AccountManagement is ignored.
func (p *print) AccountManagement(handle pam.Handle, flags int, args []string) error {
	p.print(handle, flags, args)
	return pam.ErrorIgnore
}

// OpenSession is ignored.
func (p *print) OpenSession(handle pam.Handle, flags int, args []string) error {
	p.print(handle, flags, args)
	return pam.ErrorIgnore
}

// CloseSession is ignored.
func (p *print) CloseSession(handle pam.Handle, flags int, args []string) error {
	p.print(handle, flags, args)
	return pam.ErrorIgnore
}

// ChangeAuthToken is ignored.
func (p *print) ChangeAuthToken(handle pam.Handle, flags int, args []string) error {
	p.print(handle, flags, args)
	return pam.ErrorIgnore
}

func main() {
	flagSet.Usage()
}
