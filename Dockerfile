FROM golang:latest

WORKDIR /opt/token_bank/sync_node_peers

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir -p /opt/token_bank/sync_node_peers
RUN mkdir -p /opt/token_bank/sync_node_peers/conf
RUN mkdir -p /opt/token_bank/sync_node_peers/bin

ADD tb.sync_node_peers.s /opt/token_bank/sync_node_peers/
ADD conf/* /opt/token_bank/sync_node_peers/conf/
ADD bin/* /opt/token_bank/sync_node_peers/bin/

ENTRYPOINT ["sh ./bin/start.sh"]