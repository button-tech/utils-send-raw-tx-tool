#!/usr/bin/env python3.7
import sys
import os

def send_grams(raw_tx):

    h = str(hash(raw_tx))
    filename = h[1:]
    text = '''
    B{'''+ raw_tx +'''}
    "'''+ filename +'''.boc" B>file
    '''

    try:
        with open(f'/app/wrappers/{filename}.fift', "w") as f:
            f.write(text)
        os.system(f"/app/liteclient-build/crypto/fift -I /app/lite-client/crypto/fift/lib/ /app/wrappers/{filename}.fift")
        os.system(f"mv /app/{filename}.boc /app/wrappers")
        os.system(f"rm /app/wrappers/{filename}.fift")
        os.system(f"/app/wrappers/sendfile {filename}")
        os.system(f"rm /app/wrappers/{filename}.boc")
        return True
    except:
        return False


result = send_grams(sys.argv[1])

if result == False:
    print("error")
else:
    print("success")