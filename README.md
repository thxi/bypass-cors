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
// or
let resp = await fetch('http://localhost:3228/google.com');
```

## Deploy

`docker`

## TODO:

- [ ] deploy to heroku
- [x] add unit tests
- [ ] add a good readme
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
