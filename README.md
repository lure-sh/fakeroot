# fakeroot

[![Go Reference](https://pkg.go.dev/badge/lure.sh/fakeroot.svg)](https://pkg.go.dev/lure.sh/fakeroot)

A pure-Go implementation of fakeroot using Linux user namespaces.

### What is fakeroot?

Fakeroot is a utility that runs commands in an environment where they appear to have root privileges even though they don't. The [original `fakeroot` command](https://salsa.debian.org/clint/fakeroot/) does this by intercepting calls to libc functions like `stat()`, `chmod()`, `chown()`, etc. and replacing them with ones that return values that make it seem like the user is root.

### How is this library different?

Instead of injecting custom libc functions, this library uses the Linux kernel's built-in isolation features to make a sort of container where the user is root. That means even programs that don't use libc (such as Go programs), or programs with a statically-linked libc, will believe they're running as root.

You can also nest this type of fakeroot up to 32 times, unlike the original libc-based one, which doesn't support nesting at all.

However, this approach will only work on Linux kernels new enough (3.8+) and on distros that don't disable this functionality. Most modern Linux systems support it though, so it should work in most cases.

### Why?

Many utilities depend on file permissions and user ownership. For instance, the tar command creates files within a tar archive with the same permissions as the original files. This means that if the files were owned by a specific user, they will retain that ownership when the tar archive is extracted. This can become problematic when building packages because it could lead to system files in a package being owned by non-root users. By making it seem as if the current user is root and therefore all the files are owned by root, fakeroot tricks utilities like `tar` into making its files owned by root.

Also, many utilities may require root privileges for certain operations but might return errors even when the specific task doesn't necessarily need those elevated permissions. Fakeroot can be used to execute these programs without actually granting them root privileges, which provides some extra security.

### nsfakeroot

This repo includdes a command-line utility called `nsfakeroot`. To install it, run the following command:

```bash
go install lure.sh/fakeroot/cmd/nsfakeroot@latest
```

Running `nsfakeroot` on its own will start your login shell in the fakeroot environment. If you provide arguments, those will be used as the command.

Examples:

```bash
nsfakeroot        # -> (login shell)
nsfakeroot whoami # -> root
nsfakeroot id -u  # -> 0
nsfakeroot id -g  # -> 0
```