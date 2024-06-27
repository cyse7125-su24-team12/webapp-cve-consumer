#!/bin/sh

# Function to check PostgreSQL readiness
check_postgres() {
    echo "Checking if the PostgreSQL database is up and running at host: $DB_HOST"

    # Assuming DB_HOST, DB_PORT, DB_USER, and DB_NAME are environment variables
    pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "PostgreSQL is reachable."
        # exit 0 
    else
        echo "PostgreSQL is not reachable."
        exit 1
    fi
}

# Function to check Kafka readiness
check_kafka() {
    echo "Checking if Kafka topic $TOPIC_NAME exists at $KAFKA_HOST"
    # exit 1 # Temporarily disabling Kafka check
    # Assuming KAFKA_HOST, KAFKA_USER, KAFKA_PASSWORD, and TOPIC_NAME are environment variables
    # Prepare a client properties file with credentials
    echo "security.protocol=SASL_PLAINTEXT
    sasl.mechanism=PLAIN
    sasl.jaas.config=org.apache.kafka.common.security.scram.ScramLoginModule required username=\"$KAFKA_USER\" password=\"$KAFKA_PASSWORD\";" > /tmp/client.properties
   
    # echo "Starting Kafka check..."
    # kafka-topics.sh --bootstrap-server "$KAFKA_HOST" --command-config /tmp/client.properties --list
  

    kafka-topics.sh --bootstrap-server "$KAFKA_HOST" --command-config /tmp/client.properties --list | grep -w "$TOPIC_NAME" > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Kafka topic $TOPIC_NAME is accessible."
        exit 0
    else
        echo "Kafka topic $TOPIC_NAME is not accessible."
        exit 1
    fi
}

# Main script execution starts here

# Check PostgreSQL
check_postgres

# Check Kafka
check_kafka





# # Assuming DB_HOST, DB_PORT, DB_USER, and DB_NAME are environment variables
# pg_isready -h "$DB_HOST" -p" $DB_PORT" -U "$DB_USER" -d "$DB_NAME" > /dev/null 2>&1

# if [ $? -eq 0 ]; then
#     echo "PostgreSQL is reachable."
#     exit 0
# else
#     echo "PostgreSQL is not reachable."
#     exit 1
# fi
