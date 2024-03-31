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

## Inputs

Only one input is required for this action:

`score_threshold`: The minimum score required for a dependency to be considered
high quality. Anything below this score will fail the action.