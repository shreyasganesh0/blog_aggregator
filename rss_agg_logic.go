package main

import(
    "encoding/xml"
    "context"
    "net/http"
    "io"
)

type RSSFeed struct{
    Channel struct{
        Title string `xml:"title"`
        Link string `xml:"link"`
        Description string `xml:"description"`
        Items []RSSItem `xml:"item"`
    } `xml:"channel"`
}

type RSSItem struct{
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, fetchUrl string) (*RSSFeed, error){

    var rss_feed RSSFeed;
    req, err := http.NewRequestWithContext(ctx, "GET", fetchUrl, nil);

    if err != nil {
        return nil, err;
    }

    req.Header.Set("User-Agent", "gator");

    client := &http.Client{};

    resp, err1 := client.Do(req);

    if err1 != nil{
        return nil, err1;
    }

    byte_resp, err2 := io.ReadAll(resp.Body);

    if err2 != nil{
        return nil, err2;
    }

    if err := xml.Unmarshal(byte_resp, &rss_feed); err != nil{
        return &rss_feed, err;
    }

    return &rss_feed, nil;
}



        
    

