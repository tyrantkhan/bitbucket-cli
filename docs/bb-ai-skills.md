# bb — Bitbucket Cloud CLI skills

This file teaches AI assistants how to use the `bb` CLI to interact with Bitbucket Cloud. Add it to your project's context (e.g., CLAUDE.md, .cursorrules, or similar) to enable your AI to manage PRs, repos, and pipelines.

## Overview

`bb` is a CLI for Bitbucket Cloud. It auto-detects workspace and repo from git remotes, so most commands work without extra flags when run inside a Bitbucket-hosted repo.

## Authentication

The user must authenticate before using bb:

```sh
bb auth login       # interactive — choose OAuth or API token
bb auth status      # verify authentication
```

## Common workflows

### Pull request workflow

```sh
# Check status of your PRs
bb pr status

# List open PRs
bb pr list

# View a specific PR with comments
bb pr view 42 --comments

# Create a PR from the current branch
bb pr create --title "Add feature X" --description "Details here" --destination main

# Create a PR and close source branch after merge
bb pr create --title "Fix bug" --close-source-branch

# Review a PR
bb pr diff 42                # view the diff
bb pr activity 42            # view activity feed
bb pr approve 42             # approve it
bb pr comment 42 --body "LGTM"

# Inline comment on a specific file and line
bb pr comment 42 --body "This needs a nil check" --file src/handler.go --line 55

# Merge with a strategy
bb pr merge 42 --strategy squash
bb pr merge 42 --strategy merge_commit --message "Merge feature X"
```

### Repository operations

```sh
# List repos in workspace
bb repo list

# View repo details
bb repo view myrepo

# Create a new private repo
bb repo create --name myrepo --description "My new project" --private

# Clone a repo
bb repo clone myrepo
bb repo clone myrepo --protocol ssh
```

### Pipeline operations

```sh
# List recent pipelines
bb pipeline list

# Run a pipeline on the current branch
bb pipeline run

# Run a custom pipeline with variables
bb pipeline run --custom deploy --variable ENV=staging --variable REGION=us-east-1

# View pipeline details
bb pipeline view {uuid}

# Follow live logs
bb pipeline logs {uuid} --follow

# Stop a running pipeline
bb pipeline stop {uuid}
```

## Key flags

- `-w, --workspace` — override workspace (auto-detected from git remote)
- `-r, --repo` — override repo (auto-detected from git remote)
- `--format json` — output JSON instead of tables (useful for parsing)
- `--limit N` — limit number of results (default: 30)
- `--web` — open the resource in the browser

## JSON output for scripting

When you need to parse bb output programmatically, use `--format json`:

```sh
bb pr list --format json
bb repo list --format json
bb pipeline list --format json
```

## Tips

- bb auto-detects workspace and repo from your git remote — no need to specify `-w` and `-r` when inside a repo
- Source branch defaults to the current git branch when creating PRs
- Pipeline `run` defaults to the current branch
- `bb pr status` shows PRs you authored and PRs where you're a reviewer
- `--follow` on pipeline logs tails the output in real time
- All commands support `--help` for detailed usage
