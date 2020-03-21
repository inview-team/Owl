from clickhouse_driver import connect

conn = connect('clickhouse://localhost')
cursor = conn.cursor()

cursor.execute('SELECT * FROM metrics')
metrics = cursor.fetchall()

print(metrics)
#df = pd.DataFrame()
#for metric in metrics:
#    print(metric)
