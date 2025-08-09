"""Exceptions for template operations."""


class TemplateError(Exception):
    """Base exception for template-related errors."""
    pass


class TemplateNotFoundError(TemplateError):
    """Raised when a template cannot be found."""
    pass


class TemplateLoadError(TemplateError):
    """Raised when a template fails to load or parse."""
    pass


class TemplateValidationError(TemplateError):
    """Raised when a template fails validation."""
    pass


class TemplateDirectoryNotFoundError(TemplateError):
    """Raised when the templates directory does not exist."""
    pass
