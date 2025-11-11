package constants

// CLI indicators
const (
	SuccessIndicator = "[SUCCESS]"
	FailureIndicator = "[FAILED]"
	OptVar           = "--var"
)

// Environment variable names
const (
	EnvSecretName = "ARCHER_SECRET"
	EnvVarPrefix  = "ARCHER_VAR_"
)

// ANSI color codes
const (
	ColorRed   = "\033[91m"
	ColorCyan  = "\033[96m"
	ColorReset = "\033[0m"
)

// Security warnings
const (
	WarningSecretInCLI = ColorRed + "[WARNING] Secrets passed as CLI arguments are exposed in shell history, process lists, and logs.\n" + ColorReset
)
