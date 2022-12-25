# Gittlz
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/karashiiro/gittlz)](https://github.com/karashiiro/gittlz/blob/main/go.mod)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/karashiiro/gittlz/build.yml)](https://github.com/karashiiro/gittlz/actions/workflows/build.yml)
[![GitHub](https://img.shields.io/github/license/karashiiro/gittlz)](https://github.com/karashiiro/gittlz/blob/main/LICENSE)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/karashiiro/gittlz)](https://hub.docker.com/r/karashiiro/gittlz)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/karashiiro/gittlz)](https://hub.docker.com/r/karashiiro/gittlz)

A Git server for the laziest of us. Write and test your Git utilities without any hassle.

*Gittlz is still in active development. Expect some issues while using it.*

- [Gittlz](#gittlz)
  - [Do you need Gittlz?](#do-you-need-gittlz)
  - [Usage](#usage)
  - [Authentication](#authentication)
    - [SSH password authentication](#ssh-password-authentication)
    - [SSH key authentication](#ssh-key-authentication)
    - [HTTP URL authentication](#http-url-authentication)
    - [HTTP basic authentication](#http-basic-authentication)
  - [Containerless](#containerless)
  - [Architecture](#architecture)

## Do you need Gittlz?
If all you need is a no-auth Git *remote* (not necessarily a server), consider trying Git's
[Local protocol](https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols#_local_protocol) first.

Gittlz is meant to work in place of a live Git host for development purposes, and not to act as a
production server in any form. If that's what you were looking for, try [Gitea](https://gitea.io/en-us/),
[Gogs](https://gogs.io), [OneDev](https://github.com/theonedev/onedev), or
[Soft Serve](https://github.com/charmbracelet/soft-serve), among the many projects floating around.

This will likely be repeated several times throughout this documentation:
*Do not use Gittlz as a production Git server.*

## Usage
Gittlz requires no configuration by default - just point a Git client at it and get started:

```sh
docker run --rm -it --name=gittlz -p 6177:6177 -p 9418:9418 karashiiro/gittlz:latest
```

If you want to use a persistent directory for repositories, mount it to `/srv/git`:

```sh
docker run --rm -it --name=gittlz -v /path/to/repos:/srv/git:rw -p 6177:6177 -p 9418:9418 karashiiro/gittlz:latest
```

Repositories should be [bare repositories](https://git-scm.com/book/en/v2/Git-on-the-Server-Getting-Git-on-a-Server)
on the server. The Gittlz CLI abstracts away this setup process:

```sh
CGO_ENABLED=0 go install github.com/karashiiro/gittlz@v0.3.0
gittlz create-repo repo
```

Then, you can clone repositories from a Git client outside the container:

```sh
git clone git://localhost/repo.git
```

And that's it! Enjoy your Gittlz.

The Gittlz [Docker image](https://hub.docker.com/repository/docker/karashiiro/gittlz) makes this setup process
nearly as simple as it can be. The image is based on Alpine Linux, but it includes a full Git installation, which
can be used to manually perform operations inside the container. `sh` is available as a basic shell for manual
repository setup, if needed.

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
docker run --rm -it -p 6177:6177 -p 22:22 karashiiro/gittlz:latest gittlz serve --protocol=ssh --password=password
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
docker run --rm -it -p 6177:6177 -p 22:22 karashiiro/gittlz:latest gittlz serve --protocol=ssh
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
docker run --rm -it -p 6177:6177 -p 80:80 karashiiro/gittlz:latest gittlz serve --protocol=http --username=gitt --password=lz
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

## Containerless
The Gittlz container attempts to abstract configuration as much as possible, without sacrificing
maintainability or debuggability. However, Gittlz is also just a CLI application, and can be built
and run in other environments.

Building Gittlz from sources is simple, just disable `cgo` (optional) and install it like any other
Go application. In `sh`-like shells, this is done as follows:

```sh
CGO_ENABLED=0 go install github.com/karashiiro/gittlz@v0.3.0
```

Gittlz has runtime dependencies on the standard `git` toolkit and `git-http-backend`. `git-http-backend`
is a CGI script sometimes offered as part of a separate `git-daemon` package. For Windows users,
[Git for Windows](https://gitforwindows.org) includes everything needed in a single installer.

Refer to the `--help` commands such as `gittlz --help` and `gittlz serve --help` for configuration
options.

## Architecture
Git's server functionality is mostly usable out of the box, and the official handbook even dedicates
an entire [chapter](https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols) to describing how
to configure and use it. However, that's more configuration than anyone should want to do if they only
want a disposable HTTP or SSH Git server, and don't care about security at all.

With this being the case, Gittlz is just a very thin wrapper around Git itself, with the exception of
the handling for the SSH protocol. Each protocol has a different strategy used to wrap it.

| Protocol | Strategy                                                                                                                                                                                                                                                  |
| -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Git      | `git daemon` is launched as a subprocess. That's it.                                                                                                                                                                                                      |
| HTTP     | [`net/http/cgi`](https://pkg.go.dev/net/http/cgi) (yes, that's part of the Go standard library) is used to interface with `git-http-backend`. Gittlz adds some authentication middleware to simulate a typical managed Git provider.                      |
| SSH      | [charmbracelet/wish](https://github.com/charmbracelet/wish) is used to set up a simple SSH server in front of Git. Gittlz's implementation is almost an exact copy of Wish's [Git example](https://github.com/charmbracelet/wish/tree/main/examples/git). |

The [`serve`](https://github.com/karashiiro/gittlz/blob/main/cmd/serve.go) command is used to select
which protocol is used at runtime.

Finally, a control API is put on top for simpler repository creation within the container.
