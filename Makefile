
docker-build:
	docker build \
		-t bypass-cors \
		--label latest \
		--label `git rev-parse HEAD` \
		.

docker-run:
	docker run \
		-p 3228:3228 \
		bypass-cors -pp
