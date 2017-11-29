import json
import re
import sys
from telegram.ext import Updater, MessageHandler, Filters
import logging

config = {}
ruleset = {}
updater = None

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)

def error(bot, update, error):
    """Log Errors caused by Updates."""
    logger.warning('Update "%s" caused error "%s"', update, error)

def compare_match(rule, msg):
    for pattern in rule["patterns"]:
        if msg == pattern:
            return rule["response"].format(pattern)
    return ""

def compare_contain(rule, msg):
    for pattern in rule["patterns"]:
        if pattern in msg:
            return rule["response"].format(pattern)
    return ""

def compare_regex(rule, msg):
    for pattern in rule["patterns"]:
        match = re.search(pattern, msg)
        if match:
            return rule["response"].format(match.group(0))
    return ""

def get_response(msg):
    response = ""

    for rule in ruleset:
        if rule["type"] == "match":
            response = compare_match(rule, msg)
        elif rule["type"] == "contain":
            response = compare_contain(rule, msg)
        elif rule["type"] == "regex":
            response = compare_regex(rule, msg)

        if response != "":
            break
    return response

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
    global ruleset
    config_file = open("config.json", "r")
    config = json.loads(config_file.read())

    ruleset_file = open("ruleset.json", "r")
    ruleset = json.loads(ruleset_file.read())

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
