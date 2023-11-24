# easyMail
Super easy mail server written in go

## Installation
```sh
go get github.com/darkard2003/easyMail@latest
```

## Usage
Using easyMail is very simple. You will need a email, and a application specific password. In gmail you can generate this password under security and then two factor authentication.

### Importing
```go
import "github.com/darkard2003/easyMail"
```


### Creating a mail server
```go 
server := easyMail.NewMailServer("email", "password", easyMail.GMAIL)
```

### Creating a mail
```go
mail := easyMail.NewMail(
    "From",
    []string{"To1", "To2", "To3"},
    "Subject",
    "Body",
    false, // Is HTML
)
```

### Attaching a file
```go
err := mail.AddAttachment("path/to/file")
```

### Sending a mail
```go
err := server.SendMail(mail)
```