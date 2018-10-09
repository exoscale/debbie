Debbie down(load)er
===================

**debbie** is a small utility that will download a given large file, store it in
the filesystem, and compute its md5sum checksum, retrying if anything fails.

Usage
-----

```
./debbie -url=https://... -destination=/home/blah/file.img -md5sum=aaaffff...
```

Building
--------

To build the debbie binary, all you need is:

```go build```

All standard debian packaging tools should work to build packages, such as for example:

```debuild -S```

This projects tries to follow [the Debian golang packagingguide](https://go-team.pages.debian.net/packaging.html)
