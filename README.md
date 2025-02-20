# Demo App: outerspace-go

This is a demo app that uses the SpaceX API to fetch data about SpaceX launches, rockets, and capsules.

# Running

## How to run the go app locally

This will run a server on port `:8080` with the API endpoints.
```
go run main.go
```

## How to run the app in Kubernetes

## How to run the app in docker

```

```

# Building

## How to build the go binary locally

This will build the app and tag it with the current version.
```
go build
```

## How to build the docker image locally

Note this is just a local image, you need to tag and push to a container registry to share with others.
```
docker build . -t outerspace-go:latest
```
