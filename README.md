# json2smtp

An email proxy: input: **json**, output: **smtp call**

For a legacy project I needed to have a **proxy** that reads **json** input and execute a **smtp** call in order to **send emails**. So I created a small proxy for emails in go (golang)

Read more about why I needed it here: https://www.c2kb.com/json2smtp

# Table of Contents
1. [How it works](#howitworks)
2. [The **json** struct object](#jsonstruct)
    1. [Simple object](#simpleobject)
    2. [Full with smtp data](#fullobject)
         1. [Adding attachments](#attachments)
3. [How to install](#install)
     1. [Build and compile](#buildandcompile)
     2. [Execute the proxy examples](#execute)
          1. [Run the proxy](#example1)
          2. [Run in the background](#runinbackground)
          3. [Simple execute](#simpleexecute)
     3. [Recommended architecture](#recommended)
4. [Download binaries](#download)
5. [Todo](#todo)
6. [Libraries used](#libraries)
7. [Thank you](#thankyou)


## How it works: <a name="howitworks"></a>
Simple calling diagram

![Simple architecture of calling the json2smtp email proxy server with json and smtp calls](https://www.c2kb.com/json2smtp/json2smtp_architecture_1.jpg)

## The **json** struct object <a name="jsonstruct"></a>
### Simple object: <a name="simpleobject"></a>
	curl -X POST \
	 -H "Content-Type: application/json" \
	 -d '{ \
	"from": "john doe <john@example.com>", \
	"to": ["kermit@muppets.com", "oneperson@example.com"], \
	"cc": ["email1@example.com"], \
	"bcc": ["secret@example.com"], \
	"subject": "email subject line", \
	"message": "message body in text/html to be sent", \
	"attachments": {"filename.pdf": "base64 file encoded", "anotherfilename.txt": "base64 file encoded"}, \
	 }' \
	 http://localhost:8080/


### Full with smtp data: <a name="fullobject"></a>
	curl -X POST \
	 -H "Content-Type: application/json" \
	 -d '{ \
	"from": "john doe <john@example.com>", \
	"to": ["kermit@muppets.com", "oneperson@example.com"], \
	"cc": ["email1@example.com"], \
	"bcc": ["secret@example.com"], \
	"subject": "email subject line", \
	"message": "message body in text/html to be sent", \
	"attachments": {"filename.pdf": "base64 file encoded", "anotherfilename.txt": "base64 file encoded"}, \
	"smtphost": "smtp.example.com - optional parameter", \
	"smtpport": 587 - optional paramater, \
	"smtpuser": "username - optional parameter", \
	"smtppassword": "password - optional parameter" \
	 }' \
	 http://localhost:8080/

#### Attachments <a name="attachments"></a>
In order to send attachments with your json email struct you need to construct an object of base64 encoded string of your binary file.

## How to install: <a name="install"></a>
Download the code and run it

	git clone https://github.com/caviv/json2smtp.git
	go run ./
	go run ./ --help

### Build and compile <a name="buildandcompile"></a>
Download  the code compile it and run with help command

	git clone https://github.com/caviv/json2smtp.git
	go build ./
	./json2smtp --help

### Execute the proxy examples <a name="execute"></a>
Command line help:

	json2smtp utility https://www.c2kb.com/json2smtp v1.0.1 2023-11-13
	Get json input and calls smtp - function as a json to smtp proxy
	Options:
	  -port int
	    	the port to listen on (default 8080)
	  -smtphost string
	    	smtp host, e.g. smtp.example.com
	  -smtpoverride
	    	true - allows to pass smtp parameters in the json call, false will always use the config smtp data (default true)
	  -smtppassword string
	    	password for the smtp user
	  -smtpport int
	    	the port to listen on (default 587)
	  -smtpuser string
	    	username for the smtp

#### Run the proxy: <a name="example1"></a>
	json2smtp --port=8200 --smtphost='smtp.example.com' --smtpport=587 --smtpuser='username' --smtppassword='password' --smtpoverride=false

#### Run in the background: <a name="runinbackground"></a>
	nohup json2smtp --port=8200 --smtphost='smtp.example.com' --smtpport=587 --smtpuser='username' --smtppassword='password' >> logfile.log 2>&1 &

#### Simple execute: <a name="simpleexecute"></a>
In  this way the calling client will have to pass the smtp server details in each call because we don't set the smtp default server for the proxy. The default port to listen on is 8080.

    json2smtp 

### Recommended architecture <a name="recommended"></a>
![Calling json2smtp proxy behind a caddy web server for https / tls](https://www.c2kb.com/json2smtp/json2smtp_architecture_2.jpg)

## Download binaries: <a name="download"></a>
https://www.c2kb.com/json2smtp

## Todo: <a name="todo"></a>
* Replace the mail.v2 package to a newer one - as recommended here: https://www.reddit.com/r/golang/comments/17w2l9s/comment/k9fve4r/
* Create delayed sending of emails / asynchronic sending - I think that could be useful

## Libraries used <a name="libraries"></a>
This external libraries are used in the project:

require  gopkg.in/mail.v2  v2.3.1
require  gopkg.in/alexcesaro/quotedprintable.v3  v3.0.0

## Thank you <a name="thankyou"></a>
contact me on support at c2kb.com