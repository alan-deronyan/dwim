import requests
import time

class FlureeClient:
    def __init__(self, url, ledger):
        self.url = url
        self.ledger = ledger

    def transact(self, context_map, inserts, deletes):

        # transaction body
        transaction = {
            "ledger": self.ledger,
        }

        if inserts is not None:
            transaction["inserts"] = inserts
        if deletes is not None:
            transaction["deletes"] = deletes
        if context_map is not None:
            transaction["context"] = context_map

        headers = {
            "Content-Type": "application/json",
            #"Authorization": api_key,
            #'Accept': 'text/plain'
        }

        response = requests.post(self.url + "/fluree/transaction",  headers=headers, json=transaction)
        if response.ok:
            data = response.json()
            print(data)
