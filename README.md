# Local development with replace directives

Local development often involves local replace directives like

```
replace aduu.dev/utils => ../go/aduu-dev-utils
```

But tests fail when they are encountered on different build machines like Github Actions.
To Avoid this I wrote this small script which can temporarily remove the local replacements, especially during git commits.

This avoids having to think about one last step before a push: removing local replace directives.

# Install

```
go get aduu.dev/tools/gogit
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
gogit replace . &&git add go.mod
```

and to `.git/hooks/post-commit` it adds

```
gogit replace --undo .
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

A backup is being written into "go.mod.b".

To reapply a backup:

```bash
gogit replace --undo .
```
