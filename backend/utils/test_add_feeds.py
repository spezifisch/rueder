#!/usr/bin/env python3

import requests
from test_login import test_login


urls = """
https://blog.fefe.de/rss.xml?html
https://blog.fefe.de/rss.xml
https://www.heise.de/rss/heise-atom.xml
https://www.heise.de/rss/heise.rdf
https://www.theguardian.com/international/rss
https://mjg59.dreamwidth.org/data/rss
https://fair.org/feed/
"""


def add_feed(endpoint: str, jwt: str, data=dict):
    headers = {
        "authorization": f"Bearer {jwt}"
    }
    r = requests.post(endpoint, json=data, headers=headers)
    if r.status_code != 200:
        raise Exception(f"request failed: {r.text}")


def run():
    jwt = test_login()

    endpoint = "http://127.0.0.1:8080/api/v1/feed"

    for url in urls.split():
        if not url:
            continue

        data = {
            "url": url,
        }
        add_feed(endpoint, jwt, data)


if __name__ == "__main__":
    run()
