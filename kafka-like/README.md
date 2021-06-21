# Kafka-like

Hopefully I am able to write a Kafka-like using 100% Python3 and HTTP interface with ~60k write per second, ~60k read per second using [wiredtiger](http://source.wiredtiger.com/1.4.2/group__wt__python.html) on my Macbook Pro.

## how-to

1. Install requirements,

```bash
pip install -r requirements.txt
```

2. Run `kafka-like`,

```bash
uvicorn main:app --reload
```

FastAPI variables should able to set using bash environment.

## Lack of

1. Replication factor.
2. Multibrokers.
3. partitions.