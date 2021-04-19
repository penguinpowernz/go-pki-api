# PKI API

This uses [Sudhi Herle](https://github.com/opencoff/go-pki)'s PKI database library to create
an API that allows you to manage your public key infrastructure via a REST API.

## Building

Use the `Makefile`:

    make build

Then you have the binary in the `bin` directory.

## Usage

The following routes are available:

### GET /servers

Lists all the servers in the database:

```
{
    "common_name":      "example.com",
    "serial":           "0x007165466a",
    "expired":          false,
    "expires_at":       1785423101,
    "expires_at_human": "2009-11-10 23:00:00 +0000 UTC",
}
```

### GET /servers/:cn/export

Export the given servers cert, key and optionally chain and CA cert.

Returns a payload like:

```
{
    "pem": "the certs",
    "key": "the private key",
    "ca":  "the CA cert"
}
```

### POST /servers/:cn

Create a new server with the given JSON payload of options:

```
{
    "domain_names": ["example.com"],
    "ips": ["127.0.0.1"],
    "validity_days": 365,
    "sign_with": "optional",
    "password": "optional",
}
```

### DELETE /servers/:cn

Delete a server with the given common name from the database.

Returns 204 when successful.

### GET /users

Lists all the users in the database:

```
{
    "common_name":      "bob@example.com",
    "serial":           "0x007165466a",
    "expired":          false,
    "expires_at":       1785423101,
    "expires_at_human": "2009-11-10 23:00:00 +0000 UTC",
}
```

### GET /users/:cn/export[?ca=true][&chain=true]

Export the given users cert, key and optionally chain and CA cert.

Returns a payload like:

```
{
    "pem": "the certs",
    "key": "the private key",
    "ca":  "the CA cert"
}
```

### POST /users/:cn

Create a new server with the given JSON payload of options:

```
{
    "email": "optional (will use CN if that is an email)"
    "validity_days": 365,
    "sign_with": "optional",
    "password": "optional",
}
```

### DELETE /users/:cn

Delete a user with the given common name from the database.

Returns 204 when successful.

### GET /crl/:days

Generate a CRL with the validity of it given in days. Outputs the CRL
in plain text for easy curling to a location.