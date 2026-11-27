// Command tracelogs aggregates logs for a single trace from every
// logs/<service>_log.jsonl file into one merged, time-sorted JSONL file.
//
// Each app server writes its own log file, so one request's trace is spread
// across several of them (webserver, workerproducer, workerconsumer, ...).
// This command finds every line whose trace_id matches, tags it with the
// service it came from (derived from the log file name), sorts all lines by
// their "time" field, and writes the result to one file.
//
// Usage: tracelogs <trace-id> [out-file]
//
// out-file defaults to logs/trace_logs_search_result.jsonl. Lines that are not
// valid JSON, or that carry no "time" field, are copied through unchanged and
// sorted last so nothing is silently dropped.
//
// Run with `make trace-logs TRACE_ID=<id>`. Exits non-zero when no matching
// log line is found or the trace id is missing.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	ansiReset = "\033[0m"
	ansiRed   = "\033[31m"
	ansiGreen = "\033[32m"
	ansiCyan  = "\033[36m"

	defaultLogsDir = "logs"
	defaultOutName = "trace_logs_search_result.jsonl"
)

// logLine is the subset of a logrus JSON entry this command reasons about.
type logLine struct {
	Service string
	Time    time.Time
	HasTime bool
	Raw     map[string]any
	Order   int
}

func main() {
	if len(os.Args) < 2 || strings.TrimSpace(os.Args[1]) == "" {
		fmt.Fprintf(os.Stderr, "%susage: tracelogs <trace-id> [out-file]%s\n", ansiRed, ansiReset)
		os.Exit(2)
	}
	traceID := os.Args[1]

	outPath := filepath.Join(defaultLogsDir, defaultOutName)
	if len(os.Args) >= 3 && strings.TrimSpace(os.Args[2]) != "" {
		outPath = os.Args[2]
	}

	files, err := filepath.Glob(filepath.Join(defaultLogsDir, "*_log.jsonl"))
	if err != nil {
		fail("listing log files: %v", err)
	}
	if len(files) == 0 {
		fail("no log files found under %s/ (run the servers first)", defaultLogsDir)
	}

	var matched []logLine
	order := 0
	for _, path := range files {
		service := serviceFromPath(path)
		lines, err := searchFile(path, traceID, service)
		if err != nil {
			fail("reading %s: %v", path, err)
		}
		for i := range lines {
			lines[i].Order = order
			order++
		}
		matched = append(matched, lines...)
	}

	if len(matched) == 0 {
		fail("no log line with trace_id %q in %d file(s)", traceID, len(files))
	}

	sort.SliceStable(matched, func(i, j int) bool {
		a, b := matched[i], matched[j]
		if a.HasTime && b.HasTime {
			if !a.Time.Equal(b.Time) {
				return a.Time.Before(b.Time)
			}
			return a.Order < b.Order
		}
		// Lines without a parsable time sort last, keeping their file order.
		return !a.HasTime && b.HasTime
	})

	if err := writeOutput(outPath, matched); err != nil {
		fail("writing %s: %v", outPath, err)
	}

	fmt.Printf("%s✔%s wrote %s%d%s line(s) for trace %s%s%s to %s%s%s\n",
		ansiGreen, ansiReset,
		ansiCyan, len(matched), ansiReset,
		ansiCyan, traceID, ansiReset,
		ansiCyan, outPath, ansiReset)
}

// searchFile returns the lines of path whose trace_id equals traceID, each
// tagged with the service derived from the file name.
func searchFile(path, traceID, service string) ([]logLine, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var out []logLine
	scanner := bufio.NewScanner(f)
	// Log lines can be long (GORM SQL + spans); raise the default 64 KiB cap.
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}

		var obj map[string]any
		if err := json.Unmarshal([]byte(raw), &obj); err != nil {
			continue
		}
		if id, _ := obj["trace_id"].(string); id != traceID {
			continue
		}

		line := logLine{Service: service, Raw: obj}
		if ts, _ := obj["time"].(string); ts != "" {
			if t, err := time.Parse(time.RFC3339, ts); err == nil {
				line.Time = t
				line.HasTime = true
			}
		}
		out = append(out, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// serviceFromPath maps logs/webserver_log.jsonl to "webserver".
func serviceFromPath(path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, ".jsonl")
	base = strings.TrimSuffix(base, "_log")
	return base
}

func writeOutput(path string, lines []logLine) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for _, line := range lines {
		out := make(map[string]any, len(line.Raw)+1)
		for k, v := range line.Raw {
			out[k] = v
		}
		out["service"] = line.Service
		if err := enc.Encode(out); err != nil {
			_ = f.Close()
			return err
		}
	}
	if err := w.Flush(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s✘ %s%s\n", ansiRed, fmt.Sprintf(format, args...), ansiReset)
	os.Exit(1)
}
