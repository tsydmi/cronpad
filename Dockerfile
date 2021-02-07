FROM node:lts-alpine as frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/babel.config.js frontend/vue.config.js ./
COPY frontend/public public
COPY frontend/src src
RUN npm run build


FROM golang:1.15-alpine AS backend-builder
ENV BUILD_DIR=$GOPATH/src/github.com/ts-dmitry/cronpad/backend/
COPY backend $BUILD_DIR
WORKDIR $BUILD_DIR
RUN CGO_ENABLED=0 go build -o /cronpad



FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates supervisor nginx
RUN addgroup cronpad && adduser --disabled-password cronpad -G cronpad

ENV CRONPAD_URL="http://cronpad:9000" \
    KEYCLOAK_URL="http://keycloak:8080"

# Prepare cronpad dir
RUN mkdir -p /cronpad/frontend && \
    mkdir -p /cronpad/backend && \
    chown -R cronpad.cronpad /cronpad

# Copy backend
COPY --from=backend-builder --chown=cronpad:cronpad /cronpad /cronpad/backend/

# Copy frontend
COPY --from=frontend-builder --chown=cronpad:cronpad /app/dist /cronpad/frontend

# Script to adjust frontend config to environment
COPY --chown=cronpad:cronpad frontend/docker/update-environment.sh /cronpad/frontend/
RUN chmod +x /cronpad/frontend/update-environment.sh

# Make sure files/folders needed by the nginx processes are accessable when they run under the cronpad user
RUN chown -R cronpad.cronpad /cronpad && \
    chown -R cronpad.cronpad /var/lib/nginx && \
    chown -R cronpad.cronpad /var/log/nginx && \
    touch /var/run/nginx.pid && \
    chown -R cronpad.cronpad /var/run/nginx.pid

USER cronpad

COPY --chown=cronpad:cronpad frontend/docker/nginx.conf /etc/nginx/nginx.conf
COPY --chown=cronpad:cronpad supervisord.conf /etc/supervisor/conf.d/supervisord.conf

HEALTHCHECK --interval=30s --timeout=3s --retries=6 CMD wget --no-verbose --tries=1 --spider http://localhost:9000/api/health || exit 1

ENTRYPOINT ["supervisord", "--nodaemon", "--configuration", "etc/supervisor/conf.d/supervisord.conf"]