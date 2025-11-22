import uuid

from fastapi import FastAPI
from fastapi.requests import Request

from base_models import LogEntryBaseModel

app = FastAPI()


@app.post("/session/create")
async def start_session(request: Request) -> str:
    return str(uuid.uuid4())


@app.post("/session/upload")
async def upload_session(request: Request, logs: list[LogEntryBaseModel]): ...
