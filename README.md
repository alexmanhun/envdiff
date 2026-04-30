# envdiff

> Compare `.env` files across environments and report missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff && go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-file> <compare-file> [compare-file...]
```

### Example

```bash
envdiff .env.example .env.production .env.staging
```

**Sample output:**

```
[.env.production]
  ✗ MISSING   : STRIPE_SECRET_KEY
  ✗ MISSING   : SENTRY_DSN

[.env.staging]
  ✗ MISSING   : STRIPE_SECRET_KEY
  ~ MISMATCH  : LOG_LEVEL (expected: "debug", got: "info")
```

### Flags

| Flag | Description |
|------|-------------|
| `--strict` | Exit with non-zero status if any differences are found |
| `--keys-only` | Only report missing keys, ignore value mismatches |
| `--json` | Output results as JSON |

---

## Why envdiff?

Keeping `.env` files in sync across environments is error-prone. `envdiff` makes it easy to catch missing variables before they cause production issues — great for CI pipelines and onboarding checks.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE) © 2024 yourusername