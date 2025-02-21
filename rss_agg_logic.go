package main

import(
    "encoding/xml"
    "context"
    "net/http"
    "io"
    "fmt"
    "time"
    "database/sql"
    "github.com/google/uuid"
    "github.com/shreyasganesh0/blog_aggregator/database" 

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
    Description sql.NullString `xml:"description"`
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

func scrapeFeeds(s *state)  error{
    ctx := context.Background();
    url, err := s.queries.GetNextFeedToFetch(ctx);
    if  err != nil{
        return err;
    }
    sql_time := sql.NullTime{
        Time: time.Now(),
        Valid: true,
    }
    query_args := database.MarkFeedFetchedParams{
        LastFetchedAt: sql_time,
        UpdatedAt: time.Now(),
        Url: url,
    };
    err1 := s.queries.MarkFeedFetched(ctx, query_args);
    if err1 != nil{
        return err1;
    }
    var rss_feed_p *RSSFeed;
    rss_feed_p, err2 := fetchFeed(ctx, url);
    if err2 != nil{
        return err2;
    }
    
    feed_id, err3 := s.queries.FetchFeedId(ctx, url);
    if err3 != nil{
        return err3;
    }
    
    for _, rss_item := range rss_feed_p.Channel.Items {
        layout := time.RFC1123Z; 
        pub_date, timerr := time.Parse(layout, rss_item.PubDate);
        if timerr != nil{
            fmt.Printf("time error : %v", timerr);
            continue;
        }
       post_args := database.CreatePostParams{
            ID: uuid.New(),        
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Title: rss_item.Title,
            Url: rss_item.Link,
            Description: sql.NullString(rss_item.Description), 
            PublishedAt: pub_date,
            FeedID: feed_id, 
        };

        _, err := s.queries.CreatePost(ctx, post_args);
        if err != nil{
            return err;
        }
    }
    return nil;
}

