# oauth2_proxy

A reverse proxy and static file server that provides authentication using Providers (Google, GitHub, and others)
to validate accounts by email, domain or group.

**Note:** This repository was forked from [bitly/OAuth2_Proxy](https://github.com/bitly/oauth2_proxy) on 27/11/2018.
Versions v3.0.0 and up are from this fork and will have diverged from any changes in the original fork.
A list of changes can be seen in the [CHANGELOG](CHANGELOG.md).

[![Build Status](https://secure.travis-ci.org/pusher/oauth2_proxy.svg?branch=master)](http://travis-ci.org/pusher/oauth2_proxy)

![Sign In Page](https://cloud.githubusercontent.com/assets/45028/4970624/7feb7dd8-6886-11e4-93e0-c9904af44ea8.png)

## Epitech Fork

This fork adds support for an Epitech provider, based on AzureAD and Epitech Intranet groups.

Enjoy ;)

### Build

```
docker build -t samber/epitech-oauth2-proxy:4.0.0 .
```

### Setup

1- Register a new webapp here => https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-register-app

2- Copy Tenant ID + Client ID

3- Navigate to the "Certificates & secrets" page and create a new "client secret".

4- Navigate to "Authentication" page and set the redirect url to `http://localhost:80/oauth2/callback`

5- Navigate to "Manifest" page and set `groupMembershipClaims` to `All`.

```
docker run -d -p 80:80 -p 443:443 \
           samber/epitech-oauth2-proxy \
                 -upstream=http://very-private-webapp:80 \
                 -http-address=0.0.0.0:80 \
                 -redirect-url=http://localhost:80/oauth2/callback \
                 -scope='profile User.Read' \
                 -email-domain=* \
                 -cookie-domain=localhost \
                 -cookie-secure=false \
                 -cookie-secret=somerandomstring1234567890 \
                 -provider=epitech \
                 -azure-tenant ******************************** \
                 -client-id ****************************** \
                 -client-secret '****************************' \
                 -epitech-group adm -epitech-group dpr -epitech-group dpra -epitech-group ape \
                 -epitech-auth-token auth-***********************
```

⚠️  For prod environment, remove the `-cookie-secure=false` argument ;)

## Installation

1.  Choose how to deploy:

    a. Download [Prebuilt Binary](https://github.com/pusher/oauth2_proxy/releases) (current release is `v4.0.0`)

    b. Build with `$ go get github.com/pusher/oauth2_proxy` which will put the binary in `$GOROOT/bin`

    c. Using the prebuilt docker image [quay.io/pusher/oauth2_proxy](https://quay.io/pusher/oauth2_proxy) (AMD64, ARMv6 and ARM64 tags available)

Prebuilt binaries can be validated by extracting the file and verifying it against the `sha256sum.txt` checksum file provided for each release starting with version `v3.0.0`.

```
sha256sum -c sha256sum.txt 2>&1 | grep OK
oauth2_proxy-4.0.0.linux-amd64: OK
```

2.  [Select a Provider and Register an OAuth Application with a Provider](https://pusher.github.io/oauth2_proxy/auth-configuration)
3.  [Configure OAuth2 Proxy using config file, command line options, or environment variables](https://pusher.github.io/oauth2_proxy/configuration)
4.  [Configure SSL or Deploy behind a SSL endpoint](https://pusher.github.io/oauth2_proxy/tls-configuration) (example provided for Nginx)

## Docs

Read the docs on our [Docs site](https://pusher.github.io/oauth2_proxy).

![OAuth2 Proxy Architecture](https://cloud.githubusercontent.com/assets/45028/8027702/bd040b7a-0d6a-11e5-85b9-f8d953d04f39.png)

## Getting Involved

If you would like to reach out to the maintainers, come talk to us in the `#oauth2_proxy` channel in the [Gophers slack](http://gophers.slack.com/).

## Contributing

Please see our [Contributing](CONTRIBUTING.md) guidelines.
