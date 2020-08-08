# bypass-cors

A simple http server to bypass CORS origin request

Features:

- all http methods are supported
- forwards:
  - query parameters
  - headers
  - body
- follows 3xx http redirects \*
- [rs/zerolog](https://github.com/rs/zerolog) for logging
- [rs/cors](https://github.com/rs/cors) for handling cors

## Example

An example application is deployed to [heroku]()

```javascript
let resp = await fetch('http://localhost:3228/https://google.com');
// or, if the application is running locally
let resp = await fetch('http://localhost:3228/google.com');
```

## Run

### Flags

```
-p string
      server port (default "3228")
-pp
      enable pretty print
```

### Locally

```bash
go run .
# or with flags
go run . -pp -p 8080
```

### Docker

To build the image locally:

```bash
make docker-build
```

To run it locally:

```bash
make docker-run
# or, if you want to set the flags yourself
docker run -p 1337:1337 bypass-cors -p 1337 -pp
```

## TODO:

- [ ] deploy to heroku
- [x] add unit tests
- [x] add a good readme
- [x] profiling
- [x] remove todos
- [x] dockerfile
- [x] close req bodies and handle errors

\* These 3xx codes are followd:

- 301 (Moved Permanently)
- 302 (Found)
- 303 (See Other)
- 307 (Temporary Redirect)
- 308 (Permanent Redirect)
