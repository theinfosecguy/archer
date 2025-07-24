"""HTTP request related type definitions."""

from typing import TypeAlias, Dict, Any, Optional

# HTTP request configuration types
RequestKwargs: TypeAlias = Dict[str, Any]
HeadersDict: TypeAlias = Dict[str, str]
QueryParamsDict: TypeAlias = Dict[str, str]
OptionalQueryParamsDict: TypeAlias = Optional[Dict[str, str]]

# HTTP response types  
StatusCodeList: TypeAlias = list[int]
ResponseData: TypeAlias = Dict[str, Any]
