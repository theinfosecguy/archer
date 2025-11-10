# Archer

<div align="center">
  <img width="250" height="250" alt="archer-logo" src="https://github.com/user-attachments/assets/d7b71c8b-eaef-46cf-adfe-a8bb6430b181" />
</div>



A fast, lightweight command-line tool for validating API secrets using YAML templates. Written in Go.

## Installation

### Download Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/theinfosecguy/archer/releases):

```bash
# Linux
curl -L https://github.com/theinfosecguy/archer/releases/latest/download/archer-linux-amd64 -o archer
chmod +x archer
sudo mv archer /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/theinfosecguy/archer/releases/latest/download/archer-darwin-amd64 -o archer
chmod +x archer
sudo mv archer /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/theinfosecguy/archer/releases/latest/download/archer-darwin-arm64 -o archer
chmod +x archer
sudo mv archer /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/theinfosecguy/archer.git
cd archer
make build
# Binary will be created as ./archer
```

### Using Go Install

```bash
go install github.com/theinfosecguy/archer/cmd/archer@latest
```

### Using Docker

```bash
# Build the image
docker build -t archer:latest .

# Run commands
docker run --rm archer:latest list
docker run --rm archer:latest info github
docker run --rm archer:latest validate github ghp_xxxxx

# Use with custom templates
docker run --rm -v $(pwd)/my-template.yaml:/app/my-template.yaml \
  archer:latest validate --template-file /app/my-template.yaml ghp_xxxxxxxxxxxxxxxxxxxx
```

## Security Warning

**⚠️ IMPORTANT: Passing secrets as command-line arguments is NOT secure!**

When you pass secrets directly on the command line, they are exposed in:
- Shell history files (`~/.bash_history`, `~/.zsh_history`)
- Process lists visible to other users (`ps`, `top`, `htop`)
- System logs and monitoring tools
- CI/CD logs and build outputs

**Always use environment variables for production and automation.**

## Usage

### Basic Commands

```bash
# List all available templates
archer list

# Get information about a template
archer info github
```

### Validating Secrets

#### Recommended: Using Environment Variables (Secure)

**Single Mode:**
```bash
# Set the secret as an environment variable
export ARCHER_SECRET="ghp_xxxxxxxxxxxxxxx"

# Validate without exposing the secret
archer validate github
```

**Multipart Mode:**
```bash
# Set multiple variables with ARCHER_VAR_ prefix
export ARCHER_VAR_BASE_URL="https://myblog.com"
export ARCHER_VAR_API_TOKEN="xxxxx"

# Validate without exposing secrets
archer validate ghost
```

#### Not Recommended: Command-Line Arguments (Insecure)

If you must use command-line arguments (not recommended), Archer will show a security warning:

**Single Mode:**
```bash
archer validate github ghp_xxxxxxxxxxxxxxx
# [WARNING] Passing secrets as command-line arguments is not secure...
```

**Multipart Mode:**
```bash
archer validate ghost --var base-url=https://myblog.com --var api-token=xxxxx
# [WARNING] Passing secrets as command-line arguments is not secure...
```

### Examples

#### Single Mode Examples (Using Environment Variables)

```bash
# GitHub
export ARCHER_SECRET="ghp_xxxxxxxxxxxxxxxxxxxx"
archer validate github

# OpenAI
export ARCHER_SECRET="sk-xxxxxxxxxxxxxxxx"
archer validate openai

# Slack
export ARCHER_SECRET="xoxb-xxxxxxxx-xxxxxxxx-xxxxxxxxxxxx"
archer validate slack

# Stripe
export ARCHER_SECRET="sk_test_xxxxxxxxxxxxxxxxxxxx"
archer validate stripe
```

#### Multipart Mode Examples (Using Environment Variables)

```bash
# Ghost CMS
export ARCHER_VAR_BASE_URL="https://myblog.com"
export ARCHER_VAR_API_TOKEN="xxxxx"
archer validate ghost

# Custom multipart template
export ARCHER_VAR_API_KEY="xxx"
export ARCHER_VAR_API_SECRET="yyy"
export ARCHER_VAR_ENDPOINT="https://api.example.com"
archer validate myservice
```

#### Custom Templates

```bash
# Use custom template file with environment variable
export ARCHER_SECRET="sk_xxxxxxxxxxxxx"
archer validate myapi --template-file ./custom-api.yaml

# Get info about custom template
archer info --template-file ./custom-template.yaml
```

### CI/CD Integration

For CI/CD pipelines, use environment variables from your secrets management:

```bash
# GitHub Actions
- name: Validate API Key
  env:
    ARCHER_SECRET: ${{ secrets.API_KEY }}
  run: archer validate github

# GitLab CI
validate:
  script:
    - export ARCHER_SECRET="${API_KEY}"
    - archer validate github

# Jenkins
sh 'ARCHER_SECRET=${API_KEY} archer validate github'
```

## Supported Services

Archer includes built-in templates for 26+ services:

| Service | Template Name | Mode |
|---------|--------------|------|
| **Airtable** | `airtable` | single |
| **Asana** | `asana` | single |
| **CircleCI** | `circleci` | single |
| **ClickUp** | `clickup` | single |
| **Codacy** | `codacy` | single |
| **Datadog** | `datadog` | single |
| **DigitalOcean** | `digitalocean` | single |
| **Discord** | `discord` | single |
| **Figma** | `figma` | single |
| **Ghost** | `ghost` | multipart |
| **GitHub** | `github` | single |
| **GitLab** | `gitlab` | single |
| **Heroku** | `heroku` | single |
| **JotForm** | `jotform` | single |
| **Linear** | `linear` | single |
| **Miro** | `miro` | single |
| **New Relic** | `newrelic` | single |
| **Notion** | `notion` | single |
| **npm** | `npm` | single |
| **OpenAI** | `openai` | single |
| **Postman** | `postman` | single |
| **Sentry** | `sentry` | single |
| **Slack** | `slack` | single |
| **Stripe** | `stripe` | single |
| **Supabase** | `supabase` | single |
| **Vercel** | `vercel` | single |

Run `archer list` to see all available templates.

## Creating Custom Templates

Create a YAML file with the following structure:

```yaml
name: myservice
description: "Validates MyService API keys"

api_url: "https://api.myservice.com/validate"
method: GET

# Mode: "single" or "multipart"
mode: single

request:
  headers:
    Authorization: "Bearer ${SECRET}"
    Content-Type: "application/json"
  timeout: 10

success_criteria:
  status_code: [200]
  required_fields:
    - "user_id"
    - "valid"

error_handling:
  max_retries: 2
  retry_delay: 1
  error_messages:
    401: "Invalid API key"
    403: "API key lacks required permissions"
```

For multipart templates (multiple variables):

```yaml
name: myservice
description: "Validates MyService with API key and secret"

api_url: "https://api.myservice.com/validate"
method: POST

mode: multipart
required_variables:
  - API_KEY
  - API_SECRET

request:
  headers:
    X-API-Key: "${API_KEY}"
  json_data:
    secret: "${API_SECRET}"
  timeout: 10

success_criteria:
  status_code: [200]
```

## Development

### Setup

```bash
# Clone the repository
git clone https://github.com/theinfosecguy/archer.git
cd archer

# Install development tools (goimports, etc.)
make install-tools
```

### Build

```bash
# Build for current platform
make build

# Build for all platforms
make release-all

# Build specific platform
make release-linux
make release-darwin
make release-windows
```

### Code Quality

```bash
# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Organize imports only
make imports

# Check if imports are properly formatted
make check-imports
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

---

**Get started now:**

```bash
archer list
```
