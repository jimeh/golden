package golden

import "os"

// EnvVarUpdateFunc checks if the UPDATE_GOLDEN environment variable is set to a
// truthy value. Truthy values are "yes", "on", "true", "t", and "1".
func EnvVarUpdateFunc() bool {
	env := os.Getenv("UPDATE_GOLDEN")
	for _, v := range []string{"yes", "on", "true", "t", "1"} {
		if env == v {
			return true
		}
	}

	return false
}
