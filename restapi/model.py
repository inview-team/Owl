from flask_sqlalchemy import SQLAlchemy
db = SQLAlchemy()

class Alarms(db.Model):
    __tablename__ = "alarms"

    id = db.Column('id', db.Integer, primary_key=True)
    time = db.Column('time', db.String)
    info = db.Column('info', db.String)


    def __init__(self, id, time, info):
        self.id=id
        self.time=time
        self.info=info

    def serialize(self):
        print(self.__dict__)
        return {
            'id': self.id,
            'time': self.time,
            'info': self.info,
        }

class Logs(db.Model):
    __tablename__ = "logs"

    id = db.Column('id', db.Integer, primary_key=True)
    time = db.Column('time', db.String)
    info = db.Column('info', db.String)

    def __init__(self, id, time, info):
        self.id = id
        self.time = time
        self.info = info

    def serialize(self):
        print(self.__dict__)
        return {
            'id': self.id,
            'time': self.time,
            'info': self.info,
        }

class Settings(db.Model):
    __tablename__ = "settings"
    id=db.Column('id', db.Integer, primary_key=True)
    name=db.Column('name', db.String)
    from_bord=db.Column('from_bord', db.Integer)
    to_bord = db.Column('to_bord', db.Integer)


    def __init__(self, id, name, from_bord, to_bord):
        self.id=id
        self.name = name
        self.from_bord = from_bord
        self.to_bord = to_bord

    def serialize(self):
        return {
            'id': self.id,
            'metric': self.name,
            'from': self.from_bord,
            'to': self.to_bord
        }

class Telegram(db.Model):
    __tablename__ = "telegram"
    id=db.Column('id', db.Integer, primary_key=True)
    chat_id = db.Column('chat_id', db.String)

    def __init__(self, id, chat_id):
        self.id = id
        self.chat_id = chat_id

    def serialize(self):
        return {
            'id': self.id,
            'chat_id': self.chat_id
        }

def init_db():
    db.create_all()

if __name__ == '__main__':
    init_db()
