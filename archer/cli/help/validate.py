"""Help text for validate command."""

import textwrap


def get_help() -> str:
    """Get help text for validate command."""
    return textwrap.dedent("""\
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
