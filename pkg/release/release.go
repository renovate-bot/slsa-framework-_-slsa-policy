package release

import (
	"fmt"
	"io"

	"github.com/laurentsimon/slsa-policy/pkg/release/internal"
	"github.com/laurentsimon/slsa-policy/pkg/utils/iterator"
)

// Policy defines the release policy.
type Policy struct {
	policy *internal.Policy
}

// New creates a release policy.
func New(org io.Reader, projects iterator.ReaderIterator) (*Policy, error) {
	policy, err := internal.New(org, projects)
	if err != nil {
		return nil, err
	}
	return &Policy{
		policy: policy,
	}, nil
}

func Hello() error {
	fmt.Println("hey")
	return nil
}

// TODO:
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// Evaluate evaluates the policy.
// func (p *Policy) Evaluate(options ...Options) error {
// 	return p.policy.Evaluate(sourceURI, imageURI, builderID)
// }

// func (p *Policy) Evaluate() results.Verification {
// 	return p.policy.Evaluate(sourceURI, imageURI, builderID)
// }
