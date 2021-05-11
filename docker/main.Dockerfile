FROM build

FROM alpine:latest

COPY --from=build /app/main /app/

COPY --from=build /app/config.json /app/

WORKDIR /app

CMD ./main