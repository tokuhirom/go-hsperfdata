# hsperfdata-go

HotSpot JVM's performance data analyzer written in Golang.

## Why?

jstat/jps is slow. And these commands aren't programmable.

## INSTALL

    go get github.com/tokuhirom/go-hsperfdata/hsps/
    go get github.com/tokuhirom/go-hsperfdata/hsstat/

## SYNOPSIS

    $ hsps
    13223
    21916 org.jetbrains.jps.cmdline.Launcher

    $ hsstat 21916
    sun.rt._sync_Inflations=13
    sun.rt._sync_Deflations=11
    sun.rt._sync_ContendedLockAttempts=65
    ...
    sun.classloader.findClassTime=937707942
    sun.urlClassLoader.readClassBytesTime=164844351
    sun.zip.zipFiles=135
    sun.zip.zipFile.openTime=23649506

## Benchmarking

```
$ time hsps
13223
21916 org.jetbrains.jps.cmdline.Launcher
hsps  0.00s user 0.00s system 83% cpu 0.010 total

$ time jps
13223
21916 Launcher
93597 Jps
jps  0.41s user 0.09s system 104% cpu 0.479 total
```

## LICENSE

    The MIT License (MIT)
    Copyright © 2015 Tokuhiro Matsuno, http://64p.org/ <tokuhirom@gmail.com>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the “Software”), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.
