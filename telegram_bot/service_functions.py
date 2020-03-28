import requests
import json


def load_setting():
    getData = False
    while getData != True:
        try:
            r = requests.get(url='http://127.0.0.1:1337/settings')
            result = r.json()
            getData = True
        except:
            print('Error')
    return result
