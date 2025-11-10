package constants

// CLI indicators
const (
	SuccessIndicator = "[SUCCESS]"
	FailureIndicator = "[FAILED]"
	OptVar           = "--var"
)

// Environment variable names
const (
	EnvSecretName       = "ARCHER_SECRET"
	EnvVarPrefix        = "ARCHER_VAR_"
)

// Security warnings
const (
	WarningSecretInCLI = "[WARNING] Passing secrets as command-line arguments is not secure.\n" +
		"           Secrets will be exposed in:\n" +
		"           - Shell history (~/.bash_history, ~/.zsh_history)\n" +
		"           - Process lists (ps, top, htop)\n" +
		"           - System logs and monitoring tools\n" +
		"           \n" +
		"           Recommended: Use environment variables instead:\n" +
		"           - Single mode: export ARCHER_SECRET=\"your-secret\"\n" +
		"           - Multipart mode: export ARCHER_VAR_API_KEY=\"value\"\n"
)
