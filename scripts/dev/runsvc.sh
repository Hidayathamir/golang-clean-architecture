#!/usr/bin/env bash
# Run a service and tee its output to a JSONL log file.
#
# The log is truncated on every start so a dev restart begins from an empty file
# instead of growing without bound. Truncation lives here, not in the Makefile:
# `tuistory launch` on an existing session does not run `rm` before the restart,
# so deleting the file there would orphan the still-open descriptor.
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
>"$logfile"
exec go run "$pkg" >>"$logfile" 2>&1