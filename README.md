Backend server for image uploading

It has swagger on route /swagger/

basic routes:
POST /upload - to upload JPEG image
GET /get - to get image by its id (uuid token)


DOCKER:

to run in Windows use the following command:

docker run -dp 30001:80 -e DB_HOST=host.docker.internal imloader  

