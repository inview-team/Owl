from clickhouse_driver import Client
import json
import requests

nodes = {"ns=2;i=9": "pressure", "ns=2;i=10": "humidity", "ns=2;i=11": "roomTemperature", "ns=2;i=12": "workingAreaTemperature", "ns=2;i=13": "pH", "ns=2;i=14": "weight", "ns=2;i=15": "fluidFlow", "ns=2;i=16": "co2"}
min = {"pressure": 1001, "humidity": 25, "roomTemperature":15, "workingAreaTemperature":200, "pH":3, "weight":100, "fluidFlow":10, "co2":0}
max = {"pressure": 1099, "humidity": 40, "roomTemperature":46, "workingAreaTemperature":500, "pH":8, "weight":900, "fluidFlow":50, "co2":30}

while True:
    client = Client('clickhouse-svc')
    metrics = client.execute('select * from metrics')

    for mt in metrics:
        currNode = nodes.get(mt[0])
        cmin = min.get(currNode)
        cmax = max.get(currNode)
        if (float(mt[2]) < cmin ) or (float(mt[2]) > cmax):
            # send alarm
            headers = {'Content-type': 'application/json',
                    'Accept':'text/plain',
                    'Content-Encoding':'utf-8'}
            data = {'time': mt[1], 'info': currNode}
            r = requests.post(url = 'http://127.0.0.1:5000/alarms', data = json.dumps(data), headers = headers)
            print(r)
