from pathlib import Path
from typing import Dict, Optional
import logging
import yaml
from archer.models import SecretTemplate
from archer.constants import (
    DEFAULT_TEMPLATES_DIR,
    TEMPLATE_FILE_EXTENSION,
    ENCODING_UTF8,
    READ_MODE,
    TEMPLATE_LOADER_INITIALIZED,
    TEMPLATE_LOADED,
    TEMPLATE_CACHED,
    TEMPLATE_NOT_CACHED,
)

logger = logging.getLogger(__name__)


class TemplateLoader:
    """Loads and manages secret validation templates."""

    def __init__(self, templates_dir: str = DEFAULT_TEMPLATES_DIR):
        self.templates_dir = Path(templates_dir)
        self._templates: Dict[str, SecretTemplate] = {}
        logger.info(TEMPLATE_LOADER_INITIALIZED.format(dir=self.templates_dir.absolute()))

    def load_template(self, template_name: str) -> Optional[SecretTemplate]:
        """Load a specific template by name."""
        template_path = self.templates_dir / f"{template_name}{TEMPLATE_FILE_EXTENSION}"
        logger.debug(f"Attempting to load template from '{template_path.absolute()}'")

        if not template_path.exists():
            logger.error(f"Template file not found: '{template_path.absolute()}'")
            return None

        try:
            with open(template_path, READ_MODE, encoding=ENCODING_UTF8) as f:
                data = yaml.safe_load(f)
            logger.debug(f"Successfully parsed YAML from template file '{template_name}{TEMPLATE_FILE_EXTENSION}'")

            if not data:
                logger.error(f"Template file '{template_name}{TEMPLATE_FILE_EXTENSION}' is empty or contains no valid YAML")
                return None

            template = SecretTemplate(**data)
            self._templates[template_name] = template
            logger.info(TEMPLATE_LOADED.format(name=template_name, description=template.description))
            return template

        except FileNotFoundError:
            logger.error(f"Template file not found: '{template_path.absolute()}'")
            return None
        except PermissionError:
            logger.error(f"Permission denied reading template file: '{template_path.absolute()}'")
            return None
        except yaml.YAMLError as e:
            logger.error(f"YAML parsing failed for template '{template_name}': {str(e)}")
            return None
        except ValueError as e:
            logger.error(f"Template validation failed for '{template_name}': {str(e)}")
            return None
        except Exception as e:
            logger.error(f"Unexpected error loading template '{template_name}': {str(e)}")
            return None

    def get_template(self, template_name: str) -> Optional[SecretTemplate]:
        """Get a template, loading it if not already cached."""
        if template_name in self._templates:
            logger.debug(TEMPLATE_CACHED.format(name=template_name))
            return self._templates[template_name]

        logger.debug(TEMPLATE_NOT_CACHED.format(name=template_name))
        return self.load_template(template_name)
