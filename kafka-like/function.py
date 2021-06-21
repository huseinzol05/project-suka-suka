from wired import WTDict
from expiringdict import ExpiringDict
from pathlib import Path
import os

HOME = os.environ.get(
    'CACHE_KAFKA_LIKE',
    os.path.join(str(Path.home()), 'KAFKA_LIKE'),
)

if not os.path.exists(HOME):
    os.makedirs(HOME)

WT, WT_OFFSET, ED = {}, {}, {}


def set_topic_dict(topic):
    if topic not in WT:
        path = f'{HOME}/{topic}'
        if not os.path.exists(path):
            os.makedirs(path)
        WT[topic] = WTDict(path)


def set_offset_topic_dict(topic):
    if topic not in WT_OFFSET:
        path = f'{HOME}/{topic}-offset'
        if not os.path.exists(path):
            os.makedirs(path)
        WT_OFFSET[topic] = WTDict(path)


def set_consumergroup_dict(topic):
    pass


def write_wiredtiger(topic, key, value):
    set_topic_dict(topic)
    WT[topic][key] = value


def read_wiredtiger(topic, key, default=None):
    set_topic_dict(topic)
    return WT[topic].get(key, default=default)


def insert(topic, values):
    pass


def check_heartbeat(topic, consumer_group, consumer_id, heartbeat_second=5):


def read(topic, batch_size, consumer_group, consumer_id):
    pass


from tqdm import tqdm
import time

for i in tqdm(range(10000)):
    write_wiredtiger('test', str(i), value=i)

for i in tqdm(range(10000)):
    print(read_wiredtiger('test', str(i)))

print(WT)
print(read_wiredtiger('test', 'a'))
