import click
from typing import Optional, Dict, Tuple

from archer.core.variables import parse_var_args, format_var_name_for_cli
from archer.templates import TemplateLoader
from archer.constants import MODE_SINGLE, MODE_MULTIPART, FAILURE_INDICATOR, OPT_VAR


class ArgumentHandler:
    """Handles CLI argument validation and processing."""

    def __init__(self):
        self.template_loader = TemplateLoader()

    def validate_and_process(self, template_name: str, secret: Optional[str], 
                           template_file: Optional[str], var_args: tuple) -> Tuple[Dict[str, str], any]:
        """Validate arguments and return processed variables and template."""
        # Load template to determine mode
        if template_file:
            template = self.template_loader.get_template(template_file)
        else:
            template = self.template_loader.get_template(template_name)

        if not template:
            error_msg = f"Template file '{template_file}' not found or invalid." if template_file else f"Template '{template_name}' not found."
            raise click.ClickException(error_msg)

        # Process based on template mode
        if template.mode == MODE_SINGLE:
            return self._handle_single_mode(template, secret, var_args, template_file, template_name)
        elif template.mode == MODE_MULTIPART:
            return self._handle_multipart_mode(template, secret, var_args, template_file, template_name)
        else:
            raise click.ClickException(f"Template mode '{template.mode}' is not supported.")

    def _handle_single_mode(self, template, secret: Optional[str], var_args: tuple, 
                          template_file: Optional[str], template_name: str) -> Tuple[Dict[str, str], any]:
        """Handle single mode validation."""
        if not secret:
            click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' is in single mode. Please provide a secret.")
            if template_file:
                click.echo(f"Usage: archer validate --template-file {template_file} <secret>")
            else:
                click.echo(f"Usage: archer validate {template_name} <secret>")
            raise click.ClickException("Secret required for single mode template")

        if var_args:
            click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' is in single mode. Variables are not supported.")
            if template_file:
                click.echo(f"Usage: archer validate --template-file {template_file} <secret>")
            else:
                click.echo(f"Usage: archer validate {template_name} <secret>")
            raise click.ClickException("Variables not supported for single mode template")

        variables = {'SECRET': secret}
        return variables, template

    def _handle_multipart_mode(self, template, secret: Optional[str], var_args: tuple,
                             template_file: Optional[str], template_name: str) -> Tuple[Dict[str, str], any]:
        """Handle multipart mode validation."""
        if secret:
            click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' is in multipart mode. Please use --var arguments.")
            self._show_multipart_usage(template, template_file, template_name)
            raise click.ClickException("Use --var arguments for multipart template")

        if not var_args:
            click.echo(f"{FAILURE_INDICATOR} Template '{template.name}' requires variables.")
            self._show_multipart_usage(template, template_file, template_name)
            raise click.ClickException("Variables required for multipart template")

        try:
            variables = parse_var_args(list(var_args))
            self._validate_required_variables(template, variables)
            self._check_unexpected_variables(template, variables)
            return variables, template
        except ValueError as e:
            click.echo(f"{FAILURE_INDICATOR} Variable parsing error: {str(e)}")
            raise click.ClickException("Invalid variable format")

    def _show_multipart_usage(self, template, template_file: Optional[str], template_name: str):
        """Show usage examples for multipart mode."""
        var_examples = []
        if template.required_variables:
            for var in template.required_variables:
                cli_name = format_var_name_for_cli(var)
                var_examples.append(f"{OPT_VAR} {cli_name}=<value>")

        if template_file:
            click.echo(f"Usage: archer validate --template-file {template_file} {' '.join(var_examples)}")
        else:
            click.echo(f"Usage: archer validate {template_name} {' '.join(var_examples)}")

    def _validate_required_variables(self, template, variables: Dict[str, str]):
        """Validate all required variables are provided."""
        missing_vars = []
        if template.required_variables:
            for required_var in template.required_variables:
                if required_var not in variables:
                    cli_name = format_var_name_for_cli(required_var)
                    missing_vars.append(cli_name)

        if missing_vars:
            click.echo(f"{FAILURE_INDICATOR} Missing required variables: {', '.join(missing_vars)}")
            raise click.ClickException("Missing required variables")

    def _check_unexpected_variables(self, template, variables: Dict[str, str]):
        """Check for unexpected variables."""
        unexpected_vars = []
        if template.required_variables:
            for provided_var in variables.keys():
                if provided_var not in template.required_variables:
                    cli_name = format_var_name_for_cli(provided_var)
                    unexpected_vars.append(cli_name)

        if unexpected_vars:
            click.echo(f"{FAILURE_INDICATOR} Unexpected variables: {', '.join(unexpected_vars)}")
            if template.required_variables:
                expected_vars = [format_var_name_for_cli(var) for var in template.required_variables]
                click.echo(f"Expected variables: {', '.join(expected_vars)}")
            raise click.ClickException("Unexpected variables provided")
