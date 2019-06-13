FROM ddfddf/randgenx

RUN mkdir -p /root/result

COPY ./randgen-server /root

ENV CONFPATH=/root/conf
ENV RMPATH=/root/randgenx
ENV RESULTPATH=/root/result

WORKDIR /root

ENTRYPOINT ["./randgen-server"]