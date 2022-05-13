##  Sub-task 8.1 (Get familiar with Linux)

### 1. Run next code in storage dir


`docker run -it -v $(pwd):/usr/src/storage --name storage_script --rm ubuntu`

### Now you are in docker container. You need to install curl
### 2. Install curl

`apt update; apt install curl -y`

### 3. change dir

`cd /usr/src/storage`


### 4. Run bash script

`bash sender.sh your_data_file`

## Your data delimiter in file is a new line