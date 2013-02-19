package main
import "flag"
import "github.com/c9s/go-appcast/appcast"
import "github.com/c9s/go-rss/rss"
import "os"
import "time"


/*
$ appcast -pubDate "Date string" \
		-description " description" \
		-url "http://host.com/app/App.zip" \
		-version "109" \
		-dsaSignature "BAFJW4B6B1K1JyW30nbkBwainOzrN6EQuAh" \
		-title "Release 1.4" \
		path/to/app.zip
*/

type path string

func PathExists(p string) (bool, error) {
    _, err := os.Stat(p)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func main() {

	// required options
	url         := *flag.String("url","","file url")
	description := *flag.String("description","","description")
	title       := *flag.String("title","","title")
	version     := *flag.String("version"            , "", "sparkle:version (CFBundleVersion: build number)")
	file        := *flag.String("file","","application file")

	// optional options
	pubDate            := *flag.String("pubDate","","pubDate")
	versionShortString := *flag.String("versionShortString" , "", "sparkle:versionShortString (Release Version)")
	dsaSignature       := *flag.String("dsaSignature"       , "", "sparkle:dsaSignature")


	if url == "" {
		panic("-url is required.")
	}
	if description == "" {
		panic("-description is required.")
	}
	if title == "" {
		panic("-title is required.")
	}
	if version == "" {
		panic("-version is required.")
	}


	if pubDate == "" {
		pubDate = time.Now().Format(time.RFC822Z)
	}


	flag.Parse()
	appcastFile := flag.Arg(0)

	if file == "" {
		panic("-file is required.")
	}
	if appcastFile == "" {
		panic("file argument is required to create an enclosure.")
	}

	if ok, _ := PathExists(appcastFile) ; ! ok {
		panic( appcastFile + " does not exist: " )
	}

	if ok, _ := PathExists(file) ; ! ok {
		panic( file + " does not exist: " )
	}


	item := appcast.Item{}
	item.PubDate = rss.Date(pubDate)
	// item.PubDate = rss.Date(pubDate)
	en , err := appcast.CreateItemEnclosureFromFile(file)
	if err != nil {
		panic(err)
	}

	en.SparkleVersion = version
	if versionShortString != "" {
		en.SparkleVersionShortString = versionShortString
	}
	if dsaSignature != "" {
		en.SparkleDSASignature = dsaSignature
	}

	item.Title = title
	item.Description = description
	item.Enclosure = *en
	apprss, err := appcast.ReadFile(appcastFile)

	appcast.WriteFile(appcastFile,apprss)

	_ = err
}