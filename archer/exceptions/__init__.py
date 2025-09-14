"""Centralized exception classes for Archer CLI."""


class CLIError(Exception):
    """Base exception for CLI-related errors."""
    pass


class JSONWriteError(CLIError):
    """Raised when JSON output cannot be written."""
    pass
