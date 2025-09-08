import click
import textwrap

from archer.cli.utils import discover_templates, load_template_safely
from archer.constants import MODE_SINGLE


@click.command()
def list_templates() -> None:
    """List all available templates from the default templates directory.

    Displays all built-in templates with their mode indicators and descriptions.
    Use this command to discover available validation templates before validation.
    """
    
    help_text = textwrap.dedent("""\
    \b
    Options:
      --help  Show this message and exit

    \b
    Output Format:
      Each template is displayed with:
      - Template name (left-aligned for easy reading)
      - Mode indicator: [single] or [multipart] 
      - Description of the API service

    \b
    Mode Indicators:
      [single]     - Use with: archer validate <template> <secret>
      [multipart]  - Use with: archer validate <template> --var key=value
      [invalid]    - Template has configuration errors

    \b
    Examples:
      $ archer list
      Available templates (25):
      
        github          [single]     - GitHub Personal Access Token validation
        ghost           [multipart]  - Ghost CMS API validation
        openai          [single]     - OpenAI API key validation

    \b
    Next Steps:
      1. Choose a template from the list
      2. Get detailed info: archer info <template_name>
      3. Run validation: archer validate <template_name> [arguments]

    See 'archer validate --help' for detailed usage patterns.
    """)
    
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
