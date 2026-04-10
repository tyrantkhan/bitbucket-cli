# Changelog

## [0.0.12](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.11...v0.0.12) (2026-04-10)


### Features

* **repo:** [#54](https://github.com/tyrantkhan/bitbucket-cli/issues/54) add --details flag to repo list and repo move command ([#57](https://github.com/tyrantkhan/bitbucket-cli/issues/57)) ([c5c9090](https://github.com/tyrantkhan/bitbucket-cli/commit/c5c9090c0bf77796b0ad45162c94a3a0f68b793e))
* **workspace:** [#53](https://github.com/tyrantkhan/bitbucket-cli/issues/53) add workspace members command ([#55](https://github.com/tyrantkhan/bitbucket-cli/issues/55)) ([6b076f0](https://github.com/tyrantkhan/bitbucket-cli/commit/6b076f0be2c4d20eaefca4825da948d93b01b50b))

## [0.0.11](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.10...v0.0.11) (2026-03-21)


### Features

* **repo:** [#50](https://github.com/tyrantkhan/bitbucket-cli/issues/50) add project column and --project/--exclude-project filters to repo list ([#51](https://github.com/tyrantkhan/bitbucket-cli/issues/51)) ([c8cb514](https://github.com/tyrantkhan/bitbucket-cli/commit/c8cb5142578bc35615f7291b1b1f4e46c46e65fc))

## [0.0.10](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.9...v0.0.10) (2026-03-09)


### Features

* [#34](https://github.com/tyrantkhan/bitbucket-cli/issues/34) add repo URL to version output ([#35](https://github.com/tyrantkhan/bitbucket-cli/issues/35)) ([a0f187a](https://github.com/tyrantkhan/bitbucket-cli/commit/a0f187a39f392b7eb17faa990410a4be1c133948))
* [#37](https://github.com/tyrantkhan/bitbucket-cli/issues/37) add comment to update check ([59b2219](https://github.com/tyrantkhan/bitbucket-cli/commit/59b2219117699ac22524ac95b884211b3f0e70dd))
* [#37](https://github.com/tyrantkhan/bitbucket-cli/issues/37) test release-please without bullet prefix ([#38](https://github.com/tyrantkhan/bitbucket-cli/issues/38)) ([59b2219](https://github.com/tyrantkhan/bitbucket-cli/commit/59b2219117699ac22524ac95b884211b3f0e70dd))
* [#40](https://github.com/tyrantkhan/bitbucket-cli/issues/40) add --format json support to all output-producing commands ([#44](https://github.com/tyrantkhan/bitbucket-cli/issues/44)) ([5f8c4eb](https://github.com/tyrantkhan/bitbucket-cli/commit/5f8c4eb96463299143644604b487869c1a61eca8))
* **pr:** [#39](https://github.com/tyrantkhan/bitbucket-cli/issues/39) add --parent flag to pr comment for threaded replies ([38bcbf4](https://github.com/tyrantkhan/bitbucket-cli/commit/38bcbf4824e345e3533217aff9851a890bd9132b))
* **pr:** [#39](https://github.com/tyrantkhan/bitbucket-cli/issues/39) add --parent flag to pr comment for threaded replies ([#42](https://github.com/tyrantkhan/bitbucket-cli/issues/42)) ([38bcbf4](https://github.com/tyrantkhan/bitbucket-cli/commit/38bcbf4824e345e3533217aff9851a890bd9132b))


### Bug Fixes

* [#37](https://github.com/tyrantkhan/bitbucket-cli/issues/37) add inline comment to repo URL ([59b2219](https://github.com/tyrantkhan/bitbucket-cli/commit/59b2219117699ac22524ac95b884211b3f0e70dd))
* Do not cast int64 to int ([38bcbf4](https://github.com/tyrantkhan/bitbucket-cli/commit/38bcbf4824e345e3533217aff9851a890bd9132b))

## [0.0.9](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.8...v0.0.9) (2026-03-08)


### Bug Fixes

* [#29](https://github.com/tyrantkhan/bitbucket-cli/issues/29) add completion command to cmd menu and remove release auto-merge ([#30](https://github.com/tyrantkhan/bitbucket-cli/issues/30)) ([22d72f0](https://github.com/tyrantkhan/bitbucket-cli/commit/22d72f0daf02bc9030b73ba4df046ef7a56af692))


### Refactoring

* **auth:** [#32](https://github.com/tyrantkhan/bitbucket-cli/issues/32) restyle login form with lipgloss info box ([#33](https://github.com/tyrantkhan/bitbucket-cli/issues/33)) ([4d9312f](https://github.com/tyrantkhan/bitbucket-cli/commit/4d9312fe3ee9e6bf4914f8bac03e12bc7628437b))

## [0.0.8](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.7...v0.0.8) (2026-03-08)


### Features

* **search:** [#26](https://github.com/tyrantkhan/bitbucket-cli/issues/26) add bb search code command ([#27](https://github.com/tyrantkhan/bitbucket-cli/issues/27)) ([8b80689](https://github.com/tyrantkhan/bitbucket-cli/commit/8b80689e80410e7a5d0caa95e578f0b4affa8efd))

## [0.0.7](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.6...v0.0.7) (2026-03-08)


### Features

* **config:** [#18](https://github.com/tyrantkhan/bitbucket-cli/issues/18) add config command group and update notifier ([#22](https://github.com/tyrantkhan/bitbucket-cli/issues/22)) ([9b85e0a](https://github.com/tyrantkhan/bitbucket-cli/commit/9b85e0aeae72ba851b091fa91cdb63154b89a508))

## [0.0.6](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.5...v0.0.6) (2026-03-08)


### Bug Fixes

* [#19](https://github.com/tyrantkhan/bitbucket-cli/issues/19) improve unknown command error handling ([#20](https://github.com/tyrantkhan/bitbucket-cli/issues/20)) ([4ec9f9b](https://github.com/tyrantkhan/bitbucket-cli/commit/4ec9f9b53589cedb0ccf0acc04205ecbaec3d3af))

## [0.0.5](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.4...v0.0.5) (2026-03-08)


### Features

* **pr:** [#14](https://github.com/tyrantkhan/bitbucket-cli/issues/14) add ready, draft, and edit commands ([#15](https://github.com/tyrantkhan/bitbucket-cli/issues/15)) ([31926ae](https://github.com/tyrantkhan/bitbucket-cli/commit/31926ae3d665e7d8af8ad4941f3675fb9116a380))

## [0.0.4](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.3...v0.0.4) (2026-03-08)


### Refactoring

* **auth:** improve login TUI and document shell completions ([#8](https://github.com/tyrantkhan/bitbucket-cli/issues/8)) ([9859c80](https://github.com/tyrantkhan/bitbucket-cli/commit/9859c806314a5a88cd00d46717eecd68d811502e))

## [0.0.3](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.2...v0.0.3) (2026-03-07)


### Bug Fixes

* migrate homebrew tap from broken formula to cask ([e8197b0](https://github.com/tyrantkhan/bitbucket-cli/commit/e8197b0b939490078966b029033f7072d30d1408))

## [0.0.2](https://github.com/tyrantkhan/bitbucket-cli/compare/v0.0.1...v0.0.2) (2026-03-07)


### Features

* add OAuth and API token auth, shell completion, and repo validation ([2feefda](https://github.com/tyrantkhan/bitbucket-cli/commit/2feefdad1018b4d158e900b2ad59821582b34f01)), closes [#1](https://github.com/tyrantkhan/bitbucket-cli/issues/1)
* add pr status, help topics, UX improvements, and docs ([a974725](https://github.com/tyrantkhan/bitbucket-cli/commit/a97472580a17bf2c00aa1ca9521ad3ac36ddc32d))
* initial implementation of bb - Bitbucket Cloud CLI ([3fc293a](https://github.com/tyrantkhan/bitbucket-cli/commit/3fc293a3349d1594756e2174451419829c4e75ee))


### Bug Fixes

* bump Go to 1.25.8 to resolve stdlib vulnerabilities ([a8672ec](https://github.com/tyrantkhan/bitbucket-cli/commit/a8672ec1f46de82c36a8768c0974a5c73655c6d8))
* bump Go to 1.26.1 to resolve crypto/x509 vulnerabilities ([da55075](https://github.com/tyrantkhan/bitbucket-cli/commit/da550750f0e5bd463ec621957ff4066352cdfc2d))
* migrate golangci-lint config to v2 exclusions format ([912022d](https://github.com/tyrantkhan/bitbucket-cli/commit/912022dbcd4e85728dae482384740c75615b6406))
* move gofmt to formatters section for golangci-lint v2 ([6431493](https://github.com/tyrantkhan/bitbucket-cli/commit/643149384418296d8b9f4a5fb27875757cba1383))
* remove gosimple linter and fix lefthook hooks ([2bb0a23](https://github.com/tyrantkhan/bitbucket-cli/commit/2bb0a231bb1ab6d47452f272ad619289b27244db))
* resolve all golangci-lint errors (errcheck, gofmt, ineffassign) ([e9c74cd](https://github.com/tyrantkhan/bitbucket-cli/commit/e9c74cd2fef72e4772f4e1d1d0f27b9b762d300c))
