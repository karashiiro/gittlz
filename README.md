# Gittlz
A simple Git server for the laziest people. Hack on it, destroy it, and test against it.

## Do you need Gittlz?
If all you need is a no-auth Git *remote* (not necessarily a server), consider trying Git's
[Local protocol](https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols#_local_protocol) first.

Gittlz is meant to work in place of a live Git host for development purposes, and not to act as a
production server in any form. If that's what you were looking for, try [Gitea](https://gitea.io/en-us/),
[Gogs](https://gogs.io), [OneDev](https://github.com/theonedev/onedev), or
[Soft Serve](https://github.com/charmbracelet/soft-serve), among the many projects floating around.

This will likely be repeated several times throughout this documentation:
*Do not use Gittlz as a production Git server.*

## Setup
Gittlz requires no configuration by default - just point a Git client at it and get started.

```sh
docker run --rm -it -p 9418:9418 karashiiro/gittlz:latest
```

If you want to use a persistent directory for repositories, mount it to `/srv/git`:

```sh
docker run --rm -it -v /path/to/repos:/srv/git:rw -p 9418:9418 karashiiro/gittlz:latest
```

Repositories should be [bare repositories](https://git-scm.com/book/en/v2/Git-on-the-Server-Getting-Git-on-a-Server)
on the server. To create a new bare repository, run:

```sh
git init --bare repo.git
```

Then, you can clone repositories from a Git client outside the container:

```sh
git clone git://localhost/repo.git
```

## Authentication
Gittlz comes preconfigured with no authentication whatsoever. All of the optional authentication methods provided
are intentionally insecure - Gittlz favors convenience over security where possible.

*Do not use Gittlz as a production Git server.*

The following forms of authentication are configurable:

* SSH password authentication
* SSH key authentication
* HTTP URL authentication
* HTTP basic authentication

This covers the majority of authentication schemes used by Git hosting providers.

### SSH password authentication
Start the server with a command override, replacing the port mapping and password options as needed:

```sh
docker run --rm -it -p 22:22 karashiiro/gittlz:latest gittlz serve --protocol=ssh --password=password
```

Then, clone repositories by providing the password interactively:

```sh
git clone ssh://localhost/repo.git
# Cloning into 'repo'...
# you@localhost's password: password
```

It is not possible to use this authentication method non-interactively.

### SSH key authentication
Note that Gittlz will not validate the SSH key used to access the server. This is intentional, as key
configuration has little to do with a Git server's public interface. That said, if your use case requires
SSH key auth failures, open an issue describing your intended workflow.

Start the server with a command override, replacing the port mapping as needed:

```sh
docker run --rm -it -p 22:22 karashiiro/gittlz:latest gittlz serve --protocol=ssh
```

Then, clone repositories with a Git client:

```sh
git clone ssh://localhost/repo.git
```

You may also want to override your Git client's SSH command to avoid host key verification errors. This
is done by setting the `GIT_SSH_COMMAND` environment variable, which is shell-specific. In `sh`-like shells,
this can simply be prepended to the command:

```sh
GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no" git clone ssh://localhost/repo.git
```

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
