import os
from flask import Flask, jsonify, request
from model import db, init_db, Alarms, Logs, Settings, Telegram
from flask_cors import CORS
from dotenv import load_dotenv, find_dotenv
import telebot

# from clickhouse_driver import connect

'''nodes = {"ns=2;i=9": "pressure", "ns=2;i=10": "humidity", "ns=2;i=11": "roomTemperature", "ns=2;i=12": "workingAreaTemperatur", "ns=2;i=13": "pH", "ns=2;i=14": "weight", "ns=2;i=15": "fluidFlow", "ns=2;i=16": "co2"}

conn = connect('clickhouse://clickhouse-svc')
cursor = conn.cursor()
'''

load_dotenv(find_dotenv())

application = Flask(__name__)
application.config.from_object(__name__)

CORS(application, resources={r'/*': {'origins': '*'}})

database_uri = 'mysql+pymysql://{dbuser}:{dbpass}@{dbhost}/{dbname}'.format(
    dbuser=os.environ['DBUSER'],
    dbpass=os.environ['DBPASS'],
    dbhost=os.environ['DBHOST'],
    dbname=os.environ['DBNAME']
)

application.config.update(
    SQLALCHEMY_DATABASE_URI=database_uri,
    SQLALCHEMY_TRACK_MODIFICATIONS=False,
)

db.init_app(application)
with application.app_context():
    init_db()

'''
@app.route('/projects', methods=['GET','POST'])
def get_metrics():
    cursor.execute('SELECT * FROM metrics WHERE timestamp > now() - 100')
    metrics = cursor.fetchall()
'''

# Logs routes

@application.route('/alarms', methods=['GET', 'POST'])
def get_alarms():
    if request.method == 'POST':
        # Add to database
        time = str(request.json.get('time'))
        info = str(request.json.get('info'))
        alarm_request = Alarms(id=None, time=time, info=info)
        db.session.add(alarm_request)
        db.session.commit()

        # Send to telegram
        token = os.environ['TOKEN']
        bot = telebot.TeleBot(token)

        records = Telegram.query.all()
        for record in records:
            bot.send_message(
                record.chat_id,
                info
            )


        return jsonify({'message': 'Alarm added'})
    else:
        records = Alarms.query.all()
        return jsonify({'alarms': [record.serialize() for record in records]})


@application.route('/logs', methods=['GET', 'POST'])
def get_logs():
    if request.method == 'POST':
        time = str(request.json.get('time'))
        info = str(request.json.get('info'))

        log_request = Logs(id=None, time=time, info=info)
        db.session.add(log_request)
        db.session.commit()
    else:
        records = Logs.query.all()
        print(records)
        return jsonify({'logs': [record.serialize() for record in records]})

# Settings routes

@application.route('/settings', methods=['GET'])
def get_settings():
    records = Settings.query.all()
    return jsonify({'settings': [record.serialize() for record in records]})


@application.route('/settings/<metric>', methods=['GET'])
def get_one_metric(metric):
    setting = Settings.query.filter_by(name=metric).first()
    return jsonify({'id': setting.id, 'name': setting.name, 'from': setting.from_bord, 'to': setting.to_bord})


@application.route('/settings_update/<metric_id>', methods=['PUT'])
def update_settings(metric_id):
    responce = request.get_json()
    id = int(metric_id)
    setting = Settings.query.get(id)
    setting.name = responce['metric']
    setting.from_bord = responce['from']
    setting.to_bord = responce['to']
    db.session.commit()


# Telegram routes
@application.route('/admin', methods=['POST'])
def add_admin(admin):
    chat_id = str(request.json.get('chat_id'))
    admin_request = Telegram(id=None, chat_id=chat_id)
    db.session.add(admin_request)
    db.session.commit()
    return jsonify({'message': 'Admin added'})

if __name__ == "__main__":
    application.run(host="0.0.0.0", port=1337)
