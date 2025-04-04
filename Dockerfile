FROM python:3.11-slim-bookworm AS python-and-curl
RUN apt-get update && apt-get -y --no-install-recommends install curl

# Install all dependencies from pyproject.toml
FROM python-and-curl AS dependencies-installer
WORKDIR /app
# Install C++ build tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    g++ \
    && rm -rf /var/lib/apt/lists/*

RUN <<EOF
    curl -LsSf https://rye-up.com/get | RYE_INSTALL_OPTION="--yes" RYE_TOOLCHAIN=/usr/local/bin/python3 bash
    ln -s /root/.rye/shims/rye /usr/local/bin/rye
    rye pin 3.11
EOF
COPY pyproject.toml .
RUN --mount=type=cache,target=/root/.cache rye sync --no-dev

FROM gcr.io/distroless/python3-debian12:nonroot
LABEL org.opencontainers.image.authors="Satont <satontworldwide@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/twirapp/language-processor"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="Language Processor"
LABEL org.opencontainers.image.description="Simple HTTP Server for translate texts and detect languages"
LABEL org.opencontainers.image.vendor="TwirApp"

WORKDIR /app
ENV PATH=/app/.venv/bin:$PATH
ENV PYTHONPATH="/app/.venv/lib/python3.11/site-packages/:${PYTHONPATH:-}"
USER nonroot

ADD --chmod=004 https://dl.fbaipublicfiles.com/fasttext/supervised-models/lid.176.bin .

COPY --from=dependencies-installer /app/.venv .venv
COPY ./app app
CMD ["/app/.venv/bin/uvicorn", "--host", "0.0.0.0", "app.server:app"]