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

def get_one_metric(metric):
    getData = False
    while getData != True:
        try:
            r = requests.get(url='http://127.0.0.1:1337/settings/{}'.format(metric))
            result = r.json()
            getData = True
        except:
            print('Error')
    return result

def add_new_admin(chat_id):
    url='http://127.0.0.1:1337/admin'
    headers = {'Content-type': 'application/json',  # Определение типа данных
               'Accept': 'text/plain',
               'Content-Encoding': 'utf-8'}
    data = {
        'chat_id': chat_id
    }
    r=requests.post(url, data=json.dumps(data), headers=headers)
    result = r.json()
    return result
