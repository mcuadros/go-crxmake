go-crxmake [![Latest Stable Version](http://img.shields.io/github/release/mcuadros/go-crxmake.svg?style=flat)](https://github.com/mcuadros/go-crxmake/releases)
==============================

Tool for building Chrome/Chromium extensions from a extension folder. Following
the CRX Package Format [specs](https://developer.chrome.com/extensions/crx)

Installation
------------

```
wget https://github.com/mcuadros/go-crxmake/releases/download/v0.1.0/go-crxmake_v0.1.0_linux_amd64.tar.gz
tar -xvzf go-crxmake_v0.1.0_linux_amd64.tar.gz
cp go-crxmake_v0_v0.1.0_linux_amd64/crxmake /usr/local/bin/
```

browse the [`releases`](https://github.com/mcuadros/go-crxmake/releases) section to see other archs and versions


Usage
-----

```sh
Usage:
  crxmake [OPTIONS] [folder] [output]

Help Options:
  -h, --help    Show this help message

Arguments:
  folder:       folder where the extension is located.
  output:       output file name
```

License
-------

MIT, see [LICENSE](LICENSE)
