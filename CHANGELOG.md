# Changelog

`feat-add-test-config`:
- added `config.parseSourceDir` function that parses logs `sourceDir` into a workable path
- added test for `config.parseSourceDir`
- added test for `config.Init`. Changed `config.Init()` signature to accept directories parameters explicitly rather than on the package level
- added test for `config.ReadConfig`. Changed `config.ReadConfig` signature to accept directories pramaters explicitly rather than on the package level
