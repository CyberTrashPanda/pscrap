
# pscrap
  
### Description
pscrap is a Go Pastebin scraping tool, utilizing the Pastebin scrapping API.

The pastes are stored in a MongoDB database and it includes a quite simple python3 telegram bot to monitor database usage from anywhere.

It uses regular expressions to validate and exclude pastes, moreover, the structure of the regex configuration file is quite simple.

| Property Name | Description |
| ------------- | ------------ |
| name | Used as a filename to save the pastes. |
| regex | The main regular expression to match. |
| secondary_regex | An array of secondary regular expressions to match. |
| blacklist_regex | An array of regular expressions to avoid.


The telegram bot will only reply to the user specified in the **bot.json** file.



### TODO
* ~~Store pastes in a database~~.
* ~~Create a small Telegram bot to provide information about stored pastes.~~


### Tips
1) Configure both **db.json** files.
2) Configure **bot.json** with your *API_TOKEN* and *username*.



### Screenshots
Main Pastebin scrapping tool.

![screenshot_1](https://i.imgur.com/XtNFnjx.png)

Telegram Bot (Owner)

![screenshot_2](https://i.imgur.com/WR1FB3i.png)

Telegram Bot (When you are not the bot owner)

![screenshot_3](https://i.imgur.com/ZebQ1yf.png)