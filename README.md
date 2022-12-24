# (WIP) Ditz
A zero-maintenance Git server for high-maintenance people. Suitable for random experiments, projects you don't
care about losing, and integration testing.

## Do you need Ditz?
If all you need is a no-auth Git *remote* (not necessarily a server), consider trying Git's
[Local protocol](https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols#_local_protocol) first.

## Setup
Ditz requires no configuration by default - just point a Git client at it and get started.

```sh
docker run command here
```

```sh
git clone git://localhost/repo.git
```

## Authentication
Ditz comes preconfigured with no authentication whatsoever. However, the following forms of authentication are configurable:

* SSH password authentication
* SSH key authentication
* HTTP header authentication
* HTTP URL authentication

This covers the majority of authentication schemes used by Git hosting providers.