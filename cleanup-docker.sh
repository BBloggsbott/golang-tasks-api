# Remove all stopped containers
docker container prune -f

# Remove the specific container if it exists
docker rm -f task-api-mysql
docker rm -f task-api-server

# Remove dangling images
docker image prune -f

# If still issues, remove the volume
docker volume rm task-api_mysql_data