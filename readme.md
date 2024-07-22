Mr. Tasker it's a pet project
Mr. Tasker worried about my knowledge

Mr. Tasker can make a lot of different task.

struct according [to](https://github.com/golang-standards/project-layout/blob/master/README_ua.md)

### Plan

- [x] init git repo
- [x] add makefile for building, runing, testing
- [x] add simple webserver witch return /Get localhost:8080/status 200ok "Hello from mr. Tasker" (use http.ServerMux)
- [x] add graceful shutdown
- [x] add concurrency example
- [ ] add logger
Task1 `image-reverter`:
CMD app for downloading, saving and converting image.
Contains 2 implementation sync and async

| version | image count | time              |
|---------|-------------|-------------------|
| sync    | 23          | 2069 milliseconds |
| async   | 23          | 269 milliseconds  |

To run this task please run `make run-image-reverter` from the root dir

Additional:

- [ ] add more static analysis checks
- [ ] add more sophisticated graceful shutdown ([link](https://habr.com/ru/articles/771626/))
