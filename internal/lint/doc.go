// Package lint implements style and quality checks for .env file key-value pairs.
//
// Rules enforced:
//
//   - Keys must be fully uppercase (e.g. DATABASE_URL, not database_url).
//   - Values must not have leading or trailing whitespace.
//   - Values must not be empty.
//   - Keys must not contain spaces.
//
// Usage:
//
//	env, _ := parser.ParseFile(".env")
//	result := lint.Check(env)
//	if result.HasIssues() {
//		for _, issue := range result.Issues {
//			fmt.Println(issue)
//		}
//	}
package lint
