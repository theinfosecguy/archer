"""I/O utilities for Archer."""

import json
import os
from pathlib import Path
from typing import Any
from datetime import datetime

from archer.exceptions import JSONWriteError


class DateTimeEncoder(json.JSONEncoder):
    """JSON encoder that handles datetime objects."""
    def default(self, obj):
        if isinstance(obj, datetime):
            return obj.isoformat()
        return super().default(obj)


def write_json_file(path: str, data: Any) -> None:
    """Write JSON data (pretty) to a file path.

    Overwrites existing file. Creates parent directories if they do not exist.
    Raises JSONWriteError on failure.

    Args:
        path: Target file path
        data: Data to serialize as JSON

    Raises:
        JSONWriteError: If file cannot be written
    """
    try:
        target = Path(path)
        parent = target.parent
        parent.mkdir(parents=True, exist_ok=True)

        # Use atomic write via temp file
        tmp_path = target.with_suffix(target.suffix + ".tmp")
        with tmp_path.open('w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2, cls=DateTimeEncoder)
            f.write('\n')
        os.replace(tmp_path, target)
    except Exception as e:
        raise JSONWriteError(f"Failed to write JSON to {path}: {e}")
