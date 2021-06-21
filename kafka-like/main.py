from fastapi import FastAPI

app = FastAPI()


@app.get('/')
def root():
    return {'message': 'this is a fake kafka, called `kafka-like`.'}


@app.get('/{topic}/')
def topic_describe(topic: str):
    pass


@app.get('/{topic}/insert')
def topic_insert(topic: str):
    pass


@app.get('/{consumer}/')
def consumer_describe(consumer: str):
    pass


@app.get('/{consumer}/consume')
def consumer_consume(consumer: str):
    pass
