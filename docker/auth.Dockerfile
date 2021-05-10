FROM build

FROM alpine:latest

COPY --from=build /app/auth /app/

COPY --from=build /app/config.json /app/

WORKDIR /app

CMD ./auth