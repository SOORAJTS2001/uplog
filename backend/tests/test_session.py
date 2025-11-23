import uuid


def test_session_create(session_id):
    # If the fixture didn't raise, the test passes
    assert isinstance(uuid.UUID(session_id), uuid.UUID)


def test_session_upload(dummy_app):
    res = dummy_app.post(
        "/session/upload",
        params={"session_id": "b082a753-0f43-439d-92c4-0f5b72724b84"},
        json=[
            {"message": "test", "timestamp": "2025-11-22 08:41:14.509427+00:00", "level": "INFO"},
        ],
    )
    assert res.status_code == 200
