import re
import click
from archer.cli.utils import load_template_safely
from archer.core.variables import format_var_name_for_cli


@click.command()
@click.argument('template_name')
def info(template_name: str) -> None:
    """Show detailed information about a template."""
    template = load_template_safely(template_name)

    if not template:
        click.echo(f"❌ Template '{template_name}' not found or invalid.")
        raise click.ClickException("Template not found")

    click.echo(f"Template: {template.name}")
    click.echo(f"Description: {template.description}")
    click.echo(f"Mode: {template.mode}")
    click.echo(f"API URL: {template.api_url}")
    click.echo(f"Method: {template.method}")
    click.echo()

    # Show usage information based on mode
    if template.mode == 'single':
        click.echo("Usage:")
        click.echo(f"  archer validate {template_name} <secret>")
        click.echo()
    elif template.mode == 'multipart':
        click.echo("Required Variables:")
        for var in template.required_variables:
            cli_name = format_var_name_for_cli(var)
            click.echo(f"  {var} (--var {cli_name}=<value>)")
        click.echo()
        
        click.echo("Usage:")
        var_examples = []
        for var in template.required_variables:
            cli_name = format_var_name_for_cli(var)
            var_examples.append(f"--var {cli_name}=<value>")
        click.echo(f"  archer validate {template_name} {' '.join(var_examples)}")
        click.echo()

    click.echo("Request Headers:")
    for key, value in template.request.headers.items():
        if template.mode == 'single':
            masked_value = value.replace("${SECRET}", "***MASKED***") if "${SECRET}" in value else value
        else:
            # For multipart, mask any variable
            def mask_variable(match):
                var_name = match.group(1)
                return f"${{***{var_name}_MASKED***}}"
            masked_value = re.sub(r'\$\{([^}]+)\}', mask_variable, value)
        click.echo(f"  {key}: {masked_value}")

    if template.request.query_params:
        click.echo()
        click.echo("Query Parameters:")
        for key, value in template.request.query_params.items():
            if template.mode == 'single':
                masked_value = value.replace("${SECRET}", "***MASKED***") if "${SECRET}" in value else value
            else:
                # For multipart, mask any variable
                def mask_variable(match):
                    var_name = match.group(1)
                    return f"${{***{var_name}_MASKED***}}"
                masked_value = re.sub(r'\$\{([^}]+)\}', mask_variable, value)
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