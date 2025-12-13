from sqlalchemy.orm import declarative_base
from sqlalchemy import Boolean, Column, String, Integer, DateTime, LargeBinary
from sqlalchemy.sql import func

Base = declarative_base()


class AnonymousUsers(Base):
    __tablename__ = "anonymous_user"
    __table_args__ = {"schema": "uplog"}
    user_id = Column(String, primary_key=True)
    hashed_ip = Column(LargeBinary, index=True)
    sessions_alive = Column(Integer, index=True, default=0)
    sessions_removed = Column(Integer, index=True, default=0)
    created_at = Column(
        DateTime(timezone=True),  # it should be utc
    )
    last_updated_timestamp = Column(
        DateTime(timezone=True),  # tell SQLAlchemy this is tz-aware
        server_default=func.now(),
        onupdate=func.now(),
    )


class VerifiedUsers(Base):
    __tablename__ = "verified_user"
    __table_args__ = {"schema": "uplog"}

    api_key = Column(String, index=True, unique=True)
    user_id = Column(String, primary_key=True)
    user_name = Column(String, index=True)
    email = Column(String, index=True)
    google_uid = Column(String, index=True)
    hashed_ip = Column(LargeBinary, index=True)
    sessions_alive = Column(Integer, index=True)
    sessions_removed = Column(Integer, index=True)
    created_at = Column(
        DateTime(timezone=True),  # it should be utc
    )
    last_updated_timestamp = Column(
        DateTime(timezone=True),  # tell SQLAlchemy this is tz-aware
        server_default=func.now(),
        onupdate=func.now(),
    )


class Sessions(Base):
    __tablename__ = "sessions"
    __table_args__ = {"schema": "uplog"}

    session_id = Column(String, primary_key=True)
    enable_sharing = Column(Boolean)
    user = Column(String, index=True)
    stream_name = Column(String, index=True)
    subject_name = Column(String, index=True, unique=True)
    log_line_count = Column(Integer, default=0)
    expires_at = Column(
        DateTime(timezone=True),  # it should be utc
    )
    created_at = Column(
        DateTime(timezone=True),  # it should be utc
    )
    last_updated_timestamp = Column(
        DateTime(timezone=True),  # tell SQLAlchemy this is tz-aware
        server_default=func.now(),
        onupdate=func.now(),
    )
    tag = Column(String, index=True)
