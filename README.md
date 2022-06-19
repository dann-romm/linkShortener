# linkShortener

[//]: # ([![forthebadge]&#40;https://forthebadge.com/images/badges/made-with-go.svg&#41;]&#40;https://forthebadge.com&#41; [![forthebadge]&#40;http://forthebadge.com/images/badges/built-with-love.svg&#41;]&#40;http://forthebadge.com&#41;)

a microservice for shortening links

# Usage

To run service with postgresql database:
`make` or `make postgres`

To run service with saving data in RAM:
`make inmemory`

> Note: after stopping service data will save on disk

To stop application use `make stop`

### Service provides an API

create new short link method:

`POST localhost/link?full_link=<full_link>`

get full url by short version:

`GET localhost/link?short_link=<short_link>`

### The output is a json data like this:
```json
{
    "full_link": "google.com"
}
```
```json
{
    "short_link": "WDBsSk9hYr"
}
```

and in case of error
```json
{
    "error": "link not found"
}
```
