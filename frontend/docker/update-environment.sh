#!/bin/sh
printf '{"VUE_APP_BACKEND_URL":"%s","VUE_APP_KEYCLOAK_URL":"%s"}\n' "$CRONPAD_URL/api/v1" "$KEYCLOAK_URL/auth" > /cronpad/frontend/environment.json