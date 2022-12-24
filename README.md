# (WIP) gittlz
A zero-maintenance Git server for the laziest people. Suitable for random experiments, projects you don't
care about losing, and integration testing.

## Do you need gittlz?
If all you need is a no-auth Git *remote* (not necessarily a server), consider trying Git's
[Local protocol](https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols#_local_protocol) first.

## Setup
gittlz requires no configuration by default - just point a Git client at it and get started.

```sh
docker run --rm -it -p 9418:9418 karashiiro/gittlz:latest
```

If you want to use a persistent directory for repositories, mount it to `/srv/git`:

```sh
docker run --rm -it -v /path/to/repos:/srv/git:rw -p 9418:9418 karashiiro/gittlz:latest
```

Then, you can clone repositories from a Git client outside the container:

```sh
git clone git://localhost/repo.git
```

## Authentication
gittlz comes preconfigured with no authentication whatsoever. However, the following forms of authentication are configurable:

* (todo) SSH password authentication
* (todo) SSH key authentication
* HTTP URL authentication
* HTTP basic authentication

This covers the majority of authentication schemes used by Git hosting providers.

### HTTP URL authentication
See [HTTP basic authentication](#http-basic-authentication). The same setup applies, but the username and
password can be embedded in the URL:

```sh
git clone http://gitt:lz@localhost/repo.git
```

This authentication scheme is both insecure and deprecated by many clients. Some Git clients will even
silently strip the credentials out of the URL. cURL automatically converts it into an `Authorization`
header.

Be prepared to debug issues yourself.

### HTTP basic authentication
Start the server with a command override, replacing the port mapping, username, and password options as needed:

```sh
docker run --rm -it -p 80:80 karashiiro/gittlz:latest gittlz serve --protocol=http --username=gitt --password=lz
```

Then, make sure to base64-encode the username and password somewhere locally. Most operating systems and shells
have a means of doing this. In Powershell, for example:

```powershell
$gittlzAuth = "gitt:lz"
$B64gittlzAuth = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($gittlzAuth))
```

Finally, add the `http.extraHeader` option to all of your Git commands:

```sh
git -c http.extraHeader="Authorization: Basic $B64gittlzAuth" clone http://localhost/repo.git
```
