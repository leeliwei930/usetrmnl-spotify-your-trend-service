# ========================================
# Stage 1: Builder
# ========================================
FROM python:3.12-alpine AS builder

# Install build dependencies
# gcc, musl-dev: for compiling Python packages with C extensions
# libffi-dev: for cryptography package
# linux-headers: for some packages that need kernel headers
RUN apk add --no-cache \
	gcc \
	musl-dev \
	libffi-dev \
	linux-headers \
	make \
	g++

# Set working directory
WORKDIR /build

# Copy requirements file
COPY requirements.txt .

# Create virtual environment and install dependencies
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN pip install --no-cache-dir --upgrade pip && \
	pip install --no-cache-dir -r requirements.txt

# ========================================
# Stage 2: Runtime
# ========================================
FROM python:3.12-alpine

# Install runtime dependencies only
RUN apk add --no-cache libffi

# Copy virtual environment from builder
COPY --from=builder /opt/venv /opt/venv

# Set working directory
WORKDIR /app

# Copy application code
COPY . .

# Set environment variables
ENV PATH="/opt/venv/bin:$PATH" \
	PYTHONUNBUFFERED=1 \
	PYTHONDONTWRITEBYTECODE=1

# Expose port
EXPOSE 8000

# Run the application
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
