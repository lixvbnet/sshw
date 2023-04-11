# sshw

ssh client wrapper for automatic login.

- Select / Search / Login
  ![](./resources/01.gif)

- Use Command Line Arguments
  ![](./resources/02.gif)

- Run / Start Command
  ![](./resources/03.gif)

- And if you need to manage a bunch of hosts under the same domain, such as

  ```
  hostA.your.domain.com
  hostB.your.domain.com
  hostC.your.domain.com
  ...
  ```

  you can write your domain in the config file (see `config` section below), and then you can just specify ***short hostname*** in the command line.

  ```shell
  sshw root@hostA
  ```

   `sshw` will concatenate domain name for you. FQDN won't be affected. (Note: make sure there is no `.` in your ***short hostname*** otherwise `sshw` would consider it as an FQDN.)

- AND you need to set separate default passwords for different users? Just config them in "settings - logins" section (see config section below. With that example config, when you connect with user admin such as `sshw admin@host` or `sshw -u admin` , the password for admin user, in this case "password", will be used; but if a user is not configured, default password will be used.)

## install

```shell
go install github.com/lixvbnet/sshw@latest
```

or download pre-compiled binaries from [Releases](https://github.com/lixvbnet/sshw/releases) page.

## config

put config file `.sshw.yml` in home directory or in the same directory of the binary.

[config example](./sshlib/config_example.yml):

```yaml
settings:
  domain: example.com
  logins:
    - user: admin
      password: password
    - user: testuser
      password: testpass
    - user: anotheruser
      password: anotherpass

default:
  user: root
  password: 123456


nodes:
  - name: nodeA
    alias: nodeA
    host: 192.168.1.10

  - name: nodeB
    alias: nodeB,secondNode,11
    host: 192.168.1.11
    user: root
    password: Password
```

## usage

```shell
sshw [options] [target] [command]
```

where `target` is of the form `[user[:pass]@]host`. A `host` can be an FQDN, short hostname, configured node name or alias. Run `sshw -h` for a full list of available options.

---
This is a complete rewrite, except the UI interface, of [yinheli/sshw](https://github.com/yinheli/sshw) for extensibility. And the following features were added:

- Support setting default user/password
- Support short hostnames by setting up domain ( concatenated as `hostname.domain` )
- Support command line arguments
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
- Or even using both (The flags have top priority)
  ```shell
  # will prompt for password
  sshw -pass user@host
  ```
- Run command on host
  ```shell
  # run command and get its output
  sshw -pass user@host "echo hello"
  # request a terminal and start command
  sshw -t -pass user@host "/bin/bash"
  ```
- Support more ciphers and kex algorithms


