import json
import logging
from typing import Dict, Any
import httpx
from jsonpath_ng import parse

from archer.models import SecretTemplate
from archer.templates import TemplateLoader
from archer.core.variables import (
    validate_variables_provided,
    process_url,
    process_headers,
    process_query_params,
    process_data,
    process_json_data
)
from archer.types import (
    ValidationResult,
)

logger = logging.getLogger(__name__)


class SecretValidator:
    """Validates secrets using configured templates."""

    def __init__(self, templates_dir: str = "templates"):
        self.template_loader = TemplateLoader(templates_dir)
        logger.info(f"SecretValidator initialized with templates directory '{templates_dir}'")

    async def validate_secret(self, template_name: str, secret: str) -> ValidationResult:
        """Validate a secret using the specified template (single mode)."""
        logger.info(f"Starting secret validation using template '{template_name}' in single mode")

        template = self.template_loader.get_template(template_name)
        if not template:
            logger.error(f"Validation failed: template '{template_name}' not found in templates directory")
            return {"valid": False, "error": f"Template '{template_name}' not found"}

        if template.mode != 'single':
            logger.error(f"Template '{template_name}' is not in single mode")
            return {"valid": False, "error": f"Template '{template_name}' is not a single mode template"}

        logger.debug(f"Loaded template '{template.name}': {template.description}")
        
        # For single mode, create variables dict with SECRET
        variables = {"SECRET": secret}
        return await self._validate_with_template(template, variables)

    async def validate_secret_multipart(self, template_name: str, variables: Dict[str, str]) -> ValidationResult:
        """Validate secrets using the specified multipart template."""
        logger.info(f"Starting secret validation using template '{template_name}' in multipart mode")

        template = self.template_loader.get_template(template_name)
        if not template:
            logger.error(f"Validation failed: template '{template_name}' not found in templates directory")
            return {"valid": False, "error": f"Template '{template_name}' not found"}

        if template.mode != 'multipart':
            logger.error(f"Template '{template_name}' is not in multipart mode")
            return {"valid": False, "error": f"Template '{template_name}' is not a multipart mode template"}

        logger.debug(f"Loaded template '{template.name}': {template.description}")
        
        # Validate all required variables are provided
        missing_vars = validate_variables_provided(template.required_variables, variables)
        if missing_vars:
            error_msg = f"Missing required variables: {', '.join(missing_vars)}"
            logger.error(error_msg)
            return {"valid": False, "error": error_msg}

        return await self._validate_with_template(template, variables)

    async def _validate_with_template(self, template: SecretTemplate, variables: Dict[str, str]) -> ValidationResult:
        """Validate using the template configuration with provided variables."""
        # Process URL for variable injection
        request_url, masked_url = process_url(template.api_url, variables)

        # Process headers for variable injection
        request_headers, masked_headers = process_headers(template.request.headers, variables)

        # Process query parameters for variable injection  
        request_query_params, masked_query_params = process_query_params(template.request.query_params, variables)

        # Process data for variable injection
        request_data, masked_data = process_data(template.request.data, variables)

        # Process JSON data for variable injection
        request_json_data, masked_json_data = process_json_data(template.request.json_data, variables)

        # Log what we're about to do (with masked values)
        logger.debug(f"Preparing {template.method} request to {masked_url} with headers: {masked_headers}")
        if masked_query_params:
            logger.debug(f"Query parameters (masked): {masked_query_params}")
        if masked_data:
            logger.debug(f"Request data (masked): {masked_data}")
        if masked_json_data:
            logger.debug(f"Request JSON data (masked): {masked_json_data}")

        # Prepare request kwargs
        request_kwargs = {
            "method": template.method,
            "url": request_url,
            "headers": request_headers,
            "timeout": template.request.timeout
        }

        # Add query parameters if present
        if request_query_params:
            request_kwargs["params"] = request_query_params

        # Add request body if present
        if request_data:
            request_kwargs["content"] = request_data
        elif request_json_data:
            request_kwargs["json"] = request_json_data

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

    def _check_response(self, response: httpx.Response, template: SecretTemplate) -> ValidationResult:
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
