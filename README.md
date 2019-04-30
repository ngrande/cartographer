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
   |
   +--text.md
+--index.md
+--pictures/
   |
   +--fremen.jpg
   +--sandworm.jpg

Cartographer will serve everything as it is with ~/base as its root.
So you reach ~/base/index.html by going to http://localhost:8080/index.html
Or call localhost:8080/pictures/sandword.jpg to get a nice picture!
