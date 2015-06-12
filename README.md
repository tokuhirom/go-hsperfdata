# hsperfdata-go

HotSpot JVM's performance data analyzer written in Golang.

## Why?

jstat/jps is slow. And these commands aren't programmable.

## SYNOPSIS

    $ hsperfdata-go ps
    13223
    21916 org.jetbrains.jps.cmdline.Launcher

    $ hsperfdata-go stat 21916
    sun.rt._sync_Inflations=13
    sun.rt._sync_Deflations=11
    sun.rt._sync_ContendedLockAttempts=65
    ...
    sun.classloader.findClassTime=937707942
    sun.urlClassLoader.readClassBytesTime=164844351
    sun.zip.zipFiles=135
    sun.zip.zipFile.openTime=23649506

