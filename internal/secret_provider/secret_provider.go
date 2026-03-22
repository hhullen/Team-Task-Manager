package secretprovider

import (
	"path/filepath"
	"sync"
	"team-task-manager/internal/supports"
)

type SecretProvider struct {
	mtx       sync.RWMutex
	cache     map[string]string
	secretDir string
}

func NewSecretProvider(secretDir string) *SecretProvider {
	return &SecretProvider{
		cache:     make(map[string]string),
		secretDir: secretDir,
	}
}

func (sp *SecretProvider) ReadSecret(key string) (string, error) {
	sp.mtx.RLock()
	v, ok := sp.cache[key]
	sp.mtx.RUnlock()
	if ok {
		return v, nil
	}

	v, err := supports.ReadSecretFile(filepath.Join(sp.secretDir, key))
	if err != nil {
		return "", err
	}

	sp.mtx.Lock()
	sp.cache[key] = v
	sp.mtx.Unlock()

	return v, nil
}
