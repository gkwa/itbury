# itbury

## Summary

File System Event Watcher with Rate Limiting

## Purpose

The purpose of this project is to provide a solution for watching file
system events in a specified directory and executing a function in
response to those events, while also applying rate limiting to control
the frequency of function execution.

## Motivation

Efficiently backup Obsidian notes using git.

## Log
```bash
[mtm@taylors-MacBook-Pro-2:itbury(master)]$ make && ./itbury test1  -vv
level=DEBUG source=cmd/root.go:123 msg=setup verbose=true
level=DEBUG source=test1/test1.go:76 msg="function called within time limit after file system change"
level=DEBUG source=test1/test1.go:55 msg="Function call suppressed due to rate limiting"
level=DEBUG source=test1/test1.go:55 msg="Function call suppressed due to rate limiting"
level=DEBUG source=test1/test1.go:55 msg="Function call suppressed due to rate limiting"
level=DEBUG source=test1/test1.go:76 msg="function called within time limit after file system change"
level=DEBUG source=test1/test1.go:55 msg="Function call suppressed due to rate limiting"
```
