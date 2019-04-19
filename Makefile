# Log into Bintray's public registry prior to using this Makefile
# docker login -u <user_login> -p <API_KEY> 

DOCKERREPO       := $(DOCKERREGISTRY)/blog-ser-ssl:1.0.1

all: blogservice container push

blogservice:
	docker run --rm -v $(shell pwd):/go/src/github.com/skckadiyala/blog-svc -w /go/src/github.com/skckadiyala/blog-svc -e CGO_ENABLED=0 golang:1.12.4-stretch go build -a -installsuffix cgo -o blogserver cmd/server/main.go

container: blogservice
	mkdir -p build/Dockerbuild && \
	cp -f blogserver build/Dockerbuild/ &&\
	cp -f start-service.sh build/Dockerbuild/ &&\
	cp -r third_party/swagger-ui/ build/Dockerbuild/swagger-ui/ &&\
	chmod +x build/Dockerbuild/blogserver &&\
	cp -f build/Dockerfile build/Dockerbuild/ &&\
	rm -f blogserver

	docker build -t $(DOCKERREPO) build/Dockerbuild/ 

push: container
	docker push $(DOCKERREPO)

