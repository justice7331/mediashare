# 🖥️ mediashare
A ShareX media server made in Go.
## 🧳 Requirements
* Go (latest)

## 🏗️ Building
```
go build ./src/mediashare
```

## 📘 mediashare Setup
Change the password variable in `main.go` to a password you would like to use, then build. Next create an SSL certificate and name the files `certificate.pem` & `certificate.key`.

## 📗 ShareX Setup
Create a custom uploader by going to Destinations -> Custom Uploader Settings... and clicking the new button. Set the destination type to "Image Uploader" and name it whatever you want. In the request tab set the method to "POST", set the url to `https://your.domain/upload`, and set the body to "Form data (multipart/form-data)". Add a row to the body form named "password" and set the value to what you did in `main.go`. Set the file form name to "media". Now go to the response tab and set the URL to `https://your.domain/$json:extension$/$json:filename$`, set the Thumbnail URL to `https://your.domain/$json:filename$.$json:extension$`, and set the error message to `$json:message$`. Now go to Destinations -> Image Uploader and set it to "Custom image uploader". Your custom uploader should look something like this:
![Request Tab](https://i.imgur.com/9i72Z4J.png)
![Response Tab](https://i.imgur.com/pyT7pcR.png)

## 🕹️ Usage
```
./mediashare
```
