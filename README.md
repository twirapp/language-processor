# Language processor
Simple HTTP Server for translate texts and detect languages

# Installation
> [!WARNING]
> The project is designed to run on CPU, if you want to use GPU you will have to replace torch dependency in `pyproject.toml`.

> [!TIP]
> It works well on hetzner Shared vCPU AMD server with 4cpu, 8gb ram, handles messages under 100ms. Maybe less resources needed, check this out yourself.

## Local Run
> [!IMPORTANT]
>
> Minimum requirement [Python 3.9](https://www.python.org/downloads).
>
> This project uses [Rye](https://rye.astral.sh) for dependency management, but it is also possible to install dependencies via pip. This is not necessary.

1. Clone the repository

    ```bash
    git clone https://github.com/twirapp/language-processor.git && cd language-processor
    ```
2. Download model
   Run this in the project root directory:
    ```bash
    wget https://dl.fbaipublicfiles.com/fasttext/supervised-models/lid.176.bin
    ```
3. Install dependencies

   This will automatically create the virtual environment in the `.venv` directory and install the required dependencies
    ```bash
    rye sync
    ```
    <details>
    <summary>(not recommended) alternative install via pip</summary>
    Create a virtual environment and activate:

    ```bash
    python3 -m venv .venv && source .venv/bin/activate
    ```
   Install only the required dependencies:

    ```bash
    pip3 install --no-deps -r requirements.lock
    ```
    </details>
4. Run the server

   With autoload:
    ```bash
    rye run dev-server
    ```
   Without autoload:
    ```bash
    rye run server
    ```
    <details>
    <summary>Without Rye</summary>

   With autoload:
    ```bash
    uvicorn app.server:app --reload
    ```

   Without autoload:
    ```bash
    uvicorn app.server:app
    ```
    </details>

## Docker Hub
You can pull the pre-built Docker image from Docker Hub:
```bash
docker pull twirapp/language-processor
```

And run it with the command:
```
docker run --rm -p 8000:8000 --name language-processor twirapp/language-processor
```

## Docker Build
1. Clone the repository

  ```bash
  git clone https://github.com/twirapp/language-processor.git && cd language-processor
  ```
2. Build the Docker image

  ```bash
  docker build -t language-processor .
  ```
3. Run the container

  ```bash
  docker run --rm -p 8000:8000 --name language-processor language-processor
  ```

## Docker Compose
Create a `docker-compose.yml` file with the following content:
```yml
services:
  language-processor:
    image: twirapp/language-processor
    ports:
      - "8000:8000"
    environment:
      TOXICITY_THRESHOLD: 0
      # WEB_CONCURRENCY: 1 # uvicorn workers count
```

Then run:
```bash
docker compose up -d
```
