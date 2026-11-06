#!/bin/bash
set -euo pipefail

export DATA_DIR="${DATA_DIR:-/data}"
export HTTP_ADDR="${HTTP_ADDR:-:${PORT:-8080}}"

for setting in UPLOAD_DIR SOURCE_DIR STATE_DIR; do
    if [[ -n "${!setting:-}" ]]; then
        echo "Railway launcher does not support $setting overrides" >&2
        exit 1
    fi
done

# A fresh platform volume is root-owned. Drop privileges before application work.
if [[ "$(id -u)" == 0 ]]; then
    if [[ "$DATA_DIR" != /data ]]; then
        echo 'Railway DATA_DIR must be /data' >&2
        exit 1
    fi
    mkdir -p /data/uploads /data/sources /data/state
    chown app:app /data /data/uploads /data/sources /data/state
    exec su-exec app:app "$0" "$@"
fi

if [[ "$#" -gt 0 ]]; then
    exec /usr/local/bin/ai-doc-gen "$@"
fi

/usr/local/bin/ai-doc-gen migrate

# Both children use one stable encryption key before any request is admitted.
if [[ -z "${LLM_ENCRYPTION_KEY:-}" && ! -s "$DATA_DIR/state/llm-encryption.key" ]]; then
    umask 077
    key_tmp=$(mktemp "$DATA_DIR/state/key.XXXXXX")
    head -c 32 /dev/urandom | base64 > "$key_tmp"
    mv "$key_tmp" "$DATA_DIR/state/llm-encryption.key"
fi

web_pid=''
worker_pid=''
stop_children() {
    trap '' TERM INT
    [[ -z "$web_pid" ]] || kill -TERM "$web_pid" 2>/dev/null || true
    [[ -z "$worker_pid" ]] || kill -TERM "$worker_pid" 2>/dev/null || true
    wait || true
}
trap 'stop_children; exit 0' TERM INT

/usr/local/bin/ai-doc-gen worker &
worker_pid=$!
/usr/local/bin/ai-doc-gen serve &
web_pid=$!

set +e
wait -n "$worker_pid" "$web_pid"
child_status=$?
set -e
stop_children
# Unexpected clean child termination also requires platform recovery.
if [[ "$child_status" == 0 ]]; then
    child_status=1
fi
exit "$child_status"
