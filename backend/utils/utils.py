from fastapi import Request
import hashlib
from settings import IP_HASH


def compute_hmac(salt: str, data: str):
    return hashlib.pbkdf2_hmac(
        "sha256",  # Hashing algorithm
        data.encode("utf-8"),  # Encode password to bytes
        salt.encode("utf-8"),  # The salt
        100000,  # Number of iterations (cost factor)
    )


def get_hashed_client_ip(request: Request) -> str:
    x_forwarded_for = request.headers.get("x-forwarded-for")
    client_ip = x_forwarded_for.split(",")[0].strip() if x_forwarded_for else request.client.host
    return compute_hmac(IP_HASH, client_ip)
