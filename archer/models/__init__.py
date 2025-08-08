import re
from typing import Dict, List, Optional, Any
from pydantic import BaseModel, field_validator, model_validator


class RequestConfig(BaseModel):
    """Request configuration for API calls."""
    headers: Dict[str, str]
    timeout: int = 10
    data: Optional[str] = None
    json_data: Optional[Dict[str, Any]] = None
    query_params: Optional[Dict[str, str]] = None

    @field_validator('json_data')
    @classmethod
    def validate_mutual_exclusion(cls, v, info):
        if v and info.data.get('data'):
            raise ValueError("Cannot specify both 'data' and 'json_data'")
        return v


class SuccessCriteria(BaseModel):
    """Success criteria for validating API responses."""
    status_code: List[int]
    required_fields: Optional[List[str]] = None


class ErrorHandling(BaseModel):
    """Error handling configuration."""
    max_retries: int = 2
    retry_delay: int = 1
    error_messages: Dict[int, str] = {}


class SecretTemplate(BaseModel):
    """Template for secret validation."""
    name: str
    description: str
    api_url: str
    method: str = "GET"
    mode: Optional[str] = "single"
    required_variables: Optional[List[str]] = None
    request: RequestConfig
    success_criteria: SuccessCriteria
    error_handling: ErrorHandling

    @field_validator('mode')
    @classmethod
    def validate_mode(cls, v):
        """Validate that mode is either 'single' or 'multipart'."""
        if v not in ['single', 'multipart']:
            raise ValueError("mode must be either 'single' or 'multipart'")
        return v

    @field_validator('required_variables')
    @classmethod
    def validate_required_variables_format(cls, v):
        """Validate that required_variables are in UPPER_SNAKE_CASE format."""
        if v is None:
            return v
        
        pattern = re.compile(r'^[A-Z][A-Z0-9_]*$')
        for var in v:
            if not pattern.match(var):
                raise ValueError(f"Variable '{var}' must be in UPPER_SNAKE_CASE format")
        return v

    @model_validator(mode='after')
    def validate_multipart_requirements(self):
        """Validate multipart mode requirements and variable usage consistency."""
        mode = self.mode
        required_variables = self.required_variables
        
        # If mode is multipart, required_variables is mandatory
        if mode == 'multipart' and not required_variables:
            raise ValueError("required_variables is mandatory when mode is 'multipart'")
        
        # If mode is single, required_variables should not be present
        if mode == 'single' and required_variables:
            raise ValueError("required_variables should not be specified when mode is 'single'")
        
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
        variable_pattern = re.compile(r'\$\{([^}]+)\}')
        used_variables = set()
        
        for content in template_content:
            matches = variable_pattern.findall(content)
            used_variables.update(matches)
        
        # Validate based on mode
        if mode == 'single':
            # Only ${SECRET} should be used
            if used_variables and used_variables != {'SECRET'}:
                invalid_vars = used_variables - {'SECRET'}
                if invalid_vars:
                    raise ValueError(f"In single mode, only ${{SECRET}} is allowed. Found: {', '.join(invalid_vars)}")
        
        elif mode == 'multipart':
            # No ${SECRET} should be used
            if 'SECRET' in used_variables:
                raise ValueError("${SECRET} is not allowed in multipart mode. Use custom variables instead.")
            
            # All required_variables must be used
            if required_variables:
                required_set = set(required_variables)
                if not used_variables.issubset(required_set):
                    undefined_vars = used_variables - required_set
                    raise ValueError(f"Template uses undefined variables: {', '.join(undefined_vars)}. Add them to required_variables.")
                
                if not required_set.issubset(used_variables):
                    unused_vars = required_set - used_variables
                    raise ValueError(f"Required variables not used in template: {', '.join(unused_vars)}")
        
        return self