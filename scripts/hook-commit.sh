#!/bin/sh
# The PreToolUse hook on Bash. Most commands are not commits: only those that
# mention one start the harness (and build it the first time).
input=$(cat)
case "$input" in
*commit*) ;;
*) exit 0 ;;
esac
printf '%s' "$input" | exec "$(dirname "$0")/mcp.sh" hook-commit
