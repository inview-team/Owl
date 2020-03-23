from clickhouse_driver import Client
import requests

nodes = {"ns=2;i=9": "pressure", "ns=2;i=10": "humidity", "ns=2;i=11": "roomTemperature", "ns=2;i=12": "workingAreaTemperatur", "ns=2;i=13": "pH", "ns=2;i=14": "weight", "ns=2;i=15": "fluidFlow", "ns=2;i=16": "co2"}
min = {"pressure": 1001, "humidity": 25, "roomTemperature":15, "workingAreaTemperatur":200, "pH":3, "weight":100, "fluidFlow":10, "co2":0}
max = {"pressure": 1001, "humidity": 25, "roomTemperature":15, "workingAreaTemperatur":200, "pH":3, "weight":100, "fluidFlow":10, "co2":0}

while True:
    client = Client('10.244.0.113')
    metrics = client.execute('select * from metrics')

    for mt in metrics:
        currNode = nodes.get(mt[0])
        cmin = min.get(currNode)
        cmax = max.get(currNode)
        if (float(mt[2]) < cmin ) or (float(mt[2]) > cmax):
            # send alarm
            data = {'': currNode, '': float(mt[2])}
            r = requests.post(url = 'localhost:5000/alarms', data = data)
