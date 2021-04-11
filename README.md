# 🍅 pomo

Kudos to [kevinschoon](https://github.com/kevinschoon). This is an adaptation of his code. The core logic remains the same.

I am using this fork as a demo of my coding skills to apply for the job of Software Development.

I believe that the best way to show my coding skills is to pick a project that I can refactor to show some good practices and quality software.

#
## Main objectives:

- [ ] Develop Unit Tests (Ongoing)
- [x] Update libraries
- [x] Change the cli to the industry standard [cobra](https://github.com/spf13/cobra) and 
- [x] Make logging uniform accross the solution
- [x] Separate client and Server logic
- [ ] Make client and server communicate using GRPC
- [ ] Develop a RESTFul interface for the server
- [ ] Create a web client
- [ ] Define deployment environment 
- [ ] Generate charts/burndown

### Folder Structure
```
pkg
├── core                  <- Main structures and interfaces.
│   ├── models
|   │   ├── bindata.go
|   │   ├── definitions.go
|   │   └── protocol.go
│   ├── client.go
│   └── server.go
├── conf                  <- Base configuration and loggers.
│   ├── defaults.go
│   ├── logger.go
│   └── version.go
├── client                <- Client implementation.
│   ├── unix
│   |   └── unix.go
│   └── rest
│       └── rest.go
├── runner                <- Runner implementation.
│   ├── util.go
│   ├── runner.go
│   └── ui.go
└── server                <- Server implementation.
    ├── unix
    |   └── unix.go
    ├── rest
    |   └── rest.go
    └── store
        └── store.go
```

