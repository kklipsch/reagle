FROM golang:1.11.0-alpine3.8 as basebuild

ARG REAGLE_LOCAL_LOCATION
ARG REAGLE_LOCAL_USER 
ARG REAGLE_LOCAL_PASSWORD 

ARG PROJECT_PATH=github.com/kklipsch/reagle
ARG CGO_ENABLED=0

#we use mod for getting our vendored dependencies but not for the build
#so to update deps you'd need to have go mod installed
ARG GO111MODULE=off

WORKDIR /go/src/$PROJECT_PATH
COPY . .

RUN go vet $PROJECT_PATH/local
RUN go test -v ./... 

