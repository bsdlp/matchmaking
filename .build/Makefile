all: doit

doit:
	docker build -t build .
	docker run --env-file ./env build
