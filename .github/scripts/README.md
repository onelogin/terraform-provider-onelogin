# Version Bump Automation

This directory contains scripts for automated version bumping as part of the GitHub Actions workflow.

## Files

### `bump_version.py`
Python script that handles version bumping in the project files:
- Reads current version from `GNUmakefile`
- Supports semantic versioning (major, minor, patch)
- Updates version in both `GNUmakefile` and `README.md`
- Outputs the new version for GitHub Actions using the modern `GITHUB_OUTPUT` approach

### Usage
```bash
python3 .github/scripts/bump_version.py [major|minor|patch]
```

## Key Features

1. **Modern GitHub Actions Output**: Uses `$GITHUB_OUTPUT` environment file instead of the deprecated `::set-output` command
2. **Semantic Versioning**: Properly handles major, minor, and patch version bumps
3. **Multi-file Updates**: Updates version references in both GNUmakefile and README.md
4. **Error Handling**: Validates input and provides clear error messages

## GitHub Actions Integration

The script is designed to work with the `version-bump.yml` workflow, which:
- Accepts manual trigger with bump type selection
- Runs the version bump script
- Creates a pull request with the changes using `peter-evans/create-pull-request`

The workflow avoids manual git operations and lets the action handle branch creation, commits, and PR creation to prevent conflicts.