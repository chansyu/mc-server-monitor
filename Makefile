server:
	@docker compose up

watch-static:
	@npx tailwindcss -i ./ui/static/css/styles.css -o ./ui/static/css/output.css --watch

build:
	@npx tailwindcss -i ./ui/static/css/styles.css -o ./ui/static/css/output.css --minify
	@docker build -t mc-server-monitor-web:multistage -f ./cmd/web/Dockerfile .
	@docker build -t mc-server-monitor-logs:multistage -f ./cmd/logs/Dockerfile .

format:
	@gofmt -w .
	@npx prettier --write .