# ToDoApp
## The Best Way to Organize 

[![N|Solid](https://cldup.com/dTxpPi9lDf.thumb.png)](https://nodesource.com/products/nsolid)

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

ToDoApp is a lightweight app to maintain your daily tasks.

## Features

- Login with gmail and reset your password using your registered email address.
- Password encryption using HS264.
- Upload and download attachments with your tasks.
- Remainds about the pending tasks on a given day through emails.
- Lists down similar tasks for a given user.


## Docker

Dillinger is very easy to install and deploy in a Docker container.

By default, the Docker will expose port 8080, so change this within the
Dockerfile if necessary. When ready, simply use the Dockerfile to
build the image.

```sh
cd ToDoApp
docker build -t <user>/ToDoApp:${package.json.version} .
```

This will create the ToDoApp image and pull in the necessary dependencies.


Once done, run the Docker image and map the port to whatever you wish on
your host. 
```sh
docker run -d -p 8000:myPORT --restart=always --cap-add=SYS_ADMIN --name=ToDoApp <Junaid>r:${package.json.version}
```

All Rights Reserved.
