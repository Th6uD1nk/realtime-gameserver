# RTGS Client

## Run the client

```bash
go run rtgs-client
````

## Build the client for pc

```bash
CGO_CPPFLAGS="-v" go build -x ./...
```

## Build the client for mobile

```bash
gomobile build -target=android -androidapi=33 -tags=mobile .
```


