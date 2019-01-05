import json
from pymongo import MongoClient
from hurry.filesize import size
from telegram.ext import Updater, CommandHandler
from sys import exit

DB_CONFIG = "db.json"
BOT_CONFIG = "bot.json"


def read_config(filename):
	try:
		fh = open(filename, "r")
		data = fh.read()
		fh.close()
	except Exception as e:
		print("[-] Could not read config '{}'.".format(filename))
		exit(0)
	try:
		json_data = json.loads(data)
	except ValueError as e:
		print("[-] File '{}' does not contain valid JSON data.\n".format(filename))
		exit(0)
	return json_data
		
def get_mongo_stats():
	collections = []
	db_config = read_config(DB_CONFIG)
	
	try:
		mclient = MongoClient("mongodb://{}".format(db_config["host"]))
		mclient.server_info()
	except Exception as e:
		print("[-] Could not connect to host '{}'.\n".format(db_config["host"]))
		exit(0)
	db = mclient[db_config["dbname"]]
	db_size = size(db.command("dbstats")['storageSize'])
	text = "Database Size: {}\n".format(db_size)
	for col in db.list_collection_names():
		collections.append([col, db[col].estimated_document_count()])
	if collections:
		text = "\n{}Name\tCount\n".format(text)
		for col in collections:
			entry = "{}\t{}\n".format(col[0], col[1])
			text = "{}{}".format(text, entry)
	else:
		text = "{}\n[-] No collections\n".format(text)
	return text
	
def get_bot_token():
	config = read_config(BOT_CONFIG)
	return config["api_token"]

def get_bot_master():
	config = read_config(BOT_CONFIG)
	return config["master"]

def reply_pscrape_stats(bot, update):

	if update.message.from_user.username != get_bot_master():
		update.message.reply_text("You are not allowed to use this bot.\n")
	else:
		update.message.reply_text(get_mongo_stats())


if __name__ == "__main__":
	updater = Updater(get_bot_token())
	updater.dispatcher.add_handler(CommandHandler('pscrape_stats', reply_pscrape_stats))
	updater.start_polling()
	updater.idle()

