# Setting up theSketchDB

Initial Server Setup: https://www.digitalocean.com/community/tutorials/initial-server-setup-with-ubuntu-18-04

Install Docker (Step 1&2): https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04

Install Docker Compose (Step 1): https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-18-04

Use this tutorial now to set up go app, nginx:
https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-with-docker-and-nginx-on-ubuntu-18-04


To analyze the assets image (or any base image for that matter):

```bash
docker build --target assets -t my-assets .
docker run -it my-assets /bin/sh
```

run docker system prune -a every once in a while (it can also help if things don't seem to be copying correctly between images)
