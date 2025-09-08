import click
import textwrap
from archer.cli.utils import load_template_safely
from archer.core.variables import format_var_name_for_cli
from archer.constants import (
    MODE_SINGLE,
    MODE_MULTIPART,
    OPT_VAR,
    FAILURE_INDICATOR,
    SECRET_VARIABLE_NAME,
    VARIABLE_PATTERN,
    MASKED_VARIABLE_FORMAT,
    MASKED_VARIABLE_PREFIX,
    MASKED_VARIABLE_SUFFIX,
)


@click.command()
@click.argument('template_name')
@click.option('--template-file', type=click.Path(exists=True, file_okay=True, dir_okay=False), help='Load template from specific file instead of built-in')
def info(template_name: str, template_file: str) -> None:
    """Show detailed information about a template.

    Displays comprehensive information about a template's configuration,
    required variables, API endpoints, and usage examples.
    """
    
    help_text = textwrap.dedent("""\
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

    if template_file:
        template = load_template_safely(template_file)
    else:
        template = load_template_safely(template_name)

    if not template:
        if template_file:
            click.echo(f"{FAILURE_INDICATOR} Template file '{template_file}' not found or invalid.")
        else:
            click.echo(f"{FAILURE_INDICATOR} Template '{template_name}' not found or invalid.")
        raise click.ClickException("Template not found")

    click.echo(f"Template: {template.name}")
    click.echo(f"Description: {template.description}")
    click.echo(f"Mode: {template.mode}")
    click.echo(f"API URL: {template.api_url}")
    click.echo(f"Method: {template.method}")
    click.echo()

    # Show usage information based on mode
    if template.mode == MODE_SINGLE:
        click.echo("Usage:")
        if template_file:
            click.echo(f"  archer validate --template-file {template_file} <secret>")
        else:
            click.echo(f"  archer validate {template_name} <secret>")
        click.echo()
    elif template.mode == MODE_MULTIPART:
        click.echo("Required Variables:")
        if template.required_variables:
            for var in template.required_variables:
                cli_name = format_var_name_for_cli(var)
                click.echo(f"  {var} ({OPT_VAR} {cli_name}=<value>)")
        click.echo()

        click.echo("Usage:")
        var_examples = []
        if template.required_variables:
            for var in template.required_variables:
                cli_name = format_var_name_for_cli(var)
                var_examples.append(f"{OPT_VAR} {cli_name}=<value>")

        if template_file:
            click.echo(f"  archer validate --template-file {template_file} {' '.join(var_examples)}")
        else:
            click.echo(f"  archer validate {template_name} {' '.join(var_examples)}")
        click.echo()

    click.echo("Request Headers:")
    for key, value in template.request.headers.items():
        if template.mode == MODE_SINGLE:
            secret_placeholder = f"${{{SECRET_VARIABLE_NAME}}}"
            masked_secret = f"{MASKED_VARIABLE_PREFIX}{SECRET_VARIABLE_NAME}{MASKED_VARIABLE_SUFFIX}"
            masked_value = value.replace(secret_placeholder, masked_secret) if secret_placeholder in value else value
        else:
            # For multipart, mask any variable
            def mask_variable(match):
                var_name = match.group(1)
                return MASKED_VARIABLE_FORMAT.format(var_name=var_name)
            masked_value = VARIABLE_PATTERN.sub(mask_variable, value)
        click.echo(f"  {key}: {masked_value}")

    if template.request.query_params:
        click.echo()
        click.echo("Query Parameters:")
        for key, value in template.request.query_params.items():
            if template.mode == MODE_SINGLE:
                secret_placeholder = f"${{{SECRET_VARIABLE_NAME}}}"
                masked_secret = f"{MASKED_VARIABLE_PREFIX}{SECRET_VARIABLE_NAME}{MASKED_VARIABLE_SUFFIX}"
                masked_value = value.replace(secret_placeholder, masked_secret) if secret_placeholder in value else value
            else:
                # For multipart, mask any variable
                def mask_variable(match):
                    var_name = match.group(1)
                    return MASKED_VARIABLE_FORMAT.format(var_name=var_name)
                masked_value = VARIABLE_PATTERN.sub(mask_variable, value)
            click.echo(f"  {key}: {masked_value}")

    click.echo()
    click.echo(f"Timeout: {template.request.timeout}s")
    click.echo()

    click.echo("Success Criteria:")
    click.echo(f"  Status Codes: {template.success_criteria.status_code}")
    if template.success_criteria.required_fields:
        click.echo(f"  Required Fields: {', '.join(template.success_criteria.required_fields)}")

    click.echo()
    click.echo("Error Handling:")
    click.echo(f"  Max Retries: {template.error_handling.max_retries}")
    click.echo(f"  Retry Delay: {template.error_handling.retry_delay}s")
    if template.error_handling.error_messages:
        click.echo("  Error Messages:")
        for code, message in template.error_handling.error_messages.items():
            click.echo(f"    {code}: {message}")