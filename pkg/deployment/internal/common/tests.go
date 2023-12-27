package common

import (
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/laurentsimon/slsa-policy/pkg/deployment/internal/options"
	"github.com/laurentsimon/slsa-policy/pkg/errs"
	"github.com/laurentsimon/slsa-policy/pkg/utils/intoto"
	"github.com/laurentsimon/slsa-policy/pkg/utils/iterator"
)

func AsPointer[K interface{}](o K) *K {
	return &o
}

// Bytes iterator.
func NewNamedBytesIterator(values [][]byte, uniqueID bool) iterator.NamedReadCloserIterator {
	return &bytesIterator{values: values, index: -1, uniqueID: uniqueID}
}

type bytesIterator struct {
	values   [][]byte
	index    int
	uniqueID bool
	err      error
}

func (iter *bytesIterator) Next() (string, io.ReadCloser) {
	if iter.err != nil {
		return "", nil
	}
	iter.index++
	if iter.uniqueID {
		return fmt.Sprintf("policy_id%d", iter.index), io.NopCloser(bytes.NewReader(iter.values[iter.index]))
	}
	return fmt.Sprintf("policy_id%d", 0), io.NopCloser(bytes.NewReader(iter.values[iter.index]))
}

func (iter *bytesIterator) HasNext() bool {
	if iter.err != nil {
		return false
	}
	return iter.index+1 < len(iter.values)
}

func (iter *bytesIterator) Error() error {
	return nil
}

// Attestation verifier.
func NewAttestationVerifier(digests intoto.DigestSet, packageURI, env, releaserID string) options.AttestationVerifier {
	return &attesationVerifier{digests: digests, packageURI: packageURI, releaserID: releaserID, env: env}
}

type attesationVerifier struct {
	packageURI string
	releaserID string
	env        string
	digests    intoto.DigestSet
}

func (v *attesationVerifier) VerifyReleaseAttestation(digests intoto.DigestSet, packageURI string, env []string, releaserID string) (*string, error) {
	if packageURI == v.packageURI && releaserID == v.releaserID &&
		mapEq(digests, v.digests) &&
		((v.env != "" && len(env) > 0 && slices.Contains(env, v.env)) ||
			(v.env == "" && len(env) == 0)) {
		if v.env == "" {
			return nil, nil
		}
		return &v.env, nil
	}
	return nil, fmt.Errorf("%w: cannot verify package URI (%q) releaser ID (%q) env (%q)", errs.ErrorVerification, packageURI, releaserID, env)
}

func mapEq(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		vv, exists := m2[k]
		if !exists {
			return false
		}
		if vv != v {
			return false
		}
	}
	return true
}