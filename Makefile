all: build_linux build_image deploy_image

build_linux:
	GOOS=linux GOARCH=amd64 go build -o admission-webhook

build_image:
	docker buildx build --platform=linux/amd64 -t jusene/admission-webhook:v1 . --push

deploy:
	kubectl delete -f yaml/deployment.yaml
	kubectl label ns default admission-webhook-
	kubectl apply -f yaml/deployment.yaml
	kubectl label ns default admission-webhook=enabled