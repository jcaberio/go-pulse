package internal

type PublishRequest struct {
	Async        bool `json:"async"`
	FullReload   bool `json:"fullReload"`
	Rolling      bool `json:"rolling"`
	SkipRecovery bool `json:"skipRecovery"`
}
