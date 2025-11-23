import uuid

import pytest
from fastapi.testclient import TestClient

from app import app


@pytest.fixture(scope="session")
def dummy_app():
    return TestClient(app)


@pytest.fixture(scope="session")
def session_id(dummy_app):
    res = dummy_app.post("/session/create")
    sid = res.json()["session_id"]
    uuid.UUID(sid)  # raises ValueError if invalid
    return sid
