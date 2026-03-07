# bb — Bitbucket Cloud CLI

## Project overview

A CLI for Bitbucket Cloud written in Go. The binary is `bb`. Module path is `github.com/tyrantkhan/bb`.

## Inspiration

bb is heavily inspired by the [GitHub CLI (gh)](https://github.com/cli/cli). When exploring how to implement a feature, reference gh's codebase and patterns — it's the gold standard for a git-host CLI. The `gh` command is available in the shell for reference.

## Structure

- `main.go` — entry point
- `cmd/` — command definitions (auth, repo, pr, pipeline)
- `internal/api/` — Bitbucket API client, pagination, error handling
- `internal/auth/` — OAuth flow, token refresh, credential storage
- `internal/config/` — YAML config, XDG paths
- `internal/cmdutil/` — factory pattern, shared flags, workspace/repo resolution
- `internal/models/` — data structs (repo, PR, pipeline, user, etc.)
- `internal/output/` — table/JSON formatting, colors, markdown rendering
- `internal/git/` — git remote parsing, clone execution
- `docs/` — command reference and AI skills file

## Key patterns

- **Factory pattern**: all commands get dependencies via `cmdutil.Factory` injected through context
- **Smart defaults**: workspace and repo are resolved from git remotes, falling back to config
- **OAuth auto-refresh**: expired tokens are refreshed transparently in `Factory.APIClient()`
- **Generic pagination**: `api.Paginate[T]` handles all paginated endpoints
- **Interactive fallback**: missing inputs trigger TUI prompts via `huh`

## Dependencies

- `urfave/cli/v3` — CLI framework
- `charm.land/huh/v2` — interactive forms
- `charm.land/lipgloss/v2` — terminal styling
- `charm.land/glamour/v2` — markdown rendering
- `gopkg.in/yaml.v3` — config parsing

## Build

```sh
make build    # builds ./bb with version info from git tags
make test     # go test ./...
make lint     # golangci-lint
```

## Conventions

- Don't add Co-Authored-By to commits
- Keep commands in their own files under `cmd/<group>/`
- All API paths use `/2.0/` prefix (Bitbucket Cloud v2 API)
- Slug validation: `^[a-zA-Z0-9._-]+$`
- Credentials stored at `~/.config/bb/credentials.json` with 0600 permissions
- Config stored at `~/.config/bb/config.yml`

## When modifying commands

- Update `docs/commands.md` if flags or behavior change
- Update `docs/bb-ai-skills.md` if new workflows are added
- The `bb reference` command auto-generates from command definitions — no manual update needed
