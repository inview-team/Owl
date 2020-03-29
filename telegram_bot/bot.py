import telebot
import os
from dotenv import find_dotenv,load_dotenv
from service_functions import load_setting
load_dotenv(find_dotenv())


token=os.environ['TOKEN']
bot = telebot.TeleBot(token)

@bot.message_handler(content_types=['text'])
def get_text_messages(message):
    if message.text == "Привет":
        bot.send_message(message.from_user.id, "Уху Уху")
    elif message.text == "/help":
        bot.send_message(message.from_user.id, "Привет, я сова")
    else:
        bot.send_message(message.from_user.id, "I don't know")

@bot.message_handler(commands=['settings'])
def get_settings(message):
    result=load_setting()
    settings = result['settings']
    keyboard = telebot.types.InlineKeyboardMarkup()
    for i in range(len(settings)):
        keyboard.row(
            telebot.types.InlineKeyboardButton(text=settings[i]['metric'], callback_data='get_info')
        )
    bot.send_message(
        message.chat.id,
        'Get info about Settings',
        reply_markup=keyboard
    )

@bot.callback_query_handler(func=lambda call:True)
def iq_callback(query):
    data = query.data
    print(data)

bot.polling(none_stop=True, interval=0)


