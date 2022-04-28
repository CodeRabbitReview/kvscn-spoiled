##  Sub-task 8.1 (Get familiar with Linux)

### 1. Run next code in storage dir


`docker run --rm -it --name storage -v ${PWD}:/usr/src/storage ubuntu`

### Now you are in docker container. You need to install curl
### 2. Install curl


`apt update; update install curl`

### 3. Run bash script

`bash sender.sh your_data_file`
