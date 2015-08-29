go-crxmake [![Build Status](https://travis-ci.org/mcuadros/go-crxmake.svg?branch=master)](https://travis-ci.org/mcuadros/go-crxmake) [![Latest Stable Version](http://img.shields.io/github/release/mcuadros/go-crxmake.svg?style=flat)](https://github.com/mcuadros/go-crxmake/releases)
==============================

Tool for building Chrome/Chromium extensions from a extension folder. Following
the CRX Package Format [specs](https://developer.chrome.com/extensions/crx)

Installation
------------

### Binaries
```
wget https://github.com/mcuadros/go-crxmake/releases/download/v0.2.0/crxmake_v0.2.0_linux_amd64.tar.gz
tar -xvzf crxmake_v0.2.0_linux_amd64.tar.gz
cp-crxmake_v0.2.0_linux_amd64/crxmake /usr/local/bin/
```

browse the [`releases`](https://github.com/mcuadros/go-crxmake/releases) section to see other archs and versions


### From sources
```
go get -u github.com/mcuadros/go-crxmake/...
```

Usage
-----

```sh
Usage:
  crxmake [OPTIONS] [folder] [output]
      --key-file= private key file.

Help Options:
  -h, --help      Show this help message

Arguments:
  folder:         folder where the extension is located.
  output:         output file name.
```

License
-------

MIT, see [LICENSE](LICENSE)
