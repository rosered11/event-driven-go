# Use the official Confluent Kafka image as the base image
FROM confluentinc/cp-kafka:latest

# Copy your Java Agent JAR into the container
COPY jmx_prometheus_javaagent-0.19.0.jar /usr/share/jmx_exporter/jmx_prometheus_javaagent-0.12.0.jar

# Set the KAFKA_OPTS environment variable to include the Java Agent
# ENV KAFKA_OPTS="-javaagent:/opt/kafka/libs/your-agent.jar=config=/opt/kafka/config/agent-config.json"
