FROM node:lts-alpine as builder
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY babel.config.js vue.config.js ./
COPY public public
COPY src src
RUN npm run build


FROM nginx:1.19.6-alpine
ENV VUE_APP_BACKEND_URL="http://localhost:9000/api/v1" \
    VUE_APP_KEYCLOAK_URL="http://localhost:8080/auth"
RUN mkdir /app
COPY --from=builder /app/dist /app
COPY docker/nginx.conf /etc/nginx/nginx.conf

COPY docker/update-environment.sh /docker-entrypoint.d/
RUN chmod +x /docker-entrypoint.d/*