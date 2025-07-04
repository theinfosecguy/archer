import asyncio
import click
from archer.core import SecretValidator


@click.command()
@click.argument('template_name')
@click.argument('secret')
def validate_secret(template_name: str, secret: str) -> None:
    """Validate a secret using the specified template."""

    async def _validate() -> None:
        validator = SecretValidator()
        result = await validator.validate_secret(template_name, secret)

        if result["valid"]:
            click.echo(f"✅ {result['message']}")
        else:
            click.echo(f"❌ {result['error']}")
            raise click.ClickException("Secret validation failed")

    asyncio.run(_validate())


if __name__ == "__main__":
    validate_secret()
