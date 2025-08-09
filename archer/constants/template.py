"""Template-related constants for the Archer secret validation system."""

import re

# Template modes
MODE_SINGLE = "single"
MODE_MULTIPART = "multipart"

# File system
DEFAULT_TEMPLATES_DIR = "templates"
TEMPLATE_FILE_EXTENSION = ".yaml"
TEMPLATE_FILE_EXTENSIONS = [".yaml", ".yml"]

# Variable handling
VARIABLE_PATTERN = re.compile(r'\$\{([^}]+)\}')
SECRET_VARIABLE_NAME = "SECRET"

# Variable name validation patterns
UPPER_SNAKE_CASE_PATTERN = re.compile(r'^[A-Z][A-Z0-9_]*$')
KEBAB_CASE_PATTERN = re.compile(r'^[a-z][a-z0-9-]*$')
