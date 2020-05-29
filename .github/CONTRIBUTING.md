## Contributing

[legal]: https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license
[license]: ../LICENSE.md
[code-of-conduct]: CODE_OF_CONDUCT.md

Thanks for your interest in contributing to the gotodo.

We accept pull requests for bug fixes and features. We'd also love to hear about ideas for new features as issues.

Please do:

* open an issue if things aren't working as expected
* open an issue to propose a significant change
* open a pull request to fix a bug
* open a pull request to fix documentation about a command

Please avoid:

* adding installation instructions specifically for your OS/package manager

## Building the project

Prerequisites:
- Go 1.14

Build with `go build -o bin/gotodo ./cmd/gotodo`

Run the new binary as: `./bin/gotodo`

Run tests with `go test ./...`

## Submitting a pull request

1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and ensure tests pass
1. Submit a pull request

Contributions to this project are [released][legal] to the public under the [project's open source license][license].

Please note that this project adheres to a [Contributor Code of Conduct][code-of-conduct]. By participating in this project you agree to abide by its terms.

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
