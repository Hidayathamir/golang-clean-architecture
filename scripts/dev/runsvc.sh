#!/usr/bin/env bash
# Run a service and tee its output to a JSONL log file.
#
# Why this exists: `tuistory launch -- sh -c 'go run X >> log'` does not work,
# because tuistory re-parses the command string like a shell, mangling the quoted
# `sh -c` argument. Executing this script as a single path avoids that, and the
# redirect lives here instead of in the Makefile.
#
# usage: runsvc.sh <package> <logfile>
set -euo pipefail

pkg="$1"
logfile="$2"

mkdir -p "$(dirname "$logfile")"
exec go run "$pkg" >>"$logfile" 2>&1