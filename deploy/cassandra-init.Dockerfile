FROM cassandra:4.1

COPY cmd/cassandra/init.cql /init.cql
COPY cmd/cassandra/init.sh /init.sh
