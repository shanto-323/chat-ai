migration_new: 
	@if [ -z "$(name)" ]; then \
		@echo "Error: name parameter is required" \
        @echo "Usage: task migrations:new name=migration_name" \
		exit 1; \
	fi
	@echo "Creating migration file for $(name)"
	@tern new -m ./internal/database/migrations ${name}


confirm:
	@read -p "Press [y/Y] to continue: " value; \
	value=$$(echo $$value | tr '[:upper:]' '[:lower:]'); \
	if [ "$$value" = "y" ]; then \
		echo "Running target..."; \
	else \
		exit 1; \
	fi
