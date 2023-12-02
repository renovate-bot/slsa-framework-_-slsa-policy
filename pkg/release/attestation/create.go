package attestation

import (
	"encoding/json"
	"fmt"

	"github.com/laurentsimon/slsa-policy/pkg/errs"

	"github.com/laurentsimon/slsa-policy/pkg/utils/intoto"
)

type Creation struct {
	attestation
	safeMode bool
}

type CreationOption func(*Creation) error

// NOTE: See https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
func CreationNew(subject intoto.Subject, authorID string, result ReleaseResult, options ...CreationOption) (*Creation, error) {
	if err := subject.Validate(); err != nil {
		return nil, err
	}
	// Validate the digests.
	att := Creation{
		attestation: attestation{
			Header: intoto.Header{
				Type:          statementType,
				PredicateType: predicateType,
				Subjects:      []intoto.Subject{subject},
			},
			Predicate: predicate{
				ReleaseResult: result,
				CreationTime:  intoto.Now(),
				Author: intoto.Author{
					ID: authorID,
				},
			},
		},
	}
	for _, option := range options {
		err := option(&att)
		if err != nil {
			return nil, err
		}
	}
	return &att, nil
}

func (a *Creation) ToBytes() ([]byte, error) {
	content, err := json.Marshal(*&a.attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %v", err)
	}
	return content, nil
}

// TODO: unit tests.
func SetSafeMode() CreationOption {
	return func(a *Creation) error {
		return a.setSafeMode()
	}
}

func (a *Creation) setSafeMode() error {
	a.safeMode = true
	return nil
}

func (a *Creation) isSafeMode() bool {
	return a.safeMode
}

func SetAuthorVersion(version string) CreationOption {
	return func(a *Creation) error {
		return a.setAuthorVersion(version)
	}
}

func (a *Creation) setAuthorVersion(version string) error {
	a.attestation.Predicate.Author.Version = version
	return nil
}

func SetEnvironment(env string) CreationOption {
	return func(a *Creation) error {
		return a.setEnvironment(env)
	}
}

func (a *Creation) setEnvironment(env string) error {
	if a.isSafeMode() {
		return fmt.Errorf("%w: safe mode enabled, cannot edit environment", errs.ErrorInvalidInput)
	}
	if a.attestation.Header.Subjects[0].Annotations == nil {
		a.attestation.Header.Subjects[0].Annotations = make(map[string]interface{})
	}
	a.attestation.Header.Subjects[0].Annotations[environmentAnnotation] = env
	return nil
}

func SetPolicy(policy map[string]intoto.Policy) CreationOption {
	return func(a *Creation) error {
		return a.setPolicy(policy)
	}
}

func (a *Creation) setPolicy(policy map[string]intoto.Policy) error {
	a.attestation.Predicate.Policy = policy
	return nil
}

func SetSlsaBuildLevel(level int) CreationOption {
	return func(a *Creation) error {
		return a.setSlsaBuildLevel(level)
	}
}

func (a *Creation) setSlsaBuildLevel(level int) error {
	if a.isSafeMode() {
		return fmt.Errorf("%w: safe mode enabled, cannot edit SLSA build level", errs.ErrorInvalidInput)
	}
	if !a.isResultAllowed() {
		return fmt.Errorf("%w: level cannot be set for %q result", errs.ErrorInvalidInput, a.attestation.Predicate.ReleaseResult)
	}
	if level < 0 {
		return fmt.Errorf("%w: level (%v) is negative", errs.ErrorInvalidInput, level)
	}
	if level > 4 {
		return fmt.Errorf("%w: level (%v) is too large", errs.ErrorInvalidInput, level)
	}
	if a.attestation.Predicate.ReleaseProperties == nil {
		a.attestation.Predicate.ReleaseProperties = make(map[string]interface{})
	}
	a.attestation.Predicate.ReleaseProperties[buildLevelProperty] = level
	return nil
}

func (a *Creation) isResultAllowed() bool {
	return a.attestation.Predicate.ReleaseResult == ReleaseResultAllow
}