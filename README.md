# oauth2 with golang

``` 
git clone XXX
cd XXX
source .env
go run main.go
```

http://localhost:8080/

## Requirements

- ClientID
- ClientSecret
- AuthURL (is part of oauth2.github package, github.Endpoint)
- TokenURL (is part of oauth2.github package, github.Endpoint)
- RedirectURL (optional, http://localhost/oauth2/receive)
