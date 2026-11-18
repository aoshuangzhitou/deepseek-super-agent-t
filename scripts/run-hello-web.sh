#!/bin/sh
# Build and run the hello web service in the background.
#
# Usage:
#   scripts/run-hello-web.sh            # listen on 127.0.0.1:8123
#   scripts/run-hello-web.sh 127.0.0.1:9000
#
# The PID and log live in .run/ so the server can be stopped with:
#   kill "$(cat .run/hello-web.pid)"
set -eu

addr="${1:-127.0.0.1:8123}"
root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
run_dir="$root/.run"

mkdir -p "$run_dir"
cd "$root"

go build -o "$run_dir/hello-web" ./cmd/hello-web

if [ -f "$run_dir/hello-web.pid" ] && kill -0 "$(cat "$run_dir/hello-web.pid")" 2>/dev/null; then
	echo "hello-web already running with PID $(cat "$run_dir/hello-web.pid")"
	exit 1
fi

nohup "$run_dir/hello-web" -addr "$addr" >"$run_dir/hello-web.log" 2>&1 &
echo $! >"$run_dir/hello-web.pid"

sleep 1
if ! kill -0 "$(cat "$run_dir/hello-web.pid")" 2>/dev/null; then
	echo "hello-web failed to start; see $run_dir/hello-web.log" >&2
	exit 1
fi

echo "hello-web listening on http://$addr/ (PID $(cat "$run_dir/hello-web.pid"))"
