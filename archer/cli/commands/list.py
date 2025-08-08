import click
from archer.cli.utils import discover_templates, load_template_safely


@click.command()
def list_templates() -> None:
    """List all available templates."""
    templates = discover_templates()

    if not templates:
        click.echo("No templates found.")
        return

    click.echo(f"Available templates ({len(templates)}):")
    click.echo()

    for template_name in sorted(templates):
        template = load_template_safely(template_name)
        if template:
            mode_indicator = f"[{template.mode}]" if template.mode else "[single]"
            click.echo(f"  {template.name:<15} {mode_indicator:<12} - {template.description}")
        else:
            click.echo(f"  {template_name:<15} {'[invalid]':<12} - [Invalid template]")
