package pulse

import "time"

// Options stores the required parameters to be used by the client for authenticating with Feedzai Pulse API.
type Options struct {
	// Username is the active directory username.
	Username string
	// Password is the active directory password.
	Password string
	// BaseURL is the URL of the Feedzai Pulse website, for example
	// https://feedzai-pulse-stg.voyagerinnovation.com
	BaseURL  string
	// Timeout specifies a time limit for requests made by the client.
	Timeout  time.Duration
}
