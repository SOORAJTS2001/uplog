from pydantic import BaseModel


class SessionResponseBaseModel(BaseModel):
    session_id: str
