FROM scratch
EXPOSE 8080
ENTRYPOINT ["/cli-apps"]
COPY ./build/linux /