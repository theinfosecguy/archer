.PHONY: help validate-github validate-openai validate-gitlab validate-slack list info

help:
	@echo "Archer Secret Validation Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make validate-github TOKEN=your_token"
	@echo "  make validate-openai TOKEN=your_token"
	@echo "  make validate-gitlab TOKEN=your_token"
	@echo "  make validate-slack TOKEN=your_token"
	@echo ""
	@echo "Other commands:"
	@echo "  make list                    # List all templates"
	@echo "  make info TEMPLATE=name      # Show template info"

validate-github:
	@if [ -z "$(TOKEN)" ]; then echo "Usage: make validate-github TOKEN=your_token"; exit 1; fi
	uv run archer validate github $(TOKEN)

validate-openai:
	@if [ -z "$(TOKEN)" ]; then echo "Usage: make validate-openai TOKEN=your_token"; exit 1; fi
	uv run archer validate openai $(TOKEN)

validate-gitlab:
	@if [ -z "$(TOKEN)" ]; then echo "Usage: make validate-gitlab TOKEN=your_token"; exit 1; fi
	uv run archer validate gitlab $(TOKEN)

validate-slack:
	@if [ -z "$(TOKEN)" ]; then echo "Usage: make validate-slack TOKEN=your_token"; exit 1; fi
	uv run archer validate slack $(TOKEN)

list:
	uv run archer list

info:
	@if [ -z "$(TEMPLATE)" ]; then echo "Usage: make info TEMPLATE=template_name"; exit 1; fi
	uv run archer info $(TEMPLATE)
