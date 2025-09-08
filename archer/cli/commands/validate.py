import click
import asyncio
import logging
import textwrap

from typing import Optional

from archer.core import SecretValidator
from archer.core.variables import parse_var_args, format_var_name_for_cli
from archer.log_setup import setup_logging

from archer.templates import TemplateLoader
from archer.constants import (
    MODE_SINGLE,
    MODE_MULTIPART,
    SUCCESS_INDICATOR,
    FAILURE_INDICATOR,
    OPT_VAR,
    MISSING_REQUIRED_VARIABLES,
    UNEXPECTED_VARIABLES,
    SECRET_VALID,
    VALIDATION_SUCCESS,
    VALIDATION_FAILED,
)

logger = logging.getLogger(__name__)


@click.command()
@click.argument('template_name')
@click.argument('secret', required=False)
@click.option('--template-file', type=click.Path(exists=True, file_okay=True, dir_okay=False), help='Load template from specific file')
@click.option('--var', 'var_args', multiple=True, help='Variable in format key=value (for multipart templates)')
@click.option('--verbose', '-v', is_flag=True, help='Enable verbose logging')
@click.option('--debug', '-d', is_flag=True, help='Enable debug logging')
def validate(template_name: str, secret: Optional[str], template_file: Optional[str], var_args: tuple, verbose: bool, debug: bool) -> None:
    """Validate a secret using the specified template.

    Validates API secrets against endpoints based on template configuration.
    Template mode determines the required arguments and usage pattern.
    """
    
    help_text = textwrap.dedent("""\
    \b
    Arguments:
      TEMPLATE_NAME  Name of the template to use for validation
      SECRET         Secret to validate (only for single mode templates)

    \b
    Options:
      --template-file FILE      Load template from specific file instead of built-in
      --var TEXT               Variable in format key=value (for multipart templates only)
      --verbose, -v            Enable verbose logging
      --debug, -d              Enable debug logging
      --help                   Show this message and exit

    \b
    Usage Examples:

    \b
    1. SINGLE MODE TEMPLATES:
       Validate simple API tokens that require only the secret value.
       
       archer validate github ghp_xxxxxxxxxxxxxxxxxxxx
       archer validate openai sk-xxxxxxxxxxxxxxxxxxxxxxxx
       archer validate npm npm_xxxxxxxxxxxxxxxxxxxxxxxx

    \b
    2. MULTIPART MODE TEMPLATES:
       Validate APIs requiring multiple parameters using --var options.
       
       archer validate ghost --var base-url=https://myblog.com --var api-token=xxxxx
       archer validate stripe --var secret-key=sk_test_xxxxx --var publishable-key=pk_test_xxxxx
       archer validate supabase --var project-url=https://xxx.supabase.co --var anon-key=eyJxxx

    \b
    3. CUSTOM TEMPLATE FILES:
       Use your own YAML template files for validation.
       
       archer validate myapi --template-file ./custom-api.yaml sk_xxxxxxxxxxxxx
       archer validate custom --template-file ./multipart.yaml --var token=xxx --var url=https://api.example.com

    \b
    Variable Format Requirements:
      For multipart templates, variables must be provided in kebab-case format:
      
      ✓ Correct:   --var api-token=xxx --var base-url=https://example.com
      ✗ Incorrect: --var apiToken=xxx --var base_url=https://example.com
      
      Use 'archer info <template>' to see required variables for any template.
      Use 'archer list' to see which templates are single vs multipart mode.
    """)
    
    setup_logging(verbose, debug)
    logger.info(f"Starting secret validation process for template '{template_name}'")

    async def _validate() -> None:
        # Validate arguments
        if not template_file and not template_name:
            click.echo(f"{FAILURE_INDICATOR} Either template name or --template-file must be provided.")
            click.echo("Usage: archer validate <template_name> <secret>")
            click.echo("   or: archer validate --template-file <path> <secret>")
            raise click.ClickException("Missing template specification")

        # Load template to determine mode
        if template_file:
            # Load from specific file
            template_loader = TemplateLoader()
            template = template_loader.get_template(template_file)
            validation_identifier = template_file
        else:
            # Default directory
            template_loader = TemplateLoader()
            template = template_loader.get_template(template_name)
            validation_identifier = template_name

        if not template:
            if template_file:
                click.echo(f"{FAILURE_INDICATOR} Template file '{template_file}' not found or invalid.")
            else:
                click.echo(f"{FAILURE_INDICATOR} Template '{template_name}' not found.")
            raise click.ClickException("Template not found")

        # Use the same template loader for validation
        validator = SecretValidator()



        # Handle based on template mode
        if template.mode == MODE_SINGLE:
            if not secret:
                click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' is in single mode. Please provide a secret.")
                if template_file:
                    click.echo(f"Usage: archer validate --template-file {template_file} <secret>")
                else:
                    click.echo(f"Usage: archer validate {template_name} <secret>")
                raise click.ClickException("Secret required for single mode template")

            if var_args:
                click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' is in single mode. Variables are not supported.")
                if template_file:
                    click.echo(f"Usage: archer validate --template-file {template_file} <secret>")
                else:
                    click.echo(f"Usage: archer validate {template_name} <secret>")
                raise click.ClickException("Variables not supported for single mode template")

            result = await validator.validate_secret(validation_identifier, secret)

        elif template.mode == MODE_MULTIPART:
            if secret:
                click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' is in multipart mode. Please use --var arguments.")

                # Show expected variables
                var_examples = []
                if template.required_variables:
                    for var in template.required_variables:
                        cli_name = format_var_name_for_cli(var)
                        var_examples.append(f"{OPT_VAR} {cli_name}=<value>")

                # Show appropriate usage based on template source
                if template_file:
                    click.echo(f"Usage: archer validate --template-file {template_file} {' '.join(var_examples)}")
                else:
                    click.echo(f"Usage: archer validate {template_name} {' '.join(var_examples)}")
                raise click.ClickException("Use --var arguments for multipart template")

            if not var_args:
                click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' requires variables.")

                # Show expected variables
                var_examples = []
                if template.required_variables:
                    for var in template.required_variables:
                        cli_name = format_var_name_for_cli(var)
                        var_examples.append(f"{OPT_VAR} {cli_name}=<value>")

                # Show appropriate usage based on template source
                if template_file:
                    click.echo(f"Usage: archer validate --template-file {template_file} {' '.join(var_examples)}")
                else:
                    click.echo(f"Usage: archer validate {template_name} {' '.join(var_examples)}")
                raise click.ClickException("Variables required for multipart template")

            try:
                # Parse variables
                variables = parse_var_args(list(var_args))

                # Validate all required variables are provided
                missing_vars = []
                if template.required_variables:
                    for required_var in template.required_variables:
                        if required_var not in variables:
                            cli_name = format_var_name_for_cli(required_var)
                            missing_vars.append(cli_name)

                if missing_vars:
                    click.echo(f"{FAILURE_INDICATOR} {MISSING_REQUIRED_VARIABLES.format(vars=', '.join(missing_vars))}")
                    raise click.ClickException("Missing required variables")

                # Check for unexpected variables
                unexpected_vars = []
                if template.required_variables:
                    for provided_var in variables.keys():
                        if provided_var not in template.required_variables:
                            cli_name = format_var_name_for_cli(provided_var)
                            unexpected_vars.append(cli_name)

                if unexpected_vars:
                    click.echo(f"{FAILURE_INDICATOR} {UNEXPECTED_VARIABLES.format(vars=', '.join(unexpected_vars))}")
                    if template.required_variables:
                        expected_vars = [format_var_name_for_cli(var) for var in template.required_variables]
                        click.echo(f"Expected variables: {', '.join(expected_vars)}")
                    raise click.ClickException("Unexpected variables provided")

                result = await validator.validate_secret_multipart(validation_identifier, variables)

            except ValueError as e:
                click.echo(f"{FAILURE_INDICATOR} Variable parsing error: {str(e)}")
                raise click.ClickException("Invalid variable format")
        else:
            click.echo(f"{FAILURE_INDICATOR} Template mode '{template.mode}' is not supported.")
            raise click.ClickException("Unsupported template mode")

        if result.get("valid", False):
            click.echo(f"{SUCCESS_INDICATOR} {result.get('message', SECRET_VALID)}")
            logger.info(VALIDATION_SUCCESS)
        else:
            error_msg = result.get('error', 'Unknown validation error')
            click.echo(f"{FAILURE_INDICATOR} {error_msg}")
            logger.error(VALIDATION_FAILED.format(error=error_msg))
            raise click.ClickException("Secret validation failed")

    asyncio.run(_validate())
