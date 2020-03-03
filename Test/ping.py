import socket
import time
# Buffers
writingbuffer = bytearray()
readingbuffer = bytearray()

# Connect non-blocking to the given server with ip and port


def tcp_connect(ip, port):
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.connect((ip, port))
    server_socket.setblocking(0)
    return server_socket

### Writing functions ###


def clearbuffer():
    global writingbuffer
    writingbuffer = bytearray(0)
    writebyte()
    writebyte()
# Write 1 byte to the writing buffer


def writebyte(byte=0):
    global writingbuffer
    writingbuffer.append(byte)
# Write a 2-byte short to the writing buffer


def writeshort(value):
    global writingbuffer
    a = value.to_bytes(2, byteorder='big')
    writingbuffer = writingbuffer + a

# Write a 4-byte integer to the writing buffer


def writeint(value):
    global writingbuffer
    a = value.to_bytes(4, byteorder='big')
    writingbuffer = writingbuffer + a

# Write a string to the writing buffer
# NOTE the string will end in a null terminating byte(0)


def writestring(text):
    global writingbuffer
    a = bytes(text, encoding='ascii')
    writingbuffer = writingbuffer + a
    writebyte()

# Send the writing buffer to the given socket


def sendmessage(sock):
    global writingbuffer
    size = len(writingbuffer)
    writingbuffer[0] = size - 2
    sock.send(writingbuffer)

### Reading functions ###

# Receive message from given sock
# returns size of message or 0 if no message was received


def receivemessage(sock):
    global readingbuffer

    buffer = bytearray(256)
    try:
        nbytes = sock.recv_into(buffer)
        readingbuffer = readingbuffer + buffer[:nbytes]
        #print('messaged received size ' + str(len(readingbuffer)))
    except socket.error as e:
        size = len(readingbuffer)
        # Read first 2 bytes that dont contain message
        if (size > 2):
            readbyte()
            readbyte()
        return size
    # Read first 2 bytes
    if (nbytes > 0):
        readbyte()
        readbyte()
    return nbytes

# Read 1 byte from the reading buffer


def readbyte():
    global readingbuffer
    return readingbuffer.pop(0)

# Read a string from the reading buffer
# NOTE readstring creates a string from continously reading bytes
# until it reaches null terminating byte (0)


def readstring():
    char = readbyte()
    word = bytearray()
    while char != 0:
        word.append(char)
        char = readbyte()
    return word.decode('ascii')

# Read a 4 byte integer from the reading buffer


def readint():
    number = bytearray(0)
    for x in range(4):
        number.append(readbyte())
    return int.from_bytes(number, byteorder='big')

# Read a 2 byte integer from the reading buffer


def readshort():
    number = bytearray(0)
    for x in range(2):
        number.append(readbyte())
    return int.from_bytes(number, byteorder='little')

# Print both the reading and writing buffer
# Mainly for debugging


def printbuffer():
    global readingbuffer, writingbuffer
    print('Reading buffer')
    print(readingbuffer)
    print('Writing buffer')
    print(writingbuffer)


if __name__ == '__main__':
    sock = tcp_connect('127.0.0.1', 3223)

    clearbuffer()
    writebyte(255)
    writestring("HI this is a string")
    writeint(12345)
    sendmessage(sock)
    count = 0
    while True:
        time.sleep(.2)
        clearbuffer()
        writebyte(count)
        writestring("HI this is in loop")
        writeint(12345)
        sendmessage(sock)
        count += 1
        if (receivemessage(sock) > 0):
            print(readbyte())
            print(readint())
            print(readstring())
