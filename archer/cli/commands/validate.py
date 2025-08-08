import asyncio
import logging
import click
from typing import Optional
from archer.core import SecretValidator
from archer.core.variables import parse_var_args, format_var_name_for_cli
from archer.log_setup import setup_logging
from archer.templates import TemplateLoader

logger = logging.getLogger(__name__)


@click.command()
@click.argument('template_name')
@click.argument('secret', required=False)
@click.option('--var', 'var_args', multiple=True, help='Variable in format key=value (for multipart templates)')
@click.option('--verbose', '-v', is_flag=True, help='Enable verbose logging')
@click.option('--debug', '-d', is_flag=True, help='Enable debug logging')
def validate(template_name: str, secret: Optional[str], var_args: tuple, verbose: bool, debug: bool) -> None:
    """Validate a secret using the specified template.

    For single mode templates:
        archer validate github <secret>

    For multipart mode templates:
        archer validate ghost --var base-url=https://myblog.com --var api-token=<token>
    """
    setup_logging(verbose, debug)
    logger.info(f"Starting secret validation process for template '{template_name}'")

    async def _validate() -> None:
        # Load template to determine mode
        template_loader = TemplateLoader()
        template = template_loader.get_template(template_name)

        if not template:
            click.echo(f"❌ Template '{template_name}' not found.")
            raise click.ClickException("Template not found")

        validator = SecretValidator()

        # Handle based on template mode
        if template.mode == 'single':
            if not secret:
                click.echo(f"❌ Template '{template_name}' is in single mode. Please provide a secret.")
                click.echo(f"Usage: archer validate {template_name} <secret>")
                raise click.ClickException("Secret required for single mode template")

            if var_args:
                click.echo(f"❌ Template '{template_name}' is in single mode. Variables are not supported.")
                click.echo(f"Usage: archer validate {template_name} <secret>")
                raise click.ClickException("Variables not supported for single mode template")

            result = await validator.validate_secret(template_name, secret)

        elif template.mode == 'multipart':
            if secret:
                click.echo(f"❌ Template '{template_name}' is in multipart mode. Please use --var arguments.")

                # Show expected variables
                var_examples = []
                for var in template.required_variables:
                    cli_name = format_var_name_for_cli(var)
                    var_examples.append(f"--var {cli_name}=<value>")

                click.echo(f"Usage: archer validate {template_name} {' '.join(var_examples)}")
                raise click.ClickException("Use --var arguments for multipart template")

            if not var_args:
                click.echo(f"❌ Template '{template_name}' requires variables.")

                # Show expected variables
                var_examples = []
                for var in template.required_variables:
                    cli_name = format_var_name_for_cli(var)
                    var_examples.append(f"--var {cli_name}=<value>")

                click.echo(f"Usage: archer validate {template_name} {' '.join(var_examples)}")
                raise click.ClickException("Variables required for multipart template")

            try:
                # Parse variables
                variables = parse_var_args(var_args)

                # Validate all required variables are provided
                missing_vars = []
                for required_var in template.required_variables:
                    if required_var not in variables:
                        cli_name = format_var_name_for_cli(required_var)
                        missing_vars.append(cli_name)

                if missing_vars:
                    click.echo(f"❌ Missing required variables: {', '.join(missing_vars)}")
                    raise click.ClickException("Missing required variables")

                # Check for unexpected variables
                unexpected_vars = []
                for provided_var in variables.keys():
                    if provided_var not in template.required_variables:
                        cli_name = format_var_name_for_cli(provided_var)
                        unexpected_vars.append(cli_name)

                if unexpected_vars:
                    click.echo(f"❌ Unexpected variables: {', '.join(unexpected_vars)}")
                    expected_vars = [format_var_name_for_cli(var) for var in template.required_variables]
                    click.echo(f"Expected variables: {', '.join(expected_vars)}")
                    raise click.ClickException("Unexpected variables provided")

                result = await validator.validate_secret_multipart(template_name, variables)

            except ValueError as e:
                click.echo(f"❌ Variable parsing error: {str(e)}")
                raise click.ClickException("Invalid variable format")

        if result["valid"]:
            click.echo(f"✅ {result['message']}")
            logger.info("Secret validation process completed successfully")
        else:
            click.echo(f"❌ {result['error']}")
            logger.error(f"Secret validation process failed: {result['error']}")
            raise click.ClickException("Secret validation failed")

    asyncio.run(_validate())
