ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-alpine AS build
ARG APP_NAME=healthcheckr
RUN apk update --no-cache
RUN apk add --no-cache \
  ca-certificates \
  g++ \
  git \
  make
RUN go install github.com/swaggo/swag/cmd/swag@master
WORKDIR /go/src/${APP_NAME}
COPY ./go.mod ./go.sum ./
RUN go mod download -x
COPY ./.git ./.git
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./Makefile .
COPY ./main.go .
ENV CGO_ENABLED=0
# docs not available as of 2024-03-29
# RUN make docs-swaggo
RUN make deps
RUN make binary
RUN sha256sum ./bin/${APP_NAME} > ./bin/${APP_NAME}.sha256
RUN chmod +x ./bin/${APP_NAME}

FROM scratch AS final
ARG APP_NAME=healthcheckr
ENV APP_NAME=${APP_NAME}
COPY --from=build /go/src/${APP_NAME}/bin/${APP_NAME} /entrypoint
COPY --from=build /go/src/${APP_NAME}/bin/${APP_NAME}.sha256 /${APP_NAME}.sha256
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# CMD [ /bin/bash ]
ENTRYPOINT [ "/entrypoint" ]
