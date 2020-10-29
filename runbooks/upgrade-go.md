# Upgrade version of go

- Update the version in `go.mod`
- Update the version in Dockerfiles
  - The Dockerfiles in `dockerfiles/`
  - Custom Dockerfiles in service directories
- The version in the GitHub workflows in `.github/workflows/`
- Search for any other references to the old version, e.g. `git grep '1\1.14'`.
