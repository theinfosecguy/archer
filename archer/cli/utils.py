from pathlib import Path
from typing import List, Optional
from archer.models import SecretTemplate
from archer.templates import TemplateLoader


def discover_templates(templates_dir: str = "templates") -> List[str]:
    """Discover all available template files."""
    templates_path = Path(templates_dir)
    if not templates_path.exists():
        return []

    return [f.stem for f in templates_path.glob("*.yaml")]


def load_template_safely(template_name: str, templates_dir: str = "templates") -> Optional[SecretTemplate]:
    """Load a template with error handling."""
    try:
        loader = TemplateLoader(templates_dir)
        return loader.get_template(template_name)
    except Exception:
        return None
