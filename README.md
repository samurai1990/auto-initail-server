# auto-initail-server

this app is helper for admin server that initialize pure server.this tools does the following:
- copy ssh key
- change ssh port
- add AllowUsers
- adduser debian(default)
- disable PasswordAuthentication
- install ufw
- allow ssh port in ufw

## Installation

Use golang 1.21.4.

```bash
make
```

## Usage

```bash
Usage: ./auto-initail-server [flags]
        Example: auto-initail-server -c conf.yaml -f ~/.ssh/rsa_pub

  -c string
        yaml file path
  -f string
        PuplicKey

```

## License

[MIT](https://choosealicense.com/licenses/mit/)