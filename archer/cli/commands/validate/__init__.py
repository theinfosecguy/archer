import click
import asyncio
import logging
from datetime import datetime, timezone
from typing import Optional, Dict

from archer.log_setup import setup_logging
from archer.cli.help import get_validate_help
from .argument_handler import ArgumentHandler
from .core import ValidationOrchestrator
from .output_handler import OutputHandler

logger = logging.getLogger(__name__)


@click.command()
@click.argument('template_name')
@click.argument('secret', required=False)
@click.option('--template-file', type=click.Path(exists=True, file_okay=True, dir_okay=False), help='Load template from specific file')
@click.option('--var', 'var_args', multiple=True, help='Variable in format key=value (for multipart templates)')
@click.option('--verbose', '-v', is_flag=True, help='Enable verbose logging')
@click.option('--debug', '-d', is_flag=True, help='Enable debug logging')
@click.option('--output-json', '-o', 'output_json', type=click.Path(dir_okay=False), help='Write structured validation result to JSON file (pretty, overwrite)')
@click.option('--json-only', is_flag=True, help='Suppress normal terminal success output when writing JSON (errors still shown)')
def validate(template_name: str, secret: Optional[str], template_file: Optional[str], var_args: tuple, 
           verbose: bool, debug: bool, output_json: Optional[str], json_only: bool) -> None:
    """Validate a secret using the specified template.

    Validates API secrets against endpoints based on template configuration.
    Template mode determines the required arguments and usage pattern.
    """

    setup_logging(verbose, debug)
    logger.info(f"Starting secret validation process for template '{template_name}'")

    start_time = datetime.now(timezone.utc)
    template_ref = None
    variables: Dict[str, str] = {}
    masked_url = None
    masked_headers = None
    masked_query_params = None

    # Initialize handlers
    arg_handler = ArgumentHandler()
    orchestrator = ValidationOrchestrator()
    output_handler = OutputHandler()

    async def _validate() -> None:
        nonlocal template_ref, variables, masked_url, masked_headers, masked_query_params

        # Validate and process arguments
        variables, template_ref = arg_handler.validate_and_process(
            template_name, secret, template_file, var_args)

        # Execute validation
        result = await orchestrator.execute_validation(
            template_name, secret, template_file, variables)

        # Build masked artifacts for JSON output
        if template_ref and variables:
            masked_url, masked_headers, masked_query_params = orchestrator.build_masked_artifacts(
                template_ref, variables)

        # Handle output
        if result.get("valid", False):
            output_handler.handle_success(result, output_json, json_only)
        else:
            output_handler.handle_error(result, output_json, json_only)
            raise click.ClickException("Secret validation failed")

        return result

    try:
        result = asyncio.run(_validate())
    except click.ClickException as e:
        if output_json:
            output_handler.write_json_output(
                output_json, {}, template_name, template_file, template_ref,
                variables, masked_url, masked_headers, masked_query_params,
                start_time, str(e), json_only)
        raise

    if output_json:
        output_handler.write_json_output(
            output_json, result, template_name, template_file, template_ref,
            variables, masked_url, masked_headers, masked_query_params,
            start_time, None, json_only)
