# Archer

<div align="center">
  <img width="250" height="250" alt="archer-logo" src="https://github.com/user-attachments/assets/d7b71c8b-eaef-46cf-adfe-a8bb6430b181" />
</div>

A command-line tool for validating API secrets using YAML templates.

## Installation

### For end users:
```bash
pip install archer
```

### Using Docker:
```bash
# Build the image
docker build -t archer:latest .

# Run commands
docker run --rm archer:latest list
docker run --rm archer:latest info github
docker run --rm archer:latest validate github ghp_xxxxx

# Use with custom templates
docker run --rm -v $(pwd)/my-template.yaml:/app/my-template.yaml \
  archer:latest validate github ghp_xxxxxxxxxxxxxxxxxxxx
```

### For development:
```bash
git clone <repository>
cd archer
uv sync
```

## Usage

### Basic usage:
```bash
archer validate <template_name> <secret>
archer list
archer info <template_name>
```

### Custom templates:
```bash
archer validate --template-file ./my-template.yaml <secret>
archer info --template-file ./my-template.yaml
```

### Examples
```bash
# Basic usage
archer validate github ghp_xxxxxxxxxxxxxxx
archer validate openai sk-xxxxxxxxxxxxxxxx

# Custom template file
archer validate --template-file ./my-github-template.yaml ghp_xxxxxxxxxxxxxxx
archer info --template-file ./my-github-template.yaml
```

## Supported Secrets

- **airtable** - Airtable API keys
- **asana** - Asana personal access tokens
- **circleci** - CircleCI API tokens
- **clickup** - ClickUp API tokens
- **codacy** - Codacy API tokens
- **digitalocean** - DigitalOcean API tokens
- **discord** - Discord bot tokens
- **figma** - Figma personal access tokens
- **github** - GitHub personal access tokens
- **gitlab** - GitLab personal access tokens
- **heroku** - Heroku API keys
- **jotform** - JotForm API keys
- **linear** - Linear API keys
- **notion** - Notion integration tokens
- **npm** - npm access tokens
- **openai** - OpenAI API keys
- **postman** - Postman API keys
- **slack** - Slack bot tokens
- **stripe** - Stripe API keys
- **supabase** - Supabase API keys
- **vercel** - Vercel API tokens
