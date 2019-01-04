FROM scratch

COPY main/ladder /bin/ladder

ENTRYPOINT ["ladder"]