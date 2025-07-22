import json
import logging
from typing import Dict, Any
import httpx
from jsonpath_ng import parse
from archer.models import SecretTemplate
from archer.templates import TemplateLoader

logger = logging.getLogger(__name__)


class SecretValidator:
    """Validates secrets using configured templates."""

    def __init__(self, templates_dir: str = "templates"):
        self.template_loader = TemplateLoader(templates_dir)
        logger.info(f"SecretValidator initialized with templates directory '{templates_dir}'")

    async def validate_secret(self, template_name: str, secret: str) -> Dict[str, Any]:
        """Validate a secret using the specified template."""
        logger.info(f"Starting secret validation using template '{template_name}'")

        template = self.template_loader.get_template(template_name)
        if not template:
            logger.error(f"Validation failed: template '{template_name}' not found in templates directory")
            return {"valid": False, "error": f"Template '{template_name}' not found"}

        logger.debug(f"Loaded template '{template.name}': {template.description}")
        return await self._validate_with_template(template, secret)

    async def _validate_with_template(self, template: SecretTemplate, secret: str) -> Dict[str, Any]:
        """Validate secret using the template configuration."""
        # Replace ${SECRET} placeholder in headers
        headers = {}
        for key, value in template.request.headers.items():
            headers[key] = value.replace("${SECRET}", "***MASKED***" if "${SECRET}" in value else value)

        logger.debug(f"Preparing {template.method} request to {template.api_url} with masked headers: {headers}")

        # Prepare actual headers with secret
        actual_headers = {}
        for key, value in template.request.headers.items():
            actual_headers[key] = value.replace("${SECRET}", secret)

        # Prepare request kwargs
        request_kwargs = {
            "method": template.method,
            "url": template.api_url,
            "headers": actual_headers,
            "timeout": template.request.timeout
        }

        # Add request body if present
        if template.request.data:
            request_kwargs["content"] = template.request.data
        elif template.request.json_data:
            request_kwargs["json"] = template.request.json_data

        async with httpx.AsyncClient() as client:
            try:
                logger.debug(f"Sending HTTP request with {template.request.timeout}s timeout")
                response = await client.request(**request_kwargs)

                logger.info(f"API request completed with status code {response.status_code}")

                # Log response content in debug mode
                if logger.isEnabledFor(logging.DEBUG):
                    try:
                        response_content = response.text
                        logger.debug(f"Response content: {response_content}")
                    except Exception as e:
                        logger.debug(f"Could not read response content: {str(e)}")

                return self._check_response(response, template)

            except httpx.TimeoutException:
                logger.error(f"API request timed out after {template.request.timeout} seconds")
                return {"valid": False, "error": "Request timeout"}
            except Exception as e:
                logger.error(f"API request failed with exception: {str(e)}")
                return {"valid": False, "error": f"Request failed: {str(e)}"}

    def _check_response(self, response: httpx.Response, template: SecretTemplate) -> Dict[str, Any]:
        """Check if response meets success criteria."""
        logger.debug(f"Validating response against template success criteria")

        # Check status code
        if response.status_code not in template.success_criteria.status_code:
            logger.warning(f"Status code validation failed: got {response.status_code}, expected one of {template.success_criteria.status_code}")
            error_msg = template.error_handling.error_messages.get(
                response.status_code, f"HTTP {response.status_code}"
            )
            return {"valid": False, "error": error_msg}

        logger.debug(f"Status code validation passed: {response.status_code} is in expected range")

        # Check required fields if specified
        if template.success_criteria.required_fields:
            logger.debug(f"Checking {len(template.success_criteria.required_fields)} required fields in JSON response")
            try:
                response_data = response.json()
                for field_path in template.success_criteria.required_fields:
                    jsonpath_expr = parse(field_path)
                    if not jsonpath_expr.find(response_data):
                        logger.warning(f"Required field validation failed: '{field_path}' not found in response")
                        return {"valid": False, "error": f"Required field '{field_path}' not found"}
                    else:
                        logger.debug(f"Required field validation passed: '{field_path}' found in response")
            except json.JSONDecodeError:
                logger.error("Response validation failed: API returned invalid JSON")
                return {"valid": False, "error": "Invalid JSON response"}

        logger.info("Secret validation completed successfully - all criteria met")
        return {"valid": True, "message": "Secret is valid"}
