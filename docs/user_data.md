# rueder user data design

## improvements

* maybe `feed_seq` should be used for `readArticles` in `feedstate`
* we need a way to sync folder changes between connected clients of the same user
* we need a way to sync feedstates and labels between connected clients of the
  same user

## current state

### frontend

Frontend stores its data in local storage.

#### feedstate

Contains a list of read articles per feed. This is used to calculate the
*unread count* bubble and for the *unread article* indicators.
The data is compacted when you read all articles or press "Mark All Read",
the list of `readArticles` is cleared and `readAllUntil` is saved
to indicate that all articles have been read up until the given sequence number.

**Note:** Feeds and articles are identified by UUIDs (like most of our objects).
Additionally there are 2 different *sequence numbers* for articles: `seq` and
`feed_seq`.

`seq` is a running ID that behaves pretty much like an auto-increment field in
MySQL for the article table, i.e. `seq=3` says "this is the third article in our
 database" (in the order we added the articles).

`feed_seq` is a running ID *per feed*, i.e. `feed_seq=3` says "this is the third
article of the given feed".

example:

```json
{
    "feeds": {
        "91751be3-8b9c-4ccd-a02c-652df9cfda04": {
            "totalArticles": 58,
            "readArticles": [19, 18, 17, 20, 15, 21, 26, 28, 38, 37, 36, 53, 51, 58, 57, 56],
            "readAllUntil": 0
        },
        "5279c064-1c4a-4499-9862-bc5ef891dd0d": {
            "totalArticles": 551,
            "readArticles": [
                164, 162, 163, 161, 160, 180, 181, 179, 177, 176, 219, 233, 251, 250, 403, 404, 497, 495, 515, 514, 516,
                513, 511
            ],
            "readAllUntil": 0
        },
        "86808460-2440-42c0-9518-4c18e35e81f6": { "totalArticles": 8, "readArticles": [], "readAllUntil": 8 },
        "1a235fe7-234f-45e7-9dd8-0f87523a6ad0": {
            "totalArticles": 226,
            "readArticles": [223, 224, 225],
            "readAllUntil": 0
        },
        "4371eb27-4e91-43f1-94b8-1ab40e33782e": {
            "totalArticles": 31,
            "readArticles": [24, 23, 4, 5, 6, 14, 18, 20, 21, 16, 15, 17, 22, 19, 25, 10, 9, 7],
            "readAllUntil": 3
        },
        "a6961a5d-31fa-4460-8c15-0d5ab6a584c7": { "totalArticles": 206, "readArticles": [92, 91], "readAllUntil": 0 },
        "763476bb-214b-4317-bb2f-eb0c6c999c75": { "totalArticles": 25, "readArticles": [25], "readAllUntil": 0 },
        "8a6607f5-413a-43d1-ae57-73219ac7e3e3": { "totalArticles": 58, "readArticles": [], "readAllUntil": 0 },
        "e5a284fe-717a-451d-a569-b1e7a4ea65c3": { "totalArticles": 554, "readArticles": [213], "readAllUntil": 0 },
        "49d6c093-bdf8-401a-a937-b81e050824dd": { "totalArticles": 644, "readArticles": [], "readAllUntil": 0 },
        "c84a1674-ee88-4382-8901-6d418f2a14f2": {
            "totalArticles": 25,
            "readArticles": [25, 24, 22, 23],
            "readAllUntil": 0
        },
        "b5b7bff8-6599-4f95-b3c4-9a01ca91425e": { "totalArticles": 19, "readArticles": [], "readAllUntil": 0 },
        "042701f2-c5fc-4156-beba-7a46ba7f2e69": { "totalArticles": 14, "readArticles": [], "readAllUntil": 10 }
    }
}
```

#### labels

Contains all labelled articles (with full content for a kind of local mirror).
Also label appearance.

example:

```json
{
    "articles": {
        "b6a0c0b4-1516-497f-9b7c-d1b1d42b2793": {
            "id": "b6a0c0b4-1516-497f-9b7c-d1b1d42b2793",
            "seq": -1,
            "feed_seq": -1,
            "title": "Oh, ist Krieg? Ohne die SPD? Nee, das geht nicht. Da ...",
            "time": "2022-05-07T10:39:21Z",
            "feed_title": "Fefes Blog",
            "feed_icon": null,
            "teaser": "Oh, ist Krieg? Ohne die SPD? Nee, das geht nicht. Da muss die SPD schnell aufspringen![...]"
        },
        "e7cfb0f9-1db2-42b6-b627-3ecfde269c52": {
            "id": "e7cfb0f9-1db2-42b6-b627-3ecfde269c52",
            "seq": -1,
            "feed_seq": -1,
            "title": "ZTA doesn't solve all problems, but partial implementations solve fewer",
            "time": "2022-03-31T23:06:44Z",
            "feed_title": "Matthew Garrett",
            "feed_icon": null,
            "teaser": "Traditional network access controls work by assuming that something is trustworthy based on some other factor[...]"
        }
    },
    "labels": {
        "important": {
            "name": "important",
            "color": "#d0f806",
            "articleIDs": ["b6a0c0b4-1516-497f-9b7c-d1b1d42b2793", "e7cfb0f9-1db2-42b6-b627-3ecfde269c52"]
        },
        "favorite": { "name": "favorite", "color": "#F59E0B", "articleIDs": [] },
        "read later": { "name": "read later", "color": "#10B981", "articleIDs": [] },
        "sfsdf": { "name": "sfsdf", "color": "#cd3535", "articleIDs": ["b6a0c0b4-1516-497f-9b7c-d1b1d42b2793"] },
        "sf": { "name": "sf", "color": "#aaa656", "articleIDs": ["b6a0c0b4-1516-497f-9b7c-d1b1d42b2793"] }
    }
}
```

#### session

Contains the JWT which is used by all authenticated API requests.

example:

```json
{
    "loggedIn": true,
    "jwtToken": "eyJhbGciOiJIUzUxMiIsInR5c[...]N6Y3MPag8cBw"
}
```

### backend

The user's folders (including feeds) are stored in backend.
Manipulations (adding/removing feeds/folders) are done in frontend which then
just sends the updated full JSON to the backend.
The backend does some validation (enforcing feed/folder count limits,
sanity/completeness checks) and saves the data.
The backend also extracts subscribed feeds for its global list of feeds of all
users that need to be fetched.

example:

```json
[
    {
        "id": "78f8cef4-3a9e-4921-9f9e-e671ceae87e1",
        "title": "Default",
        "feeds": [
            {
                "id": "8a6607f5-413a-43d1-ae57-73219ac7e3e3",
                "title": "Fefes Blog",
                "url": "https://blog.fefe.de/rss.xml?html",
                "site_url": "https://blog.fefe.de/",
                "article_count": 58
            },
            {
                "id": "e5a284fe-717a-451d-a569-b1e7a4ea65c3",
                "title": "heise online News",
                "url": "https://www.heise.de/rss/heise-atom.xml",
                "site_url": "https://www.heise.de/",
                "article_count": 554
            },
            {
                "id": "49d6c093-bdf8-401a-a937-b81e050824dd",
                "title": "The Guardian",
                "icon": "https://assets.guim.co.uk/images/guardian-logo-rss.c45beb1bafa34b347ac333af2e6fe23f.png",
                "url": "https://www.theguardian.com/international/rss",
                "site_url": "https://www.theguardian.com/international",
                "article_count": 644
            },
            {
                "id": "c84a1674-ee88-4382-8901-6d418f2a14f2",
                "title": "Matthew Garrett",
                "url": "https://mjg59.dreamwidth.org/data/rss",
                "site_url": "https://mjg59.dreamwidth.org/",
                "article_count": 25
            }
        ]
    },
    {
        "id": "4a2817ec-d50b-4bbb-a640-1d1d7e942fb3",
        "title": "Unnamed Folder 1"
    },
    {
        "id": "0fdf7374-d057-47a5-96ed-30372e774940",
        "title": "Unnamed Folder 25"
    },
    {
        "id": "6577e9c9-0bab-4a30-86a5-88ed1d73c5ff",
        "title": "Unnamed Folder 26"
    },
    {
        "id": "1cf390c0-0985-430e-9e50-999b51980a36",
        "title": "Unnamed Folder 27",
        "feeds": [
            {
                "id": "042701f2-c5fc-4156-beba-7a46ba7f2e69",
                "title": "Wiktionary Word of the day",
                "url": "https://en.wiktionary.org/w/api.php?action=featuredfeed\u0026feed=wotd\u0026feedformat=atom",
                "site_url": "https://en.wiktionary.org/wiki/Wiktionary:Main_Page",
                "article_count": 14
            }
        ]
    }
]
```
