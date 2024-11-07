package verification

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sigstore/sigstore-go/pkg/fulcio/certificate"
)

var (
	GitHubOIDCIssuer       = "https://token.actions.githubusercontent.com"
	GitHubTenantOIDCIssuer = "https://token.actions.%s.ghe.com"
)

func VerifyCertExtensions(results []*AttestationProcessingResult, ec EnforcementCriteria) error {
	if len(results) == 0 {
		return errors.New("no attestations proccessing results")
	}

	var lastErr error
	for _, attestation := range results {
		err := verifyCertExtensions(*attestation.VerificationResult.Signature.Certificate, ec.Certificate)
		if err == nil {
			// if at least one attestation is verified, we're good as verification
			// is defined as successful if at least one attestation is verified
			return nil
		}
		lastErr = err
	}

	// if we have exited the for loop without returning early due to successful
	// verification, we need to return an error
	return lastErr
}

func verifyCertExtensions(verified, expected certificate.Summary) error {
	if !strings.EqualFold(expected.SourceRepositoryOwnerURI, verified.SourceRepositoryOwnerURI) {
		return fmt.Errorf("expected SourceRepositoryOwnerURI to be %s, got %s", expected.SourceRepositoryOwnerURI, verified.SourceRepositoryOwnerURI)
	}

	// if repo is set, compare the SourceRepositoryURI fields
	if expected.SourceRepositoryURI != "" && !strings.EqualFold(expected.SourceRepositoryURI, verified.SourceRepositoryURI) {
		return fmt.Errorf("expected SourceRepositoryURI to be %s, got %s", expected.SourceRepositoryURI, verified.SourceRepositoryURI)
	}

	// compare the OIDC issuers. If not equal, return an error depending
	// on if there is a partial match
	if !strings.EqualFold(expected.Issuer, verified.Issuer) {
		if strings.Index(verified.Issuer, expected.Issuer+"/") == 0 {
			return fmt.Errorf("expected Issuer to be %s, got %s -- if you have a custom OIDC issuer policy for your enterprise, use the --cert-oidc-issuer flag with your expected issuer", expected.Issuer, verified.Issuer)
		}
		return fmt.Errorf("expected Issuer to be %s, got %s", expected.Issuer, verified.Issuer)
	}

	return nil
}
