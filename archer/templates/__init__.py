from pathlib import Path
from typing import Dict, Optional
import logging
import yaml
from archer.models import SecretTemplate

logger = logging.getLogger(__name__)


class TemplateLoader:
    """Loads and manages secret validation templates."""

    def __init__(self, templates_dir: str = "templates"):
        self.templates_dir = Path(templates_dir)
        self._templates: Dict[str, SecretTemplate] = {}
        logger.info(f"TemplateLoader initialized with directory '{self.templates_dir.absolute()}'")

    def load_template(self, template_name: str) -> Optional[SecretTemplate]:
        """Load a specific template by name."""
        template_path = self.templates_dir / f"{template_name}.yaml"
        logger.debug(f"Attempting to load template from '{template_path.absolute()}'")

        if not template_path.exists():
            logger.error(f"Template file not found: '{template_path.absolute()}'")
            return None

        try:
            with open(template_path, 'r') as f:
                data = yaml.safe_load(f)
            logger.debug(f"Successfully parsed YAML from template file '{template_name}.yaml'")

            template = SecretTemplate(**data)
            self._templates[template_name] = template
            logger.info(f"Template '{template_name}' loaded successfully: {template.description}")
            return template

        except yaml.YAMLError as e:
            logger.error(f"YAML parsing failed for template '{template_name}': {str(e)}")
            return None
        except Exception as e:
            logger.error(f"Template validation failed for '{template_name}': {str(e)}")
            return None

    def get_template(self, template_name: str) -> Optional[SecretTemplate]:
        """Get a template, loading it if not already cached."""
        if template_name in self._templates:
            logger.debug(f"Template '{template_name}' found in cache")
            return self._templates[template_name]

        logger.debug(f"Template '{template_name}' not cached, loading from disk")
        return self.load_template(template_name)
