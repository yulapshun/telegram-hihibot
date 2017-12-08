import json
import re
import sys
from telegram.ext import Updater, MessageHandler, Filters
import logging

config = {}
responses = {}
regex = None
updater = None

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)

def error(bot, update, error):
    """Log Errors caused by Updates."""
    logger.warning('Update "%s" caused error "%s"', update, error)

def get_response(msg):
    global responses
    match = regex.match(msg)
    if match:
        return responses[match.lastgroup].format(match.group())
    return ""

def process_msg(bot, update):
    msg = update.message.text

    response = get_response(msg)

    if response != "":
        update.message.reply_text(response, quote=True)

def run_cli():
    msg = ""
    while msg != "quit":
        msg = input("Say something: ")
        response = get_response(msg) or "Nothing"
        print(response)
    print("Bye")
    exit(0)

def run_polling():
    updater.start_polling()
    updater.idle()

def run_webhook():
    updater.start_webhook(listen="0.0.0.0",
                          port=config["listenPort"],
                          url_path=config["listenPath"])
    updater.bot.set_webhook(config["webhookUrl"])
    updater.idle()

def main():
    global config
    global responses
    global regex
    config_file = open("config.json", "r")
    config = json.loads(config_file.read())

    ruleset_file = open("ruleset.json", "r")
    ruleset = json.loads(ruleset_file.read())

    template = "(?P<{0}>{1})"
    arr = []
    count = 0
    for rule in ruleset:
        key = "g" + str(count)
        arr.append(template.format(key, rule["regex"]))
        responses[key] = rule["response"]
        count = count + 1
    regex = re.compile("|".join(arr))

    mode = "--bot"
    if len(sys.argv) == 2:
        mode = sys.argv[1]
    if mode not in ["--bot", "--interactive"]:
        sys.stderr.write("Invalid argument: {}\n".format(mode))
        exit(1)

    if mode == "--interactive":
        run_cli()

    global updater
    updater = Updater(config["token"])
    dp = updater.dispatcher

    dp.add_handler(MessageHandler(Filters.text, process_msg))
    dp.add_error_handler(error)

    if config["useWebhook"]:
        run_webhook()
    else:
        run_polling()


if __name__ == '__main__':
    main()
