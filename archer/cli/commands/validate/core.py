import asyncio
import logging
from datetime import datetime, timezone
from typing import Optional, Dict, Any

from archer.core import SecretValidator
from archer.core.variables import process_url, process_headers, process_query_params
from archer.templates import TemplateLoader
from archer.constants import MODE_SINGLE, MODE_MULTIPART

logger = logging.getLogger(__name__)


class ValidationOrchestrator:
    """Orchestrates the validation process."""

    def __init__(self):
        self.template_loader = TemplateLoader()
        self.validator = SecretValidator()

    async def execute_validation(self, template_name: str, secret: Optional[str], 
                               template_file: Optional[str], variables: Dict[str, str]) -> Dict[str, Any]:
        """Execute the validation process."""
        # Load template
        if template_file:
            template = self.template_loader.get_template(template_file)
            validation_identifier = template_file
        else:
            template = self.template_loader.get_template(template_name)
            validation_identifier = template_name

        if not template:
            raise ValueError(f"Template not found: {template_file or template_name}")

        # Execute validation based on mode
        if template.mode == MODE_SINGLE:
            result = await self.validator.validate_secret(validation_identifier, secret)
        elif template.mode == MODE_MULTIPART:
            result = await self.validator.validate_secret_multipart(validation_identifier, variables)
        else:
            raise ValueError(f"Unsupported template mode: {template.mode}")

        return result

    def build_masked_artifacts(self, template, variables: Dict[str, str]):
        """Build masked URL, headers, and query params for output."""
        try:
            _, masked_url = process_url(template.api_url, variables)
            _, masked_headers = process_headers(template.request.headers, variables)
            _, masked_query_params = process_query_params(template.request.query_params, variables)
            return masked_url, masked_headers, masked_query_params
        except Exception:
            return None, None, None
