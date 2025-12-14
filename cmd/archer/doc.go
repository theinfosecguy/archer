/*
Archer is a fast, lightweight command-line tool for validating API secrets using YAML templates.

# Overview

Archer helps security researchers and developers quickly validate API keys and tokens
against 26+ popular services. It uses YAML templates to define validation requests
and supports both single-value secrets and multi-variable configurations.

# Installation

Install using Go:

	go install github.com/theinfosecguy/archer/cmd/archer@latest

Or download pre-built binaries from the releases page.

# Usage

List all available templates:

	archer list

Get information about a specific template:

	archer info github

Validate a secret (single mode):

	export ARCHER_SECRET="ghp_xxxxxxxxxxxxxxx"
	archer validate github

Validate with multiple variables (multipart mode):

	export ARCHER_VAR_BASE_URL="https://myblog.com"
	export ARCHER_VAR_API_TOKEN="xxxxx"
	archer validate ghost

# Security Warning

Always use environment variables instead of passing secrets as command-line arguments.
Secrets passed via CLI are exposed in shell history, process lists, and logs.

# Supported Services

Archer includes built-in templates for 26+ services including:
Airtable, Asana, CircleCI, ClickUp, Codacy, Datadog, DigitalOcean,
Discord, Figma, Ghost, GitHub, GitLab, Heroku, JotForm, Linear, Miro,
New Relic, Notion, npm, OpenAI, Postman, Sentry, Slack, Stripe,
Supabase, and Vercel.

Run 'archer list' to see all available templates with their validation modes.

# Template Modes

Templates operate in two modes:

Single Mode: Validates a single secret value using the ARCHER_SECRET environment variable.

Multipart Mode: Validates configurations requiring multiple variables (like base URLs
and tokens). Variables are passed using ARCHER_VAR_* environment variables.

Use 'archer info <template>' to see the required variables for each template.

# Creating Custom Templates

Templates are YAML files defining the HTTP request for validation.
See the examples/ directory and existing templates/ for reference.

Custom templates can be used by placing them in the templates directory
or by specifying a custom template path.

# Commands

	list        List all available templates
	validate    Validate a secret using a template
	info        Display information about a template
	version     Print version information

Use "archer [command] --help" for more information about a command.
*/
package main
