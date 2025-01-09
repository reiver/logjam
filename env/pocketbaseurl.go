package env

import (
	"os"
)
const PocketBaseURLEnvVar string = "POCKETBASE_URL"

var PocketBaseURL string = os.Getenv(PocketBaseURLEnvVar)
