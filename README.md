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

### For development:
```bash
git clone <repository>
cd archer
uv sync
```

## Usage

### With uv (development):
```bash
uv run archer validate <template_name> <secret>
uv run archer list
uv run archer info <template_name>
```

### With pip install:
```bash
archer validate <template_name> <secret>
archer list
archer info <template_name>
```

### Examples
```bash
# Development (with uv)
uv run archer validate github ghp_xxxxxxxxxxxxxxx
uv run archer validate openai sk-xxxxxxxxxxxxxxxx
uv run archer validate slack xoxb-xxxxxxxxxx

# Production (with pip)
archer validate github ghp_xxxxxxxxxxxxxxx
archer validate openai sk-xxxxxxxxxxxxxxxx
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
