"""Template loading and management for the Archer secret validation system.

This module provides functionality to load templates from either:
- Individual template files (e.g., ./my-template.yaml)
- Template directories (e.g., templates/github.yaml)


"""

from .loader import TemplateLoader, is_file_path, load_template_from_file, load_template_from_directory
from .discovery import discover_templates, discover_templates_in_directory, get_template_identifier_display_name
from .exceptions import (
    TemplateError,
    TemplateNotFoundError,
    TemplateLoadError,
    TemplateValidationError,
    TemplateDirectoryNotFoundError,
)

__all__ = [
    # Main loader class
    'TemplateLoader',

    # Core loading functions
    'is_file_path',
    'load_template_from_file',
    'load_template_from_directory',

    # Discovery functions
    'discover_templates',
    'discover_templates_in_directory',
    'get_template_identifier_display_name',

    # Exceptions
    'TemplateError',
    'TemplateNotFoundError',
    'TemplateLoadError',
    'TemplateValidationError',
    'TemplateDirectoryNotFoundError',
]
