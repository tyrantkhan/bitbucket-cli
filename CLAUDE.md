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

## Branch naming

Format: `<type>/<issue>-<description>`

- **Type** must be one of: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`, `hotfix`, `release`
- **Issue number** is required — always create an issue first if one doesn't exist
- **Description** must be lowercase alphanumeric with dots, hyphens, or underscores
- Standalone branches `main`, `master`, `develop` are also allowed

Examples:
- `feat/1-add-oauth`
- `fix/12-login-bug`
- `chore/5-update-deps`

## Commit & PR title format

Conventional Commits required for both commits and PR titles (PR titles become squash merge commits).

Format: `<type>[(<scope>)]: #<issue> <description>`

- **Scope** is optional but preferred — use the module name (`auth`, `pr`, `repo`, `pipeline`, `api`, `config`)
- **Issue number** is required — always create an issue first if one doesn't exist

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`

PR title must use the highest-priority conventional commit type from its commits, so release-please picks up version-bumping changes:

- Priority: `feat` > `fix` > `refactor` > `perf` > all others
- Example: if a PR has both `docs` and `fix` commits, the title must use `fix`
- CI enforces this via the `pr-title-priority` job in the `pr-title` workflow

Examples:
- `feat(auth): #7 add spinner to login flow`
- `fix(pr): #12 handle empty reviewer list`
- `docs: #15 add shell completions to README`
- `ci: #16 add PR title lint workflow`

## Conventions

- Don't add Co-Authored-By to commits
- Keep commands in their own files under `cmd/<group>/`
- All API paths use `/2.0/` prefix (Bitbucket Cloud v2 API)
- Slug validation: `^[a-zA-Z0-9._-]+$`
- Credentials stored at `~/.config/bb/credentials.json` with 0600 permissions
- Config stored at `~/.config/bb/config.yml`

## When modifying commands

- Update `docs/commands.md` if flags or behavior change
- Update `docs/bb-skill.md` if new workflows are added
- The `bb reference` command auto-generates from command definitions — no manual update needed
