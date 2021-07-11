import json
from flask import Flask
from flask import render_template

app = Flask(__name__)


@app.route("/")
def hello():
    return render_template('./index.html')


@app.route("/test")
def test():
    return render_template('./selangor.html')


@app.route("/selangor")
def selangor():
    with open('data/selangor-test.json') as fopen:
        data = json.load(fopen)
    return json.dumps(data)


if __name__ == "__main__":
    app.run(debug=True, port=5000)
