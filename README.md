# (WIP) Ditz
A zero-maintenance Git server for high-maintenance people. Suitable for random experiments, projects you don't
care about losing, and integration testing.

## Do you need Ditz?
If all you need is a no-auth Git *remote* (not necessarily a server), consider trying Git's
[Local protocol](https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols#_local_protocol) first.

## Setup
Ditz requires no configuration by default - just point a Git client at it and get started.

```sh
docker run --rm -it -p 9418:9418 karashiiro/ditz:latest
```

If you want to use a persistent directory for repositories, mount it to `/srv/git`:

```sh
docker run --rm -it -v /path/to/repos:/srv/git:rw -p 9418:9418 karashiiro/ditz:latest
```

Then, you can clone repositories from a Git client outside the container:

```sh
git clone git://localhost/repo.git
```

## Authentication
Ditz comes preconfigured with no authentication whatsoever. However, the following forms of authentication are configurable:

* (todo) SSH password authentication
* (todo) SSH key authentication
* HTTP basic authentication
* (todo) HTTP URL authentication

This covers the majority of authentication schemes used by Git hosting providers.

### HTTP basic authentication
Start the server with a command override, replacing the port mapping, username, and password options as needed:

```sh
docker run --rm -it -p 80:80 karashiiro/ditz:latest ditz serve --protocol=http --username=ditz --password=y
```

Then, make sure to base64-encode the username and password somewhere locally. Most operating systems and shells
have a means of doing this. In Powershell, for example:

```powershell
$DitzAuth = "ditz:y"
$B64DitzAuth = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($DitzAuth))
```

Finally, add the `http.extraHeader` option to all of your Git commands:

```sh
git -c http.extraHeader="Authorization: Basic $B64DitzAuth" clone http://localhost/repo.git
```
