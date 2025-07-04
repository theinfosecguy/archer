from typing import Dict, List, Optional
from pydantic import BaseModel


class RequestConfig(BaseModel):
    """Request configuration for API calls."""
    headers: Dict[str, str]
    timeout: int = 10


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
