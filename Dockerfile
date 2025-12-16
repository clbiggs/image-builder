FROM scratch
COPY imagebuilder /
ENTRYPOINT ["/imagebuilder"]
