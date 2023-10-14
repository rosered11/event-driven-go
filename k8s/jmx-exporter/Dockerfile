FROM openjdk:11-jre-slim

# Set environment variables
ENV JMX_EXPORTER_VERSION=0.16.0
ENV JMX_EXPORTER_CONFIG_PATH=/etc/jmx_exporter/

# Download jmx_exporter JAR
RUN mkdir -p /jmx_exporter \
    && wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_httpserver/$JMX_EXPORTER_VERSION/jmx_prometheus_httpserver-$JMX_EXPORTER_VERSION-jar-with-dependencies.jar -O /jmx_exporter/jmx_prometheus_httpserver.jar

# Create a directory for the jmx_exporter config
RUN mkdir -p $JMX_EXPORTER_CONFIG_PATH

# Copy your jmx_exporter configuration file (e.g., jmx_exporter_config.yml)
COPY jmx_exporter_config.yml $JMX_EXPORTER_CONFIG_PATH

CMD ["java", "-javaagent:/jmx_exporter/jmx_prometheus_httpserver.jar=9404:$JMX_EXPORTER_CONFIG_PATH/jmx_exporter_config.yml", "-jar", "/path/to/your-application.jar"]
