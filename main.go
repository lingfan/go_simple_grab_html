package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "strconv"
    "strings"
)

func getListDetail(m string, s []byte) []string {
    r := bytes.NewReader(s)
    doc, err := goquery.NewDocumentFromReader(r)
    checkError(err)

    arr := []string{}

    dir := "./detail/" + m
    err0 := os.MkdirAll(dir, 0777)
    checkError(err0)

    doc.Find("div.cont dl").Each(func(i int, s *goquery.Selection) {
        title := s.Find("dt a").Text()
        link, _ := s.Find("dt a").Attr("href")

        url := "http://www.safe10000.com" + link

        file := dir + "/" + path.Base(url)
        //desc := s.Find("dd").Text()
        arr = append(arr, url)
        c := getContent(url, file)

        fmt.Printf("Review %d: %s - %s - %s\n", i, url, title, file)
        parseDetailHtml(m, c, path.Base(url))
    })
    return arr
}

func checkError(e error) {
    if e != nil {
        log.Println(e)
        //panic(e)
    }
}

func loadHtml(m string, maxPage int) {

    if maxPage >= 1 {
        dir := "./html/" + m
        err0 := os.MkdirAll(dir, 0777)
        checkError(err0)
        for i := 1; i <= maxPage; i++ {
            p := strconv.Itoa(i)
            file := dir + "/" + p + ".html"
            url := "http://www.safe10000.com/news/" + m + "-" + p
            listContentByte := getContent(url, file)
            urls := getListDetail(m, listContentByte)
            fmt.Println("%v", urls)
        }

    }

}

func getContent(url string, file string) []byte {
    contentByte, err := ioutil.ReadFile(file)
    if err == nil {
        return contentByte
    }
    fmt.Printf("url: %s \n", url)
    resp, err := http.Get(url)
    checkError(err)
    defer resp.Body.Close()

    htmlData, err := ioutil.ReadAll(resp.Body)
    fmt.Printf("size: (%d)\n", len(htmlData))
    checkError(err)

    contentByte1 := []byte(htmlData)
    err1 := ioutil.WriteFile(file, contentByte1, 0777)
    checkError(err1)
    return contentByte1
}

func saveImg() {

}

func parseDetailHtml(m string, s []byte, name string) map[string]string {
    ret := map[string]string{}
    dir := "./data/" + m
    err0 := os.MkdirAll(dir, 0777)
    checkError(err0)

    file := dir + "/" + name
    _, err := ioutil.ReadFile(file)
    if err == nil {
        return ret
    }

    r := bytes.NewReader(s)
    doc, err := goquery.NewDocumentFromReader(r)
    checkError(err)
    title := doc.Find("div#detail h1").Text()
    dateStr := doc.Find("div#detail div.related").Text()
    tmpArr := strings.Split(dateStr, "来源：")
    date := tmpArr[0]
    from := tmpArr[1]
    desc := doc.Find("div#detail div.related2").Text()
    content, err := doc.Find("div#detail div.text").Html()

    ret["title"] = title
    ret["date"] = strings.TrimSpace(date)
    ret["from"] = strings.TrimSpace(from)
    ret["desc"] = strings.TrimSpace(desc)
    ret["content"] = strings.TrimSpace(strings.Replace(content, "<p style=\"padding-left:35px;\"><img width=\"500\" heigth=\"273\" src=\"/skin/safe/image/wxtzx.jpg\"/></p>", "", -1))
    fmt.Println("%v", ret)
    b, err := json.Marshal(ret)
    checkError(err)
    fmt.Println("%v", string(b))

    err1 := ioutil.WriteFile(file, b, 0777)
    checkError(err1)
    return ret
}

func main() {
    mod := map[string]int{
        //"afyj-156":  443,
        "cpcp-1079": 74,
        "cphq-1080": 74,
        "afht-159":  45,
        "jsqy-423":  48,
        "yyal-424":  126,
        "jjfa-1398": 27,
    }
    //m := strings.Join(os.Args[1:2], "+")
    //page := strings.Join(os.Args[2:3], "+")
    //fmt.Printf("%v",os.Args)
    for m, maxPage := range mod {
        fmt.Printf("%s %d\n", m, maxPage)
        loadHtml(m, maxPage)
    }

    //loadUrl()
}
