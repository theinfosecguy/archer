"""Type definitions for the Archer secret validation system."""

from .core import (
    StringDict,
    OptionalStringDict,
    ProcessedHeaders,
    ProcessedParams,
    ProcessedUrl,
    ValidationResult,
)

from .http import (
    RequestKwargs,
    HeadersDict,
    QueryParamsDict,
    OptionalQueryParamsDict,
    StatusCodeList,
    ResponseData,
)

from .templates import (
    FieldPath,
    RequiredFieldsList,
    ErrorMessageCode,
    TemplateName,
)

__all__ = [
    # Core types
    "StringDict",
    "OptionalStringDict", 
    "ProcessedHeaders",
    "ProcessedParams",
    "ProcessedUrl",
    "ValidationResult",
    
    # HTTP types
    "RequestKwargs",
    "HeadersDict", 
    "QueryParamsDict",
    "OptionalQueryParamsDict",
    "StatusCodeList",
    "ResponseData",
    
    # Template types
    "FieldPath",
    "RequiredFieldsList",
    "ErrorMessageCode",
    "TemplateName",
]
