# gorsh
Go Reverse Shell

## Introduction 

Gorsh is a library that can be embedded into any Go program, or operate standalone.  It either listens for an incoming SSH connection, make an outgoing TLS connection to a rendezvous server if incoming networking is restricted, or use any other given connection. Gorsh has a minimal built-in shell, and can be configured at compile-time, so no additional files are required for its operation.

## Use Cases

My original use-case was exporing the AWS Lambda environment, where bundling an SSH server is annoying, and getting inbound SSH to them is even more so. But deploying new code to them is easy!

I expect this may be useful in other situations, like debugging minimal containers, embedded devices, or other limited execution environments. 

## Security

This code is not yet production ready, and should be considered a research and hobby project only.  Please contact me before using this code.
