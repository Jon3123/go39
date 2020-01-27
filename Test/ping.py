import socket

# create an INET, STREAMing socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# now connect to the web server on port 80 - the normal http port
s.connect(("127.0.0.1", 3223))

s.send(b'Hi')
a = bytearray()
a.append(0)
s.send(a)
