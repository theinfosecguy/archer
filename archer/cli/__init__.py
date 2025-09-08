import click
import textwrap

from archer.cli.commands.validate import validate
from archer.cli.commands.list import list_templates
from archer.cli.commands.info import info
from archer.constants import VERSION

CONTEXT_SETTINGS = dict(
    help_option_names=["-h", "--help"],
    max_content_width=100,  # adjust to taste
)

HELP = "Archer — validate secrets against APIs using YAML templates."

EPILOG = textwrap.dedent("""\
\b
Usage patterns

\b
1) Single-mode templates (simple API tokens)
   archer validate github ghp_xxxxxxxxxxxxxxxxxxxx
   archer validate openai sk-xxxxxxxxxxxxxxxxxxxxxxxx
   archer validate slack xoxb-xxxxxxxx-xxxxxxxx-xxxxxxxxxxxx

\b
2) Multipart templates (multiple parameters)
   archer validate ghost --var base-url=https://myblog.com --var api-token=xxxxx
   archer validate stripe --var secret-key=sk_test_xxxxx --var publishable-key=pk_test_xxxxx

\b
3) Custom template files
   archer validate myapi --template-file ./custom-api.yaml sk_xxxxxxxxxxxxx
   archer validate custom --template-file ./multipart.yaml --var token=xxx --var url=https://api.example.com

\b
4) Template information
   archer list
   archer info github
   archer info --template-file ./custom.yaml

Get started with: archer list
"""
)

@click.group(context_settings=CONTEXT_SETTINGS, help=HELP, epilog=EPILOG)
@click.version_option(version=VERSION)
def cli() -> None:
    pass

cli.add_command(validate)
cli.add_command(list_templates, name="list")
cli.add_command(info)
