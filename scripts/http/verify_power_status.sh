#!/usr/bin/env bash
# Verify by checking Power Status
# verify_power_status.sh deviceip

set -euo pipefail

exec 1>&2

host=${HOST:-http://localhost:9090}
# shellcheck disable=SC2153 # unset variable comes from environment
access_id="${ACCESS_ID}"
# shellcheck disable=SC2153 # unset variable comes from environment
access_secret="${ACCESS_SECRET}"
# shellcheck disable=SC2153 # unset variable comes from environment
user="${DEVICE_USERNAME}"
# shellcheck disable=SC2153 # unset variable comes from environment
pass="${DEVICE_PASSWORD}"

content_type=application/json
uri="/devices/$1/power"
manufacturer="${DEVICE_MANUFACTURER:-supermicro}"

d=$(TZ=GMT date "+%a, %d %b %Y %T %Z")
sig=$(echo -n "GET,$content_type,,$uri,$d" | openssl dgst -sha1 -binary -hmac "${access_secret}" | base64)

# shellcheck disable=SC2015 # if and/else expression is correct
curl -fs \
	-H "X-IPMI-Username: $user" \
	-H "X-IPMI-Password: $pass" \
	-H "X-DEVICE-MANUFACTURER: $manufacturer" \
	-H "Authorization: APIAuth $access_id:$sig" \
	-H "Content-Type: $content_type" \
	-H "Date: $d" \
	"$host$uri" | jq -r -S .state | grep -iE '^(on|off)$' && exit 0 || :

echo Could not verify power status is working properly
exit 1
