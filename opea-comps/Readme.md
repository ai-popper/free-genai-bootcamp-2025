## Running ollama third party Service
   
### choosing a model

you can get a model_id that ollama will launch from the [ollama library](https://ollama.com/library).

eg. https://ollama.com/library/llama3.2:1b
   

   
### Getting the Host IP

#### Linux
```sh
apt install net-tools
ifconfig
```
or you can try this way `=$(hostname -I | awk '{print $1}')`

HOST_IP==$(hostname -I | awk '{print $1}') NO_PROXY=localhost LLM_ENDPOINT_PORT=8008 LLM_MODEL_ID="llama3.2:1b"
docker-compose up

### Ollama API

Once the ollama server is running we can make API calls to ollama API.

https://github.com/ollama/ollama/blob/main/docs/api.md


## Download(Pull) a model

curl http://localhost:8008/api/pull -d '{
  "model": "llama3.2:1b"
}'


## Generate a Request

curl http://localhost:8008/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt": "Why is the sky blue?"
}'

# Technical uncertainty

Q: will the model be downloaded into the container ? does that mean the model
 will be deleted when the container stops running?

> The model will download into the continer, and vanish when the container stops running . 
you need to mount a local drive and there is probably more work to be done.

Q: Which port is being mapped 8008->11434

> In this case 8008 is the port that host machine will access, 
and 11434 is the port inside the container that ollama server is listening on.

Q: if we pass the LLM_MODEL_ID to the ollama server will it download the model
when on start?

> It does not appear so, the ollama CLI might be running the multiple
 APIs so you need to call the pull api before trying to generate text.