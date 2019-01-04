package main


import (
	"regexp"
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
	"os"
)

const pastebinURL = "https://scrape.pastebin.com/api_scraping.php?limit=100"
const timeout = 60 * time.Second
const outfileFormat = "paste_%s_%s.txt"
const regexFile = "regex.json"

/* The pastebin json */
type Paste struct {
	ScrapeURL 	string		`json:"scrape_url"`
	FullURL 	string		`json:"full_url"`
	Date 		int64		`json:"date"`
	Key 		string		`json:"key"`
	Size 		int64		`json:"size"`
	Title 		string		`json:"title"`
	Syntax 		string		`json:"syntax"`
	User 		string		`json:"user"`
}

type RE struct {
	Name			string		`json:"name"`
	Regex			string		`json:"regex"`
	SecondaryRegex	[]string 	`json:"secondary_regex"`
	BlacklistRegex	[]string	`json:"blacklist_regex"`
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func getPastes() []Paste{
	var httpRes *http.Response
	var err error
	var pastes []Paste //= make([]Paste, 0)
	var data []byte
	//var pasteCount int
	//var counter int

	httpRes, err = http.Get(pastebinURL)
	
	if err != nil{
		return pastes
	}
	
	defer httpRes.Body.Close()
	data, err = ioutil.ReadAll(httpRes.Body)
	
	if err != nil {
		return pastes
	}
	json.Unmarshal(data, &pastes)
	/*pasteCount = len(pastes)
	for counter = 0; counter < pasteCount; counter++ {
		fmt.Printf("Title: %s\n Name: %s\n", pastes[counter].Title, pastes[counter].Key)
	}*/
	return pastes
}

func getBlacklist(pastes []Paste, blacklist []string) []string{
	var pasteCount int 
	var counter int

	pasteCount = len(pastes)
	for counter = 0; counter < pasteCount; counter++ {
		blacklist = append(blacklist, pastes[counter].Key)
	} 
	return blacklist
}

func checkBlacklist(blacklist []string, pastes []Paste) []Paste {
	var newPastes []Paste
	var counter int
	var pastesCount int

	pastesCount = len(pastes)

	for counter = 0; counter < pastesCount; counter++{
		if !stringInSlice(pastes[counter].Key, blacklist){
			newPastes = append(newPastes, pastes[counter])
		}
	}
	return newPastes
}

func readRegex(filename string) []RE {
	var err error
	var data []byte
	var re []RE

	data, err = ioutil.ReadFile(filename)
	if err != nil{
		fmt.Printf("[-] Can't read regex file.\n")
		os.Exit(0)
	}
	err = json.Unmarshal(data, &re)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}
	return re
}

func hasRegex(regexes []RE, data []byte) (bool, string){
	var re *regexp.Regexp;
	var regexesCount int
	var counter int
	var secLen int
	var secCount int
	var blacklistLen int
	var blacklistCount int
	
	regexesCount = len(regexes)

	for counter = 0; counter < regexesCount; counter++ {
		
		re = regexp.MustCompile(regexes[counter].Regex)
		if re.Find(data) != nil{
			blacklistLen = len(regexes[counter].BlacklistRegex)
			for blacklistCount = 0; blacklistCount < blacklistLen; blacklistCount++ {

				/* check for our blacklist */
				re = regexp.MustCompile(regexes[counter].BlacklistRegex[blacklistCount])
				if re.Find(data) != nil {
					return false, ""
				}
			}

			/* check for secondary regexes */
			secLen = len(regexes[counter].SecondaryRegex)
			if secLen > 0 {
				for secCount = 0; secCount < secLen; secCount++ {
					re = regexp.MustCompile(regexes[counter].SecondaryRegex[secCount])
					if re.Find(data) != nil{
						return true, regexes[counter].Name
					}
				}
			}else {
				return true, regexes[counter].Name
			}
		}
	}
	return false, ""
}

func savePaste(paste Paste, reName string,  data []byte){
	var filename string
	var err error

	filename = fmt.Sprintf(outfileFormat, strings.ToLower(reName), paste.Key)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("[-] Can't write to file '%s'.\n", filename)
		os.Exit(0)
	}else{
		fmt.Printf("[+] Dumped '%s'\n", filename)
	}
}

func checkPaste(paste Paste, regexes []RE){
	var httpRes *http.Response
	var err error
	var data []byte
	var reBool bool
	var reStr string


	httpRes, err = http.Get(paste.ScrapeURL)
	if err != nil{
		return
	}
	defer httpRes.Body.Close()
	data, err = ioutil.ReadAll(httpRes.Body)
	
	if err != nil {
		return
	}
	reBool, reStr = hasRegex(regexes, data)
	if reBool {
		savePaste(paste, reStr, data)
	}
}

func main(){
	
	var blacklist []string
	var regexes []RE
	var pastes []Paste
	var pastesCount int
	var counter int

	regexes = readRegex(regexFile)
	for true{
		pastes = getPastes()
		pastes = checkBlacklist(blacklist, pastes)
		blacklist = getBlacklist(pastes, blacklist)
		pastesCount = len(pastes)
		for counter = 0; counter < pastesCount; counter++{
			checkPaste(pastes[counter], regexes)
		}
		fmt.Printf("[ ] Sleeping for 1m.\n")
		time.Sleep(timeout)
	}
		
}
