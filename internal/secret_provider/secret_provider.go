package secretprovider

import (
	"sync"
	"team-task-manager/internal/supports"
)

const (
	defaultSecretsDir          = "./secrets/"
	defaultContainerSecretsDir = "/run/secrets/"
)

type SecretProvider struct {
	mtx       sync.RWMutex
	cache     map[string]string
	secretDir string
}

func NewSecretProvider() *SecretProvider {
	dir := defaultSecretsDir
	if supports.IsInContainer() {
		dir = defaultContainerSecretsDir
	}
	return &SecretProvider{
		cache:     make(map[string]string),
		secretDir: dir,
	}
}

func (sp *SecretProvider) ReadSecret(key string) (string, error) {
	sp.mtx.RLock()
	v, ok := sp.cache[key]
	sp.mtx.RUnlock()
	if ok {
		return v, nil
	}

	v, err := supports.ReadSecretFile(supports.Concat(sp.secretDir, key))
	if err != nil {
		return "", err
	}

	sp.mtx.Lock()
	sp.cache[key] = v
	sp.mtx.Unlock()

	return v, nil
}
