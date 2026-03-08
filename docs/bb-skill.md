---
name: bb
description: Use the bb CLI to interact with Bitbucket Cloud — manage PRs, repos, and pipelines. Use when the user asks about Bitbucket pull requests, repositories, pipelines, or wants to perform Bitbucket operations.
allowed-tools: Bash(bb *)
---

You have access to `bb`, a CLI for Bitbucket Cloud. Use it via the Bash tool to help the user with Bitbucket operations.

bb auto-detects workspace and repo from git remotes, so most commands work without flags when inside a Bitbucket-hosted repo.

## Quick examples

### Pull requests

```sh
bb pr status                          # PRs for current branch, authored by you, requesting your review
bb pr list                            # list open PRs
bb pr view 42 --comments              # view PR with threaded comments
bb pr create --title "Add X" --destination main
bb pr diff 42                         # view diff
bb pr approve 42                      # approve
bb pr comment 42 --body "LGTM"        # add comment
bb pr comment 42 --body "Fix this" --file src/handler.go --line 55  # inline comment
bb pr merge 42 --strategy squash      # merge
```

### Repositories

```sh
bb repo list                          # list repos in workspace
bb repo view myrepo                   # view details
bb repo create --name myrepo --private
bb repo clone myrepo --protocol ssh
```

### Pipelines

```sh
bb pipeline list                      # list recent pipelines
bb pipeline run                       # run on current branch
bb pipeline run --custom deploy --variable ENV=staging
bb pipeline view {uuid}               # view details
bb pipeline logs {uuid} --follow      # tail live logs
bb pipeline stop {uuid}               # stop running pipeline
```

## Command reference

### Global flags

| Flag | Short | Description |
|---|---|---|
| `--workspace` | `-w` | Bitbucket workspace slug |
| `--repo` | `-r` | Repository slug |
| `--format` | | Output format: `table` (default), `json` |
| `--limit` | | Maximum number of results (default: 30) |
| `--web` | | Open resource in browser |

### Auth

| Command | Description |
|---|---|
| `bb auth login` | Authenticate with Bitbucket Cloud |
| `bb auth logout` | Remove stored credentials |
| `bb auth status` | Show authentication status |

### Repositories

| Command | Description |
|---|---|
| `bb repo list` | List repositories in a workspace |
| `bb repo view [slug]` | View repository details |
| `bb repo create` | Create a new repository |
| `bb repo clone <slug>` | Clone a repository |

### Pull Requests

| Command | Description |
|---|---|
| `bb pr status` | Show PRs for current branch, authored by you, requesting your review |
| `bb pr list` | List pull requests (`--state OPEN\|MERGED\|DECLINED`) |
| `bb pr view <id>` | View PR details (`--comments` for threads) |
| `bb pr create` | Create a pull request |
| `bb pr merge <id>` | Merge a PR (`--strategy merge_commit\|squash\|fast_forward`) |
| `bb pr approve <id>` | Approve a pull request |
| `bb pr decline <id>` | Decline a pull request |
| `bb pr comment <id>` | Add a comment (`--body`, `--file`, `--line` for inline) |
| `bb pr diff <id>` | Show diff (`--stat` for summary) |
| `bb pr activity <id>` | Show activity feed |

### Pipelines

| Command | Description |
|---|---|
| `bb pipeline list` | List recent pipelines |
| `bb pipeline view <uuid>` | View pipeline details |
| `bb pipeline run` | Run a pipeline (`--branch`, `--custom`, `--variable KEY=VALUE`) |
| `bb pipeline stop <uuid>` | Stop a running pipeline |
| `bb pipeline logs <uuid>` | View step logs (`--step`, `--follow`) |

### Shell Completions

```sh
bb completion bash   # output bash completion script
bb completion zsh    # output zsh completion script
bb completion fish   # output fish completion script
bb completion pwsh   # output PowerShell completion script
```

## Tips

- **Always use `--format json`** when reading bb output — the default table format is designed for humans, not machines. JSON output is structured and reliable to parse.
- Source branch defaults to current git branch when creating PRs
- Pipeline `run` defaults to current branch
- User must run `bb auth login` before first use
