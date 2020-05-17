# GitHub Actions

These are the workflows that run on each commit to test the correctness and integrity of the new code.

To test locally, use [act](https://github.com/nektos/act).

The standard image that `act` uses does not contain the tools needed to run the workflows, but the more complete node images seem to work ok. For example, to run the `lint` job from `go.yml`:

```shell
act -j lint -P ubuntu-latest=node:12.16-buster
```
