## gokrazy-selfupdate

A collection of modules for managing on-device self updates for the gokrazy devices.

This repository contains:
- The selfupdate protocol specification, under [/protocol](./protocol)
- The selfupdate versioned api that implements the interfaces of the specification, under [/api](./api)
- A selfupdate client that implements the protocol, to be deployed on the grokrazy device, under [/client](./client)
- A selfupdate server that implements the protocol, to be deployed as a central service for updates, under [/server](./server)
