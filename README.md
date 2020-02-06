# go-birthday-calendar

## Usage

```bash
This is a rest API service to remind people that their birthday is coming up or wish them 'Happy birthday' on their birthday.

Environment Variables:

It is possible to configure the service by using environment variables. To do so, export a variable that matches the desired flag option that meets the following rules:

- convert to upper case
- replace dashes (-) with underscores (_)
- prefix with BIRTHDAY_

Examples

- port -> BIRTHDAY_PORT
- tls-port -> BIRTHDAY_TLS_PORT

Configuration Reload:

If a configuration file is used, the service will monitor for changes to the file and will update the running configuration.

Usage:
  go-birthday-calendar [flags]

Flags:
      --config string      config file (default is $HOME/.go-birthday-calendar.yaml)
      --db-conn string     custom database connection string
      --db-engine string   database engine to use, supported engines: mysql, mssql, postgres, sqlite3 (default "sqlite3")
      --db-host string     database name (default "localhost")
      --db-name string     database port, default depends on db-engine (default ":memory:")
      --db-pass string     database password, if required
      --db-port int        database name
      --db-user string     database username, if required
  -h, --help               help for go-birthday-calendar
      --log-file string    log file to us, defaults to stdin (default "-")
      --log-level string   log level (default "INFO")
      --path string        path to serve from, e.g. /hello
      --port int           port to serve HTTP requests (default 8080)
      --redirect           redirect HTTP to HTTPS
      --tls-cert string    path to server certificate file, include intermediary certificates (default "server.crt")
      --tls-enable         enable TLS
      --tls-port int       port to serve TLS/HTTPS requests (default 8443)
```

## To-Do

- [x] Create configurion file with env and cli flag overrides
- [x] Implement database using backend determined by configuration
- [x] Create router to use with server
- [x] Configure and Enable HTTP server
- [ ] Configure and Enable TLS server
- [ ] Configure and Enable HTTP -> HTTPS redirect
- [ ] Add test code
- [ ] Create Dockerfile
- [ ] Create minimal Chart

## Requirements

1. Build an application that serves the following HTTP-based APIs:

    **Description:** Saves/updates the given user's name and date of birth in the database  
    **Request:** PUT /hello/Morty { "dateOfBirth": "2000-01-01" }  
    **Response:** 204 No Content

    **Description:** Return a hello/birthday message for the given user  
    **Request:** GET /hello/Morty  
    **Response:** 200 OK

    a. when Morty’s birthday is in 5 days:

    ```json
    { "message": "Hello, Morty! Your birthday is in 5 days" }
    ```

    b. when Morty's birthday is today:

    ```json
    { "message": "Hello, Morty! Happy birthday" }
    ```

    Note #1: Use storage/database of your choice  
    Note #2: At BEAT, Go is the programming language of choice. We'll consider it a plus.

1. Create the build, package and deploy logic for your application. It should be as simple as possible to run the application in any environment (locally or in the Cloud). Write instructions so we can run and test it.

1. Produce a system diagram of your solution deployed in a production-grade
environment.

    The application must:  
    1. Be designed to run on Kubernetes  
    1. Scale  
    1. Be highly available
    1. Be well monitored
    1. Be as simple as it can.

1. Please explain the design choices you made.

    Implicit requirements:
    1. The code produced by you is expected to be of high quality
    1. Use common sense
