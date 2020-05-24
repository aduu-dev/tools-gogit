# gogit

[![API reference](https://img.shields.io/badge/godoc-reference-5272B4)](https://pkg.go.dev/aduu.dev/tools/gogit?tab=overview) ![Test branches](https://github.com/aduu-dev/tools-gogit/workflows/Test%20branches/badge.svg?branch=master)
---

Local development often involves local replace directives like

```
replace aduu.dev/utils => ../go/aduu-dev-utils
```

When remote machines (think Github Actions or build machines) encounter local replace directives then tests fail.
This can be avoided by removing the replace directives temporarily, by hand or with git commit hooks.

gogit allows installing pre-commit and post-commit hooks and also offers manual commands 
which can temporarily back up go.mod, remove the replace directives and then later restore the go.mod from the backup.

# Install

```
go get -u aduu.dev/tools/gogit
```

# Examples

You are in a git repository root and you want to temporarily remove local go.mod replace directives during commit.

## Automatically during commit

The replace commands can be installed into a pre-commit hook.

To install the git pre-commit and post-commit hooks:
```
gogit install-hooks .
```

It basically adds to `.git/hooks/pre-commit` 

```
gogit replace --replace-only-if-staged .
```

and to `.git/hooks/post-commit` it adds

```
gogit replace --replace-only-if-staged --undo .
```

It adds comments to those lines to remember which lines it wrote.
So applying `gogit install-hooks .` twice in a row is idempotent (does not add the line twice).

The base command `gogit` can be replaced with a flag for install-hooks: `--base-command=my-command`

## Removing gogit install hooks

The counter-part to `goit install-hooks .`:

```
gogit remove-hooks .
```

Note that it does not delete the git hooks, but rather only removes the line with the comment it inserted itself.

## Manual

This way you can temporarily remove local replace directives:

```bash
gogit replace .
```

A backup is being written into `go.mod.b`.

To reapply a backup:

```bash
gogit replace --undo .
```
