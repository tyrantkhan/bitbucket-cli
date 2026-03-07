# Changelog

## 1.0.0 (2026-03-07)


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
