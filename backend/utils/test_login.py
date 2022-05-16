#!/usr/bin/env python3

import requests


def login(endpoint: str, user: str, passwd: str) -> str:
    data = {
        "username": user,
        "password": passwd,
    }
    r = requests.post(endpoint, json=data)
    if r.status_code != 200:
        raise Exception(f"login failed: {r.text}")
    return r.text


def test_login() -> str:
    login_endpoint = "http://127.0.0.1:8082/login"
    user = "bob"
    passwd = "secret"
    return login(login_endpoint, user, passwd)


if __name__ == "__main__":
    jwt = test_login()
    print("Bearer", jwt)
