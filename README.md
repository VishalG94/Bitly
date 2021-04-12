# Cloud Project - Bitly

### Application
Bitly is a URL shortining Appliaction. 

The main components of the appliation are as follows.

Control Panel: The Control Panel receives request for URL shortining. Control Panel creates the short link and place the create request in the Message Queue.

Trend Server: The Trend Server that kept track of “trending links” statistics.

Link Redirect Server: The link redirect server redirects the short url to the expanded url.

Main Database: For this project MySQL DB is used as the main relational database. This database server stores all the shortlink information.

NoSQL Cache Database: This Database is an implementation of the 5-Node AP Key-Value Pair DB.

### Architecture Diagram

<img width="900" alt="ClassDiagram" src="https://github.com/VishalG94/Bitly/blob/main/Artifacts/Architecture_Diagram/Final_Architecture_Diagram.png">

### Youtube link

https://www.youtube.com/watch?v=8cWiiPTZz3M&ab_channel=VishalGadapa

### Heroku Deployment

<img width="900" alt="ClassDiagram" src="https://github.com/VishalG94/Bitly/blob/main/Artifacts/Screenshots/Heroku.png">


### FrontEnd

<img width="900" alt="ClassDiagram" src="https://github.com/VishalG94/Bitly/blob/main/Artifacts/Screenshots/Frontend.png">


### Sampe API Request Responses 

### Ping Check

#### Request 
curl --location --request GET 'http://34.221.223.237:8000/ts/ping'

#### Response 
{
    "Test": "Trend Server API version 1.0 alive!"
}

#### Request 
curl --location --request GET 'http://34.221.223.237:8000/cp/ping'

#### Response 
{
    "Test": "Control Panel API version 1.0 alive!"
}

#### Request 
curl --location --request GET 'http://34.221.223.237:8000/lr/ping'

#### Response 
{
    "Test": "Link Redirect API version 1.0 alive!"
}

#### Request 
curl --location --request GET 'http://34.221.223.237:8000/cs/ping'

#### Response 
{
    "Test": "Core Server API version 1.0 alive!"
}

#### Request
curl 10.0.1.199:3002/ping

#### Response
{
  "Test": "Base Count API version 1.0 alive!"
}

### Control Panel - Create Short Link

#### Request
curl --location --request POST '34.221.223.237:8000/cp/shortlink' \
--header 'Content-Type: application/json' \
--data-raw '{
    "URL":"https://sjsu.instructure.com/courses/1374205/assignments/5537602"
}
'

#### Response
{
    "Id": "54935348-2686-4fd8-8aa1-81b056d15e78",
    "ShortLink": "http://34.221.223.237:8000/lr/264Z",
    "URL": "https://sjsu.instructure.com/courses/1374205/assignments/5537602",
    "Count": 0
}

### Link Redirect

#### Request
curl --location --request GET 'http://34.221.223.237:8000/lr/lr/1g3h' \
--header 'apikey: secretkey'

#### Response
<img width="900" alt="ClassDiagram" src="https://github.com/VishalG94/Bitly/blob/main/Artifacts/Screenshots/Screen%20Shot%202020-12-09%20at%2010.50.19%20PM.png">


Trend Server - Trending Links

#### Request:
curl --location --request GET 'http://34.221.223.237:8000/ts/shortlinktrend'

#### Response:
[
    {
        "Id": "feffee73-5ce7-431b-9ba2-a89bbcb15f06",
        "ShortLink": "http://34.221.223.237:8000/lr/q1u",
        "URL": "https://www.netflix.com/",
        "Count": 8
    },
    {
        "Id": "67c93e49-c164-463e-a540-e442a2153baa",
        "ShortLink": "http://34.221.223.237:8000/lr/264X",
        "URL": "https://www.wired.com/category/culture/",
        "Count": 6
    },
    {
        "Id": "600ab828-0a8b-4a32-9d97-fa90e41a3277",
        "ShortLink": "http://34.221.223.237:8000/lr/1g3h",
        "URL": "https://www.hackerrank.com/skills-verification/problem_solving_basic",
        "Count": 6
    },
    {
        "Id": "7908e8dd-35cb-4328-a59f-b1999c70b1a2",
        "ShortLink": "http://34.221.223.237:8000/lr/Q0w",
        "URL": "https://www.instacart.com/store/walmart/storefront",
        "Count": 4
    },
    {
        "Id": "5f065ea6-db58-4ef3-a90b-308e019e5fa3",
        "ShortLink": "http://34.221.223.237:8000/lr/1g3j",
        "URL": "https://www.hackerrank.com/contests",
        "Count": 3
    },
    {
        "Id": "57811cc4-dd06-475c-810d-a90e607466ba",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2o",
        "URL": "https://signup.heroku.com/account",
        "Count": 3
    },
    {
        "Id": "4cf13d05-9877-42da-ac36-581fb4c8f0f1",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2k",
        "URL": "https://www.tinder.com/",
        "Count": 3
    },
    {
        "Id": "d784fc44-c914-47a7-9411-7a19c3efbaba",
        "ShortLink": "http://34.221.223.237:8000/lr/Q0v",
        "URL": "https://www.instacart.com/",
        "Count": 3
    },
    {
        "Id": "5212168c-1b75-4485-a6a6-44c8df1d6b3a",
        "ShortLink": "http://34.221.223.237:8000/lr/1g3f",
        "URL": "https://www.wired.com/category/business/",
        "Count": 2
    },
    {
        "Id": "2c4a86a6-b0aa-405e-97a0-6ba78b9468a8",
        "ShortLink": "http://34.221.223.237:8000/lr/q1t",
        "URL": "https://www.netflix.com/browse",
        "Count": 2
    },
    {
        "Id": "9961e3d3-4bf5-4e7d-b0b3-9a4123f98735",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2l",
        "URL": "http://www2.cs.uic.edu/~jbell/CourseNotes/OperatingSystems/",
        "Count": 2
    },
    {
        "Id": "ab70b60d-ead0-4839-8c63-86efc72a9585",
        "ShortLink": "http://34.221.223.237:8000/lr/264Y",
        "URL": "https://www.wired.com/story/one-mans-search-for-dna-data-that-could-save-his-life/",
        "Count": 2
    },
    {
        "Id": "b71b0e68-de5b-480b-9c90-71beaff4647b",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2m",
        "URL": "https://www.youtube.com/",
        "Count": 1
    },
    {
        "Id": "870a7ab3-1d82-47d5-a374-8dbcae78a0d6",
        "ShortLink": "http://34.221.223.237:8000/lr/1g3e",
        "URL": "https://www.wired.com/category/gear/",
        "Count": 1
    },
    {
        "Id": "af92e8cd-719e-4861-9a84-1a6112a06f42",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2q",
        "URL": "https://www.amazon.com/gcx/Home-Holiday-Guide/gfhz/events/ref=cg_GH20T1CG_1a1_w?categoryId=HHG20-Hub&pf_rd_m=ATVPDKIKX0DER&pf_rd_s=desktop-top-slot-1&pf_rd_r=HDCBJCX4MYH9T7SYGC56&pf_rd_t=0&pf_rd_p=2d682585-ad94-4073-9583-b37ed7b987c7&pf_rd_i=gf-landing",
        "Count": 1
    },
    {
        "Id": "f335556b-d8b1-4bda-a712-2f7b445adff8",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2p",
        "URL": "https://www.wired.com/account/sign-in",
        "Count": 1
    },
    {
        "Id": "2884a9d2-abb4-4c06-964f-ee3e552b13a7",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2j",
        "URL": "https://www.facebook.com/",
        "Count": 1
    },
    {
        "Id": "72a5ebe3-1c00-4a21-ab4c-99b65e81f694",
        "ShortLink": "http://34.221.223.237:8000/lr/Q0x",
        "URL": "https://www.instacart.com/store/walmart/",
        "Count": 1
    },
    {
        "Id": "09c659e4-ed8c-4c81-ae05-9ff44d5afd47",
        "ShortLink": "http://34.221.223.237:8000/lr/q1v",
        "URL": "https://www.google.com/",
        "Count": 1
    },
    {
        "Id": "6771892e-87c7-42bb-8b96-f5d1055a46ca",
        "ShortLink": "http://34.221.223.237:8000/lr/1G2n",
        "URL": "https://docs.konghq.com/2.0.x/db-less-and-declarative-config/",
        "Count": 1
    }
]

### Base Count Server
 
### Request
curl 10.0.1.199:3002/basecount

### Response
{
  "Min": 600001,
  "Max": 700000
}

