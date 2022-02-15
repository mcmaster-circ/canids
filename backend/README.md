# canids-v2-backend

## Running
```sh
go run main.go
```

## New Dependencies
When new Golang dependencies are added to the project, please run the following.

```
go mod vendor
```

In the event that an upstream repository is no longer available, this ensures
that required dependencies are always available locally.
