package datastruct

type Status struct {
	Message string `json:"status,omitempty" example:"status message"`
}

func (s Status) GetStatus() string {
	return s.Message
}

type CachedStatus struct {
	Cached bool `json:"cached" schema:"cached" example:"false"`
}

func (r *CachedStatus) SetCached(is bool) {
	r.Cached = is
}

type AvoidCacheFlag struct {
	Flag bool `schema:"avoid_cache" json:"avoid_cache" example:"true"`
}

func (a *AvoidCacheFlag) AvoidCache() bool {
	return a.Flag
}

const (
	StatusOK                   = "Success"
	StatusNotFound             = "resource not found"
	StatusAlreadyExists        = "resource already exists"
	StatusUserAlreadyExists    = "user already exists"
	StatusWrongLoginOrPassword = "wrong login or password"
	StatusInvalidToken         = "invalid token"
	StatusRefreshTokenReset    = "refresh token have been reset"
	StatusServiceError         = "service failed exec request"
)
