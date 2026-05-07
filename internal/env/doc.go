// Package env provides utilities for serialising and deserialising
// environment variable maps to and from common text formats.
//
// Supported write formats:
//
//	[FormatPlain]     – KEY="VALUE"  (shell-compatible quoted form)
//	[FormatExport]    – export KEY="VALUE"  (sourceable by bash/zsh)
//	[FormatDockerEnv] – KEY=VALUE  (unquoted, for Docker --env-file)
//
// FromSlice converts an os.Environ-style []string into a map[string]string,
// which is the canonical representation used throughout envdiff.
package env
