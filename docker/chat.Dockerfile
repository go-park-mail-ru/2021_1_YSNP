FROM build

FROM alpine:latest

COPY --from=build /app/chat /app/

COPY --from=build /app/config.json /app/

WORKDIR /app

CMD ./chat