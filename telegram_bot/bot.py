import telebot
import os
from dotenv import find_dotenv,load_dotenv
from service_functions import load_setting, get_one_metric, add_new_admin
load_dotenv(find_dotenv())


token=os.environ['TOKEN']
bot = telebot.TeleBot(token)


@bot.message_handler(commands=['settings'])
def settings(message):
    result = load_setting()
    settings = result['settings']
    keyboard = telebot.types.InlineKeyboardMarkup()
    for i in range(len(settings)):
        keyboard.row(
            telebot.types.InlineKeyboardButton(text=settings[i]['metric'], callback_data=settings[i]['metric'])
        )
    bot.send_message(
        message.chat.id,
        'Get info about Settings',
        reply_markup=keyboard
    )

@bot.message_handler(commands=['add_admin'])
def add_admin(message):
    add_new_admin(message.chat.id)

@bot.callback_query_handler(func=lambda call:True)
def iq_callback(query):
    data = query.data
    print(data)
    result=get_one_metric(data)
    answer = '<b>' + data + ':</b>\n\n' + \
             'From: ' + str(result['from']) + '\n' + \
             'To: ' + str(result['to'])
    bot.send_message(
        query.message.chat.id,
        answer,
        parse_mode='HTML'
    )


bot.polling(none_stop=True, interval=0)


