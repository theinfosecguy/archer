"""Template loader with support for both directory and file-based templates."""

import logging
from pathlib import Path
from typing import Optional
import yaml

from archer.models import SecretTemplate
from archer.constants import (
    DEFAULT_TEMPLATES_DIR,
    TEMPLATE_FILE_EXTENSION,
    ENCODING_UTF8,
    READ_MODE,
    TEMPLATE_LOADER_INITIALIZED,
    TEMPLATE_LOADED,
)
from .exceptions import (
    TemplateNotFoundError,
    TemplateLoadError,
    TemplateValidationError,
    TemplateDirectoryNotFoundError,
)

logger = logging.getLogger(__name__)


def is_file_path(identifier: str) -> bool:
    """Detect if identifier is a file path vs template name.

    Args:
        identifier: Template identifier (name or file path)

    Returns:
        True if identifier appears to be a file path, False if template name
    """
    return (
        '/' in identifier or
        '\\' in identifier or
        identifier.endswith('.yaml') or
        identifier.endswith('.yml')
    )


def load_template_from_file(file_path: Path) -> SecretTemplate:
    """Load a template from a specific file path.

    Args:
        file_path: Path to the template file

    Returns:
        Loaded SecretTemplate instance

    Raises:
        TemplateNotFoundError: If file doesn't exist
        TemplateLoadError: If file can't be read or parsed
        TemplateValidationError: If template data is invalid
    """
    logger.debug(f"Attempting to load template from '{file_path.absolute()}'")

    if not file_path.exists():
        raise TemplateNotFoundError(f"Template file not found: '{file_path.absolute()}'")

    try:
        with open(file_path, READ_MODE, encoding=ENCODING_UTF8) as f:
            data = yaml.safe_load(f)
        logger.debug(f"Successfully parsed YAML from template file '{file_path.name}'")

        if not data:
            raise TemplateLoadError(f"Template file '{file_path.name}' is empty or contains no valid YAML")

        template = SecretTemplate(**data)
        logger.info(TEMPLATE_LOADED.format(name=file_path.stem, description=template.description))
        return template

    except FileNotFoundError:
        raise TemplateNotFoundError(f"Template file not found: '{file_path.absolute()}'")
    except PermissionError:
        raise TemplateLoadError(f"Permission denied reading template file: '{file_path.absolute()}'")
    except yaml.YAMLError as e:
        raise TemplateLoadError(f"YAML parsing failed for template '{file_path.name}': {str(e)}")
    except ValueError as e:
        raise TemplateValidationError(f"Template validation failed for '{file_path.name}': {str(e)}")
    except Exception as e:
        raise TemplateLoadError(f"Unexpected error loading template '{file_path.name}': {str(e)}")


def load_template_from_directory(template_name: str, templates_dir: Path) -> SecretTemplate:
    """Load a template by name from a templates directory.

    Args:
        template_name: Name of the template (without extension)
        templates_dir: Directory containing template files

    Returns:
        Loaded SecretTemplate instance

    Raises:
        TemplateDirectoryNotFoundError: If templates directory doesn't exist
        TemplateNotFoundError: If template file doesn't exist
        TemplateLoadError: If template can't be loaded
        TemplateValidationError: If template data is invalid
    """
    if not templates_dir.exists():
        raise TemplateDirectoryNotFoundError(f"Templates directory not found: '{templates_dir.absolute()}'")

    template_path = templates_dir / f"{template_name}{TEMPLATE_FILE_EXTENSION}"
    return load_template_from_file(template_path)


class TemplateLoader:
    """Loads templates from default directory or individual files."""

    def __init__(self, default_templates_dir: str = DEFAULT_TEMPLATES_DIR):
        """Initialize template loader.

        Args:
            default_templates_dir: Default directory to look for templates
        """
        self.default_templates_dir = Path(default_templates_dir)
        logger.info(TEMPLATE_LOADER_INITIALIZED.format(dir=self.default_templates_dir.absolute()))

    def get_template(self, template_identifier: str) -> Optional[SecretTemplate]:
        """Get a template by name or file path.

        Args:
            template_identifier: Template name (for directory lookup) or file path

        Returns:
            SecretTemplate instance if found and valid, None otherwise
        """
        try:
            if is_file_path(template_identifier):
                # Direct file path
                file_path = Path(template_identifier)
                return load_template_from_file(file_path)
            else:
                # Template name - look in default directory
                return load_template_from_directory(template_identifier, self.default_templates_dir)

        except (TemplateNotFoundError, TemplateLoadError, TemplateValidationError,
                TemplateDirectoryNotFoundError) as e:
            logger.error(str(e))
            return None
        except Exception as e:
            logger.error(f"Unexpected error loading template '{template_identifier}': {str(e)}")
            return None

    def load_template(self, template_identifier: str) -> Optional[SecretTemplate]:
        """Load a template by name or file path (alias for get_template for backward compatibility).

        Args:
            template_identifier: Template name or file path

        Returns:
            SecretTemplate instance if found and valid, None otherwise
        """
        return self.get_template(template_identifier)
