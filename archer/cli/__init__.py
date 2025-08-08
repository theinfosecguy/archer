import click
from archer.cli.commands.validate import validate
from archer.cli.commands.list import list_templates
from archer.cli.commands.info import info


@click.group()
@click.version_option(version="0.1.0")
def cli() -> None:
    """Archer - Secret validation system using YAML templates.

    Validate secrets against various APIs using configurable templates.

    Examples:
        # Single mode templates
        archer validate github ghp_xxxxxxxxxxxx
        archer validate openai sk-xxxxxxxxxxxxxxxx

        # Multipart mode templates
        archer validate ghost --var base-url=https://myblog.com --var api-token=xxxxx

        # Other commands
        archer list
        archer info github
        archer info ghost
    """
    pass


cli.add_command(validate)
cli.add_command(list_templates, name="list")
cli.add_command(info)
