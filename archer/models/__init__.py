import re
from typing import Dict, List, Optional, Any
from pydantic import BaseModel, field_validator, model_validator
from archer.constants import (
    MODE_SINGLE,
    MODE_MULTIPART,
    DEFAULT_TIMEOUT,
    DEFAULT_MAX_RETRIES,
    DEFAULT_RETRY_DELAY,
    METHOD_GET,
    UPPER_SNAKE_CASE_PATTERN,
    SECRET_VARIABLE_NAME,
    VARIABLE_PATTERN,
    MODE_VALIDATION_ERROR,
    MULTIPART_REQUIRES_VARIABLES,
    SINGLE_MODE_NO_VARIABLES,
    MUTUAL_EXCLUSION_ERROR,
    UPPER_SNAKE_CASE_ERROR,
    SECRET_NOT_ALLOWED_MULTIPART,
    INVALID_VARIABLES_SINGLE,
    UNDEFINED_VARIABLES,
    UNUSED_REQUIRED_VARIABLES,
)


class RequestConfig(BaseModel):
    """Request configuration for API calls."""
    headers: Dict[str, str]
    timeout: int = DEFAULT_TIMEOUT
    data: Optional[str] = None
    json_data: Optional[Dict[str, Any]] = None
    query_params: Optional[Dict[str, str]] = None

    @field_validator('json_data')
    @classmethod
    def validate_mutual_exclusion(cls, v, info):
        if v and info.data.get('data'):
            raise ValueError(MUTUAL_EXCLUSION_ERROR)
        return v


class SuccessCriteria(BaseModel):
    """Success criteria for validating API responses."""
    status_code: List[int]
    required_fields: Optional[List[str]] = None


class ErrorHandling(BaseModel):
    """Error handling configuration."""
    max_retries: int = DEFAULT_MAX_RETRIES
    retry_delay: int = DEFAULT_RETRY_DELAY
    error_messages: Dict[int, str] = {}


class SecretTemplate(BaseModel):
    """Template for secret validation."""
    name: str
    description: str
    api_url: str
    method: str = METHOD_GET
    mode: Optional[str] = MODE_SINGLE
    required_variables: Optional[List[str]] = None
    request: RequestConfig
    success_criteria: SuccessCriteria
    error_handling: ErrorHandling

    @field_validator('mode')
    @classmethod
    def validate_mode(cls, v):
        """Validate that mode is either 'single' or 'multipart'."""
        if v not in [MODE_SINGLE, MODE_MULTIPART]:
            raise ValueError(MODE_VALIDATION_ERROR)
        return v

    @field_validator('required_variables')
    @classmethod
    def validate_required_variables_format(cls, v):
        """Validate that required_variables are in UPPER_SNAKE_CASE format."""
        if v is None:
            return v

        for var in v:
            if not UPPER_SNAKE_CASE_PATTERN.match(var):
                raise ValueError(UPPER_SNAKE_CASE_ERROR.format(var=var))
        return v

    @model_validator(mode='after')
    def validate_multipart_requirements(self):
        """Validate multipart mode requirements and variable usage consistency."""
        mode = self.mode
        required_variables = self.required_variables

        # If mode is multipart, required_variables is mandatory
        if mode == MODE_MULTIPART and not required_variables:
            raise ValueError(MULTIPART_REQUIRES_VARIABLES)

        # If mode is single, required_variables should not be present
        if mode == MODE_SINGLE and required_variables:
            raise ValueError(SINGLE_MODE_NO_VARIABLES)

        # Extract all template content for variable validation
        template_content = []
        template_content.append(self.api_url)

        # Add headers
        for header_value in self.request.headers.values():
            template_content.append(header_value)

        # Add query params
        if self.request.query_params:
            for param_value in self.request.query_params.values():
                template_content.append(param_value)

        # Add data
        if self.request.data:
            template_content.append(self.request.data)

        # Add json_data (convert to string representation)
        if self.request.json_data:
            template_content.append(str(self.request.json_data))

        # Find all variables used in template
        used_variables = set()

        for content in template_content:
            matches = VARIABLE_PATTERN.findall(content)
            used_variables.update(matches)

        # Validate based on mode
        if mode == MODE_SINGLE:
            # Only ${SECRET} should be used
            secret_var_set = {SECRET_VARIABLE_NAME}
            if used_variables and used_variables != secret_var_set:
                invalid_vars = used_variables - secret_var_set
                if invalid_vars:
                    raise ValueError(INVALID_VARIABLES_SINGLE.format(vars=', '.join(invalid_vars)))

        elif mode == MODE_MULTIPART:
            # No ${SECRET} should be used
            if SECRET_VARIABLE_NAME in used_variables:
                raise ValueError(SECRET_NOT_ALLOWED_MULTIPART)

            # All required_variables must be used
            if required_variables:
                required_set = set(required_variables)
                if not used_variables.issubset(required_set):
                    undefined_vars = used_variables - required_set
                    raise ValueError(UNDEFINED_VARIABLES.format(vars=', '.join(undefined_vars)))

                if not required_set.issubset(used_variables):
                    unused_vars = required_set - used_variables
                    raise ValueError(UNUSED_REQUIRED_VARIABLES.format(vars=', '.join(unused_vars)))

        return self
