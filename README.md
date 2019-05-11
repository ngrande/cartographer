# Cartographer
Ease navigation in an intuitive way, for the web

## Design
An easy to setup web server. One simply has to create a directory (root) and everything within it will be served by the Cartographer.

---

Simple setup:

    ~/base
    |
    +--index.html
    +--notes/
    |  |
    |  +--text.md
    +--index.md
    +--pictures/
       |
       +--fremen.jpg
       +--sandworm.jpg

Cartographer will serve everything as it is with ~/base as its root.
So you reach ~/base/index.html by going to http://localhost:8080/index.html
Or call localhost:8080/pictures/sandword.jpg to get a nice picture!

---

## Requirements
The following utils have to be installed on the machine to ensure a flawless experience:
+ pandoc

## Features
### (Auto) Markdown conversion
Every file with the suffix ```.md``` will be treated as a markdown text file. Those files are automatically converted to HTML using [pandoc](https://pandoc.org/).
Currently markdown files are not supported by the template engine.

### Template Engine
Don't repeat yourself. We all know it, thus also this project knows it and ships with a (admittedly) very simplistic template engine.
Besides the root directory for the web content you have to create another directory where you store the templates.
Let's say ```~/templates``` and in here you simply create files, like a normal HTML. In there you can put keywords wich have to be surrounded by ```$``` characters.
For example you could create an ```default.html``` and put in there a keyword ```$TITLE$``` (this can be placed multiple times in the same file).
Now to make use of this template you have to add the whole template filename as suffix to your file which shall use the template: ```index.default.html``` (we call it a template implementation).
This file is now expected to contain a section for every keyword from ```default.html``` - let's call these sections ```content```.
The content has to have the following format: ```<$KEYWORD$>Your text comes here<$KEYWORD$>``` whereas ```KEYWORD``` is a placeholder for your keyword (i.e. TITLE).
The keywords in the template file will then be replaced by the content specified in your template implementation.
The result will then be stored in memory and the URL will point to the name of the template implementation using the last suffix without the template name: ```index.default.html -> index.html```.
