ifneq (,$(wildcard ./.env))
    include .env
    export
endif

frontendLocation := ./frontend-bundle
distDir := ./dist

ifneq ($(ENVIRONMENT),PRODUCTION)
dev:
	@echo "Starting frontend and backend in development mode..."
	@npm run dev --prefix $(frontendLocation) & \
	air . & \
	wait

else
create_dist_dir:
	@if [ ! -d "$(distDir)" ]; then \
		echo "Creating $(distDir) directory..."; \
		mkdir -p $(distDir); \
	fi

copy_static_files: create_dist_dir
	@echo "Copying static files..."
	@if [ -f .env ]; then cp .env $(distDir)/; fi
	@if [ -d static ]; then cp -R static $(distDir)/; fi

build: create_dist_dir copy_static_files
	@echo "Building frontend and backend for production..."
	@npm run build --prefix $(frontendLocation)
	@go build -o $(distDir)/main
	@echo "Build completed. Output is in $(distDir) directory."
endif

run-cli:
ifneq ($(ENVIRONMENT),PRODUCTION)
	@make dev
else
	@make build
	@echo "Starting the application in production mode..."
	@$(distDir)/main
endif

run-docker:
	docker compose up