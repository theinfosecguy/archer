"""Template discovery utilities for finding and listing available templates."""

import logging
from pathlib import Path
from typing import List

from archer.constants import DEFAULT_TEMPLATES_DIR, TEMPLATE_FILE_EXTENSIONS
from .exceptions import TemplateDirectoryNotFoundError
from .loader import is_file_path

logger = logging.getLogger(__name__)


def discover_templates_in_directory(templates_dir: str = DEFAULT_TEMPLATES_DIR) -> List[str]:
    """Discover all available template files in a directory.

    Args:
        templates_dir: Directory path to search for templates

    Returns:
        List of template names (without file extensions)

    Raises:
        TemplateDirectoryNotFoundError: If the templates directory doesn't exist
    """
    templates_path = Path(templates_dir)

    if not templates_path.exists():
        raise TemplateDirectoryNotFoundError(f"Templates directory not found: '{templates_path.absolute()}'")

    if not templates_path.is_dir():
        raise TemplateDirectoryNotFoundError(f"Templates path is not a directory: '{templates_path.absolute()}'")

    template_files = []
    for extension in TEMPLATE_FILE_EXTENSIONS:
        template_files.extend(templates_path.glob(f"*{extension}"))
    template_names = [f.stem for f in template_files]

    logger.debug(f"Discovered {len(template_names)} templates in '{templates_path.absolute()}'")
    return sorted(template_names)


def discover_templates(templates_dir: str = DEFAULT_TEMPLATES_DIR) -> List[str]:
    """Discover templates from the default templates directory.

    Args:
        templates_dir: Directory path to search for templates

    Returns:
        List of template names
    """
    try:
        return discover_templates_in_directory(templates_dir)
    except TemplateDirectoryNotFoundError as e:
        logger.error(str(e))
        return []
    except Exception as e:
        logger.error(f"Unexpected error discovering templates from '{templates_dir}': {str(e)}")
        return []


def get_template_identifier_display_name(identifier: str) -> str:
    """Get a display-friendly name for a template identifier.

    For file paths, returns just the filename without extension.
    For template names, returns the name as-is.

    Args:
        identifier: Template identifier (name or file path)

    Returns:
        Display-friendly name
    """
    if is_file_path(identifier):
        return Path(identifier).stem
    else:
        return identifier
