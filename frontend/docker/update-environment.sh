#!/bin/sh
printf '{"VUE_APP_BACKEND_URL":"%s","VUE_APP_KEYCLOAK_URL":"%s"}\n' "$VUE_APP_BACKEND_URL" "$VUE_APP_KEYCLOAK_URL" > /app/environment.json