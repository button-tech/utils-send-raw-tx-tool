#!/usr/bin/env python3
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import uvicorn
import uuid
from starlette.middleware.cors import CORSMiddleware
import os
import asyncio

fift = "./liteclient-build/crypto/fift -I ./ton/crypto/fift/lib/"
lite_client = "./liteclient-build/lite-client/lite-client"


class Boc(BaseModel):
    data: str

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/ton/send")
async def send_boc(boc: Boc):
    result = await create_and_send_boc(boc.data)
    if result == "error":
        raise HTTPException(status_code=500, detail="err")
    return result


async def create_and_send_boc(hexData):
    fileName = str(uuid.uuid4().hex)

    text = '''
     B{''' + hexData + '''}
     "''' + fileName + '''.boc"

     tuck

     B>file
     ."(Saved to file " type .")" cr
     '''

    try:
        with open(f'{fileName}.fif', "w") as f:
            f.write(text)
    except:
        return "err"

    os.system(f'{fift} {fileName}.fif')
    await cli_call("sendfile " + fileName + ".boc")
    os.remove(f'./{fileName}.boc')
    os.remove(f'./{fileName}.fif')

    return {"result": "ok"}

async def cli_call(cmd):
    proc = await asyncio.create_subprocess_exec(
        lite_client, f'-c {cmd}',
        stderr=asyncio.subprocess.PIPE)

    data = await proc.stderr.read()

    data = data.decode('ascii').rstrip()

    await proc.wait()

    return data

if __name__ == "__main__":
   uvicorn.run("main:app", host="0.0.0.0", port=3000, workers=8, loop="asyncio")
