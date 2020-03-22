package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var (
	sourceListOfUrls = []string{
		"https://www.casacochecurro.com/post-sitemap1.xml",
		"https://www.casacochecurro.com/post-sitemap10.xml",
		"https://www.casacochecurro.com/post-sitemap11.xml",
		"https://www.casacochecurro.com/post-sitemap12.xml",
		"https://www.casacochecurro.com/post-sitemap13.xml",
		"https://www.casacochecurro.com/post-sitemap14.xml",
		"https://www.casacochecurro.com/post-sitemap15.xml",
		"https://www.casacochecurro.com/post-sitemap16.xml",
	}

	amountFilesToGenerate    = 2000
	expectedAmountUrlsAmount = 1000
	destinationFolder        = "source_files/files/"
)

func main() {
	allUrls := goThroughAllUrls(sourceListOfUrls)
	allNewFiles := generateName(destinationFolder, "tst", amountFilesToGenerate)
	err := populateFilesWithData(allNewFiles, allUrls)
	check(err)
}

func goThroughAllUrls(srcList []string) []string {
	consolidatedSetOfURLs := make([]string, 0)

	for _, url := range srcList {
		result := parseLinkAndGetUrls(expectedAmountUrlsAmount, url)
		consolidatedSetOfURLs = append(consolidatedSetOfURLs, result...)
	}

	return consolidatedSetOfURLs
}

func populateFilesWithData(listGeneratedFiles, listUrls []string) error {
	if len(listGeneratedFiles)*4 > len(listUrls) {
		return errors.New("errExceededAmountOfUrls")
	}

	for i := 0; i < cap(listGeneratedFiles); i++ {
		file, err := os.OpenFile(listGeneratedFiles[i], os.O_WRONLY, 0666)
		check(err)

		_, err = file.WriteString(generateFileContent(listUrls))
		check(err)

		listUrls = listUrls[4:]
		err = file.Close()
		check(err)
	}

	return nil
}

func generateFileContent(givenSliceOfUrls []string) string {
	newContent := "{Block1}\n" +
		"aaa|bbb|ccc|" + givenSliceOfUrls[0] + "\n" +
		"{Block1}\n" +
		"{Block2}\n" +
		givenSliceOfUrls[1] + "|111|222|333\n" +
		givenSliceOfUrls[2] + "|444|555|666\n" +
		givenSliceOfUrls[3] + "|777|888|999\n" +
		"{Block2}\n"

	return newContent
}

func generateName(path, inName string, amount int) []string {
	var numerator int

	newListingFiles := make([]string, amount)

	for i := 0; i < amount; i++ {
		numerator = i
		newName := path + inName + strconv.Itoa(numerator) + ".lst"
		file, err := os.Create(newName)
		check(err)

		newListingFiles[i] = file.Name()
		err = file.Close()
		check(err)
	}

	return newListingFiles
}

func parseLinkAndGetUrls(expectedAmountOfUrls int, inURL string) []string {
	nc := http.Client{}
	response, err := nc.Get(inURL)
	check(err)

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	check(err)

	var cnt int

	var f func(*html.Node)

	sliceUrls := make([]string, expectedAmountOfUrls)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			if strings.HasPrefix(n.Data, "http") {
				if cnt >= len(sliceUrls) {
					sliceUrls = append(sliceUrls, n.Data)
				}

				sliceUrls[cnt] = n.Data
				cnt++
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return sliceUrls
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
