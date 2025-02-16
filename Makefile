.PHONY: prod test rebuild

DC = docker-compose

prod:
	-$(DC) down
	$(DC) up -d

test:
	-$(DC) -f docker-compose.test.yml down -v
	$(DC) -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from tests
	$(DC) -f docker-compose.test.yml down -v

rebuild:
	-$(DC) down -v
	$(DC) up -d --build
