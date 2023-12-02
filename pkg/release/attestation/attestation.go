package attestation

import (
	"github.com/laurentsimon/slsa-policy/pkg/utils/intoto"
)

type predicate struct {
	Author               intoto.Author            `json:"author"`
	CreationTime         string                   `json:"creationTime"`
	Policy               map[string]intoto.Policy `json:"policy,omitempty"`
	ReleaseResult        ReleaseResult            `json:"releaseResult"`
	ReleaseProperties    properties               `json:"releaseProperties,omitempty"`
	DependencyProperties map[string]properties    `json:"dependencyProperties,omitempty"`
}

type attestation struct {
	intoto.Header
	Predicate predicate `json:"predicate"`
}

type properties map[string]interface{}

const (
	statementType      = "https://in-toto.io/Statement/v1"
	predicateType      = "https://slsa.dev/release_attestation/v1"
	buildLevelProperty = "slsa.dev/build/level"
	// TODO: make these public for users to be able to construct subjects.
	environmentAnnotation = "environment"
	versionAnnotation     = "version"
)

type ReleaseResult string

const (
	ReleaseResultAllow ReleaseResult = "ALLOW"
	ReleaseResultDeny  ReleaseResult = "DENY"
)