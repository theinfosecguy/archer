import json
from typing import Dict, Any
import httpx
from jsonpath_ng import parse
from archer.models import SecretTemplate
from archer.templates import TemplateLoader


class SecretValidator:
    """Validates secrets using configured templates."""

    def __init__(self, templates_dir: str = "templates"):
        self.template_loader = TemplateLoader(templates_dir)

    async def validate_secret(self, template_name: str, secret: str) -> Dict[str, Any]:
        """Validate a secret using the specified template."""
        template = self.template_loader.get_template(template_name)
        if not template:
            return {"valid": False, "error": f"Template '{template_name}' not found"}

        return await self._validate_with_template(template, secret)

    async def _validate_with_template(self, template: SecretTemplate, secret: str) -> Dict[str, Any]:
        """Validate secret using the template configuration."""
        # Replace ${SECRET} placeholder in headers
        headers = {}
        for key, value in template.request.headers.items():
            headers[key] = value.replace("${SECRET}", secret)

        async with httpx.AsyncClient() as client:
            try:
                response = await client.request(
                    method=template.method,
                    url=template.api_url,
                    headers=headers,
                    timeout=template.request.timeout
                )

                return self._check_response(response, template)

            except httpx.TimeoutException:
                return {"valid": False, "error": "Request timeout"}
            except Exception as e:
                return {"valid": False, "error": f"Request failed: {str(e)}"}

    def _check_response(self, response: httpx.Response, template: SecretTemplate) -> Dict[str, Any]:
        """Check if response meets success criteria."""
        # Check status code
        if response.status_code not in template.success_criteria.status_code:
            error_msg = template.error_handling.error_messages.get(
                response.status_code, f"HTTP {response.status_code}"
            )
            return {"valid": False, "error": error_msg}

        # Check required fields if specified
        if template.success_criteria.required_fields:
            try:
                response_data = response.json()
                for field_path in template.success_criteria.required_fields:
                    jsonpath_expr = parse(field_path)
                    if not jsonpath_expr.find(response_data):
                        return {"valid": False, "error": f"Required field '{field_path}' not found"}
            except json.JSONDecodeError:
                return {"valid": False, "error": "Invalid JSON response"}

        return {"valid": True, "message": "Secret is valid"}
