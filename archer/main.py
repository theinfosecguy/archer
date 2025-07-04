import asyncio
import logging
import click
from archer.core import SecretValidator
from archer.logging import setup_logging



@click.command()
@click.argument('template_name')
@click.argument('secret')
@click.option('--verbose', '-v', is_flag=True, help='Enable verbose logging')
@click.option('--debug', '-d', is_flag=True, help='Enable debug logging')
def validate_secret(template_name: str, secret: str, verbose: bool, debug: bool) -> None:
    """Validate a secret using the specified template."""
    setup_logging(verbose, debug)
    logger = logging.getLogger(__name__)

    logger.info(f"Starting secret validation process for template '{template_name}'")

    async def _validate() -> None:
        validator = SecretValidator()
        result = await validator.validate_secret(template_name, secret)

        if result["valid"]:
            click.echo(f"✅ {result['message']}")
            logger.info("Secret validation process completed successfully")
        else:
            click.echo(f"❌ {result['error']}")
            logger.error(f"Secret validation process failed: {result['error']}")
            raise click.ClickException("Secret validation failed")

    asyncio.run(_validate())


if __name__ == "__main__":
    validate_secret()
