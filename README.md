# sshw

ssh client wrapper for automatic login.

This is a fork of [yinheli/sshw](https://github.com/yinheli/sshw), with the following features added:

- Support setting default user/password

- Support short hostname by setting up domain ( concatenated as `hostname.domain` )

- Command line support

  ```shell
  sshw user:pass@host
  # use default password
  sshw user@host
  # use default user/password
  sshw host
  ```

- Support using flags before selecting host from list

  ```shell
  # use specified user and default password for selected host
  sshw -u admin
  # use specified user and password (will prompt for input) for selected host
  sshw -u admin -pass
  # use specified port for selected host
  sshw -p 33
  ```
## install
```shell
go get https://github.com/lixvbnet/sshw
```

## config

put config file in `~/.sshw` or `~/.sshw.yml` or `./.sshw` or `./.sshw.yml` .

[config example](./sshlib/config_example.yml):

```yaml
settings:
  domain: example.com
  logins:
    - user: admin
      password: password

default:
  user: root
  password: 123456


nodes:
  - name: nodeA
    alias: nodeA
    host: 192.168.1.10

  - name: nodeB
    alias: nodeB
    host: 192.168.1.11
    user: root
    password: Password
```