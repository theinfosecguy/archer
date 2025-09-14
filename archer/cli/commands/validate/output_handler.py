import click
import logging
from datetime import datetime, timezone
from typing import Optional, Dict, Any

from archer.models.validation_output import ValidationResultJSON, ValidationRequestMeta, ValidationResponseMeta
from archer.utils.io import write_json_file
from archer.exceptions import JSONWriteError
from archer.constants import SUCCESS_INDICATOR, FAILURE_INDICATOR, SECRET_VALID, VALIDATION_SUCCESS, VALIDATION_FAILED, VERSION

logger = logging.getLogger(__name__)


class OutputHandler:
    """Handles terminal and JSON output."""

    def handle_success(self, result: Dict[str, Any], output_json: Optional[str], json_only: bool):
        """Handle successful validation output."""
        if not (output_json and json_only):
            click.echo(f"{SUCCESS_INDICATOR} {result.get('message', SECRET_VALID)}")
        logger.info(VALIDATION_SUCCESS)

    def handle_error(self, result: Dict[str, Any], output_json: Optional[str], json_only: bool):
        """Handle error output."""
        error_msg = result.get('error', 'Unknown validation error')
        if not (output_json and json_only):
            click.echo(f"{FAILURE_INDICATOR} {error_msg}")
        logger.error(VALIDATION_FAILED.format(error=error_msg))

    def write_json_output(self, output_json: str, result: Dict[str, Any], template_name: str,
                         template_file: Optional[str], template_ref, variables: Dict[str, str],
                         masked_url: Optional[str], masked_headers: Optional[Dict],
                         masked_query_params: Optional[Dict], start_time: datetime,
                         error: Optional[str] = None, json_only: bool = False):
        """Write JSON output file."""
        end_time = datetime.now(timezone.utc)

        validation_json = ValidationResultJSON(
            version=VERSION,
            valid=result.get('valid', False) if not error else False,
            message=result.get('message') if not error else None,
            error=result.get('error') if not error else error,
            request=ValidationRequestMeta(
                template=template_file if template_file else template_name,
                resolved_template_name=template_ref.name if template_ref else None,
                mode=template_ref.mode if template_ref else None,
                source=("file" if template_file else ("builtin" if template_ref else None)),
                method=template_ref.method if template_ref else None,
                api_url_masked=masked_url,
                headers_masked=masked_headers,
                query_params_masked=masked_query_params,
                variables_provided=list(variables.keys()) if variables else None,
                started_at=start_time,
                finished_at=end_time,
                duration_ms=(end_time - start_time).total_seconds() * 1000.0,
            ),
            response=ValidationResponseMeta(
                status_code=None,
                required_fields_checked=None,
                failed_required_field=None,
                error=result.get('error') if not error else error,
            )
        )

        try:
            write_json_file(output_json, validation_json.model_dump())
        except JSONWriteError as we:
            error_msg = f"Failed writing JSON output: {we}"
            if not json_only:
                click.echo(f"{FAILURE_INDICATOR} {error_msg}")
            else:
                click.echo(error_msg)
            raise click.ClickException(f"JSON write failed: {we}")
