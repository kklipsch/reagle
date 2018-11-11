FROM golang:1.11.0-alpine3.8 as builder

ARG REAGLE_LOCAL_LOCATION
ARG REAGLE_LOCAL_USER 
ARG REAGLE_LOCAL_PASSWORD 
ARG REAGLE_IMPROVED_FIRMWARE
ARG REAGLE_MODEL_ID_NAME
ARG TEST_FLAG

ARG PROJECT_PATH=github.com/kklipsch/reagle
ARG CGO_ENABLED=0

#we use mod for getting our vendored dependencies but not for the build
#so to update deps you'd need to have go mod installed
ARG GO111MODULE=off

WORKDIR /go/src/$PROJECT_PATH
COPY . .

RUN go vet $PROJECT_PATH/local
RUN go vet $PROJECT_PATH/cmd/reagled
RUN go test $TEST_FLAG -v ./... 

RUN mkdir -p /out
RUN go build -o /out/reagled $PROJECT_PATH/cmd/reagled

FROM alpine:3.8

ENV REAGLE_LOCAL_LOCATION "localhost"
ENV REAGLE_LOCAL_USER "fake" 
ENV REAGLE_LOCAL_PASSWORD "fail"
ENV REAGLED_ADDRESS ":9000"
ENV REAGLE_IMPROVED_FIRMWARE "yes"
ENV REAGLE_MODEL_ID_NAME ""
EXPOSE 9000

WORKDIR /root/
COPY --from=builder /out/reagled /usr/local/bin/reagled
RUN apk --no-cache add ca-certificates

ENTRYPOINT ["reagled"]
