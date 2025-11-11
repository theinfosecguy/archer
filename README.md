# Archer

<div align="center">
  <img width="250" height="250" alt="archer-logo" src="https://github.com/user-attachments/assets/d7b71c8b-eaef-46cf-adfe-a8bb6430b181" />
</div>



A fast, lightweight command-line tool for validating API secrets using YAML templates. Written in Go.

## Installation

**Download Binary:**

```bash
# Linux
curl -L https://github.com/theinfosecguy/archer/releases/latest/download/archer-linux-amd64 -o archer
chmod +x archer
sudo mv archer /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/theinfosecguy/archer/releases/latest/download/archer-darwin-arm64 -o archer
chmod +x archer
sudo mv archer /usr/local/bin/
```

**Or build from source:**

```bash
git clone https://github.com/theinfosecguy/archer.git
cd archer
make build
```

## Security Warning

**⚠️ Always use environment variables instead of passing secrets as command-line arguments.**

Secrets passed via CLI are exposed in shell history, process lists, and logs.

## Usage

### Basic Commands

```bash
# List all available templates
archer list

# Get information about a template
archer info github
```

### Validating Secrets

**Single Mode:**
```bash
export ARCHER_SECRET="ghp_xxxxxxxxxxxxxxx"
archer validate github
```

**Multipart Mode:**
```bash
export ARCHER_VAR_BASE_URL="https://myblog.com"
export ARCHER_VAR_API_TOKEN="xxxxx"
archer validate ghost
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

## Development

```bash
# Clone and build
git clone https://github.com/theinfosecguy/archer.git
cd archer
make build

# Run tests
make test

# Build for all platforms
make release-all
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
