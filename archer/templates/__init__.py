from pathlib import Path
from typing import Dict, Optional
import yaml
from archer.models import SecretTemplate


class TemplateLoader:
    """Loads and manages secret validation templates."""

    def __init__(self, templates_dir: str = "templates"):
        self.templates_dir = Path(templates_dir)
        self._templates: Dict[str, SecretTemplate] = {}

    def load_template(self, template_name: str) -> Optional[SecretTemplate]:
        """Load a specific template by name."""
        template_path = self.templates_dir / f"{template_name}.yaml"

        if not template_path.exists():
            return None

        with open(template_path, 'r') as f:
            data = yaml.safe_load(f)

        template = SecretTemplate(**data)
        self._templates[template_name] = template
        return template

    def get_template(self, template_name: str) -> Optional[SecretTemplate]:
        """Get a template, loading it if not already cached."""
        if template_name not in self._templates:
            return self.load_template(template_name)
        return self._templates[template_name]
