package publish

import (
	"fmt"
	"io"

	"github.com/slsa-framework/slsa-policy/pkg/errs"
	"github.com/slsa-framework/slsa-policy/pkg/publish/internal"
	"github.com/slsa-framework/slsa-policy/pkg/publish/internal/options"
	"github.com/slsa-framework/slsa-policy/pkg/utils/intoto"
	"github.com/slsa-framework/slsa-policy/pkg/utils/iterator"
)

// AttestationVerifier defines an interface to verify attestations.
type AttestationVerifier interface {
	// Build attestation verification.
	VerifyBuildAttestation(digests intoto.DigestSet, policyPackageName, builderID, sourceURI string) error
}

// AttestationVerificationOption defines the configuration to verify
// build attestations.
type AttestationVerificationOption struct {
	Verifier AttestationVerifier
	// We can add attestation-specific options here.
}

// RequestOption contains options from the caller.
type RequestOption struct {
	Environment *string
}

// Policy defines the publish policy.
type Policy struct {
	policy        *internal.Policy
	validator     options.PolicyValidator
	packageHelper PackageHelper
}

// PolicyOption defines a policy option.
type PolicyOption func(*Policy) error

// This is a helpder class to forward calls between the internal
// classes and the caller.
type internal_verifier struct {
	opts AttestationVerificationOption
}

func (i *internal_verifier) VerifyBuildAttestation(digests intoto.DigestSet, policyPackageName, builderID, sourceURI string) error {
	if i.opts.Verifier == nil {
		return fmt.Errorf("%w: verifier is nil", errs.ErrorInvalidInput)
	}
	return i.opts.Verifier.VerifyBuildAttestation(digests, policyPackageName, builderID, sourceURI)
}

// This is a class to forward calls between internal
// classes and the caller for the PolicyValidator interface.
type internal_validator struct {
	validator PolicyValidator
}

func (i *internal_validator) ValidatePackage(pkg options.ValidationPackage) error {
	if i.validator == nil {
		return nil
	}
	return i.validator.ValidatePackage(ValidationPackage{
		Name: pkg.Name,
		Environment: ValidationEnvironment{
			// NOTE: make a copy of the array.
			AnyOf: append([]string{}, pkg.Environment.AnyOf...),
		},
	})
}

// New creates a publish policy.
func PolicyNew(org io.ReadCloser, projects iterator.ReadCloserIterator, packageHelper PackageHelper, opts ...PolicyOption) (*Policy, error) {
	// Initialize a policy with caller options.
	p := new(Policy)
	for _, option := range opts {
		err := option(p)
		if err != nil {
			return nil, err
		}
	}
	policy, err := internal.PolicyNew(org, projects, p.validator)
	if err != nil {
		return nil, err
	}
	p.policy = policy
	if packageHelper == nil {
		return nil, fmt.Errorf("%w: package hepler is nil", errs.ErrorInvalidInput)
	}
	p.packageHelper = packageHelper
	return p, nil
}

// SetValidator sets a custom validator.
func SetValidator(validator PolicyValidator) PolicyOption {
	return func(p *Policy) error {
		return p.setValidator(validator)
	}
}

func (p *Policy) setValidator(validator PolicyValidator) error {
	// Construct an internal validator
	// using the caller's public validator interface.
	p.validator = &internal_validator{
		validator: validator,
	}
	return nil
}

// Evaluate evalues the publish policy.
func (p *Policy) Evaluate(digests intoto.DigestSet, policyPackageName string, reqOpts RequestOption,
	opts AttestationVerificationOption) PolicyEvaluationResult {
	level, err := p.policy.Evaluate(digests, policyPackageName,
		options.Request{
			Environment: reqOpts.Environment,
		},
		options.BuildVerification{
			Verifier: &internal_verifier{
				opts: opts,
			},
		},
	)
	if err != nil {
		return PolicyEvaluationResult{
			err:       err,
			evaluated: true,
		}
	}

	// Translate the policy package names to a package descriptor.
	packageDesc, err := p.packageHelper.PackageDescriptor(policyPackageName)
	if err != nil {
		return PolicyEvaluationResult{
			err:       err,
			evaluated: true,
		}
	}
	return PolicyEvaluationResult{
		level:       level,
		err:         err,
		packageDesc: packageDesc,
		digests:     digests,
		environment: reqOpts.Environment,
		evaluated:   true,
	}
}

// Utility function for cosign integration.
func PredicateType() string {
	return predicateType
}
