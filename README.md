
# pscrap
  
### Description
pscrap is just a small pastebin scraping tool utilizing the pastebin scrapping API.

It uses regular expressions to validate or exclude pastes.


The structure of the config file is quite simple.

| Property Name | Description |
| ------------- | ------------ |
| name | Used as a filename to save the pastes. |
| regex | The main regular expression to match. |
| secondary_regex | An array of secondary regular expressions to match. |
| blacklist_regex | An array of regular expressions to avoid.



### TODO
* Store pastes in a database.
* Create a small Telegram bot to provide information about stored pastes.


### Screenshot
![screenshot_1](https://i.imgur.com/VjZSxH8.png)