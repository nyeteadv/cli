package set

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/cmd/secret/shared"
)

type SecretPayload struct {
	EncryptedValue string   `json:"encrypted_value"`
	Visibility     string   `json:"visibility,omitempty"`
	Repositories   []string `json:"selected_repository_ids,omitempty"`
	KeyID          string   `json:"key_id"`
}

type PubKey struct {
	Raw [32]byte
	ID  string `json:"key_id"`
	Key string
}

type repositoryNWO struct {
	owner string
	name  string
}

func getPubKey(client *api.Client, host, path string) (*PubKey, error) {
	pk := PubKey{}
	err := client.REST(host, "GET", path, nil, &pk)
	if err != nil {
		return nil, err
	}

	if pk.Key == "" {
		return nil, fmt.Errorf("failed to find public key at %s/%s", host, path)
	}

	decoded, err := base64.StdEncoding.DecodeString(pk.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	copy(pk.Raw[:], decoded[0:32])
	return &pk, nil
}

func getOrgPublicKey(client *api.Client, host, orgName string) (*PubKey, error) {
	return getPubKey(client, host, fmt.Sprintf("orgs/%s/actions/secrets/public-key", orgName))
}

func getUserPublicKey(client *api.Client, host string) (*PubKey, error) {
	return getPubKey(client, host, "user/codespaces/secrets/public-key")
}

func getRepoPubKey(client *api.Client, repo ghrepo.Interface) (*PubKey, error) {
	return getPubKey(client, repo.RepoHost(), fmt.Sprintf("repos/%s/actions/secrets/public-key",
		ghrepo.FullName(repo)))
}

func getEnvPubKey(client *api.Client, repo ghrepo.Interface, envName string) (*PubKey, error) {
	return getPubKey(client, repo.RepoHost(), fmt.Sprintf("repos/%s/environments/%s/secrets/public-key",
		ghrepo.FullName(repo), envName))
}

func putSecret(client *api.Client, host, path string, payload SecretPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize: %w", err)
	}
	requestBody := bytes.NewReader(payloadBytes)

	return client.REST(host, "PUT", path, requestBody, nil)
}

func putOrgSecret(client *api.Client, host string, pk *PubKey, opts SetOptions, eValue string) error {
	secretName := opts.SecretName
	orgName := opts.OrgName
	visibility := opts.Visibility

	var repositoryIDs []string
	var err error
	if orgName != "" && visibility == shared.Selected {
		repos := make([]ghrepo.Interface, 0, len(opts.RepositoryNames))
		for _, repo := range opts.RepositoryNames {
			repos = append(repos, ghrepo.New(opts.OrgName, repo))
		}
		repositoryIDs, err = mapRepoToID(client, host, repos)
		if err != nil {
			return fmt.Errorf("failed to look up IDs for repositories %v: %w", opts.RepositoryNames, err)
		}
	}

	payload := SecretPayload{
		EncryptedValue: eValue,
		KeyID:          pk.ID,
		Repositories:   repositoryIDs,
		Visibility:     visibility,
	}
	path := fmt.Sprintf("orgs/%s/actions/secrets/%s", orgName, secretName)

	return putSecret(client, host, path, payload)
}

func putUserSecret(client *api.Client, host string, pk *PubKey, opts SetOptions, eValue string) error {
	secretName := opts.SecretName
	visibility := opts.Visibility

	var repositoryIDs []string
	var err error
	if visibility == shared.Selected {
		repos := make([]ghrepo.Interface, 0, len(opts.RepositoryNames))
		for _, repo := range opts.RepositoryNames {
			// For user secrets, repository names should be fully qualifed (e.g. "owner/repo")
			repoNWO, err := ghrepo.FromFullName(repo)
			if err != nil {
				return err
			}
			repos = append(repos, repoNWO)
		}

		repositoryIDs, err = mapRepoToID(client, host, repos)
		if err != nil {
			return fmt.Errorf("failed to look up IDs for repositories %v: %w", opts.RepositoryNames, err)
		}
	}

	payload := SecretPayload{
		EncryptedValue: eValue,
		KeyID:          pk.ID,
		Repositories:   repositoryIDs,
		Visibility:     visibility,
	}
	path := fmt.Sprintf("user/codespaces/secrets/%s", secretName)
	return putSecret(client, host, path, payload)
}

func putEnvSecret(client *api.Client, pk *PubKey, repo ghrepo.Interface, envName string, secretName, eValue string) error {
	payload := SecretPayload{
		EncryptedValue: eValue,
		KeyID:          pk.ID,
	}
	path := fmt.Sprintf("repos/%s/environments/%s/secrets/%s", ghrepo.FullName(repo), envName, secretName)
	return putSecret(client, repo.RepoHost(), path, payload)
}

func putRepoSecret(client *api.Client, pk *PubKey, repo ghrepo.Interface, secretName, eValue string) error {
	payload := SecretPayload{
		EncryptedValue: eValue,
		KeyID:          pk.ID,
	}
	path := fmt.Sprintf("repos/%s/actions/secrets/%s", ghrepo.FullName(repo), secretName)
	return putSecret(client, repo.RepoHost(), path, payload)
}

// This does similar logic to `api.RepoNetwork`, but without the overfetching.
func mapRepoToID(client *api.Client, host string, repositories []ghrepo.Interface) ([]string, error) {
	queries := make([]string, 0, len(repositories))
	for i, repo := range repositories {
		queries = append(queries, fmt.Sprintf(`
			repo_%03d: repository(owner: %q, name: %q) {
				databaseId
			}
		`, i, repo.RepoOwner(), repo.RepoName()))
	}

	query := fmt.Sprintf(`query MapRepositoryNames { %s }`, strings.Join(queries, ""))

	graphqlResult := make(map[string]*struct {
		DatabaseID int `json:"databaseId"`
	})

	if err := client.GraphQL(host, query, nil, &graphqlResult); err != nil {
		return nil, fmt.Errorf("failed to look up repositories: %w", err)
	}

	repoKeys := make([]string, 0, len(repositories))
	for k := range graphqlResult {
		repoKeys = append(repoKeys, k)
	}
	sort.Strings(repoKeys)

	result := make([]string, len(repositories))
	for i, k := range repoKeys {
		result[i] = strconv.Itoa(graphqlResult[k].DatabaseID)
	}

	return result, nil
}
