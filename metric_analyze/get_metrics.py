from clickhouse_driver import connect

nodes = {"ns=2;i=9": "pressure", "ns=2;i=10": "humidity", "ns=2;i=11": "roomTemperature", "ns=2;i=12": "workingAreaTemperatur", "ns=2;i=13": "pH", "ns=2;i=14": "weight", "ns=2;i=15": "fluidFlow", "ns=2;i=16": "co2"}

conn = connect('clickhouse://clickhouse-svc')
cursor = conn.cursor()

while True:
    cursor.execute('SELECT * FROM metrics WHERE timestamp > now() - 100')
    metrics = cursor.fetchall()

    for mt in metrics:
        print(mt[0], ": ", mt[2], "\n")
