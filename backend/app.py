import uuid
import asyncio
from fastapi import FastAPI
from fastapi.requests import Request
from base_models import LogEntryBaseModel

app = FastAPI()
app.state.lock = asyncio.Lock()
app.state.total = 0


@app.post("/user/create")
async def create_user(request: Request) -> dict:
    return {"user_id": str(uuid.uuid4())}


@app.post("/session/create")
async def create_session(request: Request) -> dict:
    print("Request with headers", request.headers)
    return {"session_id": str(uuid.uuid4())}


@app.post("/session/upload")
async def upload_session(
    request: Request, session_id: str, tag: str, logs: list[LogEntryBaseModel]
):
    # print("Request with headers",request.headers)
    print("Current data", logs[0].message, "|", "Last data", logs[-1].message)
    async with app.state.lock:
        app.state.total += len(logs)
    # await asyncio.sleep(random.randint(2,10))
    print(app.state.total)
