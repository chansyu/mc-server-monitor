server:
	@docker compose up

build:
	@npx tailwindcss -i ./ui/static/css/styles.css -o ./ui/static/css/output.css

watch-static:
	@npx tailwindcss -i ./ui/static/css/styles.css -o ./ui/static/css/output.css --watch