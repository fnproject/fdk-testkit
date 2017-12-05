# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ENV D=/go/src/github.com/fnproject/fdk-testkit
ADD . $D
RUN cd $D && go test -c -i  && cp fdk-testkit /tmp/

# final stage
FROM fnproject/dind
WORKDIR /app
COPY --from=build-env /tmp/fdk-testkit /app/fdk-testkit
CMD ["./fdk-testkit"]
