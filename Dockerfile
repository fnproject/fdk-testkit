# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
RUN go get -u github.com/golang/dep/cmd/dep
ENV D=/go/src/github.com/fnproject/fdk-testkit
ADD . $D
WORKDIR $D
RUN $GOPATH/bin/dep ensure
RUN go test -c -i  &&  cp fdk-testkit.test /tmp/

# final stage
FROM fnproject/dind
WORKDIR /app
COPY --from=build-env /tmp/fdk-testkit.test /app/fdk-testkit
CMD ["./fdk-testkit"]
