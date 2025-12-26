FROM python:3.12-slim

# Basic Python hygiene
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV UV_PROJECT_ENVIRONMENT="/usr/local/"

WORKDIR /app

# System deps (minimal)
RUN apt-get update \
    && apt-get install -y --no-install-recommends curl \
    && rm -rf /var/lib/apt/lists/*

# Install uv (fast, modern, correct)
RUN pip install --no-cache-dir uv

# Copy dependency metadata first (layer caching)
COPY pyproject.toml uv.lock ./

# Install dependencies (system-wide, no venv)
RUN uv sync --frozen --no-dev

# Copy application code
COPY . .

EXPOSE 8000

# Your entrypoint is app.py
CMD ["python", "-m", "uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
