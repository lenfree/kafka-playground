kafka-playground
----------------

This would have all dirty once off administration tasks to deal with Kafka
and written in Go.

listTopics
===========

This returns all topics except for specific names.

```bash
$ cp env.example .env
```

Usage:
======

```bash
$ go run listTopics.go
```

To delete topics:

```bash
$  go run listTopics.go > listTopics.txt
$ cat listTopics.txt|while read line; do
    kafka-topics --zookeeper 127.0.0.1:2183/<cluster_name> --delete --topic $line;
  done
```
