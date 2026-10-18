#!/bin/sh
# Starts the harness's MCP server. A binary built next to the source
# (.local/bin, make build) is used as is; otherwise one is built from this
# source with Go, once per source, into ~/.terminu-harness/bin.
set -e
root=$(cd "$(dirname "$0")/.." && pwd)
if [ -x "$root/.local/bin/terminu-harness" ]; then
	exec "$root/.local/bin/terminu-harness" "$@"
fi
home=${TERMINU_HARNESS__HOME:-$HOME/.terminu-harness}
sum=$(cd "$root" && find cmd internal go.mod go.sum -type f | LC_ALL=C sort | xargs cat | cksum | cut -d ' ' -f 1)
bin="$home/bin/terminu-harness-$sum"
if [ ! -x "$bin" ]; then
	if ! command -v go >/dev/null 2>&1; then
		echo "terminu-harness: building the MCP server needs Go (https://go.dev/dl/)" >&2
		exit 1
	fi
	mkdir -p "$home/bin"
	(cd "$root" && go build -o "$bin.tmp" ./cmd/terminu-harness) >&2
	mv "$bin.tmp" "$bin"
fi
exec "$bin" "$@"
