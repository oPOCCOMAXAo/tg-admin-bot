build-docker:
	docker build -t poccomaxa/tg-admin-bot:latest .

upload-docker:
	docker push poccomaxa/tg-admin-bot:latest

update-prod:
	docker tag poccomaxa/tg-admin-bot:latest poccomaxa/tg-admin-bot:prod
	docker push poccomaxa/tg-admin-bot:prod

prod: build-docker upload-docker update-prod
