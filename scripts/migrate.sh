#!/bin/bash

# 加载 .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

COMMAND=$1
NAME=$2

case "$COMMAND" in

    up)
        migrate -path migrations \
            -database "$DATABASE_URL" \
            up
        ;;

    down)
        COUNT=${NAME:-1}

        read -p "Rolling back $COUNT migration(s). Continue? [y/N] " CONFIRM

        if [[ "$CONFIRM" == "y" || "$CONFIRM" == "Y" ]]; then
            migrate -path migrations \
                -database "$DATABASE_URL" \
                down "$COUNT"
        fi
        ;;

    create)
        migrate create \
            -ext sql \
            -dir migrations \
            -seq "$NAME"
        ;;

    force)
        migrate -path migrations \
            -database "$DATABASE_URL" \
            force $NAME \
        ;;

    *)
        echo "Usage:"
        echo "./migrate.sh up"
        echo "./migrate.sh down [count]"
        echo "./migrate.sh create migration_name"
        ;;
esac