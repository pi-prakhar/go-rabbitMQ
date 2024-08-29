# Use the official RabbitMQ base image
FROM rabbitmq:3-management

# Optional: Copy a custom RabbitMQ configuration file
# COPY rabbitmq.conf /etc/rabbitmq/rabbitmq.conf

# Expose RabbitMQ ports
EXPOSE 5672 5672

# Define the default command to run RabbitMQ server
CMD ["rabbitmq-server"]
