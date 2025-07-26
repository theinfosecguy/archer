import json
import logging
from typing import Dict, Any
import httpx
from jsonpath_ng import parse

from archer.models import SecretTemplate
from archer.templates import TemplateLoader
from archer.types import (
    StringDict,
    OptionalStringDict,
    ProcessedHeaders,
    ProcessedParams,
    ProcessedUrl,
    ValidationResult,
)

logger = logging.getLogger(__name__)


class SecretValidator:
    """Validates secrets using configured templates."""

    def __init__(self, templates_dir: str = "templates"):
        self.template_loader = TemplateLoader(templates_dir)
        logger.info(f"SecretValidator initialized with templates directory '{templates_dir}'")

    def _inject_secret_into_string(self, value: str, secret: str) -> str:
        """Inject secret into a string if it contains the placeholder."""
        return value.replace("${SECRET}", secret) if "${SECRET}" in value else value

    def _mask_secret_in_string(self, value: str) -> str:
        """Mask secret in a string if it contains the placeholder."""
        return value.replace("${SECRET}", "***MASKED***") if "${SECRET}" in value else value

    def _process_headers(self, headers: StringDict, secret: str) -> ProcessedHeaders:
        """Process headers for both request use and masked logging."""
        request_headers = {}
        masked_headers = {}

        for key, value in headers.items():
            request_headers[key] = self._inject_secret_into_string(value, secret)
            masked_headers[key] = self._mask_secret_in_string(value)

        return request_headers, masked_headers

    def _process_query_params(self, query_params: OptionalStringDict, secret: str) -> ProcessedParams:
        """Process query parameters for both request use and masked logging."""
        if not query_params:
            return None, None

        request_params = {}
        masked_params = {}

        for key, value in query_params.items():
            request_params[key] = self._inject_secret_into_string(value, secret)
            masked_params[key] = self._mask_secret_in_string(value)

        return request_params, masked_params

    def _process_url(self, url: str, secret: str) -> ProcessedUrl:
        """Process URL for both request use and masked logging."""
        request_url = self._inject_secret_into_string(url, secret)
        masked_url = self._mask_secret_in_string(url)
        return request_url, masked_url

    def _inject_variables_into_string(self, value: str, variables: Dict[str, str]) -> str:
        """Inject multiple variables into a string."""
        result = value
        for var_name, var_value in variables.items():
            placeholder = f"${{{var_name}}}"
            result = result.replace(placeholder, var_value)
        return result

    def _mask_variables_in_string(self, value: str, variables: Dict[str, str]) -> str:
        """Mask multiple variables in a string."""
        result = value
        for var_name in variables.keys():
            placeholder = f"${{{var_name}}}"
            result = result.replace(placeholder, "***MASKED***")
        return result

    def _process_headers_multipart(self, headers: StringDict, variables: Dict[str, str]) -> ProcessedHeaders:
        """Process headers with multiple variables."""
        request_headers = {}
        masked_headers = {}

        for key, value in headers.items():
            request_headers[key] = self._inject_variables_into_string(value, variables)
            masked_headers[key] = self._mask_variables_in_string(value, variables)

        return request_headers, masked_headers

    def _process_query_params_multipart(self, query_params: OptionalStringDict, variables: Dict[str, str]) -> ProcessedParams:
        """Process query parameters with multiple variables."""
        if not query_params:
            return None, None

        request_params = {}
        masked_params = {}

        for key, value in query_params.items():
            request_params[key] = self._inject_variables_into_string(value, variables)
            masked_params[key] = self._mask_variables_in_string(value, variables)

        return request_params, masked_params

    def _process_url_multipart(self, url: str, variables: Dict[str, str]) -> ProcessedUrl:
        """Process URL with multiple variables."""
        request_url = self._inject_variables_into_string(url, variables)
        masked_url = self._mask_variables_in_string(url, variables)
        return request_url, masked_url

    def _process_body_data_multipart(self, data: Any, variables: Dict[str, str]) -> Any:
        """Process request body data with multiple variables."""
        if isinstance(data, str):
            return self._inject_variables_into_string(data, variables)
        elif isinstance(data, dict):
            return {k: self._process_body_data_multipart(v, variables) for k, v in data.items()}
        elif isinstance(data, list):
            return [self._process_body_data_multipart(item, variables) for item in data]
        else:
            return data

    async def validate_secret(self, template_name: str, secret: str) -> ValidationResult:
        """Validate a secret using the specified template."""
        logger.info(f"Starting secret validation using template '{template_name}'")

        template = self.template_loader.get_template(template_name)
        if not template:
            logger.error(f"Validation failed: template '{template_name}' not found in templates directory")
            return {"valid": False, "error": f"Template '{template_name}' not found"}

        logger.debug(f"Loaded template '{template.name}': {template.description}")
        return await self._validate_with_template(template, secret)

    async def validate_secret_multipart(self, template_name: str, variables: Dict[str, str]) -> ValidationResult:
        """Validate a secret using multiple variables - NEW METHOD."""
        logger.info(f"Starting multipart secret validation using template '{template_name}' with {len(variables)} variable(s)")

        template = self.template_loader.get_template(template_name)
        if not template:
            logger.error(f"Validation failed: template '{template_name}' not found in templates directory")
            return {"valid": False, "error": f"Template '{template_name}' not found"}

        logger.debug(f"Loaded template '{template.name}': {template.description}")
        
        # Validate that all required variables are provided for multipart templates
        if template.mode == 'multipart':
            missing_vars = set(template.required_variables) - set(variables.keys())
            if missing_vars:
                error_msg = f"Missing required variables: {', '.join(missing_vars)}"
                logger.error(f"Validation failed: {error_msg}")
                return {"valid": False, "error": error_msg}
            
            logger.debug(f"Using variables: {list(variables.keys())}")
        
        return await self._validate_with_template_multipart(template, variables)

    async def _validate_with_template(self, template: SecretTemplate, secret: str) -> ValidationResult:
        """Validate secret using the template configuration."""
        # Process URL for secret injection
        request_url, masked_url = self._process_url(template.api_url, secret)

        # Process headers for secret injection
        request_headers, masked_headers = self._process_headers(template.request.headers, secret)

        # Process query parameters for secret injection  
        request_query_params, masked_query_params = self._process_query_params(template.request.query_params, secret)

        # Log what we're about to do (with masked values)
        logger.debug(f"Preparing {template.method} request to {masked_url} with headers: {masked_headers}")
        if masked_query_params:
            logger.debug(f"Query parameters (masked): {masked_query_params}")

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

    async def _validate_with_template_multipart(self, template: SecretTemplate, variables: Dict[str, str]) -> ValidationResult:
        """Validate variables using the template configuration - NEW METHOD."""
        # Process URL for variable injection
        request_url, masked_url = self._process_url_multipart(template.api_url, variables)

        # Process headers for variable injection
        request_headers, masked_headers = self._process_headers_multipart(template.request.headers, variables)

        # Process query parameters for variable injection  
        request_query_params, masked_query_params = self._process_query_params_multipart(template.request.query_params, variables)

        # Log what we're about to do (with masked values)
        logger.debug(f"Preparing {template.method} request to {masked_url} with headers: {masked_headers}")
        if masked_query_params:
            logger.debug(f"Query parameters (masked): {masked_query_params}")

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

        # Add request body if present (with variable injection)
        if template.request.data:
            request_kwargs["content"] = self._process_body_data_multipart(template.request.data, variables)
        elif template.request.json_data:
            request_kwargs["json"] = self._process_body_data_multipart(template.request.json_data, variables)

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
