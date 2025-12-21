import uuid
import asyncio
from fastapi import FastAPI, Header
from fastapi.requests import Request
from base_models import LogEntryBaseModel
from models import Base, AnonymousUsers, Sessions
from settings import async_engine, NATS_URL, NATS_MAX_MSG_PER_SUBJECT, NATS_MSG_TTL_IN_SECONDS
from utils import get_hashed_client_ip
from contextlib import asynccontextmanager
from collections.abc import AsyncGenerator
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker
from fastapi import Depends
from nats.js import api
from datetime import datetime, timezone, timedelta
import nats
from fastapi.responses import StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
import json


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator:
    nc = await nats.connect(NATS_URL)
    app.state.js = nc.jetstream()
    async with async_engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    yield


async def subject_consumer(subject_name: str):
    sub = await app.state.js.subscribe(subject_name)
    temp_data = []
    try:
        while True:
            try:
                msg = await sub.next_msg(timeout=5)
            except TimeoutError:
                # keep connection alive, add sleep here
                yield ": heartbeat\n\n"
                continue

            data = msg.data.decode("utf-8")
            if len(temp_data) == 10:
                yield f"data: {json.dumps(temp_data)}\n\n"
                temp_data = []
            temp_data.append(json.loads(data))
            await asyncio.sleep(0.1)
            await msg.ack()

    except asyncio.CancelledError:
        # Client disconnected
        await sub.unsubscribe()
        raise


AsyncSessionLocal = async_sessionmaker(
    async_engine,
    expire_on_commit=False,  # Prevents objects from expiring after commit, useful for returning them
    class_=AsyncSession,  # This is crucial for async sessions
)


# --- Database Dependency (Async) ---
async def get_db_session() -> AsyncGenerator[AsyncSession, None]:
    """Dependency that provides an async database session."""
    async with AsyncSessionLocal() as session:
        yield session
        await session.close()  # Ensure session is closed


app = FastAPI(lifespan=lifespan)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:8080"],
    allow_credentials=True,
    allow_methods=["*"],  # Allows all HTTP methods (GET, POST, PUT, DELETE, etc.)
    allow_headers=["*"],  # Allows all headers
)


@app.post("/user/create")
async def create_user(request: Request, db: AsyncSession = Depends(get_db_session)) -> dict:
    client_ip = get_hashed_client_ip(request)
    user_id = str(uuid.uuid4())
    config = api.StreamConfig(
        name="stream-" + user_id,
        subjects=[f"subject.{user_id}.*"],
        max_msgs_per_subject=NATS_MAX_MSG_PER_SUBJECT,
        max_age=NATS_MSG_TTL_IN_SECONDS,
    )
    await app.state.js.add_stream(config)
    user = AnonymousUsers(
        user_id=user_id,
        hashed_ip=client_ip,
        created_at=datetime.now(timezone.utc),  # noqa: UP017
    )
    db.add(user)
    await db.commit()
    print("Stream created")
    return {"user_id": user_id}


@app.post("/session/create")
async def create_session(
    request: Request,
    user_id: str = Header(..., alias="User-Id"),
    db: AsyncSession = Depends(get_db_session),
) -> dict:
    print("Request with headers", request.headers)
    session_id = str(uuid.uuid4())
    stream_name = "stream-" + user_id
    subject_name = "subject." + user_id + "." + session_id
    session = Sessions(
        session_id=session_id,
        enable_sharing=False,
        user=user_id,
        subject_name=subject_name,
        stream_name=stream_name,
        expires_at=datetime.now(timezone.utc) + timedelta(days=2),  # noqa: UP017
        created_at=datetime.now(timezone.utc),  # noqa: UP017
    )
    db.add(session)
    await db.commit()
    return {"session_id": session_id}


@app.post("/session/upload")
async def upload_session(
    request: Request,
    session_id: str,
    tag: str,
    logs: list[LogEntryBaseModel],
    user_id: str = Header(..., alias="User-Id"),
):
    await asyncio.gather(*[
        app.state.js.publish(
            subject=f"subject.{user_id}.{session_id}", payload=log.model_dump_json().encode()
        )
        for log in logs
    ])


@app.get("/session/consume")
async def consume_session(
    request: Request,
    session_id: str,
):
    return StreamingResponse(subject_consumer(session_id), media_type="text/event-stream")
