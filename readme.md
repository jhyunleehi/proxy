```
C:\Gocode\src>curl -v http://localhost:8081/
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> GET / HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.83.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sat, 06 Aug 2022 03:35:20 GMT
< Content-Length: 5
< Content-Type: text/plain; charset=utf-8
<
hello* Connection #0 to host localhost left intact

C:\Gocode\src>
```