from typing import Dict, List, Optional, Any, Literal
from pydantic import BaseModel, field_validator


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
    mode: Optional[Literal["single", "multipart"]] = "single"
    required_variables: Optional[List[str]] = None
    api_url: str
    method: str = "GET"
    request: RequestConfig
    success_criteria: SuccessCriteria
    error_handling: ErrorHandling

    @field_validator('required_variables')
    @classmethod
    def validate_required_variables(cls, v, info):
        mode = info.data.get('mode')
        if mode == 'multipart':
            if not v:
                raise ValueError("required_variables is mandatory when mode is 'multipart'")
            for var in v:
                if not var.isupper() or not var.replace('_', '').isalnum():
                    raise ValueError(f"Variable '{var}' must be uppercase snake_case")
        return v