# Command Reference

## Global flags

These flags are available on most commands:

| Flag | Short | Description |
|---|---|---|
| `--workspace` | `-w` | Bitbucket workspace slug |
| `--repo` | `-r` | Repository slug |
| `--format` | | Output format: `table` (default), `json` |
| `--limit` | | Maximum number of results (default: 30) |
| `--web` | | Open resource in browser |

---

## Auth

### `bb auth login`

Authenticate with Bitbucket Cloud. Supports two methods:

**OAuth (browser flow):**

```sh
bb auth login --web
```

Opens your browser to authorize with Bitbucket. A local callback server receives the authorization code and exchanges it for tokens. Tokens auto-refresh when they expire.

**API token:**

```sh
bb auth login --api-token
```

Prompts for your Atlassian email and an API token generated at https://id.atlassian.com/manage-profile/security/api-tokens.

Required token scopes: Repositories (read/write), Pull Requests (read/write), Pipelines (read/write).

**Flags:**

| Flag | Description |
|---|---|
| `--web` | Use OAuth browser flow |
| `--api-token` | Use API token authentication |
| `--username` | Atlassian account email (skips prompt) |
| `--token` | API token value (skips prompt) |
| `--client-id` | Override default OAuth consumer key |
| `--client-secret` | Override default OAuth consumer secret |

### `bb auth logout`

Remove stored credentials. Prompts for confirmation.

### `bb auth status`

Display current authentication status including user, auth method, token expiry, and default workspace.

---

## Repositories

### `bb repo list`

List repositories in a workspace.

```sh
bb repo list
bb repo list -w myworkspace --limit 50
bb repo list --format json
```

### `bb repo view [slug]`

View repository details including description, visibility, language, default branch, and clone URLs.

```sh
bb repo view myrepo
bb repo view --web          # open in browser
bb repo view --format json
```

### `bb repo create`

Create a new repository. Interactive form for missing values.

```sh
bb repo create
bb repo create --name myrepo --private --project PROJ
```

| Flag | Description |
|---|---|
| `--name` | Repository name |
| `--description` | Repository description |
| `--private` | Make private (default: true) |
| `--project` | Project key |

### `bb repo clone <slug> [directory]`

Clone a repository.

```sh
bb repo clone myrepo
bb repo clone myrepo ./local-dir
bb repo clone myrepo --protocol ssh
```

| Flag | Description |
|---|---|
| `--protocol` | `https` (default) or `ssh` |

---

## Pull Requests

### `bb pr list`

List pull requests filtered by state.

```sh
bb pr list
bb pr list --state MERGED
bb pr list --state DECLINED --limit 10
```

| Flag | Description |
|---|---|
| `--state` | `OPEN` (default), `MERGED`, `DECLINED` |

### `bb pr view <id>`

View pull request details including title, state, author, branches, description, reviewers, and review status.

```sh
bb pr view 42
bb pr view 42 --comments    # include comment threads
bb pr view 42 --web         # open in browser
```

| Flag | Description |
|---|---|
| `--comments` | Show PR comments |

### `bb pr create`

Create a pull request. Auto-detects source branch from git.

```sh
bb pr create
bb pr create --title "Add feature" --destination develop
bb pr create --source feature/x --close-source-branch
```

| Flag | Description |
|---|---|
| `--title` | PR title |
| `--description` | PR description |
| `--source` | Source branch (default: current branch) |
| `--destination` | Destination branch (default: `main`) |
| `--close-source-branch` | Close source branch after merge |
| `--reviewer` | Reviewer UUID (repeatable) |

### `bb pr merge <id>`

Merge a pull request. Prompts for confirmation.

```sh
bb pr merge 42
bb pr merge 42 --strategy squash
bb pr merge 42 --strategy fast_forward --message "Ship it"
```

| Flag | Description |
|---|---|
| `--strategy` | `merge_commit` (default), `squash`, `fast_forward` |
| `--message` | Merge commit message |

### `bb pr approve <id>`

Approve a pull request.

```sh
bb pr approve 42
```

### `bb pr decline <id>`

Decline a pull request. Prompts for confirmation.

```sh
bb pr decline 42
```

### `bb pr comment <id>`

Add a comment to a pull request. Supports inline comments on specific files and lines.

```sh
bb pr comment 42 --body "Looks good!"
bb pr comment 42 --body "Fix this" --file src/main.go --line 15
```

| Flag | Description |
|---|---|
| `--body` | Comment text (required) |
| `--file` | File path for inline comment |
| `--line` | Line number for inline comment |

### `bb pr diff <id>`

Show the diff of a pull request with color-coded output.

```sh
bb pr diff 42
bb pr diff 42 --stat    # file-level summary only
```

| Flag | Description |
|---|---|
| `--stat` | Show file-level summary |

### `bb pr activity <id>`

Show pull request activity feed including approvals, state changes, and comments.

```sh
bb pr activity 42
bb pr activity 42 --limit 5
```

### `bb pr status`

Show status of pull requests relevant to you (created by you, reviewing, etc.).

```sh
bb pr status
```

---

## Pipelines

### `bb pipeline list`

List pipelines sorted by most recent.

```sh
bb pipeline list
bb pipeline list --limit 5
bb pipeline list --format json
```

### `bb pipeline view <uuid>`

View pipeline details including state, branch, commit, creator, duration, and steps.

```sh
bb pipeline view {uuid}
bb pipeline view {uuid} --web
```

### `bb pipeline run`

Trigger a pipeline. Auto-detects branch from git.

```sh
bb pipeline run
bb pipeline run --branch main
bb pipeline run --custom deploy --variable ENV=staging
bb pipeline run --custom deploy --variable ENV=prod --variable REGION=us-east-1
```

| Flag | Description |
|---|---|
| `--branch` | Branch to run on (default: current branch) |
| `--custom` | Custom pipeline name |
| `--variable` | `KEY=VALUE` variable (repeatable) |

### `bb pipeline stop <uuid>`

Stop a running pipeline. Prompts for confirmation.

```sh
bb pipeline stop {uuid}
```

### `bb pipeline logs <uuid>`

View logs for a pipeline step. Interactive step picker if `--step` is not provided.

```sh
bb pipeline logs {uuid}
bb pipeline logs {uuid} --step {step-uuid}
bb pipeline logs {uuid} --follow    # live tail
```

| Flag | Description |
|---|---|
| `--step` | Step UUID (interactive picker if omitted) |
| `--follow` | Follow live log output |
