# https://stackoverflow.com/questions/47233562/key-value-store-in-python-for-possibly-100-gb-of-data-without-client-server
# miraculously very fast

import ujson as json
from wiredtiger import wiredtiger_open


WT_NOT_FOUND = -31803


class WTDict:
    """Create a wiredtiger backed dictionary"""

    def __init__(self, path, config='create'):
        self._cnx = wiredtiger_open(path, config)
        self._session = self._cnx.open_session()
        # define key value table
        self._session.create('table:keyvalue', 'key_format=S,value_format=S')
        self._keyvalue = self._session.open_cursor('table:keyvalue')

    def __enter__(self):
        return self

    def close(self):
        self._cnx.close()

    def __exit__(self, *args, **kwargs):
        self.close()

    def _loads(self, value):
        return json.loads(value)

    def _dumps(self, value):
        return json.dumps(value)

    def __delitem__(self, key):
        self._session.begin_transaction()
        self._keyvalue.remove(key)
        self._session.commit_transaction()

    def __contains__(self, key):
        self._session.begin_transaction()
        self._keyvalue.set_key(key)
        found = True
        if self._keyvalue.search() == WT_NOT_FOUND:
            found = False

        self._session.commit_transaction()

        return found

    def get(self, key, default=None):
        if self.__contains__(key):
            return self.__getitem__(key)
        else:
            return default

    def __getitem__(self, key):
        self._session.begin_transaction()
        self._keyvalue.set_key(key)
        if self._keyvalue.search() == WT_NOT_FOUND:
            raise KeyError()
        out = self._loads(self._keyvalue.get_value())
        self._session.commit_transaction()
        return out

    def __setitem__(self, key, value):
        self._session.begin_transaction()
        self._keyvalue.set_key(key)
        self._keyvalue.set_value(self._dumps(value))
        self._keyvalue.insert()
        self._session.commit_transaction()
