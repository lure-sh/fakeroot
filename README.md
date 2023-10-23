# fakeroot

A pure-Go implementation of fakeroot using Linux user namespaces.

### What is fakeroot?

Fakeroot is a utility that runs commands in an environment where they appear to have root privileges even though they don't. The [original `fakeroot` command](https://salsa.debian.org/clint/fakeroot/) does this by using `LD_PRELOAD` to inject custom wrappers around libc functions that behave as if they're running as the root user. Basically, it intercepts calls to functions like `stat()`, `chmod()`, `chown()`, etc. and replaces them with ones that return values that make it seem like the user is root.

### How is this library different?

Instead of injecting custom libc functions, this library uses Linux's user namespaces. Basically, rather than pretending that the user is root, this library uses the Linux kernel's built-in isolation features to make it seem as if the user is actually root. That means even programs that don't use libc (such as Go programs), or programs with a statically-linked libc, will believe they're running as root. However, this approach will only work on Linux kernels new enough (3.8+) and on distros that don't disable this functionality. Most modern Linux systems support it though, so it should work in most cases.

### Why?

Fakeroot is very useful for building packages, as various utilities depend on file permissions and users. For example, the `tar` command. It creates files inside the tar archive with the same permissions as the original files. That means if the files were owned by a particular user, they will still be owned by that user when the tar archive is extracted. This is problematic for package building because it means you can end up with system files in a package, owned by non-root users. Fakeroot is used to trick utilities like `tar` into making files owned as root.

Many utilities require root privileges for some operations but return errors even if the specific thing you're doing doesn't require them. Fakeroot can also be used to execute these programs without actually giving them root privileges, which provides extra security.
