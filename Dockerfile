# Use Python 3.11 slim image as base
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Install system dependencies (if needed)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy project files
COPY pyproject.toml ./
COPY README.md ./
COPY archer/ ./archer/
COPY templates/ ./templates/
COPY examples/ ./examples/

# Install Python dependencies
RUN pip install --no-cache-dir -e .

# Create a non-root user to run the application
RUN useradd -m -u 1000 archer && \
    chown -R archer:archer /app

# Switch to non-root user
USER archer

# Set the entrypoint to the archer CLI
ENTRYPOINT ["archer"]

# Default command (show help)
CMD ["--help"]
