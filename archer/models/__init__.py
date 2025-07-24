from typing import Dict, List, Optional, Any
from pydantic import BaseModel, validator


class RequestConfig(BaseModel):
    """Request configuration for API calls."""
    headers: Dict[str, str]
    timeout: int = 10
    data: Optional[str] = None
    json_data: Optional[Dict[str, Any]] = None
    query_params: Optional[Dict[str, str]] = None

    @validator('json_data')
    def validate_mutual_exclusion(cls, v, values):
        if v and values.get('data'):
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
    request: RequestConfig
    success_criteria: SuccessCriteria
    error_handling: ErrorHandling
