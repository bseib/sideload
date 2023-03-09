## Build and Release

To do a build and release, update the version by tagging the repo, then
run `goreleaser release --rm-dist`:

```
git tag -a v0.1.0 -m "Sideload Release v0.1.0"
git push origin v0.1.0
goreleaser release --rm-dist
```

Requires:

```
export GITHUB_TOKEN="YOUR_GH_TOKEN"
```

See: https://goreleaser.com/quick-start/
