from clickhouse_driver import connect

conn = connect('clickhouse://clickhouse-svc')
cursor = conn.cursor()

cursor.execute('SELECT * FROM metrics ')
metrics = cursor.fetchall()

print(metrics)
