---
name: bb
version: 1.1.0
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
bb pr comment 42 --body "Good point" --parent 764369882            # threaded reply
bb pr merge 42 --strategy squash      # merge
bb pr ready 42                        # mark draft as ready
bb pr draft 42                        # convert to draft
bb pr edit 42 --title "New title"     # edit PR fields
```

### Repositories

```sh
bb repo list                          # list repos in workspace
bb repo list --project PROJ           # filter by project
bb repo list --exclude-project PROJ   # exclude a project
bb repo view myrepo                   # view details
bb repo create --name myrepo --private
bb repo clone myrepo --protocol ssh
```

### Search

```sh
bb search code "handleError"                          # search across workspace
bb search code "func main" --repo myrepo --extension go  # filter by repo and extension
bb search code "TODO" --language python               # filter by language
bb search code "handleError" --format json            # JSON output
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
| `--format` | | Output format: `table` (default), `json` (supported on most commands) |
| `--limit` | | Maximum number of results (default: 30) |
| `--web` | | Open resource in browser |

### Auth

| Command | Description |
|---|---|
| `bb auth login` | Authenticate with Bitbucket Cloud |
| `bb auth login --web --client-id KEY --client-secret SECRET` | Authenticate with a custom OAuth consumer |
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
| `bb pr comment <id>` | Add a comment (`--body`, `--file`, `--line` for inline, `--parent` for threaded replies) |
| `bb pr diff <id>` | Show diff (`--stat` for summary) |
| `bb pr activity <id>` | Show activity feed |
| `bb pr ready <id>` | Mark draft PR as ready for review |
| `bb pr draft <id>` | Convert PR to draft |
| `bb pr edit <id>` | Edit title, description, or reviewers (`--title`, `--description`, `--reviewer`, `--add-reviewer`, `--remove-reviewer`) |

### Search

| Command | Description |
|---|---|
| `bb search code <query>` | Search for code across repos (`--repo`, `--extension`, `--language`, `--path`) |

### Pipelines

| Command | Description |
|---|---|
| `bb pipeline list` | List recent pipelines |
| `bb pipeline view <uuid>` | View pipeline details |
| `bb pipeline run` | Run a pipeline (`--branch`, `--custom`, `--variable KEY=VALUE`) |
| `bb pipeline stop <uuid>` | Stop a running pipeline |
| `bb pipeline logs <uuid>` | View step logs (`--step`, `--follow`) |

## PR review workflow

When the user asks you to review a PR, follow these steps:

### Step 1: Identify the PR

```sh
# If you have a PR number:
bb pr view 42 --format json

# If you have a branch name, list open PRs and find it:
bb pr list --format json
```

**Do NOT** use `bb pr list --source <branch>` — there is no `--source` filter flag.

### Step 2: Read the diff

**Use `bb pr diff` when reviewing PRs.** It shows exactly what Bitbucket sees — the diff between source and destination branches as the API returns it. Use `git diff` for local work: checking uncommitted changes, comparing arbitrary branches/commits, or exploring history. But for PR-specific review, `bb pr diff` is the right tool.

```sh
bb pr diff 42              # full colored diff (best for reading)
bb pr diff 42 --stat       # file-level summary first
bb pr diff 42 --format json  # structured JSON with per-file additions/deletions/patches
```

Use the plain `bb pr diff 42` output to read the diff (it's human-readable). Use `--format json` only if you need to programmatically parse file names or stats.

### Step 3: Read existing comments

```sh
bb pr view 42 --comments --format json
```

Note the comment IDs — you'll need them if you want to reply to a thread.

### Step 4: Read the changed files

Use `Read` or `Grep` tools to read the full files for context beyond the diff. Don't rely only on the diff — understand the surrounding code.

### Step 5: Post comments

There are **three types of comments** — use the right one:

**General comment** — for overall feedback, summaries, or approval/decline rationale:
```sh
bb pr comment 42 --body "Overall LGTM, a few minor suggestions below."
```

**Inline comment** — for feedback on a specific line in a specific file. Use `--file` and `--line` together. The `--line` value must be a line number visible in the diff (a new-side line number):
```sh
bb pr comment 42 --body "This could be null" --file src/handler.go --line 55
```

**Threaded reply** — to reply to an existing comment. Use `--parent` with the comment ID (from `bb pr view --comments --format json` or from a previous `bb pr comment` result):
```sh
bb pr comment 42 --body "Good point, fixed" --parent 764369882
```

**Decision guide:**
- Giving overall feedback or a review summary? → **General comment**
- Pointing out an issue on a specific line of code? → **Inline comment** (`--file` + `--line`)
- Responding to someone else's comment or your own? → **Threaded reply** (`--parent`)
- Posting multiple review points? → One **general comment** for the summary, then **inline comments** for each specific issue

### Step 6: Approve, request changes, or leave as-is

```sh
bb pr approve 42           # approve the PR
# No "request changes" action in Bitbucket API — just leave comments
```

## Tips

- **Use `--format json`** when you need to parse bb output programmatically. For reading diffs or viewing PR details, the default table/text format is fine.
- Source branch defaults to current git branch when creating PRs
- Pipeline `run` defaults to current branch
- User must run `bb auth login` before first use
- **Use `bb pr diff` for PR reviews**, `git diff` for local/branch comparisons
