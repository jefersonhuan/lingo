# Lingo CLI

Lingo is a CLI MongoDB utility that helps you transfer data from one server to another 
without any other package (even mongodump and mongorestore) and on top of that it is
powered by Go's concurrency.

### Usage

##### Add server

To add a server you need to simply execute de `add` command:

```shell script
$ lingo add
Host: [default localhost] 
Port: [default 27017] 
User: [default blank] 
Password: [default blank] 
Please, insert an identifier for this server
ID: local
Do you want to test the connection right now? [y/n] y
Testing connection...
Successfully connected to local
Successfully saved server
``` 

Or, if you already have a URI (or need more complete options):

```shell script
$ lingo add --from-uri "mongodb://localhost:27017" --name "local"
```

#### Transfer

To transfer data between servers you only need to indicate a `from` (source) 
and a `to` (target)

```shell script
$ lingo transfer --from local --to new-qa
```

And watch the magic :)

#### Possible issues

Windows support is a maybe.

# TODO
- Transfer indexes
- Enables SSH Tunneling
