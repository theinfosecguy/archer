"""Structured, strongly-typed JSON output models for `archer validate` command.

These models describe the shape of the JSON written when users invoke the
`validate` command with the --output-json option.
"""

from datetime import datetime
from typing import Dict, List, Optional, Literal
from pydantic import BaseModel, Field


class ValidationRequestMeta(BaseModel):
    """Metadata about the validation request."""
    template: str = Field(description="Template identifier provided by user (name or file path)")
    resolved_template_name: Optional[str] = Field(description="Template name as defined inside the template file")
    mode: Optional[Literal["single", "multipart"]] = Field(description="Template mode if template resolved")
    source: Optional[Literal["builtin", "file"]] = Field(description="Where the template was loaded from")
    method: Optional[str] = Field(description="HTTP method used for validation request")
    api_url_masked: Optional[str] = Field(description="Masked API URL with variables hidden")
    headers_masked: Optional[Dict[str, str]] = Field(description="Masked request headers")
    query_params_masked: Optional[Dict[str, str]] = Field(description="Masked query parameters")
    variables_provided: Optional[List[str]] = Field(description="Names of variables provided (values omitted)")
    started_at: datetime = Field(description="UTC timestamp when validation started")
    finished_at: datetime = Field(description="UTC timestamp when validation finished (even if failed early)")
    duration_ms: float = Field(description="Total duration in milliseconds")


class ValidationResponseMeta(BaseModel):
    """Metadata about the validation response."""
    status_code: Optional[int] = Field(description="HTTP status code returned by the endpoint if request executed")
    required_fields_checked: Optional[List[str]] = Field(description="List of JSONPath fields checked, if any")
    failed_required_field: Optional[str] = Field(description="First required field that was missing, if applicable")
    error: Optional[str] = Field(description="Low-level error encountered before or during request execution")


class ValidationResultJSON(BaseModel):
    """Top-level JSON output for validate command."""
    command: Literal["validate"] = "validate"
    version: str = Field(description="Archer package version")
    valid: bool = Field(description="Indicates whether the secret validation succeeded")
    message: Optional[str] = Field(description="Success message when valid is true")
    error: Optional[str] = Field(description="Error message when valid is false")
    request: ValidationRequestMeta
    response: ValidationResponseMeta
