import uuid


def test_session_create(session_id):
    # If the fixture didn't raise, the test passes
    assert isinstance(uuid.UUID(session_id), uuid.UUID)


def test_session_upload(dummy_app):
    res = dummy_app.post(
        "/session/upload",
        json=[
            {"message": "test", "timestamp": "2025-11-22 08:41:14.509427+00:00", "level": "INFO"},
        ],
    )
    assert res.status_code == 200
