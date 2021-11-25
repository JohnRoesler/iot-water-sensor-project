FROM scratch
ADD ./main /main
EXPOSE 8888
USER nobody:nogroup
ENTRYPOINT /main