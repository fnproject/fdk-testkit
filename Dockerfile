# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ENV D=/go/src/github.com/fnproject/fdk-testkit
RUN go get -u github.com/golang/dep/cmd/dep
ADD Gopkg.* $D/
RUN cd $D && dep ensure --vendor-only
ADD . $D
RUN cd $D && go test -c -i  &&  cp fdk-testkit.test /tmp/

# final stage
FROM fnproject/dind
WORKDIR /app
COPY --from=build-env /tmp/fdk-testkit.test /app/fdk-testkit
CMD ["./fdk-testkit"]
