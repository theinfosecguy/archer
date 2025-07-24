"""Core type definitions for the Archer validation system."""

from typing import Dict, Optional, Tuple, Any, TypeAlias

# Basic type aliases
StringDict: TypeAlias = Dict[str, str]
OptionalStringDict: TypeAlias = Optional[Dict[str, str]]
ValidationResult: TypeAlias = Dict[str, Any]

# Processing result type aliases
ProcessedHeaders: TypeAlias = Tuple[StringDict, StringDict]  # (request_headers, masked_headers)
ProcessedParams: TypeAlias = Tuple[OptionalStringDict, OptionalStringDict]  # (request_params, masked_params)
ProcessedUrl: TypeAlias = Tuple[str, str]  # (request_url, masked_url)
