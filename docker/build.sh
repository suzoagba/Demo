# Build docker image with tag forum
docker image build -t forum .
# If the above command gave an error
if [ $? -ne 0 ]; then
    echo "Image build failed, try again"
    exit 1;
else
    while true; do
        # Read input from terminal
        read -p "Clean up any unused objects [y/N] " yn
        case $yn in
            # Y or y followed by any letters
            [Yy]* ) docker system prune && exit;;
            [Nn]* ) exit;;
            * ) exit;;
        esac
    done
fi