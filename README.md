## cogmote hub

This project is the central database for cogmoteGO

## Getting started

First, create the dedicated `secrets` folder in the project root and switch into it:

```shell
mkdir -p secrets
cd secrets
```

### Generate the JWT key pair

Run the following commands in the current `secrets` directory to generate a PEM-formatted Ed25519 key pair:

```shell
openssl genpkey -algorithm Ed25519 -out jwt-private-key.pem
openssl pkey -in jwt-private-key.pem -pubout -out jwt-public-key.pem
```

### Create the database password

Create `password.txt` in the same directory and store the database password inside.
