# countdown-generator
Public Go Library that allows you to create and serve countdown-gif images which can be used to embed into e-mails. It generates up-to-date new gif for every request so when a user opens the email, the correct remaining time is shown.

![Image](https://i.ibb.co/3dqXkFm/Screen-Shot-2018-12-05-at-14-35-50.png "")

# Prerequisites

You need to have go installed in your environment.

# Installation,

Go to your go directory ($GOPATH),if you don't have a github.com directory, create github.com folder inside your go directory

```
cd $GOPATH/github.com
```
 
clone this project
```
git clone https://github.com/otsimo/countdown-generator
```

You need freetype to run countdown generator
```
go get github.com/golang/freetype
```

Run server
```
go run main.go
```

Your server will be available on http://localhost:8090

# Usage
- countdownArial
- countdownOpenSans
- countdownPTM

Example url : 
http://localhost:8090/countdownArial?expires=2019-09-06T17:30:05&fg=000000&bg=ff2851&fontSize=88


Parameters:
* expires:
Countdown time in GMT format
Ex: 2019-09-06T17:30:05 

* fg
font color in hex format without #
Ex: 000000

* bg
background color in hex format without #
Ex: ff2851

