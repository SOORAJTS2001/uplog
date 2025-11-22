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
