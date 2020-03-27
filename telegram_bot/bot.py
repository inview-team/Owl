import telebot

bot = telebot.TeleBot('')


@bot.message_handler(content_types=['text'])
def get_text_messages(message):
    if message.text == "Привет":
        bot.send_message(message.from_user.id, "Уху Уху")
    elif message.text == "/help":
        bot.send_message(message.from_user.id, "Привет, я сова")
    else:
        bot.send_message(message.from_user.id, "I don't know")

bot.polling(none_stop=True, interval=0)
