# Contributing to `vrsn`

## Setting up for local development

1. Fork the repo.
2. Run `make build` to ensure you can build the binary before making any changes.
3. Checkout a new branch for your changes.
4. Once made don't forget to run `make fmt` and `make lint` to ensure your
changes are inline with the repo's code standards.
5. Make sure you have added any new unit tests for new functionality.
6. Run the tests with `make test`.
7. Create a pull request 🎉

## Adding support for a new version file type

Supported version file types are stored in the `versionFileMatchers()` function in
`internal/files/version.go`.

This function maps the version file name to the values required to extract the
semantic version from a file.

The `lineMatcher` is a function to determine if the current line of the file is
the one with the version on.

`versionRegex` is a regexp expression to extract the version from the line.
Some key things to note:

- The group for the actual version must be a capture group named `semver`. The
following should always work for that part: `(?P<semver>\d+.\d+.\d)`.
- The version should always be the 3rd group in the expression. This is
important as the updater function uses the groups by index.
If you don't need more groups in the expression you can pad it out with `(.*)`,
see the expression used for `VERSION` files as an example.

### Adding unit tests

Unit tests are essential to keep everything working properly as the code
changes. `internal/files/version_(reader|writer)_test.go` contains tests for
every supported file type.
If you are adding a new one you should add a valid example of that file type
to the `internal/files/testdata/all` directory and an example that does not
include the version to `internal/files/testdata/no-version`.

You should  add a unit test for successfully reading the value from the file
and one for throwing an error when the version cannot be found. Both can be
added to the existing table tests in the same format as the current tests.
