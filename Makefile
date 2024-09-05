
build:
	go build -o ./dist/bin/flight-metrics .

docker:
	docker build -t gitsang/flight-metrics:latest .

push:
	docker push gitsang/flight-metrics:latest

run:
	docker rm -f flight-metrics
	docker run -d \
		--name flight-metrics \
		--network host \
		-v $(PWD)/configs:/configs \
		-e HTTP_PROXY=http://127.0.0.1:7890 \
		-e HTTPS_PROXY=http://127.0.0.1:7890 \
		gitsang/flight-metrics:latest \
			-c /configs/config.yml
