# Communicate from `/tmp/.java_pid$$`

## Protocol

## Request

    1 byte PROTOCOL_VERSION
    1 byte '\0'
    n byte command
    1 byte '\0'
    n byte arg1
    1 byte '\0'
    n byte arg2
    1 byte '\0'
    n byte arg3
    1 byte '\0'

## Response

Human readable response... sigh.

## How to create unix socket?

Normally, JVM creates `/tmp/.java_pid$$` by default.
But if it disabled, you can request to create the unix socket to existed process.

  1. create `/proc/$PID/cwd/.attach_pid$PID`.
  2. Send SIGQUIT to `$PID`


