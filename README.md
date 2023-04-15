### Theory

#### Proxy server

A proxy server is a server application that provides
a gateway between users and the internet.
It acts as an intermediary between a client
requesting a resource and the server providing it.
Proxy servers also help to improve network performance,
as it caches frequently requested resources,
allowing for faster loading times.
It also could be used for anonymizing purposes and for
browsing sites that are unavailable in some countries.

#### Reverse proxy

A reverse proxy is a server that accepts requests from
a client, forwards the request to one of many other
servers, and then returns the response to the client.
It is typically positioned at a network's edge and can
be used to protect web servers from attacks, as well as
providing performance and reliability benefits.
Additionally, reverse proxies can be used to modify
request headers and fine-tune buffering of responses.

### Build

```bash
$ make server
```

```bash
$ make proxy
```

### Features

- [ ] Caching for content
- [ ] Load balancing
- [ ] Content filtering
- [ ] Compression