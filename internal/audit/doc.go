// Package audit provides a persistent audit log for envdiff comparison runs.
//
// Each time envdiff is invoked with the --audit flag, a summary Entry is
// appended to a JSON log file (default: .envdiff-audit.json). This allows
// teams to track when environment drift was detected and between which files.
//
// Usage:
//
//	entry := audit.Entry{
//		Timestamp:  time.Now().UTC(),
//		LeftFile:   leftPath,
//		RightFile:  rightPath,
//		Missing:    missingCount,
//		Mismatched: mismatchedCount,
//		HasDiff:    result.HasDiff(),
//	}
//	if err := audit.Append(".envdiff-audit.json", entry); err != nil {
//		log.Fatal(err)
//	}
//
// The log file is created automatically on first use.
package audit
