.PHONY: prod test rebuild

DC = docker-compose

prod:
	-$(DC) down
	$(DC) up -d

test:
	-go mod tidy
	-$(DC) -f docker-compose.test.yml down -v
	$(DC) -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from tests

rebuild:
	-$(DC) down -v
	$(DC) up -d --build
