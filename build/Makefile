ROOT=$(shell (dirname $(shell pwd )))

BINARY_PATH=$(ROOT)/build/binary

# MAIN=$(ROOT)/main.go
MAIN=$(ROOT)/cmd/$(NAME)/main.go

PACKAGE_CMD=rxt/cmd/rxsc/app/cmd

VERSION=0.0.2
NAME=rxsc
ENV=local

# 本地运行
run:
	go run $(MAIN) server --name=$(NAME) -v=$(VERSION) -e=$(ENV)
build:
	go build -ldflags \
	"-X '$(PACKAGE_CMD).version=$(VERSION)' \
	-X '$(PACKAGE_CMD).env=$(ENV)' \
	-X '$(PACKAGE_CMD).name=$(NAME)'" \
	-o $(BINARY_PATH)/$(NAME) \
	$(MAIN)

# 制作镜像
deploy:docker-deploy
docker-build:
	docker run -it --rm -v $(ROOT):/go/src mengino/go:build make build

docker-release:docker-build
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) -t $(NAME):$(VERSION) .

docker-deploy:docker-release
	docker rm -f $(NAME)
	docker run -dp 8080:8080 --name $(NAME) $(NAME):$(VERSION) $(NAME) server
