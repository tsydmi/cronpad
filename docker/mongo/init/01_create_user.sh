#!/bin/bash
set -e;

MONGO_INITDB_USER_ROLE="${MONGO_INITDB_USER_ROLE:-readWrite}"

if [ -n "${MONGO_INITDB_USERNAME:-}" ] && [ -n "${MONGO_INITDB_PASSWORD:-}" ]; then
	"${mongo[@]}" "$MONGO_INITDB_DATABASE" <<-EOJS
		db.createUser({
			user: $(_js_escape "$MONGO_INITDB_USERNAME"),
			pwd: $(_js_escape "$MONGO_INITDB_PASSWORD"),
			roles: [ { role: $(_js_escape "$MONGO_INITDB_USER_ROLE"), db: $(_js_escape "$MONGO_INITDB_DATABASE") } ]
			})
	EOJS
else
    exit 1
fi