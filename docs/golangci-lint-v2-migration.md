# golangci-lint v2 Migration Guide

## Overview

This document explains the migration from golangci-lint v1 to v2 configuration schema for the LightShare backend.

## Migration Date

**Date:** 2025-11-23
**Version:** golangci-lint v2.6
**Configuration File:** `backend/.golangci.yml`

## Why We Migrated

The GitHub Actions workflow was upgraded to use `golangci/golangci-lint-action@v8.0.0` with golangci-lint `v2.6`, which requires the v2 configuration schema. The v1 configuration was causing validation errors in CI.

## Schema Changes

### 1. Version Field Type

**v1 Syntax:**
```yaml
version: 2
```

**v2 Syntax:**
```yaml
version: "2"
```

**Change:** The version field must be a string, not a number.

---

### 2. Linter Settings Restructure

**v1 Syntax:**
```yaml
linters:
  enable:
    - dupl
    - gocyclo

linters-settings:
  dupl:
    threshold: 100
  gocyclo:
    min-complexity: 15
```

**v2 Syntax:**
```yaml
linters:
  enable:
    - dupl
    - gocyclo

  settings:
    dupl:
      threshold: 100
    gocyclo:
      min-complexity: 15
```

**Change:** The top-level `linters-settings` section has been moved under `linters.settings`.

---

### 3. Issue Exclusions Restructure

**v1 Syntax:**
```yaml
issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - dupl
        - gosec
    - linters:
        - gocritic
      text: "unnecessaryDefer"

  max-same-issues: 0
  new: false
```

**v2 Syntax:**
```yaml
linters:
  exclusions:
    rules:
      - path: _test\.go
        linters:
          - dupl
          - gosec
      - linters:
          - gocritic
        text: "unnecessaryDefer"

issues:
  max-same-issues: 0
  new: false
```

**Change:** The `issues.exclude-rules` section has been moved to `linters.exclusions.rules`. Other issue properties remain under `issues`.

---

## Complete Configuration Comparison

### Before (v1 Schema)

```yaml
version: 2

run:
  timeout: 5m
  tests: true

linters:
  enable:
    - bodyclose
    - dupl
    # ... other linters

linters-settings:
  dupl:
    threshold: 100
  # ... other settings

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - dupl
  max-same-issues: 0
```

### After (v2 Schema)

```yaml
version: "2"

run:
  timeout: 5m
  tests: true

linters:
  enable:
    - bodyclose
    - dupl
    # ... other linters

  settings:
    dupl:
      threshold: 100
    # ... other settings

  exclusions:
    rules:
      - path: _test\.go
        linters:
          - dupl

issues:
  max-same-issues: 0
```

---

## Validation Errors Fixed

The migration resolved these validation errors:

1. **Version type error:**
   ```
   "version" does not validate with "/properties/version/type":
   got number, want string
   ```

2. **Additional properties error:**
   ```
   "" does not validate with "/additionalProperties":
   additional properties 'linters-settings' not allowed
   ```

3. **Exclude rules error:**
   ```
   "issues" does not validate with "/properties/issues/additionalProperties":
   additional properties 'exclude-rules' not allowed
   ```

---

## Verifying the Configuration

To verify the golangci-lint configuration locally:

```bash
cd backend
golangci-lint config verify
```

To run linting:

```bash
cd backend
golangci-lint run
```

---

## CI/CD Integration

The configuration is used in the GitHub Actions workflow:

**File:** `.github/workflows/backend.yml`

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v8.0.0
  with:
    version: v2.6
    working-directory: ./backend
    args: --timeout=5m
```

---

## Additional Resources

- [golangci-lint v2 Migration Guide](https://golangci-lint.run/docs/product/migration-guide/)
- [golangci-lint Configuration Documentation](https://golangci-lint.run/docs/configuration/)
- [Welcome to golangci-lint v2 Blog Post](https://ldez.github.io/blog/2025/03/23/golangci-lint-v2/)

---

## Enabled Linters

Our configuration enables the following linters:

- `bodyclose` - Checks whether HTTP response body is closed successfully
- `copyloopvar` - Detects places where loop variables are copied
- `dogsled` - Checks assignments with too many blank identifiers
- `dupl` - Tool for code clone detection
- `errcheck` - Checks for unchecked errors
- `gochecknoinits` - Checks that no init functions are present
- `goconst` - Finds repeated strings that could be replaced by a constant
- `gocritic` - Provides diagnostics that check for bugs, performance and style issues
- `gocyclo` - Computes and checks cyclomatic complexity
- `goprintffuncname` - Checks that printf-like functions are named with `f` at the end
- `gosec` - Inspects source code for security problems
- `govet` - Reports suspicious constructs
- `ineffassign` - Detects ineffectual assignments
- `misspell` - Finds commonly misspelled English words in comments
- `nakedret` - Finds naked returns in functions greater than a specified length
- `noctx` - Finds sending HTTP request without context.Context
- `nolintlint` - Reports ill-formed or insufficient nolint directives
- `revive` - Fast, configurable, extensible, flexible linter
- `staticcheck` - Advanced Go linter
- `unconvert` - Removes unnecessary type conversions
- `unparam` - Reports unused function parameters
- `unused` - Checks for unused constants, variables, functions and types
- `whitespace` - Detects leading and trailing whitespace

---

## Settings Configuration

### dupl
- `threshold: 100` - Tokens count to trigger issue

### gocyclo
- `min-complexity: 15` - Minimal code complexity to report

### goconst
- `min-len: 2` - Minimum string length to check
- `min-occurrences: 2` - Minimum occurrences to trigger

### gocritic
- Enabled tags: diagnostic, experimental, opinionated, performance, style

### govet
- `enable-all: true` - Enable all available checks

### misspell
- `locale: US` - Use US English

### revive
- 24 rules enabled for comprehensive code review

---

## Exclusions

Test files (`*_test.go`) are excluded from:
- `dupl` - Duplicate code detection
- `gosec` - Security checks

The `gocritic` linter ignores "unnecessaryDefer" warnings.

---

## Future Considerations

- Monitor golangci-lint releases for new linters and features
- Consider enabling additional linters as the codebase matures
- Review and adjust complexity thresholds based on team feedback
- Keep documentation updated with configuration changes
