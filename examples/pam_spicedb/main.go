package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/squat/pam"
	pamc "github.com/squat/pam/c"
)

var flagSet = flag.NewFlagSet("pam_spicedb", flag.ContinueOnError)

func init() {
	flagSet.Bool("tls", true, "Whether to enable TLS")
	flagSet.Bool("insecure-skip-verify", false, "Whether to skip TLS verification")
	flagSet.String("endpoint", "", "The SpiceDB URL")
	flagSet.String("token-file", "", "Path to a file containing the SpiceDB token")
	flagSet.String("permission", "", "SpiceDB permission")
	flagSet.String("resource-type", "", "SpiceDB resource type")
	flagSet.String("resource-id", "", "SpiceDB resource ID")
	flagSet.String("subject-type", "", "SpiceDB subject type")
	pamc.Register(new(spicedb))
}

var _ pam.ServiceModule = (*spicedb)(nil)

// spicedb is a ServiceModule that validates user accounts against SpiceDB.
type spicedb struct {
	client *authzed.Client
	lock   sync.Mutex
}

func (s *spicedb) connect(args []string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.client != nil {
		return nil
	}

	if err := flagSet.Parse(args); err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	endpoint := flagSet.Lookup("endpoint").Value.String()
	tokenFile := flagSet.Lookup("token-file").Value.String()
	tokenBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return fmt.Errorf("failed to read token file: %w", err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	opts := []grpc.DialOption{}
	tls := flagSet.Lookup("tls").Value.(flag.Getter).Get().(bool)
	if !tls {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		opts = append(opts, grpcutil.WithInsecureBearerToken(token))
	} else {
		verify := lo.Ternary(
			flagSet.Lookup("insecure-skip-verify").Value.(flag.Getter).Get().(bool),
			grpcutil.SkipVerifyCA,
			grpcutil.VerifyCA,
		)

		certsOpt, err := grpcutil.WithSystemCerts(verify)
		if err != nil {
			return err
		}

		opts = append(opts, certsOpt)
		opts = append(opts, grpcutil.WithBearerToken(token))
	}

	client, err := authzed.NewClient(endpoint, opts...)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	s.client = client
	return nil
}

// Authenticate is ignored.
func (p *spicedb) Authenticate(handle pam.Handle, flags int, args []string) error {
	return pam.ErrorIgnore
}

// SetCredentials is ignored.
func (p *spicedb) SetCredentials(handle pam.Handle, flags int, args []string) error {
	return pam.ErrorIgnore
}

// AccountManagement is checks if the user has the specified relation in SpiceDB.
func (p *spicedb) AccountManagement(handle pam.Handle, flags int, args []string) error {
	if err := p.connect(args); err != nil {
		return pam.ErrorService
	}

	sshAuthInfo := handle.GetEnv("SSH_AUTH_INFO_0")
	if sshAuthInfo == "" {
		return pam.ErrorAuth
	}
	sshAuthInfoLines := strings.Split(sshAuthInfo, "\n")
	ids := make([]string, 0, len(sshAuthInfoLines))
	for _, sshAuthInfoLine := range sshAuthInfoLines {
		parts := strings.Split(sshAuthInfoLine, " ")
		if len(parts) != 3 {
			continue
		}
		ids = append(ids, parts[2])
	}
	for _, id := range ids {
		response, err := p.client.CheckPermission(context.Background(), &v1.CheckPermissionRequest{
			Consistency: &v1.Consistency{Requirement: &v1.Consistency_FullyConsistent{FullyConsistent: true}},
			Permission:  flagSet.Lookup("permission").Value.String(),
			Resource: &v1.ObjectReference{
				ObjectId:   flagSet.Lookup("resource-id").Value.String(),
				ObjectType: flagSet.Lookup("resource-type").Value.String(),
			},
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectId:   id,
					ObjectType: flagSet.Lookup("subject-type").Value.String(),
				},
			},
		})
		if err != nil {
			return pam.ErrorService
		}
		if err := response.Validate(); err == nil && response.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
			return nil
		}
	}
	return pam.ErrorAuth
}

// OpenSession is ignored.
func (p *spicedb) OpenSession(handle pam.Handle, flags int, args []string) error {
	return pam.ErrorIgnore
}

// CloseSession is ignored.
func (p *spicedb) CloseSession(handle pam.Handle, flags int, args []string) error {
	return pam.ErrorIgnore
}

// ChangeAuthToken is ignored.
func (p *spicedb) ChangeAuthToken(handle pam.Handle, flags int, args []string) error {
	return pam.ErrorIgnore
}

func main() {
	flagSet.Usage()
}
