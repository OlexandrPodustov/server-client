//in order to match folder structure this file should be placed
//into folder source_files
package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	sourceListOfUrls = []string{
		"https://www.casacochecurro.com/post-sitemap1.xml",  //1266
		"https://www.casacochecurro.com/post-sitemap10.xml", //
		"https://www.casacochecurro.com/post-sitemap11.xml", //
		"https://www.casacochecurro.com/post-sitemap12.xml", //
		"https://www.casacochecurro.com/post-sitemap13.xml", //
		"https://www.casacochecurro.com/post-sitemap14.xml", //1577
		"https://www.casacochecurro.com/post-sitemap15.xml", //1426
		"https://www.casacochecurro.com/post-sitemap16.xml", //288
	}

	amountFilesToGenerate             int    = 2000
	expectedAmountUrlsInEachSourceUrl int    = 1000
	destinationFolder                 string = "source_files/files/"
	errExceededAmountOfUrls                  = errors.New("errExceededAmountOfUrls")
)

func main() {

	allUrls := goThroughAllUrls(sourceListOfUrls)
	allNewFiles := generateName(destinationFolder, "tst", amountFilesToGenerate)
	err := populateFilesWithData(allNewFiles, allUrls)
	check(err)
	//printSlice(allUrls)
	//printSliceWithoutGaps(allUrls)
	//println("total url", len(allUrls))
}

func goThroughAllUrls(srcList []string) []string {
	consolidatedSetOfURLs := make([]string, 0)

	for _, url := range srcList {
		resultSetFromOneUrl, err := parseLinkAndGetUrls(expectedAmountUrlsInEachSourceUrl, url)
		check(err)
		consolidatedSetOfURLs = append(consolidatedSetOfURLs, resultSetFromOneUrl...)

	}

	return consolidatedSetOfURLs
}

func populateFilesWithData(listGeneratedFiles, listUrls []string) error {
	//some debugging info
	//println("length of files", len(listGeneratedFiles), "capasity of files", cap(listGeneratedFiles))
	//println("length of urls", len(listUrls), "capasity of files", cap(listUrls))
	if len(listGeneratedFiles)*4 > len(listUrls) {
		return errExceededAmountOfUrls
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
	var newContent string
	newContent = "{Block1}\n" +
		"aaa|bbb|ccc|" + string(givenSliceOfUrls[0]) + "\n" +
		"{Block1}\n" +
		"{Block2}\n" +
		string(givenSliceOfUrls[1]) + "|111|222|333\n" +
		string(givenSliceOfUrls[2]) + "|444|555|666\n" +
		string(givenSliceOfUrls[3]) + "|777|888|999\n" +
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

func parseLinkAndGetUrls(expectedAmountOfUrls int, inURL string) ([]string, error) {
	nc := http.Client{}
	response, err := nc.Get(inURL)
	check(err)

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	check(err)

	var sliceUrls = make([]string, expectedAmountOfUrls)
	var cnt int

	var f func(*html.Node)
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

	return sliceUrls, nil
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func printSliceWithoutGaps(inputSlice []string) {
	var realAmount int
	for _, val := range inputSlice {
		if val != "" {
			realAmount++
			fmt.Printf("%+v\n", val)
		}
	}
	fmt.Printf("%+v\n", realAmount)
}

func printSlice(inputSlice []string) {
	for i, val := range inputSlice {
		fmt.Printf("%+v, %+v\n", i, val)
	}
}
