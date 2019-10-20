#!/usr/bin/env python3.7
import sys
import subprocess
import os

def send_grams(raw_tx, workdir):

    h = str(hash(raw_tx))
    filename = h[1:]
    text = '''
    B{'''+ raw_tx +'''}
    "'''+ filename +'''.boc" B>file
    '''

    try:
        with open(f'{workdir}wrappers/{filename}.fift', "w") as f:
            f.write(text)
        os.system(f"{workdir}liteclient-build/crypto/fift -I {workdir}lite-client/crypto/fift/lib/ {workdir}wrappers/{filename}.fift")
        os.system(f"mv {workdir}{filename}.boc {workdir}wrappers")
        os.system(f"rm {workdir}wrappers/{filename}.fift")
        subprocess.getoutput(f"{workdir}wrappers/sendfile {workdir} {filename}")
        os.system(f"rm {workdir}wrappers/{filename}.boc")
        return True
    except:
        return False


result = send_grams(sys.argv[1], sys.argv[2])

if result == False:
    print("error")
else:
    print("success")