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

Alter topic retention period rather than broker wide config to free up diskspace.
Below is to 7days only:

```bash
$ kafka-configs --zookeeper \
    localhost:2182/<cluster> \
    --alter --entity-type topics \
    --entity-name metrics \
    --add-config retention.ms=604800000
```
