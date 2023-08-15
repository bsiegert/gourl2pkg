gourl2pkg – add Go packages to pkgsrc easily
============================================

gourl2pkg is a tool similar to [pkgtools/url2pkg](http://pkgsrc.se/pkgtools/url2pkg). It allows you to add packages for software and libraries written in [Go](http://golang.org) to pkgsrc.

**This tool is obsolete. It was useful for `GOPATH` style builds.**

Go 1.11 introduced a different way of building Go code, called modules. You recognize a module from the `go.mod` file at the top level. Today, module builds are the only sensible way of building Go code. This tool is not needed for modules, `url2pkg` does a decent job on them.
