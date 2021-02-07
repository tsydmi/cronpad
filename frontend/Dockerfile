FROM node:lts-alpine as frontend-builder
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY babel.config.js vue.config.js ./
COPY public public
COPY src src
RUN npm run build


FROM nginx:1.19.6-alpine

RUN addgroup cronpad && adduser --disabled-password cronpad -G cronpad

ENV VUE_APP_BACKEND_URL="http://localhost:9000/api/v1" \
    VUE_APP_KEYCLOAK_URL="http://localhost:8080/auth"

# Prepare cronpad dir
RUN mkdir -p /cronpad/frontend && \
    chown -R cronpad.cronpad /docker-entrypoint.d/

# Copy frontend
COPY --from=frontend-builder --chown=cronpad:cronpad /app/dist /cronpad/frontend

# Script to adjust frontend config to environment
COPY --chown=cronpad:cronpad docker/update-environment.sh /docker-entrypoint.d/
COPY --chown=cronpad:cronpad docker/nginx.conf /etc/nginx/nginx.conf

# Make sure files/folders needed by the nginx processes are accessable when they run under the cronpad user
RUN chown -R cronpad.cronpad /cronpad && \
    chown -R cronpad.cronpad /var/cache/nginx && \
    chown -R cronpad.cronpad /var/log/nginx && \
    touch /var/run/nginx.pid && \
    chown -R cronpad.cronpad /var/run/nginx.pid && \
    chown -R cronpad.cronpad /docker-entrypoint.d/ && \
    chmod +x /docker-entrypoint.d/*

USER cronpad