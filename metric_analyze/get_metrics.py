from clickhouse_driver import connect

nodes = {"ns=2;i=9": "pressure", "ns=2;i=10": "humidity", "ns=2;i=11": "roomTemperature", "ns=2;i=12": "workingAreaTemperatur", "ns=2;i=13": "pH", "ns=2;i=14": "weight", "ns=2;i=15": "fluidFlow", "ns=2;i=16": "co2"}

conn = connect('clickhouse://178.128.116.174:30000')
cursor = conn.cursor()

cursor.execute('SELECT * FROM metrics')
metrics = cursor.fetchall()
print(metrics)
