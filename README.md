# Trusty Dependency Analysis Action

This action takes any added dependencies within a pull request and assesses their 
quality using the [Trusty](https://trustypkg.dev/) API. If any dependencies are
found to be below a certain threshold (See details below), the action will fail.

If any dependencies are malicious, or deprecated, the action will also fail.

Full Language Support (inline with Trusty):

* Python
* JavaScript
* Java
* Rust
* Go

## Usage

To use this action, you can add the following to your workflow:

```yaml
name: TrustyPkg Dependency Check

on:
  pull_request:
    branches:
      - main

jobs:
  trusty_pkg_check:
    runs-on: ubuntu-latest
    name: Check Dependencies with TrustyPkg
    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: TrustyPkg Action
        uses: stacklok/trusty-action@v0.0.1
        with:
          score_threshold: 5
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

Only one input is available for this action:

`score_threshold`: The minimum score required for a dependency to be considered
high quality. Anything below this score will fail the action.