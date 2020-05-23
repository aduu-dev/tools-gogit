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
go get aduu.dev/tools/gogit/gogitcmd/gogit
```

# Examples

You are in a git repository root and you want to temporarily remove local go.mod replace directives during commit.

## Automatically during commit

The replace commands can be installed into a pre-commit hook.

To install the git pre-commit and post-commit hooks:
```
gogit install-hooks .
```

What is basically being inserted is


.git/hooks/pre-commit:

```
#!/bin/bash

gogit replace .
git add go.mod
```

.git/hooks/post-commit
```
#!/bin/bash

gogit replace --undo .
```

The base command `gogit` can be replaced with a flag for install-hooks: `--base-command=my-command`

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
