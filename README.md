# sideload

Sideload files (which are not source controlled) in and out of a project.

Inevitably we end up with files in my project that we do not want in source control, but still would like to manage
in some fashion. A list of files is kept in a `.sideload-config` file. When you run `sideload restore` in your project
directory, the list of tracked files found in `.sideload.config` is copied from `~/.sideload/storage/{project_name}/...`
into the local project directory tree.

The `.sideload-config` file itself *would* be kept in source control. This file also communicates to other developers
what external files are expected to be present for development.

# Installation

Download the `sideload` binary for your architecture and add the executable to your `PATH`.
See [Releases](https://github.com/bseib/sideload/releases)

Or build and install from its sources:

```
go install github.com/bseib/sideload@latest
```

(Note: doing a `go install` build and install will not inject the version and build date into the binary, meaning
`sideload version` will report empty values.)

# Quick Start Examples

First time use, and setting up `sideload` in your project dir:

```
$ sideload init
SIDELOADHOME environment variable has not been set. Using default: '/home/bseib/.sideload'
Directory '/home/bseib/.sideload' does not exist. Create it now? [Y/n]
Creating directory '/home/bseib/.sideload'.
Config file '.sideload-config' does not exist. Create it now? [Y/n]
Creating config file '.sideload-config'.
```

Now edit your `.sideload-config` file to add files you want to have tracked by `sideload`. You might want to add those
same filenames to your `.gitignore` file if you don't want them in source control.

When those files exist, `sideload status` tell you what would happen if you ran `sideload restore` or `sideload store`

```
$ sideload status
Would restore with 'sideload restore' from '/home/bseib/.sideload/storage/foo_project':
   -->  testdir/zxcv.txt

Would store with 'sideload store' into '/home/bseib/.sideload/storage/foo_project':
  <--   testdir/asdf
  <--   testdir/qwer
```

The `sideload store` command copies files from your local directory structure to `~/.sideload/storage`.

```
$ sideload store
  stored  testdir/asdf
  stored  testdir/qwer
```

The `sideload restore` command copies files from `~/.sideload/storage` to your local directory structure.

```
$ sideload restore
  restored  testdir/zxcv.txt
```

Then the `sideload status` will show all files being equal:

```
$ sideload status
No changes:
   ==   testdir/asdf
   ==   testdir/qwer
   ==   testdir/zxcv.txt
```

# Contributing

This is a quick first cut at some code, and certainly could use some improvements. Some ideas:
  - Warning if tracked files are not covered by `.gitignore` would be useful.
  - Adding/Removing tracked files via command line might be handy, but perhaps overkill.
  - Identifying dead wood in storage

Just open a [github issue](https://github.com/bseib/sideload/issues) or pull request.

# License

Licensed under the [Apache License Version 2.0](LICENSE).

