##  Sub-task 8.1 (Get familiar with Linux)

### 1. Run next code in storage dir

`docker run --rm -it --name storage_script12 -v ${PWD}:/usr/src/storage ubuntu`

### Now you are in docker container. You need to install curl
### 2. Install curl

`apt update; apt install curl -y`

### 3. change dir

`cd /usr/src/storage`


### 4. Run bash script

`bash sender.sh your_data_file your machine(default mac)`

## Your data delimiter in file is a new line