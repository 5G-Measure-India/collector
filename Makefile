GOCMD=go
GOBUILD=${GOCMD} build

all: build

deps: go.mod go.sum
	${GOCMD} mod download

clean:
	${GOCMD} clean

build:
	${GOBUILD} .

release:
	CGOENABLED=1 ${GOBUILD} -ldflags="-s -w" -trimpath .

docker:
	docker build -t collector:dev .
