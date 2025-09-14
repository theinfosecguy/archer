import click

from archer.utils import discover_templates, load_template_safely
from archer.cli.help import get_list_help
from archer.constants import MODE_SINGLE


@click.command()
def list_templates() -> None:
    """List all available templates from the default templates directory.

    Displays all built-in templates with their mode indicators and descriptions.
    Use this command to discover available validation templates before validation.
    """
    
    templates = discover_templates()

    if not templates:
        click.echo("No templates found.")
        return

    click.echo(f"Available templates ({len(templates)}):")
    click.echo()

    for template_name in sorted(templates):
        template = load_template_safely(template_name)

        if template:
            mode_indicator = f"[{template.mode}]" if template.mode else f"[{MODE_SINGLE}]"
            click.echo(f"  {template.name:<15} {mode_indicator:<12} - {template.description}")
        else:
            click.echo(f"  {template_name:<15} {'[invalid]':<12} - [Invalid template]")
