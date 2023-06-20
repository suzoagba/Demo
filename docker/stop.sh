# Use 'sudo systemctl restart docker.socket docker.service' to restart Docker
# Stop container (use 'kill' to forcefully stop)
docker stop forum
# Remove stopped container
docker rm forum
# Remove Docker image
docker rmi -f forum
echo "Docker removed"