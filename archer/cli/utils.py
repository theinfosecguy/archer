from typing import List, Optional
from archer.models import SecretTemplate
from archer.templates import TemplateLoader, discover_templates as find_templates
from archer.constants import DEFAULT_TEMPLATES_DIR


def discover_templates(templates_source: str = DEFAULT_TEMPLATES_DIR) -> List[str]:
    """Discover all available templates from a directory or file.

    Args:
        templates_source: Directory path or single file path

    Returns:
        List of template identifiers
    """
    return find_templates(templates_source)


def load_template_safely(template_identifier: str) -> Optional[SecretTemplate]:
    """Load a template with error handling.

    Args:
        template_identifier: Template name or file path

    Returns:
        SecretTemplate instance if successful, None otherwise
    """
    try:
        loader = TemplateLoader()
        return loader.get_template(template_identifier)
    except Exception:
        return None
