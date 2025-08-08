"""Variable processing utilities for single and multipart templates."""

import re
import logging
from typing import Dict, Optional, Tuple, Any

logger = logging.getLogger(__name__)

# Compile the pattern once at module level
VARIABLE_PATTERN = re.compile(r'\$\{([^}]+)\}')


def inject_variables(content: str, variables: Dict[str, str]) -> str:
    """Inject variables into content string."""
    if not content:
        return content
    
    def replace_variable(match):
        var_name = match.group(1)
        if var_name in variables:
            return variables[var_name]
        else:
            logger.warning(f"Variable '{var_name}' not found in provided variables")
            return match.group(0)  # Return original placeholder if variable not found
    
    return VARIABLE_PATTERN.sub(replace_variable, content)


def mask_variables(content: str) -> str:
    """Mask all variables in content for safe logging."""
    if not content:
        return content
    
    def mask_variable(match):
        var_name = match.group(1)
        return f"${{***{var_name}_MASKED***}}"
    
    return VARIABLE_PATTERN.sub(mask_variable, content)


def process_headers(headers: Dict[str, str], variables: Dict[str, str]) -> Tuple[Dict[str, str], Dict[str, str]]:
    """Process headers for both request use and masked logging."""
    request_headers = {}
    masked_headers = {}
    
    for key, value in headers.items():
        request_headers[key] = inject_variables(value, variables)
        masked_headers[key] = mask_variables(value)
    
    return request_headers, masked_headers


def process_query_params(query_params: Optional[Dict[str, str]], variables: Dict[str, str]) -> Tuple[Optional[Dict[str, str]], Optional[Dict[str, str]]]:
    """Process query parameters for both request use and masked logging."""
    if not query_params:
        return None, None
    
    request_params = {}
    masked_params = {}
    
    for key, value in query_params.items():
        request_params[key] = inject_variables(value, variables)
        masked_params[key] = mask_variables(value)
    
    return request_params, masked_params


def process_url(url: str, variables: Dict[str, str]) -> Tuple[str, str]:
    """Process URL for both request use and masked logging."""
    request_url = inject_variables(url, variables)
    masked_url = mask_variables(url)
    return request_url, masked_url


def process_data(data: Optional[str], variables: Dict[str, str]) -> Tuple[Optional[str], Optional[str]]:
    """Process data string for both request use and masked logging."""
    if not data:
        return None, None
    
    request_data = inject_variables(data, variables)
    masked_data = mask_variables(data)
    return request_data, masked_data


def process_json_data(json_data: Optional[Dict[str, Any]], variables: Dict[str, str]) -> Tuple[Optional[Dict[str, Any]], Optional[Dict[str, Any]]]:
    """Process JSON data for both request use and masked logging."""
    if not json_data:
        return None, None
    
    def process_json_recursive(obj: Any) -> Tuple[Any, Any]:
        if isinstance(obj, str):
            request_value = inject_variables(obj, variables)
            masked_value = mask_variables(obj)
            return request_value, masked_value
        elif isinstance(obj, dict):
            request_dict = {}
            masked_dict = {}
            for key, value in obj.items():
                req_val, mask_val = process_json_recursive(value)
                request_dict[key] = req_val
                masked_dict[key] = mask_val
            return request_dict, masked_dict
        elif isinstance(obj, list):
            request_list = []
            masked_list = []
            for item in obj:
                req_item, mask_item = process_json_recursive(item)
                request_list.append(req_item)
                masked_list.append(mask_item)
            return request_list, masked_list
        else:
            # For non-string values (numbers, bools, etc.), return as-is
            return obj, obj
    
    return process_json_recursive(json_data)


def get_variables_from_template_content(content: str) -> set:
    """Extract all variable names from template content."""
    if not content:
        return set()
    
    matches = VARIABLE_PATTERN.findall(content)
    return set(matches)


def validate_variables_provided(required_variables: list, provided_variables: Dict[str, str]) -> list:
    """Validate that all required variables are provided. Returns list of missing variables."""
    if not required_variables:
        return []
    
    missing = []
    for var in required_variables:
        if var not in provided_variables or not provided_variables[var].strip():
            missing.append(var)
    
    return missing


def parse_var_args(var_args: list) -> Dict[str, str]:
    """Parse --var key=value arguments into a dictionary."""
    variables = {}
    
    for var_arg in var_args:
        if '=' not in var_arg:
            raise ValueError(f"Invalid variable format: '{var_arg}'. Use --var key=value")
        
        # Split only on first '=' to allow '=' in values
        key, value = var_arg.split('=', 1)
        
        # Validate key format (kebab-case)
        if not re.match(r'^[a-z][a-z0-9-]*$', key):
            raise ValueError(f"Variable name '{key}' must be in kebab-case format (e.g., 'api-token', 'base-url')")
        
        # Convert kebab-case to UPPER_SNAKE_CASE
        upper_snake_key = key.upper().replace('-', '_')
        variables[upper_snake_key] = value
    
    return variables


def format_var_name_for_cli(upper_snake_name: str) -> str:
    """Convert UPPER_SNAKE_CASE to kebab-case for CLI display."""
    return upper_snake_name.lower().replace('_', '-')
