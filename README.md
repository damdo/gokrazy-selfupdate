## gokrazy-selfupdate

A collection of modules for managing on-device self updates for the gokrazy devices.

This repository contains:
- The **selfupdate protocol** specification, under [/protocol](./protocol)
- The **selfupdate versioned api** that implements the interfaces of the specification, under [/api](./api)
- A **selfupdate client** that implements the protocol, to be deployed on the grokrazy device, under [/client](./client)

selfupdate server:
- For a reference **selfupdate server** that implements the protocol and other useful tooling around building for and managing gokrazy devices see: [gokrazy-operator](https://github.com/damdo/gokrazy-operator)
