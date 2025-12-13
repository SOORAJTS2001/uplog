from datetime import datetime
from enum import Enum

from pydantic import BaseModel


class LogLevelEnum(Enum):
    error = "ERROR"
    warn = "WARN"
    info = "INFO"
    debug = "DEBUG"


class LogEntryBaseModel(BaseModel):
    message: str
    timestamp: datetime
    level: LogLevelEnum

    def to_dict(self) -> dict:
        return {
            "message": self.message,
            "timestamp": self.timestamp.isoformat(),  # convert to string
            "level": self.level.value,  # Enum → string
        }
