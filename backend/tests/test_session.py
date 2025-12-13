# import uuid


def test_sample():
    c = 4 // 2
    assert c == 2


# def test_session_create(session_id):
#     # If the fixture didn't raise, the test passes
#     #TODO complete this
#     assert 1==1
#     # assert isinstance(uuid.UUID(session_id), uuid.UUID)


# def test_session_upload(dummy_app):
#     # TODO do full test after flow completiion
#     assert 1==1
# res = dummy_app.post(
#     "/session/upload",
#     params={"session_id": "b082a753-0f43-439d-92c4-0f5b72724b84", "tag": "sample-tag"},
#     json=[
#         {
#             "message": "test",
#             "timestamp": "2025-11-22 08:41:14.509427+00:00",
#             "level": "INFO",
#         },
#     ],
# )
# assert res.status_code == 200
