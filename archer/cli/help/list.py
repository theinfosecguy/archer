"""Help text for list command."""

import textwrap


def get_help() -> str:
    """Get help text for list command."""
    return textwrap.dedent("""\
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
