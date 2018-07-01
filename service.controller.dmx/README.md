# service.controller.dmx

Controls ADJ Megapar Profile lights.

Manual: http://www.fullcompass.com/common/files/14961-AmericanDJMegaParProfileUserManual.pdf

---
Why normal Python libraries don't work
https://stackoverflow.com/a/22508876/3105582

Which OLA plugin works
http://www.martinjhiggins.co.uk/hashtag-lighting/
TL;DR disable `Enttec Open DMX` and `Serial USB`, enable `FTDI USB DMX`

The docker image has the OLA python libraries available at
`/usr/local/lib/python2.7/dist-packages/ola`

Raspbian has them at
`/usr/lib/python2.7/dist-packages/ola`
