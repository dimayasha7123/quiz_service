FROM gomicro/goose

WORKDIR /migrations/
ADD ./migrations/*.sql ./
ADD ./goose_script.sh ./

RUN chmod +x ./goose_script.sh

ENTRYPOINT "./goose_script.sh"
