# Tempo Backend Interview Assignment (Go)

## Requirements

To run this project you need to have the following installed:

1. [Go](https://golang.org/doc/install) version 1.18
2. [Docker](https://docs.docker.com/get-docker/) version 20
3. [Docker Compose](https://docs.docker.com/compose/install/) version 1.29
4. [GNU Make](https://www.gnu.org/software/make/)
5. [mock](https://github.com/golang/mock)

    Install the latest version with:
    ```
    go install github.com/golang/mock/mockgen@latest
    ```

## Running

To run the project, run the following command:

```
copy .env.example to .env
```

```
make migrate
```

```
make run-server
```

To see the api docs, you can access on 
```
http://localhost:8080/docs/swagger/index.html#
```

## Testing

To run unit test and integration test, run the following command:

```
make test-all
```
