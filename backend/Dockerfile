FROM golang:1.15-alpine AS builder
ENV BUILD_DIR=$GOPATH/src/github.com/ts-dmitry/cronpad/backend/
COPY . $BUILD_DIR
WORKDIR $BUILD_DIR
RUN CGO_ENABLED=0 go build -o /cronpad


FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN addgroup cronpad && adduser --disabled-password cronpad -G cronpad
USER cronpad
COPY --from=builder --chown=cronpad:cronpad /cronpad .
HEALTHCHECK --interval=30s --timeout=3s CMD curl --fail http://localhost:9000/api/health || exit 1
CMD ["./cronpad"]