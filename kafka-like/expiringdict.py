"""
MIT License

Copyright (c) 2019 David Parker

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"""

try:
    from collections.abc import MutableMapping
except ImportError:  # Will allow usage with virtually any python 3 version
    from collections import MutableMapping

from threading import Timer, Thread, Lock
from sortedcontainers import SortedKeyList
from time import time, sleep


class ExpiringDict(MutableMapping):
    def __init__(self, ttl=None, interval=0.100, *args, **kwargs):
        """
        Create an ExpiringDict class, optionally passing in a time-to-live
        number in seconds that will act globally as an expiration time for keys.
        If omitted, the dict will work like a normal dict by default, expiring
        only keys explicity set via the `.ttl` method.
        """
        self.__store = dict(*args, **kwargs)
        self.__keys = SortedKeyList(key=lambda x: x[0])
        self.__ttl = ttl
        self.__lock = Lock()
        self.__interval = interval

        Thread(target=self.__worker, daemon=True).start()

    def flush(self):
        now = time()
        max_index = 0
        with self.__lock:
            for index, (timestamp, key) in enumerate(self.__keys):
                if timestamp > now:
                    max_index = index
                    break
                try:
                    del self.__store[key]
                except KeyError:
                    pass
            del self.__keys[0:max_index]

    def __worker(self):
        while True:
            self.flush()
            sleep(self.__interval)

    def __setitem__(self, key, value):
        """
        Set `value` of `key` in dict. `key` will be automatically
        deleted if the `ttl` option was provided for this object.
        """
        if self.__ttl:
            self.__set_with_expire(key, value, self.__ttl)
        else:
            self.__store[key] = value

    def ttl(self, key, value, ttl):
        """
        Set `value` of `key` in dict to expire after `ttl` seconds.
        Overrides object level `ttl` setting.
        """
        self.__set_with_expire(key, value, ttl)

    def __set_with_expire(self, key, value, ttl):
        self.__lock.acquire()
        self.__keys.add((time() + ttl, key))
        self.__store[key] = value
        self.__lock.release()

    def __delitem__(self, key):
        del self.__store[key]

    def __getitem__(self, key):
        return self.__store[key]

    def __iter__(self):
        return iter(self.__store)

    def __len__(self):
        return len(self.__store)
