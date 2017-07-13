package stack

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/jetstack/tarmak/pkg/tarmak/config"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

type StateStack struct {
	*Stack
}

var _ interfaces.Stack = &StateStack{}

func newStateStack(s *Stack, conf *config.StackState) (*StateStack, error) {
	s.name = config.StackNameState
	return &StateStack{
		Stack: s,
	}, nil
}

func (s *StateStack) Variables() map[string]interface{} {
	output := map[string]interface{}{}
	state := s.Stack.conf.State
	if state.BucketPrefix != "" {
		output["bucket_prefix"] = state.BucketPrefix
	}
	if state.PublicZone != "" {
		output["public_zone"] = state.PublicZone
	}

	return output
}

func (s *StateStack) VerifyPost() error {
	return s.verifyDNSDelegation()
}

func (s *StateStack) verifyDNSDelegation() error {

	tries := 5
	for {
		host := strings.Join([]string{utils.RandStringRunes(16), "_tarmak", s.conf.State.PublicZone}, ".")

		result, err := net.LookupTXT(host)
		if err == nil {
			if reflect.DeepEqual([]string{"tarmak delegation works"}, result) {
				return nil
			} else {
				s.log.Warn("error checking delegation to public zone: ", err)
			}
		} else {
			s.log.Warn("error checking delegation to public zone: ", err)
		}

		if tries == 0 {
			return fmt.Errorf("Failed 5 times")
		}
		tries -= 1
		time.Sleep(2 * time.Second)
	}
}
