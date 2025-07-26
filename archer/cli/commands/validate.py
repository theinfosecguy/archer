import asyncio
import logging
import click
from typing import Dict
from archer.core import SecretValidator
from archer.logging import setup_logging

logger = logging.getLogger(__name__)


def parse_var_option(ctx, param, value) -> Dict[str, str]:
    """Parse --var options into a dictionary."""
    if not value:
        return {}
    
    variables = {}
    for var_pair in value:
        if '=' not in var_pair:
            raise click.BadParameter(f"Variable '{var_pair}' must be in format 'name=value'")
        
        name, val = var_pair.split('=', 1)
        name = name.upper().replace('-', '_')
        variables[name] = val
    
    return variables


@click.command()
@click.argument('template_name')
@click.argument('secret', required=False) 
@click.option('--var', multiple=True, callback=parse_var_option, 
              help='Variable for multipart templates (format: name=value). Can be used multiple times.')
@click.option('--verbose', '-v', is_flag=True, help='Enable verbose logging')
@click.option('--debug', '-d', is_flag=True, help='Enable debug logging')
def validate(template_name: str, secret: str, var: Dict[str, str], verbose: bool, debug: bool) -> None:
    """Validate a secret using the specified template.
    
    For single-part templates:
        archer validate airtable YOUR_TOKEN
        
    For multi-part templates:
        archer validate ghost --var base-url=https://myblog.ghost.io --var api-token=abc123
    """
    setup_logging(verbose, debug)
    logger.info(f"Starting secret validation process for template '{template_name}'")

    async def _validate() -> None:
        validator = SecretValidator()
        
        # If variables provided, use multipart validation
        if var:
            if secret:
                raise click.ClickException("Cannot use both positional secret and --var options together")
            result = await validator.validate_secret_multipart(template_name, var)
        else:
            # Use single secret validation
            if not secret:
                raise click.ClickException("Secret argument is required when not using --var options")
            result = await validator.validate_secret(template_name, secret)

        if result["valid"]:
            click.echo(f"✅ {result['message']}")
            logger.info("Secret validation process completed successfully")
        else:
            click.echo(f"❌ {result['error']}")
            logger.error(f"Secret validation process failed: {result['error']}")
            raise click.ClickException("Secret validation failed")

    asyncio.run(_validate())
