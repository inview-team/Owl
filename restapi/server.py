import os
from flask import Flask, jsonify, request
from model import db,init_db, Alarms, Logs, Settings
from flask_cors import CORS
from dotenv import load_dotenv, find_dotenv

from clickhouse_driver import connect

'''nodes = {"ns=2;i=9": "pressure", "ns=2;i=10": "humidity", "ns=2;i=11": "roomTemperature", "ns=2;i=12": "workingAreaTemperatur", "ns=2;i=13": "pH", "ns=2;i=14": "weight", "ns=2;i=15": "fluidFlow", "ns=2;i=16": "co2"}

conn = connect('clickhouse://clickhouse-svc')
cursor = conn.cursor()
'''

load_dotenv(find_dotenv())

app = Flask(__name__)
app.config.from_object(__name__)

CORS(app, resources={r'/*': {'origins': '*'}})

database_uri = 'mysql+pymysql://{dbuser}:{dbpass}@{dbhost}/{dbname}'.format(
    dbuser=os.environ['DBUSER'],
    dbpass=os.environ['DBPASS'],
    dbhost=os.environ['DBHOST'],
    dbname=os.environ['DBNAME']
)

app.config.update(
    SQLALCHEMY_DATABASE_URI=database_uri,
    SQLALCHEMY_TRACK_MODIFICATIONS=False,
)


db.init_app(app)
with app.app_context():
    init_db()


'''
@app.route('/projects', methods=['GET','POST'])
def get_metrics():
    cursor.execute('SELECT * FROM metrics WHERE timestamp > now() - 100')
    metrics = cursor.fetchall()
'''

@app.route('/alarms',methods=['GET','POST'])
def get_alarms():
    if request.method == 'POST':
        time = str(request.json.get('time'))
        info = str(request.json.get('info'))

        alarm_request = Alarms(id=None, time=time, info=info)
        db.session.add(alarm_request)
        db.session.commit()
        return jsonify({'message':'Alarm added'})
    else:
        records = Alarms.query.all()
        return jsonify({'alarms': [record.serialize() for record in records]})
    
@app.route('/logs',methods=['GET','POST'])
def get_logs():
    if request.method == 'POST':
        time =str(request.json.get('time'))
        info =str(request.json.get('info'))

        log_request = Logs(id=None,time=time,info=info)
        db.session.add(log_request)
        db.session.commit()
    else:
        records = Logs.query.all()
        print(records)
        return jsonify({'logs': [record.serialize() for record in records]})


@app.route('/settings', methods=['GET'])
def get_settings():
    records = Settings.query.all()
    return jsonify({'settings':[record.serialize() for record in records]})

@app.route('/settings_update/<metric_id>',methods=['PUT'])
def update_settings(metric_id):
    responce = request.get_json()
    id=int(metric_id)

    setting = Settings.query.get(id)
    setting.name = responce['metric']
    setting.from_bord = responce['from']
    setting.to_bord = responce['to']
    db.session.commit()

if __name__ == "__main__":
    app.run()
