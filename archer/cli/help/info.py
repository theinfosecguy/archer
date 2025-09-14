"""Help text for info command."""

import textwrap


def get_help() -> str:
    """Get help text for info command."""
    return textwrap.dedent("""\
    \b
    Arguments:
      TEMPLATE_NAME  Name of the template to inspect

    \b
    Options:
      --template-file FILE  Load template from specific file instead of built-in
      --help               Show this message and exit

    \b
    Displayed Information:
      - Template name and description
      - Validation mode (single or multipart)
      - API endpoint and HTTP method
      - Required variables (for multipart templates)
      - Usage examples with correct syntax
      - Request headers and query parameters (with masked values)
      - Success criteria and error handling configuration

    \b
    Examples:
      # Get info for built-in templates
      archer info github
      archer info ghost
      archer info openai
      
      # Get info for custom template file
      archer info --template-file ./custom-api.yaml
      archer info mytemplate --template-file ./templates/custom.yaml
      
      Note: When using --template-file, TEMPLATE_NAME is ignored

    \b
    Usage Flow:
      1. Run 'archer list' to see available templates
      2. Run 'archer info <template>' to understand requirements
      3. Use the displayed usage example to validate your secret
      
      For single mode templates:
      archer validate <template> <your_secret>
      
      For multipart mode templates:
      archer validate <template> --var key1=value1 --var key2=value2

    The displayed usage examples show the exact command syntax needed
    for validation, including all required variables in kebab-case format.
    """)
