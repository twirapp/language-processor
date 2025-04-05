# Language processor
Simple HTTP Server for translate texts and detect languages

## Local Run
> [!IMPORTANT]
>
> Minimum requirement [Go 1.24](https://go.dev/dl).


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
   ```bash
   go mod download
   ```
4. Run the server
   ```bash
   go run cmd/main.go
   ```

## Docker Hub
You can pull the pre-built Docker image from Docker Hub:
```bash
docker pull twirapp/language-processor
```

And run it with the command:
```
docker run --rm -p 3012:3012 --name language-processor twirapp/language-processor
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
  docker run --rm -p 3012:3012 --name language-processor language-processor
  ```

## Docker Compose
Create a `docker-compose.yml` file with the following content:
```yml
services:
  language-processor:
    image: twirapp/language-processor
    ports:
      - "3012:3012"
```

Then run:
```bash
docker compose up -d
```
